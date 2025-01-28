// utils/retry/retry_test.go
package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

var errTemporary = errors.New("temporary error")

func TestWithExponentialBackoff(t *testing.T) {
	tests := []struct {
		name          string
		operation     func(context.Context) (string, error)
		config        RetryConfig
		expectedRes   string
		expectedErr   error
		expectedCalls int
	}{
		{
			name: "successful_first_attempt",
			operation: func(ctx context.Context) (string, error) {
				return "success", nil
			},
			config: RetryConfig{
				MaxAttempts:  3,
				InitialDelay: 10 * time.Millisecond,
				MaxDelay:     50 * time.Millisecond,
				Multiplier:   2.0,
			},
			expectedRes:   "success",
			expectedErr:   nil,
			expectedCalls: 1,
		},
		{
			name:      "success_after_retries",
			operation: createRetryingOperation(3), // We'll define this helper function
			config: RetryConfig{
				MaxAttempts:  3,
				InitialDelay: 10 * time.Millisecond,
				MaxDelay:     50 * time.Millisecond,
				Multiplier:   2.0,
			},
			expectedRes:   "success",
			expectedErr:   nil,
			expectedCalls: 3,
		},
		{
			name:      "max_retries_exceeded",
			operation: createFailingOperation(),
			config: RetryConfig{
				MaxAttempts:  3,
				InitialDelay: 10 * time.Millisecond,
				MaxDelay:     50 * time.Millisecond,
				Multiplier:   2.0,
			},
			expectedRes:   "",
			expectedErr:   &MaxRetriesExceededError{Attempts: 3},
			expectedCalls: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := WithExponentialBackoff(context.Background(), tt.operation, tt.config)

			// Check result
			if result != tt.expectedRes {
				t.Errorf("Expected result %v, got %v", tt.expectedRes, result)
			}

			// Check error
			if tt.expectedErr != nil {
				var maxRetriesErr *MaxRetriesExceededError
				if !errors.As(err, &maxRetriesErr) {
					t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	}
}

// Helper functions to create test operations
func createRetryingOperation(successAttempt int) func(context.Context) (string, error) {
	attempts := 0
	return func(ctx context.Context) (string, error) {
		attempts++
		if attempts < successAttempt {
			return "", errTemporary
		}
		return "success", nil
	}
}

func createFailingOperation() func(context.Context) (string, error) {
	return func(ctx context.Context) (string, error) {
		return "", errTemporary
	}
}

// Add context cancellation test
func TestWithExponentialBackoffContextCancellation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	config := RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 1 * time.Second, // Long enough to trigger timeout
		MaxDelay:     5 * time.Second,
		Multiplier:   2.0,
	}

	operation := func(ctx context.Context) (string, error) {
		return "", errTemporary
	}

	_, err := WithExponentialBackoff(ctx, operation, config)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Expected deadline exceeded error, got %v", err)
	}
}
