package gorm_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/repository/gorm"
	"github.com/EdlanioJ/kbu-store/app/utils/sample"
	"github.com/stretchr/testify/assert"
)

func TestAccountRepository(t *testing.T) {
	t.Parallel()
	db, mock := dbMock()
	repo := gorm.NewAccountRepository(db)

	t.Run("Store", func(t *testing.T) {
		account := sample.NewAccount()
		query := `INSERT INTO "accounts" ("id","created_at","updated_at","balance") VALUES ($1,$2,$3,$4)`
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(account.ID, account.CreatedAt, sqlmock.AnyArg(), account.Balance).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Store(context.TODO(), account)
		assert.NoError(t, err)
	})

	t.Run("FindByID", func(t *testing.T) {
		account := sample.NewAccount()
		row := sqlmock.
			NewRows([]string{"id", "balance", "created_at", "updated_at"}).
			AddRow(account.ID, account.Balance, account.CreatedAt, account.UpdatedAt)
		query := `SELECT * FROM "accounts" WHERE id = $1 ORDER BY "accounts"."id" LIMIT 1`

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(account.ID).
			WillReturnRows(row)

		res, err := repo.FindByID(context.TODO(), account.ID)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, res.ID, account.ID)
	})

	t.Run("Update", func(t *testing.T) {
		account := sample.NewAccount()
		query := `UPDATE "accounts" SET "created_at"=$1,"updated_at"=$2,"balance"=$3 WHERE "id" = $4`

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(account.CreatedAt, sqlmock.AnyArg(), account.Balance, account.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Update(context.TODO(), account)
		assert.NoError(t, err)
	})
	t.Run("Delete", func(t *testing.T) {
		account := sample.NewAccount()
		query := `DELETE FROM "accounts" WHERE id = $1`

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(account.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Delete(context.TODO(), account.ID)
		assert.NoError(t, err)
	})
}
