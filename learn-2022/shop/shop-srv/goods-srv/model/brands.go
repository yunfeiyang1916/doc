package model

type Brands struct {
	BaseModel
	Name string `gorm:"type:varchar(50);not null" json:"name"`
	Logo string `gorm:"type:varchar(200);not null" json:"logo"`
}
