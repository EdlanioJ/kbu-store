package pg_test

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/app/domain"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/repository/pg"
	"github.com/EdlanioJ/kbu-store/app/utils/sample"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func Test_AccountRepo_Store(t *testing.T) {
	a := sample.NewAccount()
	testCases := []struct {
		name        string
		arg         *domain.Account
		expectedErr bool
		prepare     func(mock sqlmock.Sqlmock)
	}{
		{
			name:        "failure_exec_query_returns_error",
			arg:         a,
			expectedErr: true,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO accounts (id,created_at,updated_at,balance) VALUES ($1,$2,$3,$4)`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(a.ID, a.CreatedAt, a.UpdatedAt, a.Balance).WillReturnError(errors.New("unexpected error"))
			},
		},
		{
			name:        "failure_get_affected_row_returns_error",
			arg:         a,
			expectedErr: true,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO accounts (id,created_at,updated_at,balance) VALUES ($1,$2,$3,$4)`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(a.ID, a.CreatedAt, a.UpdatedAt, a.Balance).WillReturnResult(sqlmock.NewErrorResult(errors.New("unexpected error")))
			},
		},
		{
			name:        "failure_get_invalid_numbe_of_afected_rows",
			arg:         a,
			expectedErr: true,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO accounts (id,created_at,updated_at,balance) VALUES ($1,$2,$3,$4)`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(a.ID, a.CreatedAt, a.UpdatedAt, a.Balance).WillReturnResult(sqlmock.NewResult(1, 2))
			},
		},
		{
			name: "success",
			arg:  a,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO accounts (id,created_at,updated_at,balance) VALUES ($1,$2,$3,$4)`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(a.ID, a.CreatedAt, a.UpdatedAt, a.Balance).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			repo := pg.NewAccountRepository(db)
			tc.prepare(mock)
			err = repo.Store(context.TODO(), tc.arg)
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_AccountRepo_FindByID(t *testing.T) {
	id := uuid.NewV4().String()
	a := sample.NewAccount()
	testCases := []struct {
		name        string
		arg         string
		expectedErr bool
		prepare     func(mock sqlmock.Sqlmock)
	}{
		{
			name:        "failure_exec_query_returns_error",
			arg:         id,
			expectedErr: true,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `SELECT * FROM accounts WHERE id = $1`
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(id).WillReturnError(errors.New("unexpected error"))
			},
		},
		{
			name: "success",
			arg:  id,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `SELECT * FROM accounts WHERE id = $1`
				row := sqlmock.
					NewRows([]string{"id", "created_at", "updated_at", "balance"}).
					AddRow(a.ID, a.CreatedAt, a.UpdatedAt, a.Balance)

				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(id).WillReturnRows(row)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			repo := pg.NewAccountRepository(db)
			tc.prepare(mock)
			res, err := repo.FindByID(context.TODO(), tc.arg)
			if tc.expectedErr {
				assert.Nil(t, res)
				assert.Error(t, err)
			} else {
				assert.NotNil(t, res)
				assert.NoError(t, err)
			}
		})
	}
}

func Test_AccountRepo_Update(t *testing.T) {
	a := sample.NewAccount()
	testCases := []struct {
		name        string
		arg         *domain.Account
		expectedErr bool
		prepare     func(mock sqlmock.Sqlmock)
	}{
		{
			name:        "failure_exec_query_returns_error",
			arg:         a,
			expectedErr: true,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `UPDATE accounts SET created_at=$1,updated_at=$2,balance=$3 WHERE id = $4`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(a.CreatedAt, a.UpdatedAt, a.Balance, a.ID).WillReturnError(errors.New("unexpected error"))
			},
		},
		{
			name:        "failure_get_affected_row_returns_error",
			arg:         a,
			expectedErr: true,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `UPDATE accounts SET created_at=$1,updated_at=$2,balance=$3 WHERE id = $4`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(a.CreatedAt, a.UpdatedAt, a.Balance, a.ID).WillReturnResult(sqlmock.NewErrorResult(errors.New("unexpected error")))
			},
		},
		{
			name:        "failure_returns_invalid_number_of_affected_row",
			arg:         a,
			expectedErr: true,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `UPDATE accounts SET created_at=$1,updated_at=$2,balance=$3 WHERE id = $4`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(a.CreatedAt, a.UpdatedAt, a.Balance, a.ID).WillReturnResult(sqlmock.NewResult(1, 2))
			},
		},
		{
			name: "success",
			arg:  a,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `UPDATE accounts SET created_at=$1,updated_at=$2,balance=$3 WHERE id = $4`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(a.CreatedAt, a.UpdatedAt, a.Balance, a.ID).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			repo := pg.NewAccountRepository(db)
			tc.prepare(mock)
			err = repo.Update(context.TODO(), tc.arg)
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_AccountRepo_Delete(t *testing.T) {
	id := uuid.NewV4().String()
	testCases := []struct {
		name        string
		arg         string
		expectedErr bool
		prepare     func(mock sqlmock.Sqlmock)
	}{
		{
			name:        "failure_exec_query_returns_error",
			arg:         id,
			expectedErr: true,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `DELETE FROM accounts WHERE id = $1`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(id).WillReturnError(errors.New("unexpected error"))
			},
		},
		{
			name:        "failure_get_affected_row_returns_error",
			arg:         id,
			expectedErr: true,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `DELETE FROM accounts WHERE id = $1`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(id).WillReturnResult(sqlmock.NewErrorResult(errors.New("unexpected error")))
			},
		},
		{
			name:        "failure_returns_invalid_number_of_affected_row",
			arg:         id,
			expectedErr: true,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `DELETE FROM accounts WHERE id = $1`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 2))
			},
		},
		{
			name: "success",
			arg:  id,
			prepare: func(mock sqlmock.Sqlmock) {
				query := `DELETE FROM accounts WHERE id = $1`
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			repo := pg.NewAccountRepository(db)
			tc.prepare(mock)
			err = repo.Delete(context.TODO(), tc.arg)

			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
