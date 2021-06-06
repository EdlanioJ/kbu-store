package gorm

import (
	"context"

	"github.com/EdlanioJ/kbu-store/dto"
	"gorm.io/gorm"
)

type gormDeleteAccount struct {
	db *gorm.DB
}

func NewGormDeleteAccount(db *gorm.DB) *gormDeleteAccount {
	return &gormDeleteAccount{
		db: db,
	}
}

func (r *gormDeleteAccount) Exec(ctx context.Context, id string) (err error) {
	err = r.db.WithContext(ctx).
		Table("accounts").
		Delete(&dto.CategoryDBModel{}, "id = ?", id).
		Error

	return
}
