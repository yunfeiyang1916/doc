package tcp

import (
	"bufio"
	"context"
	"io"
	"net"
	"sync"
	"time"

	"github.com/yunfeiyang1916/doc/go-redis/lib/logger"

	"github.com/yunfeiyang1916/doc/go-redis/lib/sync/wait"

	"github.com/yunfeiyang1916/doc/go-redis/lib/sync/atomic"
)

// 客户端
type EchoClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

func (e *EchoClient) Close() error {
	// 等待10秒之后关闭
	e.Waiting.WaitWithTimeout(10 * time.Second)
	_ = e.Conn.Close()
	return nil
}

// 回声(客户端发什么内容就回复什么内容)服务处理器
type EchoHandler struct {
	activeConn sync.Map
	// 是否正在关闭
	closing atomic.Boolean
}

func MakeHandler() *EchoHandler {
	return &EchoHandler{}
}

// 业务处理
func (h *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	if h.closing.Get() {
		_ = conn.Close()
	}
	// 包装成client
	client := &EchoClient{
		Conn: conn,
	}
	// 以client为键，不需要值
	h.activeConn.Store(client, struct{}{})
	// 使用 bufio 标准库提供的缓冲区功能
	reader := bufio.NewReader(conn)
	for {
		// ReadString 会一直阻塞直到遇到分隔符 '\n'
		// 遇到分隔符后会返回上次遇到分隔符或连接建立后收到的所有数据, 包括分隔符本身
		// 若在遇到分隔符之前遇到异常, ReadString 会返回已收到的数据和错误信息
		msg, err := reader.ReadString('\n')
		if err != nil {
			// 表示客户端退出
			if err == io.EOF {
				logger.Info("Connecting close")
				h.activeConn.Delete(client)
			} else {
				logger.Warn(err)
			}
			return
		}
		// 要处理业务，客户端先不要关闭我
		client.Waiting.Add(1)
		b := []byte(msg)
		_, _ = conn.Write(b)
		client.Waiting.Done()
	}
}

func (h *EchoHandler) Close() error {
	logger.Info("handler shutting down")
	h.closing.Set(true)
	// 关闭所有已建立的链接
	h.activeConn.Range(func(key, value any) bool {
		client := key.(*EchoClient)
		_ = client.Close()
		// 继续处理下一个循环
		return true
	})
	return nil
}
