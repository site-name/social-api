// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	goprices "github.com/site-name/go-prices"
	measurement "github.com/sitename/sitename/modules/measurement"
	mock "github.com/stretchr/testify/mock"

	shipping "github.com/sitename/sitename/model/shipping"
)

// ShippingMethodStore is an autogenerated mock type for the ShippingMethodStore type
type ShippingMethodStore struct {
	mock.Mock
}

// ApplicableShippingMethods provides a mock function with given fields: price, channelID, weight, countryCode, productIDs
func (_m *ShippingMethodStore) ApplicableShippingMethods(price *goprices.Money, channelID string, weight *measurement.Weight, countryCode string, productIDs []string) ([]*shipping.ShippingMethod, error) {
	ret := _m.Called(price, channelID, weight, countryCode, productIDs)

	var r0 []*shipping.ShippingMethod
	if rf, ok := ret.Get(0).(func(*goprices.Money, string, *measurement.Weight, string, []string) []*shipping.ShippingMethod); ok {
		r0 = rf(price, channelID, weight, countryCode, productIDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*shipping.ShippingMethod)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*goprices.Money, string, *measurement.Weight, string, []string) error); ok {
		r1 = rf(price, channelID, weight, countryCode, productIDs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateIndexesIfNotExists provides a mock function with given fields:
func (_m *ShippingMethodStore) CreateIndexesIfNotExists() {
	_m.Called()
}

// Get provides a mock function with given fields: methodID
func (_m *ShippingMethodStore) Get(methodID string) (*shipping.ShippingMethod, error) {
	ret := _m.Called(methodID)

	var r0 *shipping.ShippingMethod
	if rf, ok := ret.Get(0).(func(string) *shipping.ShippingMethod); ok {
		r0 = rf(methodID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*shipping.ShippingMethod)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(methodID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ModelFields provides a mock function with given fields:
func (_m *ShippingMethodStore) ModelFields() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// Upsert provides a mock function with given fields: method
func (_m *ShippingMethodStore) Upsert(method *shipping.ShippingMethod) (*shipping.ShippingMethod, error) {
	ret := _m.Called(method)

	var r0 *shipping.ShippingMethod
	if rf, ok := ret.Get(0).(func(*shipping.ShippingMethod) *shipping.ShippingMethod); ok {
		r0 = rf(method)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*shipping.ShippingMethod)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*shipping.ShippingMethod) error); ok {
		r1 = rf(method)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}