package usecase

import (
	"context"
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
)

type fetchTags struct {
	fetchTagsRepo  domain.FetchTagsRepository
	contextTimeout time.Duration
}

func NewFetchTags(ft domain.FetchTagsRepository, tc time.Duration) *fetchTags {
	return &fetchTags{
		fetchTagsRepo:  ft,
		contextTimeout: tc,
	}
}

func (u *fetchTags) Exec(c context.Context, sort string, page, limit int) (res []*domain.Tag, total int64, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if limit <= 0 {
		limit = 10
	}
	if sort == "" {
		sort = "count DESC"
	}

	if page <= 0 {
		page = 1
	}

	res, total, err = u.fetchTagsRepo.Exec(ctx, sort, page, limit)
	if err != nil {
		total = 0
		return
	}

	return
}
