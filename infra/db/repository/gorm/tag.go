package gorm

import (
	"context"
	"fmt"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/infra/db/model"
	"gorm.io/gorm"
)

type tagsRepository struct {
	db *gorm.DB
}

func NewTagsRepository(db *gorm.DB) *tagsRepository {
	return &tagsRepository{
		db: db,
	}
}

func (r *tagsRepository) GetAll(ctx context.Context, sort string, page, limit int) (res []*domain.Tag, total int64, err error) {
	var tags []*model.Tag
	offset := (page - 1) * limit
	sql := fmt.Sprintf(`
	SELECT *, count(1) FROM (
		SELECT unnest(tags) AS tag FROM stores
	) tags
		GROUP BY tag ORDER BY %s LIMIT %d OFFSET %d
	`, sort, limit, offset)

	err = r.db.
		Raw(sql).
		Scan(&tags).
		Raw(`
			SELECT count(1) FROM (
				SELECT *, count(1) FROM(
					SELECT unnest(tags) AS tag FROM stores
				) tags
					GROUP BY tag
			) stores
	`).Scan(&total).Error
	res = make([]*domain.Tag, 0)
	for _, value := range tags {
		res = append(res, value.ToTagDomain())
	}
	return
}

func (r *tagsRepository) GetAllByCategory(ctx context.Context, categoryID, sort string, page, limit int) (res []*domain.Tag, total int64, err error) {
	var tags []*model.Tag
	offset := (page - 1) * limit

	sql := fmt.Sprintf(`
	SELECT *, count(1) FROM (
		SELECT unnest(tags) AS tag FROM	stores
		WHERE category_id = ?
	) stores
	GROUP BY tag ORDER BY %s LIMIT %d OFFSET %d
	`, sort, limit, offset)
	err = r.db.
		Raw(sql, categoryID).
		Scan(&tags).
		Raw(`
			SELECT count(1) FROM (
				SELECT *, count(1) FROM (
					SELECT unnest(tags) AS tag FROM	stores
					WHERE category_id = ?
				) tags
				GROUP BY tag
			) stores
	`, categoryID).Scan(&total).Error

	res = make([]*domain.Tag, 0)
	for _, value := range tags {
		res = append(res, value.ToTagDomain())
	}
	return
}
