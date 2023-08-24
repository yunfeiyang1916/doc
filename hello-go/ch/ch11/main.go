// 接口与反射
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strconv"
)

// 形状接口
type Shaper interface {
	//计算面积
	Area() float32
}

// 正方形
type Square struct {
	//边
	side float32
}

// 矩形
type Rect struct {
	width  float32
	height float32
}

// 实现接口方法
func (this *Square) Area() float32 {
	return this.side * this.side
}
func (this *Rect) Area() float32 {
	return this.width * this.height
}

// 可以以接口作为参数，接口值是一个多字数据结构，实际上是占据两个字长，一个用来存储它包含的类型，另一个用来存储它包含的数据或者指向数据的指针
// 它的值是nil，本质上是一个指针，虽然完全不是一回事。指向接口值的指针是非法的
func printArea(sh Shaper) {
	fmt.Printf("address=%p,sh type is %T,sh.Area()=%f\n", sh, sh, sh.Area())
}

// 形状接口测试
func ShaperTest() {
	//值类型
	//sqValue := Square{side: 20}
	//不能给接口赋值值类型
	//sh = sqValue
	//可以定义一个接口类型
	var sh Shaper
	//还可以这样定义一个接口类型
	var sh2 Shaper = Shaper(&Square{side: 20})
	sq := new(Square)
	sq.side = 10
	//可以将实例赋值给接口
	sh = sq
	fmt.Printf("sh Type is %T\n", sh)
	fmt.Printf("sh.Area()=%f\n", sh.Area())
	fmt.Printf("sh2 Type is %T\n", sh2)
	fmt.Printf("sh2.Area()=%f\n", sh2.Area())
	//矩形
	r := &Rect{width: 10, height: 20}
	//如果未实现接口方法，编译时会报错
	sh = r
	fmt.Printf("sh Type is %T\n", sh)
	fmt.Printf("sh.Area()=%f\n", sh.Area())
	//接口切片
	shSlicp := []Shaper{sq, r}
	fmt.Println(shSlicp)
	for _, v := range shSlicp {
		//v = new(Square)
		printArea(v)
	}
	fmt.Println(shSlicp)
}

// 标准接口
func standardTest() {
	//读接口
	var r io.Reader
	r = os.Stdin
	r = bufio.NewReader(r)
	r = new(bytes.Buffer)
	f, _ := os.Open("text.txt")
	r = bufio.NewReader(f)
}

// 接口嵌套接口
type ReadWrite interface {
	Read() string
	Write(str string) bool
}
type Lock interface {
	Lock()
	UnLock()
}
type File interface {
	ReadWrite
	Lock
	Print()
}
type FileStrcut struct {
	str string
}

func (this *FileStrcut) Read() string {
	return this.str
}
func (this *FileStrcut) Write(str string) bool {
	this.str = str
	return true
}
func (this *FileStrcut) Lock() {}
func (this *FileStrcut) UnLock() {

}
func (this *FileStrcut) Print() {
	fmt.Printf("Address=%p, type=%T,str=%s\n", this, this, this.Read())
}

// 内嵌接口测试
func embeddedTest() {
	//接口的值可以赋值给另一个接口变量，只要底层类型实现了必要的方法，这个转换是在允许时进行检查的，转换失败会导致一个运行时错误
	var f File
	var rw ReadWrite
	fs := &FileStrcut{}

	f = fs
	rw = fs
	rw.Write("张三")
	f.Print()
}

// 类型断言测试
func typeTest() {
	var sh Shaper
	sq := &Square{side: 10}
	sh = sq
	if t, ok := sh.(*Square); ok {
		printArea(t)
	}
	r := &Rect{}
	r.height = 10
	r.width = 20
	sh = r
	if t, ok := sh.(*Rect); ok {
		printArea(t)
	}
	if t, ok := sh.(Shaper); ok {
		fmt.Println("sh实现了Shaper接口：", t)
	}
	sh = nil
	//这种类型断言智能用于switch
	//t := sh.(type)
	//fmt.Println(t)
	//可以使用type-switch
	switch t := sh.(type) {
	case *Square:
		printArea(t)
	case *Rect:
		printArea(t)
	case nil:
		fmt.Println("t是nil")
	default:
		fmt.Println("未知类型")
	}
	//类型分类函数
	classifier(13, -14.3, "你好", complex(1, 2), nil, false)

}

