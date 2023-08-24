// 把map变量限制在一个单独的监控协程中的缓存
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

// 函数类型
type Func func(key string) (interface{}, error)

// 结果
type result struct {
	value interface{}
	err   error
}

// 入口
type entry struct {
	//结果
	res result
	//当结果准备好的时候需要关闭的通道
	ready chan struct{}
}

// 请求
type request struct {
	key string
	//只读的响应通道
	response chan<- result
}

// 缓存
type Memo struct {
	//请求通道
	requests chan request
}

func (memo *Memo) Get(key string) (interface{}, error) {
	response := make(chan result)
	memo.requests <- request{key, response}
	res := <-response
	return res.value, res.err
}
func (this *Memo) server(f Func) {
	cache := make(map[string]*entry)
	for req := range this.requests {
		e := cache[req.key]
		if e == nil {
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e
			go e.call(f, req.key)
		}
		go e.deliver(req.response)
	}
}
func (this *entry) call(f Func, key string) {
	this.res.value, this.res.err = f(key)
	close(this.ready)
}
func (e *entry) deliver(response chan<- result) {
	// Wait for the ready condition.
	<-e.ready
	// Send the result to the client.
	response <- e.res
}
func New(f Func) *Memo {
	memo := &Memo{requests: make(chan request)}
	go memo.server(f)
	return memo
}

// 要下载的url集合
var urls = []string{"https://e.360.cn", "https://baidu.com", "https://12306.cn", "https;//abc.com"}

// 请求并且获取http响应body
func httpGetBody(url string) (interface{}, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
func main() {
	now := time.Now()
	m := New(httpGetBody)
	//使url集合重复一遍
	urls = append(urls, urls...)
	//使url集合重复一遍
	urls = append(urls, urls...)
	wg := sync.WaitGroup{}
	for _, url := range urls {
		wg.Add(1)
		//将url传进去，防止闭包
		go func(url string) {
			defer wg.Done()
			start := time.Now()
			value, err := m.Get(url)
			if err != nil {
				log.Println(err)
			}
			v, ok := value.([]byte)
			var length int
			if ok {
				length = len(v)
			}
			fmt.Printf("%s\t%s\t%d\tbytes\n", url, time.Since(start), length)
		}(url)
	}
	wg.Wait()
	fmt.Println("执行总耗时：", time.Since(now))
}
