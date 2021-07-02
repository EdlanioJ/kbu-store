package pg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EdlanioJ/kbu-store/domain"
)

type tagRepository struct {
	db *sql.DB
}

func NewTagRepository(db *sql.DB) *tagRepository {
	return &tagRepository{
		db: db,
	}
}

func (r *tagRepository) getAll(ctx context.Context, query string, args ...interface{}) (res []*domain.Tag, err error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return
	}
	res = make([]*domain.Tag, 0)
	for rows.Next() {
		t := &domain.Tag{}
		err = rows.Scan(
			&t.Name,
			&t.Count,
		)

		if err != nil {
			return nil, err
		}
		defer rows.Close()
		res = append(res, t)
	}
	return
}

func (r *tagRepository) GetAll(ctx context.Context, sort string, page, limit int) (res []*domain.Tag, total int64, err error) {
	offset := (page - 1) * limit
	query := fmt.Sprintf(`
	SELECT *, count(1) FROM ( SELECT unnest(tags) AS tag FROM stores ) tags GROUP BY tag ORDER BY %s LIMIT %d OFFSET %d
	`, sort, limit, offset)

	queryCount := `SELECT count(1) FROM ( SELECT *, count(1) FROM( SELECT unnest(tags) AS tag FROM stores ) tags GROUP BY tag ) stores`

	res, err = r.getAll(ctx, query)
	if err != nil {
		return
	}
	err = r.db.QueryRowContext(ctx, queryCount).Scan(&total)
	if err != nil {
		res = make([]*domain.Tag, 0)
		return
	}
	return
}

func (r *tagRepository) GetAllByCategory(ctx context.Context, categoryID, sort string, page, limit int) (res []*domain.Tag, total int64, err error) {
	offset := (page - 1) * limit
	query := fmt.Sprintf(`SELECT *, count(1) FROM ( SELECT unnest(tags) AS tag FROM stores WHERE category_id = $1 ) stores GROUP BY tag ORDER BY %s LIMIT %d OFFSET %d
	`, sort, limit, offset)
	queryCount := `SELECT count(1) FROM ( SELECT *, count(1) FROM ( SELECT unnest(tags) AS tag FROM stores WHERE category_id = $1 ) tags GROUP BY tag ) stores`
	res, err = r.getAll(ctx, query, categoryID)
	if err != nil {
		return
	}
	err = r.db.QueryRowContext(ctx, queryCount, categoryID).Scan(&total)
	if err != nil {
		res = make([]*domain.Tag, 0)
		return
	}
	return
}
