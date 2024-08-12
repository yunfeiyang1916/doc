package model

import "time"

// 购物车
type ShoppingCart struct {
	BaseModel
	// 在购物车列表中我们需要查询当前用户的购物车记录
	User int32 `gorm:"type:int;index"`
	// 加索引：我们需要查询时候， 1. 会影响插入性能 2. 会占用磁盘
	Goods int32 `gorm:"type:int;index"`
	Nums  int32 `gorm:"type:int"`
	// 是否选中
	Checked bool
}

func (ShoppingCart) TableName() string {
	return "shoppingcart"
}

type OrderInfo struct {
	BaseModel

	User int32 `gorm:"type:int;index"`
	// 订单号，我们平台自己生成的订单号
	OrderSn string `gorm:"type:varchar(30);index"`
	// 支付方式
	PayType string `gorm:"type:varchar(20) comment 'alipay(支付宝)， wechat(微信)'"`

	// status大家可以考虑使用iota来做
	Status string `gorm:"type:varchar(20)  comment 'PAYING(待支付), TRADE_SUCCESS(成功)， TRADE_CLOSED(超时关闭), WAIT_BUYER_PAY(交易创建), TRADE_FINISHED(交易结束)'"`
	// 交易号就是支付宝的订单号 查账
	TradeNo string `gorm:"type:varchar(100) comment '交易号'"`
	// 订单金额
	OrderMount float32
	// 支付时间
	PayTime *time.Time `gorm:"type:datetime"`

	// 收货地址
	Address string `gorm:"type:varchar(100)"`
	// 收货人
	SignerName string `gorm:"type:varchar(20)"`
	// 收货人手机
	SingerMobile string `gorm:"type:varchar(11)"`
	// 留言备注
	Post string `gorm:"type:varchar(20)"`
}

func (OrderInfo) TableName() string {
	return "orderinfo"
}

type OrderGoods struct {
	BaseModel
	// 订单id
	Order int32 `gorm:"type:int;index"`
	// 商品id
	Goods int32 `gorm:"type:int;index"`

	// 把商品的信息保存下来了 ， 字段冗余， 高并发系统中我们一般都不会遵循三范式  做镜像 记录
	GoodsName  string `gorm:"type:varchar(100);index"`
	GoodsImage string `gorm:"type:varchar(200)"`
	GoodsPrice float32
	Nums       int32 `gorm:"type:int"`
}

func (OrderGoods) TableName() string {
	return "ordergoods"
}
