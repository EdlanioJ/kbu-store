package pg_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/infra/db/repository/pg"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func getStore() *domain.Store {
	c := new(domain.Category)
	c.ID = uuid.NewV4().String()

	store := &domain.Store{
		Base: domain.Base{
			ID:        uuid.NewV4().String(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Time{},
		},
		Name:        "Store 001",
		Description: "Store description 001",
		Status:      domain.StoreStatusActive,
		ExternalID:  uuid.NewV4().String(),
		AccountID:   uuid.NewV4().String(),
		Category:    c,
		Tags:        []string{"tag 001", "tag 002"},
		Position: domain.Position{
			Lat: -8.8867698,
			Lng: 13.4771186,
		},
	}
	return store
}

func Test_StoreRepo_Create(t *testing.T) {
	s := getStore()
	testCases := []struct {
		name          string
		arg           *domain.Store
		builtSts      func(mock sqlmock.Sqlmock)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "fail on exec query",
			arg:  s,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO stores (id,created_at,updated_at,name,description,status,external_id,account_id,category_id,tags,lat,lng) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Description, s.Status, s.ExternalID, s.AccountID, s.Category.ID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail on get affected row",
			arg:  s,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO stores (id,created_at,updated_at,name,description,status,external_id,account_id,category_id,tags,lat,lng) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Description, s.Status, s.ExternalID, s.AccountID, s.Category.ID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng).WillReturnResult(sqlmock.NewErrorResult(errors.New("unexpected error")))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail on invalid number of affected row",
			arg:  s,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO stores (id,created_at,updated_at,name,description,status,external_id,account_id,category_id,tags,lat,lng) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Description, s.Status, s.ExternalID, s.AccountID, s.Category.ID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng).WillReturnResult(sqlmock.NewResult(1, 2))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			arg:  s,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO stores (id,created_at,updated_at,name,description,status,external_id,account_id,category_id,tags,lat,lng) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Description, s.Status, s.ExternalID, s.AccountID, s.Category.ID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			repo := pg.NewStoreRepository(db)
			tc.builtSts(mock)
			err = repo.Create(context.TODO(), tc.arg)
			tc.checkResponse(t, err)
		})
	}
}

func Test_StoreRepo_GetById(t *testing.T) {
	id := uuid.NewV4().String()
	s := getStore()
	testCases := []struct {
		name          string
		arg           string
		builtSts      func(mock sqlmock.Sqlmock)
		checkResponse func(t *testing.T, res *domain.Store, err error)
	}{
		{
			name: "fail on get list",
			arg:  id,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `SELECT * FROM stores WHERE id = $1 ORDER BY id`
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(id).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, res *domain.Store, err error) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "fail on copy rows",
			arg:  id,
			builtSts: func(mock sqlmock.Sqlmock) {
				row := sqlmock.
					NewRows([]string{"uuid", "s_created_at", "s_updated_at", "s_name", "s_status", "s_lng"}).
					AddRow(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Status, s.Position.Lng)

				query := `SELECT * FROM stores WHERE id = $1 ORDER BY id`
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(id).WillReturnRows(row)
			},
			checkResponse: func(t *testing.T, res *domain.Store, err error) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "fail on returns empty",
			arg:  id,
			builtSts: func(mock sqlmock.Sqlmock) {
				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "external_id", "tags", "lat", "lng"})

				query := `SELECT * FROM stores WHERE id = $1 ORDER BY id`
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(id).WillReturnRows(row)
			},
			checkResponse: func(t *testing.T, res *domain.Store, err error) {
				assert.Error(t, err)
				assert.Nil(t, res)
				assert.EqualError(t, err, domain.ErrNotFound.Error())
			},
		},
		{
			name: "success",
			arg:  id,
			builtSts: func(mock sqlmock.Sqlmock) {
				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "external_id", "tags", "lat", "lng"}).
					AddRow(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Status, s.Description, s.AccountID, s.Category.ID, s.ExternalID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng)

				query := `SELECT * FROM stores WHERE id = $1 ORDER BY id`
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(id).WillReturnRows(row)
			},
			checkResponse: func(t *testing.T, res *domain.Store, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res, s)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			repo := pg.NewStoreRepository(db)
			tc.builtSts(mock)
			res, err := repo.GetById(context.TODO(), tc.arg)
			tc.checkResponse(t, res, err)
		})
	}
}

