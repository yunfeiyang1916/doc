package rpcserver

import (
	"context"
	"shop-v2/gmicro/registry"
	"shop-v2/gmicro/server/rpcserver/clientinterceptors"
	"shop-v2/gmicro/server/rpcserver/resolver/discovery"
	"shop-v2/pkg/log"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	grpcinsecure "google.golang.org/grpc/credentials/insecure"
)

type ClientOption func(o *clientOptions)
type clientOptions struct {
	//设置目标地址
	endpoint string
	timeout  time.Duration
	//discovery接口
	discovery registry.Discovery
	//可传递拦截器
	unaryInts []grpc.UnaryClientInterceptor
	//stream
	streamInts []grpc.StreamClientInterceptor
	//用户可自己设置grpc连接的结构体
	rpcOpts []grpc.DialOption
	//根据Name生成负载均衡的策略
	balancerName string

	logger log.Logger
	//是否启用链路追踪
	enableTracing bool

	enableMetrics bool
}

// 设置是否开启链路追踪
func WithEnableTracing(enable bool) ClientOption {
	return func(o *clientOptions) {
		o.enableTracing = enable
	}
}

// 设置是否开启指标Metrics
func WithClientMetrics(enable bool) ClientOption {
	return func(o *clientOptions) {
		o.enableMetrics = enable
	}
}

// 设置地址
func WithEndpoint(endpoint string) ClientOption {
	return func(o *clientOptions) {
		o.endpoint = endpoint
	}
}

// 设置超时时间
func WithClientTimeout(timeout time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.timeout = timeout
	}
}

// 设置服务发现
func WithDiscovery(d registry.Discovery) ClientOption {
	return func(o *clientOptions) {
		o.discovery = d
	}
}

// 设置拦截器
func WithClientUnaryInterceptor(in ...grpc.UnaryClientInterceptor) ClientOption {
	return func(o *clientOptions) {
		o.unaryInts = in
	}
}

// 设置stream拦截器
func WithClientStreamInterceptor(in ...grpc.StreamClientInterceptor) ClientOption {
	return func(o *clientOptions) {
		o.streamInts = in
	}
}

// 设置grpc的dial选项
func WithClientOptions(opts ...grpc.DialOption) ClientOption {
	return func(o *clientOptions) {
		o.rpcOpts = opts
	}
}

// 设置负载均衡器
func WithBalancerName(name string) ClientOption {
	return func(o *clientOptions) {
		o.balancerName = name
	}
}

// 拨号
func DialInsecure(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, true, opts...)
}
func Dial(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, false, opts...)
}

func dial(ctx context.Context, insecure bool, opts ...ClientOption) (*grpc.ClientConn, error) {
	options := clientOptions{
		timeout:       2000 * time.Millisecond,
		balancerName:  "round_robin",
		enableTracing: true,
	}
	for _, opt := range opts {
		opt(&options)
	}
	//TODO 客户端默认拦截器
	ints := []grpc.UnaryClientInterceptor{
		clientinterceptors.TimeoutInterceptor(options.timeout),
	}
	// 可给用户自己设置需不需链路追踪
	if options.enableTracing {
		ints = append(ints, otelgrpc.UnaryClientInterceptor())
	}
	// 可给用户自己设置需不需指标
	if options.enableMetrics {
		//ints = append(ints, clientinterceptors.PrometheusInterceptor())
	}
	streamInts := []grpc.StreamClientInterceptor{}

	if len(options.unaryInts) > 0 {
		ints = append(ints, options.unaryInts...)
	}
	if len(options.streamInts) > 0 {
		streamInts = append(streamInts, options.streamInts...)
	}
	// 可以由用户端自己传递 这些默认的
	grpcOpts := []grpc.DialOption{
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "` + options.balancerName + `"}`),
		grpc.WithChainUnaryInterceptor(ints...),
		grpc.WithChainStreamInterceptor(streamInts...),
	}
	// 服务发现的选项
	if options.discovery != nil {
		grpcOpts = append(grpcOpts, grpc.WithResolvers(discovery.NewBuilder(options.discovery, discovery.WithInsecure(insecure))))
	}
	if insecure {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(grpcinsecure.NewCredentials()))
	}
	if len(options.rpcOpts) > 0 {
		grpcOpts = append(grpcOpts, options.rpcOpts...)
	}
	return grpc.DialContext(ctx, options.endpoint, grpcOpts...)
}
