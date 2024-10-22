// Code generated by mockery v2.23.2. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	model "github.com/sitename/sitename/model"
	mock "github.com/stretchr/testify/mock"

	model_helper "github.com/sitename/sitename/model_helper"
)

// AttributePageStore is an autogenerated mock type for the AttributePageStore type
type AttributePageStore struct {
	mock.Mock
}

// Get provides a mock function with given fields: pageID
func (_m *AttributePageStore) Get(pageID string) (*model.AttributePage, error) {
	ret := _m.Called(pageID)

	var r0 *model.AttributePage
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*model.AttributePage, error)); ok {
		return rf(pageID)
	}
	if rf, ok := ret.Get(0).(func(string) *model.AttributePage); ok {
		r0 = rf(pageID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AttributePage)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(pageID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByOption provides a mock function with given fields: option
func (_m *AttributePageStore) GetByOption(option model_helper.AttributePageFilterOption) (*model.AttributePage, error) {
	ret := _m.Called(option)

	var r0 *model.AttributePage
	var r1 error
	if rf, ok := ret.Get(0).(func(model_helper.AttributePageFilterOption) (*model.AttributePage, error)); ok {
		return rf(option)
	}
	if rf, ok := ret.Get(0).(func(model_helper.AttributePageFilterOption) *model.AttributePage); ok {
		r0 = rf(option)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AttributePage)
		}
	}

	if rf, ok := ret.Get(1).(func(model_helper.AttributePageFilterOption) error); ok {
		r1 = rf(option)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: page
func (_m *AttributePageStore) Save(page model.AttributePage) (*model.AttributePage, error) {
	ret := _m.Called(page)

	var r0 *model.AttributePage
	var r1 error
	if rf, ok := ret.Get(0).(func(model.AttributePage) (*model.AttributePage, error)); ok {
		return rf(page)
	}
	if rf, ok := ret.Get(0).(func(model.AttributePage) *model.AttributePage); ok {
		r0 = rf(page)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AttributePage)
		}
	}

	if rf, ok := ret.Get(1).(func(model.AttributePage) error); ok {
		r1 = rf(page)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewAttributePageStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewAttributePageStore creates a new instance of AttributePageStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAttributePageStore(t mockConstructorTestingTNewAttributePageStore) *AttributePageStore {
	mock := &AttributePageStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
