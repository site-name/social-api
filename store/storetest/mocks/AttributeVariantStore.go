// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	attribute "github.com/sitename/sitename/model/attribute"
	mock "github.com/stretchr/testify/mock"
)

// AttributeVariantStore is an autogenerated mock type for the AttributeVariantStore type
type AttributeVariantStore struct {
	mock.Mock
}

// CreateIndexesIfNotExists provides a mock function with given fields:
func (_m *AttributeVariantStore) CreateIndexesIfNotExists() {
	_m.Called()
}

// Get provides a mock function with given fields: attributeVariantID
func (_m *AttributeVariantStore) Get(attributeVariantID string) (*attribute.AttributeVariant, error) {
	ret := _m.Called(attributeVariantID)

	var r0 *attribute.AttributeVariant
	if rf, ok := ret.Get(0).(func(string) *attribute.AttributeVariant); ok {
		r0 = rf(attributeVariantID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*attribute.AttributeVariant)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(attributeVariantID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByOption provides a mock function with given fields: option
func (_m *AttributeVariantStore) GetByOption(option *attribute.AttributeVariantFilterOption) (*attribute.AttributeVariant, error) {
	ret := _m.Called(option)

	var r0 *attribute.AttributeVariant
	if rf, ok := ret.Get(0).(func(*attribute.AttributeVariantFilterOption) *attribute.AttributeVariant); ok {
		r0 = rf(option)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*attribute.AttributeVariant)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*attribute.AttributeVariantFilterOption) error); ok {
		r1 = rf(option)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: attributeVariant
func (_m *AttributeVariantStore) Save(attributeVariant *attribute.AttributeVariant) (*attribute.AttributeVariant, error) {
	ret := _m.Called(attributeVariant)

	var r0 *attribute.AttributeVariant
	if rf, ok := ret.Get(0).(func(*attribute.AttributeVariant) *attribute.AttributeVariant); ok {
		r0 = rf(attributeVariant)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*attribute.AttributeVariant)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*attribute.AttributeVariant) error); ok {
		r1 = rf(attributeVariant)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}