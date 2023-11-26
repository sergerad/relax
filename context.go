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
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		// Wait for signal
		<-signalChan
		// Cancel context
		cancel()
	}()
	return ctx
}
