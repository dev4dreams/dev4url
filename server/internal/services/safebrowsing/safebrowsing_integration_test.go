package safebrowsing

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

func loadEnvFile(path string) (map[string]string, error) {
	envMap := make(map[string]string)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split on first = only
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		// Remove quotes if they exist
		value = strings.Trim(value, `"'`)

		envMap[key] = value
	}

	return envMap, scanner.Err()
}

func skipIfNoAPIKey(t *testing.T) string {
	// First try to get from environment
	apiKey := os.Getenv("GCP_SAFE_BROWSING_API_KEY")
	if apiKey != "" {
		return apiKey
	}

	// If not in environment, try to load from .env file
	envPath := "../../../.env"
	envVars, err := loadEnvFile(envPath)
	if err != nil {
		t.Skipf("Skipping integration test: Failed to load .env file from %s: %v", envPath, err)
	}

	apiKey, exists := envVars["GCP_SAFE_BROWSING_API_KEY"]
	if !exists || apiKey == "" {
		t.Skip("Skipping integration test: GCP_SAFE_BROWSING_API_KEY not found in .env file")
	}

	return apiKey
}

func TestIntegration_SafeBrowsingService_CheckURL(t *testing.T) {
	apiKey := skipIfNoAPIKey(t)
	service := NewSafeBrowsingService(apiKey)

	tests := []struct {
		name          string
		url           string
		expectError   bool
		errorContains string // Expected substring in error message
		expectMatches bool
	}{
		{
			name:          "safe url - google.com",
			url:           "https://www.google.com",
			expectError:   false,
			expectMatches: false,
		},
		{
			name:          "safe url - github.com",
			url:           "https://github.com",
			expectError:   false,
			expectMatches: false,
		},
		{
			name:          "test malware url",
			url:           "http://malware.testing.google.test/testing/malware/*",
			expectError:   true,
			expectMatches: true,
		},
		{
			name:          "invalid url scheme",
			url:           "ftp://example.com",
			expectError:   true,
			errorContains: "invalid URL scheme",
			expectMatches: false,
		},
		{
			name:          "invalid url format",
			url:           "not-a-valid-url",
			expectError:   true,
			errorContains: "invalid URL",
			expectMatches: false,
		},
		{
			name:          "missing scheme",
			url:           "example.com",
			expectError:   true,
			errorContains: "missing scheme",
			expectMatches: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := service.CheckURL(tt.url)

			// Check error expectations
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error containing %q, got %q", tt.errorContains, err.Error())
				}
				return // Don't check response if we expected an error
			}

			// If we didn't expect an error, but got one
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Check response
			if !tt.expectError {
				if response == nil {
					t.Fatal("Expected response but got nil")
				}
				hasMatches := len(response.Matches) > 0
				if hasMatches != tt.expectMatches {
					t.Errorf("Expected matches=%v, got matches=%v", tt.expectMatches, hasMatches)
				}

				if hasMatches {
					t.Logf("Threat matches found: %+v", response.Matches)
				}
			}
		})
	}
}

func TestIntegration_SafeBrowsingService_IsURLSafe(t *testing.T) {
	apiKey := skipIfNoAPIKey(t)
	service := NewSafeBrowsingService(apiKey)

	tests := []struct {
		name          string
		url           string
		expectError   bool
		errorContains string
		expectSafe    bool
	}{
		{
			name:        "safe url - microsoft.com",
			url:         "https://www.microsoft.com",
			expectError: false,
			expectSafe:  true,
		},
		{
			name:        "test social engineering url",
			url:         "http://social-engineering.testing.google.test/testing/social/*",
			expectError: true,
			expectSafe:  false,
		},
		{
			name:          "invalid url format",
			url:           "not-a-valid-url",
			expectError:   true,
			errorContains: "invalid URL",
			expectSafe:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isSafe, err := service.IsURLSafe(tt.url)

			// Check error expectations
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error containing %q, got %q", tt.errorContains, err.Error())
				}
				return
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if isSafe != tt.expectSafe {
				t.Errorf("Expected safe=%v, got safe=%v", tt.expectSafe, isSafe)
			}
		})
	}
}

func TestIntegration_RateLimiting(t *testing.T) {
	apiKey := skipIfNoAPIKey(t)
	service := NewSafeBrowsingService(apiKey)

	// Test multiple rapid requests to check rate limiting
	for i := 0; i < 5; i++ {
		_, err := service.IsURLSafe("https://www.example.com")
		if err != nil {
			t.Errorf("Request %d failed: %v", i+1, err)
		}
	}
}
