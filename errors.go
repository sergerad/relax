package relax

import (
	"fmt"
)

var (
	// PanicError is returned when a panic is recovered
	// during the execution of a goroutine
	PanicError = fmt.Errorf("recovered from panic")
)
