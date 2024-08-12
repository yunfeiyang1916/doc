package dict

import "sync"

type SyncDict struct {
	m sync.Map
}

func MakeSyncDict() *SyncDict {
	return &SyncDict{}
}

func (dict *SyncDict) Get(key string) (val interface{}, exists bool) {
	return dict.m.Load(key)
}

func (dict *SyncDict) Len() int {
	length := 0
	dict.m.Range(func(key, value any) bool {
		length++
		return true
	})
	return length
}

func (dict *SyncDict) Put(key string, val interface{}) (result int) {
	_, ok := dict.m.Load(key)
	dict.m.Store(key, val)
	// 说明是更新
	if ok {
		return 0
	}
	return 1
}

// 如果不存在才设置
func (dict *SyncDict) PutIfAbsent(key string, val interface{}) (result int) {
	_, ok := dict.m.Load(key)
	// 已存在，直接返回
	if ok {
		return 0
	}
	dict.m.Store(key, val)
	return 1
}

// 如果存在才设置
func (dict *SyncDict) PutIfExists(key string, val interface{}) (result int) {
	_, ok := dict.m.Load(key)
	// 已存在才更改
	if ok {
		dict.m.Store(key, val)
		return 1
	}

	return 0
}

func (dict *SyncDict) Remove(key string) (val interface{}, result int) {
	val, ok := dict.m.Load(key)
	dict.m.Delete(key)
	if ok {
		return val, 1
	}
	return nil, 0
}

func (dict *SyncDict) ForEach(consumer Consumer) {
	if consumer != nil {
		dict.m.Range(func(key, value any) bool {
			consumer(key.(string), value)
			return true
		})
	}
}

func (dict *SyncDict) Keys() []string {
	result := make([]string, 0)
	dict.m.Range(func(key, value any) bool {
		result = append(result, key.(string))
		return true
	})
	return result
}

// 返回随机键，键可能重复
func (dict *SyncDict) RandomKeys(limit int) []string {
	result := make([]string, 0)
	for i := 0; i < limit; i++ {
		dict.m.Range(func(key, value any) bool {
			result = append(result, key.(string))
			// 退出sync.map的Range循环
			return false
		})
	}
	return result
}

func (dict *SyncDict) RandomDistinctKeys(limit int) []string {
	result := make([]string, 0)
	i := 0
	dict.m.Range(func(key, value any) bool {
		result = append(result, key.(string))
		i++
		if i == limit {
			// 已经循环limit次了，退出循环
			return false
		}
		return true
	})
	return result
}

func (dict *SyncDict) Clear() {
	*dict = *MakeSyncDict()
}
