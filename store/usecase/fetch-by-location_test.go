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

func Test_FetchStoreByCloseLocationUsecase(t *testing.T) {
	t.Run("fail on fetch store by location and status repo", func(t *testing.T) {
		is := assert.New(t)
		dm := testMock()

		fetchStoreByLocationRepo := new(mocks.FetchStoreByCloseLocationRepository)
		location := dm.store.Position
		distance := 10
		status := domain.StoreStatusActive
		fetchStoreByLocationRepo.On("Exec", mock.Anything, location.Lat, location.Lng, distance, mock.AnythingOfType("int"), mock.AnythingOfType("int"), status, mock.AnythingOfType("string")).
			Return(nil, int64(0), errors.New("Unexpexted Error")).
			Once()

		u := usecase.NewFetchStoreByCloseLocation(fetchStoreByLocationRepo, nil, time.Second*2)

		list, total, err := u.Exec(context.TODO(), location.Lat, location.Lng, distance, status, 0, 0, "")

		is.Len(list, 0)
		is.Equal(total, int64(0))
		is.Error(err)
		fetchStoreByLocationRepo.AssertExpectations(t)
	})

	t.Run("fail on get category by id repo", func(t *testing.T) {
		is := assert.New(t)
		dm := testMock()

		mockListStore := make([]*domain.Store, 0)
		mockListStore = append(mockListStore, dm.store)
		fetchStoreByLocationRepo := new(mocks.FetchStoreByCloseLocationRepository)
		getCategoryRepo := new(mocks.GetCategoryByIDRepository)
		location := dm.store.Position
		distance := 10
		status := domain.StoreStatusActive
		fetchStoreByLocationRepo.
			On("Exec", mock.Anything, location.Lat, location.Lng, distance, mock.AnythingOfType("int"), mock.AnythingOfType("int"), status, mock.AnythingOfType("string")).
			Return(mockListStore, int64(1), nil).
			Once()
		getCategoryRepo.
			On("Exec", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, errors.New("Unexpexted Error"))
		u := usecase.NewFetchStoreByCloseLocation(fetchStoreByLocationRepo, getCategoryRepo, time.Second*2)

		list, total, err := u.Exec(context.TODO(), location.Lat, location.Lng, distance, status, 0, 0, "")

		is.Len(list, 0)
		is.Equal(total, int64(0))
		is.Error(err)
		fetchStoreByLocationRepo.AssertExpectations(t)
		getCategoryRepo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		is := assert.New(t)
		dm := testMock()

		mockListStore := make([]*domain.Store, 0)
		mockListStore = append(mockListStore, dm.store)
		fetchStoreByLocationRepo := new(mocks.FetchStoreByCloseLocationRepository)
		getCategoryRepo := new(mocks.GetCategoryByIDRepository)
		location := dm.store.Position
		distance := 10
		status := domain.StoreStatusActive
		fetchStoreByLocationRepo.
			On("Exec", mock.Anything, location.Lat, location.Lng, distance, mock.AnythingOfType("int"), mock.AnythingOfType("int"), status, mock.AnythingOfType("string")).
			Return(mockListStore, int64(1), nil).
			Once()
		getCategoryRepo.
			On("Exec", mock.Anything, mock.AnythingOfType("string")).
			Return(dm.storType, nil)
		u := usecase.NewFetchStoreByCloseLocation(fetchStoreByLocationRepo, getCategoryRepo, time.Second*2)

		list, total, err := u.Exec(context.TODO(), location.Lat, location.Lng, distance, status, 0, 0, "")

		is.Len(list, 1)
		is.Equal(total, int64(1))
		is.NoError(err)
		fetchStoreByLocationRepo.AssertExpectations(t)
		getCategoryRepo.AssertExpectations(t)
	})
}
