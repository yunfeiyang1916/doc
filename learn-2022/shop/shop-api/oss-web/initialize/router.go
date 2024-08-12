package initialize

import (
	"fmt"
	"net/http"
	"shop/shop-api/oss-web/middlewares"
	"shop/shop-api/oss-web/router"

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

	r.LoadHTMLFiles(fmt.Sprintf("templates/index.html"))
	// 配置静态文件夹路径 第一个参数是api，第二个是文件夹路径
	r.StaticFS("/static", http.Dir(fmt.Sprintf("static")))
	// GET：请求方式；/hello：请求的路径
	// 当客户端以GET方法请求/hello路径时，会执行后面的匿名函数
	r.GET("", func(c *gin.Context) {
		// c.JSON：返回JSON格式的数据
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "posts/index",
		})
	})

	apiGroup := r.Group("/oss/v1")
	router.InitOssRouter(apiGroup)
	return r
}
