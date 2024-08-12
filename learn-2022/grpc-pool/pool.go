package grpc_pool

import (
	"errors"
	"fmt"
	"log"
	"math"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc"
)

// ErrClosed 用于表示当连接池通过pool.Close()方法被关闭时产生的错误
var ErrClosed = errors.New("pool is closed")

// 线程安全的连接池接口
type Pool interface {
	// 用于从池中获取一个连接。关闭连接会将其放回池中。
	// 如果在池被销毁或已满时关闭连接将被视为错误。连接不为空时，可以保证conn.Value()不为空
	Get() (Conn, error)

	// 于关闭池及其所有连接。关闭后，池不再可用。不能并发调用Close()和Get()方法，否则会导致panic
	Close() error

	// 返回池的当前状态
	Status() string
}

// 连接池，用于管理和复用网络连接
type pool struct {
	// 连接池的选项
	opt Options
	// 用于随机获取连接
	index uint32

	// 当前物理连接
	current int32

	// 表示逻辑连接
	// 逻辑连接等于物理连接乘以MaxConcurrentStreams
	// logic connection = physical connection * MaxConcurrentStreams
	ref int32

	// 创建的所有物理连接的切片
	conns []*conn

	// 用于创建连接的服务器地址
	address string

	// 表示连接池已关闭的原子变量
	closed int32

	// control the atomic var current's concurrent read write.
	sync.RWMutex
}

// 创建连接池
func New(address string, option Options) (Pool, error) {
	if address == "" {
		return nil, errors.New("invalid address settings")
	}
	if option.Dial == nil {
		return nil, errors.New("invalid dial settings")
	}
	if option.MaxIdle <= 0 || option.MaxActive <= 0 || option.MaxIdle > option.MaxActive {
		return nil, errors.New("invalid maximum settings")
	}
	if option.MaxConcurrentStreams <= 0 {
		return nil, errors.New("invalid maximun settings")
	}

	p := &pool{
		index:   0,
		current: int32(option.MaxIdle),
		ref:     0,
		opt:     option,
		conns:   make([]*conn, option.MaxActive),
		address: address,
		closed:  0,
	}

	for i := 0; i < p.opt.MaxIdle; i++ {
		c, err := p.opt.Dial(address)
		if err != nil {
			p.Close()
			return nil, fmt.Errorf("dial is not able to fill the pool: %s", err)
		}
		p.conns[i] = p.wrapConn(c, false)
	}
	log.Printf("new pool success: %v\n", p.Status())

	return p, nil
}

// 递增逻辑连接
func (p *pool) incrRef() int32 {
	newRef := atomic.AddInt32(&p.ref, 1)
	if newRef == math.MaxInt32 {
		panic(fmt.Sprintf("overflow ref: %d", newRef))
	}
	return newRef
}

// 递减逻辑连接
func (p *pool) decrRef() {
	newRef := atomic.AddInt32(&p.ref, -1)
	if newRef < 0 && atomic.LoadInt32(&p.closed) == 0 {
		panic(fmt.Sprintf("negative ref: %d", newRef))
	}
	if newRef == 0 && atomic.LoadInt32(&p.current) > int32(p.opt.MaxIdle) {
		p.Lock()
		if atomic.LoadInt32(&p.ref) == 0 {
			log.Printf("shrink pool: %d ---> %d, decrement: %d, maxActive: %d\n",
				p.current, p.opt.MaxIdle, p.current-int32(p.opt.MaxIdle), p.opt.MaxActive)
			atomic.StoreInt32(&p.current, int32(p.opt.MaxIdle))
			p.deleteFrom(p.opt.MaxIdle)
		}
		p.Unlock()
	}
}

func (p *pool) reset(index int) {
	conn := p.conns[index]
	if conn == nil {
		return
	}
	conn.reset()
	p.conns[index] = nil
}

func (p *pool) deleteFrom(begin int) {
	for i := begin; i < p.opt.MaxActive; i++ {
		p.reset(i)
	}
}

func (p *pool) wrapConn(cc *grpc.ClientConn, once bool) *conn {
	return &conn{
		cc:   cc,
		pool: p,
		once: once,
	}
}

// 从连接池中获取一个连接
func (p *pool) Get() (Conn, error) {
	// 递增逻辑引用计数
	nextRef := p.incrRef()
	p.RLock()
	current := atomic.LoadInt32(&p.current)
	p.RUnlock()
	if current == 0 {
		return nil, ErrClosed
	}
	// 如果当前逻辑引用计数小于等于最大并发数乘以当前物理连接数，则从连接数组中返回一个连接
	if nextRef <= current*int32(p.opt.MaxConcurrentStreams) {
		next := atomic.AddUint32(&p.index, 1) % uint32(current)
		return p.conns[next], nil
	}

	// 如果连接池已达到最大活动连接数
	if current == int32(p.opt.MaxActive) {
		// 如果重用连接，则从连接数组中返回一个连接
		if p.opt.Reuse {
			next := atomic.AddUint32(&p.index, 1) % uint32(current)
			return p.conns[next], nil
		}
		// 创建临时连接
		c, err := p.opt.Dial(p.address)
		return p.wrapConn(c, true), err
	}

	// 创建一批连接放到连接池中
	p.Lock()
	current = atomic.LoadInt32(&p.current)
	// 如果连接池未达到最大活动连接数，并且当前逻辑引用计数大于最大并发数乘以当前物理连接数，则创建新连接
	if current < int32(p.opt.MaxActive) && nextRef > current*int32(p.opt.MaxConcurrentStreams) {
		// 先尝试增加连接数到最大活动连接数的两倍，或者剩余的空闲数
		increment := current
		if current+increment > int32(p.opt.MaxActive) {
			increment = int32(p.opt.MaxActive) - current
		}
		var i int32
		var err error
		for i = 0; i < increment; i++ {
			// 并不会真正建立连接，只是返回一个连接器对象，等真正调用grpc函数时才真正建立连接，所以这里不会阻塞锁
			c, er := p.opt.Dial(p.address)
			if er != nil {
				err = er
				break
			}
			p.reset(int(current + i))
			p.conns[current+i] = p.wrapConn(c, false)
		}
		current += i
		log.Printf("grow pool: %d ---> %d, increment: %d, maxActive: %d\n",
			p.current, current, increment, p.opt.MaxActive)
		atomic.StoreInt32(&p.current, current)
		if err != nil {
			p.Unlock()
			return nil, err
		}
	}
	p.Unlock()
	next := atomic.AddUint32(&p.index, 1) % uint32(current)
	return p.conns[next], nil
}

// 关闭连接池
func (p *pool) Close() error {
	atomic.StoreInt32(&p.closed, 1)
	atomic.StoreUint32(&p.index, 0)
	atomic.StoreInt32(&p.current, 0)
	atomic.StoreInt32(&p.ref, 0)
	p.deleteFrom(0)
	log.Printf("close pool success: %v\n", p.Status())
	return nil
}

// 连接池状态
func (p *pool) Status() string {
	return fmt.Sprintf("address:%s, index:%d, current:%d, ref:%d. option:%v",
		p.address, p.index, p.current, p.ref, p.opt)
}
