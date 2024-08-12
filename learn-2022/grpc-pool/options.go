package grpc_pool

import (
	"context"
	"time"

	"google.golang.org/grpc/keepalive"

	"google.golang.org/grpc"
)

const (
	// 创建连接的超时时间，默认为5秒
	DialTimeout = 5 * time.Second

	// BackoffMaxDelay provided maximum delay when backing off after failed connection attempts.
	// 连接失败后退出的最大延迟时间，设置为3秒
	BackoffMaxDelay = 3 * time.Second

	// 如果客户端在该时间段内没有看到任何活动，则会向服务器发送ping以检查传输是否仍然可用，设置为10秒
	KeepAliveTime = time.Duration(10) * time.Second

	// 客户端在发送ping进行保活检查后等待活动的超时时间，如果在此时间段内仍未看到活动，则关闭连接，设置为3秒
	KeepAliveTimeout = time.Duration(3) * time.Second

	// 设置初始窗口大小为1GB，以提供系统的吞吐量
	InitialWindowSize = 1 << 30

	// 设置初始连接窗口大小为1GB，以提供系统的吞吐量
	InitialConnWindowSize = 1 << 30

	// 设置最大发送消息大小为4GB，如果发送的消息大小超过此值，gRPC将报告错误
	MaxSendMsgSize = 4 << 30

	// 设置最大接收消息大小为4GB，如果接收到的消息大小超过此值，gRPC将报告错误
	MaxRecvMsgSize = 4 << 30
)

// 用于设置创建gRPC连接池的参数
type Options struct {
	// 用于创建和配置连接的函数
	Dial func(address string) (*grpc.ClientConn, error)

	// 连接池中空闲连接的最大数量
	MaxIdle int

	// 连接池中同时分配的最大连接数。当为零时，连接池中的连接数没有限制
	MaxActive int

	// 每个单连接上并发流的数量限制
	MaxConcurrentStreams int

	// 连接是否可重用，如果为true，并且连接池已达到MaxActive限制，则Get()方法会重用连接并返回。如果为false，并且连接池已达到MaxActive限制，则会创建一个一次性连接并返回。
	Reuse bool
}

var DefaultOptions = Options{
	Dial:                 Dial,
	MaxIdle:              8,
	MaxActive:            64,
	MaxConcurrentStreams: 64,
	Reuse:                true,
}

// 该函数用于创建一个带有预定义配置的gRPC连接器对象
func Dial(address string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DialTimeout)
	defer cancel()
	// 此时并不会真的建立连接，只是返回了一个连接器对象
	return grpc.DialContext(ctx, address, grpc.WithInsecure(),
		grpc.WithBackoffMaxDelay(BackoffMaxDelay),
		grpc.WithInitialWindowSize(InitialWindowSize),
		grpc.WithInitialConnWindowSize(InitialConnWindowSize),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(MaxSendMsgSize)),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MaxRecvMsgSize)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                KeepAliveTime,
			Timeout:             KeepAliveTimeout,
			PermitWithoutStream: true,
		}))
}

// 该函数用于创建一个简单的、带有预定义配置的gRPC连接
func DialTest(address string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DialTimeout)
	defer cancel()
	return grpc.DialContext(ctx, address, grpc.WithInsecure())
}
