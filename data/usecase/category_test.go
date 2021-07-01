package usecase_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/EdlanioJ/kbu-store/data/usecase"
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/domain/mocks"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CategoryUsecase_Create(t *testing.T) {

	testCases := []struct {
		name          string
		arg           string
		builtSts      func(categoryRepo *mocks.CategoryRepository)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "fail on new category",
			arg:  "",
			builtSts: func(_ *mocks.CategoryRepository) {
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail on categories's repo",
			arg:  "New Category",
			builtSts: func(categoryRepo *mocks.CategoryRepository) {
				categoryRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			arg:  "New Category",
			builtSts: func(categoryRepo *mocks.CategoryRepository) {
				categoryRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			categoryRepo := new(mocks.CategoryRepository)
			tc.builtSts(categoryRepo)
			u := usecase.NewCategoryUsecase(categoryRepo, time.Second*2)
			fmt.Println(tc.arg)
			err := u.Create(context.TODO(), tc.arg)
			tc.checkResponse(t, err)
			categoryRepo.AssertExpectations(t)
		})
	}
}

func Test_CategoryUsecase_GetById(t *testing.T) {
	testCases := []struct {
		name          string
		arg           string
		builtSts      func(categoryRepo *mocks.CategoryRepository)
		checkResponse func(t *testing.T, res *domain.Category, err error)
	}{
		{
			name: "fail",
			arg:  uuid.NewV4().String(),
			builtSts: func(categoryRepo *mocks.CategoryRepository) {
				categoryRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, res *domain.Category, err error) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "success",
			arg:  uuid.NewV4().String(),
			builtSts: func(categoryRepo *mocks.CategoryRepository) {
				category := getCategory()
				categoryRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(category, nil).Once()
			},
			checkResponse: func(t *testing.T, res *domain.Category, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			categoryRepo := new(mocks.CategoryRepository)
			tc.builtSts(categoryRepo)
			u := usecase.NewCategoryUsecase(categoryRepo, time.Second*2)
			fmt.Println(tc.arg)
			res, err := u.GetById(context.TODO(), tc.arg)
			tc.checkResponse(t, res, err)
			categoryRepo.AssertExpectations(t)
		})
	}
}

