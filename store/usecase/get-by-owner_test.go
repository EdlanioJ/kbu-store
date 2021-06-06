package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/EdlanioJ/kbu-store/domain/mocks"
	"github.com/EdlanioJ/kbu-store/store/usecase"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_GetStoreByOwnerUsecase(t *testing.T) {
	t.Run("fail on get store repo", func(t *testing.T) {
		is := assert.New(t)
		getStoreByOwnerRepo := new(mocks.GetStoreByOwnerRepository)

		id := uuid.NewV4().String()
		ownerID := uuid.NewV4().String()
		getStoreByOwnerRepo.
			On("Exec", mock.Anything, id, ownerID).
			Return(nil, errors.New("Unexpexted Error")).
			Once()
		u := usecase.NewGetStoreByOwner(getStoreByOwnerRepo, nil, time.Second*2)

		res, err := u.Exec(context.TODO(), id, ownerID)
		is.Error(err)
		is.Nil(res)

		getStoreByOwnerRepo.AssertExpectations(t)
	})

	t.Run("fail on get store type's repo", func(t *testing.T) {
		is := assert.New(t)
		dm := testMock()
		getStoreByOwnerRepo := new(mocks.GetStoreByOwnerRepository)
		getCategoryRepo := new(mocks.GetCategoryByIDUsecase)
		id := uuid.NewV4().String()
		ownerID := uuid.NewV4().String()
		getStoreByOwnerRepo.
			On("Exec", mock.Anything, id, ownerID).
			Return(dm.store, nil).
			Once()
		getCategoryRepo.
			On("Exec", mock.Anything, mock.AnythingOfType(("string"))).
			Return(nil, errors.New("Unexpexted Error")).
			Once()
		u := usecase.NewGetStoreByOwner(getStoreByOwnerRepo, getCategoryRepo, time.Second*2)

		res, err := u.Exec(context.TODO(), id, ownerID)
		is.Error(err)
		is.Nil(res)

		getStoreByOwnerRepo.AssertExpectations(t)
		getCategoryRepo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		is := assert.New(t)
		dm := testMock()
		getStoreByOwnerRepo := new(mocks.GetStoreByOwnerRepository)
		getCategoryRepo := new(mocks.GetCategoryByIDUsecase)

		getStoreByOwnerRepo.
			On("Exec", mock.Anything, dm.store.ID, dm.store.ExternalID).
			Return(dm.store, nil).
			Once()
		getCategoryRepo.
			On("Exec", mock.Anything, mock.AnythingOfType(("string"))).
			Return(dm.storType, nil).
			Once()
		u := usecase.NewGetStoreByOwner(getStoreByOwnerRepo, getCategoryRepo, time.Second*2)

		res, err := u.Exec(context.TODO(), dm.store.ID, dm.store.ExternalID)
		is.NoError(err)
		is.NotNil(res)

		getStoreByOwnerRepo.AssertExpectations(t)
		getCategoryRepo.AssertExpectations(t)
	})
}
