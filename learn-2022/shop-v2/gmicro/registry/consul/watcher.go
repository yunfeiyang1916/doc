package consul

import (
	"context"
	"shop-v2/gmicro/registry"
)

type watcher struct {
	event chan struct{}
	set   *serviceSet

	// for cancel
	ctx    context.Context
	cancel context.CancelFunc
}

func (w *watcher) Next() (services []*registry.ServiceInstance, err error) {
	select {
	case <-w.ctx.Done():
		err = w.ctx.Err()
		return
	case <-w.event:
	}
	//消费 代表数据已经加载完成  ↓开始拿   然后放到全局变量里
	ss, ok := w.set.services.Load().([]*registry.ServiceInstance)

	if ok {
		//放到全局变量   然后观察者（resolver的Next）就能观察到
		services = append(services, ss...)
	}
	return
}

func (w *watcher) Stop() error {
	w.cancel()
	w.set.lock.Lock()
	defer w.set.lock.Unlock()
	delete(w.set.watcher, w)
	return nil
}
