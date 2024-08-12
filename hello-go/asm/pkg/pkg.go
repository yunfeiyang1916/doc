package pkg

import "fmt"

var Id int

var Name string

var NameData [8]byte

// 因为汇编中不允许调用内建函数，可以将内建函数包装下，然后在汇编调用即可
func Print(a ...any) (n int, err error) {
	return fmt.Print(a...)
}

func Println(a ...any) (n int, err error) {
	return fmt.Println(a...)
}

var helloworld = "你好，世界"

func PrintHello()
