package relax

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRelax_WithPanic_MainRoutineCompletes(t *testing.T) {
	g, ctx := Main()

	g.Go(func() error {
		defer Routine(ctx)
		return func() error {
			select {
			case <-ctx.Done():
				t.Logf("blocking routine is complete")
			}
			return nil
		}()
	})

	g.Go(func() error {
		childCtx, _ := context.WithCancel(ctx)
		defer Routine(childCtx)
		panic(1)
	})

	if err := g.Wait(); err != nil {
		assert.FailNowf(t, "error should be nil: %s", err.Error())
	}

	t.Log("main goroutine is complete")
}

func TestRelax_WithError_MainRoutineCompletes(t *testing.T) {
	g, ctx := Main()

	g.Go(func() error {
		defer Routine(ctx)
		return func() error {
			select {
			case <-ctx.Done():
				t.Logf("blocking routine is complete")
			}
			return nil
		}()
	})

	g.Go(func() error {
		childCtx, _ := context.WithCancel(ctx)
		defer Routine(childCtx)
		return errors.New("bad")
	})

	if err := g.Wait(); err == nil {
		assert.FailNow(t, "error should not be nil")
	}

	t.Log("main goroutine is complete")
}
