package handler

import (
	"context"
	"fmt"
	"os"
	"shop/shop-srv/user-srv/initialize"
	"shop/shop-srv/user-srv/proto"
	"testing"
)

func initConfig() {
	os.Setenv("shop_srv_config", "../config-debug.yaml")
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
}

func TestUserServer_CreateUser(t *testing.T) {
	initConfig()
	userServer := &UserService{}
	user, err := userServer.CreateUser(context.Background(), &proto.CreateUserInfo{
		NickName: "张三",
		Password: "123456",
		Mobile:   "18612922641",
	})
	fmt.Println(user, err)
}
