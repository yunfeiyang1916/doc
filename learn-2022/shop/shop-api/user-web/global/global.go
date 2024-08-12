package global

import (
	"shop/shop-api/user-web/config"
	"shop/shop-srv/user-srv/proto"

	ut "github.com/go-playground/universal-translator"
)

var (
	Trans        ut.Translator
	ServerConfig = &config.ServerConfig{}
	// 配置中心，从本地配置文读取
	NacosConfig   = &config.NacosConfig{}
	UserSrvClient proto.UserClient
)
