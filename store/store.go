//go:generate go run layer_generators/main.go

package store

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/mattermost/gorp"
	"github.com/site-name/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/app"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/model/audit"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/cluster"
	"github.com/sitename/sitename/model/compliance"
	"github.com/sitename/sitename/model/csv"
	"github.com/sitename/sitename/model/file"
	"github.com/sitename/sitename/model/giftcard"
	"github.com/sitename/sitename/model/menu"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/model/plugins"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/model/shop"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/model/wishlist"
)

// Store is database gateway of the system
type Store interface {
	Context() context.Context                                                                                          // Context gets context
	Close()                                                                                                            // Close closes databases
	LockToMaster()                                                                                                     // LockToMaster constraints all queries to be performed on master
	UnlockFromMaster()                                                                                                 // UnlockFromMaster makes all datasources available
	DropAllTables()                                                                                                    // DropAllTables drop all tables in databases
	SetContext(context context.Context)                                                                                // set context
	GetDbVersion(numerical bool) (string, error)                                                                       // GetDbVersion returns version in use of database
	GetMaster() *gorp.DbMap                                                                                            // GetMaster get master datasource
	GetReplica() *gorp.DbMap                                                                                           // GetMaster gets slave datasource
	CommonMetaDataIndex(tableName string)                                                                              // CommonMetaDataIndex create indexes for tables that have fields `metadata` and `privatemetadata`
	CommonSeoMaxLength(table *gorp.TableMap)                                                                           // CommonSeoMaxLength is common method for settings max lengths for tables's `seotitle` and `seodescription`
	CreateIndexIfNotExists(indexName, tableName, columnName string) bool                                               // CreateIndexIfNotExists creates indexes for tables
	GetAllConns() []*gorp.DbMap                                                                                        // GetAllConns returns all datasources available in use
	GetQueryBuilder() squirrel.StatementBuilderType                                                                    // GetQueryBuilder create squirrel sql query builder
	CreateFullTextIndexIfNotExists(indexName string, tableName string, columnName string) bool                         //
	IsUniqueConstraintError(err error, indexName []string) bool                                                        //
	DBFromContext(ctx context.Context) *gorp.DbMap                                                                     //
	CreateForeignKeyIfNotExists(tableName, columnName, refTableName, refColumnName string, onDeleteCascade bool) error //
	CreateFullTextFuncIndexIfNotExists(indexName string, tableName string, function string) bool                       //
	MarkSystemRanUnitTests()

	User() UserStore                                                   // account
	Address() AddressStore                                             //
	UserTermOfService() UserTermOfServiceStore                         //
	UserAddress() UserAddressStore                                     //
	CustomerEvent() CustomerEventStore                                 //
	StaffNotificationRecipient() StaffNotificationRecipientStore       //
	CustomerNote() CustomerNoteStore                                   //
	System() SystemStore                                               // system
	Job() JobStore                                                     // job
	Session() SessionStore                                             // session
	Preference() PreferenceStore                                       // preference
	Token() TokenStore                                                 // token
	Status() StatusStore                                               // status
	Role() RoleStore                                                   // role
	UserAccessToken() UserAccessTokenStore                             // user access token
	TermsOfService() TermsOfServiceStore                               // term of service
	ClusterDiscovery() ClusterDiscoveryStore                           // cluster
	Audit() AuditStore                                                 // audit
	App() AppStore                                                     // app
	AppToken() AppTokenStore                                           //
	Channel() ChannelStore                                             // channel
	Checkout() CheckoutStore                                           // checkout
	CheckoutLine() CheckoutLineStore                                   //
	CsvExportEvent() CsvExportEventStore                               // csv
	CsvExportFile() CsvExportFileStore                                 //
	DiscountVoucher() DiscountVoucherStore                             // discount
	VoucherChannelListing() VoucherChannelListingStore                 //
	VoucherTranslation() VoucherTranslationStore                       //
	DiscountSale() DiscountSaleStore                                   //
	DiscountSaleTranslation() DiscountSaleTranslationStore             //
	DiscountSaleChannelListing() DiscountSaleChannelListingStore       //
	OrderDiscount() OrderDiscountStore                                 //
	VoucherCategory() VoucherCategoryStore                             //
	VoucherCollection() VoucherCollectionStore                         //
	VoucherProduct() VoucherProductStore                               //
	VoucherCustomer() VoucherCustomerStore                             //
	GiftCard() GiftCardStore                                           // giftcard
	GiftCardOrder() GiftCardOrderStore                                 //
	GiftCardCheckout() GiftCardCheckoutStore                           //
	InvoiceEvent() InvoiceEventStore                                   // invoice
	Menu() MenuStore                                                   // menu
	MenuItem() MenuItemStore                                           //
	MenuItemTranslation() MenuItemTranslationStore                     //
	Fulfillment() FulfillmentStore                                     // order
	FulfillmentLine() FulfillmentLineStore                             //
	OrderEvent() OrderEventStore                                       //
	Order() OrderStore                                                 //
	OrderLine() OrderLineStore                                         //
	Page() PageStore                                                   // page
	PageType() PageTypeStore                                           //
	PageTranslation() PageTranslationStore                             //
	Payment() PaymentStore                                             // payment
	PaymentTransaction() PaymentTransactionStore                       //
	Category() CategoryStore                                           // product
	CategoryTranslation() CategoryTranslationStore                     //
	ProductType() ProductTypeStore                                     //
	Product() ProductStore                                             //
	ProductTranslation() ProductTranslationStore                       //
	ProductChannelListing() ProductChannelListingStore                 //
	ProductVariant() ProductVariantStore                               //
	ProductVariantTranslation() ProductVariantTranslationStore         //
	ProductVariantChannelListing() ProductVariantChannelListingStore   //
	DigitalContent() DigitalContentStore                               //
	DigitalContentUrl() DigitalContentUrlStore                         //
	ProductMedia() ProductMediaStore                                   //
	VariantMedia() VariantMediaStore                                   //
	CollectionProduct() CollectionProductStore                         //
	Collection() CollectionStore                                       //
	CollectionChannelListing() CollectionChannelListingStore           //
	CollectionTranslation() CollectionTranslationStore                 //
	ShippingMethodTranslation() ShippingMethodTranslationStore         // shipping
	ShippingMethodChannelListing() ShippingMethodChannelListingStore   //
	ShippingMethodPostalCodeRule() ShippingMethodPostalCodeRuleStore   //
	ShippingMethod() ShippingMethodStore                               //
	ShippingZone() ShippingZoneStore                                   //
	ShippingZoneChannel() ShippingZoneChannelStore                     //
	Warehouse() WarehouseStore                                         // warehouse
	Stock() StockStore                                                 //
	Allocation() AllocationStore                                       //
	WarehouseShippingZone() WarehouseShippingZoneStore                 //
	Wishlist() WishlistStore                                           // wishlist
	WishlistItem() WishlistItemStore                                   //
	WishlistProductVariant() WishlistProductVariantStore               //
	PluginConfiguration() PluginConfigurationStore                     // plugin
	Compliance() ComplianceStore                                       // Compliance
	Attribute() AttributeStore                                         // attribute
	AttributeTranslation() AttributeTranslationStore                   //
	AttributeValue() AttributeValueStore                               //
	AttributeValueTranslation() AttributeValueTranslationStore         //
	AssignedPageAttributeValue() AssignedPageAttributeValueStore       //
	AssignedPageAttribute() AssignedPageAttributeStore                 //
	AttributePage() AttributePageStore                                 //
	AssignedVariantAttributeValue() AssignedVariantAttributeValueStore //
	AssignedVariantAttribute() AssignedVariantAttributeStore           //
	AttributeVariant() AttributeVariantStore                           //
	AssignedProductAttributeValue() AssignedProductAttributeValueStore //
	AssignedProductAttribute() AssignedProductAttributeStore           //
	AttributeProduct() AttributeProductStore                           //
	FileInfo() FileInfoStore                                           // upload session
	UploadSession() UploadSessionStore                                 //
	Plugin() PluginStore                                               //
	Shop() ShopStore                                                   // shop
	ShopTranslation() ShopTranslationStore                             //
	ShopStaff() ShopStaffStore                                         //
}

