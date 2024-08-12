package grpc_pool_channel

import (
	"context"
	"errors"
	"sync"
	"time"

	"google.golang.org/grpc"
)

var (
	// ErrClosed 当客户端连接池已关闭时返回的错误
	ErrClosed = errors.New("grpc pool: client pool is closed")
	// ErrTimeout 当客户端连接池连接超时时返回的错误
	ErrTimeout = errors.New("grpc pool: client pool timed out")
	// ErrAlreadyClosed 当客户端连接已被关闭时返回的错误
	ErrAlreadyClosed = errors.New("grpc pool: the connection was already closed")
	// ErrFullPool 当池已满时关闭一个连接返回的错误
	ErrFullPool = errors.New("grpc pool: closing a ClientConn into a full pool")
)

// 用于创建一个gRPC客户端的
type Factory func() (*grpc.ClientConn, error)

// 该函数类型用于创建一个grpc客户端，它接受一个context参数，该参数可以从Get或NewWithContext方法中传递。
type FactoryWithContext func(ctx context.Context) (*grpc.ClientConn, error)

// grpc客户端连接的包装器
type ClientConn struct {
	*grpc.ClientConn
	pool *Pool
	// 记录连接上一次被使用的时间
	timeUsed time.Time
	// 记录连接初始化的时间
	timeInitiated time.Time
	// 标识连接是否处于不健康的状态
	unhealthy bool
}

// grpc客户端连接池
type Pool struct {
	// 一个ClientConn类型的通道，用于存储客户端连接
	clients chan ClientConn
	// 用于创建客户端连接
	factory FactoryWithContext
	// 表示空闲连接的超时时间
	idleTimeout time.Duration
	// 表示客户端连接的最大生命周期
	maxLifeDuration time.Duration
	mu              sync.RWMutex
}

func (p *Pool) getClients() chan ClientConn {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.clients
}

// 该函数用于清空连接池中的所有客户端，并关闭连接池
func (p *Pool) Close() {
	p.mu.Lock()
	clients := p.clients
	p.clients = nil
	p.mu.Unlock()

	if clients == nil {
		return
	}

	close(clients)
	for client := range clients {
		if client.ClientConn == nil {
			continue
		}
		client.ClientConn.Close()
	}
}

// 该函数用于判断客户端连接池是否已关闭
func (p *Pool) IsClosed() bool {
	return p == nil || p.getClients() == nil
}

// 它会尝试从连接池中获取一个可用的连接，如果连接池未达到容量限制，则会创建一个新的连接。
// 如果连接池已满，则会等待直到有连接变为可用或超时。如果超时时间为0，则表示无限等待。
func (p *Pool) Get(ctx context.Context) (*ClientConn, error) {
	clients := p.getClients()
	if clients == nil {
		return nil, ErrClosed
	}
	wrapper := ClientConn{pool: p}
	select {
	// 从队列里面取出一个客户端连接
	case wrapper = <-clients:
	// All good
	case <-ctx.Done():
		// it would better returns ctx.Err()
		return nil, ErrTimeout
	}
	// 如果包装器空闲时间过长，请关闭连接并创建一个新连接。
	// 可以安全地假设没有新的客户端，因为我们获取的客户端是通道中的第一个客户端
	idleTimeout := p.idleTimeout
	if wrapper.ClientConn != nil && idleTimeout > 0 &&
		wrapper.timeUsed.Add(idleTimeout).Before(time.Now()) {

		wrapper.ClientConn.Close()
		wrapper.ClientConn = nil
	}

	var err error
	if wrapper.ClientConn == nil {
		wrapper.ClientConn, err = p.factory(ctx)
		if err != nil {
			// 如果出现错误，我们希望在管道中放回空客户端连接的占位符
			clients <- ClientConn{
				pool: p,
			}
		}
		// 这是一个新的连接，重置它的初始时间
		wrapper.timeInitiated = time.Now()
	}

	return &wrapper, err
}

