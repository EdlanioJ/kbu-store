package usecase

import (
	"context"
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
)

type fetchStoreByTabs struct {
	fetchStoreByTabsRepo domain.FetchStoreByTagsRepository
	fillCategoryDetails  *fillCategoryDetails
	contextTimeout       time.Duration
}

func NewFetchStoreByTabs(s domain.FetchStoreByTagsRepository, c domain.GetCategoryByIDRepository, t time.Duration) *fetchStoreByTabs {
	f := newFillCategoryDetails(c)

	return &fetchStoreByTabs{
		fetchStoreByTabsRepo: s,
		fillCategoryDetails:  f,
		contextTimeout:       t,
	}
}

func (s *fetchStoreByTabs) Exec(c context.Context, tags []string, sort string, limit, page int) (res []*domain.Store, total int64, err error) {
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

	res, total, err = s.fetchStoreByTabsRepo.Exec(ctx, tags, sort, limit, page)
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
