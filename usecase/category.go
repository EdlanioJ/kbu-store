package usecase

import (
	"context"
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
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

func (u *CategoryUsecase) Create(c context.Context, name string) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	category, err := domain.NewCategory(name)
	if err != nil {
		return
	}

	err = u.categoryRepo.Store(ctx, category)
	return
}

func (u *CategoryUsecase) Update(c context.Context, Category *domain.Category) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	existedCategory, err := u.categoryRepo.FindByID(ctx, Category.ID)
	if err != nil {
		return
	}

	Category.Status = existedCategory.Status
	return u.categoryRepo.Update(ctx, Category)
}
