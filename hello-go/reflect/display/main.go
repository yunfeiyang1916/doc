// 打印聚合数据类型的显示。
// 构建一个用于调式用的Display函数，给定一个聚合类型x，打印这个值对应的完整的结构，同时记录每个发现的每个元素的路径
package main

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func main() {
	//formatTest()
	displayTest()
}
func displayTest() {
	type M struct {
		ID       int
		Name     string
		Age      int
		Sex      bool
		Children []M
		Class    *string
		Dic      map[int]string
		I        interface{}
	}
	str := "呵呵哒"
	m := M{ID: 1, Name: "张三", Age: 50, Sex: true, Class: &str, Dic: map[int]string{1: "李四", 2: "王五"}, I: str}
	Display("m", m)
	//Display("f", os.Stdout)
}

func Display(name string, v interface{}) {
	fmt.Printf("Display %s (%T):\n", name, v)
	display(name, reflect.ValueOf(v))
}

// 打印聚合类型的元素路径及完整结构
func display(path string, v reflect.Value) {
	switch v.Kind() {
	case reflect.Invalid:
		fmt.Printf("%s =无效的类型\n", path)
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			display(fmt.Sprintf("%s[%d]", path, i), v.Index(i))
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fieldPath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name)
			display(fieldPath, v.Field(i))
		}
	case reflect.Map:
		for _, k := range v.MapKeys() {
			display(fmt.Sprintf("%s[%s]", path, formatAtom(k)), v.MapIndex(k))
		}
	case reflect.Ptr:
		if v.IsNil() {
			fmt.Printf("%s=nil\n", path)
		} else {
			display(fmt.Sprintf("(*s)", path), v.Elem())
		}
	case reflect.Interface:
		if v.IsNil() {
			fmt.Printf("%s=nil\n", path)
		} else {
			fmt.Printf("%s.type=%s\n", path, v.Elem().Type())
			display(path+".value", v.Elem())
		}
	default: // basic types, channels, funcs
		fmt.Printf("%s=%s\n", path, formatAtom(v))
	}
}

func formatTest() {
	var x int64 = 64
	var d time.Duration = 1 * time.Nanosecond
	Format(x)
	Format(d)
	Format(&x)
	Format([]int64{x})
	Format([...]time.Duration{d})
	Format(&[]time.Duration{d})
}

// 打印
func Format(v interface{}) {
	fmt.Println(formatAtom(reflect.ValueOf(v)))
}

// 格式原子在不检查其内部结构的情况下格式化一个值
func formatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "无效的类型"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Slice, reflect.Map:
		return v.Type().String() + " 0x" + strconv.FormatUint(uint64(v.Pointer()), 16)
	default: // reflect.Array, reflect.Struct, reflect.Interface
		return v.Type().String() + " value"
	}
}