func Test_StoreRepo_GetByIdAndOwner(t *testing.T) {
	type args struct {
		id    string
		owner string
	}

	s := getStore()
	a := args{
		id:    uuid.NewV4().String(),
		owner: uuid.NewV4().String(),
	}
	testCases := []struct {
		name          string
		args          args
		builtSts      func(mock sqlmock.Sqlmock)
		checkResponse func(t *testing.T, res *domain.Store, err error)
	}{
		{
			name: "fail on get all",
			args: a,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `SELECT * FROM stores WHERE id = $1 AND external_id = $2 ORDER BY id LIMIT 1`
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(a.id, a.owner).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, res *domain.Store, err error) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "fail on returns empty",
			args: a,
			builtSts: func(mock sqlmock.Sqlmock) {
				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "external_id", "tags", "lat", "lng"})

				query := `SELECT * FROM stores WHERE id = $1 AND external_id = $2 ORDER BY id LIMIT 1`
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(a.id, a.owner).WillReturnRows(row)
			},
			checkResponse: func(t *testing.T, res *domain.Store, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, domain.ErrNotFound.Error())
				assert.Nil(t, res)
			},
		},
		{
			name: "success",
			args: a,
			builtSts: func(mock sqlmock.Sqlmock) {
				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "external_id", "tags", "lat", "lng"}).
					AddRow(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Status, s.Description, s.AccountID, s.Category.ID, s.ExternalID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng)
				query := `SELECT * FROM stores WHERE id = $1 AND external_id = $2 ORDER BY id LIMIT 1`
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(a.id, a.owner).WillReturnRows(row)
			},
			checkResponse: func(t *testing.T, res *domain.Store, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res, s)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			repo := pg.NewStoreRepository(db)
			tc.builtSts(mock)
			res, err := repo.GetByIdAndOwner(context.TODO(), tc.args.id, tc.args.owner)
			tc.checkResponse(t, res, err)
		})
	}
}

func Test_StoreRepo_GetAll(t *testing.T) {
	type args struct {
		page  int
		limit int
		sort  string
	}

	s := getStore()
	a := args{
		page:  1,
		limit: 10,
		sort:  "created_at",
	}
	testCases := []struct {
		name          string
		args          args
		builtSts      func(mock sqlmock.Sqlmock)
		checkResponse func(t *testing.T, res []*domain.Store, count int64, err error)
	}{
		{
			name: "fail on get all",
			args: a,
			builtSts: func(mock sqlmock.Sqlmock) {
				offset := (a.page - 1) * a.limit
				query := fmt.Sprintf(`SELECT * FROM stores ORDER BY %s OFFSET %d LIMIT %d`, a.sort, offset, a.limit)
				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, count, int64(0))
				assert.Len(t, res, 0)
			},
		},
		{
			name: "fail on get count",
			args: a,
			builtSts: func(mock sqlmock.Sqlmock) {
				offset := (a.page - 1) * a.limit
				query := fmt.Sprintf(`SELECT * FROM stores ORDER BY %s OFFSET %d LIMIT %d`, a.sort, offset, a.limit)
				countQuery := `SELECT count(1) FROM stores`

				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "external_id", "tags", "lat", "lng"}).
					AddRow(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Status, s.Description, s.AccountID, s.Category.ID, s.ExternalID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng)

				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(row)
				mock.ExpectQuery(regexp.QuoteMeta(countQuery)).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, count, int64(0))
				assert.Len(t, res, 0)
			},
		},
		{
			name: "success",
			args: a,
			builtSts: func(mock sqlmock.Sqlmock) {
				offset := (a.page - 1) * a.limit
				query := fmt.Sprintf(`SELECT * FROM stores ORDER BY %s OFFSET %d LIMIT %d`, a.sort, offset, a.limit)
				countQuery := `SELECT count(1) FROM stores`

				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "external_id", "tags", "lat", "lng"}).
					AddRow(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Status, s.Description, s.AccountID, s.Category.ID, s.ExternalID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng)

				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(row)
				countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(countQuery)).WillReturnRows(countRow)
			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, count, int64(1))
				assert.Len(t, res, 1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			repo := pg.NewStoreRepository(db)
			tc.builtSts(mock)
			res, count, err := repo.GetAll(context.TODO(), tc.args.sort, tc.args.limit, tc.args.page)
			tc.checkResponse(t, res, count, err)
		})
	}
}

