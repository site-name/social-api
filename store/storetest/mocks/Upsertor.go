// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import mock "github.com/stretchr/testify/mock"

// Upsertor is an autogenerated mock type for the Upsertor type
type Upsertor struct {
	mock.Mock
}

// Insert provides a mock function with given fields: list
func (_m *Upsertor) Insert(list ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, list...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(...interface{}) error); ok {
		r0 = rf(list...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: list
func (_m *Upsertor) Update(list ...interface{}) (int64, error) {
	var _ca []interface{}
	_ca = append(_ca, list...)
	ret := _m.Called(_ca...)

	var r0 int64
	if rf, ok := ret.Get(0).(func(...interface{}) int64); ok {
		r0 = rf(list...)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(...interface{}) error); ok {
		r1 = rf(list...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}