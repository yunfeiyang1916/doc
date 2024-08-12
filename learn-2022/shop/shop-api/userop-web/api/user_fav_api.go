package api

import (
	"context"
	"net/http"
	"shop/shop-api/userop-web/forms"
	"shop/shop-api/userop-web/global"
	goodsProto "shop/shop-srv/goods-srv/proto"
	"shop/shop-srv/userop-srv/proto"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetUserFavList(c *gin.Context) {
	userId, _ := c.Get("userId")
	conn, err := global.UserOpSrvConnPool.Get()
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	defer conn.Close()
	userFavRsp, err := proto.NewUserFavClient(conn.Value()).GetFavList(context.Background(), &proto.UserFavRequest{
		UserId: int32(userId.(uint)),
	})
	if err != nil {
		zap.S().Errorw("获取收藏列表失败")
		HandleGrpcErrorToHttp(err, c)
		return
	}

	ids := make([]int32, 0)
	for _, item := range userFavRsp.Data {
		ids = append(ids, item.GoodsId)
	}

	if len(ids) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"total": 0,
		})
		return
	}
	goodsConn, err := global.GoodsSrvConnPool.Get()
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	defer goodsConn.Close()
	//请求商品服务
	goods, err := goodsProto.NewGoodsClient(conn.Value()).BatchGetGoods(context.Background(), &goodsProto.BatchGoodsIdInfo{
		Id: ids,
	})
	if err != nil {
		zap.S().Errorw("[List] 批量查询【商品列表】失败")
		HandleGrpcErrorToHttp(err, c)
		return
	}

	reMap := map[string]interface{}{
		"total": userFavRsp.Total,
	}

	goodsList := make([]interface{}, 0)
	for _, item := range userFavRsp.Data {
		data := gin.H{
			"id": item.GoodsId,
		}

		for _, good := range goods.Data {
			if item.GoodsId == good.Id {
				data["name"] = good.Name
				data["shop_price"] = good.ShopPrice
			}
		}

		goodsList = append(goodsList, data)
	}
	reMap["data"] = goodsList
	c.JSON(http.StatusOK, reMap)
}

func NewUserFav(c *gin.Context) {
	userFavForm := forms.UserFavForm{}
	if err := c.ShouldBindJSON(&userFavForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	//缺少一步， 这个时候应该去商品服务查询一下这个是否存在
	userId, _ := c.Get("userId")
	conn, err := global.UserOpSrvConnPool.Get()
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	defer conn.Close()
	_, err = proto.NewUserFavClient(conn.Value()).AddUserFav(context.Background(), &proto.UserFavRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: userFavForm.GoodsId,
	})

	if err != nil {
		zap.S().Errorw("添加收藏记录失败")
		HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func DeleteUserFav(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	userId, _ := c.Get("userId")
	conn, err := global.UserOpSrvConnPool.Get()
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	defer conn.Close()
	_, err = proto.NewUserFavClient(conn.Value()).DeleteUserFav(context.Background(), &proto.UserFavRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(i),
	})
	if err != nil {
		zap.S().Errorw("删除收藏记录失败")
		HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}

func UserFavDetail(c *gin.Context) {
	goodsId := c.Param("id")
	goodsIdInt, err := strconv.ParseInt(goodsId, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	userId, _ := c.Get("userId")
	conn, err := global.UserOpSrvConnPool.Get()
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	defer conn.Close()
	_, err = proto.NewUserFavClient(conn.Value()).GetUserFavDetail(context.Background(), &proto.UserFavRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(goodsIdInt),
	})
	if err != nil {
		zap.S().Errorw("查询收藏状态失败")
		HandleGrpcErrorToHttp(err, c)
		return
	}

	c.Status(http.StatusOK)
}
