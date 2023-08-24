package main

import "fmt"

func main() {
	i := 1
	arr := []int{1, 2, 3, 5}
	fmt.Printf("%p \n%p", &i, arr)
}
