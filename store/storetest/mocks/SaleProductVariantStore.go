// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	product_and_discount "github.com/sitename/sitename/model/product_and_discount"
	mock "github.com/stretchr/testify/mock"
)

// SaleProductVariantStore is an autogenerated mock type for the SaleProductVariantStore type
type SaleProductVariantStore struct {
	mock.Mock
}

// CreateIndexesIfNotExists provides a mock function with given fields:
func (_m *SaleProductVariantStore) CreateIndexesIfNotExists() {
	_m.Called()
}

// FilterByOption provides a mock function with given fields: options
func (_m *SaleProductVariantStore) FilterByOption(options *product_and_discount.SaleProductVariantFilterOption) ([]*product_and_discount.SaleProductVariant, error) {
	ret := _m.Called(options)

	var r0 []*product_and_discount.SaleProductVariant
	if rf, ok := ret.Get(0).(func(*product_and_discount.SaleProductVariantFilterOption) []*product_and_discount.SaleProductVariant); ok {
		r0 = rf(options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*product_and_discount.SaleProductVariant)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*product_and_discount.SaleProductVariantFilterOption) error); ok {
		r1 = rf(options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Upsert provides a mock function with given fields: relation
func (_m *SaleProductVariantStore) Upsert(relation *product_and_discount.SaleProductVariant) (*product_and_discount.SaleProductVariant, error) {
	ret := _m.Called(relation)

	var r0 *product_and_discount.SaleProductVariant
	if rf, ok := ret.Get(0).(func(*product_and_discount.SaleProductVariant) *product_and_discount.SaleProductVariant); ok {
		r0 = rf(relation)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*product_and_discount.SaleProductVariant)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*product_and_discount.SaleProductVariant) error); ok {
		r1 = rf(relation)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}