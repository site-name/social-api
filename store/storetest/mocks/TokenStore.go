// Code generated by mockery v2.23.2. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	model "github.com/sitename/sitename/model"
	mock "github.com/stretchr/testify/mock"
)

// TokenStore is an autogenerated mock type for the TokenStore type
type TokenStore struct {
	mock.Mock
}

// Cleanup provides a mock function with given fields:
func (_m *TokenStore) Cleanup() {
	_m.Called()
}

// Delete provides a mock function with given fields: token
func (_m *TokenStore) Delete(token string) error {
	ret := _m.Called(token)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllTokensByType provides a mock function with given fields: tokenType
func (_m *TokenStore) GetAllTokensByType(tokenType model.TokenType) ([]*model.Token, error) {
	ret := _m.Called(tokenType)

	var r0 []*model.Token
	var r1 error
	if rf, ok := ret.Get(0).(func(model.TokenType) ([]*model.Token, error)); ok {
		return rf(tokenType)
	}
	if rf, ok := ret.Get(0).(func(model.TokenType) []*model.Token); ok {
		r0 = rf(tokenType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Token)
		}
	}

	if rf, ok := ret.Get(1).(func(model.TokenType) error); ok {
		r1 = rf(tokenType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByToken provides a mock function with given fields: token
func (_m *TokenStore) GetByToken(token string) (*model.Token, error) {
	ret := _m.Called(token)

	var r0 *model.Token
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*model.Token, error)); ok {
		return rf(token)
	}
	if rf, ok := ret.Get(0).(func(string) *model.Token); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Token)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveAllTokensByType provides a mock function with given fields: tokenType
func (_m *TokenStore) RemoveAllTokensByType(tokenType string) error {
	ret := _m.Called(tokenType)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(tokenType)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Save provides a mock function with given fields: recovery
func (_m *TokenStore) Save(recovery *model.Token) error {
	ret := _m.Called(recovery)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Token) error); ok {
		r0 = rf(recovery)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewTokenStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewTokenStore creates a new instance of TokenStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTokenStore(t mockConstructorTestingTNewTokenStore) *TokenStore {
	mock := &TokenStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}