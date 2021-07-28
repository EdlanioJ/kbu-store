package usecases_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/EdlanioJ/kbu-store/app/domain"
	"github.com/EdlanioJ/kbu-store/app/usecases"
	"github.com/EdlanioJ/kbu-store/app/utils/mocks"
	"github.com/EdlanioJ/kbu-store/app/utils/sample"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CategoryUsecase_Create(t *testing.T) {
	testCases := []struct {
		name        string
		arg         *domain.Category
		expectedErr bool
		prepare     func(categoryRepo *mocks.CategoryRepository)
	}{
		{
			name:        "failure_store_category_returns_error",
			arg:         sample.NewCategory(),
			expectedErr: true,
			prepare: func(categoryRepo *mocks.CategoryRepository) {
				categoryRepo.On("Store", mock.Anything, mock.Anything).Return(errors.New("Unexpexted Error")).Once()
			},
		},
		{
			name: "success",
			arg:  sample.NewCategory(),
			prepare: func(categoryRepo *mocks.CategoryRepository) {
				categoryRepo.On("Store", mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			categoryRepo := new(mocks.CategoryRepository)
			tc.prepare(categoryRepo)
			u := usecases.NewCategoryUsecase(categoryRepo, time.Second*2)
			fmt.Println(tc.arg)
			err := u.Create(context.TODO(), tc.arg)
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			categoryRepo.AssertExpectations(t)
		})
	}
}

func Test_CategoryUsecase_Update(t *testing.T) {
	testCases := []struct {
		name          string
		arg           *domain.Category
		expectedErr   bool
		prepare       func(categoryRepo *mocks.CategoryRepository)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name:        "failure_find_category_by_id_returns_error",
			arg:         sample.NewCategory(),
			expectedErr: true,
			prepare: func(categoryRepo *mocks.CategoryRepository) {
				categoryRepo.On("FindByID", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:        "failure_update_category_returns_error",
			arg:         sample.NewCategory(),
			expectedErr: true,
			prepare: func(categoryRepo *mocks.CategoryRepository) {
				category := sample.NewCategory()
				categoryRepo.On("FindByID", mock.Anything, mock.AnythingOfType("string")).Return(category, nil).Once()
				categoryRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("Unexpexted Error")).Once()
			},
		},
		{
			name: "success",
			arg:  sample.NewCategory(),
			prepare: func(categoryRepo *mocks.CategoryRepository) {
				category := sample.NewCategory()
				categoryRepo.On("FindByID", mock.Anything, mock.AnythingOfType("string")).Return(category, nil).Once()
				categoryRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			categoryRepo := new(mocks.CategoryRepository)
			tc.prepare(categoryRepo)
			u := usecases.NewCategoryUsecase(categoryRepo, time.Second*2)
			err := u.Update(context.TODO(), tc.arg)
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			categoryRepo.AssertExpectations(t)
		})
	}
}
