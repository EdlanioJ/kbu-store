package gorm_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/repository/gorm"
	"github.com/EdlanioJ/kbu-store/app/utils/sample"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestStoreRepository(t *testing.T) {
	t.Parallel()
	db, mock := dbMock()
	repo := gorm.NewStoreRepository(db)

	t.Run("Create", func(t *testing.T) {
		store := sample.NewStore()
		query := `INSERT INTO "stores" ("id","created_at","updated_at","name","description","status","user_id","account_id","category_id","tags","lat","lng") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(store.ID, store.CreatedAt, sqlmock.AnyArg(), store.Name, store.Description, store.Status, store.UserID, store.AccountID, store.CategoryID, pq.StringArray(store.Tags), store.Position.Lat, store.Position.Lng).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Create(context.TODO(), store)
		assert.NoError(t, err)
	})
	t.Run("FindByID", func(t *testing.T) {
		store := sample.NewStore()
		query := `SELECT * FROM "stores" WHERE id = $1 ORDER BY "stores"."id"`

		row := sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "external_id", "lat", "lng"}).
			AddRow(store.ID, store.CreatedAt, store.UpdatedAt, store.Name, store.Status, store.Description, store.AccountID, store.CategoryID, store.Name, store.Position.Lat, store.Position.Lng)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(store.ID).
			WillReturnRows(row)

		res, err := repo.FindByID(context.TODO(), store.ID)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
	t.Run("FindByName", func(t *testing.T) {
		store := sample.NewStore()
		query := `SELECT * FROM "stores" WHERE name = $1 ORDER BY "stores"."id"`

		row := sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "external_id", "lat", "lng"}).
			AddRow(store.ID, store.CreatedAt, store.UpdatedAt, store.Name, store.Status, store.Description, store.AccountID, store.CategoryID, store.Name, store.Position.Lat, store.Position.Lng)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(store.Name).
			WillReturnRows(row)

		res, err := repo.FindByName(context.TODO(), store.Name)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
	t.Run("FindAll", func(t *testing.T) {
		store := sample.NewStore()
		page := 2
		limit := 10
		sort := "created_at DESC"
		query := fmt.Sprintf(`SELECT * FROM "stores" ORDER BY %s LIMIT %d`, sort, limit)
		queryCount := `SELECT count(*) FROM "stores"`

		countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)
		row := sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "external_id", "lat", "lng"}).
			AddRow(store.ID, store.CreatedAt, store.UpdatedAt, store.Name, store.Status, store.Description, store.AccountID, store.CategoryID, store.Name, store.Position.Lat, store.Position.Lng)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(row)
		mock.ExpectQuery(regexp.QuoteMeta(queryCount)).
			WillReturnRows(countRow)

		list, total, err := repo.FindAll(context.TODO(), sort, limit, page)
		assert.NoError(t, err)
		assert.Equal(t, total, int64(1))
		assert.Len(t, list, 1)
	})
	t.Run("Update", func(t *testing.T) {
		store := sample.NewStore()
		query := `UPDATE "stores" SET "created_at"=$1,"updated_at"=$2,"name"=$3,"description"=$4,"status"=$5,"user_id"=$6,"account_id"=$7,"category_id"=$8,"tags"=$9,"lat"=$10,"lng"=$11 WHERE "id" = $12`

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(store.CreatedAt, sqlmock.AnyArg(), store.Name, store.Description, store.Status, store.UserID, store.AccountID, store.CategoryID, pq.StringArray(store.Tags), store.Position.Lat, store.Position.Lng, store.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Update(context.TODO(), store)
		assert.NoError(t, err)
	})
	t.Run("Delete", func(t *testing.T) {
		store := sample.NewStore()
		query := `DELETE FROM "stores" WHERE id = $1`

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(store.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Delete(context.TODO(), store.ID)
		assert.NoError(t, err)
	})
}
