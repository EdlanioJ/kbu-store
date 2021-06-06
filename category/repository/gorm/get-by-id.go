package gorm

import (
	"context"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"gorm.io/gorm"
)

type gormGetCategoryByID struct {
	db *gorm.DB
}

func NewGormGetCategoryByID(db *gorm.DB) *gormGetCategoryByID {
	return &gormGetCategoryByID{
		db: db,
	}
}

func (r *gormGetCategoryByID) Exec(ctx context.Context, id string) (res *domain.Category, err error) {
	category := &dto.CategoryDBModel{}
	err = r.db.
		WithContext(ctx).
		Table("categories").
		First(category, "id = ?", id).
		Error

	res = category.ParserToCategoryDomain()
	return
}
