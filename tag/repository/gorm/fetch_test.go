package gorm_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/tag/repository/gorm"
	"github.com/stretchr/testify/assert"
)

func Test_GormFetchTag(t *testing.T) {
	is := assert.New(t)
	db, mock, tag := testMock()

	repo := gorm.NewGormFetchTags(db)

	sort := "total DESC"
	page := 1
	limit := 10
	offset := (page - 1) * limit
	query := fmt.Sprintf(`
	SELECT *, count(1) FROM ( SELECT unnest(tags) AS tag FROM stores ) tags GROUP BY tag ORDER BY %s LIMIT %d OFFSET %d
	`, sort, limit, offset)
	queryCount := `SELECT count(1) FROM ( SELECT *, count(1) FROM( SELECT unnest(tags) AS tag FROM stores ) tags GROUP BY tag ) stores`

	countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)
	row := sqlmock.
		NewRows([]string{"tag", "count"}).
		AddRow(tag.Name, tag.Count)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WillReturnRows(row)

	mock.ExpectQuery(regexp.QuoteMeta(queryCount)).
		WillReturnRows(countRow)

	res, total, err := repo.Exec(context.TODO(), sort, page, limit)

	is.NoError(err)
	is.Equal(total, int64(1))
	is.Len(res, 1)
	is.Equal(res[0], tag)
}
