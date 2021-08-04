package usecases

import (
	"context"
	"time"

	"github.com/EdlanioJ/kbu-store/app/domain"
	"github.com/opentracing/opentracing-go"
)

type CategoryUsecase struct {
	categoryRepo   domain.CategoryRepository
	contextTimeout time.Duration
}

func NewCategoryUsecase(c domain.CategoryRepository, t time.Duration) *CategoryUsecase {
	return &CategoryUsecase{
		categoryRepo:   c,
		contextTimeout: t,
	}
}

func (u *CategoryUsecase) Create(c context.Context, category *domain.Category) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	span, ctx := opentracing.StartSpanFromContext(ctx, "categoryUsecãse.Store")
	defer span.Finish()

	err = u.categoryRepo.Store(ctx, category)
	return
}

func (u *CategoryUsecase) Update(c context.Context, category *domain.Category) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	span, ctx := opentracing.StartSpanFromContext(ctx, "categoryUsecãse.Update")
	defer span.Finish()

	_, err = u.categoryRepo.FindByID(ctx, category.ID)
	if err != nil {
		return
	}

	return u.categoryRepo.Update(ctx, category)
}
