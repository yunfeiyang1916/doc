package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        int32     `gorm:"primarykey;type:int" json:"id"`
	CreatedAt time.Time `gorm:"column:add_time" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:update_time" json:"updated_at"`
	// gorm中的软删除
	DeletedAt gorm.DeletedAt `json:"-"`
	// 自定义一个软删除标记
	IsDeleted bool `gorm:"column:is_deleted" json:"-"`
}

/*
1. 密文 2. 密文不可反解
 1. 对称加密
 2. 非对称加密
 3. md5 信息摘要算法
    密码如果不可以反解，用户找回密码
*/
type User struct {
	BaseModel
	Mobile   string     `gorm:"index:idx_mobile;unique;type:varchar(11);not null"`
	Password string     `gorm:"type:varchar(100);not null"`
	NickName string     `gorm:"type:varchar(20)"`
	Birthday *time.Time `gorm:"type:datetime"`
	Gender   string     `gorm:"column:gender;default:male;type:varchar(6) comment 'female表示女, male表示男'"`
	Role     int        `gorm:"column:role;default:1;type:int comment '1表示普通用户, 2表示管理员'"`
}
