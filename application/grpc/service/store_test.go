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

func Test_StoreGrpcService_Create(t *testing.T) {
	a := &pb.CreateStoreRequest{}
	testCases := []struct {
		name          string
		arg           *pb.CreateStoreRequest
		builtSts      func(storeUsecase *mocks.StoreUsecase)
		checkResponse func(t *testing.T, res *empty.Empty, err error)
	}{
		{
			name: "fail",
			arg:  a,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Create",
					mock.Anything,
					mock.AnythingOfType("string"),
					mock.AnythingOfType("string"),
					mock.AnythingOfType("string"),
					mock.AnythingOfType("string"),
					mock.Anything,
					mock.AnythingOfType("float64"),
					mock.AnythingOfType("float64"),
				).Return(errors.New("Unexpected Error"))
			},
			checkResponse: func(t *testing.T, res *empty.Empty, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "success",
			arg:  a,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Create",
					mock.Anything,
					a.Name,
					a.Description,
					a.CategoryID,
					a.ExternalID,
					a.Tags,
					a.Latitude,
					a.Longitude,
				).Return(nil)
			},
			checkResponse: func(t *testing.T, res *empty.Empty, err error) {
				assert.NotNil(t, res)
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		usecase := new(mocks.StoreUsecase)
		tc.builtSts(usecase)
		s := service.NewStoreServer(usecase)
		res, err := s.Create(context.TODO(), tc.arg)
		tc.checkResponse(t, res, err)
	}
}
