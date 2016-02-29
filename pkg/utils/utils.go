package utils

import (
	"fmt"
	"time"
)

var (
	PollingInterval = 30
	MaxAttempts     = 180
)

func WaitFor(f func() bool) error {
	return WaitForSpecific(f, MaxAttempts, time.Duration(PollingInterval)*time.Second)
}

func WaitForSpecific(f func() bool, maxAttempts int, interval time.Duration) error {
	return WaitForSpecificOrError(func() (bool, error) {
		return f(), nil
	}, maxAttempts, interval)
}

func WaitForSpecificOrError(f func() (bool, error), maxAttempts int, interval time.Duration) error {
	for i := 0; i < maxAttempts; i++ {
		done, err := f()

		if err != nil {
			return err
		}

		if done {
			return nil
		}

		time.Sleep(interval)
	}

	return fmt.Errorf("Maximum number of retries (%d) exceeded", maxAttempts)
}
