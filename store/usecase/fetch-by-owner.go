package usecase

import (
	"context"
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
)

type fetchStoreByOwner struct {
	fetchStoreByOwnerRepo domain.FetchStoreByOwnerRepository
	fillCategoryDetails   *fillCategoryDetails
	contextTimeout        time.Duration
}

func NewFetchStoreByOwner(
	s domain.FetchStoreByStatusRepository,
	c domain.GetCategoryByIDRepository,
	t time.Duration,
) *fetchStoreByOwner {
	f := newFillCategoryDetails(c)

	return &fetchStoreByOwner{
		fetchStoreByOwnerRepo: s,
		fillCategoryDetails:   f,
		contextTimeout:        t,
	}
}

func (s *fetchStoreByOwner) Exec(c context.Context, ownerID, sort string, limit, page int) (res []*domain.Store, total int64, err error) {
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

	res, total, err = s.fetchStoreByOwnerRepo.Exec(ctx, ownerID, sort, limit, page)
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
