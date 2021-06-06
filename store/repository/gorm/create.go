package gorm

import (
	"context"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"gorm.io/gorm"
)

type gormCreateStore struct {
	db *gorm.DB
}

func NewGormCreateStore(db *gorm.DB) *gormCreateStore {
	return &gormCreateStore{
		db: db,
	}
}

func (r *gormCreateStore) Add(ctx context.Context, store *domain.Store) (err error) {
	storeModel := &dto.StoreDBModel{}
	storeModel.ParserToDBModel(store)

	err = r.db.WithContext(ctx).
		Create(storeModel).
		Error
	return
}
