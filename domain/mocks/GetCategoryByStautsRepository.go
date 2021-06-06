// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/EdlanioJ/kbu-store/domain"
	mock "github.com/stretchr/testify/mock"
)

// GetCategoryByStautsRepository is an autogenerated mock type for the GetCategoryByStautsRepository type
type GetCategoryByStautsRepository struct {
	mock.Mock
}

// Exec provides a mock function with given fields: ctx, id, status
func (_m *GetCategoryByStautsRepository) Exec(ctx context.Context, id string, status string) (*domain.Category, error) {
	ret := _m.Called(ctx, id, status)

	var r0 *domain.Category
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *domain.Category); ok {
		r0 = rf(ctx, id, status)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Category)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, id, status)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
