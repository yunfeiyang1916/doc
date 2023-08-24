// 测试变量作用域
package main

import (
	"fmt"
	"os"
	"runtime"
	"unsafe"
)

var a = "O"

const Pi = 3

func init() {
	fmt.Println("调用ch04的初始化函数")
	//fmt.Println(c + v)
}

// 测试变量作用域
func varScope() {
	n()
	m()
	n()
}

func n() {
	println(a)
}
func m() {
	a := "N"
	println(a)
}

// 类型别名
type IZ int

// 类型测试
func typeTest() {
	var i int = 1
	var i32 int32 = 32
	var i64 int64 = 64
	if i < int(i32) {
		fmt.Printf("i64=%d", i64)
	}
	if i < Pi {
		fmt.Printf("Pi=%d", Pi)
	}
	fmt.Println(unsafe.Sizeof(i))
	//fmt.println("测试int类似是否能与int32直接比较", i < i32)
	//fmt.println("测试int类似是否能与int32直接比较", i < i32)
}

// 常量测试
const (
	Monday, Tuesday, Wednesday = 1, 2, 3
	Thursday, Friday, Saturday = 4, 5, 6
)

// 如果有常量没有赋值，则使用上一个常量的值，比如e=d,f=e
const (
	d IZ = 1
	e
	f
)

// h=g=iota,j=h=iota。 iota从零开始，每当iota在新的一行使用时值自动+1，在每遇到一个新的常量块或单个常量生命时，iota都会重置为0
const (
	g = iota
	h
	j
)

// iota被重置为0
const k = iota

// 常量测试
func constTest() {
	fmt.Printf("常量d=%v e=%v f=%v g=%v h=%v j=%v k=%v", d, e, f, g, h, j, k)
	fmt.Println()
}

// 变量测试
var l int

func varTest() {
	fmt.Printf("变量l=%v e=%v f=%v g=%v h=%v j=%v k=%v", d, e, f, g, h, j, k)
	fmt.Println()
	getOS()
}

// 获取操作系统信息
func getOS() {
	var goOS string = runtime.GOOS //os.Getenv("OS")
	fmt.Printf("操作系统是：%s\n", goOS)
	path := os.Getenv("PATH")
	fmt.Println("Path=", path)
}

// 交换俩值
func swap(a *int, b *int) {
	*a, *b = *b, *a
}

// 测试交换值
func swapTest() {
	defer fmt.Println("结束交换值测试")
	a, b := 1, 2
	fmt.Printf("初始值a=%v b= %v\n", a, b)
	//交换俩值
	swap(&a, &b)
	fmt.Printf("交换后a=%v b=%v\n", a, b)
}