// shop
type (
	ShopStaffStore interface {
		CreateIndexesIfNotExists()
		Save(shopStaff *shop.ShopStaffRelation) (*shop.ShopStaffRelation, error)             // Save inserts given shopStaff into database then returns it with an error
		Get(shopStaffID string) (*shop.ShopStaffRelation, error)                             // Get finds a shop staff with given id then returns it with an error
		FilterByShopAndStaff(shopID string, staffID string) (*shop.ShopStaffRelation, error) // FilterByShopAndStaff finds a relation ship with given shopId and staffId
	}
	ShopStore interface {
		CreateIndexesIfNotExists()
		Upsert(shop *shop.Shop) (*shop.Shop, error) // Upsert depends on shop's Id to decide to update/insert the given shop.
		Get(shopID string) (*shop.Shop, error)      // Get finds a shop with given id and returns it
	}
	ShopTranslationStore interface {
		CreateIndexesIfNotExists()
		Upsert(translation *shop.ShopTranslation) (*shop.ShopTranslation, error) // Upsert depends on translation's Id then decides to update or insert
		Get(id string) (*shop.ShopTranslation, error)                            // Get finds a shop translation with given id then return it with an error
	}
)

// Plugin
type PluginStore interface {
	CreateIndexesIfNotExists()
	SaveOrUpdate(keyVal *plugins.PluginKeyValue) (*plugins.PluginKeyValue, error)
	CompareAndSet(keyVal *plugins.PluginKeyValue, oldValue []byte) (bool, error)
	CompareAndDelete(keyVal *plugins.PluginKeyValue, oldValue []byte) (bool, error)
	SetWithOptions(pluginID string, key string, value []byte, options plugins.PluginKVSetOptions) (bool, error)
	Get(pluginID, key string) (*plugins.PluginKeyValue, error)
	Delete(pluginID, key string) error
	DeleteAllForPlugin(PluginID string) error
	DeleteAllExpired() error
	List(pluginID string, page, perPage int) ([]string, error)
}

type UploadSessionStore interface {
	CreateIndexesIfNotExists()
	Save(session *file.UploadSession) (*file.UploadSession, error)
	Update(session *file.UploadSession) error
	Get(id string) (*file.UploadSession, error)
	GetForUser(userID string) ([]*file.UploadSession, error)
	Delete(id string) error
}

// fileinfo
type FileInfoStore interface {
	CreateIndexesIfNotExists()
	Save(info *file.FileInfo) (*file.FileInfo, error)
	Upsert(info *file.FileInfo) (*file.FileInfo, error)
	Get(id string) (*file.FileInfo, error)
	GetFromMaster(id string) (*file.FileInfo, error)
	GetByIds(ids []string) ([]*file.FileInfo, error)
	GetByPath(path string) (*file.FileInfo, error)
	GetForUser(userID string) ([]*file.FileInfo, error)
	GetWithOptions(page, perPage int, opt *file.GetFileInfosOptions) ([]*file.FileInfo, error)
	InvalidateFileInfosForPostCache(postID string, deleted bool)
	PermanentDelete(fileID string) error
	PermanentDeleteBatch(endTime int64, limit int64) (int64, error)
	PermanentDeleteByUser(userID string) (int64, error)
	SetContent(fileID, content string) error
	ClearCaches()
	CountAll() (int64, error)

	// Search(paramsList []*model.SearchParams, userID, teamID string, page, perPage int) (*model.FileInfoList, error)
	// GetFilesBatchForIndexing(startTime, endTime int64, limit int) ([]*model.FileForIndexing, error)
	// AttachToPost(fileID string, postID string, creatorID string) error
	// DeleteForPost(postID string) (string, error)
	// GetForPost(postID string, readFromMaster, includeDeleted, allowFromCache bool) ([]*model.FileInfo, error)
}

