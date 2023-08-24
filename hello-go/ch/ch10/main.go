// 结构体测试
package main

import (
	"fmt"
	"project/hello/ch10/structPack"
	"unsafe"
)

// 结构体值测试
func pointValueTest() {
	var p structPack.Point
	p.Y = 10
	name := structPack.Name
	fmt.Println(name)
	fmt.Printf("p 类型为：%T 地址：%p\n", p, &p)
	fmt.Println(p)
}

// 结构体指针测试
func pointPointerTest() {
	var p *structPack.Point = new(structPack.Point)
	p.X = 1
	p.Y = 2
	fmt.Printf("p 类型为：%T 地址：%p\n", p, &p)
	fmt.Println(*p)
	//混合字面量语法，底层仍然会调用new()
	p1 := &structPack.Point{X: 2, Y: 10}
	fmt.Println(p1)
}

// 人员结构体测试
func personTest() {
	//值类型
	var p1 structPack.Person
	p1.FirstName = "zhang"
	p1.LastName = "shan"
	structPack.UpPerson(&p1)
	fmt.Printf("FirstName=%s LastName=%s\n", p1.FirstName, p1.LastName)
	//指针测试
	p2 := new(structPack.Person)
	p2.FirstName = "li"
	p2.LastName = "si"
	structPack.UpPerson(p2)
	fmt.Printf("FirstName=%s LastName=%s\n", p2.FirstName, p2.LastName)
	//混合字母量语法
	p3 := &structPack.Person{FirstName: "wang", LastName: "wu"}
	structPack.UpPerson(p3)
	fmt.Printf("FirstName=%s LastName=%s\n", p3.FirstName, p3.LastName)
}

// 矩形测试
func rectTest() {

	var rect structPack.Rect
	rect.Height.X = 1
	rect.Height.Y = 2
	rect.Width.X = 4
	rect.Width.Y = 8
	//fmt.Println(rect)
	//fmt.Printf("%p %p %p %p %p %p %p\n", &rect, &(rect.Width), &(rect.Width.X), &(rect.Width.Y), &(rect.Height), &(rect.Height.X), &(rect.Height.Y))

	var rect2 structPack.Rect2
	rect2.Height = new(structPack.Point)
	rect2.Height.X = 1
	rect2.Height.Y = 2
	rect2.Width = new(structPack.Point)
	rect2.Width.X = 4
	rect2.Width.Y = 8
	fmt.Println(rect)
	fmt.Printf("%p %p %p %p %p %p %p\n", &rect2, rect2.Width, &(rect2.Width.X), &(rect2.Width.Y), rect2.Height, &(rect2.Height.X), &(rect2.Height.Y))
}

// 结构体转换
func structConvert() {
	p1 := structPack.Point{X: 1, Y: 2}
	var p2 = structPack.Po{X: 2, Y: 4}
	//将Point别名转成Point类型
	var p3 structPack.Point = structPack.Point(p2)
	fmt.Printf("p1类型为%T\n", p1)
	fmt.Printf("p2类型为%T\n", p2)
	fmt.Printf("p3类型为%T\n", p3)
	fmt.Println(p1, p2, p3)
}

// 工厂方法测试
func factoryTest() {
	file := structPack.NewFile(1, "新文件")
	fmt.Println(file)
	m := structPack.NewMatrix(1, 2)
	m.X *= 2
	m.Y *= 2
	fmt.Println(m)
	fmt.Println(unsafe.Sizeof(m))
}

func main() {
	//pointValueTest()
	//pointPointerTest()
	//personTest()
	//rectTest()
	//structConvert()
	//factoryTest()
	//structPack.NewAndMake()
	//structPack.TagTest()
	//structPack.AnonymousTest()
	//structPack.NamingConflicts()
	//structPack.MethodTest()
	//structPack.MethodTest2()
	//structPack.CombFuncTest()
	//structPack.MultipleInherit()
	structPack.StringFormat()
}
