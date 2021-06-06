package gorm_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/category/repository/gorm"
	"github.com/stretchr/testify/assert"
)

func Test_GormGetCategoryByIDRepository(t *testing.T) {
	is := assert.New(t)
	db, mock, category := setupTest()

	repo := gorm.NewGormGetCategoryByID(db)

	query := `SELECT * FROM "categories" WHERE id = $1 ORDER BY "categories"."id" LIMIT 1`
	row := sqlmock.
		NewRows([]string{"id", "created_at", "updated_at", "name", "status"}).
		AddRow(category.ID, category.CreatedAt, category.UpdatedAt, category.Name, category.Status)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(category.ID).
		WillReturnRows(row)

	res, err := repo.Exec(context.TODO(), category.ID)
	is.NoError(err)
	is.NotNil(res)
}
