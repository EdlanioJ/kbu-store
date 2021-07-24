package gorm_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/repository/gorm"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func Test_GormStoreRepository(t *testing.T) {
	t.Run("should test store repo on create", func(t *testing.T) {
		is := assert.New(t)
		db, mock := dbMock()
		store := getStore()
		repo := gorm.NewStoreRepository(db)

		query := `INSERT INTO "stores" ("id","created_at","updated_at","name","description","status","user_id","account_id","category_id","tags","lat","lng") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(store.ID, store.CreatedAt, sqlmock.AnyArg(), store.Name, store.Description, store.Status, store.UserID, store.AccountID, store.CategoryID, pq.StringArray(store.Tags), store.Position.Lat, store.Position.Lng).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Create(context.TODO(), store)
		is.NoError(err)
	})

	t.Run("should test store repo on get by id", func(t *testing.T) {
		is := assert.New(t)
		db, mock := dbMock()
		store := getStore()
		repo := gorm.NewStoreRepository(db)

		query := `SELECT * FROM "stores" WHERE id = $1 ORDER BY "stores"."id"`

		row := sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "external_id", "lat", "lng"}).
			AddRow(store.ID, store.CreatedAt, store.UpdatedAt, store.Name, store.Status, store.Description, store.AccountID, store.CategoryID, store.Name, store.Position.Lat, store.Position.Lng)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(store.ID).
			WillReturnRows(row)

		res, err := repo.FindByID(context.TODO(), store.ID)
		is.NoError(err)
		is.NotNil(res)
	})

	t.Run("should test store repo on get by name", func(t *testing.T) {
		is := assert.New(t)
		db, mock := dbMock()
		store := getStore()
		repo := gorm.NewStoreRepository(db)

		query := `SELECT * FROM "stores" WHERE name = $1 ORDER BY "stores"."id"`

		row := sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "external_id", "lat", "lng"}).
			AddRow(store.ID, store.CreatedAt, store.UpdatedAt, store.Name, store.Status, store.Description, store.AccountID, store.CategoryID, store.Name, store.Position.Lat, store.Position.Lng)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(store.Name).
			WillReturnRows(row)

		res, err := repo.FindByName(context.TODO(), store.Name)
		is.NoError(err)
		is.NotNil(res)
	})
	t.Run("should test store repo on find all", func(t *testing.T) {
		is := assert.New(t)
		db, mock := dbMock()
		store := getStore()
		repo := gorm.NewStoreRepository(db)

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

		is.NoError(err)
		is.Equal(total, int64(1))
		is.Len(list, 1)
	})

	t.Run("should test store repo on update", func(t *testing.T) {
		is := assert.New(t)
		db, mock := dbMock()
		store := getStore()
		repo := gorm.NewStoreRepository(db)

		query := `UPDATE "stores" SET "created_at"=$1,"updated_at"=$2,"name"=$3,"description"=$4,"status"=$5,"user_id"=$6,"account_id"=$7,"category_id"=$8,"tags"=$9,"lat"=$10,"lng"=$11 WHERE "id" = $12`

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(store.CreatedAt, sqlmock.AnyArg(), store.Name, store.Description, store.Status, store.UserID, store.AccountID, store.CategoryID, pq.StringArray(store.Tags), store.Position.Lat, store.Position.Lng, store.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Update(context.TODO(), store)

		is.NoError(err)
	})

	t.Run("should test store repo on delete", func(t *testing.T) {
		is := assert.New(t)
		db, mock := dbMock()
		store := getStore()
		repo := gorm.NewStoreRepository(db)

		query := `DELETE FROM "stores" WHERE id = $1`

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(store.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Delete(context.TODO(), store.ID)

		is.NoError(err)
	})
}
