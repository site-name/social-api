package model

import "github.com/Masterminds/squirrel"

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
