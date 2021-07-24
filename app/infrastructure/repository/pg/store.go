package pg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EdlanioJ/kbu-store/app/domain"
	"github.com/lib/pq"
)

type storeRepository struct {
	db *sql.DB
}

func NewStoreRepository(db *sql.DB) *storeRepository {
	return &storeRepository{
		db: db,
	}
}

func (r *storeRepository) getAll(ctx context.Context, query string, args ...interface{}) (res []*domain.Store, err error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return
	}
	res = make([]*domain.Store, 0)
	for rows.Next() {
		s := &domain.Store{}

		var lat float64
		var lng float64
		var tags pq.StringArray
		err = rows.Scan(
			&s.ID,
			&s.CreatedAt,
			&s.UpdatedAt,
			&s.Name,
			&s.Status,
			&s.Description,
			&s.AccountID,
			&s.CategoryID,
			&s.UserID,
			&tags,
			&lat,
			&lng,
		)
		if err != nil {
			return nil, err
		}
		s.Tags = tags
		s.Position = domain.Position{
			Lat: lat,
			Lng: lng,
		}
		defer rows.Close()
		res = append(res, s)
	}
	return
}

func (r *storeRepository) Store(ctx context.Context, s *domain.Store) (err error) {
	query := `INSERT INTO stores (id,created_at,updated_at,name,description,status,user_id,account_id,category_id,tags,lat,lng) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	res, err := r.db.ExecContext(ctx, query, s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Description, s.Status, s.UserID, s.AccountID, s.CategoryID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng)
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

func (r *storeRepository) FindByID(ctx context.Context, id string) (res *domain.Store, err error) {
	query := `SELECT * FROM stores WHERE id = $1 ORDER BY id`
	list, err := r.getAll(ctx, query, id)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return nil, domain.ErrNotFound
	}
	return
}

func (r *storeRepository) FindByName(ctx context.Context, name string) (res *domain.Store, err error) {
	query := `SELECT * FROM stores WHERE name = $1 ORDER BY id`
	list, err := r.getAll(ctx, query, name)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return nil, domain.ErrNotFound
	}
	return
}

func (r *storeRepository) FindAll(ctx context.Context, sort string, limit, page int) (res []*domain.Store, total int64, err error) {
	offset := (page - 1) * limit

	query := fmt.Sprintf(`SELECT * FROM stores ORDER BY %s OFFSET %d LIMIT %d`, sort, offset, limit)
	countQuery := `SELECT count(1) FROM stores`

	res, err = r.getAll(ctx, query)
	if err != nil {
		return
	}

	err = r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		res = make([]*domain.Store, 0)
		return
	}

	return
}

func (r *storeRepository) Update(ctx context.Context, s *domain.Store) (err error) {
	query := `UPDATE stores SET created_at=$1,updated_at=$2,name=$3,description=$4,status=$5,user_id=$6,account_id=$7,category_id=$8,tags=$9,lat=$10,lng=$11 WHERE id = $12`
	res, err := r.db.ExecContext(ctx, query, s.CreatedAt, s.UpdatedAt, s.Name, s.Description, s.Status, s.UserID, s.AccountID, s.CategoryID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng, s.ID)
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

func (r *storeRepository) Delete(ctx context.Context, id string) (err error) {
	query := `DELETE FROM stores WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return
	}

	if affect != 1 {
		err = fmt.Errorf("Weird Behavior. Total Affected: %d", affect)
		return
	}
	return
}
