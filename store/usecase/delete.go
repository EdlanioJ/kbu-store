package usecase

import (
	"context"
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
)

type deleteStore struct {
	getStoreByIDRepo  domain.GetStoreByIDRepository
	deleteStoreRepo   domain.DeleteStoreRepository
	deleteAccountRepo domain.DeleteAccountRepository
	contextTimeout    time.Duration
}

func NewDeleteStore(sd domain.DeleteStoreRepository, da domain.DeleteAccountRepository, s domain.GetStoreByIDRepository, t time.Duration) *deleteStore {
	return &deleteStore{
		getStoreByIDRepo:  s,
		deleteStoreRepo:   sd,
		deleteAccountRepo: da,
		contextTimeout:    t,
	}
}

func (s *deleteStore) Exec(c context.Context, id string) (err error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	store, err := s.getStoreByIDRepo.Exec(ctx, id)
	if err != nil {
		return
	}

	if store.ID == "" {
		return domain.ErrNotFound
	}

	err = s.deleteAccountRepo.Exec(ctx, store.Account.ID)
	if err != nil {
		return
	}

	return s.deleteStoreRepo.Exec(ctx, id)
}
