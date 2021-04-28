//go:generate go run layer_generators/main.go

package store

import (
	"context"

	"github.com/sitename/sitename/model"
)

type StoreResult struct {
	Data interface{}

	// NErr a temporary field used by the new code for the AppError migration. This will later become Err when the entire store is migrated.
	NErr error
}

type Store interface {
	User() UserStore
	System() SystemStore
	Job() JobStore
	Session() SessionStore
	Preference() PreferenceStore
	Token() TokenStore
	Context() context.Context
	Status() StatusStore
	Close()
	LockToMaster()
	UnlockFromMaster()
	DropAllTables()
	SetContext(context context.Context)
	Role() RoleStore
	GetDbVersion(numerical bool) (string, error)
	UserAccessToken() UserAccessTokenStore
	TermsOfService() TermsOfServiceStore
	ClusterDiscovery() ClusterDiscoveryStore
}

type ClusterDiscoveryStore interface {
	Save(discovery *model.ClusterDiscovery) error
	Delete(discovery *model.ClusterDiscovery) (bool, error)
	Exists(discovery *model.ClusterDiscovery) (bool, error)
	GetAll(discoveryType, clusterName string) ([]*model.ClusterDiscovery, error)
	SetLastPingAt(discovery *model.ClusterDiscovery) error
	Cleanup() error
}

type TermsOfServiceStore interface {
	Save(termsOfService *model.TermsOfService) (*model.TermsOfService, error)
	GetLatest(allowFromCache bool) (*model.TermsOfService, error)
	Get(id string, allowFromCache bool) (*model.TermsOfService, error)
}

type PreferenceStore interface {
	Save(preferences *model.Preferences) error
	GetCategory(userID, category string) (model.Preferences, error)
	Get(userID, category, name string) (*model.Preference, error)
	GetAll(userID string) (model.Preferences, error)
	Delete(userID, category, name string) error
	DeleteCategory(userID string, category string) error
	DeleteCategoryAndName(category string, name string) error
	PermanentDeleteByUser(userID string) error
	CleanupFlagsBatch(limit int64) (int64, error)
}

type JobStore interface {
	Save(job *model.Job) (*model.Job, error)
	UpdateOptimistically(job *model.Job, currentStatus string) (bool, error)
	UpdateStatus(id string, status string) (*model.Job, error)
	UpdateStatusOptimistically(id string, currentStatus string, newStatus string) (bool, error)
	Get(id string) (*model.Job, error)
	GetAllPage(offset int, limit int) ([]*model.Job, error)
	GetAllByType(jobType string) ([]*model.Job, error)
	GetAllByTypePage(jobType string, offset int, limit int) ([]*model.Job, error)
	GetAllByTypesPage(jobTypes []string, offset int, limit int) ([]*model.Job, error)
	GetAllByStatus(status string) ([]*model.Job, error)
	GetNewestJobByStatusAndType(status string, jobType string) (*model.Job, error)
	GetNewestJobByStatusesAndType(statuses []string, jobType string) (*model.Job, error)
	GetCountByStatusAndType(status string, jobType string) (int64, error)
	Delete(id string) (string, error)
}

type StatusStore interface {
	SaveOrUpdate(status *model.Status) error
	Get(userID string) (*model.Status, error)
	GetByIds(userIds []string) ([]*model.Status, error)
	ResetAll() error
	GetTotalActiveUsersCount() (int64, error)
	UpdateLastActivityAt(userID string, lastActivityAt int64) error
}

