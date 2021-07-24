package gorm

import (
	"context"

	"github.com/EdlanioJ/kbu-store/app/db/model"
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
	accountEntity := &model.Account{}
	accountEntity.FromAccountDomain(account)

	err = r.db.WithContext(ctx).
		Table("accounts").
		Create(accountEntity).
		Error
	return
}

func (r *accountRepository) FindByID(ctx context.Context, id string) (res *domain.Account, err error) {
	account := &model.Account{}

	err = r.db.WithContext(ctx).
		Table("accounts").
		First(account, "id = ?", id).
		Error
	res = account.ToAccountDomain()

	return
}

func (r *accountRepository) Update(ctx context.Context, account *domain.Account) (err error) {
	accountModel := &model.Account{}
	accountModel.FromAccountDomain(account)

	err = r.db.WithContext(ctx).
		Table("accounts").
		Save(accountModel).
		Error
	return
}

func (r *accountRepository) Delete(ctx context.Context, id string) (err error) {
	err = r.db.WithContext(ctx).
		Table("accounts").
		Delete(&model.Account{}, "id = ?", id).
		Error

	return
}
