package global

import (
	"grpc_pool"
	"shop/shop-api/userop-web/config"

	ut "github.com/go-playground/universal-translator"
)

var (
	Trans        ut.Translator
	ServerConfig = &config.ServerConfig{}
	// 配置中心，从本地配置文读取
	NacosConfig       = &config.NacosConfig{}
	UserOpSrvConnPool grpc_pool.Pool
	GoodsSrvConnPool  grpc_pool.Pool
)
