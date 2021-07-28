package usecases_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/EdlanioJ/kbu-store/app/domain"
	"github.com/EdlanioJ/kbu-store/app/usecases"
	"github.com/EdlanioJ/kbu-store/app/utils/mocks"
	"github.com/EdlanioJ/kbu-store/app/utils/sample"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_StoreUsecase_Create(t *testing.T) {
	args := sample.NewStoreUsecaseRequest()
	validCategory := sample.NewCategory()
	validCategory.Status = domain.CategoryStatusActive
	invalidCategory := sample.NewCategory()
	invalidArgs := sample.NewStoreUsecaseRequest()
	invalidArgs.Name = ""
	type fields struct {
		storeRepo    *mocks.StoreRepository
		accountRepo  *mocks.AccountRepository
		categoryRepo *mocks.CategoryRepository
		msgProducer  *mocks.MessengerProducer
	}
	testCases := []struct {
		name        string
		args        sample.StoreUsecaseCreateRequest
		expectedErr bool
		prepare     func(f fields)
	}{
		{
			name:        "failure_find_category_by_id_returns_error",
			args:        args,
			expectedErr: true,
			prepare: func(f fields) {
				f.categoryRepo.On("FindByID", mock.Anything, args.CategoryID).Return(nil, errors.New("Unexpexted Error")).Once()
			},
		},
		{
			name:        "failure_find_category_by_id_returns_invalid_category",
			args:        args,
			expectedErr: true,
			prepare: func(f fields) {
				f.categoryRepo.On("FindByID", mock.Anything, args.CategoryID).Return(invalidCategory, nil).Once()
			},
		},
		{
			name:        "failure_create_account_returns_error",
			args:        args,
			expectedErr: true,
			prepare: func(f fields) {
				f.categoryRepo.On("FindByID", mock.Anything, args.CategoryID).Return(validCategory, nil).Once()
				f.accountRepo.On("Store", mock.Anything, mock.Anything).Return(errors.New("Unexpexted Error")).Once()
			},
		},
		{
			name:        "failure_create_store_returns_error",
			args:        args,
			expectedErr: true,
			prepare: func(f fields) {
				f.categoryRepo.On("FindByID", mock.Anything, args.CategoryID).Return(validCategory, nil).Once()
				f.accountRepo.On("Store", mock.Anything, mock.Anything).Return(nil).Once()
				f.storeRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("Unexpexted Error")).Once()
			},
		},
		{
			name:        "failure_produce_returns_error",
			args:        args,
			expectedErr: true,
			prepare: func(f fields) {

				f.categoryRepo.On("FindByID", mock.Anything, args.CategoryID).Return(validCategory, nil).Once()
				f.accountRepo.On("Store", mock.Anything, mock.Anything).Return(nil).Once()
				f.storeRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
				f.msgProducer.On("Publish", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(errors.New("Unexpexted Error"))
			},
		},
		{
			name:        "success",
			args:        args,
			expectedErr: false,
			prepare: func(f fields) {
				f.categoryRepo.On("FindByID", mock.Anything, args.CategoryID).Return(validCategory, nil).Once()
				f.accountRepo.On("Store", mock.Anything, mock.Anything).Return(nil).Once()
				f.storeRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
				f.msgProducer.On("Publish", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			accountRepo := new(mocks.AccountRepository)
			categoryRepo := new(mocks.CategoryRepository)
			storeRepo := new(mocks.StoreRepository)
			msgProducer := new(mocks.MessengerProducer)
			f := fields{storeRepo, accountRepo, categoryRepo, msgProducer}
			tc.prepare(f)
			u := usecases.NewStoreUsecase(storeRepo, accountRepo, categoryRepo, msgProducer, time.Second*2)

			err := u.Store(context.TODO(), tc.args.Name, tc.args.Description, tc.args.CategoryID, tc.args.UserID, tc.args.Tags, 0, 0)
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			accountRepo.AssertExpectations(t)
			categoryRepo.AssertExpectations(t)
			storeRepo.AssertExpectations(t)
		})
	}
}

func Test_StoreUsecase_Get(t *testing.T) {
	testCases := []struct {
		name        string
		arg         string
		expectedErr bool
		prepare     func(storeRepo *mocks.StoreRepository)
	}{
		{
			name:        "Failure_find_store_by_id_returns_error",
			arg:         uuid.NewV4().String(),
			expectedErr: true,
			prepare: func(storeRepo *mocks.StoreRepository) {
				storeRepo.On("FindByID", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("Unexpexted Error")).Once()
			},
		},
		{
			name: "success",
			arg:  uuid.NewV4().String(),
			prepare: func(storeRepo *mocks.StoreRepository) {
				s := sample.NewStore()
				storeRepo.On("FindByID", mock.Anything, mock.AnythingOfType("string")).Return(s, nil).Once()
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			storeRepo := new(mocks.StoreRepository)
			tc.prepare(storeRepo)
			u := usecases.NewStoreUsecase(storeRepo, nil, nil, nil, time.Second*2)
			res, err := u.Get(context.TODO(), tc.arg)

			if tc.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			}
			storeRepo.AssertExpectations(t)
		})
	}
}

