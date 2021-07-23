package domain

import (
	"context"
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

// Account struct
type Account struct {
	Base    `valid:"required"`
	Balance decimal.Decimal `json:"balance,omitempty" valid:"-"`
}

type (
	AccountRepository interface {
		Store(ctx context.Context, account *Account) error
		FindByID(ctx context.Context, id string) (*Account, error)
		Update(ctx context.Context, account *Account) error
		Delete(ctx context.Context, id string) error
	}
)

// Account entity validator
func (a *Account) isValid() (err error) {
	_, err = govalidator.ValidateStruct(a)

	if err != nil {
		return
	}

	if a.Balance.LessThan(decimal.NewFromFloat(0)) {
		return errors.New("balance must positive number")
	}

	return
}

// NewAccount creates an *Account struct
func NewAccount(balance float64) (account *Account, err error) {
	account = &Account{
		Balance: decimal.NewFromFloat(balance),
	}

	account.ID = uuid.NewV4().String()
	account.CreatedAt = time.Now()
	err = account.isValid()

	if err != nil {
		return nil, err
	}

	return
}

func (a *Account) Deposit(amount float64) (err error) {
	if govalidator.IsNegative(amount) {
		return errors.New("value must be a positive number")
	}
	a.Balance = a.Balance.Add(decimal.NewFromFloat(amount))

	err = a.isValid()

	return
}

func (a *Account) Withdow(amount float64) (err error) {
	if govalidator.IsNegative(amount) {
		return errors.New("value must be a positive number")
	}
	a.Balance = a.Balance.Sub(decimal.NewFromFloat(amount))

	err = a.isValid()

	return
}
