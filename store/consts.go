package store

const (
	ChannelExistsError                  = "store.sql_channel.save_channel.exists.app_error"
	UserSearchOptionNamesOnly           = "names_only"
	UserSearchOptionNamesOnlyNoFullName = "names_only_no_full_name"
	UserSearchOptionAllNoFullName       = "all_no_full_name"
	UserSearchOptionAllowInactive       = "allow_inactive"
	FeatureTogglePrefix                 = "feature_enabled_"
	UUID_MAX_LENGTH                     = 36 // max length for all tables's Id fields, since google's uuid generates ids have length of 36
)

type StoreResult struct {
	Data interface{}
	NErr error // NErr a temporary field used by the new code for the AppError migration. This will later become Err when the entire store is migrated.
}
