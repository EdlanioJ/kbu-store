package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/EdlanioJ/kbu-store/app/domain"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/grpc/pb"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/grpc/service"
	"github.com/EdlanioJ/kbu-store/app/utils/mocks"
	"github.com/EdlanioJ/kbu-store/app/utils/sample"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_StoreGrpcService_Create(t *testing.T) {
	t.Parallel()
	arg := sample.NewPBCreateStoreRequest()
	testCases := []struct {
		name        string
		arg         *pb.CreateStoreRequest
		prepare     func(storeUsecase *mocks.StoreUsecase)
		expectedErr bool
	}{
		{
			name:        "failure_usecase_returns_error",
			arg:         arg,
			expectedErr: true,
			prepare: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Store",
					mock.Anything,
					mock.Anything,
				).Return(errors.New("Unexpected Error"))
			},
		},
		{
			name:        "success",
			arg:         arg,
			expectedErr: false,
			prepare: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Store",
					mock.Anything,
					mock.Anything,
				).Return(nil)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			usecase := new(mocks.StoreUsecase)
			if tc.prepare != nil {
				tc.prepare(usecase)
			}
			s := service.NewStoreServer(usecase)
			res, err := s.Create(context.TODO(), tc.arg)
			if tc.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NotNil(t, res)
				assert.NoError(t, err)
			}
		})
	}
}

func Test_StoreGrpcService_Get(t *testing.T) {
	t.Parallel()
	arg := sample.NewPBStoreRequest()
	emptyId := sample.NewPBStoreRequest()
	emptyId.Id = ""
	invalidId := sample.NewPBStoreRequest()
	invalidId.Id = "invalid_id"
	store := sample.NewStore()
	testCases := []struct {
		name        string
		arg         *pb.StoreRequest
		prepare     func(storeUsecase *mocks.StoreUsecase)
		expectedErr bool
	}{
		{
			name:        "failure_empty_id",
			arg:         emptyId,
			expectedErr: true,
		},
		{
			name:        "failure_invalid_id",
			arg:         invalidId,
			expectedErr: true,
		},
		{
			name:        "failure_usecase_returns_error",
			arg:         arg,
			expectedErr: true,
			prepare: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.
					On("Get", mock.Anything, arg.Id).
					Return(nil, errors.New("Unexpected Error"))
			},
		},
		{
			name:        "success",
			arg:         arg,
			expectedErr: false,
			prepare: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.
					On("Get", mock.Anything, arg.Id).
					Return(store, nil)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			usecase := new(mocks.StoreUsecase)
			if tc.prepare != nil {
				tc.prepare(usecase)
			}
			s := service.NewStoreServer(usecase)
			res, err := s.Get(context.TODO(), tc.arg)
			if tc.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			}
		})
	}
}

func Test_StoreGrpcService_List(t *testing.T) {
	t.Parallel()
	arg := sample.NewPBListStoreRequest()
	store := sample.NewStore()
	testCases := []struct {
		name        string
		arg         *pb.ListStoreRequest
		prepare     func(storeUsecase *mocks.StoreUsecase)
		expectedErr bool
	}{
		{
			name:        "failure_usecase_returns_error",
			arg:         arg,
			expectedErr: true,
			prepare: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.
					On("Index", mock.Anything, arg.Sort, int(arg.Limit), int(arg.Page)).
					Return(nil, int64(0), errors.New("Unexpected Error"))
			},
		},
		{
			name:        "success",
			arg:         arg,
			expectedErr: false,
			prepare: func(storeUsecase *mocks.StoreUsecase) {
				stores := make(domain.Stores, 0)

				stores = append(stores, store)
				storeUsecase.
					On("Index", mock.Anything, arg.Sort, int(arg.Limit), int(arg.Page)).
					Return(stores, int64(1), nil)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			usecase := new(mocks.StoreUsecase)
			if tc.prepare != nil {
				tc.prepare(usecase)
			}
			s := service.NewStoreServer(usecase)
			res, err := s.List(context.TODO(), tc.arg)
			if tc.expectedErr {
				assert.Nil(t, res)
				assert.Error(t, err)
			} else {
				assert.NotNil(t, res)
				assert.Equal(t, res.Total, int64(1))
				assert.NoError(t, err)
			}
		})
	}
}

