// Code generated by mockery v2.23.2. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	context "context"

	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"

	model "github.com/sitename/sitename/model"

	store "github.com/sitename/sitename/store"
)

// UserStore is an autogenerated mock type for the UserStore type
type UserStore struct {
	mock.Mock
}

// AddRelations provides a mock function with given fields: transaction, userID, relations, customerNoteOnUser
func (_m *UserStore) AddRelations(transaction *gorm.DB, userID string, relations interface{}, customerNoteOnUser bool) *model.AppError {
	ret := _m.Called(transaction, userID, relations, customerNoteOnUser)

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func(*gorm.DB, string, interface{}, bool) *model.AppError); ok {
		r0 = rf(transaction, userID, relations, customerNoteOnUser)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// AnalyticsActiveCount provides a mock function with given fields: time, options
func (_m *UserStore) AnalyticsActiveCount(time int64, options model.UserCountOptions) (int64, error) {
	ret := _m.Called(time, options)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(int64, model.UserCountOptions) (int64, error)); ok {
		return rf(time, options)
	}
	if rf, ok := ret.Get(0).(func(int64, model.UserCountOptions) int64); ok {
		r0 = rf(time, options)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(int64, model.UserCountOptions) error); ok {
		r1 = rf(time, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AnalyticsActiveCountForPeriod provides a mock function with given fields: startTime, endTime, options
func (_m *UserStore) AnalyticsActiveCountForPeriod(startTime int64, endTime int64, options model.UserCountOptions) (int64, error) {
	ret := _m.Called(startTime, endTime, options)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(int64, int64, model.UserCountOptions) (int64, error)); ok {
		return rf(startTime, endTime, options)
	}
	if rf, ok := ret.Get(0).(func(int64, int64, model.UserCountOptions) int64); ok {
		r0 = rf(startTime, endTime, options)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(int64, int64, model.UserCountOptions) error); ok {
		r1 = rf(startTime, endTime, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AnalyticsGetExternalUsers provides a mock function with given fields: hostDomain
func (_m *UserStore) AnalyticsGetExternalUsers(hostDomain string) (bool, error) {
	ret := _m.Called(hostDomain)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (bool, error)); ok {
		return rf(hostDomain)
	}
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(hostDomain)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(hostDomain)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AnalyticsGetGuestCount provides a mock function with given fields:
func (_m *UserStore) AnalyticsGetGuestCount() (int64, error) {
	ret := _m.Called()

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func() (int64, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AnalyticsGetInactiveUsersCount provides a mock function with given fields:
func (_m *UserStore) AnalyticsGetInactiveUsersCount() (int64, error) {
	ret := _m.Called()

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func() (int64, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AnalyticsGetSystemAdminCount provides a mock function with given fields:
func (_m *UserStore) AnalyticsGetSystemAdminCount() (int64, error) {
	ret := _m.Called()

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func() (int64, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ClearAllCustomRoleAssignments provides a mock function with given fields:
func (_m *UserStore) ClearAllCustomRoleAssignments() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ClearCaches provides a mock function with given fields:
func (_m *UserStore) ClearCaches() {
	_m.Called()
}

// Count provides a mock function with given fields: options
func (_m *UserStore) Count(options model.UserCountOptions) (int64, error) {
	ret := _m.Called(options)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(model.UserCountOptions) (int64, error)); ok {
		return rf(options)
	}
	if rf, ok := ret.Get(0).(func(model.UserCountOptions) int64); ok {
		r0 = rf(options)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(model.UserCountOptions) error); ok {
		r1 = rf(options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FilterByOptions provides a mock function with given fields: ctx, options
func (_m *UserStore) FilterByOptions(ctx context.Context, options *model.UserFilterOptions) ([]*model.User, error) {
	ret := _m.Called(ctx, options)

	var r0 []*model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.UserFilterOptions) ([]*model.User, error)); ok {
		return rf(ctx, options)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.UserFilterOptions) []*model.User); ok {
		r0 = rf(ctx, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.UserFilterOptions) error); ok {
		r1 = rf(ctx, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllProfiles provides a mock function with given fields: options
func (_m *UserStore) GetAllProfiles(options *model.UserGetOptions) ([]*model.User, error) {
	ret := _m.Called(options)

	var r0 []*model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.UserGetOptions) ([]*model.User, error)); ok {
		return rf(options)
	}
	if rf, ok := ret.Get(0).(func(*model.UserGetOptions) []*model.User); ok {
		r0 = rf(options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.UserGetOptions) error); ok {
		r1 = rf(options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByOptions provides a mock function with given fields: ctx, options
func (_m *UserStore) GetByOptions(ctx context.Context, options *model.UserFilterOptions) (*model.User, error) {
	ret := _m.Called(ctx, options)

	var r0 *model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.UserFilterOptions) (*model.User, error)); ok {
		return rf(ctx, options)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.UserFilterOptions) *model.User); ok {
		r0 = rf(ctx, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.UserFilterOptions) error); ok {
		r1 = rf(ctx, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEtagForAllProfiles provides a mock function with given fields:
func (_m *UserStore) GetEtagForAllProfiles() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetEtagForProfiles provides a mock function with given fields: teamID
func (_m *UserStore) GetEtagForProfiles(teamID string) string {
	ret := _m.Called(teamID)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(teamID)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetForLogin provides a mock function with given fields: loginID, allowSignInWithUsername, allowSignInWithEmail
func (_m *UserStore) GetForLogin(loginID string, allowSignInWithUsername bool, allowSignInWithEmail bool) (*model.User, error) {
	ret := _m.Called(loginID, allowSignInWithUsername, allowSignInWithEmail)

	var r0 *model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(string, bool, bool) (*model.User, error)); ok {
		return rf(loginID, allowSignInWithUsername, allowSignInWithEmail)
	}
	if rf, ok := ret.Get(0).(func(string, bool, bool) *model.User); ok {
		r0 = rf(loginID, allowSignInWithUsername, allowSignInWithEmail)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(string, bool, bool) error); ok {
		r1 = rf(loginID, allowSignInWithUsername, allowSignInWithEmail)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetKnownUsers provides a mock function with given fields: userID
func (_m *UserStore) GetKnownUsers(userID string) ([]string, error) {
	ret := _m.Called(userID)

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]string, error)); ok {
		return rf(userID)
	}
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProfileByIds provides a mock function with given fields: ctx, userIds, options, allowFromCache
func (_m *UserStore) GetProfileByIds(ctx context.Context, userIds []string, options *store.UserGetByIdsOpts, allowFromCache bool) ([]*model.User, error) {
	ret := _m.Called(ctx, userIds, options, allowFromCache)

	var r0 []*model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []string, *store.UserGetByIdsOpts, bool) ([]*model.User, error)); ok {
		return rf(ctx, userIds, options, allowFromCache)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []string, *store.UserGetByIdsOpts, bool) []*model.User); ok {
		r0 = rf(ctx, userIds, options, allowFromCache)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []string, *store.UserGetByIdsOpts, bool) error); ok {
		r1 = rf(ctx, userIds, options, allowFromCache)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSystemAdminProfiles provides a mock function with given fields:
func (_m *UserStore) GetSystemAdminProfiles() (map[string]*model.User, error) {
	ret := _m.Called()

	var r0 map[string]*model.User
	var r1 error
	if rf, ok := ret.Get(0).(func() (map[string]*model.User, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() map[string]*model.User); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUnreadCount provides a mock function with given fields: userID
func (_m *UserStore) GetUnreadCount(userID string) (int64, error) {
	ret := _m.Called(userID)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (int64, error)); ok {
		return rf(userID)
	}
	if rf, ok := ret.Get(0).(func(string) int64); ok {
		r0 = rf(userID)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUsersBatchForIndexing provides a mock function with given fields: startTime, endTime, limit
func (_m *UserStore) GetUsersBatchForIndexing(startTime int64, endTime int64, limit int) ([]*model.UserForIndexing, error) {
	ret := _m.Called(startTime, endTime, limit)

	var r0 []*model.UserForIndexing
	var r1 error
	if rf, ok := ret.Get(0).(func(int64, int64, int) ([]*model.UserForIndexing, error)); ok {
		return rf(startTime, endTime, limit)
	}
	if rf, ok := ret.Get(0).(func(int64, int64, int) []*model.UserForIndexing); ok {
		r0 = rf(startTime, endTime, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.UserForIndexing)
		}
	}

	if rf, ok := ret.Get(1).(func(int64, int64, int) error); ok {
		r1 = rf(startTime, endTime, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InferSystemInstallDate provides a mock function with given fields:
func (_m *UserStore) InferSystemInstallDate() (int64, error) {
	ret := _m.Called()

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func() (int64, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InvalidateProfileCacheForUser provides a mock function with given fields: userID
func (_m *UserStore) InvalidateProfileCacheForUser(userID string) {
	_m.Called(userID)
}

// PermanentDelete provides a mock function with given fields: userID
func (_m *UserStore) PermanentDelete(userID string) error {
	ret := _m.Called(userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveRelations provides a mock function with given fields: transaction, userID, relations, customerNoteOnUser
func (_m *UserStore) RemoveRelations(transaction *gorm.DB, userID string, relations interface{}, customerNoteOnUser bool) *model.AppError {
	ret := _m.Called(transaction, userID, relations, customerNoteOnUser)

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func(*gorm.DB, string, interface{}, bool) *model.AppError); ok {
		r0 = rf(transaction, userID, relations, customerNoteOnUser)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// ResetAuthDataToEmailForUsers provides a mock function with given fields: service, userIDs, includeDeleted, dryRun
func (_m *UserStore) ResetAuthDataToEmailForUsers(service string, userIDs []string, includeDeleted bool, dryRun bool) (int, error) {
	ret := _m.Called(service, userIDs, includeDeleted, dryRun)

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(string, []string, bool, bool) (int, error)); ok {
		return rf(service, userIDs, includeDeleted, dryRun)
	}
	if rf, ok := ret.Get(0).(func(string, []string, bool, bool) int); ok {
		r0 = rf(service, userIDs, includeDeleted, dryRun)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(string, []string, bool, bool) error); ok {
		r1 = rf(service, userIDs, includeDeleted, dryRun)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ResetLastPictureUpdate provides a mock function with given fields: userID
func (_m *UserStore) ResetLastPictureUpdate(userID string) error {
	ret := _m.Called(userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Save provides a mock function with given fields: user
func (_m *UserStore) Save(user *model.User) (*model.User, error) {
	ret := _m.Called(user)

	var r0 *model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.User) (*model.User, error)); ok {
		return rf(user)
	}
	if rf, ok := ret.Get(0).(func(*model.User) *model.User); ok {
		r0 = rf(user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.User) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ScanFields provides a mock function with given fields: user
func (_m *UserStore) ScanFields(user *model.User) []interface{} {
	ret := _m.Called(user)

	var r0 []interface{}
	if rf, ok := ret.Get(0).(func(*model.User) []interface{}); ok {
		r0 = rf(user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]interface{})
		}
	}

	return r0
}

// Search provides a mock function with given fields: term, options
func (_m *UserStore) Search(term string, options *model.UserSearchOptions) ([]*model.User, error) {
	ret := _m.Called(term, options)

	var r0 []*model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(string, *model.UserSearchOptions) ([]*model.User, error)); ok {
		return rf(term, options)
	}
	if rf, ok := ret.Get(0).(func(string, *model.UserSearchOptions) []*model.User); ok {
		r0 = rf(term, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(string, *model.UserSearchOptions) error); ok {
		r1 = rf(term, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: user, allowRoleUpdate
func (_m *UserStore) Update(user *model.User, allowRoleUpdate bool) (*model.UserUpdate, error) {
	ret := _m.Called(user, allowRoleUpdate)

	var r0 *model.UserUpdate
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.User, bool) (*model.UserUpdate, error)); ok {
		return rf(user, allowRoleUpdate)
	}
	if rf, ok := ret.Get(0).(func(*model.User, bool) *model.UserUpdate); ok {
		r0 = rf(user, allowRoleUpdate)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.UserUpdate)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.User, bool) error); ok {
		r1 = rf(user, allowRoleUpdate)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateAuthData provides a mock function with given fields: userID, service, authData, email, resetMfa
func (_m *UserStore) UpdateAuthData(userID string, service string, authData *string, email string, resetMfa bool) (string, error) {
	ret := _m.Called(userID, service, authData, email, resetMfa)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string, *string, string, bool) (string, error)); ok {
		return rf(userID, service, authData, email, resetMfa)
	}
	if rf, ok := ret.Get(0).(func(string, string, *string, string, bool) string); ok {
		r0 = rf(userID, service, authData, email, resetMfa)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string, string, *string, string, bool) error); ok {
		r1 = rf(userID, service, authData, email, resetMfa)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateFailedPasswordAttempts provides a mock function with given fields: userID, attempts
func (_m *UserStore) UpdateFailedPasswordAttempts(userID string, attempts int) error {
	ret := _m.Called(userID, attempts)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, int) error); ok {
		r0 = rf(userID, attempts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateLastPictureUpdate provides a mock function with given fields: userID
func (_m *UserStore) UpdateLastPictureUpdate(userID string) error {
	ret := _m.Called(userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateMfaActive provides a mock function with given fields: userID, active
func (_m *UserStore) UpdateMfaActive(userID string, active bool) error {
	ret := _m.Called(userID, active)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, bool) error); ok {
		r0 = rf(userID, active)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateMfaSecret provides a mock function with given fields: userID, secret
func (_m *UserStore) UpdateMfaSecret(userID string, secret string) error {
	ret := _m.Called(userID, secret)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(userID, secret)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdatePassword provides a mock function with given fields: userID, newPassword
func (_m *UserStore) UpdatePassword(userID string, newPassword string) error {
	ret := _m.Called(userID, newPassword)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(userID, newPassword)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateUpdateAt provides a mock function with given fields: userID
func (_m *UserStore) UpdateUpdateAt(userID string) (int64, error) {
	ret := _m.Called(userID)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (int64, error)); ok {
		return rf(userID)
	}
	if rf, ok := ret.Get(0).(func(string) int64); ok {
		r0 = rf(userID)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// VerifyEmail provides a mock function with given fields: userID, email
func (_m *UserStore) VerifyEmail(userID string, email string) (string, error) {
	ret := _m.Called(userID, email)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (string, error)); ok {
		return rf(userID, email)
	}
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(userID, email)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(userID, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewUserStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserStore creates a new instance of UserStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserStore(t mockConstructorTestingTNewUserStore) *UserStore {
	mock := &UserStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}