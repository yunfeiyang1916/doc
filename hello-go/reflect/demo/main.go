// 反射测试
package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"time"
	"unsafe"
)

type Person struct {
	ID      int
	Name    string
	Address []string
}

func main() {
	//test1()
	//test2()
	//canAddrTest()
	//setTest()
	//performTest()
	interfaceSetTest()
}

// 反射测试
func test1() {
	slice := []int{1, 2, 3, 4}
	t := reflect.TypeOf(slice)
	v := reflect.ValueOf(slice)
	fmt.Println(t, v)
	var i interface{}
	i = slice
	t = reflect.TypeOf(i)
	v = reflect.ValueOf(i)
	r := reflect.Indirect(v)
	fmt.Println(t.Kind(), v.Kind(), r.Kind())
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	fmt.Printf("%p %v\n", slice, strconv.FormatInt(int64(sliceHeader.Data), 16))
	fmt.Println(sliceHeader, *sliceHeader)
	fmt.Println(reflect.ValueOf(&slice).Kind(), reflect.Indirect(reflect.ValueOf(&slice)).Kind())
}

// 测试2
func test2() {
	var w io.Writer = os.Stdout
	fmt.Println(reflect.TypeOf(w), reflect.ValueOf(w).Interface())
	x, ok := w.(io.Writer)
	fmt.Println(x, ok)
	switch w.(type) {
	case *os.File:
		fmt.Println("w是*os.File类型")
	case io.Writer:
		fmt.Println("w实现了io.Writer接口")
	}
}

// 可寻址测试
func canAddrTest() {
	x := 2
	a := reflect.ValueOf(2)
	b := reflect.ValueOf(x)
	c := reflect.ValueOf(&x)
	d := c.Elem()
	fmt.Println(d.Interface())
	fmt.Printf("a可寻址? %v\n", a.CanAddr())
	fmt.Printf("b可寻址? %v\n", b.CanAddr())
	fmt.Printf("c可寻址? %v\n", c.CanAddr())
	fmt.Printf("d可寻址? %v\n", d.CanAddr())
	d.Set(reflect.ValueOf(123))
	fmt.Println(x)
	px := d.Addr().Interface().(*int)
	*px = 900
	fmt.Println(x)
}

// 赋值测试
func setTest() {
	//基础类型
	x := 1
	rx := reflect.ValueOf(&x).Elem()
	fmt.Println(reflect.ValueOf(&x).Type(), rx.Type())
	rx.SetInt(2)
	fmt.Println(x)
	rx.Set(reflect.ValueOf(3))
	fmt.Println(x)
	//rx.SetString("hehe")
	//接口
	var y interface{}
	ry := reflect.ValueOf(&y).Elem()
	fmt.Println(reflect.ValueOf(&y).Type(), ry.Type())
	//接口类型反射赋值只能使用Set,而不能使用SetInt、SetString等
	//ry.SetInt(2)
	ry.Set(reflect.ValueOf(56))
	fmt.Println(y)
	ry.Set(reflect.ValueOf("你好啊"))
	fmt.Println(y)
	//反射可以读取未导出的变量，但是不能修改值
	stdout := reflect.ValueOf(os.Stdout).Elem()
	fmt.Println(stdout.Type())
	name := stdout.FieldByName("name")
	fmt.Println(name.String())
	fmt.Println(name.CanAddr(), name.CanSet())
	//name.SetString("你好啊")
	//切片
	s := []int{1, 2, 3}
	rs := reflect.ValueOf(s)
	for i := 0; i < rs.Len(); i++ {
		rs.Index(i).SetInt(rs.Index(i).Int() * 2)
	}
	fmt.Println(s)
	//map
	/*m := map[string]string{"key1": "v1", "key2": "v2"}
	rm := reflect.ValueOf(&m).Elem()
	//fmt.Println(rm.Elem())
	for _, k := range rm.MapKeys() {
		rm.MapIndex(k).SetString(rm.MapIndex(k).String() + "_new")
	}
	fmt.Println(m)*/
	//结构体
	p := Person{ID: 1, Name: "张三", Address: []string{"北京", "上海", "郑州"}}
	rp := reflect.Indirect(reflect.ValueOf(&p)) //reflect.ValueOf(&p).Elem()
	fmt.Println(rp.Type(), rp)
	for i := 0; i < rp.NumField(); i++ {
		field := rp.Field(i)
		switch field.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			field.SetInt(field.Int() * 2)
		case reflect.String:
			field.SetString(field.String() + "_new")
		case reflect.Array, reflect.Slice:
			for j := 0; j < field.Len(); j++ {
				field.Index(j).SetString(field.Index(j).String() + "_new")
			}
		}
	}
	fmt.Println(p)
}

// 性能测试
func performTest() {
	//一千万次
	n := 10000000
	start := time.Now()
	var p Person
	v := reflect.ValueOf(&p).Elem()
	for i := 0; i < n; i++ {
		//v.FieldByName("Name").SetString("李四")
		//v.FieldByName("ID").SetInt(123)
		//v.FieldByName("Address").Set(reflect.ValueOf([]string{"北京", "上海", "郑州"}))
		for j := 0; j < v.NumField(); j++ {
			field := v.Field(j)
			switch field.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				field.SetInt(123)
			case reflect.String:
				field.SetString("呵呵哒")
			case reflect.Array, reflect.Slice:
				field.Set(reflect.ValueOf([]string{"北京", "上海", "郑州"}))
			}
		}
	}
	fmt.Println(p)
	elapsed := time.Since(start).Seconds() * 1000
	fmt.Printf("反射类型耗时：%v 毫秒\n", elapsed)
	start = time.Now()
	for i := 0; i < n; i++ {
		p.ID = 123
		p.Name = "李四"
		p.Address = []string{"北京", "上海", "郑州"}
	}
	elapsed = time.Since(start).Seconds() * 1000
	fmt.Printf("正常赋值耗时：%v 毫秒\n", elapsed)
}

// 接口赋值
func interfaceSetTest() {
	list := []int{5, 4, 3, 2, 1}
	fmt.Printf("%p %v\n", list, list)
	interfaceSet(&list)
	fmt.Printf("%p %v\n", list, list)
}

func interfaceSet(list interface{}) {
	data := []interface{}{1, 2, 3, 4, 5}
	ind := reflect.Indirect(reflect.ValueOf(list))
	ind.Set(reflect.ValueOf(data))
}
