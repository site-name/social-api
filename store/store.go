//go:generate go run layer_generators/main.go

package store

import (
	"context"
	"database/sql"
	"database/sql/driver"
	timemodule "time"

	"github.com/Masterminds/squirrel"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store/store_iface"
)

// Store is database gateway of the system
type Store interface {
	Context() context.Context           // Context gets context
	SetContext(context context.Context) // set context
	Close()                             // Close closes databases
	LockToMaster()                      // LockToMaster constraints all queries to be performed on master
	UnlockFromMaster()                  // UnlockFromMaster makes all datasources available
	GetInternalReplicaDBs() []*sql.DB   // GetInternalReplicaDBs allows access to the raw replica DB handles for the multi-product architecture.
	ReplicaLagTime() error
	ReplicaLagAbs() error
	CheckIntegrity() <-chan model.IntegrityCheckResult
	DropAllTables()                              // DropAllTables drop all tables in databases
	GetDbVersion(numerical bool) (string, error) // GetDbVersion returns version in use of database
	GetMasterX() store_iface.SqlxExecutor        // GetMaster get master datasource
	GetReplicaX() store_iface.SqlxExecutor       // GetMaster gets slave datasource
	// GetQueryBuilder create squirrel sql query builder.
	//
	// NOTE: Don't pass much placeholder format since only the first passed is applied.
	// Ellipsis operator is a trick to support no argument passing.
	//
	// If no placeholder format is passed, defaut to squirrel.Dollar ($)
	GetQueryBuilder(placeholderFormats ...squirrel.PlaceholderFormat) squirrel.StatementBuilderType
	IsUniqueConstraintError(err error, indexName []string) bool //
	MarkSystemRanUnitTests()                                    //
	FinalizeTransaction(transaction driver.Tx)                  // FinalizeTransaction ensures a transaction is closed after use, rolling back if not already committed.
	DBXFromContext(ctx context.Context) store_iface.SqlxExecutor

	User() UserStore                                                   // account
	Address() AddressStore                                             //
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
	ChannelShop() ChannelShopStore                                     //
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
	VoucherProductVariant() VoucherProductVariantStore                 //
	SaleProductVariant() SaleProductVariantStore                       //
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
	PreorderAllocation() PreorderAllocationStore                       //
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
	Vat() VatStore                                                     //
	OpenExchangeRate() OpenExchangeRateStore                           // external services
}

// shop
type (
	ShopStaffStore interface {
		Save(shopStaff *model.ShopStaff) (*model.ShopStaff, error)                         // Save inserts given shopStaff into database then returns it with an error
		Get(shopStaffID string) (*model.ShopStaff, error)                                  // Get finds a shop staff with given id then returns it with an error
		FilterByOptions(options *model.ShopStaffFilterOptions) ([]*model.ShopStaff, error) // FilterByShopAndStaff finds a relation ship with given shopId and staffId
		GetByOptions(options *model.ShopStaffFilterOptions) (*model.ShopStaff, error)
	}
	ShopStore interface {
		ModelFields(prefix string) util.AnyArray[string]
		ScanFields(shop *model.Shop) []interface{}
		Upsert(shop *model.Shop) (*model.Shop, error)                            // Upsert depends on shop's Id to decide to update/insert the given model.
		Get(shopID string) (*model.Shop, error)                                  // Get finds a shop with given id and returns it
		FilterByOptions(options *model.ShopFilterOptions) ([]*model.Shop, error) // FilterByOptions finds and returns shops with given options
		GetByOptions(options *model.ShopFilterOptions) (*model.Shop, error)      // GetByOptions finds and returns 1 shop with given options
	}
	ShopTranslationStore interface {
		Upsert(translation *model.ShopTranslation) (*model.ShopTranslation, error) // Upsert depends on translation's Id then decides to update or insert
		Get(id string) (*model.ShopTranslation, error)                             // Get finds a shop translation with given id then return it with an error
	}
	VatStore interface {
		Upsert(transaction store_iface.SqlxTxExecutor, vats []*model.Vat) ([]*model.Vat, error)
		FilterByOptions(options *model.VatFilterOptions) ([]*model.Vat, error)
	}
)

// Plugin
type PluginStore interface {
	SaveOrUpdate(keyVal *model.PluginKeyValue) (*model.PluginKeyValue, error)
	CompareAndSet(keyVal *model.PluginKeyValue, oldValue []byte) (bool, error)
	CompareAndDelete(keyVal *model.PluginKeyValue, oldValue []byte) (bool, error)
	SetWithOptions(pluginID string, key string, value []byte, options model.PluginKVSetOptions) (bool, error)
	Get(pluginID, key string) (*model.PluginKeyValue, error)
	Delete(pluginID, key string) error
	DeleteAllForPlugin(PluginID string) error
	DeleteAllExpired() error
	List(pluginID string, page, perPage int) ([]string, error)
}

type UploadSessionStore interface {
	Save(session *model.UploadSession) (*model.UploadSession, error)
	Update(session *model.UploadSession) error
	Get(id string) (*model.UploadSession, error)
	GetForUser(userID string) ([]*model.UploadSession, error)
	Delete(id string) error
}

