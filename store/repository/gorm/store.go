package gorm

import (
	"context"
	"fmt"

	"github.com/EdlanioJ/kbu-store/db/model"
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/lib/pq"
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
	storeModel.Parser(store)

	err = r.db.WithContext(ctx).
		Create(storeModel).
		Error
	return
}

func (r *storeRepository) GetById(ctx context.Context, id string) (res *domain.Store, err error) {
	store := &model.Store{}

	err = r.db.WithContext(ctx).
		Where("id = ?", id).
		First(store).
		Error

	res = store.ToStoreDomain()
	return
}

func (r *storeRepository) GetByIdAndOwner(ctx context.Context, id string, externalID string) (res *domain.Store, err error) {
	store := &model.Store{}

	err = r.db.WithContext(ctx).
		Where("id = ? AND external_id = ?", id, externalID).
		First(store).
		Error

	res = store.ToStoreDomain()
	return
}

func (r *storeRepository) GetAll(ctx context.Context, sort string, limit, page int) (res []*domain.Store, total int64, err error) {
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

func (r *storeRepository) GetAllByCategory(ctx context.Context, categoryID, sort string, limit, page int) (res []*domain.Store, total int64, err error) {
	var stores []*model.Store

	err = r.db.WithContext(ctx).
		Where("category_id = ?", categoryID).
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

func (r *storeRepository) GetAllByLocation(ctx context.Context, lat, lng float64, distance, limit, page int, status, sort string) (res []*domain.Store, total int64, err error) {
	var stores []*model.Store

	calc := fmt.Sprintf(`(((acos(sin((%[1]v * pi() / 180)) * sin((lat * pi() / 180)) +	cos((%[1]v * pi() / 180)) *	cos((lat * pi() / 180)) * cos(((%[2]v - lng) * pi() / 180))))  * 180 / pi()) * 60 * 1.1515 * 1.609344	) AS distance`, lat, lng)
	subquery := r.db.Table("stores").Select(
		[]string{"*", calc},
	)

	err = r.db.WithContext(ctx).
		Table("(?) stores", subquery).
		Where("distance <= ? AND status = ?", distance, status).
		Offset((page - 1) * limit).
		Limit(limit).
		Order(sort).
		Find(&stores).
		Count(&total).
		Error

	res = make([]*domain.Store, 0)
	for _, value := range stores {
		res = append(res, value.ToStoreDomain())
	}
	return
}

func (r *storeRepository) GetAllByOwner(ctx context.Context, externalID, sort string, limit, page int) (res []*domain.Store, total int64, err error) {
	var stores []*model.Store

	err = r.db.WithContext(ctx).
		Where("external_id = ?", externalID).
		Offset((page - 1) * limit).
		Limit(limit).
		Order(sort).
		Find(&stores).
		Count(&total).
		Error

	res = make([]*domain.Store, 0)
	for _, value := range stores {
		res = append(res, value.ToStoreDomain())
	}
	return
}

func (r *storeRepository) GetAllByStatus(ctx context.Context, status, sort string, limit, page int) (res []*domain.Store, total int64, err error) {
	var stores []*model.Store

	err = r.db.WithContext(ctx).
		Where("status = ?", status).
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

func (r *storeRepository) GetAllByTags(ctx context.Context, tags []string, sort string, limit, page int) (res []*domain.Store, total int64, err error) {
	var stores []*model.Store

	err = r.db.WithContext(ctx).
		Where("tags && ?", pq.StringArray(tags)).
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
	storeEntity.Parser(store)

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
