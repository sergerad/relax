# Relax

*Relax* is the antonym for the verb *panic*.

In the context of Go programs, relax means to make critical failures less severe so that graceful shutdown is never jeopardized.

Relaxed Go programs start and finish gracefully, even in the case of SIGINT, SIGTERM, and concurrent panics.

This module is intended to aid in the development of relaxed Go programs.


## Use Case

This module is only intended for use in a program which:
1. launches concurrent, long-running goroutines; and
2. intends to gracefully shutdown the entire program as soon one of the goroutines either:
    * returns an error; or
    * panics.

## Usage

Import the pkg:

```Go
import (
	"github.com/sergerad/relax"
)
```

Instantiate the error group and root context in the main goroutine.

```Go
func main() {
	g, ctx := relax.Main()
```

This ensures that the main goroutine is relaxed against SIGINT and SIGTERM. It also allows us to relax child goroutines through the root context.

Run goroutines through the error group. Make sure to relax the goroutines so that panics do not interfere with graceful shutdown:

```Go
	g.Go(func() error {
		defer relax.Routine(ctx)
		return myLongRunningFunc(ctx)
```

Note that every goroutine in the program (apart from main goroutine), must call `relax.Routine(ctx)` in order to prevent panics from interfering with graceful shutdown. This is because panics cause programs to crash when all funcs in the respective goroutine have returned. Read about panic and recover [here](https://go.dev/blog/defer-panic-and-recover) for more detail.

Finally, in the main goroutine, make sure to wait for the error group:

```Go
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
```

For more detailed usage, see the tests.
