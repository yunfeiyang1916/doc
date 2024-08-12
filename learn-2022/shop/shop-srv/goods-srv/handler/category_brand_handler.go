package handler

import (
	"context"
	"shop/shop-srv/goods-srv/global"
	"shop/shop-srv/goods-srv/model"
	"shop/shop-srv/goods-srv/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gorm.io/gorm"

	"github.com/jinzhu/copier"

	"github.com/golang/protobuf/ptypes/empty"
)

func (g *GoodsService) CategoryBrandList(ctx context.Context, request *proto.CategoryBrandFilterRequest) (*proto.CategoryBrandListResponse, error) {
	var (
		count int64
		list  []model.CategoryBrands
		resp  proto.CategoryBrandListResponse
	)
	if err := global.DB.Model(&model.CategoryBrands{}).Count(&count).Error; err != nil {
		return nil, err
	}
	resp.Total = int32(count)
	if err := global.DB.Preload("Category").Preload("Brands").Scopes(Paginate(int(request.Pages), int(request.PagePerNums))).Find(&list).Error; err != nil {
		return nil, err
	}
	for _, v := range list {
		obj := proto.CategoryBrandResponse{}
		copier.Copy(&obj.Brand, v.Brands)
		copier.Copy(&obj.Category, v.Category)
		resp.Data = append(resp.Data, &obj)
	}
	return &resp, nil
}

// 获取商品分类下的品牌列表
func (g *GoodsService) GetCategoryBrandList(ctx context.Context, request *proto.CategoryInfoRequest) (*proto.BrandListResponse, error) {
	var entity model.Category
	if err := global.DB.First(&entity, request.Id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "商品分类不存在")
		}
		return nil, err
	}
	var list []model.CategoryBrands
	if err := global.DB.Preload("Brands").Where(&model.CategoryBrands{CategoryID: request.Id}).Find(&list).Error; err != nil {
		return nil, err
	}
	var resp proto.BrandListResponse
	resp.Total = int32(len(list))
	for _, v := range list {
		obj := proto.BrandInfoResponse{}
		copier.Copy(&obj, v.Brands)
		resp.Data = append(resp.Data, &obj)
	}
	return &resp, nil
}

func (g *GoodsService) CreateCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*proto.CategoryBrandResponse, error) {
	var count int64
	if err := global.DB.Where(&model.Category{}, request.CategoryId).Count(&count).Error; err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, status.Error(codes.NotFound, "商品分类不存在")
	}
	if err := global.DB.Where(&model.Brands{}, request.BrandId).Count(&count).Error; err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, status.Error(codes.NotFound, "品牌不存在")
	}
	var entity = &model.CategoryBrands{CategoryID: request.CategoryId, BrandsID: request.BrandId}
	if err := global.DB.Create(entity).Error; err != nil {
		return nil, err
	}
	return &proto.CategoryBrandResponse{Id: entity.ID}, nil
}

func (g *GoodsService) DeleteCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*empty.Empty, error) {
	if err := global.DB.Where(&model.CategoryBrands{}, request.Id).Delete(&model.CategoryBrands{}).Error; err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (g *GoodsService) UpdateCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*empty.Empty, error) {
	var count int64
	if err := global.DB.Where(&model.Category{}, request.CategoryId).Count(&count).Error; err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, status.Error(codes.NotFound, "商品分类不存在")
	}
	if err := global.DB.Where(&model.Brands{}, request.BrandId).Count(&count).Error; err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, status.Error(codes.NotFound, "品牌不存在")
	}
	// 查询分类品牌是否存在
	var entity model.CategoryBrands
	if err := global.DB.First(&entity, request.Id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "分类品牌不存在")
		}
		return nil, err
	}
	entity.CategoryID = request.CategoryId
	entity.BrandsID = request.BrandId
	if err := global.DB.Save(&entity).Error; err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
