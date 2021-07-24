package model

import (
	"github.com/EdlanioJ/kbu-store/app/domain"
	"github.com/shopspring/decimal"
)

type Account struct {
	Base
	Balance decimal.Decimal `gorm:"type:decimal(20,8)"`
	Store   *Store          `gorm:"ForeignKey:AccountID"`
}

func (Account) TableName() string {
	return "accounts"
}

func (a *Account) FromAccountDomain(d *domain.Account) {
	a.ID = d.ID
	a.Balance = d.Balance
	a.CreatedAt = d.CreatedAt
	a.UpdatedAt = d.UpdatedAt
}

func (a *Account) ToAccountDomain() (res *domain.Account) {
	res = new(domain.Account)

	res.ID = a.ID
	res.Balance = a.Balance
	res.CreatedAt = a.CreatedAt
	res.UpdatedAt = a.UpdatedAt

	return res
}
