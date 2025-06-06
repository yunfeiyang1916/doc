package consul

import (
	"shop-v2/gmicro/registry"
	"sync"
	"sync/atomic"
)

type serviceSet struct {
	serviceName string
	watcher     map[*watcher]struct{}
	services    *atomic.Value
	lock        sync.RWMutex
}

func (s *serviceSet) broadcast(ss []*registry.ServiceInstance) {
	//原子操作   数据全部写完 另一端才能拿到  保证线程安全
	//我们平时些struct的时候
	s.services.Store(ss)
	s.lock.RLock()
	defer s.lock.RUnlock()
	for k := range s.watcher {
		select {
		//把消息放到event event消费方在watcher的Next
		case k.event <- struct{}{}:
		default:
		}
	}
}
