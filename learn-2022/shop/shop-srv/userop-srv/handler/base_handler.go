package handler

import (
	"shop/shop-srv/userop-srv/proto"

	"gorm.io/gorm"
)

// 使用gorm的作用域来复用通用的分页逻辑
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

type UserOpService struct {
	proto.UnimplementedAddressServer
	proto.UnimplementedUserFavServer
	proto.UnimplementedMessageServer
}
