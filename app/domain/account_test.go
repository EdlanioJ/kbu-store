package domain_test

import (
	"testing"

	"github.com/EdlanioJ/kbu-store/app/domain"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestAccount(t *testing.T) {
	t.Parallel()

	t.Run("new_account", func(t *testing.T) {
		account := domain.NewAccount()
		assert.NotNil(t, account)
		assert.NotEmpty(t, account.ID)
	})
	t.Run("deposit", func(t *testing.T) {
		t.Parallel()
		t.Run("failure_negative_amount", func(t *testing.T) {
			account := domain.NewAccount()
			amount := -10.00
			err := account.Deposit(amount)
			assert.Error(t, err)
			assert.True(t, account.Balance.Equal(decimal.NewFromFloat(0)))
		})

		t.Run("success", func(t *testing.T) {
			account := domain.NewAccount()
			amount := 10.00
			err := account.Deposit(amount)
			assert.NoError(t, err)
			assert.True(t, account.Balance.Equal(decimal.NewFromFloat(10)))
		})
	})
	t.Run("withdow", func(t *testing.T) {
		t.Parallel()
		t.Run("failure_negative_amount", func(t *testing.T) {
			account := domain.NewAccount()
			amount := -10.00
			err := account.Withdow(amount)
			assert.Error(t, err)
			assert.True(t, account.Balance.Equal(decimal.NewFromFloat(0)))
		})

		t.Run("failure_balance_less_than_amount", func(t *testing.T) {
			account := domain.NewAccount()
			amount := 10.00
			err := account.Withdow(amount)
			assert.Error(t, err)
			assert.True(t, account.Balance.Equal(decimal.NewFromFloat(0)))
		})
		t.Run("success", func(t *testing.T) {
			account := domain.NewAccount()
			account.Balance = decimal.NewFromFloat(100)
			amount := 10.00
			err := account.Withdow(amount)
			assert.NoError(t, err)
			assert.True(t, account.Balance.Equal(decimal.NewFromFloat(90)))

		})
	})
}
