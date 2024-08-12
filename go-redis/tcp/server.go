package tcp

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/yunfeiyang1916/doc/go-redis/lib/logger"

	"github.com/yunfeiyang1916/doc/go-redis/interface/tcp"
)

// tcp服务配置
type Config struct {
	Address string
}

func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) error {
	closeChan := make(chan struct{})
	// 系统信号量
	sigChan := make(chan os.Signal)
	// 监听系统退出信号
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigChan
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			// 发送关闭信息
			closeChan <- struct{}{}
		}
	}()
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("bind: %s, start listening...", cfg.Address))
	ListenAndServe(listener, handler, closeChan)
	return nil
}

func ListenAndServe(listener net.Listener, handler tcp.Handler, closeChan <-chan struct{}) {
	go func() {
		// 收到系统关闭信号
		<-closeChan
		logger.Info("shutting down")
		// 回收资源
		_ = listener.Close()
		_ = handler.Close()
	}()
	defer func() {
		// 回收资源
		_ = listener.Close()
		_ = handler.Close()
	}()
	ctx := context.Background()
	// 用于等待所有链接关闭
	var waitDone sync.WaitGroup
	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		logger.Info("accepted link")
		waitDone.Add(1)
		go func() {
			defer func() {
				waitDone.Done()
			}()
			handler.Handle(ctx, conn)
		}()
	}
	waitDone.Wait()
}
