// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/EdlanioJ/kbu-store/domain"
	mock "github.com/stretchr/testify/mock"
)

// StoreUsecase is an autogenerated mock type for the StoreUsecase type
type StoreUsecase struct {
	mock.Mock
}

// Active provides a mock function with given fields: ctx, id
func (_m *StoreUsecase) Active(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Block provides a mock function with given fields: ctx, id
func (_m *StoreUsecase) Block(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Create provides a mock function with given fields: ctx, name, description, CategoryID, externalID, tags, lat, lng
func (_m *StoreUsecase) Create(ctx context.Context, name string, description string, CategoryID string, externalID string, tags []string, lat float64, lng float64) error {
	ret := _m.Called(ctx, name, description, CategoryID, externalID, tags, lat, lng)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, string, []string, float64, float64) error); ok {
		r0 = rf(ctx, name, description, CategoryID, externalID, tags, lat, lng)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, id
func (_m *StoreUsecase) Delete(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Disable provides a mock function with given fields: ctx, id
func (_m *StoreUsecase) Disable(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAll provides a mock function with given fields: ctx, sort, limit, page
func (_m *StoreUsecase) GetAll(ctx context.Context, sort string, limit int, page int) ([]*domain.Store, int64, error) {
	ret := _m.Called(ctx, sort, limit, page)

	var r0 []*domain.Store
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) []*domain.Store); ok {
		r0 = rf(ctx, sort, limit, page)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Store)
		}
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(context.Context, string, int, int) int64); ok {
		r1 = rf(ctx, sort, limit, page)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, int, int) error); ok {
		r2 = rf(ctx, sort, limit, page)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetAllByCategory provides a mock function with given fields: ctx, categoryID, sort, limit, page
func (_m *StoreUsecase) GetAllByCategory(ctx context.Context, categoryID string, sort string, limit int, page int) ([]*domain.Store, int64, error) {
	ret := _m.Called(ctx, categoryID, sort, limit, page)

	var r0 []*domain.Store
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int, int) []*domain.Store); ok {
		r0 = rf(ctx, categoryID, sort, limit, page)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Store)
		}
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(context.Context, string, string, int, int) int64); ok {
		r1 = rf(ctx, categoryID, sort, limit, page)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, string, int, int) error); ok {
		r2 = rf(ctx, categoryID, sort, limit, page)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetAllByCloseLocation provides a mock function with given fields: ctx, lat, lng, distance, status, limit, page, sort
func (_m *StoreUsecase) GetAllByCloseLocation(ctx context.Context, lat float64, lng float64, distance int, status string, limit int, page int, sort string) ([]*domain.Store, int64, error) {
	ret := _m.Called(ctx, lat, lng, distance, status, limit, page, sort)

	var r0 []*domain.Store
	if rf, ok := ret.Get(0).(func(context.Context, float64, float64, int, string, int, int, string) []*domain.Store); ok {
		r0 = rf(ctx, lat, lng, distance, status, limit, page, sort)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Store)
		}
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(context.Context, float64, float64, int, string, int, int, string) int64); ok {
		r1 = rf(ctx, lat, lng, distance, status, limit, page, sort)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, float64, float64, int, string, int, int, string) error); ok {
		r2 = rf(ctx, lat, lng, distance, status, limit, page, sort)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetAllByOwner provides a mock function with given fields: ctx, owner, sort, limit, page
func (_m *StoreUsecase) GetAllByOwner(ctx context.Context, owner string, sort string, limit int, page int) ([]*domain.Store, int64, error) {
	ret := _m.Called(ctx, owner, sort, limit, page)

	var r0 []*domain.Store
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int, int) []*domain.Store); ok {
		r0 = rf(ctx, owner, sort, limit, page)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Store)
		}
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(context.Context, string, string, int, int) int64); ok {
		r1 = rf(ctx, owner, sort, limit, page)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, string, int, int) error); ok {
		r2 = rf(ctx, owner, sort, limit, page)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetAllByStatus provides a mock function with given fields: ctx, status, sort, limit, page
func (_m *StoreUsecase) GetAllByStatus(ctx context.Context, status string, sort string, limit int, page int) ([]*domain.Store, int64, error) {
	ret := _m.Called(ctx, status, sort, limit, page)

	var r0 []*domain.Store
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int, int) []*domain.Store); ok {
		r0 = rf(ctx, status, sort, limit, page)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Store)
		}
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(context.Context, string, string, int, int) int64); ok {
		r1 = rf(ctx, status, sort, limit, page)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, string, int, int) error); ok {
		r2 = rf(ctx, status, sort, limit, page)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetAllByTags provides a mock function with given fields: ctx, tags, sort, limit, page
func (_m *StoreUsecase) GetAllByTags(ctx context.Context, tags []string, sort string, limit int, page int) ([]*domain.Store, int64, error) {
	ret := _m.Called(ctx, tags, sort, limit, page)

	var r0 []*domain.Store
	if rf, ok := ret.Get(0).(func(context.Context, []string, string, int, int) []*domain.Store); ok {
		r0 = rf(ctx, tags, sort, limit, page)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Store)
		}
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(context.Context, []string, string, int, int) int64); ok {
		r1 = rf(ctx, tags, sort, limit, page)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, []string, string, int, int) error); ok {
		r2 = rf(ctx, tags, sort, limit, page)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetById provides a mock function with given fields: ctx, id
func (_m *StoreUsecase) GetById(ctx context.Context, id string) (*domain.Store, error) {
	ret := _m.Called(ctx, id)

	var r0 *domain.Store
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.Store); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Store)
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

// GetByIdAndOwner provides a mock function with given fields: ctx, id, externalID
func (_m *StoreUsecase) GetByIdAndOwner(ctx context.Context, id string, externalID string) (*domain.Store, error) {
	ret := _m.Called(ctx, id, externalID)

	var r0 *domain.Store
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *domain.Store); ok {
		r0 = rf(ctx, id, externalID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Store)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, id, externalID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, store
func (_m *StoreUsecase) Update(ctx context.Context, store *domain.Store) error {
	ret := _m.Called(ctx, store)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Store) error); ok {
		r0 = rf(ctx, store)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
