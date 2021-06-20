package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/domain/mocks"
	"github.com/EdlanioJ/kbu-store/store/usecase"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_TagsUsecase_GetAll(t *testing.T) {

	type args struct {
		page  int
		limit int
		sort  string
	}

	tag := &domain.Tag{
		Name:  "tag001",
		Count: 2,
	}
	testCases := []struct {
		name          string
		args          args
		builtSts      func(tagRepo *mocks.TagRepository)
		checkResponse func(t *testing.T, res []*domain.Tag, count int64, err error)
	}{
		{
			name: "fail",
			args: args{
				page:  0,
				limit: 0,
				sort:  "",
			},
			builtSts: func(tagRepo *mocks.TagRepository) {
				tagRepo.On("GetAll", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
					Return(nil, int64(0), errors.New("Unexpexted Error")).
					Once()
			},
			checkResponse: func(t *testing.T, res []*domain.Tag, count int64, err error) {
				assert.Len(t, res, 0)
				assert.Equal(t, count, int64(0))
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			args: args{
				page:  0,
				limit: 0,
				sort:  "",
			},
			builtSts: func(tagRepo *mocks.TagRepository) {
				tags := make([]*domain.Tag, 0)
				tags = append(tags, tag)

				tagRepo.On("GetAll", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
					Return(tags, int64(1), nil).
					Once()
			},
			checkResponse: func(t *testing.T, res []*domain.Tag, count int64, err error) {
				assert.Len(t, res, 1)
				assert.Equal(t, count, int64(1))
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tagsRepo := new(mocks.TagRepository)
			tc.builtSts(tagsRepo)
			u := usecase.NewtagUsecase(tagsRepo, time.Second*2)
			res, count, err := u.GetAll(context.TODO(), tc.args.sort, tc.args.page, tc.args.limit)
			tc.checkResponse(t, res, count, err)
			tagsRepo.AssertExpectations(t)
		})
	}
}

func Test_TagsUsecase_GetAllByCategory(t *testing.T) {

	type args struct {
		categoryId string
		page       int
		limit      int
		sort       string
	}

	tag := &domain.Tag{
		Name:  "tag001",
		Count: 2,
	}
	testCases := []struct {
		name          string
		args          args
		builtSts      func(tagRepo *mocks.TagRepository)
		checkResponse func(t *testing.T, res []*domain.Tag, count int64, err error)
	}{
		{
			name: "fail",
			args: args{
				categoryId: uuid.NewV4().String(),
				page:       0,
				limit:      0,
				sort:       "",
			},
			builtSts: func(tagRepo *mocks.TagRepository) {
				tagRepo.On("GetAllByCategory", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
					Return(nil, int64(0), errors.New("Unexpexted Error")).
					Once()
			},
			checkResponse: func(t *testing.T, res []*domain.Tag, count int64, err error) {
				assert.Len(t, res, 0)
				assert.Equal(t, count, int64(0))
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			args: args{
				page:  0,
				limit: 0,
				sort:  "",
			},
			builtSts: func(tagRepo *mocks.TagRepository) {
				tags := make([]*domain.Tag, 0)
				tags = append(tags, tag)

				tagRepo.On("GetAllByCategory", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
					Return(tags, int64(1), nil).
					Once()
			},
			checkResponse: func(t *testing.T, res []*domain.Tag, count int64, err error) {
				assert.Len(t, res, 1)
				assert.Equal(t, count, int64(1))
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tagsRepo := new(mocks.TagRepository)
			tc.builtSts(tagsRepo)
			u := usecase.NewtagUsecase(tagsRepo, time.Second*2)
			res, count, err := u.GetAllByCategory(context.TODO(), tc.args.categoryId, tc.args.sort, tc.args.page, tc.args.limit)
			tc.checkResponse(t, res, count, err)
			tagsRepo.AssertExpectations(t)
		})
	}
}
