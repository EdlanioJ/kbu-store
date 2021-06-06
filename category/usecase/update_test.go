package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/EdlanioJ/kbu-store/category/usecase"
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_UpdateCategoryUsecase(t *testing.T) {
	t.Run("fail on get store type", func(t *testing.T) {
		is := assert.New(t)
		category := setupStore()

		getCategoryByIDRepo := new(mocks.GetCategoryByIDRepository)
		getCategoryByIDRepo.On("Exec", mock.Anything, category.ID).Return(nil, errors.New("Unexpexted Error"))
		u := usecase.NewUpdateCategory(nil, getCategoryByIDRepo, time.Second*2)
		err := u.Exec(context.TODO(), category)

		is.Error(err)
		getCategoryByIDRepo.AssertExpectations(t)
	})

	t.Run("fail on get store type returns empty", func(t *testing.T) {
		is := assert.New(t)
		category := setupStore()

		getCategoryByIDRepo := new(mocks.GetCategoryByIDRepository)
		getCategoryByIDRepo.On("Exec", mock.Anything, category.ID).Return(&domain.Category{}, nil)
		u := usecase.NewUpdateCategory(nil, getCategoryByIDRepo, time.Second*2)
		err := u.Exec(context.TODO(), category)

		is.Error(err)
		is.EqualError(err, domain.ErrNotFound.Error())
		getCategoryByIDRepo.AssertExpectations(t)
	})

	t.Run("fail on store type's repo update", func(t *testing.T) {
		is := assert.New(t)
		category := setupStore()

		getCategoryByIDRepo := new(mocks.GetCategoryByIDRepository)
		updateCategoryRepo := new(mocks.UpdateCategoryRepository)

		getCategoryByIDRepo.On("Exec", mock.Anything, category.ID).Return(category, nil)
		updateCategoryRepo.On("Exec", mock.Anything, mock.Anything).Return(errors.New("Unexpexted Error"))
		u := usecase.NewUpdateCategory(updateCategoryRepo, getCategoryByIDRepo, time.Second*2)
		err := u.Exec(context.TODO(), category)

		is.Error(err)
		getCategoryByIDRepo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		is := assert.New(t)
		category := setupStore()

		getCategoryByIDRepo := new(mocks.GetCategoryByIDRepository)
		updateCategoryRepo := new(mocks.UpdateCategoryRepository)

		getCategoryByIDRepo.On("Exec", mock.Anything, category.ID).Return(category, nil)
		updateCategoryRepo.On("Exec", mock.Anything, mock.Anything).Return(nil)
		u := usecase.NewUpdateCategory(updateCategoryRepo, getCategoryByIDRepo, time.Second*2)
		err := u.Exec(context.TODO(), category)

		is.NoError(err)
		getCategoryByIDRepo.AssertExpectations(t)
	})
}
