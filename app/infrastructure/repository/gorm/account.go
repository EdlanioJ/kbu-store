package gorm

import (
	"context"

	"github.com/EdlanioJ/kbu-store/app/domain"
	"gorm.io/gorm"
)

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *accountRepository {
	return &accountRepository{
		db: db,
	}
}

func (r *accountRepository) Store(ctx context.Context, account *domain.Account) (err error) {
	err = r.db.WithContext(ctx).
		Table("accounts").
		Create(account).
		Error
	return
}

func (r *accountRepository) FindByID(ctx context.Context, id string) (res *domain.Account, err error) {
	account := &domain.Account{}

	err = r.db.WithContext(ctx).
		Table("accounts").
		First(account, "id = ?", id).
		Error
	res = account

	return
}

func (r *accountRepository) Update(ctx context.Context, account *domain.Account) (err error) {
	err = r.db.WithContext(ctx).
		Table("accounts").
		Save(account).
		Error
	return
}

func (r *accountRepository) Delete(ctx context.Context, id string) (err error) {
	err = r.db.WithContext(ctx).
		Table("accounts").
		Delete(&domain.Account{}, "id = ?", id).
		Error

	return
}
