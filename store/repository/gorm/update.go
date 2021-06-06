package gorm

import (
	"context"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"gorm.io/gorm"
)

type gormUpdateStore struct {
	db *gorm.DB
}

func NewGormUpdateStore(db *gorm.DB) *gormUpdateStore {
	return &gormUpdateStore{
		db: db,
	}
}

func (r *gormUpdateStore) Exec(ctx context.Context, store *domain.Store) (err error) {
	storeEntity := &dto.StoreDBModel{}
	storeEntity.ParserToDBModel(store)

	err = r.db.WithContext(ctx).
		Save(storeEntity).
		Error

	return
}
