package gorm

import (
	"context"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"gorm.io/gorm"
)

type gormUpdateAccount struct {
	db *gorm.DB
}

func NewGormUpdateAccount(db *gorm.DB) *gormUpdateAccount {
	return &gormUpdateAccount{
		db: db,
	}
}

func (r *gormUpdateAccount) Exec(ctx context.Context, account *domain.Account) (err error) {
	accountModel := &dto.AccountDBModel{}
	accountModel.ParserToDBModel(account)

	err = r.db.WithContext(ctx).
		Table("accounts").
		Save(accountModel).
		Error
	return
}
