// Code generated by mockery v2.23.2. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import mock "github.com/stretchr/testify/mock"

// AppStore is an autogenerated mock type for the AppStore type
type AppStore struct {
	mock.Mock
}

type mockConstructorTestingTNewAppStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewAppStore creates a new instance of AppStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAppStore(t mockConstructorTestingTNewAppStore) *AppStore {
	mock := &AppStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
