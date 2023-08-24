// 控制结构相关测试
package main

import (
	"fmt"
	"strconv"
)

// if else测试
func ifelseTest() {
	fmt.Printf("10,20中的最大值是：%d\n", isGreater(10, 20))
	first := 10
	if first <= 0 {
		fmt.Printf("first<=0\n")
	} else if first > 0 && first < 5 {
		fmt.Printf("first>0&&first<5\n")
	} else {
		fmt.Printf("first>=5\n")
	}
	if cond := 5; cond > 10 {
		fmt.Printf("cond>10\n")
	} else if cond := 19; cond > 5 && cond <= 10 {
		fmt.Printf("cond>5&&cond<=10\n")
	} else {
		str := fmt.Sprintf("cond=%d\n", cond)
		fmt.Printf(str)
	}
	//判断多值返回的
	i, err := strconv.Atoi("abc")
	if err != nil {
		fmt.Printf("abc转换成整型失败，失败原因：%v\n", err.Error())
		//os.Exit(-1)
	} else {
		fmt.Println(i)
	}
	if x := 23; x < 10 {
		fmt.Printf("x<10\n")
	} else if x > 20 {
		fmt.Printf("x>20\n")
	}
}

// 比较俩数大小
func isGreater(x int, y int) int {
	if x > y {
		return x
	}
	return y
}

// switch测试
func switchTest() {
	i := 0
	//不需要写break，switch会在匹配到case时自动跳出
	switch i {
	case 0, 1, 2:
		fmt.Println("0,1,2")
	case 3:
		fmt.Println(3)
	default:
		fmt.Println("default")
	}
	//可以使用fallthrough穿透
	switch i {
	case 0:
		fmt.Println(0)
		fallthrough
	case 1:
		fmt.Println(1)
		fallthrough
	default:
		fmt.Println("default")
	}
	//第二种形式，不提供任何被判断值，然后在每个case分支中进行条件测试
	switch {
	case i < 0:
		fmt.Println("i<0")
	case i == 0:
		fmt.Println("i==0")
	case i > 0:
		fmt.Println("i>0")
	default:
		fmt.Println("default")
	}
	//第三种形式，包含一个初始化语句，然后在每个case分支中进行条件测试
	switch j := 10; {
	case j < 10:
		fmt.Println("j<10")
	case j == 10:
		fmt.Println("j==10")

	}
}

// 循环测试
func forTest() {
	for i := 0; i < 5; i++ {
		fmt.Printf("%d ", i)
	}
	fmt.Println()
	str := "hehe 中文"
	for i := 0; i < len(str); i++ {
		fmt.Printf("%c ", str[i])
	}
	fmt.Println()
	for i := 1; i <= 20; i++ {
		switch {
		case i%3 == 0 && i%5 == 0:
			fmt.Printf("FizzBuzz ")
		case i%3 == 0:
			fmt.Printf("Fizz ")
		case i%5 == 0:
			fmt.Printf("Buzz ")
		default:
			fmt.Printf("%d ", i)
		}
	}
	fmt.Println()
	//for第二种形式，没有初始化语句和修饰语句，因此";;"是多余的了
	var i int
	for i >= 0 {
		fmt.Printf("%d ", i)
		i--
	}
	fmt.Println()
	//for第三种形式，无限循环
	for {
		if i > 10 {
			break
		}
		fmt.Printf("%d ", i)
		i++
	}
	fmt.Println()
	//for第四种形式，for-rang，类似c#的foreach，i为索引，val为值，并且val是值拷贝
	for i, s := range str {
		fmt.Printf("i=%d s=%c ", i, s)
	}
	fmt.Println()
}

// 标签测试
func labelTest() {
label1:
	for i := 0; i <= 5; i++ {
		for j := 0; j <= 5; j++ {
			if j == 4 {
				break label1
			}
			fmt.Printf("i=%d j=%d ", i, j)
		}
	}
	fmt.Println()
	//模拟循环
	i := 0
HERE:
	if i == 5 {
		return
	}
	fmt.Printf("%d ", i)
	i++
	goto HERE
	//正常用法一般是将标签放到goto后面
}

func main() {
	ifelseTest()
	//switchTest()
	//forTest()
	//labelTest()
}
