package global

import (
	"grpc_pool"
	"shop/shop-api/goods-web/config"
	"shop/shop-srv/goods-srv/proto"

	ut "github.com/go-playground/universal-translator"
)

var (
	Trans        ut.Translator
	ServerConfig = &config.ServerConfig{}
	// 配置中心，从本地配置文读取
	NacosConfig      = &config.NacosConfig{}
	GoodsSrvClient   proto.GoodsClient
	GoodsSrvConnPool grpc_pool.Pool
)
