package main

import (
	"fmt"
	"shop/shop-api/user-web/global"
	"shop/shop-api/user-web/initialize"
	"shop/shop-api/user-web/utils/register/consul"

	ut "github.com/go-playground/universal-translator"

	myvalidator "shop/shop-api/user-web/validator"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", myvalidator.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	zap.S().Infof("启动服务，端口：%d", global.ServerConfig.Port)
	if err := r.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动失败：", err.Error())
	}
}
