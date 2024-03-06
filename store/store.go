//go:generate go run layer_generators/main.go

package store

import (
	"context"
	timemodule "time"

	"github.com/mattermost/squirrel"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// Store is database gateway of the system
type Store interface {
	Context() context.Context           // Context gets context
	SetContext(context context.Context) // set context
	Close()                             // Close closes databases
	LockToMaster()                      // LockToMaster constraints all queries to be performed on master
	UnlockFromMaster()                  // UnlockFromMaster makes all datasources available
	ReplicaLagTime() error
	ReplicaLagAbs() error
	CheckIntegrity() <-chan model_helper.IntegrityCheckResult
	DropAllTables()                                // DropAllTables drop all tables in databases
	GetDbVersion(numerical bool) (string, error)   // GetDbVersion returns version in use of database
	FinalizeTransaction(tx boil.ContextTransactor) // FinalizeTransaction tries to rollback given transaction, if an error occur and is not of type sql.ErrTxDone, it prints out the error

	GetMaster() ContextRunner
	GetReplica() boil.ContextExecutor

	// GetQueryBuilder create squirrel sql query builder.
	//
	// NOTE: Don't pass much placeholder format since only the first passed is applied.
	// Ellipsis operator is a trick to support no argument passing.
	//
	// If no placeholder format is passed, defaut to squirrel.Dollar ($)
	GetQueryBuilder(placeholderFormats ...squirrel.PlaceholderFormat) squirrel.StatementBuilderType
	IsUniqueConstraintError(err error, indexNames []string) bool
	MarkSystemRanUnitTests() //
	DBXFromContext(ctx context.Context) boil.ContextExecutor

	User() UserStore                                                   // account
	Address() AddressStore                                             //
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
	VoucherCustomer() VoucherCustomerStore                             //
	Giftcard() GiftCardStore                                           // giftcard
	GiftcardEvent() GiftcardEventStore                                 //
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
	CollectionProduct() CollectionProductStore                         //
	Collection() CollectionStore                                       //
	CollectionChannelListing() CollectionChannelListingStore           //
	CollectionTranslation() CollectionTranslationStore                 //
	ShippingMethodTranslation() ShippingMethodTranslationStore         // shipping
	ShippingMethodChannelListing() ShippingMethodChannelListingStore   //
	ShippingMethodPostalCodeRule() ShippingMethodPostalCodeRuleStore   //
	ShippingMethod() ShippingMethodStore                               //
	ShippingZone() ShippingZoneStore                                   //
	Warehouse() WarehouseStore                                         // warehouse
	Stock() StockStore                                                 //
	Allocation() AllocationStore                                       //
	PreorderAllocation() PreorderAllocationStore                       //
	Wishlist() WishlistStore                                           // wishlist
	WishlistItem() WishlistItemStore                                   //
	PluginConfiguration() PluginConfigurationStore                     // plugin
	Compliance() ComplianceStore                                       // Compliance
	Attribute() AttributeStore                                         // attribute
	AttributeTranslation() AttributeTranslationStore                   //
	AttributeValue() AttributeValueStore                               //
	AttributeValueTranslation() AttributeValueTranslationStore         //
	AssignedPageAttributeValue() AssignedPageAttributeValueStore       //
	AssignedPageAttribute() AssignedPageAttributeStore                 //
	AttributePage() AttributePageStore                                 //
	AssignedProductAttributeValue() AssignedProductAttributeValueStore //
	AssignedProductAttribute() AssignedProductAttributeStore           //
	FileInfo() FileInfoStore                                           // upload session
	UploadSession() UploadSessionStore                                 //
	Plugin() PluginStore                                               //
	ShopTranslation() ShopTranslationStore                             // shop
	ShopStaff() ShopStaffStore                                         //
	Vat() VatStore                                                     //
	OpenExchangeRate() OpenExchangeRateStore                           // external services

	// AssignedVariantAttributeValue() AssignedVariantAttributeValueStore //
	// AssignedVariantAttribute() AssignedVariantAttributeStore           //
	// AttributeVariant() AttributeVariantStore                           //
	// AttributeProduct() AttributeProductStore                           //
}

// shop
type (
	ShopStaffStore interface {
		Upsert(shopStaff model.ShopStaff) (*model.ShopStaff, error)                                // Save inserts given shopStaff into database then returns it with an error
		Get(shopStaffID string) (*model.ShopStaff, error)                                          // Get finds a shop staff with given id then returns it with an error
		FilterByOptions(options model_helper.ShopStaffFilterOptions) (model.ShopStaffSlice, error) // FilterByShopAndStaff finds a relation ship with given shopId and staffId
		GetByOptions(options model_helper.ShopStaffFilterOptions) (*model.ShopStaff, error)
	}
	ShopTranslationStore interface {
		Upsert(translation model.ShopTranslation) (*model.ShopTranslation, error) // Upsert depends on translation's Id then decides to update or insert
		Get(id string) (*model.ShopTranslation, error)                            // Get finds a shop translation with given id then return it with an error
	}
	VatStore interface {
		Upsert(tx boil.ContextTransactor, vats model.VatSlice) (model.VatSlice, error)
		FilterByOptions(options ...qm.QueryMod) (model.VatSlice, error)
	}
)

// Plugin
type PluginStore interface {
	Upsert(keyVal model.PluginKeyValue) (*model.PluginKeyValue, error)
	CompareAndSet(keyVal model.PluginKeyValue, oldValue []byte) (bool, error)
	CompareAndDelete(keyVal model.PluginKeyValue, oldValue []byte) (bool, error)
	SetWithOptions(pluginID string, key string, value []byte, options model_helper.PluginKVSetOptions) (bool, error)
	Get(pluginID, key string) (*model.PluginKeyValue, error)
	Delete(pluginID, key string) error
	DeleteAllForPlugin(PluginID string) error
	DeleteAllExpired() error
	List(pluginID string, page, perPage int) ([]string, error)
}

type UploadSessionStore interface {
	Upsert(session model.UploadSession) (*model.UploadSession, error)
	Get(id string) (*model.UploadSession, error)
	FindAll(options model_helper.UploadSessionFilterOption) (model.UploadSessionSlice, error)
	Delete(id string) error
}

// fileinfo
type FileInfoStore interface {
	Upsert(info model.FileInfo) (*model.FileInfo, error)
	Get(id string, fromMaster bool) (*model.FileInfo, error)
	GetWithOptions(options model_helper.FileInfoFilterOption) (model.FileInfoSlice, error) // Leave perPage and page nil to get all result
	InvalidateFileInfosForPostCache(postID string, deleted bool)
	PermanentDelete(fileID string) error
	PermanentDeleteBatch(endTime int64, limit int64) (int64, error)
	PermanentDeleteByUser(userID string) (int64, error)
	ClearCaches()
	CountAll() (int64, error)

	// Search(paramsList []*model.SearchParams, userID, teamID string, page, perPage int) (*model.FileInfoList, error)
	// GetFilesBatchForIndexing(startTime, endTime int64, limit int) ([]*model.FileForIndexing, error)
	// AttachToPost(fileID string, postID string, creatorID string) error
	// DeleteForPost(postID string) (string, error)
	// GetForPost(postID string, readFromMaster, includeDeleted, allowFromCache bool) ([]*model.FileInfo, error)
}

type (
	AttributeStore interface {
		Delete(tx boil.ContextTransactor, ids []string) (int64, error)
		Upsert(attr model.Attribute) (*model.Attribute, error)                                  // Upsert inserts or updates given model then returns it
		FilterbyOption(option model_helper.AttributeFilterOption) (model.AttributeSlice, error) // FilterbyOption returns a list of attributes by given option
		GetProductTypeAttributes(productTypeID string, unassigned bool, filter model_helper.AttributeFilterOption) (model.AttributeSlice, error)
		GetPageTypeAttributes(pageTypeID string, unassigned bool) (model.AttributeSlice, error)
		CountByOptions(options model_helper.AttributeFilterOption) (int64, error)
	}
	AttributeTranslationStore interface {
	}
	AttributeValueStore interface {
		Count(options model_helper.AttributeValueFilterOptions) (int64, error)
		Delete(tx boil.ContextTransactor, ids []string) (int64, error)
		Upsert(tx boil.ContextTransactor, values model.AttributeValueSlice) (model.AttributeValueSlice, error)
		Get(attributeID string) (*model.AttributeValue, error)                                               // Get finds an model value with given id then returns it with an error
		FilterByOptions(options model_helper.AttributeValueFilterOptions) (model.AttributeValueSlice, error) // FilterByOptions finds and returns all matched model values based on given options
	}
	AttributeValueTranslationStore interface {
	}
	AssignedPageAttributeValueStore interface {
		Upsert(tx boil.ContextTransactor, assignedPageAttrValue model.AssignedPageAttributeValueSlice) (model.AssignedPageAttributeValueSlice, error) // Save insert given value into database then returns it with an error
		Get(assignedPageAttrValueID string) (*model.AssignedPageAttributeValue, error)                                                                // Get try finding an value with given id then returns it with an error
		SelectForSort(assignmentID string) (model.AssignedPageAttributeValueSlice, model.AttributeValueSlice, error)                                  // SelectForSort uses inner join to find two list: []*assignedPageAttributeValue and []*attributeValue. With given assignedPageAttributeID
	}
	AssignedPageAttributeStore interface {
		Upsert(assignedPageAttr model.AssignedPageAttribute) (*model.AssignedPageAttribute, error)                        // Save inserts given assigned page model into database and returns it with an error
		Get(id string) (*model.AssignedPageAttribute, error)                                                              // Get returns an assigned page model with an error
		FilterByOptions(options model_helper.AssignedPageAttributeFilterOption) (model.AssignedPageAttributeSlice, error) // GetByOption try to find an assigned page model with given option. If nothing found, creats new instance with that option and returns such value with an error
	}
	AttributePageStore interface {
		Save(page model.AttributePage) (*model.AttributePage, error)
		Get(pageID string) (*model.AttributePage, error)
		GetByOption(option model_helper.AttributePageFilterOption) (*model.AttributePage, error)
	}
	// AssignedVariantAttributeValueStore interface {
	// 	Save(assignedVariantAttrValue model.AssignedVariantAttributeValue) (*model.AssignedVariantAttributeValue, error)                                                 // Save inserts new value into database then returns it with an error
	// 	Get(id string) (*model.AssignedVariantAttributeValue, error)                                                                                                     // Get try finding a value with given id then returns it with an error
	// 	SaveInBulk(assignmentID string, attributeValueIDs []string) (model.AssignedVariantAttributeValueSlice, error)                                                    // SaveInBulk save multiple values into database then returns them
	// 	SelectForSort(assignmentID string) (assignedVariantAttributeValues model.AssignedVariantAttributeValueSlice, attributeValues []*model.AttributeValue, err error) // SelectForSort
	// 	UpdateInBulk(attributeValues model.AssignedVariantAttributeValueSlice) error                                                                                     // UpdateInBulk use transaction to update given values, then returns an error to indicate if the operation was successful or not
	// 	FilterByOptions(options model.AssignedVariantAttributeValueFilterOptions) (model.AssignedVariantAttributeValueSlice, error)
	// }
	// AssignedVariantAttributeStore interface {
	// Save(assignedVariantAttribute model.AssignedVariantAttribute) (*model.AssignedVariantAttribute, error)         // Save insert new instance into database then returns it with an error
	// Get(id string) (*model.AssignedVariantAttribute, error)                                                        // Get find assigned variant model from database then returns it with an error
	// GetWithOption(option model.AssignedVariantAttributeFilterOption) (*model.AssignedVariantAttribute, error)      // GetWithOption try finding an assigned variant model with given option. If nothing found, it creates instance with given option. Finally it returns expected value with an error
	// FilterByOption(option model.AssignedVariantAttributeFilterOption) (model.AssignedVariantAttributeSlice, error) // FilterByOption finds and returns a list of assigned variant attributes filtered by given options
	// }
	// AttributeVariantStore interface {
	// 	Save(attributeVariant model.AttributeVariant) (*model.AttributeVariant, error)
	// 	Get(attributeVariantID string) (*model.AttributeVariant, error)
	// 	GetByOption(option model.AttributeVariantFilterOption) (*model.AttributeVariant, error) // GetByOption finds 1 model variant with given option.
	// 	FilterByOptions(options model.AttributeVariantFilterOption) ([]*model.AttributeVariant, error)
	// }
	AssignedProductAttributeValueStore interface {
		Save(assignedProductAttrValue model.AssignedProductAttributeValue) (*model.AssignedProductAttributeValue, error) // Save inserts given instance into database then returns it with an error
		Get(assignedProductAttrValueID string) (*model.AssignedProductAttributeValue, error)                             // Get try finding an instance with given id then returns the value with an error
		SelectForSort(assignmentID string) (model.AssignedProductAttributeValueSlice, model.AttributeValueSlice, error)  // SelectForSort finds all `*AssignedProductAttributeValue` and related `*AttributeValues` with given `assignmentID`, then returns them with an error.
		FilterByOptions(options model_helper.AssignedProductAttributeValueFilterOptions) (model.AssignedProductAttributeValueSlice, error)
	}
	AssignedProductAttributeStore interface {
		// Save(assignedProductAttribute model.AssignedProductAttribute) (*model.AssignedProductAttribute, error)           // Save inserts new assgignedProductAttribute into database and returns it with an error
		// Get(id string) (*model.AssignedProductAttribute, error)                                                          // Get finds and returns an assignedProductAttribute with en error
		GetWithOption(option model_helper.AssignedProductAttributeFilterOption) (*model.AssignedProductAttribute, error) // GetWithOption try finding an `AssignedProductAttribute` with given `option`. If nothing found, it creates new instance then returns it with an error
		FilterByOptions(options model_helper.AssignedProductAttributeFilterOption) (model.AssignedProductAttributeSlice, error)
	}
	// AttributeProductStore interface {
	// 	Save(attributeProduct model.AttributeProduct) (*model.AttributeProduct, error)                // Save inserts given model product relationship into database then returns it and an error
	// 	Get(attributeProductID string) (*model.AttributeProduct, error)                               // Get finds an attributeProduct relationship and returns it with an error
	// 	GetByOption(option model.AttributeProductFilterOption) (*model.AttributeProduct, error)       // GetByOption returns an attributeProduct with given condition
	// 	FilterByOptions(option model.AttributeProductFilterOption) ([]*model.AttributeProduct, error) // FilterByOptions returns attributeProducts with given condition
	// }
	CustomProductAttributeStore interface {
		Upsert(tx boil.ContextTransactor, record model.CustomProductAttribute) (*model.CustomProductAttribute, error)
		Delete(tx boil.ContextTransactor, ids []string) (int64, error)
		FilterByOptions(options model_helper.CustomProductAttributeFilterOptions) (model.CustomProductAttributeSlice, error)
	}
)

// model
type ComplianceStore interface {
	Upsert(model model.Compliance) (*model.Compliance, error)
	Get(id string) (*model.Compliance, error)
	GetAll(offset, limit int) (model.ComplianceSlice, error)
	ComplianceExport(model model.Compliance, cursor model_helper.ComplianceExportCursor, limit int) ([]*model_helper.CompliancePost, model_helper.ComplianceExportCursor, error)
	MessageExport(cursor model_helper.MessageExportCursor, limit int) ([]*model_helper.MessageExport, model_helper.MessageExportCursor, error)
}

// plugin
type PluginConfigurationStore interface {
	Upsert(config model.PluginConfiguration) (*model.PluginConfiguration, error)                                              // Upsert inserts or updates given plugin configuration and returns it
	Get(id string) (*model.PluginConfiguration, error)                                                                        // Get finds a plugin configuration with given id then returns it
	FilterPluginConfigurations(options model_helper.PluginConfigurationFilterOptions) (model.PluginConfigurationSlice, error) // FilterPluginConfigurations finds and returns a list of configs with given options then returns them
}

// model
type (
	WishlistStore interface {
		Upsert(wishList model.Wishlist) (*model.Wishlist, error)                       // Upsert inserts or update given model and returns it
		GetByOption(option model_helper.WishlistFilterOption) (*model.Wishlist, error) // GetByOption finds and returns a slice of wishlists by given option
	}
	WishlistItemStore interface {
		GetById(id string) (*model.WishlistItem, error)                                                               // GetById returns a model item wish given id
		BulkUpsert(tx boil.ContextTransactor, wishlistItems model.WishlistItemSlice) (model.WishlistItemSlice, error) // Upsert inserts or updates given model item then returns it
		FilterByOption(option model_helper.WishlistItemFilterOption) (model.WishlistItemSlice, error)                 // FilterByOption finds and returns a slice of model items filtered using given options
		GetByOption(option model_helper.WishlistItemFilterOption) (*model.WishlistItem, error)                        // GetByOption finds and returns a model item filtered by given option
		Delete(tx boil.ContextTransactor, ids []string) (int64, error)                                                // DeleteItemsByOption finds and deletes model items that satisfy given filtering options and returns number of items deleted
	}
)

type (
	WarehouseStore interface {
		WarehouseShipingZonesByCountryCodeAndChannelID(countryCode, channelID string) (model.WarehouseShippingZoneSlice, error)
		Delete(tx boil.ContextTransactor, ids []string) error
		Upsert(model model.Warehouse) (*model.Warehouse, error)                                 // Save inserts given model into database then returns it.
		FilterByOprion(option model_helper.WarehouseFilterOption) (model.WarehouseSlice, error) // FilterByOprion returns a slice of warehouses with given option
		GetByOption(option model_helper.WarehouseFilterOption) (*model.Warehouse, error)        // GetByOption finds and returns a model filtered given option
		WarehouseByStockID(stockID string) (*model.Warehouse, error)                            // WarehouseByStockID returns 1 model by given stock id
		ApplicableForClickAndCollectNoQuantityCheck(checkoutLines model.CheckoutLineSlice, country model.CountryCode) (model.WarehouseSlice, error)
		ApplicableForClickAndCollectCheckoutLines(checkoutLines model.CheckoutLineSlice, country model.CountryCode) (model.WarehouseSlice, error)
		ApplicableForClickAndCollectOrderLines(orderLines model.OrderLineSlice, country model.CountryCode) (model.WarehouseSlice, error)
	}
	StockStore interface {
		Delete(tx boil.ContextTransactor, ids []string) (int64, error)
		Get(stockID string) (*model.Stock, error)                                                                                      // Get finds and returns stock with given stockID. Returned error could be either (nil, *ErrNotFound, error)
		FilterForCountryAndChannel(options model_helper.StockFilterOptionsForCountryAndChannel) (model.StockSlice, error)              // FilterForCountryAndChannel finds and returns stocks with given options
		FilterVariantStocksForCountry(options model_helper.StockFilterOptionsForCountryAndChannel) (model.StockSlice, error)           // FilterVariantStocksForCountry finds and returns stocks with given options
		FilterProductStocksForCountryAndChannel(options model_helper.StockFilterOptionsForCountryAndChannel) (model.StockSlice, error) // FilterProductStocksForCountryAndChannel finds and returns stocks with given options
		ChangeQuantity(stockID string, quantity int) error                                                                             // ChangeQuantity reduce or increase the quantity of given stock
		FilterByOption(options model_helper.StockFilterOption) (int64, model.StockSlice, error)                                        // FilterByOption finds and returns a slice of stocks that satisfy given option
		BulkUpsert(tx boil.ContextTransactor, stocks model.StockSlice) (model.StockSlice, error)                                       // BulkUpsert performs upserts or inserts given stocks, then returns them
		FilterForChannel(options model_helper.StockFilterForChannelOption) (squirrel.Sqlizer, model.StockSlice, error)                 // FilterForChannel finds and returns stocks that satisfy given options
	}
	AllocationStore interface {
		BulkUpsert(tx boil.ContextTransactor, allocations model.AllocationSlice) (model.AllocationSlice, error) // BulkUpsert performs update, insert given allocations then returns them afterward
		Get(id string) (*model.Allocation, error)                                                               // Get find and returns allocation with given id
		FilterByOption(option model_helper.AllocationFilterOption) (model.AllocationSlice, error)               // FilterbyOption finds and returns a list of allocations based on given option
		Delete(tx boil.ContextTransactor, ids []string) error                                                   // BulkDelete perform bulk deletes given allocations.
		CountAvailableQuantityForStock(stock model.Stock) (int, error)                                          // CountAvailableQuantityForStock counts and returns available quantity of given stock
	}
	PreorderAllocationStore interface {
		BulkCreate(tx boil.ContextTransactor, preorderAllocations model.PreorderAllocationSlice) (model.PreorderAllocationSlice, error) // BulkCreate bulk inserts given preorderAllocations and returns them
		FilterByOption(options model_helper.PreorderAllocationFilterOption) (model.PreorderAllocationSlice, error)                      // FilterByOption finds and returns a list of preorder allocations filtered using given options
		Delete(tx boil.ContextTransactor, ids []string) error                                                                           // Delete deletes preorder-allocations by given ids
	}
)

type (
	ShippingZoneStore interface {
		Delete(tx boil.ContextTransactor, ids []string) (int64, error)
		Upsert(tx boil.ContextTransactor, shippingZone model.ShippingZone) (*model.ShippingZone, error) // Upsert depends on given model zone's Id to decide update or insert the zone
		Get(id string) (*model.ShippingZone, error)                                                     // Get finds 1 model zone for given shippingZoneID
		FilterByOption(option model_helper.ShippingZoneFilterOption) (model.ShippingZoneSlice, error)   // FilterByOption finds a list of model zones based on given option
		CountByOptions(options model_helper.ShippingZoneFilterOption) (int64, error)
		// ToggleRelations(tx boil.ContextTransactor, zones model.ShippingZoneSlice, warehouseIds, channelIds []string, delete bool) error // NOTE: relations must be []*Channel or []*Warehouse
	}
	ShippingMethodStore interface {
		Upsert(tx boil.ContextTransactor, method model.ShippingMethod) (*model.ShippingMethod, error)                                                                                         // Upsert bases on given method's Id to decide update or insert it
		Get(id string) (*model.ShippingMethod, error)                                                                                                                                         // Get finds and returns a model method with given id
		ApplicableShippingMethods(price *goprices.Money, channelID string, weight *measurement.Weight, countryCode model.CountryCode, productIDs []string) (model.ShippingMethodSlice, error) // ApplicableShippingMethods finds all model methods with given conditions
		GetbyOption(options model_helper.ShippingMethodFilterOption) (*model.ShippingMethod, error)                                                                                           // GetbyOption finds and returns a model method that satisfy given options
		FilterByOptions(options model_helper.ShippingMethodFilterOption) (model.ShippingMethodSlice, error)
		Delete(tx boil.ContextTransactor, ids []string) error
	}
	ShippingMethodPostalCodeRuleStore interface {
		Delete(tx boil.ContextTransactor, ids []string) error
		Save(tx boil.ContextTransactor, rules model.ShippingMethodPostalCodeRuleSlice) (model.ShippingMethodPostalCodeRuleSlice, error)
		FilterByOptions(options model_helper.ShippingMethodPostalCodeRuleFilterOptions) (model.ShippingMethodPostalCodeRuleSlice, error)
	}
	ShippingMethodChannelListingStore interface {
		Delete(tx boil.ContextTransactor, ids []string) error
		Upsert(tx boil.ContextTransactor, listings model.ShippingMethodChannelListingSlice) (model.ShippingMethodChannelListingSlice, error) // Upsert depends on given listing's Id to decide whether to save or update the listing
		Get(id string) (*model.ShippingMethodChannelListing, error)                                                                          // Get finds a model method channel listing with given listingID
		FilterByOption(option model_helper.ShippingMethodChannelListingFilterOption) (model.ShippingMethodChannelListingSlice, error)        // FilterByOption returns a list of model method channel listings based on given option. result sorted by creation time ASC
	}
	ShippingMethodTranslationStore interface {
	}
)

// product
type (
	CollectionTranslationStore interface {
	}
	CollectionChannelListingStore interface {
		Delete(tx boil.ContextTransactor, ids []string) error
		Upsert(tx boil.ContextTransactor, relations model.CollectionChannelListingSlice) (model.CollectionChannelListingSlice, error)
		FilterByOptions(options model_helper.CollectionChannelListingFilterOptions) (model.CollectionChannelListingSlice, error)
	}
	CollectionStore interface {
		Upsert(collection model.Collection) (*model.Collection, error)                                          // Upsert depends on given collection's Id property to decide update or insert the collection
		Get(collectionID string) (*model.Collection, error)                                                     // Get finds and returns collection with given collectionID
		FilterByOption(option model_helper.CollectionFilterOptions) (model_helper.CustomCollectionSlice, error) // FilterByOption finds and returns a list of collections satisfy the given option
		Delete(tx boil.ContextTransactor, ids []string) error
	}
	CollectionProductStore interface {
		// Delete(tx boil.ContextTransactor, options *model.CollectionProductFilterOptions) error
		// BulkSave(tx boil.ContextTransactor, relations []*model.CollectionProduct) ([]*model.CollectionProduct, error)
		// FilterByOptions(options *model.CollectionProductFilterOptions) ([]*model.CollectionProduct, error)
	}
	ProductMediaStore interface {
		Upsert(tx boil.ContextTransactor, medias model.ProductMediumSlice) (model.ProductMediumSlice, error) // Upsert depends on given media's Id property to decide insert or update it
		Get(id string) (*model.ProductMedium, error)                                                         // Get finds and returns 1 product media with given id
		FilterByOption(option model_helper.ProductMediaFilterOption) (model.ProductMediumSlice, error)       // FilterByOption finds and returns a list of product medias with given id
		Delete(tx boil.ContextTransactor, ids []string) (int64, error)
	}
	DigitalContentUrlStore interface {
		Upsert(contentURL model.DigitalContentURL) (*model.DigitalContentURL, error) // Upsert inserts or updates given digital content url into database then returns it
		Get(id string) (*model.DigitalContentURL, error)                             // Get finds and returns a digital content url with given id
		FilterByOptions(options model_helper.DigitalContentUrlFilterOptions) (model.DigitalContentURLSlice, error)
	}
	DigitalContentStore interface {
		Delete(tx boil.ContextTransactor, ids []string) error
		Save(content model.DigitalContent) (*model.DigitalContent, error)
		GetByOption(option model_helper.DigitalContentFilterOption) (*model.DigitalContent, error)
		FilterByOption(option model_helper.DigitalContentFilterOption) (model.DigitalContentSlice, error)
	}
	ProductVariantChannelListingStore interface {
		Get(variantChannelListingID string) (*model.ProductVariantChannelListing, error)                                                                   // Get finds and returns 1 product variant channel listing based on given variantChannelListingID
		FilterbyOption(option model.ProductVariantChannelListingFilterOption) ([]*model.ProductVariantChannelListing, error)                               // FilterbyOption finds and returns all product variant channel listings filterd using given option
		Upsert(tx boil.ContextTransactor, variantChannelListings model.ProductVariantChannelListingSlice) (model.ProductVariantChannelListingSlice, error) // BulkUpsert performs bulk upsert given product variant channel listings then returns them
	}
	ProductVariantTranslationStore interface {
		// Upsert(translation *model.ProductVariantTranslation) (*model.ProductVariantTranslation, error)                  // Upsert inserts or updates given translation then returns it
		// Get(translationID string) (*model.ProductVariantTranslation, error)                                             // Get finds and returns 1 product variant translation with given id
		// FilterByOption(option *model.ProductVariantTranslationFilterOption) ([]*model.ProductVariantTranslation, error) // FilterByOption finds and returns product variant translations filtered using given options
	}
	ProductVariantStore interface {
		Delete(tx boil.ContextTransactor, ids []string) (int64, error)
		FindVariantsAvailableForPurchase(variantIds []string, channelID string) (model.ProductVariantSlice, error)
		Upsert(tx boil.ContextTransactor, variant model.ProductVariant) (*model.ProductVariant, error) // Save inserts product variant instance to database
		Get(id string) (*model.ProductVariant, error)                                                  // Get returns a product variant with given id
		GetWeight(productVariantID string) (*measurement.Weight, error)                                // GetWeight returns weight of given product variant
		GetByOrderLineID(orderLineID string) (*model.ProductVariant, error)                            // GetByOrderLineID finds and returns a product variant by given orderLineID
		FilterByOption(option *model.ProductVariantFilterOption) (model.ProductVariantSlice, error)    // FilterByOption finds and returns product variants based on given option
	}
	ProductChannelListingStore interface {
		Upsert(tx boil.ContextTransactor, listings model.ProductChannelListingSlice) (model.ProductChannelListingSlice, error) // BulkUpsert performs bulk upsert on given product channel listings
		Get(channelListingID string) (*model.ProductChannelListing, error)                                                     // Get try finding a product channel listing, then returns it with an error
		FilterByOption(option model_helper.ProductChannelListingFilterOption) (model.ProductChannelListingSlice, error)        // FilterByOption filter a list of product channel listings by given option. Then returns them with an error
	}
	ProductTranslationStore interface {
		// Upsert(translation *model.ProductTranslation) (*model.ProductTranslation, error)                  // Upsert inserts or update given translation
		// Get(translationID string) (*model.ProductTranslation, error)                                      // Get finds and returns a product translation by given id
		// FilterByOption(option *model.ProductTranslationFilterOption) ([]*model.ProductTranslation, error) // FilterByOption finds and returns product translations filtered using given options
	}
	ProductTypeStore interface {
		// ToggleProductTypeRelations(tx boil.ContextTransactor, productTypeID string, productAttributes, variantAttributes model.AttributeSlice, isDelete bool) error
		// Delete(tx boil.ContextTransactor, ids []string) (int64, error)
		// FilterbyOption(options *model.ProductTypeFilterOption) (int64, []*model.ProductType, error)
		// Save(tx boil.ContextTransactor, productType *model.ProductType) (*model.ProductType, error) // Save try inserting new product type into database then returns it
		// FilterProductTypesByCheckoutToken(checkoutToken string) ([]*model.ProductType, error)       // FilterProductTypesByCheckoutToken is used to check if a model requires model
		// ProductTypesByProductIDs(productIDs []string) ([]*model.ProductType, error)                 // ProductTypesByProductIDs returns all product types belong to given products
		// ProductTypeByProductVariantID(variantID string) (*model.ProductType, error)                 // ProductTypeByProductVariantID finds and returns 1 product type that is related to given product variant
		// GetByOption(options *model.ProductTypeFilterOption) (*model.ProductType, error)             // GetByOption finds and returns a product type with given options
		// Count(options *model.ProductTypeFilterOption) (int64, error)
	}
	CategoryTranslationStore interface{}
	CategoryStore            interface {
		Upsert(category model.Category) (*model.Category, error)                                  // Upsert depends on given category's Id field to decide update or insert it
		Get(ctx context.Context, categoryID string, allowFromCache bool) (*model.Category, error) // Get finds and returns a category with given id
		FilterByOption(option model_helper.CategoryFilterOption) (model.CategorySlice, error)     // FilterByOption finds and returns a list of categories satisfy given option
	}
	ProductStore interface {
		Save(tx boil.ContextTransactor, product model.Product) (*model.Product, error)
		FilterByOption(option model_helper.ProductFilterOption) (model.ProductSlice, error)                                                                             // FilterByOption finds and returns all products that satisfy given option
		PublishedProducts(channelSlug string) (model.ProductSlice, error)                                                                                               // FilterPublishedProducts finds and returns products that belong to given channel slug and are published
		NotPublishedProducts(channelID string) (model.ProductSlice, error)                                                                                              // NotPublishedProducts finds all not published products belong to given channel
		PublishedWithVariants(channelIdOrSlug string) squirrel.SelectBuilder                                                                                            // PublishedWithVariants finds and returns products.
		VisibleToUserProductsQuery(channelIdOrSlug string, userHasOneOfProductpermissions bool) squirrel.SelectBuilder                                                  // FilterVisibleToUserProduct finds and returns all products that are visible to requesting user.
		SelectForUpdateDiscountedPricesOfCatalogues(tx boil.ContextTransactor, productIDs, categoryIDs, collectionIDs, variantIDs []string) (model.ProductSlice, error) // SelectForUpdateDiscountedPricesOfCatalogues finds and returns product based on given ids lists.
		AdvancedFilterQueryBuilder(input model_helper.ExportProductsFilterOptions) squirrel.SelectBuilder                                                               // AdvancedFilterQueryBuilder advancedly finds products, filtered using given options
		FilterByQuery(query squirrel.SelectBuilder) (model.ProductSlice, error)                                                                                         // FilterByQuery finds and returns products with given query, limit, createdAtGt
		CountByCategoryIDs(categoryIDs []string) ([]*model_helper.ProductCountByCategoryID, error)
	}
)

// model
type (
	PaymentStore interface {
		Upsert(tx boil.ContextTransactor, model model.Payment) (*model.Payment, error)                                    // Save save model instance into database
		CancelActivePaymentsOfCheckout(checkoutToken string) error                                                        // CancelActivePaymentsOfCheckout inactivate all payments that belong to given model and in active status
		FilterByOption(option model_helper.PaymentFilterOptions) (model.PaymentSlice, error)                              // FilterByOption finds and returns a list of payments that satisfy given option
		UpdatePaymentsOfCheckout(tx boil.ContextTransactor, checkoutToken string, option model_helper.PaymentPatch) error // UpdatePaymentsOfCheckout updates payments of given model
		PaymentOwnedByUser(userID, paymentID string) (bool, error)
	}
	PaymentTransactionStore interface {
		Upsert(tx boil.ContextTransactor, paymentTransaction model.PaymentTransaction) (*model.PaymentTransaction, error) // Save inserts new model transaction into database
		Get(id string) (*model.PaymentTransaction, error)                                                                 // Get returns a model transaction with given id
		FilterByOption(option model_helper.PaymentTransactionFilterOpts) ([]*model.PaymentTransaction, error)             // FilterByOption finds and returns a list of transactions with given option
	}
)

// page
type (
	PageTypeStore interface {
	}
	PageTranslationStore interface {
	}
	PageStore interface {
		FilterByOptions(options model_helper.PageFilterOptions) (model.PageSlice, error)
	}
)

// order
type (
	OrderLineStore interface {
		Upsert(tx boil.ContextTransactor, orderLine model.OrderLine) (*model.OrderLine, error)   // Upsert depends on given orderLine's Id to decide to update or save it
		Get(id string) (*model.OrderLine, error)                                                 // Get returns a order line with id of given id
		FilterbyOption(option model_helper.OrderLineFilterOptions) (model.OrderLineSlice, error) // FilterbyOption finds and returns order lines by given option
	}
	OrderStore interface {
		Delete(tx boil.ContextTransactor, ids []string) (int64, error)
		Get(id string) (*model.Order, error)                                                         // Get find order in database with given id
		FilterByOption(option model_helper.OrderFilterOption) (model_helper.CustomOrderSlice, error) // FilterByOption returns a list of orders, filtered by given option
		BulkUpsert(tx boil.ContextTransactor, orders model.OrderSlice) (model.OrderSlice, error)
	}
	OrderEventStore interface {
		// Save(tx boil.ContextTransactor, orderEvent *model.OrderEvent) (*model.OrderEvent, error) // Save inserts given order event into database then returns it
		// Get(orderEventID string) (*model.OrderEvent, error)                                      // Get finds order event with given id then returns it
		// FilterByOptions(options *model.OrderEventFilterOptions) ([]*model.OrderEvent, error)
	}
	FulfillmentLineStore interface {
		Upsert(fulfillmentLine model.FulfillmentLine) (*model.FulfillmentLine, error)
		Get(id string) (*model.FulfillmentLine, error)
		Delete(tx boil.ContextTransactor, ids []string) error
		FilterByOptions(option model_helper.FulfillmentLineFilterOption) (model.FulfillmentLineSlice, error)
	}
	FulfillmentStore interface {
		Upsert(tx boil.ContextTransactor, fulfillment model.Fulfillment) (*model.Fulfillment, error) // Upsert depends on given fulfillment's Id to decide update or insert it
		Get(id string) (*model.Fulfillment, error)                                                   // Get finds and return a fulfillment by given id
		FilterByOption(option model_helper.FulfillmentFilterOption) (model.FulfillmentSlice, error)  // FilterByOption finds and returns a slice of fulfillments by given option
		Delete(tx boil.ContextTransactor, ids []string) error                                        // BulkDeleteFulfillments deletes given fulfillments
	}
)

// menu
type (
	MenuItemTranslationStore interface {
	}
	MenuStore interface {
		Delete(tx boil.ContextTransactor, ids []string) (int64, error)
		Upsert(menu model.Menu) (*model.Menu, error)
		GetByOptions(options model_helper.MenuFilterOptions) (*model.Menu, error)
		FilterByOptions(options model_helper.MenuFilterOptions) (model.MenuSlice, error)
	}
	MenuItemStore interface {
		Delete(ids []string) (int64, error)
		Upsert(menuItem model.MenuItem) (*model.MenuItem, error) // Save insert given menu item into database and returns it
		GetByOptions(options model_helper.MenuItemFilterOptions) (*model.MenuItem, error)
		FilterByOptions(options model_helper.MenuItemFilterOptions) (model.MenuItemSlice, error)
	}
)

// invoice
type (
	InvoiceEventStore interface {
		Upsert(invoiceEvent model.InvoiceEvent) (*model.InvoiceEvent, error) // Upsert depends on given invoice event's Id to update/insert it
		Get(id string) (*model.InvoiceEvent, error)                          // Get finds and returns 1 invoice event
	}
	InvoiceStore interface {
		Upsert(invoice model.Invoice) (*model.Invoice, error)                          // Upsert depends on given invoice Id to update/insert it
		GetbyOptions(options model_helper.InvoiceFilterOption) (*model.Invoice, error) // Get finds and returns 1 invoice
		FilterByOptions(options model_helper.InvoiceFilterOption) (model.InvoiceSlice, error)
		Delete(tx boil.ContextTransactor, ids []string) error
	}
)

// giftcard related stores
type (
	GiftCardStore interface {
		Delete(tx boil.ContextTransactor, ids []string) error
		BulkUpsert(tx boil.ContextTransactor, giftCards model.GiftcardSlice) (model.GiftcardSlice, error) // BulkUpsert depends on given giftcards's Id properties then perform according operation
		GetById(id string) (*model.Giftcard, error)                                                       // GetById returns a giftcard instance that has id of given id
		FilterByOption(option model_helper.GiftcardFilterOption) (model.GiftcardSlice, error)             // FilterByOption finds giftcards wth option
		DeactivateOrderGiftcards(tx boil.ContextTransactor, orderID string) ([]string, error)
	}
	GiftcardEventStore interface {
		Get(id string) (*model.GiftcardEvent, error)                                                         // Get finds and returns a giftcard event found by given id
		Upsert(tx boil.ContextTransactor, events model.GiftcardEventSlice) (model.GiftcardEventSlice, error) // BulkUpsert upserts and returns given giftcard events
		FilterByOptions(options model_helper.GiftCardEventFilterOption) (model.GiftcardEventSlice, error)    // FilterByOptions finds and returns a list of giftcard events with given options
	}
)

// discount
type (
	OrderDiscountStore interface {
		Upsert(tx boil.ContextTransactor, orderDiscount model.OrderDiscount) (*model.OrderDiscount, error) // Upsert depends on given order discount's Id property to decide to update/insert it
		Get(orderDiscountID string) (*model.OrderDiscount, error)                                          // Get finds and returns an order discount with given id
		FilterbyOption(option model_helper.OrderDiscountFilterOption) (model.OrderDiscountSlice, error)    // FilterbyOption filters order discounts that satisfy given option, then returns them
		BulkDelete(ids []string) error                                                                     // BulkDelete perform bulk delete all given order discount ids
	}
	DiscountSaleTranslationStore interface {
	}
	DiscountSaleChannelListingStore interface {
		Delete(tx boil.ContextTransactor, ids []string) error
		Upsert(tx boil.ContextTransactor, listings model.SaleChannelListingSlice) (model.SaleChannelListingSlice, error)
		Get(id string) (*model.SaleChannelListing, error) // Get finds and returns sale channel listing with given id
		FilterByOptions(option model_helper.SaleChannelListingFilterOption) (model.SaleChannelListingSlice, error)
	}
	VoucherTranslationStore interface {
		Upsert(translation model.VoucherTranslation) (*model.VoucherTranslation, error)                           // Save inserts given translation into database and returns it
		Get(id string) (*model.VoucherTranslation, error)                                                         // Get finds and returns a voucher translation with given id
		FilterByOption(option model_helper.VoucherTranslationFilterOption) (model.VoucherTranslationSlice, error) // FilterByOption returns a list of voucher translations filtered using given options
		GetByOption(option model_helper.VoucherTranslationFilterOption) (*model.VoucherTranslation, error)        // GetByOption finds and returns 1 voucher translation by given options
	}
	DiscountSaleStore interface {
		Delete(tx boil.ContextTransactor, ids []string) (int64, error)
		Upsert(tx boil.ContextTransactor, sale model.Sale) (*model.Sale, error)                         // Upsert bases on sale's Id to decide to update or insert given sale
		Get(saleID string) (*model.Sale, error)                                                         // Get finds and returns a sale with given saleID
		FilterSalesByOption(option model_helper.SaleFilterOption) (model_helper.CustomSaleSlice, error) // FilterSalesByOption filter sales by option
	}
	VoucherChannelListingStore interface {
		Upsert(tx boil.ContextTransactor, voucherChannelListings model.VoucherChannelListingSlice) (model.VoucherChannelListingSlice, error) // upsert check given listing's Id to decide whether to create or update it. Then returns a listing with an error
		Get(voucherChannelListingID string) (*model.VoucherChannelListing, error)                                                            // Get finds a listing with given id, then returns it with an error
		FilterbyOption(option model_helper.VoucherChannelListingFilterOption) (model.VoucherChannelListingSlice, error)                      // FilterbyOption finds and returns a list of voucher channel listing relationship instances filtered by given option
		Delete(tx boil.ContextTransactor, ids []string) error
	}
	DiscountVoucherStore interface {
		Upsert(voucher model.Voucher) (*model.Voucher, error)                                                    // Upsert saves or updates given voucher then returns it with an error
		Get(id string) (*model.Voucher, error)                                                                   // Get finds a voucher with given id, then returns it with an error
		FilterVouchersByOption(option model_helper.VoucherFilterOption) (model_helper.CustomVoucherSlice, error) // FilterVouchersByOption finds vouchers bases on given option.
		ExpiredVouchers(date timemodule.Time) (model.VoucherSlice, error)                                        // ExpiredVouchers finds and returns vouchers that are expired before given date
		Delete(tx boil.ContextTransactor, ids []string) (int64, error)
	}
	VoucherCustomerStore interface {
		Save(voucherCustomer model.VoucherCustomer) (*model.VoucherCustomer, error)                           // Save inserts given voucher customer instance into database ands returns it
		Delete(ids []string) error                                                                            // DeleteInBulk deletes given voucher-customers with given id
		GetByOption(options model_helper.VoucherCustomerFilterOption) (*model.VoucherCustomer, error)         // GetByOption finds and returns a voucher customer with given options
		FilterByOptions(options model_helper.VoucherCustomerFilterOption) (model.VoucherCustomerSlice, error) // FilterByOptions finds and returns a slice of voucher customers by given options
	}
)

// csv
type (
	CsvExportEventStore interface {
		Save(event model.ExportEvent) (*model.ExportEvent, error)                                  // Save inserts given export event into database then returns it
		FilterByOption(options model_helper.ExportEventFilterOption) ([]*model.ExportEvent, error) // FilterByOption finds and returns a list of export events filtered using given option
	}
	CsvExportFileStore interface {
		Save(file model.ExportFile) (*model.ExportFile, error) // Save inserts given export file into database then returns it
		Get(id string) (*model.ExportFile, error)              // Get finds and returns an export file found using given id
	}
)

// model
type (
	CheckoutLineStore interface {
		Upsert(checkoutLines model.CheckoutLineSlice) (model.CheckoutLineSlice, error)                        // Upsert checks whether to update or insert given model line then performs according operation
		Get(id string) (*model.CheckoutLine, error)                                                           // Get returns a model line with given id
		DeleteLines(tx boil.ContextTransactor, checkoutLineIDs []string) error                                // DeleteLines deletes all model lines with given uuids
		TotalWeightForCheckoutLines(checkoutLineIDs []string) (*measurement.Weight, error)                    // TotalWeightForCheckoutLines calculate total weight for given model lines
		CheckoutLinesByOption(option model_helper.CheckoutLineFilterOptions) (model.CheckoutLineSlice, error) // CheckoutLinesByOption finds and returns model lines filtered using given option
		// CheckoutLinesByCheckoutWithPrefetch(checkoutID string) (model.CheckoutLineSlice, model.ProductVariantSlice, model.ProductSlice, error)
	}
	CheckoutStore interface {
		Upsert(tx boil.ContextTransactor, checkouts model.CheckoutSlice) (model.CheckoutSlice, error)              // Upsert depends on given model's Token property to decide to update or insert it
		FetchCheckoutLinesAndPrefetchRelatedValue(checkout model.Checkout) (model_helper.CheckoutLineInfos, error) // FetchCheckoutLinesAndPrefetchRelatedValue Fetch model lines as CheckoutLineInfo objects.
		GetByOption(option model_helper.CheckoutFilterOptions) (*model.Checkout, error)                            // GetByOption finds and returns 1 model based on given option
		FilterByOption(option model_helper.CheckoutFilterOptions) (model.CheckoutSlice, error)                     // FilterByOption finds and returns a list of model based on given option
		Delete(tx boil.ContextTransactor, ids []string) error                                                      // DeleteCheckoutsByOption deletes model row(s) from database, filtered using given option.  It returns an error indicating if the operation was performed successfully.
		CountCheckouts(options model_helper.CheckoutFilterOptions) (int64, error)
	}
)

// channel
type ChannelStore interface {
	Get(id string) (*model.Channel, error)
	Upsert(tx boil.ContextTransactor, channel model.Channel) (*model.Channel, error)
	FilterByOptions(conds model_helper.ChannelFilterOptions) (model.ChannelSlice, error)
	DeleteChannels(tx boil.ContextTransactor, ids []string) error
}

// app
type (
	AppTokenStore interface {
		// Save(appToken *model.AppToken) (*model.AppToken, error)
	}
	AppStore interface {
		// Save(app *model.App) (*model.App, error)
	}
)

type ClusterDiscoveryStore interface {
	Save(discovery model.ClusterDiscovery) error
	Delete(discovery model.ClusterDiscovery) (bool, error)
	Exists(discovery model.ClusterDiscovery) (bool, error)
	GetAll(discoveryType, clusterName string) (model.ClusterDiscoverySlice, error)
	SetLastPingAt(discovery model.ClusterDiscovery) error
	Cleanup() error
}

type AuditStore interface {
	Save(audit model.Audit) error
	Get(userID string, offset int, limit int) (model.AuditSlice, error)
	PermanentDeleteByUser(userID string) error
}

type TermsOfServiceStore interface {
	Save(termsOfService model.TermsOfService) (*model.TermsOfService, error)
	GetLatest(allowFromCache bool) (*model.TermsOfService, error)
	Get(id string, allowFromCache bool) (*model.TermsOfService, error)
}

type PreferenceStore interface {
	Save(preferences model.PreferenceSlice) error
	GetCategory(userID, category string) (model.PreferenceSlice, error)
	Get(userID, category, name string) (*model.Preference, error)
	GetAll(userID string) (model.PreferenceSlice, error)
	Delete(userID, category, name string) error
	DeleteCategory(userID string, category string) error
	DeleteCategoryAndName(category string, name string) error
	PermanentDeleteByUser(userID string) error
	DeleteUnusedFeatures()
	// CleanupFlagsBatch(limit int64) (int64, error)
}

type JobStore interface {
	Save(job model.Job) (*model.Job, error)
	UpdateOptimistically(job model.Job, currentStatus model.JobStatus) (bool, error)
	UpdateStatus(id string, status model.JobStatus) (*model.Job, error)
	UpdateStatusOptimistically(id string, currentStatus model.JobStatus, newStatus model.JobStatus) (bool, error) // update job status from current status to new status
	Get(mods ...qm.QueryMod) (*model.Job, error)
	FindAll(mods ...qm.QueryMod) (model.JobSlice, error)
	Count(mods ...qm.QueryMod) (int64, error)
	Delete(id string) (string, error)
	// GetAllPage(offset int, limit int) ([]*model.Job, error)
	// GetAllByType(jobType string) ([]*model.Job, error)
	// GetAllByTypePage(jobType string, offset int, limit int) ([]*model.Job, error)
	// GetAllByTypesPage(jobTypes []string, offset int, limit int) ([]*model.Job, error)
	// GetAllByStatus(status model.JobStatus) ([]*model.Job, error)
	// GetNewestJobByStatusAndType(status model.JobStatus, jobType string) (*model.Job, error)
	// GetNewestJobByStatusesAndType(statuses []model.JobStatus, jobType string) (*model.Job, error) // GetNewestJobByStatusesAndType get 1 job from database that has status is one of given statuses, and job type is given jobType. order by created time
	// GetCountByStatusAndType(status string, jobType string) (int64, error)
}

type StatusStore interface {
	Upsert(status model.Status) (*model.Status, error)
	Get(userID string) (*model.Status, error)
	GetByIds(userIds []string) (model.StatusSlice, error)
	ResetAll() error
	GetTotalActiveUsersCount() (int64, error)
	UpdateLastActivityAt(userID string, lastActivityAt int64) error
}

// account stores
type (
	AddressStore interface {
		Upsert(tx boil.ContextTransactor, address model.Address) (*model.Address, error)
		Get(addressID string) (*model.Address, error)                         // Get returns an Address with given addressID is exist
		DeleteAddresses(tx boil.ContextTransactor, addressIDs []string) error // DeleteAddress deletes given address and returns an error
		FilterByOption(option model_helper.AddressFilterOptions) (model.AddressSlice, error)
	}
	UserStore interface {
		ClearCaches()
		Save(user model.User) (*model.User, error)                                      // Save takes an user struct and save into database
		Update(user model.User, allowRoleUpdate bool) (*model_helper.UserUpdate, error) // Update update given user
		UpdateLastPictureUpdate(userID string, updateMillis int64) error
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
		GetEtagForProfiles() string
		UpdateFailedPasswordAttempts(userID string, attempts int) error
		GetSystemAdminProfiles() (map[string]*model.User, error)
		PermanentDelete(userID string) error // PermanentDelete completely delete user from the system
		AnalyticsGetInactiveUsersCount() (int64, error)
		AnalyticsGetExternalUsers(hostDomain string) (bool, error)
		AnalyticsGetSystemAdminCount() (int64, error)
		AnalyticsGetGuestCount() (int64, error)
		ClearAllCustomRoleAssignments() error
		InferSystemInstallDate() (int64, error)
		GetUsersBatchForIndexing(startTime, endTime int64, limit int) ([]*model_helper.UserForIndexing, error)
		Count(options model_helper.UserCountOptions) (int64, error)
		AnalyticsActiveCountForPeriod(startTime int64, endTime int64, options model_helper.UserCountOptions) (int64, error)
		GetAllProfiles(options model_helper.UserGetOptions) (model.UserSlice, error)
		Search(term string, options *model_helper.UserSearchOptions) (model.UserSlice, error)
		AnalyticsActiveCount(time int64, options model_helper.UserCountOptions) (int64, error)
		GetProfileByIds(ctx context.Context, userIds []string, options UserGetByIdsOpts, allowFromCache bool) (model.UserSlice, error)
		IsEmpty() (bool, error)
		Get(ctx context.Context, id string) (*model.User, error)
		Find(options model_helper.UserFilterOptions) (model.UserSlice, error)
	}
	TokenStore interface {
		Save(token model.Token) (*model.Token, error)
		Delete(token string) error
		GetByToken(token string) (*model.Token, error)
		Cleanup() error
		GetAllTokensByType(tokenType model_helper.TokenType) (model.TokenSlice, error)
	}
	UserAccessTokenStore interface {
		Save(token model.UserAccessToken) (*model.UserAccessToken, error)
		DeleteAllForUser(userID string) error
		Delete(tokenID string) error
		Get(tokenID string) (*model.UserAccessToken, error)
		GetAll(conds ...qm.QueryMod) (model.UserAccessTokenSlice, error)
		GetByToken(tokenString string) (*model.UserAccessToken, error)
		Search(term string) (model.UserAccessTokenSlice, error)
		UpdateTokenEnable(tokenID string) error
		UpdateTokenDisable(tokenID string) error
	}
	CustomerEventStore interface {
		Upsert(tx boil.ContextTransactor, customemrEvent model.CustomerEvent) (*model.CustomerEvent, error)
		Get(id string) (*model.CustomerEvent, error)
		Count() (int64, error)
		FilterByOptions(queryMods ...qm.QueryMod) (model.CustomerEventSlice, error)
	}
	StaffNotificationRecipientStore interface {
		Save(notificationRecipient model.StaffNotificationRecipient) (*model.StaffNotificationRecipient, error)
		FilterByOptions(options ...qm.QueryMod) (model.StaffNotificationRecipientSlice, error)
	}
	CustomerNoteStore interface {
		Upsert(note model.CustomerNote) (*model.CustomerNote, error) // Save insert given customer note into database and returns it
		Get(id string) (*model.CustomerNote, error)                  // Get find customer note with given id and returns it
	}
	SessionStore interface {
		Get(ctx context.Context, sessionIDOrToken string) (*model.Session, error)
		Save(session model.Session) (*model.Session, error)
		GetSessions(userID string) (model.SessionSlice, error)
		GetSessionsWithActiveDeviceIds(userID string) (model.SessionSlice, error)
		GetSessionsExpired(thresholdMillis int64, mobileOnly bool, unnotifiedOnly bool) (model.SessionSlice, error)
		UpdateExpiredNotify(sessionid string, notified bool) error
		Remove(sessionIDOrToken string) error
		RemoveAllSessions() error
		PermanentDeleteSessionsByUser(userID string) error
		UpdateExpiresAt(sessionID string, time int64) error
		UpdateLastActivityAt(sessionID string, time int64) error                    // UpdateLastActivityAt
		UpdateRoles(userID string, roles string) (string, error)                    // UpdateRoles updates roles for all sessions that have userId of given userID,
		UpdateDeviceId(id string, deviceID string, expiresAt int64) (string, error) // UpdateDeviceId updates device id for sessions
		UpdateProps(session model.Session) error                                    // UpdateProps update session's props
		AnalyticsSessionCount() (int64, error)                                      // AnalyticsSessionCount counts numbers of sessions
		Cleanup(expiryTime int64, batchSize int64)                                  // Cleanup is called periodicly to remove sessions that are expired
	}
)

type SystemStore interface {
	Save(system model.System) error
	SaveOrUpdate(system model.System) error
	Update(system model.System) error
	Get() (map[string]string, error)
	GetByName(name string) (*model.System, error)
	PermanentDeleteByName(name string) (*model.System, error)
	InsertIfExists(system model.System) (*model.System, error)
	SaveOrUpdateWithWarnMetricHandling(system model.System) error
}

type RoleStore interface {
	Upsert(role model.Role) (*model.Role, error)
	Get(roleID string) (*model.Role, error)
	GetAll() (model.RoleSlice, error)
	GetByName(ctx context.Context, name string) (*model.Role, error)
	GetByNames(names []string) (model.RoleSlice, error)
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
	BulkUpsert(rates model.OpenExchangeRateSlice) (model.OpenExchangeRateSlice, error) // BulkUpsert performs bulk update/insert to given exchange rates
	GetAll() (model.OpenExchangeRateSlice, error)                                      // GetAll returns all exchange currency rates
}

type UserGetByIdsOpts struct {
	IsAdmin bool  // IsAdmin tracks whether or not the request is being made by an administrator. Does nothing when provided by a client.
	Since   int64 // Since filters the users based on their UpdateAt timestamp.
	// Restrict to search in a list of teams and channels. Does nothing when provided by a client.
	// ViewRestrictions *model.ViewUsersRestrictions
}
