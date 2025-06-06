package mock

import (
	"context"
	"shop-v2/internal/user/data"
	"shop-v2/pkg/common/meta"
)

func NewUserMock() *userMock {
	return &userMock{}
}

type userMock struct{}

func (u userMock) List(ctx context.Context, orderby []string, opts meta.ListMeta) (*data.UserDOList, error) {
	return &data.UserDOList{}, nil
}
