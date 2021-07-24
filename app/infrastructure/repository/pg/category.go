package pg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EdlanioJ/kbu-store/app/domain"
)

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *categoryRepository {
	return &categoryRepository{
		db: db,
	}
}

func (r *categoryRepository) Store(ctx context.Context, c *domain.Category) (err error) {
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

func (r *categoryRepository) FindByID(ctx context.Context, id string) (res *domain.Category, err error) {
	query := `SELECT * FROM categories WHERE id = $1 ORDER BY id LIMIT 1`
	row := r.db.QueryRowContext(ctx, query, id)

	res = &domain.Category{}
	err = row.Scan(
		&res.ID,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.Name,
		&res.Status,
	)
	if err != nil {
		return nil, err
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
