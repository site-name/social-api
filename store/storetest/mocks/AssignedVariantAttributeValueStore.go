// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	attribute "github.com/sitename/sitename/model/attribute"
	mock "github.com/stretchr/testify/mock"
)

// AssignedVariantAttributeValueStore is an autogenerated mock type for the AssignedVariantAttributeValueStore type
type AssignedVariantAttributeValueStore struct {
	mock.Mock
}

// CreateIndexesIfNotExists provides a mock function with given fields:
func (_m *AssignedVariantAttributeValueStore) CreateIndexesIfNotExists() {
	_m.Called()
}

// Get provides a mock function with given fields: assignedVariantAttrValueID
func (_m *AssignedVariantAttributeValueStore) Get(assignedVariantAttrValueID string) (*attribute.AssignedVariantAttributeValue, error) {
	ret := _m.Called(assignedVariantAttrValueID)

	var r0 *attribute.AssignedVariantAttributeValue
	if rf, ok := ret.Get(0).(func(string) *attribute.AssignedVariantAttributeValue); ok {
		r0 = rf(assignedVariantAttrValueID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*attribute.AssignedVariantAttributeValue)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(assignedVariantAttrValueID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: assignedVariantAttrValue
func (_m *AssignedVariantAttributeValueStore) Save(assignedVariantAttrValue *attribute.AssignedVariantAttributeValue) (*attribute.AssignedVariantAttributeValue, error) {
	ret := _m.Called(assignedVariantAttrValue)

	var r0 *attribute.AssignedVariantAttributeValue
	if rf, ok := ret.Get(0).(func(*attribute.AssignedVariantAttributeValue) *attribute.AssignedVariantAttributeValue); ok {
		r0 = rf(assignedVariantAttrValue)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*attribute.AssignedVariantAttributeValue)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*attribute.AssignedVariantAttributeValue) error); ok {
		r1 = rf(assignedVariantAttrValue)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveInBulk provides a mock function with given fields: assignmentID, attributeValueIDs
func (_m *AssignedVariantAttributeValueStore) SaveInBulk(assignmentID string, attributeValueIDs []string) ([]*attribute.AssignedVariantAttributeValue, error) {
	ret := _m.Called(assignmentID, attributeValueIDs)

	var r0 []*attribute.AssignedVariantAttributeValue
	if rf, ok := ret.Get(0).(func(string, []string) []*attribute.AssignedVariantAttributeValue); ok {
		r0 = rf(assignmentID, attributeValueIDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*attribute.AssignedVariantAttributeValue)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, []string) error); ok {
		r1 = rf(assignmentID, attributeValueIDs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SelectForSort provides a mock function with given fields: assignmentID
func (_m *AssignedVariantAttributeValueStore) SelectForSort(assignmentID string) ([]*attribute.AssignedVariantAttributeValue, []*attribute.AttributeValue, error) {
	ret := _m.Called(assignmentID)

	var r0 []*attribute.AssignedVariantAttributeValue
	if rf, ok := ret.Get(0).(func(string) []*attribute.AssignedVariantAttributeValue); ok {
		r0 = rf(assignmentID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*attribute.AssignedVariantAttributeValue)
		}
	}

	var r1 []*attribute.AttributeValue
	if rf, ok := ret.Get(1).(func(string) []*attribute.AttributeValue); ok {
		r1 = rf(assignmentID)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]*attribute.AttributeValue)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string) error); ok {
		r2 = rf(assignmentID)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// UpdateInBulk provides a mock function with given fields: attributeValues
func (_m *AssignedVariantAttributeValueStore) UpdateInBulk(attributeValues []*attribute.AssignedVariantAttributeValue) error {
	ret := _m.Called(attributeValues)

	var r0 error
	if rf, ok := ret.Get(0).(func([]*attribute.AssignedVariantAttributeValue) error); ok {
		r0 = rf(attributeValues)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}