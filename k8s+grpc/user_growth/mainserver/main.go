package main

import (
	"log"
	"net"
	"user_growth/conf"
	"user_growth/dbhelper"
	"user_growth/pb"
	"user_growth/ugserver"

	"google.golang.org/grpc/reflection"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
)

func initDb() {
	//time.Local=time.UTC
	conf.LoadConfigs()
	dbhelper.InitDb()
}

func main() {
	initDb()
	// 监听端口
	lis, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatalf("failed to listen:%s", err.Error())
	}
	// 初始化服务
	s := grpc.NewServer()
	// 注册服务
	pb.RegisterUserCoinServer(s, &ugserver.UgCoinServer{})
	pb.RegisterUserGradeServer(s, &ugserver.UgGradeServer{})

	// 注册反射服务，可以使用grpcurl工具调用
	reflection.Register(s)
	// 启动服务
	log.Printf("server listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server:%v", err)
	}
}
