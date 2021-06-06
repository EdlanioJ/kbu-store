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

func Test_DeleteUsecase(t *testing.T) {
	t.Run("fail error on store's repo get by id", func(t *testing.T) {
		is := assert.New(t)

		id := uuid.NewV4().String()
		getStoreByIDRepo := new(mocks.GetStoreByIDRepository)
		getStoreByIDRepo.On("Exec", mock.Anything, id).Return(nil, errors.New("Unexpexted Error")).Once()

		u := usecase.NewDeleteStore(nil, nil, getStoreByIDRepo, time.Second*2)

		err := u.Exec(context.TODO(), id)

		is.Error(err)
	})

	t.Run("fail error on store's repo returns empty value", func(t *testing.T) {
		is := assert.New(t)

		id := uuid.NewV4().String()
		getStoreByIDRepo := new(mocks.GetStoreByIDRepository)
		getStoreByIDRepo.On("Exec", mock.Anything, id).Return(&domain.Store{}, nil).Once()

		u := usecase.NewDeleteStore(nil, nil, getStoreByIDRepo, time.Second*2)

		err := u.Exec(context.TODO(), id)

		is.Error(err)
		is.EqualError(err, domain.ErrNotFound.Error())
	})

	t.Run("fail on delete account", func(t *testing.T) {
		is := assert.New(t)
		dm := testMock()
		id := uuid.NewV4().String()
		getStoreByIDRepo := new(mocks.GetStoreByIDRepository)
		deleteAccountRepo := new(mocks.DeleteAccountRepository)
		getStoreByIDRepo.On("Exec", mock.Anything, id).Return(dm.store, nil).Once()
		deleteAccountRepo.On("Exec", mock.Anything, dm.store.Account.ID).Return(errors.New("Unexpexted Error"))
		u := usecase.NewDeleteStore(nil, deleteAccountRepo, getStoreByIDRepo, time.Second*2)

		err := u.Exec(context.TODO(), id)

		is.Error(err)
	})

	t.Run("fail on delete store", func(t *testing.T) {
		is := assert.New(t)
		dm := testMock()
		id := uuid.NewV4().String()
		getStoreByIDRepo := new(mocks.GetStoreByIDRepository)
		deleteStoreRepo := new(mocks.DeleteStoreRepository)
		deleteAccountRepo := new(mocks.DeleteAccountRepository)
		getStoreByIDRepo.On("Exec", mock.Anything, id).Return(dm.store, nil).Once()
		deleteAccountRepo.On("Exec", mock.Anything, dm.store.Account.ID).Return(nil)
		deleteStoreRepo.On("Exec", mock.Anything, id).Return(errors.New("Unexpexted Error"))
		u := usecase.NewDeleteStore(deleteStoreRepo, deleteAccountRepo, getStoreByIDRepo, time.Second*2)

		err := u.Exec(context.TODO(), id)

		is.Error(err)
	})

	t.Run("success", func(t *testing.T) {
		is := assert.New(t)
		dm := testMock()
		id := uuid.NewV4().String()
		getStoreByIDRepo := new(mocks.GetStoreByIDRepository)
		deleteStoreRepo := new(mocks.DeleteStoreRepository)

		deleteAccountRepo := new(mocks.DeleteAccountRepository)
		getStoreByIDRepo.On("Exec", mock.Anything, id).Return(dm.store, nil).Once()
		deleteAccountRepo.On("Exec", mock.Anything, dm.store.Account.ID).Return(nil)
		deleteStoreRepo.On("Exec", mock.Anything, id).Return(nil)
		u := usecase.NewDeleteStore(deleteStoreRepo, deleteAccountRepo, getStoreByIDRepo, time.Second*2)

		err := u.Exec(context.TODO(), id)

		is.NoError(err)
	})
}
