// 基于select的多路复用
// time.Tick与time.After的区别，Tick是每隔一段时间都向管道发送消息，即使是函数返回后（容易造成协程泄露，适用于用完后程序结束运行）。
// After只发送一次
package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	//countdown()
	//countdown2()
	//countdown3()
	timeAfter()
}

// 倒计时
func countdown() {
	fmt.Println("开始倒计时...")
	//time.Tick会周期性的向管道发送消息
	tick := time.Tick(1 * time.Second)
	for i := 10; i > 0; i-- {
		fmt.Println(i)
		//从通道读取，会阻塞一秒
		<-tick
	}
	fmt.Println("倒计时结束")
}

// 支持中断的倒计时
func countdown2() {
	fmt.Println("开始倒计时，按任意键可中断...")
	//放弃通道
	abort := make(chan struct{})
	go func() {
		//读一个字节
		os.Stdin.Read(make([]byte, 1))
		abort <- struct{}{}
	}()
	select {
	case <-time.After(10 * time.Second):
		fmt.Println("倒计时结束")
	case <-abort:
		fmt.Println("倒计时中断")
	}
}

// 支持中断并输出倒计时
func countdown3() {
	fmt.Println("开始倒计时，按任意键可中断...")
	tick := time.Tick(1 * time.Second)
	//放弃通道
	abort := make(chan struct{})
	go func() {
		//从标准输入读取任意键
		os.Stdin.Read(make([]byte, 1))
		abort <- struct{}{}
	}()
	for i := 10; i > 0; i-- {
		select {
		case <-tick:
			fmt.Println(i)
		case <-abort:
			fmt.Println("倒计时中断")
			return
		}
	}
	fmt.Println("倒计时结束")
}

// 试验time.Tick在接收函数返回后，是否仍然向管道发送消息，time.After只发送一次
func timeAfter() {
	//c := time.Tick(1 * time.Second)
	c := time.After(1 * time.Second)
	for i := 10; i > 0; i-- {
		fmt.Println(i, <-c)
	}
}
