package dict

// Consumer is used to traversal dict, if it returns false the traversal will be break
type Consumer func(key string, val interface{}) bool

// Dict is interface of a key-value data structure
type Dict interface {
	Get(key string) (val interface{}, exists bool)
	Len() int
	Put(key string, val interface{}) (result int)
	// 如果不存在才设置
	PutIfAbsent(key string, val interface{}) (result int)
	// 如果存在才设置
	PutIfExists(key string, val interface{}) (result int)
	Remove(key string) (val interface{}, result int)
	ForEach(consumer Consumer)
	Keys() []string
	// 返回随机键，键可能重复
	RandomKeys(limit int) []string
	// 返回随机不重复的键
	RandomDistinctKeys(limit int) []string
	Clear()
}