func Test_StoreGrpcService_Activate(t *testing.T) {
	t.Parallel()
	arg := sample.NewPBStoreRequest()
	emptyId := sample.NewPBStoreRequest()
	emptyId.Id = ""
	invalidId := sample.NewPBStoreRequest()
	invalidId.Id = "invalid_id"

	testCases := []struct {
		name        string
		arg         *pb.StoreRequest
		prepare     func(storeUsecase *mocks.StoreUsecase)
		expectedErr bool
	}{
		{
			name:        "failure_empty_id",
			arg:         emptyId,
			expectedErr: true,
		},
		{
			name:        "failure_invalid_id",
			arg:         invalidId,
			expectedErr: true,
		},
		{
			name:        "failure_usecase_returns_error",
			arg:         arg,
			expectedErr: true,
			prepare: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.
					On("Active", mock.Anything, arg.GetId()).
					Return(errors.New("Unexpected Error"))
			},
		},
		{
			name:        "success",
			arg:         arg,
			expectedErr: false,
			prepare: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.
					On("Active", mock.Anything, arg.GetId()).
					Return(nil)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			usecase := new(mocks.StoreUsecase)
			if tc.prepare != nil {
				tc.prepare(usecase)
			}
			s := service.NewStoreServer(usecase)
			res, err := s.Activate(context.TODO(), tc.arg)
			if tc.expectedErr {
				assert.Nil(t, res)
				assert.Error(t, err)
			} else {
				assert.NotNil(t, res)
				assert.NoError(t, err)
			}
		})
	}
}

func Test_StoreGrpcService_Block(t *testing.T) {
	t.Parallel()
	arg := sample.NewPBStoreRequest()
	emptyId := sample.NewPBStoreRequest()
	emptyId.Id = ""
	invalidId := sample.NewPBStoreRequest()
	invalidId.Id = "invalid_id"

	testCases := []struct {
		name        string
		arg         *pb.StoreRequest
		prepare     func(storeUsecase *mocks.StoreUsecase)
		expectedErr bool
	}{
		{
			name:        "failure_empty_id",
			arg:         emptyId,
			expectedErr: true,
		},
		{
			name:        "failure_invalid_id",
			arg:         invalidId,
			expectedErr: true,
		},
		{
			name:        "failure_usecase_returns_error",
			arg:         arg,
			expectedErr: true,
			prepare: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.
					On("Block", mock.Anything, arg.GetId()).
					Return(errors.New("Unexpected Error"))
			},
		},
		{
			name:        "success",
			arg:         arg,
			expectedErr: false,
			prepare: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.
					On("Block", mock.Anything, arg.GetId()).
					Return(nil)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			usecase := new(mocks.StoreUsecase)
			if tc.prepare != nil {
				tc.prepare(usecase)
			}
			s := service.NewStoreServer(usecase)
			res, err := s.Block(context.TODO(), tc.arg)
			if tc.expectedErr {
				assert.Nil(t, res)
				assert.Error(t, err)
			} else {
				assert.NotNil(t, res)
				assert.NoError(t, err)
			}
		})
	}
}