func Test_StoreRepo_GetAllByCategory(t *testing.T) {
	type args struct {
		page       int
		limit      int
		sort       string
		categoryId string
	}

	s := getStore()
	a := args{
		page:       1,
		limit:      10,
		sort:       "created_at",
		categoryId: uuid.NewV4().String(),
	}
	testCases := []struct {
		name          string
		args          args
		builtSts      func(mock sqlmock.Sqlmock)
		checkResponse func(t *testing.T, res []*domain.Store, count int64, err error)
	}{
		{
			name: "fail on get all",
			args: a,
			builtSts: func(mock sqlmock.Sqlmock) {
				offset := (a.page - 1) * a.limit
				query := fmt.Sprintf(`SELECT * FROM stores WHERE category_id = $1 ORDER BY %s OFFSET %d LIMIT %d`, a.sort, offset, a.limit)
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(a.categoryId).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, count, int64(0))
				assert.Len(t, res, 0)
			},
		},
		{
			name: "fail on get count",
			args: a,
			builtSts: func(mock sqlmock.Sqlmock) {
				offset := (a.page - 1) * a.limit
				query := fmt.Sprintf(`SELECT * FROM stores WHERE category_id = $1 ORDER BY %s OFFSET %d LIMIT %d`, a.sort, offset, a.limit)
				queryCount := `SELECT count(1) FROM stores WHERE category_id = $1`

				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "external_id", "tags", "lat", "lng"}).
					AddRow(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Status, s.Description, s.AccountID, s.Category.ID, s.ExternalID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng)

				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(a.categoryId).WillReturnRows(row)
				mock.ExpectQuery(regexp.QuoteMeta(queryCount)).WithArgs(a.categoryId).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, count, int64(0))
				assert.Len(t, res, 0)
			},
		},
		{
			name: "success",
			args: a,
			builtSts: func(mock sqlmock.Sqlmock) {
				offset := (a.page - 1) * a.limit
				query := fmt.Sprintf(`SELECT * FROM stores WHERE category_id = $1 ORDER BY %s OFFSET %d LIMIT %d`, a.sort, offset, a.limit)
				queryCount := `SELECT count(1) FROM stores WHERE category_id = $1`

				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "external_id", "tags", "lat", "lng"}).
					AddRow(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Status, s.Description, s.AccountID, s.Category.ID, s.ExternalID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng)

				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(a.categoryId).WillReturnRows(row)
				countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(queryCount)).WithArgs(a.categoryId).WillReturnRows(countRow)
			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, count, int64(1))
				assert.Len(t, res, 1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			repo := pg.NewStoreRepository(db)
			tc.builtSts(mock)
			res, count, err := repo.GetAllByCategory(context.TODO(), tc.args.categoryId, tc.args.sort, tc.args.limit, tc.args.page)
			tc.checkResponse(t, res, count, err)
		})
	}
}

