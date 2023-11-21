package relax

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

var (
	// PanicError is returned when a panic is recovered
	// during the execution of a goroutine
	PanicError = fmt.Errorf("recovered from panic")
)

// ErrorGroup is a wrapper around golang.org/x/sync/errgroup.Group that
// recovers from panics and returns them as errors
type ErrorGroup struct {
	*errgroup.Group
}

// NewErrorGroup instantiates an ErrorGroup and corresponding context.
// This function should be used the same way as errgroup.WithContext() from
// golang.org/x/sync/errgroup
func NewErrorGroup(ctx context.Context) (*ErrorGroup, context.Context) {
	errgroup, groupCtx := errgroup.WithContext(ctx)
	return &ErrorGroup{
		Group: errgroup,
	}, groupCtx
}

// Go runs a provided func in a goroutine while ensuring that
// any panic is recovered and returned as an error
func (g *ErrorGroup) Go(f func() error) {
	g.Group.Go(func() (err error) {
		// Define a recover func that converts a panic to an error
		recoverFunc := func() {
			if r := recover(); r != nil {
				// Assign the panic content to returned error
				err = fmt.Errorf("%w: %v", PanicError, r)
			}
		}
		// Handle panics
		defer recoverFunc()
		// Call the provided func
		return f()
	})
}
