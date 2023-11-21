package main

import (
	"context"
	"fmt"

	"github.com/sergerad/relax"
)

func main() {
	// Instantiate the main context and error group
	g, ctx := relax.NewErrorGroup(relax.Context())

	// Launch goroutine that blocks on context
	g.Go(func() error {
		<-ctx.Done()
		fmt.Println("blocking routine done")
		return nil
	})

	// Launch goroutine that resembles a long running processor
	g.Go(func() error {
		return exampleProcessor(ctx)
	})

	// Wait for errgroup
	if err := g.Wait(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("shutting down")
}

func exampleProcessor(ctx context.Context) error {
	panic("processor failed")
}
