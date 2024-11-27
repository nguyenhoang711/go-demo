package goroutineleak

import "fmt"


func NotCorrectCode() {
    ch := make(chan int)
    go func() {
        val := <-ch  // Goroutine blocked here indefinitely
        fmt.Println(val)
    }()
}

func CorrectCode(ch <-chan int) {
	go func() {
		val := <-ch // goroutine receive value from channel --> exit
		fmt.Println(val)
	}()
}