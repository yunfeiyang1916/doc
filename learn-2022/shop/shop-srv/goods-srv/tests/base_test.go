package tests

import (
	"os"
	"shop/shop-srv/goods-srv/initialize"
	"shop/shop-srv/goods-srv/proto"

	"google.golang.org/grpc"
)

var (
	conn   *grpc.ClientConn
	client proto.GoodsClient
)

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50052", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client = proto.NewGoodsClient(conn)
}

func InitConfig() {
	os.Setenv("shop_srv_config", "../config-debug.yaml")
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	initialize.InitEs()
}
