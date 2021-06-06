package gorm

import (
	"context"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"gorm.io/gorm"
)

type gormFetchCategoryByStatus struct {
	db *gorm.DB
}

func NewGormFetchCategoryByStatus(db *gorm.DB) *gormFetchCategoryByStatus {
	return &gormFetchCategoryByStatus{
		db: db,
	}
}

func (r *gormFetchCategoryByStatus) Exec(ctx context.Context, status, sort string, page, limit int) (res []*domain.Category, total int64, err error) {
	var categories []*dto.CategoryDBModel

	offset := (page - 1) * limit
	err = r.db.
		Table("categories").
		Limit(limit).
		Offset(offset).
		WithContext(ctx).
		Order(sort).
		Find(&categories, "status = ?", status).
		Count(&total).
		Error

	for _, value := range categories {
		res = append(res, value.ParserToCategoryDomain())
	}

	return
}
