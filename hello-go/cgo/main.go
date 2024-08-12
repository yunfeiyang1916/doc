package main

/*
int sum(int a,int b){
	return a+b;
}
*/
import "C"
import "fmt"

// 定义的c函数要放到注释里，并且放到import之前，并且c函数与import之间不能有空格
func main() {
	fmt.Println(C.sum(1, 1))
}
