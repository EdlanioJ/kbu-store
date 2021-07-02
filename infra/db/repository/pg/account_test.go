package pg_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/infra/db/repository/pg"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func getAccount() *domain.Account {
	account, _ := domain.NewAccount(20.75)
	return account
}

func Test_AccountRepo_Create(t *testing.T) {
	a := getAccount()
	testCases := []struct {
		name          string
		arg           *domain.Account
		builtSts      func(mock sqlmock.Sqlmock)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "fail on exec query",
			arg:  a,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO accounts (id,created_at,updated_at,balance) VALUES ($1,$2,$3,$4)`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(a.ID, a.CreatedAt, a.UpdatedAt, a.Balance).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail on get affected row",
			arg:  a,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO accounts (id,created_at,updated_at,balance) VALUES ($1,$2,$3,$4)`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(a.ID, a.CreatedAt, a.UpdatedAt, a.Balance).WillReturnResult(sqlmock.NewErrorResult(errors.New("unexpected error")))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail on invalid number of affected row",
			arg:  a,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO accounts (id,created_at,updated_at,balance) VALUES ($1,$2,$3,$4)`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(a.ID, a.CreatedAt, a.UpdatedAt, a.Balance).WillReturnResult(sqlmock.NewResult(1, 2))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			arg:  a,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO accounts (id,created_at,updated_at,balance) VALUES ($1,$2,$3,$4)`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(a.ID, a.CreatedAt, a.UpdatedAt, a.Balance).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			repo := pg.NewAccountRepository(db)
			tc.builtSts(mock)
			err = repo.Create(context.TODO(), tc.arg)
			tc.checkResponse(t, err)
		})
	}
}

func Test_AccountRepo_GetById(t *testing.T) {
	id := uuid.NewV4().String()
	a := getAccount()
	testCases := []struct {
		name          string
		arg           string
		builtSts      func(mock sqlmock.Sqlmock)
		checkResponse func(t *testing.T, res *domain.Account, err error)
	}{
		{
			name: "fail on exec query",
			arg:  id,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `SELECT * FROM accounts WHERE id = $1`
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(id).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, res *domain.Account, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			arg:  id,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `SELECT * FROM accounts WHERE id = $1`
				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "balance"}).
					AddRow(a.ID, a.CreatedAt, a.UpdatedAt, a.Balance)

				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(id).WillReturnRows(row)
			},
			checkResponse: func(t *testing.T, res *domain.Account, err error) {
				fmt.Println(a)
				assert.NotNil(t, res)
				assert.NoError(t, err)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			repo := pg.NewAccountRepository(db)
			tc.builtSts(mock)
			res, err := repo.GetById(context.TODO(), tc.arg)
			tc.checkResponse(t, res, err)
		})
	}
}

func Test_AccountRepo_Update(t *testing.T) {
	a := getAccount()
	testCases := []struct {
		name          string
		arg           *domain.Account
		builtSts      func(mock sqlmock.Sqlmock)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "fail on exec query",
			arg:  a,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `UPDATE accounts SET created_at=$1,updated_at=$2,balance=$3 WHERE id = $4`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(a.CreatedAt, a.UpdatedAt, a.Balance, a.ID).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail on get affected row",
			arg:  a,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `UPDATE accounts SET created_at=$1,updated_at=$2,balance=$3 WHERE id = $4`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(a.CreatedAt, a.UpdatedAt, a.Balance, a.ID).WillReturnResult(sqlmock.NewErrorResult(errors.New("unexpected error")))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail on invalid number of affected row",
			arg:  a,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `UPDATE accounts SET created_at=$1,updated_at=$2,balance=$3 WHERE id = $4`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(a.CreatedAt, a.UpdatedAt, a.Balance, a.ID).WillReturnResult(sqlmock.NewResult(1, 2))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			arg:  a,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `UPDATE accounts SET created_at=$1,updated_at=$2,balance=$3 WHERE id = $4`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(a.CreatedAt, a.UpdatedAt, a.Balance, a.ID).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			repo := pg.NewAccountRepository(db)
			tc.builtSts(mock)
			err = repo.Update(context.TODO(), tc.arg)
			tc.checkResponse(t, err)
		})
	}
}

func Test_AccountRepo_Delete(t *testing.T) {
	id := uuid.NewV4().String()
	testCases := []struct {
		name          string
		arg           string
		builtSts      func(mock sqlmock.Sqlmock)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "fail on exec query",
			arg:  id,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `DELETE FROM accounts WHERE id = $1`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(id).WillReturnError(errors.New("unexpected error"))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail on get affected row",
			arg:  id,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `DELETE FROM accounts WHERE id = $1`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(id).WillReturnResult(sqlmock.NewErrorResult(errors.New("unexpected error")))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail on invalid number of affected row",
			arg:  id,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `DELETE FROM accounts WHERE id = $1`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 2))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			arg:  id,
			builtSts: func(mock sqlmock.Sqlmock) {
				query := `DELETE FROM accounts WHERE id = $1`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			checkResponse: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			repo := pg.NewAccountRepository(db)
			tc.builtSts(mock)
			err = repo.Delete(context.TODO(), tc.arg)
			tc.checkResponse(t, err)
		})
	}
}
