package relax

import (
	"fmt"

	"errors"
)

// Routine is a handle to a goroutine's error response.
// Panics that occur in this routine are converted into errors.
type Routine struct {
	errChan chan error
}

// Wait blocks until the goroutine corresponding to the Routine instance returns an error.
// Once an error is returned, all subsequent calls to Wait will return nil.
func (r *Routine) Wait() error {
	return <-r.errChan
}

// Release will call the handler against an error returned by this routine.
// Once a routine is released, waiting on it will return nil.
func (r *Routine) Release(handler func(error)) {
	go func() {
		handler(r.Wait())
	}()
}

// Go launches a goroutine that will return an error if the provided func panics or
// an error is returned by the provided func.
func Go(f func() error) *Routine {
	routine := &Routine{
		errChan: make(chan error, 1),
	}
	go func() {
		// Always close
		defer close(routine.errChan)
		// Handle panics
		defer func() {
			if r := recover(); r != nil {
				routine.errChan <- recoverError(r)
			}
		}()
		// Call the provided func
		routine.errChan <- f()
	}()
	return routine
}

// recoverError will transform a recovered panic datum to an error
func recoverError(r any) error {
	switch x := r.(type) {
	case error:
		return errors.Join(x, PanicError)
	case nil:
		return nil
	default:
		return fmt.Errorf("%w: %v", PanicError, r)
	}
}
