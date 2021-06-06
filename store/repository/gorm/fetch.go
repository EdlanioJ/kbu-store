package gorm

import (
	"context"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"gorm.io/gorm"
)

type gormFetchStore struct {
	db *gorm.DB
}

func NewGormFetchStore(db *gorm.DB) *gormFetchStore {
	return &gormFetchStore{
		db: db,
	}
}

func (r *gormFetchStore) Exec(ctx context.Context, sort string, limit, page int) (res []*domain.Store, total int64, err error) {
	var stores []*dto.StoreDBModel

	err = r.db.WithContext(ctx).
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
