package goroutineleak

import (
	"fmt"
	"sync"
)

func NotCorrectInfiniteCreation() {
	for i := 0; i < 10000; i++ {
        go func(i int) {
            // Do some work without proper termination or exit
            fmt.Println(i)
        }(i)
    }
}

func CorrectCreation() {
	var wg sync.WaitGroup
	for i := 0; i< 10000; i++ {
		wg.Add(1) // increment the WaitGroupC counter
		go func(i int) {
			defer wg.Done()
			fmt.Println(i)
		}(i)
	}
	wg.Wait() // wait all goroutines to finish
}