package gorm

import (
	"context"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type gormFetchStoreByTags struct {
	db *gorm.DB
}

func NewGormFetchStoreByTags(db *gorm.DB) *gormFetchStoreByTags {
	return &gormFetchStoreByTags{
		db: db,
	}
}

func (r *gormFetchStoreByTags) Exec(ctx context.Context, tags []string, sort string, limit, page int) (res []*domain.Store, total int64, err error) {
	var stores []*dto.StoreDBModel

	err = r.db.WithContext(ctx).
		Where("tags && ?", pq.StringArray(tags)).
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
