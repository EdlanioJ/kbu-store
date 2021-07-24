package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/EdlanioJ/kbu-store/app/domain"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/grpc/pb"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/grpc/service"
	"github.com/EdlanioJ/kbu-store/app/utils/mocks"
	"github.com/golang/protobuf/ptypes/empty"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func getStore() *domain.Store {
	store := &domain.Store{
		Base: domain.Base{
			ID:        uuid.NewV4().String(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:        "Store 001",
		Description: "store description 001",
		Status:      domain.StoreStatusPending,
		UserID:      uuid.NewV4().String(),
		AccountID:   uuid.NewV4().String(),
		CategoryID:  uuid.NewV4().String(),
		Position: domain.Position{
			Lat: -8.8368200,
			Lng: 13.2343200,
		},
	}

	return store
}

func Test_StoreGrpcService_Create(t *testing.T) {
	a := &pb.CreateStoreRequest{
		Name:        "store 001",
		Description: "store description 001",
		CategoryID:  uuid.NewV4().String(),
		ExternalID:  uuid.NewV4().String(),
		Tags:        []string{"tag001", "tag002"},
		Latitude:    -8.8368200,
		Longitude:   13.2343200,
	}
	testCases := []struct {
		name          string
		arg           *pb.CreateStoreRequest
		builtSts      func(storeUsecase *mocks.StoreUsecase)
		checkResponse func(t *testing.T, res *empty.Empty, err error)
	}{
		{
			name: "should fail if usecase returns error",
			arg:  a,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Store",
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
			name: "should succeed",
			arg:  a,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Store",
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

func Test_StoreGrpcService_Get(t *testing.T) {
	a := &pb.StoreRequest{
		Id: uuid.NewV4().String(),
	}
	store := getStore()
	testCases := []struct {
		name          string
		arg           *pb.StoreRequest
		builtSts      func(storeUsecase *mocks.StoreUsecase)
		checkResponse func(t *testing.T, res *pb.Store, err error)
	}{
		{
			name: "should fail if id is empty",
			arg: &pb.StoreRequest{
				Id: "",
			},
			builtSts: func(_ *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, res *pb.Store, err error) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "should fail if id is not a valid uuidv4",
			arg: &pb.StoreRequest{
				Id: "invalid_id",
			},
			builtSts: func(_ *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, res *pb.Store, err error) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "should fail if usecase returns error",
			arg:  a,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.
					On("Get", mock.Anything, a.Id).
					Return(nil, errors.New("Unexpected Error"))
			},
			checkResponse: func(t *testing.T, res *pb.Store, err error) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "should succeed",
			arg:  a,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.
					On("Get", mock.Anything, a.Id).
					Return(store, nil)
			},
			checkResponse: func(t *testing.T, res *pb.Store, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			},
		},
	}

	for _, tc := range testCases {
		usecase := new(mocks.StoreUsecase)
		tc.builtSts(usecase)
		s := service.NewStoreServer(usecase)
		res, err := s.Get(context.TODO(), tc.arg)
		tc.checkResponse(t, res, err)
	}
}

func Test_StoreGrpcService_List(t *testing.T) {
	a := &pb.ListStoreRequest{
		Page:  1,
		Limit: 10,
		Sort:  "created_at",
	}
	store := getStore()
	testCases := []struct {
		name          string
		arg           *pb.ListStoreRequest
		builtSts      func(storeUsecase *mocks.StoreUsecase)
		checkResponse func(t *testing.T, res *pb.ListStoreResponse, err error)
	}{
		{
			name: "should fails if usecase returns error",
			arg:  a,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.
					On("Index", mock.Anything, a.Sort, int(a.Limit), int(a.Page)).
					Return(nil, int64(0), errors.New("Unexpected Error"))
			},
			checkResponse: func(t *testing.T, res *pb.ListStoreResponse, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "should succeed",
			arg:  a,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				stores := make(domain.Stores, 0)

				stores = append(stores, store)
				storeUsecase.
					On("Index", mock.Anything, a.Sort, int(a.Limit), int(a.Page)).
					Return(stores, int64(1), nil)
			},
			checkResponse: func(t *testing.T, res *pb.ListStoreResponse, err error) {
				assert.NotNil(t, res)
				assert.Equal(t, res.Total, int64(1))
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		usecase := new(mocks.StoreUsecase)
		tc.builtSts(usecase)
		s := service.NewStoreServer(usecase)
		res, err := s.List(context.TODO(), tc.arg)
		tc.checkResponse(t, res, err)
	}
}

func Test_StoreGrpcService_Activate(t *testing.T) {
	a := &pb.StoreRequest{
		Id: uuid.NewV4().String(),
	}
	testCases := []struct {
		name          string
		arg           *pb.StoreRequest
		builtSts      func(storeUsecase *mocks.StoreUsecase)
		checkResponse func(t *testing.T, res *empty.Empty, err error)
	}{
		{
			name: "should fail if id is empty",
			arg: &pb.StoreRequest{
				Id: "",
			},
			builtSts: func(storeUsecase *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, res *empty.Empty, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "should fail if invalid id",
			arg: &pb.StoreRequest{
				Id: "invaalid_id",
			},
			builtSts: func(storeUsecase *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, res *empty.Empty, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "should fail if usecase returns error",
			arg:  a,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.
					On("Active", mock.Anything, a.GetId()).
					Return(errors.New("Unexpected Error"))
			},
			checkResponse: func(t *testing.T, res *empty.Empty, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "should succeed",
			arg:  a,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.
					On("Active", mock.Anything, a.GetId()).
					Return(nil)
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
		res, err := s.Activate(context.TODO(), tc.arg)
		tc.checkResponse(t, res, err)
	}
}

func Test_StoreGrpcService_Block(t *testing.T) {
	a := &pb.StoreRequest{
		Id: uuid.NewV4().String(),
	}
	testCases := []struct {
		name          string
		arg           *pb.StoreRequest
		builtSts      func(storeUsecase *mocks.StoreUsecase)
		checkResponse func(t *testing.T, res *empty.Empty, err error)
	}{
		{
			name: "should fail if id is empty",
			arg: &pb.StoreRequest{
				Id: "",
			},
			builtSts: func(storeUsecase *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, res *empty.Empty, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "should fail if invalid id",
			arg: &pb.StoreRequest{
				Id: "invaalid_id",
			},
			builtSts: func(storeUsecase *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, res *empty.Empty, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "should fail if usecase returns error",
			arg:  a,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.
					On("Block", mock.Anything, a.GetId()).
					Return(errors.New("Unexpected Error"))
			},
			checkResponse: func(t *testing.T, res *empty.Empty, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "should succeed",
			arg:  a,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.
					On("Block", mock.Anything, a.GetId()).
					Return(nil)
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
		res, err := s.Block(context.TODO(), tc.arg)
		tc.checkResponse(t, res, err)
	}
}

func Test_StoreGrpcService_Disable(t *testing.T) {
	a := &pb.StoreRequest{
		Id: uuid.NewV4().String(),
	}
	testCases := []struct {
		name          string
		arg           *pb.StoreRequest
		builtSts      func(storeUsecase *mocks.StoreUsecase)
		checkResponse func(t *testing.T, res *empty.Empty, err error)
	}{
		{
			name: "should fail if id is empty",
			arg: &pb.StoreRequest{
				Id: "",
			},
			builtSts: func(storeUsecase *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, res *empty.Empty, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "should fail if invalid id",
			arg: &pb.StoreRequest{
				Id: "invaalid_id",
			},
			builtSts: func(storeUsecase *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, res *empty.Empty, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "should fail if usecase returns error",
			arg:  a,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.
					On("Disable", mock.Anything, a.GetId()).
					Return(errors.New("Unexpected Error"))
			},
			checkResponse: func(t *testing.T, res *empty.Empty, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "should succeed",
			arg:  a,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.
					On("Disable", mock.Anything, a.GetId()).
					Return(nil)
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
		res, err := s.Disable(context.TODO(), tc.arg)
		tc.checkResponse(t, res, err)
	}
}

func Test_StoreGrpcService_Update(t *testing.T) {
	a := &pb.UpdateStoreRequest{
		ID:          uuid.NewV4().String(),
		Name:        "store 001",
		Description: "store description 001",
		CategoryID:  uuid.NewV4().String(),
		Tags:        []string{"tag001", "tag002"},
		Latitude:    -8.8368200,
		Longitude:   13.2343200,
	}
	testCases := []struct {
		name          string
		arg           *pb.UpdateStoreRequest
		builtSts      func(storeUsecase *mocks.StoreUsecase)
		checkResponse func(t *testing.T, res *empty.Empty, err error)
	}{
		{
			name: "should fail if store id request param is empty",
			arg: &pb.UpdateStoreRequest{
				ID: "",
			},
			builtSts: func(storeUsecase *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, res *empty.Empty, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "should fail if store id request param is not uuidv4",
			arg: &pb.UpdateStoreRequest{
				ID: "invalid_id",
			},
			builtSts: func(storeUsecase *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, res *empty.Empty, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "should fail if category id request param is not uuidv4",
			arg: &pb.UpdateStoreRequest{
				ID:         uuid.NewV4().String(),
				CategoryID: "invaild_id",
			},
			builtSts: func(storeUsecase *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, res *empty.Empty, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "should fail if latitude request param is not valid",
			arg: &pb.UpdateStoreRequest{
				ID:       uuid.NewV4().String(),
				Latitude: 12345,
			},
			builtSts: func(storeUsecase *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, res *empty.Empty, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "should fail if longitude request param is not valid",
			arg: &pb.UpdateStoreRequest{
				ID:        uuid.NewV4().String(),
				Longitude: 12345,
			},
			builtSts: func(storeUsecase *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, res *empty.Empty, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "should fail if usecase returns error",
			arg:  a,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Update", mock.Anything, mock.Anything).Return(errors.New("Unexpected Error"))
			},
			checkResponse: func(t *testing.T, res *empty.Empty, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},

		{
			name: "should succeed",
			arg:  a,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Update", mock.Anything, mock.Anything).Return(nil)
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
		res, err := s.Update(context.TODO(), tc.arg)
		tc.checkResponse(t, res, err)
	}
}

func Test_StoreGrpcService_Delete(t *testing.T) {
	a := &pb.StoreRequest{
		Id: uuid.NewV4().String(),
	}
	testCases := []struct {
		name          string
		arg           *pb.StoreRequest
		builtSts      func(storeUsecase *mocks.StoreUsecase)
		checkResponse func(t *testing.T, res *empty.Empty, err error)
	}{
		{
			name: "should fail if invalid id",
			arg: &pb.StoreRequest{
				Id: "invaalid_id",
			},
			builtSts: func(storeUsecase *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, res *empty.Empty, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "should fail if usecase returns error",
			arg:  a,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.
					On("Delete", mock.Anything, a.GetId()).
					Return(errors.New("Unexpected Error"))
			},
			checkResponse: func(t *testing.T, res *empty.Empty, err error) {
				assert.Nil(t, res)
				assert.Error(t, err)
			},
		},
		{
			name: "should succeed",
			arg:  a,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.
					On("Delete", mock.Anything, a.GetId()).
					Return(nil)
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
		res, err := s.Delete(context.TODO(), tc.arg)
		tc.checkResponse(t, res, err)
	}
}
