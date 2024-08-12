package consistenthash

import (
	"hash/crc32"
	"sort"
)

//******* 一致性哈希 *******//

type HashFunc func(data []byte) uint32

type Map struct {
	// 哈希函数
	hashFunc HashFunc
	// 各个节点的哈希值，排序后的
	keys []int
	// 以哈希值为键，node节点地址为值的map
	hashMap map[int]string
}

func New(fn HashFunc) *Map {
	m := &Map{
		hashFunc: fn,
		keys:     nil,
		hashMap:  make(map[int]string),
	}
	if m.hashFunc == nil {
		m.hashFunc = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) IsEmpty() bool {
	return len(m.keys) == 0
}

// 将节点加入集群
func (m *Map) AddNode(keys ...string) {
	for _, key := range keys {
		if key == "" {
			continue
		}
		hash := int(m.hashFunc([]byte(key)))
		m.keys = append(m.keys, hash)
		m.hashMap[hash] = key
	}
	// 排序
	sort.Ints(m.keys)
}

// 返回指定值所选的节点
func (m *Map) PickNode(key string) string {
	if m.IsEmpty() {
		return ""
	}
	hash := int(m.hashFunc([]byte(key)))
	// 顺时针选取离key最近的点，使用二分查找，也就是找第一个比key大的值
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	// 说明未找到，则取第一个
	if idx == len(m.keys) {
		idx = 0
	}
	// 取出目标哈希值
	targetHash := m.keys[idx]
	return m.hashMap[targetHash]
}
