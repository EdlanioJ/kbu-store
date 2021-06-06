package usecase

import (
	"context"
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
)

type createCategory struct {
	createCategoryRepository domain.CreateCategoryRepository
	contextTimeout           time.Duration
}

func NewCreateCategory(r domain.CreateCategoryRepository, t time.Duration) *createCategory {
	return &createCategory{
		createCategoryRepository: r,
		contextTimeout:           t,
	}
}

func (st *createCategory) Add(c context.Context, name string) (err error) {
	ctx, cancel := context.WithTimeout(c, st.contextTimeout)
	defer cancel()

	Category, err := domain.NewCategory(name)
	if err != nil {
		return
	}

	err = st.createCategoryRepository.Add(ctx, Category)
	return
}
