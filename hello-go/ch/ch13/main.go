// 错误处理与测试
package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"project/hello/ch13/even"
	"strconv"
	"strings"
)

var errNotFound error = errors.New("Not found errors")

// 错误测试
func errorTest() {
	fmt.Printf("error:%v\n", errNotFound)
	//用fmt.Errorf创建错误对象
	err := fmt.Errorf("math:square root of negative number %g", -123.5)
	fmt.Printf("error:%v\n", err)
}

// 异常测试
func panicTest() {
	fmt.Println("开始运行")
	defer func() {
		log.Println("处理异常")
		if err := recover(); err != nil {
			log.Printf("run time panic:%v\n", err)
		}
	}()
	file, err := os.Open("123.txt")
	if err != nil {
		fmt.Println(err)
		panic(err)
		//return后面的defer不会在执行
		return
	}
	defer func() {
		fmt.Println("进入defer")
		file.Close()
	}()
	//panic("抛出异常")
	fmt.Println("运行结束")
}

// 自定义包中错误处理和panicking
// 解析错误
type ParseError struct {
	//将索引放入空格分隔的单词列表中
	Index int
	//产生解析错误的单词
	Word string
	//错误
	Err error
}

func (e *ParseError) String() string {
	return fmt.Sprintf("pkg parse:error parsing %q as int", e.Word)
}

// 将字符串转成整型数组
func Parse(input string) (numbers []int, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("pkg:%v", r)
			}
		}
	}()
	fields := strings.Fields(input)
	numbers = fields2numbers(fields)
	return
}

// 将字符串数组转成整型数组
func fields2numbers(fields []string) (numbers []int) {
	if len(fields) == 0 {
		panic("no words to parse")
	}
	for i, v := range fields {
		num, err := strconv.Atoi(v)
		if err != nil {
			panic(&ParseError{Index: i, Word: v, Err: err})
		}
		numbers = append(numbers, num)
	}
	return
}
func ParseTest() {
	str := []string{
		"1 2 3 4 5",
		"100 50 25 12.5 6.25",
		"2 + 2 = 4",
		"1st class",
		""}
	for _, v := range str {
		fmt.Printf("Parsing %q:\n ", v)
		nums, err := Parse(v)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(nums)
	}
}

func f() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	fmt.Println("Calling g.")
	g(0)
	fmt.Println("Returned normally from g.")
}
func g(i int) {
	if i > 3 {
		fmt.Println("Panicking!")
		panic(fmt.Sprintf("%v", i))
	}
	defer fmt.Println("Defer in g", i)
	fmt.Println("Printing in g", i)
	g(i + 1)
}

// 启动外部命令和程序
func startProcess() {
	env := os.Environ()
	procAttr := &os.ProcAttr{Env: env, Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}}
	pid, err := os.StartProcess("explorer", nil, procAttr)
	if err != nil {
		fmt.Printf("Error %v starting process!\n", err)
		//os.Exit(1)
	}
	fmt.Printf("The process id is %v\n", pid)
	//使用命令运行
	cmd := exec.Command("explorer.exe")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error %v executing command!", err)
		os.Exit(1)
	}
	fmt.Printf("The command is %v\n", cmd)
}

// 奇偶数测试
func evenTest() {
	for i := 0; i <= 100; i++ {
		fmt.Printf("%d is even %v?\n", i, even.Even(i))
	}
}

func main() {
	//errorTest()
	//panicTest()
	//ParseTest()
	//f()
	//startProcess()
	evenTest()
}
