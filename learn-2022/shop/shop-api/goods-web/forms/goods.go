package forms

type GoodsForm struct {
	Name string `form:"name" json:"name" binding:"required,min=2,max=100"`
	// 商品编号
	GoodsSn string `form:"goods_sn" json:"goods_sn" binding:"required,min=2,lt=20"`
	// 库存数量
	Stocks     int32 `form:"stocks" json:"stocks" binding:"required,min=1"`
	CategoryId int32 `form:"category" json:"category" binding:"required"`
	// 市场价
	MarketPrice float32 `form:"market_price" json:"market_price" binding:"required,min=0"`
	// 销售价
	ShopPrice float32 `form:"shop_price" json:"shop_price" binding:"required,min=0"`
	// 商品简介
	GoodsBrief string   `form:"goods_brief" json:"goods_brief" binding:"required,min=3"`
	Images     []string `form:"images" json:"images" binding:"required,min=1"`
	DescImages []string `form:"desc_images" json:"desc_images" binding:"required,min=1"`
	// 是否包邮
	ShipFree *bool `form:"ship_free" json:"ship_free" binding:"required"`
	// 正面图片
	FrontImage string `form:"front_image" json:"front_image" binding:"required,url"`
	// 品牌
	Brand int32 `form:"brand" json:"brand" binding:"required"`
}

type GoodsStatusForm struct {
	// 是否新品
	IsNew *bool `form:"new" json:"new" binding:"required"`
	// 是否热销
	IsHot *bool `form:"hot" json:"hot" binding:"required"`
	// 是否在售
	OnSale *bool `form:"sale" json:"sale" binding:"required"`
}
