// 客户端程序
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	//打开连接
	conn, err := net.Dial("tcp", "localhost:5000")
	if err != nil {
		//由于目标计算机积极拒绝而无法创建连接
		fmt.Println("Error dialing", err.Error())
		return
	}
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Println("First,what is your name?")
	clientName, _ := inputReader.ReadString('\n')
	//windows平台下用\r\n
	trimmedClient := strings.Trim(clientName, "\r\n")
	//给服务器发送信息直到程序退出
	for {
		fmt.Println("What to send to the server?Type Q to quit.")
		input, _ := inputReader.ReadString('\n')
		trimmedInput := strings.Trim(input, "\r\n")
		if trimmedInput == "Q" {
			return
		}
		_, err = conn.Write([]byte(trimmedClient + " says: " + trimmedInput))
	}
}
