package pg_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/app/domain"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/repository/pg"
	"github.com/EdlanioJ/kbu-store/app/utils/sample"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func Test_StoreRepo_Store(t *testing.T) {
	s := sample.NewStore()
	testCases := []struct {
		name        string
		arg         *domain.Store
		expectedErr bool
		prepare     func(mock sqlmock.Sqlmock)
	}{
		{
			name:        "failure_exec_query_returns_error",
			arg:         s,
			expectedErr: true,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO stores (id,created_at,updated_at,name,description,status,user_id,account_id,category_id,tags,lat,lng) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Description, s.Status, s.UserID, s.AccountID, s.CategoryID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng).WillReturnError(errors.New("unexpected error"))
			},
		},
		{
			name:        "failure_get_affected_row_returns_error",
			arg:         s,
			expectedErr: true,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO stores (id,created_at,updated_at,name,description,status,user_id,account_id,category_id,tags,lat,lng) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Description, s.Status, s.UserID, s.AccountID, s.CategoryID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng).WillReturnResult(sqlmock.NewErrorResult(errors.New("unexpected error")))
			},
		},
		{
			name:        "failure_returns_invalid_number_of_affected_row",
			arg:         s,
			expectedErr: true,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO stores (id,created_at,updated_at,name,description,status,user_id,account_id,category_id,tags,lat,lng) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Description, s.Status, s.UserID, s.AccountID, s.CategoryID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng).WillReturnResult(sqlmock.NewResult(1, 2))
			},
		},
		{
			name: "success",
			arg:  s,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO stores (id,created_at,updated_at,name,description,status,user_id,account_id,category_id,tags,lat,lng) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Description, s.Status, s.UserID, s.AccountID, s.CategoryID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			repo := pg.NewStoreRepository(db)
			tc.prepare(mock)
			err = repo.Store(context.TODO(), tc.arg)
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_StoreRepo_FindByID(t *testing.T) {
	id := uuid.NewV4().String()
	s := sample.NewStore()
	testCases := []struct {
		name        string
		arg         string
		expectedErr bool
		prepare     func(mock sqlmock.Sqlmock)
	}{
		{
			name:        "failure_get_list_get_returns_error",
			arg:         id,
			expectedErr: true,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `SELECT * FROM stores WHERE id = $1 ORDER BY id`
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(id).WillReturnError(errors.New("unexpected error"))
			},
		},
		{
			name:        "failure_copy_rows_returns_error",
			arg:         id,
			expectedErr: true,
			prepare: func(mock sqlmock.Sqlmock) {
				row := sqlmock.
					NewRows([]string{"uuid", "s_created_at", "s_updated_at", "s_name", "s_status", "s_lng"}).
					AddRow(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Status, s.Position.Lng)

				query := `SELECT * FROM stores WHERE id = $1 ORDER BY id`
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(id).WillReturnRows(row)
			},
		},
		{
			name:        "failure_returns_empty_list",
			arg:         id,
			expectedErr: true,
			prepare: func(mock sqlmock.Sqlmock) {
				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "user_id", "tags", "lat", "lng"})

				query := `SELECT * FROM stores WHERE id = $1 ORDER BY id`
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(id).WillReturnRows(row)
			},
		},
		{
			name: "success",
			arg:  id,
			prepare: func(mock sqlmock.Sqlmock) {
				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "user_id", "tags", "lat", "lng"}).
					AddRow(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Status, s.Description, s.AccountID, s.CategoryID, s.UserID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng)

				query := `SELECT * FROM stores WHERE id = $1 ORDER BY id`
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(id).WillReturnRows(row)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			repo := pg.NewStoreRepository(db)
			tc.prepare(mock)
			res, err := repo.FindByID(context.TODO(), tc.arg)
			if tc.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res, s)
			}
		})
	}
}

