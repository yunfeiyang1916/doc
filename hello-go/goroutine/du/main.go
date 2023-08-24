// 并发的目录遍历，类似unix的du命令，使用协程同步
package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func main() {
	du1()
}

var wg sync.WaitGroup

// 限制并发数量的通道信号量
var sema = make(chan struct{}, 20)

// 用于发出退出消息的通道
var done = make(chan struct{})

// 是否已取消
func canceled() bool {
	//使用done取消消息
	/*select {
	case <-done:
		return true
	default:
		return false
	}*/
	//根据上下文判断
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

// 使用上下文控制协程取消
var (
	ctx    context.Context
	cancel context.CancelFunc
)

func du1() {
	ctx, cancel = context.WithCancel(context.Background())
	flag.Parse()
	roots := flag.Args()
	//如果未输入目录，则输出当前目录
	if len(roots) == 0 {
		roots = []string{"e:/"}
	}
	//总文件大小消息通道
	fileSizes := make(chan int64)
	for _, dir := range roots {
		wg.Add(1)
		go walkDir(dir, fileSizes)
	}
	go func() {
		//读任意一个字符
		os.Stdin.Read(make([]byte, 1))
		//close(done)
		//发出取消消息
		cancel()
	}()

	go func() {
		//文件计算完后需要关闭通道
		wg.Wait()
		close(fileSizes)
	}()
	//每500毫秒打印一次
	tick := time.Tick(500 * time.Millisecond)
	var nfiles, nbytes int64
loop:
	//需要从总文件大小通道循环读取，否则会造成死锁
	for {
		select {
		/*case d, ok := <-done: //是否取消
		fmt.Println(d, ok)
		fmt.Println("收到取消信号")
		//休眠两秒后结束
		time.Sleep(5 * time.Second)
		return*/
		case d, ok := <-ctx.Done(): //是否取消
			fmt.Println(d, ok)
			fmt.Println("收到取消信号")
			//休眠两秒后结束
			time.Sleep(5 * time.Second)
			return
		case size, ok := <-fileSizes:
			//通道是否关闭
			if !ok {
				break loop ///退出标签，如果只使用break只会退出select
			}
			nfiles++
			nbytes += size
			//printDiskUsage(nfiles, nbytes)
		case <-tick: //每500毫秒打印一次
			printDiskUsage(nfiles, nbytes)
		}
	}
	printDiskUsage(nfiles, nbytes)
}

// 递归遍历目录，统计所有文件大小
func walkDir(dir string, fileSizes chan<- int64) {
	//获取信号量，是否可以执行
	sema <- struct{}{}
	defer func() {
		wg.Done()
		//释放占用的信号量
		<-sema
	}()
	//是否已取消，如果是则结束
	if canceled() {
		return
	}
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("读取%s出错:%s\n", dir, err)
		return
	}
	for _, file := range fs {
		if file.IsDir() {
			subdir := filepath.Join(dir, file.Name())
			wg.Add(1)
			//递归调用
			go walkDir(subdir, fileSizes)
		} else {
			fileSizes <- file.Size()
		}
	}
}

// 打印磁盘使用情况
// @param	nfiles	文件总数
// @params	nbytes	文件总大小
func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files %1.f GB\n", nfiles, float64(nbytes)/1e9)
}
