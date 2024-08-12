package tests

import (
	"context"
	"shop/shop-srv/userop-srv/proto"
	"sync"
	"testing"

	"google.golang.org/grpc"
)

var (
	conn   *grpc.ClientConn
	client proto.InventoryClient
)

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50053", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client = proto.NewInventoryClient(conn)
}

func TestSetInv(t *testing.T) {
	Init()
	// 初始化一批库存
	for i := 421; i <= 840; i++ {
		info := &proto.GoodsInvInfo{
			GoodsId: int32(i),
			Num:     100,
		}
		_, err := client.SetInv(context.Background(), info)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestInvDetail(t *testing.T) {
	Init()
	info := &proto.GoodsInvInfo{
		GoodsId: 421,
	}
	r, err := client.InvDetail(context.Background(), info)
	if err != nil {
		panic(err)
	}
	t.Logf("库存详情：%+v", r)
}

func TestSell(t *testing.T) {
	Init()
	info := &proto.SellInfo{
		GoodsInfo: []*proto.GoodsInvInfo{
			{GoodsId: 421, Num: 1},
			{GoodsId: 422, Num: 1},
		},
	}
	wg := sync.WaitGroup{}
	gNum := 20
	for i := 0; i < gNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r, err := client.Sell(context.Background(), info)
			if err != nil {
				panic(err)
			}
			t.Logf("售卖成功：%+v", r)
		}()
	}
	wg.Wait()
}

func TestReback(t *testing.T) {
	Init()
	info := &proto.SellInfo{
		GoodsInfo: []*proto.GoodsInvInfo{
			{GoodsId: 421, Num: 99},
			{GoodsId: 425, Num: 99},
		},
	}
	r, err := client.Reback(context.Background(), info)
	if err != nil {
		panic(err)
	}
	t.Logf("归还成功：%+v", r)
}
