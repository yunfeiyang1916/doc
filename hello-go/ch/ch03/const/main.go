// 常量
package main

import (
	"fmt"
)

// 编译器为没有明确的基础类型的数字常量提供比基础类型更高精度的算术运算，你可以认为至少有256位的运算精度
// 有六种未明确类型的常量类型，无类型的布尔型、无类型的整数、无类型的浮点数、无类型的复数、无类型的字符、无类型的字符串
// 无类型的常量可以直接用于任意相同类型（不同位数）的地方
const PI = 3.14159265358979323846264338327950288419716939937510582097494459

func main() {
	var f float32 = PI
	var f2 float64 = PI
	var c complex128 = PI
	fmt.Printf("%T %v\n", PI, PI)
	fmt.Printf("%T %v\n", f, f)
	fmt.Printf("%T %v\n", f2, f2)
	fmt.Printf("%T %v\n", c, c)
}
