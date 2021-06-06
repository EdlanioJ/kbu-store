package gorm

import (
	"context"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"gorm.io/gorm"
)

type gormUpdateCategory struct {
	db *gorm.DB
}

func NewGormUpdateCategory(db *gorm.DB) *gormUpdateCategory {
	return &gormUpdateCategory{
		db: db,
	}
}

func (r *gormUpdateCategory) Exec(ctx context.Context, category *domain.Category) (err error) {
	CategoryEntity := &dto.CategoryDBModel{}
	CategoryEntity.ParserToDBModel(category)

	err = r.db.WithContext(ctx).
		Table("categories").
		Save(CategoryEntity).
		Error
	return
}
