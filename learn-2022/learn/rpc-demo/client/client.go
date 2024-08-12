package main

import (
	"bytes"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func main() {
	//tcpClient()
	httpClient()
}

func tcpClient() {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Panic(err)
	}
	var reply string
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
	if err = client.Call("HelloService.Hello", "张三", &reply); err != nil {
		log.Panic(err)
	} else {
		log.Println(reply)
	}
}

func httpClient() {
	reader := bytes.NewReader([]byte(`{"id":0,"params":["张三"],"method":"HelloService.Hello"}`))
	resp, err := http.Post("http://127.0.0.1:8000/jsonrpc", "application/json", reader)
	log.Println("resp:", resp, "   err:", err)
}
