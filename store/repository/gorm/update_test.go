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

func Test_GormUpdateStoreRepository(t *testing.T) {
	is := assert.New(t)
	db, mock, store := testMock()

	repo := gorm.NewGormUpdateStore(db)
	query := `UPDATE "stores" SET "created_at"=$1,"updated_at"=$2,"name"=$3,"description"=$4,"status"=$5,"external_id"=$6,"account_id"=$7,"category_id"=$8,"tags"=$9,"lat"=$10,"lng"=$11 WHERE "id" = $12`

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(store.CreatedAt, sqlmock.AnyArg(), store.Name, store.Description, store.Status, store.ExternalID, store.Account.ID, store.Category.ID, pq.StringArray(store.Tags), store.Position.Lat, store.Position.Lng, store.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Exec(context.TODO(), store)

	is.NoError(err)
}
