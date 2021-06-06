package gorm_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/store/repository/gorm"
	"github.com/stretchr/testify/assert"
)

func Test_GormGetStoreByOwnerRepository(t *testing.T) {
	is := assert.New(t)
	db, mock, store := testMock()

	repo := gorm.NewGormGetStoreByOwner(db)

	query := `SELECT * FROM "stores" WHERE id = $1 AND external_id = $2 ORDER BY "stores"."id" LIMIT 1`

	row := sqlmock.
		NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "external_id", "lat", "lng"}).
		AddRow(store.ID, store.CreatedAt, store.UpdatedAt, store.Name, store.Status, store.Description, store.Account.ID, store.Category.ID, store.ExternalID, store.Position.Lat, store.Position.Lng)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(store.ID, store.ExternalID).
		WillReturnRows(row)

	res, err := repo.Exec(context.TODO(), store.ID, store.ExternalID)
	is.NoError(err)
	is.NotNil(res)
}
