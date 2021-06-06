package usecase

import (
	"context"
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
)

type fetchStoreByStatus struct {
	fetchStoreByStatusRepo domain.FetchStoreByStatusRepository
	fillCategoryDetails    *fillCategoryDetails
	contextTimeout         time.Duration
}

func NewFetchStoreByStatus(s domain.FetchStoreByStatusRepository, c domain.GetCategoryByIDRepository, t time.Duration) *fetchStoreByStatus {
	f := newFillCategoryDetails(c)
	return &fetchStoreByStatus{
		fetchStoreByStatusRepo: s,
		fillCategoryDetails:    f,
		contextTimeout:         t,
	}
}

func (s *fetchStoreByStatus) Exec(c context.Context, status, sort string, limit, page int) (res []*domain.Store, total int64, err error) {
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

	res, total, err = s.fetchStoreByStatusRepo.Exec(ctx, status, sort, limit, page)
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
