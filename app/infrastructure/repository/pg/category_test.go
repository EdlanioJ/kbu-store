package pg_test

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/app/domain"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/repository/pg"
	"github.com/EdlanioJ/kbu-store/app/utils/sample"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func Test_CategoryRepo_Store(t *testing.T) {
	c := sample.NewCategory()
	testCases := []struct {
		name        string
		arg         *domain.Category
		expectedErr bool
		prepare     func(mock sqlmock.Sqlmock)
	}{
		{
			name:        "failur_exec_query_returns_error",
			arg:         c,
			expectedErr: true,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO categories (id,created_at,updated_at,name,status) VALUES ($1,$2,$3,$4,$5)`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(c.ID, c.CreatedAt, c.UpdatedAt, c.Name, c.Status).WillReturnError(errors.New("unexpected error"))
			},
		},
		{
			name:        "fail_get_affected_row_returns_error",
			arg:         c,
			expectedErr: true,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO categories (id,created_at,updated_at,name,status) VALUES ($1,$2,$3,$4,$5)`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(c.ID, c.CreatedAt, c.UpdatedAt, c.Name, c.Status).WillReturnResult(sqlmock.NewErrorResult(errors.New("unexpected error")))
			},
		},
		{
			name:        "failure_returns_invalid_number_of_affected_row",
			arg:         c,
			expectedErr: true,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO categories (id,created_at,updated_at,name,status) VALUES ($1,$2,$3,$4,$5)`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(c.ID, c.CreatedAt, c.UpdatedAt, c.Name, c.Status).WillReturnResult(sqlmock.NewResult(1, 2))
			},
		},
		{
			name: "success",
			arg:  c,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO categories (id,created_at,updated_at,name,status) VALUES ($1,$2,$3,$4,$5)`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(c.ID, c.CreatedAt, c.UpdatedAt, c.Name, c.Status).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			repo := pg.NewCategoryRepository(db)
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

func Test_CategoryRepo_FindByID(t *testing.T) {
	id := uuid.NewV4().String()
	c := sample.NewCategory()
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
				query := `SELECT * FROM categories WHERE id = $1 ORDER BY id LIMIT 1`
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(id).WillReturnError(errors.New("unexpected error"))
			},
		},
		{
			name: "success",
			arg:  id,
			prepare: func(mock sqlmock.Sqlmock) {
				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status"}).
					AddRow(c.ID, c.CreatedAt, c.UpdatedAt, c.Name, c.Status)
				query := `SELECT * FROM categories WHERE id = $1 ORDER BY id LIMIT 1`
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
			repo := pg.NewCategoryRepository(db)
			tc.prepare(mock)
			res, err := repo.FindByID(context.TODO(), tc.arg)
			if tc.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res, c)
			}
		})
	}
}

func Test_CategoryRepo_Update(t *testing.T) {
	c := sample.NewCategory()
	testCases := []struct {
		name        string
		arg         *domain.Category
		expectedErr bool
		prepare     func(mock sqlmock.Sqlmock)
	}{
		{
			name:        "failure_exec_query_returns_error",
			arg:         c,
			expectedErr: true,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `UPDATE categories SET created_at=$1, updated_at=$2, name=$3, status=$4 WHERE id = $5`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(c.CreatedAt, c.UpdatedAt, c.Name, c.Status, c.ID).WillReturnError(errors.New("unexpected error"))
			},
		},
		{
			name:        "failure_get_affected_row_returns_error",
			arg:         c,
			expectedErr: true,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `UPDATE categories SET created_at=$1, updated_at=$2, name=$3, status=$4 WHERE id = $5`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(c.CreatedAt, c.UpdatedAt, c.Name, c.Status, c.ID).WillReturnResult(sqlmock.NewErrorResult(errors.New("unexpected error")))
			},
		},
		{
			name:        "failure_returns_invalid_number_of_affected_row",
			arg:         c,
			expectedErr: true,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `UPDATE categories SET created_at=$1, updated_at=$2, name=$3, status=$4 WHERE id = $5`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(c.CreatedAt, c.UpdatedAt, c.Name, c.Status, c.ID).WillReturnResult(sqlmock.NewResult(1, 2))
			},
		},
		{
			name: "success",
			arg:  c,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `UPDATE categories SET created_at=$1, updated_at=$2, name=$3, status=$4 WHERE id = $5`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(c.CreatedAt, c.UpdatedAt, c.Name, c.Status, c.ID).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			repo := pg.NewCategoryRepository(db)
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