// 类型断言测试2
func typeTest2() {
	//var list []interface{}
	list := make([]interface{}, 0)
	fmt.Println(list == nil)
	list = append(list, Rect{})
	list = append(list, Square{})
	//list = nil
	for _, v := range list {
		if t, ok := v.(Rect); ok {
			fmt.Printf("T=%T\n", t)
		} else if t, ok := v.(Square); ok {
			fmt.Printf("T=%T\n", t)
		}
	}

	for _, v := range list {
		//这种类型断言智能用于switch
		switch t := v.(type) {
		case Rect:
			fmt.Printf("T=%T\n", t)
		case Square:
			fmt.Printf("T=%T\n", t)
		}
	}
}

// 类型分类函数
func classifier(items ...interface{}) {
	for i, x := range items {
		switch t := x.(type) {
		case bool:
			fmt.Printf("type=%T Param #%d is a bool\n", t, i)
		case float64:
			fmt.Printf("type=%T Param #%d is a float64\n", t, i)
		case int, int64:
			fmt.Printf("type=%T Param #%d is a int\n", t, i)
		case nil:
			fmt.Printf("type=%T Param #%d is a nil\n", t, i)
		case string:
			fmt.Printf("type=%T Param #%d is a nil\n", t, i)
		default:
			fmt.Printf("type=%T Param #%d is a unknown\n", t, i)
		}
	}
}

// 使用方法集与接口
type List []int

// 定义在值上的方法，实现接口Lener
func (l List) Len() int {
	return len(l)
}

// 定义在指针上的方法，实现接口Appender
func (l *List) Append(val int) {
	*l = append(*l, val)
}

// 追加器接口
type Appender interface {
	Append(int)
}

func CountInto(a Appender, start, end int) {
	fmt.Printf("Appender add=%p type=%T\n", a, a)
	for i := start; i < end; i++ {
		a.Append(i)
	}
}

type Lener interface {
	Len() int
}

func LongEnough(l Lener) bool {
	fmt.Printf("Lener add=%p type=%T\n", l, l)
	return l.Len()*10 > 42
}

// 方法集测试
func methodTest() {
	l := List{1, 2, 3, 4, 5, 6}
	fmt.Printf("l add=%p type=%T\n", l, l)
	lp := &l
	fmt.Printf("lp add=%p type=%T\n", lp, lp)
	//因为实现接口方法的接收器是一个指针类型，所以需要传指针
	CountInto(lp, 7, 9)
	fmt.Printf("l add=%p type=%T value=%v\n", l, l, l)
	fmt.Printf("lp add=%p type=%T\n", lp, lp)
	//因为实现接口方法的接收器是一个值类型，所以需要传值，也可以传指针，指针会被自动解引用
	le := LongEnough(l)
	le = LongEnough(lp)
	fmt.Println(le)
}

// 使用Sorter接口排序
// 整型切片
type IntArray []int

// 长度计算次数
var LenCount int

// 长度
func (p IntArray) Len() int {
	LenCount++
	return len(p)
}

// 小于
func (p IntArray) Less(i, j int) bool { return p[i] < p[j] }

// 交换
func (p IntArray) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// 字符串切片
type StringArray []string

// 长度
func (p StringArray) Len() int { return len(p) }

// 小于
func (p StringArray) Less(i, j int) bool { return p[i] < p[j] }

// 交换
func (p StringArray) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// 个人类型
type Person struct {
	id   int
	name string
}

// 个人切片
type Persons []Person

func (p Persons) Len() int { return len(p) }

func (p Persons) Less(i, j int) bool { return p[i].id < p[j].id }

func (p Persons) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// 切片封装
type PersonArray struct {
	data []*Person
}

