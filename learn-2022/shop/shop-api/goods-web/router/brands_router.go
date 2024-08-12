package router

import (
	"shop/shop-api/goods-web/api"

	"github.com/gin-gonic/gin"
)

func InitBrandRouter(r *gin.RouterGroup) {
	brandRouter := r.Group("brands")
	{
		brandRouter.GET("", api.BrandList)          // 品牌列表页
		brandRouter.DELETE("/:id", api.DeleteBrand) // 删除品牌
		brandRouter.POST("", api.NewBrand)          //新建品牌
		brandRouter.PUT("/:id", api.UpdateBrand)    //修改品牌信息
	}

	categoryBrandRouter := r.Group("categorybrands")
	{
		categoryBrandRouter.GET("", api.CategoryBrandList)          // 类别品牌列表页
		categoryBrandRouter.DELETE("/:id", api.DeleteCategoryBrand) // 删除类别品牌
		categoryBrandRouter.POST("", api.NewCategoryBrand)          //新建类别品牌
		categoryBrandRouter.PUT("/:id", api.UpdateCategoryBrand)    //修改类别品牌
		categoryBrandRouter.GET("/:id", api.GetCategoryBrandList)   //获取分类的品牌
	}
}
