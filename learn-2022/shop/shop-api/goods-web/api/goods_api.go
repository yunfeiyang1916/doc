package api

import (
	"net/http"
	"shop/shop-api/goods-web/forms"
	"shop/shop-api/goods-web/global"
	"shop/shop-srv/goods-srv/proto"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetGoodsList(c *gin.Context) {
	zap.S().Info("获取商品列表页")
	request := &proto.GoodsFilterRequest{}
	// 最低价格
	priceMin := c.DefaultQuery("pmin", "0")
	priceMinInt, _ := strconv.Atoi(priceMin)
	request.PriceMin = int32(priceMinInt)
	// 最高价格
	priceMax := c.DefaultQuery("pmax", "0")
	priceMaxInt, _ := strconv.Atoi(priceMax)
	request.PriceMax = int32(priceMaxInt)
	// 是否热销
	isHot := c.DefaultQuery("ih", "0")
	if isHot == "1" {
		request.IsHot = true
	}
	// 是否新品
	isNew := c.DefaultQuery("in", "0")
	if isNew == "1" {
		request.IsNew = true
	}
	// 是否tab栏
	isTab := c.DefaultQuery("it", "0")
	if isTab == "1" {
		request.IsTab = true
	}
	// 分类
	categoryId := c.DefaultQuery("c", "0")
	categoryIdInt, _ := strconv.Atoi(categoryId)
	request.TopCategory = int32(categoryIdInt)
	// 页码
	pages := c.DefaultQuery("p", "0")
	pagesInt, _ := strconv.Atoi(pages)
	request.Pages = int32(pagesInt)
	// 每页数量
	perNums := c.DefaultQuery("pnum", "0")
	perNumsInt, _ := strconv.Atoi(perNums)
	request.PagePerNums = int32(perNumsInt)
	// 关键词
	keywords := c.DefaultQuery("q", "")
	request.KeyWords = keywords
	// 品牌
	brandId := c.DefaultQuery("b", "0")
	brandIdInt, _ := strconv.Atoi(brandId)
	request.Brand = int32(brandIdInt)

	//r, err := global.GoodsSrvClient.GoodsList(c.Request.Context(), request)
	//if err != nil {
	//	zap.S().Errorw("[List] 查询 【商品列表】失败")
	//	HandleGrpcErrorToHttp(err, c)
	//	return
	//}
	conn, err := global.GoodsSrvConnPool.Get()
	if err != nil {
		zap.S().Errorw("[List] 查询 【商品列表】失败")
		HandleGrpcErrorToHttp(err, c)
		return
	}
	defer conn.Close()
	r, err := proto.NewGoodsClient(conn.Value()).GoodsList(c.Request.Context(), request)
	if err != nil {
		zap.S().Errorw("[List] 查询 【商品列表】失败")
		HandleGrpcErrorToHttp(err, c)
		return
	}

	reMap := map[string]interface{}{
		"total": r.Total,
	}

	goodsList := make([]interface{}, 0)
	for _, value := range r.Data {
		m := map[string]interface{}{
			"id":          value.Id,
			"name":        value.Name,
			"goods_brief": value.GoodsBrief,
			"desc":        value.GoodsDesc,
			"ship_free":   value.ShipFree,
			"images":      value.Images,
			"desc_images": value.DescImages,
			"front_image": value.GoodsFrontImage,
			"shop_price":  value.ShopPrice,
			"is_hot":      value.IsHot,
			"is_new":      value.IsNew,
			"on_sale":     value.OnSale,
		}
		if value.Brands != nil {
			m["brands"] = map[string]interface{}{
				"id":   value.Brands.Id,
				"name": value.Brands.Name,
				"logo": value.Brands.Logo,
			}
		}
		if value.Category != nil {
			m["category"] = map[string]interface{}{
				"id":   value.Category.Id,
				"name": value.Category.Name,
			}
		}
		goodsList = append(goodsList, m)
	}
	reMap["data"] = goodsList

	c.JSON(http.StatusOK, reMap)
}

func NewGoods(c *gin.Context) {
	goodsForm := forms.GoodsForm{}
	if err := c.ShouldBindJSON(&goodsForm); err != nil {
		HandleValidatorError(c, err)
		return
	}
	conn, err := global.GoodsSrvConnPool.Get()
	if err != nil {
		zap.S().Errorw("[New] 创建 【商品】失败")
		HandleGrpcErrorToHttp(err, c)
		return
	}
	defer conn.Close()
	r, err := proto.NewGoodsClient(conn.Value()).CreateGoods(c.Request.Context(), &proto.CreateGoodsInfo{
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		ShipFree:        *goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	})
	if err != nil {
		zap.S().Errorw("[New] 创建 【商品】失败")
		HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, r)
}

func GoodsDetail(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	conn, err := global.GoodsSrvConnPool.Get()
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	defer conn.Close()
	r, err := proto.NewGoodsClient(conn.Value()).GetGoodsDetail(c.Request.Context(), &proto.GoodInfoRequest{Id: int32(i)})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	rsp := map[string]interface{}{
		"id":          r.Id,
		"name":        r.Name,
		"goods_brief": r.GoodsBrief,
		"desc":        r.GoodsDesc,
		"ship_free":   r.ShipFree,
		"images":      r.Images,
		"desc_images": r.DescImages,
		"front_image": r.GoodsFrontImage,
		"shop_price":  r.ShopPrice,
		"ctegory": map[string]interface{}{
			"id":   r.Category.Id,
			"name": r.Category.Name,
		},

		"is_hot":  r.IsHot,
		"is_new":  r.IsNew,
		"on_sale": r.OnSale,
	}
	if r.Brands != nil {
		rsp["brand"] = map[string]interface{}{
			"id":   r.Brands.Id,
			"name": r.Brands.Name,
			"logo": r.Brands.Logo,
		}
	}
	if r.Category != nil {
		rsp["ctegory"] = map[string]interface{}{
			"id":   r.Category.Id,
			"name": r.Category.Name,
		}
	}
	c.JSON(http.StatusOK, rsp)
}

func DeleteGoods(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	conn, err := global.GoodsSrvConnPool.Get()
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	defer conn.Close()
	_, err = proto.NewGoodsClient(conn.Value()).DeleteGoods(c.Request.Context(), &proto.DeleteGoodsInfo{Id: int32(i)})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}

	c.Status(http.StatusOK)
	return
}

func UpdateGoods(c *gin.Context) {
	goodsForm := forms.GoodsForm{}
	if err := c.ShouldBindJSON(&goodsForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)

	conn, err := global.GoodsSrvConnPool.Get()
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	defer conn.Close()
	if _, err = proto.NewGoodsClient(conn.Value()).UpdateGoods(c.Request.Context(), &proto.CreateGoodsInfo{
		Id:              int32(i),
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		ShipFree:        *goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	}); err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}

func UpdateGoodsStatus(c *gin.Context) {
	goodsStatusForm := forms.GoodsStatusForm{}
	if err := c.ShouldBindJSON(&goodsStatusForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	conn, err := global.GoodsSrvConnPool.Get()
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	defer conn.Close()
	if _, err = proto.NewGoodsClient(conn.Value()).UpdateGoods(c.Request.Context(), &proto.CreateGoodsInfo{
		Id:     int32(i),
		IsHot:  *goodsStatusForm.IsHot,
		IsNew:  *goodsStatusForm.IsNew,
		OnSale: *goodsStatusForm.OnSale,
	}); err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "修改成功",
	})
}
