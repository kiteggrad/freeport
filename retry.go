package freeport

import "github.com/cenkalti/backoff/v4"

type Retrier struct {
	generator *Generator
	backoff   backoff.BackOff
}

type retryOperation func(port int) error

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
func (r *Retrier) Retry(o retryOperation) (port int, err error) {
	err = backoff.Retry(func() error {
		port, err = r.generator.Get()
		if err != nil {
			return err
		}

		err = o(port)
		if err != nil {
			return err
		}

		return nil
	}, r.backoff)
	if err != nil {
		return 0, err
	}

	return port, nil
}
