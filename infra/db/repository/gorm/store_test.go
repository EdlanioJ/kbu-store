package gorm_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/infra/db/repository/gorm"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func Test_GormStoreRepository(t *testing.T) {
	t.Run("store repo -> create", func(t *testing.T) {
		is := assert.New(t)
		db, mock := dbMock()
		store := getStore()
		repo := gorm.NewStoreRepository(db)

		query := `INSERT INTO "stores" ("id","created_at","updated_at","name","description","status","external_id","account_id","category_id","tags","lat","lng") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(store.ID, store.CreatedAt, sqlmock.AnyArg(), store.Name, store.Description, store.Status, store.ExternalID, store.AccountID, store.Category.ID, pq.StringArray(store.Tags), store.Position.Lat, store.Position.Lng).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Create(context.TODO(), store)
		is.NoError(err)
	})

	t.Run("store repo -> get by id", func(t *testing.T) {
		is := assert.New(t)
		db, mock := dbMock()
		store := getStore()
		repo := gorm.NewStoreRepository(db)

		query := `SELECT * FROM "stores" WHERE id = $1 ORDER BY "stores"."id"`

		row := sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "external_id", "lat", "lng"}).
			AddRow(store.ID, store.CreatedAt, store.UpdatedAt, store.Name, store.Status, store.Description, store.AccountID, store.Category.ID, store.ExternalID, store.Position.Lat, store.Position.Lng)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(store.ID).
			WillReturnRows(row)

		res, err := repo.GetById(context.TODO(), store.ID)
		is.NoError(err)
		is.NotNil(res)
	})

	t.Run("store repo -> get by id and owner", func(t *testing.T) {
		is := assert.New(t)
		db, mock := dbMock()
		store := getStore()
		repo := gorm.NewStoreRepository(db)

		query := `SELECT * FROM "stores" WHERE id = $1 AND external_id = $2 ORDER BY "stores"."id" LIMIT 1`

		row := sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "external_id", "lat", "lng"}).
			AddRow(store.ID, store.CreatedAt, store.UpdatedAt, store.Name, store.Status, store.Description, store.AccountID, store.Category.ID, store.ExternalID, store.Position.Lat, store.Position.Lng)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(store.ID, store.ExternalID).
			WillReturnRows(row)

		res, err := repo.GetByIdAndOwner(context.TODO(), store.ID, store.ExternalID)
		is.NoError(err)
		is.NotNil(res)
	})

	t.Run("store repo -> get all", func(t *testing.T) {
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
			AddRow(store.ID, store.CreatedAt, store.UpdatedAt, store.Name, store.Status, store.Description, store.AccountID, store.Category.ID, store.ExternalID, store.Position.Lat, store.Position.Lng)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(row)
		mock.ExpectQuery(regexp.QuoteMeta(queryCount)).
			WillReturnRows(countRow)
		list, total, err := repo.GetAll(context.TODO(), sort, limit, page)

		is.NoError(err)
		is.Equal(total, int64(1))
		is.Len(list, 1)
	})

	t.Run("store repo -> get all by category", func(t *testing.T) {
		is := assert.New(t)
		db, mock := dbMock()
		store := getStore()
		repo := gorm.NewStoreRepository(db)

		page := 2
		limit := 10
		sort := "created_at DESC"
		query := fmt.Sprintf(`SELECT * FROM "stores" WHERE category_id = $1 ORDER BY %s LIMIT %d`, sort, limit)
		queryCount := `SELECT count(*) FROM "stores"`

		countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)
		row := sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "external_id", "lat", "lng"}).
			AddRow(store.ID, store.CreatedAt, store.UpdatedAt, store.Name, store.Status, store.Description, store.AccountID, store.Category.ID, store.ExternalID, store.Position.Lat, store.Position.Lng)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(store.Category.ID).
			WillReturnRows(row)
		mock.ExpectQuery(regexp.QuoteMeta(queryCount)).
			WillReturnRows(countRow)
		list, total, err := repo.GetAllByCategory(context.TODO(), store.Category.ID, sort, limit, page)

		is.NoError(err)
		is.Equal(total, int64(1))
		is.Len(list, 1)
	})

	t.Run("store repo -> get all by location", func(t *testing.T) {
		is := assert.New(t)
		db, mock := dbMock()
		store := getStore()
		repo := gorm.NewStoreRepository(db)

		page := 2
		limit := 10
		sort := "created_at DESC"
		distance := 1
		offset := (page - 1) * limit

		query := fmt.Sprintf(`SELECT * FROM (SELECT *,(((acos(sin((%[1]v * pi() / 180)) * sin((lat * pi() / 180)) + cos((%[1]v * pi() / 180)) * cos((lat * pi() / 180)) * cos(((%[2]v - lng) * pi() / 180)))) * 180 / pi()) * 60 * 1.1515 * 1.609344 ) AS distance FROM "stores") stores WHERE distance <= $1 AND status = $2 ORDER BY %[3]v LIMIT %[4]v OFFSET %[5]v`, store.Position.Lat, store.Position.Lng, sort, limit, offset)
		row := sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "external_id", "lat", "lng"}).
			AddRow(store.ID, store.CreatedAt, store.UpdatedAt, store.Name, store.Status, store.Description, store.AccountID, store.Category.ID, store.ExternalID, store.Position.Lat, store.Position.Lng)

		queryCount := `SELECT count(*) FROM (SELECT *,(((acos(sin((-8.8867698 * pi() / 180)) * sin((lat * pi() / 180)) + cos((-8.8867698 * pi() / 180)) * cos((lat * pi() / 180)) * cos(((13.4771186 - lng) * pi() / 180)))) * 180 / pi()) * 60 * 1.1515 * 1.609344 ) AS distance FROM "stores") stores WHERE distance <= $1 AND status = $2`
		countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(
				distance,
				store.Status,
			).
			WillReturnRows(row)
		mock.ExpectQuery(regexp.QuoteMeta(queryCount)).
			WithArgs(
				distance,
				store.Status,
			).
			WillReturnRows(countRow)

		list, total, err := repo.GetAllByLocation(context.TODO(), store.Position.Lat, store.Position.Lng, distance, limit, page, store.Status, sort)

		is.NoError(err)
		is.Equal(total, int64(1))
		is.Len(list, 1)
	})

	t.Run("store repo -> get all by owner", func(t *testing.T) {
		is := assert.New(t)
		db, mock := dbMock()
		store := getStore()
		repo := gorm.NewStoreRepository(db)

		page := 2
		limit := 10
		sort := "created_at DESC"
		query := fmt.Sprintf(`SELECT * FROM "stores" WHERE external_id = $1 ORDER BY %s LIMIT %d`, sort, limit)
		queryCount := `SELECT count(*) FROM "stores"`

		countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)
		row := sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "external_id", "lat", "lng"}).
			AddRow(store.ID, store.CreatedAt, store.UpdatedAt, store.Name, store.Status, store.Description, store.AccountID, store.Category.ID, store.ExternalID, store.Position.Lat, store.Position.Lng)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(store.ExternalID).
			WillReturnRows(row)
		mock.ExpectQuery(regexp.QuoteMeta(queryCount)).
			WillReturnRows(countRow)
		list, total, err := repo.GetAllByOwner(context.TODO(), store.ExternalID, sort, limit, page)

		is.NoError(err)
		is.Equal(total, int64(1))
		is.Len(list, 1)
	})

	t.Run("store repo -> get all by status", func(t *testing.T) {
		is := assert.New(t)
		db, mock := dbMock()
		store := getStore()
		repo := gorm.NewStoreRepository(db)

		page := 2
		limit := 10
		sort := "created_at DESC"
		query := fmt.Sprintf(`SELECT * FROM "stores" WHERE status = $1 ORDER BY %s LIMIT %d`, sort, limit)
		queryCount := `SELECT count(*) FROM "stores"`

		countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)
		row := sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "external_id", "lat", "lng"}).
			AddRow(store.ID, store.CreatedAt, store.UpdatedAt, store.Name, store.Status, store.Description, store.AccountID, store.Category.ID, store.ExternalID, store.Position.Lat, store.Position.Lng)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(store.Status).
			WillReturnRows(row)
		mock.ExpectQuery(regexp.QuoteMeta(queryCount)).
			WillReturnRows(countRow)
		list, total, err := repo.GetAllByStatus(context.TODO(), store.Status, sort, limit, page)

		is.NoError(err)
		is.Equal(total, int64(1))
		is.Len(list, 1)
	})

	t.Run("store repo -> get all by tags", func(t *testing.T) {
		is := assert.New(t)
		db, mock := dbMock()
		store := getStore()
		repo := gorm.NewStoreRepository(db)

		page := 2
		limit := 10
		sort := "created_at DESC"
		query := fmt.Sprintf(`SELECT * FROM "stores" WHERE tags && $1 ORDER BY %s LIMIT %d`, sort, limit)
		queryCount := `SELECT count(*) FROM "stores"`

		countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)
		row := sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "external_id", "lat", "lng"}).
			AddRow(store.ID, store.CreatedAt, store.UpdatedAt, store.Name, store.Status, store.Description, store.AccountID, store.Category.ID, store.ExternalID, store.Position.Lat, store.Position.Lng)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(pq.StringArray(store.Tags)).
			WillReturnRows(row)
		mock.ExpectQuery(regexp.QuoteMeta(queryCount)).
			WillReturnRows(countRow)
		list, total, err := repo.GetAllByTags(context.TODO(), store.Tags, sort, limit, page)

		is.NoError(err)
		is.Equal(total, int64(1))
		is.Len(list, 1)
	})

	t.Run("store repo -> update", func(t *testing.T) {
		is := assert.New(t)
		db, mock := dbMock()
		store := getStore()
		repo := gorm.NewStoreRepository(db)

		query := `UPDATE "stores" SET "created_at"=$1,"updated_at"=$2,"name"=$3,"description"=$4,"status"=$5,"external_id"=$6,"account_id"=$7,"category_id"=$8,"tags"=$9,"lat"=$10,"lng"=$11 WHERE "id" = $12`

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(store.CreatedAt, sqlmock.AnyArg(), store.Name, store.Description, store.Status, store.ExternalID, store.AccountID, store.Category.ID, pq.StringArray(store.Tags), store.Position.Lat, store.Position.Lng, store.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Update(context.TODO(), store)

		is.NoError(err)
	})

	t.Run("store repo -> delete", func(t *testing.T) {
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
