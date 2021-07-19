package usecase

import (
	"context"
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/interfaces"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type StoreUsecase struct {
	storeRepo      domain.StoreRepository
	accountRepo    domain.AccountRepository
	msgProducer    interfaces.MessengerProducer
	categoryRepo   domain.CategoryRepository
	contextTimeout time.Duration
}

func NewStoreUsecase(s domain.StoreRepository, a domain.AccountRepository, c domain.CategoryRepository, t time.Duration) *StoreUsecase {
	return &StoreUsecase{
		storeRepo:      s,
		accountRepo:    a,
		categoryRepo:   c,
		contextTimeout: t,
	}
}

func (u *StoreUsecase) fillCategoryDetails(c context.Context, data domain.Stores) (domain.Stores, error) {
	g, ctx := errgroup.WithContext(c)

	mapCategories := map[string]*domain.Category{}
	for _, store := range data {
		mapCategories[store.Category.ID] = &domain.Category{}
	}

	chanCategory := make(chan *domain.Category)
	for categoryID := range mapCategories {
		categoryID := categoryID
		g.Go(func() error {
			res, err := u.categoryRepo.GetById(ctx, categoryID)
			if err != nil {
				return err
			}
			chanCategory <- res
			return nil
		})
	}

	go func() {
		err := g.Wait()
		if err != nil {
			logrus.Error(err)
		}

		close(chanCategory)
	}()

	for category := range chanCategory {
		if category != (&domain.Category{}) {
			mapCategories[category.ID] = category
		}
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	for index, item := range data {
		if a, ok := mapCategories[item.Category.ID]; ok {
			data[index].Category = a
		}
	}
	return data, nil
}

func (u *StoreUsecase) Store(c context.Context, name, description, categoryID, externalID string, tags []string, lat, lng float64) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer func() {
		cancel()
		u.msgProducer.Close()
	}()

	account, err := domain.NewAccount(0)
	if err != nil {
		return err
	}

	err = u.accountRepo.Create(ctx, account)
	if err != nil {
		return err
	}
	category, err := u.categoryRepo.GetByIdAndStatus(ctx, categoryID, domain.CategoryStatusActive)
	if err != nil {
		return err
	}

	store, err := domain.NewStore(name, description, externalID, category, account.ID, tags, lat, lng)
	if err != nil {
		return err
	}

	err = u.storeRepo.Create(ctx, store)
	if err != nil {
		return err
	}

	storeJson := store.ToJson()
	return u.msgProducer.Publish(ctx, string(storeJson), "stores.create")
}

func (u *StoreUsecase) Get(c context.Context, id string) (res *domain.Store, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err = u.storeRepo.FindByID(ctx, id)
	if err != nil {
		return
	}

	category, err := u.categoryRepo.GetById(ctx, res.Category.ID)
	if err != nil {
		return nil, err
	}

	res.Category = category
	return
}

func (u *StoreUsecase) Index(c context.Context, sort string, limit, page int) (res domain.Stores, total int64, err error) {
	if limit <= 0 {
		limit = 10
	}
	if sort == "" {
		sort = "created_at DESC"
	}
	if page <= 0 {
		page = 1
	}

	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, total, err = u.storeRepo.FindAll(ctx, sort, limit, page)
	if err != nil {
		return nil, 0, err
	}

	res, err = u.fillCategoryDetails(ctx, res)
	if err != nil {
		total = 0
		return
	}
	return
}

func (u *StoreUsecase) Block(c context.Context, id string) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer func() {
		cancel()
		u.msgProducer.Close()
	}()

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
	return u.msgProducer.Publish(ctx, string(storeJson), "stores.update")
}

func (u *StoreUsecase) Active(c context.Context, id string) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer func() {
		cancel()
		u.msgProducer.Close()
	}()

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
	return u.msgProducer.Publish(ctx, string(storeJson), "stores.update")
}

func (u *StoreUsecase) Disable(c context.Context, id string) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer func() {
		cancel()
		u.msgProducer.Close()
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
	return u.msgProducer.Publish(ctx, string(storeJson), "stores.update")
}

func (u *StoreUsecase) Update(c context.Context, store *domain.Store) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer func() {
		cancel()
		u.msgProducer.Close()
	}()
	_, err = u.storeRepo.FindByID(ctx, store.ID)
	if err != nil {
		return err
	}

	store.UpdatedAt = time.Now()

	err = store.Disable()
	if err != nil {
		return err
	}

	err = u.storeRepo.Update(ctx, store)
	if err != nil {
		return err
	}

	storeJson := store.ToJson()
	return u.msgProducer.Publish(ctx, string(storeJson), "stores.update")
}

func (u *StoreUsecase) Delete(c context.Context, id string) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer func() {
		cancel()
		u.msgProducer.Close()
	}()

	store, err := u.storeRepo.FindByID(ctx, id)
	if err != nil {
		return
	}

	err = u.accountRepo.Delete(ctx, store.AccountID)
	if err != nil {
		return
	}

	err = u.storeRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	storeJson := store.ToJson()
	return u.msgProducer.Publish(ctx, string(storeJson), "stores.delete")
}
