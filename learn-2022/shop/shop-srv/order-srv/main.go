package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"shop/shop-srv/order-srv/global"
	"shop/shop-srv/order-srv/handler"
	"shop/shop-srv/order-srv/initialize"
	"shop/shop-srv/order-srv/proto"
	"shop/shop-srv/order-srv/utils"
	"syscall"

	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"go.uber.org/zap"
)

func main() {
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	initialize.InitRedis()
	// 初始化srv的连接池
	initialize.InitSrvConn()
	if global.ServerConfig.Port == 0 {
		global.ServerConfig.Port, _ = utils.GetFreePort()
	}
	zap.S().Info(global.ServerConfig)

	server := grpc.NewServer()
	proto.RegisterOrderServer(server, &handler.OrderService{})
	// 注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", global.ServerConfig.Host, global.ServerConfig.Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	// 注册到consul并启动服务
	//registerConsulAndStartServer(server, lis)
	startServer(server, lis)
}

// 直接启动服务
func startServer(server *grpc.Server, lis net.Listener) {
	err := server.Serve(lis)
	if err != nil {
		panic("failed to start grpc:" + err.Error())
	}
}

// 注册到consul并启动服务
func registerConsulAndStartServer(server *grpc.Server, lis net.Listener) {
	// 注册到consul
	consulCfg := api.DefaultConfig()
	consulCfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	consulClient, err := api.NewClient(consulCfg)
	if err != nil {
		panic("failed to new consul client:" + err.Error())
	}
	// 生成consul需要的健康检查对象
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("host.docker.internal:%d", global.ServerConfig.Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s",
	}
	serviceID := uuid.NewV4().String()
	registration := &api.AgentServiceRegistration{
		Name:    global.ServerConfig.Name,
		ID:      serviceID,
		Address: "localhost",
		Port:    global.ServerConfig.Port,
		Tags:    global.ServerConfig.Tags,
		Check:   check,
	}
	if err = consulClient.Agent().ServiceRegister(registration); err != nil {
		panic("failed to consul ServiceRegister:" + err.Error())
	}

	go func() {
		startServer(server, lis)
	}()
	// 接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	// 从consul注销服务
	if err = consulClient.Agent().ServiceDeregister(serviceID); err != nil {
		zap.S().Info("注销失败：", err.Error())
	} else {
		zap.S().Info("注销成功")
	}
}
