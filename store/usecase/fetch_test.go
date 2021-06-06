package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/domain/mocks"
	"github.com/EdlanioJ/kbu-store/store/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_FetchStoreUsesase(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		is := assert.New(t)
		fetchStoreRepo := new(mocks.FetchStoreRepository)

		fetchStoreRepo.
			On("Exec", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
			Return(nil, int64(0), errors.New("Unexpexted Error")).
			Once()

		u := usecase.NewFetchStore(fetchStoreRepo, nil, time.Second*2)

		list, total, err := u.Exec(context.TODO(), "", 0, 0)

		is.Len(list, 0)
		is.Equal(total, int64(0))
		is.Error(err)
		fetchStoreRepo.AssertExpectations(t)
	})

	t.Run("fail on category repo get by id", func(t *testing.T) {
		dm := testMock()
		is := assert.New(t)
		fetchStoreRepo := new(mocks.FetchStoreRepository)
		getCategoryRepo := new(mocks.GetCategoryByIDRepository)

		sort := "created_at"
		page := 1
		limit := 10
		mockListStore := make([]*domain.Store, 0)
		mockListStore = append(mockListStore, dm.store)

		fetchStoreRepo.
			On("Exec", mock.Anything, sort, limit, page).
			Return(mockListStore, int64(1), nil).
			Once()
		getCategoryRepo.
			On("Exec", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, errors.New("Unexpexted Error"))

		u := usecase.NewFetchStore(fetchStoreRepo, getCategoryRepo, time.Second*2)

		list, total, err := u.Exec(context.Background(), sort, limit, page)

		is.Len(list, 0)
		is.Equal(total, int64(0))
		is.Error(err)

		fetchStoreRepo.AssertExpectations(t)
		getCategoryRepo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		dm := testMock()
		is := assert.New(t)
		fetchStoreRepo := new(mocks.FetchStoreRepository)
		getCategoryRepo := new(mocks.GetCategoryByIDRepository)

		sort := "created_at"
		page := 1
		limit := 10
		mockListStore := make([]*domain.Store, 0)
		mockListStore = append(mockListStore, dm.store)

		fetchStoreRepo.
			On("Exec", mock.Anything, sort, limit, page).
			Return(mockListStore, int64(1), nil).
			Once()
		getCategoryRepo.
			On("Exec", mock.Anything, mock.AnythingOfType("string")).
			Return(dm.storType, nil)

		u := usecase.NewFetchStore(fetchStoreRepo, getCategoryRepo, time.Second*2)

		list, total, err := u.Exec(context.TODO(), sort, limit, page)

		is.Len(list, 1)
		is.Equal(total, int64(1))
		is.NoError(err)

		fetchStoreRepo.AssertExpectations(t)
		getCategoryRepo.AssertExpectations(t)
	})
}
