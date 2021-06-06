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

func Test_CreateUsecase(t *testing.T) {
	t.Run("fail on create account", func(t *testing.T) {
		is := assert.New(t)
		createAccountRepo := new(mocks.CreateAccountRepository)

		name := "store 001"
		description := "store description 001"
		typeID := uuid.NewV4().String()
		externalID := uuid.NewV4().String()
		tags := []string{"tag001", "tag002"}

		createAccountRepo.On("Add", mock.Anything, mock.Anything).Return(errors.New("Unexpexted Error"))
		u := usecase.NewCreateStore(nil, createAccountRepo, nil, time.Second*2)

		err := u.Add(context.TODO(), name, description, typeID, externalID, tags, 0, 0)
		is.Error(err)
		createAccountRepo.AssertExpectations(t)
	})

	t.Run("fail on store type's repo get by id", func(t *testing.T) {
		is := assert.New(t)
		createAccountRepo := new(mocks.CreateAccountRepository)
		getCategoryByIDRepo := new(mocks.GetCategoryByIDRepository)

		name := "store 001"
		description := "store description 001"
		typeID := uuid.NewV4().String()
		externalID := uuid.NewV4().String()
		tags := []string{"tag001", "tag002"}

		createAccountRepo.On("Add", mock.Anything, mock.Anything).Return(nil)
		getCategoryByIDRepo.On("Exec", mock.Anything, typeID).Return(nil, errors.New("Unexpexted Error"))
		u := usecase.NewCreateStore(nil, createAccountRepo, getCategoryByIDRepo, time.Second*2)

		err := u.Add(context.TODO(), name, description, typeID, externalID, tags, 0, 0)
		is.Error(err)
		createAccountRepo.AssertExpectations(t)
		getCategoryByIDRepo.AssertExpectations(t)
	})

	t.Run("fail on new store", func(t *testing.T) {
		is := assert.New(t)
		dm := testMock()
		createAccountRepo := new(mocks.CreateAccountRepository)
		getCategoryByIDRepo := new(mocks.GetCategoryByIDRepository)

		name := "store 001"
		description := "store description 001"
		typeID := uuid.NewV4().String()
		externalID := "invalid_id"
		tags := []string{"tag001", "tag002"}

		createAccountRepo.On("Add", mock.Anything, mock.Anything).Return(nil)
		getCategoryByIDRepo.On("Exec", mock.Anything, typeID).Return(dm.storType, nil)
		u := usecase.NewCreateStore(nil, createAccountRepo, getCategoryByIDRepo, time.Second*2)
		err := u.Add(context.TODO(), name, description, typeID, externalID, tags, 0, 0)

		is.Error(err)
		createAccountRepo.AssertExpectations(t)
		getCategoryByIDRepo.AssertExpectations(t)
	})

	t.Run("fail on store's repo create", func(t *testing.T) {
		is := assert.New(t)
		dm := testMock()
		createStoreRepo := new(mocks.CreateStoreRepository)
		createAccountRepo := new(mocks.CreateAccountRepository)
		getCategoryByIDRepo := new(mocks.GetCategoryByIDRepository)

		name := "store 001"
		description := "store description 001"
		typeID := uuid.NewV4().String()
		externalID := uuid.NewV4().String()
		tags := []string{"tag001", "tag002"}

		createAccountRepo.On("Add", mock.Anything, mock.Anything).Return(nil)
		getCategoryByIDRepo.On("Exec", mock.Anything, typeID).Return(dm.storType, nil)
		createStoreRepo.On("Add", mock.Anything, mock.Anything).Return(errors.New("Unexpexted Error"))
		u := usecase.NewCreateStore(createStoreRepo, createAccountRepo, getCategoryByIDRepo, time.Second*2)
		err := u.Add(context.TODO(), name, description, typeID, externalID, tags, 0, 0)

		is.Error(err)
		createAccountRepo.AssertExpectations(t)
		getCategoryByIDRepo.AssertExpectations(t)
		createStoreRepo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		is := assert.New(t)
		dm := testMock()
		createStoreRepo := new(mocks.CreateStoreRepository)
		createAccountRepo := new(mocks.CreateAccountRepository)
		getCategoryByIDRepo := new(mocks.GetCategoryByIDRepository)

		name := "store 001"
		description := "store description 001"
		typeID := uuid.NewV4().String()
		externalID := uuid.NewV4().String()
		tags := []string{"tag001", "tag002"}

		createAccountRepo.On("Add", mock.Anything, mock.Anything).Return(nil)
		getCategoryByIDRepo.On("Exec", mock.Anything, typeID).Return(dm.storType, nil)
		createStoreRepo.On("Add", mock.Anything, mock.Anything).Return(nil)
		u := usecase.NewCreateStore(createStoreRepo, createAccountRepo, getCategoryByIDRepo, time.Second*2)
		err := u.Add(context.TODO(), name, description, typeID, externalID, tags, 0, 0)

		is.NoError(err)
		createAccountRepo.AssertExpectations(t)
		getCategoryByIDRepo.AssertExpectations(t)
		createStoreRepo.AssertExpectations(t)
	})
}
