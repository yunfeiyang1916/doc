package db

import (
	"context"
	"shop-v2/internal/user/data"
	"shop-v2/pkg/common/meta"

	"gorm.io/gorm"
)

func NewUserDB(db *gorm.DB) *userDB {
	return &userDB{db: db}
}

type userDB struct {
	db *gorm.DB
}

func (u userDB) List(ctx context.Context, orderby []string, opts meta.ListMeta) (*data.UserDOList, error) {
	//TODO implement me
	panic("implement me")
}
