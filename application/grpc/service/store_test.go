package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/EdlanioJ/kbu-store/application/grpc/pb"
	"github.com/EdlanioJ/kbu-store/application/grpc/service"
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/domain/mocks"
	"github.com/golang/protobuf/ptypes/empty"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func getStore() *domain.Store {
	c := new(domain.Category)
	c.ID = uuid.NewV4().String()
	c.Name = "store 001"
	c.Status = domain.CategoryStatusActive

	store := &domain.Store{
		Base: domain.Base{
			ID:        uuid.NewV4().String(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:        "Store 001",
		Description: "store description 001",
		Status:      domain.StoreStatusPending,
		ExternalID:  uuid.NewV4().String(),
		AccountID:   uuid.NewV4().String(),
		Category:    c,
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
			name: "should succeed",
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

func Test_StoreGrpcService_GetById(t *testing.T) {
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
					On("GetById", mock.Anything, a.Id).
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
					On("GetById", mock.Anything, a.Id).
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
		res, err := s.GetById(context.TODO(), tc.arg)
		tc.checkResponse(t, res, err)
	}
}

func Test_StoreGrpcService_GetByIdAndOwner(t *testing.T) {
	a := &pb.GetStoreByIdAndOwnerRequest{
		ID:    uuid.NewV4().String(),
		Owner: uuid.NewV4().String(),
	}
	store := getStore()
	testCases := []struct {
		name          string
		arg           *pb.GetStoreByIdAndOwnerRequest
		builtSts      func(storeUsecase *mocks.StoreUsecase)
		checkResponse func(t *testing.T, res *pb.Store, err error)
	}{
		{
			name: "should fail if id is not a valid uuidv4",
			arg: &pb.GetStoreByIdAndOwnerRequest{
				ID: "invalid_id",
			},
			builtSts: func(_ *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, res *pb.Store, err error) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "should fail if owner is not a valid uuidv4",
			arg: &pb.GetStoreByIdAndOwnerRequest{
				ID:    uuid.NewV4().String(),
				Owner: "invalid_id",
			},
			builtSts: func(_ *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, res *pb.Store, err error) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "should fail if usecase fails",
			arg:  a,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.
					On("GetByIdAndOwner", mock.Anything, a.ID, a.Owner).
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
					On("GetByIdAndOwner", mock.Anything, a.ID, a.Owner).
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
		res, err := s.GetByIdAndOwner(context.TODO(), tc.arg)
		tc.checkResponse(t, res, err)
	}
}