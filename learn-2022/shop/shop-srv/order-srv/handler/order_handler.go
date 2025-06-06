package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	goodsProto "shop/shop-srv/goods-srv/proto"
	inventoryProto "shop/shop-srv/inventory-srv/proto"
	"shop/shop-srv/order-srv/global"
	"shop/shop-srv/order-srv/model"
	"shop/shop-srv/order-srv/proto"
	"time"

	"github.com/apache/rocketmq-client-go/v2/consumer"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go.uber.org/zap"

	"github.com/apache/rocketmq-client-go/v2/primitive"

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
	//tracer := otel.Tracer("OrderService.CartItemList")
	//ctx, parentSpan := tracer.Start(ctx, "OrderService.CartItemList")
	//time.Sleep(time.Second)
	//ctx, span1 := tracer.Start(ctx, "span1")
	//span1.End()
	//_, span2 := tracer.Start(ctx, "span2")
	//span2.End()
	//parentSpan.End()
	var (
		list []model.ShoppingCart
		rsp  proto.CartItemListResponse
	)
	if r := global.DB.WithContext(ctx).Where(&model.ShoppingCart{User: info.Id}).Find(&list); r.Error != nil {
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
				User:    request.UserId,
				Goods:   request.GoodsId,
				Nums:    request.Nums,
				Checked: request.Checked,
			}
		} else {
			return nil, err
		}
	} else {
		//如果记录已经存在，则合并购物车记录, 更新操作
		obj.Nums += request.Nums
		obj.Checked = request.Checked
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

// 订单的事务侦听器
type OrderListener struct {
	Ctx context.Context
	Err error
	model.OrderInfo
}

// 当发送事务预处理消息成功时，调用此方法执行本地事务
func (o *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	var (
		ctx         = o.Ctx
		goodsIds    []int32
		shopCharts  []*model.ShoppingCart
		goodsNumMap = make(map[int32]int32)
	)
	if err := global.DB.Where(&model.ShoppingCart{User: o.User, Checked: true}).Find(&shopCharts).Error; err != nil {
		o.Err = err
		// 此时还没有扣减库存，所以直接回滚事务消息即可
		return primitive.RollbackMessageState
	}
	if len(shopCharts) == 0 {
		o.Err = status.Errorf(codes.InvalidArgument, "没有选中的商品")
		return primitive.RollbackMessageState
	}
	for _, v := range shopCharts {
		goodsIds = append(goodsIds, v.Goods)
		goodsNumMap[v.Goods] = v.Nums
	}
	// 跨服务调用
	conn, err := global.GoodsSrvConnPool.Get()
	if err != nil {
		o.Err = err
		return primitive.RollbackMessageState
	}
	defer conn.Close()
	goodsSrvClient := goodsProto.NewGoodsClient(conn.Value())
	goodsResp, err := goodsSrvClient.BatchGetGoods(ctx, &goodsProto.BatchGoodsIdInfo{Id: goodsIds})
	if err != nil {
		o.Err = err
		return primitive.RollbackMessageState
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
	o.OrderMount = orderAmount

	// 调用库存服务，进行库存扣减
	inventoryConn, err := global.InventorySrvConnPool.Get()
	if err != nil {
		// todo 需要判断下具体的失败原因，如果是网络问题，则有可能扣减成功了，此时就需要commit消息
		o.Err = err
		return primitive.CommitMessageState
	}
	defer inventoryConn.Close()
	inventoryClient := inventoryProto.NewInventoryClient(inventoryConn.Value())
	if _, err = inventoryClient.Sell(ctx, &inventoryProto.SellInfo{OrderSn: o.OrderSn, GoodsInfo: goodsInvInfos}); err != nil {
		o.Err = err
		return primitive.RollbackMessageState
	}
	// 生成订单表
	order := o.OrderInfo
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
		if err = tx.Where(&model.ShoppingCart{User: o.User, Checked: true}).Unscoped().Delete(&model.ShoppingCart{}).Error; err != nil {
			return err
		}
		//发送延时消息
		//p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.0.104:9876"}))
		//if err != nil {
		//	panic("生成producer失败")
		//}
		//
		////不要在一个进程中使用多个producer， 但是不要随便调用shutdown因为会影响其他的producer
		//if err = p.Start(); err != nil {panic("启动producer失败")}
		//
		//msg = primitive.NewMessage("order_timeout", msg.Body)
		//msg.WithDelayTimeLevel(3)
		//_, err = p.SendSync(context.Background(), msg)
		//if err != nil {
		//	zap.S().Errorf("发送延时消息失败: %v\n", err)
		//	tx.Rollback()
		//	o.Code = codes.Internal
		//	o.Detail = "发送延时消息失败"
		//	return primitive.CommitMessageState
		//}
		return nil
	}); tErr != nil {
		o.Err = err
		// 订单创建失败，因为此时已经扣减了库存了，所以需要确认发送归还库存消息
		return primitive.CommitMessageState
	}
	// 订单创建成功，不需要归还库存了，所以需要回滚消息
	return primitive.RollbackMessageState
}

