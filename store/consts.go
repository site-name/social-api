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
	OrderLineTableName       = "Orderlines"
	OrderTableName           = "Orders"
	FulfillmentLineTableName = "FulfillmentLines"
	FulfillmentTableName     = "Fulfillments"
	OrderEventTableName      = "OrderEvents"
)

// account-related table names
const (
	AddressTableName                    = "Addresses"
	UserTableName                       = "Users"
	CustomerEventTableName              = "CustomerEvents"
	StaffNotificationRecipientTableName = "StaffNotificationRecipients"
	CustomerNoteTableName               = "CustomerNotes"
	TokenTableName                      = "Tokens"
)

// channel-related table names
const (
	ChannelTableName = "Channels"
)

// giftcard-related table names
const (
	GiftcardTableName         = "GiftCards"
	OrderGiftCardTableName    = "OrderGiftCards"
	GiftcardCheckoutTableName = "GiftcardCheckouts"
)

// payment-related table names
const (
	PaymentTableName     = "Payments"
	TransactionTableName = "Transactions"
)

// plugin-related table names
const (
	PluginKeyValueStoreTableName = "PluginKeyValueStore"
)

// preference table names
const PreferenceTableName = "Preferences"

// role related table names
const RoleTableName = "Roles"

// csv-related table names
const (
	CsvExportEventTablename = "ExportEvents"
	CsvExportFileTablename  = "ExportFiles"
)
