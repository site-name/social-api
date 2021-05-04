package account

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
