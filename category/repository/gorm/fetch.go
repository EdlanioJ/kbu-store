package gorm

import (
	"context"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"gorm.io/gorm"
)

type gormFetchCategory struct {
	db *gorm.DB
}

func NewGormFetchCategory(db *gorm.DB) *gormFetchCategory {
	return &gormFetchCategory{
		db: db,
	}
}

func (r *gormFetchCategory) Exec(ctx context.Context, sort string, page, limit int) (res []*domain.Category, total int64, err error) {
	var categories []*dto.CategoryDBModel

	offset := (page - 1) * limit
	err = r.db.
		WithContext(ctx).
		Table("categories").
		Limit(limit).
		Offset(offset).
		Order(sort).
		Find(&categories).
		Count(&total).
		Error

	for _, value := range categories {
		res = append(res, value.ParserToCategoryDomain())
	}
	return
}
