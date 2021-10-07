//go:generate go run layer_generators/main.go

package store

import (
	"context"
	"database/sql/driver"
	timemodule "time"

	"github.com/Masterminds/squirrel"
	"github.com/mattermost/gorp"
	goprices "github.com/site-name/go-prices"
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
	"github.com/sitename/sitename/model/external_services"
	"github.com/sitename/sitename/model/file"
	"github.com/sitename/sitename/model/giftcard"
	"github.com/sitename/sitename/model/invoice"
	"github.com/sitename/sitename/model/menu"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/model/plugins"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/model/shop"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/model/wishlist"
	"github.com/sitename/sitename/modules/measurement"
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
	MarkSystemRanUnitTests()                                                                                           //
	FinalizeTransaction(transaction driver.Tx)                                                                         // finalizeTransaction ensures a transaction is closed after use, rolling back if not already committed.

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
	SaleCategoryRelation() SaleCategoryRelationStore                   //
	SaleProductRelation() SaleProductRelationStore                     //
	SaleCollectionRelation() SaleCollectionRelationStore               //
	GiftCard() GiftCardStore                                           // giftcard
	GiftcardEvent() GiftcardEventStore                                 //
	GiftCardOrder() GiftCardOrderStore                                 //
	GiftCardCheckout() GiftCardCheckoutStore                           //
	InvoiceEvent() InvoiceEventStore                                   // invoice
	Invoice() InvoiceStore                                             //
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
	ShippingMethodExcludedProduct() ShippingMethodExcludedProductStore //
	Warehouse() WarehouseStore                                         // warehouse
	Stock() StockStore                                                 //
	Allocation() AllocationStore                                       //
	WarehouseShippingZone() WarehouseShippingZoneStore                 //
	PreorderAllocation() PreorderAllocationStore
	Wishlist() WishlistStore                                           // wishlist
	WishlistItem() WishlistItemStore                                   //
	WishlistItemProductVariant() WishlistItemProductVariantStore       //
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
	OpenExchangeRate() OpenExchangeRateStore                           // external services tables
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
		Save(attr *attribute.Attribute) (*attribute.Attribute, error)                           // Save insert given attribute into database then returns it with an error. Returned can be wither *AppError or *NewErrInvalidInput or system error
		Get(id string) (*attribute.Attribute, error)                                            // Get try finding an attribute with given id then returns it with an error. Returned error can be either *store.ErrNotFound or system error
		GetBySlug(slug string) (*attribute.Attribute, error)                                    // GetBySlug finds an attribute with given slug, then returns it with an error. Returned error can be wither *ErrNotFound or system error
		FilterbyOption(option *attribute.AttributeFilterOption) ([]*attribute.Attribute, error) // FilterbyOption returns a list of attributes by given option
	}
	AttributeTranslationStore interface {
		CreateIndexesIfNotExists()
	}
	AttributeValueStore interface {
		CreateIndexesIfNotExists()
		ModelFields() []string
		Save(attribute *attribute.AttributeValue) (*attribute.AttributeValue, error) // Save inserts given attribute value into database, then returns inserted value and an error
		Get(attributeID string) (*attribute.AttributeValue, error)                   // Get finds an attribute value with given id then returns it with an error
		GetAllByAttributeID(attributeID string) ([]*attribute.AttributeValue, error) // GetAllByAttributeID finds all attribute values that belong to given attribute then returns them withh an error
	}
	AttributeValueTranslationStore interface {
		CreateIndexesIfNotExists()
	}
	AssignedPageAttributeValueStore interface {
		CreateIndexesIfNotExists()
		ModelFields() []string
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
		Save(assignedVariantAttribute *attribute.AssignedVariantAttribute) (*attribute.AssignedVariantAttribute, error)       // Save insert new instance into database then returns it with an error
		Get(id string) (*attribute.AssignedVariantAttribute, error)                                                           // Get find assigned variant attribute from database then returns it with an error
		GetWithOption(option *attribute.AssignedVariantAttributeFilterOption) (*attribute.AssignedVariantAttribute, error)    // GetWithOption try finding an assigned variant attribute with given option. If nothing found, it creates instance with given option. Finally it returns expected value with an error
		FilterByOption(option *attribute.AssignedVariantAttributeFilterOption) ([]*attribute.AssignedVariantAttribute, error) // FilterByOption finds and returns a list of assigned variant attributes filtered by given options
	}
	AttributeVariantStore interface {
		CreateIndexesIfNotExists()
		Save(attributeVariant *attribute.AttributeVariant) (*attribute.AttributeVariant, error)
		Get(attributeVariantID string) (*attribute.AttributeVariant, error)
		GetByOption(option *attribute.AttributeVariantFilterOption) (*attribute.AttributeVariant, error) // GetByOption finds 1 attribute variant with given option.
	}
	AssignedProductAttributeValueStore interface {
		CreateIndexesIfNotExists()
		ModelFields() []string
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
		Save(attributeProduct *attribute.AttributeProduct) (*attribute.AttributeProduct, error)          // Save inserts given attribute product relationship into database then returns it and an error
		Get(attributeProductID string) (*attribute.AttributeProduct, error)                              // Get finds an attributeProduct relationship and returns it with an error
		GetByOption(option *attribute.AttributeProductFilterOption) (*attribute.AttributeProduct, error) // GetByOption returns an attributeProduct with given condition
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
		GetById(id string) (*wishlist.Wishlist, error)                                 // GetById returns a wishlist with given id
		Upsert(wishList *wishlist.Wishlist) (*wishlist.Wishlist, error)                // Upsert inserts or update given wishlist and returns it
		GetByOption(option *wishlist.WishlistFilterOption) (*wishlist.Wishlist, error) // GetByOption finds and returns a slice of wishlists by given option
	}
	WishlistItemStore interface {
		CreateIndexesIfNotExists()
		GetById(selector *gorp.Transaction, id string) (*wishlist.WishlistItem, error)                                  // GetById returns a wishlist item wish given id
		BulkUpsert(transaction *gorp.Transaction, wishlistItems wishlist.WishlistItems) (wishlist.WishlistItems, error) // Upsert inserts or updates given wishlist item then returns it
		FilterByOption(option *wishlist.WishlistItemFilterOption) ([]*wishlist.WishlistItem, error)                     // FilterByOption finds and returns a slice of wishlist items filtered using given options
		GetByOption(option *wishlist.WishlistItemFilterOption) (*wishlist.WishlistItem, error)                          // GetByOption finds and returns a wishlist item filtered by given option
		DeleteItemsByOption(transaction *gorp.Transaction, option *wishlist.WishlistItemFilterOption) (int64, error)    // DeleteItemsByOption finds and deletes wishlist items that satisfy given filtering options and returns number of items deleted
	}
	WishlistItemProductVariantStore interface {
		CreateIndexesIfNotExists()
		Save(wishlistVariant *wishlist.WishlistItemProductVariant) (*wishlist.WishlistItemProductVariant, error)                                    // Save inserts new wishlist product variant relation into database and returns it
		BulkUpsert(transaction *gorp.Transaction, relations []*wishlist.WishlistItemProductVariant) ([]*wishlist.WishlistItemProductVariant, error) // BulkUpsert does bulk update/insert given relations
		GetById(selector *gorp.Transaction, id string) (*wishlist.WishlistItemProductVariant, error)                                                // GetByID returns a wishlist item product variant with given id
		DeleteRelation(relation *wishlist.WishlistItemProductVariant) (int64, error)                                                                // DeleteRelation deletes a product variant-wishlist item relation and counts numeber of relations left in database
	}
)

