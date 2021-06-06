package gorm

import (
	"context"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"gorm.io/gorm"
)

type gormGetCategoryByStatus struct {
	db *gorm.DB
}

func NewGormGetCategoryByStatus(db *gorm.DB) *gormGetCategoryByStatus {
	return &gormGetCategoryByStatus{
		db: db,
	}
}

func (r *gormGetCategoryByStatus) Exec(ctx context.Context, id, status string) (res *domain.Category, err error) {
	category := &dto.CategoryDBModel{}
	err = r.db.WithContext(ctx).
		Table("categories").
		First(category, "id = ? AND status = ?", id, status).
		Error

	res = category.ParserToCategoryDomain()
	return
}
