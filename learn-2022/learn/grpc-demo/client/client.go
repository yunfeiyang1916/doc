package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc/codes"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/yunfeiyang1916/doc/learn-2022/learn/grpc-demo/share"

	"google.golang.org/grpc/metadata"

	"github.com/yunfeiyang1916/doc/learn-2022/learn/grpc-demo/pb"
)

// 自定义rpc的认证
type customRPCCredentials struct{}

func (c *customRPCCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	log.Println("uri:", uri)
	return map[string]string{"app_id": share.AppId, "app_key": share.AppKey}, nil
}
func (c *customRPCCredentials) RequireTransportSecurity() bool {
	return false
}

func main() {
	// 客户端的一元拦截器
	interceptor := grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		defer func(start time.Time) {
			log.Printf("客户端调用执行耗时：%s\n", time.Since(start))
		}(time.Now())
		// 利用metadata传输数据进行认证，或者直接使用customRPCCredentials
		//md := metadata.New(map[string]string{"app_id": share.AppId, "app_key": share.AppKey})
		//ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	})

	conn, err := grpc.Dial("127.0.0.1:8000", grpc.WithInsecure(), interceptor, grpc.WithPerRPCCredentials(&customRPCCredentials{}), grpc.WithUnaryInterceptor(retry.UnaryClientInterceptor()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ctx := context.Background()
	md := metadata.Pairs("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	ctx = metadata.NewOutgoingContext(ctx, md)

	// 增加超时控制
	//ctx, _ = context.WithTimeout(ctx, 1*time.Second)
	client := pb.NewGreeterClient(conn)
	// 设置重试
	resp, err := client.SayHello(ctx, &pb.HelloReq{Name: "张三"}, retry.WithMax(3), retry.WithPerRetryTimeout(1*time.Second), retry.WithCodes(codes.Unknown, codes.DeadlineExceeded, codes.Unavailable))
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Println("status.FromError error,err=", err)
		} else {
			log.Println(st.Code(), "  ", st.Message())
		}
	}
	log.Println("SayHello:", resp, "err:", err)

	person, err := client.SayPerson(ctx, &pb.Person{Id: 1000, Name: "张三"})
	log.Println("SayPerson:", person, "err:", err)

	// 单、双向流测试
	//streamTest(client)

}

func streamTest(client pb.GreeterClient) {
	// 从服务端获取流
	res, err := client.GetStream(context.Background(), &pb.StreamReq{Data: "张三"})
	if err != nil {
		log.Fatalln(err)
	}
	for {
		// 这里类似于socket的接收
		r, err := res.Recv()
		if err != nil {
			log.Println(err)
			break
		}
		log.Println("GetStream:", r.GetData())
	}
	// 向服务端推送流
	c, err := client.PostStream(context.Background())
	if err != nil {
		log.Panic(err)
	}
	// 发送十次
	for i := 1; i <= 10; i++ {
		if err := c.Send(&pb.StreamReq{Data: fmt.Sprintf("%v,%d", time.Now().Unix(), i)}); err != nil {
			log.Panic(err)
		}
		time.Sleep(time.Second)
	}

	// 双向流
	rr, err := client.Stream(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(2)
	// 接收服务端数据流
	go func() {
		defer wg.Done()
		for {
			r, err := rr.Recv()
			if err != nil {
				log.Println("接收错误：", err)
				break
			}
			fmt.Println("接收到服务端的消息：", r.Data)
		}
	}()
	// 向客服务端发送数据流
	go func() {
		defer wg.Done()
		i := 0
		for {
			i++
			if err := rr.Send(&pb.StreamReq{Data: "我是客户端" + strconv.Itoa(i)}); err != nil {
				log.Println("发送错误：", err)
				break
			}
			time.Sleep(time.Second)
		}
	}()
	wg.Wait()
}
