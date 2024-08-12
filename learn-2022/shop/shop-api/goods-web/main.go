package main

import (
	"fmt"
	"os"
	"os/signal"
	"shop/shop-api/goods-web/global"
	"shop/shop-api/goods-web/initialize"
	"shop/shop-api/goods-web/utils/register/consul"
	"syscall"

	uuid "github.com/satori/go.uuid"

	"go.uber.org/zap"
)

func main() {
	// 1 初始化logger
	initialize.InitLogger()

	// 2 初始化配置文件
	initialize.InitConfig()

	// 3 初始化路由
	r := initialize.Routers()

	// 4 初始化翻译
	if err := initialize.InitTrans("zh"); err != nil {
		panic(err)
	}
	// 5 初始化srv的连接
	initialize.InitSrvConn()

	// 注册服务到注册中心
	registryClient := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId := uuid.NewV4().String()
	if err := registryClient.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId); err != nil {
		zap.S().Error("服务注册失败：", err.Error())
	}
	zap.S().Infof("启动服务，端口：%d", global.ServerConfig.Port)
	go func() {
		if err := r.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.S().Panic("启动失败：", err.Error())
		}
	}()
	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err := registryClient.Deregister(serviceId); err != nil {
		zap.S().Error("注销失败:", err.Error())
	} else {
		zap.S().Info("注销成功:")
	}
}
