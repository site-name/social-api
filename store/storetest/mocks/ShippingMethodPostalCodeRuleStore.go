// Code generated by mockery v2.23.2. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	model "github.com/sitename/sitename/model"
	mock "github.com/stretchr/testify/mock"
	gorm "gorm.io/gorm"
)

// ShippingMethodPostalCodeRuleStore is an autogenerated mock type for the ShippingMethodPostalCodeRuleStore type
type ShippingMethodPostalCodeRuleStore struct {
	mock.Mock
}

// Delete provides a mock function with given fields: transaction, ids
func (_m *ShippingMethodPostalCodeRuleStore) Delete(transaction *gorm.DB, ids ...string) error {
	_va := make([]interface{}, len(ids))
	for _i := range ids {
		_va[_i] = ids[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, transaction)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(*gorm.DB, ...string) error); ok {
		r0 = rf(transaction, ids...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FilterByOptions provides a mock function with given fields: options
func (_m *ShippingMethodPostalCodeRuleStore) FilterByOptions(options *model.ShippingMethodPostalCodeRuleFilterOptions) ([]*model.ShippingMethodPostalCodeRule, error) {
	ret := _m.Called(options)

	var r0 []*model.ShippingMethodPostalCodeRule
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.ShippingMethodPostalCodeRuleFilterOptions) ([]*model.ShippingMethodPostalCodeRule, error)); ok {
		return rf(options)
	}
	if rf, ok := ret.Get(0).(func(*model.ShippingMethodPostalCodeRuleFilterOptions) []*model.ShippingMethodPostalCodeRule); ok {
		r0 = rf(options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.ShippingMethodPostalCodeRule)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.ShippingMethodPostalCodeRuleFilterOptions) error); ok {
		r1 = rf(options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: transaction, rules
func (_m *ShippingMethodPostalCodeRuleStore) Save(transaction *gorm.DB, rules model.ShippingMethodPostalCodeRules) (model.ShippingMethodPostalCodeRules, error) {
	ret := _m.Called(transaction, rules)

	var r0 model.ShippingMethodPostalCodeRules
	var r1 error
	if rf, ok := ret.Get(0).(func(*gorm.DB, model.ShippingMethodPostalCodeRules) (model.ShippingMethodPostalCodeRules, error)); ok {
		return rf(transaction, rules)
	}
	if rf, ok := ret.Get(0).(func(*gorm.DB, model.ShippingMethodPostalCodeRules) model.ShippingMethodPostalCodeRules); ok {
		r0 = rf(transaction, rules)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(model.ShippingMethodPostalCodeRules)
		}
	}

	if rf, ok := ret.Get(1).(func(*gorm.DB, model.ShippingMethodPostalCodeRules) error); ok {
		r1 = rf(transaction, rules)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ScanFields provides a mock function with given fields: rule
func (_m *ShippingMethodPostalCodeRuleStore) ScanFields(rule *model.ShippingMethodPostalCodeRule) []interface{} {
	ret := _m.Called(rule)

	var r0 []interface{}
	if rf, ok := ret.Get(0).(func(*model.ShippingMethodPostalCodeRule) []interface{}); ok {
		r0 = rf(rule)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]interface{})
		}
	}

	return r0
}

type mockConstructorTestingTNewShippingMethodPostalCodeRuleStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewShippingMethodPostalCodeRuleStore creates a new instance of ShippingMethodPostalCodeRuleStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewShippingMethodPostalCodeRuleStore(t mockConstructorTestingTNewShippingMethodPostalCodeRuleStore) *ShippingMethodPostalCodeRuleStore {
	mock := &ShippingMethodPostalCodeRuleStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}