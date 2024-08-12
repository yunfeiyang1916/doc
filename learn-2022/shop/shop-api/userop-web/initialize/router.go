package initialize

import (
	"net/http"
	"shop/shop-api/userop-web/middlewares"
	"shop/shop-api/userop-web/router"

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

	apiGroup := r.Group("/up/v1")
	router.InitRouter(apiGroup)
	return r
}
