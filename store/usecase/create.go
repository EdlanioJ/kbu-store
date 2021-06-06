package usecase

import (
	"context"
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
)

type createStore struct {
	createStoreRepo     domain.CreateStoreRepository
	createAccountRepo   domain.CreateAccountRepository
	getCategoryByIDRepo domain.GetCategoryByIDRepository
	contextTimeout      time.Duration
}

func NewCreateStore(s domain.CreateStoreRepository, a domain.CreateAccountRepository, st domain.GetCategoryByIDRepository, t time.Duration) *createStore {
	return &createStore{
		createStoreRepo:     s,
		createAccountRepo:   a,
		getCategoryByIDRepo: st,
		contextTimeout:      t,
	}
}

func (s *createStore) Add(c context.Context, name, description, CategoryID, externalID string, tags []string, lat, lng float64) (err error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	account, err := domain.NewAccount(0)
	if err != nil {
		return err
	}

	err = s.createAccountRepo.Add(ctx, account)
	if err != nil {
		return err
	}
	Category, err := s.getCategoryByIDRepo.Exec(ctx, CategoryID)
	if err != nil {
		return err
	}

	store, err := domain.NewStore(name, description, externalID, Category, account, tags, lat, lng)
	if err != nil {
		return err
	}

	return s.createStoreRepo.Add(ctx, store)
}