func Test_StoreRepo_GetAllByOwner(t *testing.T) {
	type args struct {
		page    int
		limit   int
		sort    string
		ownerId string
	}

	s := getStore()
	a := args{
		page:    1,
		limit:   10,
		sort:    "created_at",
		ownerId: uuid.NewV4().String(),
	}
	testCases := []struct {
		name          string
		args          args
		builtSts      func(mock sqlmock.Sqlmock)
		checkResponse func(t *testing.T, res []*domain.Store, count int64, err error)
	}{
		{
			name: "fail on get all",
			args: a,
			builtSts: func(mock sqlmock.Sqlmock) {
				offset := (a.page - 1) * a.limit
				query := fmt.Sprintf(`SELECT * FROM stores WHERE external_id = $1 ORDER BY %s OFFSET %d LIMIT %d`, a.sort, offset, a.limit)
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(a.ownerId).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, count, int64(0))
				assert.Len(t, res, 0)
			},
		},
		{
			name: "fail on get count",
			args: a,
			builtSts: func(mock sqlmock.Sqlmock) {
				offset := (a.page - 1) * a.limit
				query := fmt.Sprintf(`SELECT * FROM stores WHERE external_id = $1 ORDER BY %s OFFSET %d LIMIT %d`, a.sort, offset, a.limit)
				queryCount := `SELECT count(1) FROM stores WHERE external_id = $1`

				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "external_id", "tags", "lat", "lng"}).
					AddRow(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Status, s.Description, s.AccountID, s.Category.ID, s.ExternalID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng)

				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(a.ownerId).WillReturnRows(row)
				mock.ExpectQuery(regexp.QuoteMeta(queryCount)).WithArgs(a.ownerId).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, count, int64(0))
				assert.Len(t, res, 0)
			},
		},
		{
			name: "success",
			args: a,
			builtSts: func(mock sqlmock.Sqlmock) {
				offset := (a.page - 1) * a.limit
				query := fmt.Sprintf(`SELECT * FROM stores WHERE external_id = $1 ORDER BY %s OFFSET %d LIMIT %d`, a.sort, offset, a.limit)
				queryCount := `SELECT count(1) FROM stores WHERE external_id = $1`

				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "external_id", "tags", "lat", "lng"}).
					AddRow(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Status, s.Description, s.AccountID, s.Category.ID, s.ExternalID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng)

				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(a.ownerId).WillReturnRows(row)
				countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(queryCount)).WithArgs(a.ownerId).WillReturnRows(countRow)
			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, count, int64(1))
				assert.Len(t, res, 1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			repo := pg.NewStoreRepository(db)
			tc.builtSts(mock)
			res, count, err := repo.GetAllByOwner(context.TODO(), tc.args.ownerId, tc.args.sort, tc.args.limit, tc.args.page)
			tc.checkResponse(t, res, count, err)
		})
	}
}

func Test_StoreRepo_GetAllByStatus(t *testing.T) {
	type args struct {
		page   int
		limit  int
		sort   string
		status string
	}

	s := getStore()
	a := args{
		page:   1,
		limit:  10,
		sort:   "created_at",
		status: domain.StoreStatusActive,
	}
	testCases := []struct {
		name          string
		args          args
		builtSts      func(mock sqlmock.Sqlmock)
		checkResponse func(t *testing.T, res []*domain.Store, count int64, err error)
	}{
		{
			name: "fail on get all",
			args: a,
			builtSts: func(mock sqlmock.Sqlmock) {
				offset := (a.page - 1) * a.limit
				query := fmt.Sprintf(`SELECT * FROM stores WHERE status = $1 ORDER BY %s OFFSET %d LIMIT %d`, a.sort, offset, a.limit)
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(a.status).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, count, int64(0))
				assert.Len(t, res, 0)
			},
		},
		{
			name: "fail on get count",
			args: a,
			builtSts: func(mock sqlmock.Sqlmock) {
				offset := (a.page - 1) * a.limit
				query := fmt.Sprintf(`SELECT * FROM stores WHERE status = $1 ORDER BY %s OFFSET %d LIMIT %d`, a.sort, offset, a.limit)
				queryCount := `SELECT count(1) FROM stores WHERE status = $1`

				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "external_id", "tags", "lat", "lng"}).
					AddRow(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Status, s.Description, s.AccountID, s.Category.ID, s.ExternalID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng)

				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(a.status).WillReturnRows(row)
				mock.ExpectQuery(regexp.QuoteMeta(queryCount)).WithArgs(a.status).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, count, int64(0))
				assert.Len(t, res, 0)
			},
		},
		{
			name: "success",
			args: a,
			builtSts: func(mock sqlmock.Sqlmock) {
				offset := (a.page - 1) * a.limit
				query := fmt.Sprintf(`SELECT * FROM stores WHERE status = $1 ORDER BY %s OFFSET %d LIMIT %d`, a.sort, offset, a.limit)
				queryCount := `SELECT count(1) FROM stores WHERE status = $1`

				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "external_id", "tags", "lat", "lng"}).
					AddRow(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Status, s.Description, s.AccountID, s.Category.ID, s.ExternalID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng)

				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(a.status).WillReturnRows(row)
				countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(queryCount)).WithArgs(a.status).WillReturnRows(countRow)
			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, count, int64(1))
				assert.Len(t, res, 1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			repo := pg.NewStoreRepository(db)
			tc.builtSts(mock)
			res, count, err := repo.GetAllByStatus(context.TODO(), tc.args.status, tc.args.sort, tc.args.limit, tc.args.page)
			tc.checkResponse(t, res, count, err)
		})
	}
}

