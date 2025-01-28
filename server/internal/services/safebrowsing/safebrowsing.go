package safebrowsing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultBaseURL = "https://safebrowsing.googleapis.com/v4/threatMatches:find"

type SafeBrowsingChecker interface {
	IsURLSafe(url string) (bool, error)
	CheckURL(url string) (*ThreatResponse, error)
}

// SafeBrowsingService handles communication with the Google Safe Browsing API
type SafeBrowsingService struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewSafeBrowsingService creates a new instance of SafeBrowsingService
func NewSafeBrowsingService(apiKey string) *SafeBrowsingService {
	return &SafeBrowsingService{
		apiKey:  apiKey,
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 6 * time.Second,
		},
	}
}

// validateURL checks if the provided URL is valid
func validateURL(urlStr string) error {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	// Check if the URL has a scheme and host
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return fmt.Errorf("invalid URL: missing scheme or host")
	}

	// Check if the scheme is http or https
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("invalid URL scheme: must be http or https")
	}

	// Block reserved TLDs like .test
	if strings.HasSuffix(parsedURL.Host, ".test") {
		fmt.Println("Find error for debuging: ", parsedURL)
		return fmt.Errorf("URLs with .test domains are not allowed")
	}

	// Prevent javascript: protocol URLs that could execute code
	if strings.HasPrefix(strings.ToLower(urlStr), "javascript:") {
		return fmt.Errorf("URLs start with javascript domains are execute code which is not allowed")
	}
	// Prevent data: URLs that could contain malicious base64 encoded content
	if strings.HasPrefix(strings.ToLower(urlStr), "data:") {
		return fmt.Errorf("invalid execution code which start with data: is not allowed")
	}

	return nil
}

// CheckURL checks a single URL against the Safe Browsing API
func (s *SafeBrowsingService) CheckURL(url string) (*ThreatResponse, error) {
	// Validate URL before making API call
	if err := validateURL(url); err != nil {
		return nil, err
	}

	request := ThreatRequest{
		Client: ClientInfo{
			ClientID:      "yoururlshortener",
			ClientVersion: "1.0.0",
		},
		ThreatInfo: ThreatInfo{
			ThreatTypes:      []string{"MALWARE", "SOCIAL_ENGINEERING", "UNWANTED_SOFTWARE", "POTENTIALLY_HARMFUL_APPLICATION"},
			PlatformTypes:    []string{"ANY_PLATFORM"},
			ThreatEntryTypes: []string{"URL"},
			ThreatEntries: []ThreatEntry{
				{URL: url},
			},
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	fullURL := fmt.Sprintf("%s?key=%s", s.baseURL, s.apiKey)
	req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var threatResponse ThreatResponse
	if err := json.NewDecoder(resp.Body).Decode(&threatResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &threatResponse, nil
}

// IsURLSafe returns true if the URL is safe, false if it's potentially dangerous
func (s *SafeBrowsingService) IsURLSafe(url string) (bool, error) {
	response, err := s.CheckURL(url)

	if err != nil {
		return false, err
	}
	// If there are no matches, the URL is safe
	return len(response.Matches) == 0, nil
}
