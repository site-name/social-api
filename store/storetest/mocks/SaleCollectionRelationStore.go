// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	product_and_discount "github.com/sitename/sitename/model/product_and_discount"
	mock "github.com/stretchr/testify/mock"
)

// SaleCollectionRelationStore is an autogenerated mock type for the SaleCollectionRelationStore type
type SaleCollectionRelationStore struct {
	mock.Mock
}

// CreateIndexesIfNotExists provides a mock function with given fields:
func (_m *SaleCollectionRelationStore) CreateIndexesIfNotExists() {
	_m.Called()
}

// FilterByOption provides a mock function with given fields: option
func (_m *SaleCollectionRelationStore) FilterByOption(option *product_and_discount.SaleCollectionRelationFilterOption) ([]*product_and_discount.SaleCollectionRelation, error) {
	ret := _m.Called(option)

	var r0 []*product_and_discount.SaleCollectionRelation
	if rf, ok := ret.Get(0).(func(*product_and_discount.SaleCollectionRelationFilterOption) []*product_and_discount.SaleCollectionRelation); ok {
		r0 = rf(option)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*product_and_discount.SaleCollectionRelation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*product_and_discount.SaleCollectionRelationFilterOption) error); ok {
		r1 = rf(option)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: relationID
func (_m *SaleCollectionRelationStore) Get(relationID string) (*product_and_discount.SaleCollectionRelation, error) {
	ret := _m.Called(relationID)

	var r0 *product_and_discount.SaleCollectionRelation
	if rf, ok := ret.Get(0).(func(string) *product_and_discount.SaleCollectionRelation); ok {
		r0 = rf(relationID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*product_and_discount.SaleCollectionRelation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(relationID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: relation
func (_m *SaleCollectionRelationStore) Save(relation *product_and_discount.SaleCollectionRelation) (*product_and_discount.SaleCollectionRelation, error) {
	ret := _m.Called(relation)

	var r0 *product_and_discount.SaleCollectionRelation
	if rf, ok := ret.Get(0).(func(*product_and_discount.SaleCollectionRelation) *product_and_discount.SaleCollectionRelation); ok {
		r0 = rf(relation)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*product_and_discount.SaleCollectionRelation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*product_and_discount.SaleCollectionRelation) error); ok {
		r1 = rf(relation)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}