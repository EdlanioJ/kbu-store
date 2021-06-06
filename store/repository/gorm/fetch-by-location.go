package gorm

import (
	"context"
	"fmt"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"gorm.io/gorm"
)

type gormFetchStoreByCloseLocation struct {
	db *gorm.DB
}

func NewGormFetchStoreByCloseLocation(db *gorm.DB) *gormFetchStoreByCloseLocation {
	return &gormFetchStoreByCloseLocation{
		db: db,
	}
}

func (r *gormFetchStoreByCloseLocation) Exec(ctx context.Context, lat, lng float64, distance, limit, page int, status, sort string) (res []*domain.Store, total int64, err error) {
	var stores []*dto.StoreDBModel

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
		res = append(res, value.ParserToStoreDomain())
	}
	return
}
