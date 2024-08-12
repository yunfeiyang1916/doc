package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type GormList []string

func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

// 实现sql.Scanner接口
func (g *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), g)
}

type BaseModel struct {
	ID        int32     `gorm:"primarykey;type:int" json:"id"`
	CreatedAt time.Time `gorm:"column:add_time" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:update_time" json:"updated_at"`
	// gorm中的软删除
	DeletedAt gorm.DeletedAt `json:"-"`
	// 自定义一个软删除标记
	IsDeleted bool `gorm:"column:is_deleted" json:"-"`
}
