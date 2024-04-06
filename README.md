# Relax
#### *...there's no need to panic.*

![Coverage](https://img.shields.io/badge/Coverage-100.0%25-brightgreen)

Sometimes we don't want a single panic to result in the abrupt termination of our entire application.

For example, an unrecovered panic might be problematic:
* If it causes your API server to abruptly terminate many parallel connections or leave dangling resources in your data store;
* If your application writes state to a filesystem and an interruption between writes may produce irrecoverable state;
* If you CLI application is expected to always output a particular format (E.G. JSON); or
* When running tests, such as integration tests using Go's [coverage capabilities](https://go.dev/testing/coverage/#panicprof), where you may lose your prized test results.

Instead of crashing, relaxed Go programs always shutdown gracefully, even in the case of SIGINT, SIGTERM, and concurrent panics.

#### What can this module help with?

Relaxed Go programs will only shutdown after all running operations and connections have completed and closed, respectively.

This module is intended to aid in the development of relaxed Go programs.

The main challenge in achieving this is to ensure that panicking goroutines are recovered and lead to the graceful shutdown of the program. A panic can only be recovered inside the goroutine within which the panic occurred. Read about panic and recover [here](https://go.dev/blog/defer-panic-and-recover) for more detail.

The other challenge to recovering from panics is that it requires a bit of complex boilerplate code. In order to convert a panic to an error, you need to defer a `recover()` call inside a closure. This closure needs to assign the panic content to an error variable which is a named return variable in order for it to be returned by the func that has recovered from the panic. See [here](https://golang.org/ref/spec#Defer_statements) for more detail on that.

With the `relax` module, recovering from a panicking goroutine looks just like a normal func call.

## Usage

Import the pkg:

```Go
import (
	"github.com/sergerad/relax"
)
```

Instantiate the main context in the main goroutine.

```Go
func main() {
	mainCtx := relax.Context()
```

This ensures that SIGINT/SIGTERM will cause all contexts used in the application to be `Done()`.

If you have multiple, long running processes to run in your program, you can use `RoutineGroup` to launch them.

```Go
	group, ctx := relax.NewGroup(mainCtx)
```

You can use the `RoutineGroup` to launch goroutines which will return an error if they encounter a panic.
```Go
	group.Go(func() error {
		[]int{}[0] = 1 // Panic for example
	})
```

Finally, in the main goroutine, make sure to wait for the error group:

```Go
	if err := group.Wait(); err != nil {
		// Handle the error...
	}
```

When you only have a single goroutine to run, you can use `Routine` instead of `RoutineGroup`:

```Go
	routine := relax.Go(func() error {
		[]int{}[0] = 1 // Panic for example
	})
	if err := routine.Wait(); err != nil {
		// Handle the error...
	}
```

If you don't wish to wait for the result, you can register a callback:

```Go
	routine.Release(func(err error) {
		// Handle the error
	})
```

For more detailed usage, see [examples](./examples/) and the `*_test.go` files.