// warehouse
type (
	WarehouseStore interface {
		CreateIndexesIfNotExists()
		ModelFields() []string
		ScanFields(wh warehouse.WareHouse) []interface{}
		Save(warehouse *warehouse.WareHouse) (*warehouse.WareHouse, error)                      // Save inserts given warehouse into database then returns it.
		Get(id string) (*warehouse.WareHouse, error)                                            // Get try findings warehouse with given id, returns it. returned error could be wither (nil, *ErrNotFound, error)
		FilterByOprion(option *warehouse.WarehouseFilterOption) ([]*warehouse.WareHouse, error) // FilterByOprion returns a slice of warehouses with given option
		GetByOption(option *warehouse.WarehouseFilterOption) (*warehouse.WareHouse, error)      // GetByOption finds and returns a warehouse filtered given option
		WarehouseByStockID(stockID string) (*warehouse.WareHouse, error)                        // WarehouseByStockID returns 1 warehouse by given stock id
	}
	StockStore interface {
		CreateIndexesIfNotExists()
		ScanFields(stock warehouse.Stock) []interface{}
		ModelFields() []string
		Get(stockID string) (*warehouse.Stock, error)                                                                                                          // Get finds and returns stock with given stockID. Returned error could be either (nil, *ErrNotFound, error)
		FilterForCountryAndChannel(transaction *gorp.Transaction, options *warehouse.StockFilterForCountryAndChannel) ([]*warehouse.Stock, error)              // FilterForCountryAndChannel finds and returns stocks with given options
		FilterVariantStocksForCountry(transaction *gorp.Transaction, options *warehouse.StockFilterForCountryAndChannel) ([]*warehouse.Stock, error)           // FilterVariantStocksForCountry finds and returns stocks with given options
		FilterProductStocksForCountryAndChannel(transaction *gorp.Transaction, options *warehouse.StockFilterForCountryAndChannel) ([]*warehouse.Stock, error) // FilterProductStocksForCountryAndChannel finds and returns stocks with given options
		ChangeQuantity(stockID string, quantity int) error                                                                                                     // ChangeQuantity reduce or increase the quantity of given stock
		FilterByOption(transaction *gorp.Transaction, options *warehouse.StockFilterOption) ([]*warehouse.Stock, error)                                        // FilterByOption finds and returns a slice of stocks that satisfy given option
		BulkUpsert(transaction *gorp.Transaction, stocks []*warehouse.Stock) ([]*warehouse.Stock, error)                                                       // BulkUpsert performs upserts or inserts given stocks, then returns them
		FilterForChannel(options *warehouse.StockFilterForChannelOption) ([]*warehouse.Stock, error)                                                           // FilterForChannel finds and returns stocks that satisfy given options
	}
	AllocationStore interface {
		CreateIndexesIfNotExists()
		BulkUpsert(transaction *gorp.Transaction, allocations []*warehouse.Allocation) ([]*warehouse.Allocation, error)          // BulkUpsert performs update, insert given allocations then returns them afterward
		Get(allocationID string) (*warehouse.Allocation, error)                                                                  // Get find and returns allocation with given id
		FilterByOption(transaction *gorp.Transaction, option *warehouse.AllocationFilterOption) ([]*warehouse.Allocation, error) // FilterbyOption finds and returns a list of allocations based on given option
		BulkDelete(transaction *gorp.Transaction, allocationIDs []string) error                                                  // BulkDelete perform bulk deletes given allocations.
		CountAvailableQuantityForStock(stock *warehouse.Stock) (int, error)                                                      // CountAvailableQuantityForStock counts and returns available quantity of given stock
	}
	WarehouseShippingZoneStore interface {
		CreateIndexesIfNotExists()
		ModelFields() []string
		Save(warehouseShippingZone *warehouse.WarehouseShippingZone) (*warehouse.WarehouseShippingZone, error) // Save inserts given warehouse-shipping zone relation into database
	}
	PreorderAllocationStore interface {
		CreateIndexesIfNotExists()
		ModelFields() []string
		ScanFields(preorderAllocation warehouse.PreorderAllocation) []interface{}
		FilterByOption(options *warehouse.PreorderAllocationFilterOption) ([]*warehouse.PreorderAllocation, error) // FilterByOption finds and returns a list of preorder allocations filtered using given options
	}
)

