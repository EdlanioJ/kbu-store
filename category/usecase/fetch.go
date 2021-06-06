package usecase

import (
	"context"
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
)

type fetchCategory struct {
	fetchCategoryRepository domain.FetchCategoryRepository
	contextTimeout          time.Duration
}

func NewFetchCategory(r domain.FetchCategoryRepository, t time.Duration) *fetchCategory {
	return &fetchCategory{
		fetchCategoryRepository: r,
		contextTimeout:          t,
	}
}

func (st *fetchCategory) Exec(c context.Context, sort string, page, limit int) (res []*domain.Category, total int64, err error) {
	if limit <= 0 {
		limit = 10
	}
	if sort == "" {
		sort = "created_at"
	}
	if page <= 0 {
		page = 1
	}
	ctx, cancel := context.WithTimeout(c, st.contextTimeout)
	defer cancel()

	res, total, err = st.fetchCategoryRepository.Exec(ctx, sort, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return
}
