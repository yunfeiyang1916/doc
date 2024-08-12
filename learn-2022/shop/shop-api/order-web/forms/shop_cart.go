package forms

type ShopCartItemForm struct {
	GoodsId int32 `json:"goods" binding:"required"`
	UserId  int32 `json:"user_id" binding:"required"`
	Nums    int32 `json:"nums" binding:"required,min=1"`
}

type ShopCartItemUpdateForm struct {
	Nums    int32 `json:"nums" binding:"required,min=1"`
	UserId  int32 `json:"user_id" binding:"required"`
	Checked *bool `json:"checked"`
}