// shipping
type (
	ShippingZoneStore interface {
		CreateIndexesIfNotExists()
		ModelFields() []string
		ScanFields(shippingZone shipping.ShippingZone) []interface{}
		Upsert(shippingZone *shipping.ShippingZone) (*shipping.ShippingZone, error)                 // Upsert depends on given shipping zone's Id to decide update or insert the zone
		Get(shippingZoneID string) (*shipping.ShippingZone, error)                                  // Get finds 1 shipping zone for given shippingZoneID
		FilterByOption(option *shipping.ShippingZoneFilterOption) ([]*shipping.ShippingZone, error) // FilterByOption finds a list of shipping zones based on given option
	}
	ShippingMethodStore interface {
		CreateIndexesIfNotExists()
		ModelFields() []string
		Upsert(method *shipping.ShippingMethod) (*shipping.ShippingMethod, error)                                                                                                   // Upsert bases on given method's Id to decide update or insert it
		Get(methodID string) (*shipping.ShippingMethod, error)                                                                                                                      // Get finds and returns a shipping method with given id
		ApplicableShippingMethods(price *goprices.Money, channelID string, weight *measurement.Weight, countryCode string, productIDs []string) ([]*shipping.ShippingMethod, error) // ApplicableShippingMethods finds all shipping methods with given conditions
		GetbyOption(options *shipping.ShippingMethodFilterOption) (*shipping.ShippingMethod, error)                                                                                 // GetbyOption finds and returns a shipping method that satisfy given options
	}
	ShippingMethodPostalCodeRuleStore interface {
		CreateIndexesIfNotExists()
		ModelFields() []string
	}
	ShippingMethodChannelListingStore interface {
		CreateIndexesIfNotExists()
		Upsert(listing *shipping.ShippingMethodChannelListing) (*shipping.ShippingMethodChannelListing, error)                      // Upsert depends on given listing's Id to decide whether to save or update the listing
		Get(listingID string) (*shipping.ShippingMethodChannelListing, error)                                                       // Get finds a shipping method channel listing with given listingID
		FilterByOption(option *shipping.ShippingMethodChannelListingFilterOption) ([]*shipping.ShippingMethodChannelListing, error) // FilterByOption returns a list of shipping method channel listings based on given option. result sorted by creation time ASC
	}
	ShippingMethodTranslationStore interface {
		CreateIndexesIfNotExists()
	}
	ShippingZoneChannelStore interface {
		CreateIndexesIfNotExists()
	}
	ShippingMethodExcludedProductStore interface {
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
		ModelFields() []string
		Upsert(collection *product_and_discount.Collection) (*product_and_discount.Collection, error)                   // Upsert depends on given collection's Id property to decide update or insert the collection
		Get(collectionID string) (*product_and_discount.Collection, error)                                              // Get finds and returns collection with given collectionID
		FilterByOption(option *product_and_discount.CollectionFilterOption) ([]*product_and_discount.Collection, error) // FilterByOption finds and returns a list of collections satisfy the given option
	}
	CollectionProductStore interface {
		CreateIndexesIfNotExists()
	}
	VariantMediaStore interface {
		CreateIndexesIfNotExists()
	}
	ProductMediaStore interface {
		CreateIndexesIfNotExists()
		Upsert(media *product_and_discount.ProductMedia) (*product_and_discount.ProductMedia, error)                        // Upsert depends on given media's Id property to decide insert or update it
		Get(id string) (*product_and_discount.ProductMedia, error)                                                          // Get finds and returns 1 product media with given id
		FilterByOption(option *product_and_discount.ProductMediaFilterOption) ([]*product_and_discount.ProductMedia, error) // FilterByOption finds and returns a list of product medias with given id
	}
	DigitalContentUrlStore interface {
		CreateIndexesIfNotExists()
		Upsert(contentURL *product_and_discount.DigitalContentUrl) (*product_and_discount.DigitalContentUrl, error) // Upsert inserts or updates given digital content url into database then returns it
		Get(id string) (*product_and_discount.DigitalContentUrl, error)                                             // Get finds and returns a digital content url with given id
	}
	DigitalContentStore interface {
		CreateIndexesIfNotExists()
		ModelFields() []string
		ScanFields(content product_and_discount.DigitalContent) []interface{}
		Save(content *product_and_discount.DigitalContent) (*product_and_discount.DigitalContent, error)                    // Save inserts given digital content into database then returns it
		GetByOption(option *product_and_discount.DigitalContenetFilterOption) (*product_and_discount.DigitalContent, error) // GetByOption finds and returns 1 digital content filtered using given option
	}
	ProductVariantChannelListingStore interface {
		CreateIndexesIfNotExists()
		ModelFields() []string
		ScanFields(listing product_and_discount.ProductVariantChannelListing) []interface{}
		Save(variantChannelListing *product_and_discount.ProductVariantChannelListing) (*product_and_discount.ProductVariantChannelListing, error)                                         // Save insert given value into database then returns it with an error
		Get(variantChannelListingID string) (*product_and_discount.ProductVariantChannelListing, error)                                                                                    // Get finds and returns 1 product variant channel listing based on given variantChannelListingID
		FilterbyOption(transaction *gorp.Transaction, option *product_and_discount.ProductVariantChannelListingFilterOption) ([]*product_and_discount.ProductVariantChannelListing, error) // FilterbyOption finds and returns all product variant channel listings filterd using given option
	}
	ProductVariantTranslationStore interface {
		CreateIndexesIfNotExists()
		Upsert(translation *product_and_discount.ProductVariantTranslation) (*product_and_discount.ProductVariantTranslation, error)                  // Upsert inserts or updates given translation then returns it
		Get(translationID string) (*product_and_discount.ProductVariantTranslation, error)                                                            // Get finds and returns 1 product variant translation with given id
		FilterByOption(option *product_and_discount.ProductVariantTranslationFilterOption) ([]*product_and_discount.ProductVariantTranslation, error) // FilterByOption finds and returns product variant translations filtered using given options
	}
	ProductVariantStore interface {
		CreateIndexesIfNotExists()
		ModelFields() []string
		ScanFields(variant product_and_discount.ProductVariant) []interface{}
		Save(variant *product_and_discount.ProductVariant) (*product_and_discount.ProductVariant, error)                        // Save inserts product variant instance to database
		Get(id string) (*product_and_discount.ProductVariant, error)                                                            // Get returns a product variant with given id
		GetWeight(productVariantID string) (*measurement.Weight, error)                                                         // GetWeight returns weight of given product variant
		GetByOrderLineID(orderLineID string) (*product_and_discount.ProductVariant, error)                                      // GetByOrderLineID finds and returns a product variant by given orderLineID
		FilterByOption(option *product_and_discount.ProductVariantFilterOption) ([]*product_and_discount.ProductVariant, error) // FilterByOption finds and returns product variants based on given option
	}
	ProductChannelListingStore interface {
		CreateIndexesIfNotExists()
		ModelFields() []string
		BulkUpsert(listings []*product_and_discount.ProductChannelListing) ([]*product_and_discount.ProductChannelListing, error)             // BulkUpsert performs bulk upsert on given product channel listings
		Get(channelListingID string) (*product_and_discount.ProductChannelListing, error)                                                     // Get try finding a product channel listing, then returns it with an error
		FilterByOption(option *product_and_discount.ProductChannelListingFilterOption) ([]*product_and_discount.ProductChannelListing, error) // FilterByOption filter a list of product channel listings by given option. Then returns them with an error
	}
	ProductTranslationStore interface {
		CreateIndexesIfNotExists()
		Upsert(translation *product_and_discount.ProductTranslation) (*product_and_discount.ProductTranslation, error)                  // Upsert inserts or update given translation
		Get(translationID string) (*product_and_discount.ProductTranslation, error)                                                     // Get finds and returns a product translation by given id
		FilterByOption(option *product_and_discount.ProductTranslationFilterOption) ([]*product_and_discount.ProductTranslation, error) // FilterByOption finds and returns product translations filtered using given options
	}
	ProductTypeStore interface {
		CreateIndexesIfNotExists()
		ModelFields() []string
		Save(productType *product_and_discount.ProductType) (*product_and_discount.ProductType, error)                // Save try inserting new product type into database then returns it
		FilterProductTypesByCheckoutID(checkoutToken string) ([]*product_and_discount.ProductType, error)             // FilterProductTypesByCheckoutID is used to check if a checkout requires shipping
		ProductTypesByProductIDs(productIDs []string) ([]*product_and_discount.ProductType, error)                    // ProductTypesByProductIDs returns all product types belong to given products
		ProductTypeByProductVariantID(variantID string) (*product_and_discount.ProductType, error)                    // ProductTypeByProductVariantID finds and returns 1 product type that is related to given product variant
		GetByOption(options *product_and_discount.ProductTypeFilterOption) (*product_and_discount.ProductType, error) // GetByOption finds and returns a product type with given options
	}
	CategoryTranslationStore interface {
		CreateIndexesIfNotExists()
	}
	CategoryStore interface {
		CreateIndexesIfNotExists()
		Upsert(category *product_and_discount.Category) (*product_and_discount.Category, error)                     // Upsert depends on given category's Id field to decide update or insert it
		Get(categoryID string) (*product_and_discount.Category, error)                                              // Get finds and returns a category with given id
		GetByOption(option *product_and_discount.CategoryFilterOption) (*product_and_discount.Category, error)      // GetByOption finds and returns 1 category satisfy given option
		FilterByOption(option *product_and_discount.CategoryFilterOption) ([]*product_and_discount.Category, error) // FilterByOption finds and returns a list of categories satisfy given option
	}
	ProductStore interface {
		CreateIndexesIfNotExists()
		ModelFields() []string
		ScanFields(prd product_and_discount.Product) []interface{}
		Save(prd *product_and_discount.Product) (*product_and_discount.Product, error)
		GetByOption(option *product_and_discount.ProductFilterOption) (*product_and_discount.Product, error)      // GetByOption finds and returns 1 product that satisfies given option
		FilterByOption(option *product_and_discount.ProductFilterOption) ([]*product_and_discount.Product, error) // FilterByOption finds and returns all products that satisfy given option
		PublishedProducts(channelSlug string) ([]*product_and_discount.Product, error)                            // FilterPublishedProducts finds and returns products that belong to given channel slug and are published
		NotPublishedProducts(channelSlug string) ([]*struct {
			product_and_discount.Product
			IsPublished     bool
			PublicationDate *timemodule.Time
		}, error) // FilterNotPublishedProducts finds all not published products belong to given channel
		PublishedWithVariants(channelSlug string) ([]*product_and_discount.Product, error)                                                                      // PublishedWithVariants finds and returns products.
		VisibleToUserProducts(channelSlug string, requesterIsStaff bool) ([]*product_and_discount.Product, error)                                               // FilterVisibleToUserProduct finds and returns all products that are visible to requesting user.
		SelectForUpdateDiscountedPricesOfCatalogues(productIDs []string, categoryIDs []string, collectionIDs []string) ([]*product_and_discount.Product, error) // SelectForUpdateDiscountedPricesOfCatalogues finds and returns product based on given ids lists.
	}
)

