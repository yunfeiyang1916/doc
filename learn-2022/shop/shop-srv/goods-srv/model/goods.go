package model

import (
	"context"
	"shop/shop-srv/goods-srv/global"
	"strconv"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type Goods struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;not null;" json:"category_id"`
	Category   Category
	BrandsID   int32 `gorm:"type:int;not null;" json:"brands_id"`
	Brands     Brands

	// 是否上架
	OnSale bool `gorm:"default:false;not null" json:"on_sale"`
	// 是否包邮
	ShipFree bool `gorm:"default:false;not null" json:"ship_free"`
	// 是否新品
	IsNew bool `gorm:"default:false;not null" json:"is_new"`
	// 是否热卖
	IsHot bool `gorm:"default:false;not null" json:"is_hot"`

	// 商品名称
	Name string `gorm:"type:varchar(100);not null" json:"name"`
	// 商品货号
	GoodsSn string `gorm:"type:varchar(50);not null" json:"goods_sn"`
	// 点击数
	ClickNum int32 `gorm:"type:int;not null;default:0" json:"click_num"`
	// 销量
	SoldNum int32 `gorm:"type:int;not null;default:0" json:"sold_num"`
	// 收藏数
	FavNum int32 `gorm:"type:int;not null;default:0" json:"fav_num"`
	// 市场价
	MarketPrice float32 `gorm:"type:decimal(10,2);not null;default:0" json:"market_price"`
	// 销售价
	ShopPrice float32 `gorm:"type:decimal(10,2);not null;default:0" json:"shop_price"`
	// 商品简介
	GoodsBrief string `gorm:"type:varchar(100);not null" json:"goods_brief"`
	// 图片集合
	Images GormList `gorm:"type:json;not null" json:"images"`
	// 简介图片集合
	DescImages GormList `gorm:"type:json;not null" json:"desc_images"`
	// 封面图
	GoodsFrontImage string `gorm:"type:varchar(200);not null" json:"goods_front_image"`
}

func (g *Goods) AfterCreate(tx *gorm.DB) (err error) {
	var esGoods EsGoods
	copier.Copy(&esGoods, g)
	id := strconv.Itoa(int(g.ID))
	_, err = global.EsClient.Index().Index(esGoods.GetIndexName()).BodyJson(esGoods).Id(id).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (g *Goods) AfterUpdate(tx *gorm.DB) (err error) {
	var esGoods EsGoods
	copier.Copy(&esGoods, g)
	id := strconv.Itoa(int(g.ID))
	_, err = global.EsClient.Update().Index(esGoods.GetIndexName()).Doc(esGoods).Id(id).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (g *Goods) AfterDelete(tx *gorm.DB) (err error) {
	id := strconv.Itoa(int(g.ID))
	_, err = global.EsClient.Delete().Index(EsGoods{}.GetIndexName()).Id(id).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
