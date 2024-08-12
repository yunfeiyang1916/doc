package tests

import (
	"context"
	"shop/shop-srv/goods-srv/global"
	"shop/shop-srv/goods-srv/model"
	"shop/shop-srv/goods-srv/proto"
	"strconv"
	"testing"

	"github.com/jinzhu/copier"
)

func TestInitEs(t *testing.T) {
	InitConfig()
	var goodsList []model.Goods
	if err := global.DB.Find(&goodsList).Error; err != nil {
		panic(err)
	}
	for _, goods := range goodsList {
		id := strconv.Itoa(int(goods.ID))
		var esGoods model.EsGoods
		copier.Copy(&esGoods, &goods)
		if _, err := global.EsClient.Index().Index(esGoods.GetIndexName()).BodyJson(esGoods).Id(id).Do(context.Background()); err != nil {
			panic(err)
		}
	}
}

func TestGoodsList(t *testing.T) {
	Init()
	res, err := client.GoodsList(context.Background(), &proto.GoodsFilterRequest{
		Pages:       1,
		PagePerNums: 10,
	})
	if err != nil {
		panic(err)
	}
	t.Logf("%+v", res)
}
