package middlewares

import (
	"net/http"
	"shop/shop-api/goods-web/models"

	"github.com/gin-gonic/gin"
)

func IsAdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get("claims")
		currentUser := claims.(*models.CustomClaims)
		// 2是管理员
		if currentUser.AuthorityId != 2 {
			c.JSON(http.StatusForbidden, gin.H{"msg": "无权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}
