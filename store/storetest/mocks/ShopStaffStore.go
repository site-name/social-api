// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	shop "github.com/sitename/sitename/model/shop"
	mock "github.com/stretchr/testify/mock"
)

// ShopStaffStore is an autogenerated mock type for the ShopStaffStore type
type ShopStaffStore struct {
	mock.Mock
}

// CreateIndexesIfNotExists provides a mock function with given fields:
func (_m *ShopStaffStore) CreateIndexesIfNotExists() {
	_m.Called()
}

// FilterByShopAndStaff provides a mock function with given fields: shopID, staffID
func (_m *ShopStaffStore) FilterByShopAndStaff(shopID string, staffID string) (*shop.ShopStaffRelation, error) {
	ret := _m.Called(shopID, staffID)

	var r0 *shop.ShopStaffRelation
	if rf, ok := ret.Get(0).(func(string, string) *shop.ShopStaffRelation); ok {
		r0 = rf(shopID, staffID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*shop.ShopStaffRelation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(shopID, staffID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: shopStaffID
func (_m *ShopStaffStore) Get(shopStaffID string) (*shop.ShopStaffRelation, error) {
	ret := _m.Called(shopStaffID)

	var r0 *shop.ShopStaffRelation
	if rf, ok := ret.Get(0).(func(string) *shop.ShopStaffRelation); ok {
		r0 = rf(shopStaffID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*shop.ShopStaffRelation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(shopStaffID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: shopStaff
func (_m *ShopStaffStore) Save(shopStaff *shop.ShopStaffRelation) (*shop.ShopStaffRelation, error) {
	ret := _m.Called(shopStaff)

	var r0 *shop.ShopStaffRelation
	if rf, ok := ret.Get(0).(func(*shop.ShopStaffRelation) *shop.ShopStaffRelation); ok {
		r0 = rf(shopStaff)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*shop.ShopStaffRelation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*shop.ShopStaffRelation) error); ok {
		r1 = rf(shopStaff)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}