package api

import (
	"context"
	"net/http"
	"shop/shop-api/userop-web/forms"
	"shop/shop-api/userop-web/global"
	"shop/shop-srv/userop-srv/proto"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetMessageList(c *gin.Context) {
	request := &proto.MessageRequest{}

	//userId, _ := ctx.Get("userId")
	//claims, _ := ctx.Get("claims")
	//model := claims.(*models.CustomClaims)
	//if model.AuthorityId == 1 {
	//	request.UserId = int32(userId.(uint))
	//}
	conn, err := global.UserOpSrvConnPool.Get()
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	defer conn.Close()
	rsp, err := proto.NewMessageClient(conn.Value()).MessageList(context.Background(), request)
	if err != nil {
		zap.S().Errorw("获取留言失败")
		HandleGrpcErrorToHttp(err, c)
		return
	}

	reMap := map[string]interface{}{
		"total": rsp.Total,
	}
	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		rMap := make(map[string]interface{})
		rMap["id"] = value.Id
		rMap["user_id"] = value.UserId
		rMap["type"] = value.MessageType
		rMap["subject"] = value.Subject
		rMap["message"] = value.Message
		rMap["file"] = value.File

		result = append(result, rMap)
	}
	reMap["data"] = result

	c.JSON(http.StatusOK, reMap)
}

func NewMessage(c *gin.Context) {
	userId, _ := c.Get("userId")

	messageForm := forms.MessageForm{}
	if err := c.ShouldBindJSON(&messageForm); err != nil {
		HandleValidatorError(c, err)
		return
	}
	conn, err := global.UserOpSrvConnPool.Get()
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	defer conn.Close()
	rsp, err := proto.NewMessageClient(conn.Value()).CreateMessage(context.Background(), &proto.MessageRequest{
		UserId:      int32(userId.(uint)),
		MessageType: messageForm.MessageType,
		Subject:     messageForm.Subject,
		Message:     messageForm.Message,
		File:        messageForm.File,
	})

	if err != nil {
		zap.S().Errorw("添加留言失败")
		HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})
}
