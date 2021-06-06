package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/EdlanioJ/kbu-store/category/usecase"
	"github.com/EdlanioJ/kbu-store/domain/mocks"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_GetCategoryByIDUsecase(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		is := assert.New(t)
		getCategoryByIDRepo := new(mocks.GetCategoryByIDRepository)
		id := uuid.NewV4().String()

		getCategoryByIDRepo.On("Exec", mock.Anything, id).Return(nil, errors.New("Unexpexted Error")).Once()
		u := usecase.NewGetCategoryByID(getCategoryByIDRepo, time.Second*2)

		res, err := u.Exec(context.TODO(), id)

		is.Nil(res)
		is.Error(err)
		getCategoryByIDRepo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		is := assert.New(t)
		category := setupStore()
		getCategoryByIDRepo := new(mocks.GetCategoryByIDRepository)
		id := uuid.NewV4().String()

		getCategoryByIDRepo.On("Exec", mock.Anything, id).Return(category, nil).Once()
		u := usecase.NewGetCategoryByID(getCategoryByIDRepo, time.Second*2)

		res, err := u.Exec(context.TODO(), id)

		is.NotNil(res)
		is.NoError(err)
		getCategoryByIDRepo.AssertExpectations(t)
	})
}
