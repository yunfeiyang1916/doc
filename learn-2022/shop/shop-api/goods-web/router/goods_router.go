package router

import (
	"shop/shop-api/goods-web/api"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitGoodsRouter(r *gin.RouterGroup) {
	zap.S().Info("配置商品相关的路由")
	router := r.Group("goods")
	router.GET("", api.GetGoodsList)
	router.GET("/:id", api.GoodsDetail)
	router.POST("", api.NewGoods)
	router.PUT("", api.UpdateGoods)
	router.PATCH("/:id", api.UpdateGoodsStatus)
	router.DELETE("/:id", api.DeleteGoods)
}