func Test_StoreRepo_FindByName(t *testing.T) {
	name := "store 001"
	s := sample.NewStore()
	testCases := []struct {
		name        string
		arg         string
		expectedErr bool
		prepare     func(mock sqlmock.Sqlmock)
	}{
		{
			name:        "failure_exec_returns_error",
			arg:         name,
			expectedErr: true,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `SELECT * FROM stores WHERE name = $1 ORDER BY id`
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(name).WillReturnError(errors.New("unexpected error"))
			},
		},
		{
			name:        "failure_exec_returns_empty_list",
			arg:         name,
			expectedErr: true,
			prepare: func(mock sqlmock.Sqlmock) {
				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "user_id", "tags", "lat", "lng"})

				query := `SELECT * FROM stores WHERE name = $1 ORDER BY id`
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(name).WillReturnRows(row)
			},
		},
		{
			name: "success",
			arg:  name,
			prepare: func(mock sqlmock.Sqlmock) {
				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "user_id", "tags", "lat", "lng"}).
					AddRow(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Status, s.Description, s.AccountID, s.CategoryID, s.UserID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng)

				query := `SELECT * FROM stores WHERE name = $1 ORDER BY id`
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(name).WillReturnRows(row)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			repo := pg.NewStoreRepository(db)
			tc.prepare(mock)
			res, err := repo.FindByName(context.TODO(), tc.arg)
			if tc.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res, s)
			}
		})
	}
}
func Test_StoreRepo_FindAll(t *testing.T) {
	s := sample.NewStore()

	page := 1
	limit := 10
	sort := "created_at"

	testCases := []struct {
		name        string
		page        int
		limit       int
		sort        string
		expectedErr bool
		prepare     func(mock sqlmock.Sqlmock)
	}{
		{
			name:        "failure_get_all_returns_error",
			page:        page,
			limit:       limit,
			sort:        sort,
			expectedErr: true,
			prepare: func(mock sqlmock.Sqlmock) {
				offset := (page - 1) * limit
				query := fmt.Sprintf(`SELECT * FROM stores ORDER BY %s OFFSET %d LIMIT %d`, sort, offset, limit)
				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(errors.New("unexpected error"))
			},
		},
		{
			name:        "failure_get_count_returns_error",
			page:        page,
			limit:       limit,
			sort:        sort,
			expectedErr: true,
			prepare: func(mock sqlmock.Sqlmock) {
				offset := (page - 1) * limit
				query := fmt.Sprintf(`SELECT * FROM stores ORDER BY %s OFFSET %d LIMIT %d`, sort, offset, limit)
				countQuery := `SELECT count(1) FROM stores`

				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "user_id", "tags", "lat", "lng"}).
					AddRow(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Status, s.Description, s.AccountID, s.CategoryID, s.UserID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng)

				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(row)
				mock.ExpectQuery(regexp.QuoteMeta(countQuery)).WillReturnError(errors.New("unexpected error"))
			},
		},
		{
			name:  "success",
			page:  page,
			limit: limit,
			sort:  sort,
			prepare: func(mock sqlmock.Sqlmock) {
				offset := (page - 1) * limit
				query := fmt.Sprintf(`SELECT * FROM stores ORDER BY %s OFFSET %d LIMIT %d`, sort, offset, limit)
				countQuery := `SELECT count(1) FROM stores`

				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "user_id", "tags", "lat", "lng"}).
					AddRow(s.ID, s.CreatedAt, s.UpdatedAt, s.Name, s.Status, s.Description, s.AccountID, s.CategoryID, s.UserID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng)

				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(row)
				countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(countQuery)).WillReturnRows(countRow)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			repo := pg.NewStoreRepository(db)
			tc.prepare(mock)
			res, count, err := repo.FindAll(context.TODO(), tc.sort, tc.limit, tc.page)

			if tc.expectedErr {
				assert.Error(t, err)
				assert.Equal(t, count, int64(0))
				assert.Len(t, res, 0)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, count, int64(1))
				assert.Len(t, res, 1)
			}
		})
	}
}

