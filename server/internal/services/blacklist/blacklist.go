// internal/services/blacklist.go
package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

type BlacklistService struct {
	patterns    []string
	regexRules  []*regexp.Regexp
	mu          sync.RWMutex
	lastUpdated time.Time
}

// SafeBrowsingRequest represents the request structure for Google Safe Browsing API
type SafeBrowsingRequest struct {
	Client struct {
		ClientID      string `json:"clientId"`
		ClientVersion string `json:"clientVersion"`
	} `json:"client"`
	ThreatInfo struct {
		ThreatTypes      []string `json:"threatTypes"`
		PlatformTypes    []string `json:"platformTypes"`
		ThreatEntryTypes []string `json:"threatEntryTypes"`
		ThreatEntries    []struct {
			URL string `json:"url"`
		} `json:"threatEntries"`
	} `json:"threatInfo"`
}

// SafeBrowsingResponse represents the response from Google Safe Browsing API
type SafeBrowsingResponse struct {
	Matches []struct {
		ThreatType      string `json:"threatType"`
		PlatformType    string `json:"platformType"`
		ThreatEntryType string `json:"threatEntryType"`
		Threat          struct {
			URL string `json:"url"`
		} `json:"threat"`
	} `json:"matches"`
}

// internal/services/blacklist.go
func NewBlacklistService() *BlacklistService {
	bs := &BlacklistService{
		patterns: []string{},
		regexRules: []*regexp.Regexp{
			regexp.MustCompile(`(?i)paypal\.com\..*`),
			regexp.MustCompile(`(?i)google\.com\..*`),
			regexp.MustCompile(`(?i)(\.pw|\.top|\.xyz)$`),
			// Updated IP address patterns
			regexp.MustCompile(`^https?://(\d{1,3}\.){3}\d{1,3}`), // IPv4
			regexp.MustCompile(`^https?://\[?([0-9a-fA-F:]+)\]?`), // IPv6
			regexp.MustCompile(`(?i)(admin|login|signin|banking|secure)\d+`),
		},
	}

	go bs.startPeriodicUpdates()
	return bs
}

func (bs *BlacklistService) updateBlacklists() error {
	safeBrowsingAPI := "https://safebrowsing.googleapis.com/v4/threatMatches:find"
	apiKey := os.Getenv("GCP_SAFE_BROWSING_API_KEY")
	fmt.Println("apiKey check: ", apiKey)
	// Prepare request body
	reqBody := SafeBrowsingRequest{
		Client: struct {
			ClientID      string `json:"clientId"`
			ClientVersion string `json:"clientVersion"`
		}{
			ClientID:      "dev4url",
			ClientVersion: "1.0.0",
		},
		ThreatInfo: struct {
			ThreatTypes      []string `json:"threatTypes"`
			PlatformTypes    []string `json:"platformTypes"`
			ThreatEntryTypes []string `json:"threatEntryTypes"`
			ThreatEntries    []struct {
				URL string `json:"url"`
			} `json:"threatEntries"`
		}{
			ThreatTypes:      []string{"MALWARE", "SOCIAL_ENGINEERING", "UNWANTED_SOFTWARE", "POTENTIALLY_HARMFUL_APPLICATION"},
			PlatformTypes:    []string{"ANY_PLATFORM"},
			ThreatEntryTypes: []string{"URL"},
			ThreatEntries: []struct {
				URL string `json:"url"`
			}{{URL: "http://example.com"}}, // You might want to batch check multiple URLs
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("POST", safeBrowsingAPI+"?key="+apiKey, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var sbResp SafeBrowsingResponse
	if err := json.NewDecoder(resp.Body).Decode(&sbResp); err != nil {
		return err
	}

	bs.mu.Lock()
	defer bs.mu.Unlock()

	// Update patterns from API response
	for _, match := range sbResp.Matches {
		bs.patterns = append(bs.patterns, match.Threat.URL)
	}

	bs.lastUpdated = time.Now()
	return nil
}

func (bs *BlacklistService) startPeriodicUpdates() {
	ticker := time.NewTicker(6 * time.Hour) // Update every 6 hours
	for range ticker.C {
		bs.updateBlacklists()
	}
}

func (bs *BlacklistService) IsURLBlacklisted(urlStr string) bool {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	// Debug output
	fmt.Printf("\nChecking URL: %s\n", urlStr)

	// Check against static patterns
	for _, pattern := range bs.patterns {
		if strings.Contains(urlStr, pattern) {
			fmt.Printf("Matched pattern: %s\n", pattern)
			return true
		}
	}

	// Check against regex rules
	for _, regex := range bs.regexRules {
		if regex.MatchString(urlStr) {
			fmt.Printf("Matched regex: %s\n", regex.String())
			return true
		}
	}

	fmt.Printf("No matches found\n")
	return false
}
