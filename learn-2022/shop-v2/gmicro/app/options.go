package app

import (
	"net/url"
	"os"
	"shop-v2/gmicro/registry"
	"shop-v2/gmicro/server/rpcserver"
	"time"
)

type Option func(*options)

type options struct {
	id        string
	name      string
	endpoints []*url.URL

	// 需要监听的信号
	sigs []os.Signal

	registrarTimeout time.Duration
	// 允许用户传入自己的实现
	registrar registry.Registrar

	// stop超时时间
	stopTimeout time.Duration

	// 传递rpc服务
	rpcServer *rpcserver.Server
}

func WithOptions(endpoints ...*url.URL) Option {
	return func(o *options) {
		o.endpoints = endpoints
	}
}

func WithID(id string) Option {
	return func(o *options) {
		o.id = id
	}
}

func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

func WithSigs(sigs []os.Signal) Option {
	return func(o *options) {
		o.sigs = sigs
	}
}

func WithRegistrar(registrar registry.Registrar) Option {
	return func(o *options) {
		o.registrar = registrar
	}
}

func WithRPCServer(server *rpcserver.Server) Option {
	return func(o *options) {
		o.rpcServer = server
	}
}
