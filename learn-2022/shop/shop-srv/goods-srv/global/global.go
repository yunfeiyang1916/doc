package global

import (
	"shop/shop-srv/goods-srv/config"

	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	ServerConfig = &config.ServerConfig{}
	// 配置中心，从本地配置文读取
	NacosConfig = &config.NacosConfig{}
	EsClient    *elastic.Client
)
