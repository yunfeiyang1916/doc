// 并发非阻塞缓存
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

// 缓存
type Memo struct {
	//值不存在时的取值函数
	f Func
	//字典
	cache map[string]*entry
	//互斥锁
	sync.Mutex
}

// 入口
type entry struct {
	//结果
	res result
	//当结果准备好的时候需要关闭的通道
	ready chan struct{}
}

// 实例化一个新的缓存
func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]*entry)}
}

// 从缓存中读取，缓存中不存在时自动填充
func (memo *Memo) Get(key string) (interface{}, error) {
	memo.Lock()
	e := memo.cache[key]
	if e == nil {
		e = &entry{ready: make(chan struct{})}
		memo.cache[key] = e
		memo.Unlock()
		e.res.value, e.res.err = memo.f(key)
		close(e.ready)
	} else {
		memo.Unlock()
		//等待数据准备好
		<-e.ready
	}
	return e.res.value, e.res.err
}

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

// 要下载的url集合
var urls = []string{"https://e.360.cn", "https://baidu.com", "https://12306.cn", "https;//abc.com"}

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
