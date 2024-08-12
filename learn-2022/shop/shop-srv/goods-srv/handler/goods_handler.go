package handler

import (
	"context"
	"fmt"
	"shop/shop-srv/goods-srv/global"
	"shop/shop-srv/goods-srv/model"
	"shop/shop-srv/goods-srv/proto"
	"strconv"

	"github.com/olivere/elastic/v7"

	"github.com/jinzhu/copier"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/golang/protobuf/ptypes/empty"
)

type GoodsService struct {
	proto.UnimplementedGoodsServer
}

func (g *GoodsService) GoodsList(ctx context.Context, request *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {

	//使用es的目的是搜索出商品的id来，通过id拿到具体的字段信息是通过mysql来完成
	//我们使用es是用来做搜索的， 是否应该将所有的mysql字段全部在es中保存一份
	//es用来做搜索，这个时候我们一般只把搜索和过滤的字段信息保存到es中
	//es可以用来当做mysql使用， 但是实际上mysql和es之间是互补的关系， 一般mysql用来做存储使用，es用来做搜索使用
	//es想要提高性能， 就要将es的内存设置的够大， 1k 2k

	// 使用es中的match bool 复合查询
	q := elastic.NewBoolQuery()
	tx := global.DB.WithContext(ctx).Model(&model.Goods{})
	if request.KeyWords != "" {
		//tx = tx.Where("name like ?", "%"+request.KeyWords+"%")
		q = q.Must(elastic.NewMultiMatchQuery(request.KeyWords, "name", "goods_brief"))
	}
	if request.IsHot {
		//tx = tx.Where(&model.Goods{IsHot: request.IsHot})
		// 使用filter不会计算权重
		q = q.Filter(elastic.NewTermQuery("is_hot", request.IsHot))
	}
	if request.IsNew {
		//tx = tx.Where(&model.Goods{IsNew: request.IsNew})
		q = q.Filter(elastic.NewTermQuery("is_new", request.IsNew))
	}
	if request.PriceMin > 0 {
		//tx = tx.Where("shop_price >= ?", request.PriceMin)
		q = q.Filter(elastic.NewRangeQuery("shop_price").Gte(request.PriceMin))
	}
	if request.PriceMax > 0 {
		//tx = tx.Where("shop_price <= ?", request.PriceMax)
		q = q.Filter(elastic.NewRangeQuery("shop_price").Lte(request.PriceMax))
	}
	if request.Brand > 0 {
		//tx = tx.Where(&model.Goods{BrandsID: request.Brand})
		q = q.Filter(elastic.NewTermQuery("brands_id", request.Brand))
	}
	if request.TopCategory > 0 {
		categoryIds := make([]interface{}, 0)
		var category model.Category
		if err := global.DB.Where("id = ?", request.TopCategory).First(&category).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, status.Error(codes.NotFound, "分类不存在")
			}
			return nil, err
		}
		// 通过分类筛选
		var subQuery string
		// 一级分类
		if category.Level == 1 {
			subQuery = fmt.Sprintf("select id from category where parent_category_id in (select id from category WHERE parent_category_id=%d)", request.TopCategory)
		} else if category.Level == 2 {
			subQuery = fmt.Sprintf("select id from category WHERE parent_category_id=%d", request.TopCategory)
		} else if category.Level == 3 {
			subQuery = fmt.Sprintf("select id from category WHERE id=%d", request.TopCategory)
		}
		//tx = tx.Where(fmt.Sprintf("category_id in (%s)", subQuery))
		var results []model.Category
		if err := global.DB.Model(model.Category{}).Raw(subQuery).Scan(&results).Error; err != nil {
			return nil, err
		}
		for _, re := range results {
			categoryIds = append(categoryIds, re.ID)
		}
		// 生成terms查询
		q = q.Filter(elastic.NewTermsQuery("category_id", categoryIds...))
	}

	// 分页
	if request.Pages == 0 {
		request.Pages = 1
	}

	switch {
	case request.PagePerNums > 100:
		request.PagePerNums = 100
	case request.PagePerNums <= 0:
		request.PagePerNums = 10
	}

	res, err := global.EsClient.Search(model.EsGoods{}.GetIndexName()).Query(q).From(int(request.Pages)).Size(int(request.PagePerNums)).Do(context.Background())
	if err != nil {
		return nil, err
	}

	var total int64
	//if err := tx.Count(&total).Error; err != nil {
	//	return nil, err
	//}
	total = res.Hits.TotalHits.Value
	if total == 0 {
		return &proto.GoodsListResponse{Total: 0, Data: nil}, nil
	}
	// 组装goods ids
	goodsIds := make([]int32, 0, len(res.Hits.Hits))
	for _, hit := range res.Hits.Hits {
		id, _ := strconv.Atoi(hit.Id)
		goodsIds = append(goodsIds, int32(id))
	}
	var goodsList []model.Goods
	//if err := tx.Preload("Category").Preload("Brands").Scopes(Paginate(int(request.Pages), int(request.PagePerNums))).Find(&goodsList).Error; err != nil {
	//	return nil, err
	//}
	if err = tx.Preload("Category").Preload("Brands").Where("id in (?)", goodsIds).Find(&goodsList).Error; err != nil {
		return nil, err
	}
	var resp = proto.GoodsListResponse{Total: int32(total)}
	for _, v := range goodsList {
		var obj proto.GoodsInfoResponse
		copier.Copy(&obj, &v)
		resp.Data = append(resp.Data, &obj)
	}
	return &resp, nil
}

