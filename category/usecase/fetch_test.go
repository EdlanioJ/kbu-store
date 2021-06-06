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

func Test_FetchCategoryUsecase(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		is := assert.New(t)

		fetchCategoryRepository := new(mocks.FetchCategoryRepository)
		fetchCategoryRepository.On("Exec", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(nil, int64(0), errors.New("Unexpexted Error")).Once()

		u := usecase.NewFetchCategory(fetchCategoryRepository, time.Second*2)
		list, total, err := u.Exec(context.TODO(), "", 0, 0)

		is.Len(list, 0)
		is.Equal(total, int64(0))
		is.Error(err)
		fetchCategoryRepository.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		is := assert.New(t)
		category := setupStore()

		mockListStore := make([]*domain.Category, 0)
		mockListStore = append(mockListStore, category)

		fetchCategoryRepository := new(mocks.FetchCategoryRepository)
		fetchCategoryRepository.
			On("Exec", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
			Return(mockListStore, int64(1), nil).
			Once()

		u := usecase.NewFetchCategory(fetchCategoryRepository, time.Second*2)
		list, total, err := u.Exec(context.TODO(), "", 0, 0)

		is.Len(list, 1)
		is.Equal(total, int64(1))
		is.NoError(err)
		fetchCategoryRepository.AssertExpectations(t)
	})
}
