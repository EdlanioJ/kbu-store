package gorm_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/category/repository/gorm"
	"github.com/stretchr/testify/assert"
)

func Test_GormFetchCategoryByStatusRepository(t *testing.T) {
	is := assert.New(t)
	db, mock, Category := setupTest()

	repo := gorm.NewGormFetchCategoryByStatus(db)

	page := 2
	limit := 10
	sort := "created_at DESC"
	query := fmt.Sprintf(`SELECT * FROM "categories" WHERE status = $1 ORDER BY %s LIMIT %d`, sort, limit)
	queryCount := `SELECT count(1) FROM "categories"`

	countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)
	row := sqlmock.
		NewRows([]string{"id", "created_at", "updated_at", "name", "status"}).
		AddRow(Category.ID, Category.CreatedAt, Category.UpdatedAt, Category.Name, Category.Status)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(Category.Status).
		WillReturnRows(row)
	mock.ExpectQuery(regexp.QuoteMeta(queryCount)).
		WillReturnRows(countRow)

	list, total, err := repo.Exec(context.TODO(), Category.Status, sort, page, limit)

	is.NoError(err)
	is.Equal(total, int64(1))
	is.Len(list, 1)

}
