// Code generated by mockery v2.23.2. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	boil "github.com/volatiletech/sqlboiler/v4/boil"

	model "github.com/sitename/sitename/model"

	model_helper "github.com/sitename/sitename/model_helper"
)

// PaymentTransactionStore is an autogenerated mock type for the PaymentTransactionStore type
type PaymentTransactionStore struct {
	mock.Mock
}

// FilterByOption provides a mock function with given fields: option
func (_m *PaymentTransactionStore) FilterByOption(option model_helper.PaymentTransactionFilterOpts) ([]*model.PaymentTransaction, error) {
	ret := _m.Called(option)

	var r0 []*model.PaymentTransaction
	var r1 error
	if rf, ok := ret.Get(0).(func(model_helper.PaymentTransactionFilterOpts) ([]*model.PaymentTransaction, error)); ok {
		return rf(option)
	}
	if rf, ok := ret.Get(0).(func(model_helper.PaymentTransactionFilterOpts) []*model.PaymentTransaction); ok {
		r0 = rf(option)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.PaymentTransaction)
		}
	}

	if rf, ok := ret.Get(1).(func(model_helper.PaymentTransactionFilterOpts) error); ok {
		r1 = rf(option)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: id
func (_m *PaymentTransactionStore) Get(id string) (*model.PaymentTransaction, error) {
	ret := _m.Called(id)

	var r0 *model.PaymentTransaction
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*model.PaymentTransaction, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *model.PaymentTransaction); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.PaymentTransaction)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Upsert provides a mock function with given fields: tx, paymentTransaction
func (_m *PaymentTransactionStore) Upsert(tx boil.ContextTransactor, paymentTransaction model.PaymentTransaction) (*model.PaymentTransaction, error) {
	ret := _m.Called(tx, paymentTransaction)

	var r0 *model.PaymentTransaction
	var r1 error
	if rf, ok := ret.Get(0).(func(boil.ContextTransactor, model.PaymentTransaction) (*model.PaymentTransaction, error)); ok {
		return rf(tx, paymentTransaction)
	}
	if rf, ok := ret.Get(0).(func(boil.ContextTransactor, model.PaymentTransaction) *model.PaymentTransaction); ok {
		r0 = rf(tx, paymentTransaction)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.PaymentTransaction)
		}
	}

	if rf, ok := ret.Get(1).(func(boil.ContextTransactor, model.PaymentTransaction) error); ok {
		r1 = rf(tx, paymentTransaction)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewPaymentTransactionStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewPaymentTransactionStore creates a new instance of PaymentTransactionStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPaymentTransactionStore(t mockConstructorTestingTNewPaymentTransactionStore) *PaymentTransactionStore {
	mock := &PaymentTransactionStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
