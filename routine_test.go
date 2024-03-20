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

func TestRecoverError(t *testing.T) {
	var tests = []struct {
		name        string
		panicDatum  any
		expectedErr error
	}{
		{"nil", nil, nil},
		{"string", "fail", PanicError},
		{"int", 0, PanicError},
		{"error", errors.New(""), PanicError},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := recoverError(test.panicDatum)
			if !errors.Is(err, test.expectedErr) {
				t.Errorf("expected %v, got %v", test.expectedErr, err)
			}
		})
	}
}
