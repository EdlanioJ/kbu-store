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

func Test_Activate(t *testing.T) {
	t.Run("fail on find store", func(t *testing.T) {
		is := assert.New(t)
		getStoreByIDRepo := new(mocks.GetStoreByIDRepository)

		id := uuid.NewV4().String()
		getStoreByIDRepo.
			On("Exec", mock.Anything, id).
			Return(nil, errors.New("Unexpexted Error")).
			Once()

		u := usecase.NewActivateStore(getStoreByIDRepo, nil, time.Second*2)

		err := u.Exec(context.TODO(), id)
		is.Error(err)
	})

	t.Run("fail on empty store", func(t *testing.T) {
		is := assert.New(t)
		getStoreByIDRepo := new(mocks.GetStoreByIDRepository)

		id := uuid.NewV4().String()
		getStoreByIDRepo.
			On("Exec", mock.Anything, id).
			Return(&domain.Store{}, nil).
			Once()

		u := usecase.NewActivateStore(getStoreByIDRepo, nil, time.Second*2)
		err := u.Exec(context.TODO(), id)
		is.Error(err)
		is.EqualError(err, domain.ErrNotFound.Error())
	})
	t.Run("fail on status already active", func(t *testing.T) {
		is := assert.New(t)
		dm := testMock()

		store := dm.store
		store.Status = domain.StoreStatusActive
		getStoreByIDRepo := new(mocks.GetStoreByIDRepository)

		id := uuid.NewV4().String()
		getStoreByIDRepo.
			On("Exec", mock.Anything, id).
			Return(store, nil).
			Once()

		u := usecase.NewActivateStore(getStoreByIDRepo, nil, time.Second*2)
		err := u.Exec(context.TODO(), id)
		is.Error(err)
		is.EqualError(err, domain.ErrActived.Error())
	})

	t.Run("fail on activate", func(t *testing.T) {
		is := assert.New(t)
		dm := testMock()
		store := dm.store
		store.ExternalID = "id"
		getStoreByIDRepo := new(mocks.GetStoreByIDRepository)

		id := uuid.NewV4().String()
		getStoreByIDRepo.
			On("Exec", mock.Anything, id).
			Return(store, nil).
			Once()

		u := usecase.NewActivateStore(getStoreByIDRepo, nil, time.Second*2)
		err := u.Exec(context.TODO(), id)
		is.Error(err)
	})

	t.Run("fail on update", func(t *testing.T) {
		is := assert.New(t)
		getStoreByIDRepo := new(mocks.GetStoreByIDRepository)
		updateStoreRepo := new(mocks.UpdateStoreRepository)
		dm := testMock()
		store := dm.store

		id := uuid.NewV4().String()
		getStoreByIDRepo.
			On("Exec", mock.Anything, id).
			Return(store, nil).
			Once()
		updateStoreRepo.
			On("Exec", mock.Anything, mock.Anything).
			Return(errors.New("Unexpexted Error")).
			Once()

		u := usecase.NewActivateStore(getStoreByIDRepo, updateStoreRepo, time.Second*2)
		err := u.Exec(context.TODO(), id)

		is.Error(err)
	})

	t.Run("success", func(t *testing.T) {
		is := assert.New(t)
		getStoreByIDRepo := new(mocks.GetStoreByIDRepository)
		updateStoreRepo := new(mocks.UpdateStoreRepository)
		dm := testMock()
		store := dm.store

		id := uuid.NewV4().String()
		getStoreByIDRepo.
			On("Exec", mock.Anything, id).
			Return(store, nil).
			Once()
		updateStoreRepo.
			On("Exec", mock.Anything, mock.Anything).
			Return(nil).
			Once()

		u := usecase.NewActivateStore(getStoreByIDRepo, updateStoreRepo, time.Second*2)
		err := u.Exec(context.TODO(), id)

		is.NoError(err)
	})
}