// attribute
type (
	AttributeStore interface {
		CreateIndexesIfNotExists()
		Save(attr *attribute.Attribute) (*attribute.Attribute, error)    // Save insert given attribute into database then returns it with an error. Returned can be wither *AppError or *NewErrInvalidInput or system error
		Get(id string) (*attribute.Attribute, error)                     // Get try finding an attribute with given id then returns it with an error. Returned error can be either *store.ErrNotFound or system error
		GetAttributesByIds(ids []string) ([]*attribute.Attribute, error) // GetAttributesByIds try finding all attributes with given `ids` then returns them. Returned error can be wither *store.ErrNotFound or system error
		GetProductAndVariantHeaders(ids []string) ([]string, error)      // GetProductAndVariantHeaders
		GetBySlug(slug string) (*attribute.Attribute, error)             // GetBySlug finds an attribute with given slug, then returns it with an error. Returned error can be wither *ErrNotFound or system error
	}
	AttributeTranslationStore interface {
		CreateIndexesIfNotExists()
	}
	AttributeValueStore interface {
		CreateIndexesIfNotExists()
		Save(attribute *attribute.AttributeValue) (*attribute.AttributeValue, error) // Save inserts given attribute value into database, then returns inserted value and an error
		Get(attributeID string) (*attribute.AttributeValue, error)                   // Get finds an attribute value with given id then returns it with an error
		GetAllByAttributeID(attributeID string) ([]*attribute.AttributeValue, error) // GetAllByAttributeID finds all attribute values that belong to given attribute then returns them withh an error
	}
	AttributeValueTranslationStore interface {
		CreateIndexesIfNotExists()
	}
	AssignedPageAttributeValueStore interface {
		CreateIndexesIfNotExists()
		Save(assignedPageAttrValue *attribute.AssignedPageAttributeValue) (*attribute.AssignedPageAttributeValue, error)                                                 // Save insert given value into database then returns it with an error
		Get(assignedPageAttrValueID string) (*attribute.AssignedPageAttributeValue, error)                                                                               // Get try finding an value with given id then returns it with an error
		SaveInBulk(assignmentID string, attributeValueIDs []string) ([]*attribute.AssignedPageAttributeValue, error)                                                     // SaveInBulk inserts multiple values into database then returns them with an error
		SelectForSort(assignmentID string) (assignedPageAttributeValues []*attribute.AssignedPageAttributeValue, attributeValues []*attribute.AttributeValue, err error) // SelectForSort uses inner join to find two list: []*assignedPageAttributeValue and []*attributeValue. With given assignedPageAttributeID
		UpdateInBulk(attributeValues []*attribute.AssignedPageAttributeValue) error                                                                                      // UpdateInBulk use transaction to update all given assigned page attribute values
	}
	AssignedPageAttributeStore interface {
		CreateIndexesIfNotExists()
		Save(assignedPageAttr *attribute.AssignedPageAttribute) (*attribute.AssignedPageAttribute, error)          // Save inserts given assigned page attribute into database and returns it with an error
		Get(id string) (*attribute.AssignedPageAttribute, error)                                                   // Get returns an assigned page attribute with an error
		GetByOption(option *attribute.AssignedPageAttributeFilterOption) (*attribute.AssignedPageAttribute, error) // GetByOption try to find an assigned page attribute with given option. If nothing found, creats new instance with that option and returns such value with an error
	}
	AttributePageStore interface {
		CreateIndexesIfNotExists()
		Save(page *attribute.AttributePage) (*attribute.AttributePage, error)
		Get(pageID string) (*attribute.AttributePage, error)
		GetByOption(option *attribute.AttributePageFilterOption) (*attribute.AttributePage, error)
	}
	AssignedVariantAttributeValueStore interface {
		CreateIndexesIfNotExists()
		Save(assignedVariantAttrValue *attribute.AssignedVariantAttributeValue) (*attribute.AssignedVariantAttributeValue, error)                                              // Save inserts new value into database then returns it with an error
		Get(assignedVariantAttrValueID string) (*attribute.AssignedVariantAttributeValue, error)                                                                               // Get try finding a value with given id then returns it with an error
		SaveInBulk(assignmentID string, attributeValueIDs []string) ([]*attribute.AssignedVariantAttributeValue, error)                                                        // SaveInBulk save multiple values into database then returns them
		SelectForSort(assignmentID string) (assignedVariantAttributeValues []*attribute.AssignedVariantAttributeValue, attributeValues []*attribute.AttributeValue, err error) // SelectForSort
		UpdateInBulk(attributeValues []*attribute.AssignedVariantAttributeValue) error                                                                                         // UpdateInBulk use transaction to update given values, then returns an error to indicate if the operation was successful or not
	}
	AssignedVariantAttributeStore interface {
		CreateIndexesIfNotExists()
		Save(assignedVariantAttribute *attribute.AssignedVariantAttribute) (*attribute.AssignedVariantAttribute, error)    // Save insert new instance into database then returns it with an error
		Get(id string) (*attribute.AssignedVariantAttribute, error)                                                        // Get find assigned variant attribute from database then returns it with an error
		GetWithOption(option *attribute.AssignedVariantAttributeFilterOption) (*attribute.AssignedVariantAttribute, error) // GetWithOption try finding an assigned variant attribute with given option. If nothing found, it creates instance with given option. Finally it returns expected value with an error
	}
	AttributeVariantStore interface {
		CreateIndexesIfNotExists()
		Save(attributeVariant *attribute.AttributeVariant) (*attribute.AttributeVariant, error)
		Get(attributeVariantID string) (*attribute.AttributeVariant, error)
		GetByOption(option *attribute.AttributeVariantFilterOption) (*attribute.AttributeVariant, error) // GetByOption finds 1 attribute variant with given option.
	}
	AssignedProductAttributeValueStore interface {
		CreateIndexesIfNotExists()
		Save(assignedProductAttrValue *attribute.AssignedProductAttributeValue) (*attribute.AssignedProductAttributeValue, error) // Save inserts given instance into database then returns it with an error
		Get(assignedProductAttrValueID string) (*attribute.AssignedProductAttributeValue, error)                                  // Get try finding an instance with given id then returns the value with an error
		SaveInBulk(assignmentID string, attributeValueIDs []string) ([]*attribute.AssignedProductAttributeValue, error)           // SaveInBulk save multiple values into database
		SelectForSort(assignmentID string) ([]*attribute.AssignedProductAttributeValue, []*attribute.AttributeValue, error)       // SelectForSort finds all `*AssignedProductAttributeValue` and related `*AttributeValues` with given `assignmentID`, then returns them with an error.
		UpdateInBulk(attributeValues []*attribute.AssignedProductAttributeValue) error                                            // UpdateInBulk use transaction to update the given values. Returned error can be `*store.ErrInvalidInput` or `system error`
	}
	AssignedProductAttributeStore interface {
		CreateIndexesIfNotExists()
		Save(assignedProductAttribute *attribute.AssignedProductAttribute) (*attribute.AssignedProductAttribute, error)    // Save inserts new assgignedProductAttribute into database and returns it with an error
		Get(id string) (*attribute.AssignedProductAttribute, error)                                                        // Get finds and returns an assignedProductAttribute with en error
		GetWithOption(option *attribute.AssignedProductAttributeFilterOption) (*attribute.AssignedProductAttribute, error) // GetWithOption try finding an `AssignedProductAttribute` with given `option`. If nothing found, it creates new instance then returns it with an error
	}
	AttributeProductStore interface {
		CreateIndexesIfNotExists()
		Save(attributeProduct *attribute.AttributeProduct) (*attribute.AttributeProduct, error)       // Save inserts given attribute product relationship into database then returns it and an error
		Get(attributeProductID string) (*attribute.AttributeProduct, error)                           // Get finds an attributeProduct relationship and returns it with an error
		GetByOption(option *attribute.AttributeProductGetOption) (*attribute.AttributeProduct, error) // GetByOption returns an attributeProduct with given condition
	}
)

// compliance
type ComplianceStore interface {
	CreateIndexesIfNotExists()
	Save(compliance *compliance.Compliance) (*compliance.Compliance, error)
	Update(compliance *compliance.Compliance) (*compliance.Compliance, error)
	Get(id string) (*compliance.Compliance, error)
	GetAll(offset, limit int) (compliance.Compliances, error)
	ComplianceExport(compliance *compliance.Compliance, cursor compliance.ComplianceExportCursor, limit int) ([]*compliance.CompliancePost, compliance.ComplianceExportCursor, error)
	MessageExport(cursor compliance.MessageExportCursor, limit int) ([]*compliance.MessageExport, compliance.MessageExportCursor, error)
}

//plugin
type PluginConfigurationStore interface {
	CreateIndexesIfNotExists()
}

