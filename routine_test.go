package relax

import (
	"bytes"
	"errors"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	errTestSentinel = errors.New("test sentinel")
)

func TestRoutine_NilError(t *testing.T) {
	r := Go(func() error {
		return nil
	})
	assert.NoError(t, r.Wait())
}

func TestRoutine_Multicall(t *testing.T) {
	var tests = []struct {
		name        string
		f           func() error
		expectedErr error
	}{
		{"nil", func() error { return nil }, nil},
		{"err", func() error { return errTestSentinel }, errTestSentinel},
		{"panic", func() error { panic("test panic") }, PanicError},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := Go(func() error {
				return test.f()
			})
			for i := 0; i < 5; i++ {
				err := r.Wait()
				if i == 0 {
					if !errors.Is(err, test.expectedErr) {
						t.Errorf("expected %v, got %v. iteration %d", test.expectedErr, err, i)
					}
				} else {
					if err != nil {
						t.Errorf("expected nil, got %v. iteration %d", err, i)
					}
				}
			}
		})
	}
}

func TestRoutine_Release(t *testing.T) {
	var tests = []struct {
		name string
		err  error
	}{
		{"nil", nil},
		{"error", errTestSentinel},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Go(func() error {
				return test.err
			}).Release(func(err error) {
				if !errors.Is(err, test.err) {
					t.Errorf("expected %v, got %v", test.err, err)
				}
			})

		})
	}
}

func TestRoutine_Panic_NotNilError(t *testing.T) {
	var tests = []struct {
		name            string
		panicDatum      any
		expectedContent string
	}{
		{"empty", "", ""},
		{"non-empty", "test panic", "test panic"},
		{"error", errTestSentinel, errTestSentinel.Error()},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := Go(func() error {
				panic(test.panicDatum)
			})
			err := r.Wait()
			require.Error(t, err)
			switch x := test.panicDatum.(type) {
			case error:
				assert.Contains(t, err.Error(), x.Error())
			default:
				assert.Contains(t, err.Error(), test.expectedContent)
			}
			assert.True(t, errors.Is(err, PanicError))
		})
	}
}

func TestRoutine_Stack_WritesStack(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	r := &Routine{
		stackBuff: make([]byte, 1<<16),
	}
	stack := r.stack()
	assert.NotEmpty(t, stack)
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
