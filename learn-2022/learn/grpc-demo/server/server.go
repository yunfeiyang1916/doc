package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/yunfeiyang1916/doc/learn-2022/learn/grpc-demo/share"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"

	"github.com/yunfeiyang1916/doc/learn-2022/learn/grpc-demo/pb"
)

type GreeterService struct {
	pb.UnimplementedGreeterServer
}

type Validator interface {
	Validate() error
}

// 服务端流模式，从服务端源源不断的获取
func (g *GreeterService) GetStream(req *pb.StreamReq, server pb.Greeter_GetStreamServer) error {
	i := 0
	for {
		i++
		if err := server.Send(&pb.StreamResp{Data: fmt.Sprintf("%v,%d", time.Now().Unix(), i)}); err != nil {
			return err
		}
		if i > 10 {
			break
		}
		time.Sleep(time.Second)
	}
	return nil
}

func (g *GreeterService) PostStream(server pb.Greeter_PostStreamServer) error {
	for {
		r, err := server.Recv()
		if err != nil {
			log.Println(err)
			break
		}
		log.Println(r.Data)
	}
	return nil
}

func (g *GreeterService) Stream(server pb.Greeter_StreamServer) error {
	wg := sync.WaitGroup{}
	wg.Add(2)
	// 接收客户端数据流
	go func() {
		defer wg.Done()
		for {
			r, err := server.Recv()
			if err != nil {
				log.Println("接收错误：", err)
				break
			}
			fmt.Println("接收到客户端的消息：", r.Data)
		}
	}()
	// 向客户端发送数据流
	go func() {
		defer wg.Done()
		i := 0
		for {
			i++
			if err := server.Send(&pb.StreamResp{Data: "我是服务器" + strconv.Itoa(i)}); err != nil {
				log.Println("发送错误：", err)
				break
			}
			time.Sleep(time.Second)
		}
	}()
	wg.Wait()
	return nil
}

func (g *GreeterService) SayHello(ctx context.Context, req *pb.HelloReq) (*pb.HelloResp, error) {
	log.Println("收到请求：SayHello. 参数：" + req.Name)
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		log.Println("收到metadata:", md)
	} else {
		log.Println("没有metadata数据")
	}
	time.Sleep(time.Second)
	return &pb.HelloResp{Reply: "Hello " + req.Name}, nil
}

// 带验证器
func (g *GreeterService) SayPerson(ctx context.Context, req *pb.Person) (*pb.Person, error) {
	return req, nil
}

func main() {
	// 一元拦截器，（除了一元拦截器，还有流拦截器）
	interceptor := grpc.UnaryInterceptor(func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func(start time.Time) {
			log.Printf("执行耗时：%s\n", time.Since(start))
		}(time.Now())
		// 先做验证
		if r, ok := req.(Validator); ok {
			if err := r.Validate(); err != nil {
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
		}

		// 登录认证
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return resp, status.Error(codes.Unauthenticated, "无token认证信息")
		}
		var (
			appId  string
			appKey string
		)
		if val, ok := md["app_id"]; ok && len(val) > 0 {
			appId = val[0]
		}
		if val, ok := md["app_key"]; ok && len(val) > 0 {
			appKey = val[0]
		}
		if appId != share.AppId || appKey != share.AppKey {
			return resp, status.Error(codes.Unauthenticated, "token认证失败")
		}
		return handler(ctx, req)
	})
	// 初始化grpc服务
	g := grpc.NewServer(interceptor)
	pb.RegisterGreeterServer(g, &GreeterService{})
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	if err = g.Serve(listener); err != nil {
		panic("failed to server:" + err.Error())
	}
}
