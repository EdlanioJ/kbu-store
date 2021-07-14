package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/EdlanioJ/kbu-store/application/grpc/pb"
	"github.com/EdlanioJ/kbu-store/application/grpc/service"
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/domain/mocks"
	"github.com/golang/protobuf/ptypes/empty"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func getCategory() *domain.Category {
	category, _ := domain.NewCategory("Store type 001")

	return category
}
func Test_CategoryGrpcService_Create(t *testing.T) {
	a := &pb.CreateCategoryRequest{
		Name: "New Category",
	}
	testCases := []struct {
		name          string
		arg           *pb.CreateCategoryRequest
		builtSts      func(categoryUsecase *mocks.CategoryUsecase)
		checkResponse func(t *testing.T, res *empty.Empty, err error)
	}{
		{
			name: "fail",
			arg:  a,
			builtSts: func(categoryUsecase *mocks.CategoryUsecase) {
				categoryUsecase.On("Create", mock.Anything, a.Name).Return(errors.New("Unexpected Error"))
			},
			checkResponse: func(t *testing.T, res *empty.Empty, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			arg:  a,
			builtSts: func(categoryUsecase *mocks.CategoryUsecase) {
				categoryUsecase.On("Create", mock.Anything, a.Name).Return(nil)
			},
			checkResponse: func(t *testing.T, res *empty.Empty, err error) {
				assert.NotNil(t, res)
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			usecase := new(mocks.CategoryUsecase)
			tc.builtSts(usecase)
			s := service.NewCategotyServer(usecase)
			res, err := s.Create(context.TODO(), tc.arg)
			tc.checkResponse(t, res, err)
		})
	}
}

func Test_CategoryGrpcService_GetById(t *testing.T) {
	a := &pb.CategoryRequest{
		Id: uuid.NewV4().String(),
	}
	testCases := []struct {
		name          string
		arg           *pb.CategoryRequest
		builtSts      func(categoryUsecase *mocks.CategoryUsecase)
		checkResponse func(t *testing.T, res *pb.Category, err error)
	}{
		{
			name: "fail on id validation",
			arg: &pb.CategoryRequest{
				Id: "invalid_id",
			},
			builtSts: func(_ *mocks.CategoryUsecase) {},
			checkResponse: func(t *testing.T, res *pb.Category, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "fail on usecase",
			arg:  a,
			builtSts: func(usecase *mocks.CategoryUsecase) {
				usecase.On("GetById", mock.Anything, a.Id).Return(nil, errors.New("Unexpected Error")).Once()
			},
			checkResponse: func(t *testing.T, res *pb.Category, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			arg:  a,
			builtSts: func(usecase *mocks.CategoryUsecase) {
				usecase.On("GetById", mock.Anything, a.Id).Return(getCategory(), nil).Once()
			},
			checkResponse: func(t *testing.T, res *pb.Category, err error) {
				assert.NotNil(t, res)
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			usecase := new(mocks.CategoryUsecase)
			tc.builtSts(usecase)
			s := service.NewCategotyServer(usecase)
			res, err := s.GetById(context.TODO(), tc.arg)
			tc.checkResponse(t, res, err)
		})
	}
}

func Test_CategoryGrpcService_GetByIdAndStatus(t *testing.T) {
	a := &pb.GetCategoryByIdAndStatusRequest{
		Id:     uuid.NewV4().String(),
		Status: pb.GetCategoryByIdAndStatusRequest_active,
	}
	testCases := []struct {
		name          string
		arg           *pb.GetCategoryByIdAndStatusRequest
		builtSts      func(categoryUsecase *mocks.CategoryUsecase)
		checkResponse func(t *testing.T, res *pb.Category, err error)
	}{
		{
			name: "fail on id validation",
			arg: &pb.GetCategoryByIdAndStatusRequest{
				Id:     "invalid_id",
				Status: pb.GetCategoryByIdAndStatusRequest_pending,
			},
			builtSts: func(_ *mocks.CategoryUsecase) {},
			checkResponse: func(t *testing.T, res *pb.Category, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "fail on usecase",
			arg:  a,
			builtSts: func(usecase *mocks.CategoryUsecase) {
				usecase.On("GetByIdAndStatus", mock.Anything, a.Id, a.Status.String()).Return(nil, errors.New("Unexpected Error")).Once()
			},
			checkResponse: func(t *testing.T, res *pb.Category, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			arg:  a,
			builtSts: func(usecase *mocks.CategoryUsecase) {
				usecase.On("GetByIdAndStatus", mock.Anything, a.Id, a.Status.String()).Return(getCategory(), nil).Once()
			},
			checkResponse: func(t *testing.T, res *pb.Category, err error) {
				assert.NotNil(t, res)
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			usecase := new(mocks.CategoryUsecase)
			tc.builtSts(usecase)
			s := service.NewCategotyServer(usecase)
			res, err := s.GetByIdAndStatus(context.TODO(), tc.arg)
			tc.checkResponse(t, res, err)
		})
	}
}

func Test_CategoryGrpcService_GetAll(t *testing.T) {
	a := &pb.GetAllCategoryRequest{
		Page:  1,
		Limit: 10,
		Sort:  "created_at",
	}
	testCases := []struct {
		name          string
		arg           *pb.GetAllCategoryRequest
		builtSts      func(categoryUsecase *mocks.CategoryUsecase)
		checkResponse func(t *testing.T, res *pb.ListCategoryResponse, err error)
	}{
		{
			name: "fail",
			arg:  a,
			builtSts: func(categoryUsecase *mocks.CategoryUsecase) {
				categoryUsecase.On("GetAll", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(nil, int64(0), errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, res *pb.ListCategoryResponse, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			arg:  a,
			builtSts: func(categoryUsecase *mocks.CategoryUsecase) {
				categories := make([]*domain.Category, 0)
				categories = append(categories, getCategory())
				categoryUsecase.On("GetAll", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(categories, int64(1), nil).Once()
			},
			checkResponse: func(t *testing.T, res *pb.ListCategoryResponse, err error) {
				assert.NotNil(t, res)
				assert.Equal(t, res.Total, int64(1))
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			usecase := new(mocks.CategoryUsecase)
			tc.builtSts(usecase)
			s := service.NewCategotyServer(usecase)
			res, err := s.GetAll(context.TODO(), tc.arg)
			tc.checkResponse(t, res, err)
		})
	}
}

func Test_CategoryGrpcService_GetAllByStatus(t *testing.T) {
	a := &pb.GetAllCategoryByStatusRequest{
		Status: pb.GetAllCategoryByStatusRequest_active,
		Page:   1,
		Limit:  10,
		Sort:   "created_at",
	}
	testCases := []struct {
		name          string
		arg           *pb.GetAllCategoryByStatusRequest
		builtSts      func(categoryUsecase *mocks.CategoryUsecase)
		checkResponse func(t *testing.T, res *pb.ListCategoryResponse, err error)
	}{
		{
			name: "fail",
			arg:  a,
			builtSts: func(categoryUsecase *mocks.CategoryUsecase) {
				categoryUsecase.On("GetAllByStatus", mock.Anything, a.Status.String(), a.Sort, mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(nil, int64(0), errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, res *pb.ListCategoryResponse, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			arg:  a,
			builtSts: func(categoryUsecase *mocks.CategoryUsecase) {
				categories := make([]*domain.Category, 0)
				categories = append(categories, getCategory())
				categoryUsecase.On("GetAllByStatus", mock.Anything, a.Status.String(), a.Sort, mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(categories, int64(1), nil).Once()
			},
			checkResponse: func(t *testing.T, res *pb.ListCategoryResponse, err error) {
				assert.NotNil(t, res)
				assert.Equal(t, res.Total, int64(1))
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			usecase := new(mocks.CategoryUsecase)
			tc.builtSts(usecase)
			s := service.NewCategotyServer(usecase)
			res, err := s.GetAllByStatus(context.TODO(), tc.arg)
			tc.checkResponse(t, res, err)
		})
	}
}

func Test_CategoryGrpcService_Activate(t *testing.T) {
	a := &pb.CategoryRequest{
		Id: uuid.NewV4().String(),
	}
	testCases := []struct {
		name          string
		arg           *pb.CategoryRequest
		builtSts      func(categoryUsecase *mocks.CategoryUsecase)
		checkResponse func(t *testing.T, res *empty.Empty, err error)
	}{
		{
			name: "fail on id validation",
			arg: &pb.CategoryRequest{
				Id: "invalid_id",
			},
			builtSts: func(_ *mocks.CategoryUsecase) {},
			checkResponse: func(t *testing.T, res *empty.Empty, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "fail on usecase",
			arg:  a,
			builtSts: func(categoryUsecase *mocks.CategoryUsecase) {
				categoryUsecase.On("Activate", mock.Anything, a.Id).Return(errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, res *empty.Empty, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			arg:  a,
			builtSts: func(categoryUsecase *mocks.CategoryUsecase) {
				categoryUsecase.On("Activate", mock.Anything, a.Id).Return(nil).Once()
			},
			checkResponse: func(t *testing.T, res *empty.Empty, err error) {
				assert.NotNil(t, res)
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			usecase := new(mocks.CategoryUsecase)
			tc.builtSts(usecase)
			s := service.NewCategotyServer(usecase)
			res, err := s.Activate(context.TODO(), tc.arg)
			tc.checkResponse(t, res, err)
		})
	}
}

func Test_CategoryGrpcService_Disable(t *testing.T) {
	a := &pb.CategoryRequest{
		Id: uuid.NewV4().String(),
	}
	testCases := []struct {
		name          string
		arg           *pb.CategoryRequest
		builtSts      func(categoryUsecase *mocks.CategoryUsecase)
		checkResponse func(t *testing.T, res *empty.Empty, err error)
	}{
		{
			name: "fail on id validation",
			arg: &pb.CategoryRequest{
				Id: "invalid_id",
			},
			builtSts: func(_ *mocks.CategoryUsecase) {},
			checkResponse: func(t *testing.T, res *empty.Empty, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "fail on usecase",
			arg:  a,
			builtSts: func(categoryUsecase *mocks.CategoryUsecase) {
				categoryUsecase.On("Disable", mock.Anything, a.Id).Return(errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, res *empty.Empty, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			arg:  a,
			builtSts: func(categoryUsecase *mocks.CategoryUsecase) {
				categoryUsecase.On("Disable", mock.Anything, a.Id).Return(nil).Once()
			},
			checkResponse: func(t *testing.T, res *empty.Empty, err error) {
				assert.NotNil(t, res)
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			usecase := new(mocks.CategoryUsecase)
			tc.builtSts(usecase)
			s := service.NewCategotyServer(usecase)
			res, err := s.Disable(context.TODO(), tc.arg)
			tc.checkResponse(t, res, err)
		})
	}
}
