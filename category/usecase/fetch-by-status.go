package usecase

import (
	"context"
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
)

type fetchCategoryByStatus struct {
	fetchCategoryBysStatusRepository domain.FetchCategoryByStatusRepository
	contextTimeout                   time.Duration
}

func NewFetchCategoryByStatus(r domain.FetchCategoryByStatusRepository, t time.Duration) *fetchCategoryByStatus {
	return &fetchCategoryByStatus{
		fetchCategoryBysStatusRepository: r,
		contextTimeout:                   t,
	}
}

func (st *fetchCategoryByStatus) Exec(c context.Context, status, sort string, page, limit int) (res []*domain.Category, total int64, err error) {
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
	res, total, err = st.fetchCategoryBysStatusRepository.Exec(ctx, status, sort, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return
}
