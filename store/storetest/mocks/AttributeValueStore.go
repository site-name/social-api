// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	attribute "github.com/sitename/sitename/model/attribute"
	mock "github.com/stretchr/testify/mock"
)

// AttributeValueStore is an autogenerated mock type for the AttributeValueStore type
type AttributeValueStore struct {
	mock.Mock
}

// CreateIndexesIfNotExists provides a mock function with given fields:
func (_m *AttributeValueStore) CreateIndexesIfNotExists() {
	_m.Called()
}

// Get provides a mock function with given fields: attributeID
func (_m *AttributeValueStore) Get(attributeID string) (*attribute.AttributeValue, error) {
	ret := _m.Called(attributeID)

	var r0 *attribute.AttributeValue
	if rf, ok := ret.Get(0).(func(string) *attribute.AttributeValue); ok {
		r0 = rf(attributeID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*attribute.AttributeValue)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(attributeID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllByAttributeID provides a mock function with given fields: attributeID
func (_m *AttributeValueStore) GetAllByAttributeID(attributeID string) ([]*attribute.AttributeValue, error) {
	ret := _m.Called(attributeID)

	var r0 []*attribute.AttributeValue
	if rf, ok := ret.Get(0).(func(string) []*attribute.AttributeValue); ok {
		r0 = rf(attributeID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*attribute.AttributeValue)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(attributeID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ModelFields provides a mock function with given fields:
func (_m *AttributeValueStore) ModelFields() []string {
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

// Save provides a mock function with given fields: _a0
func (_m *AttributeValueStore) Save(_a0 *attribute.AttributeValue) (*attribute.AttributeValue, error) {
	ret := _m.Called(_a0)

	var r0 *attribute.AttributeValue
	if rf, ok := ret.Get(0).(func(*attribute.AttributeValue) *attribute.AttributeValue); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*attribute.AttributeValue)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*attribute.AttributeValue) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}