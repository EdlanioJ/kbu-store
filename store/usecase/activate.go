package usecase

import (
	"context"
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
)

type activateStore struct {
	getStoreByIDRepo domain.GetStoreByIDRepository
	updateStoreRepo  domain.UpdateStoreRepository
	contextTimeout   time.Duration
}

func NewActivateStore(getStore domain.GetStoreByIDRepository, updateStore domain.UpdateStoreRepository, t time.Duration) *activateStore {
	return &activateStore{
		getStoreByIDRepo: getStore,
		updateStoreRepo:  updateStore,
		contextTimeout:   t,
	}
}

func (s *activateStore) Exec(c context.Context, id string) (err error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	store, err := s.getStoreByIDRepo.Exec(ctx, id)
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
	return s.updateStoreRepo.Exec(ctx, store)
}
