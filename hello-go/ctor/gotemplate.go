/*Go程序的一般结构
 * 1、先定义包名
 * 2、在导入依赖包
 * 3、开始对常量、变量和类型定义或声明
 * 4、如果存在init函数，则定义该函数（这是一个特殊的函数，每个含有该函数的包都会首先执行这个函数）
 * 5、如果当前包是main包，则定义main函数
 * 6、定义其余函数，首先是类型的方法，接着是按照main函数中先后调用的顺序来定义相关函数，如果函数很多，可以按照字母顺序来进行排序
 */
package ctor

import "fmt"

const c string = "我是常量"

var v string = "我是变量"

type T struct{}

// 初始化函数
func init() {
	fmt.Println("调用ctor的初始化函数")
}

func f1() {
	fmt.Println(c)
	fmt.Println(v)
}

func Run() {
	fmt.Println("开始运行Go程序的一般结构")
}
