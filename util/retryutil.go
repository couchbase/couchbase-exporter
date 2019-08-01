package util

import (
	"context"
	"fmt"
	"time"
)

type RetryError struct {
	n int
	e error
}

type RetryOkError error

func (e *RetryError) Error() string {
	return fmt.Sprintf("still failing after %d retries: %v", e.n, e.e)
}

func IsRetryFailure(err error) bool {
	_, ok := err.(*RetryError)
	return ok
}

type ConditionFunc func() error
type RetryFunc func() (bool, error)

// Retry retries f every interval until after maxRetries.
// The interval won't be affected by how long f takes.
// For example, if interval is 3s, f takes 1s, another f will be called 2s later.
// However, if f takes longer than interval, it will be delayed.
func Retry(ctx context.Context, interval time.Duration, maxRetries int, f RetryFunc) error {
	if maxRetries <= 0 {
		return fmt.Errorf("maxRetries (%d) should be > 0", maxRetries)
	}
	tick := time.NewTicker(interval)
	defer tick.Stop()

	var err error
	var ok bool

	for i := 0; ; i++ {
		ok, err = f()
		if err != nil {
			// Ignore error's when expected during retryOnErr
			_, shouldRetry := err.(RetryOkError)
			if !shouldRetry {
				return err
			}
		}
		if ok {
			return nil
		}
		if i == maxRetries {
			break
		}
		select {
		case <-tick.C:
		case <-ctx.Done():
			return fmt.Errorf("%v: %v", ctx.Err(), err)
		}
	}
	return &RetryError{n: maxRetries, e: err}
}