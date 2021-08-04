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
	ProductCategoryTableName                 = "ProductCategories"
	ProductCategoryTranslationTableName      = "ProductCategoryTranslations"
	ProductChannelListingTableName           = "ProductChannelListings"
	ProductCollectionChannelListingTableName = "ProductCollectionChannelListings"
	ProductCollectionProductTableName        = "ProductCollections"
	ProductCollectionTableName               = "Collections"
	ProductCollectionTranslationTableName    = "ProductCollectionTranslations"
	ProductDigitalContentTableName           = "DigitalContents"
	ProductDigitalContentURLTableName        = "DigitalContentURLs"
	ProductMediaTableName                    = "ProductMedias"
	ProductTableName                         = "Products"
	ProductTranslationTableName              = "ProductTranslations"
	ProductTypeTableName                     = "ProductTypes"
	ProductVariantChannelListingTableName    = "ProductVariantChannelListings"
	ProductVariantMediaTbaleName             = "ProductVariantMedias"
	ProductVariantTableName                  = "ProductVariants"
	ProductVariantTranslactionTableName      = "ProductVariantTranslations"
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
	UserAddressTableName                = "UserAddresses"
	TermsOfServiceTableName             = "TermsOfServices"
	StatusTableName                     = "Status"
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

// attribute-related table names
const (
	BaseAssignedAttributeTableName         = "BaseAssignedAttributes"
	AttributeTableName                     = "Attributes"
	AttributeTranslationTableName          = "AttributeTranslations"
	AttributeValueTableName                = "AttributeValues"
	AttributeValueTranslationTableName     = "AttributeValueTranslations"
	AssignedPageAttributeValueTableName    = "AssignedPageAttributeValues"
	AssignedPageAttributeTableName         = "AssignedPageAttributes"
	AttributePageTableName                 = "AttributePages"
	AssignedVariantAttributeValueTableName = "AssignedVariantAttributeValues"
	AssignedVariantAttributeTableName      = "AssignedVariantAttributes"
	AttributeVariantTableName              = "AttributeVariants"
	AssignedProductAttributeValueTableName = "AssignedProductAttributeValues"
	AssignedProductAttributeTableName      = "AssignedProductAttributes"
	AttributeProductTableName              = "AttributeProducts"
)

// discount-related table names
const (
	VoucherTableName               = "Vouchers"
	VoucherCategoryTableName       = "VoucherCategories"
	VoucherCollectionTableName     = "VoucherCollections"
	VoucherProductTableName        = "VoucherProducts"
	VoucherChannelListingTableName = "VoucherChannelListings"
	VoucherCustomerTableName       = "VoucherCustomers"
	SaleChannelListingTableName    = "SaleChannelListings"
	SaleTableName                  = "Sales"
	SaleTranslationTableName       = "SaleTranslations"
)

// shop-related table names
const (
	ShopTableName            = "Shops"
	ShopTranslationTableName = "ShopTranslations"
	ShopStaffTableName       = "ShopStaffs"
)

// menu-related table names
const (
	MenuTableName     = "Menus"
	MenuItemTableName = "MenuItems"
)

// shipping-related table names
const (
	ShippingMethodTableName                = "ShippingMethods"
	ShippingZoneTableName                  = "ShippingZones"
	ShippingZoneChannelTableName           = "ShippingZoneChannels"
	ShippingMethodTranslationTableName     = "ShippingMethodTranslations"
	ShippingMethodPostalCodeRuleTableName  = "ShippingMethodPostalCodeRules"
	ShippingMethodChannelListingTableName  = "ShippingMethodChannelListings"
	ShippingMethodExcludedProductTableName = "ShippingMethodExcludedProducts"
)

// job-related table names
const (
	JobTableName = "Jobs"
)

// file-related table names
const (
	FileInfoTableName      = "FileInfos"
	UploadSessionTableName = "UploadSessions"
)
