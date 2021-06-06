package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/domain/mocks"
	"github.com/EdlanioJ/kbu-store/store/usecase"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_FetchStoreByOwnerUsecase(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		is := assert.New(t)
		fetchStoreByOwnerRepo := new(mocks.FetchStoreByOwnerRepository)

		ownerID := uuid.NewV4().String()
		fetchStoreByOwnerRepo.
			On("Exec", mock.Anything, ownerID, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
			Return(nil, int64(0), errors.New("Unexpexted Error")).
			Once()

		u := usecase.NewFetchStoreByOwner(fetchStoreByOwnerRepo, nil, time.Second*2)

		list, total, err := u.Exec(context.TODO(), ownerID, "", 0, 0)

		is.Len(list, 0)
		is.Equal(total, int64(0))
		is.Error(err)
		fetchStoreByOwnerRepo.AssertExpectations(t)
	})

	t.Run("fail on store type repo get by id", func(t *testing.T) {
		dm := testMock()
		is := assert.New(t)
		fetchStoreByOwnerRepo := new(mocks.FetchStoreByOwnerRepository)
		getCategoryRepo := new(mocks.GetCategoryByIDRepository)

		sort := "created_at"
		page := 1
		limit := 10
		mockListStore := make([]*domain.Store, 0)
		mockListStore = append(mockListStore, dm.store)

		ownerID := uuid.NewV4().String()
		fetchStoreByOwnerRepo.
			On("Exec", mock.Anything, ownerID, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
			Return(mockListStore, int64(1), nil).
			Once()
		getCategoryRepo.
			On("Exec", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, errors.New("Unexpexted Error"))

		u := usecase.NewFetchStoreByOwner(fetchStoreByOwnerRepo, getCategoryRepo, time.Second*2)

		list, total, err := u.Exec(context.TODO(), ownerID, sort, limit, page)

		is.Len(list, 0)
		is.Equal(total, int64(0))
		is.Error(err)

		fetchStoreByOwnerRepo.AssertExpectations(t)
		getCategoryRepo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		dm := testMock()
		is := assert.New(t)
		fetchStoreByOwnerRepo := new(mocks.FetchStoreByOwnerRepository)
		getCategoryRepo := new(mocks.GetCategoryByIDRepository)

		sort := "created_at"
		page := 1
		limit := 10
		mockListStore := make([]*domain.Store, 0)
		mockListStore = append(mockListStore, dm.store)

		ownerID := uuid.NewV4().String()
		fetchStoreByOwnerRepo.
			On("Exec", mock.Anything, ownerID, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
			Return(mockListStore, int64(1), nil).
			Once()
		getCategoryRepo.
			On("Exec", mock.Anything, mock.AnythingOfType("string")).
			Return(dm.storType, nil)

		u := usecase.NewFetchStoreByOwner(fetchStoreByOwnerRepo, getCategoryRepo, time.Second*2)

		list, total, err := u.Exec(context.TODO(), ownerID, sort, limit, page)

		is.Len(list, 1)
		is.Equal(total, int64(1))
		is.NoError(err)

		fetchStoreByOwnerRepo.AssertExpectations(t)
		getCategoryRepo.AssertExpectations(t)
	})
}
