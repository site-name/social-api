// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	attribute "github.com/sitename/sitename/model/attribute"
	mock "github.com/stretchr/testify/mock"
)

// AssignedVariantAttributeStore is an autogenerated mock type for the AssignedVariantAttributeStore type
type AssignedVariantAttributeStore struct {
	mock.Mock
}

// CreateIndexesIfNotExists provides a mock function with given fields:
func (_m *AssignedVariantAttributeStore) CreateIndexesIfNotExists() {
	_m.Called()
}

// Get provides a mock function with given fields: id
func (_m *AssignedVariantAttributeStore) Get(id string) (*attribute.AssignedVariantAttribute, error) {
	ret := _m.Called(id)

	var r0 *attribute.AssignedVariantAttribute
	if rf, ok := ret.Get(0).(func(string) *attribute.AssignedVariantAttribute); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*attribute.AssignedVariantAttribute)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWithOption provides a mock function with given fields: option
func (_m *AssignedVariantAttributeStore) GetWithOption(option *attribute.AssignedVariantAttributeFilterOption) (*attribute.AssignedVariantAttribute, error) {
	ret := _m.Called(option)

	var r0 *attribute.AssignedVariantAttribute
	if rf, ok := ret.Get(0).(func(*attribute.AssignedVariantAttributeFilterOption) *attribute.AssignedVariantAttribute); ok {
		r0 = rf(option)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*attribute.AssignedVariantAttribute)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*attribute.AssignedVariantAttributeFilterOption) error); ok {
		r1 = rf(option)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: assignedVariantAttribute
func (_m *AssignedVariantAttributeStore) Save(assignedVariantAttribute *attribute.AssignedVariantAttribute) (*attribute.AssignedVariantAttribute, error) {
	ret := _m.Called(assignedVariantAttribute)

	var r0 *attribute.AssignedVariantAttribute
	if rf, ok := ret.Get(0).(func(*attribute.AssignedVariantAttribute) *attribute.AssignedVariantAttribute); ok {
		r0 = rf(assignedVariantAttribute)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*attribute.AssignedVariantAttribute)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*attribute.AssignedVariantAttribute) error); ok {
		r1 = rf(assignedVariantAttribute)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}