func Test_StoreUsecase_Index(t *testing.T) {
	type args struct {
		page  int
		limit int
		sort  string
	}
	s := getStore()

	a := args{
		page:  1,
		limit: 10,
		sort:  "created_at",
	}
	testCases := []struct {
		name          string
		args          args
		builtSts      func(storeRepo *mocks.StoreRepository)
		checkResponse func(t *testing.T, res domain.Stores, count int64, err error)
	}{
		{
			name: "should fail if find all stores returns error",
			args: args{
				page:  0,
				limit: 0,
				sort:  "",
			},
			builtSts: func(storeRepo *mocks.StoreRepository) {
				storeRepo.On("FindAll", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
					Return(nil, int64(0), errors.New("Unexpexted Error")).
					Once()
			},
			checkResponse: func(t *testing.T, res domain.Stores, count int64, err error) {
				assert.Len(t, res, 0)
				assert.Equal(t, count, int64(0))
				assert.Error(t, err)
			},
		},
		{
			name: "should succeed",
			args: a,
			builtSts: func(storeRepo *mocks.StoreRepository) {
				stores := make(domain.Stores, 0)
				stores = append(stores, s)
				storeRepo.On("FindAll", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
					Return(stores, int64(1), nil).
					Once()
			},
			checkResponse: func(t *testing.T, res domain.Stores, count int64, err error) {
				assert.Len(t, res, 1)
				assert.Equal(t, count, int64(1))
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storeRepo := new(mocks.StoreRepository)

			tc.builtSts(storeRepo)
			u := usecases.NewStoreUsecase(storeRepo, nil, nil, nil, time.Second*2)

			res, count, err := u.Index(context.TODO(), tc.args.sort, tc.args.limit, tc.args.page)
			tc.checkResponse(t, res, count, err)

			storeRepo.AssertExpectations(t)
		})
	}
}

func Test_StoreUsecase_Block(t *testing.T) {
	type fields struct {
		storeRepo   *mocks.StoreRepository
		msgProducer *mocks.MessengerProducer
	}
	testCases := []struct {
		name        string
		arg         string
		expectedErr bool
		prepare     func(f fields)
	}{
		{
			name:        "failure_find_store_by_id_returns_error",
			arg:         uuid.NewV4().String(),
			expectedErr: true,
			prepare: func(f fields) {
				f.storeRepo.On("FindByID", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("Unexpexted Error")).Once()
			},
		},
		{
			name:        "failure_block_method_returns_error",
			arg:         uuid.NewV4().String(),
			expectedErr: true,
			prepare: func(f fields) {
				store := sample.NewStore()
				store.Status = domain.StoreStatusBlock
				f.storeRepo.
					On("FindByID", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
			},
		},
		{
			name:        "failure_update_store_returns_error",
			arg:         uuid.NewV4().String(),
			expectedErr: true,
			prepare: func(f fields) {
				store := sample.NewStore()
				store.Status = domain.StoreStatusActive
				f.storeRepo.
					On("FindByID", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
				f.storeRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("Unexpexted Error")).Once()
			},
		},
		{
			name:        "failure_publish_msg_returns_error",
			arg:         uuid.NewV4().String(),
			expectedErr: true,
			prepare: func(f fields) {
				store := sample.NewStore()
				store.Status = domain.StoreStatusActive
				f.storeRepo.
					On("FindByID", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
				f.storeRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
				f.msgProducer.On("Publish", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(errors.New("Unexpected Error"))
			},
		},
		{
			name: "success",
			arg:  uuid.NewV4().String(),
			prepare: func(f fields) {
				store := sample.NewStore()
				store.Status = domain.StoreStatusActive
				f.storeRepo.
					On("FindByID", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
				f.storeRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
				f.msgProducer.On("Publish", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			storeRepo := new(mocks.StoreRepository)
			msgProducer := new(mocks.MessengerProducer)
			f := fields{storeRepo, msgProducer}
			tc.prepare(f)
			u := usecases.NewStoreUsecase(storeRepo, nil, nil, msgProducer, time.Second*2)
			err := u.Block(context.TODO(), tc.arg)
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_StoreUsecase_Active(t *testing.T) {
	type fields struct {
		storeRepo   *mocks.StoreRepository
		msgProducer *mocks.MessengerProducer
	}
	testCases := []struct {
		name        string
		arg         string
		expectedErr bool
		prepare     func(f fields)
	}{
		{
			name:        "failure_find_store_by_id_returns_error",
			arg:         uuid.NewV4().String(),
			expectedErr: true,
			prepare: func(f fields) {
				f.storeRepo.On("FindByID", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("Unexpected Error")).Once()
			},
		},
		{
			name:        "failure_activate_method_returns_error",
			arg:         uuid.NewV4().String(),
			expectedErr: true,
			prepare: func(f fields) {
				store := sample.NewStore()
				store.Status = domain.StoreStatusActive
				f.storeRepo.
					On("FindByID", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
			},
		},
		{
			name:        "failure_update_store_returns_error",
			arg:         uuid.NewV4().String(),
			expectedErr: true,
			prepare: func(f fields) {
				store := sample.NewStore()
				store.Status = domain.StoreStatusPending
				f.storeRepo.
					On("FindByID", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
				f.storeRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("Unexpected Error")).Once()
			},
		},
		{
			name:        "failure_publish_msg_returns_error",
			arg:         uuid.NewV4().String(),
			expectedErr: true,
			prepare: func(f fields) {
				store := sample.NewStore()
				store.Status = domain.StoreStatusPending
				f.storeRepo.
					On("FindByID", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
				f.storeRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
				f.msgProducer.On("Publish", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(errors.New("Unexpected Error"))
			},
		},
		{
			name: "success",
			arg:  uuid.NewV4().String(),
			prepare: func(f fields) {
				store := sample.NewStore()
				store.Status = domain.StoreStatusPending
				f.storeRepo.
					On("FindByID", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
				f.storeRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
				f.msgProducer.On("Publish", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			storeRepo := new(mocks.StoreRepository)
			msgProducer := new(mocks.MessengerProducer)
			f := fields{storeRepo, msgProducer}
			tc.prepare(f)
			u := usecases.NewStoreUsecase(storeRepo, nil, nil, msgProducer, time.Second*2)
			err := u.Active(context.TODO(), tc.arg)
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			storeRepo.AssertExpectations(t)
			msgProducer.AssertExpectations(t)
		})
	}
}

func Test_StoreUsecase_Disable(t *testing.T) {
	type fields struct {
		storeRepo   *mocks.StoreRepository
		msgProducer *mocks.MessengerProducer
	}
	testCases := []struct {
		name        string
		arg         string
		expectedErr bool
		prepare     func(f fields)
	}{
		{
			name:        "failure_find_store_by_id_returns_error",
			arg:         uuid.NewV4().String(),
			expectedErr: true,
			prepare: func(f fields) {
				f.storeRepo.On("FindByID", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("Unexpexted Error")).Once()

			},
		},
		{
			name:        "failure_disable_method_returns_error",
			arg:         uuid.NewV4().String(),
			expectedErr: true,
			prepare: func(f fields) {
				store := sample.NewStore()
				store.Status = domain.StoreStatusDisable
				f.storeRepo.
					On("FindByID", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
			},
		},
		{
			name:        "failure_update_store_returns_error",
			arg:         uuid.NewV4().String(),
			expectedErr: true,
			prepare: func(f fields) {
				store := sample.NewStore()
				store.Status = domain.StoreStatusActive
				f.storeRepo.
					On("FindByID", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
				f.storeRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("Unexpexted Error")).Once()
			},
		},
		{
			name:        "failure_publish_msg_returns_error",
			arg:         uuid.NewV4().String(),
			expectedErr: true,
			prepare: func(f fields) {
				store := sample.NewStore()
				store.Status = domain.StoreStatusActive
				f.storeRepo.
					On("FindByID", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
				f.storeRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
				f.msgProducer.On("Publish", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(errors.New("Unexpected Error"))
			},
		},
		{
			name: "success",
			arg:  uuid.NewV4().String(),
			prepare: func(f fields) {
				store := sample.NewStore()
				store.Status = domain.StoreStatusActive
				f.storeRepo.
					On("FindByID", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
				f.storeRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
				f.msgProducer.On("Publish", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			storeRepo := new(mocks.StoreRepository)
			msgProducer := new(mocks.MessengerProducer)
			f := fields{storeRepo, msgProducer}
			tc.prepare(f)
			u := usecases.NewStoreUsecase(storeRepo, nil, nil, msgProducer, time.Second*2)
			err := u.Disable(context.TODO(), tc.arg)
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			storeRepo.AssertExpectations(t)
			msgProducer.AssertExpectations(t)
		})
	}
}

func Test_StoreUsecase_Update(t *testing.T) {

	type fields struct {
		storeRepo   *mocks.StoreRepository
		msgProducer *mocks.MessengerProducer
	}
	testCases := []struct {
		name        string
		arg         *domain.Store
		expectedErr bool
		prepare     func(f fields)
	}{
		{
			name:        "failure_find_store_by_id_returns_error",
			arg:         sample.NewStore(),
			expectedErr: true,
			prepare: func(f fields) {
				f.storeRepo.On("FindByID", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("Unexpexted Error")).Once()
			},
		},
		{
			name:        "failure_update_store_returns_error",
			arg:         sample.NewStore(),
			expectedErr: true,
			prepare: func(f fields) {
				foundStore := sample.NewStore()
				f.storeRepo.On("FindByID", mock.Anything, mock.AnythingOfType("string")).Return(foundStore, nil).Once()
				f.storeRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("Unexpected Error")).Once()
			},
		},
		{
			name:        "failure_publish_msg_returns_error",
			arg:         sample.NewStore(),
			expectedErr: true,
			prepare: func(f fields) {
				foundStore := sample.NewStore()
				f.storeRepo.On("FindByID", mock.Anything, mock.AnythingOfType("string")).Return(foundStore, nil).Once()
				f.storeRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
				f.msgProducer.On("Publish", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(errors.New("Unexpected Error"))
			},
		},
		{
			name: "success",
			arg:  sample.NewStore(),
			prepare: func(f fields) {
				foundStore := sample.NewStore()
				f.storeRepo.On("FindByID", mock.Anything, mock.AnythingOfType("string")).Return(foundStore, nil).Once()
				f.storeRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
				f.msgProducer.On("Publish", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			storeRepo := new(mocks.StoreRepository)
			msgProducer := new(mocks.MessengerProducer)
			f := fields{storeRepo, msgProducer}
			tc.prepare(f)
			u := usecases.NewStoreUsecase(storeRepo, nil, nil, msgProducer, time.Second*2)
			err := u.Update(context.TODO(), tc.arg)
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			storeRepo.AssertExpectations(t)
			msgProducer.AssertExpectations(t)
		})
	}
}

func Test_StoreUsecase_Delete(t *testing.T) {
	type fields struct {
		storeRepo   *mocks.StoreRepository
		accountRepo *mocks.AccountRepository
		msgProducer *mocks.MessengerProducer
	}

	testCases := []struct {
		name        string
		arg         string
		expectedErr bool
		prepare     func(f fields)
	}{
		{
			name:        "failure_find_store_by_id_returns_error",
			arg:         uuid.NewV4().String(),
			expectedErr: true,
			prepare: func(f fields) {
				f.storeRepo.On("FindByID", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("Unexpexted Error")).Once()
			},
		},
		{
			name:        "failure_delete_store_returns_error",
			arg:         uuid.NewV4().String(),
			expectedErr: true,
			prepare: func(f fields) {
				store := sample.NewStore()
				store.Status = domain.StoreStatusActive
				f.storeRepo.
					On("FindByID", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
				f.storeRepo.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(errors.New("Unexpexted Error")).Once()
			},
		},
		{
			name:        "failure_delete_account_returns_error",
			arg:         uuid.NewV4().String(),
			expectedErr: true,
			prepare: func(f fields) {
				store := sample.NewStore()
				store.Status = domain.StoreStatusActive
				f.storeRepo.
					On("FindByID", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
				f.storeRepo.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(nil).Once()
				f.accountRepo.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(errors.New("Unexpexted Error")).Once()
			},
		},
		{
			name:        "failure_publish_msg_returns_error",
			arg:         uuid.NewV4().String(),
			expectedErr: true,
			prepare: func(f fields) {
				store := sample.NewStore()
				store.Status = domain.StoreStatusActive
				f.storeRepo.
					On("FindByID", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
				f.storeRepo.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(nil).Once()
				f.accountRepo.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(nil).Once()
				f.msgProducer.On("Publish", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(errors.New("Unexpected Error"))
			},
		},
		{
			name: "success",
			arg:  uuid.NewV4().String(),
			prepare: func(f fields) {
				store := sample.NewStore()
				store.Status = domain.StoreStatusActive
				f.storeRepo.
					On("FindByID", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
				f.storeRepo.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(nil).Once()
				f.accountRepo.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(nil).Once()
				f.msgProducer.On("Publish", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			storeRepo := new(mocks.StoreRepository)
			accountRepo := new(mocks.AccountRepository)
			msgProducer := new(mocks.MessengerProducer)
			f := fields{storeRepo, accountRepo, msgProducer}
			tc.prepare(f)
			u := usecases.NewStoreUsecase(storeRepo, accountRepo, nil, msgProducer, time.Second*2)
			err := u.Delete(context.TODO(), tc.arg)
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			storeRepo.AssertExpectations(t)
			accountRepo.AssertExpectations(t)
			msgProducer.AssertExpectations(t)
		})
	}
}
