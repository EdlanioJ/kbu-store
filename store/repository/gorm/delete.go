package gorm

import (
	"context"

	"github.com/EdlanioJ/kbu-store/dto"
	"gorm.io/gorm"
)

type gormDeleteStore struct {
	db *gorm.DB
}

func NewGormDeleteStore(db *gorm.DB) *gormDeleteStore {
	return &gormDeleteStore{
		db: db,
	}
}

func (r *gormDeleteStore) Exec(ctx context.Context, id string) (err error) {
	err = r.db.WithContext(ctx).
		Delete(&dto.StoreDBModel{}, "id = ?", id).
		Error

	return
}
