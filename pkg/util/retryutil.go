//  Copyright (c) 2019 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package util

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const (
	retryErrorValue string = "maxRetries should be > 0"
)

type RetryError struct {
	N int
	E error
}

type RetryOkError error

var ErrMaxRetry = fmt.Errorf(retryErrorValue)

func (e *RetryError) Error() string {
	return fmt.Sprintf("still failing after %d retries: %v", e.N, e.E)
}

func IsRetryFailure(err error) bool {
	var e *RetryError
	ok := errors.As(err, &e)

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
		return ErrMaxRetry
	}

	tick := time.NewTicker(interval)

	defer tick.Stop()

	var err error

	var ok bool

	for i := 0; ; i++ {
		ok, err = f()
		if err != nil {
			// Ignore error's when expected during retryOnErr
			var e *RetryError
			shouldRetry := errors.As(err, &e)

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
			return fmt.Errorf("%v: %w", ctx.Err(), err)
		}
	}

	return &RetryError{N: maxRetries, E: err}
}
