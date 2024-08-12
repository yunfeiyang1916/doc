package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type HelloService struct {
}

func (h *HelloService) Hello(req string, reply *string) error {
	*reply = "Hello " + req
	return nil
}

// go rpc的测试
func main() {
	//tcpServer()
	httpServer()
}

func tcpServer() {
	// 1.实例化一个server
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Panic(err)
	}
	// 2.注册处理逻辑handler
	if err = rpc.RegisterName("HelloService", &HelloService{}); err != nil {
		log.Panic(err)
	}
	// 3.启动服务
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		// 这里使用的是gob的编码协议
		//rpc.ServeConn(conn)
		// 替换为json的编码
		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

// 使用http协议
func httpServer() {
	if err := rpc.RegisterName("HelloService", &HelloService{}); err != nil {
		log.Panic(err)
	}
	http.HandleFunc("/jsonrpc", func(w http.ResponseWriter, r *http.Request) {
		var conn io.ReadWriteCloser = struct {
			io.Writer
			io.ReadCloser
		}{
			Writer:     w,
			ReadCloser: r.Body,
		}
		rpc.ServeCodec(jsonrpc.NewServerCodec(conn))

	})
	http.ListenAndServe(":8000", nil)
}
