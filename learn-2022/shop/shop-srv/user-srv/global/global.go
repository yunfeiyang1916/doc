package global

import (
	"shop/shop-srv/user-srv/config"

	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	ServerConfig = &config.ServerConfig{}
	// 配置中心，从本地配置文读取
	NacosConfig = &config.NacosConfig{}
)
