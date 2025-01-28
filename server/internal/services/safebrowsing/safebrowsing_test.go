package safebrowsing

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		expectError bool
	}{
		{
			name:        "valid http url",
			url:         "http://example.com",
			expectError: false,
		},
		{
			name:        "valid https url",
			url:         "https://example.com/path?query=value",
			expectError: false,
		},
		{
			name:        "invalid scheme",
			url:         "ftp://example.com",
			expectError: true,
		},
		{
			name:        "missing scheme",
			url:         "example.com",
			expectError: true,
		},
		{
			name:        "invalid format",
			url:         "not-a-valid-url",
			expectError: true,
		},
		{
			name:        "empty url",
			url:         "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateURL(tt.url)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestSafeBrowsingService_CheckURL(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		mockResponse   interface{}
		mockStatusCode int
		expectError    bool
		expectMatches  int
	}{
		{
			name: "safe url",
			url:  "https://example.com",
			mockResponse: ThreatResponse{
				Matches: []ThreatMatch{},
			},
			mockStatusCode: http.StatusOK,
			expectError:    false,
			expectMatches:  0,
		},
		{
			name: "unsafe url",
			url:  "http://malicious.example.com",
			mockResponse: ThreatResponse{
				Matches: []ThreatMatch{
					{
						ThreatType:      "MALWARE",
						PlatformType:    "ANY_PLATFORM",
						ThreatEntryType: "URL",
						Threat: ThreatEntry{
							URL: "http://malicious.example.com",
						},
						CacheDuration: "300.000s",
					},
				},
			},
			mockStatusCode: http.StatusOK,
			expectError:    false,
			expectMatches:  1,
		},
		{
			name:           "api error",
			url:            "https://example.com",
			mockResponse:   map[string]interface{}{"error": "API Error"},
			mockStatusCode: http.StatusInternalServerError,
			expectError:    true,
			expectMatches:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request method
				if r.Method != http.MethodPost {
					t.Errorf("Expected POST request, got %s", r.Method)
				}

				// Verify content type
				if r.Header.Get("Content-Type") != "application/json" {
					t.Errorf("Expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
				}

				// Set response status code
				w.WriteHeader(tt.mockStatusCode)

				// Write mock response
				json.NewEncoder(w).Encode(tt.mockResponse)
			}))
			defer server.Close()

			// Create service with test server URL
			service := &SafeBrowsingService{
				apiKey:  "test-api-key",
				baseURL: server.URL,
				httpClient: &http.Client{
					Timeout: 5 * time.Second,
				},
			}

			// Make request
			response, err := service.CheckURL(tt.url)

			// Check error expectation
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// If we expect a successful response, verify the matches
			if !tt.expectError {
				if response == nil {
					t.Fatal("Expected response but got nil")
				}
				if len(response.Matches) != tt.expectMatches {
					t.Errorf("Expected %d matches, got %d", tt.expectMatches, len(response.Matches))
				}
			}
		})
	}
}

func TestSafeBrowsingService_IsURLSafe(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		mockResponse   interface{}
		mockStatusCode int
		expectError    bool
		expectSafe     bool
	}{
		{
			name: "safe url",
			url:  "https://example.com",
			mockResponse: ThreatResponse{
				Matches: []ThreatMatch{},
			},
			mockStatusCode: http.StatusOK,
			expectError:    false,
			expectSafe:     true,
		},
		{
			name: "unsafe url",
			url:  "http://malicious.example.com",
			mockResponse: ThreatResponse{
				Matches: []ThreatMatch{
					{
						ThreatType:      "MALWARE",
						PlatformType:    "ANY_PLATFORM",
						ThreatEntryType: "URL",
						Threat: ThreatEntry{
							URL: "http://malicious.example.com",
						},
					},
				},
			},
			mockStatusCode: http.StatusOK,
			expectError:    false,
			expectSafe:     false,
		},
		{
			name:           "api error",
			url:            "https://example.com",
			mockResponse:   map[string]interface{}{"error": "API Error"},
			mockStatusCode: http.StatusInternalServerError,
			expectError:    true,
			expectSafe:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockStatusCode)
				json.NewEncoder(w).Encode(tt.mockResponse)
			}))
			defer server.Close()

			// Create service with test server URL
			service := &SafeBrowsingService{
				apiKey:  "test-api-key",
				baseURL: server.URL,
				httpClient: &http.Client{
					Timeout: 5 * time.Second,
				},
			}

			// Make request
			isSafe, err := service.IsURLSafe(tt.url)

			// Check error expectation
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Check safety status
			if !tt.expectError && isSafe != tt.expectSafe {
				t.Errorf("Expected safe=%v, got %v", tt.expectSafe, isSafe)
			}
		})
	}
}

func TestNewSafeBrowsingService(t *testing.T) {
	apiKey := "test-api-key"
	service := NewSafeBrowsingService(apiKey)

	if service == nil {
		t.Fatal("Expected non-nil service")
	}

	if service.apiKey != apiKey {
		t.Errorf("Expected apiKey=%s, got %s", apiKey, service.apiKey)
	}

	if service.baseURL != defaultBaseURL {
		t.Errorf("Expected baseURL=%s, got %s", defaultBaseURL, service.baseURL)
	}

	if service.httpClient == nil {
		t.Error("Expected non-nil HTTP client")
	}

	if service.httpClient.Timeout != 10*time.Second {
		t.Errorf("Expected timeout=10s, got %v", service.httpClient.Timeout)
	}
}
