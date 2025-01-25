package utils

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strings"
)

// URLValidator handles URL validation with configurable rules
type URLValidator struct {
	config *Config
}

type URLValidatorInterface interface {
	ValidateURL(ctx context.Context, urlStr string) *ValidationResult
}

// Config holds validation configuration
type Config struct {
	MaxURLLength    int      `json:"maxUrlLength"`
	AllowedDomains  []string `json:"allowedDomains"`
	BlockedPatterns []string `json:"blockedPatterns"`
	BlockedDomains  []string `json:"blockedDomains"`
}

// ValidationResult contains the validation outcome and any errors
type ValidationResult struct {
	IsValid bool     `json:"isValid"`
	Errors  []string `json:"errors,omitempty"`
}

// DefaultConfig provides sensible default settings
func DefaultConfig() *Config {
	return &Config{
		MaxURLLength: 2048,
		BlockedPatterns: []string{
			"javascript:", "data:", "vbscript:",
			"<script", "alert(", "prompt(",
			"onload=", "onerror=",
			"eval(", "exec(", "file:", "ftp:",
			"confirm(", "../", "\\\\",
			".php?", ".asp?", "eval(", "exec(",
			"--", "DROP ", "UNION ", "%00", "0x00",
		},
		BlockedDomains: []string{
			"example.com",
			"test.com",
			"localhost",
		},
	}
}

// NewURLValidator creates a new validator instance
func NewURLValidator(config *Config) *URLValidator {
	if config == nil {
		config = DefaultConfig()
	}
	return &URLValidator{
		config: config,
	}
}

// ValidateURL performs comprehensive URL validation
func (v *URLValidator) ValidateURL(ctx context.Context, urlStr string) *ValidationResult {
	result := &ValidationResult{
		IsValid: true,
		Errors:  make([]string, 0),
	}

	// Perform all validations
	if err := v.validateBasics(urlStr); err != nil {
		result.Errors = append(result.Errors, err.Error())
	}

	if err := v.validateSecurity(urlStr); err != nil {
		result.Errors = append(result.Errors, err.Error())
	}

	if err := v.validateDomain(urlStr); err != nil {
		result.Errors = append(result.Errors, err.Error())
	}

	// Set final validity
	result.IsValid = len(result.Errors) == 0

	return result
}

// validateBasics checks fundamental URL properties
func (v *URLValidator) validateBasics(urlStr string) error {
	if strings.TrimSpace(urlStr) == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	if len(urlStr) > v.config.MaxURLLength {
		return fmt.Errorf("URL exceeds maximum length of %d characters", v.config.MaxURLLength)
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("URL scheme must be http or https")
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("URL must have a host")
	}

	return nil
}

// validateSecurity performs security-related checks
func (v *URLValidator) validateSecurity(urlStr string) error {
	urlLower := strings.ToLower(urlStr)

	// Check for suspicious patterns
	for _, pattern := range v.config.BlockedPatterns {
		if strings.Contains(urlLower, strings.ToLower(pattern)) {
			return fmt.Errorf("URL contains suspicious pattern: %s", pattern)
		}
	}

	// Check for control characters
	for _, r := range urlStr {
		if r < 32 || r == 127 {
			return fmt.Errorf("URL contains invalid control characters")
		}
	}

	return nil
}

// validateDomain performs domain-specific validation
func (v *URLValidator) validateDomain(urlStr string) error {
	parsedURL, _ := url.Parse(urlStr)
	hostname := parsedURL.Hostname()

	// Block localhost in all forms
	if hostname == "localhost" || strings.HasSuffix(hostname, ".localhost") {
		return fmt.Errorf("localhost URLs are not allowed")
	}
	// Check for IP addresses
	if ip := net.ParseIP(parsedURL.Hostname()); ip != nil {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified() {
			return fmt.Errorf("IP-based URLs with private/local addresses are not allowed")
		}
		return nil
	}

	// Check blocked domains
	for _, blockedDomain := range v.config.BlockedDomains {
		if strings.Contains(parsedURL.Hostname(), blockedDomain) {
			return fmt.Errorf("domain is blocked")
		}
	}

	// Check allowed domains if configured
	if len(v.config.AllowedDomains) > 0 {
		allowed := false
		for _, allowedDomain := range v.config.AllowedDomains {
			if parsedURL.Hostname() == allowedDomain || strings.HasSuffix(parsedURL.Hostname(), "."+allowedDomain) {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("domain not in allowed list")
		}
	}

	return nil
}
