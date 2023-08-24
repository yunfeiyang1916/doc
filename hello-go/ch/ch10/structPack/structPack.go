// 结构体包
package structPack

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// 名称
var Name = "structPack"

// 点结构体
type Point struct {
	X int
	Y int
}

// 别名
type Po Point

// 矩形，里面嵌套了其他结构体，在内存中是连续的
type Rect struct {
	Width  Point
	Height Point
}

// 矩形2，里面是结构体指针，在内存中非连续
type Rect2 struct {
	Width  *Point
	Height *Point
}

// 人员
type Person struct {
	//第一个名
	FirstName string
	//第二个名
	LastName string
}

// 将人员名称转为大写
func UpPerson(person *Person) {
	person.FirstName = strings.ToUpper(person.FirstName)
	person.LastName = strings.ToUpper(person.LastName)
}

// 使用工厂方法创建结构体
type File struct {
	fd   int    //文件描述符
	name string //文件名
}

// 实例化一个文件对象
func NewFile(fd int, name string) *File {
	if fd < 0 {
		return nil
	}
	return &File{fd: fd, name: name}
}

// 强制使用工厂方法，也就是将一个结构体定义成私有的
// 矩阵
type matrix struct {
	X int
	Y int
}

// 实例化矩阵对象
func NewMatrix(x int, y int) *matrix {
	return &matrix{X: x, Y: y}
}

// map和struct vs new()和make()
// map[string]string别名
type Foo map[string]string

// 结构体
type Bar struct {
	X int
	Y int
}

// map和struct vs new()和make()
func NewAndMake() {
	//结构体正常编译
	s := new(Bar)
	s.X = 1
	s.Y = 2
	fmt.Println(s)
	//结构体不能编译
	//s2 := make(Bar)
	//map正常编译，但是只返回了一个nil指针
	m := new(Foo)
	fmt.Println(m)
	//map正常编译并且分配了一块内存，可以直接使用
	m2 := make(Foo)
	m2["x"] = "x"
	fmt.Println(m2)
}

// 带标签的结构体
type TagType struct {
	X int "This is tag"
	Y int `这是标签`
}

// 反射标签
func refTag(tag TagType, ix int) {
	tType := reflect.TypeOf(tag)
	f := tType.Field(ix)
	fmt.Println(f.Tag)
}

// 标签测试
func TagTest() {
	t := TagType{X: 10, Y: 20}
	for i := 0; i < 2; i++ {
		refTag(t, i)
	}
}

// 匿名字段和内嵌结构体，使用内嵌或组合来实现继承
// 内部结构体
type Inner struct {
	X int
	Y int
}

// 外部结构体
type Outer struct {
	Z     int
	int   //匿名字段，同类型的匿名字段只能出现一次
	Inner //匿名字段，内嵌结构体
}

// 匿名结构体测试
func AnonymousTest() {
	//使用结构字面量赋值
	o := Outer{Inner: Inner{X: 1, Y: 2}, Z: 3, int: 4}
	fmt.Println(o)
	//使用new()赋值
	op := new(Outer)
	op.X = 10
	op.Inner.X = op.X * 10
	op.Y = 20
	op.Z = 30
	op.int = 40
	fmt.Println(op)
}

// 命名冲突
type A struct {
	X int
}
type B struct {
	X, Y int
}
type C struct {
	A
	B
}
type D struct {
	B
	Y int
}

// 命名冲突
func NamingConflicts() {
	var c C
	//这样使用就会命名冲突
	//c.X = 5
	c.A.X = 5
	fmt.Println(c)
	var d D
	//外层名字会覆盖内层名字
	d.Y = 10
	d.B.Y = 120
	fmt.Println(d)
}

// 结构体方法
// point求和
func (recv *Point) Sum() int {
	return recv.X + recv.Y
}

// 允许不同接收类型方法名相同的重载（这还能叫重载么？）
func (this Rect) Sum() int {
	return this.Height.X + this.Width.Y
}
func (this *Point) SumAndParam(x int) int {
	return this.X + this.Y + x
}
func (this *Point) Sub() {
	this.X -= 10
	this.Y -= 20
}

// int别名
type Int int

// 别名
type PAlias Point

// 别名方法重载，类型和作用在它上面的方法必须在同一个包里，可以不是一个文件
func (this *PAlias) Sum() int {
	return this.X + this.Y
}

// 给整型加一个方法
func (this Int) ToString() string {
	return string(this)
}

// 时间别名，给非在本包中的类型增加方法
type MyTime struct {
	time.Time
}

func (t MyTime) first3Chars() string {
	return t.String()[0:3]
}

