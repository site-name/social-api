// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	account "github.com/sitename/sitename/model/account"
	mock "github.com/stretchr/testify/mock"
)

// UserTermOfServiceStore is an autogenerated mock type for the UserTermOfServiceStore type
type UserTermOfServiceStore struct {
	mock.Mock
}

// CreateIndexesIfNotExists provides a mock function with given fields:
func (_m *UserTermOfServiceStore) CreateIndexesIfNotExists() {
	_m.Called()
}

// Delete provides a mock function with given fields: userID, termsOfServiceId
func (_m *UserTermOfServiceStore) Delete(userID string, termsOfServiceId string) error {
	ret := _m.Called(userID, termsOfServiceId)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(userID, termsOfServiceId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByUser provides a mock function with given fields: userID
func (_m *UserTermOfServiceStore) GetByUser(userID string) (*account.UserTermsOfService, error) {
	ret := _m.Called(userID)

	var r0 *account.UserTermsOfService
	if rf, ok := ret.Get(0).(func(string) *account.UserTermsOfService); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*account.UserTermsOfService)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: userTermsOfService
func (_m *UserTermOfServiceStore) Save(userTermsOfService *account.UserTermsOfService) (*account.UserTermsOfService, error) {
	ret := _m.Called(userTermsOfService)

	var r0 *account.UserTermsOfService
	if rf, ok := ret.Get(0).(func(*account.UserTermsOfService) *account.UserTermsOfService); ok {
		r0 = rf(userTermsOfService)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*account.UserTermsOfService)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*account.UserTermsOfService) error); ok {
		r1 = rf(userTermsOfService)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}