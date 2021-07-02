package pg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EdlanioJ/kbu-store/domain"
)

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *categoryRepository {
	return &categoryRepository{
		db: db,
	}
}

func (r *categoryRepository) getAll(ctx context.Context, query string, args ...interface{}) (res []*domain.Category, err error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return
	}
	res = make([]*domain.Category, 0)
	for rows.Next() {
		c := &domain.Category{}
		err = rows.Scan(
			&c.ID,
			&c.CreatedAt,
			&c.UpdatedAt,
			&c.Name,
			&c.Status,
		)

		if err != nil {
			return nil, err
		}
		defer rows.Close()
		res = append(res, c)
	}
	return
}
func (r *categoryRepository) Create(ctx context.Context, c *domain.Category) (err error) {
	query := `INSERT INTO categories (id,created_at,updated_at,name,status) VALUES ($1,$2,$3,$4,$5)`
	res, err := r.db.ExecContext(ctx, query, c.ID, c.CreatedAt, c.UpdatedAt, c.Name, c.Status)
	if err != nil {
		return
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return
	}

	if affect != 1 {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", affect)
		return
	}
	return
}

func (r *categoryRepository) GetAll(ctx context.Context, sort string, page, limit int) (res []*domain.Category, total int64, err error) {
	offset := (page - 1) * limit
	countQuery := `SELECT count(1) FROM categories`
	selectQuery := fmt.Sprintf(`SELECT * FROM "categories" ORDER BY %s OFFSET %d LIMIT %d`, sort, offset, limit)

	res, err = r.getAll(ctx, selectQuery)
	if err != nil {
		return
	}
	err = r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		res = make([]*domain.Category, 0)
		return
	}
	return
}

func (r *categoryRepository) GetAllByStatus(ctx context.Context, status, sort string, page, limit int) (res []*domain.Category, total int64, err error) {
	offset := (page - 1) * limit
	countQuery := `SELECT count(1) FROM categories WHERE status = $1`
	query := fmt.Sprintf(`SELECT * FROM categories WHERE status = $1 ORDER BY %s OFFSET %d LIMIT %d`, sort, offset, limit)
	res, err = r.getAll(ctx, query, status)
	if err != nil {
		return
	}

	err = r.db.QueryRowContext(ctx, countQuery, status).Scan(&total)
	if err != nil {
		res = make([]*domain.Category, 0)
		return
	}
	return
}

func (r *categoryRepository) GetById(ctx context.Context, id string) (res *domain.Category, err error) {
	query := `SELECT * FROM categories WHERE id = $1 ORDER BY id LIMIT 1`
	list, err := r.getAll(ctx, query, id)

	if len(list) > 0 {
		res = list[0]
	} else {
		return nil, domain.ErrNotFound
	}
	return
}

func (r *categoryRepository) GetByIdAndStatus(ctx context.Context, id, status string) (res *domain.Category, err error) {
	query := `SELECT * FROM categories WHERE id = $1 AND status = $2 ORDER BY id LIMIT 1`
	list, err := r.getAll(ctx, query, id, status)

	if len(list) > 0 {
		res = list[0]
	} else {
		return nil, domain.ErrNotFound
	}
	return
}

func (r *categoryRepository) Update(ctx context.Context, c *domain.Category) (err error) {
	query := `UPDATE categories SET created_at=$1, updated_at=$2, name=$3, status=$4 WHERE id = $5`
	res, err := r.db.ExecContext(ctx, query, c.CreatedAt, c.UpdatedAt, c.Name, c.Status, c.ID)
	if err != nil {
		return
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", affect)
		return
	}
	return
}
