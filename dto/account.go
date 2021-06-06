package dto

import (
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/shopspring/decimal"
)

type AccountDBModel struct {
	Base
	Balance decimal.Decimal `gorm:"type:decimal(20,8)"`
	Store   *StoreDBModel   `gorm:"ForeignKey:AccountID"`
}

func (AccountDBModel) TableName() string {
	return "accounts"
}

func (a *AccountDBModel) ParserToDBModel(d *domain.Account) {
	a.ID = d.ID
	a.Balance = d.Balance
	a.CreatedAt = d.CreatedAt
	a.UpdatedAt = d.UpdatedAt
}

func (a *AccountDBModel) ParserToAccountDomain() (res *domain.Account) {
	res = new(domain.Account)

	res.ID = a.ID
	res.Balance = a.Balance
	res.CreatedAt = a.CreatedAt
	res.UpdatedAt = a.UpdatedAt

	return res
}
