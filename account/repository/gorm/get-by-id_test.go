package gorm_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/account/repository/gorm"
	"github.com/stretchr/testify/assert"
)

func Test_GormGetAccountByIDRepository(t *testing.T) {
	is := assert.New(t)
	db, mock, accountMock := setupTest()

	repo := gorm.NewGormGetAccountByID(db)

	row := sqlmock.
		NewRows([]string{"id", "balance", "created_at", "updated_at"}).
		AddRow(accountMock.ID, accountMock.Balance, accountMock.CreatedAt, accountMock.UpdatedAt)

	query := `SELECT * FROM "accounts" WHERE id = $1 ORDER BY "accounts"."id" LIMIT 1`

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(accountMock.ID).
		WillReturnRows(row)

	res, err := repo.Exec(context.TODO(), accountMock.ID)

	is.NoError(err)
	is.NotNil(res)

	is.Equal(res, accountMock)
}