func Test_StoreRepo_Update(t *testing.T) {
	s := sample.NewStore()
	testCases := []struct {
		name        string
		arg         *domain.Store
		expectedErr bool
		prepare     func(mock sqlmock.Sqlmock)
	}{
		{
			name:        "failure_exec_query_returns_error",
			arg:         s,
			expectedErr: true,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `UPDATE stores SET created_at=$1,updated_at=$2,name=$3,description=$4,status=$5,user_id=$6,account_id=$7,category_id=$8,tags=$9,lat=$10,lng=$11 WHERE id = $12`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(s.CreatedAt, s.UpdatedAt, s.Name, s.Description, s.Status, s.UserID, s.AccountID, s.CategoryID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng, s.ID).WillReturnError(errors.New("unexpected error"))
			},
		},
		{
			name:        "failure_get_affected_row_returns_error",
			arg:         s,
			expectedErr: true,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `UPDATE stores SET created_at=$1,updated_at=$2,name=$3,description=$4,status=$5,user_id=$6,account_id=$7,category_id=$8,tags=$9,lat=$10,lng=$11 WHERE id = $12`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(s.CreatedAt, s.UpdatedAt, s.Name, s.Description, s.Status, s.UserID, s.AccountID, s.CategoryID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng, s.ID).WillReturnResult(sqlmock.NewErrorResult(errors.New("unexpected error")))
			},
		},
		{
			name:        "failure_returns_nvalid_number_of_affected_rows",
			arg:         s,
			expectedErr: true,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `UPDATE stores SET created_at=$1,updated_at=$2,name=$3,description=$4,status=$5,user_id=$6,account_id=$7,category_id=$8,tags=$9,lat=$10,lng=$11 WHERE id = $12`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(s.CreatedAt, s.UpdatedAt, s.Name, s.Description, s.Status, s.UserID, s.AccountID, s.CategoryID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng, s.ID).WillReturnResult(sqlmock.NewResult(1, 2))
			},
		},
		{
			name: "success",
			arg:  s,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `UPDATE stores SET created_at=$1,updated_at=$2,name=$3,description=$4,status=$5,user_id=$6,account_id=$7,category_id=$8,tags=$9,lat=$10,lng=$11 WHERE id = $12`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(s.CreatedAt, s.UpdatedAt, s.Name, s.Description, s.Status, s.UserID, s.AccountID, s.CategoryID, pq.StringArray(s.Tags), s.Position.Lat, s.Position.Lng, s.ID).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			repo := pg.NewStoreRepository(db)
			tc.prepare(mock)
			err = repo.Update(context.TODO(), tc.arg)
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_StoreRepo_Delete(t *testing.T) {
	id := uuid.NewV4().String()

	testCases := []struct {
		name        string
		arg         string
		expectedErr bool
		prepare     func(mock sqlmock.Sqlmock)
	}{
		{
			name:        "failure_exec_query_returns_error",
			arg:         id,
			expectedErr: true,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `DELETE FROM stores WHERE id = $1`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(id).WillReturnError(errors.New("unexpected error"))
			},
		},
		{
			name:        "failure_get_affected_row_returns_an error",
			arg:         id,
			expectedErr: true,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `DELETE FROM stores WHERE id = $1`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(id).WillReturnResult(sqlmock.NewErrorResult(errors.New("unexpected error")))
			},
		},
		{
			name:        "failure_returns_invalid_number_of_affected_rows",
			arg:         id,
			expectedErr: true,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `DELETE FROM stores WHERE id = $1`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 2))
			},
		},
		{
			name: "should succeed",
			arg:  id,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `DELETE FROM stores WHERE id = $1`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			repo := pg.NewStoreRepository(db)
			tc.prepare(mock)
			err = repo.Delete(context.TODO(), tc.arg)
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
