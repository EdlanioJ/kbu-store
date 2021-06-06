package usecase

import (
	"context"
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
)

type getCategoryByID struct {
	getCategoryByIDRepository domain.GetCategoryByIDRepository
	contextTimeout            time.Duration
}

func NewGetCategoryByID(r domain.GetCategoryByIDRepository, t time.Duration) *getCategoryByID {
	return &getCategoryByID{
		getCategoryByIDRepository: r,
		contextTimeout:            t,
	}
}

func (st *getCategoryByID) Exec(c context.Context, id string) (res *domain.Category, err error) {
	ctx, cancel := context.WithTimeout(c, st.contextTimeout)
	defer cancel()

	res, err = st.getCategoryByIDRepository.Exec(ctx, id)
	return
}
