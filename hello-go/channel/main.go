package main

import (
	"fmt"
	"time"
)

func main() {
	//defaultTest()
	multCaseTest()
}

// 多个条件同时满足
func multCaseTest() {
	ch := make(chan int)
	go func() {
		// 每隔一秒写一次
		for range time.Tick(time.Second) {
			ch <- 1
		}
	}()

	for {
		select {
		case <-ch:
			fmt.Println("case1")
		case <-ch:
			fmt.Println("case2")
		}
	}
}

func defaultTest() {
	ch := make(chan int)
	select {
	case <-ch:
		fmt.Println("case1")
	case ch <- 1:
		fmt.Println("case2")
	default:
		fmt.Println("default")
	}
}
