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
