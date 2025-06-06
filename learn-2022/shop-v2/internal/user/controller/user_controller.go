package controller

import (
	v1 "shop-v2/api/user/v1"
	"shop-v2/internal/user/service"
)

type userServer struct {
	v1.UnimplementedUserServer
	srv service.UserSrv
}

func NewUserServer(srv service.UserSrv) *userServer {
	return &userServer{srv: srv}
}
