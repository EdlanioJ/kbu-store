package usecase_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/domain/mocks"
	"github.com/EdlanioJ/kbu-store/usecase"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_StoreUsecase_Create(t *testing.T) {
	type args struct {
		name        string
		description string
		categoryID  string
		externalID  string
		tags        []string
	}

	s := getStore()
	c := getCategory()

	a := args{
		name:        s.Name,
		description: s.Description,
		categoryID:  s.Category.ID,
		externalID:  s.UserID,
		tags:        s.Tags,
	}
	testCases := []struct {
		name          string
		args          args
		builtSts      func(storeRepo *mocks.StoreRepository, accountRepo *mocks.AccountRepository, categoryRepo *mocks.CategoryRepository)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "should fail if create account repo returns error",
			args: a,
			builtSts: func(_ *mocks.StoreRepository, accountRepo *mocks.AccountRepository, _ *mocks.CategoryRepository) {
				accountRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "should fail if get category by id and status returns error",
			args: a,
			builtSts: func(_ *mocks.StoreRepository, accountRepo *mocks.AccountRepository, categoryRepo *mocks.CategoryRepository) {
				accountRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
				categoryRepo.On("GetByIdAndStatus", mock.Anything, s.Category.ID, domain.CategoryStatusActive).Return(nil, errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "should fail if new store returns error",
			args: args{
				name:        "",
				description: s.Description,
				categoryID:  s.Category.ID,
				externalID:  s.UserID,
				tags:        s.Tags,
			},
			builtSts: func(_ *mocks.StoreRepository, accountRepo *mocks.AccountRepository, categoryRepo *mocks.CategoryRepository) {
				accountRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
				categoryRepo.On("GetByIdAndStatus", mock.Anything, s.Category.ID, domain.CategoryStatusActive).Return(c, nil).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "should fail if create store repo returns error",
			args: a,
			builtSts: func(storeRepo *mocks.StoreRepository, accountRepo *mocks.AccountRepository, categoryRepo *mocks.CategoryRepository) {
				accountRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
				categoryRepo.On("GetByIdAndStatus", mock.Anything, s.Category.ID, domain.CategoryStatusActive).Return(c, nil).Once()
				storeRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "should succeed",
			args: a,
			builtSts: func(storeRepo *mocks.StoreRepository, accountRepo *mocks.AccountRepository, categoryRepo *mocks.CategoryRepository) {
				accountRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
				categoryRepo.On("GetByIdAndStatus", mock.Anything, s.Category.ID, domain.CategoryStatusActive).Return(c, nil).Once()
				storeRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, ts := range testCases {
		t.Run(ts.name, func(t *testing.T) {
			accountRepo := new(mocks.AccountRepository)
			categoryRepo := new(mocks.CategoryRepository)
			storeRepo := new(mocks.StoreRepository)

			ts.builtSts(storeRepo, accountRepo, categoryRepo)
			u := usecase.NewStoreUsecase(storeRepo, accountRepo, categoryRepo, time.Second*2)

			err := u.Store(context.TODO(), ts.args.name, ts.args.description, ts.args.categoryID, ts.args.externalID, ts.args.tags, 0, 0)
			ts.checkResponse(t, err)
			accountRepo.AssertExpectations(t)
			categoryRepo.AssertExpectations(t)
			storeRepo.AssertExpectations(t)
		})
	}
}

func Test_StoreUsecase_Get(t *testing.T) {
	s := getStore()
	c := getCategory()

	testCases := []struct {
		name          string
		arg           string
		builtSts      func(storeRepo *mocks.StoreRepository, categoryRepo *mocks.CategoryRepository)
		checkResponse func(t *testing.T, res *domain.Store, err error)
	}{
		{
			name: "should fail if find store by id returns error",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository, _ *mocks.CategoryRepository) {
				storeRepo.On("FindByID", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, res *domain.Store, err error) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "should fail if get category by id returns error",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository, categoryRepo *mocks.CategoryRepository) {
				storeRepo.On("FindByID", mock.Anything, mock.AnythingOfType("string")).Return(s, nil).Once()
				categoryRepo.On("GetById", mock.Anything, s.Category.ID).Return(nil, errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, res *domain.Store, err error) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "should succeed",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository, categoryRepo *mocks.CategoryRepository) {
				storeRepo.On("FindByID", mock.Anything, mock.AnythingOfType("string")).Return(s, nil).Once()
				categoryRepo.On("GetById", mock.Anything, s.Category.ID).Return(c, nil).Once()
			},
			checkResponse: func(t *testing.T, res *domain.Store, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			categoryRepo := new(mocks.CategoryRepository)
			storeRepo := new(mocks.StoreRepository)

			tc.builtSts(storeRepo, categoryRepo)
			u := usecase.NewStoreUsecase(storeRepo, nil, categoryRepo, time.Second*2)

			res, err := u.Get(context.TODO(), tc.arg)
			tc.checkResponse(t, res, err)

			categoryRepo.AssertExpectations(t)
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
	c := getCategory()

	a := args{
		page:  1,
		limit: 10,
		sort:  "created_at",
	}
	testCases := []struct {
		name          string
		args          args
		builtSts      func(storeRepo *mocks.StoreRepository, categoryRepo *mocks.CategoryRepository)
		checkResponse func(t *testing.T, res domain.Stores, count int64, err error)
	}{
		{
			name: "should fail if find all stores returns error",
			args: args{
				page:  0,
				limit: 0,
				sort:  "",
			},
			builtSts: func(storeRepo *mocks.StoreRepository, _ *mocks.CategoryRepository) {
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
			name: "should fail if get category by id returns error",
			args: a,
			builtSts: func(storeRepo *mocks.StoreRepository, categoryRepo *mocks.CategoryRepository) {
				stores := make(domain.Stores, 0)
				stores = append(stores, s)
				storeRepo.On("FindAll", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
					Return(stores, int64(1), nil).
					Once()
				categoryRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).
					Return(nil, errors.New("Unexpexted Error"))

			},
			checkResponse: func(t *testing.T, res domain.Stores, count int64, err error) {
				fmt.Println(err)
				assert.Len(t, res, 0)
				assert.Equal(t, count, int64(0))
				assert.Error(t, err)
			},
		},
		{
			name: "should succeed",
			args: a,
			builtSts: func(storeRepo *mocks.StoreRepository, categoryRepo *mocks.CategoryRepository) {
				stores := make(domain.Stores, 0)
				stores = append(stores, s)
				storeRepo.On("FindAll", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
					Return(stores, int64(1), nil).
					Once()
				categoryRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).
					Return(c, nil)
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
			categoryRepo := new(mocks.CategoryRepository)
			storeRepo := new(mocks.StoreRepository)

			tc.builtSts(storeRepo, categoryRepo)
			u := usecase.NewStoreUsecase(storeRepo, nil, categoryRepo, time.Second*2)

			res, count, err := u.Index(context.TODO(), tc.args.sort, tc.args.limit, tc.args.page)
			tc.checkResponse(t, res, count, err)

			categoryRepo.AssertExpectations(t)
			storeRepo.AssertExpectations(t)
		})
	}
}

func Test_StoreUsecase_Block(t *testing.T) {
	s := getStore()
	testCases := []struct {
		name          string
		arg           string
		builtSts      func(storeRepo *mocks.StoreRepository)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "should fail if find store by id returns error",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository) {
				storeRepo.On("FindByID", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "should fail if block returns error",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository) {
				store := s
				store.Status = domain.StoreStatusActive
				store.UserID = "invalid_id"
				storeRepo.
					On("FindByID", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "should fails if update store returns error",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository) {
				store := s
				store.Status = domain.StoreStatusActive
				storeRepo.
					On("FindByID", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
				storeRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "should succeed",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository) {
				store := s
				store.Status = domain.StoreStatusActive
				store.UserID = uuid.NewV4().String()
				storeRepo.
					On("FindByID", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
				storeRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storeRepo := new(mocks.StoreRepository)
			tc.builtSts(storeRepo)
			u := usecase.NewStoreUsecase(storeRepo, nil, nil, time.Second*2)
			err := u.Block(context.TODO(), tc.arg)
			tc.checkResponse(t, err)
		})
	}
}

func Test_StoreUsecase_Active(t *testing.T) {
	s := getStore()
	testCases := []struct {
		name          string
		arg           string
		builtSts      func(storeRepo *mocks.StoreRepository)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "should fail if find store by id returns error",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository) {
				storeRepo.On("FindByID", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "should fail activate returns error",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository) {
				store := s
				store.Status = domain.StoreStatusPending
				store.UserID = "invalid_id"
				storeRepo.
					On("FindByID", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "should fails if update store returns error",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository) {
				store := s
				store.Status = domain.StoreStatusPending
				storeRepo.
					On("FindByID", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
				storeRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "shuould succeed",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository) {
				store := s
				store.Status = domain.StoreStatusPending
				store.UserID = uuid.NewV4().String()
				storeRepo.
					On("FindByID", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
				storeRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storeRepo := new(mocks.StoreRepository)
			tc.builtSts(storeRepo)
			u := usecase.NewStoreUsecase(storeRepo, nil, nil, time.Second*2)
			err := u.Active(context.TODO(), tc.arg)
			tc.checkResponse(t, err)
		})
	}
}

func Test_StoreUsecase_Disable(t *testing.T) {
	s := getStore()

	testCases := []struct {
		name          string
		arg           string
		builtSts      func(storeRepo *mocks.StoreRepository)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "fail on find store",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository) {
				storeRepo.On("FindByID", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail on disable",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository) {
				store := s
				store.Status = domain.StoreStatusActive
				store.UserID = "invalid_id"
				storeRepo.
					On("FindByID", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail on update",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository) {
				store := s
				store.Status = domain.StoreStatusActive
				storeRepo.
					On("FindByID", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
				storeRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository) {
				store := s
				store.Status = domain.StoreStatusActive
				store.UserID = uuid.NewV4().String()
				storeRepo.
					On("FindByID", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
				storeRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storeRepo := new(mocks.StoreRepository)
			tc.builtSts(storeRepo)
			u := usecase.NewStoreUsecase(storeRepo, nil, nil, time.Second*2)
			err := u.Disable(context.TODO(), tc.arg)
			tc.checkResponse(t, err)
		})
	}
}

func Test_StoreUsecase_Update(t *testing.T) {
	s := getStore()

	testCases := []struct {
		name          string
		arg           *domain.Store
		builtSts      func(storeRepo *mocks.StoreRepository)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "fail error on store's repo",
			arg:  s,
			builtSts: func(storeRepo *mocks.StoreRepository) {
				storeRepo.On("FindByID", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail error on store's repo update",
			arg:  s,
			builtSts: func(storeRepo *mocks.StoreRepository) {
				storeRepo.On("FindByID", mock.Anything, mock.AnythingOfType("string")).Return(s, nil).Once()
				storeRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			arg:  s,
			builtSts: func(storeRepo *mocks.StoreRepository) {
				storeRepo.On("FindByID", mock.Anything, mock.AnythingOfType("string")).Return(s, nil).Once()
				storeRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storeRepo := new(mocks.StoreRepository)
			tc.builtSts(storeRepo)
			u := usecase.NewStoreUsecase(storeRepo, nil, nil, time.Second*2)
			err := u.Update(context.TODO(), tc.arg)
			tc.checkResponse(t, err)
		})
	}
}

func Test_StoreUsecase_Delete(t *testing.T) {

	s := getStore()

	testCases := []struct {
		name          string
		arg           string
		builtSts      func(storeRepo *mocks.StoreRepository, accountRepo *mocks.AccountRepository)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "fail on find store",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository, _ *mocks.AccountRepository) {
				storeRepo.On("FindByID", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{

			name: "fail on delete account",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository, accountRepo *mocks.AccountRepository) {
				store := s
				store.Status = domain.StoreStatusActive
				storeRepo.
					On("FindByID", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
				accountRepo.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail on update",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository, accountRepo *mocks.AccountRepository) {
				store := s
				store.Status = domain.StoreStatusActive
				storeRepo.
					On("FindByID", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
				accountRepo.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(nil).Once()
				storeRepo.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository, accountRepo *mocks.AccountRepository) {
				store := s
				store.Status = domain.StoreStatusActive
				store.UserID = uuid.NewV4().String()
				storeRepo.
					On("FindByID", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
				accountRepo.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(nil).Once()
				storeRepo.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(nil).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storeRepo := new(mocks.StoreRepository)
			accountRepo := new(mocks.AccountRepository)
			tc.builtSts(storeRepo, accountRepo)
			u := usecase.NewStoreUsecase(storeRepo, accountRepo, nil, time.Second*2)
			err := u.Delete(context.TODO(), tc.arg)
			tc.checkResponse(t, err)
		})
	}
}
