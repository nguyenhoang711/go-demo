package goroutineleak

import (
	"context"
	"fmt"
	"sync"
)

func worker(ctx context.Context, wg *sync.WaitGroup, doSomeWork <-chan int) {
    defer wg.Done()

    for {
        select {
        case <-ctx.Done():
            fmt.Println("Worker stopping")
            return
        case work := <-doSomeWork:
            fmt.Println(work)
        }
    }
}

