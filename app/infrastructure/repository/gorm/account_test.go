package gorm_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/repository/gorm"
	"github.com/stretchr/testify/assert"
)

func Test_GormAccountRepository(t *testing.T) {
	t.Run("account repo -> create", func(t *testing.T) {
		is := assert.New(t)
		db, mock := dbMock()
		account := getAccount()
		repo := gorm.NewAccountRepository(db)

		query := `INSERT INTO "accounts" ("id","created_at","updated_at","balance") VALUES ($1,$2,$3,$4)`

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(account.ID, account.CreatedAt, sqlmock.AnyArg(), account.Balance).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Store(context.TODO(), account)
		is.NoError(err)
	})
	t.Run("account repo -> get by id", func(t *testing.T) {
		is := assert.New(t)
		db, mock := dbMock()
		account := getAccount()
		repo := gorm.NewAccountRepository(db)

		row := sqlmock.
			NewRows([]string{"id", "balance", "created_at", "updated_at"}).
			AddRow(account.ID, account.Balance, account.CreatedAt, account.UpdatedAt)

		query := `SELECT * FROM "accounts" WHERE id = $1 ORDER BY "accounts"."id" LIMIT 1`

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(account.ID).
			WillReturnRows(row)

		res, err := repo.FindByID(context.TODO(), account.ID)

		is.NoError(err)
		is.NotNil(res)
		is.Equal(res, account)
	})

	t.Run("account repo -> update", func(t *testing.T) {
		is := assert.New(t)
		db, mock := dbMock()
		account := getAccount()
		repo := gorm.NewAccountRepository(db)

		query := `UPDATE "accounts" SET "created_at"=$1,"updated_at"=$2,"balance"=$3 WHERE "id" = $4`

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(account.CreatedAt, sqlmock.AnyArg(), account.Balance, account.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Update(context.TODO(), account)
		is.NoError(err)
	})
}