func Test_CategoryUsecase_GetByIdAndStatus(t *testing.T) {
	type args struct {
		id     string
		status string
	}
	testCases := []struct {
		name          string
		args          args
		builtSts      func(categoryRepo *mocks.CategoryRepository)
		checkResponse func(t *testing.T, res *domain.Category, err error)
	}{
		{
			name: "fail",
			args: args{
				id:     uuid.NewV4().String(),
				status: domain.CategoryStatusActive,
			},
			builtSts: func(categoryRepo *mocks.CategoryRepository) {
				categoryRepo.On("GetByIdAndStatus", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, res *domain.Category, err error) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "success",
			args: args{
				id:     uuid.NewV4().String(),
				status: domain.CategoryStatusActive,
			},
			builtSts: func(categoryRepo *mocks.CategoryRepository) {
				categoryRepo.On("GetByIdAndStatus", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(getCategory(), nil).Once()
			},
			checkResponse: func(t *testing.T, res *domain.Category, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			categoryRepo := new(mocks.CategoryRepository)
			tc.builtSts(categoryRepo)
			u := usecase.NewCategoryUsecase(categoryRepo, time.Second*2)
			res, err := u.GetByIdAndStatus(context.TODO(), tc.args.id, tc.args.status)
			tc.checkResponse(t, res, err)
			categoryRepo.AssertExpectations(t)
		})
	}
}

func Test_CategoryUsecase_GetAll(t *testing.T) {
	type args struct {
		sort  string
		page  int
		limit int
	}
	testCases := []struct {
		name          string
		args          args
		builtSts      func(categoryRepo *mocks.CategoryRepository)
		checkResponse func(t *testing.T, res []*domain.Category, count int64, err error)
	}{
		{
			name: "fail",
			args: args{
				sort:  "",
				page:  0,
				limit: 0,
			},
			builtSts: func(categoryRepo *mocks.CategoryRepository) {
				categoryRepo.On("GetAll", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(nil, int64(0), errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, res []*domain.Category, count int64, err error) {
				assert.Len(t, res, 0)
				assert.Equal(t, count, int64(0))
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			args: args{
				sort:  "",
				page:  0,
				limit: 0,
			},
			builtSts: func(categoryRepo *mocks.CategoryRepository) {
				category := getCategory()
				categories := make([]*domain.Category, 0)
				categories = append(categories, category)
				categoryRepo.On("GetAll", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(categories, int64(1), nil).Once()
			},
			checkResponse: func(t *testing.T, res []*domain.Category, count int64, err error) {
				assert.Len(t, res, 1)
				assert.Equal(t, count, int64(1))
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			categoryRepo := new(mocks.CategoryRepository)
			tc.builtSts(categoryRepo)
			u := usecase.NewCategoryUsecase(categoryRepo, time.Second*2)
			res, count, err := u.GetAll(context.TODO(), tc.args.sort, tc.args.page, tc.args.limit)
			tc.checkResponse(t, res, count, err)
			categoryRepo.AssertExpectations(t)
		})
	}
}

func Test_CategoryUsecase_GetAllByStatus(t *testing.T) {
	type args struct {
		status string
		sort   string
		page   int
		limit  int
	}
	testCases := []struct {
		name          string
		args          args
		builtSts      func(categoryRepo *mocks.CategoryRepository)
		checkResponse func(t *testing.T, res []*domain.Category, count int64, err error)
	}{
		{
			name: "fail",
			args: args{
				status: domain.CategoryStatusActive,
				sort:   "",
				page:   0,
				limit:  0,
			},
			builtSts: func(categoryRepo *mocks.CategoryRepository) {
				categoryRepo.On("GetAllByStatus", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(nil, int64(0), errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, res []*domain.Category, count int64, err error) {
				assert.Len(t, res, 0)
				assert.Equal(t, count, int64(0))
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			args: args{
				status: domain.CategoryStatusInactive,
				sort:   "",
				page:   0,
				limit:  0,
			},
			builtSts: func(categoryRepo *mocks.CategoryRepository) {
				category := getCategory()
				categories := make([]*domain.Category, 0)
				categories = append(categories, category)
				categoryRepo.On("GetAllByStatus", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(categories, int64(1), nil).Once()
			},
			checkResponse: func(t *testing.T, res []*domain.Category, count int64, err error) {
				assert.Len(t, res, 1)
				assert.Equal(t, count, int64(1))
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			categoryRepo := new(mocks.CategoryRepository)
			tc.builtSts(categoryRepo)
			u := usecase.NewCategoryUsecase(categoryRepo, time.Second*2)
			res, count, err := u.GetAllByStatus(context.TODO(), tc.args.status, tc.args.sort, tc.args.page, tc.args.limit)
			tc.checkResponse(t, res, count, err)
			categoryRepo.AssertExpectations(t)
		})
	}
}

func Test_CategoryUsecase_Update(t *testing.T) {
	category := getCategory()
	type args struct {
		category *domain.Category
	}
	testCases := []struct {
		name          string
		args          args
		builtSts      func(categoryRepo *mocks.CategoryRepository)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "fail on get category",
			args: args{
				category: category,
			},
			builtSts: func(categoryRepo *mocks.CategoryRepository) {
				categoryRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail on get category returns empty",
			args: args{
				category: category,
			},
			builtSts: func(categoryRepo *mocks.CategoryRepository) {
				categoryRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(&domain.Category{}, nil).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, domain.ErrNotFound.Error())
			},
		},
		{
			name: "fail on update",
			args: args{
				category: category,
			},
			builtSts: func(categoryRepo *mocks.CategoryRepository) {
				categoryRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(category, nil).Once()
				categoryRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			categoryRepo := new(mocks.CategoryRepository)
			tc.builtSts(categoryRepo)
			u := usecase.NewCategoryUsecase(categoryRepo, time.Second*2)
			err := u.Update(context.TODO(), tc.args.category)
			tc.checkResponse(t, err)
			categoryRepo.AssertExpectations(t)
		})
	}
}

func Test_CategoryRepo_Activate(t *testing.T) {
	testCases := []struct {
		name          string
		arg           string
		builtSts      func(categoryRepo *mocks.CategoryRepository)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "fail on get category by id",
			arg:  uuid.NewV4().String(),
			builtSts: func(categoryRepo *mocks.CategoryRepository) {
				categoryRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail on get category by id returns empty",
			arg:  uuid.NewV4().String(),
			builtSts: func(categoryRepo *mocks.CategoryRepository) {
				categoryRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(&domain.Category{}, nil).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, domain.ErrNotFound.Error())
			},
		},
		{
			name: "fail on get category by id returns actived category",
			arg:  uuid.NewV4().String(),
			builtSts: func(categoryRepo *mocks.CategoryRepository) {
				c := getCategory()
				c.Status = domain.CategoryStatusActive
				categoryRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(c, nil).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, domain.ErrActived.Error())
			},
		},
		{
			name: "fail on activate",
			arg:  uuid.NewV4().String(),
			builtSts: func(categoryRepo *mocks.CategoryRepository) {
				c := getCategory()
				c.Name = ""
				categoryRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(c, nil).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail on update",
			arg:  uuid.NewV4().String(),
			builtSts: func(categoryRepo *mocks.CategoryRepository) {
				c := getCategory()
				categoryRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(c, nil).Once()
				categoryRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			arg:  uuid.NewV4().String(),
			builtSts: func(categoryRepo *mocks.CategoryRepository) {
				c := getCategory()
				categoryRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(c, nil).Once()
				categoryRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		categoryRepo := new(mocks.CategoryRepository)
		tc.builtSts(categoryRepo)
		u := usecase.NewCategoryUsecase(categoryRepo, time.Second*2)
		err := u.Activate(context.TODO(), tc.arg)
		tc.checkResponse(t, err)
	}
}

func Test_CategoryRepo_Disable(t *testing.T) {
	testCases := []struct {
		name          string
		arg           string
		builtSts      func(categoryRepo *mocks.CategoryRepository)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "fail on get category by id",
			arg:  uuid.NewV4().String(),
			builtSts: func(categoryRepo *mocks.CategoryRepository) {
				categoryRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail on get category by id returns empty",
			arg:  uuid.NewV4().String(),
			builtSts: func(categoryRepo *mocks.CategoryRepository) {
				categoryRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(&domain.Category{}, nil).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, domain.ErrNotFound.Error())
			},
		},
		{
			name: "fail on get category by id returns blocked category",
			arg:  uuid.NewV4().String(),
			builtSts: func(categoryRepo *mocks.CategoryRepository) {
				c := getCategory()
				c.Status = domain.CategoryStatusInactive
				categoryRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(c, nil).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, domain.ErrBlocked.Error())
			},
		},
		{
			name: "fail on get category by id returns pending category",
			arg:  uuid.NewV4().String(),
			builtSts: func(categoryRepo *mocks.CategoryRepository) {
				c := getCategory()
				c.Status = domain.CategoryStatusPending
				categoryRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(c, nil).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, domain.ErrIsPending.Error())
			},
		},
		{
			name: "fail on activate",
			arg:  uuid.NewV4().String(),
			builtSts: func(categoryRepo *mocks.CategoryRepository) {
				c := getCategory()
				c.Status = domain.CategoryStatusActive
				c.Name = ""
				categoryRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(c, nil).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail on update",
			arg:  uuid.NewV4().String(),
			builtSts: func(categoryRepo *mocks.CategoryRepository) {
				c := getCategory()
				c.Status = domain.CategoryStatusActive
				categoryRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(c, nil).Once()
				categoryRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			arg:  uuid.NewV4().String(),
			builtSts: func(categoryRepo *mocks.CategoryRepository) {
				c := getCategory()
				c.Status = domain.CategoryStatusActive
				categoryRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(c, nil).Once()
				categoryRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		categoryRepo := new(mocks.CategoryRepository)
		tc.builtSts(categoryRepo)
		u := usecase.NewCategoryUsecase(categoryRepo, time.Second*2)
		err := u.Disable(context.TODO(), tc.arg)
		tc.checkResponse(t, err)
	}
}
