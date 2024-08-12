package model

// 库存
type Inventory struct {
	BaseModel
	// 产品id
	Goods int32 `gorm:"type:int;index"`
	// 库存
	Stocks int32 `gorm:"type:int"`
	// 分布式锁的乐观锁
	Version int32 `gorm:"type:int"`
}
