package gorm

import (
	"context"

	"github.com/EdlanioJ/kbu-store/app/domain"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *categoryRepository {
	return &categoryRepository{
		db: db,
	}
}

func (r *categoryRepository) Store(ctx context.Context, category *domain.Category) (err error) {
	err = r.db.WithContext(ctx).
		Table("categories").
		Create(category).
		Error
	return
}

func (r *categoryRepository) FindByID(ctx context.Context, id string) (res *domain.Category, err error) {
	category := &domain.Category{}
	err = r.db.
		WithContext(ctx).
		Table("categories").
		First(category, "id = ?", id).
		Error

	res = category
	return
}

func (r *categoryRepository) Update(ctx context.Context, category *domain.Category) (err error) {
	err = r.db.WithContext(ctx).
		Table("categories").
		Save(category).
		Error
	return
}
