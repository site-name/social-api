// Code generated by mockery v2.23.2. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	model "github.com/sitename/sitename/model"
	mock "github.com/stretchr/testify/mock"
	gorm "gorm.io/gorm"
)

// VatStore is an autogenerated mock type for the VatStore type
type VatStore struct {
	mock.Mock
}

// FilterByOptions provides a mock function with given fields: options
func (_m *VatStore) FilterByOptions(options *model.VatFilterOptions) ([]*model.Vat, error) {
	ret := _m.Called(options)

	var r0 []*model.Vat
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.VatFilterOptions) ([]*model.Vat, error)); ok {
		return rf(options)
	}
	if rf, ok := ret.Get(0).(func(*model.VatFilterOptions) []*model.Vat); ok {
		r0 = rf(options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Vat)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.VatFilterOptions) error); ok {
		r1 = rf(options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Upsert provides a mock function with given fields: transaction, vats
func (_m *VatStore) Upsert(transaction *gorm.DB, vats []*model.Vat) ([]*model.Vat, error) {
	ret := _m.Called(transaction, vats)

	var r0 []*model.Vat
	var r1 error
	if rf, ok := ret.Get(0).(func(*gorm.DB, []*model.Vat) ([]*model.Vat, error)); ok {
		return rf(transaction, vats)
	}
	if rf, ok := ret.Get(0).(func(*gorm.DB, []*model.Vat) []*model.Vat); ok {
		r0 = rf(transaction, vats)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Vat)
		}
	}

	if rf, ok := ret.Get(1).(func(*gorm.DB, []*model.Vat) error); ok {
		r1 = rf(transaction, vats)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewVatStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewVatStore creates a new instance of VatStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewVatStore(t mockConstructorTestingTNewVatStore) *VatStore {
	mock := &VatStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}