// payment
type (
	PaymentStore interface {
		CreateIndexesIfNotExists()
		ScanFields(payMent payment.Payment) []interface{}
		Save(transaction *gorp.Transaction, payment *payment.Payment) (*payment.Payment, error)                           // Save save payment instance into database
		Get(transaction *gorp.Transaction, id string, lockForUpdate bool) (*payment.Payment, error)                       // Get returns a payment with given id. `lockForUpdate` is true if you want to add "FOR UPDATE" to sql
		Update(transaction *gorp.Transaction, payment *payment.Payment) (*payment.Payment, error)                         // Update updates given payment and returns new updated payment
		CancelActivePaymentsOfCheckout(checkoutToken string) error                                                        // CancelActivePaymentsOfCheckout inactivate all payments that belong to given checkout and in active status
		FilterByOption(option *payment.PaymentFilterOption) ([]*payment.Payment, error)                                   // FilterByOption finds and returns a list of payments that satisfy given option
		UpdatePaymentsOfCheckout(transaction *gorp.Transaction, checkoutToken string, option *payment.PaymentPatch) error // UpdatePaymentsOfCheckout updates payments of given checkout
	}
	PaymentTransactionStore interface {
		CreateIndexesIfNotExists()
		Save(transaction *gorp.Transaction, paymentTransaction *payment.PaymentTransaction) (*payment.PaymentTransaction, error) // Save inserts new payment transaction into database
		Get(id string) (*payment.PaymentTransaction, error)                                                                      // Get returns a payment transaction with given id
		Update(transaction *payment.PaymentTransaction) (*payment.PaymentTransaction, error)                                     // Update updates given transaction and returns updated one
		FilterByOption(option *payment.PaymentTransactionFilterOpts) ([]*payment.PaymentTransaction, error)                      // FilterByOption finds and returns a list of transactions with given option
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
		ScanFields(orderLine order.OrderLine) []interface{}
		ModelFields() []string
		Upsert(transaction *gorp.Transaction, orderLine *order.OrderLine) (*order.OrderLine, error)          // Upsert depends on given orderLine's Id to decide to update or save it
		Get(id string) (*order.OrderLine, error)                                                             // Get returns a order line with id of given id
		BulkDelete(orderLineIDs []string) error                                                              // BulkDelete delete all given order lines. NOTE: validate given ids are valid uuids before calling me
		FilterbyOption(option *order.OrderLineFilterOption) ([]*order.OrderLine, error)                      // FilterbyOption finds and returns order lines by given option
		BulkUpsert(transaction *gorp.Transaction, orderLines []*order.OrderLine) ([]*order.OrderLine, error) // BulkUpsert performs upsert multiple order lines in once
	}
	OrderStore interface {
		CreateIndexesIfNotExists()
		ModelFields() []string
		ScanFields(holder order.Order) []interface{}
		Save(transaction *gorp.Transaction, order *order.Order) (*order.Order, error)   // Save insert an order into database and returns that order if success
		Get(id string) (*order.Order, error)                                            // Get find order in database with given id
		Update(transaction *gorp.Transaction, order *order.Order) (*order.Order, error) // Update update order
		FilterByOption(option *order.OrderFilterOption) ([]*order.Order, error)         // FilterByOption returns a list of orders, filtered by given option
		BulkUpsert(orders []*order.Order) ([]*order.Order, error)                       // BulkUpsert performs bulk upsert given orders
	}
	OrderEventStore interface {
		CreateIndexesIfNotExists()
		Save(transaction *gorp.Transaction, orderEvent *order.OrderEvent) (*order.OrderEvent, error) // Save inserts given order event into database then returns it
		Get(orderEventID string) (*order.OrderEvent, error)                                          // Get finds order event with given id then returns it
	}
	FulfillmentLineStore interface {
		CreateIndexesIfNotExists()
		Save(fulfillmentLine *order.FulfillmentLine) (*order.FulfillmentLine, error)
		Get(id string) (*order.FulfillmentLine, error)
		FilterbyOption(option *order.FulfillmentLineFilterOption) ([]*order.FulfillmentLine, error)                            // FilterbyOption finds and returns a list of fulfillment lines by given option
		BulkUpsert(transaction *gorp.Transaction, fulfillmentLines []*order.FulfillmentLine) ([]*order.FulfillmentLine, error) // BulkUpsert upsert given fulfillment lines
		DeleteFulfillmentLinesByOption(transaction *gorp.Transaction, option *order.FulfillmentLineFilterOption) error         // DeleteFulfillmentLinesByOption filters fulfillment lines by given option, then deletes them
	}
	FulfillmentStore interface {
		CreateIndexesIfNotExists()
		ModelFields() []string
		ScanFields(holder order.Fulfillment) []interface{}
		Upsert(transaction *gorp.Transaction, fulfillment *order.Fulfillment) (*order.Fulfillment, error)                  // Upsert depends on given fulfillment's Id to decide update or insert it
		Get(id string) (*order.Fulfillment, error)                                                                         // Get finds and return a fulfillment by given id
		GetByOption(transaction *gorp.Transaction, option *order.FulfillmentFilterOption) (*order.Fulfillment, error)      // GetByOption returns 1 fulfillment, filtered by given option
		FilterByOption(transaction *gorp.Transaction, option *order.FulfillmentFilterOption) ([]*order.Fulfillment, error) // FilterByOption finds and returns a slice of fulfillments by given option
		DeleteByOptions(transaction *gorp.Transaction, options *order.FulfillmentFilterOption) error                       // DeleteByOptions deletes fulfillment database records that satisfy given option. It returns an error indicates if there is a problem occured during deletion process
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

// invoice
type (
	InvoiceEventStore interface {
		CreateIndexesIfNotExists()
		Upsert(invoiceEvent *invoice.InvoiceEvent) (*invoice.InvoiceEvent, error) // Upsert depends on given invoice event's Id to update/insert it
		Get(invoiceEventID string) (*invoice.InvoiceEvent, error)                 // Get finds and returns 1 invoice event
	}
	InvoiceStore interface {
		CreateIndexesIfNotExists()
		Upsert(invoice *invoice.Invoice) (*invoice.Invoice, error) // Upsert depends on given invoice Id to update/insert it
		Get(invoiceID string) (*invoice.Invoice, error)            // Get finds and returns 1 invoice
	}
)

// giftcard related stores
type (
	GiftCardStore interface {
		CreateIndexesIfNotExists()
		BulkUpsert(transaction *gorp.Transaction, giftCards ...*giftcard.GiftCard) ([]*giftcard.GiftCard, error)           // BulkUpsert depends on given giftcards's Id properties then perform according operation
		GetById(id string) (*giftcard.GiftCard, error)                                                                     // GetById returns a giftcard instance that has id of given id
		FilterByOption(transaction *gorp.Transaction, option *giftcard.GiftCardFilterOption) ([]*giftcard.GiftCard, error) // FilterByOption finds giftcards wth option
	}
	GiftcardEventStore interface {
		CreateIndexesIfNotExists()
		Save(event *giftcard.GiftCardEvent) (*giftcard.GiftCardEvent, error)                                            // Save insdert given giftcard event into database then returns it
		Get(id string) (*giftcard.GiftCardEvent, error)                                                                 // Get finds and returns a giftcard event found by given id
		BulkUpsert(transaction *gorp.Transaction, events ...*giftcard.GiftCardEvent) ([]*giftcard.GiftCardEvent, error) // BulkUpsert upserts and returns given giftcard events
	}
	GiftCardOrderStore interface {
		CreateIndexesIfNotExists()
		Save(giftcardOrder *giftcard.OrderGiftCard) (*giftcard.OrderGiftCard, error)                                            // Save inserts new giftcard-order relation into database then returns it
		Get(id string) (*giftcard.OrderGiftCard, error)                                                                         // Get returns giftcard-order relation table with given id
		BulkUpsert(transaction *gorp.Transaction, orderGiftcards ...*giftcard.OrderGiftCard) ([]*giftcard.OrderGiftCard, error) // BulkUpsert upserts given order-giftcard relations and returns it
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
		Upsert(transaction *gorp.Transaction, orderDiscount *product_and_discount.OrderDiscount) (*product_and_discount.OrderDiscount, error) // Upsert depends on given order discount's Id property to decide to update/insert it
		Get(orderDiscountID string) (*product_and_discount.OrderDiscount, error)                                                              // Get finds and returns an order discount with given id
		FilterbyOption(option *product_and_discount.OrderDiscountFilterOption) ([]*product_and_discount.OrderDiscount, error)                 // FilterbyOption filters order discounts that satisfy given option, then returns them
		BulkDelete(orderDiscountIDs []string) error                                                                                           // BulkDelete perform bulk delete all given order discount ids
	}
	DiscountSaleTranslationStore interface {
		CreateIndexesIfNotExists()
	}
	DiscountSaleChannelListingStore interface {
		CreateIndexesIfNotExists()
		ModelFields() []string
		Save(saleChannelListing *product_and_discount.SaleChannelListing) (*product_and_discount.SaleChannelListing, error) // Save insert given instance into database then returns it
		Get(saleChannelListingID string) (*product_and_discount.SaleChannelListing, error)                                  // Get finds and returns sale channel listing with given id
		// SaleChannelListingsWithOption finds a list of sale channel listings plus foreign channel slugs
		SaleChannelListingsWithOption(option *product_and_discount.SaleChannelListingFilterOption) (
			[]*struct {
				product_and_discount.SaleChannelListing
				ChannelSlug string
			},
			error,
		)
	}
	VoucherTranslationStore interface {
		CreateIndexesIfNotExists()
		Save(translation *product_and_discount.VoucherTranslation) (*product_and_discount.VoucherTranslation, error)                    // Save inserts given translation into database and returns it
		Get(id string) (*product_and_discount.VoucherTranslation, error)                                                                // Get finds and returns a voucher translation with given id
		FilterByOption(option *product_and_discount.VoucherTranslationFilterOption) ([]*product_and_discount.VoucherTranslation, error) // FilterByOption returns a list of voucher translations filtered using given options
		GetByOption(option *product_and_discount.VoucherTranslationFilterOption) (*product_and_discount.VoucherTranslation, error)      // GetByOption finds and returns 1 voucher translation by given options
	}
	DiscountSaleStore interface {
		CreateIndexesIfNotExists()
		Upsert(sale *product_and_discount.Sale) (*product_and_discount.Sale, error)                              // Upsert bases on sale's Id to decide to update or insert given sale
		Get(saleID string) (*product_and_discount.Sale, error)                                                   // Get finds and returns a sale with given saleID
		FilterSalesByOption(option *product_and_discount.SaleFilterOption) ([]*product_and_discount.Sale, error) // FilterSalesByOption filter sales by option
	}
	VoucherChannelListingStore interface {
		CreateIndexesIfNotExists()
		Upsert(voucherChannelListing *product_and_discount.VoucherChannelListing) (*product_and_discount.VoucherChannelListing, error)        // upsert check given listing's Id to decide whether to create or update it. Then returns a listing with an error
		Get(voucherChannelListingID string) (*product_and_discount.VoucherChannelListing, error)                                              // Get finds a listing with given id, then returns it with an error
		FilterbyOption(option *product_and_discount.VoucherChannelListingFilterOption) ([]*product_and_discount.VoucherChannelListing, error) // FilterbyOption finds and returns a list of voucher channel listing relationship instances filtered by given option
	}
	DiscountVoucherStore interface {
		CreateIndexesIfNotExists()
		ModelFields() []string
		ScanFields(voucher product_and_discount.Voucher) []interface{}
		Upsert(voucher *product_and_discount.Voucher) (*product_and_discount.Voucher, error)                              // Upsert saves or updates given voucher then returns it with an error
		Get(voucherID string) (*product_and_discount.Voucher, error)                                                      // Get finds a voucher with given id, then returns it with an error
		FilterVouchersByOption(option *product_and_discount.VoucherFilterOption) ([]*product_and_discount.Voucher, error) // FilterVouchersByOption finds vouchers bases on given option.
		ExpiredVouchers(date *timemodule.Time) ([]*product_and_discount.Voucher, error)                                   // ExpiredVouchers finds and returns vouchers that are expired before given date
		GetByOptions(options *product_and_discount.VoucherFilterOption) (*product_and_discount.Voucher, error)            // GetByOptions finds and returns 1 voucher filtered using given options
	}
	VoucherCategoryStore interface {
		CreateIndexesIfNotExists()
		Upsert(voucherCategory *product_and_discount.VoucherCategory) (*product_and_discount.VoucherCategory, error) // Upsert saves or updates given voucher category then returns it with an error
		Get(voucherCategoryID string) (*product_and_discount.VoucherCategory, error)                                 // Get finds a voucher category with given id, then returns it with an error
	}
	VoucherCollectionStore interface {
		CreateIndexesIfNotExists()
		Upsert(voucherCollection *product_and_discount.VoucherCollection) (*product_and_discount.VoucherCollection, error) // Upsert saves or updates given voucher collection then returns it with an error
		Get(voucherCollectionID string) (*product_and_discount.VoucherCollection, error)                                   // Get finds a voucher collection with given id, then returns it with an error
	}
	VoucherProductStore interface {
		CreateIndexesIfNotExists()
		Upsert(voucherProduct *product_and_discount.VoucherProduct) (*product_and_discount.VoucherProduct, error) // Upsert saves or updates given voucher product then returns it with an error
		Get(voucherProductID string) (*product_and_discount.VoucherProduct, error)                                // Get finds a voucher product with given id, then returns it with an error
	}
	VoucherCustomerStore interface {
		CreateIndexesIfNotExists()
		Save(voucherCustomer *product_and_discount.VoucherCustomer) (*product_and_discount.VoucherCustomer, error)                  // Save inserts given voucher customer instance into database ands returns it
		DeleteInBulk(relations []*product_and_discount.VoucherCustomer) error                                                       // DeleteInBulk deletes given voucher-customers with given id
		GetByOption(options *product_and_discount.VoucherCustomerFilterOption) (*product_and_discount.VoucherCustomer, error)       // GetByOption finds and returns a voucher customer with given options
		FilterByOptions(options *product_and_discount.VoucherCustomerFilterOption) ([]*product_and_discount.VoucherCustomer, error) // FilterByOptions finds and returns a slice of voucher customers by given options
	}
	SaleCategoryRelationStore interface {
		CreateIndexesIfNotExists()
		Save(relation *product_and_discount.SaleCategoryRelation) (*product_and_discount.SaleCategoryRelation, error)                               // Save inserts given sale-category relation into database
		Get(relationID string) (*product_and_discount.SaleCategoryRelation, error)                                                                  // Get returns 1 sale-category relation with given id
		SaleCategoriesByOption(option *product_and_discount.SaleCategoryRelationFilterOption) ([]*product_and_discount.SaleCategoryRelation, error) // SaleCategoriesByOption returns a slice of sale-category relations with given option
	}
	SaleProductRelationStore interface {
		CreateIndexesIfNotExists()
		Save(relation *product_and_discount.SaleProductRelation) (*product_and_discount.SaleProductRelation, error)                             // Save inserts given sale-product relation into database then returns it
		Get(relationID string) (*product_and_discount.SaleProductRelation, error)                                                               // Get finds and returns a sale-product relation with given id
		SaleProductsByOption(option *product_and_discount.SaleProductRelationFilterOption) ([]*product_and_discount.SaleProductRelation, error) // SaleProductsByOption returns a slice of sale-product relations, filtered by given option
	}
	SaleCollectionRelationStore interface {
		CreateIndexesIfNotExists()
		Save(relation *product_and_discount.SaleCollectionRelation) (*product_and_discount.SaleCollectionRelation, error)                       // Save insert given sale-collection relation into database
		Get(relationID string) (*product_and_discount.SaleCollectionRelation, error)                                                            // Get finds and returns a sale-collection relation with given id
		FilterByOption(option *product_and_discount.SaleCollectionRelationFilterOption) ([]*product_and_discount.SaleCollectionRelation, error) // FilterByOption returns a list of collections filtered based on given option
	}
)

// csv
type (
	CsvExportEventStore interface {
		CreateIndexesIfNotExists()
		Save(event *csv.ExportEvent) (*csv.ExportEvent, error)                           // Save inserts given export event into database then returns it
		FilterByOption(options *csv.ExportEventFilterOption) ([]*csv.ExportEvent, error) // FilterByOption finds and returns a list of export events filtered using given option
	}
	CsvExportFileStore interface {
		CreateIndexesIfNotExists()
		Save(file *csv.ExportFile) (*csv.ExportFile, error) // Save inserts given export file into database then returns it
		Get(id string) (*csv.ExportFile, error)             // Get finds and returns an export file found using given id
	}
)

// checkout
type (
	CheckoutLineStore interface {
		CreateIndexesIfNotExists()
		ModelFields() []string
		ScanFields(line checkout.CheckoutLine) []interface{}
		Upsert(checkoutLine *checkout.CheckoutLine) (*checkout.CheckoutLine, error)          // Upsert checks whether to update or insert given checkout line then performs according operation
		Get(id string) (*checkout.CheckoutLine, error)                                       // Get returns a checkout line with given id
		CheckoutLinesByCheckoutID(checkoutID string) ([]*checkout.CheckoutLine, error)       // CheckoutLinesByCheckoutID returns a list of checkout lines that belong to given checkout
		DeleteLines(transaction *gorp.Transaction, checkoutLineIDs []string) error           // DeleteLines deletes all checkout lines with given uuids
		BulkUpdate(checkoutLines []*checkout.CheckoutLine) error                             // BulkUpdate receives a list of modified checkout lines, updates them in bulk.
		BulkCreate(checkoutLines []*checkout.CheckoutLine) ([]*checkout.CheckoutLine, error) // BulkCreate takes a list of raw checkout lines, save them into database then returns them fully with an error
		// CheckoutLinesByCheckoutWithPrefetch finds all checkout lines belong to given checkout
		//
		// and prefetch all related product variants, products
		//
		// this borrows the idea from Django's prefetch_related() method
		CheckoutLinesByCheckoutWithPrefetch(checkoutID string) ([]*checkout.CheckoutLine, []*product_and_discount.ProductVariant, []*product_and_discount.Product, error)
		TotalWeightForCheckoutLines(checkoutLineIDs []string) (*measurement.Weight, error)                 // TotalWeightForCheckoutLines calculate total weight for given checkout lines
		CheckoutLinesByOption(option *checkout.CheckoutLineFilterOption) ([]*checkout.CheckoutLine, error) // CheckoutLinesByOption finds and returns checkout lines filtered using given option
	}
	CheckoutStore interface {
		CreateIndexesIfNotExists()
		Get(token string) (*checkout.Checkout, error)                                                             // Get finds a checkout with given token (checkouts use tokens(uuids) as primary keys)
		Upsert(ckout *checkout.Checkout) (*checkout.Checkout, error)                                              // Upsert depends on given checkout's Token property to decide to update or insert it
		FetchCheckoutLinesAndPrefetchRelatedValue(ckout *checkout.Checkout) ([]*checkout.CheckoutLineInfo, error) // FetchCheckoutLinesAndPrefetchRelatedValue Fetch checkout lines as CheckoutLineInfo objects.
		GetByOption(option *checkout.CheckoutFilterOption) (*checkout.Checkout, error)                            // GetByOption finds and returns 1 checkout based on given option
		FilterByOption(option *checkout.CheckoutFilterOption) ([]*checkout.Checkout, error)                       // FilterByOption finds and returns a list of checkout based on given option
		DeleteCheckoutsByOption(transaction *gorp.Transaction, option *checkout.CheckoutFilterOption) error       // DeleteCheckoutsByOption deletes checkout row(s) from database, filtered using given option.  It returns an error indicating if the operation was performed successfully.
	}
)

// channel
type ChannelStore interface {
	CreateIndexesIfNotExists()
	ModelFields() []string
	ScanFields(chanNel channel.Channel) []interface{}
	Save(ch *channel.Channel) (*channel.Channel, error)
	Get(id string) (*channel.Channel, error)                                        // Get returns channel by given id
	GetRandomActiveChannel() (*channel.Channel, error)                              // GetRandomActiveChannel get an abitrary channel that is active
	FilterByOption(option *channel.ChannelFilterOption) ([]*channel.Channel, error) // FilterByOption returns a list of channels with given option
	GetbyOption(option *channel.ChannelFilterOption) (*channel.Channel, error)      // GetbyOption finds and returns 1 channel filtered using given options
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
		ModelFields() model.StringArray
		ScanFields(addr account.Address) []interface{}
		CreateIndexesIfNotExists()                                                                // CreateIndexesIfNotExists creates indexes for table if needed
		Save(transaction *gorp.Transaction, address *account.Address) (*account.Address, error)   // Save saves address into database
		Get(addressID string) (*account.Address, error)                                           // Get returns an Address with given addressID is exist
		Update(transaction *gorp.Transaction, address *account.Address) (*account.Address, error) // Update update given address and returns it
		DeleteAddresses(addressIDs []string) error                                                // DeleteAddress deletes given address and returns an error
		FilterByOption(option *account.AddressFilterOption) ([]*account.Address, error)           // FilterByOption finds and returns a list of address(es) filtered by given option
	}
	UserTermOfServiceStore interface {
		CreateIndexesIfNotExists()                                                                //
		GetByUser(userID string) (*account.UserTermsOfService, error)                             // GetByUser returns a term of service with given user id
		Save(userTermsOfService *account.UserTermsOfService) (*account.UserTermsOfService, error) // Save inserts new user term of service to database
		Delete(userID, termsOfServiceId string) error                                             // Delete deletes from database an usder term of service with given userId and term of service id
	}
	UserStore interface {
		ClearCaches()
		CreateIndexesIfNotExists()
		ModelFields() []string
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
		GetKnownUsers(userID string) ([]string, error)
		Count(options account.UserCountOptions) (int64, error)
		AnalyticsActiveCountForPeriod(startTime int64, endTime int64, options account.UserCountOptions) (int64, error)
		GetAllProfiles(options *account.UserGetOptions) ([]*account.User, error)
		Search(term string, options *account.UserSearchOptions) ([]*account.User, error)
		AnalyticsActiveCount(time int64, options account.UserCountOptions) (int64, error)
		GetProfileByIds(ctx context.Context, userIds []string, options *UserGetByIdsOpts, allowFromCache bool) ([]*account.User, error)
		GetProfilesByUsernames(usernames []string) ([]*account.User, error)
		GetProfiles(options *account.UserGetOptions) ([]*account.User, error)
		GetUnreadCount(userID string) (int64, error)         // TODO: consider me
		UserByOrderID(orderID string) (*account.User, error) // UserByOrderID finds and returns an user who whose order is given

		// PromoteGuestToUser(userID string) error
		// DemoteUserToGuest(userID string) (*account.User, error)
		// DeactivateGuests() ([]string, error)
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

type OpenExchangeRateStore interface {
	CreateIndexesIfNotExists()
	BulkUpsert(rates []*external_services.OpenExchangeRate) ([]*external_services.OpenExchangeRate, error) // BulkUpsert performs bulk update/insert to given exchange rates
	GetAll() ([]*external_services.OpenExchangeRate, error)                                                // GetAll returns all exchange currency rates
}

type UserGetByIdsOpts struct {
	IsAdmin bool  // IsAdmin tracks whether or not the request is being made by an administrator. Does nothing when provided by a client.
	Since   int64 // Since filters the users based on their UpdateAt timestamp.
	// Restrict to search in a list of teams and channels. Does nothing when provided by a client.
	// ViewRestrictions *model.ViewUsersRestrictions
}
