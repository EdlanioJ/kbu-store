package gorm

import (
	"context"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"gorm.io/gorm"
)

type gormGetAccountByID struct {
	db *gorm.DB
}

func NewGormGetAccountByID(db *gorm.DB) *gormGetAccountByID {
	return &gormGetAccountByID{
		db: db,
	}
}

func (r *gormGetAccountByID) Exec(ctx context.Context, id string) (res *domain.Account, err error) {
	account := &dto.AccountDBModel{}

	err = r.db.WithContext(ctx).
		Table("accounts").
		First(account, "id = ?", id).
		Error
	res = account.ParserToAccountDomain()

	return
}
