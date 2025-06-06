package registry

import "context"

// 服务注册接口
type Registrar interface {
	//注册
	Register(ctx context.Context, service *ServiceInstance) error
	//注销
	Deregister(ctx context.Context, service *ServiceInstance) error
}

// 服务发现接口
type Discovery interface {
	// 获取服务实例 通过Name   如果使用id的话只有一个实例，通过Name可以获取多个实例进行负载
	GetService(ctx context.Context, serviceName string) ([]*ServiceInstance, error)
	// 创建服务监听器
	Watch(ctx context.Context, serviceName string) (Watcher, error)
}

type Watcher interface {
	// 获取服务实例，next在下面的情况下会返回服务
	// 1. 第一次监听时，如果服务实例列表不为空，则返回服务实例列表
	// 2. 如果服务实例发生变化，则返回服务实例列表
	// 3. 如果上面两种情况都不满足，则会阻塞到context deadline或者cancel
	Next() ([]*ServiceInstance, error)

	// 主动放弃监听
	Stop() error
}

// 服务实例
type ServiceInstance struct {
	// 注册到注册中心的服务ID
	ID string `json:"id"`
	// 服务名称
	Name string `json:"name"`
	// 服务版本
	Version string `json:"version"`
	// 服务元数据
	Metadata map[string]string `json:"metadata"`

	//http://127.0.0.1:8080
	//grpc://127.0.0.1:9000
	Endpoints []string `json:"endpoints"`
}
