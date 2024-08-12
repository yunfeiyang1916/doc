package handler

import (
	"context"
	"fmt"
	"shop/shop-srv/inventory-srv/global"
	"shop/shop-srv/inventory-srv/model"
	"shop/shop-srv/inventory-srv/proto"

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
	for _, v := range info.GoodsInfo {
		var obj model.Inventory
		// clause.Locking{Strength: "UPDATE"} 实际上会生成select * from table where xx=xx for update
		// 1、因为mysql每个单独的SQL语句都自动提交事务的，如果不在事务中使用for update或者在连接中关闭自动提交事务，for update就不生效，关闭mysql自动提交语句为：set autocommit=0;select @@autocommit;
		// 2、for update只会锁insert/update/delete/for update等执行语句，不会锁纯select等查询语句
		// 3、where 条件如果有索引的话则只会是行锁
		// 4、where 条件没有索引那么会升级成表锁，所以for update要用在where 条件有索引的行
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
