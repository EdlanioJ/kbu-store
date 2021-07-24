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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CategoryUsecase_Create(t *testing.T) {
	a := getCategory()

	testCases := []struct {
		name          string
		arg           *domain.Category
		builtSts      func(categoryRepo *mocks.CategoryRepository)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "fail on new category",
			arg: &domain.Category{
				Name: "Category",
			},
			builtSts: func(_ *mocks.CategoryRepository) {
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail on categories's repo",
			arg:  a,
			builtSts: func(categoryRepo *mocks.CategoryRepository) {
				categoryRepo.On("Store", mock.Anything, mock.Anything).Return(errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			arg:  a,
			builtSts: func(categoryRepo *mocks.CategoryRepository) {
				categoryRepo.On("Store", mock.Anything, mock.Anything).Return(nil).Once()
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
				categoryRepo.On("FindByID", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "fail on update",
			args: args{
				category: category,
			},
			builtSts: func(categoryRepo *mocks.CategoryRepository) {
				categoryRepo.On("FindByID", mock.Anything, mock.AnythingOfType("string")).Return(category, nil).Once()
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
