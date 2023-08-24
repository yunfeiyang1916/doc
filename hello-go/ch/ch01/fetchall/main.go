// 并发批量调用下载
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// 协程同步组
var waitGroup = sync.WaitGroup{}

func main() {
	//使用waitGroup同步协程
	//useWaitGroupSync()
	//使用管道同步协程
	useChanSync()
}

// 使用waitGroup同步协程
func useWaitGroupSync() {
	start := time.Now()
	//ch := make(chan string)
	for _, url := range os.Args[1:] {
		if !strings.HasPrefix(url, "http://") {
			url = "http://" + url
		}
		go fetch(url)
		waitGroup.Add(1)
	}
	waitGroup.Wait()
	fmt.Printf("总耗时%.2f秒.\n", time.Since(start).Seconds())
}

// 获取请求
func fetch(url string) {
	//waitGroup.Add(1)
	defer waitGroup.Done()
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}
	defer resp.Body.Close()
	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Fprintf(os.Stdout, "请求:%s 耗时:%.2f秒 响应状态码:%d 内容长度：%7d字节\n", url, time.Since(start).Seconds(), resp.StatusCode, nbytes)
}

// 使用管道同步
func useChanSync() {
	start := time.Now()
	ch := make(chan string)
	for _, url := range os.Args[1:] {
		if !strings.HasPrefix(url, "http://") {
			url = "http://" + url
		}
		go fetch2(url, ch)
	}
	for range os.Args[1:] {
		fmt.Println(<-ch)
	}
	fmt.Printf("总耗时%.2f秒.\n", time.Since(start).Seconds())
}

// 获取请求
func fetch2(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- err.Error()
		return
	}
	defer resp.Body.Close()
	nbytes, err := io.Copy(os.Stdout, resp.Body) //io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		ch <- err.Error()
		return
	}
	ch <- fmt.Sprintf("请求:%s 耗时:%.2f秒 响应状态码:%d 内容长度：%7d字节", url, time.Since(start).Seconds(), resp.StatusCode, nbytes)
}
