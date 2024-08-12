package main

import "fmt"

// 逃逸分析，当将指针返回时会逃逸
func a() *int {
	x := 1
	return &x
}

// 空接口做实参不一定会逃逸，只有在使用了反射时才会逃逸（因为反射会要求变量分配在堆上）
func b(v interface{}) {

}

func c(i interface{}) {
	// fmt.Println 内部用到了反射，所以会逃逸
	fmt.Println(i)
}
