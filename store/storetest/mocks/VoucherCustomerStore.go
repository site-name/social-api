// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	product_and_discount "github.com/sitename/sitename/model/product_and_discount"
	mock "github.com/stretchr/testify/mock"
)

// VoucherCustomerStore is an autogenerated mock type for the VoucherCustomerStore type
type VoucherCustomerStore struct {
	mock.Mock
}

// CreateIndexesIfNotExists provides a mock function with given fields:
func (_m *VoucherCustomerStore) CreateIndexesIfNotExists() {
	_m.Called()
}

// DeleteInBulk provides a mock function with given fields: relations
func (_m *VoucherCustomerStore) DeleteInBulk(relations []*product_and_discount.VoucherCustomer) error {
	ret := _m.Called(relations)

	var r0 error
	if rf, ok := ret.Get(0).(func([]*product_and_discount.VoucherCustomer) error); ok {
		r0 = rf(relations)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FilterByEmailAndCustomerEmail provides a mock function with given fields: voucherID, email
func (_m *VoucherCustomerStore) FilterByEmailAndCustomerEmail(voucherID string, email string) ([]*product_and_discount.VoucherCustomer, error) {
	ret := _m.Called(voucherID, email)

	var r0 []*product_and_discount.VoucherCustomer
	if rf, ok := ret.Get(0).(func(string, string) []*product_and_discount.VoucherCustomer); ok {
		r0 = rf(voucherID, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*product_and_discount.VoucherCustomer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(voucherID, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: id
func (_m *VoucherCustomerStore) Get(id string) (*product_and_discount.VoucherCustomer, error) {
	ret := _m.Called(id)

	var r0 *product_and_discount.VoucherCustomer
	if rf, ok := ret.Get(0).(func(string) *product_and_discount.VoucherCustomer); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*product_and_discount.VoucherCustomer)
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

// Save provides a mock function with given fields: voucherCustomer
func (_m *VoucherCustomerStore) Save(voucherCustomer *product_and_discount.VoucherCustomer) (*product_and_discount.VoucherCustomer, error) {
	ret := _m.Called(voucherCustomer)

	var r0 *product_and_discount.VoucherCustomer
	if rf, ok := ret.Get(0).(func(*product_and_discount.VoucherCustomer) *product_and_discount.VoucherCustomer); ok {
		r0 = rf(voucherCustomer)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*product_and_discount.VoucherCustomer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*product_and_discount.VoucherCustomer) error); ok {
		r1 = rf(voucherCustomer)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}