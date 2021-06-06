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

func Test_UpdateUsecase(t *testing.T) {
	t.Run("fail error on store's repo get by id", func(t *testing.T) {
		is := assert.New(t)
		dm := testMock()
		store := dm.store

		getStoreByIDRepo := new(mocks.GetStoreByIDRepository)
		getStoreByIDRepo.On("Exec", mock.Anything, store.ID).Return(nil, errors.New("Unexpexted Error")).Once()

		u := usecase.NewUpdateStore(getStoreByIDRepo, nil, time.Second*2)
		err := u.Exec(context.TODO(), store)
		is.Error(err)
		getStoreByIDRepo.AssertExpectations(t)
	})

	t.Run("fail error on store's repo get by id returns empty struct", func(t *testing.T) {
		is := assert.New(t)
		dm := testMock()
		store := dm.store

		getStoreByIDRepo := new(mocks.GetStoreByIDRepository)
		getStoreByIDRepo.On("Exec", mock.Anything, store.ID).Return(&domain.Store{}, nil).Once()

		u := usecase.NewUpdateStore(getStoreByIDRepo, nil, time.Second*2)
		err := u.Exec(context.TODO(), store)
		is.Error(err)
		is.EqualError(err, domain.ErrNotFound.Error())
		getStoreByIDRepo.AssertExpectations(t)
	})
	t.Run("fail error on store's repo update", func(t *testing.T) {
		is := assert.New(t)
		dm := testMock()
		store := dm.store

		getStoreByIDRepo := new(mocks.GetStoreByIDRepository)
		updateStoreRepo := new(mocks.UpdateStoreRepository)

		getStoreByIDRepo.On("Exec", mock.Anything, store.ID).Return(store, nil).Once()
		updateStoreRepo.On("Exec", mock.Anything, mock.Anything).Return(errors.New("Unexpexted Error"))

		u := usecase.NewUpdateStore(getStoreByIDRepo, updateStoreRepo, time.Second*2)
		err := u.Exec(context.TODO(), store)
		is.Error(err)
		getStoreByIDRepo.AssertExpectations(t)
		updateStoreRepo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		is := assert.New(t)
		dm := testMock()
		store := dm.store

		getStoreByIDRepo := new(mocks.GetStoreByIDRepository)
		updateStoreRepo := new(mocks.UpdateStoreRepository)

		getStoreByIDRepo.On("Exec", mock.Anything, store.ID).Return(store, nil).Once()
		updateStoreRepo.On("Exec", mock.Anything, mock.Anything).Return(nil).Once()

		u := usecase.NewUpdateStore(getStoreByIDRepo, updateStoreRepo, time.Second*2)
		err := u.Exec(context.TODO(), store)
		is.NoError(err)
		getStoreByIDRepo.AssertExpectations(t)
		updateStoreRepo.AssertExpectations(t)
	})
}
