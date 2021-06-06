package gorm_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/category/repository/gorm"
	"github.com/stretchr/testify/assert"
)

func Test_GormCreateCategoryRepository(t *testing.T) {
	is := assert.New(t)
	db, mock, category := setupTest()

	repo := gorm.NewGormCreateCategory(db)

	query := `INSERT INTO "categories" ("id","created_at","updated_at","name","status") VALUES ($1,$2,$3,$4,$5)`

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(category.ID, category.CreatedAt, sqlmock.AnyArg(), category.Name, category.Status).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Add(context.TODO(), category)

	is.NoError(err)
}