func Test_StoreRepo_GetAllByLocation(t *testing.T) {
	type args struct {
		page     int
		limit    int
		sort     string
		status   string
		distance int
		lat      float64
		lng      float64
	}

	s := getStore()
	a := args{
		page:     1,
		limit:    10,
		sort:     "created_at",
		status:   domain.StoreStatusActive,
		distance: 10,
		lat:      -8.8867698,
		lng:      13.4771186,
	}
	testCases := []struct {
		name          string
		args          args
		builtSts      func(mock sqlmock.Sqlmock)
		checkResponse func(t *testing.T, res []*domain.Store, count int64, err error)
	}{
		{
			name: "fail on get all",
			args: a,
			builtSts: func(mock sqlmock.Sqlmock) {
				offset := (a.page - 1) * a.limit
				query := fmt.Sprintf(`SELECT * FROM (SELECT *,(((acos(sin((%[1]v * pi() / 180)) * sin((lat * pi() / 180)) + cos((%[1]v * pi() / 180)) * cos((lat * pi() / 180)) * cos(((%[2]v - lng) * pi() / 180)))) * 180 / pi()) * 60 * 1.1515 * 1.609344 ) AS distance FROM stores) stores WHERE distance <= $1 AND status = $2 ORDER BY %[3]v LIMIT %[4]v OFFSET %[5]v`, a.lat, a.lng, a.sort, a.limit, offset)
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(a.status).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, count, int64(0))
				assert.Len(t, res, 0)
			},
		},
		{
			name: "fail on get count",
			args: a,
			builtSts: func(mock sqlmock.Sqlmock) {
				offset := (a.page - 1) * a.limit
				query := fmt.Sprintf(`SELECT * FROM (SELECT *,(((acos(sin((%[1]v * pi() / 180)) * sin((lat * pi() / 180)) + cos((%[1]v * pi() / 180)) * cos((lat * pi() / 180)) * cos(((%[2]v - lng) * pi() / 180)))) * 180 / pi()) * 60 * 1.1515 * 1.609344 ) AS distance FROM stores) stores WHERE distance <= $1 AND status = $2 ORDER BY %[3]v LIMIT %[4]v OFFSET %[5]v`, a.lat, a.lng, a.sort, a.limit, offset)
				queryCount := fmt.Sprintf(`SELECT count(1) FROM (SELECT *,(((acos(sin((%[1]v * pi() / 180)) * sin((lat * pi() / 180)) + cos((%[1]v * pi() / 180)) * cos((lat * pi() / 180)) * cos(((%[2]v - lng) * pi() / 180)))) * 180 / pi()) * 60 * 1.1515 * 1.609344 ) AS distance FROM stores) stores WHERE distance <= $1 AND status = $2`, a.lat, a.lng)
				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "external_id", "tags", "lat", "lng"}).
					AddRow(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Status, s.Description, s.AccountID, s.Category.ID, s.ExternalID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng)

				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(a.distance, a.status).WillReturnRows(row)
				mock.ExpectQuery(regexp.QuoteMeta(queryCount)).WithArgs(a.distance, a.status).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, count, int64(0))
				assert.Len(t, res, 0)
			},
		},
		{
			name: "success",
			args: a,
			builtSts: func(mock sqlmock.Sqlmock) {
				offset := (a.page - 1) * a.limit
				query := fmt.Sprintf(`SELECT * FROM (SELECT *,(((acos(sin((%[1]v * pi() / 180)) * sin((lat * pi() / 180)) + cos((%[1]v * pi() / 180)) * cos((lat * pi() / 180)) * cos(((%[2]v - lng) * pi() / 180)))) * 180 / pi()) * 60 * 1.1515 * 1.609344 ) AS distance FROM stores) stores WHERE distance <= $1 AND status = $2 ORDER BY %[3]v LIMIT %[4]v OFFSET %[5]v`, a.lat, a.lng, a.sort, a.limit, offset)
				queryCount := fmt.Sprintf(`SELECT count(1) FROM (SELECT *,(((acos(sin((%[1]v * pi() / 180)) * sin((lat * pi() / 180)) + cos((%[1]v * pi() / 180)) * cos((lat * pi() / 180)) * cos(((%[2]v - lng) * pi() / 180)))) * 180 / pi()) * 60 * 1.1515 * 1.609344 ) AS distance FROM stores) stores WHERE distance <= $1 AND status = $2`, a.lat, a.lng)
				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "external_id", "tags", "lat", "lng"}).
					AddRow(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Status, s.Description, s.AccountID, s.Category.ID, s.ExternalID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng)

				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(a.distance, a.status).WillReturnRows(row)
				countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(queryCount)).WithArgs(a.distance, a.status).WillReturnRows(countRow)
			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, count, int64(1))
				assert.Len(t, res, 1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			repo := pg.NewStoreRepository(db)
			tc.builtSts(mock)
			res, count, err := repo.GetAllByLocation(context.TODO(), tc.args.lat, tc.args.lng, tc.args.distance, tc.args.limit, tc.args.page, tc.args.status, tc.args.sort)
			tc.checkResponse(t, res, count, err)
		})
	}
}

