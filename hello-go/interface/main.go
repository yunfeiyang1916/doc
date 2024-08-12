package main

import (
	"fmt"
	"unsafe"
)

type TestStruct struct {
}

func NilOrNot(v interface{}) bool {
	return v == nil
}

func main() {
	//var s *TestStruct
	//fmt.Println(s, "  ", s == nil)
	//fmt.Println(NilOrNot(s), "  ", NilOrNot(nil))
	fmt.Printf("string size: %d align: %d\n", unsafe.Sizeof(""), unsafe.Alignof(""))
	fmt.Printf("size: %d align: %d\n", unsafe.Sizeof(Align{}), unsafe.Alignof(Align{}))
	fmt.Printf("size: %d align: %d\n", unsafe.Sizeof(User{}), unsafe.Alignof(User{}))
}

// 测试内存对齐的结构体
type Align struct {
	a bool
	b string
	z struct{}
	c int16
}

type User struct {
	a int32
	b []int32
	c string
	d bool
	e struct{}
}
