package retry

import "fmt"

type MaxRetriesExceededError struct {
	Attempts int
}

func (e *MaxRetriesExceededError) Error() string {
	return fmt.Sprintf("maximum number of retries (%d) exceeded", e.Attempts)
}
