package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/EdlanioJ/kbu-store/category/usecase"
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_FetchCategoryByStatusUsecase(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		is := assert.New(t)

		status := domain.CategoryStatusPending
		fetchCategoryByStatusRepository := new(mocks.FetchCategoryByStatusRepository)
		fetchCategoryByStatusRepository.On("Exec", mock.Anything, status, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(nil, int64(0), errors.New("Unexpexted Error")).Once()

		u := usecase.NewFetchCategoryByStatus(fetchCategoryByStatusRepository, time.Second*2)
		list, total, err := u.Exec(context.TODO(), status, "", 0, 0)

		is.Len(list, 0)
		is.Equal(total, int64(0))
		is.Error(err)
		fetchCategoryByStatusRepository.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		is := assert.New(t)
		category := setupStore()

		mockListCategory := make([]*domain.Category, 0)
		mockListCategory = append(mockListCategory, category)

		status := domain.CategoryStatusPending
		fetchCategoryByStatusRepository := new(mocks.FetchCategoryByStatusRepository)
		fetchCategoryByStatusRepository.
			On("Exec", mock.Anything, status, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
			Return(mockListCategory, int64(1), nil).
			Once()

		u := usecase.NewFetchCategoryByStatus(fetchCategoryByStatusRepository, time.Second*2)
		list, total, err := u.Exec(context.TODO(), status, "", 0, 0)

		is.Len(list, 1)
		is.Equal(total, int64(1))
		is.NoError(err)
		fetchCategoryByStatusRepository.AssertExpectations(t)
	})
}
