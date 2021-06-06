package gorm_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/account/repository/gorm"
	"github.com/stretchr/testify/assert"
)

func Test_GormCreateAccountRepository(t *testing.T) {
	is := assert.New(t)
	db, mock, accountMock := setupTest()

	repo := gorm.NewGormCreateAccount(db)

	query := `INSERT INTO "accounts" ("id","created_at","updated_at","balance") VALUES ($1,$2,$3,$4)`

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(accountMock.ID, accountMock.CreatedAt, sqlmock.AnyArg(), accountMock.Balance).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Add(context.TODO(), accountMock)

	is.NoError(err)
}
