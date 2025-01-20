package utils

import (
	"context"
	"testing"
)

func TestURLValidator(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		config   *Config
		url      string
		expected bool
		errors   []string
	}{
		{
			name:     "Valid HTTPS URL",
			config:   DefaultConfig(),
			url:      "https://google.com/path?param=value",
			expected: true,
			errors:   nil,
		},
		{
			name:     "Empty URL",
			config:   DefaultConfig(),
			url:      "",
			expected: false,
			errors:   []string{"URL cannot be empty"},
		},
		{
			name: "URL Exceeding Max Length",
			config: &Config{
				MaxURLLength:    20,
				BlockedPatterns: []string{},
				BlockedDomains:  []string{},
			},
			url:      "https://verylongdomainname.com/path",
			expected: false,
			errors:   []string{"URL exceeds maximum length of 20 characters"},
		},
		{
			name:     "Invalid Scheme",
			config:   DefaultConfig(),
			url:      "ftp://example.com",
			expected: false,
			errors: []string{
				"URL scheme must be http or https",
				"URL contains suspicious pattern: ftp:",
				"domain is blocked",
			},
		},
		{
			name:     "Blocked Domain",
			config:   DefaultConfig(),
			url:      "https://example.com",
			expected: false,
			errors:   []string{"domain is blocked"},
		},
		{
			name:     "Suspicious Pattern",
			config:   DefaultConfig(),
			url:      "https://domain.com/path?script=<script>alert('xss')</script>",
			expected: false,
			errors:   []string{"URL contains suspicious pattern: <script"},
		},
		{
			name: "Allowed Domains Test",
			config: &Config{
				MaxURLLength:    2048,
				AllowedDomains:  []string{"trusted.com"},
				BlockedPatterns: DefaultConfig().BlockedPatterns,
				BlockedDomains:  DefaultConfig().BlockedDomains,
			},
			url:      "https://untrusted.com",
			expected: false,
			errors:   []string{"domain not in allowed list"},
		},
		{
			name:     "Local IP Not Allowed",
			config:   DefaultConfig(),
			url:      "http://127.0.0.1/admin",
			expected: false,
			errors:   []string{"IP-based URLs with private/local addresses are not allowed"},
		},
		{
			name:     "Private IP Not Allowed",
			config:   DefaultConfig(),
			url:      "http://192.168.1.1/admin",
			expected: false,
			errors:   []string{"IP-based URLs with private/local addresses are not allowed"},
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewURLValidator(tt.config)
			result := validator.ValidateURL(context.Background(), tt.url)

			if result.IsValid != tt.expected {
				t.Errorf("ValidateURL() got = %v, want %v", result.IsValid, tt.expected)
			}

			if tt.errors != nil {
				if len(result.Errors) != len(tt.errors) {
					t.Errorf("ValidateURL() got %d errors, want %d errors. Got errors: %v, want errors: %v",
						len(result.Errors), len(tt.errors), result.Errors, tt.errors)
				}

				for i, expectedErr := range tt.errors {
					if i < len(result.Errors) && result.Errors[i] != expectedErr {
						t.Errorf("ValidateURL() error = %v, want %v", result.Errors[i], expectedErr)
					}
				}
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.MaxURLLength != 2048 {
		t.Errorf("DefaultConfig() MaxURLLength = %d, want %d", config.MaxURLLength, 2048)
	}

	if len(config.BlockedPatterns) == 0 {
		t.Error("DefaultConfig() BlockedPatterns is empty")
	}

	if len(config.BlockedDomains) == 0 {
		t.Error("DefaultConfig() BlockedDomains is empty")
	}
}

func TestNewURLValidator(t *testing.T) {
	t.Run("With Custom Config", func(t *testing.T) {
		config := &Config{
			MaxURLLength:   1000,
			BlockedDomains: []string{"blocked.com"},
		}
		validator := NewURLValidator(config)
		if validator.config != config {
			t.Error("NewURLValidator() did not set custom config correctly")
		}
	})

	t.Run("With Nil Config", func(t *testing.T) {
		validator := NewURLValidator(nil)
		if validator.config == nil {
			t.Error("NewURLValidator() did not set default config when nil was provided")
		}
	})
}
