package concurrency

import (
	"fmt"
	"time"
)

func CheckTimeoutWithSelect() {
	c1 := make(chan string)
	c2 := make(chan string)
	go func() {
		time.Sleep(1 * time.Second)
		c1 <- "from 1"
	}()

	go func() {
		time.Sleep(3 * time.Second)
		c2 <- "from 2"
	}()

	select {
	case msg1 := <-c1:
		fmt.Println("Message 1", msg1)
	case msg2 := <-c2:
		fmt.Println("Message 2", msg2)
	case <-time.After(2 * time.Second):
		fmt.Println("timeout")
	}
	var input string
	fmt.Scanln(&input)
}
