package gorm

import (
	"context"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"gorm.io/gorm"
)

type gormCreateAccount struct {
	db *gorm.DB
}

func NewGormCreateAccount(db *gorm.DB) *gormCreateAccount {
	return &gormCreateAccount{
		db: db,
	}
}

func (r *gormCreateAccount) Add(ctx context.Context, account *domain.Account) (err error) {
	accountEntity := &dto.AccountDBModel{}
	accountEntity.ParserToDBModel(account)

	err = r.db.WithContext(ctx).
		Table("accounts").
		Create(accountEntity).
		Error
	return
}
