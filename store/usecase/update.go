package usecase

import (
	"context"
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
)

type updateStore struct {
	getStoreByIDRepo domain.GetStoreByIDRepository
	updateStoreRepo  domain.UpdateStoreRepository
	contextTimeout   time.Duration
}

func NewUpdateStore(s domain.GetStoreByIDRepository, su domain.UpdateStoreRepository, t time.Duration) *updateStore {
	return &updateStore{
		getStoreByIDRepo: s,
		updateStoreRepo:  su,
		contextTimeout:   t,
	}
}

func (s *updateStore) Exec(c context.Context, store *domain.Store) (err error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	existedStore, err := s.getStoreByIDRepo.Exec(ctx, store.ID)
	if err != nil {
		return err
	}

	if existedStore.ID == "" {
		return domain.ErrNotFound
	}

	store.Status = existedStore.Status
	store.UpdatedAt = time.Now()

	return s.updateStoreRepo.Exec(ctx, store)
}
