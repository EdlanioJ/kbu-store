package gorm

import (
	"context"

	"github.com/EdlanioJ/kbu-store/app/domain"
	"github.com/opentracing/opentracing-go"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "accountRepository.Create")
	defer span.Finish()

	err = r.db.WithContext(ctx).
		Table("accounts").
		Create(account).
		Error
	return
}

func (r *accountRepository) FindByID(ctx context.Context, id string) (res *domain.Account, err error) {
	account := &domain.Account{}

	span, ctx := opentracing.StartSpanFromContext(ctx, "accountRepository.FindByID")
	defer span.Finish()

	err = r.db.WithContext(ctx).
		Table("accounts").
		First(account, "id = ?", id).
		Error
	res = account

	return
}

func (r *accountRepository) Update(ctx context.Context, account *domain.Account) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "accountRepository.Update")
	defer span.Finish()

	err = r.db.WithContext(ctx).
		Table("accounts").
		Save(account).
		Error
	return
}

func (r *accountRepository) Delete(ctx context.Context, id string) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "accountRepository.Delete")
	defer span.Finish()

	err = r.db.WithContext(ctx).
		Table("accounts").
		Delete(&domain.Account{}, "id = ?", id).
		Error

	return
}
