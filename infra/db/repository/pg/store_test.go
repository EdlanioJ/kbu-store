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
		UserID:      uuid.NewV4().String(),
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
			name: "should fail if exec query returns an error",
			arg:  s,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO stores (id,created_at,updated_at,name,description,status,user_id,account_id,category_id,tags,lat,lng) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Description, s.Status, s.UserID, s.AccountID, s.Category.ID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "should fail if get affected row returns an error",
			arg:  s,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO stores (id,created_at,updated_at,name,description,status,user_id,account_id,category_id,tags,lat,lng) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Description, s.Status, s.UserID, s.AccountID, s.Category.ID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng).WillReturnResult(sqlmock.NewErrorResult(errors.New("unexpected error")))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "should fail if returns an invalid number of affected row",
			arg:  s,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO stores (id,created_at,updated_at,name,description,status,user_id,account_id,category_id,tags,lat,lng) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Description, s.Status, s.UserID, s.AccountID, s.Category.ID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng).WillReturnResult(sqlmock.NewResult(1, 2))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "should succeed",
			arg:  s,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO stores (id,created_at,updated_at,name,description,status,user_id,account_id,category_id,tags,lat,lng) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Description, s.Status, s.UserID, s.AccountID, s.Category.ID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng).WillReturnResult(sqlmock.NewResult(1, 1))
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

func Test_StoreRepo_FindByID(t *testing.T) {
	id := uuid.NewV4().String()
	s := getStore()
	testCases := []struct {
		name          string
		arg           string
		builtSts      func(mock sqlmock.Sqlmock)
		checkResponse func(t *testing.T, res *domain.Store, err error)
	}{
		{
			name: "should fail if get list get returns an error",
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
			name: "should fail if copy rows returns an error",
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
			name: "shoul fail if returns an empty list",
			arg:  id,
			builtSts: func(mock sqlmock.Sqlmock) {
				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "user_id", "tags", "lat", "lng"})

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
			name: "should succeed",
			arg:  id,
			builtSts: func(mock sqlmock.Sqlmock) {
				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "user_id", "tags", "lat", "lng"}).
					AddRow(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Status, s.Description, s.AccountID, s.Category.ID, s.UserID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng)

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
			res, err := repo.FindByID(context.TODO(), tc.arg)
			tc.checkResponse(t, res, err)
		})
	}
}

func Test_StoreRepo_FindByName(t *testing.T) {
	name := "store 001"
	s := getStore()
	testCases := []struct {
		name          string
		arg           string
		builtSts      func(mock sqlmock.Sqlmock)
		checkResponse func(t *testing.T, res *domain.Store, err error)
	}{
		{
			name: "should fail if exec returns an error",
			arg:  name,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `SELECT * FROM stores WHERE name = $1 ORDER BY id`
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(name).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, res *domain.Store, err error) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "should fail if exec returns an empty list",
			arg:  name,
			builtSts: func(mock sqlmock.Sqlmock) {
				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "user_id", "tags", "lat", "lng"})

				query := `SELECT * FROM stores WHERE name = $1 ORDER BY id`
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(name).WillReturnRows(row)
			},
			checkResponse: func(t *testing.T, res *domain.Store, err error) {
				assert.Error(t, err)
				assert.Nil(t, res)
				assert.EqualError(t, err, domain.ErrNotFound.Error())
			},
		},
		{
			name: "should succeed",
			arg:  name,
			builtSts: func(mock sqlmock.Sqlmock) {
				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "user_id", "tags", "lat", "lng"}).
					AddRow(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Status, s.Description, s.AccountID, s.Category.ID, s.UserID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng)

				query := `SELECT * FROM stores WHERE name = $1 ORDER BY id`
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(name).WillReturnRows(row)
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
			res, err := repo.FindByName(context.TODO(), tc.arg)
			tc.checkResponse(t, res, err)
		})
	}
}
func Test_StoreRepo_FindAll(t *testing.T) {
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
			name: "should fail if get all returns an error",
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
			name: "should fail if get count returns an error",
			args: a,
			builtSts: func(mock sqlmock.Sqlmock) {
				offset := (a.page - 1) * a.limit
				query := fmt.Sprintf(`SELECT * FROM stores ORDER BY %s OFFSET %d LIMIT %d`, a.sort, offset, a.limit)
				countQuery := `SELECT count(1) FROM stores`

				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "user_id", "tags", "lat", "lng"}).
					AddRow(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Status, s.Description, s.AccountID, s.Category.ID, s.UserID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng)

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
			name: "should succeed",
			args: a,
			builtSts: func(mock sqlmock.Sqlmock) {
				offset := (a.page - 1) * a.limit
				query := fmt.Sprintf(`SELECT * FROM stores ORDER BY %s OFFSET %d LIMIT %d`, a.sort, offset, a.limit)
				countQuery := `SELECT count(1) FROM stores`

				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "user_id", "tags", "lat", "lng"}).
					AddRow(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Status, s.Description, s.AccountID, s.Category.ID, s.UserID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng)

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
			res, count, err := repo.FindAll(context.TODO(), tc.args.sort, tc.args.limit, tc.args.page)
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
			name: "should fail if exec query returns an error",
			arg:  s,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `UPDATE stores SET created_at=$1,updated_at=$2,name=$3,description=$4,status=$5,user_id=$6,account_id=$7,category_id=$8,tags=$9,lat=$10,lng=$11 WHERE id = $12`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(s.CreatedAt, s.UpdatedAt, s.Name, s.Description, s.Status, s.UserID, s.AccountID, s.Category.ID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng, s.ID).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "should fail if get affected row returns an error",
			arg:  s,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `UPDATE stores SET created_at=$1,updated_at=$2,name=$3,description=$4,status=$5,user_id=$6,account_id=$7,category_id=$8,tags=$9,lat=$10,lng=$11 WHERE id = $12`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(s.CreatedAt, s.UpdatedAt, s.Name, s.Description, s.Status, s.UserID, s.AccountID, s.Category.ID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng, s.ID).WillReturnResult(sqlmock.NewErrorResult(errors.New("unexpected error")))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "should fail if returns an invalid number of affected rows",
			arg:  s,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `UPDATE stores SET created_at=$1,updated_at=$2,name=$3,description=$4,status=$5,user_id=$6,account_id=$7,category_id=$8,tags=$9,lat=$10,lng=$11 WHERE id = $12`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(s.CreatedAt, s.UpdatedAt, s.Name, s.Description, s.Status, s.UserID, s.AccountID, s.Category.ID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng, s.ID).WillReturnResult(sqlmock.NewResult(1, 2))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "should succeed",
			arg:  s,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `UPDATE stores SET created_at=$1,updated_at=$2,name=$3,description=$4,status=$5,user_id=$6,account_id=$7,category_id=$8,tags=$9,lat=$10,lng=$11 WHERE id = $12`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(s.CreatedAt, s.UpdatedAt, s.Name, s.Description, s.Status, s.UserID, s.AccountID, s.Category.ID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng, s.ID).WillReturnResult(sqlmock.NewResult(1, 1))
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
			name: "should fail if exec query returns an error",
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
			name: "should fail if get affected row returns an error",
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
			name: "should fail if returns an invalid number of affected rows",
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
			name: "should succeed",
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
