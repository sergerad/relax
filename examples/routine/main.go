package main

import (
	"fmt"

	"github.com/sergerad/relax"
)

func main() {
	// Start a routine
	routine := relax.Go(func() error {
		panic(1)
	})

	// Wait for routine
	if err := routine.Wait(); err != nil {
		fmt.Println(err)
	}
}
