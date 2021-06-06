package usecase

import (
	"context"
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
)

type fetchStoreByCloseLocation struct {
	fetchStoreByLocationRepo domain.FetchStoreByCloseLocationRepository
	fillCategoryDetails      *fillCategoryDetails
	contextTimeout           time.Duration
}

func NewFetchStoreByCloseLocation(
	s domain.FetchStoreByCloseLocationRepository,
	c domain.GetCategoryByIDRepository,
	t time.Duration,
) *fetchStoreByCloseLocation {
	f := newFillCategoryDetails(c)

	return &fetchStoreByCloseLocation{
		fetchStoreByLocationRepo: s,
		fillCategoryDetails:      f,
		contextTimeout:           t,
	}
}

func (s *fetchStoreByCloseLocation) Exec(c context.Context, lat, lng float64, distance int, status string, limit, page int, sort string) (res []*domain.Store, total int64, err error) {
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

	res, total, err = s.fetchStoreByLocationRepo.Exec(ctx, lat, lng, distance, limit, page, status, sort)
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
