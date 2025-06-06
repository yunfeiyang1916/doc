package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"shop/shop-srv/user-srv/global"
	"shop/shop-srv/user-srv/handler"
	"shop/shop-srv/user-srv/initialize"
	"shop/shop-srv/user-srv/proto"
	"shop/shop-srv/user-srv/utils"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	port := flag.Int("port", 0, "端口号")
	flag.Parse()
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	if *port == 0 {
		*port, _ = utils.GetFreePort()
	}
	global.ServerConfig.Port = *port
	zap.S().Info(global.ServerConfig)

	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserService{})
	// 注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	// 注册到consul
	//consulCfg := api.DefaultConfig()
	//consulCfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	//consulClient, err := api.NewClient(consulCfg)
	//if err != nil {
	//	panic("failed to new consul client:" + err.Error())
	//}
	//// 生成consul需要的健康检查对象
	//check := &api.AgentServiceCheck{
	//	GRPC:                           fmt.Sprintf("host.docker.internal:%d", global.ServerConfig.Port),
	//	Timeout:                        "5s",
	//	Interval:                       "5s",
	//	DeregisterCriticalServiceAfter: "15s",
	//}
	//serviceID := uuid.NewV4().String()
	//registration := &api.AgentServiceRegistration{
	//	Name:    global.ServerConfig.Name,
	//	ID:      serviceID,
	//	Address: "localhost",
	//	Port:    global.ServerConfig.Port,
	//	Tags:    global.ServerConfig.Tags,
	//	Check:   check,
	//}
	//if err = consulClient.Agent().ServiceRegister(registration); err != nil {
	//	panic("failed to consul ServiceRegister:" + err.Error())
	//}
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", global.ServerConfig.Host, global.ServerConfig.Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()
	// 接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	// 从consul注销服务
	//if err = consulClient.Agent().ServiceDeregister(serviceID); err != nil {
	//	zap.S().Info("注销失败：", err.Error())
	//} else {
	//	zap.S().Info("注销成功")
	//}
}
