package account

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
