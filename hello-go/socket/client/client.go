// 客户端
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	tcpClient()
	//udpClient()
}

// tcp客户端
func tcpClient() {
	fmt.Println("startTime=", time.Now())
	for i := 1; i < 2; i++ {
		conn, err := net.Dial("tcp", "127.0.0.1:8883")
		if err != nil {
			log.Printf("%d net.Dial error,err=%s \n", i, err)
			continue
		}
		log.Println(i, ":connect to server ok")
		fmt.Println(conn.Write([]byte("你好")))
		//conn.Close()
	}
}

// udp客户端
func udpClient() {
	service := "127.0.0.1:7777"
	updAddress, err := net.ResolveUDPAddr("udp", service)
	checkError(err)
	conn, err := net.DialUDP("udp", nil, updAddress)
	checkError(err)
	defer conn.Close()
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Printf("请输入客户端名称：")
	clientName, _ := inputReader.ReadString('\n')
	//去除\r\n
	trimmedClient := strings.Trim(clientName, "\r\n")
	conn.Write([]byte(trimmedClient))
	for {
		input, _ := inputReader.ReadString('\n')
		//去除\r\n
		trimmedInput := strings.Trim(input, "\r\n")
		conn.Write([]byte(fmt.Sprintf("%s 发送消息：%s", trimmedClient, trimmedInput)))
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "致命错误：%s", err.Error())
		os.Exit(1)
	}
}
