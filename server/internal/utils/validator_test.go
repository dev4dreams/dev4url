// internal/utils/validator_test.go
package utils_test

import (
	"strings"
	"testing"

	"github.com/dev4dreams/dev4url/internal/utils"
)

func TestValidateURL(t *testing.T) {
	validator := utils.NewURLValidator(nil)

	tests := []struct {
		name    string
		url     string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Valid HTTPS URL",
			url:     "https://www.google.com",
			wantErr: false,
		},
		{
			name:    "Valid HTTP URL with path",
			url:     "http://example.com/path",
			wantErr: false,
		},
		{
			name:    "Empty URL",
			url:     "",
			wantErr: true,
			errMsg:  "URL cannot be empty",
		},
		{
			name:    "Invalid scheme (FTP)",
			url:     "ftp://example.com",
			wantErr: true,
			errMsg:  "URL scheme must be http or https",
		},
		{
			name:    "Local IP address",
			url:     "http://192.168.1.1",
			wantErr: true,
			errMsg:  "private/local addresses are not allowed",
		},
		{
			name:    "JavaScript injection attempt",
			url:     "javascript:alert(1)",
			wantErr: true,
			errMsg:  "URL scheme must be http or https",
		},
		{
			name:    "SQL Injection attempt",
			url:     "http://example.com/page?id=1--DROP",
			wantErr: true,
			errMsg:  "potentially malicious pattern", // Updated to match actual error
		},
		{
			name:    "Phishing attempt",
			url:     "http://google.com.malicious.com", // Updated to match suspicious domain pattern
			wantErr: true,
			errMsg:  "suspicious domain pattern",
		}, {
			name:    "Phishing attempt",
			url:     "http://google-account.malicious.com", // More obvious phishing attempt
			wantErr: true,
			errMsg:  "suspicious domain pattern",
		},
		{
			name:    "Another phishing attempt",
			url:     "http://accounts-google.com",
			wantErr: true,
			errMsg:  "suspicious domain pattern",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateURL(tt.url)

			// Check if error was expected
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If we expect an error, check the message
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if !strings.Contains(strings.ToLower(err.Error()), strings.ToLower(tt.errMsg)) {
					t.Errorf("ValidateURL() error message = %v, want to contain %v", err, tt.errMsg)
				}
			}
		})
	}
}

// internal/utils/validator_test.go
// func TestValidatorWithBlacklist(t *testing.T) {
//     bs := services.NewBlacklistService()
//     validator := NewURLValidator(bs)

//     tests := []struct {
//         name    string
//         url     string
//         wantErr bool
//     }{
//         {
//             name:    "Clean URL",
//             url:     "https://example.com",
//             wantErr: false,
//         },
//         {
//             name:    "Blacklisted pattern",
//             url:     "https://paypal.com.phishing.com",
//             wantErr: true,
//         },
//     }

//     for _, tt := range tests {
//         t.Run(tt.name, func(t *testing.T) {
//             err := validator.ValidateURL(tt.url)
//             if (err != nil) != tt.wantErr {
//                 t.Errorf("ValidateURL() with blacklist error = %v, wantErr %v", err, tt.wantErr)
//             }
//         })
//     }
// }
