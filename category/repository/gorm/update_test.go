package gorm_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/category/repository/gorm"
	"github.com/stretchr/testify/assert"
)

func Test_GormUpdateCategoryRepository(t *testing.T) {
	is := assert.New(t)
	db, mock, category := setupTest()

	repo := gorm.NewGormUpdateCategory(db)

	query := `UPDATE "categories" SET "created_at"=$1,"updated_at"=$2,"name"=$3,"status"=$4 WHERE "id" = $5`

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(category.CreatedAt, sqlmock.AnyArg(), category.Name, category.Status, category.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Exec(context.TODO(), category)

	is.NoError(err)
}
