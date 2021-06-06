package gorm

import (
	"context"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"gorm.io/gorm"
)

type gormFetchStoreByStauts struct {
	db *gorm.DB
}

func NewGormFetchStoreByStauts(db *gorm.DB) *gormFetchStoreByStauts {
	return &gormFetchStoreByStauts{
		db: db,
	}
}

func (r *gormFetchStoreByStauts) Exec(ctx context.Context, status, sort string, limit, page int) (res []*domain.Store, total int64, err error) {
	var stores []*dto.StoreDBModel

	err = r.db.WithContext(ctx).
		Where("status = ?", status).
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
