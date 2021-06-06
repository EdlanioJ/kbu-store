package gorm

import (
	"context"
	"fmt"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"gorm.io/gorm"
)

type gormFetchTags struct {
	db *gorm.DB
}

func NewGormFetchTags(db *gorm.DB) *gormFetchTags {
	return &gormFetchTags{
		db: db,
	}
}

func (r *gormFetchTags) Exec(ctx context.Context, sort string, page, limit int) (res []*domain.Tag, total int64, err error) {
	var tags []*dto.TagDBModel
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
	for _, value := range tags {
		res = append(res, value.ParserToTagDomain())
	}
	return
}
