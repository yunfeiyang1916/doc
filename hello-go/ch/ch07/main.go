// 数组、切片
package main

import (
	"bytes"
	"fmt"
	"sort"
)

// 传值数组
func valueArray(a [10]int) {
	fmt.Printf("a 的类型为：%T\n", a)
	for i := 0; i < len(a); i++ {
		a[i] = 2 * i
		fmt.Printf("%d ", a[i])
	}
	fmt.Println()
}

// 传数组指针
func pointerArray(a *[10]int) {
	fmt.Printf("a 的类型为：%T\n", a)
	for i := 0; i < len(a); i++ {
		a[i] = 2 * a[i]
		fmt.Printf("%d ", a[i])
	}
	fmt.Println()
}

// 数组测试
func arrTest() {
	var arr [10]int = [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	//使用new生成的是数组指针
	var arrP *[10]int = new([10]int)
	var slice = [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Printf("arr Type is %T\n", arr)
	fmt.Printf("arrP Type is %T\n", arrP)
	fmt.Printf("slice Type is %T\n", slice)
	//fmt.Println(arr)
	//fmt.Println(arrP)
	//fmt.Println(slice)
	//传值数组
	valueArray(arr)
	fmt.Println(arr)
	for i := 0; i < len(arrP); i++ {
		fmt.Printf("%d ", arrP[i])
	}
	fmt.Println()
	//传指针数组
	pointerArray(&arr)
	for _, v := range arr {
		fmt.Printf("%d ", v)
	}
	fmt.Println()
}

// 数组测试2
func arrTest2() {
	//初始化数组（数组常量）
	var arr = [10]int{1, 2, 3}
	fmt.Println(arr)
	fmt.Printf("容量是：%d\n", cap(arr))
	var strArr = [10]string{3: "呵呵", 5: "你好", 9: "hello"}
	fmt.Println(strArr)
	//多维数组
	var arr2 = [3][3]int{{1, 2, 3}, {4, 5, 6}}
	for _, a := range arr2 {
		//fmt.Printf("a的类型是%T\n", a)
		for _, b := range a {
			fmt.Printf("%d ", b)
		}
		fmt.Println()
	}
}

// 切片声明
func sliceDeclared() {
	//先声明一个数组
	var arr = [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	//var arrP = &arr
	var slice1 = arr[:len(arr)-1]
	fmt.Printf("slice1类型是%T 长度是%d 容量是%d\n", slice1, len(slice1), cap(slice1))
	fmt.Println(slice1)
	slice1[4] = slice1[4] * 2
	//扩展到上限
	slice1 = slice1[:cap(slice1)]
	fmt.Printf("slice1类型是%T 长度是%d 容量是%d\n", slice1, len(slice1), cap(slice1))
	fmt.Println(slice1)
	fmt.Println("数组arr=", arr)
	var slice3 = slice1[4:5]
	fmt.Printf("slice3类型是%T 长度是%d 容量是%d\n", slice3, len(slice3), cap(slice3))
	fmt.Println("slice3=", slice3)
	//可这么声明切片
	var slice2 = []int{1, 2, 3, 4} //[4]int{1, 2, 3, 4}[:]
	fmt.Printf("slice2类型是%T 长度是%d 容量是%d\n", slice2, len(slice2), cap(slice2))
	var sp = &slice2
	fmt.Printf("sp类型是%T 长度是%d 容量是%d\n", sp, len(*sp), cap(*sp))
	for i := 0; i < len(*sp); i++ {
		fmt.Printf("%d ", (*sp)[i])
	}
}

// 切片声明2
func sliceDeclared2() {
	ss := []int{1, 2, 3}
	fmt.Printf("&ss=%p ss=%p 类型=%T 长度=%d 容量=%d \n", &ss, ss, ss, len(ss), cap(ss))
	sliceDeclared2Func(ss)
	ss = append(ss, 4)
	fmt.Printf("&ss=%p ss=%p 类型=%T 长度=%d 容量=%d \n", &ss, ss, ss, len(ss), cap(ss))
	sliceDeclared2Func(ss)
	copy(ss, []int{5, 6, 7})
	fmt.Println(ss)
	fmt.Printf("&ss=%p ss=%p 类型=%T 长度=%d 容量=%d \n", &ss, ss, ss, len(ss), cap(ss))
}

// 切片声明2函数调用
func sliceDeclared2Func(ss []int) {
	fmt.Printf("&ss=%p ss=%p 类型=%T 长度=%d 容量=%d \n", &ss, ss, ss, len(ss), cap(ss))
}

// 切片函数
func sliceFunc(slice []int) {
	slice = []int{1, 2, 3, 4}
	for k := range slice {
		slice[k] *= 2
	}
	fmt.Println(slice)
}

// 切片函数测试
func sliceTest() {
	var slice = []int{1, 2, 3, 4}
	sliceFunc(slice)
	fmt.Println(slice)
}

// make与new测试
func makeTest() {
	//make(T)返回一个类型为T的初始值，它只适用于3种内建的引用类型：切片、map和channel
	//make原型：func make([]T, len, cap)，其中 cap 是可选参数
	//当不需要先声明数组时可以使用make创建切片
	var slice1 []int = make([]int, 5 /*, 10*/)
	fmt.Printf("slice1类型是%T 长度是%d 容量是%d\n", slice1, len(slice1), cap(slice1))
	fmt.Println(slice1)
	//new(T)为每个新的类型T分配一片内存，初始化为0并且返回类型为*T的内存地址：这种方法返回一个指向类型为T，值为0的地址指针，
	//它适用于值类型如数组和结构体，它相当于&T{}
	//使用new分配切片
	var sp = new([10]int)
	sp[0] = 10
	(*sp)[1] = 20
	var slice2 = sp[:5] //var slice2 = new([10]int)[:5]

	fmt.Printf("slice2类型是%T 长度是%d 容量是%d\n", slice2, len(slice2), cap(slice2))
	fmt.Println(slice2)
}

// bytes包测试
func buffTest() {
	var buff bytes.Buffer
	buff.WriteString("张三")
	buff.WriteString(",李四")
	fmt.Println(buff.String())
}

// 切片For-Range测试
func sliceForRangeTest() {
	var slice1 = []int{1, 2, 3, 4, 5}
	fmt.Println(slice1)
	for i := 0; i < len(slice1); i++ {
		slice1[i] = 2 * slice1[i]
	}
	fmt.Println(slice1)
	//v只是slice1某个索引位置的值的一个拷贝，不能用来修改slice2该索引位置的值
	for _, v := range slice1 {
		v = 2 * v
	}
	fmt.Println(slice1)
}

// 切片重组测试
func resliceTest() {
	var slice1 = make([]int, 0, 10)
	fmt.Printf("slice1类型是%T 长度：%d 容量：%d \n", slice1, len(slice1), cap(slice1))
	fmt.Println(slice1)
	slice1 = slice1[5:5]
	fmt.Printf("slice1类型是%T 长度：%d 容量：%d \n", slice1, len(slice1), cap(slice1))
	fmt.Println(slice1)
	for i := 0; i < cap(slice1); i++ {
		slice1 = slice1[:i+1]
		slice1[i] = i
	}
	for _, v := range slice1 {
		fmt.Printf("%d ", v)
	}
	fmt.Println()
}

// 切片拷贝测试
func copyTest() {
	from := []int{1, 2, 3, 4, 5, 6}
	to := make([]int, 10)
	n := copy(to, from)
	fmt.Println(to)
	fmt.Println(n)
	// arrFrom := [5]int{1, 2, 3, 4, 5}
	// arrTo := [6]int{5: 6}
	// //copy不支持数组
	// n = copy(arrTo, arrFrom)
	// fmt.Println(arrTo)
	// fmt.Println(n)
}

// 追加测试
func appendTest() {
	from := []int{1, 2, 3, 4, 5, 6}
	to := make([]int, 5, 10)
	to = append(to, from...)
	to = append(to, 7, 8, 9)
	fmt.Println(to)
	fmt.Printf("to类型是%T 长度：%d 容量：%d \n", to, len(to), cap(to))
	slice1 := make([]byte, 5, 10) //[]byte{1, 2, 3, 4, 5}
	slice1 = AppendByte(slice1, 6, 7, 8)
	fmt.Printf("slice1类型是%T 长度：%d 容量：%d \n", slice1, len(slice1), cap(slice1))
	fmt.Println(slice1)
}

// 追加字节
func AppendByte(slice []byte, data ...byte) []byte {
	m := len(slice)
	n := m + len(data)
	if n > cap(slice) {
		newSlice := make([]byte, (n+1)*2)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0:n]
	copy(slice[m:n], data)
	return slice
}

// 传值
func passValue(slice []int) {
	fmt.Printf("地址：%p\n", slice)
	for i := 0; i < len(slice); i++ {
		slice[i] *= 2
	}
	from := []int{1, 2, 3, 4, 5}
	slice = slice[:len(slice)-1]
	fmt.Printf("len=%d cap=%d\n", len(slice), cap(slice))
	//copy(slice, from)
	slice = append(slice[0:5], from...)
	fmt.Println(slice)
	fmt.Printf("len=%d cap=%d\n", len(slice), cap(slice))
	fmt.Printf("地址：%p\n", slice)
}

// 传值测试
func passValueTest() {
	var slice1 []int //make([]int, 6, 12)
	fmt.Printf("地址：%p 长度%d 是否为nil:%v\n", slice1, len(slice1), slice1 == nil)
	fmt.Println(slice1)
	passValue(slice1)
	fmt.Printf("len=%d cap=%d\n", len(slice1), cap(slice1))
	fmt.Println(slice1)
}

// 字符串切片
func strSlice() {
	var str = "世界" //"Hello World!"
	slice1 := []byte(str)
	for _, v := range slice1 {
		fmt.Printf("%U->%c ", v, v)
	}
	fmt.Println()
	str = string(slice1)
	fmt.Println(slice1)
	fmt.Println(str)
	slice2 := []int32(str)
	for _, v := range slice2 {
		fmt.Printf("%U->%c ", v, v)
	}
	fmt.Println()
	str = string(slice2)
	fmt.Println(slice2)
	fmt.Println(str)
	r := []rune(str)
	for _, v := range r {
		fmt.Printf("%U->%c ", v, v)
	}
	fmt.Println()
	str = string(r)
	fmt.Println(r)
	fmt.Println(str)
}

// 截取字符串
func substrTest() {
	str := "张三，你好啊"
	str1 := str[0:1]
	fmt.Println(str1)
	str2 := str[1:2]
	fmt.Println(str2)
}

// 字节数组比较函数
func Compare(a, b []byte) int {
	for i := 0; i < len(a) && i < len(b); i++ {
		switch {
		case a[i] > b[i]:
			return 1
		case a[i] < b[i]:
			return -1
		}
	}
	//数组的长度可能不同
	switch {
	case len(a) > len(b):
		return 1
	case len(a) < len(b):
		return -1
	}
	//走到这说明数组相等
	return 0
}

// 字节数组比较测试
func compareTest() {
	slice1 := []byte{1, 3, 3, 4}
	slice2 := []byte{1, 2, 3, 4}
	fmt.Println(Compare(slice1, slice2))
}

// 整型排序
func sortInt(slice []int) {
	for i := 0; i < len(slice)-1; i++ {
		for j := i + 1; j < len(slice); j++ {
			if slice[i] > slice[j] {
				//并行交换值
				slice[i], slice[j] = slice[j], slice[i]
			}
		}
	}
}

// 排序测试
func sortTest() {
	var slice1 = []int{1, 2, 10, 9, 5, 6, 7, 8, 1, 2, 3}
	//sortInt(slice1)
	sort.Ints(slice1)
	fmt.Println(sort.SearchInts(slice1, 2))
	for _, v := range slice1 {
		fmt.Printf("%d ", v)
	}
	fmt.Println()
}

func main() {
	//arrTest()
	//arrTest2()
	//sliceDeclared()
	sliceDeclared2()
	//sliceTest()
	//makeTest()
	//buffTest()
	//sliceForRangeTest()
	//resliceTest()
	//copyTest()
	//appendTest()
	//passValueTest()
	//strSlice()
	//substrTest()
	//compareTest()
	//sortTest()
}
