package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/EdlanioJ/kbu-store/category/usecase"
	"github.com/EdlanioJ/kbu-store/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_TreateCategoryUsecase(t *testing.T) {
	t.Run("fail on new category", func(t *testing.T) {
		is := assert.New(t)

		name := ""
		u := usecase.NewCreateCategory(nil, time.Second*2)

		err := u.Add(context.TODO(), name)

		is.Error(err)
	})

	t.Run("fail on new categories's repo", func(t *testing.T) {
		is := assert.New(t)
		createCategoryRepo := new(mocks.CreateCategoryRepository)
		name := "Store 001"

		createCategoryRepo.On("Add", mock.Anything, mock.Anything).Return(errors.New("Unexpexted Error"))
		u := usecase.NewCreateCategory(createCategoryRepo, time.Second*2)

		err := u.Add(context.TODO(), name)

		is.Error(err)
	})
}
