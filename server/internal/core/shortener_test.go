package core

import (
	"fmt"
	"sync"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	tests := []struct {
		name      string
		workerID  int64
		wantError bool
	}{
		{
			name:      "Valid worker ID",
			workerID:  1,
			wantError: false,
		},
		{
			name:      "Worker ID zero",
			workerID:  0,
			wantError: false,
		},
		{
			name:      "Invalid worker ID (negative)",
			workerID:  -1,
			wantError: true,
		},
		{
			name:      "Invalid worker ID (too large)",
			workerID:  256, // exceeds 8 bits
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generator, err := NewGenerator(tt.workerID)
			if tt.wantError {
				if err == nil {
					t.Errorf("NewGenerator() expected error for workerID %d", tt.workerID)
				}
			} else {
				if err != nil {
					t.Errorf("NewGenerator() unexpected error: %v", err)
				}
				if generator == nil {
					t.Error("NewGenerator() returned nil generator")
				}
			}
		})
	}
}

func TestGenerateShortURL(t *testing.T) {
	generator, err := NewGenerator(1)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	// Test URL length
	shortURL, err := generator.GenerateShortURL()
	fmt.Println("new short url : ", shortURL)
	if err != nil {
		t.Errorf("GenerateShortURL() unexpected error: %v", err)
	}
	if len(shortURL) != 7 {
		t.Errorf("GenerateShortURL() got length %d, want 7", len(shortURL))
	}

	// Test character set
	for _, char := range shortURL {
		if !isValidCharacter(char) {
			t.Errorf("GenerateShortURL() contains invalid character: %c", char)
		}
	}
}

func TestURLUniqueness(t *testing.T) {
	generator, err := NewGenerator(1)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	urlSet := make(map[string]bool)
	iterations := 1000

	for i := 0; i < iterations; i++ {
		shortURL, err := generator.GenerateShortURL()
		if err != nil {
			t.Errorf("GenerateShortURL() iteration %d error: %v", i, err)
			continue
		}

		if urlSet[shortURL] {
			t.Errorf("Duplicate URL generated: %s", shortURL)
		}
		urlSet[shortURL] = true
	}
}

func TestConcurrentGeneration(t *testing.T) {
	generator, err := NewGenerator(1)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	var wg sync.WaitGroup
	urlSet := sync.Map{}
	goroutines := 10
	iterations := 100

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				shortURL, err := generator.GenerateShortURL()
				if err != nil {
					t.Errorf("GenerateShortURL() error: %v", err)
					continue
				}

				// Check for duplicates
				if _, loaded := urlSet.LoadOrStore(shortURL, true); loaded {
					t.Errorf("Duplicate URL generated in concurrent execution: %s", shortURL)
				}
			}
		}()
	}

	wg.Wait()
}

func TestNoAmbiguousCharacters(t *testing.T) {
	generator, err := NewGenerator(1)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	ambiguousChars := []rune{'l', 'I', '0', 'O'}
	iterations := 1000

	for i := 0; i < iterations; i++ {
		shortURL, err := generator.GenerateShortURL()
		if err != nil {
			t.Errorf("GenerateShortURL() error: %v", err)
			continue
		}

		for _, char := range shortURL {
			for _, ambiguous := range ambiguousChars {
				if char == ambiguous {
					t.Errorf("Found ambiguous character %c in URL %s", char, shortURL)
				}
			}
		}
	}
}

// Helper function to check if a character is in our allowed set
func isValidCharacter(char rune) bool {
	for _, validChar := range alphabet {
		if char == validChar {
			return true
		}
	}
	return false
}

func TestURLLength(t *testing.T) {
	generator, err := NewGenerator(1)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	// Test multiple URLs to ensure consistent length
	for i := 0; i < 1000; i++ {
		shortURL, err := generator.GenerateShortURL()
		if err != nil {
			t.Errorf("GenerateShortURL() unexpected error: %v", err)
			continue
		}

		if len(shortURL) != 7 {
			t.Errorf("URL length = %d, want 7. URL: %s", len(shortURL), shortURL)
		}
	}
}
