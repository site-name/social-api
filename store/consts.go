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

type Upsertor interface {
	Insert(list ...interface{}) error
	Update(list ...interface{}) (int64, error)
}

type Selector interface {
	SelectOne(holder interface{}, query string, args ...interface{}) error
	Select(i interface{}, query string, args ...interface{}) ([]interface{}, error)
}

type SelectUpsertor interface {
	Upsertor
	Selector
}

var TableOrderingMap map[string]string

func init() {
	TableOrderingMap = map[string]string{
		OrderLineTableName:       "CreateAt ASC",  // order
		OrderTableName:           "CreateAt DESC", //
		FulfillmentLineTableName: "",              //
		FulfillmentTableName:     "CreateAt ASC",  //
		OrderEventTableName:      "CreateAt ASC",  //

		ProductCategoryTableName:                 "",                            // product
		ProductCategoryTranslationTableName:      "LanguageCode ASC",            //
		ProductChannelListingTableName:           "CreateAt ASC",                //
		ProductCollectionChannelListingTableName: "CreateAt ASC",                //
		CollectionProductRelationTableName:       "",                            //
		ProductCollectionTableName:               "Slug ASC",                    //
		ProductCollectionTranslationTableName:    "LanguageCode ASC",            //
		ProductDigitalContentTableName:           "",                            //
		ProductDigitalContentURLTableName:        "",                            //
		ProductMediaTableName:                    "SortOrder ASC, CreateAt ASC", //
		ProductTableName:                         "Slug ASC",                    //
		ProductTranslationTableName:              "LanguageCode ASC",            //
		ProductTypeTableName:                     "Slug ASC",                    //
		ProductVariantChannelListingTableName:    "CreateAt ASC",                //
		ProductVariantMediaTableName:             "",                            //
		ProductVariantTableName:                  "SortOrder ASC, Sku ASC",      //
		ProductVariantTranslactionTableName:      "",                            //

		CheckoutLineTableName: "CreatAt ASC",                // checkout
		CheckoutTableName:     "CreatAt ASC, UpdateAt DESC", //

		ChannelTableName: "Slug ASC", //channel

		WishlistItemTableName:           "", // wishlist
		WishlistProductVariantTableName: "", //
		WishlistTableName:               "", //

		StockTableName:                 "CreateAt ASC", // warehouse
		WarehouseTableName:             "Slug DESC",    //
		WarehouseShippingZoneTableName: "",             //
		AllocationTableName:            "CreateAt ASC", //

		AddressTableName:                    "",               // account
		UserTableName:                       "Email ASC",      //
		CustomerEventTableName:              "Date ASC",       //
		StaffNotificationRecipientTableName: "StaffEmail ASC", //
		CustomerNoteTableName:               "Date ASC",       //
		TokenTableName:                      "",               //
		UserAddressTableName:                "",               //
		TermsOfServiceTableName:             "",               //
		StatusTableName:                     "",               //

		GiftcardTableName:         "Code ASC", // giftcard
		OrderGiftCardTableName:    "",         //
		GiftcardCheckoutTableName: "",         //

		PaymentTableName:     "CreateAt ASC", // payment
		TransactionTableName: "CreateAt ASC", //

		PluginKeyValueStoreTableName: "PluginKeyValueStore", //

		PreferenceTableName: "Preferences", //

		RoleTableName: "Roles", //

		CsvExportEventTablename: "ExportEvents", // event
		CsvExportFileTablename:  "ExportFiles",  //

		BaseAssignedAttributeTableName:         "BaseAssignedAttributes",                 // attribute
		AttributeTableName:                     "StorefrontSearchPosition ASC, Slug ASC", //
		AttributeTranslationTableName:          "AttributeTranslations",                  //
		AttributeValueTableName:                "AttributeValues",                        //
		AttributeValueTranslationTableName:     "AttributeValueTranslations",             //
		AssignedPageAttributeValueTableName:    "AssignedPageAttributeValues",            //
		AssignedPageAttributeTableName:         "AssignedPageAttributes",                 //
		AttributePageTableName:                 "AttributePages",                         //
		AssignedVariantAttributeValueTableName: "AssignedVariantAttributeValues",         //
		AssignedVariantAttributeTableName:      "AssignedVariantAttributes",              //
		AttributeVariantTableName:              "AttributeVariants",                      //
		AssignedProductAttributeValueTableName: "AssignedProductAttributeValues",         //
		AssignedProductAttributeTableName:      "AssignedProductAttributes",              //
		AttributeProductTableName:              "AttributeProducts",                      //

		VoucherTableName:                "Code ASC",                         // discount
		VoucherCategoryTableName:        "VoucherCategories",                //
		VoucherCollectionTableName:      "VoucherCollections",               //
		VoucherProductTableName:         "VoucherProducts",                  //
		VoucherChannelListingTableName:  "CreateAt ASC",                     //
		VoucherCustomerTableName:        "VoucherID ASC, CustomerEmail ASC", //
		SaleChannelListingTableName:     "CreateAt ASC",                     //
		SaleTableName:                   "Name ASC, CreateAt ASC",           //
		SaleTranslationTableName:        "LanguageCode ASC, Name ASC",       //
		VoucherTranslationTableName:     "LanguageCode ASC, CreateAt ASC",   //
		SaleCategoryRelationTableName:   "CreateAt ASC",                     //
		SaleProductRelationTableName:    "CreateAt ASC",                     //
		SaleCollectionRelationTableName: "CreateAt ASC",                     //
		OrderDiscountTableName:          "",                                 //

		ShopTableName:            "Shops",            // shop
		ShopTranslationTableName: "ShopTranslations", //
		ShopStaffTableName:       "ShopStaffs",       //

		MenuTableName:                "CreateAt ASC",     // menu
		MenuItemTableName:            "SortOrder ASC",    //
		MenuItemTranslationTableName: "LanguageCode ASC", //

		ShippingMethodTableName:                "",                               // shipping
		ShippingZoneTableName:                  "CreateAt ASC",                   //
		ShippingZoneChannelTableName:           "ShippingZoneChannels",           //
		ShippingMethodTranslationTableName:     "",                               //
		ShippingMethodPostalCodeRuleTableName:  "",                               //
		ShippingMethodChannelListingTableName:  "CreateAt ASC",                   //
		ShippingMethodExcludedProductTableName: "ShippingMethodExcludedProducts", //

		JobTableName: "Jobs", // job

		FileInfoTableName:      "FileInfos",      // file
		UploadSessionTableName: "UploadSessions", //

		PageTableName:            "Slug ASC",         // page
		PageTranslationtableName: "LanguageCode ASC", //
		PageTypeTableName:        "Slug ASC",         //

		InvoiceEventTableName: "CreateAt ASC", // invoice
		InvoiceTableName:      "CreateAt ASC", //

		OpenExchangeRateTableName: "ToCurrency ASC", // external services
	}
}

