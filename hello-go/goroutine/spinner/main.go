// 打印菲波那契数列的第45个元素值，使用动画等待
// 主函数结束时，所有的协程都会被直接打断，程序退出。
package main

import (
	"fmt"
	"time"
)

func main() {
	//每100毫秒打印一次动画
	go spinner(100 * time.Millisecond)
	const n = 45
	fibN := fib(n)
	fmt.Printf("\nFibonacci(%d)=%d\n", n, fibN)
}

// 打印动画
func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

// 菲波那契数列函数
func fib(x int) int {
	if x < 2 {
		return x
	}
	return fib(x-1) + fib(x-2)
}
