package relax

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"golang.org/x/sync/errgroup"
)

const (
	rootCancelKey = "rootCancel"
)

// Main should be called at the start of the main goroutine.
// It sets up the root context and cancels it on SIGINT/SIGTERM.
// The returned error group should be used to launch goroutines that
// are intended to shutdown as soon as any error is returned from any
// goroutines launched via the same error group.
func Main() (*errgroup.Group, context.Context) {
	// Set up root cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// Cancel root context on SIGINT/SIGTERM
		quitCh := make(chan os.Signal, 1)
		signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM)
		<-quitCh
		cancel()
	}()

	// Return the error group and associated context
	return errgroup.WithContext(context.WithValue(ctx, rootCancelKey, cancel))
}

// Routine should be deferred at the start of any goroutine.
// It will recover from panic, cancel the context and exit the
// goroutine. Assuming all the program's goroutines are using
// the cancelled context, or a contexts derived from it, the entire
// program will shutdown gracefully.
func Routine(ctx context.Context) {
	// Find root cancel func if it exists in the context
	if cancel, ok := ctx.Value(rootCancelKey).(context.CancelFunc); ok {
		if r := recover(); r != nil {
			fmt.Println("recovered from panic:", r)
			cancel()
			// Ensure current goroutine ends
			runtime.Goexit()
		}
	}
}