func (p *PersonArray) Len() int { return len(p.data) }

func (p *PersonArray) Less(i, j int) bool { return p.data[i].id < p.data[j].id }

func (p *PersonArray) Swap(i, j int) { p.data[i], p.data[j] = p.data[j], p.data[i] }

func (p *PersonArray) String() string {
	var str string
	for _, v := range p.data {
		str += fmt.Sprintf("id=%d name=%s\n", v.id, v.name)
	}
	return str
}

// 排序，实现sort.Interface排序接口
func Sort(data sort.Interface) {
	sum := 0
	//冒泡排序
	for i := 1; i < data.Len(); i++ {
		for j := 0; j < data.Len()-1; j++ {
			sum++
			if data.Less(j+1, j) {
				data.Swap(j, j+1)
			}
		}
	}
	//快速排序
	// for i := 0; i < data.Len(); i++ {
	// 	for j := i + 1; j < data.Len(); j++ {
	// 		sum++
	// 		if data.Less(j, i) {
	// 			data.Swap(i, j)
	// 		}
	// 	}
	// }
	fmt.Printf("排序循环%d次，调用获取长度%d次\n", sum, LenCount)
}
func SortTest() {

	intArray := IntArray{10, 1, 3, 4, 98, 29, 34, 90}
	fmt.Println(intArray)
	Sort(intArray)
	fmt.Println(intArray)
	strArray := StringArray{"张三", "nihao", "bhao", "嘿嘿", "傻帽"}
	Sort(strArray)
	fmt.Println(strArray)
	perArray := []Person{{id: 23, name: "张三"}, {id: 25, name: "李四"}, {id: 12, name: "王五"}, {id: 2, name: "燕小六"}, {id: 3, name: "鬼脚七"}}
	ps := Persons(perArray)
	Sort(ps)
	perArray = ps
	fmt.Println(perArray)
	pData := []*Person{&Person{id: 23, name: "张三"}, {id: 25, name: "李四"}, {id: 12, name: "王五"}, {id: 2, name: "燕小六"}, {id: 3, name: "鬼脚七"}}
	pa := PersonArray{data: pData}
	Sort(&pa)
	for _, v := range pa.data {
		fmt.Printf("type=%T add=%p id=%d name=%s\n", v, v, v.id, v.name)
	}
}

// 空接口，表示任何类型，类似java/c#中的object
type Any interface{}

// 矢量，包含空接口切片的容器,里面放的每个元素可以是不同类型的变量，要得到它们的原始类型（unboxing：拆箱）需要用到类型的断言
type Vector struct {
	data []Any
}

func (p *Vector) At(i int) Any { return p.data[i] }

// 设置值
func (p *Vector) Set(i int, val Any) { p.data[i] = val }

// 空接口测试
func emptyInterface() {
	//var anyPointer *Any
	//x := 123
	//指向接口值的指针是非法的
	//anyPointer = &x
	var any Any
	any = 123
	fmt.Printf("type=%T add=%p val=%d\n", any, &any, any)
	any = "呵呵哒"
	fmt.Printf("type=%T add=%p val=%s\n", any, any, any)
	p := Person{id: 1, name: "张三"}
	any = p
	fmt.Printf("type=%T add=%p val=%v\n", any, any, any)
	any = &p
	fmt.Printf("type=%T add=%p val=%v\n", any, any, any)
	any = 3456
	switch t := any.(type) {
	case int:
		fmt.Printf("type=%T add=%p val=%v\n", t, &t, t)
	case Person:
		fmt.Printf("type=%T add=%p val=%v\n", t, t, t)
	case *Person:
		fmt.Printf("type=%T add=%p val=%v\n", t, t, t)
	}
	TypeSwitch()
}

