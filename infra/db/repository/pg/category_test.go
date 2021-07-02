package pg_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/infra/db/repository/pg"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func Test_CategoryRepo_Create(t *testing.T) {
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
			err = repo.Create(context.TODO(), tc.arg)
			tc.checkResponse(t, err)
		})
	}
}

func Test_CategoryRepo_GetAll(t *testing.T) {
	type args struct {
		page  int
		limit int
		sort  string
	}
	a := args{
		page:  1,
		limit: 10,
		sort:  "created_at",
	}
	c := getCategory()
	testCases := []struct {
		name          string
		args          args
		builtSts      func(mock sqlmock.Sqlmock)
		checkResponse func(t *testing.T, res []*domain.Category, count int64, err error)
	}{
		{
			name: "fail on get all",
			args: a,
			builtSts: func(mock sqlmock.Sqlmock) {
				offset := (a.page - 1) * a.limit
				selectQuery := fmt.Sprintf(`SELECT * FROM "categories" ORDER BY %s OFFSET %d LIMIT %d`, a.sort, offset, a.limit)
				mock.ExpectQuery(regexp.QuoteMeta(selectQuery)).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, res []*domain.Category, count int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, count, int64(0))
				assert.Len(t, res, 0)
			},
		},
		{
			name: "fail on copy rows",
			args: a,
			builtSts: func(mock sqlmock.Sqlmock) {
				offset := (a.page - 1) * a.limit
				selectQuery := fmt.Sprintf(`SELECT * FROM "categories" ORDER BY %s OFFSET %d LIMIT %d`, a.sort, offset, a.limit)
				row := sqlmock.
					NewRows([]string{"uuid", "cname", "cstatus"}).
					AddRow(c.ID, c.Name, c.Status)

				mock.ExpectQuery(regexp.QuoteMeta(selectQuery)).WillReturnRows(row)
			},
			checkResponse: func(t *testing.T, res []*domain.Category, count int64, err error) {
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
				selectQuery := fmt.Sprintf(`SELECT * FROM "categories" ORDER BY %s OFFSET %d LIMIT %d`, a.sort, offset, a.limit)
				countQuery := `SELECT count(1) FROM categories`
				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status"}).
					AddRow(c.ID, c.CreatedAt, c.UpdatedAt, c.Name, c.Status)
				mock.ExpectQuery(regexp.QuoteMeta(selectQuery)).WillReturnRows(row)
				mock.ExpectQuery(regexp.QuoteMeta(countQuery)).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, res []*domain.Category, count int64, err error) {
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
				selectQuery := fmt.Sprintf(`SELECT * FROM "categories" ORDER BY %s OFFSET %d LIMIT %d`, a.sort, offset, a.limit)
				countQuery := `SELECT count(1) FROM categories`
				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status"}).
					AddRow(c.ID, c.CreatedAt, c.UpdatedAt, c.Name, c.Status)
				countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(selectQuery)).WillReturnRows(row)
				mock.ExpectQuery(regexp.QuoteMeta(countQuery)).WillReturnRows(countRow)
			},
			checkResponse: func(t *testing.T, res []*domain.Category, count int64, err error) {
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
			repo := pg.NewCategoryRepository(db)
			tc.builtSts(mock)
			res, total, err := repo.GetAll(context.TODO(), tc.args.sort, tc.args.page, tc.args.limit)
			tc.checkResponse(t, res, total, err)
		})
	}
}

