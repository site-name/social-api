// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	product_and_discount "github.com/sitename/sitename/model/product_and_discount"
	mock "github.com/stretchr/testify/mock"
)

// VoucherCategoryStore is an autogenerated mock type for the VoucherCategoryStore type
type VoucherCategoryStore struct {
	mock.Mock
}

// CreateIndexesIfNotExists provides a mock function with given fields:
func (_m *VoucherCategoryStore) CreateIndexesIfNotExists() {
	_m.Called()
}

// Get provides a mock function with given fields: voucherCategoryID
func (_m *VoucherCategoryStore) Get(voucherCategoryID string) (*product_and_discount.VoucherCategory, error) {
	ret := _m.Called(voucherCategoryID)

	var r0 *product_and_discount.VoucherCategory
	if rf, ok := ret.Get(0).(func(string) *product_and_discount.VoucherCategory); ok {
		r0 = rf(voucherCategoryID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*product_and_discount.VoucherCategory)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(voucherCategoryID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProductCategoriesByVoucherID provides a mock function with given fields: voucherID
func (_m *VoucherCategoryStore) ProductCategoriesByVoucherID(voucherID string) ([]*product_and_discount.Category, error) {
	ret := _m.Called(voucherID)

	var r0 []*product_and_discount.Category
	if rf, ok := ret.Get(0).(func(string) []*product_and_discount.Category); ok {
		r0 = rf(voucherID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*product_and_discount.Category)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(voucherID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Upsert provides a mock function with given fields: voucherCategory
func (_m *VoucherCategoryStore) Upsert(voucherCategory *product_and_discount.VoucherCategory) (*product_and_discount.VoucherCategory, error) {
	ret := _m.Called(voucherCategory)

	var r0 *product_and_discount.VoucherCategory
	if rf, ok := ret.Get(0).(func(*product_and_discount.VoucherCategory) *product_and_discount.VoucherCategory); ok {
		r0 = rf(voucherCategory)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*product_and_discount.VoucherCategory)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*product_and_discount.VoucherCategory) error); ok {
		r1 = rf(voucherCategory)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}