package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/sergerad/relax"
)

func main() {
	// Main context that is cancelled on SIGINT/SIGTERM
	ctx := relax.Context()

	// Indicate how to exit the app gracefully
	println("Either of the following commands would exit the app gracefully:")
	println("kill -SIGINT", os.Getpid())
	println("kill -SIGTERM", os.Getpid())

	// Exit the app automatically. Comment this out to try manually
	// sending signals
	go func() {
		cmd := exec.Command("kill", "-SIGINT", fmt.Sprint(os.Getpid()))
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for context to be cancelled
	<-ctx.Done()
	println("Exiting gracefully")
}
