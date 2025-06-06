package app

import (
	"context"
	"net/url"
	"os"
	"os/signal"
	"shop-v2/gmicro/registry"
	"shop-v2/gmicro/server"
	"shop-v2/pkg/log"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/google/uuid"
)

type App struct {
	opts options

	mu sync.Mutex
	// 受保护
	instance *registry.ServiceInstance
	cancel   func()
}

func New(opts ...Option) *App {
	o := options{
		sigs:             []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		registrarTimeout: 10 * time.Second,
		stopTimeout:      10 * time.Second,
	}
	if id, err := uuid.NewUUID(); err == nil {
		o.id = id.String()
	}

	for _, opt := range opts {
		opt(&o)
	}
	return &App{opts: o}
}

// 启动整个服务
func (a *App) Run() error {
	// 注册的信息
	instance, err := a.buildInstance()
	if err != nil {
		return err
	}
	//这个变量可能被其他的goroutine访问
	a.mu.Lock()
	a.instance = instance
	a.mu.Unlock()

	// 启动服务
	var servers []server.Server
	if a.opts.rpcServer != nil {
		servers = append(servers, a.opts.rpcServer)
	}

	// app在stop的时候想要通知到服务下进行cancel
	// 这时候我们自己生成一个context，把cancel方法注入到app当中，这时候在stop的时候cancel方法就能通知到服务中
	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel
	eg, ctx := errgroup.WithContext(ctx)
	wg := sync.WaitGroup{}

	for _, srv := range servers {
		// 赋值临时变量，防止go协程闭包
		sr := srv
		eg.Go(func() error {
			// 等待stop信号
			<-ctx.Done()
			sctx, cancel := context.WithTimeout(context.Background(), a.opts.stopTimeout)
			defer cancel()
			return sr.Stop(sctx)
		})
		wg.Add(1)
		eg.Go(func() error {
			// 启动服务
			wg.Done()
			log.Info("start server")
			return sr.Start(ctx)
		})
	}
	wg.Wait()

	// 是否需要注册服务
	if a.opts.registrar != nil {
		//rctx, rcancel := context.WithTimeout(context.Background(), a.opts.registrarTimeout)
		//defer rcancel()
		//err = a.opts.registrar.Register(rctx, instance)
		//if err != nil {
		//	log.Errorf("registrar service error: %s", err)
		//	return err
		//}
	}

	// 监听退出信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opts.sigs...)
	// <-c
	// 由于a.cancel()执行的很快 导致整个goroutine程序退出  所以放到goroutine里监听chan。达到一个阻塞的效果
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c:
			return a.Stop()
		}
	})
	if err = eg.Wait(); err != nil {
		return err
	}
	return nil
}

func (a *App) Stop() error {
	a.mu.Lock()
	instance := a.instance
	a.mu.Unlock()
	log.Info("start deregister service")
	if a.opts.registrar != nil && instance != nil {
		rctx, rcancel := context.WithTimeout(context.Background(), a.opts.stopTimeout)
		defer rcancel()
		err := a.opts.registrar.Deregister(rctx, instance)
		if err != nil {
			log.Errorf("deregister service error: %s", err)
			return err
		}
	}
	// 自己生成的context生成cancel后往服务中传递，所以能通知到所有的服务下的context
	if a.cancel != nil {
		log.Infof("start cancel context")
		a.cancel()
	}
	return nil
}

func (a *App) buildInstance() (*registry.ServiceInstance, error) {
	endpoints := make([]string, 0)
	for _, ep := range a.opts.endpoints {
		endpoints = append(endpoints, ep.String())
	}
	// 从rpcserver中获取
	if a.opts.rpcServer != nil {
		if a.opts.rpcServer.Endpoint() != nil {
			endpoints = append(endpoints, a.opts.rpcServer.Endpoint().String())
		} else {
			u := url.URL{
				Scheme: "grpc",
				Host:   a.opts.rpcServer.Address(),
			}
			endpoints = append(endpoints, u.String())
		}
	}

	return &registry.ServiceInstance{
		ID:        a.opts.id,
		Name:      a.opts.name,
		Endpoints: endpoints,
	}, nil
}
