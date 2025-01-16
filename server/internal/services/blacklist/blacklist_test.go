// internal/services/blacklist_test.go
package services_test

import (
	"testing"

	services "github.com/dev4dreams/dev4url/internal/services/blacklist"
)

func TestBlacklistService(t *testing.T) {
	bs := services.NewBlacklistService()

	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{
			name:     "Clean URL",
			url:      "https://example.com",
			expected: false,
		},
		{
			name:     "PayPal phishing attempt",
			url:      "https://paypal.com.suspicious.com",
			expected: true,
		},
		{
			name:     "Suspicious TLD",
			url:      "https://login.xyz",
			expected: true,
		},
		{
			name:     "Admin phishing page",
			url:      "https://admin123.com",
			expected: true,
		},
		{
			name:     "Raw IP address",
			url:      "http://192.168.1.1",
			expected: true,
		},
		// Add more specific test cases
		{
			name:     "IPv4 with path",
			url:      "http://192.168.1.1/login",
			expected: true,
		},
		{
			name:     "IPv6 address",
			url:      "http://[2001:db8::1]",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bs.IsURLBlacklisted(tt.url)
			if result != tt.expected {
				t.Errorf("IsURLBlacklisted() = %v, want %v\nURL: %s", result, tt.expected, tt.url)
			}
		})
	}
}
