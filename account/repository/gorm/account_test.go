package gorm_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/account/repository/gorm"
	"github.com/stretchr/testify/assert"
)

func Test_GormAccountRepository(t *testing.T) {
	t.Run("account repo -> create", func(t *testing.T) {
		is := assert.New(t)
		db, mock, accountMock := testMock()
		repo := gorm.NewAccountRepository(db)

		query := `INSERT INTO "accounts" ("id","created_at","updated_at","balance") VALUES ($1,$2,$3,$4)`

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(accountMock.ID, accountMock.CreatedAt, sqlmock.AnyArg(), accountMock.Balance).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Create(context.TODO(), accountMock)
		is.NoError(err)
	})
	t.Run("account repo -> get by id", func(t *testing.T) {
		is := assert.New(t)
		db, mock, accountMock := testMock()
		repo := gorm.NewAccountRepository(db)

		row := sqlmock.
			NewRows([]string{"id", "balance", "created_at", "updated_at"}).
			AddRow(accountMock.ID, accountMock.Balance, accountMock.CreatedAt, accountMock.UpdatedAt)

		query := `SELECT * FROM "accounts" WHERE id = $1 ORDER BY "accounts"."id" LIMIT 1`

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(accountMock.ID).
			WillReturnRows(row)

		res, err := repo.GetById(context.TODO(), accountMock.ID)

		is.NoError(err)
		is.NotNil(res)
		is.Equal(res, accountMock)
	})

	t.Run("account repo -> update", func(t *testing.T) {
		is := assert.New(t)
		db, mock, accountMock := testMock()
		repo := gorm.NewAccountRepository(db)

		query := `UPDATE "accounts" SET "created_at"=$1,"updated_at"=$2,"balance"=$3 WHERE "id" = $4`

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(accountMock.CreatedAt, sqlmock.AnyArg(), accountMock.Balance, accountMock.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Update(context.TODO(), accountMock)
		is.NoError(err)
	})

	t.Run("account repo -> delete", func(t *testing.T) {
		is := assert.New(t)
		db, mock, accountMock := testMock()
		repo := gorm.NewAccountRepository(db)

		query := `DELETE FROM "accounts" WHERE id = $1`

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(accountMock.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Delete(context.TODO(), accountMock.ID)
		is.NoError(err)
	})
}
