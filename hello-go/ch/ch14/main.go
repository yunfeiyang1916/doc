// 协程与通道
package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"
)

// 长时间等待
func longWait() {
	fmt.Println("Beginning longWait()")
	//休眠5秒，单位为纳秒
	time.Sleep(5 * 1e9)
	fmt.Println("End of longWait()")
}

// 短时间等待
func shortWait() {
	fmt.Println("Beginning shortWait()")
	//休眠2秒，单位为纳秒
	time.Sleep(2 * 1e9)
	fmt.Println("End of shortWait()")
}

// 协程测试
func goroutineTest() {
	var numCores = flag.Int("n", 2, "number of CPU cores to use")
	flag.Parse()
	//设置并发时的线程数量
	runtime.GOMAXPROCS(*numCores)
	fmt.Println("numCores=", *numCores)
	fmt.Println("in goroutineTest()")
	go longWait()
	go shortWait()
	fmt.Println("About to sleep in goroutineTest()")
	//休眠10秒，单位为纳秒
	time.Sleep(10 * 1e9)
	fmt.Println("At the end of goroutineTest()")
}

// 使用通道在协程间通信
// 发送数据
func sendData(ch chan string) {
	fmt.Println("开始发送数据")
	fmt.Println("发送张三")
	ch <- "张三"
	fmt.Println("发送李四")
	ch <- "李四"
	fmt.Println("发送王五")
	ch <- "王五"
	fmt.Println("发送燕小六")
	ch <- "燕小六"
	fmt.Println("发送鬼脚七")
	ch <- "鬼脚七"
	fmt.Println("结束发送数据")
}

// 读取数据
func getData(ch chan string) {
	fmt.Println("开始接收数据")
	//通道是同步阻塞的，也就是发送的消息没有被消费就会阻塞住
	//ch <- "我插一条数据"
	//fmt.Println("会执行吗")
	var input string
	//休眠两秒
	//time.Sleep(2e9)
	for {
		input = <-ch
		fmt.Println("接收到" + input)
	}
	fmt.Println("结束发送数据")
}
func channelTest() {
	var ch chan string
	ch = make(chan string)
	go sendData(ch)
	go getData(ch)
	//休眠1秒
	time.Sleep(1e9)
}

// 通道阻塞，通道是同步发送数据的
// 循环发送
func pump(ch chan int) {
	for i := 0; ; i++ {
		ch <- i
	}
}

// 循环消费
func suck(ch chan int) {
	for {
		fmt.Println(<-ch)
	}
}
func channelBlock() {
	ch := make(chan int)
	go pump(ch)
	//fmt.Println(<-ch)
	go suck(ch)
	time.Sleep(1e9)
	fmt.Println("主程序结束")
}

// 死锁，发送消息阻塞住了
func deadLock() {
	ch := make(chan int)
	//go func() { ch <- 2 }()
	ch <- 2
	go func() {
		fmt.Println(<-ch)
	}()
	time.Sleep(1e9)
}

// 带缓冲的通道
func channelBuffer() {
	//声明一个带缓冲的通道
	//在缓冲满载（缓冲被全部使用）之前，给一个带缓冲的通道发送数据是不会阻塞的，而从通道读取数据也不会阻塞，直到缓冲空了
	ch := make(chan int, 100)
	//因为缓冲没有满载，所以不会阻塞，这样就不会产生死锁了
	ch <- 2
	go func() { fmt.Println(<-ch) }()
	time.Sleep(1e9)
}

// 使用协程回报计算结果
func sumTest() {
	ch := make(chan int)
	go func() int {
		var s int
		for i := 0; i < 100; i++ {
			s += i
		}
		ch <- s
		return s
	}()
	fmt.Println(<-ch)
}

// 使用通道实现信号量
// 通道信号量测试
func semaphoreTest() {
	n := 1000000
	data := make([]int, n)
	res := make([]int, n)
	sem := make(chan int, n)
	for i := 0; i < n; i++ {
		data[i] = i
	}
	for k, v := range data {
		go func(k int, v int) {
			res[k] = k + v
			//发出信号
			sem <- k
		}(k, v)
	}
	for i := 0; i < n; i++ {
		//等待执行结果
		fmt.Println(<-sem, "  ", res[i])

	}
}

// 实现一个信号量
type Empty interface{}

// 信号量
type Semaphore chan Empty

// 请求n个资源
func (s Semaphore) P(n int) {
	e := new(Empty)
	for i := 0; i < n; i++ {
		//发送一条消息来阻塞
		s <- e
	}
}

// 释放n个资源
func (s Semaphore) V(n int) {
	for i := 0; i < n; i++ {
		//消费一条消息
		<-s
	}
}

// 实现互斥锁
// 加锁
func (s Semaphore) Lock() {
	s.P(1)
}

// 解锁
func (s Semaphore) Unlock() {
	s.V(1)
}

// 实现信号等待
// 等待n个信号
func (s Semaphore) Wait(n int) {
	s.P(n)
}

// 释放一个信号
func (s Semaphore) Signal() {
	s.V(1)
}
func semaphoreTest2() {
	sem := make(Semaphore, 100)
	go func() {
		sem.Wait(100)
		fmt.Println("只要有一个信号被释放了，我就可以执行了")
	}()
	go func() {
		for i := 0; i < 100; i++ {
			sem.Signal()
			fmt.Printf("释放第%d个信号！\n", i)
		}
	}()
	time.Sleep(1e9)
}

