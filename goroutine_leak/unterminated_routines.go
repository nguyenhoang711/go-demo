package goroutineleak

import (
	"context"
	"fmt"
	"sync"
)

func NotCorrectUnterminated() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel() // not exit goroutine properly before program end

    go func() {
        select {
        case <-ctx.Done():  // Goroutine will terminate only if ctx.Done() is triggered
            fmt.Println("Goroutine exiting")
        }
    }()
}

func CorrectTerminatedGoroutine() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		select {
		case <- ctx.Done():
			fmt.Println("Goroutine exiting")
		}
	}()
	//simulate work
	fmt.Println("Doing some work...")
	cancel()

	// wait goroutine finish
	wg.Wait()
}