// 方法测试
func MethodTest() {
	p := &Point{X: 10, Y: 20}
	s := p.Sum()
	fmt.Printf("Point.Sum()=%d\n", s)
	s = p.SumAndParam(100)
	fmt.Printf("Point.SumAndParam(100)=%d\n", s)
	p2 := Point{X: 20, Y: 40}
	fmt.Println(p2.Sum())
	p.Sub()
	p2.Sub()
	fmt.Printf("p.sub()=")
	fmt.Println(p)
	fmt.Printf("p2.sub()=")
	fmt.Println(p2)
	p3 := PAlias{X: 20, Y: 40}
	fmt.Printf("PAlias.Sum()=%d\n", p3.Sum())
	t := MyTime{Time: time.Now()}
	fmt.Println(t.first3Chars())
}

// 使用方法访问或设置私有字段
type People struct {
	name string
}

// 读取name值
func (this *People) Name() string {
	return this.name
}

// 设置name
func (this *People) SetName(str string) {
	this.name = str
}

// 内嵌类型的方法会被继承，这样来实现方法继承
type Point2 struct {
	Point
	Name string
}

func MethodTest2() {
	p := new(People)
	p.SetName("张三")
	//p.name = "张三"
	fmt.Printf("p.name=%v\n", p.Name())
	p2 := new(Point2)
	p2.Sub()
	fmt.Println(p2)
}

// 接口
type Engine interface {
	Start()
	Stop()
}

type Car struct {
	Engine
}

func (c *Car) GoToWorkIn() {
	c.Start()
	c.Stop()
}

// 在类型中嵌入功能
// 记录日志功能
type Log struct {
	msg string
}

// 顾客，使用组合方式包含日志功能，包含一个所需功能类型的具名字段
type Customer struct {
	Name string
	//日志功能
	log *Log
}

// 内嵌（匿名地）所需功能类型
type Customer2 struct {
	Name string
	Log
}

// 增加日志
func (this *Log) Add(s string) {
	this.msg += "\n" + s
}

// 或者日志消息
func (this *Log) String() string {
	return this.msg
}

// 获取Customer上的Log对象
func (this *Customer) Log() *Log {
	return this.log
}

// 获取字符串方法
func (this *Customer2) String() string {
	return fmt.Sprintf("%s\nLog:%s", this.Name, fmt.Sprintln(this.Log))
}

// 组合方式实现功能测试
func CombFuncTest() {
	c := new(Customer)
	c.Name = "张三"
	c.log = new(Log)
	c.log.msg = "1-Yes we can!"
	c.log.Add("2-c.log.Add")
	c.Log().Add("3-c.Log().Add")
	fmt.Println(c)
	fmt.Println(c.Log())
	//内嵌所需功能
	c2 := &Customer2{Name: "李四", Log: Log{msg: "第一条日志"}}
	c2.Log.Add("第二条日志")
	c2.Add("第三条日志")
	fmt.Println(c2)
}

// 多重继承
type Camera struct {
}

// 提取一张照片
func (c *Camera) TakeAPicture() string {
	return "Click"
}

type Phone struct{}

// 打电话
func (p *Phone) Call() string {
	return "Ring Ring"
}

// 多继承
type CameraPhone struct {
	Camera
	Phone
}

// 基类
type Base struct{}

// Base的方法
func (Base) Magic() {
	fmt.Println("Base.Magic()")
}
func (this Base) MoreMagic() {
	this.Magic()
	this.Magic()
}

type Voodoo struct {
	Base
}

// Voodoo的方法
func (Voodoo) Magic() {
	fmt.Println("Voodoo.Magic()")
}

// 多继承测试
func MultipleInherit() {
	c := new(CameraPhone)
	fmt.Println("多继承测试")
	fmt.Printf("Camera.TakeAPicture()=%v\n", c.TakeAPicture())
	fmt.Printf("Camera.Call()=%v\n", c.Call())

	v := new(Voodoo)
	//会调用Voodoo.Magic()
	v.Magic()
	//会调用Base.MoreMagic(),内部调用Base.Magic()
	v.MoreMagic()
}

// String()格式化
type TwoInts struct {
	a int
	b int
}

func (this *TwoInts) String() string {
	return fmt.Sprintf("(%s/%s)", strconv.Itoa(this.a), strconv.Itoa(this.b))
}

// 字符串格式化
func StringFormat() {
	t := new(TwoInts)
	t.a = 10
	t.b = 20
	fmt.Printf("输出字符串：%s\n", t)
	fmt.Println("Println输出", t)
	fmt.Printf("输出类型%T\n", t)
	fmt.Printf("实例完整输出%#v\n", t)
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("占用内存%dkb\n", m.Alloc/1024)
}
