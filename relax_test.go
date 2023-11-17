package relax

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestErrGroup_WithPanic_MainRoutineCompletes(t *testing.T) {
	mainCtx, cancel := MainContext()
	g, ctx := errgroup.WithContext(mainCtx)
	g.Go(func() error {
		defer Recover(cancel)
		<-ctx.Done()
		t.Logf("blocking routine is complete")
		return nil
	})
	g.Go(func() error {
		defer Recover(cancel)
		panic(1)
	})
	if err := g.Wait(); err != nil {
		assert.FailNowf(t, "error should be nil: %s", err.Error())
	}
	<-ctx.Done()
	t.Log("main goroutine is complete")
}

func TestErrGroup_WithError_MainRoutineCompletes(t *testing.T) {
	mainCtx, cancel := MainContext()
	g, ctx := errgroup.WithContext(mainCtx)

	g.Go(func() error {
		defer Recover(cancel)
		<-ctx.Done()
		t.Logf("blocking routine is complete")
		return nil
	})
	g.Go(func() error {
		defer Recover(cancel)
		return fmt.Errorf("bad")
	})
	if err := g.Wait(); err == nil {
		assert.FailNow(t, "error should not be nil")
	}
	<-ctx.Done()
	t.Log("main goroutine is complete")
}
