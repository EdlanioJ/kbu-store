package usecase

import (
	"context"
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
)

type fetchTagsByCategory struct {
	fetchTagsByCategoryRepo domain.FetchTagsByCategoryRepository
	contextTimeout          time.Duration
}

func NewFetchTagsByCategory(tf domain.FetchTagsByCategoryRepository, tc time.Duration) *fetchTagsByCategory {
	return &fetchTagsByCategory{
		fetchTagsByCategoryRepo: tf,
		contextTimeout:          tc,
	}
}

func (u *fetchTagsByCategory) Exec(c context.Context, categoryID, sort string, page, limit int) (res []*domain.Tag, total int64, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if limit <= 0 {
		limit = 10
	}
	if sort == "" {
		sort = "total DESC"
	}

	if page <= 0 {
		page = 1
	}

	res, total, err = u.fetchTagsByCategoryRepo.Exec(ctx, categoryID, sort, page, limit)
	if err != nil {
		total = 0
		return
	}
	return
}
