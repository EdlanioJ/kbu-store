package pg_test

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/infra/db/repository/pg"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func getCategory() *domain.Category {
	id := uuid.NewV4().String()
	name := "Type 001"
	category, _ := domain.NewCategory(id, name, domain.CategoryStatusInactive)

	return category
}

func Test_CategoryRepo_Store(t *testing.T) {
	c := getCategory()
	testCases := []struct {
		name          string
		arg           *domain.Category
		builtSts      func(mock sqlmock.Sqlmock)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "fail on exec query",
			arg:  c,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO categories (id,created_at,updated_at,name,status) VALUES ($1,$2,$3,$4,$5)`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(c.ID, c.CreatedAt, c.UpdatedAt, c.Name, c.Status).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail on get affected row",
			arg:  c,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO categories (id,created_at,updated_at,name,status) VALUES ($1,$2,$3,$4,$5)`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(c.ID, c.CreatedAt, c.UpdatedAt, c.Name, c.Status).WillReturnResult(sqlmock.NewErrorResult(errors.New("unexpected error")))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail on invalid number of affected row",
			arg:  c,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO categories (id,created_at,updated_at,name,status) VALUES ($1,$2,$3,$4,$5)`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(c.ID, c.CreatedAt, c.UpdatedAt, c.Name, c.Status).WillReturnResult(sqlmock.NewResult(1, 2))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			arg:  c,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO categories (id,created_at,updated_at,name,status) VALUES ($1,$2,$3,$4,$5)`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(c.ID, c.CreatedAt, c.UpdatedAt, c.Name, c.Status).WillReturnResult(sqlmock.NewResult(1, 1))
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
			repo := pg.NewCategoryRepository(db)
			tc.builtSts(mock)
			err = repo.Store(context.TODO(), tc.arg)
			tc.checkResponse(t, err)
		})
	}
}

func Test_CategoryRepo_FindByID(t *testing.T) {
	id := uuid.NewV4().String()
	c := getCategory()
	testCases := []struct {
		name          string
		arg           string
		builtSts      func(mock sqlmock.Sqlmock)
		checkResponse func(t *testing.T, res *domain.Category, err error)
	}{
		{
			name: "should fail if exec query returns err",
			arg:  id,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `SELECT * FROM categories WHERE id = $1 ORDER BY id LIMIT 1`
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(id).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, res *domain.Category, err error) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "should succeed",
			arg:  id,
			builtSts: func(mock sqlmock.Sqlmock) {
				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status"}).
					AddRow(c.ID, c.CreatedAt, c.UpdatedAt, c.Name, c.Status)
				query := `SELECT * FROM categories WHERE id = $1 ORDER BY id LIMIT 1`
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(id).WillReturnRows(row)
			},
			checkResponse: func(t *testing.T, res *domain.Category, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res, c)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			repo := pg.NewCategoryRepository(db)
			tc.builtSts(mock)
			res, err := repo.FindByID(context.TODO(), tc.arg)
			tc.checkResponse(t, res, err)
		})
	}
}

func Test_CategoryRepo_Update(t *testing.T) {
	c := getCategory()
	testCases := []struct {
		name          string
		arg           *domain.Category
		builtSts      func(mock sqlmock.Sqlmock)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "fail on exec query",
			arg:  c,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `UPDATE categories SET created_at=$1, updated_at=$2, name=$3, status=$4 WHERE id = $5`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(c.CreatedAt, c.UpdatedAt, c.Name, c.Status, c.ID).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail on get affected row",
			arg:  c,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `UPDATE categories SET created_at=$1, updated_at=$2, name=$3, status=$4 WHERE id = $5`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(c.CreatedAt, c.UpdatedAt, c.Name, c.Status, c.ID).WillReturnResult(sqlmock.NewErrorResult(errors.New("unexpected error")))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail on invalid number of affected row",
			arg:  c,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `UPDATE categories SET created_at=$1, updated_at=$2, name=$3, status=$4 WHERE id = $5`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(c.CreatedAt, c.UpdatedAt, c.Name, c.Status, c.ID).WillReturnResult(sqlmock.NewResult(1, 2))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			arg:  c,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `UPDATE categories SET created_at=$1, updated_at=$2, name=$3, status=$4 WHERE id = $5`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(c.CreatedAt, c.UpdatedAt, c.Name, c.Status, c.ID).WillReturnResult(sqlmock.NewResult(1, 1))
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
			repo := pg.NewCategoryRepository(db)
			tc.builtSts(mock)
			err = repo.Update(context.TODO(), tc.arg)
			tc.checkResponse(t, err)
		})
	}
}
