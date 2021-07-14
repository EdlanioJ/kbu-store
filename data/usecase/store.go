package usecase

import (
	"context"
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type StoreUsecase struct {
	storeRepo      domain.StoreRepository
	accountRepo    domain.AccountRepository
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

func (u *StoreUsecase) fillCategoryDetails(c context.Context, data []*domain.Store) ([]*domain.Store, error) {
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

func (u *StoreUsecase) Create(c context.Context, name, description, categoryID, externalID string, tags []string, lat, lng float64) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

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

	return u.storeRepo.Create(ctx, store)
}

func (u *StoreUsecase) GetById(c context.Context, id string) (res *domain.Store, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err = u.storeRepo.GetById(ctx, id)
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

func (u *StoreUsecase) GetByIdAndOwner(c context.Context, id string, ownerID string) (res *domain.Store, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err = u.storeRepo.GetByIdAndOwner(ctx, id, ownerID)
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

func (u *StoreUsecase) GetAll(c context.Context, sort string, limit, page int) (res []*domain.Store, total int64, err error) {
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

	res, total, err = u.storeRepo.GetAll(ctx, sort, limit, page)
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

func (u *StoreUsecase) GetAllByCategory(c context.Context, categoryID, sort string, limit, page int) (res []*domain.Store, total int64, err error) {
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

	res, total, err = u.storeRepo.GetAllByCategory(ctx, categoryID, sort, limit, page)
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

func (u *StoreUsecase) GetAllByOwner(c context.Context, ownerID, sort string, limit, page int) (res []*domain.Store, total int64, err error) {
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

	res, total, err = u.storeRepo.GetAllByOwner(ctx, ownerID, sort, limit, page)
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

func (u *StoreUsecase) GetAllByStatus(c context.Context, status, sort string, limit, page int) (res []*domain.Store, total int64, err error) {
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

	res, total, err = u.storeRepo.GetAllByStatus(ctx, status, sort, limit, page)
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

func (u *StoreUsecase) GetAllByTags(c context.Context, tags []string, sort string, limit, page int) (res []*domain.Store, total int64, err error) {
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

	res, total, err = u.storeRepo.GetAllByTags(ctx, tags, sort, limit, page)
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

func (u *StoreUsecase) GetAllByCloseLocation(c context.Context, lat, lng float64, distance int, status string, limit, page int, sort string) (res []*domain.Store, total int64, err error) {
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

	res, total, err = u.storeRepo.GetAllByLocation(ctx, lat, lng, distance, limit, page, status, sort)
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
	defer cancel()

	store, err := u.storeRepo.GetById(ctx, id)
	if err != nil {
		return err
	}

	if store.ID == "" {
		return domain.ErrNotFound
	}

	if store.Status == domain.StoreStatusBlock {
		return domain.ErrBlocked
	}

	if store.Status == domain.StoreStatusPending {
		return domain.ErrIsPending
	}

	err = store.Block()
	if err != nil {
		return err
	}

	return u.storeRepo.Update(ctx, store)
}

func (u *StoreUsecase) Active(c context.Context, id string) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	store, err := u.storeRepo.GetById(ctx, id)
	if err != nil {
		return err
	}

	if store.ID == "" {
		return domain.ErrNotFound
	}

	if store.Status == domain.StoreStatusActive {
		return domain.ErrActived
	}

	err = store.Activate()
	if err != nil {
		return err
	}
	return u.storeRepo.Update(ctx, store)
}

func (u *StoreUsecase) Disable(c context.Context, id string) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	store, err := u.storeRepo.GetById(ctx, id)
	if err != nil {
		return err
	}

	if store.ID == "" {
		return domain.ErrNotFound
	}

	if store.Status == domain.StoreStatusInactive {
		return domain.ErrInactived
	}

	if store.Status == domain.StoreStatusBlock {
		return domain.ErrBlocked
	}
	err = store.Inactivate()
	if err != nil {
		return err
	}

	return u.storeRepo.Update(ctx, store)
}

func (u *StoreUsecase) Update(c context.Context, store *domain.Store) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	existedStore, err := u.storeRepo.GetById(ctx, store.ID)
	if err != nil {
		return err
	}

	if existedStore.ID == "" {
		return domain.ErrNotFound
	}

	store.UpdatedAt = time.Now()

	return u.storeRepo.Update(ctx, store)
}

func (u *StoreUsecase) Delete(c context.Context, id string) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	store, err := u.storeRepo.GetById(ctx, id)
	if err != nil {
		return
	}

	if store.ID == "" {
		return domain.ErrNotFound
	}

	err = u.accountRepo.Delete(ctx, store.AccountID)
	if err != nil {
		return
	}

	return u.storeRepo.Delete(ctx, id)
}
