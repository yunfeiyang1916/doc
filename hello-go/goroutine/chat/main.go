// 聊天服务
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// 只发送消息的通道类型
type client chan<- string

var (
	//进入的客户端
	entering = make(chan client)
	//离开的客户端
	leaving = make(chan client)
	//发送的消息
	messages = make(chan string)
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatalln(err)
	}
	//开启广播
	go boradcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn(conn)
	}
}

// 广播
func boradcaster() {
	//所有的客户端连接
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-messages:
			//将接收到的消息发送给所有客户端
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering: //接收新客户端的连接
			clients[cli] = true
		case cli := <-leaving: //接收到有客户端断开连接
			delete(clients, cli)
			//关闭该客户端的通道
			close(cli)
		}
	}
}

// 处理连接
func handleConn(conn net.Conn) {
	ch := make(chan string)
	//处理发送给客户端的消息
	go clientWriter(conn, ch)
	//客户端信息
	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + " has arrived"
	//发送客户端新连接进入消息
	entering <- ch
	//从连接读取
	input := bufio.NewScanner(conn)
	//忽略连接错误
	for input.Scan() {
		messages <- who + ":" + input.Text()
	}
	//发送客户端断开连接消息
	leaving <- ch
	messages <- who + ":" + "has left"
	//关闭连接
	conn.Close()
}

// 向客户端写消息
func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}
