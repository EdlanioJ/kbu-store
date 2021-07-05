package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/EdlanioJ/kbu-store/application/grpc/pb"
	"github.com/EdlanioJ/kbu-store/application/grpc/service"
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/domain/mocks"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_TagGrpcService_GetAll(t *testing.T) {
	a := &pb.TagListRequest{
		Page:  1,
		Limit: 10,
		Sort:  "created_at",
	}
	testCases := []struct {
		name          string
		arg           *pb.TagListRequest
		builtSts      func(tagUsecase *mocks.TagUsecase)
		checkResponse func(t *testing.T, res *pb.TagListResponse, err error)
	}{
		{
			name: "fail",
			arg:  a,
			builtSts: func(tagUsecase *mocks.TagUsecase) {
				tagUsecase.On("GetAll", mock.Anything, a.Sort, int(a.Page), int(a.Limit)).Return(nil, int64(0), errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, res *pb.TagListResponse, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			arg:  a,
			builtSts: func(tagUsecase *mocks.TagUsecase) {
				tag := &domain.Tag{
					Name:  "tag001",
					Count: 10,
				}
				tags := make([]*domain.Tag, 0)
				tags = append(tags, tag)
				tagUsecase.On("GetAll", mock.Anything, a.Sort, int(a.Page), int(a.Limit)).Return(tags, int64(1), nil).Once()
			},
			checkResponse: func(t *testing.T, res *pb.TagListResponse, err error) {
				assert.NotNil(t, res)
				assert.Equal(t, res.Total, int64(1))
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			usecase := new(mocks.TagUsecase)
			tc.builtSts(usecase)
			s := service.NewTagServer(usecase)
			res, err := s.GetAll(context.TODO(), tc.arg)
			tc.checkResponse(t, res, err)
		})
	}
}

func Test_TagGrpcService_GetAllByCategory(t *testing.T) {
	a := &pb.TagListByCategoryRequest{
		CategoryId: uuid.NewV4().String(),
		Page:       1,
		Limit:      10,
		Sort:       "created_at",
	}

	testCases := []struct {
		name          string
		arg           *pb.TagListByCategoryRequest
		builtSts      func(tagUsecase *mocks.TagUsecase)
		checkResponse func(t *testing.T, res *pb.TagListResponse, err error)
	}{
		{
			name: "fail on category id validation",
			arg: &pb.TagListByCategoryRequest{
				CategoryId: "invalid_id",
			},
			builtSts: func(_ *mocks.TagUsecase) {},
			checkResponse: func(t *testing.T, res *pb.TagListResponse, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "fail on usecase",
			arg:  a,
			builtSts: func(tagUsecase *mocks.TagUsecase) {
				tagUsecase.On("GetAllByCategory", mock.Anything, a.CategoryId, a.Sort, int(a.Page), int(a.Limit)).Return(nil, int64(0), errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, res *pb.TagListResponse, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			arg:  a,
			builtSts: func(tagUsecase *mocks.TagUsecase) {
				tag := &domain.Tag{
					Name:  "tag001",
					Count: 10,
				}
				tags := make([]*domain.Tag, 0)
				tags = append(tags, tag)
				tagUsecase.On("GetAllByCategory", mock.Anything, a.CategoryId, a.Sort, int(a.Page), int(a.Limit)).Return(tags, int64(1), nil).Once()
			},
			checkResponse: func(t *testing.T, res *pb.TagListResponse, err error) {
				assert.NotNil(t, res)
				assert.Equal(t, res.Total, int64(1))
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			usecase := new(mocks.TagUsecase)
			tc.builtSts(usecase)
			s := service.NewTagServer(usecase)
			res, err := s.GetAllByCategory(context.TODO(), tc.arg)
			tc.checkResponse(t, res, err)
		})
	}
}
