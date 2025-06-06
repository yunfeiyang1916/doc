package model

import (
	"database/sql/driver"
	"encoding/json"
)

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

// 库存售卖明细
type StockSellDetail struct {
	OrderSn string          `gorm:"type:varchar(200);index:idx_order_sn,unique;"`
	Status  int32           `gorm:"type:varchar(200)"` //1 表示已扣减 2. 表示已归还
	Detail  GoodsDetailList `gorm:"type:varchar(2000)"`
}

type GoodsDetailList []GoodsDetail

func (g GoodsDetailList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (g *GoodsDetailList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), g)
}

// 商品售卖详情
type GoodsDetail struct {
	GoodsId int32
	Num     int32
}