func Test_StoreRepo_GetAllByTags(t *testing.T) {
	type args struct {
		page  int
		limit int
		sort  string
		tags  []string
	}

	s := getStore()
	a := args{
		page:  1,
		limit: 10,
		sort:  "created_at",
		tags:  []string{"tag 001", "tag 002"},
	}
	testCases := []struct {
		name          string
		args          args
		builtSts      func(mock sqlmock.Sqlmock)
		checkResponse func(t *testing.T, res []*domain.Store, count int64, err error)
	}{
		{
			name: "fail on get all",
			args: a,
			builtSts: func(mock sqlmock.Sqlmock) {
				offset := (a.page - 1) * a.limit
				query := fmt.Sprintf(`SELECT * FROM "stores" WHERE tags && $1 ORDER BY %s OFFSET %d LIMIT %d`, a.sort, offset, a.limit)
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(pq.StringArray(a.tags)).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, count, int64(0))
				assert.Len(t, res, 0)
			},
		},
		{
			name: "fail on get count",
			args: a,
			builtSts: func(mock sqlmock.Sqlmock) {
				offset := (a.page - 1) * a.limit
				query := fmt.Sprintf(`SELECT * FROM "stores" WHERE tags && $1 ORDER BY %s OFFSET %d LIMIT %d`, a.sort, offset, a.limit)
				queryCount := `SELECT count(1) FROM "stores" WHERE tags && $1`

				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "external_id", "tags", "lat", "lng"}).
					AddRow(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Status, s.Description, s.AccountID, s.Category.ID, s.ExternalID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng)

				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(pq.StringArray(a.tags)).WillReturnRows(row)
				mock.ExpectQuery(regexp.QuoteMeta(queryCount)).WithArgs(pq.StringArray(a.tags)).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, count, int64(0))
				assert.Len(t, res, 0)
			},
		},
		{
			name: "success",
			args: a,
			builtSts: func(mock sqlmock.Sqlmock) {
				offset := (a.page - 1) * a.limit
				query := fmt.Sprintf(`SELECT * FROM "stores" WHERE tags && $1 ORDER BY %s OFFSET %d LIMIT %d`, a.sort, offset, a.limit)
				queryCount := `SELECT count(1) FROM "stores" WHERE tags && $1`

				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "external_id", "tags", "lat", "lng"}).
					AddRow(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Status, s.Description, s.AccountID, s.Category.ID, s.ExternalID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng)

				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(pq.StringArray(a.tags)).WillReturnRows(row)
				countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(queryCount)).WithArgs(pq.StringArray(a.tags)).WillReturnRows(countRow)
			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, count, int64(1))
				assert.Len(t, res, 1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			repo := pg.NewStoreRepository(db)
			tc.builtSts(mock)
			res, count, err := repo.GetAllByTags(context.TODO(), tc.args.tags, tc.args.sort, tc.args.limit, tc.args.page)
			tc.checkResponse(t, res, count, err)
		})
	}
}

