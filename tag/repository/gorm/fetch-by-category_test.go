package gorm_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/tag/repository/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func Test_GormFetchTagsByCategory(t *testing.T) {
	is := assert.New(t)
	db, mock, tag := testMock()

	repo := gorm.NewGormFetchTagsByCategory(db)

	sort := "total DESC"
	page := 1
	limit := 10
	offset := (page - 1) * limit
	categoryID := uuid.NewV4().String()

	query := fmt.Sprintf(`SELECT *, count(1) FROM ( SELECT unnest(tags) AS tag FROM stores WHERE category_id = $1 ) stores GROUP BY tag ORDER BY %s LIMIT %d OFFSET %d
	`, sort, limit, offset)

	queryCount := `SELECT count(1) FROM ( SELECT *, count(1) FROM ( SELECT unnest(tags) AS tag FROM stores WHERE category_id = $1 ) tags GROUP BY tag ) stores`
	countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)
	row := sqlmock.
		NewRows([]string{"tag", "count"}).
		AddRow(tag.Name, tag.Count)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(categoryID).
		WillReturnRows(row)
	mock.ExpectQuery(regexp.QuoteMeta(queryCount)).
		WithArgs(categoryID).
		WillReturnRows(countRow)

	res, total, err := repo.Exec(context.TODO(), categoryID, sort, page, limit)

	is.NoError(err)
	is.Equal(total, int64(1))
	is.Len(res, 1)
	is.Equal(res[0], tag)
}
