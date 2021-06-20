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

	err = u.categoryRepo.Create(ctx, category)
	return
}

func (u *CategoryUsecase) GetById(c context.Context, id string) (res *domain.Category, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err = u.categoryRepo.GetById(ctx, id)
	return
}

func (u *CategoryUsecase) GetByIdAndStatus(c context.Context, id, status string) (res *domain.Category, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err = u.categoryRepo.GetByIdAndStatus(ctx, id, status)
	return
}

func (u *CategoryUsecase) GetAll(c context.Context, sort string, page, limit int) (res []*domain.Category, total int64, err error) {
	if limit <= 0 {
		limit = 10
	}
	if sort == "" {
		sort = "created_at"
	}
	if page <= 0 {
		page = 1
	}
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, total, err = u.categoryRepo.GetAll(ctx, sort, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return
}

func (u *CategoryUsecase) GetAllByStatus(c context.Context, status, sort string, page, limit int) (res []*domain.Category, total int64, err error) {
	if limit <= 0 {
		limit = 10
	}
	if sort == "" {
		sort = "created_at"
	}
	if page <= 0 {
		page = 1
	}
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, total, err = u.categoryRepo.GetAllByStatus(ctx, status, sort, page, limit)
	if err != nil {
		return nil, 0, err
	}
	return
}

func (u *CategoryUsecase) Update(c context.Context, Category *domain.Category) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	existedCategory, err := u.categoryRepo.GetById(ctx, Category.ID)
	if err != nil {
		return
	}

	if existedCategory.ID == "" {
		return domain.ErrNotFound
	}

	Category.Status = existedCategory.Status
	return u.categoryRepo.Update(ctx, Category)
}