func Test_StoreRepo_Update(t *testing.T) {
	s := getStore()
	testCases := []struct {
		name          string
		arg           *domain.Store
		builtSts      func(mock sqlmock.Sqlmock)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "fail on exec query",
			arg:  s,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `UPDATE stores SET created_at=$1,updated_at=$2,name=$3,description=$4,status=$5,external_id=$6,account_id=$7,category_id=$8,tags=$9,lat=$10,lng=$11 WHERE id = $12`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(s.CreatedAt, s.UpdatedAt, s.Name, s.Description, s.Status, s.ExternalID, s.AccountID, s.Category.ID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng, s.ID).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail on get affected rows",
			arg:  s,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `UPDATE stores SET created_at=$1,updated_at=$2,name=$3,description=$4,status=$5,external_id=$6,account_id=$7,category_id=$8,tags=$9,lat=$10,lng=$11 WHERE id = $12`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(s.CreatedAt, s.UpdatedAt, s.Name, s.Description, s.Status, s.ExternalID, s.AccountID, s.Category.ID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng, s.ID).WillReturnResult(sqlmock.NewErrorResult(errors.New("unexpected error")))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail on invalid number of afected row",
			arg:  s,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `UPDATE stores SET created_at=$1,updated_at=$2,name=$3,description=$4,status=$5,external_id=$6,account_id=$7,category_id=$8,tags=$9,lat=$10,lng=$11 WHERE id = $12`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(s.CreatedAt, s.UpdatedAt, s.Name, s.Description, s.Status, s.ExternalID, s.AccountID, s.Category.ID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng, s.ID).WillReturnResult(sqlmock.NewResult(1, 2))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			arg:  s,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `UPDATE stores SET created_at=$1,updated_at=$2,name=$3,description=$4,status=$5,external_id=$6,account_id=$7,category_id=$8,tags=$9,lat=$10,lng=$11 WHERE id = $12`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(s.CreatedAt, s.UpdatedAt, s.Name, s.Description, s.Status, s.ExternalID, s.AccountID, s.Category.ID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng, s.ID).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			repo := pg.NewStoreRepository(db)
			tc.builtSts(mock)
			err = repo.Update(context.TODO(), tc.arg)
			tc.checkResponse(t, err)
		})
	}
}

func Test_StoreRepo_Delete(t *testing.T) {
	id := uuid.NewV4().String()

	testCases := []struct {
		name          string
		arg           string
		builtSts      func(mock sqlmock.Sqlmock)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "fail on exec query",
			arg:  id,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `DELETE FROM stores WHERE id = $1`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(id).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail on get affected rows",
			arg:  id,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `DELETE FROM stores WHERE id = $1`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(id).WillReturnResult(sqlmock.NewErrorResult(errors.New("unexpected error")))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail on invalid number of afected row",
			arg:  id,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `DELETE FROM stores WHERE id = $1`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 2))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			arg:  id,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `DELETE FROM stores WHERE id = $1`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			repo := pg.NewStoreRepository(db)
			tc.builtSts(mock)
			err = repo.Delete(context.TODO(), tc.arg)
			tc.checkResponse(t, err)
		})
	}
}
