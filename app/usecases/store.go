package usecases

import (
	"context"
	"time"

	"github.com/EdlanioJ/kbu-store/app/domain"
	"github.com/EdlanioJ/kbu-store/app/interfaces"
)

type StoreUsecase struct {
	storeRepo      domain.StoreRepository
	accountRepo    domain.AccountRepository
	categoryRepo   domain.CategoryRepository
	msgProducer    interfaces.MessengerProducer
	contextTimeout time.Duration
}

func NewStoreUsecase(s domain.StoreRepository, a domain.AccountRepository, c domain.CategoryRepository, m interfaces.MessengerProducer, t time.Duration) *StoreUsecase {
	return &StoreUsecase{
		storeRepo:      s,
		accountRepo:    a,
		categoryRepo:   c,
		msgProducer:    m,
		contextTimeout: t,
	}
}

func (u *StoreUsecase) Store(c context.Context, name, description, categoryID, externalID string, tags []string, lat, lng float64) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	category, err := u.categoryRepo.FindByID(ctx, categoryID)
	if err != nil {
		return err
	}

	if category.Status != domain.CategoryStatusActive {
		return domain.ErrNotFound
	}

	account, err := domain.NewAccount(0)
	if err != nil {
		return err
	}

	err = u.accountRepo.Store(ctx, account)
	if err != nil {
		return err
	}

	store, err := domain.NewStore(name, description, externalID, category.ID, account.ID, tags, lat, lng)
	if err != nil {
		return err
	}

	err = u.storeRepo.Create(ctx, store)
	if err != nil {
		return err
	}

	storeJson := store.ToJson()
	return u.msgProducer.Publish(ctx, string(storeJson), "store.new")
}

func (u *StoreUsecase) Get(c context.Context, id string) (res *domain.Store, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err = u.storeRepo.FindByID(ctx, id)
	if err != nil {
		return
	}

	return
}

func (u *StoreUsecase) Index(c context.Context, sort string, limit, page int) (res domain.Stores, total int64, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if limit <= 0 {
		limit = 10
	}
	if sort == "" {
		sort = "created_at DESC"
	}
	if page <= 0 {
		page = 1
	}

	res, total, err = u.storeRepo.FindAll(ctx, sort, limit, page)
	if err != nil {
		total = 0
		return
	}

	return
}

func (u *StoreUsecase) Block(c context.Context, id string) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	store, err := u.storeRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	err = store.Block()
	if err != nil {
		return err
	}

	err = u.storeRepo.Update(ctx, store)
	if err != nil {
		return err
	}

	storeJson := store.ToJson()
	return u.msgProducer.Publish(ctx, string(storeJson), "store.update")
}

func (u *StoreUsecase) Active(c context.Context, id string) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	store, err := u.storeRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	err = store.Activate()
	if err != nil {
		return err
	}

	err = u.storeRepo.Update(ctx, store)
	if err != nil {
		return err
	}

	storeJson := store.ToJson()
	return u.msgProducer.Publish(ctx, string(storeJson), "store.update")
}

func (u *StoreUsecase) Disable(c context.Context, id string) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer func() {
		cancel()
	}()

	store, err := u.storeRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	err = store.Disable()
	if err != nil {
		return err
	}

	err = u.storeRepo.Update(ctx, store)
	if err != nil {
		return err
	}

	storeJson := store.ToJson()
	return u.msgProducer.Publish(ctx, string(storeJson), "store.update")
}

func (u *StoreUsecase) Update(c context.Context, store *domain.Store) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	_, err = u.storeRepo.FindByID(ctx, store.ID)
	if err != nil {
		return err
	}

	store.UpdatedAt = time.Now()

	err = u.storeRepo.Update(ctx, store)
	if err != nil {
		return err
	}

	storeJson := store.ToJson()
	return u.msgProducer.Publish(ctx, string(storeJson), "store.update")
}

func (u *StoreUsecase) Delete(c context.Context, id string) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	store, err := u.storeRepo.FindByID(ctx, id)
	if err != nil {
		return
	}

	err = u.storeRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	err = u.accountRepo.Delete(ctx, store.AccountID)
	if err != nil {
		return
	}

	storeJson := store.ToJson()
	return u.msgProducer.Publish(ctx, string(storeJson), "store.delete")
}