package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/EdlanioJ/kbu-store/data/usecase"
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/domain/mocks"
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
		externalID:  s.ExternalID,
		tags:        s.Tags,
	}
	testCases := []struct {
		name          string
		args          args
		builtSts      func(storeRepo *mocks.StoreRepository, accountRepo *mocks.AccountRepository, categoryRepo *mocks.CategoryRepository)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "fail on create account",
			args: a,
			builtSts: func(_ *mocks.StoreRepository, accountRepo *mocks.AccountRepository, _ *mocks.CategoryRepository) {
				accountRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail on category's repo get by id",
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
			name: "fail on new store",
			args: args{
				name:        "",
				description: s.Description,
				categoryID:  s.Category.ID,
				externalID:  s.ExternalID,
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
			name: "fail on store's repo create",
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
			name: "success",
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

			err := u.Create(context.TODO(), ts.args.name, ts.args.description, ts.args.categoryID, ts.args.externalID, ts.args.tags, 0, 0)
			ts.checkResponse(t, err)
			accountRepo.AssertExpectations(t)
			categoryRepo.AssertExpectations(t)
			storeRepo.AssertExpectations(t)
		})
	}
}

func Test_StoreUsecase_GetById(t *testing.T) {

	s := getStore()
	c := getCategory()

	testCases := []struct {
		name          string
		arg           string
		builtSts      func(storeRepo *mocks.StoreRepository, categoryRepo *mocks.CategoryRepository)
		checkResponse func(t *testing.T, res *domain.Store, err error)
	}{
		{
			name: "fail on store repo",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository, _ *mocks.CategoryRepository) {
				storeRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, res *domain.Store, err error) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "fail on category's repo",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository, categoryRepo *mocks.CategoryRepository) {
				storeRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(s, nil).Once()
				categoryRepo.On("GetById", mock.Anything, s.Category.ID).Return(nil, errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, res *domain.Store, err error) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "success",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository, categoryRepo *mocks.CategoryRepository) {
				storeRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(s, nil).Once()
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

			res, err := u.GetById(context.TODO(), tc.arg)
			tc.checkResponse(t, res, err)

			categoryRepo.AssertExpectations(t)
			storeRepo.AssertExpectations(t)
		})
	}
}

