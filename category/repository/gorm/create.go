package gorm

import (
	"context"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"gorm.io/gorm"
)

type gormCreateCategory struct {
	db *gorm.DB
}

func NewGormCreateCategory(db *gorm.DB) *gormCreateCategory {
	return &gormCreateCategory{
		db: db,
	}
}

func (r *gormCreateCategory) Add(ctx context.Context, category *domain.Category) (err error) {
	categoryModel := &dto.CategoryDBModel{}
	categoryModel.ParserToDBModel(category)

	err = r.db.WithContext(ctx).
		Table("categories").
		Create(categoryModel).
		Error
	return
}
