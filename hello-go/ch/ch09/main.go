// 标准包测试
package main

import (
	"fmt"
	"math/big"
	"project/hello/ch09/pack1"
	"regexp"
	"strconv"
	"sync"
)

// 正则测试
func regexpText() {
	str := "张三：2578.34 李四：4567.23 王五：5632.18"
	//正则
	pat := "[0-9]+.[0-9]+"
	if ok, _ := regexp.Match(pat, []byte(str)); ok {
		fmt.Println("完全匹配")
	}
	re, _ := regexp.Compile(pat)
	//将匹配到的部分替换为"##.#
	str2 := re.ReplaceAllString(str, "##.#")
	fmt.Println(str2)
	//参数为函数时
	str3 := re.ReplaceAllStringFunc(str, func(s string) string {
		v, _ := strconv.ParseFloat(s, 32)
		return strconv.FormatFloat(v*2, 'f', 2, 32)
	})
	fmt.Println(str3)
}

// 锁结构体
type Info struct {
	//互斥锁
	mu sync.Mutex
	//读写锁
	rw sync.RWMutex
}

// 锁同步
func syncTest() {
	info := new(Info)
	info.mu.Lock()
	//todo
	info.mu.Unlock()
	info.rw.RLock()
	//todo
	info.rw.RUnlock()
}

// 精密计算
func bigTest() {
	im := big.NewInt(1)
	in := im
	io := big.NewInt(4)
	ip := big.NewInt(2)
	//ip.Mul(im, in).Add(ip, im).Div(ip, io)
	ip.Add(in, im).Mul(ip, io)
	fmt.Println(ip)
}

// 自定义包测试
func defPack() {
	str := pack1.RetrunStr()

	fmt.Printf("%d %f %s\n", pack1.I, pack1.Pi, str)
}
func main() {
	//regexpText()
	//bigTest()
	defPack()
}
