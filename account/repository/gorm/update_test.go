package gorm_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/account/repository/gorm"
	"github.com/stretchr/testify/assert"
)

func Test_GormUpdateAccountRepository(t *testing.T) {
	is := assert.New(t)
	db, mock, accountMock := setupTest()

	repo := gorm.NewGormUpdateAccount(db)
	query := `UPDATE "accounts" SET "created_at"=$1,"updated_at"=$2,"balance"=$3 WHERE "id" = $4`

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(accountMock.CreatedAt, sqlmock.AnyArg(), accountMock.Balance, accountMock.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Exec(context.TODO(), accountMock)

	is.NoError(err)
}
