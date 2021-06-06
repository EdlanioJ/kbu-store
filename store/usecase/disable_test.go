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

func Test_Disable(t *testing.T) {
	t.Run("fail on find store", func(t *testing.T) {
		is := assert.New(t)
		getStoreByIDRepo := new(mocks.GetStoreByIDRepository)

		id := uuid.NewV4().String()
		getStoreByIDRepo.
			On("Exec", mock.Anything, id).
			Return(nil, errors.New("Unexpexted Error")).
			Once()

		u := usecase.NewDisableStore(getStoreByIDRepo, nil, time.Second*2)

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

		u := usecase.NewDisableStore(getStoreByIDRepo, nil, time.Second*2)
		err := u.Exec(context.TODO(), id)
		is.Error(err)
		is.EqualError(err, domain.ErrNotFound.Error())
	})
	t.Run("fail on status already disabled", func(t *testing.T) {
		is := assert.New(t)
		dm := testMock()

		store := dm.store
		store.Status = domain.StoreStatusInactive
		getStoreByIDRepo := new(mocks.GetStoreByIDRepository)

		id := uuid.NewV4().String()
		getStoreByIDRepo.
			On("Exec", mock.Anything, id).
			Return(store, nil).
			Once()

		u := usecase.NewDisableStore(getStoreByIDRepo, nil, time.Second*2)
		err := u.Exec(context.TODO(), id)
		is.Error(err)
		is.EqualError(err, domain.ErrInactived.Error())
	})

	t.Run("fail on status still pending", func(t *testing.T) {
		is := assert.New(t)
		dm := testMock()

		store := dm.store
		store.Status = domain.StoreStatusBlock
		getStoreByIDRepo := new(mocks.GetStoreByIDRepository)

		id := uuid.NewV4().String()
		getStoreByIDRepo.
			On("Exec", mock.Anything, id).
			Return(store, nil).
			Once()

		u := usecase.NewDisableStore(getStoreByIDRepo, nil, time.Second*2)
		err := u.Exec(context.TODO(), id)
		is.Error(err)
	})
	t.Run("fail on disable", func(t *testing.T) {
		is := assert.New(t)
		dm := testMock()
		store := dm.store
		store.Status = domain.StoreStatusActive
		store.ExternalID = "id"
		getStoreByIDRepo := new(mocks.GetStoreByIDRepository)

		id := uuid.NewV4().String()
		getStoreByIDRepo.
			On("Exec", mock.Anything, id).
			Return(store, nil).
			Once()

		u := usecase.NewDisableStore(getStoreByIDRepo, nil, time.Second*2)
		err := u.Exec(context.TODO(), id)
		is.Error(err)
	})

	t.Run("fail on update", func(t *testing.T) {
		is := assert.New(t)
		getStoreByIDRepo := new(mocks.GetStoreByIDRepository)
		updateStoreRepo := new(mocks.UpdateStoreRepository)
		dm := testMock()
		store := dm.store

		store.Status = domain.StoreStatusActive
		id := uuid.NewV4().String()
		getStoreByIDRepo.
			On("Exec", mock.Anything, id).
			Return(store, nil).
			Once()
		updateStoreRepo.
			On("Exec", mock.Anything, mock.Anything).
			Return(errors.New("Unexpexted Error")).
			Once()

		u := usecase.NewDisableStore(getStoreByIDRepo, updateStoreRepo, time.Second*2)
		err := u.Exec(context.TODO(), id)

		is.Error(err)
	})

	t.Run("success", func(t *testing.T) {
		is := assert.New(t)
		getStoreByIDRepo := new(mocks.GetStoreByIDRepository)
		updateStoreRepo := new(mocks.UpdateStoreRepository)
		dm := testMock()
		store := dm.store
		store.Status = domain.StoreStatusActive
		id := uuid.NewV4().String()
		getStoreByIDRepo.
			On("Exec", mock.Anything, id).
			Return(store, nil).
			Once()
		updateStoreRepo.
			On("Exec", mock.Anything, mock.Anything).
			Return(nil).
			Once()

		u := usecase.NewDisableStore(getStoreByIDRepo, updateStoreRepo, time.Second*2)
		err := u.Exec(context.TODO(), id)

		is.NoError(err)
	})
}
