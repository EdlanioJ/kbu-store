package gorm

import (
	"context"

	"github.com/EdlanioJ/kbu-store/app/domain"
	"gorm.io/gorm"
)

type storeRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) *storeRepository {
	return &storeRepository{
		db: db,
	}
}

func (r *storeRepository) Create(ctx context.Context, store *domain.Store) (err error) {
	err = r.db.WithContext(ctx).
		Table("stores").
		Create(store).
		Error
	return
}

func (r *storeRepository) FindByID(ctx context.Context, id string) (res *domain.Store, err error) {
	res = &domain.Store{}

	err = r.db.WithContext(ctx).
		Table("stores").
		Where("id = ?", id).
		First(res).
		Error

	return
}

func (r *storeRepository) FindByName(ctx context.Context, name string) (res *domain.Store, err error) {
	res = &domain.Store{}

	err = r.db.WithContext(ctx).
		Table("stores").
		Where("name = ?", name).
		First(res).
		Error

	return
}

func (r *storeRepository) FindAll(ctx context.Context, sort string, limit, page int) (res domain.Stores, total int64, err error) {
	var stores []*domain.Store

	err = r.db.WithContext(ctx).
		Table("stores").
		Offset((page - 1) * limit).
		Limit(limit).
		Order(sort).
		Find(&stores).
		Count(&total).Error

	res = stores
	return
}

func (r *storeRepository) Update(ctx context.Context, store *domain.Store) (err error) {
	err = r.db.WithContext(ctx).
		Table("stores").
		Save(store).
		Error

	return
}

func (r *storeRepository) Delete(ctx context.Context, id string) (err error) {
	err = r.db.WithContext(ctx).
		Table("stores").
		Delete(&domain.Store{}, "id = ?", id).
		Error

	return
}
