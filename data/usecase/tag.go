package usecase

import (
	"context"
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
)

type TagUsecase struct {
	tagRepo        domain.TagRepository
	contextTimeout time.Duration
}

func NewtagUsecase(tr domain.TagRepository, t time.Duration) *TagUsecase {
	return &TagUsecase{
		tagRepo:        tr,
		contextTimeout: t,
	}
}

func (u *TagUsecase) GetAll(c context.Context, sort string, page, limit int) (res []*domain.Tag, total int64, err error) {
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

	res, total, err = u.tagRepo.GetAll(ctx, sort, page, limit)
	if err != nil {
		total = 0
		return
	}

	return
}

func (u *TagUsecase) GetAllByCategory(c context.Context, categoryID, sort string, page, limit int) (res []*domain.Tag, total int64, err error) {
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

	res, total, err = u.tagRepo.GetAllByCategory(ctx, categoryID, sort, page, limit)
	if err != nil {
		total = 0
		return
	}
	return
}
