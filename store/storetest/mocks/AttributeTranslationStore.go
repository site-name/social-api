// Code generated by mockery v2.23.2. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import mock "github.com/stretchr/testify/mock"

// AttributeTranslationStore is an autogenerated mock type for the AttributeTranslationStore type
type AttributeTranslationStore struct {
	mock.Mock
}

type mockConstructorTestingTNewAttributeTranslationStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewAttributeTranslationStore creates a new instance of AttributeTranslationStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAttributeTranslationStore(t mockConstructorTestingTNewAttributeTranslationStore) *AttributeTranslationStore {
	mock := &AttributeTranslationStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
