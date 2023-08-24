// 浮点数
package main

import (
	"fmt"
	"math"
	"reflect"
)

func main() {
	//divZero()
	inexactFloat()
}

// 除0测试
func divZero() {
	var f float64
	//正无穷
	pInf := 2 / f
	//负无穷
	nInf := -2 / f
	fmt.Println(f, -f, 1/f, -1/f, f/f, f/2)
	fmt.Println(reflect.TypeOf(pInf), reflect.TypeOf(nInf))
	fmt.Println(1/f == pInf, -1/f == nInf, pInf == nInf)
	nan := math.NaN()
	fmt.Println(nan == nan, nan > nan, nan < nan)
}

// 不精准浮点数
func inexactFloat() {
	//float32的有效bit位只有23个，其他的bit位用于指数和符号；当整数打印23位能表达的范围时，float32的表示将出现误差
	//float32可以提供大约6个十进制数的精度，而float64则可以提供约15个十进制数的精度
	var f float32 = 1 << 24 //16777216
	f2 := f + 1
	fmt.Println(f == f2)
	fmt.Printf("%f\n", f2)
	fmt.Printf("%f\n", f+0.11)
	fmt.Printf("%f\n", float64(16777216+1))
	var f64 float64 = 1 << 56
	fmt.Printf("%f\n", f64)
	f3 := f64 + 1
	fmt.Printf("%f\n", f3)
	fmt.Println(f3 == f64)
}
