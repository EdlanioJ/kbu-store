package gorm

import (
	"context"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"gorm.io/gorm"
)

type gormGetStoreByID struct {
	db *gorm.DB
}

func NewGormGetStoreByID(db *gorm.DB) *gormGetStoreByID {
	return &gormGetStoreByID{
		db: db,
	}
}

func (r *gormGetStoreByID) Exec(ctx context.Context, id string) (res *domain.Store, err error) {
	store := &dto.StoreDBModel{}

	err = r.db.WithContext(ctx).
		Where("id = ?", id).
		First(store).
		Error

	res = store.ParserToStoreDomain()
	return
}
