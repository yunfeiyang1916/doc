package tests

import (
	"context"
	"shop/shop-srv/order-srv/proto"
	"testing"

	"google.golang.org/grpc"
)

var (
	conn   *grpc.ClientConn
	client proto.OrderClient
)

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50054", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client = proto.NewOrderClient(conn)
}

func TestCreateCartItem(t *testing.T) {
	Init()
	r, err := client.CreateCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  1,
		Nums:    1,
		GoodsId: 421,
	})
	if err != nil {
		panic(err)
	}
	t.Logf("创建购物车成功：%+v", r)
}

func TestCartItemList(t *testing.T) {
	Init()
	rsp, err := client.CartItemList(context.Background(), &proto.UserInfo{Id: 1})
	if err != nil {
		panic(err)
	}
	t.Logf("购物车列表：%+v", rsp)
}

func TestUpdateCartItem(t *testing.T) {
	Init()
	rsp, err := client.UpdateCartItem(context.Background(), &proto.CartItemRequest{
		Id:      3,
		Checked: true,
	})
	if err != nil {
		panic(err)
	}
	t.Logf("更新购物车成功：%+v", rsp)
}

func TestCreateOrder(t *testing.T) {
	Init()
	rsp, err := client.CreateOrder(context.Background(), &proto.OrderRequest{
		UserId:  1,
		Address: "北京市",
		Name:    "张三",
		Mobile:  "1888888887",
		Post:    "请尽快发货",
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("创建订单成功：%+v", rsp)
}

func TestOrderDetail(t *testing.T) {
	Init()
	rsp, err := client.OrderDetail(context.Background(), &proto.OrderRequest{
		Id: 1,
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("订单详情：%+v", rsp)
}