// fileinfo
type FileInfoStore interface {
	Upsert(info *model.FileInfo) (*model.FileInfo, error)
	Get(id string) (*model.FileInfo, error)
	GetFromMaster(id string) (*model.FileInfo, error)
	GetByIds(ids []string) ([]*model.FileInfo, error)
	GetByPath(path string) (*model.FileInfo, error)
	GetForUser(userID string) ([]*model.FileInfo, error)
	GetWithOptions(page, perPage *int, opt *model.GetFileInfosOptions) ([]*model.FileInfo, error) // Leave perPage and page nil to get all result
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

// model
type (
	AttributeStore interface {
		ModelFields(prefix string) util.AnyArray[string]
		Delete(ids ...string) (int64, error)
		ScanFields(attr *model.Attribute) []interface{}
		Upsert(attr *model.Attribute) (*model.Attribute, error)                       // Upsert inserts or updates given model then returns it
		GetByOption(option *model.AttributeFilterOption) (*model.Attribute, error)    // GetByOption finds and returns 1 model
		FilterbyOption(option *model.AttributeFilterOption) (model.Attributes, error) // FilterbyOption returns a list of attributes by given option
		GetProductTypeAttributes(productTypeID string, unassigned bool, filter *model.AttributeFilterOption) (model.Attributes, error)
		GetPageTypeAttributes(pageTypeID string, unassigned bool) (model.Attributes, error)
	}
	AttributeTranslationStore interface {
	}
	AttributeValueStore interface {
		ScanFields(attributeValue *model.AttributeValue) []interface{}
		ModelFields(prefix string) util.AnyArray[string]
		Count(options *model.AttributeValueFilterOptions) (int64, error)
		Delete(ids ...string) (int64, error)
		Upsert(av *model.AttributeValue) (*model.AttributeValue, error)
		BulkUpsert(transaction store_iface.SqlxTxExecutor, values model.AttributeValues) (model.AttributeValues, error)
		Get(attributeID string) (*model.AttributeValue, error)                                    // Get finds an model value with given id then returns it with an error
		FilterByOptions(options model.AttributeValueFilterOptions) (model.AttributeValues, error) // FilterByOptions finds and returns all matched model values based on given options
	}
	AttributeValueTranslationStore interface {
	}
	AssignedPageAttributeValueStore interface {
		ModelFields(prefix string) util.AnyArray[string]
		Save(assignedPageAttrValue *model.AssignedPageAttributeValue) (*model.AssignedPageAttributeValue, error)                                                 // Save insert given value into database then returns it with an error
		Get(assignedPageAttrValueID string) (*model.AssignedPageAttributeValue, error)                                                                           // Get try finding an value with given id then returns it with an error
		SaveInBulk(assignmentID string, attributeValueIDs []string) ([]*model.AssignedPageAttributeValue, error)                                                 // SaveInBulk inserts multiple values into database then returns them with an error
		SelectForSort(assignmentID string) (assignedPageAttributeValues []*model.AssignedPageAttributeValue, attributeValues []*model.AttributeValue, err error) // SelectForSort uses inner join to find two list: []*assignedPageAttributeValue and []*attributeValue. With given assignedPageAttributeID
		UpdateInBulk(attributeValues []*model.AssignedPageAttributeValue) error                                                                                  // UpdateInBulk use transaction to update all given assigned page model values
	}
	AssignedPageAttributeStore interface {
		ModelFields(prefix string) util.AnyArray[string]
		Save(assignedPageAttr *model.AssignedPageAttribute) (*model.AssignedPageAttribute, error)          // Save inserts given assigned page model into database and returns it with an error
		Get(id string) (*model.AssignedPageAttribute, error)                                               // Get returns an assigned page model with an error
		GetByOption(option *model.AssignedPageAttributeFilterOption) (*model.AssignedPageAttribute, error) // GetByOption try to find an assigned page model with given option. If nothing found, creats new instance with that option and returns such value with an error
	}
	AttributePageStore interface {
		ModelFields(prefix string) util.AnyArray[string]
		Save(page *model.AttributePage) (*model.AttributePage, error)
		Get(pageID string) (*model.AttributePage, error)
		GetByOption(option *model.AttributePageFilterOption) (*model.AttributePage, error)
	}
	AssignedVariantAttributeValueStore interface {
		ModelFields(prefix string) util.AnyArray[string]
		ScanFields(assignedVariantAttributeValue *model.AssignedVariantAttributeValue) []interface{}
		Save(assignedVariantAttrValue *model.AssignedVariantAttributeValue) (*model.AssignedVariantAttributeValue, error)                                              // Save inserts new value into database then returns it with an error
		Get(id string) (*model.AssignedVariantAttributeValue, error)                                                                                                   // Get try finding a value with given id then returns it with an error
		SaveInBulk(assignmentID string, attributeValueIDs []string) ([]*model.AssignedVariantAttributeValue, error)                                                    // SaveInBulk save multiple values into database then returns them
		SelectForSort(assignmentID string) (assignedVariantAttributeValues []*model.AssignedVariantAttributeValue, attributeValues []*model.AttributeValue, err error) // SelectForSort
		UpdateInBulk(attributeValues []*model.AssignedVariantAttributeValue) error                                                                                     // UpdateInBulk use transaction to update given values, then returns an error to indicate if the operation was successful or not
		FilterByOptions(options *model.AssignedVariantAttributeValueFilterOptions) ([]*model.AssignedVariantAttributeValue, error)
	}
	AssignedVariantAttributeStore interface {
		Save(assignedVariantAttribute *model.AssignedVariantAttribute) (*model.AssignedVariantAttribute, error)       // Save insert new instance into database then returns it with an error
		Get(id string) (*model.AssignedVariantAttribute, error)                                                       // Get find assigned variant model from database then returns it with an error
		GetWithOption(option *model.AssignedVariantAttributeFilterOption) (*model.AssignedVariantAttribute, error)    // GetWithOption try finding an assigned variant model with given option. If nothing found, it creates instance with given option. Finally it returns expected value with an error
		FilterByOption(option *model.AssignedVariantAttributeFilterOption) ([]*model.AssignedVariantAttribute, error) // FilterByOption finds and returns a list of assigned variant attributes filtered by given options
	}
	AttributeVariantStore interface {
		ModelFields(prefix string) util.AnyArray[string]
		Save(attributeVariant *model.AttributeVariant) (*model.AttributeVariant, error)
		Get(attributeVariantID string) (*model.AttributeVariant, error)
		GetByOption(option *model.AttributeVariantFilterOption) (*model.AttributeVariant, error) // GetByOption finds 1 model variant with given option.
		FilterByOptions(options *model.AttributeVariantFilterOption) ([]*model.AttributeVariant, error)
	}
	AssignedProductAttributeValueStore interface {
		ModelFields(prefix string) util.AnyArray[string]
		ScanFields(assignedProductAttributeValue *model.AssignedProductAttributeValue) []interface{}
		Save(assignedProductAttrValue *model.AssignedProductAttributeValue) (*model.AssignedProductAttributeValue, error) // Save inserts given instance into database then returns it with an error
		Get(assignedProductAttrValueID string) (*model.AssignedProductAttributeValue, error)                              // Get try finding an instance with given id then returns the value with an error
		SaveInBulk(assignmentID string, attributeValueIDs []string) ([]*model.AssignedProductAttributeValue, error)       // SaveInBulk save multiple values into database
		SelectForSort(assignmentID string) ([]*model.AssignedProductAttributeValue, []*model.AttributeValue, error)       // SelectForSort finds all `*AssignedProductAttributeValue` and related `*AttributeValues` with given `assignmentID`, then returns them with an error.
		UpdateInBulk(attributeValues []*model.AssignedProductAttributeValue) error                                        // UpdateInBulk use transaction to update the given values. Returned error can be `*store.ErrInvalidInput` or `system error`
		FilterByOptions(options *model.AssignedProductAttributeValueFilterOptions) ([]*model.AssignedProductAttributeValue, error)
	}
	AssignedProductAttributeStore interface {
		Save(assignedProductAttribute *model.AssignedProductAttribute) (*model.AssignedProductAttribute, error)    // Save inserts new assgignedProductAttribute into database and returns it with an error
		Get(id string) (*model.AssignedProductAttribute, error)                                                    // Get finds and returns an assignedProductAttribute with en error
		GetWithOption(option *model.AssignedProductAttributeFilterOption) (*model.AssignedProductAttribute, error) // GetWithOption try finding an `AssignedProductAttribute` with given `option`. If nothing found, it creates new instance then returns it with an error
		FilterByOptions(options *model.AssignedProductAttributeFilterOption) ([]*model.AssignedProductAttribute, error)
	}
	AttributeProductStore interface {
		Save(attributeProduct *model.AttributeProduct) (*model.AttributeProduct, error)                // Save inserts given model product relationship into database then returns it and an error
		Get(attributeProductID string) (*model.AttributeProduct, error)                                // Get finds an attributeProduct relationship and returns it with an error
		GetByOption(option *model.AttributeProductFilterOption) (*model.AttributeProduct, error)       // GetByOption returns an attributeProduct with given condition
		FilterByOptions(option *model.AttributeProductFilterOption) ([]*model.AttributeProduct, error) // FilterByOptions returns attributeProducts with given condition
	}
)

// model
type ComplianceStore interface {
	Save(model *model.Compliance) (*model.Compliance, error)
	Update(model *model.Compliance) (*model.Compliance, error)
	Get(id string) (*model.Compliance, error)
	GetAll(offset, limit int) (model.Compliances, error)
	ComplianceExport(model *model.Compliance, cursor model.ComplianceExportCursor, limit int) ([]*model.CompliancePost, model.ComplianceExportCursor, error)
	MessageExport(cursor model.MessageExportCursor, limit int) ([]*model.MessageExport, model.MessageExportCursor, error)
}

// plugin
type PluginConfigurationStore interface {
	GetByOptions(options *model.PluginConfigurationFilterOptions) (*model.PluginConfiguration, error)                // GetByOptions finds and returns 1 plugin configuration with given options
	Upsert(config *model.PluginConfiguration) (*model.PluginConfiguration, error)                                    // Upsert inserts or updates given plugin configuration and returns it
	Get(id string) (*model.PluginConfiguration, error)                                                               // Get finds a plugin configuration with given id then returns it
	FilterPluginConfigurations(options model.PluginConfigurationFilterOptions) ([]*model.PluginConfiguration, error) // FilterPluginConfigurations finds and returns a list of configs with given options then returns them
}

// model
type (
	WishlistStore interface {
		Upsert(wishList *model.Wishlist) (*model.Wishlist, error)                // Upsert inserts or update given model and returns it
		GetByOption(option *model.WishlistFilterOption) (*model.Wishlist, error) // GetByOption finds and returns a slice of wishlists by given option
	}
	WishlistItemStore interface {
		GetById(selector store_iface.SqlxTxExecutor, id string) (*model.WishlistItem, error)                               // GetById returns a model item wish given id
		BulkUpsert(transaction store_iface.SqlxTxExecutor, wishlistItems model.WishlistItems) (model.WishlistItems, error) // Upsert inserts or updates given model item then returns it
		FilterByOption(option *model.WishlistItemFilterOption) ([]*model.WishlistItem, error)                              // FilterByOption finds and returns a slice of model items filtered using given options
		GetByOption(option *model.WishlistItemFilterOption) (*model.WishlistItem, error)                                   // GetByOption finds and returns a model item filtered by given option
		DeleteItemsByOption(transaction store_iface.SqlxTxExecutor, option *model.WishlistItemFilterOption) (int64, error) // DeleteItemsByOption finds and deletes model items that satisfy given filtering options and returns number of items deleted
	}
	WishlistItemProductVariantStore interface {
		Save(wishlistVariant *model.WishlistItemProductVariant) (*model.WishlistItemProductVariant, error)                                             // Save inserts new model product variant relation into database and returns it
		BulkUpsert(transaction store_iface.SqlxTxExecutor, relations []*model.WishlistItemProductVariant) ([]*model.WishlistItemProductVariant, error) // BulkUpsert does bulk update/insert given relations
		GetById(selector store_iface.SqlxTxExecutor, id string) (*model.WishlistItemProductVariant, error)                                             // GetByID returns a model item product variant with given id
		DeleteRelation(relation *model.WishlistItemProductVariant) (int64, error)                                                                      // DeleteRelation deletes a product variant-model item relation and counts numeber of relations left in database
	}
)

// model
type (
	WarehouseStore interface {
		Delete(transaction store_iface.SqlxTxExecutor, ids ...string) error
		Update(warehouse *model.WareHouse) (*model.WareHouse, error)
		ModelFields(prefix string) util.AnyArray[string]
		ScanFields(wh *model.WareHouse) []interface{}
		Save(model *model.WareHouse) (*model.WareHouse, error)                          // Save inserts given model into database then returns it.
		FilterByOprion(option *model.WarehouseFilterOption) ([]*model.WareHouse, error) // FilterByOprion returns a slice of warehouses with given option
		GetByOption(option *model.WarehouseFilterOption) (*model.WareHouse, error)      // GetByOption finds and returns a model filtered given option
		WarehouseByStockID(stockID string) (*model.WareHouse, error)                    // WarehouseByStockID returns 1 model by given stock id
		ApplicableForClickAndCollectNoQuantityCheck(checkoutLines model.CheckoutLines, country model.CountryCode) (model.Warehouses, error)
		ApplicableForClickAndCollectCheckoutLines(checkoutLines model.CheckoutLines, country model.CountryCode) (model.Warehouses, error)
		ApplicableForClickAndCollectOrderLines(orderLines model.OrderLines, country model.CountryCode) (model.Warehouses, error)
	}
	StockStore interface {
		ScanFields(stock *model.Stock) []interface{}
		ModelFields(prefix string) util.AnyArray[string]
		Get(stockID string) (*model.Stock, error)                                                                       // Get finds and returns stock with given stockID. Returned error could be either (nil, *ErrNotFound, error)
		FilterForCountryAndChannel(options *model.StockFilterForCountryAndChannel) ([]*model.Stock, error)              // FilterForCountryAndChannel finds and returns stocks with given options
		FilterVariantStocksForCountry(options *model.StockFilterForCountryAndChannel) ([]*model.Stock, error)           // FilterVariantStocksForCountry finds and returns stocks with given options
		FilterProductStocksForCountryAndChannel(options *model.StockFilterForCountryAndChannel) ([]*model.Stock, error) // FilterProductStocksForCountryAndChannel finds and returns stocks with given options
		ChangeQuantity(stockID string, quantity int) error                                                              // ChangeQuantity reduce or increase the quantity of given stock
		FilterByOption(options *model.StockFilterOption) ([]*model.Stock, error)                                        // FilterByOption finds and returns a slice of stocks that satisfy given option
		BulkUpsert(transaction store_iface.SqlxTxExecutor, stocks []*model.Stock) ([]*model.Stock, error)               // BulkUpsert performs upserts or inserts given stocks, then returns them
		FilterForChannel(options *model.StockFilterForChannelOption) (squirrel.Sqlizer, []*model.Stock, error)          // FilterForChannel finds and returns stocks that satisfy given options
	}
	AllocationStore interface {
		BulkUpsert(transaction store_iface.SqlxTxExecutor, allocations []*model.Allocation) ([]*model.Allocation, error)          // BulkUpsert performs update, insert given allocations then returns them afterward
		Get(allocationID string) (*model.Allocation, error)                                                                       // Get find and returns allocation with given id
		FilterByOption(transaction store_iface.SqlxTxExecutor, option *model.AllocationFilterOption) ([]*model.Allocation, error) // FilterbyOption finds and returns a list of allocations based on given option
		BulkDelete(transaction store_iface.SqlxTxExecutor, allocationIDs []string) error                                          // BulkDelete perform bulk deletes given allocations.
		CountAvailableQuantityForStock(stock *model.Stock) (int, error)                                                           // CountAvailableQuantityForStock counts and returns available quantity of given stock
	}
	WarehouseShippingZoneStore interface {
		Delete(transaction store_iface.SqlxTxExecutor, options *model.WarehouseShippingZoneFilterOption) error
		ModelFields(prefix string) util.AnyArray[string]
		Save(transaction store_iface.SqlxTxExecutor, warehouseShippingZones []*model.WarehouseShippingZone) ([]*model.WarehouseShippingZone, error) // Save inserts given model-model zone relation into database
		FilterByCountryCodeAndChannelID(countryCode, channelID string) ([]*model.WarehouseShippingZone, error)
		FilterByOptions(options *model.WarehouseShippingZoneFilterOption) ([]*model.WarehouseShippingZone, error)
	}
	PreorderAllocationStore interface {
		ModelFields(prefix string) util.AnyArray[string]
		BulkCreate(transaction store_iface.SqlxTxExecutor, preorderAllocations []*model.PreorderAllocation) ([]*model.PreorderAllocation, error) // BulkCreate bulk inserts given preorderAllocations and returns them
		ScanFields(preorderAllocation *model.PreorderAllocation) []interface{}
		FilterByOption(options *model.PreorderAllocationFilterOption) ([]*model.PreorderAllocation, error) // FilterByOption finds and returns a list of preorder allocations filtered using given options
		Delete(transaction store_iface.SqlxTxExecutor, preorderAllocationIDs ...string) error              // Delete deletes preorder-allocations by given ids
	}
)

// model
type (
	ShippingZoneStore interface {
		ModelFields(prefix string) util.AnyArray[string]
		ScanFields(shippingZone *model.ShippingZone) []interface{}
		Upsert(shippingZone *model.ShippingZone) (*model.ShippingZone, error)                 // Upsert depends on given model zone's Id to decide update or insert the zone
		Get(shippingZoneID string) (*model.ShippingZone, error)                               // Get finds 1 model zone for given shippingZoneID
		FilterByOption(option *model.ShippingZoneFilterOption) ([]*model.ShippingZone, error) // FilterByOption finds a list of model zones based on given option
		CountByOptions(options *model.ShippingZoneFilterOption) (int64, error)
	}
	ShippingMethodStore interface {
		ModelFields(prefix string) util.AnyArray[string]
		Upsert(method *model.ShippingMethod) (*model.ShippingMethod, error)                                                                                                                 // Upsert bases on given method's Id to decide update or insert it
		Get(methodID string) (*model.ShippingMethod, error)                                                                                                                                 // Get finds and returns a model method with given id
		ApplicableShippingMethods(price *goprices.Money, channelID string, weight *measurement.Weight, countryCode model.CountryCode, productIDs []string) ([]*model.ShippingMethod, error) // ApplicableShippingMethods finds all model methods with given conditions
		GetbyOption(options *model.ShippingMethodFilterOption) (*model.ShippingMethod, error)                                                                                               // GetbyOption finds and returns a model method that satisfy given options
		FilterByOptions(options *model.ShippingMethodFilterOption) ([]*model.ShippingMethod, error)
	}
	ShippingMethodPostalCodeRuleStore interface {
		ModelFields(prefix string) util.AnyArray[string]
		ScanFields(rule *model.ShippingMethodPostalCodeRule) []interface{}
		FilterByOptions(options *model.ShippingMethodPostalCodeRuleFilterOptions) ([]*model.ShippingMethodPostalCodeRule, error)
	}
	ShippingMethodChannelListingStore interface {
		BulkDelete(transaction store_iface.SqlxTxExecutor, ids []string) error
		Upsert(listing *model.ShippingMethodChannelListing) (*model.ShippingMethodChannelListing, error)                      // Upsert depends on given listing's Id to decide whether to save or update the listing
		Get(listingID string) (*model.ShippingMethodChannelListing, error)                                                    // Get finds a model method channel listing with given listingID
		FilterByOption(option *model.ShippingMethodChannelListingFilterOption) ([]*model.ShippingMethodChannelListing, error) // FilterByOption returns a list of model method channel listings based on given option. result sorted by creation time ASC
	}
	ShippingMethodTranslationStore interface {
	}
	ShippingZoneChannelStore interface {
		BulkDelete(transaction store_iface.SqlxTxExecutor, relations []*model.ShippingZoneChannel) error
		BulkSave(transaction store_iface.SqlxTxExecutor, relations []*model.ShippingZoneChannel) ([]*model.ShippingZoneChannel, error)
		FilterByOptions(options *model.ShippingZoneChannelFilterOptions) ([]*model.ShippingZoneChannel, error)
	}
	ShippingMethodExcludedProductStore interface {
		Save(instance *model.ShippingMethodExcludedProduct) (*model.ShippingMethodExcludedProduct, error) // Save inserts given ShippingMethodExcludedProduct into database then returns it
		FilterByOptions(options *model.ShippingMethodExcludedProductFilterOptions) ([]*model.ShippingMethodExcludedProduct, error)
	}
)

// product
type (
	CollectionTranslationStore interface {
	}
	CollectionChannelListingStore interface {
		Delete(transaction store_iface.SqlxTxExecutor, options *model.CollectionChannelListingFilterOptions) error
		Upsert(transaction store_iface.SqlxTxExecutor, relations ...*model.CollectionChannelListing) ([]*model.CollectionChannelListing, error)
		FilterByOptions(options *model.CollectionChannelListingFilterOptions) ([]*model.CollectionChannelListing, error)
	}
	CollectionStore interface {
		ModelFields(prefix string) util.AnyArray[string]
		Upsert(collection *model.Collection) (*model.Collection, error)                   // Upsert depends on given collection's Id property to decide update or insert the collection
		Get(collectionID string) (*model.Collection, error)                               // Get finds and returns collection with given collectionID
		FilterByOption(option *model.CollectionFilterOption) ([]*model.Collection, error) // FilterByOption finds and returns a list of collections satisfy the given option
		ScanFields(col *model.Collection) []interface{}
		Delete(ids ...string) error
	}
	CollectionProductStore interface {
		Delete(transaction store_iface.SqlxTxExecutor, options *model.CollectionProductFilterOptions) error
		BulkSave(transaction store_iface.SqlxTxExecutor, relations []*model.CollectionProduct) ([]*model.CollectionProduct, error)
		FilterByOptions(options *model.CollectionProductFilterOptions) ([]*model.CollectionProduct, error)
	}
	VariantMediaStore interface {
		FilterByOptions(options *model.VariantMediaFilterOptions) ([]*model.VariantMedia, error)
	}
	ProductMediaStore interface {
		Upsert(media *model.ProductMedia) (*model.ProductMedia, error)                        // Upsert depends on given media's Id property to decide insert or update it
		Get(id string) (*model.ProductMedia, error)                                           // Get finds and returns 1 product media with given id
		FilterByOption(option *model.ProductMediaFilterOption) ([]*model.ProductMedia, error) // FilterByOption finds and returns a list of product medias with given id
	}
	DigitalContentUrlStore interface {
		Upsert(contentURL *model.DigitalContentUrl) (*model.DigitalContentUrl, error) // Upsert inserts or updates given digital content url into database then returns it
		Get(id string) (*model.DigitalContentUrl, error)                              // Get finds and returns a digital content url with given id
		FilterByOptions(options *model.DigitalContentUrlFilterOptions) ([]*model.DigitalContentUrl, error)
	}
	DigitalContentStore interface {
		ModelFields(prefix string) util.AnyArray[string]
		ScanFields(content *model.DigitalContent) []interface{}
		Save(content *model.DigitalContent) (*model.DigitalContent, error)                        // Save inserts given digital content into database then returns it
		GetByOption(option *model.DigitalContentFilterOption) (*model.DigitalContent, error)      // GetByOption finds and returns 1 digital content filtered using given option
		FilterByOption(option *model.DigitalContentFilterOption) ([]*model.DigitalContent, error) //
	}
	ProductVariantChannelListingStore interface {
		ModelFields(prefix string) util.AnyArray[string]
		ScanFields(listing *model.ProductVariantChannelListing) []interface{}
		Save(variantChannelListing *model.ProductVariantChannelListing) (*model.ProductVariantChannelListing, error)                                                    // Save insert given value into database then returns it with an error
		Get(variantChannelListingID string) (*model.ProductVariantChannelListing, error)                                                                                // Get finds and returns 1 product variant channel listing based on given variantChannelListingID
		FilterbyOption(transaction store_iface.SqlxTxExecutor, option *model.ProductVariantChannelListingFilterOption) ([]*model.ProductVariantChannelListing, error)   // FilterbyOption finds and returns all product variant channel listings filterd using given option
		BulkUpsert(transaction store_iface.SqlxTxExecutor, variantChannelListings []*model.ProductVariantChannelListing) ([]*model.ProductVariantChannelListing, error) // BulkUpsert performs bulk upsert given product variant channel listings then returns them
	}
	ProductVariantTranslationStore interface {
		Upsert(translation *model.ProductVariantTranslation) (*model.ProductVariantTranslation, error)                  // Upsert inserts or updates given translation then returns it
		Get(translationID string) (*model.ProductVariantTranslation, error)                                             // Get finds and returns 1 product variant translation with given id
		FilterByOption(option *model.ProductVariantTranslationFilterOption) ([]*model.ProductVariantTranslation, error) // FilterByOption finds and returns product variant translations filtered using given options
	}
	ProductVariantStore interface {
		ModelFields(prefix string) util.AnyArray[string]
		ScanFields(variant *model.ProductVariant) []interface{}
		Save(transaction store_iface.SqlxTxExecutor, variant *model.ProductVariant) (*model.ProductVariant, error)   // Save inserts product variant instance to database
		Get(id string) (*model.ProductVariant, error)                                                                // Get returns a product variant with given id
		GetWeight(productVariantID string) (*measurement.Weight, error)                                              // GetWeight returns weight of given product variant
		GetByOrderLineID(orderLineID string) (*model.ProductVariant, error)                                          // GetByOrderLineID finds and returns a product variant by given orderLineID
		FilterByOption(option *model.ProductVariantFilterOption) ([]*model.ProductVariant, error)                    // FilterByOption finds and returns product variants based on given option
		Update(transaction store_iface.SqlxTxExecutor, variant *model.ProductVariant) (*model.ProductVariant, error) // Update updates given product variant and returns it
	}
	ProductChannelListingStore interface {
		ModelFields(prefix string) util.AnyArray[string]
		BulkUpsert(transaction store_iface.SqlxTxExecutor, listings []*model.ProductChannelListing) ([]*model.ProductChannelListing, error) // BulkUpsert performs bulk upsert on given product channel listings
		Get(channelListingID string) (*model.ProductChannelListing, error)                                                                  // Get try finding a product channel listing, then returns it with an error
		FilterByOption(option *model.ProductChannelListingFilterOption) ([]*model.ProductChannelListing, error)                             // FilterByOption filter a list of product channel listings by given option. Then returns them with an error
	}
	ProductTranslationStore interface {
		Upsert(translation *model.ProductTranslation) (*model.ProductTranslation, error)                  // Upsert inserts or update given translation
		Get(translationID string) (*model.ProductTranslation, error)                                      // Get finds and returns a product translation by given id
		FilterByOption(option *model.ProductTranslationFilterOption) ([]*model.ProductTranslation, error) // FilterByOption finds and returns product translations filtered using given options
	}
	ProductTypeStore interface {
		ModelFields(prefix string) util.AnyArray[string]
		FilterbyOption(options *model.ProductTypeFilterOption) ([]*model.ProductType, error)
		Save(productType *model.ProductType) (*model.ProductType, error)                      // Save try inserting new product type into database then returns it
		FilterProductTypesByCheckoutToken(checkoutToken string) ([]*model.ProductType, error) // FilterProductTypesByCheckoutToken is used to check if a model requires model
		ProductTypesByProductIDs(productIDs []string) ([]*model.ProductType, error)           // ProductTypesByProductIDs returns all product types belong to given products
		ProductTypeByProductVariantID(variantID string) (*model.ProductType, error)           // ProductTypeByProductVariantID finds and returns 1 product type that is related to given product variant
		GetByOption(options *model.ProductTypeFilterOption) (*model.ProductType, error)       // GetByOption finds and returns a product type with given options
		Count(options *model.ProductTypeFilterOption) (int64, error)
	}
	CategoryTranslationStore interface{}
	CategoryStore            interface {
		ModelFields(prefix string) util.AnyArray[string]
		Upsert(category *model.Category) (*model.Category, error)                                 // Upsert depends on given category's Id field to decide update or insert it
		Get(ctx context.Context, categoryID string, allowFromCache bool) (*model.Category, error) // Get finds and returns a category with given id
		GetByOption(option *model.CategoryFilterOption) (*model.Category, error)                  // GetByOption finds and returns 1 category satisfy given option
		FilterByOption(option *model.CategoryFilterOption) ([]*model.Category, error)             // FilterByOption finds and returns a list of categories satisfy given option
	}
	ProductStore interface {
		ModelFields(prefix string) util.AnyArray[string]
		ScanFields(prd *model.Product) []interface{}
		Save(prd *model.Product) (*model.Product, error)
		GetByOption(option *model.ProductFilterOption) (*model.Product, error)      // GetByOption finds and returns 1 product that satisfies given option
		FilterByOption(option *model.ProductFilterOption) ([]*model.Product, error) // FilterByOption finds and returns all products that satisfy given option
		PublishedProducts(channelSlug string) ([]*model.Product, error)             // FilterPublishedProducts finds and returns products that belong to given channel slug and are published
		NotPublishedProducts(channelSlug string) ([]*struct {
			model.Product
			IsPublished     bool
			PublicationDate *timemodule.Time
		}, error) // FilterNotPublishedProducts finds all not published products belong to given channel
		PublishedWithVariants(channel_SlugOrID string) squirrel.SelectBuilder                                                              // PublishedWithVariants finds and returns products.
		VisibleToUserProductsQuery(channel_SlugOrID string, userHasOneOfProductpermissions bool) squirrel.SelectBuilder                    // FilterVisibleToUserProduct finds and returns all products that are visible to requesting user.
		SelectForUpdateDiscountedPricesOfCatalogues(productIDs, categoryIDs, collectionIDs, variantIDs []string) ([]*model.Product, error) // SelectForUpdateDiscountedPricesOfCatalogues finds and returns product based on given ids lists.
		AdvancedFilterQueryBuilder(input *model.ExportProductsFilterOptions) squirrel.SelectBuilder                                        // AdvancedFilterQueryBuilder advancedly finds products, filtered using given options
		FilterByQuery(query squirrel.SelectBuilder) (model.Products, error)                                                                // FilterByQuery finds and returns products with given query, limit, createdAtGt
		CountByCategoryIDs(categoryIDs []string) ([]*model.ProductCountByCategoryID, error)
	}
)

// model
type (
	PaymentStore interface {
		ScanFields(payMent *model.Payment) []interface{}
		Save(transaction store_iface.SqlxTxExecutor, model *model.Payment) (*model.Payment, error)                               // Save save model instance into database
		Get(transaction store_iface.SqlxTxExecutor, id string, lockForUpdate bool) (*model.Payment, error)                       // Get returns a model with given id. `lockForUpdate` is true if you want to add "FOR UPDATE" to sql
		Update(transaction store_iface.SqlxTxExecutor, model *model.Payment) (*model.Payment, error)                             // Update updates given model and returns new updated model
		CancelActivePaymentsOfCheckout(checkoutToken string) error                                                               // CancelActivePaymentsOfCheckout inactivate all payments that belong to given model and in active status
		FilterByOption(option *model.PaymentFilterOption) ([]*model.Payment, error)                                              // FilterByOption finds and returns a list of payments that satisfy given option
		UpdatePaymentsOfCheckout(transaction store_iface.SqlxTxExecutor, checkoutToken string, option *model.PaymentPatch) error // UpdatePaymentsOfCheckout updates payments of given model
		PaymentOwnedByUser(userID, paymentID string) (bool, error)
	}
	PaymentTransactionStore interface {
		Save(transaction store_iface.SqlxTxExecutor, paymentTransaction *model.PaymentTransaction) (*model.PaymentTransaction, error) // Save inserts new model transaction into database
		Get(id string) (*model.PaymentTransaction, error)                                                                             // Get returns a model transaction with given id
		Update(transaction *model.PaymentTransaction) (*model.PaymentTransaction, error)                                              // Update updates given transaction and returns updated one
		FilterByOption(option *model.PaymentTransactionFilterOpts) ([]*model.PaymentTransaction, error)                               // FilterByOption finds and returns a list of transactions with given option
	}
)

// page
type (
	PageTypeStore interface {
	}
	PageTranslationStore interface {
	}
	PageStore interface {
		FilterByOptions(options *model.PageFilterOptions) ([]*model.Page, error)
	}
)

// order
type (
	OrderLineStore interface {
		ScanFields(orderLine *model.OrderLine) []interface{}
		ModelFields(prefix string) util.AnyArray[string]
		Upsert(transaction store_iface.SqlxTxExecutor, orderLine *model.OrderLine) (*model.OrderLine, error)          // Upsert depends on given orderLine's Id to decide to update or save it
		Get(id string) (*model.OrderLine, error)                                                                      // Get returns a order line with id of given id
		BulkDelete(orderLineIDs []string) error                                                                       // BulkDelete delete all given order lines. NOTE: validate given ids are valid uuids before calling me
		FilterbyOption(option *model.OrderLineFilterOption) ([]*model.OrderLine, error)                               // FilterbyOption finds and returns order lines by given option
		BulkUpsert(transaction store_iface.SqlxTxExecutor, orderLines []*model.OrderLine) ([]*model.OrderLine, error) // BulkUpsert performs upsert multiple order lines in once
	}
	OrderStore interface {
		ModelFields(prefix string) util.AnyArray[string]
		ScanFields(holder *model.Order) []interface{}
		Get(id string) (*model.Order, error)                                                              // Get find order in database with given id
		FilterByOption(option *model.OrderFilterOption) ([]*model.Order, error)                           // FilterByOption returns a list of orders, filtered by given option
		BulkUpsert(transaction store_iface.SqlxTxExecutor, orders []*model.Order) ([]*model.Order, error) // BulkUpsert performs bulk upsert given orders
	}
	OrderEventStore interface {
		Save(transaction store_iface.SqlxTxExecutor, orderEvent *model.OrderEvent) (*model.OrderEvent, error) // Save inserts given order event into database then returns it
		Get(orderEventID string) (*model.OrderEvent, error)                                                   // Get finds order event with given id then returns it
		FilterByOptions(options *model.OrderEventFilterOptions) ([]*model.OrderEvent, error)
	}
	FulfillmentLineStore interface {
		ModelFields(prefix string) util.AnyArray[string]
		Save(fulfillmentLine *model.FulfillmentLine) (*model.FulfillmentLine, error)
		Get(id string) (*model.FulfillmentLine, error)
		FilterbyOption(option *model.FulfillmentLineFilterOption) ([]*model.FulfillmentLine, error)                                     // FilterbyOption finds and returns a list of fulfillment lines by given option
		BulkUpsert(transaction store_iface.SqlxTxExecutor, fulfillmentLines []*model.FulfillmentLine) ([]*model.FulfillmentLine, error) // BulkUpsert upsert given fulfillment lines
		DeleteFulfillmentLinesByOption(transaction store_iface.SqlxTxExecutor, option *model.FulfillmentLineFilterOption) error         // DeleteFulfillmentLinesByOption filters fulfillment lines by given option, then deletes them
	}
	FulfillmentStore interface {
		ModelFields(prefix string) util.AnyArray[string]
		ScanFields(holder *model.Fulfillment) []interface{}
		Upsert(transaction store_iface.SqlxTxExecutor, fulfillment *model.Fulfillment) (*model.Fulfillment, error)                  // Upsert depends on given fulfillment's Id to decide update or insert it
		Get(id string) (*model.Fulfillment, error)                                                                                  // Get finds and return a fulfillment by given id
		GetByOption(transaction store_iface.SqlxTxExecutor, option *model.FulfillmentFilterOption) (*model.Fulfillment, error)      // GetByOption returns 1 fulfillment, filtered by given option
		FilterByOption(transaction store_iface.SqlxTxExecutor, option *model.FulfillmentFilterOption) ([]*model.Fulfillment, error) // FilterByOption finds and returns a slice of fulfillments by given option
		BulkDeleteFulfillments(transaction store_iface.SqlxTxExecutor, fulfillments model.Fulfillments) error                       // BulkDeleteFulfillments deletes given fulfillments
	}
)

// menu
type (
	MenuItemTranslationStore interface {
	}
	MenuStore interface {
		Save(menu *model.Menu) (*model.Menu, error) // Save insert given menu into database and returns it
		GetByOptions(options *model.MenuFilterOptions) (*model.Menu, error)
		FilterByOptions(options *model.MenuFilterOptions) ([]*model.Menu, error)
	}
	MenuItemStore interface {
		Save(menuItem *model.MenuItem) (*model.MenuItem, error) // Save insert given menu item into database and returns it
		GetByOptions(options *model.MenuItemFilterOptions) (*model.MenuItem, error)
		FilterByOptions(options *model.MenuItemFilterOptions) ([]*model.MenuItem, error)
	}
)

// invoice
type (
	InvoiceEventStore interface {
		Upsert(invoiceEvent *model.InvoiceEvent) (*model.InvoiceEvent, error) // Upsert depends on given invoice event's Id to update/insert it
		Get(invoiceEventID string) (*model.InvoiceEvent, error)               // Get finds and returns 1 invoice event
	}
	InvoiceStore interface {
		Upsert(invoice *model.Invoice) (*model.Invoice, error) // Upsert depends on given invoice Id to update/insert it
		Get(invoiceID string) (*model.Invoice, error)          // Get finds and returns 1 invoice
		FilterByOptions(options *model.InvoiceFilterOptions) ([]*model.Invoice, error)
		Delete(transaction store_iface.SqlxTxExecutor, ids []string) error
	}
)

// giftcard related stores
type (
	GiftCardStore interface {
		DeleteGiftcards(transaction store_iface.SqlxTxExecutor, ids []string) error
		BulkUpsert(transaction store_iface.SqlxTxExecutor, giftCards ...*model.GiftCard) ([]*model.GiftCard, error) // BulkUpsert depends on given giftcards's Id properties then perform according operation
		GetById(id string) (*model.GiftCard, error)                                                                 // GetById returns a giftcard instance that has id of given id
		FilterByOption(option *model.GiftCardFilterOption) ([]*model.GiftCard, error)                               // FilterByOption finds giftcards wth option
		// DeactivateOrderGiftcards update giftcards
		// which have giftcard events with type == 'bought', parameters.order_id == given order id
		// by setting their IsActive model to false
		DeactivateOrderGiftcards(orderID string) ([]string, error)
	}
	GiftcardEventStore interface {
		Save(event *model.GiftCardEvent) (*model.GiftCardEvent, error)                                                     // Save insdert given giftcard event into database then returns it
		Get(id string) (*model.GiftCardEvent, error)                                                                       // Get finds and returns a giftcard event found by given id
		BulkUpsert(transaction store_iface.SqlxTxExecutor, events ...*model.GiftCardEvent) ([]*model.GiftCardEvent, error) // BulkUpsert upserts and returns given giftcard events
		FilterByOptions(options *model.GiftCardEventFilterOption) ([]*model.GiftCardEvent, error)                          // FilterByOptions finds and returns a list of giftcard events with given options
	}
	GiftCardOrderStore interface {
		Save(giftcardOrder *model.OrderGiftCard) (*model.OrderGiftCard, error)                                                     // Save inserts new giftcard-order relation into database then returns it
		Get(id string) (*model.OrderGiftCard, error)                                                                               // Get returns giftcard-order relation table with given id
		BulkUpsert(transaction store_iface.SqlxTxExecutor, orderGiftcards ...*model.OrderGiftCard) ([]*model.OrderGiftCard, error) // BulkUpsert upserts given order-giftcard relations and returns it
		FilterByOptions(options *model.OrderGiftCardFilterOptions) ([]*model.OrderGiftCard, error)
	}
	GiftCardCheckoutStore interface {
		Save(giftcardOrder *model.GiftCardCheckout) (*model.GiftCardCheckout, error) // Save inserts new giftcard-model relation into database then returns it
		Get(id string) (*model.GiftCardCheckout, error)                              // Get returns giftcard-model relation table with given id
		Delete(giftcardID string, checkoutID string) error                           // Delete deletes a giftcard-model relation with given id
	}
)

// discount
type (
	OrderDiscountStore interface {
		Upsert(transaction store_iface.SqlxTxExecutor, orderDiscount *model.OrderDiscount) (*model.OrderDiscount, error) // Upsert depends on given order discount's Id property to decide to update/insert it
		Get(orderDiscountID string) (*model.OrderDiscount, error)                                                        // Get finds and returns an order discount with given id
		FilterbyOption(option *model.OrderDiscountFilterOption) ([]*model.OrderDiscount, error)                          // FilterbyOption filters order discounts that satisfy given option, then returns them
		BulkDelete(orderDiscountIDs []string) error                                                                      // BulkDelete perform bulk delete all given order discount ids
	}
	DiscountSaleTranslationStore interface {
	}
	DiscountSaleChannelListingStore interface {
		ModelFields(prefix string) util.AnyArray[string]
		Save(saleChannelListing *model.SaleChannelListing) (*model.SaleChannelListing, error) // Save insert given instance into database then returns it
		Get(saleChannelListingID string) (*model.SaleChannelListing, error)                   // Get finds and returns sale channel listing with given id
		// SaleChannelListingsWithOption finds a list of sale channel listings plus foreign channel slugs
		SaleChannelListingsWithOption(option *model.SaleChannelListingFilterOption) ([]*model.SaleChannelListing, error)
	}
	VoucherTranslationStore interface {
		Save(translation *model.VoucherTranslation) (*model.VoucherTranslation, error)                    // Save inserts given translation into database and returns it
		Get(id string) (*model.VoucherTranslation, error)                                                 // Get finds and returns a voucher translation with given id
		FilterByOption(option *model.VoucherTranslationFilterOption) ([]*model.VoucherTranslation, error) // FilterByOption returns a list of voucher translations filtered using given options
		GetByOption(option *model.VoucherTranslationFilterOption) (*model.VoucherTranslation, error)      // GetByOption finds and returns 1 voucher translation by given options
	}
	DiscountSaleStore interface {
		Upsert(sale *model.Sale) (*model.Sale, error)                              // Upsert bases on sale's Id to decide to update or insert given sale
		Get(saleID string) (*model.Sale, error)                                    // Get finds and returns a sale with given saleID
		FilterSalesByOption(option *model.SaleFilterOption) ([]*model.Sale, error) // FilterSalesByOption filter sales by option
	}
	VoucherChannelListingStore interface {
		Upsert(voucherChannelListing *model.VoucherChannelListing) (*model.VoucherChannelListing, error)        // upsert check given listing's Id to decide whether to create or update it. Then returns a listing with an error
		Get(voucherChannelListingID string) (*model.VoucherChannelListing, error)                               // Get finds a listing with given id, then returns it with an error
		FilterbyOption(option *model.VoucherChannelListingFilterOption) ([]*model.VoucherChannelListing, error) // FilterbyOption finds and returns a list of voucher channel listing relationship instances filtered by given option
	}
	DiscountVoucherStore interface {
		ModelFields(prefix string) util.AnyArray[string]
		ScanFields(voucher *model.Voucher) []interface{}
		Upsert(voucher *model.Voucher) (*model.Voucher, error)                              // Upsert saves or updates given voucher then returns it with an error
		Get(voucherID string) (*model.Voucher, error)                                       // Get finds a voucher with given id, then returns it with an error
		FilterVouchersByOption(option *model.VoucherFilterOption) ([]*model.Voucher, error) // FilterVouchersByOption finds vouchers bases on given option.
		ExpiredVouchers(date *timemodule.Time) ([]*model.Voucher, error)                    // ExpiredVouchers finds and returns vouchers that are expired before given date
		GetByOptions(options *model.VoucherFilterOption) (*model.Voucher, error)            // GetByOptions finds and returns 1 voucher filtered using given options
	}
	VoucherCategoryStore interface {
		Upsert(voucherCategory *model.VoucherCategory) (*model.VoucherCategory, error) // Upsert saves or updates given voucher category then returns it with an error
		Get(voucherCategoryID string) (*model.VoucherCategory, error)                  // Get finds a voucher category with given id, then returns it with an error
		FilterByOptions(options *model.VoucherCategoryFilterOption) ([]*model.VoucherCategory, error)
	}
	VoucherCollectionStore interface {
		Upsert(voucherCollection *model.VoucherCollection) (*model.VoucherCollection, error) // Upsert saves or updates given voucher collection then returns it with an error
		Get(voucherCollectionID string) (*model.VoucherCollection, error)                    // Get finds a voucher collection with given id, then returns it with an error
		FilterByOptions(options *model.VoucherCollectionFilterOptions) ([]*model.VoucherCollection, error)
	}
	VoucherProductStore interface {
		Upsert(voucherProduct *model.VoucherProduct) (*model.VoucherProduct, error) // Upsert saves or updates given voucher product then returns it with an error
		Get(voucherProductID string) (*model.VoucherProduct, error)                 // Get finds a voucher product with given id, then returns it with an error
		FilterByOptions(options *model.VoucherProductFilterOptions) ([]*model.VoucherProduct, error)
	}
	VoucherCustomerStore interface {
		Save(voucherCustomer *model.VoucherCustomer) (*model.VoucherCustomer, error)                  // Save inserts given voucher customer instance into database ands returns it
		DeleteInBulk(options *model.VoucherCustomerFilterOption) error                                // DeleteInBulk deletes given voucher-customers with given id
		GetByOption(options *model.VoucherCustomerFilterOption) (*model.VoucherCustomer, error)       // GetByOption finds and returns a voucher customer with given options
		FilterByOptions(options *model.VoucherCustomerFilterOption) ([]*model.VoucherCustomer, error) // FilterByOptions finds and returns a slice of voucher customers by given options
	}
	SaleCategoryRelationStore interface {
		Save(relation *model.SaleCategoryRelation) (*model.SaleCategoryRelation, error)                               // Save inserts given sale-category relation into database
		Get(relationID string) (*model.SaleCategoryRelation, error)                                                   // Get returns 1 sale-category relation with given id
		SaleCategoriesByOption(option *model.SaleCategoryRelationFilterOption) ([]*model.SaleCategoryRelation, error) // SaleCategoriesByOption returns a slice of sale-category relations with given option
	}
	SaleProductRelationStore interface {
		Save(relation *model.SaleProductRelation) (*model.SaleProductRelation, error)                             // Save inserts given sale-product relation into database then returns it
		Get(relationID string) (*model.SaleProductRelation, error)                                                // Get finds and returns a sale-product relation with given id
		SaleProductsByOption(option *model.SaleProductRelationFilterOption) ([]*model.SaleProductRelation, error) // SaleProductsByOption returns a slice of sale-product relations, filtered by given option
	}
	SaleCollectionRelationStore interface {
		Save(relation *model.SaleCollectionRelation) (*model.SaleCollectionRelation, error)                       // Save insert given sale-collection relation into database
		Get(relationID string) (*model.SaleCollectionRelation, error)                                             // Get finds and returns a sale-collection relation with given id
		FilterByOption(option *model.SaleCollectionRelationFilterOption) ([]*model.SaleCollectionRelation, error) // FilterByOption returns a list of collections filtered based on given option
	}
	VoucherProductVariantStore interface {
		FilterByOptions(options *model.VoucherProductVariantFilterOption) ([]*model.VoucherProductVariant, error)
	}
	SaleProductVariantStore interface {
		Upsert(relation *model.SaleProductVariant) (*model.SaleProductVariant, error)                      // Upsert inserts/updates given sale-product variant relation into database, then returns it
		FilterByOption(options *model.SaleProductVariantFilterOption) ([]*model.SaleProductVariant, error) // FilterByOption finds and returns a list of sale-product variants filtered using given options
	}
)

// csv
type (
	CsvExportEventStore interface {
		Save(event *model.ExportEvent) (*model.ExportEvent, error)                           // Save inserts given export event into database then returns it
		FilterByOption(options *model.ExportEventFilterOption) ([]*model.ExportEvent, error) // FilterByOption finds and returns a list of export events filtered using given option
	}
	CsvExportFileStore interface {
		Save(file *model.ExportFile) (*model.ExportFile, error) // Save inserts given export file into database then returns it
		Get(id string) (*model.ExportFile, error)               // Get finds and returns an export file found using given id
	}
)

// model
type (
	CheckoutLineStore interface {
		ModelFields(prefix string) util.AnyArray[string]
		ScanFields(line *model.CheckoutLine) []interface{}
		Upsert(checkoutLine *model.CheckoutLine) (*model.CheckoutLine, error)               // Upsert checks whether to update or insert given model line then performs according operation
		Get(id string) (*model.CheckoutLine, error)                                         // Get returns a model line with given id
		CheckoutLinesByCheckoutID(checkoutID string) ([]*model.CheckoutLine, error)         // CheckoutLinesByCheckoutID returns a list of model lines that belong to given model
		DeleteLines(transaction store_iface.SqlxTxExecutor, checkoutLineIDs []string) error // DeleteLines deletes all model lines with given uuids
		BulkUpdate(checkoutLines []*model.CheckoutLine) error                               // BulkUpdate receives a list of modified model lines, updates them in bulk.
		BulkCreate(checkoutLines []*model.CheckoutLine) ([]*model.CheckoutLine, error)      // BulkCreate takes a list of raw model lines, save them into database then returns them fully with an error
		// CheckoutLinesByCheckoutWithPrefetch finds all model lines belong to given model
		//
		// and prefetch all related product variants, products
		//
		// this borrows the idea from Django's prefetch_related() method
		CheckoutLinesByCheckoutWithPrefetch(checkoutID string) ([]*model.CheckoutLine, []*model.ProductVariant, []*model.Product, error)
		TotalWeightForCheckoutLines(checkoutLineIDs []string) (*measurement.Weight, error)           // TotalWeightForCheckoutLines calculate total weight for given model lines
		CheckoutLinesByOption(option *model.CheckoutLineFilterOption) ([]*model.CheckoutLine, error) // CheckoutLinesByOption finds and returns model lines filtered using given option
	}
	CheckoutStore interface {
		ModelFields(prefix string) util.AnyArray[string]
		Upsert(transaction store_iface.SqlxTxExecutor, checkouts []*model.Checkout) ([]*model.Checkout, error)    // Upsert depends on given model's Token property to decide to update or insert it
		FetchCheckoutLinesAndPrefetchRelatedValue(ckout *model.Checkout) ([]*model.CheckoutLineInfo, error)       // FetchCheckoutLinesAndPrefetchRelatedValue Fetch model lines as CheckoutLineInfo objects.
		GetByOption(option *model.CheckoutFilterOption) (*model.Checkout, error)                                  // GetByOption finds and returns 1 model based on given option
		FilterByOption(option *model.CheckoutFilterOption) ([]*model.Checkout, error)                             // FilterByOption finds and returns a list of model based on given option
		DeleteCheckoutsByOption(transaction store_iface.SqlxTxExecutor, option *model.CheckoutFilterOption) error // DeleteCheckoutsByOption deletes model row(s) from database, filtered using given option.  It returns an error indicating if the operation was performed successfully.
		CountCheckouts(options *model.CheckoutFilterOption) (int64, error)
	}
)

// channel
type ChannelStore interface {
	ModelFields(prefix string) util.AnyArray[string]
	ScanFields(chanNel *model.Channel) []interface{}
	Get(id string) (*model.Channel, error)                                      // Get returns channel by given id
	FilterByOption(option *model.ChannelFilterOption) ([]*model.Channel, error) // FilterByOption returns a list of channels with given option
	GetbyOption(option *model.ChannelFilterOption) (*model.Channel, error)      // GetbyOption finds and returns 1 channel filtered using given options
	Upsert(transaction store_iface.SqlxTxExecutor, channel *model.Channel) (*model.Channel, error)
	DeleteChannels(transaction store_iface.SqlxTxExecutor, ids []string) error
}
type ChannelShopStore interface {
	// FilterByOptions(options *model.ChannelShopRelationFilterOptions) ([]*model.ChannelShopRelation, error)
	// Save(relation *model.ChannelShopRelation) (*model.ChannelShopRelation, error)
}

// app
type (
	AppTokenStore interface {
		Save(appToken *model.AppToken) (*model.AppToken, error)
	}
	AppStore interface {
		Save(app *model.App) (*model.App, error)
	}
)

type ClusterDiscoveryStore interface {
	Save(discovery *model.ClusterDiscovery) error
	Delete(discovery *model.ClusterDiscovery) (bool, error)
	Exists(discovery *model.ClusterDiscovery) (bool, error)
	GetAll(discoveryType, clusterName string) ([]*model.ClusterDiscovery, error)
	SetLastPingAt(discovery *model.ClusterDiscovery) error
	Cleanup() error
}

type AuditStore interface {
	ModelFields(prefix string) util.AnyArray[string]
	Save(audit *model.Audit) error
	Get(userID string, offset int, limit int) (model.Audits, error)
	PermanentDeleteByUser(userID string) error
}

type TermsOfServiceStore interface {
	Save(termsOfService *model.TermsOfService) (*model.TermsOfService, error)
	GetLatest(allowFromCache bool) (*model.TermsOfService, error)
	Get(id string, allowFromCache bool) (*model.TermsOfService, error)
}

type PreferenceStore interface {
	Save(preferences model.Preferences) error
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
	SaveOrUpdate(status *model.Status) error
	Get(userID string) (*model.Status, error)
	GetByIds(userIds []string) ([]*model.Status, error)
	ResetAll() error
	GetTotalActiveUsersCount() (int64, error)
	UpdateLastActivityAt(userID string, lastActivityAt int64) error
}

// account stores
type (
	AddressStore interface {
		ModelFields(prefix string) util.AnyArray[string]
		ScanFields(addr *model.Address) []interface{}
		Upsert(transaction store_iface.SqlxTxExecutor, address *model.Address) (*model.Address, error)
		Get(addressID string) (*model.Address, error)                                      // Get returns an Address with given addressID is exist
		DeleteAddresses(transaction store_iface.SqlxTxExecutor, addressIDs []string) error // DeleteAddress deletes given address and returns an error
		FilterByOption(option *model.AddressFilterOption) ([]*model.Address, error)        // FilterByOption finds and returns a list of address(es) filtered by given option
	}
	UserStore interface {
		ClearCaches()
		ModelFields(prefix string) util.AnyArray[string]
		ScanFields(user *model.User) []interface{}
		Save(user *model.User) (*model.User, error)                               // Save takes an user struct and save into database
		Update(user *model.User, allowRoleUpdate bool) (*model.UserUpdate, error) // Update update given user
		UpdateLastPictureUpdate(userID string) error
		ResetLastPictureUpdate(userID string) error
		UpdatePassword(userID, newPassword string) error
		UpdateUpdateAt(userID string) (int64, error)
		UpdateAuthData(userID string, service string, authData *string, email string, resetMfa bool) (string, error)
		ResetAuthDataToEmailForUsers(service string, userIDs []string, includeDeleted bool, dryRun bool) (int, error)
		UpdateMfaSecret(userID, secret string) error
		UpdateMfaActive(userID string, active bool) error
		InvalidateProfileCacheForUser(userID string) // InvalidateProfileCacheForUser
		GetForLogin(loginID string, allowSignInWithUsername, allowSignInWithEmail bool) (*model.User, error)
		VerifyEmail(userID, email string) (string, error) // VerifyEmail set EmailVerified model of user to true
		GetEtagForAllProfiles() string
		GetEtagForProfiles(teamID string) string
		UpdateFailedPasswordAttempts(userID string, attempts int) error
		GetSystemAdminProfiles() (map[string]*model.User, error)
		PermanentDelete(userID string) error // PermanentDelete completely delete user from the system
		AnalyticsGetInactiveUsersCount() (int64, error)
		AnalyticsGetExternalUsers(hostDomain string) (bool, error)
		AnalyticsGetSystemAdminCount() (int64, error)
		AnalyticsGetGuestCount() (int64, error)
		ClearAllCustomRoleAssignments() error
		InferSystemInstallDate() (int64, error)
		GetUsersBatchForIndexing(startTime, endTime int64, limit int) ([]*model.UserForIndexing, error)
		GetKnownUsers(userID string) ([]string, error)
		Count(options model.UserCountOptions) (int64, error)
		AnalyticsActiveCountForPeriod(startTime int64, endTime int64, options model.UserCountOptions) (int64, error)
		GetAllProfiles(options *model.UserGetOptions) ([]*model.User, error)
		Search(term string, options *model.UserSearchOptions) ([]*model.User, error)
		AnalyticsActiveCount(time int64, options model.UserCountOptions) (int64, error)
		GetProfileByIds(ctx context.Context, userIds []string, options *UserGetByIdsOpts, allowFromCache bool) ([]*model.User, error)
		GetUnreadCount(userID string) (int64, error) // TODO: consider me
		FilterByOptions(ctx context.Context, options *model.UserFilterOptions) ([]*model.User, error)
		GetByOptions(ctx context.Context, options *model.UserFilterOptions) (*model.User, error)
	}
	TokenStore interface {
		Save(recovery *model.Token) error
		Delete(token string) error
		GetByToken(token string) (*model.Token, error)
		Cleanup()
		RemoveAllTokensByType(tokenType string) error
		GetAllTokensByType(tokenType string) ([]*model.Token, error)
	}
	UserAccessTokenStore interface {
		Save(token *model.UserAccessToken) (*model.UserAccessToken, error)
		DeleteAllForUser(userID string) error
		Delete(tokenID string) error
		Get(tokenID string) (*model.UserAccessToken, error)
		GetAll(offset int, limit int) ([]*model.UserAccessToken, error)
		GetByToken(tokenString string) (*model.UserAccessToken, error)
		GetByUser(userID string, page, perPage int) ([]*model.UserAccessToken, error)
		Search(term string) ([]*model.UserAccessToken, error)
		UpdateTokenEnable(tokenID string) error
		UpdateTokenDisable(tokenID string) error
	}
	UserAddressStore interface {
		Save(userAddress *model.UserAddress) (*model.UserAddress, error)
		DeleteForUser(userID string, addressID string) error // DeleteForUser delete the relationship between user & address
		// FilterByOptions finds and returns a list of user-address relations with given options
		FilterByOptions(options *model.UserAddressFilterOptions) ([]*model.UserAddress, error)
	}
	CustomerEventStore interface {
		ModelFields(prefix string) util.AnyArray[string]
		Save(customemrEvent *model.CustomerEvent) (*model.CustomerEvent, error)
		Get(id string) (*model.CustomerEvent, error)
		Count() (int64, error)
		FilterByOptions(options *model.CustomerEventFilterOptions) ([]*model.CustomerEvent, error)
	}
	StaffNotificationRecipientStore interface {
		Save(notificationRecipient *model.StaffNotificationRecipient) (*model.StaffNotificationRecipient, error)
		FilterByOptions(options *model.StaffNotificationRecipientFilterOptions) ([]*model.StaffNotificationRecipient, error)
	}
	CustomerNoteStore interface {
		ModelFields(prefix string) util.AnyArray[string]
		Save(note *model.CustomerNote) (*model.CustomerNote, error) // Save insert given customer note into database and returns it
		Get(id string) (*model.CustomerNote, error)                 // Get find customer note with given id and returns it
	}
)

type SystemStore interface {
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
	BulkUpsert(rates []*model.OpenExchangeRate) ([]*model.OpenExchangeRate, error) // BulkUpsert performs bulk update/insert to given exchange rates
	GetAll() ([]*model.OpenExchangeRate, error)                                    // GetAll returns all exchange currency rates
}

type UserGetByIdsOpts struct {
	IsAdmin bool  // IsAdmin tracks whether or not the request is being made by an administrator. Does nothing when provided by a client.
	Since   int64 // Since filters the users based on their UpdateAt timestamp.
	// Restrict to search in a list of teams and channels. Does nothing when provided by a client.
	// ViewRestrictions *model.ViewUsersRestrictions
}
