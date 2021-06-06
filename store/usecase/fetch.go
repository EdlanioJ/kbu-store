package usecase

import (
	"context"
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
)

type fetchStore struct {
	fetchStoreRepo      domain.FetchStoreRepository
	fillCategoryDetails *fillCategoryDetails
	contextTimeout      time.Duration
}

func NewFetchStore(s domain.FetchStoreRepository, c domain.GetCategoryByIDRepository, t time.Duration) *fetchStore {
	f := newFillCategoryDetails(c)

	return &fetchStore{
		fetchStoreRepo:      s,
		fillCategoryDetails: f,
		contextTimeout:      t,
	}
}

func (s *fetchStore) Exec(c context.Context, sort string, limit, page int) (res []*domain.Store, total int64, err error) {
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

	res, total, err = s.fetchStoreRepo.Exec(ctx, sort, limit, page)
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
