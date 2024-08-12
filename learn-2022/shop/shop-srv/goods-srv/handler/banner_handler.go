package handler

import (
	"context"
	"shop/shop-srv/goods-srv/global"
	"shop/shop-srv/goods-srv/model"
	"shop/shop-srv/goods-srv/proto"

	"gorm.io/gorm"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jinzhu/copier"

	"github.com/golang/protobuf/ptypes/empty"
)

func (g *GoodsService) BannerList(ctx context.Context, empty *empty.Empty) (*proto.BannerListResponse, error) {
	var (
		resp    = &proto.BannerListResponse{}
		banners []model.Banner
	)
	// banner数量不多，一次性全部查询出来
	if err := global.DB.Find(&banners).Error; err != nil {
		return nil, err
	}

	resp.Total = int32(len(banners))
	for _, v := range banners {
		var m proto.BannerResponse
		copier.Copy(&m, &v)
		resp.Data = append(resp.Data, &m)
	}
	return resp, nil
}

func (g *GoodsService) CreateBanner(ctx context.Context, request *proto.BannerRequest) (*proto.BannerResponse, error) {
	banner := &model.Banner{
		Image: request.Image,
		Url:   request.Url,
		Index: request.Index,
	}
	if err := global.DB.Create(banner).Error; err != nil {
		return nil, err
	}
	return &proto.BannerResponse{
		Id:    banner.ID,
		Image: banner.Image,
		Url:   banner.Url,
		Index: banner.Index,
	}, nil
}

func (g *GoodsService) DeleteBanner(ctx context.Context, request *proto.BannerRequest) (*empty.Empty, error) {
	tx := global.DB.Delete(&model.Banner{}, "id = ?", request.Id)
	if err := tx.Error; err != nil {
		return nil, err
	}
	if tx.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "该轮播图不存在")
	}
	return &empty.Empty{}, nil
}

func (g *GoodsService) UpdateBanner(ctx context.Context, request *proto.BannerRequest) (*empty.Empty, error) {
	entity := model.Banner{}
	if err := global.DB.Where("id = ?", request.Id).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "该轮播图不存在")
		}
		return nil, err
	}
	if request.Url != "" {
		entity.Url = request.Url
	}
	if request.Image != "" {
		entity.Image = request.Image
	}
	if err := global.DB.Save(&entity).Error; err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
