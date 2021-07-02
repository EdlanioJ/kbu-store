package gorm_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/infra/db/repository/gorm"
	"github.com/stretchr/testify/assert"
)

func Test_GormCategoryRepository(t *testing.T) {
	t.Run("category repo  -> create", func(t *testing.T) {
		is := assert.New(t)
		db, mock := dbMock()
		category := getCategory()

		repo := gorm.NewCategoryRepository(db)

		query := `INSERT INTO "categories" ("id","created_at","updated_at","name","status") VALUES ($1,$2,$3,$4,$5)`

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(category.ID, category.CreatedAt, sqlmock.AnyArg(), category.Name, category.Status).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Create(context.TODO(), category)

		is.NoError(err)
	})

	t.Run("category repo -> get all", func(t *testing.T) {
		is := assert.New(t)
		db, mock := dbMock()
		category := getCategory()

		repo := gorm.NewCategoryRepository(db)

		page := 2
		limit := 10
		sort := "created_at DESC"
		query := fmt.Sprintf(`SELECT * FROM "categories" ORDER BY %s LIMIT %d`, sort, limit)
		queryCount := `SELECT count(1) FROM "categories"`

		countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)
		row := sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "status"}).
			AddRow(category.ID, category.CreatedAt, category.UpdatedAt, category.Name, category.Status)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(row)
		mock.ExpectQuery(regexp.QuoteMeta(queryCount)).
			WillReturnRows(countRow)

		list, total, err := repo.GetAll(context.TODO(), sort, page, limit)

		is.NoError(err)
		is.Equal(total, int64(1))
		is.Len(list, 1)
	})

	t.Run("category repo -> get all status", func(t *testing.T) {
		is := assert.New(t)
		db, mock := dbMock()
		category := getCategory()
		repo := gorm.NewCategoryRepository(db)

		page := 2
		limit := 10
		sort := "created_at DESC"
		query := fmt.Sprintf(`SELECT * FROM "categories" WHERE status = $1 ORDER BY %s LIMIT %d`, sort, limit)
		queryCount := `SELECT count(1) FROM "categories"`

		countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)
		row := sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "status"}).
			AddRow(category.ID, category.CreatedAt, category.UpdatedAt, category.Name, category.Status)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(category.Status).
			WillReturnRows(row)
		mock.ExpectQuery(regexp.QuoteMeta(queryCount)).
			WillReturnRows(countRow)

		list, total, err := repo.GetAllByStatus(context.TODO(), category.Status, sort, page, limit)

		is.NoError(err)
		is.Equal(total, int64(1))
		is.Len(list, 1)
	})

	t.Run("category repo -> get by id", func(t *testing.T) {
		is := assert.New(t)
		db, mock := dbMock()
		category := getCategory()
		repo := gorm.NewCategoryRepository(db)

		query := `SELECT * FROM "categories" WHERE id = $1 ORDER BY "categories"."id" LIMIT 1`
		row := sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "status"}).
			AddRow(category.ID, category.CreatedAt, category.UpdatedAt, category.Name, category.Status)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(category.ID).
			WillReturnRows(row)

		res, err := repo.GetById(context.TODO(), category.ID)
		is.NoError(err)
		is.NotNil(res)
	})

	t.Run("category repo -> get by id and status", func(t *testing.T) {
		is := assert.New(t)
		db, mock := dbMock()
		category := getCategory()
		repo := gorm.NewCategoryRepository(db)

		query := `SELECT * FROM "categories" WHERE id = $1 AND status = $2 ORDER BY "categories"."id" LIMIT 1`
		row := sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "status"}).
			AddRow(category.ID, category.CreatedAt, category.UpdatedAt, category.Name, category.Status)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(category.ID, category.Status).
			WillReturnRows(row)

		res, err := repo.GetByIdAndStatus(context.TODO(), category.ID, category.Status)
		is.NoError(err)
		is.NotNil(res)
	})

	t.Run("category repo -> update", func(t *testing.T) {
		is := assert.New(t)
		db, mock := dbMock()
		category := getCategory()
		repo := gorm.NewCategoryRepository(db)

		query := `UPDATE "categories" SET "created_at"=$1,"updated_at"=$2,"name"=$3,"status"=$4 WHERE "id" = $5`

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(category.CreatedAt, sqlmock.AnyArg(), category.Name, category.Status, category.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Update(context.TODO(), category)

		is.NoError(err)
	})
}
