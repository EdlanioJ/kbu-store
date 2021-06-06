package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/domain/mocks"
	"github.com/EdlanioJ/kbu-store/tag/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_FetchTags(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		is := assert.New(t)
		fetchTagsRepo := new(mocks.FetchTagsRepository)

		fetchTagsRepo.
			On("Exec", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
			Return(nil, int64(0), errors.New("Unexpexted Error")).
			Once()

		u := usecase.NewFetchTags(fetchTagsRepo, time.Second*2)

		list, total, err := u.Exec(context.TODO(), "", 0, 0)
		is.Len(list, 0)
		is.Equal(total, int64(0))
		is.Error(err)

		fetchTagsRepo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		is := assert.New(t)
		fetchTagsRepo := new(mocks.FetchTagsRepository)
		tag := testMock()

		tags := make([]*domain.Tag, 0)
		tags = append(tags, tag)
		fetchTagsRepo.
			On("Exec", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
			Return(tags, int64(1), nil).
			Once()

		u := usecase.NewFetchTags(fetchTagsRepo, time.Second*2)

		list, total, err := u.Exec(context.TODO(), "", 0, 0)
		is.Len(list, 1)
		is.Equal(total, int64(1))
		is.NoError(err)

		fetchTagsRepo.AssertExpectations(t)
	})
}
