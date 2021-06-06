package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/EdlanioJ/kbu-store/category/usecase"
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/domain/mocks"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_GetCategoryByStatusUsecase(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		is := assert.New(t)
		getCategoryByStatusRepo := new(mocks.GetCategoryByStautsRepository)

		status := domain.CategoryStatusPending
		id := uuid.NewV4().String()

		getCategoryByStatusRepo.On("Exec", mock.Anything, id, status).Return(nil, errors.New("Unexpexted Error")).Once()
		u := usecase.NewGetCategoryByStatus(getCategoryByStatusRepo, time.Second*2)

		res, err := u.Exec(context.TODO(), id, status)

		is.Nil(res)
		is.Error(err)
		getCategoryByStatusRepo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		is := assert.New(t)
		category := setupStore()
		getCategoryByStatusRepo := new(mocks.GetCategoryByStautsRepository)

		status := domain.CategoryStatusPending
		id := uuid.NewV4().String()

		getCategoryByStatusRepo.On("Exec", mock.Anything, id, status).Return(category, nil).Once()
		u := usecase.NewGetCategoryByStatus(getCategoryByStatusRepo, time.Second*2)

		res, err := u.Exec(context.TODO(), id, status)

		is.NotNil(res)
		is.NoError(err)
		getCategoryByStatusRepo.AssertExpectations(t)
	})
}