type UserStore interface {
	Save(user *model.User) (*model.User, error)
	Update(user *model.User, allowRoleUpdate bool) (*model.UserUpdate, error)
	UpdateLastPictureUpdate(userID string) error
	ResetLastPictureUpdate(userID string) error
	UpdatePassword(userID, newPassword string) error
	UpdateUpdateAt(userID string) (int64, error)
	UpdateAuthData(userID string, service string, authData *string, email string, resetMfa bool) (string, error)
	ResetAuthDataToEmailForUsers(service string, userIDs []string, includeDeleted bool, dryRun bool) (int, error)
	UpdateMfaSecret(userID, secret string) error
	UpdateMfaActive(userID string, active bool) error
	Get(ctx context.Context, id string) (*model.User, error)
	GetMany(ctx context.Context, ids []string) ([]*model.User, error)
	GetAll() ([]*model.User, error)
	ClearCaches()
	InvalidateProfileCacheForUser(userID string)
	GetByEmail(email string) (*model.User, error)
	GetByAuth(authData *string, authService string) (*model.User, error)
	GetAllUsingAuthService(authService string) ([]*model.User, error)
	GetAllNotInAuthService(authServices []string) ([]*model.User, error)
	GetByUsername(username string) (*model.User, error)
	GetForLogin(loginID string, allowSignInWithUsername, allowSignInWithEmail bool) (*model.User, error)
	VerifyEmail(userID, email string) (string, error)
	GetEtagForAllProfiles() string
	GetEtagForProfiles(teamID string) string
	UpdateFailedPasswordAttempts(userID string, attempts int) error
	GetSystemAdminProfiles() (map[string]*model.User, error)
	PermanentDelete(userID string) error
	GetUnreadCount(userID string) (int64, error)
	AnalyticsGetInactiveUsersCount() (int64, error)
	AnalyticsGetExternalUsers(hostDomain string) (bool, error)
	AnalyticsGetSystemAdminCount() (int64, error)
	AnalyticsGetGuestCount() (int64, error)
	ClearAllCustomRoleAssignments() error
	InferSystemInstallDate() (int64, error)
	GetAllAfter(limit int, afterID string) ([]*model.User, error)
	GetUsersBatchForIndexing(startTime, endTime int64, limit int) ([]*model.UserForIndexing, error)
	PromoteGuestToUser(userID string) error
	DemoteUserToGuest(userID string) (*model.User, error)
	DeactivateGuests() ([]string, error)
	GetKnownUsers(userID string) ([]string, error)
	Count(options model.UserCountOptions) (int64, error)
	AnalyticsActiveCountForPeriod(startTime int64, endTime int64, options model.UserCountOptions) (int64, error)
	GetAllProfiles(options *model.UserGetOptions) ([]*model.User, error)
	Search(teamID string, term string, options *model.UserSearchOptions) ([]*model.User, error)
	AnalyticsActiveCount(time int64, options model.UserCountOptions) (int64, error)

	// GetTeamGroupUsers(teamID string) ([]*model.User, error)
	// GetProfileByGroupChannelIdsForUser(userID string, channelIds []string) (map[string][]*model.User, error)
	// GetEtagForProfilesNotInTeam(teamID string) string
	// GetChannelGroupUsers(channelID string) ([]*model.User, error)
	// GetUnreadCountForChannel(userID string, channelID string) (int64, error)
	// GetAnyUnreadPostCountForChannel(userID string, channelID string) (int64, error)
	// GetRecentlyActiveUsersForTeam(teamID string, offset, limit int, viewRestrictions *model.ViewUsersRestrictions) ([]*model.User, error)
	// GetNewUsersForTeam(teamID string, offset, limit int, viewRestrictions *model.ViewUsersRestrictions) ([]*model.User, error)
	// SearchNotInTeam(notInTeamID string, term string, options *model.UserSearchOptions) ([]*model.User, error)
	// SearchInChannel(channelID string, term string, options *model.UserSearchOptions) ([]*model.User, error)
	// SearchNotInChannel(teamID string, channelID string, term string, options *model.UserSearchOptions) ([]*model.User, error)
	// SearchWithoutTeam(term string, options *model.UserSearchOptions) ([]*model.User, error)
	// SearchInGroup(groupID string, term string, options *model.UserSearchOptions) ([]*model.User, error)
	// InvalidateProfilesInChannelCacheByUser(userID string)
	// InvalidateProfilesInChannelCache(channelID string)
	// GetProfilesInChannel(options *model.UserGetOptions) ([]*model.User, error)
	// GetProfilesInChannelByStatus(options *model.UserGetOptions) ([]*model.User, error)
	// GetAllProfilesInChannel(ctx context.Context, channelID string, allowFromCache bool) (map[string]*model.User, error)
	// GetProfilesNotInChannel(teamID string, channelId string, groupConstrained bool, offset int, limit int, viewRestrictions *model.ViewUsersRestrictions) ([]*model.User, error)
	// GetProfilesWithoutTeam(options *model.UserGetOptions) ([]*model.User, error)
	// GetProfilesByUsernames(usernames []string, viewRestrictions *model.ViewUsersRestrictions) ([]*model.User, error)
	// GetProfiles(options *model.UserGetOptions) ([]*model.User, error)
	// GetProfileByIds(ctx context.Context, userIds []string, options *UserGetByIdsOpts, allowFromCache bool) ([]*model.User, error)
	// AnalyticsActiveCount(time int64, options model.UserCountOptions) (int64, error)
	// GetProfilesNotInTeam(teamID string, groupConstrained bool, offset int, limit int, viewRestrictions *model.ViewUsersRestrictions) ([]*model.User, error)
	// AutocompleteUsersInChannel(teamID, channelID, term string, options *model.UserSearchOptions) (*model.UserAutocompleteInChannel, error)
}

