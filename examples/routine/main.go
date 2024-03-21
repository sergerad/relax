package main

import (
	"fmt"

	"github.com/sergerad/relax"
)

func main() {
	// Start and wait for a routine
	routine := relax.Go(func() error {
		panic(1)
	})
	if err := routine.Wait(); err != nil {
		fmt.Println(err)
	}

	// If we don't want to wait for the routine, release it
	routine = relax.Go(func() error {
		panic(2)
	})
	routine.Release(func(err error) {
		fmt.Println(err)
	})
}
