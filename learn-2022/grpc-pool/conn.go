package grpc_pool

import "google.golang.org/grpc"

// 封装了grpc.ClientConn的连接接口
type Conn interface {
	// 返回一个实际的grpc.ClientConn类型的指针
	Value() *grpc.ClientConn

	// 减少对grpc连接的引用，并不是关闭它。如果连接池已满，则直接关闭连接
	Close() error
}

type conn struct {
	cc   *grpc.ClientConn
	pool *pool
	// 是否是一次性临时连接
	once bool
}

// 返回一个实际的grpc.ClientConn类型的指针
func (c *conn) Value() *grpc.ClientConn {
	return c.cc
}

// 减少对grpc连接的引用，而不是关闭它。如果连接池已满，则直接关闭连接
func (c *conn) Close() error {
	c.pool.decrRef()
	if c.once {
		return c.reset()
	}
	return nil
}

// 重置连接，会直接关闭连接
func (c *conn) reset() error {
	cc := c.cc
	c.cc = nil
	c.once = false
	if cc != nil {
		return cc.Close()
	}
	return nil
}
