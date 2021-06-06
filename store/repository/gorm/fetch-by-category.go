package gorm

import (
	"context"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"gorm.io/gorm"
)

type gormFetchStoreByCategory struct {
	db *gorm.DB
}

func NewGormFetchStoreByCategory(db *gorm.DB) *gormFetchStoreByCategory {
	return &gormFetchStoreByCategory{
		db: db,
	}
}

func (r *gormFetchStoreByCategory) Exec(ctx context.Context, categoryID, sort string, limit, page int) (res []*domain.Store, total int64, err error) {
	var stores []*dto.StoreDBModel

	err = r.db.WithContext(ctx).
		Where("category_id = ?", categoryID).
		Offset((page - 1) * limit).
		Limit(limit).
		Order(sort).
		Find(&stores).
		Count(&total).Error

	res = make([]*domain.Store, 0)
	for _, value := range stores {
		res = append(res, value.ParserToStoreDomain())
	}
	return
}
