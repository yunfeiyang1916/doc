package main

import (
	"fmt"
	"net"
)

func main() {
	ipTest()
}

// 测试ip
func ipTest() {
	name := "127.0.0.1"
	addr := net.ParseIP(name)
	fmt.Println(addr)
}
