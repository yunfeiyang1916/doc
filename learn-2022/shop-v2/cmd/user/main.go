package main

import (
	"fmt"
	pb "shop-v2/api/user/v1"
	"shop-v2/gmicro/app"
	"shop-v2/gmicro/core/trace"
	"shop-v2/gmicro/registry"
	"shop-v2/gmicro/registry/consul"
	"shop-v2/gmicro/server/rpcserver"
	"shop-v2/internal/share/options"
	"shop-v2/internal/user/config"
	"shop-v2/internal/user/controller"
	"shop-v2/internal/user/data/mock"
	"shop-v2/internal/user/service"
	cmd "shop-v2/pkg/app"
	"shop-v2/pkg/log"

	"github.com/hashicorp/consul/api"
)

func main() {
	cfg := config.New()
	// 先创建一个cmd的app
	cmdl := cmd.NewApp("user-srv", "user-srv",
		cmd.WithDescription("user service"),
		cmd.WithOptions(cfg),
		cmd.WithRunFunc(run(cfg)),
	)
	cmdl.Run()
}

func NewRegistrar(registry *options.RegistryOptions) registry.Registrar {
	c := api.DefaultConfig()
	c.Address = registry.Address
	c.Scheme = registry.Scheme
	cli, err := api.NewClient(c)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(true))
	return r
}

func run(cfg *config.Config) cmd.RunFunc {
	return func(basename string) error {
		// 初始化log
		log.Init(cfg.Log)
		// 服务退出时落盘
		defer log.Flush()
		log.Infof("server start basename: %s", basename)
		log.Infof("log.level:%v", cfg.Log.Level)

		// 初始化open-telemetry的exporter
		trace.InitAgent(trace.Options{
			Name:     cfg.Telemetry.Name,
			Endpoint: cfg.Telemetry.Endpoint,
			Sampler:  cfg.Telemetry.Sampler,
			Batcher:  cfg.Telemetry.Batcher,
		})

		data := mock.NewUserMock()
		srv := service.NewUserService(data)
		userServer := controller.NewUserServer(srv)

		// 创建rpc服务
		rpcAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
		rpcSrv := rpcserver.NewServer(rpcserver.WithAddress(rpcAddr))
		// 服务注册
		register := NewRegistrar(cfg.Registry)
		appl := app.New(app.WithName(cfg.Server.Name),
			app.WithRPCServer(rpcSrv),
			app.WithRegistrar(register))
		if err := appl.Run(); err != nil {
			log.Fatalf("run app server error:%+v", err)
		}

		pb.RegisterUserServer(rpcSrv.Server, userServer)
		return nil
	}
}
