// Code generated by mockery v2.23.2. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	context "context"

	model "github.com/sitename/sitename/model"
	mock "github.com/stretchr/testify/mock"
)

// CategoryStore is an autogenerated mock type for the CategoryStore type
type CategoryStore struct {
	mock.Mock
}

// FilterByOption provides a mock function with given fields: option
func (_m *CategoryStore) FilterByOption(option *model.CategoryFilterOption) ([]*model.Category, error) {
	ret := _m.Called(option)

	var r0 []*model.Category
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.CategoryFilterOption) ([]*model.Category, error)); ok {
		return rf(option)
	}
	if rf, ok := ret.Get(0).(func(*model.CategoryFilterOption) []*model.Category); ok {
		r0 = rf(option)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Category)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.CategoryFilterOption) error); ok {
		r1 = rf(option)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: ctx, categoryID, allowFromCache
func (_m *CategoryStore) Get(ctx context.Context, categoryID string, allowFromCache bool) (*model.Category, error) {
	ret := _m.Called(ctx, categoryID, allowFromCache)

	var r0 *model.Category
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, bool) (*model.Category, error)); ok {
		return rf(ctx, categoryID, allowFromCache)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, bool) *model.Category); ok {
		r0 = rf(ctx, categoryID, allowFromCache)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Category)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, bool) error); ok {
		r1 = rf(ctx, categoryID, allowFromCache)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByOption provides a mock function with given fields: option
func (_m *CategoryStore) GetByOption(option *model.CategoryFilterOption) (*model.Category, error) {
	ret := _m.Called(option)

	var r0 *model.Category
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.CategoryFilterOption) (*model.Category, error)); ok {
		return rf(option)
	}
	if rf, ok := ret.Get(0).(func(*model.CategoryFilterOption) *model.Category); ok {
		r0 = rf(option)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Category)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.CategoryFilterOption) error); ok {
		r1 = rf(option)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Upsert provides a mock function with given fields: category
func (_m *CategoryStore) Upsert(category *model.Category) (*model.Category, error) {
	ret := _m.Called(category)

	var r0 *model.Category
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.Category) (*model.Category, error)); ok {
		return rf(category)
	}
	if rf, ok := ret.Get(0).(func(*model.Category) *model.Category); ok {
		r0 = rf(category)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Category)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.Category) error); ok {
		r1 = rf(category)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewCategoryStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewCategoryStore creates a new instance of CategoryStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCategoryStore(t mockConstructorTestingTNewCategoryStore) *CategoryStore {
	mock := &CategoryStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}