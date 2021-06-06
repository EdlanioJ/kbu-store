// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/EdlanioJ/kbu-store/domain"
	mock "github.com/stretchr/testify/mock"
)

// FetchCategoryByStatusRepository is an autogenerated mock type for the FetchCategoryByStatusRepository type
type FetchCategoryByStatusRepository struct {
	mock.Mock
}

// Exec provides a mock function with given fields: ctx, status, sort, page, limit
func (_m *FetchCategoryByStatusRepository) Exec(ctx context.Context, status string, sort string, page int, limit int) ([]*domain.Category, int64, error) {
	ret := _m.Called(ctx, status, sort, page, limit)

	var r0 []*domain.Category
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int, int) []*domain.Category); ok {
		r0 = rf(ctx, status, sort, page, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Category)
		}
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(context.Context, string, string, int, int) int64); ok {
		r1 = rf(ctx, status, sort, page, limit)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, string, int, int) error); ok {
		r2 = rf(ctx, status, sort, page, limit)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
