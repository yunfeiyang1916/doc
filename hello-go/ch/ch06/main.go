// 函数测试
package main

import (
	"fmt"
	"io"
	"log"
	"runtime"
	"strings"
	"time"
)

// 非命名返回双值
func doubleReturn(i int) (int, int) {
	return 2 * i, 3 * i
}

// 命名返回双值
func doubleReturn2(i int) (a, b int) {
	a, b = 2*i, 3*i
	return
	//return a,b
}

// 命名函数测试
func mingMingTest() {
	var a, b int
	a, b = doubleReturn(5)
	a, b = doubleReturn2(10)
	fmt.Printf("a=%d b=%d\n", a, b)
}

// 变参函数
func min(a ...int) int {
	minArray(a)
	if len(a) == 0 {
		return 0
	}
	min := a[0]
	for index := 0; index < len(a); index++ {
		if min > a[index] {
			min = a[index]
		}
	}
	return min
	// for _, i := range a {
	// 	if min > i {
	// 		min = i
	// 	}
	// }
	// return min
}
func minArray(s []int) {

}

// 变参测试
func minTest() {
	fmt.Println(min(1, 2, 3))
}

// 使用defer实现代码追踪
func trace(s string) string {
	fmt.Println("开始：", s)
	return s
}
func untrace(s string) {
	fmt.Println("结束：", s)
}

// 记录函数的参数与返回值
func logDefer(s string) (n int, err error) {
	defer func() {
		fmt.Printf("logDefer(%q)=%d,%v \n", s, n, err)
	}()
	return 7, io.EOF
}

// defer中修改返回值
func deferReturn() int {
	i := 1
	defer func() {
		i = 2
		fmt.Println(i)
	}()
	return i
}

// defer测试
func deferTest() {
	// fmt.Println("-----1----------")
	// defer fmt.Println("--------defer-----------")
	// fmt.Println("------2-------")
	// i := 0
	// defer fmt.Println(i)
	// i++
	//trace("defer测试")
	//defer untrace("defer测试")
	//defer untrace(trace("defer测试"))
	//fmt.Println("执行测试")
	//logDefer("呵呵")
	println(deferReturn())
}

// 计算斐波那契数列，即前两个数为1，从第三个数开始每个数均为前两个数之和
func fibonacci(n int) (res int) {
	if n <= 1 {
		res = 1
	} else {
		res = fibonacci(n-1) + fibonacci(n-2)
	}
	fmt.Printf("%d ", res)
	return res
}

// 递归测试
func recursionTest() {
	fibonacci(10)
}

// 求和
func sum(a, b int) int {
	return a + b
}

// 回调函数
func callback(f func(int, int) int) {
	if f != nil {
		s := f(1, 2)
		fmt.Println(f)
		fmt.Println(s)
		fmt.Println(strings.IndexFunc("啊578", isAscii))
	}
}

// 给定字符是否是Ascii
func isAscii(c rune) bool {
	return c < 255
}

// 匿名函数测试
func anonymityTest() {
	for i := 0; i < 4; i++ {
		g := func(a int) {
			fmt.Printf("%d ", a)
		}
		g(i)
		fmt.Printf(" - g 的类型是%T 地址是%p\n", g, g)
	}
}

// 添加，返回一个函数
func add2() func(int) int {
	return func(b int) int {
		return b + 2
	}
}

// 单参数的添加，返回一个函数
func adder(a int) func(int) int {
	return func(b int) int {
		return a + b
	}
}

// 返回一个计算累加的函数
func adder2() func(int) int {
	x := 0
	return func(b int) int {
		x += b
		return x
	}
}

// 工厂函数，返回一个用于追加后缀的函数
func makeAddSuffix(suffix string) func(string) string {
	return func(name string) string {
		if !strings.HasSuffix(name, suffix) {
			return name + suffix
		}
		return name
	}
}

// 闭包测试
func closureTest() {
	f := add2()
	f2 := adder(4)
	f3 := adder2()
	fmt.Printf("调用%T %d\n", f, f(2))
	fmt.Printf("调用%T %d\n", f2, f2(2))
	fmt.Printf("调用%T %d %d %d\n", f3, f3(2), f3(10), f3(100))
	//增加bmp后缀函数
	addBmp := makeAddSuffix(".bmp")
	//增加jpg函数
	addJpg := makeAddSuffix(".jpg")
	fmt.Printf("调用%T %s\n", addBmp, addBmp("hehe"))
	fmt.Printf("调用%T %s\n", addJpg, addJpg("hehe.jpg"))
}

// 使用闭包调试
func closureDebug() {
	start := time.Now()
	where := func() {
		_, file, line, _ := runtime.Caller(1)
		log.Printf("%s:%d", file, line)
	}
	log.Println("===开始==")
	where()
	log.Println("===执行==")
	where()
	log.Println("===结束==")
	where2 := log.Print
	log.Println("===开始==")
	where2()
	log.Println("===执行==")
	where2()
	log.Println("===结束==")
	end := time.Now()
	log.Printf("程序运行耗时：%d 秒", end.Sub(start)/100000000)
}

func main() {
	//mingMingTest()
	//minTest()
	//deferTest()
	//recursionTest()
	//callback(sum)
	//anonymityTest()
	//closureTest()
	//closureDebug()
	//fmt.Println(add.Add(1, 16))
	fmt.Println("我是主线程")
	go func() {
		fmt.Println("我是新协程")
	}()
	time.Sleep(1 * 1e6)
}
