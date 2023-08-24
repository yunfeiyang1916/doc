package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// 存储单位类型别名
type ByteSize float64

const (
	_           = iota //通过赋值给空白标识符忽略0值
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

// 存储测试
func byteTest() {
	fmt.Printf("KB=%.0f MB=%.0f GB=%.0f YB=%0.f \n", KB, MB, GB, YB)
}

// 随机数测试
func randomTest() {
	for i := 0; i < 10; i++ {
		a := rand.Int()
		fmt.Printf("%d ", a)
	}
	fmt.Println()
	for i := 0; i < 5; i++ {
		r := rand.Intn(8)
		fmt.Printf("%d ", r)
	}
	fmt.Println()
	timens := int64(time.Now().Nanosecond())
	rand.Seed(timens)
	for i := 0; i < 10; i++ {
		fmt.Printf("%2.2f ", 100*rand.Float32())
	}
	fmt.Println()
}

// 字符测试
func charTest() {
	//Unicode编码，\u表示占用俩字节，\U表示4字节
	var ch int = '\u0041'
	var ch2 int = '\u03B2'
	var ch3 int = '\U00101234'
	fmt.Printf("整型显示    %v %v %v \n", ch, ch2, ch3)
	fmt.Printf("字符显示    %c %c %c \n", ch, ch2, ch3)
	fmt.Printf("16进制显示  %X %X %X \n", ch, ch2, ch3)
	fmt.Printf("Unicode码点 %U %U %U \n", ch, ch2, ch3)
	fmt.Printf("h是整型：%v h是字母：%v ' '是否为空白字符：%v", unicode.IsDigit('h'), unicode.IsLetter('h'), unicode.IsSpace(' '))
	fmt.Println()
	str := "这是字符串"
	stringFunc(str)
	fmt.Println(str)
}

func stringFunc(str string) {
	fmt.Printf("&p=%p p=%p t=%T\n", &str, str, str)
	str = "新的值"
	x := str[0]
	fmt.Printf("x=%c\n", x)
}

// 字符串统计
func stringCountTest() {
	str := "呵呵哒，猜猜我多长"
	fmt.Printf("字节数：%v 字符数：%v", countByte(str), countChar(str))
}

// 统计字节
func countByte(str string) int {
	return len(str)
}

// 统计字符数量
func countChar(str string) int {
	return utf8.RuneCountInString(str)
}

// 前后缀测试
func prefixTest() {
	str := "呵呵哒，要来测试前后缀了"
	fmt.Printf("\"呵呵\"是否是前缀：%v \n", strings.HasPrefix(str, "呵呵"))
	fmt.Printf("\"了\"是否是后缀：%v \n", strings.HasSuffix(str, "了"))
}

// 字符串索引位置测试
func strIndexTest() {
	str := "Hello World !"
	fmt.Printf("\"W\"出现的位置：%d\n", strings.Index(str, "W"))
	fmt.Printf("\"W\"最后出现的位置：%d\n", strings.LastIndex(str, "W"))
	str = "呵呵哒，啊呵呵哒"
	fmt.Printf("\"呵呵\"出现的位置：%d\n", strings.Index(str, "呵呵"))
	fmt.Printf("\"呵呵\"最后出现的位置：%d\n", strings.LastIndex(str, "呵呵"))
}

// 统计字符出现次数
func strCountTest() {
	str := "呵呵哒，这是测试的呵啊呵"
	fmt.Printf("\"呵\"出现次数：%d \n", strings.Count(str, "呵"))
}

// 重复字符串测试
func strRepeatTest() {
	str := "你好啊"
	str2 := strings.Repeat(str, 3)
	fmt.Printf("%v", str2)
}

// 字符串大小写测试
func strUppperTest() {
	str := "This is me!"
	str2 := strings.ToUpper(str)
	str3 := strings.ToLower(str2)
	fmt.Printf("大写：%s\n", str2)
	fmt.Printf("小写：%s\n", str3)
}

// 字符串切割测试
func strSplitTest() {
	str := "This is me!"
	//使用空白符切割
	ss := strings.Fields(str)
	//使用指定字符分隔
	ss2 := strings.Split(str, "i")
	fmt.Printf("%s \n", strings.Join(ss, ","))
	fmt.Printf("%s \n", strings.Join(ss2, ","))
}

// 字符串与其他类型互转
func strconvTest() {
	str := "666"
	fmt.Printf("当前操作系统下int长度：%d", strconv.IntSize)
	//字符串转成整型
	i, _ := strconv.Atoi(str)
	fmt.Printf("转成整型 %d\n", i)
	i += 5
	str = strconv.Itoa(i)
	fmt.Printf("转成字符串 %s\n", str)
}

// 时间测试
func timeTest() {
	t := time.Now()
	fmt.Printf("原始日期：%v\n", t)
	fmt.Printf("%4d-%02d-%02d %02d:%02d:%02d\n", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	t = time.Now().UTC()
	fmt.Printf("UTC:%v\n", t)
	//一周的纳秒表示
	var week time.Duration = 7 * 24 * 60 * 60 * 1e9
	weekFromNow := t.Add(week)
	fmt.Println(weekFromNow)
	//格式化
	fmt.Printf("time.RFC822格式化：%v\n", t.Format(time.RFC822))
	fmt.Printf("time.ANSIC格式化：%v\n", t.Format(time.ANSIC))
	fmt.Println(t.Format("02 Jan 2006 15:04"))
	s := t.Format("2018-01-01 17:37:27")
	fmt.Println(t, "=>", s)
}

// 指针测试
func pointerTest() {
	var i int = 5
	fmt.Printf("i=%d &i=%p\n", i, &i)
	var p *int = &i
	fmt.Printf("p=%p *p=%d &p=%p\n", p, *p, &p)
	s := "ni hao a!"
	sp := &s
	*sp = "hello"
	fmt.Printf("s=%s sp=%p *sp=%s\n", s, sp, *sp)
	var ip *int = nil
	ip = p
	*ip = 123
}
