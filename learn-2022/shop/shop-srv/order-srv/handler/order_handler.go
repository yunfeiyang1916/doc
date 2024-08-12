package handler

import (
	"context"
	"fmt"
	"math/rand"
	goodsProto "shop/shop-srv/goods-srv/proto"
	inventoryProto "shop/shop-srv/inventory-srv/proto"
	"shop/shop-srv/order-srv/global"
	"shop/shop-srv/order-srv/model"
	"shop/shop-srv/order-srv/proto"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gorm.io/gorm"

	"github.com/jinzhu/copier"

	"github.com/golang/protobuf/ptypes/empty"
)

type OrderService struct {
	proto.UnimplementedOrderServer
}

// 获取用户的购物车信息
func (o *OrderService) CartItemList(ctx context.Context, info *proto.UserInfo) (*proto.CartItemListResponse, error) {
	var (
		list []model.ShoppingCart
		rsp  proto.CartItemListResponse
	)
	if r := global.DB.Where(&model.ShoppingCart{User: info.Id}).Find(&list); r.Error != nil {
		return nil, r.Error
	} else {
		rsp.Total = int32(len(list))
	}
	for _, shopCart := range list {
		var to proto.ShopCartInfoResponse
		copier.Copy(&to, &shopCart)
		rsp.Data = append(rsp.Data, &to)
	}
	return &rsp, nil
}

// 添加商品到购物车
func (o *OrderService) CreateCartItem(ctx context.Context, request *proto.CartItemRequest) (*proto.ShopCartInfoResponse, error) {
	// 将商品添加到购物车 1. 购物车中原本没有这件商品 - 新建一个记录 2. 这个商品之前添加到了购物车- 合并
	var obj model.ShoppingCart
	if err := global.DB.Where(&model.ShoppingCart{User: request.UserId, Goods: request.GoodsId}).First(&obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			obj = model.ShoppingCart{
				User:  request.UserId,
				Goods: request.GoodsId,
				Nums:  request.Nums,
			}
		} else {
			return nil, err
		}
	} else {
		//如果记录已经存在，则合并购物车记录, 更新操作
		obj.Nums += request.Nums
	}
	if err := global.DB.Save(&obj).Error; err != nil {
		return nil, err
	}
	return &proto.ShopCartInfoResponse{
		Id:      obj.ID,
		UserId:  obj.User,
		GoodsId: obj.Goods,
		Nums:    obj.Nums,
	}, nil
}

// 修改购物车信息
func (o *OrderService) UpdateCartItem(ctx context.Context, request *proto.CartItemRequest) (*empty.Empty, error) {
	// 更新购物车记录，更新数量和选中状态
	var obj model.ShoppingCart
	if err := global.DB.Where(&model.ShoppingCart{User: request.UserId, Goods: request.GoodsId}).First(&obj).Error; err != nil {
		return nil, err
	}
	if request.Nums > 0 {
		obj.Nums = request.Nums
	}
	obj.Checked = request.Checked
	if err := global.DB.Save(&obj).Error; err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

// 删除购物车条目
func (o *OrderService) DeleteCartItem(ctx context.Context, request *proto.CartItemRequest) (*empty.Empty, error) {
	if err := global.DB.Where(&model.ShoppingCart{User: request.UserId, Goods: request.GoodsId}).Delete(&model.ShoppingCart{}).Error; err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

// 创建订单
func (o *OrderService) CreateOrder(ctx context.Context, request *proto.OrderRequest) (*proto.OrderInfoResponse, error) {
	/*
		新建订单
			1. 从购物车中获取到选中的商品
			2. 商品的价格自己查询 - 访问商品服务 (跨微服务)
			3. 库存的扣减 - 访问库存服务 (跨微服务)
			4. 订单的基本信息表 - 订单的商品信息表
			5. 从购物车中删除已购买的记录
	*/
	var (
		goodsIds    []int32
		shopCharts  []*model.ShoppingCart
		goodsNumMap = make(map[int32]int32)
	)
	if err := global.DB.Where(&model.ShoppingCart{User: request.UserId, Checked: true}).Find(&shopCharts).Error; err != nil {
		return nil, err
	}
	if len(shopCharts) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "没有选中的商品")
	}
	for _, v := range shopCharts {
		goodsIds = append(goodsIds, v.Goods)
		goodsNumMap[v.Goods] = v.Nums
	}
	// 跨服务调用
	conn, err := global.GoodsSrvConnPool.Get()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	goodsSrvClient := goodsProto.NewGoodsClient(conn.Value())
	goodsResp, err := goodsSrvClient.BatchGetGoods(ctx, &goodsProto.BatchGoodsIdInfo{Id: goodsIds})
	if err != nil {
		return nil, err
	}

	var (
		// 订单总价
		orderAmount   float32
		orderGoods    []*model.OrderGoods
		goodsInvInfos []*inventoryProto.GoodsInvInfo
	)
	for _, v := range goodsResp.Data {
		num, ok := goodsNumMap[v.Id]
		if !ok {
			continue
		}
		orderAmount += v.ShopPrice * float32(num)
		orderGoods = append(orderGoods, &model.OrderGoods{
			Goods:      v.Id,
			GoodsName:  v.Name,
			GoodsImage: v.GoodsFrontImage,
			GoodsPrice: v.ShopPrice,
			Nums:       num,
		})
		goodsInvInfos = append(goodsInvInfos, &inventoryProto.GoodsInvInfo{
			GoodsId: v.Id,
			Num:     num,
		})
	}

	// 调用库存服务，进行库存扣减
	inventoryConn, err := global.InventorySrvConnPool.Get()
	if err != nil {
		return nil, err
	}
	defer inventoryConn.Close()
	inventoryClient := inventoryProto.NewInventoryClient(inventoryConn.Value())
	if _, err = inventoryClient.Sell(ctx, &inventoryProto.SellInfo{GoodsInfo: goodsInvInfos}); err != nil {
		return nil, err
	}

	// 生成订单表
	order := &model.OrderInfo{
		OrderSn:      GenerateOrderSn(request.UserId),
		OrderMount:   orderAmount,
		Address:      request.Address,
		SignerName:   request.Name,
		SingerMobile: request.Mobile,
		Post:         request.Post,
		User:         request.UserId,
	}
	// 开启事务
	if tErr := global.DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.Create(&order).Error; err != nil {
			return err
		}
		// 赋值订单id
		for _, v := range orderGoods {
			v.Order = order.ID
		}
		// 批量插入订单商品表
		if err = tx.Create(&orderGoods).Error; err != nil {
			return err
		}
		// 从购物车中删除已购买的记录
		if err = tx.Where(&model.ShoppingCart{User: request.UserId, Checked: true}).Unscoped().Delete(&model.ShoppingCart{}).Error; err != nil {
			return err
		}
		return nil
	}); tErr != nil {
		return nil, err
	}

	return &proto.OrderInfoResponse{
		Id:      order.ID,
		OrderSn: order.OrderSn,
		Total:   orderAmount,
	}, nil
}

