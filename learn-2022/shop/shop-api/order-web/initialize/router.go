package initialize

import (
	"net/http"
	"shop/shop-api/order-web/global"
	"shop/shop-api/order-web/middlewares"
	"shop/shop-api/order-web/router"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	r := gin.Default()
	// 配置跨域
	r.Use(middlewares.Cors(), otelgin.Middleware(global.ServerConfig.Name))
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})

	apiGroup := r.Group("/o/v1")
	router.InitOrderRouter(apiGroup)
	router.InitShopCartRouter(apiGroup)
	return r
}