// 使用循环读取通道
func forRange() {
	ch := make(chan int)
	go func() {
		for i := 0; i < 100; i++ {
			fmt.Println("发送：", i)
			ch <- i
		}
	}()
	go func() {
		for v := range ch {
			fmt.Printf("The value is %v\n", v)
		}
	}()
	time.Sleep(1e9)
}

// 通道的方向
// 顺序发送2、3、4...到ch
func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i
	}
}
func filter(in, out chan int, prime int) {
	for {
		i := <-in
		if i%prime != 0 {
			out <- 1
		}
	}
}

// 打印素数
func sieveTest() {
	ch := make(chan int)
	go generate(ch)
	for {
		prime := <-ch
		fmt.Print(prime, " ")
		ch1 := make(chan int)
		go filter(ch, ch1, prime)
		ch = ch1
	}
}

// 可关闭的通道
func channelClose() {
	ch := make(chan string)
	go func() {
		ch <- "张三"
		ch <- "李四"
		ch <- "王五"
		ch <- "燕小六"
		ch <- "鬼脚七"
		close(ch)
	}()
	getDataForClose(ch)
}

// 只接收消息的通道
func getDataForClose(ch <-chan string) {
	// for {
	// 	str, ok := <-ch
	// 	//通道是否关闭
	// 	if !ok {
	// 		break
	// 	}
	// 	fmt.Println(str)
	// }

	//使用for-range读取通道会自动检测通道是否关闭
	for str := range ch {
		fmt.Println(str)
	}
}

// 从不同并发执行的协程中取值
func selectTest() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	//定时器，每一毫秒发送一次消息
	ticker := time.NewTicker(1e6)
	defer ticker.Stop()
	go func(ch chan int) {
		for i := 0; ; i++ {
			ch <- i
		}
	}(ch1)
	go func(ch chan int) {
		for i := 0; ; i++ {
			ch <- i + 5
		}
	}(ch2)
	//select选择通道
	go func() {
		for {
			select {
			case u := <-ch1:
				fmt.Printf("ch1接收到消息：%v\n", u)
				break
			case v := <-ch2:
				fmt.Printf("ch2接收到消息%v\n", v)
			case t := <-ticker.C:
				fmt.Printf("定时器发送消息%v\n", t)
			}
		}
	}()
	time.Sleep(1e9)
}

// 定时器测试
func tickerTest() {
	//每100毫秒发送一次消息
	tick := time.Tick(1e8)
	//500毫秒后只发送一次消息
	boom := time.After(5e8)
	for i := 0; i < 10; i++ {
		select {
		case u := <-tick:
			fmt.Printf("tick:%v\n", u)
		case v := <-boom:
			fmt.Printf("boom:%v\n", v)
		}
	}
}

// 超时测试
func timeoutTest() {
	//简单超时模式
	//设置缓冲为1时发送第一条消息不会阻塞
	ch := make(chan string, 1)
	timeout := make(chan bool, 1)
	go func() {
		//休眠一秒
		time.Sleep(1e9)
		timeout <- true
	}()
	//ch <- "不会阻塞哦，会正常发送的"
	select {
	case <-ch:
		fmt.Printf("接收到ch消息\n")
	case <-timeout:
		fmt.Println("接收到超时消息")
		break
	}
	//取消耗时很长的同步调用
	ch2 := make(chan int, 1)
	go func() {
		time.Sleep(2e9)
		ch2 <- 1
	}()
	select {
	case <-ch2:
		fmt.Println("接收到ch2消息")
	case <-time.After(1e9):
		fmt.Println("ch2执行超时")
		break
	}
}

// 惰性生成器
// 返回下一个整数值
func integers() chan int {
	yield := make(chan int)
	count := 0
	go func() {
		for {
			yield <- count
			count++
		}
	}()
	return yield
}

// 生成
func generateInteger() {
	resume := integers()
	fmt.Println(<-resume)
	fmt.Println(<-resume)
	fmt.Println(<-resume)
	fmt.Println(<-resume)
	fmt.Println(<-resume)
}

// 惰性生成器函数工厂
func BuildLazyEvaluator(evalFunc EvalFunc, initState Empty) func() Empty {
	retValChan := make(chan Empty)
	loopFunc := func() {
		var actState Empty = initState
		var retVal Empty
		for {
			retVal, actState = evalFunc(actState)
			retValChan <- retVal
		}
	}
	retFunc := func() Empty {
		return <-retValChan
	}
	go loopFunc()
	return retFunc
}

// 生成整型函数的惰性生成器
func BuildLazyIntEvaluator(evalFunc EvalFunc, initState Empty) func() int {
	ef := BuildLazyEvaluator(evalFunc, initState)
	return func() int {
		return ef().(int)
	}
}

// 惰性生成函数
type EvalFunc func(Empty) (Empty, Empty)

func lazyTest() {
	evalFunc := func(state Empty) (Empty, Empty) {
		os := state.(int)
		ns := os + 2
		return os, ns
	}
	even := BuildLazyIntEvaluator(evalFunc, 0)
	for i := 0; i < 10; i++ {
		fmt.Printf("%vth even:%v\n", i, even())
	}
}

func main() {
	//goroutineTest()
	//channelTest()
	//channelBlock()
	//deadLock()
	//channelBuffer()
	//sumTest()
	//semaphoreTest()
	//semaphoreTest2()
	//forRange()
	//sieveTest()
	//channelClose()
	//selectTest()
	//tickerTest()
	//timeoutTest()
	//generateInteger()
	lazyTest()
}
