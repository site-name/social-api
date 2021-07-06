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

// AllocationsBy is used for finding stock or order line's allocations
type AllocationsBy string

// consts to know finding allocations for stock or order line
const (
	ByStock     AllocationsBy = "stock"
	ByOrderLine AllocationsBy = "order_line"
)

// all system product related table names
const (
	ProductTypeTableName    = "ProductTypes"
	ProductTableName        = "Products"
	ProductVariantTableName = "ProductVariants"
)

// wishlist-related table names
const (
	WishlistItemTableName           = "WishlistItems"
	WishlistProductVariantTableName = "WishlistItemProductVariants"
	WishlistTableName               = "Wishlists"
)

// warehouse-related table names
const (
	StockTableName                 = "Stocks"
	WarehouseTableName             = "Warehouses"
	WarehouseShippingZoneTableName = "WarehouseShippingZones"
	AllocationTableName            = "Allocations"
)

// checkout-related table names
const (
	CheckoutLineTableName = "CheckoutLines"
	CheckoutTableName     = "Checkouts"
)

// order-related table names
const (
	OrderLineTableName = "Orderlines"
)

// account-related table names
const (
	AddressTableName = "Addresses"
	UserTableName    = "Users"
)
