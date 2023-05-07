package relax

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

// MainContext should be called at the start of the main goroutine.
// It sets up the main context and cancels it on SIGINT/SIGTERM.
func MainContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// Cancel on SIGINT/SIGTERM
		kill := make(chan os.Signal, 1)
		signal.Notify(kill, syscall.SIGINT, syscall.SIGTERM)
		<-kill
		cancel()
	}()
	return ctx, cancel
}

// Recover must be deferred in any goroutine to ensure that a panic
// does not prevent graceful shutdown of the program.
// The entire goroutine will end despite recovering from panic.
func Recover(cancel context.CancelFunc) {
	if r := recover(); r != nil {
		fmt.Println("recovered from panic:", r)
		cancel()
		// Ensure current goroutine ends
		runtime.Goexit()
	}
}
