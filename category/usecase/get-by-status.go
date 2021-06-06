package usecase

import (
	"context"
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
)

type getCategoryByStatus struct {
	getCategoryByStautsRepository domain.GetCategoryByStautsRepository
	contextTimeout                time.Duration
}

func NewGetCategoryByStatus(r domain.GetCategoryByStautsRepository, t time.Duration) *getCategoryByStatus {
	return &getCategoryByStatus{
		getCategoryByStautsRepository: r,
		contextTimeout:                t,
	}
}

func (st *getCategoryByStatus) Exec(c context.Context, id, status string) (res *domain.Category, err error) {
	ctx, cancel := context.WithTimeout(c, st.contextTimeout)
	defer cancel()

	res, err = st.getCategoryByStautsRepository.Exec(ctx, id, status)
	return
}
