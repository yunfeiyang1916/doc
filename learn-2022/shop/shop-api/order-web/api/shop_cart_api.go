package api

import (
	"context"
	"net/http"
	"shop/shop-api/order-web/forms"
	"shop/shop-api/order-web/global"
	goodsProto "shop/shop-srv/goods-srv/proto"
	inventoryProto "shop/shop-srv/inventory-srv/proto"
	"shop/shop-srv/order-srv/proto"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetCartList(c *gin.Context) {
	ctx := c.Request.Context()
	//获取购物车商品
	uidStr := c.Query("uid")
	uid, err := strconv.ParseUint(uidStr, 10, 64)
	if err != nil {
		uid = 1
	}
	conn, err := global.OrderSrvConnPool.Get()
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	defer conn.Close()

	rsp, err := proto.NewOrderClient(conn.Value()).CartItemList(ctx, &proto.UserInfo{
		Id: int32(uid),
	})
	if err != nil {
		zap.S().Errorw("[List] 查询 【购物车列表】失败")
		HandleGrpcErrorToHttp(err, c)
		return
	}

	ids := make([]int32, 0)
	for _, item := range rsp.Data {
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

	//请求商品服务获取商品信息
	goodsRsp, err := goodsProto.NewGoodsClient(goodsConn.Value()).BatchGetGoods(ctx, &goodsProto.BatchGoodsIdInfo{
		Id: ids,
	})
	if err != nil {
		zap.S().Errorw("[List] 批量查询【商品列表】失败")
		HandleGrpcErrorToHttp(err, c)
		return
	}

	reMap := gin.H{
		"total": rsp.Total,
	}

	goodsList := make([]interface{}, 0)
	for _, item := range rsp.Data {
		for _, good := range goodsRsp.Data {
			if good.Id == item.GoodsId {
				tmpMap := map[string]interface{}{}
				tmpMap["id"] = item.Id
				tmpMap["goods_id"] = item.GoodsId
				tmpMap["good_name"] = good.Name
				tmpMap["good_image"] = good.GoodsFrontImage
				tmpMap["good_price"] = good.ShopPrice
				tmpMap["nums"] = item.Nums
				tmpMap["checked"] = item.Checked

				goodsList = append(goodsList, tmpMap)
			}
		}
	}
	reMap["data"] = goodsList
	c.JSON(http.StatusOK, reMap)
}

func NewCart(c *gin.Context) {
	//添加商品到购物车
	itemForm := forms.ShopCartItemForm{}
	if err := c.ShouldBindJSON(&itemForm); err != nil {
		HandleValidatorError(c, err)
		return
	}
	goodsConn, err := global.GoodsSrvConnPool.Get()
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	defer goodsConn.Close()
	//为了严谨性，添加商品到购物车之前，记得检查一下商品是否存在
	_, err = goodsProto.NewGoodsClient(goodsConn.Value()).GetGoodsDetail(context.Background(), &goodsProto.GoodInfoRequest{
		Id: itemForm.GoodsId,
	})
	if err != nil {
		zap.S().Errorw("[List] 查询【商品信息】失败")
		HandleGrpcErrorToHttp(err, c)
		return
	}

	inventoryConn, err := global.InventorySrvConnPool.Get()
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	defer inventoryConn.Close()
	//如果现在添加到购物车的数量和库存的数量不一致
	invRsp, err := inventoryProto.NewInventoryClient(inventoryConn.Value()).InvDetail(context.Background(), &inventoryProto.GoodsInvInfo{
		GoodsId: itemForm.GoodsId,
	})
	if err != nil {
		zap.S().Errorw("[List] 查询【库存信息】失败")
		HandleGrpcErrorToHttp(err, c)
		return
	}
	if invRsp.Num < itemForm.Nums {
		c.JSON(http.StatusBadRequest, gin.H{
			"nums": "库存不足",
		})
		return
	}
	conn, err := global.OrderSrvConnPool.Get()
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	defer conn.Close()
	rsp, err := proto.NewOrderClient(conn.Value()).CreateCartItem(context.Background(), &proto.CartItemRequest{
		GoodsId: itemForm.GoodsId,
		UserId:  itemForm.UserId,
		Nums:    itemForm.Nums,
	})

	if err != nil {
		zap.S().Errorw("添加到购物车失败")
		HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})
}

func UpdateCart(c *gin.Context) {
	// o/v1/421
	// id为商品id
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "url格式出错",
		})
		return
	}

	itemForm := forms.ShopCartItemUpdateForm{}
	if err = c.ShouldBindJSON(&itemForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	request := proto.CartItemRequest{
		UserId:  itemForm.UserId,
		GoodsId: int32(i),
		Nums:    itemForm.Nums,
		Checked: false,
	}
	if itemForm.Checked != nil {
		request.Checked = *itemForm.Checked
	}
	conn, err := global.OrderSrvConnPool.Get()
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	defer conn.Close()
	_, err = proto.NewOrderClient(conn.Value()).UpdateCartItem(context.Background(), &request)
	if err != nil {
		zap.S().Errorw("更新购物车记录失败")
		HandleGrpcErrorToHttp(err, c)
		return
	}
	c.Status(http.StatusOK)
}

func DeleteCart(c *gin.Context) {
	// id为商品id
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "url格式出错",
		})
		return
	}
	userIdStr := c.Query("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "user_id参数错误",
		})
		return
	}
	conn, err := global.OrderSrvConnPool.Get()
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	defer conn.Close()
	_, err = proto.NewOrderClient(conn.Value()).DeleteCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  int32(userId),
		GoodsId: int32(i),
	})
	if err != nil {
		zap.S().Errorw("删除购物车记录失败")
		HandleGrpcErrorToHttp(err, c)
		return
	}
	c.Status(http.StatusOK)
}
