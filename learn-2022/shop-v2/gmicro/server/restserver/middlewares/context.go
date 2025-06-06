package middlewares

import "github.com/gin-gonic/gin"

const (
	UsernameKey = "username"
	KeyUserID   = "userid"
	UserIP      = "ip"
)

// 为每个请求添加上下文
func Context() gin.HandlerFunc {
	return func(c *gin.Context) {
		//从c中获取到ip地址
		ip := c.ClientIP()
		//ip := c.GetHeader("X-Real-IP")
		//TODO 自己扩展
		c.Set(UserIP, ip)
		c.Next()
	}
}
