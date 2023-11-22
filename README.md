# Relax
#### *...there's no need to panic.*

In some situations, encountering a panic can be problematic. For example:
* If your REST server panics while handling a POST request, you may end up with a dangling resource in your database;
* Applications that write state to filesystems may produce irrecoverable state if a series of dependant file writes is interrupted by a panic; and
* When running tests, such as integration tests using Go's [coverage capabilities](https://go.dev/testing/coverage/#panicprof), a panic can cause you to lose your prized test results.

In contrast to panicking applications, relaxed Go programs start and finish gracefully, even in the case of SIGINT, SIGTERM, and concurrent panics.

Relaxed Go programs will only shutdown after all running operations and connections have completed and closed, respectively.

This module is intended to aid in the development of relaxed Go programs. The main challenge in achieving this is to ensure that all panicking goroutines are recovered and lead to the graceful shutdown of the program. Read about panic and recover [here](https://go.dev/blog/defer-panic-and-recover) for more detail.

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
		panic("failed")
	})
```

Finally, in the main goroutine, make sure to wait for the error group:

```Go
	if err := group.Wait(); err != nil {
		log.Fatal(err)
	}
```

When you only have a single goroutine to run, you can use `Routine` instead of `RoutineGroup`:

```Go
	routine := relax.Go(func() error {
		panic("failed")
	})
	if err := routine.Wait(); err != nil {
		log.Fatal(err)
	}
```

For more detailed usage, see [examples](./examples/) and the `*_test.go` files.
