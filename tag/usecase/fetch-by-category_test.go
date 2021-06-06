package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/domain/mocks"
	"github.com/EdlanioJ/kbu-store/tag/usecase"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_FetchTagsByCategory(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		is := assert.New(t)
		fetchTagsByCategoryRepo := new(mocks.FetchTagsByCategoryRepository)

		categoryID := uuid.NewV4().String()
		fetchTagsByCategoryRepo.
			On("Exec", mock.Anything, categoryID, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
			Return(nil, int64(0), errors.New("Unexpexted Error")).
			Once()

		u := usecase.NewFetchTagsByCategory(fetchTagsByCategoryRepo, time.Second*2)

		list, total, err := u.Exec(context.TODO(), categoryID, "", 0, 0)
		is.Len(list, 0)
		is.Equal(total, int64(0))
		is.Error(err)

		fetchTagsByCategoryRepo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		is := assert.New(t)
		fetchTagsByCategoryRepo := new(mocks.FetchTagsByCategoryRepository)
		tag := testMock()

		categoryID := uuid.NewV4().String()
		tags := make([]*domain.Tag, 0)
		tags = append(tags, tag)
		fetchTagsByCategoryRepo.
			On("Exec", mock.Anything, categoryID, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
			Return(tags, int64(1), nil).
			Once()

		u := usecase.NewFetchTagsByCategory(fetchTagsByCategoryRepo, time.Second*2)

		list, total, err := u.Exec(context.TODO(), categoryID, "", 0, 0)
		is.Len(list, 1)
		is.Equal(total, int64(1))
		is.NoError(err)

		fetchTagsByCategoryRepo.AssertExpectations(t)
	})
}
