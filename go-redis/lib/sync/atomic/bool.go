package atomic

import "sync/atomic"

// 布尔原子类型，因为atomic包下没有布尔类型，所以定义一个
type Boolean uint32

// Get reads the value atomically
func (b *Boolean) Get() bool {
	return atomic.LoadUint32((*uint32)(b)) != 0
}

// Set writes the value atomically
func (b *Boolean) Set(v bool) {
	if v {
		atomic.StoreUint32((*uint32)(b), 1)
	} else {
		atomic.StoreUint32((*uint32)(b), 0)
	}
}
