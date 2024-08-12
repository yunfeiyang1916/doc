package handler

import (
	"context"
	"shop/shop-srv/goods-srv/global"
	"shop/shop-srv/goods-srv/model"
	"shop/shop-srv/goods-srv/proto"

	"gorm.io/gorm"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/copier"
)

func (g *GoodsService) BrandList(ctx context.Context, request *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	var (
		resp   = &proto.BrandListResponse{}
		brands []model.Brands
		count  int64
	)
	if err := global.DB.Model(&model.Brands{}).Count(&count).Error; err != nil {
		return nil, err
	}
	resp.Total = int32(count)
	if count == 0 {
		return resp, nil
	}
	if err := global.DB.Scopes(Paginate(int(request.Pages), int(request.PagePerNums))).Find(&brands).Error; err != nil {
		return nil, err
	}
	for _, v := range brands {
		var m proto.BrandInfoResponse
		copier.Copy(&m, &v)
		resp.Data = append(resp.Data, &m)
	}
	return resp, nil
}

func (g *GoodsService) CreateBrand(ctx context.Context, request *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	// 是否已存在
	var count int64
	if err := global.DB.Model(&model.Brands{}).Where("name = ?", request.Name).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, status.Error(codes.InvalidArgument, "该品牌已存在")
	}
	brand := &model.Brands{
		Name: request.Name,
		Logo: request.Logo,
	}
	if err := global.DB.Create(brand).Error; err != nil {
		return nil, err
	}
	return &proto.BrandInfoResponse{Id: brand.ID}, nil
}

func (g *GoodsService) DeleteBrand(ctx context.Context, request *proto.BrandRequest) (*empty.Empty, error) {
	tx := global.DB.Delete(&model.Brands{}, "id = ?", request.Id)
	if err := tx.Error; err != nil {
		return nil, err
	}
	if tx.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "该品牌不存在")
	}
	return &empty.Empty{}, nil
}

func (g *GoodsService) UpdateBrand(ctx context.Context, request *proto.BrandRequest) (*empty.Empty, error) {
	brand := model.Brands{}
	if err := global.DB.Where("id = ?", request.Id).First(&brand).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "该品牌不存在")
		}
		return nil, err
	}
	if request.Name != "" {
		brand.Name = request.Name
	}
	if request.Logo != "" {
		brand.Logo = request.Logo
	}
	if err := global.DB.Save(&brand).Error; err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
