package pg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EdlanioJ/kbu-store/domain"
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

		var categoryID string

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
			&categoryID,
			&s.ExternalID,
			&tags,
			&lat,
			&lng,
		)
		if err != nil {
			return nil, err
		}

		c := &domain.Category{}
		c.ID = categoryID
		s.Category = c
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

func (r *storeRepository) Create(ctx context.Context, s *domain.Store) (err error) {
	query := `INSERT INTO stores (id,created_at,updated_at,name,description,status,external_id,account_id,category_id,tags,lat,lng) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	res, err := r.db.ExecContext(ctx, query, s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Description, s.Status, s.ExternalID, s.AccountID, s.Category.ID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng)
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

func (r *storeRepository) GetById(ctx context.Context, id string) (res *domain.Store, err error) {
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

func (r *storeRepository) GetByIdAndOwner(ctx context.Context, id string, externalID string) (res *domain.Store, err error) {
	query := `SELECT * FROM stores WHERE id = $1 AND external_id = $2 ORDER BY id LIMIT 1`
	list, err := r.getAll(ctx, query, id, externalID)
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

func (r *storeRepository) GetAll(ctx context.Context, sort string, limit, page int) (res []*domain.Store, total int64, err error) {
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

func (r *storeRepository) GetAllByCategory(ctx context.Context, categoryID, sort string, limit, page int) (res []*domain.Store, total int64, err error) {
	offset := (page - 1) * limit

	query := fmt.Sprintf(`SELECT * FROM stores WHERE category_id = $1 ORDER BY %s OFFSET %d LIMIT %d`, sort, offset, limit)
	queryCount := `SELECT count(1) FROM stores WHERE category_id = $1`

	res, err = r.getAll(ctx, query, categoryID)
	if err != nil {
		return
	}

	err = r.db.QueryRowContext(ctx, queryCount, categoryID).Scan(&total)
	if err != nil {
		res = make([]*domain.Store, 0)
		return
	}

	return
}

func (r *storeRepository) GetAllByLocation(ctx context.Context, lat, lng float64, distance, limit, page int, status, sort string) (res []*domain.Store, total int64, err error) {
	offset := (page - 1) * limit
	query := fmt.Sprintf(`SELECT * FROM (SELECT *,(((acos(sin((%[1]v * pi() / 180)) * sin((lat * pi() / 180)) + cos((%[1]v * pi() / 180)) * cos((lat * pi() / 180)) * cos(((%[2]v - lng) * pi() / 180)))) * 180 / pi()) * 60 * 1.1515 * 1.609344 ) AS distance FROM stores) stores WHERE distance <= $1 AND status = $2 ORDER BY %[3]v LIMIT %[4]v OFFSET %[5]v`, lat, lng, sort, limit, offset)
	queryCount := fmt.Sprintf(`SELECT count(1) FROM (SELECT *,(((acos(sin((%[1]v * pi() / 180)) * sin((lat * pi() / 180)) + cos((%[1]v * pi() / 180)) * cos((lat * pi() / 180)) * cos(((%[2]v - lng) * pi() / 180)))) * 180 / pi()) * 60 * 1.1515 * 1.609344 ) AS distance FROM stores) stores WHERE distance <= $1 AND status = $2`, lat, lng)
	res, err = r.getAll(ctx, query, distance, status)
	if err != nil {
		return
	}
	err = r.db.QueryRowContext(ctx, queryCount, distance, status).Scan(&total)
	if err != nil {
		res = make([]*domain.Store, 0)
		return
	}
	return
}

func (r *storeRepository) GetAllByOwner(ctx context.Context, externalID, sort string, limit, page int) (res []*domain.Store, total int64, err error) {
	offset := (page - 1) * limit
	query := fmt.Sprintf(`SELECT * FROM stores WHERE external_id = $1 ORDER BY %s OFFSET %d LIMIT %d`, sort, offset, limit)
	queryCount := `SELECT count(1) FROM stores WHERE external_id = $1`

	res, err = r.getAll(ctx, query, externalID)
	if err != nil {
		return
	}
	err = r.db.QueryRowContext(ctx, queryCount, externalID).Scan(&total)
	if err != nil {
		res = make([]*domain.Store, 0)
		return
	}

	return
}

func (r *storeRepository) GetAllByStatus(ctx context.Context, status, sort string, limit, page int) (res []*domain.Store, total int64, err error) {
	offset := (page - 1) * limit
	query := fmt.Sprintf(`SELECT * FROM stores WHERE status = $1 ORDER BY %s OFFSET %d LIMIT %d`, sort, offset, limit)
	queryCount := `SELECT count(1) FROM stores WHERE status = $1`

	res, err = r.getAll(ctx, query, status)
	if err != nil {
		return
	}
	err = r.db.QueryRowContext(ctx, queryCount, status).Scan(&total)
	if err != nil {
		res = make([]*domain.Store, 0)
		return
	}

	return
}

func (r *storeRepository) GetAllByTags(ctx context.Context, tags []string, sort string, limit, page int) (res []*domain.Store, total int64, err error) {
	offset := (page - 1) * limit
	query := fmt.Sprintf(`SELECT * FROM "stores" WHERE tags && $1 ORDER BY %s OFFSET %d LIMIT %d`, sort, offset, limit)
	queryCount := `SELECT count(1) FROM "stores" WHERE tags && $1`

	res, err = r.getAll(ctx, query, pq.StringArray(tags))
	if err != nil {
		return
	}
	err = r.db.QueryRowContext(ctx, queryCount, pq.StringArray(tags)).Scan(&total)
	if err != nil {
		res = make([]*domain.Store, 0)
		return
	}

	return
}

func (r *storeRepository) Update(ctx context.Context, s *domain.Store) (err error) {
	query := `UPDATE stores SET created_at=$1,updated_at=$2,name=$3,description=$4,status=$5,external_id=$6,account_id=$7,category_id=$8,tags=$9,lat=$10,lng=$11 WHERE id = $12`
	res, err := r.db.ExecContext(ctx, query, s.CreatedAt, s.UpdatedAt, s.Name, s.Description, s.Status, s.ExternalID, s.AccountID, s.Category.ID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng, s.ID)
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
