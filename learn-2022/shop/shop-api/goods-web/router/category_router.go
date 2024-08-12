package router

import (
	"shop/shop-api/goods-web/api"

	"github.com/gin-gonic/gin"
)

func InitCategoryRouter(r *gin.RouterGroup) {
	categoryRouter := r.Group("categorys").Use()
	{
		categoryRouter.GET("", api.GetCategoryList)       // 商品类别列表页
		categoryRouter.DELETE("/:id", api.DeleteCategory) // 删除分类
		categoryRouter.GET("/:id", api.GetCategoryDetail) // 获取分类详情
		categoryRouter.POST("", api.NewCategory)          //新建分类
		categoryRouter.PUT("/:id", api.UpdateCategory)    //修改分类信息
	}
}
