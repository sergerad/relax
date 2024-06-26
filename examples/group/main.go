package main

import (
	"context"
	"fmt"

	"github.com/sergerad/relax"
)

func main() {
	// Instantiate the main context and error group
	group, ctx := relax.NewGroup(relax.Context())

	// Launch goroutine that blocks on context
	group.Go(func() error {
		<-ctx.Done()
		fmt.Println("blocking routine done")
		return nil
	})

	// Launch goroutine that resembles a long running processor
	group.Go(func() error {
		return exampleProcessor(ctx)
	})

	// Wait for routine group
	if err := group.Wait(); err != nil {
		fmt.Println(err)
	}
}

func exampleProcessor(ctx context.Context) error {
	[]int{}[0] = 1
	return nil
}
