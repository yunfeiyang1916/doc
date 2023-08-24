// Panic和Recover机制
package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"
)

var user = os.Getenv("USER")

func init() {
	fmt.Println("user=", user)
}

func main() {
	go throwPanic(1)
	go throwPanic(2)
	throwPanic(3)
	fmt.Printf("cpu数量：%d 协程数量：%d\n", runtime.NumCPU(), runtime.NumGoroutine())
	time.Sleep(time.Microsecond)
}

func throwPanic(i int) {
	time.Sleep(time.Microsecond)
	defer func(n int) {
		if err := recover(); err != nil {
			fmt.Printf("捕获到序号 %d 的异常 %s\n", n, err)
		}
	}(i)
	fmt.Println("顺序", i)
	if user == "" {
		panic(strconv.Itoa(i) + " user is empty")
	}
}
