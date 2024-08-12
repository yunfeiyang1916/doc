package api

import (
	"context"
	"net/http"
	"shop/shop-api/userop-web/forms"
	"shop/shop-api/userop-web/global"
	"shop/shop-srv/userop-srv/proto"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetAddressList(c *gin.Context) {
	request := &proto.AddressRequest{}

	//claims, _ := ctx.Get("claims")
	//currentUser := claims.(*models.CustomClaims)
	//
	//if currentUser.AuthorityId != 2 {
	//	userId, _ := ctx.Get("userId")
	//	request.UserId = int32(userId.(uint))
	//}
	//request.UserId = int32(c.Query("user_id"))
	conn, err := global.UserOpSrvConnPool.Get()
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	defer conn.Close()
	rsp, err := proto.NewAddressClient(conn.Value()).GetAddressList(context.Background(), request)
	if err != nil {
		zap.S().Errorw("获取地址列表失败")
		HandleGrpcErrorToHttp(err, c)
		return
	}

	reMap := gin.H{
		"total": rsp.Total,
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		rMap := make(map[string]interface{})
		rMap["id"] = value.Id
		rMap["user_id"] = value.UserId
		rMap["province"] = value.Province
		rMap["city"] = value.City
		rMap["district"] = value.District
		rMap["address"] = value.Address
		rMap["signer_name"] = value.SignerName
		rMap["signer_mobile"] = value.SignerMobile

		result = append(result, rMap)
	}

	reMap["data"] = result

	c.JSON(http.StatusOK, reMap)
}

func NewAddress(c *gin.Context) {
	addressForm := forms.AddressForm{}
	if err := c.ShouldBindJSON(&addressForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	userId := int32(1) //int32(c.Param("user_id"))

	conn, err := global.UserOpSrvConnPool.Get()
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	defer conn.Close()

	rsp, err := proto.NewAddressClient(conn.Value()).CreateAddress(context.Background(), &proto.AddressRequest{
		UserId:       userId,
		Province:     addressForm.Province,
		City:         addressForm.City,
		District:     addressForm.District,
		Address:      addressForm.Address,
		SignerName:   addressForm.SignerName,
		SignerMobile: addressForm.SignerMobile,
	})

	if err != nil {
		zap.S().Errorw("新建地址失败")
		HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})
}

func DeleteAddress(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	conn, err := global.UserOpSrvConnPool.Get()
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	defer conn.Close()
	_, err = proto.NewAddressClient(conn.Value()).DeleteAddress(context.Background(), &proto.AddressRequest{Id: int32(i)})
	if err != nil {
		zap.S().Errorw("删除地址失败")
		HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}

func UpdateAddress(c *gin.Context) {
	addressForm := forms.AddressForm{}
	if err := c.ShouldBindJSON(&addressForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	conn, err := global.UserOpSrvConnPool.Get()
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	defer conn.Close()
	_, err = proto.NewAddressClient(conn.Value()).UpdateAddress(context.Background(), &proto.AddressRequest{
		Id:           int32(i),
		Province:     addressForm.Province,
		City:         addressForm.City,
		District:     addressForm.District,
		Address:      addressForm.Address,
		SignerName:   addressForm.SignerName,
		SignerMobile: addressForm.SignerMobile,
	})
	if err != nil {
		zap.S().Errorw("更新地址失败")
		HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