// wishlist
type (
	WishlistStore interface {
		CreateIndexesIfNotExists()
		Save(wishlist *wishlist.Wishlist) (*wishlist.Wishlist, error) // Save inserts new wishlist into database
		GetById(id string) (*wishlist.Wishlist, error)                // GetById returns a wishlist with given id
		GetByUserID(userID string) (*wishlist.Wishlist, error)        // GetByUserID returns a wishlist belong to given user
	}
	WishlistItemStore interface {
		CreateIndexesIfNotExists()
		Save(wishlistItem *wishlist.WishlistItem) (*wishlist.WishlistItem, error)      // Save insert new wishlist item into database
		GetById(id string) (*wishlist.WishlistItem, error)                             // GetById returns a wishlist item wish given id
		WishlistItemsByWishlistId(wishlistID string) ([]*wishlist.WishlistItem, error) // WishlistItemsByWishlistId returns a list of wishlist items that belong to given wishlist
	}
	WishlistProductVariantStore interface {
		CreateIndexesIfNotExists()
		Save(wishlistVariant *wishlist.WishlistProductVariant) (*wishlist.WishlistProductVariant, error) // Save inserts new wishlist product variant relation into database and returns it
		GetById(id string) (*wishlist.WishlistProductVariant, error)                                     // GetByID returns a wishlist item product variant with given id
	}
)

// warehouse
type (
	WarehouseStore interface {
		CreateIndexesIfNotExists()
		ModelFields() []string
		Save(warehouse *warehouse.WareHouse) (*warehouse.WareHouse, error) // Save inserts given warehouse into database then returns it.
		Get(id string) (*warehouse.WareHouse, error)                       // Get try findings warehouse with given id, returns it. returned error could be wither (nil, *ErrNotFound, error)
		GetWarehousesHeaders(ids []string) ([]string, error)               // GetWarehousesHeaders
	}
	StockStore interface {
		CreateIndexesIfNotExists()
		Save(stock *warehouse.Stock) (*warehouse.Stock, error) // Save inserts given stock into database and returns it. Returned error could be either (nil, *AppError, *InvalidInput)
		Get(stockID string) (*warehouse.Stock, error)          // Get finds and returns stock with given stockID. Returned error could be either (nil, *ErrNotFound, error)
		// FilterVariantStocksForCountry can returns error with type of either: (nil, *ErrNotfound, *ErrInvalidParam, server lookup error)
		FilterVariantStocksForCountry(options *warehouse.ForCountryAndChannelFilter, productVariantID string) ([]*warehouse.Stock, []*warehouse.WareHouse, []*product_and_discount.ProductVariant, error)
		// FilterProductStocksForCountryAndChannel can returns error with type of either: (nil, *ErrNotFound, *ErrinvalidParam, server lookup error)
		FilterProductStocksForCountryAndChannel(options *warehouse.ForCountryAndChannelFilter, productID string) ([]*warehouse.Stock, []*warehouse.WareHouse, []*product_and_discount.ProductVariant, error)
		FilterForCountryAndChannel(options *warehouse.ForCountryAndChannelFilter) ([]*warehouse.Stock, []*warehouse.WareHouse, []*product_and_discount.ProductVariant, error) // FilterForCountryAndChannel
	}
	AllocationStore interface {
		CreateIndexesIfNotExists()
		Save(allocation *warehouse.Allocation) (*warehouse.Allocation, error)                                        // Save inserts new allocation into database and returns it
		Get(allocationID string) (*warehouse.Allocation, error)                                                      // Get find and returns allocation with given id
		AllocationsByWhich(parentID string, toWhich warehouse.AllocationsBy) ([]*warehouse.Allocation, error)        // AllocationsByWhich finds all allocations that belong to given order line or stock
		AllocationsByParentIDs(parentIDs []string, toWhich warehouse.AllocationsBy) ([]*warehouse.Allocation, error) // AllocationsByParentIDs is similar to AllocationsByWhich but it finds for all given parent ids, not just one
	}
	WarehouseShippingZoneStore interface {
		CreateIndexesIfNotExists()
	}
)

// shipping
type (
	ShippingZoneStore interface {
		CreateIndexesIfNotExists()
	}
	ShippingMethodStore interface {
		CreateIndexesIfNotExists()
		ModelFields() []string
		Upsert(method *shipping.ShippingMethod) (*shipping.ShippingMethod, error)                                // Upsert bases on given method's Id to decide update or insert it
		Get(methodID string) (*shipping.ShippingMethod, error)                                                   // Get finds and returns a shipping method with given id
		ShippingMethodsByOption(option *shipping.ShippingMethodFilterOption) ([]*shipping.ShippingMethod, error) // ShippingMethodsByOption finds and returns a list of shipping methods that satisfy given filtering option
	}
	ShippingMethodPostalCodeRuleStore interface {
		CreateIndexesIfNotExists()
	}
	ShippingMethodChannelListingStore interface {
		CreateIndexesIfNotExists()
	}
	ShippingMethodTranslationStore interface {
		CreateIndexesIfNotExists()
	}
	ShippingZoneChannelStore interface {
		CreateIndexesIfNotExists()
	}
)

// product
type (
	CollectionTranslationStore interface {
		CreateIndexesIfNotExists()
	}
	CollectionChannelListingStore interface {
		CreateIndexesIfNotExists()
	}
	CollectionStore interface {
		CreateIndexesIfNotExists()
	}
	CollectionProductStore interface {
		CreateIndexesIfNotExists()
	}
	VariantMediaStore interface {
		CreateIndexesIfNotExists()
	}
	ProductMediaStore interface {
		CreateIndexesIfNotExists()
	}
	DigitalContentUrlStore interface {
		CreateIndexesIfNotExists()
	}
	DigitalContentStore interface {
		CreateIndexesIfNotExists()
	}
	ProductVariantChannelListingStore interface {
		CreateIndexesIfNotExists()
	}
	ProductVariantTranslationStore interface {
		CreateIndexesIfNotExists()
	}
	ProductVariantStore interface {
		CreateIndexesIfNotExists()
		ModelFields() []string
		Save(variant *product_and_discount.ProductVariant) (*product_and_discount.ProductVariant, error) // Save inserts product variant instance to database
		Get(id string) (*product_and_discount.ProductVariant, error)                                     // Get returns a product variant with given id
	}
	ProductChannelListingStore interface {
		CreateIndexesIfNotExists()
		Save(channelListing *product_and_discount.ProductChannelListing) (*product_and_discount.ProductChannelListing, error)                 // Save inserts given product channel listing into database then returns it with an error
		Get(channelListingID string) (*product_and_discount.ProductChannelListing, error)                                                     // Get try finding a product channel listing, then returns it with an error
		FilterByOption(option *product_and_discount.ProductChannelListingFilterOption) ([]*product_and_discount.ProductChannelListing, error) // FilterByOption filter a list of product channel listings by given option. Then returns them with an error
	}
	ProductTranslationStore interface {
		CreateIndexesIfNotExists()
	}
	ProductTypeStore interface {
		CreateIndexesIfNotExists()
		Save(productType *product_and_discount.ProductType) (*product_and_discount.ProductType, error)    // Save try inserting new product type into database then returns it
		Get(productTypeID string) (*product_and_discount.ProductType, error)                              // Get try finding product type with given id and returns it
		FilterProductTypesByCheckoutID(checkoutToken string) ([]*product_and_discount.ProductType, error) // FilterProductTypesByCheckoutID is used to check if a checkout requires shipping
		ProductTypesByProductIDs(productIDs []string) ([]*product_and_discount.ProductType, error)        // ProductTypesByProductIDs returns all product types belong to given products
	}
	CategoryTranslationStore interface {
		CreateIndexesIfNotExists()
	}
	CategoryStore interface {
		CreateIndexesIfNotExists()
	}
	ProductStore interface {
		CreateIndexesIfNotExists()
		ModelFields() []string
		Save(prd *product_and_discount.Product) (*product_and_discount.Product, error)
		Get(id string) (*product_and_discount.Product, error)
		GetProductsByIds(ids []string) ([]*product_and_discount.Product, error)
	}
)

