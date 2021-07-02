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

func Test_TagRepo_GetAll(t *testing.T) {
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
	tag := new(domain.Tag)
	tag.Name = "tag001"
	tag.Count = 2

	testCases := []struct {
		name          string
		args          args
		builtSts      func(mock sqlmock.Sqlmock)
		checkResponse func(t *testing.T, res []*domain.Tag, count int64, err error)
	}{
		{
			name: "fail on get all",
			args: a,
			builtSts: func(mock sqlmock.Sqlmock) {
				offset := (a.page - 1) * a.limit
				query := fmt.Sprintf(`
	SELECT *, count(1) FROM ( SELECT unnest(tags) AS tag FROM stores ) tags GROUP BY tag ORDER BY %s LIMIT %d OFFSET %d
	`, a.sort, a.limit, offset)
				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, res []*domain.Tag, count int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, count, int64(0))
				assert.Len(t, res, 0)
			},
		},
		{
			name: "fail on copy rows",
			args: a,
			builtSts: func(mock sqlmock.Sqlmock) {
				row := sqlmock.
					NewRows([]string{"name", "total"}).
					AddRow(tag.Count, tag.Name)
				offset := (a.page - 1) * a.limit
				query := fmt.Sprintf(`
	SELECT *, count(1) FROM ( SELECT unnest(tags) AS tag FROM stores ) tags GROUP BY tag ORDER BY %s LIMIT %d OFFSET %d
	`, a.sort, a.limit, offset)
				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(row)
			},
			checkResponse: func(t *testing.T, res []*domain.Tag, count int64, err error) {
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
					NewRows([]string{"tag", "count"}).
					AddRow(tag.Name, tag.Count)
				offset := (a.page - 1) * a.limit
				query := fmt.Sprintf(`
	SELECT *, count(1) FROM ( SELECT unnest(tags) AS tag FROM stores ) tags GROUP BY tag ORDER BY %s LIMIT %d OFFSET %d
	`, a.sort, a.limit, offset)
				queryCount := `SELECT count(1) FROM ( SELECT *, count(1) FROM( SELECT unnest(tags) AS tag FROM stores ) tags GROUP BY tag ) stores`
				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(row)
				mock.ExpectQuery(regexp.QuoteMeta(queryCount)).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, res []*domain.Tag, count int64, err error) {
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
					NewRows([]string{"tag", "count"}).
					AddRow(tag.Name, tag.Count)
				offset := (a.page - 1) * a.limit
				query := fmt.Sprintf(`
	SELECT *, count(1) FROM ( SELECT unnest(tags) AS tag FROM stores ) tags GROUP BY tag ORDER BY %s LIMIT %d OFFSET %d
	`, a.sort, a.limit, offset)
				countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)
				queryCount := `SELECT count(1) FROM ( SELECT *, count(1) FROM( SELECT unnest(tags) AS tag FROM stores ) tags GROUP BY tag ) stores`
				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(row)
				mock.ExpectQuery(regexp.QuoteMeta(queryCount)).WillReturnRows(countRow)
			},
			checkResponse: func(t *testing.T, res []*domain.Tag, count int64, err error) {
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
			repo := pg.NewTagRepository(db)
			tc.builtSts(mock)
			res, count, err := repo.GetAll(context.TODO(), tc.args.sort, tc.args.page, tc.args.limit)
			tc.checkResponse(t, res, count, err)
		})
	}
}

func Test_TagRepo_GetAllByCategory(t *testing.T) {
	type args struct {
		page       int
		limit      int
		sort       string
		categoryID string
	}
	a := args{
		page:       1,
		limit:      10,
		sort:       "created_at",
		categoryID: uuid.NewV4().String(),
	}
	tag := new(domain.Tag)
	tag.Name = "tag001"
	tag.Count = 2

	testCases := []struct {
		name          string
		args          args
		builtSts      func(mock sqlmock.Sqlmock)
		checkResponse func(t *testing.T, res []*domain.Tag, count int64, err error)
	}{
		{
			name: "fail on get all",
			args: a,
			builtSts: func(mock sqlmock.Sqlmock) {
				offset := (a.page - 1) * a.limit
				query := fmt.Sprintf(`SELECT *, count(1) FROM ( SELECT unnest(tags) AS tag FROM stores WHERE category_id = $1 ) stores GROUP BY tag ORDER BY %s LIMIT %d OFFSET %d
	`, a.sort, a.limit, offset)
				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, res []*domain.Tag, count int64, err error) {
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
					NewRows([]string{"tag", "count"}).
					AddRow(tag.Name, tag.Count)
				offset := (a.page - 1) * a.limit
				query := fmt.Sprintf(`SELECT *, count(1) FROM ( SELECT unnest(tags) AS tag FROM stores WHERE category_id = $1 ) stores GROUP BY tag ORDER BY %s LIMIT %d OFFSET %d
	`, a.sort, a.limit, offset)
				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(row)
				queryCount := `SELECT count(1) FROM ( SELECT *, count(1) FROM ( SELECT unnest(tags) AS tag FROM stores WHERE category_id = $1 ) tags GROUP BY tag ) stores`
				mock.ExpectQuery(regexp.QuoteMeta(queryCount)).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, res []*domain.Tag, count int64, err error) {
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
					NewRows([]string{"tag", "count"}).
					AddRow(tag.Name, tag.Count)
				offset := (a.page - 1) * a.limit
				query := fmt.Sprintf(`SELECT *, count(1) FROM ( SELECT unnest(tags) AS tag FROM stores WHERE category_id = $1 ) stores GROUP BY tag ORDER BY %s LIMIT %d OFFSET %d
	`, a.sort, a.limit, offset)
				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(row)
				countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)
				queryCount := `SELECT count(1) FROM ( SELECT *, count(1) FROM ( SELECT unnest(tags) AS tag FROM stores WHERE category_id = $1 ) tags GROUP BY tag ) stores`
				mock.ExpectQuery(regexp.QuoteMeta(queryCount)).WillReturnRows(countRow)
			},
			checkResponse: func(t *testing.T, res []*domain.Tag, count int64, err error) {
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
			repo := pg.NewTagRepository(db)
			tc.builtSts(mock)
			res, count, err := repo.GetAllByCategory(context.TODO(), tc.args.categoryID, tc.args.sort, tc.args.page, tc.args.limit)
			tc.checkResponse(t, res, count, err)
		})
	}
}
