package router

import (
	"shop/shop-api/userop-web/api"

	"github.com/gin-gonic/gin"
)

func InitRouter(router *gin.RouterGroup) {
	r := router.Group("address")
	{
		r.GET("", api.GetAddressList)
		r.POST("", api.NewAddress)
		r.PUT("/:id", api.UpdateAddress)
		r.DELETE("/:id", api.DeleteAddress)
	}
	rMessage := router.Group("message")
	{
		rMessage.GET("", api.GetMessageList)
		rMessage.POST("", api.NewMessage)
	}
	rUserFavs := router.Group("userfavs")
	{
		rUserFavs.GET("", api.GetUserFavList)
		rUserFavs.GET("/:id", api.UserFavDetail)
		rUserFavs.POST("", api.NewUserFav)
		rUserFavs.DELETE("/:id", api.DeleteUserFav)
	}
}
