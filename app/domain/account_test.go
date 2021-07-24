package domain_test

import (
	"testing"

	"github.com/EdlanioJ/kbu-store/app/domain"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	t.Run("should fail", func(t *testing.T) {
		is := assert.New(t)

		balance := -1000.00

		account, err := domain.NewAccount(balance)

		is.Nil(account)
		is.Error(err)
	})

	t.Run("should succeed", func(t *testing.T) {
		is := assert.New(t)

		balance := 1000.00

		account, err := domain.NewAccount(balance)

		is.Nil(err)
		is.NotEmpty(account)
		is.Equal(account.Balance, decimal.NewFromFloat(balance))
	})
}

func TestDeposit(t *testing.T) {
	t.Run("should fail if value is negative", func(t *testing.T) {
		is := assert.New(t)
		balance := 10.00
		amount := -10.00

		account, _ := domain.NewAccount(balance)

		err := account.Deposit(amount)

		is.Error(err)
	})

	t.Run("should fail on validation", func(t *testing.T) {
		is := assert.New(t)
		balance := 10.00
		amount := 10.00

		account, _ := domain.NewAccount(balance)

		account.ID = "1"
		err := account.Deposit(amount)

		is.Error(err)
	})

	t.Run("should succeed", func(t *testing.T) {
		is := assert.New(t)
		balance := 10.00
		amount := 10.00

		account, _ := domain.NewAccount(balance)

		err := account.Deposit(amount)

		is.Nil(err)
	})
}

func TestWithdow(t *testing.T) {
	t.Run("should fail if amount is negative", func(t *testing.T) {
		is := assert.New(t)
		balance := 10.00
		amount := -10.00

		account, _ := domain.NewAccount(balance)

		err := account.Withdow(amount)

		is.Error(err)
	})

	t.Run("should fail on validation", func(t *testing.T) {
		is := assert.New(t)
		balance := 10.00
		amount := 100.00

		account, _ := domain.NewAccount(balance)

		err := account.Withdow(amount)

		is.Error(err)
	})

	t.Run("should succeed", func(t *testing.T) {
		is := assert.New(t)
		balance := 10.00
		amount := 10.00

		account, _ := domain.NewAccount(balance)

		err := account.Withdow(amount)

		is.Nil(err)
	})
}
