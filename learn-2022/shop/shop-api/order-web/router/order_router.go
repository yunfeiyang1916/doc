package router

import (
	"shop/shop-api/order-web/api"

	"github.com/gin-gonic/gin"
)

func InitOrderRouter(router *gin.RouterGroup) {
	r := router.Group("orders")
	{
		// 订单列表
		r.GET("", api.GetOrderList)
		// 订单详情
		r.GET("/:id", api.OrderDetail)
		// 新建订单
		r.POST("", api.NewOrder)
	}

	payRouter := router.Group("pay")
	payRouter.POST("alipay/notify", api.Notify)
}
