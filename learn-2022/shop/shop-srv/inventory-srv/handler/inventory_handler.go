package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"shop/shop-srv/inventory-srv/global"
	"shop/shop-srv/inventory-srv/model"
	"shop/shop-srv/inventory-srv/proto"

	"go.uber.org/zap"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"

	"github.com/go-redsync/redsync/v4"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gorm.io/gorm"

	"github.com/golang/protobuf/ptypes/empty"
)

type InventoryService struct {
	proto.UnsafeInventoryServer
}

// 设置库存
func (i *InventoryService) SetInv(ctx context.Context, info *proto.GoodsInvInfo) (*empty.Empty, error) {
	var obj model.Inventory
	// 如果不存在，则创建
	if err := global.DB.Where(&model.Inventory{Goods: info.GoodsId}).First(&obj).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	obj.Goods = info.GoodsId
	obj.Stocks = info.Num
	if err := global.DB.Save(&obj).Error; err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

// 获取库存信息
func (i *InventoryService) InvDetail(ctx context.Context, info *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var obj model.Inventory
	if err := global.DB.Where(&model.Inventory{Goods: info.GoodsId}).First(&obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "库存不存在")
		}
		return nil, err
	}
	return &proto.GoodsInvInfo{
		GoodsId: obj.Goods,
		Num:     obj.Stocks,
	}, nil
}

// 库存预扣减
func (i *InventoryService) Sell(ctx context.Context, info *proto.SellInfo) (*empty.Empty, error) {
	tx := global.DB.Begin()
	rs := redsync.New(global.RedisPool)
	// todo 应该先查询下是否已经扣减过库存了，要实现幂等性
	// 售卖明细
	sellDetail := model.StockSellDetail{
		OrderSn: info.OrderSn,
		Status:  1,
	}
	var details = make([]model.GoodsDetail, 0, len(info.GoodsInfo))
	for _, v := range info.GoodsInfo {
		details = append(details, model.GoodsDetail{GoodsId: v.GoodsId, Num: v.Num})
		var obj model.Inventory
		// clause.Locking{Strength: "UPDATE"} 实际上会生成select * from table where xx=xx for update
		// 1、因为mysql每个单独的SQL语句都自动提交事务的，如果不在事务中使用for update或者在连接中关闭自动提交事务，for update就不生效，关闭mysql自动提交语句为：set autocommit=0;select @@autocommit;
		// 2、for update只会锁insert/update/delete/for update等执行语句，不会锁纯select等查询语句
		// 3、where 条件如果有索引的话则只会是行锁
		// 4、where 条件没有索引那么会升级成表锁，所以for update要用在where条件有索引的行
		// if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Inventory{Goods: v.GoodsId}).First(&obj).Error; err != nil {

		// 使用redis分布式锁
		key := fmt.Sprintf("goods_%d", v.GoodsId)
		mutex := rs.NewMutex(key)
		if err := mutex.Lock(); err != nil {
			return nil, status.Error(codes.Internal, "获取redis分布式锁异常")
		}
		defer mutex.Unlock()
		if err := tx.Where(&model.Inventory{Goods: v.GoodsId}).First(&obj).Error; err != nil {
			tx.Rollback()
			if err == gorm.ErrRecordNotFound {
				return nil, status.Error(codes.NotFound, "库存不存在")
			}
			return nil, err
		}
		// 判断库存是否充足
		if obj.Stocks < v.Num {
			tx.Rollback()
			return nil, status.Error(codes.ResourceExhausted, "库存不足")
		}
		// 扣减库存
		obj.Stocks -= v.Num
		if err := tx.Save(&obj).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	sellDetail.Detail = details
	if err := tx.Create(&sellDetail).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return &empty.Empty{}, nil
}

// 库存归还: 1、订单超时归还 2、订单创建失败归还 3、订单取消归还
func (i *InventoryService) Reback(ctx context.Context, info *proto.SellInfo) (*empty.Empty, error) {
	tx := global.DB.Begin()
	for _, v := range info.GoodsInfo {
		var obj model.Inventory
		if err := tx.Where(&model.Inventory{Goods: v.GoodsId}).First(&obj).Error; err != nil {
			tx.Rollback()
			if err == gorm.ErrRecordNotFound {
				return nil, status.Error(codes.NotFound, "库存不存在")
			}
			return nil, err
		}
		obj.Stocks += v.Num
		if err := tx.Save(&obj).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return &empty.Empty{}, nil
}

// 消费库存归还消息
func AutoReback(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	// todo 需要实现幂等性
	type OrderInfo struct {
		OrderSn string `json:"order_sn"`
	}
	for _, msg := range msgs {
		var orderInfo OrderInfo
		if err := json.Unmarshal(msg.Body, &orderInfo); err != nil {
			zap.S().Errorf("解析json失败： %v\n", msg.Body)
			// 消息格式不正确，返回消费成功，丢弃该消息
			return consumer.ConsumeSuccess, nil
		}

		var sellDetail model.StockSellDetail
		if err := global.DB.Where(&model.StockSellDetail{OrderSn: orderInfo.OrderSn}).First(&sellDetail).Error; err != nil {
			// 记录不存在，返回消费成功，丢弃该消息
			if err == gorm.ErrRecordNotFound {
				return consumer.ConsumeSuccess, err
			}
			// db有问题，返回稍后重试
			return consumer.ConsumeRetryLater, err
		}
		// 已归还
		if sellDetail.Status == 2 {
			return consumer.ConsumeSuccess, nil
		}
		// 开启事务
		tx := global.DB.Begin()
		for _, v := range sellDetail.Detail {
			if err := tx.Model(&model.Inventory{}).Where(&model.Inventory{Goods: v.GoodsId}).Update("stocks", gorm.Expr("stocks + ?", v.Num)).Error; err != nil {
				tx.Rollback()
				return consumer.ConsumeRetryLater, nil
			}
		}
		// 更新状态
		if err := tx.Model(&model.StockSellDetail{}).Where("order_sn", orderInfo.OrderSn).Update("status", 2).Error; err != nil {
			tx.Rollback()
			return consumer.ConsumeRetryLater, nil
		}
		tx.Commit()
	}
	return consumer.ConsumeSuccess, nil
}
