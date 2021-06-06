package usecase

import (
	"context"
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
)

type updateCategory struct {
	updateCategoryRepository  domain.UpdateCategoryRepository
	getCategoryByIDRepository domain.GetCategoryByIDRepository
	contextTimeout            time.Duration
}

func NewUpdateCategory(ur domain.UpdateCategoryRepository, gr domain.GetCategoryByIDRepository, t time.Duration) *updateCategory {
	return &updateCategory{
		updateCategoryRepository:  ur,
		getCategoryByIDRepository: gr,
		contextTimeout:            t,
	}
}

func (st *updateCategory) Exec(c context.Context, Category *domain.Category) (err error) {
	ctx, cancel := context.WithTimeout(c, st.contextTimeout)
	defer cancel()

	existedCategory, err := st.getCategoryByIDRepository.Exec(ctx, Category.ID)
	if err != nil {
		return
	}

	if existedCategory.ID == "" {
		return domain.ErrNotFound
	}

	Category.Status = existedCategory.Status
	return st.updateCategoryRepository.Exec(ctx, Category)
}