// 生成订单号
func GenerateOrderSn(userId int32) string {
	//订单号的生成规则
	/*
		年月日时分秒+用户id+2位随机数
	*/
	now := time.Now()
	rand.Seed(time.Now().UnixNano())
	orderSn := fmt.Sprintf("%d%d%d%d%d%d%d%d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Nanosecond(),
		userId, rand.Intn(90)+10,
	)
	return orderSn
}

// 订单列表
func (o *OrderService) OrderList(ctx context.Context, request *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	var (
		list  []model.OrderInfo
		total int64
		resp  proto.OrderListResponse
	)
	if err := global.DB.Model(&model.OrderInfo{}).Where(&model.OrderInfo{User: request.UserId}).Count(&total).Error; err != nil {
		return nil, err
	}
	if total == 0 {
		return &resp, nil
	}
	if err := global.DB.Scopes(Paginate(int(request.Pages), int(request.PagePerNums))).Where(&model.OrderInfo{User: request.UserId}).Find(&list).Error; err != nil {
		return nil, err
	}
	for _, order := range list {
		var to proto.OrderInfoResponse
		copier.Copy(&to, &order)
		to.Name = order.SignerName
		to.Mobile = order.SingerMobile
		to.AddTime = order.CreatedAt.Format("2006-01-02 15:04:05")
		resp.Data = append(resp.Data, &to)
	}
	return &proto.OrderListResponse{
		Total: int32(total),
		Data:  resp.Data,
	}, nil
}

// 订单详情
func (o *OrderService) OrderDetail(ctx context.Context, request *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	var (
		obj       model.OrderInfo
		rsp       proto.OrderInfoDetailResponse
		orderInfo proto.OrderInfoResponse
		goodsList []model.OrderGoods
	)
	// 这个订单的id是否是当前用户的订单， 如果在web层用户传递过来一个id的订单， web层应该先查询一下订单id是否是当前用户的
	// 在个人中心可以这样做，但是如果是后台管理系统，web层如果是后台管理系统 那么只传递order的id，如果是电商系统还需要一个用户的id
	if err := global.DB.Where(&model.OrderInfo{BaseModel: model.BaseModel{ID: request.Id}, User: request.UserId}).First(&obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "订单不存在")
		}
		return nil, err
	}
	rsp.OrderInfo = &orderInfo
	copier.Copy(&rsp.OrderInfo, &obj)
	rsp.OrderInfo.Name = obj.SignerName
	rsp.OrderInfo.Mobile = obj.SingerMobile
	if err := global.DB.Where(&model.OrderGoods{Order: obj.ID}).Find(&goodsList).Error; err != nil {
		return nil, err
	}
	for _, v := range goodsList {
		var to proto.OrderItemResponse
		copier.Copy(&to, &v)
		rsp.Goods = append(rsp.Goods, &to)
	}
	return &rsp, nil
}

// 修改订单状态
func (o *OrderService) UpdateOrderStatus(ctx context.Context, status *proto.OrderStatus) (*empty.Empty, error) {
	if err := global.DB.Model(&model.OrderInfo{}).Where("order_sn=?", status.OrderSn).Update("status", status.Status).Error; err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
