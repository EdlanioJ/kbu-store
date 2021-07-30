package domain

import (
	"context"
	"fmt"
	"math"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

// Account struct
type Account struct {
	Base
	Balance decimal.Decimal `json:"balance" gorm:"type:decimal(20,8)"`
}

type (
	AccountRepository interface {
		Store(ctx context.Context, account *Account) error
		FindByID(ctx context.Context, id string) (*Account, error)
		Update(ctx context.Context, account *Account) error
		Delete(ctx context.Context, id string) error
	}
)

// NewAccount creates an *Account struct
func NewAccount() (account *Account) {
	account = new(Account)
	account.ID = uuid.NewV4().String()
	account.CreatedAt = time.Now()

	return
}

func (a *Account) Deposit(amount float64) (err error) {
	if math.Signbit(amount) {
		return fmt.Errorf("%f must be a positive number", amount)
	}
	a.Balance = a.Balance.Add(decimal.NewFromFloat(amount))

	return
}

func (a *Account) Withdow(amount float64) (err error) {
	if math.Signbit(amount) {
		return fmt.Errorf("%f must be a positive number", amount)
	}

	if a.Balance.LessThan(decimal.NewFromFloat(amount)) {
		return fmt.Errorf("balance is less than amount: %f", amount)
	}

	a.Balance = a.Balance.Sub(decimal.NewFromFloat(amount))

	return
}