type SystemStore interface {
	Save(system *model.System) error
	SaveOrUpdate(system *model.System) error
	Update(system *model.System) error
	Get() (model.StringMap, error)
	GetByName(name string) (*model.System, error)
	PermanentDeleteByName(name string) (*model.System, error)
	InsertIfExists(system *model.System) (*model.System, error)
	SaveOrUpdateWithWarnMetricHandling(system *model.System) error
}

type TokenStore interface {
	Save(recovery *model.Token) error
	Delete(token string) error
	GetByToken(token string) (*model.Token, error)
	Cleanup()
	RemoveAllTokensByType(tokenType string) error
}

type SessionStore interface {
	Get(ctx context.Context, sessionIDOrToken string) (*model.Session, error)
	Save(session *model.Session) (*model.Session, error)
	GetSessions(userID string) ([]*model.Session, error)
	GetSessionsWithActiveDeviceIds(userID string) ([]*model.Session, error)
	GetSessionsExpired(thresholdMillis int64, mobileOnly bool, unnotifiedOnly bool) ([]*model.Session, error)
	UpdateExpiredNotify(sessionid string, notified bool) error
	Remove(sessionIDOrToken string) error
	RemoveAllSessions() error
	PermanentDeleteSessionsByUser(teamID string) error
	UpdateExpiresAt(sessionID string, time int64) error
	UpdateLastActivityAt(sessionID string, time int64) error
	UpdateRoles(userID string, roles string) (string, error)
	UpdateDeviceId(id string, deviceID string, expiresAt int64) (string, error)
	UpdateProps(session *model.Session) error
	AnalyticsSessionCount() (int64, error)
	Cleanup(expiryTime int64, batchSize int64)
}

type UserAccessTokenStore interface {
	Save(token *model.UserAccessToken) (*model.UserAccessToken, error)
	DeleteAllForUser(userID string) error
	Delete(tokenID string) error
	Get(tokenID string) (*model.UserAccessToken, error)
	GetAll(offset int, limit int) ([]*model.UserAccessToken, error)
	GetByToken(tokenString string) (*model.UserAccessToken, error)
	GetByUser(userID string, page, perPage int) ([]*model.UserAccessToken, error)
	Search(term string) ([]*model.UserAccessToken, error)
	UpdateTokenEnable(tokenID string) error
	UpdateTokenDisable(tokenID string) error
}

type RoleStore interface {
	Save(role *model.Role) (*model.Role, error)
	Get(roleID string) (*model.Role, error)
	GetAll() ([]*model.Role, error)
	GetByName(name string) (*model.Role, error)
	GetByNames(names []string) ([]*model.Role, error)
	Delete(roleID string) (*model.Role, error)
	PermanentDeleteAll() error
	ChannelHigherScopedPermissions(roleNames []string) (map[string]*model.RolePermissions, error)
	// AllChannelSchemeRoles returns all of the roles associated to channel schemes.
	// AllChannelSchemeRoles() ([]*model.Role, error)
	// ChannelRolesUnderTeamRole returns all of the non-deleted roles that are affected by updates to the
	// given role.
	// ChannelRolesUnderTeamRole(roleName string) ([]*model.Role, error)
	// HigherScopedPermissions retrieves the higher-scoped permissions of a list of role names. The higher-scope
	// (either team scheme or system scheme) is determined based on whether the team has a scheme or not.
}

type UserGetByIdsOpts struct {
	// IsAdmin tracks whether or not the request is being made by an administrator. Does nothing when provided by a client.
	IsAdmin bool

	// Restrict to search in a list of teams and channels. Does nothing when provided by a client.
	// ViewRestrictions *model.ViewUsersRestrictions

	// Since filters the users based on their UpdateAt timestamp.
	Since int64
}
