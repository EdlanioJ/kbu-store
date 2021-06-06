package gorm_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/store/repository/gorm"
	"github.com/stretchr/testify/assert"
)

func Test_GormFetchStoreByCloseLocationRepository(t *testing.T) {
	is := assert.New(t)
	db, mock, store := testMock()

	repo := gorm.NewGormFetchStoreByCloseLocation(db)
	page := 2
	limit := 10
	sort := "created_at DESC"
	distance := 1
	offset := (page - 1) * limit

	query := fmt.Sprintf(`SELECT * FROM (SELECT *,(((acos(sin((%[1]v * pi() / 180)) * sin((lat * pi() / 180)) + cos((%[1]v * pi() / 180)) * cos((lat * pi() / 180)) * cos(((%[2]v - lng) * pi() / 180)))) * 180 / pi()) * 60 * 1.1515 * 1.609344 ) AS distance FROM "stores") stores WHERE distance <= $1 AND status = $2 ORDER BY %[3]v LIMIT %[4]v OFFSET %[5]v`, store.Position.Lat, store.Position.Lng, sort, limit, offset)
	row := sqlmock.
		NewRows([]string{"id", "created_at", "updated_at", "name", "status", "description", "account_id", "category_id", "external_id", "lat", "lng"}).
		AddRow(store.ID, store.CreatedAt, store.UpdatedAt, store.Name, store.Status, store.Description, store.Account.ID, store.Category.ID, store.ExternalID, store.Position.Lat, store.Position.Lng)

	queryCount := `SELECT count(1) FROM (SELECT *,(((acos(sin((-8.8867698 * pi() / 180)) * sin((lat * pi() / 180)) + cos((-8.8867698 * pi() / 180)) * cos((lat * pi() / 180)) * cos(((13.4771186 - lng) * pi() / 180)))) * 180 / pi()) * 60 * 1.1515 * 1.609344 ) AS distance FROM "stores") stores WHERE distance <= $1 AND status = $2`
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

	list, total, err := repo.Exec(context.TODO(), store.Position.Lat, store.Position.Lng, distance, limit, page, store.Status, sort)

	is.NoError(err)
	is.Equal(total, int64(1))
	is.Len(list, 1)

}
