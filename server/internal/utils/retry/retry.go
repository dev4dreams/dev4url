// retry/retry.go
package retry

import (
	"context"
	"errors"
	"time"
)

type RetryConfig struct {
	MaxAttempts  int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
}

var ErrMaxRetriesExceeded = errors.New("maximum number of retries exceeded")

func WithExponentialBackoff[T any](
	ctx context.Context,
	operation func(context.Context) (T, error),
	config RetryConfig,
) (T, error) {
	var result T
	currentDelay := config.InitialDelay

	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		var err error
		result, err = operation(ctx)

		if err == nil {
			return result, nil
		}

		// Check if error is retryable
		if !isRetryableError(err) {
			return result, err
		}

		// Check context cancellation
		if ctx.Err() != nil {
			return result, ctx.Err()
		}

		// Wait before next retry
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		case <-time.After(currentDelay):
			currentDelay = time.Duration(float64(currentDelay) * config.Multiplier)
			if currentDelay > config.MaxDelay {
				currentDelay = config.MaxDelay
			}
		}
	}

	return result, &MaxRetriesExceededError{Attempts: config.MaxAttempts}

}

func isRetryableError(err error) bool {
	// Define which errors should trigger a retry
	// e.g., temporary network issues, rate limits, etc.
	// You can customize this based on SafeBrowsing API error responses
	return true // Implement your logic here
}

// utils/retry/retry.go

func (c RetryConfig) Validate() error {
	if c.MaxAttempts <= 0 {
		return errors.New("max attempts must be greater than 0")
	}
	if c.InitialDelay <= 0 {
		return errors.New("initial delay must be greater than 0")
	}
	if c.MaxDelay < c.InitialDelay {
		return errors.New("max delay must be greater than or equal to initial delay")
	}
	if c.Multiplier <= 1.0 {
		return errors.New("multiplier must be greater than 1.0")
	}
	return nil
}