// all system product related table names
const (
	ProductCategoryTableName                 = "Categories"                    //
	ProductCategoryTranslationTableName      = "CategoryTranslations"          //
	ProductChannelListingTableName           = "ProductChannelListings"        //
	ProductCollectionChannelListingTableName = "CollectionChannelListings"     //
	CollectionProductRelationTableName       = "ProductCollections"            //
	ProductCollectionTableName               = "Collections"                   //
	ProductCollectionTranslationTableName    = "CollectionTranslations"        //
	ProductDigitalContentTableName           = "DigitalContents"               //
	ProductDigitalContentURLTableName        = "DigitalContentURLs"            //
	ProductMediaTableName                    = "ProductMedias"                 //
	ProductTableName                         = "Products"                      //
	ProductTranslationTableName              = "ProductTranslations"           //
	ProductTypeTableName                     = "ProductTypes"                  //
	ProductVariantChannelListingTableName    = "ProductVariantChannelListings" //
	ProductVariantMediaTableName             = "VariantMedias"                 //
	ProductVariantTableName                  = "ProductVariants"               //
	ProductVariantTranslactionTableName      = "ProductVariantTranslations"    //
)

// wishlist-related table names
const (
	WishlistItemTableName           = "WishlistItems"               //
	WishlistProductVariantTableName = "WishlistItemProductVariants" //
	WishlistTableName               = "Wishlists"                   //
)

// warehouse-related table names
const (
	StockTableName                 = "Stocks"                 //
	WarehouseTableName             = "Warehouses"             //
	WarehouseShippingZoneTableName = "WarehouseShippingZones" //
	AllocationTableName            = "Allocations"            //
)

// checkout-related table names
const (
	CheckoutLineTableName = "CheckoutLines" //
	CheckoutTableName     = "Checkouts"     //
)

// order-related table names
const (
	OrderLineTableName       = "Orderlines"       //
	OrderTableName           = "Orders"           //
	FulfillmentLineTableName = "FulfillmentLines" //
	FulfillmentTableName     = "Fulfillments"     //
	OrderEventTableName      = "OrderEvents"      //
)

