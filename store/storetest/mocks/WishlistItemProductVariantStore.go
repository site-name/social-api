// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	gorp "github.com/mattermost/gorp"
	mock "github.com/stretchr/testify/mock"

	wishlist "github.com/sitename/sitename/model/wishlist"
)

// WishlistItemProductVariantStore is an autogenerated mock type for the WishlistItemProductVariantStore type
type WishlistItemProductVariantStore struct {
	mock.Mock
}

// BulkUpsert provides a mock function with given fields: transaction, relations
func (_m *WishlistItemProductVariantStore) BulkUpsert(transaction *gorp.Transaction, relations []*wishlist.WishlistItemProductVariant) ([]*wishlist.WishlistItemProductVariant, error) {
	ret := _m.Called(transaction, relations)

	var r0 []*wishlist.WishlistItemProductVariant
	if rf, ok := ret.Get(0).(func(*gorp.Transaction, []*wishlist.WishlistItemProductVariant) []*wishlist.WishlistItemProductVariant); ok {
		r0 = rf(transaction, relations)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*wishlist.WishlistItemProductVariant)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gorp.Transaction, []*wishlist.WishlistItemProductVariant) error); ok {
		r1 = rf(transaction, relations)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateIndexesIfNotExists provides a mock function with given fields:
func (_m *WishlistItemProductVariantStore) CreateIndexesIfNotExists() {
	_m.Called()
}

// DeleteRelation provides a mock function with given fields: relation
func (_m *WishlistItemProductVariantStore) DeleteRelation(relation *wishlist.WishlistItemProductVariant) (int64, error) {
	ret := _m.Called(relation)

	var r0 int64
	if rf, ok := ret.Get(0).(func(*wishlist.WishlistItemProductVariant) int64); ok {
		r0 = rf(relation)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*wishlist.WishlistItemProductVariant) error); ok {
		r1 = rf(relation)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetById provides a mock function with given fields: selector, id
func (_m *WishlistItemProductVariantStore) GetById(selector *gorp.Transaction, id string) (*wishlist.WishlistItemProductVariant, error) {
	ret := _m.Called(selector, id)

	var r0 *wishlist.WishlistItemProductVariant
	if rf, ok := ret.Get(0).(func(*gorp.Transaction, string) *wishlist.WishlistItemProductVariant); ok {
		r0 = rf(selector, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*wishlist.WishlistItemProductVariant)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gorp.Transaction, string) error); ok {
		r1 = rf(selector, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: wishlistVariant
func (_m *WishlistItemProductVariantStore) Save(wishlistVariant *wishlist.WishlistItemProductVariant) (*wishlist.WishlistItemProductVariant, error) {
	ret := _m.Called(wishlistVariant)

	var r0 *wishlist.WishlistItemProductVariant
	if rf, ok := ret.Get(0).(func(*wishlist.WishlistItemProductVariant) *wishlist.WishlistItemProductVariant); ok {
		r0 = rf(wishlistVariant)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*wishlist.WishlistItemProductVariant)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*wishlist.WishlistItemProductVariant) error); ok {
		r1 = rf(wishlistVariant)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}