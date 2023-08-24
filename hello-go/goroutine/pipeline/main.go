package main

import "fmt"

func main() {
	//pipeline1()
	//pipeline2()
	pipeline3()
}

// 使用for循环串联通道
func pipeline1() {
	//自然数通道
	naturals := make(chan int)
	//平方数通道
	squares := make(chan int)
	//计数器，用于生成0 1 2 ...形式的整数序列
	go func() {
		for x := 0; x < 100; x++ {
			//发送自然数
			naturals <- x
		}
		//关闭自然数通道，关闭后在发送数据将导致panic异常，
		// 当一个被关闭的channel中已经发送的数据都被成功接收后，后续的接收操作将不再阻塞，它们会立即返回一个零值
		close(naturals)
	}()
	//求平方
	go func() {
		for {
			x, ok := <-naturals
			if !ok {
				break
			}
			//发送平方值
			squares <- x * x
		}
		close(squares)
	}()
	//打印计算的平方值
	for {
		x, ok := <-squares
		if !ok {
			break
		}
		fmt.Printf("%d ", x)
	}
}

// 使用for-range循环串联通道
// 使用range循环是上面处理模式的简洁语法，它依次从channel接收数据，当channel被关闭并且没有值可接收时跳出循环
func pipeline2() {
	//自然数通道
	naturals := make(chan int)
	//平方数通道
	squares := make(chan int)
	//计数器，用于生成0 1 2 ...形式的整数序列
	go func() {
		for x := 0; x < 100; x++ {
			//发送自然数
			naturals <- x
		}
		//关闭自然数通道，关闭后在发送数据将导致panic异常，
		// 当一个被关闭的channel中已经发送的数据都被成功接收后，后续的接收操作将不再阻塞，它们会立即返回一个零值
		close(naturals)
	}()
	//求平方
	go func() {
		for x := range naturals {
			//发送平方值
			squares <- x * x
		}
		close(squares)
	}()
	//打印计算的平方值
	for x := range squares {
		fmt.Printf("%d ", x)
	}
}

func pipeline3() {
	//自然数通道
	naturals := make(chan int)
	//平方数通道
	squarers := make(chan int)
	go counter(naturals)
	go squarer(naturals, squarers)
	printer(squarers)
}

// 计数器，用于生成0 1 2 ...形式的整数序列
// @param out 只能向通道发送消息
func counter(out chan<- int) {
	for x := 0; x < 100; x++ {
		out <- x
	}
	close(out)
}

// 求平方
// @param	in	只能从通道接收自然数
// @param	out	只能向通道发送
func squarer(in <-chan int, out chan<- int) {
	for x := range in {
		out <- x * x
	}
	close(out)
}

// 输出接收到的平方数
// @param	in	只能从通道接收
func printer(in <-chan int) {
	for x := range in {
		fmt.Printf("%d ", x)
	}
	fmt.Println()
}
