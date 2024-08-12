package redis

type Connection interface {
	Write([]byte) (int, error)
	GetDBIndex() int
	// 切换DB
	SelectDB(int)
}
