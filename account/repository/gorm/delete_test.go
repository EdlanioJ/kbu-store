package gorm_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/account/repository/gorm"
	"github.com/stretchr/testify/assert"
)

func Test_GormDeleteAccountRepository(t *testing.T) {
	is := assert.New(t)

	db, mock, accountMock := setupTest()

	repo := gorm.NewGormDeleteAccount(db)

	query := `DELETE FROM "accounts" WHERE id = $1`

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(accountMock.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Exec(context.TODO(), accountMock.ID)

	is.NoError(err)
}
