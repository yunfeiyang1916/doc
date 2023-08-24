// 并发的目录遍历，类似unix的du命令
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func main() {
	//du1()
	//du2()
	du3()
}

// 单协程工作
func du1() {
	flag.Parse()
	roots := flag.Args()
	//如果未输入目录，则输出当前目录
	if len(roots) == 0 {
		roots = []string{"."}
	}
	fileSizes := make(chan int64)
	go func() {
		for _, root := range roots {
			walkDir(root, fileSizes)
		}
		close(fileSizes)
	}()
	var nfiles, nbytes int64
	for size := range fileSizes {
		nfiles++
		nbytes += size
	}
	printDiskUsage(nfiles, nbytes)
}

var verbose = flag.Bool("v", false, "显示运行消息")

// 每隔500毫秒打印一次结果
func du2() {
	flag.Parse()
	roots := flag.Args()
	//如果未输入目录，则输出当前目录
	if len(roots) == 0 {
		roots = []string{"."}
	}
	fileSizes := make(chan int64)
	go func() {
		for _, root := range roots {
			walkDir(root, fileSizes)
		}
		close(fileSizes)
	}()

	//每隔500毫秒打印一次结果
	var tick <-chan time.Time
	if *verbose {
		tick = time.Tick(500 * time.Millisecond)
	}
	var nfiles, nbytes int64
loop:
	for {
		select {
		case size, ok := <-fileSizes:
			//管道是否关闭
			if !ok {
				break loop //退出标签，如果只使用break只会退出select
			}
			nfiles++
			nbytes += size
		case <-tick:
			printDiskUsage(nfiles, nbytes)
		}
	}
	//最后打印总结果
	printDiskUsage(nfiles, nbytes)
}

// 限制并发数量的通道信号量
var sema = make(chan struct{}, 20)

var wg sync.WaitGroup

// 用于发送退出消息的通道
var done = make(chan struct{})

// 是否已取消
func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

// 开启多协程并发处理
func du3() {
	flag.Parse()
	roots := flag.Args()
	//如果未输入目录，则输出当前目录
	if len(roots) == 0 {
		roots = []string{"."}
	}
	//输入任意键时发出取消消息，就是关闭取消管道
	go func() {
		os.Stdin.Read(make([]byte, 1))
		close(done)
	}()
	fileSizes := make(chan int64)

	go func() {
		for _, root := range roots {
			wg.Add(1)
			walkDir(root, fileSizes)
		}
	}()
	go func() {
		//等待
		wg.Wait()
		close(fileSizes)
	}()

	//每隔500毫秒打印一次结果
	var tick <-chan time.Time
	if *verbose {
		tick = time.Tick(500 * time.Millisecond)
	}
	var nfiles, nbytes int64
loop:
	for {
		select {
		case <-done: //是否取消
			fmt.Println("收到取消信号")
			//结束之前需要把fileSizes通道的内容排空
			for range fileSizes {
				//什么也不做
			}
			return
		case size, ok := <-fileSizes:
			//管道是否关闭
			if !ok {
				break loop //退出标签，如果只使用break只会退出select
			}
			nfiles++
			nbytes += size
		case <-tick:
			printDiskUsage(nfiles, nbytes)
		}
	}
	//最后打印总结果
	printDiskUsage(nfiles, nbytes)
}

// 递归遍历目录，统计所有文件大小
func walkDir(dir string, fileSizes chan<- int64) {
	defer wg.Done()
	//是否已取消
	if cancelled() {
		return
	}
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			//fmt.Println(subdir)
			wg.Add(1)
			go walkDir(subdir, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}

// 读取给定目录下的所有目录及文件
func dirents(dir string) []os.FileInfo {
	select {
	case sema <- struct{}{}: //控制并发数量
	case <-done: //是否取消
		return nil
	}
	defer func() { <-sema }()

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du1:%v\n", err)
		return nil
	}
	return entries
}

// 打印磁盘使用情况
// @param	nfiles	文件总数
// @params	nbytes	文件总大小
func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files %1.f GB\n", nfiles, float64(nbytes)/1e9)
}
