package gorm

import (
	"context"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"gorm.io/gorm"
)

type gormFetchStoreByOwner struct {
	db *gorm.DB
}

func NewGormFetchStoreByOwner(db *gorm.DB) *gormFetchStoreByOwner {
	return &gormFetchStoreByOwner{
		db: db,
	}
}

func (r *gormFetchStoreByOwner) Exec(ctx context.Context, externalID, sort string, limit, page int) (res []*domain.Store, total int64, err error) {
	var stores []*dto.StoreDBModel

	err = r.db.WithContext(ctx).
		Where("external_id = ?", externalID).
		Offset((page - 1) * limit).
		Limit(limit).
		Order(sort).
		Find(&stores).
		Count(&total).
		Error

	res = make([]*domain.Store, 0)
	for _, value := range stores {
		res = append(res, value.ParserToStoreDomain())
	}
	return
}
