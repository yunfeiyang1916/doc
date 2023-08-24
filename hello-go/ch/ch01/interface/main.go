// 接口测试
package main

import (
	"fmt"
	"os"
)

func main() {
	var err error
	fmt.Printf("error==nil? %v\n", err == nil)
	err = foo()
	fmt.Printf("foo()==nil? %v\n", err == nil)
	err = foo2()
	fmt.Printf("foo2()==nil? %v\n", err == nil)
}

func foo() error {
	var err *os.PathError
	fmt.Printf("pathError==nil? %v\n", err == nil)
	return err
}
func foo2() error {
	var err error
	return err
}
