package freeport

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/cenkalti/backoff/v4"
)

// Get - returns free port that wasn't used before (from this package).
func Get() (port int, err error) {
	return globGenerator.Get()
}

// MustGet - returns free port that wasn't used before (from this package) or panics.
func MustGet() (port int) {
	return globGenerator.MustGet()
}

// Retry the operation o until it does not return error or BackOff stops.
// Gets a new free port that wasn't used before (from this package) for each operation.
//
// Returns last requested port or error.
//
// Uses backoff.Retry, see it's details inside.
func Retry(o func(port int) error) (port int, err error) {
	return RetryBackoff(o, nil)
}

// RetryBackoff same as `Retry`, but you can pass backoff.
func RetryBackoff(o retryOperation, bo backoff.BackOff) (port int, err error) {
	retrier := NewRetrier(nil, bo)

	port, err = retrier.Retry(o)
	if err != nil {
		return 0, err
	}

	return port, nil
}

// RetryTimeout same as `Retry`, but you can pass timeout.
func RetryTimeout(o retryOperation, timeout time.Duration) (port int, err error) {
	bo := getDefaultBackoff()
	bo.MaxElapsedTime = timeout

	retrier := NewRetrier(nil, bo)

	port, err = retrier.Retry(o)
	if err != nil {
		return 0, err
	}

	return port, nil
}

// RetryCtx same as `Retry`, but you can pass context.
func RetryCtx(ctx context.Context, o retryOperation) (port int, err error) {
	bo := backoff.WithContext(getDefaultBackoff(), ctx)

	retrier := NewRetrier(nil, bo)

	port, err = retrier.Retry(o)
	if err != nil {
		return 0, err
	}

	return port, nil
}

func tryGetFreeport() (port int, err error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, fmt.Errorf("failed to net.ResolveTCPAddr: %w", err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, fmt.Errorf("failed to net.ListenTCP: %w", err)
	}
	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port, nil
}
