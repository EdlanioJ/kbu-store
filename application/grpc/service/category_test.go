package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/EdlanioJ/kbu-store/application/grpc/pb"
	"github.com/EdlanioJ/kbu-store/application/grpc/service"
	"github.com/EdlanioJ/kbu-store/domain/mocks"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CategoryGrpcService_Create(t *testing.T) {
	a := &pb.CreateRequest{
		Name: "New Category",
	}
	testCases := []struct {
		name          string
		arg           *pb.CreateRequest
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
