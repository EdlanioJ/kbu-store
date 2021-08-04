package gorm

import (
	"context"

	"github.com/EdlanioJ/kbu-store/app/domain"
	"github.com/opentracing/opentracing-go"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "categoryRepository.Create")
	defer span.Finish()

	err = r.db.WithContext(ctx).
		Table("categories").
		Create(category).
		Error
	return
}

func (r *categoryRepository) FindByID(ctx context.Context, id string) (res *domain.Category, err error) {
	category := &domain.Category{}

	span, ctx := opentracing.StartSpanFromContext(ctx, "categoryRepository.FindByID")
	defer span.Finish()

	err = r.db.
		WithContext(ctx).
		Table("categories").
		First(category, "id = ?", id).
		Error

	res = category
	return
}

func (r *categoryRepository) Update(ctx context.Context, category *domain.Category) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "categoryRepository.Update")
	defer span.Finish()

	err = r.db.WithContext(ctx).
		Table("categories").
		Save(category).
		Error
	return
}
