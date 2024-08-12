package router

import (
	"shop/shop-api/goods-web/api"

	"github.com/gin-gonic/gin"
)

func InitBannerRouter(Router *gin.RouterGroup) {
	bannerRouter := Router.Group("banners")
	{
		bannerRouter.GET("", api.GetBannerList)       // 轮播图列表页
		bannerRouter.DELETE("/:id", api.DeleteBanner) // 删除轮播图
		bannerRouter.POST("", api.NewBanner)          //新建轮播图
		bannerRouter.PUT("/:id", api.UpdateBanner)    //修改轮播图信息
	}
}
