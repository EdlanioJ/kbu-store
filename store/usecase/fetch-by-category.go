package usecase

import (
	"context"
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
)

type fetchStoreByCategory struct {
	fetchStoreByTypeRepo domain.FetchStoreByTypeRepository
	fillCategoryDetails  *fillCategoryDetails
	contextTimeout       time.Duration
}

func NewFetchStoreByCategory(s domain.FetchStoreByTypeRepository, c domain.GetCategoryByIDRepository, t time.Duration) *fetchStoreByCategory {
	f := newFillCategoryDetails(c)

	return &fetchStoreByCategory{
		fetchStoreByTypeRepo: s,
		fillCategoryDetails:  f,
		contextTimeout:       t,
	}
}

func (s *fetchStoreByCategory) Exec(c context.Context, typeID, sort string, limit, page int) (res []*domain.Store, total int64, err error) {
	if limit <= 0 {
		limit = 10
	}
	if sort == "" {
		sort = "created_at DESC"
	}
	if page <= 0 {
		page = 1
	}

	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	res, total, err = s.fetchStoreByTypeRepo.Exec(ctx, typeID, sort, limit, page)
	if err != nil {
		return nil, 0, err
	}

	res, err = s.fillCategoryDetails.exec(ctx, res)
	if err != nil {
		total = 0
		return
	}
	return
}