// 当未收到预处理响应时，调用此方法检查并返回本地事务状态。
func (o *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	// 此时使用o可能有问题，可能是宕机后来的回调，o.OrderSn有可能是空
	if err := global.DB.Where(&model.OrderInfo{OrderSn: o.OrderSn}).First(&model.OrderInfo{}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 走到这里，说明订单表里面没有数据，说明订单创建失败了，需要归还库存（库存有可能还没扣减，所以库存服务需要兼容下）
			return primitive.CommitMessageState
		}
		return primitive.UnknowState
	}
	// 已经创建订单了，所以不需要归还库存了
	return primitive.RollbackMessageState
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
	// 生成订单表
	order := model.OrderInfo{
		OrderSn:      GenerateOrderSn(request.UserId),
		Address:      request.Address,
		SignerName:   request.Name,
		SingerMobile: request.Mobile,
		Post:         request.Post,
		User:         request.UserId,
	}
	// 使用rocketmq中的事务消息来实现分布式事务
	// 因为库存会出现库存不足的情况，所以需要发送归还消息，先发送退还消息，然后调用接口扣库存，如果库存不足，则不提交消息确认
	var orderListener = &OrderListener{
		Ctx:       ctx,
		OrderInfo: order,
	}
	p, err := rocketmq.NewTransactionProducer(orderListener, producer.WithNameServer([]string{global.ServerConfig.RocketMQConfig.GetAddr()}))
	if err != nil {
		zap.S().Errorf("生成producer失败: %s", err.Error())
		return nil, err
	}
	if err = p.Start(); err != nil {
		zap.S().Errorf("启动producer失败: %s", err.Error())
		return nil, err
	}
	//defer p.Shutdown()
	orderListener.OrderInfo = order
	jsonStr, _ := json.Marshal(order)
	// 预发送事务消息，这里会回调orderListener的ExecuteLocalTransaction方法
	_, err = p.SendMessageInTransaction(ctx, primitive.NewMessage("order_reback", jsonStr))
	if err != nil {
		zap.S().Errorf("发送事务消息失败: %s", err.Error())
		return nil, err
	}
	if orderListener.Err != nil {
		return nil, orderListener.Err
	}
	return &proto.OrderInfoResponse{
		Id:      order.ID,
		OrderSn: order.OrderSn,
		Total:   orderListener.OrderMount,
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

// 自己订阅订单超时消息
func OrderTimeout(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

	for i := range msgs {
		var orderInfo model.OrderInfo
		_ = json.Unmarshal(msgs[i].Body, &orderInfo)

		fmt.Printf("获取到订单超时消息: %v\n", time.Now())
		//查询订单的支付状态，如果已支付什么都不做，如果未支付，归还库存
		var order model.OrderInfo
		if result := global.DB.Model(model.OrderInfo{}).Where(model.OrderInfo{OrderSn: orderInfo.OrderSn}).First(&order); result.RowsAffected == 0 {
			return consumer.ConsumeSuccess, nil
		}
		if order.Status != "TRADE_SUCCESS" {
			tx := global.DB.Begin()
			//归还库存，我们可以模仿order中发送一个消息到 order_reback中去
			//修改订单的状态为已支付
			order.Status = "TRADE_CLOSED"
			tx.Save(&order)

			p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.0.104:9876"}))
			if err != nil {
				panic("生成producer失败")
			}

			if err = p.Start(); err != nil {
				panic("启动producer失败")
			}

			_, err = p.SendSync(context.Background(), primitive.NewMessage("order_reback", msgs[i].Body))
			if err != nil {
				tx.Rollback()
				fmt.Printf("发送失败: %s\n", err)
				return consumer.ConsumeRetryLater, nil
			}

			//if err = p.Shutdown(); err != nil {panic("关闭producer失败")}
			return consumer.ConsumeSuccess, nil
		}
	}
	return consumer.ConsumeSuccess, nil
}
