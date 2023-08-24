package main

import (
	"fmt"
	"sync"
)

func main() {
	//concurrentMapTest1()
	concurrentMapTest2()
}

// 并发写map会抛出：fatal error: concurrent map writes
func concurrentMapTest1() {
	m := make(map[int]int)
	for i := 0; i < 10; i++ {
		go func() {
			m[i] = i
		}()
	}
}

// 一个协程写map,多个读协程
func concurrentMapTest2() {
	m := make(map[int]int)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("进入协程")
		for i := 0; i < 10; i++ {
			m[i] = i
		}
	}()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Println(m[i])
		}(i)
	}
	wg.Wait()
}