// payment
type (
	PaymentStore interface {
		CreateIndexesIfNotExists()
		Save(payment *payment.Payment) (*payment.Payment, error)                                // Save save payment instance into database
		Get(id string) (*payment.Payment, error)                                                // Get returns a payment with given id
		GetPaymentsByOrderID(orderID string) ([]*payment.Payment, error)                        // GetPaymentsByOrderID returns all payments that belong to given order
		PaymentExistWithOptions(opts *payment.PaymentFilterOpts) (paymentExist bool, err error) // FilterWithOptions filter order's payments based on given options
		GetPaymentsByCheckoutID(checkoutID string) ([]*payment.Payment, error)                  // GetPaymentsByCheckoutID returns all payments belong to given checkout
		Update(payment *payment.Payment) (*payment.Payment, error)                              // Update updates given payment and returns new updated payment
	}
	PaymentTransactionStore interface {
		CreateIndexesIfNotExists()
		Save(transaction *payment.PaymentTransaction) (*payment.PaymentTransaction, error)   // Save inserts new payment transaction into database
		Get(id string) (*payment.PaymentTransaction, error)                                  // Get returns a payment transaction with given id
		GetAllByPaymentID(paymentID string) ([]*payment.PaymentTransaction, error)           // GetAllByPaymentID returns a slice of payment transaction(s) that belong to given payment
		Update(transaction *payment.PaymentTransaction) (*payment.PaymentTransaction, error) // Update updates given transaction and returns updated one
	}
)

// page
type (
	PageTypeStore interface {
		CreateIndexesIfNotExists()
	}
	PageTranslationStore interface {
		CreateIndexesIfNotExists()
	}
	PageStore interface {
		CreateIndexesIfNotExists()
	}
)

// order
type (
	OrderLineStore interface {
		CreateIndexesIfNotExists()
		ModelFields() []string
		Save(orderLine *order.OrderLine) (*order.OrderLine, error)  // Save save given order line instance into database and returns it
		Get(id string) (*order.OrderLine, error)                    // Get returns a order line with id of given id
		GetAllByOrderID(orderID string) ([]*order.OrderLine, error) // GetAllByOrderID returns a slice of order lines that belong to given order
		// OrderLinesByOrderWithPrefetch finds order lines belong to given order
		//
		// and preload `variants`, `products` related to these order lines
		//
		// this borrow the idea from Django's prefetch_related() method
		OrderLinesByOrderWithPrefetch(orderID string) ([]*order.OrderLine, []*product_and_discount.ProductVariant, []*product_and_discount.Product, error)
	}
	OrderStore interface {
		CreateIndexesIfNotExists()
		Save(order *order.Order) (*order.Order, error)                       // Save insert an order into database and returns that order if success
		Get(id string) (*order.Order, error)                                 // Get find order in database with given id
		Update(order *order.Order) (*order.Order, error)                     // Update update order
		UpdateTotalPaid(orderId string, newTotalPaid *decimal.Decimal) error // updateTotalPaid update total paid amount of given order
	}
	OrderEventStore interface {
		CreateIndexesIfNotExists()
	}
	FulfillmentLineStore interface {
		CreateIndexesIfNotExists()
		Save(fulfillmentLine *order.FulfillmentLine) (*order.FulfillmentLine, error)
		Get(id string) (*order.FulfillmentLine, error)
	}
	FulfillmentStore interface {
		CreateIndexesIfNotExists()
		Save(fulfillment *order.Fulfillment) (*order.Fulfillment, error)
		Get(id string) (*order.Fulfillment, error)
		FilterByExcludeStatuses(orderID string, excludeStatuses []string) (exist bool, err error) // FilterByExcludeStatuses check if there is at least 1 fulfillment belong to given order and have status differnt than given statuses.
	}
)

// menu
type (
	MenuItemTranslationStore interface {
		CreateIndexesIfNotExists()
	}
	MenuStore interface {
		CreateIndexesIfNotExists()
		Save(menu *menu.Menu) (*menu.Menu, error)  // Save insert given menu into database and returns it
		GetById(id string) (*menu.Menu, error)     // GetById returns a menu with given id
		GetByName(name string) (*menu.Menu, error) // GetByName returns a menu with given name
		GetBySlug(slug string) (*menu.Menu, error) // GetBySlug returns a menu with given slug
	}
	MenuItemStore interface {
		CreateIndexesIfNotExists()
		Save(menuItem *menu.MenuItem) (*menu.MenuItem, error) // Save insert given menu item into database and returns it
		GetById(id string) (*menu.MenuItem, error)            // GetById returns a menu item with given id
		GetByName(name string) (*menu.MenuItem, error)        // GetByName returns a menu item with given name
	}
)

type InvoiceEventStore interface {
	CreateIndexesIfNotExists()
}

// giftcard related stores
type (
	GiftCardStore interface {
		CreateIndexesIfNotExists()
		Upsert(giftCard *giftcard.GiftCard) (*giftcard.GiftCard, error)                     // Upsert depends on given giftcard's Id property then perform according operation
		GetById(id string) (*giftcard.GiftCard, error)                                      // GetById returns a giftcard instance that has id of given id
		GetAllByUserId(userID string) ([]*giftcard.GiftCard, error)                         // GetAllByUserId returns a slice aff giftcards that belong to given user
		GetAllByCheckout(checkoutID string) ([]*giftcard.GiftCard, error)                   // GetAllByCheckout returns all giftcards belong to given checkout
		FilterByOption(option *giftcard.GiftCardFilterOption) ([]*giftcard.GiftCard, error) // FilterByOption finds giftcards wth option
	}
	GiftCardOrderStore interface {
		CreateIndexesIfNotExists()
		Save(giftcardOrder *giftcard.OrderGiftCard) (*giftcard.OrderGiftCard, error) // Save inserts new giftcard-order relation into database then returns it
		Get(id string) (*giftcard.OrderGiftCard, error)                              // Get returns giftcard-order relation table with given id
	}
	GiftCardCheckoutStore interface {
		CreateIndexesIfNotExists()
		Save(giftcardOrder *giftcard.GiftCardCheckout) (*giftcard.GiftCardCheckout, error) // Save inserts new giftcard-checkout relation into database then returns it
		Get(id string) (*giftcard.GiftCardCheckout, error)                                 // Get returns giftcard-checkout relation table with given id
		Delete(giftcardID string, checkoutID string) error                                 // Delete deletes a giftcard-checkout relation with given id
	}
)

