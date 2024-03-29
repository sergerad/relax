package relax

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// RoutineGroup is a wrapper around golang.org/x/sync/errgroup.Group that
// recovers from panics and returns them as errors
type RoutineGroup struct {
	*errgroup.Group
}

// NewGroup instantiates an RoutineGroup and corresponding context.
// This function should be used the same way as errgroup.WithContext() from
// golang.org/x/sync/errgroup
func NewGroup(ctx context.Context) (*RoutineGroup, context.Context) {
	errgroup, groupCtx := errgroup.WithContext(ctx)
	return &RoutineGroup{
		Group: errgroup,
	}, groupCtx
}

// Go runs a provided func in a goroutine while ensuring that
// any panic is recovered and returned as an error
func (rg *RoutineGroup) Go(f func() error) {
	rg.Group.Go(func() (err error) {
		// Handle panics
		defer func() {
			if r := recover(); r != nil {
				err = recoverError(r)
			}
		}()
		// Call the provided func
		return f()
	})
}
