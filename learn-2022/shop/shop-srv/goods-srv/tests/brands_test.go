package tests

import (
	"context"
	"fmt"
	"shop/shop-srv/goods-srv/proto"
	"testing"
)

func Test_GetBrandList(t *testing.T) {
	Init()
	rsp, err := client.BrandList(context.Background(), &proto.BrandFilterRequest{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("总数：", rsp.Total)
	for _, v := range rsp.Data {
		fmt.Println(v.Name)
	}
}
