// 切片测试
package main

import (
	"fmt"
	"unsafe"
)

func main() {
	//sliceSet1()
	emptySlice()
}

// 切片赋值
func sliceSet1() {
	//整型切片操作共享同一个底层数组
	arr := [5]int{1, 2, 3, 4, 5}
	ss := arr[:]
	ss2 := ss[0:2]
	fmt.Printf("%p %p %p\n", &arr, ss, ss2)
	//字节切片操作也共享同一个底层数组
	byteArr := [5]byte{1, 2, 3, 4, 5}
	byteS := byteArr[:]
	byteS2 := byteS[0:2]
	fmt.Printf("%p %p %p\n", &byteArr, byteS, byteS2)

	str := "abcd1344"
	str2 := str[:2]
	fmt.Printf("%p %p \n", &str, &str2)
}

// 空切片测试
func emptySlice() {
	var s []int
	fmt.Printf("占用字节：%d 地址：%p %p 是否为nil:%v 长度：%d\n", unsafe.Sizeof(s), s, &s, s == nil, len(s))
	s2 := []int{}
	fmt.Printf("占用字节：%d 地址：%p %p 是否为nil:%v 长度：%d\n", unsafe.Sizeof(s2), s2, &s2, s2 == nil, len(s2))
	var str string = "你好啊"
	fmt.Printf("占用字节：%d 地址：%p %p 是否为空:%v 长度：%d\n", unsafe.Sizeof(str), str, &str, str == "", len(str))
}
