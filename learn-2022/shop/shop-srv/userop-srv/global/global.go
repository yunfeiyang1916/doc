package global

import (
	"shop/shop-srv/userop-srv/config"

	redsyncredis "github.com/go-redsync/redsync/v4/redis"
	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	RedisPool    redsyncredis.Pool
	ServerConfig = &config.ServerConfig{}
	// 配置中心，从本地配置文读取
	NacosConfig = &config.NacosConfig{}
)