func (g *GoodsService) BatchGetGoods(ctx context.Context, info *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	var list []model.Goods
	if err := global.DB.WithContext(ctx).Where("id in (?)", info.Id).Find(&list).Error; err != nil {
		return nil, err
	}
	var resp = proto.GoodsListResponse{Total: int32(len(list))}
	for _, v := range list {
		var obj proto.GoodsInfoResponse
		copier.Copy(&obj, &v)
		resp.Data = append(resp.Data, &obj)
	}
	return &resp, nil
}

func (g *GoodsService) CreateGoods(ctx context.Context, info *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	var category model.Category
	if err := global.DB.First(&category, info.CategoryId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "商品分类不存在")
		}
		return nil, err
	}
	var brand model.Brands
	if err := global.DB.First(&brand, info.BrandId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "品牌不存在")
		}
		return nil, err
	}
	entity := model.Goods{
		CategoryID:      info.CategoryId,
		BrandsID:        info.BrandId,
		Name:            info.Name,
		GoodsSn:         info.GoodsSn,
		MarketPrice:     info.MarketPrice,
		ShopPrice:       info.ShopPrice,
		ShipFree:        info.ShipFree,
		IsNew:           info.IsNew,
		IsHot:           info.IsHot,
		OnSale:          info.OnSale,
		GoodsBrief:      info.GoodsBrief,
		GoodsFrontImage: info.GoodsFrontImage,
		Images:          info.Images,
		DescImages:      info.DescImages,
	}
	// 在插入数据时，会调用钩子函数，钩子函数中会调用es的index方法，将数据插入es中
	// es保存失败则事务会回滚
	err := global.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&entity).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	var resp proto.GoodsInfoResponse
	copier.Copy(&resp, &entity)
	return &resp, nil
}

func (g *GoodsService) DeleteGoods(ctx context.Context, info *proto.DeleteGoodsInfo) (*empty.Empty, error) {
	// 在删除数据时，会调用钩子函数，钩子函数中会调用es的delete方法，将数据从es中删除
	err := global.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", info.Id).Delete(&model.Goods{BaseModel: model.BaseModel{ID: info.Id}}).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return status.Error(codes.NotFound, "商品不存在")
			}
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (g *GoodsService) UpdateGoods(ctx context.Context, info *proto.CreateGoodsInfo) (*empty.Empty, error) {
	var entity model.Goods
	if err := global.DB.WithContext(ctx).Where("id = ?", info.Id).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "商品不存在")
		}
		return nil, err
	}
	var count int64
	if err := global.DB.Where(&model.Category{}, info.CategoryId).Count(&count).Error; err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, status.Error(codes.NotFound, "商品分类不存在")
	}
	if err := global.DB.Where(&model.Brands{}, info.BrandId).Count(&count).Error; err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, status.Error(codes.NotFound, "品牌不存在")
	}
	var newEntity = model.Goods{
		CategoryID:      info.CategoryId,
		BrandsID:        info.BrandId,
		Name:            info.Name,
		GoodsSn:         info.GoodsSn,
		MarketPrice:     info.MarketPrice,
		ShopPrice:       info.ShopPrice,
		ShipFree:        info.ShipFree,
		IsNew:           info.IsNew,
		IsHot:           info.IsHot,
		OnSale:          info.OnSale,
		GoodsBrief:      info.GoodsBrief,
		GoodsFrontImage: info.GoodsFrontImage,
		Images:          info.Images,
		DescImages:      info.DescImages,
	}
	newEntity.ID = entity.ID
	// 在更新数据时，会调用钩子函数，钩子函数中会调用es的update方法，将数据更新到es中
	err := global.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entity).Updates(newEntity).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (g *GoodsService) GetGoodsDetail(ctx context.Context, request *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	var entity model.Goods
	if err := global.DB.WithContext(ctx).Where("id = ?", request.Id).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "商品不存在")
		}
		return nil, err
	}
	var resp proto.GoodsInfoResponse
	copier.Copy(&resp, &entity)
	return &resp, nil
}
