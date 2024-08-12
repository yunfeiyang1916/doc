package router

import (
	"shop/shop-api/user-web/api"

	"github.com/gin-gonic/gin"
)

func InitBaseRouter(r *gin.RouterGroup) {
	baseRouter := r.Group("base")
	baseRouter.GET("captcha", api.GetCaptcha)
	baseRouter.POST("send_sms", api.SendSms)
}
