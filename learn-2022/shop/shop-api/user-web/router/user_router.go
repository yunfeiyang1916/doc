package router

import (
	"shop/shop-api/user-web/api"
	"shop/shop-api/user-web/middlewares"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitUserRouter(r *gin.RouterGroup) {
	zap.S().Info("配置用户相关的路由")
	userRouter := r.Group("user")
	userRouter.GET("list", middlewares.JWTAuth() /* middlewares.IsAdminAuth()*/, api.GetUserList)
	userRouter.POST("pwd_login", api.PassWordLogin)
	userRouter.POST("register", api.Register)
}
