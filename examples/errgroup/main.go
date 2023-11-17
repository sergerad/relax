package main

import (
	"context"
	"fmt"
	"time"

	"github.com/sergerad/relax"
	"golang.org/x/sync/errgroup"
)

func main() {
	// Instantiate the main context
	mainCtx, cancel := relax.MainContext()

	// Use the main context for errgroup
	g, ctx := errgroup.WithContext(mainCtx)

	// Launch goroutine that blocks on context
	g.Go(func() error {
		// Relax panics
		defer relax.Recover(cancel)
		// Block
		select {
		case <-ctx.Done():
			fmt.Println("blocking routine done")
		}
		return nil
	})
	// Launch goroutine that resembles a long running processor
	g.Go(func() error {
		return exampleProcessor(ctx)
	})

	// Wait for errgroup
	if err := g.Wait(); err != nil {
		fmt.Println("error from group", err)
	}
	fmt.Println("shutting down")

	// Sleep to give all goroutines a chance to exit due to ctx.Done()
	time.Sleep(1)
}

func exampleProcessor(ctx context.Context) error {
	childCtx, cancel := context.WithCancel(ctx)
	go func() {
		defer relax.Recover(cancel)
		panic(1)
	}()
	go func() {
		select {
		case <-childCtx.Done():
		}
		fmt.Println("processor goroutine done")
	}()
	select {
	case <-childCtx.Done():
		fmt.Println("processor done")
		return childCtx.Err()
	}
}