// discount
type (
	OrderDiscountStore interface {
		CreateIndexesIfNotExists()
	}
	DiscountSaleTranslationStore interface {
		CreateIndexesIfNotExists()
	}
	DiscountSaleChannelListingStore interface {
		CreateIndexesIfNotExists()
	}
	VoucherTranslationStore interface {
		CreateIndexesIfNotExists()
	}
	DiscountSaleStore interface {
		CreateIndexesIfNotExists()
		Upsert(sale *product_and_discount.Sale) (*product_and_discount.Sale, error)                              // Upsert bases on sale's Id to decide to update or insert given sale
		Get(saleID string) (*product_and_discount.Sale, error)                                                   // Get finds and returns a sale with given saleID
		FilterSalesByOption(option *product_and_discount.SaleFilterOption) ([]*product_and_discount.Sale, error) // FilterSalesByOption filter sales by option
	}
	VoucherChannelListingStore interface {
		CreateIndexesIfNotExists()
		Upsert(voucherChannelListing *product_and_discount.VoucherChannelListing) (*product_and_discount.VoucherChannelListing, error) // upsert check given listing's Id to decide whether to create or update it. Then returns a listing with an error
		Get(voucherChannelListingID string) (*product_and_discount.VoucherChannelListing, error)                                       // Get finds a listing with given id, then returns it with an error
		FilterByVoucherAndChannel(voucherID string, channelID string) ([]*product_and_discount.VoucherChannelListing, error)           // FilterByVoucherAndChannel finds a list of listings that belong to given voucher and own given channel
	}
	DiscountVoucherStore interface {
		CreateIndexesIfNotExists()
		Upsert(voucher *product_and_discount.Voucher) (*product_and_discount.Voucher, error)                              // Upsert saves or updates given voucher then returns it with an error
		Get(voucherID string) (*product_and_discount.Voucher, error)                                                      // Get finds a voucher with given id, then returns it with an error
		FilterVouchersByOption(option *product_and_discount.VoucherFilterOption) ([]*product_and_discount.Voucher, error) // FilterVouchersByOption finds vouchers bases on given option.
	}
	VoucherCategoryStore interface {
		CreateIndexesIfNotExists()
		Upsert(voucherCategory *product_and_discount.VoucherCategory) (*product_and_discount.VoucherCategory, error) // Upsert saves or updates given voucher category then returns it with an error
		Get(voucherCategoryID string) (*product_and_discount.VoucherCategory, error)                                 // Get finds a voucher category with given id, then returns it with an error
		ProductCategoriesByVoucherID(voucherID string) ([]*product_and_discount.Category, error)                     // ProductCategoriesByVoucherID finds a list of product categories that have relationships with given voucher
	}
	VoucherCollectionStore interface {
		CreateIndexesIfNotExists()
		Upsert(voucherCollection *product_and_discount.VoucherCollection) (*product_and_discount.VoucherCollection, error) // Upsert saves or updates given voucher collection then returns it with an error
		Get(voucherCollectionID string) (*product_and_discount.VoucherCollection, error)                                   // Get finds a voucher collection with given id, then returns it with an error
		CollectionsByVoucherID(voucherID string) ([]*product_and_discount.Collection, error)                               // CollectionsByVoucherID finds all collections that have relationships with given voucher
	}
	VoucherProductStore interface {
		CreateIndexesIfNotExists()
		Upsert(voucherProduct *product_and_discount.VoucherProduct) (*product_and_discount.VoucherProduct, error) // Upsert saves or updates given voucher product then returns it with an error
		Get(voucherProductID string) (*product_and_discount.VoucherProduct, error)                                // Get finds a voucher product with given id, then returns it with an error
		ProductsByVoucherID(voucherID string) ([]*product_and_discount.Product, error)                            // ProductsByVoucherID finds all products that have relationships with given voucher
	}
	VoucherCustomerStore interface {
		CreateIndexesIfNotExists()
		Save(voucherCustomer *product_and_discount.VoucherCustomer) (*product_and_discount.VoucherCustomer, error) // Save inserts given voucher customer instance into database ands returns it
		Get(id string) (*product_and_discount.VoucherCustomer, error)                                              // Get finds a voucher customer with given id and returns it with an error
		FilterByVoucherAndEmail(voucherID string, email string) (*product_and_discount.VoucherCustomer, error)     // FilterByVoucherAndEmail finds a voucher customer with given voucherID and customer email then returns it with an error
	}
)

// csv
type (
	CsvExportEventStore interface {
		CreateIndexesIfNotExists()
		Save(event *csv.ExportEvent) (*csv.ExportEvent, error)
	}
	CsvExportFileStore interface {
		CreateIndexesIfNotExists()
		Save(file *csv.ExportFile) (*csv.ExportFile, error)
		Get(id string) (*csv.ExportFile, error)
	}
)

// checkout
type (
	CheckoutLineStore interface {
		CreateIndexesIfNotExists()
		ModelFields() []string
		Upsert(checkoutLine *checkout.CheckoutLine) (*checkout.CheckoutLine, error)          // Upsert checks whether to update or insert given checkout line then performs according operation
		Get(id string) (*checkout.CheckoutLine, error)                                       // Get returns a checkout line with given id
		CheckoutLinesByCheckoutID(checkoutID string) ([]*checkout.CheckoutLine, error)       // CheckoutLinesByCheckoutID returns a list of checkout lines that belong to given checkout
		DeleteLines(checkoutLineIDs []string) error                                          // DeleteLines deletes all checkout lines with given uuids
		BulkUpdate(checkoutLines []*checkout.CheckoutLine) error                             // BulkUpdate receives a list of modified checkout lines, updates them in bulk.
		BulkCreate(checkoutLines []*checkout.CheckoutLine) ([]*checkout.CheckoutLine, error) // BulkCreate takes a list of raw checkout lines, save them into database then returns them fully with an error
		// CheckoutLinesByCheckoutWithPrefetch finds all checkout lines belong to given checkout
		//
		// and prefetch all related product variants, products
		//
		// this borrows the idea from Django's prefetch_related() method
		CheckoutLinesByCheckoutWithPrefetch(checkoutID string) ([]*checkout.CheckoutLine, []*product_and_discount.ProductVariant, []*product_and_discount.Product, error)
	}
	CheckoutStore interface {
		CreateIndexesIfNotExists()
		CheckoutsByUserID(userID string, channelActive bool) ([]*checkout.Checkout, error) // CheckoutsByUserID returns a list of check outs that belong to given user and have channels active
		Get(token string) (*checkout.Checkout, error)                                      // Get finds a checkout with given token (checkouts use tokens(uuids) as primary keys)
		Upsert(ckout *checkout.Checkout) (*checkout.Checkout, error)                       // Upsert depends on given checkout's Token property to decide to update or insert it
	}
)

