package handler

import (
	"context"
	"encoding/json"
	"shop/shop-srv/goods-srv/global"
	"shop/shop-srv/goods-srv/model"
	"shop/shop-srv/goods-srv/proto"

	"github.com/jinzhu/copier"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/golang/protobuf/ptypes/empty"
)

func (g *GoodsService) GetAllCategoryList(ctx context.Context, empty *empty.Empty) (*proto.CategoryListResponse, error) {
	var list []model.Category
	// 只查一级的,预加载子类及子类的子类
	if err := global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&list).Error; err != nil {
		return nil, err
	}
	b, _ := json.Marshal(list)
	return &proto.CategoryListResponse{JsonData: string(b)}, nil
}

func (g *GoodsService) GetSubCategory(ctx context.Context, request *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	var entity model.Category
	if err := global.DB.First(&entity, request.Id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "商品分类不存在")
		}
		return nil, err
	}
	resp := &proto.SubCategoryListResponse{}
	copier.Copy(resp.Info, entity)
	var subCategories []*model.Category
	// 预加载子类
	preloads := "SubCategory"
	if entity.Level == 1 {
		// 预加载子类及子类的子类
		preloads = "SubCategory.SubCategory"
	}
	if err := global.DB.Where(&model.Category{ParentCategoryID: request.Id}).Preload(preloads).Find(&subCategories).Error; err != nil {
		return nil, err
	}
	for _, v := range subCategories {
		sub := proto.CategoryInfoResponse{}
		copier.Copy(&sub, v)
		resp.SubCategorys = append(resp.SubCategorys, &sub)
	}
	return resp, nil
}

func (g *GoodsService) CreateCategory(ctx context.Context, request *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	entity := model.Category{
		Name:             request.Name,
		Level:            request.Level,
		IsTab:            request.IsTab,
		ParentCategoryID: request.ParentCategory,
	}
	if err := global.DB.Create(&entity).Error; err != nil {
		return nil, err
	}
	resp := &proto.CategoryInfoResponse{}
	copier.Copy(resp, &entity)
	return resp, nil
}

func (g *GoodsService) DeleteCategory(ctx context.Context, request *proto.DeleteCategoryRequest) (*empty.Empty, error) {
	if err := global.DB.Delete(&model.Category{}, request.Id).Error; err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (g *GoodsService) UpdateCategory(ctx context.Context, request *proto.CategoryInfoRequest) (*empty.Empty, error) {
	var entity model.Category
	if err := global.DB.First(&entity, request.Id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "商品分类不存在")
		}
		return nil, err
	}
	if request.Name != "" {
		entity.Name = request.Name
	}
	if request.ParentCategory != 0 {
		entity.ParentCategoryID = request.ParentCategory
	}
	if request.Level != 0 {
		entity.Level = request.Level
	}
	if request.IsTab {
		entity.IsTab = request.IsTab
	}
	if err := global.DB.Save(&entity).Error; err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
