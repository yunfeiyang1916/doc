package initialize

import (
	"fmt"
	"shop/shop-api/user-web/global"
	"shop/shop-srv/user-srv/proto"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func InitSrvConn() {
	//host := global.ServerConfig.UserSrvInfo.Host
	//port := global.ServerConfig.UserSrvInfo.Port
	// 从注册中心获取
	//consulCfg := api.DefaultConfig()
	//consulCfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	//consulClient, err := api.NewClient(consulCfg)
	//if err != nil {
	//	panic("failed to new consul client:" + err.Error())
	//}
	//data, err := consulClient.Agent().ServicesWithFilter(fmt.Sprintf(`Service=="%s"`, global.ServerConfig.UserSrvInfo.Name))
	//if err != nil {
	//	panic("failed to consul ServicesWithFilter:" + err.Error())
	//}
	//// 只取第一个
	//for _, v := range data {
	//	host = v.Address
	//	port = v.Port
	//	break
	//}
	consulInfo := global.ServerConfig.ConsulInfo
	// 使用grpc-consul-resolver 直接从consul取
	userConn, err := grpc.Dial(fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
	}
	global.UserSrvClient = proto.NewUserClient(userConn)
}
