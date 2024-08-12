package router

import (
	"shop/shop-api/order-web/api"

	"github.com/gin-gonic/gin"
)

func InitShopCartRouter(router *gin.RouterGroup) {
	r := router.Group("shopcarts")
	{
		// 购物车列表
		r.GET("", api.GetCartList)
		// 添加商品到购物车
		r.POST("", api.NewCart)
		// 删除条目
		r.DELETE("/:id", api.DeleteCart)
		// 修改条目
		r.PATCH("/:id", api.UpdateCart)
	}
}
