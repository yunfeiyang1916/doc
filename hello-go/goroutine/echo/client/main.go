package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	//client1()
	clientUseChan()
}

// 使用通道同步
func clientUseChan() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatalln(err)
	}
	done := make(chan struct{})
	go func() {
		if _, err := io.Copy(os.Stdout, conn); err != nil {
			log.Println(err)
			log.Println("done")
			//发消息出去表示结束了
			done <- struct{}{}
		}
	}()
	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done
}

func client1() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	//开协程，将网络链接上的内容输出
	go mustCopy(os.Stdout, conn)
	//将标准输入的内容copy到网络链接上
	mustCopy(conn, os.Stdin)
}

// copy
func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatalln(err)
	}
}