func Test_StoreUsecase_GetByIdAndOwner(t *testing.T) {
	type args struct {
		id    string
		owner string
	}
	s := getStore()
	c := getCategory()

	testCases := []struct {
		name          string
		args          args
		builtSts      func(storeRepo *mocks.StoreRepository, categoryRepo *mocks.CategoryRepository)
		checkResponse func(t *testing.T, res *domain.Store, err error)
	}{
		{
			name: "fail on get store repo",
			args: args{
				id:    uuid.NewV4().String(),
				owner: uuid.NewV4().String(),
			},
			builtSts: func(storeRepo *mocks.StoreRepository, _ *mocks.CategoryRepository) {
				storeRepo.On("GetByIdAndOwner", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, res *domain.Store, err error) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},

		{
			name: "fail on category's repo",
			args: args{
				id:    uuid.NewV4().String(),
				owner: uuid.NewV4().String(),
			},
			builtSts: func(storeRepo *mocks.StoreRepository, categoryRepo *mocks.CategoryRepository) {
				storeRepo.On("GetByIdAndOwner", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(s, nil).Once()
				categoryRepo.On("GetById", mock.Anything, s.Category.ID).Return(nil, errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, res *domain.Store, err error) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "success",
			args: args{
				id:    uuid.NewV4().String(),
				owner: uuid.NewV4().String(),
			},
			builtSts: func(storeRepo *mocks.StoreRepository, categoryRepo *mocks.CategoryRepository) {
				storeRepo.On("GetByIdAndOwner", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(s, nil).Once()
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

			res, err := u.GetByIdAndOwner(context.TODO(), tc.args.id, tc.args.owner)
			tc.checkResponse(t, res, err)

			categoryRepo.AssertExpectations(t)
			storeRepo.AssertExpectations(t)
		})
	}
}

func Test_StoreUsecase_GetAll(t *testing.T) {
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
		checkResponse func(t *testing.T, res []*domain.Store, count int64, err error)
	}{
		{
			name: "fail on store repo",
			args: args{
				page:  0,
				limit: 0,
				sort:  "",
			},
			builtSts: func(storeRepo *mocks.StoreRepository, _ *mocks.CategoryRepository) {
				storeRepo.On("GetAll", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
					Return(nil, int64(0), errors.New("Unexpexted Error")).
					Once()
			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
				assert.Len(t, res, 0)
				assert.Equal(t, count, int64(0))
				assert.Error(t, err)
			},
		},
		{
			name: "fail on category repo",
			args: a,
			builtSts: func(storeRepo *mocks.StoreRepository, categoryRepo *mocks.CategoryRepository) {
				stores := make([]*domain.Store, 0)
				stores = append(stores, s)
				storeRepo.On("GetAll", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
					Return(stores, int64(1), nil).
					Once()
				categoryRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).
					Return(nil, errors.New("Unexpexted Error"))

			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
				assert.Len(t, res, 0)
				assert.Equal(t, count, int64(0))
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			args: a,
			builtSts: func(storeRepo *mocks.StoreRepository, categoryRepo *mocks.CategoryRepository) {
				stores := make([]*domain.Store, 0)
				stores = append(stores, s)
				storeRepo.On("GetAll", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
					Return(stores, int64(1), nil).
					Once()
				categoryRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).
					Return(c, nil)
			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
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

			res, count, err := u.GetAll(context.TODO(), tc.args.sort, tc.args.limit, tc.args.page)
			tc.checkResponse(t, res, count, err)

			categoryRepo.AssertExpectations(t)
			storeRepo.AssertExpectations(t)
		})
	}
}

func Test_StoreUsecase_GetAllByCategory(t *testing.T) {
	type args struct {
		categoryId string
		page       int
		limit      int
		sort       string
	}
	s := getStore()
	c := getCategory()

	a := args{
		categoryId: uuid.NewV4().String(),
		page:       1,
		limit:      10,
		sort:       "created_at",
	}
	testCases := []struct {
		name          string
		args          args
		builtSts      func(storeRepo *mocks.StoreRepository, categoryRepo *mocks.CategoryRepository)
		checkResponse func(t *testing.T, res []*domain.Store, count int64, err error)
	}{
		{
			name: "fail on store repo",
			args: args{
				categoryId: uuid.NewV4().String(),
				page:       0,
				limit:      0,
				sort:       "",
			},
			builtSts: func(storeRepo *mocks.StoreRepository, _ *mocks.CategoryRepository) {
				storeRepo.On("GetAllByCategory", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
					Return(nil, int64(0), errors.New("Unexpexted Error")).
					Once()
			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
				assert.Len(t, res, 0)
				assert.Equal(t, count, int64(0))
				assert.Error(t, err)
			},
		},
		{
			name: "fail on category repo",
			args: a,
			builtSts: func(storeRepo *mocks.StoreRepository, categoryRepo *mocks.CategoryRepository) {
				stores := make([]*domain.Store, 0)
				stores = append(stores, s)
				storeRepo.On("GetAllByCategory", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
					Return(stores, int64(1), nil).
					Once()
				categoryRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).
					Return(nil, errors.New("Unexpexted Error"))

			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
				assert.Len(t, res, 0)
				assert.Equal(t, count, int64(0))
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			args: a,
			builtSts: func(storeRepo *mocks.StoreRepository, categoryRepo *mocks.CategoryRepository) {
				stores := make([]*domain.Store, 0)
				stores = append(stores, s)
				storeRepo.On("GetAllByCategory", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
					Return(stores, int64(1), nil).
					Once()
				categoryRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).
					Return(c, nil)
			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
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

			res, count, err := u.GetAllByCategory(context.TODO(), tc.args.categoryId, tc.args.sort, tc.args.limit, tc.args.page)
			tc.checkResponse(t, res, count, err)

			categoryRepo.AssertExpectations(t)
			storeRepo.AssertExpectations(t)
		})
	}
}

func Test_StoreUsecase_GetAllByOwner(t *testing.T) {
	type args struct {
		ownerId string
		page    int
		limit   int
		sort    string
	}
	s := getStore()
	c := getCategory()

	a := args{
		ownerId: uuid.NewV4().String(),
		page:    1,
		limit:   10,
		sort:    "created_at",
	}
	testCases := []struct {
		name          string
		args          args
		builtSts      func(storeRepo *mocks.StoreRepository, categoryRepo *mocks.CategoryRepository)
		checkResponse func(t *testing.T, res []*domain.Store, count int64, err error)
	}{
		{
			name: "fail on store repo",
			args: args{
				ownerId: uuid.NewV4().String(),
				page:    0,
				limit:   0,
				sort:    "",
			},
			builtSts: func(storeRepo *mocks.StoreRepository, _ *mocks.CategoryRepository) {
				storeRepo.On("GetAllByOwner", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
					Return(nil, int64(0), errors.New("Unexpexted Error")).
					Once()
			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
				assert.Len(t, res, 0)
				assert.Equal(t, count, int64(0))
				assert.Error(t, err)
			},
		},
		{
			name: "fail on category repo",
			args: a,
			builtSts: func(storeRepo *mocks.StoreRepository, categoryRepo *mocks.CategoryRepository) {
				stores := make([]*domain.Store, 0)
				stores = append(stores, s)
				storeRepo.On("GetAllByOwner", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
					Return(stores, int64(1), nil).
					Once()
				categoryRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).
					Return(nil, errors.New("Unexpexted Error"))

			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
				assert.Len(t, res, 0)
				assert.Equal(t, count, int64(0))
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			args: a,
			builtSts: func(storeRepo *mocks.StoreRepository, categoryRepo *mocks.CategoryRepository) {
				stores := make([]*domain.Store, 0)
				stores = append(stores, s)
				storeRepo.On("GetAllByOwner", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
					Return(stores, int64(1), nil).
					Once()
				categoryRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).
					Return(c, nil)
			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
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

			res, count, err := u.GetAllByOwner(context.TODO(), tc.args.ownerId, tc.args.sort, tc.args.limit, tc.args.page)
			tc.checkResponse(t, res, count, err)

			categoryRepo.AssertExpectations(t)
			storeRepo.AssertExpectations(t)
		})
	}
}

func Test_StoreUsecase_GetAllByStatus(t *testing.T) {
	type args struct {
		status string
		page   int
		limit  int
		sort   string
	}
	s := getStore()
	c := getCategory()

	a := args{
		status: domain.StoreStatusActive,
		page:   1,
		limit:  10,
		sort:   "created_at",
	}

	testCases := []struct {
		name          string
		args          args
		builtSts      func(storeRepo *mocks.StoreRepository, categoryRepo *mocks.CategoryRepository)
		checkResponse func(t *testing.T, res []*domain.Store, count int64, err error)
	}{
		{
			name: "fail on store repo",
			args: args{
				status: domain.StoreStatusActive,
				page:   0,
				limit:  0,
				sort:   "",
			},
			builtSts: func(storeRepo *mocks.StoreRepository, _ *mocks.CategoryRepository) {
				storeRepo.On("GetAllByStatus", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
					Return(nil, int64(0), errors.New("Unexpexted Error")).
					Once()
			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
				assert.Len(t, res, 0)
				assert.Equal(t, count, int64(0))
				assert.Error(t, err)
			},
		},
		{
			name: "fail on category repo",
			args: a,
			builtSts: func(storeRepo *mocks.StoreRepository, categoryRepo *mocks.CategoryRepository) {
				stores := make([]*domain.Store, 0)
				stores = append(stores, s)
				storeRepo.On("GetAllByStatus", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
					Return(stores, int64(1), nil).
					Once()
				categoryRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).
					Return(nil, errors.New("Unexpexted Error"))

			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
				assert.Len(t, res, 0)
				assert.Equal(t, count, int64(0))
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			args: a,
			builtSts: func(storeRepo *mocks.StoreRepository, categoryRepo *mocks.CategoryRepository) {
				stores := make([]*domain.Store, 0)
				stores = append(stores, s)
				storeRepo.On("GetAllByStatus", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
					Return(stores, int64(1), nil).
					Once()
				categoryRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).
					Return(c, nil)
			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
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

			res, count, err := u.GetAllByStatus(context.TODO(), tc.args.status, tc.args.sort, tc.args.limit, tc.args.page)
			tc.checkResponse(t, res, count, err)

			categoryRepo.AssertExpectations(t)
			storeRepo.AssertExpectations(t)
		})
	}
}

func Test_StoreUsecase_GetAllByTags(t *testing.T) {
	type args struct {
		tags  []string
		page  int
		limit int
		sort  string
	}
	s := getStore()
	c := getCategory()

	a := args{
		tags:  s.Tags,
		page:  1,
		limit: 10,
		sort:  "created_at",
	}

	testCases := []struct {
		name          string
		args          args
		builtSts      func(storeRepo *mocks.StoreRepository, categoryRepo *mocks.CategoryRepository)
		checkResponse func(t *testing.T, res []*domain.Store, count int64, err error)
	}{
		{
			name: "fail on store repo",
			args: args{
				tags:  s.Tags,
				page:  0,
				limit: 0,
				sort:  "",
			},
			builtSts: func(storeRepo *mocks.StoreRepository, _ *mocks.CategoryRepository) {
				storeRepo.On("GetAllByTags", mock.Anything, mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
					Return(nil, int64(0), errors.New("Unexpexted Error")).
					Once()
			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
				assert.Len(t, res, 0)
				assert.Equal(t, count, int64(0))
				assert.Error(t, err)
			},
		},
		{
			name: "fail on category repo",
			args: a,
			builtSts: func(storeRepo *mocks.StoreRepository, categoryRepo *mocks.CategoryRepository) {
				stores := make([]*domain.Store, 0)
				stores = append(stores, s)
				storeRepo.On("GetAllByTags", mock.Anything, mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
					Return(stores, int64(1), nil).
					Once()
				categoryRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).
					Return(nil, errors.New("Unexpexted Error"))

			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
				assert.Len(t, res, 0)
				assert.Equal(t, count, int64(0))
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			args: a,
			builtSts: func(storeRepo *mocks.StoreRepository, categoryRepo *mocks.CategoryRepository) {
				stores := make([]*domain.Store, 0)
				stores = append(stores, s)
				storeRepo.On("GetAllByTags", mock.Anything, mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
					Return(stores, int64(1), nil).
					Once()
				categoryRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).
					Return(c, nil)
			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
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

			res, count, err := u.GetAllByTags(context.TODO(), tc.args.tags, tc.args.sort, tc.args.limit, tc.args.page)
			tc.checkResponse(t, res, count, err)

			categoryRepo.AssertExpectations(t)
			storeRepo.AssertExpectations(t)
		})
	}
}

func Test_StoreUsecase_GetAllByByCloseLocation(t *testing.T) {
	type args struct {
		lat      float64
		lng      float64
		distance int
		status   string
		page     int
		limit    int
		sort     string
	}
	s := getStore()
	c := getCategory()
	a := args{
		lat:      s.Position.Lat,
		lng:      s.Position.Lng,
		distance: 10,
		status:   domain.StoreStatusActive,
		page:     1,
		limit:    10,
		sort:     "created_at",
	}

	testCases := []struct {
		name          string
		args          args
		builtSts      func(storeRepo *mocks.StoreRepository, categoryRepo *mocks.CategoryRepository)
		checkResponse func(t *testing.T, res []*domain.Store, count int64, err error)
	}{
		{
			name: "fail on store repo",
			args: args{
				lat:      s.Position.Lat,
				lng:      s.Position.Lng,
				distance: 10,
				status:   domain.StoreStatusActive,
				page:     0,
				limit:    0,
				sort:     "",
			},
			builtSts: func(storeRepo *mocks.StoreRepository, _ *mocks.CategoryRepository) {
				storeRepo.On("GetAllByLocation", mock.Anything, mock.AnythingOfType("float64"), mock.AnythingOfType("float64"), mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).
					Return(nil, int64(0), errors.New("Unexpexted Error")).
					Once()
			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
				assert.Len(t, res, 0)
				assert.Equal(t, count, int64(0))
				assert.Error(t, err)
			},
		},
		{
			name: "fail on category repo",
			args: a,
			builtSts: func(storeRepo *mocks.StoreRepository, categoryRepo *mocks.CategoryRepository) {
				stores := make([]*domain.Store, 0)
				stores = append(stores, s)
				storeRepo.On("GetAllByLocation", mock.Anything, mock.AnythingOfType("float64"), mock.AnythingOfType("float64"), mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).
					Return(stores, int64(1), nil).
					Once()
				categoryRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).
					Return(nil, errors.New("Unexpexted Error"))

			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
				assert.Len(t, res, 0)
				assert.Equal(t, count, int64(0))
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			args: a,
			builtSts: func(storeRepo *mocks.StoreRepository, categoryRepo *mocks.CategoryRepository) {
				stores := make([]*domain.Store, 0)
				stores = append(stores, s)
				storeRepo.On("GetAllByLocation", mock.Anything, mock.AnythingOfType("float64"), mock.AnythingOfType("float64"), mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).
					Return(stores, int64(1), nil).
					Once()
				categoryRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).
					Return(c, nil)
			},
			checkResponse: func(t *testing.T, res []*domain.Store, count int64, err error) {
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

			res, count, err := u.GetAllByByCloseLocation(context.TODO(), tc.args.lat, tc.args.lng, tc.args.distance, tc.args.status, tc.args.limit, tc.args.page, tc.args.sort)
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
			name: "fail on find store",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository) {
				storeRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail on empty store",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository) {
				storeRepo.
					On("GetById", mock.Anything, mock.AnythingOfType("string")).
					Return(&domain.Store{}, nil).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, domain.ErrNotFound.Error())
			},
		},
		{
			name: "fail on status already block",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository) {
				store := s
				store.Status = domain.StoreStatusBlock
				storeRepo.
					On("GetById", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, domain.ErrBlocked.Error())
			},
		},
		{
			name: "fail on status still pending",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository) {
				store := s
				store.Status = domain.StoreStatusPending
				storeRepo.
					On("GetById", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, domain.ErrIsPending.Error())
			},
		},
		{
			name: "fail on block",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository) {
				store := s
				store.Status = domain.StoreStatusActive
				store.ExternalID = "invalid_id"
				storeRepo.
					On("GetById", mock.Anything, mock.AnythingOfType("string")).
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
					On("GetById", mock.Anything, mock.AnythingOfType("string")).
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
				store.ExternalID = uuid.NewV4().String()
				storeRepo.
					On("GetById", mock.Anything, mock.AnythingOfType("string")).
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
			name: "fail on find store",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository) {
				storeRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail on empty store",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository) {
				storeRepo.
					On("GetById", mock.Anything, mock.AnythingOfType("string")).
					Return(&domain.Store{}, nil).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, domain.ErrNotFound.Error())
			},
		},
		{
			name: "fail on status already active",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository) {
				store := s
				store.Status = domain.StoreStatusActive
				storeRepo.
					On("GetById", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, domain.ErrActived.Error())
			},
		},

		{
			name: "fail on activate",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository) {
				store := s
				store.Status = domain.StoreStatusPending
				store.ExternalID = "invalid_id"
				storeRepo.
					On("GetById", mock.Anything, mock.AnythingOfType("string")).
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
				store.Status = domain.StoreStatusPending
				storeRepo.
					On("GetById", mock.Anything, mock.AnythingOfType("string")).
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
				store.Status = domain.StoreStatusPending
				store.ExternalID = uuid.NewV4().String()
				storeRepo.
					On("GetById", mock.Anything, mock.AnythingOfType("string")).
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
				storeRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail on empty store",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository) {
				storeRepo.
					On("GetById", mock.Anything, mock.AnythingOfType("string")).
					Return(&domain.Store{}, nil).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, domain.ErrNotFound.Error())
			},
		},
		{
			name: "fail on store already disabled",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository) {
				store := s
				store.Status = domain.StoreStatusInactive
				storeRepo.
					On("GetById", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, domain.ErrInactived.Error())
			},
		},
		{
			name: "fail on store is blocked",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository) {
				store := s
				store.Status = domain.StoreStatusBlock
				storeRepo.
					On("GetById", mock.Anything, mock.AnythingOfType("string")).
					Return(store, nil).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, domain.ErrBlocked.Error())
			},
		},
		{
			name: "fail on disable",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository) {
				store := s
				store.Status = domain.StoreStatusActive
				store.ExternalID = "invalid_id"
				storeRepo.
					On("GetById", mock.Anything, mock.AnythingOfType("string")).
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
					On("GetById", mock.Anything, mock.AnythingOfType("string")).
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
				store.ExternalID = uuid.NewV4().String()
				storeRepo.
					On("GetById", mock.Anything, mock.AnythingOfType("string")).
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
				storeRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail error on store's repo get by id returns empty",
			arg:  s,
			builtSts: func(storeRepo *mocks.StoreRepository) {
				storeRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(&domain.Store{}, nil).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, domain.ErrNotFound.Error())
			},
		},
		{
			name: "fail error on store's repo update",
			arg:  s,
			builtSts: func(storeRepo *mocks.StoreRepository) {
				storeRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(s, nil).Once()
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
				storeRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(s, nil).Once()
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
				storeRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail on empty store",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository, _ *mocks.AccountRepository) {
				storeRepo.
					On("GetById", mock.Anything, mock.AnythingOfType("string")).
					Return(&domain.Store{}, nil).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, domain.ErrNotFound.Error())
			},
		},
		{

			name: "fail on delete account",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeRepo *mocks.StoreRepository, accountRepo *mocks.AccountRepository) {
				store := s
				store.Status = domain.StoreStatusActive
				storeRepo.
					On("GetById", mock.Anything, mock.AnythingOfType("string")).
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
					On("GetById", mock.Anything, mock.AnythingOfType("string")).
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
				store.ExternalID = uuid.NewV4().String()
				storeRepo.
					On("GetById", mock.Anything, mock.AnythingOfType("string")).
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
