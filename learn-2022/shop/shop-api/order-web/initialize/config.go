package initialize

import (
	"encoding/json"
	"shop/shop-api/order-web/config"
	"shop/shop-api/order-web/global"

	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/nacos-group/nacos-sdk-go/clients"

	"github.com/nacos-group/nacos-sdk-go/common/constant"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func GetEnvInfo(env string) bool {
	//刚才设置的环境变量 想要生效 我们必须得重启所有打开的goland
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

// 从本地读取nacos配置，然后从nacos中读取其他配置
func InitConfig2() {
	configFileName := "config-debug.yaml"

	v := viper.New()
	// 设置文件路径
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	var serverConfig config.ServerConfig
	if err := v.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}
	global.NacosConfig = &serverConfig.NacosConfig
	zap.S().Infof("配置信息：%v", serverConfig)
	// 从nacos中读取配置
	sc := []constant.ServerConfig{
		{
			IpAddr: serverConfig.NacosConfig.Host,
			Port:   serverConfig.NacosConfig.Port,
		},
	}
	cc := constant.ClientConfig{
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		NamespaceId:         serverConfig.NacosConfig.Namespace,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		//Username:            serverConfig.NacosConfig.User,
		//Password:            serverConfig.NacosConfig.Password,
		LogLevel: "debug",
		//Namespace:           "public",
		//RegionId:            "cn-hangzhou",
		//Endpoint:            "nacos-01.nacos-01.svc.cluster.local:8848",
		//ServerConfigs: sc,
	}
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		panic(err)
	}
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: serverConfig.NacosConfig.DataId,
		Group:  serverConfig.NacosConfig.Group,
	})
	if err != nil {
		panic(err)
	}
	zap.S().Infof("从nacos读取到的配置信息：%s", content)
	var newConfig config.ServerConfig
	if err = json.Unmarshal([]byte(content), &newConfig); err != nil {
		panic(err)
	}
	global.ServerConfig = &newConfig
	if err = configClient.ListenConfig(vo.ConfigParam{
		DataId: serverConfig.NacosConfig.DataId,
		Group:  serverConfig.NacosConfig.Group,
		OnChange: func(namespace, group, dataId, data string) {
			zap.S().Infof("配置文件产生变化：%s", dataId)
			var newConfig config.ServerConfig
			if err = json.Unmarshal([]byte(data), &newConfig); err != nil {
				panic(err)
			}
			global.ServerConfig = &newConfig
		},
	}); err != nil {
		panic(err)
	}
}

func InitConfig() {
	//debug := GetEnvInfo("MXSHOP_DEBUG")
	//configFilePrefix := "config"
	//configFileName := fmt.Sprintf("%s-pro.yaml", configFilePrefix)
	//if debug {
	//	configFileName = fmt.Sprintf("%s-debug.yaml", configFilePrefix)
	//}

	configFileName := "config-debug.yaml"

	v := viper.New()
	// 设置文件路径
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息：%v", global.ServerConfig)
	// 动态监控文件变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("配置文件产生变化：%s", e.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.ServerConfig)
		zap.S().Infof("配置信息：%v", global.ServerConfig)
	})
}
