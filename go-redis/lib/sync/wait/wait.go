package wait

import (
	"sync"
	"time"
)

// 带超时控制的WaitGroup
type Wait struct {
	wg sync.WaitGroup
}

// Add adds delta, which may be negative, to the WaitGroup counter.
func (w *Wait) Add(delta int) {
	w.wg.Add(delta)
}

// Done decrements the WaitGroup counter by one
func (w *Wait) Done() {
	w.wg.Done()
}

// Wait blocks until the WaitGroup counter is zero.
func (w *Wait) Wait() {
	w.wg.Wait()
}

// 等待，如果超时了返回true，否则返回false
func (w *Wait) WaitWithTimeout(timeout time.Duration) bool {
	c := make(chan struct{}, 1)
	go func() {
		defer close(c)
		w.Wait()
		c <- struct{}{}
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}
