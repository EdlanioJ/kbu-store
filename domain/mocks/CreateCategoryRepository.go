// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/EdlanioJ/kbu-store/domain"
	mock "github.com/stretchr/testify/mock"
)

// CreateCategoryRepository is an autogenerated mock type for the CreateCategoryRepository type
type CreateCategoryRepository struct {
	mock.Mock
}

// Add provides a mock function with given fields: ctx, Category
func (_m *CreateCategoryRepository) Add(ctx context.Context, Category *domain.Category) error {
	ret := _m.Called(ctx, Category)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Category) error); ok {
		r0 = rf(ctx, Category)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
