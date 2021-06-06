package gorm

import (
	"context"
	"fmt"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"gorm.io/gorm"
)

type gormFetchTagsByCategory struct {
	db *gorm.DB
}

func NewGormFetchTagsByCategory(db *gorm.DB) *gormFetchTagsByCategory {
	return &gormFetchTagsByCategory{
		db: db,
	}
}

func (r *gormFetchTagsByCategory) Exec(ctx context.Context, categoryID, sort string, page, limit int) (res []*domain.Tag, total int64, err error) {
	var tags []*dto.TagDBModel
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
	for _, value := range tags {
		res = append(res, value.ParserToTagDomain())
	}
	return
}
