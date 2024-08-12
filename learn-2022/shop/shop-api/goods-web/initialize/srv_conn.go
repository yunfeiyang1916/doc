package initialize

import (
	"fmt"
	"grpc_pool"
	"shop/shop-api/goods-web/global"
	"shop/shop-srv/goods-srv/proto"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func InitSrvConn() {
	host := global.ServerConfig.GoodsSrvInfo.Host
	port := global.ServerConfig.GoodsSrvInfo.Port

	// 暂时也保留全局单实例连接
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【商品服务失败】")
	}
	global.GoodsSrvClient = proto.NewGoodsClient(conn)

	// 初始化连接池
	pool, err := grpc_pool.New(fmt.Sprintf("%s:%d", host, port), grpc_pool.DefaultOptions)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 初始化商品服务连接池失败")
	}
	global.GoodsSrvConnPool = pool
}

func InitSrvConn2() {
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
	conn, err := grpc.Dial(fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
	}
	global.GoodsSrvClient = proto.NewGoodsClient(conn)
}
