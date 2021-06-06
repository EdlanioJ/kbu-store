package usecase

import (
	"context"
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
)

type getStoreByOwner struct {
	getStoreByOwnerRepo domain.GetStoreByOwnerRepository
	getCategoryRepo     domain.GetCategoryByIDRepository
	contextTimeout      time.Duration
}

func NewGetStoreByOwner(s domain.GetStoreByOwnerRepository, st domain.GetCategoryByIDRepository, t time.Duration) *getStoreByOwner {
	return &getStoreByOwner{
		getStoreByOwnerRepo: s,
		getCategoryRepo:     st,
		contextTimeout:      t,
	}
}

func (s *getStoreByOwner) Exec(c context.Context, id string, ownerID string) (res *domain.Store, err error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	res, err = s.getStoreByOwnerRepo.Exec(ctx, id, ownerID)
	if err != nil {
		return
	}

	Category, err := s.getCategoryRepo.Exec(ctx, res.Category.ID)
	if err != nil {
		return nil, err
	}

	res.Category = Category
	return
}
