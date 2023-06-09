# Relax

*Relax* - verb antonym for *panic*.

In the context of Go programs, relax means to make critical failures less severe so that graceful shutdown is never jeopardized.

Relaxed Go programs start and finish gracefully, even in the case of SIGINT, SIGTERM, and concurrent panics.

Relaxed Go programs will only shutdown after all running operations and connections have completed and closed, respectively.

This module is intended to aid in the development of relaxed Go programs. The main challenge in achieving this is to ensure that all panicking goroutines are recovered and lead to the graceful shutdown of the program. Read about panic and recover [here](https://go.dev/blog/defer-panic-and-recover) for more detail.

## Usage

Import the pkg:

```Go
import (
	"github.com/sergerad/relax"
)
```

Instantiate the error group and main context in the main goroutine.

```Go
func main() {
	mainCtx, cancel := relax.MainContext()
```

This ensures that the main goroutine is relaxed against SIGINT and SIGTERM.

If you have multiple, long running processes to run in your program, you can use errgroup to launch them.

```Go
	g, ctx := errgroup.WithContext(mainCtx)
```

Launch your goroutines and make sure to defer `relax.Recover(cancel)` so that any panics do not get in the way of graceful shutdown of the program.
```Go
	g.Go(func() error {
		defer relax.Recover(cancel)
		return myLongRunningProcess(ctx)
	})
```

You can use `relax.Recover(cancel)` for any sub-goroutine that you launch in your program at all. If you launch goroutines without deferring

Finally, in the main goroutine, make sure to wait for the error group:

```Go
	if err := g.Wait(); err != nil {
		if err != nil {
			fmt.Println("error from group", err)
		}
	}
```

For more detailed usage, see [examples](./examples/).