func Test_StoreGrpcService_Disable(t *testing.T) {
	t.Parallel()
	arg := sample.NewPBStoreRequest()
	emptyId := sample.NewPBStoreRequest()
	emptyId.Id = ""
	invalidId := sample.NewPBStoreRequest()
	invalidId.Id = "invalid_id"

	testCases := []struct {
		name        string
		arg         *pb.StoreRequest
		prepare     func(storeUsecase *mocks.StoreUsecase)
		expectedErr bool
	}{
		{
			name:        "failure_empty_id",
			arg:         emptyId,
			expectedErr: true,
		},
		{
			name:        "failure_invalid_id",
			arg:         invalidId,
			expectedErr: true,
		},
		{
			name:        "failure_usecase_returns_error",
			arg:         arg,
			expectedErr: true,
			prepare: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.
					On("Disable", mock.Anything, arg.GetId()).
					Return(errors.New("Unexpected Error"))
			},
		},
		{
			name:        "success",
			arg:         arg,
			expectedErr: false,
			prepare: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.
					On("Disable", mock.Anything, arg.GetId()).
					Return(nil)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			usecase := new(mocks.StoreUsecase)
			if tc.prepare != nil {
				tc.prepare(usecase)
			}
			s := service.NewStoreServer(usecase)
			res, err := s.Disable(context.TODO(), tc.arg)
			if tc.expectedErr {
				assert.Nil(t, res)
				assert.Error(t, err)
			} else {
				assert.NotNil(t, res)
				assert.NoError(t, err)
			}
		})
	}
}

func Test_StoreGrpcService_Update(t *testing.T) {
	t.Parallel()
	arg := sample.NewPBUpdateStoreRequest()
	emptyStoreId := sample.NewPBUpdateStoreRequest()
	emptyStoreId.ID = ""
	invalidStoreId := sample.NewPBUpdateStoreRequest()
	invalidStoreId.ID = "invalid_id"
	invalidCategoryId := sample.NewPBUpdateStoreRequest()
	invalidCategoryId.CategoryID = "invalid_id"
	invalidLatitude := sample.NewPBUpdateStoreRequest()
	invalidLatitude.Latitude = 12345
	invalidLongitude := sample.NewPBUpdateStoreRequest()
	invalidLongitude.Longitude = 12345

	testCases := []struct {
		name        string
		arg         *pb.UpdateStoreRequest
		prepare     func(storeUsecase *mocks.StoreUsecase)
		expectedErr bool
	}{
		{
			name:        "failure_usecase_returns_error",
			arg:         arg,
			expectedErr: true,
			prepare: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Update", mock.Anything, mock.Anything).Return(errors.New("Unexpected Error"))
			},
		},

		{
			name:        "success",
			arg:         arg,
			expectedErr: false,
			prepare: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			usecase := new(mocks.StoreUsecase)
			if tc.prepare != nil {
				tc.prepare(usecase)
			}
			s := service.NewStoreServer(usecase)
			res, err := s.Update(context.TODO(), tc.arg)
			if tc.expectedErr {
				assert.Nil(t, res)
				assert.Error(t, err)
			} else {
				assert.NotNil(t, res)
				assert.NoError(t, err)
			}
		})
	}
}

func Test_StoreGrpcService_Delete(t *testing.T) {
	t.Parallel()
	arg := sample.NewPBStoreRequest()
	emptyId := sample.NewPBStoreRequest()
	emptyId.Id = ""
	invalidId := sample.NewPBStoreRequest()
	invalidId.Id = "invalid_id"
	testCases := []struct {
		name          string
		arg           *pb.StoreRequest
		prepare       func(storeUsecase *mocks.StoreUsecase)
		expectedErr   bool
		checkResponse func(t *testing.T, res *empty.Empty, err error)
	}{
		{
			name:        "failure_empty_id",
			arg:         emptyId,
			expectedErr: true,
		},
		{
			name:        "failure_invalid_id",
			arg:         invalidId,
			expectedErr: true,
		},
		{
			name:        "failure_usecase_returns_error",
			arg:         arg,
			expectedErr: true,
			prepare: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.
					On("Delete", mock.Anything, arg.GetId()).
					Return(errors.New("Unexpected Error"))
			},
		},
		{
			name:        "success",
			arg:         arg,
			expectedErr: false,
			prepare: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.
					On("Delete", mock.Anything, arg.GetId()).
					Return(nil)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			usecase := new(mocks.StoreUsecase)
			if tc.prepare != nil {
				tc.prepare(usecase)
			}
			s := service.NewStoreServer(usecase)
			res, err := s.Delete(context.TODO(), tc.arg)
			if tc.expectedErr {
				assert.Nil(t, res)
				assert.Error(t, err)
			} else {
				assert.NotNil(t, res)
				assert.NoError(t, err)
			}
		})
	}
}