// channel
type ChannelStore interface {
	CreateIndexesIfNotExists()
	Save(ch *channel.Channel) (*channel.Channel, error)
	Get(id string) (*channel.Channel, error)                                         // Get returns channel by given id
	GetBySlug(slug string) (*channel.Channel, error)                                 // GetBySlug returns channel by given slug
	GetChannelsByIdsAndOrder(ids []string, order string) ([]*channel.Channel, error) //
	GetRandomActiveChannel() (*channel.Channel, error)                               // GetRandomActiveChannel get an abitrary channel that is active
}

// app
type (
	AppTokenStore interface {
		CreateIndexesIfNotExists()
		Save(appToken *app.AppToken) (*app.AppToken, error)
	}
	AppStore interface {
		CreateIndexesIfNotExists()
		Save(app *app.App) (*app.App, error)
	}
)

type ClusterDiscoveryStore interface {
	CreateIndexesIfNotExists()
	Save(discovery *cluster.ClusterDiscovery) error
	Delete(discovery *cluster.ClusterDiscovery) (bool, error)
	Exists(discovery *cluster.ClusterDiscovery) (bool, error)
	GetAll(discoveryType, clusterName string) ([]*cluster.ClusterDiscovery, error)
	SetLastPingAt(discovery *cluster.ClusterDiscovery) error
	Cleanup() error
}

type AuditStore interface {
	CreateIndexesIfNotExists()
	Save(audit *audit.Audit) error
	Get(userID string, offset int, limit int) (audit.Audits, error)
	PermanentDeleteByUser(userID string) error
}

type TermsOfServiceStore interface {
	CreateIndexesIfNotExists()
	Save(termsOfService *model.TermsOfService) (*model.TermsOfService, error)
	GetLatest(allowFromCache bool) (*model.TermsOfService, error)
	Get(id string, allowFromCache bool) (*model.TermsOfService, error)
}

type PreferenceStore interface {
	CreateIndexesIfNotExists()
	Save(preferences *model.Preferences) error
	GetCategory(userID, category string) (model.Preferences, error)
	Get(userID, category, name string) (*model.Preference, error)
	GetAll(userID string) (model.Preferences, error)
	Delete(userID, category, name string) error
	DeleteCategory(userID string, category string) error
	DeleteCategoryAndName(category string, name string) error
	PermanentDeleteByUser(userID string) error
	CleanupFlagsBatch(limit int64) (int64, error)
	DeleteUnusedFeatures()
}

type JobStore interface {
	CreateIndexesIfNotExists()
	Save(job *model.Job) (*model.Job, error)
	UpdateOptimistically(job *model.Job, currentStatus string) (bool, error)
	UpdateStatus(id string, status string) (*model.Job, error)
	UpdateStatusOptimistically(id string, currentStatus string, newStatus string) (bool, error) // update job status from current status to new status
	Get(id string) (*model.Job, error)
	GetAllPage(offset int, limit int) ([]*model.Job, error)
	GetAllByType(jobType string) ([]*model.Job, error)
	GetAllByTypePage(jobType string, offset int, limit int) ([]*model.Job, error)
	GetAllByTypesPage(jobTypes []string, offset int, limit int) ([]*model.Job, error)
	GetAllByStatus(status string) ([]*model.Job, error)
	GetNewestJobByStatusAndType(status string, jobType string) (*model.Job, error)
	GetNewestJobByStatusesAndType(statuses []string, jobType string) (*model.Job, error) // GetNewestJobByStatusesAndType get 1 job from database that has status is one of given statuses, and job type is given jobType. order by created time
	GetCountByStatusAndType(status string, jobType string) (int64, error)
	Delete(id string) (string, error)
}

type StatusStore interface {
	CreateIndexesIfNotExists()
	SaveOrUpdate(status *account.Status) error
	Get(userID string) (*account.Status, error)
	GetByIds(userIds []string) ([]*account.Status, error)
	ResetAll() error
	GetTotalActiveUsersCount() (int64, error)
	UpdateLastActivityAt(userID string, lastActivityAt int64) error
}