// 获取连接池的容量
func (p *Pool) Capacity() int {
	if p.IsClosed() {
		return 0
	}
	return cap(p.clients)
}

// 返回当前未使用的客户端数量
func (p *Pool) Available() int {
	if p.IsClosed() {
		return 0
	}
	return len(p.clients)
}

// 该函数用于将客户端连接标记为不健康，以便在关闭连接时重置连接。
func (c *ClientConn) Unhealthy() {
	c.unhealthy = true
}

// 将连接放回到连接池
func (c *ClientConn) Close() error {
	if c == nil {
		return nil
	}
	if c.ClientConn == nil {
		return ErrAlreadyClosed
	}
	if c.pool.IsClosed() {
		return ErrClosed
	}
	// If the wrapper connection has become too old, we want to recycle it. To
	// clarify the logic: if the sum of the initialization time and the max
	// duration is before Now(), it means the initialization is so old adding
	// the maximum duration couldn't put in the future. This sum therefore
	// corresponds to the cut-off point: if it's in the future we still have
	// time, if it's in the past it's too old

	// 检查一个连接是否超时，如果超时则将其标记为不健康状态。
	// 这个逻辑的目的是为了避免使用过时的连接，保证系统的安全性和稳定性。
	// 如果连接的初始化时间加上最大生命周期时长小于当前时间，则将连接标记为不健康状态。
	maxDuration := c.pool.maxLifeDuration
	if maxDuration > 0 && c.timeInitiated.Add(maxDuration).Before(time.Now()) {
		c.Unhealthy()
	}

	// 先克隆一个包装器，以便在用户使用的包装器中将ClientConn设置为nil
	wrapper := ClientConn{
		pool:       c.pool,
		ClientConn: c.ClientConn,
		timeUsed:   time.Now(),
	}
	if c.unhealthy {
		wrapper.ClientConn.Close()
		wrapper.ClientConn = nil
	} else {
		wrapper.timeInitiated = c.timeInitiated
	}
	select {
	case c.pool.clients <- wrapper:
		// All good
	default:
		return ErrFullPool
	}
	// 标记为关闭
	c.ClientConn = nil
	return nil
}

// 该函数用于创建一个具有指定初始和最大容量以及空闲客户端超时的客户端池。
// 函数还接受一个上下文参数，该参数在初始化时传递给工厂方法。
// 如果初始客户端无法创建，则返回错误。
func NewWithContext(ctx context.Context, factory FactoryWithContext, init, capacity int, idleTimeout time.Duration, maxLifeDuration ...time.Duration) (*Pool, error) {
	if capacity <= 0 {
		capacity = 1
	}
	if init < 0 {
		init = 0
	}
	if init > capacity {
		init = capacity
	}
	p := &Pool{
		clients:     make(chan ClientConn, capacity),
		factory:     factory,
		idleTimeout: idleTimeout,
	}
	if len(maxLifeDuration) > 0 {
		p.maxLifeDuration = maxLifeDuration[0]
	}
	for i := 0; i < init; i++ {
		c, err := factory(ctx)
		if err != nil {
			return nil, err
		}

		p.clients <- ClientConn{
			ClientConn:    c,
			pool:          p,
			timeUsed:      time.Now(),
			timeInitiated: time.Now(),
		}
	}
	// Fill the rest of the pool with empty clients
	for i := 0; i < capacity-init; i++ {
		p.clients <- ClientConn{
			pool: p,
		}
	}
	return p, nil
}

// 该函数用于创建一个具有指定初始和最大容量以及空闲客户端超时的客户端池。
// 如果初始客户端无法创建，则返回错误。
func New(factory Factory, init, capacity int, idleTimeout time.Duration, maxLifeDuration ...time.Duration) (*Pool, error) {
	return NewWithContext(context.Background(), func(ctx context.Context) (*grpc.ClientConn, error) { return factory() }, init, capacity, idleTimeout, maxLifeDuration...)
}
