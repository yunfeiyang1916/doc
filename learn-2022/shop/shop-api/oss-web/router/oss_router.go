package router

import (
	"shop/shop-api/oss-web/api"

	"github.com/gin-gonic/gin"
)

func InitOssRouter(r *gin.RouterGroup) {
	ossRouter := r.Group("oss")
	{
		ossRouter.GET("token", api.Token)
		ossRouter.POST("/callback", api.Callback)
	}
}