// account-related table names
const (
	AddressTableName                    = "Addresses"                   //
	UserTableName                       = "Users"                       //
	CustomerEventTableName              = "CustomerEvents"              //
	StaffNotificationRecipientTableName = "StaffNotificationRecipients" //
	CustomerNoteTableName               = "CustomerNotes"               //
	TokenTableName                      = "Tokens"                      //
	UserAddressTableName                = "UserAddresses"               //
	TermsOfServiceTableName             = "TermsOfServices"             //
	StatusTableName                     = "Status"                      //
)

// channel-related table names
const (
	ChannelTableName = "Channels"
)

// giftcard-related table names
const (
	GiftcardTableName         = "GiftCards"         //
	OrderGiftCardTableName    = "OrderGiftCards"    //
	GiftcardCheckoutTableName = "GiftcardCheckouts" //
)

// payment-related table names
const (
	PaymentTableName     = "Payments"     //
	TransactionTableName = "Transactions" //
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
	CsvExportEventTablename = "ExportEvents" //
	CsvExportFileTablename  = "ExportFiles"  //
)

// attribute-related table names
const (
	BaseAssignedAttributeTableName         = "BaseAssignedAttributes"         //
	AttributeTableName                     = "Attributes"                     //
	AttributeTranslationTableName          = "AttributeTranslations"          //
	AttributeValueTableName                = "AttributeValues"                //
	AttributeValueTranslationTableName     = "AttributeValueTranslations"     //
	AssignedPageAttributeValueTableName    = "AssignedPageAttributeValues"    //
	AssignedPageAttributeTableName         = "AssignedPageAttributes"         //
	AttributePageTableName                 = "AttributePages"                 //
	AssignedVariantAttributeValueTableName = "AssignedVariantAttributeValues" //
	AssignedVariantAttributeTableName      = "AssignedVariantAttributes"      //
	AttributeVariantTableName              = "AttributeVariants"              //
	AssignedProductAttributeValueTableName = "AssignedProductAttributeValues" //
	AssignedProductAttributeTableName      = "AssignedProductAttributes"      //
	AttributeProductTableName              = "AttributeProducts"              //
)

// discount-related table names
const (
	VoucherTableName                = "Vouchers"               //
	VoucherCategoryTableName        = "VoucherCategories"      //
	VoucherCollectionTableName      = "VoucherCollections"     //
	VoucherProductTableName         = "VoucherProducts"        //
	VoucherChannelListingTableName  = "VoucherChannelListings" //
	VoucherCustomerTableName        = "VoucherCustomers"       //
	SaleChannelListingTableName     = "SaleChannelListings"    //
	SaleTableName                   = "Sales"                  //
	SaleTranslationTableName        = "SaleTranslations"       //
	VoucherTranslationTableName     = "VoucherTranslations"    //
	SaleCategoryRelationTableName   = "SaleCategories"         //
	SaleProductRelationTableName    = "SaleProducts"           //
	SaleCollectionRelationTableName = "SaleCollections"        //
	OrderDiscountTableName          = "OrderDiscounts"         //
)

// shop-related table names
const (
	ShopTableName            = "Shops"            //
	ShopTranslationTableName = "ShopTranslations" //
	ShopStaffTableName       = "ShopStaffs"       //
)

// menu-related table names
const (
	MenuTableName                = "Menus"
	MenuItemTableName            = "MenuItems"
	MenuItemTranslationTableName = "MenuItemTranslations"
)

// shipping-related table names
const (
	ShippingMethodTableName                = "ShippingMethods"                //
	ShippingZoneTableName                  = "ShippingZones"                  //
	ShippingZoneChannelTableName           = "ShippingZoneChannels"           //
	ShippingMethodTranslationTableName     = "ShippingMethodTranslations"     //
	ShippingMethodPostalCodeRuleTableName  = "ShippingMethodPostalCodeRules"  //
	ShippingMethodChannelListingTableName  = "ShippingMethodChannelListings"  //
	ShippingMethodExcludedProductTableName = "ShippingMethodExcludedProducts" //
)

// job-related table names
const (
	JobTableName = "Jobs"
)

// file-related table names
const (
	FileInfoTableName      = "FileInfos"      //
	UploadSessionTableName = "UploadSessions" //
)

// page-related table names
const (
	PageTableName            = "Pages"            //
	PageTranslationtableName = "PageTranslations" //
	PageTypeTableName        = "PageTypes"        //
)

// invoice-related table names
const (
	InvoiceEventTableName = "InvoiceEvents" // invoice
	InvoiceTableName      = "Invoices"      //
)

const (
	OpenExchangeRateTableName = "OpenExchangeRates" // external services
)
