# golim

Go library for limiting the number of goroutines that may simultaneously be in
one or more a critical section of codes.

It includes two structures that provide this capability: one blocking and the
other non-blocking.

## Usage

### Blocking

A Blocking limiter will block until the number of goroutines running in that
limiter is below the threshold, then allows the invoking goroutine to run the
callback.

```Go
package main

import (
	"fmt"
	"sync"

	"github.com/karrick/golim"
)

func main() {
	lim, err := golim.NewBlocking(1024)
	if err != nil {
		panic(err) // TODO
	}

	const total = 1024 * 1024

	var wg sync.WaitGroup
	wg.Add(total)

	for i := 0; i < total; i++ {
		go func(i int) {
			defer wg.Done() // signal to main goroutine that this is complete

			lim.Do(func() {
				fmt.Println(i)
			})
		}(i)
	}

	wg.Wait() // Wait until all goroutines complete
}
```

### NonBlocking

A NonBlocking limiter will never block. It will either invoke the under limit
callback or the over limit callback, depending on whether the number of
goroutines in the limiter.

```Go
package main

import (
	"fmt"
	"sync"

	"github.com/karrick/golim"
)

func main() {
	lim, err := golim.NewNonBlocking(1024)
	if err != nil {
		panic(err) // TODO
	}

	const total = 1024 * 1024

	var wg sync.WaitGroup
	wg.Add(total)

	for i := 0; i < total; i++ {
		go func(i int) {
			defer wg.Done() // signal to main goroutine that this is complete

			underLimit := func() {
				fmt.Println("under", i)
			}

			overLimit := func() {
				fmt.Println("over", i)
			}

			lim.Do(underLimit, overLimit)
		}(i)
	}

	wg.Wait() // Wait until all goroutines complete
}
```
