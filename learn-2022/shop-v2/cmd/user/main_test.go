package main

import (
	"context"
	"fmt"
	v1 "shop-v2/api/user/v1"
	"shop-v2/gmicro/server/rpcserver"
	_ "shop-v2/gmicro/server/rpcserver/resolver/direct"
	"testing"
)

func TestDirectResolver(t *testing.T) {
	ctx := context.Background()
	conn, err := rpcserver.DialInsecure(ctx, rpcserver.WithEndpoint("direct:///127.0.0.1:8021"))
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer conn.Close()
	uc := v1.NewUserClient(conn)
	r, err := uc.GetUserList(ctx, &v1.PageInfo{})
	if err != nil {
		t.Fatalf("%v", err)
	}
	fmt.Println(r)
}