// type-switch联合lambda函数
func TypeSwitch() {
	testFunc := func(any interface{}) {
		switch t := any.(type) {
		case bool:
			fmt.Printf("type=%T val=%v\n", t, t)
		case int, int64:
			fmt.Printf("type=%T val=%v\n", t, t)
		case float32:
			fmt.Printf("type=%T val=%v\n", t, t)
		case string:
			fmt.Printf("type=%T val=%v\n", t, t)
		default:
			fmt.Printf("未知类型 type=%T val=%v\n", t, t)
		}
	}
	ve := Vector{data: []Any{1, "张三", "李四", 3.14, Person{id: 1, name: "张三"}}}
	for _, v := range ve.data {
		testFunc(v)
	}
}

// 通用类型的节点数据结构
type Node struct {
	le *Node
	//数据是空接口
	data interface{}
	ri   *Node
}

func (n *Node) SetData(data interface{}) {
	n.data = data
}

// 反射测试
func reflectTest() {
	var x float64 = 3.14
	t := reflect.TypeOf(x)
	fmt.Println(t)
	v := reflect.ValueOf(x)
	fmt.Println("value:", v)
	fmt.Println("type:", v.Type())
	//底层类型
	fmt.Println("kind:", v.Kind())
	fmt.Println("value:", v.Float())
	fmt.Println(v.Interface())
	//还原（接口）值
	fmt.Printf("value is %5.2e\n", v.Interface())
	y := v.Interface().(float64)
	fmt.Println(y)
	//使用反射赋值
	fmt.Println("是否可赋值：", v.CanSet())
	//v.SetFloat(3.14159) // Error: will panic: reflect.Value.SetFloat using unaddressable value
	//要想赋值需要反射时传入地址，并且使用Elem()函数
	v = reflect.ValueOf(&x)
	v = v.Elem()
	v.SetFloat(4.5)
	fmt.Println(v)
}

type NotknownType struct {
	s1, s2, s3 string
}

func (n NotknownType) String() string {
	return n.s1 + " " + n.s2 + " " + n.s3
}

type T struct {
	A int
	B string
}

func reflectTest2() {
	nts := NotknownType{"张三", "喜欢吃", "馒头"}
	v := reflect.ValueOf(nts)
	t := reflect.TypeOf(nts)
	t = v.Type()
	fmt.Println(t)
	kind := v.Kind()
	fmt.Println(kind)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fmt.Printf("Field %d:%v\n", i, field)
		//field.SetString("李四") //error:panic: reflect: reflect.Value.SetString using value obtained using unexported field
	}
	field := v.FieldByName("s3")
	fmt.Println(field)
	//result := v.Method(0).Call(nil)
	result := v.MethodByName("String").Call(nil)
	fmt.Println(result)
	v = reflect.ValueOf(&nts).Elem()
	//非导出字段不允许使用反射赋值
	//v.FieldByName("s1").SetString("你好啊")//panic: reflect: reflect.Value.SetString using value obtained using unexported field
	fmt.Println(nts)
	tt := T{24, "李四啊"}
	tv := reflect.ValueOf(&tt).Elem()
	tv.FieldByName("A").SetInt(77)
	fmt.Println(tt)
}

// String()的接口
type Stringer interface {
	String() string
}
type Celsius float64

func (this Celsius) String() string {
	return strconv.FormatFloat(float64(this), 'f', 1, 64) + " °C"
}

type Day int

var dayName = []string{"周一", "周二", "周三", "周四", "周五", "周六", "周日"}

func (this Day) String() string {
	return dayName[this]
}

// 使用反射实现print
func print(args ...interface{}) {
	for i, arg := range args {
		if i > 0 {
			os.Stdout.WriteString(" ")
		}
		switch t := arg.(type) {
		case Stringer:
			os.Stdout.WriteString(t.String())
		case int:
			os.Stdout.WriteString(strconv.Itoa(t))
		case string:
			os.Stdout.WriteString(t)
		default:
			os.Stdout.WriteString("???")
		}
	}
}

func printReflect() {
	print(Day(1), "was", Celsius(18.36))
}

func main() {
	//ShaperTest()
	//embeddedTest()
	//typeTest()
	typeTest2()
	//methodTest()
	//SortTest()
	//emptyInterface()
	//reflectTest()
	//reflectTest2()
	//printReflect()
}
