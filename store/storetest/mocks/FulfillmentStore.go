// Code generated by mockery v2.23.2. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	boil "github.com/volatiletech/sqlboiler/v4/boil"

	model "github.com/sitename/sitename/model"

	model_helper "github.com/sitename/sitename/model_helper"
)

// FulfillmentStore is an autogenerated mock type for the FulfillmentStore type
type FulfillmentStore struct {
	mock.Mock
}

// Delete provides a mock function with given fields: tx, ids
func (_m *FulfillmentStore) Delete(tx boil.ContextTransactor, ids []string) error {
	ret := _m.Called(tx, ids)

	var r0 error
	if rf, ok := ret.Get(0).(func(boil.ContextTransactor, []string) error); ok {
		r0 = rf(tx, ids)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FilterByOption provides a mock function with given fields: option
func (_m *FulfillmentStore) FilterByOption(option model_helper.FulfillmentFilterOption) (model.FulfillmentSlice, error) {
	ret := _m.Called(option)

	var r0 model.FulfillmentSlice
	var r1 error
	if rf, ok := ret.Get(0).(func(model_helper.FulfillmentFilterOption) (model.FulfillmentSlice, error)); ok {
		return rf(option)
	}
	if rf, ok := ret.Get(0).(func(model_helper.FulfillmentFilterOption) model.FulfillmentSlice); ok {
		r0 = rf(option)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(model.FulfillmentSlice)
		}
	}

	if rf, ok := ret.Get(1).(func(model_helper.FulfillmentFilterOption) error); ok {
		r1 = rf(option)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: id
func (_m *FulfillmentStore) Get(id string) (*model.Fulfillment, error) {
	ret := _m.Called(id)

	var r0 *model.Fulfillment
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*model.Fulfillment, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *model.Fulfillment); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Fulfillment)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Upsert provides a mock function with given fields: tx, fulfillment
func (_m *FulfillmentStore) Upsert(tx boil.ContextTransactor, fulfillment model.Fulfillment) (*model.Fulfillment, error) {
	ret := _m.Called(tx, fulfillment)

	var r0 *model.Fulfillment
	var r1 error
	if rf, ok := ret.Get(0).(func(boil.ContextTransactor, model.Fulfillment) (*model.Fulfillment, error)); ok {
		return rf(tx, fulfillment)
	}
	if rf, ok := ret.Get(0).(func(boil.ContextTransactor, model.Fulfillment) *model.Fulfillment); ok {
		r0 = rf(tx, fulfillment)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Fulfillment)
		}
	}

	if rf, ok := ret.Get(1).(func(boil.ContextTransactor, model.Fulfillment) error); ok {
		r1 = rf(tx, fulfillment)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewFulfillmentStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewFulfillmentStore creates a new instance of FulfillmentStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewFulfillmentStore(t mockConstructorTestingTNewFulfillmentStore) *FulfillmentStore {
	mock := &FulfillmentStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
