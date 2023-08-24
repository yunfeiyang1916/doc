// 闭包测试
package main

import (
	"fmt"
)

func main() {
	/*f1 := f(1)
	f2 := f(0)
	fmt.Printf("%p %T\n", f, f)
	fmt.Printf("%p %T %d\n", f1, f1, f1())
	fmt.Printf("%p %T %d\n", f2, f2, f2())
	fmt.Printf("%p %p\n", f1, f2)*/
	fs := [4]func(){}
	for i := 0; i < 4; i++ {
		fs[i] = func() {
			j := i
			fmt.Printf("fs[%d]()=%d\n", j, j)
		}
	}
	for _, f := range fs {
		f()
		fmt.Printf("%p\n", f)
	}
}

func f(i int) func() int {
	return func() int {
		i++
		fmt.Printf("&i=%p\n", &i)
		return i
	}
}
