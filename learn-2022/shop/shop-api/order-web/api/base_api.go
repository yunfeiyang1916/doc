package api

import (
	"net/http"
	"shop/shop-api/order-web/global"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

// 将grpc的code转成http状态码
func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"code": e.Code(), "msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": e.Code(), "msg": "内部错误:" + err.Error(),
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"code": e.Code(), "msg": "参数错误:" + err.Error(),
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": e.Code(), "msg": "grpc服务不可用:" + err.Error(),
				})
			default:
				c.JSON(http.StatusNotFound, gin.H{
					"code": e.Code(), "msg": e.Message(),
				})
			}
		}
	}
}

func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
	return
}
