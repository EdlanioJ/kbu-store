// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/EdlanioJ/kbu-store/domain"
	mock "github.com/stretchr/testify/mock"
)

// FetchStoreByTypeUsecase is an autogenerated mock type for the FetchStoreByTypeUsecase type
type FetchStoreByTypeUsecase struct {
	mock.Mock
}

// Exec provides a mock function with given fields: ctx, typeID, sort, limit, page
func (_m *FetchStoreByTypeUsecase) Exec(ctx context.Context, typeID string, sort string, limit int, page int) ([]*domain.Store, int64, error) {
	ret := _m.Called(ctx, typeID, sort, limit, page)

	var r0 []*domain.Store
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int, int) []*domain.Store); ok {
		r0 = rf(ctx, typeID, sort, limit, page)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Store)
		}
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(context.Context, string, string, int, int) int64); ok {
		r1 = rf(ctx, typeID, sort, limit, page)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, string, int, int) error); ok {
		r2 = rf(ctx, typeID, sort, limit, page)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
