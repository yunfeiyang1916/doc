package model

type Banner struct {
	BaseModel
	// 图片地址
	Image string `gorm:"type:varchar(200);not null" json:"image"`
	// 跳转页面url
	Url string `gorm:"type:varchar(200);not null" json:"url"`
	// 轮播图的顺序
	Index int32 `gorm:"type:int;not null" json:"index"`
}
