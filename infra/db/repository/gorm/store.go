package gorm

import (
	"context"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/infra/db/model"
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
	storeModel := &model.Store{}
	storeModel.FromStoreDomain(store)

	err = r.db.WithContext(ctx).
		Create(storeModel).
		Error
	return
}

func (r *storeRepository) FindByID(ctx context.Context, id string) (res *domain.Store, err error) {
	store := &model.Store{}

	err = r.db.WithContext(ctx).
		Where("id = ?", id).
		First(store).
		Error

	res = store.ToStoreDomain()
	return
}

func (r *storeRepository) FindByName(ctx context.Context, name string) (res *domain.Store, err error) {
	store := &model.Store{}

	err = r.db.WithContext(ctx).
		Where("name = ?", name).
		First(store).
		Error

	res = store.ToStoreDomain()
	return
}
func (r *storeRepository) FindAll(ctx context.Context, sort string, limit, page int) (res []*domain.Store, total int64, err error) {
	var stores []*model.Store

	err = r.db.WithContext(ctx).
		Offset((page - 1) * limit).
		Limit(limit).
		Order(sort).
		Find(&stores).
		Count(&total).Error

	res = make([]*domain.Store, 0)
	for _, value := range stores {
		res = append(res, value.ToStoreDomain())
	}
	return
}

func (r *storeRepository) Update(ctx context.Context, store *domain.Store) (err error) {
	storeEntity := &model.Store{}
	storeEntity.FromStoreDomain(store)

	err = r.db.WithContext(ctx).
		Save(storeEntity).
		Error

	return
}

func (r *storeRepository) Delete(ctx context.Context, id string) (err error) {
	err = r.db.WithContext(ctx).
		Delete(&model.Store{}, "id = ?", id).
		Error

	return
}
