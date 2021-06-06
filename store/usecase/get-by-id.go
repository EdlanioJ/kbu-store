package usecase

import (
	"context"
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
)

type getStoreByID struct {
	getStoreByIDRepo domain.GetStoreByIDRepository
	getCategoryRepo  domain.GetCategoryByIDRepository
	contextTimeout   time.Duration
}

func NewGetStoreByID(s domain.GetStoreByIDRepository, st domain.GetCategoryByIDRepository, t time.Duration) *getStoreByID {
	return &getStoreByID{
		getStoreByIDRepo: s,
		getCategoryRepo:  st,
		contextTimeout:   t,
	}
}

func (s *getStoreByID) Exec(c context.Context, id string) (res *domain.Store, err error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	res, err = s.getStoreByIDRepo.Exec(ctx, id)
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
