// 优化后的tcp服务端
package main

import (
	"fmt"
	"net"
)

const maxRead = 25

func server2() {
	// flag.Parse()
	// if flag.NArg() != 2 {
	// 	panic("usage:host port")
	// }
	host := "localhost"
	port := "5000"
	//hostAndPort := fmt.Sprintf("%s:%s", flag.Arg(0), flag.Arg(1))
	hostAndPort := fmt.Sprintf("%s:%s", host, port)
	listener := initServer(hostAndPort)
	for {
		conn, err := listener.Accept()
		checkError(err, "Accept: ")
		go connectionHandler(conn)
	}
}

// 初始化服务
func initServer(hostAndPort string) *net.TCPListener {
	serverAddr, err := net.ResolveTCPAddr("tcp", hostAndPort)
	checkError(err, "Resolving address:port failed: '"+hostAndPort+"'")
	listener, err := net.ListenTCP("tcp", serverAddr)
	checkError(err, "ListenTCP: ")
	println("Listening to: ", listener.Addr().String())
	return listener
}

// 连接处理程序
func connectionHandler(conn net.Conn) {
	connFrom := conn.RemoteAddr().String()
	println("Connection from: ", connFrom)
	sayHello(conn)
	for {
		var ibuf []byte = make([]byte, maxRead+1)
		length, err := conn.Read(ibuf[0:maxRead])
		ibuf[maxRead] = 0 // to prevent overflow
		switch err {
		case nil:
			handleMsg(length, err, ibuf)
			// case os.EAGAIN: // try again
			// 	continue
		default:
			goto DISCONNECT
		}
	}
DISCONNECT:
	err := conn.Close()
	println("Closed connection: ", connFrom)
	checkError(err, "Close: ")
}
func sayHello(to net.Conn) {
	obuf := []byte{'L', 'e', 't', '\'', 's', ' ', 'G', 'O', '!', '\n'}
	wrote, err := to.Write(obuf)
	checkError(err, "Write:wrote "+string(wrote)+" bytes.")
}

// 处理消息
func handleMsg(length int, err error, msg []byte) {
	if length > 0 {
		print("<", length, ":")
		for i := 0; ; i++ {
			if msg[i] == 0 {
				break
			}
			fmt.Printf("%c", msg[i])
		}
		print(">")
	}
}

// 检查错误
func checkError(err error, info string) {
	if err != nil {
		panic("Error: " + info + " " + err.Error())
	}
}
