// 服务端
package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	tcpServer()
	//udpServer()
}

// tcp服务端
func tcpServer() {
	listener, err := net.Listen("tcp", ":8883")
	if err != nil {
		log.Fatalf("net.Listen error,err=%s \n", err)
	}
	defer listener.Close()
	fmt.Println("listen ok")
	var i int
	//for {
	conn, err := listener.Accept()
	if err != nil {
		fmt.Fprintf(os.Stderr, "listen.Accept error,err=%s", err)
		return
	}
	i++
	fmt.Printf("接受第 %d 个新的连接\n", i)
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	fmt.Println(n, err)
	fmt.Println("读取到内容：", string(buf[:n]))

	n, err = conn.Read(buf)
	fmt.Println(n, err)
	fmt.Println("读取到内容：", string(buf[:n]))
	//}

	fmt.Println("结束")
}
