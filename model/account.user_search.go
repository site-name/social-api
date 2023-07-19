package model

import "github.com/Masterminds/squirrel"

const USER_SEARCH_MAX_LIMIT = 1000
const USER_SEARCH_DEFAULT_LIMIT = 100

// UserSearch captures the parameters provided by a client for initiating a user search.
type UserSearch struct {
	Term          string   `json:"term"`
	AllowInactive bool     `json:"allow_inactive"`
	Limit         int      `json:"limit"`
	Role          string   `json:"role"`
	Roles         []string `json:"roles"`
}

// ToJson convert a User to a json string
func (u *UserSearch) ToJSON() []byte {
	return []byte(ModelToJson(u))
}

// UserSearchOptions captures internal parameters derived from the user's permissions and a
// UserSearch request.
type UserSearchOptions struct {
	// IsAdmin tracks whether or not the search is being conducted by an administrator.
	IsAdmin bool
	// AllowEmails allows search to examine the emails of users.
	AllowEmails bool
	// AllowFullNames allows search to examine the full names of users, vs. just usernames and nicknames.
	AllowFullNames bool
	// AllowInactive configures whether or not to return inactive users in the search results.
	AllowInactive bool
	// Narrows the search to the group constrained users
	// GroupConstrained bool
	// Limit limits the total number of results returned.
	Limit int
	// Filters for the given role
	Role string
	// Filters for users that have any of the given system roles
	Roles []string
}

type UsersStats struct {
	TotalUsersCount int64 `json:"total_users_count"`
}

type UserGetOptions struct {
	Inactive bool
	// Filters the active users
	Active bool
	// Filters for the given role
	Role string
	// Filters for users matching any of the given system wide roles
	Roles []string
	// Filters for users matching any of the given channel roles, must be used with InChannelId
	// Sorting option
	Sort string
	// Restrict to search in a list of teams and channels
	// ViewRestrictions *ViewUsersRestrictions
	// Page
	Page int
	// Page size
	PerPage int
}

type UserGetByIdsOptions struct {
	// Since filters the users based on their UpdateAt timestamp.
	Since int64
}

type UserFilterOptions struct {
	Id          squirrel.Sqlizer
	Email       squirrel.Sqlizer
	Username    squirrel.Sqlizer
	FirstName   squirrel.Sqlizer
	LastName    squirrel.Sqlizer
	AuthData    squirrel.Sqlizer
	AuthService squirrel.Sqlizer

	Extra squirrel.Sqlizer // support for query AND, OR

	OrderID             squirrel.Sqlizer // INNER JOIN Orders ON Orders.UserID = Users.Id WHERE Orders.Id...
	HasNoOrder          bool             // LEFT JOIN Orders ON ... WHERE Orders.UserID IS NULL
	ExcludeBoardMembers bool             // LEFT JOIN ShopStaffs ON ... WHERE ShopStaffs.StaffID IS NULL

	Limit   int
	OrderBy string
}

// Options for counting users
type UserCountOptions struct {
	// Should include users that are bots
	// IncludeBotAccounts bool
	// Should include deleted users (of any type)
	IncludeDeleted bool
	// Exclude regular users
	ExcludeRegularUsers bool
	Roles               []string
}
