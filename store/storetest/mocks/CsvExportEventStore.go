// Code generated by mockery v2.23.2. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	model "github.com/sitename/sitename/model"
	mock "github.com/stretchr/testify/mock"

	model_helper "github.com/sitename/sitename/model_helper"
)

// CsvExportEventStore is an autogenerated mock type for the CsvExportEventStore type
type CsvExportEventStore struct {
	mock.Mock
}

// FilterByOption provides a mock function with given fields: options
func (_m *CsvExportEventStore) FilterByOption(options model_helper.ExportEventFilterOption) ([]*model.ExportEvent, error) {
	ret := _m.Called(options)

	var r0 []*model.ExportEvent
	var r1 error
	if rf, ok := ret.Get(0).(func(model_helper.ExportEventFilterOption) ([]*model.ExportEvent, error)); ok {
		return rf(options)
	}
	if rf, ok := ret.Get(0).(func(model_helper.ExportEventFilterOption) []*model.ExportEvent); ok {
		r0 = rf(options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.ExportEvent)
		}
	}

	if rf, ok := ret.Get(1).(func(model_helper.ExportEventFilterOption) error); ok {
		r1 = rf(options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: event
func (_m *CsvExportEventStore) Save(event model.ExportEvent) (*model.ExportEvent, error) {
	ret := _m.Called(event)

	var r0 *model.ExportEvent
	var r1 error
	if rf, ok := ret.Get(0).(func(model.ExportEvent) (*model.ExportEvent, error)); ok {
		return rf(event)
	}
	if rf, ok := ret.Get(0).(func(model.ExportEvent) *model.ExportEvent); ok {
		r0 = rf(event)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ExportEvent)
		}
	}

	if rf, ok := ret.Get(1).(func(model.ExportEvent) error); ok {
		r1 = rf(event)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewCsvExportEventStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewCsvExportEventStore creates a new instance of CsvExportEventStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCsvExportEventStore(t mockConstructorTestingTNewCsvExportEventStore) *CsvExportEventStore {
	mock := &CsvExportEventStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
