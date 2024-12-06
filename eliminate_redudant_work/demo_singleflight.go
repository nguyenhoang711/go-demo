package eliminateredudantwork

import (
	"fmt"
	"time"

	"golang.org/x/sync/singleflight"
)

var group singleflight.Group

func expensiveOperation(key string) (interface{}, error) {
	time.Sleep(2 * time.Second)
	return fmt.Sprintf("Data for %s", key), nil
}


func DemoSimpleCaseWithRedudantWork() {
	start := time.Now()
	for i := 0 ;i < 5; i++ {
		go func(i int) {
			val, err, _ := group.Do("my_key", func()(interface{}, error) {
				return expensiveOperation("my_key")
			})
			if err == nil {
				fmt.Printf("Goroutine %d got result %v\n", i, val)
			}
		}(i)
	}
	time.Sleep(3 * time.Second)
	fmt.Printf("time since start is ::%v\n", time.Since(start))
}