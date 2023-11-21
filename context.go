package relax

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// Context instantiates a context that is cancelled only when
// the SIGINT/SIGTERM signals are received.
// This context should be set up in main() of your application.
func Context() context.Context {
	// Instantiate cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// Cancel context on signals
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, Signals()...)
		// Wait for signal
		<-signalChan
		// Cancel context
		cancel()
	}()
	return ctx
}

// Signals returns the signals that will cause the context to be cancelled.
func Signals() []os.Signal {
	return []os.Signal{syscall.SIGINT, syscall.SIGTERM}
}
