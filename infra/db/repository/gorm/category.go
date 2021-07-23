package gorm

import (
	"context"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/infra/db/model"
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
	categoryModel := &model.Category{}
	categoryModel.FromCategoryDomain(category)

	err = r.db.WithContext(ctx).
		Table("categories").
		Create(categoryModel).
		Error
	return
}

func (r *categoryRepository) FindByID(ctx context.Context, id string) (res *domain.Category, err error) {
	category := &model.Category{}
	err = r.db.
		WithContext(ctx).
		Table("categories").
		First(category, "id = ?", id).
		Error

	res = category.ToCategoryDomain()
	return
}

func (r *categoryRepository) Update(ctx context.Context, category *domain.Category) (err error) {
	categoryEntity := &model.Category{}
	categoryEntity.FromCategoryDomain(category)

	err = r.db.WithContext(ctx).
		Table("categories").
		Save(categoryEntity).
		Error
	return
}