// account stores
type (
	AddressStore interface {
		CreateIndexesIfNotExists()                                           // CreateIndexesIfNotExists creates indexes for table if needed
		Save(address *account.Address) (*account.Address, error)             // Save saves address into database
		Get(addressID string) (*account.Address, error)                      // Get returns an Address with given addressID is exist
		GetAddressesByIDs(addressesIDs []string) ([]*account.Address, error) // GetAddressesByIDs returns a slice of Addresses with given slice of id strings
		GetAddressesByUserID(userID string) ([]*account.Address, error)      // GetAddressesByUserID returns slice of addresses belong to given user
		Update(address *account.Address) (*account.Address, error)           // Update update given address and returns it
		DeleteAddresses(addressIDs []string) error                           // DeleteAddress deletes given address and returns an error
	}
	UserTermOfServiceStore interface {
		CreateIndexesIfNotExists()                                                                //
		GetByUser(userID string) (*account.UserTermsOfService, error)                             // GetByUser returns a term of service with given user id
		Save(userTermsOfService *account.UserTermsOfService) (*account.UserTermsOfService, error) // Save inserts new user term of service to database
		Delete(userID, termsOfServiceId string) error                                             // Delete deletes from database an usder term of service with given userId and term of service id
	}
	UserStore interface {
		CreateIndexesIfNotExists()                                                    //
		Save(user *account.User) (*account.User, error)                               // Save takes an user struct and save into database
		Update(user *account.User, allowRoleUpdate bool) (*account.UserUpdate, error) // Update update given user
		UpdateLastPictureUpdate(userID string) error
		ResetLastPictureUpdate(userID string) error
		UpdatePassword(userID, newPassword string) error
		UpdateUpdateAt(userID string) (int64, error)
		UpdateAuthData(userID string, service string, authData *string, email string, resetMfa bool) (string, error)
		ResetAuthDataToEmailForUsers(service string, userIDs []string, includeDeleted bool, dryRun bool) (int, error)
		UpdateMfaSecret(userID, secret string) error
		UpdateMfaActive(userID string, active bool) error
		Get(ctx context.Context, id string) (*account.User, error)
		GetMany(ctx context.Context, ids []string) ([]*account.User, error)
		GetAll() ([]*account.User, error)
		ClearCaches()
		InvalidateProfileCacheForUser(userID string) // InvalidateProfileCacheForUser
		GetByEmail(email string) (*account.User, error)
		GetByAuth(authData *string, authService string) (*account.User, error)
		GetAllUsingAuthService(authService string) ([]*account.User, error)
		GetAllNotInAuthService(authServices []string) ([]*account.User, error)
		GetByUsername(username string) (*account.User, error)
		GetForLogin(loginID string, allowSignInWithUsername, allowSignInWithEmail bool) (*account.User, error)
		VerifyEmail(userID, email string) (string, error) // VerifyEmail set EmailVerified attribute of user to true
		GetEtagForAllProfiles() string
		GetEtagForProfiles(teamID string) string
		UpdateFailedPasswordAttempts(userID string, attempts int) error
		GetSystemAdminProfiles() (map[string]*account.User, error)
		PermanentDelete(userID string) error // PermanentDelete completely delete user from the system
		AnalyticsGetInactiveUsersCount() (int64, error)
		AnalyticsGetExternalUsers(hostDomain string) (bool, error)
		AnalyticsGetSystemAdminCount() (int64, error)
		AnalyticsGetGuestCount() (int64, error)
		ClearAllCustomRoleAssignments() error
		InferSystemInstallDate() (int64, error)
		GetAllAfter(limit int, afterID string) ([]*account.User, error)
		GetUsersBatchForIndexing(startTime, endTime int64, limit int) ([]*account.UserForIndexing, error)
		PromoteGuestToUser(userID string) error
		DemoteUserToGuest(userID string) (*account.User, error)
		DeactivateGuests() ([]string, error)
		GetKnownUsers(userID string) ([]string, error)
		Count(options account.UserCountOptions) (int64, error)
		AnalyticsActiveCountForPeriod(startTime int64, endTime int64, options account.UserCountOptions) (int64, error)
		GetAllProfiles(options *account.UserGetOptions) ([]*account.User, error)
		Search(term string, options *account.UserSearchOptions) ([]*account.User, error)
		AnalyticsActiveCount(time int64, options account.UserCountOptions) (int64, error)
		GetProfileByIds(ctx context.Context, userIds []string, options *UserGetByIdsOpts, allowFromCache bool) ([]*account.User, error)
		GetProfilesByUsernames(usernames []string) ([]*account.User, error)
		GetProfiles(options *account.UserGetOptions) ([]*account.User, error)
		GetUnreadCount(userID string) (int64, error) // TODO: consider me
	}
	TokenStore interface {
		CreateIndexesIfNotExists()
		Save(recovery *model.Token) error
		Delete(token string) error
		GetByToken(token string) (*model.Token, error)
		Cleanup()
		RemoveAllTokensByType(tokenType string) error
		GetAllTokensByType(tokenType string) ([]*model.Token, error)
	}
	UserAccessTokenStore interface {
		CreateIndexesIfNotExists()
		Save(token *account.UserAccessToken) (*account.UserAccessToken, error)
		DeleteAllForUser(userID string) error
		Delete(tokenID string) error
		Get(tokenID string) (*account.UserAccessToken, error)
		GetAll(offset int, limit int) ([]*account.UserAccessToken, error)
		GetByToken(tokenString string) (*account.UserAccessToken, error)
		GetByUser(userID string, page, perPage int) ([]*account.UserAccessToken, error)
		Search(term string) ([]*account.UserAccessToken, error)
		UpdateTokenEnable(tokenID string) error
		UpdateTokenDisable(tokenID string) error
	}
	UserAddressStore interface {
		CreateIndexesIfNotExists()
		Save(userAddress *account.UserAddress) (*account.UserAddress, error)
		DeleteForUser(userID string, addressID string) error // DeleteForUser delete the relationship between user & address
	}
	CustomerEventStore interface {
		CreateIndexesIfNotExists()
		Save(customemrEvent *account.CustomerEvent) (*account.CustomerEvent, error)
		Get(id string) (*account.CustomerEvent, error)
		Count() (int64, error)
		GetEventsByUserID(userID string) ([]*account.CustomerEvent, error) // get list of customer event belongs to given id
	}
	StaffNotificationRecipientStore interface {
		CreateIndexesIfNotExists()
		Save(notificationRecipient *account.StaffNotificationRecipient) (*account.StaffNotificationRecipient, error)
		Get(id string) (*account.StaffNotificationRecipient, error)
	}
	CustomerNoteStore interface {
		CreateIndexesIfNotExists()
		Save(note *account.CustomerNote) (*account.CustomerNote, error) // Save insert given customer note into database and returns it
		Get(id string) (*account.CustomerNote, error)                   // Get find customer note with given id and returns it
	}
)

type SystemStore interface {
	CreateIndexesIfNotExists()
	Save(system *model.System) error
	SaveOrUpdate(system *model.System) error
	Update(system *model.System) error
	Get() (model.StringMap, error)
	GetByName(name string) (*model.System, error)
	PermanentDeleteByName(name string) (*model.System, error)
	InsertIfExists(system *model.System) (*model.System, error)
	SaveOrUpdateWithWarnMetricHandling(system *model.System) error
}

// session
type SessionStore interface {
	CreateIndexesIfNotExists()
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
	UpdateLastActivityAt(sessionID string, time int64) error                    // UpdateLastActivityAt
	UpdateRoles(userID string, roles string) (string, error)                    // UpdateRoles updates roles for all sessions that have userId of given userID,
	UpdateDeviceId(id string, deviceID string, expiresAt int64) (string, error) // UpdateDeviceId updates device id for sessions
	UpdateProps(session *model.Session) error                                   // UpdateProps update session's props
	AnalyticsSessionCount() (int64, error)                                      // AnalyticsSessionCount counts numbers of sessions
	Cleanup(expiryTime int64, batchSize int64)                                  // Cleanup is called periodicly to remove sessions that are expired
}

type RoleStore interface {
	CreateIndexesIfNotExists()
	Save(role *model.Role) (*model.Role, error)
	Get(roleID string) (*model.Role, error)
	GetAll() ([]*model.Role, error)
	GetByName(ctx context.Context, name string) (*model.Role, error)
	GetByNames(names []string) ([]*model.Role, error)
	Delete(roleID string) (*model.Role, error)
	PermanentDeleteAll() error
	// ChannelHigherScopedPermissions(roleNames []string) (map[string]*model.RolePermissions, error)
	// AllChannelSchemeRoles returns all of the roles associated to channel schemes.
	// AllChannelSchemeRoles() ([]*model.Role, error)
	// ChannelRolesUnderTeamRole returns all of the non-deleted roles that are affected by updates to the given role.
	// ChannelRolesUnderTeamRole(roleName string) ([]*model.Role, error)
	// HigherScopedPermissions retrieves the higher-scoped permissions of a list of role names. The higher-scope
	// (either team scheme or system scheme) is determined based on whether the team has a scheme or not.
}

type UserGetByIdsOpts struct {
	IsAdmin bool  // IsAdmin tracks whether or not the request is being made by an administrator. Does nothing when provided by a client.
	Since   int64 // Since filters the users based on their UpdateAt timestamp.
	// Restrict to search in a list of teams and channels. Does nothing when provided by a client.
	// ViewRestrictions *model.ViewUsersRestrictions
}
