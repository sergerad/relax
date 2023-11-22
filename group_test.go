package relax

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGroup_CancelParentContext_ChildContextDone(t *testing.T) {
	// Root context to cancel
	ctx, cancel := context.WithCancel(context.Background())

	// Group to wait on root and group context Done()
	e, gCtx := NewGroup(ctx)
	e.Go(func() error {
		<-ctx.Done()
		<-gCtx.Done()
		return nil
	})

	// Cancel root context
	cancel()

	// Validate
	assert.NoError(t, e.Wait())
}

func TestGroup_Panic_Error(t *testing.T) {
	// Group
	e, ctx := NewGroup(context.Background())

	// Routine that panics
	panicMsg := "test panic"
	e.Go(func() error {
		panic(panicMsg)
	})

	// Sibling routine that blocks on Group context
	e.Go(func() error {
		<-ctx.Done()
		return nil
	})

	// Wait for all goroutines
	err := e.Wait()
	require.Error(t, err)
	// Verify panic message/error is returned
	assert.Contains(t, err.Error(), panicMsg)
	assert.True(t, errors.Is(err, PanicError))
}

func TestGroup_NoPanic_NoError(t *testing.T) {
	e, _ := NewGroup(context.Background())
	e.Go(func() error {
		return nil
	})
	assert.NoError(t, e.Wait())
}
