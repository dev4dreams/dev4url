// internal/utils/validator.go
package utils

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"unicode"

	services "github.com/dev4dreams/dev4url/internal/services/blacklist"
	"golang.org/x/net/publicsuffix"
)

// URLValidator struct with blacklist service
type URLValidator struct {
	blacklistService *services.BlacklistService
}

// NewURLValidator creates a new validator instance
func NewURLValidator(blacklistService *services.BlacklistService) *URLValidator {
	return &URLValidator{
		blacklistService: blacklistService,
	}
}

// ValidateURL combines static rules and blacklist checking
func (v *URLValidator) ValidateURL(urlStr string) error {
	// Run static validation rules first
	if err := v.validateStaticRules(urlStr); err != nil {
		return err
	}

	// If blacklist service is configured, check against it
	if v.blacklistService != nil && v.blacklistService.IsURLBlacklisted(urlStr) {
		return fmt.Errorf("URL matches known malicious pattern")
	}

	return nil
}

// validateStaticRules contains all static validation logic
func (v *URLValidator) validateStaticRules(urlStr string) error {
	if strings.TrimSpace(urlStr) == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	// Parse the URL
	u, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL format: %v", err)
	}

	// Check scheme (protocol)
	if u.Scheme == "" {
		return fmt.Errorf("URL must have a scheme (http:// or https://)")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("URL scheme must be http or https")
	}

	// Check host
	if u.Host == "" {
		return fmt.Errorf("URL must have a host")
	}

	// Check for IP addresses
	if ip := net.ParseIP(u.Hostname()); ip != nil {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified() {
			return fmt.Errorf("IP-based URLs with private/local addresses are not allowed")
		}
	}

	// Length checks
	if len(urlStr) > 2048 {
		return fmt.Errorf("URL is too long (max 2048 characters)")
	}
	if len(u.Host) > 255 {
		return fmt.Errorf("hostname is too long (max 255 characters)")
	}

	// Check for suspicious patterns
	if err := v.checkSuspiciousPatterns(urlStr); err != nil {
		return err
	}

	// Check for control characters
	for _, r := range urlStr {
		if unicode.IsControl(r) {
			return fmt.Errorf("URL contains invalid control characters")
		}
	}

	// Check for phishing attempts
	if v.containsSuspiciousDomain(u.Host) {
		return fmt.Errorf("URL contains suspicious domain pattern")
	}

	return nil
}

// checkSuspiciousPatterns checks for known malicious patterns
func (v *URLValidator) checkSuspiciousPatterns(urlStr string) error {
	suspiciousPatterns := []string{
		"javascript:", "data:", "vbscript:", "file:", "ftp:",
		"<script", "alert(", "prompt(", "confirm(",
		"onload=", "onerror=", "../", "\\\\",
		".php?", ".asp?", "eval(", "exec(",
		"--", "DROP ", "UNION ", "%00", "0x00",
	}

	lowercaseURL := strings.ToLower(urlStr)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(lowercaseURL, pattern) {
			return fmt.Errorf("URL contains potentially malicious pattern: %s", pattern)
		}
	}
	return nil
}

// func (v *URLValidator) containsSuspiciousDomain(host string) bool {
// 	parsedDomain, err := publicsuffix.EffectiveTLDPlusOne(host)
// 	if err != nil {
// 		return false
// 	}

// 	fmt.Printf("\nChecking domain: %s\n", host)
// 	fmt.Printf("Parsed domain: %s\n", parsedDomain)
// 	suspiciousDomains := []string{
// 		"google",
// 		"facebook",
// 		"apple",
// 		"microsoft",
// 		"paypal",
// 	}

// 	parsedDomainLower := strings.ToLower(parsedDomain)
// 	fmt.Printf("Lowercase parsed domain: %s\n", parsedDomainLower)

// 	// Check for exact matches first
// 	for _, domain := range suspiciousDomains {
// 		// Check if domain name contains the suspicious word but isn't the legitimate domain
// 		domainWithTLD := domain + ".com"
// 		 fmt.Printf("Checking against %s\n", domainWithTLD)
// 		if strings.Contains(parsedDomainLower, domain) && parsedDomainLower != domainWithTLD {
// 			 fmt.Printf("Found suspicious pattern! Contains %s but isn't exactly %s\n", domain, domainWithTLD)
// 			return true
// 		}
// 	}

// 	return false
// }

func (v *URLValidator) containsSuspiciousDomain(host string) bool {
	parsedDomain, err := publicsuffix.EffectiveTLDPlusOne(host)
	if err != nil {
		fmt.Printf("Error parsing domain: %v\n", err)
		return false
	}

	fmt.Printf("\nChecking domain: %s\n", host)
	fmt.Printf("Parsed domain: %s\n", parsedDomain)

	suspiciousDomains := []string{
		"google.com",
		"facebook.com",
		"apple.com",
		"microsoft.com",
		"paypal.com",
	}

	hostLower := strings.ToLower(host)
	fmt.Printf("Host lower: %s\n", hostLower)

	for _, legitimate := range suspiciousDomains {
		legitimateBase := strings.TrimSuffix(legitimate, ".com")

		// Check for exact matches of legitimate domain
		if hostLower == legitimate {
			fmt.Printf("Exact match with %s - allowing\n", legitimate)
			return false
		}

		// Check for proper subdomains
		if strings.HasSuffix(hostLower, "."+legitimate) {
			fmt.Printf("Proper subdomain of %s - allowing\n", legitimate)
			return false
		}

		// Check for suspicious patterns
		// 1. Contains the brand name followed by suspicious words
		if strings.Contains(hostLower, legitimateBase) {
			fmt.Printf("Found suspicious use of %s in %s\n", legitimateBase, hostLower)
			return true
		}

		// 2. Contains the brand name in a subdomain of another domain
		if strings.Contains(hostLower, legitimate+".") {
			fmt.Printf("Found suspicious subdomain pattern with %s\n", legitimate)
			return true
		}
	}

	return false
}
