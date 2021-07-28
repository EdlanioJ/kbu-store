package sample

import (
	"time"

	"github.com/EdlanioJ/kbu-store/app/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

func NewAccount() *domain.Account {
	account := new(domain.Account)
	account.ID = uuid.NewV4().String()
	account.Balance = decimal.NewFromFloat(1000)
	account.CreatedAt = time.Now()

	return account
}