func Test_CategoryRepo_GetAllByStatus(t *testing.T) {
	type args struct {
		status string
		sort   string
		page   int
		limit  int
	}

	a := args{
		status: domain.StoreStatusActive,
		page:   1,
		limit:  10,
		sort:   "created_at",
	}
	c := getCategory()

	testCases := []struct {
		name          string
		args          args
		builtSts      func(mock sqlmock.Sqlmock)
		checkResponse func(t *testing.T, res []*domain.Category, count int64, err error)
	}{
		{
			name: "fail on get all",
			args: a,
			builtSts: func(mock sqlmock.Sqlmock) {
				offset := (a.page - 1) * a.limit
				selectQuery := fmt.Sprintf(`SELECT * FROM categories WHERE status = $1 ORDER BY %s OFFSET %d LIMIT %d`, a.sort, offset, a.limit)
				mock.ExpectQuery(regexp.QuoteMeta(selectQuery)).WithArgs(a.status).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, res []*domain.Category, count int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, count, int64(0))
				assert.Len(t, res, 0)
			},
		},
		{
			name: "fail on get count",
			args: a,
			builtSts: func(mock sqlmock.Sqlmock) {
				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status"}).
					AddRow(c.ID, c.CreatedAt, c.UpdatedAt, c.Name, c.Status)
				offset := (a.page - 1) * a.limit
				selectQuery := fmt.Sprintf(`SELECT * FROM categories WHERE status = $1 ORDER BY %s OFFSET %d LIMIT %d`, a.sort, offset, a.limit)
				mock.ExpectQuery(regexp.QuoteMeta(selectQuery)).WithArgs(a.status).WillReturnRows(row)

				countQuery := `SELECT count(1) FROM categories WHERE status = $1`
				mock.ExpectQuery(regexp.QuoteMeta(countQuery)).WithArgs(a.status).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, res []*domain.Category, count int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, count, int64(0))
				assert.Len(t, res, 0)
			},
		},
		{
			name: "success",
			args: a,
			builtSts: func(mock sqlmock.Sqlmock) {
				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status"}).
					AddRow(c.ID, c.CreatedAt, c.UpdatedAt, c.Name, c.Status)
				offset := (a.page - 1) * a.limit
				selectQuery := fmt.Sprintf(`SELECT * FROM categories WHERE status = $1 ORDER BY %s OFFSET %d LIMIT %d`, a.sort, offset, a.limit)
				mock.ExpectQuery(regexp.QuoteMeta(selectQuery)).WithArgs(a.status).WillReturnRows(row)
				countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)

				countQuery := `SELECT count(1) FROM categories WHERE status = $1`
				mock.ExpectQuery(regexp.QuoteMeta(countQuery)).WithArgs(a.status).WillReturnRows(countRow)
			},
			checkResponse: func(t *testing.T, res []*domain.Category, count int64, err error) {
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
			repo := pg.NewCategoryRepository(db)
			tc.builtSts(mock)
			res, total, err := repo.GetAllByStatus(context.TODO(), tc.args.status, tc.args.sort, tc.args.page, tc.args.limit)
			tc.checkResponse(t, res, total, err)
		})
	}
}

func Test_CategoryRepo_GetById(t *testing.T) {
	id := uuid.NewV4().String()
	c := getCategory()
	testCases := []struct {
		name          string
		arg           string
		builtSts      func(mock sqlmock.Sqlmock)
		checkResponse func(t *testing.T, res *domain.Category, err error)
	}{
		{
			name: "fail on get all",
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
			name: "returns empty row",
			arg:  id,
			builtSts: func(mock sqlmock.Sqlmock) {
				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status"})
				query := `SELECT * FROM categories WHERE id = $1 ORDER BY id LIMIT 1`
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(id).WillReturnRows(row)
			},
			checkResponse: func(t *testing.T, res *domain.Category, err error) {
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
			res, err := repo.GetById(context.TODO(), tc.arg)
			tc.checkResponse(t, res, err)
		})
	}
}

func Test_CategoryRepo_GetByIdAndStatus(t *testing.T) {
	c := getCategory()
	type args struct {
		id     string
		status string
	}
	ma := args{
		id:     uuid.NewV4().String(),
		status: domain.CategoryStatusActive,
	}
	testCases := []struct {
		name          string
		args          args
		builtSts      func(mock sqlmock.Sqlmock)
		checkResponse func(t *testing.T, res *domain.Category, err error)
	}{
		{
			name: "fail on get all",
			args: ma,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `SELECT * FROM categories WHERE id = $1 AND status = $2 ORDER BY id LIMIT 1`
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(ma.id, ma.status).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, res *domain.Category, err error) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "returns empty row",
			args: ma,
			builtSts: func(mock sqlmock.Sqlmock) {
				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status"})
				query := `SELECT * FROM categories WHERE id = $1 AND status = $2 ORDER BY id LIMIT 1`
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(ma.id, ma.status).WillReturnRows(row)
			},
			checkResponse: func(t *testing.T, res *domain.Category, err error) {
				assert.Error(t, err)
				assert.Nil(t, res)
				assert.EqualError(t, err, domain.ErrNotFound.Error())
			},
		},
		{
			name: "success",
			args: ma,
			builtSts: func(mock sqlmock.Sqlmock) {
				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "name", "status"}).
					AddRow(c.ID, c.CreatedAt, c.UpdatedAt, c.Name, c.Status)
				query := `SELECT * FROM categories WHERE id = $1 AND status = $2 ORDER BY id LIMIT 1`
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(ma.id, ma.status).WillReturnRows(row)
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
			res, err := repo.GetByIdAndStatus(context.TODO(), tc.args.id, tc.args.status)
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
