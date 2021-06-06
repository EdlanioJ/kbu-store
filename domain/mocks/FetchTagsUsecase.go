// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/EdlanioJ/kbu-store/domain"
	mock "github.com/stretchr/testify/mock"
)

// FetchTagsUsecase is an autogenerated mock type for the FetchTagsUsecase type
type FetchTagsUsecase struct {
	mock.Mock
}

// Exec provides a mock function with given fields: ctx, sort, page, limit
func (_m *FetchTagsUsecase) Exec(ctx context.Context, sort string, page int, limit int) ([]*domain.Tag, int64, error) {
	ret := _m.Called(ctx, sort, page, limit)

	var r0 []*domain.Tag
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) []*domain.Tag); ok {
		r0 = rf(ctx, sort, page, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Tag)
		}
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(context.Context, string, int, int) int64); ok {
		r1 = rf(ctx, sort, page, limit)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, int, int) error); ok {
		r2 = rf(ctx, sort, page, limit)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
