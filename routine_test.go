package relax

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoutine_NilError(t *testing.T) {
	r := Go(func() error {
		return nil
	})
	assert.NoError(t, r.Wait())
}

func TestRoutine_Panic_NotNilError(t *testing.T) {
	panicMsg := "test panic"
	r := Go(func() error {
		panic("test panic")
	})
	err := r.Wait()
	require.Error(t, err)
	assert.Contains(t, err.Error(), panicMsg)
	assert.True(t, errors.Is(err, PanicError))
}
