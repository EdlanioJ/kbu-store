package gorm_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/store/repository/gorm"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func Test_GormFetchStoreByTagsRepository(t *testing.T) {
	is := assert.New(t)
	db, mock, store := testMock()

	repo := gorm.NewGormFetchStoreByTags(db)

	page := 2
	limit := 10
	sort := "created_at DESC"
	query := fmt.Sprintf(`SELECT * FROM "stores" WHERE tags && $1 ORDER BY %s LIMIT %d`, sort, limit)
	queryCount := `SELECT count(1) FROM "stores"`

	countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)
	row := sqlmock.
		NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "external_id", "lat", "lng"}).
		AddRow(store.ID, store.CreatedAt, store.UpdatedAt, store.Name, store.Status, store.Description, store.Account.ID, store.Category.ID, store.ExternalID, store.Position.Lat, store.Position.Lng)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(pq.StringArray(store.Tags)).
		WillReturnRows(row)
	mock.ExpectQuery(regexp.QuoteMeta(queryCount)).
		WillReturnRows(countRow)
	list, total, err := repo.Exec(context.TODO(), store.Tags, sort, limit, page)

	is.NoError(err)
	is.Equal(total, int64(1))
	is.Len(list, 1)
}
