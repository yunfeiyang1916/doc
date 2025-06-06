package data

import (
	"context"
	"shop-v2/pkg/common/meta"
)

type UserDOList struct {
	TotalCount int64 `json:"totalCount,omitempty"` //总数
	//Items      []*UserDO `json:"data"`
}

type UserStore interface {
	List(ctx context.Context, orderby []string, opts meta.ListMeta) (*UserDOList, error)
}
