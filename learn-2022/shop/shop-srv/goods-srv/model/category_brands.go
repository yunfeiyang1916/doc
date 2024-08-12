package model

type CategoryBrands struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;index:idx_category_brands,unique"`
	Category   Category
	BrandsID   int32 `gorm:"type:int;index:idx_category_brands,unique"`
	Brands     Brands
}
