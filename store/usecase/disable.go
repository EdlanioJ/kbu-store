package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
)

type disableStore struct {
	getStoreByIDRepo domain.GetStoreByIDRepository
	updateStoreRepo  domain.UpdateStoreRepository
	contextTimeout   time.Duration
}

func NewDisableStore(getStore domain.GetStoreByIDRepository, updateStore domain.UpdateStoreRepository, t time.Duration) *disableStore {
	return &disableStore{
		getStoreByIDRepo: getStore,
		updateStoreRepo:  updateStore,
		contextTimeout:   t,
	}
}

func (s *disableStore) Exec(c context.Context, id string) (err error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	store, err := s.getStoreByIDRepo.Exec(ctx, id)
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
		return errors.New("store is block")
	}
	err = store.Inactivate()
	err = store.Activate()
	if err != nil {
		return err
	}

	return s.updateStoreRepo.Exec(ctx, store)
}
