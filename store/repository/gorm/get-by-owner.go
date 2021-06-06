package gorm

import (
	"context"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"gorm.io/gorm"
)

type gormGetStoreByOwner struct {
	db *gorm.DB
}

func NewGormGetStoreByOwner(db *gorm.DB) *gormGetStoreByOwner {
	return &gormGetStoreByOwner{
		db: db,
	}
}

func (r *gormGetStoreByOwner) Exec(ctx context.Context, id string, externalID string) (res *domain.Store, err error) {
	store := &dto.StoreDBModel{}

	err = r.db.WithContext(ctx).
		Where("id = ? AND external_id = ?", id, externalID).
		First(store).
		Error

	res = store.ParserToStoreDomain()
	return
}
