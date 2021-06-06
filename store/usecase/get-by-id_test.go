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

func Test_GetStoreByIDUsecase(t *testing.T) {
	t.Run("fail on get store repo", func(t *testing.T) {
		is := assert.New(t)
		getStoreByIDRepo := new(mocks.GetStoreByIDRepository)

		id := uuid.NewV4().String()
		getStoreByIDRepo.
			On("Exec", mock.Anything, id).
			Return(nil, errors.New("Unexpexted Error")).
			Once()
		u := usecase.NewGetStoreByID(getStoreByIDRepo, nil, time.Second*2)

		res, err := u.Exec(context.TODO(), id)
		is.Error(err)
		is.Nil(res)

		getStoreByIDRepo.AssertExpectations(t)
	})

	t.Run("fail on get store type's repo", func(t *testing.T) {
		is := assert.New(t)
		dm := testMock()
		getStoreByIDRepo := new(mocks.GetStoreByIDRepository)
		getCategoryRepo := new(mocks.GetCategoryByIDUsecase)
		id := uuid.NewV4().String()
		getStoreByIDRepo.
			On("Exec", mock.Anything, id).
			Return(dm.store, nil).
			Once()
		getCategoryRepo.
			On("Exec", mock.Anything, mock.AnythingOfType(("string"))).
			Return(nil, errors.New("Unexpexted Error")).
			Once()
		u := usecase.NewGetStoreByID(getStoreByIDRepo, getCategoryRepo, time.Second*2)

		res, err := u.Exec(context.TODO(), id)
		is.Error(err)
		is.Nil(res)

		getStoreByIDRepo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		is := assert.New(t)
		dm := testMock()
		getStoreByIDRepo := new(mocks.GetStoreByIDRepository)
		getCategoryRepo := new(mocks.GetCategoryByIDUsecase)

		getStoreByIDRepo.
			On("Exec", mock.Anything, dm.store.ID).
			Return(dm.store, nil).
			Once()
		getCategoryRepo.
			On("Exec", mock.Anything, mock.AnythingOfType(("string"))).
			Return(dm.storType, nil).
			Once()
		u := usecase.NewGetStoreByID(getStoreByIDRepo, getCategoryRepo, time.Second*2)

		res, err := u.Exec(context.TODO(), dm.store.ID)
		is.NoError(err)
		is.NotNil(res)
	})
}
