package gorm_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/store/repository/gorm"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func Test_GormCreateStoreRepository(t *testing.T) {
	is := assert.New(t)
	db, mock, store := testMock()

	repo := gorm.NewGormCreateStore(db)
	query := `INSERT INTO "stores" ("id","created_at","updated_at","name","description","status","external_id","account_id","category_id","tags","lat","lng") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(store.ID, store.CreatedAt, sqlmock.AnyArg(), store.Name, store.Description, store.Status, store.ExternalID, store.Account.ID, store.Category.ID, pq.StringArray(store.Tags), store.Position.Lat, store.Position.Lng).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Add(context.TODO(), store)
	is.NoError(err)
}
