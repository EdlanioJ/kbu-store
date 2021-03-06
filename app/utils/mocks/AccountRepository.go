// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/EdlanioJ/kbu-store/app/domain"
	mock "github.com/stretchr/testify/mock"
)

// AccountRepository is an autogenerated mock type for the AccountRepository type
type AccountRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, id
func (_m *AccountRepository) Delete(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindByID provides a mock function with given fields: ctx, id
func (_m *AccountRepository) FindByID(ctx context.Context, id string) (*domain.Account, error) {
	ret := _m.Called(ctx, id)

	var r0 *domain.Account
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.Account); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Account)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store provides a mock function with given fields: ctx, account
func (_m *AccountRepository) Store(ctx context.Context, account *domain.Account) error {
	ret := _m.Called(ctx, account)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Account) error); ok {
		r0 = rf(ctx, account)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, account
func (_m *AccountRepository) Update(ctx context.Context, account *domain.Account) error {
	ret := _m.Called(ctx, account)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Account) error); ok {
		r0 = rf(ctx, account)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
