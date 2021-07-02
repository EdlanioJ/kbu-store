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

func (r *categoryRepository) Create(ctx context.Context, category *domain.Category) (err error) {
	categoryModel := &model.Category{}
	categoryModel.FromCategoryDomain(category)

	err = r.db.WithContext(ctx).
		Table("categories").
		Create(categoryModel).
		Error
	return
}

func (r *categoryRepository) GetAll(ctx context.Context, sort string, page, limit int) (res []*domain.Category, total int64, err error) {
	var categories []*model.Category

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

	res = make([]*domain.Category, 0)
	for _, value := range categories {
		res = append(res, value.ToCategoryDomain())
	}
	return
}

func (r *categoryRepository) GetAllByStatus(ctx context.Context, status, sort string, page, limit int) (res []*domain.Category, total int64, err error) {
	var categories []*model.Category

	offset := (page - 1) * limit
	err = r.db.
		Table("categories").
		Where("status = ?", status).
		Limit(limit).
		Offset(offset).
		WithContext(ctx).
		Order(sort).
		Find(&categories).
		Count(&total).
		Error

	res = make([]*domain.Category, 0)
	for _, value := range categories {
		res = append(res, value.ToCategoryDomain())
	}
	return
}

func (r *categoryRepository) GetById(ctx context.Context, id string) (res *domain.Category, err error) {
	category := &model.Category{}
	err = r.db.
		WithContext(ctx).
		Table("categories").
		First(category, "id = ?", id).
		Error

	res = category.ToCategoryDomain()
	return
}

func (r *categoryRepository) GetByIdAndStatus(ctx context.Context, id, status string) (res *domain.Category, err error) {
	category := &model.Category{}
	err = r.db.WithContext(ctx).
		Table("categories").
		Where("id = ? AND status = ?", id, status).
		First(category).
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
