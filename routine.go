package relax

import (
	"fmt"
)

// Routine is a handle to a goroutine's error response
type Routine struct {
	errChan chan error
}

// Wait blocks until the goroutine corresponding to the Routine instance returns an error
func (r *Routine) Wait() error {
	return <-r.errChan
}

// Go launches a goroutine that will return an error if the provided func panics
func Go(f func() error) *Routine {
	routine := &Routine{
		errChan: make(chan error, 1),
	}
	go func() {
		// Define a recover func that converts a panic to an error
		recoverFunc := func() {
			if r := recover(); r != nil {
				// Assign the panic content to returned error
				routine.errChan <- fmt.Errorf("%w: %v", PanicError, r)
			}
		}
		// Always close
		defer close(routine.errChan)
		// Handle panics
		defer recoverFunc()
		// Call the provided func
		routine.errChan <- f()
	}()
	return routine
}
