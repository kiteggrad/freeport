package freeport

import (
	"context"

	"github.com/cenkalti/backoff/v5"
)

type Retrier struct {
	generator *Generator
	backoff   backoff.BackOff
}

type retryOperation func(port Port) error

func getDefaultBackoff() *backoff.ExponentialBackOff {
	return backoff.NewExponentialBackOff()
}

// NewRetrier - returns new Retrier.
//   - if gen == nil - uses global generator
//   - if backoff == nil - uses default exponential backoff
func NewRetrier(gen *Generator, bo backoff.BackOff) *Retrier {
	if gen == nil {
		gen = globGenerator
	}
	if bo == nil {
		bo = getDefaultBackoff()
	}

	return &Retrier{
		generator: gen,
		backoff:   bo,
	}
}

// Retry the operation o until it does not return error or BackOff stops.
// Gets a new port for each operation.
//
// Returns last requested port or error.
//
// Uses backoff.Retry, see it's details inside.
func (r *Retrier) Retry(ctx context.Context, o retryOperation) (port Port, err error) {
	port, err = backoff.Retry(ctx, func() (port Port, err error) {
		port, err = r.generator.Get()
		if err != nil {
			return 0, err
		}

		err = o(port)
		if err != nil {
			return 0, err
		}

		return port, nil
	}, backoff.WithBackOff(r.backoff))
	if err != nil {
		return 0, err
	}

	return port, nil
}
