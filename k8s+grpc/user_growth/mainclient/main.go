package main

import (
	"context"
	"flag"
	"log"
	"sync"
	"time"
	"user_growth/pb"

	"google.golang.org/grpc/keepalive"

	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc"
)

// 连接池
var connPool = sync.Pool{
	New: func() any {
		// 连接到服务
		addr := flag.String("addr", "localhost:80", "the address to connect to")
		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithWriteBufferSize(1024 * 1024 * 1), // 默认32KB
			grpc.WithReadBufferSize(1024 * 1024 * 1),  // 默认32KB,
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                10 * time.Minute,
				Timeout:             10 * time.Second,
				PermitWithoutStream: false,
			}),
		}
		conn, err := grpc.Dial(*addr, opts...)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		return conn
	},
}

func GetConn() *grpc.ClientConn {
	return connPool.Get().(*grpc.ClientConn)
}
func CloseConn(conn *grpc.ClientConn) {
	connPool.Put(conn)
}

func main() {
	conn := GetConn()
	if conn != nil {
		// 将连接放回连接池
		defer CloseConn(conn)
	} else {
		log.Fatalf("conneection nil")
	}
	// 设置1秒超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// 新建客户端
	cCoin := pb.NewUserCoinClient(conn)
	cGrade := pb.NewUserGradeClient(conn)
	// 测试1
	r1, err1 := cCoin.ListTasks(ctx, &pb.ListTasksRequest{})
	if err1 != nil {
		log.Printf("cCoin.ListTasks error=%v\n", err1)
	} else {
		log.Printf("cCoin.ListTasks:%+v\n", r1.GetDatalist())
	}
	// 测试2
	r2, err2 := cGrade.ListGrades(ctx, &pb.ListGradesRequest{})
	if err2 != nil {
		log.Printf("cGrade.ListGrades error=%v\n", err2)
	} else {
		log.Printf("cGrade.ListGrades:%+v\n", r2.GetDatalist())
	}

	// 测试3
	r3, err3 := cCoin.UserCoinInfo(ctx, &pb.UserCoinInfoRequest{Uid: 1})
	if err3 != nil {
		log.Printf("cCoin.UserCoinInfo error=%v", err2)
	} else {
		log.Printf("cCoin.UserCoinInfo:%+v", r3.Data)
	}
}
