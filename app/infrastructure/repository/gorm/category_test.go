package gorm_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/repository/gorm"
	"github.com/EdlanioJ/kbu-store/app/utils/sample"
	"github.com/stretchr/testify/assert"
)

func TestCategoryRepository(t *testing.T) {
	t.Parallel()
	db, mock := dbMock()
	repo := gorm.NewCategoryRepository(db)

	t.Run("Create", func(t *testing.T) {
		category := sample.NewCategory()
		query := `INSERT INTO "categories" ("id","created_at","updated_at","name","status") VALUES ($1,$2,$3,$4,$5)`

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(category.ID, category.CreatedAt, sqlmock.AnyArg(), category.Name, category.Status).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Store(context.TODO(), category)
		assert.NoError(t, err)
	})

	t.Run("FindByID", func(t *testing.T) {
		category := sample.NewCategory()
		query := `SELECT * FROM "categories" WHERE id = $1 ORDER BY "categories"."id" LIMIT 1`
		row := sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "status"}).
			AddRow(category.ID, category.CreatedAt, category.UpdatedAt, category.Name, category.Status)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(category.ID).
			WillReturnRows(row)

		res, err := repo.FindByID(context.TODO(), category.ID)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
	t.Run("Update", func(t *testing.T) {
		category := sample.NewCategory()
		query := `UPDATE "categories" SET "created_at"=$1,"updated_at"=$2,"name"=$3,"status"=$4 WHERE "id" = $5`

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(category.CreatedAt, sqlmock.AnyArg(), category.Name, category.Status, category.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Update(context.TODO(), category)
		assert.NoError(t, err)
	})
}
