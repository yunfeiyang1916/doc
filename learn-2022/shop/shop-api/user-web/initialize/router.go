package initialize

import (
	"net/http"
	"shop/shop-api/user-web/middlewares"
	"shop/shop-api/user-web/router"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	r := gin.Default()
	// 配置跨域
	r.Use(middlewares.Cors())
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})

	apiGroup := r.Group("/v1")
	router.InitUserRouter(apiGroup)
	router.InitBaseRouter(apiGroup)
	return r
}
