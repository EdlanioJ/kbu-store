package usecase

import (
	"context"
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
)

type blockStore struct {
	getStoreByIDRepo domain.GetStoreByIDRepository
	updateStoreRepo  domain.UpdateStoreRepository
	contextTimeout   time.Duration
}

func NewBlockStore(getStore domain.GetStoreByIDRepository, updateStore domain.UpdateStoreRepository, t time.Duration) *blockStore {
	return &blockStore{
		getStoreByIDRepo: getStore,
		updateStoreRepo:  updateStore,
		contextTimeout:   t,
	}
}

func (s *blockStore) Exec(c context.Context, id string) (err error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	store, err := s.getStoreByIDRepo.Exec(ctx, id)
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

	return s.updateStoreRepo.Exec(ctx, store)
}
