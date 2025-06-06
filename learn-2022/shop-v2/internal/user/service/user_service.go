package service

import (
	"shop-v2/internal/user/data"
	"shop-v2/pkg/common/meta"

	"golang.org/x/net/context"
)

type UserDTOList struct {
	TotalCount int64 `json:"totalCount,omitempty"` //总数
	//Items      []*UserDTO `json:"data"`
}

type UserSrv interface {
	List(ctx context.Context, orderby []string, opts meta.ListMeta) (*UserDTOList, error)
}

// service层的管理器
type userService struct {
	userStore data.UserStore
}

func NewUserService(s data.UserStore) *userService {
	return &userService{userStore: s}
}

func (u *userService) List(ctx context.Context, orderby []string, opts meta.ListMeta) (*UserDTOList, error) {
	return nil, nil
}
