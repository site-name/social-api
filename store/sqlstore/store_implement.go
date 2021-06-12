package sqlstore

import (
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/sqlstore/account"
	"github.com/sitename/sitename/store/sqlstore/app"
	"github.com/sitename/sitename/store/sqlstore/attribute"
	"github.com/sitename/sitename/store/sqlstore/audit"
	"github.com/sitename/sitename/store/sqlstore/channel"
	"github.com/sitename/sitename/store/sqlstore/checkout"
	"github.com/sitename/sitename/store/sqlstore/cluster"
	"github.com/sitename/sitename/store/sqlstore/compliance"
	"github.com/sitename/sitename/store/sqlstore/csv"
	"github.com/sitename/sitename/store/sqlstore/discount"
	"github.com/sitename/sitename/store/sqlstore/file"
	"github.com/sitename/sitename/store/sqlstore/giftcard"
	"github.com/sitename/sitename/store/sqlstore/invoice"
	"github.com/sitename/sitename/store/sqlstore/job"
	"github.com/sitename/sitename/store/sqlstore/menu"
	"github.com/sitename/sitename/store/sqlstore/order"
	"github.com/sitename/sitename/store/sqlstore/page"
	"github.com/sitename/sitename/store/sqlstore/payment"
	"github.com/sitename/sitename/store/sqlstore/plugin"
	"github.com/sitename/sitename/store/sqlstore/preference"
	"github.com/sitename/sitename/store/sqlstore/product"
	"github.com/sitename/sitename/store/sqlstore/role"
	"github.com/sitename/sitename/store/sqlstore/session"
	"github.com/sitename/sitename/store/sqlstore/shipping"
	"github.com/sitename/sitename/store/sqlstore/status"
	"github.com/sitename/sitename/store/sqlstore/system"
	"github.com/sitename/sitename/store/sqlstore/warehouse"
	"github.com/sitename/sitename/store/sqlstore/wishlist"
)

type SqlStoreStores struct {
	user                          store.UserStore                          // account models
	address                       store.AddressStore                       //
	audit                         store.AuditStore                         // common
	cluster                       store.ClusterDiscoveryStore              //
	session                       store.SessionStore                       //
	system                        store.SystemStore                        //
	preference                    store.PreferenceStore                    //
	token                         store.TokenStore                         //
	status                        store.StatusStore                        //
	job                           store.JobStore                           //
	userAccessToken               store.UserAccessTokenStore               //
	role                          store.RoleStore                          //
	TermsOfService                store.TermsOfServiceStore                //
	app                           store.AppStore                           //
	appToken                      store.AppTokenStore                      //
	channel                       store.ChannelStore                       // channel models
	checkout                      store.CheckoutStore                      // checkout models
	checkoutLine                  store.CheckoutLineStore                  //
	csvExportEvent                store.CsvExportEventStore                // csv models
	csvExportFile                 store.CsvExportFileStore                 //
	discountVoucher               store.DiscountVoucherStore               // product and discount models
	discountVoucherChannelListing store.VoucherChannelListingStore         //
	discountVoucherTranslation    store.VoucherTranslationStore            //
	discountVoucherCustomer       store.DiscountVoucherCustomerStore       //
	discountSale                  store.DiscountSaleStore                  //
	discountSaleTranslation       store.DiscountSaleTranslationStore       //
	discountSaleChannelListing    store.DiscountSaleChannelListingStore    //
	orderDiscount                 store.OrderDiscountStore                 //
	giftCard                      store.GiftCardStore                      // gift card models
	invoiceEvent                  store.InvoiceEventStore                  // invoice models
	menu                          store.MenuStore                          // menu models
	menuItemTranslation           store.MenuItemTranslationStore           //
	order                         store.OrderStore                         // order models
	orderLine                     store.OrderLineStore                     //
	fulfillment                   store.FulfillmentStore                   //
	fulfillmentLine               store.FulfillmentLineStore               //
	orderEvent                    store.OrderEventStore                    //
	page                          store.PageStore                          // page models
	pageType                      store.PageTypeStore                      //
	pageTranslation               store.PageTranslationStore               //
	payment                       store.PaymentStore                       // payment models
	transaction                   store.PaymentTransactionStore            //
	category                      store.CategoryStore                      // product models
	categoryTranslation           store.CategoryTranslationStore           //
	productType                   store.ProductTypeStore                   //
	product                       store.ProductStore                       //
	productTranslation            store.ProductTranslationStore            //
	productChannelListing         store.ProductChannelListingStore         //
	productVariant                store.ProductVariantStore                //
	productVariantTranslation     store.ProductVariantTranslationStore     //
	productVariantChannelListing  store.ProductVariantChannelListingStore  //
	digitalContent                store.DigitalContentStore                //
	digitalContentUrl             store.DigitalContentUrlStore             //
	productMedia                  store.ProductMediaStore                  //
	variantMedia                  store.VariantMediaStore                  //
	collectionProduct             store.CollectionProductStore             //
	collection                    store.CollectionStore                    //
	collectionChannelListing      store.CollectionChannelListingStore      //
	collectionTranslation         store.CollectionTranslationStore         //
	shippingMethodTranslation     store.ShippingMethodTranslationStore     // shipping models
	shippingMethodChannelListing  store.ShippingMethodChannelListingStore  //
	shippingMethodPostalCodeRule  store.ShippingMethodPostalCodeRuleStore  //
	shippingMethod                store.ShippingMethodStore                //
	shippingZone                  store.ShippingZoneStore                  //
	warehouse                     store.WarehouseStore                     // warehouse
	stock                         store.StockStore                         //
	allocation                    store.AllocationStore                    //
	wishlist                      store.WishlistStore                      // wishlist models
	wishlistItem                  store.WishlistItemStore                  //
	pluginConfig                  store.PluginConfigurationStore           // plugin models
	compliance                    store.ComplianceStore                    // compliance models
	attribute                     store.AttributeStore                     // attribute
	attributeTranslation          store.AttributeTranslationStore          //
	attributeValue                store.AttributeValueStore                //
	attributeValueTranslation     store.AttributeValueTranslationStore     //
	assignedPageAttributeValue    store.AssignedPageAttributeValueStore    //
	assignedPageAttribute         store.AssignedPageAttributeStore         //
	attributePage                 store.AttributePageStore                 //
	assignedVariantAttributeValue store.AssignedVariantAttributeValueStore //
	assignedVariantAttribute      store.AssignedVariantAttributeStore      //
	attributeVariant              store.AttributeVariantStore              //
	assignedProductAttributeValue store.AssignedProductAttributeValueStore //
	assignedProductAttribute      store.AssignedProductAttributeStore      //
	attributeProduct              store.AttributeProductStore              //
	fileInfo                      store.FileInfoStore                      // file info models
	uploadSession                 store.UploadSessionStore                 // upload session models
}

// setup tables before performing database migration
func (store *SqlStore) setupTables() {
	// account
	store.stores.user = account.NewSqlUserStore(store, store.metrics) // metrics is already set in caller
	store.stores.address = account.NewSqlAddressStore(store)
	// general
	store.stores.audit = audit.NewSqlAuditStore(store)
	store.stores.cluster = cluster.NewSqlClusterDiscoveryStore(store)
	store.stores.session = session.NewSqlSessionStore(store)
	store.stores.system = system.NewSqlSystemStore(store)
	store.stores.preference = preference.NewSqlPreferenceStore(store)
	store.stores.token = account.NewSqlTokenStore(store)
	store.stores.status = status.NewSqlStatusStore(store)
	store.stores.job = job.NewSqlJobStore(store)
	store.stores.userAccessToken = account.NewSqlUserAccessTokenStore(store)
	store.stores.TermsOfService = account.NewSqlTermsOfServiceStore(store, store.metrics) // metrics is already set in caller
	store.stores.role = role.NewSqlRoleStore(store)
	store.stores.app = app.NewAppSqlStore(store)
	store.stores.appToken = app.NewSqlAppTokenStore(store)
	// channel
	store.stores.channel = channel.NewSqlChannelStore(store)
	// checkout
	store.stores.checkout = checkout.NewSqlCheckoutStore(store)
	store.stores.checkoutLine = checkout.NewSqlCheckoutLineStore(store)
	// csv
	store.stores.csvExportEvent = csv.NewSqlCsvExportEventStore(store)
	store.stores.csvExportFile = csv.NewSqlCsvExportFileStore(store)
	// product and discount
	store.stores.discountVoucher = discount.NewSqlVoucherStore(store)
	store.stores.discountVoucherChannelListing = discount.NewSqlVoucherChannelListingStore(store)
	store.stores.discountVoucherTranslation = discount.NewSqlVoucherTranslationStore(store)
	store.stores.discountVoucherCustomer = discount.NewSqlVoucherCustomerStore(store)
	store.stores.discountSale = discount.NewSqlDiscountSaleStore(store)
	store.stores.discountSaleChannelListing = discount.NewSqlSaleChannelListingStore(store)
	store.stores.discountSaleTranslation = discount.NewSqlDiscountSaleTranslationStore(store)
	store.stores.orderDiscount = discount.NewSqlOrderDiscountStore(store)
	// gift card
	store.stores.giftCard = giftcard.NewSqlGiftCardStore(store)
	// invoice
	store.stores.invoiceEvent = invoice.NewSqlInvoiceEventStore(store)
	// menu
	store.stores.menu = menu.NewSqlMenuStore(store)
	store.stores.menuItemTranslation = menu.NewSqlMenuItemTranslationStore(store)
	// order
	store.stores.order = order.NewSqlOrderStore(store)
	store.stores.orderLine = order.NewSqlOrderLineStore(store)
	store.stores.fulfillment = order.NewSqlFulfillmentStore(store)
	store.stores.fulfillmentLine = order.NewSqlFulfillmentLineStore(store)
	store.stores.orderEvent = order.NewSqlOrderEventStore(store)
	// page
	store.stores.page = page.NewSqlPageStore(store)
	store.stores.pageTranslation = page.NewSqlPageTranslationStore(store)
	store.stores.pageType = page.NewSqlPageTypeStore(store)
	// payment
	store.stores.payment = payment.NewSqlPaymentStore(store)
	store.stores.transaction = payment.NewSqlPaymentTransactionStore(store)
	// product
	store.stores.category = product.NewSqlCategoryStore(store)
	store.stores.categoryTranslation = product.NewSqlCategoryTranslationStore(store)
	store.stores.productType = product.NewSqlProductTypeStore(store)
	store.stores.product = product.NewSqlProductStore(store)
	store.stores.productTranslation = product.NewSqlProductTranslationStore(store)
	store.stores.productChannelListing = product.NewSqlProductChannelListingStore(store)
	store.stores.productVariant = product.NewSqlProductVariantStore(store)
	store.stores.productVariantTranslation = product.NewSqlProductVariantTranslationStore(store)
	store.stores.productVariantChannelListing = product.NewSqlProductVariantChannelListingStore(store)
	store.stores.digitalContent = product.NewSqlDigitalContentStore(store)
	store.stores.digitalContentUrl = product.NewSqlDigitalContentUrlStore(store)
	store.stores.productMedia = product.NewSqlProductMediaStore(store)
	store.stores.variantMedia = product.NewSqlVariantMediaStore(store)
	store.stores.collectionProduct = product.NewSqlCollectionProductStore(store)
	store.stores.collection = product.NewSqlCollectionStore(store)
	store.stores.collectionChannelListing = product.NewSqlCollectionChannelListingStore(store)
	store.stores.collectionTranslation = product.NewSqlCollectionTranslationStore(store)
	// shipping
	store.stores.shippingMethodTranslation = shipping.NewSqlShippingMethodTranslationStore(store)
	store.stores.shippingMethodChannelListing = shipping.NewSqlShippingMethodChannelListingStore(store)
	store.stores.shippingMethodPostalCodeRule = shipping.NewSqlShippingMethodPostalCodeRuleStore(store)
	store.stores.shippingMethod = shipping.NewSqlShippingMethodStore(store)
	store.stores.shippingZone = shipping.NewSqlShippingZoneStore(store)
	// warehouse
	store.stores.warehouse = warehouse.NewSqlWareHouseStore(store)
	store.stores.stock = warehouse.NewSqlStockStore(store)
	store.stores.allocation = warehouse.NewSqlAllocationStore(store)
	// wishlist
	store.stores.wishlist = wishlist.NewSqlWishlistStore(store)
	store.stores.wishlistItem = wishlist.NewSqlWishlistItemStore(store)
	// plugin
	store.stores.pluginConfig = plugin.NewSqlPluginConfigurationStore(store)
	// compliance
	store.stores.compliance = compliance.NewSqlComplianceStore(store)
	// attribute
	store.stores.attribute = attribute.NewSqlAttributeStore(store)
	store.stores.attributeTranslation = attribute.NewSqlAttributeTranslationStore(store)
	store.stores.attributeValue = attribute.NewSqlAttributeValueStore(store)
	store.stores.attributeValueTranslation = attribute.NewSqlAttributeValueTranslationStore(store)
	store.stores.assignedPageAttributeValue = attribute.NewSqlAssignedPageAttributeValueStore(store)
	store.stores.assignedPageAttribute = attribute.NewSqlAssignedPageAttributeStore(store)
	store.stores.attributePage = attribute.NewSqlAttributePageStore(store)
	store.stores.assignedVariantAttributeValue = attribute.NewSqlAssignedVariantAttributeValueStore(store)
	store.stores.assignedVariantAttribute = attribute.NewSqlAssignedVariantAttributeStore(store)
	store.stores.attributeVariant = attribute.NewSqlAttributeVariantStore(store)
	store.stores.assignedProductAttributeValue = attribute.NewSqlAssignedProductAttributeValueStore(store)
	store.stores.assignedProductAttribute = attribute.NewSqlAssignedProductAttributeStore(store)
	store.stores.attributeProduct = attribute.NewSqlAttributeProductStore(store)
	// file info & upload session
	store.stores.fileInfo = file.NewSqlFileInfoStore(store, store.metrics)
	store.stores.uploadSession = file.NewSqlUploadSessionStore(store)
}

// performs database indexing
func (store *SqlStore) indexingTableFields() {
	// account
	store.stores.user.CreateIndexesIfNotExists()
	store.stores.address.CreateIndexesIfNotExists()
	// common
	store.stores.audit.CreateIndexesIfNotExists()
	store.stores.session.CreateIndexesIfNotExists()
	store.stores.system.CreateIndexesIfNotExists()
	// preference
	store.stores.preference.CreateIndexesIfNotExists()
	store.stores.preference.DeleteUnusedFeatures()

	store.stores.token.CreateIndexesIfNotExists()
	store.stores.status.CreateIndexesIfNotExists()
	store.stores.job.CreateIndexesIfNotExists()
	store.stores.userAccessToken.CreateIndexesIfNotExists()
	store.stores.TermsOfService.CreateIndexesIfNotExists()
	// role
	store.stores.role.CreateIndexesIfNotExists()
	// app
	store.stores.app.CreateIndexesIfNotExists()
	store.stores.appToken.CreateIndexesIfNotExists()
	// channel
	store.stores.channel.CreateIndexesIfNotExists()
	// checkoutproduct.N	store.stores.checkout.CreateIndexesIfNotExists()
	store.stores.checkoutLine.CreateIndexesIfNotExists()
	// csv
	store.stores.csvExportEvent.CreateIndexesIfNotExists()
	store.stores.csvExportFile.CreateIndexesIfNotExists()
	// product and discount
	store.stores.discountVoucher.CreateIndexesIfNotExists()
	store.stores.discountVoucherChannelListing.CreateIndexesIfNotExists()
	store.stores.discountVoucherTranslation.CreateIndexesIfNotExists()
	store.stores.discountSale.CreateIndexesIfNotExists()
	store.stores.discountSaleChannelListing.CreateIndexesIfNotExists()
	store.stores.discountVoucherCustomer.CreateIndexesIfNotExists()
	store.stores.discountSaleTranslation.CreateIndexesIfNotExists()
	store.stores.orderDiscount.CreateIndexesIfNotExists()
	// gift card
	store.stores.giftCard.CreateIndexesIfNotExists()
	// invoice
	store.stores.invoiceEvent.CreateIndexesIfNotExists()
	// menu
	store.stores.menu.CreateIndexesIfNotExists()
	store.stores.menuItemTranslation.CreateIndexesIfNotExists()
	// order
	store.stores.order.CreateIndexesIfNotExists()
	store.stores.orderLine.CreateIndexesIfNotExists()
	store.stores.fulfillment.CreateIndexesIfNotExists()
	store.stores.fulfillmentLine.CreateIndexesIfNotExists()
	store.stores.orderEvent.CreateIndexesIfNotExists()
	// page
	store.stores.page.CreateIndexesIfNotExists()
	store.stores.pageTranslation.CreateIndexesIfNotExists()
	store.stores.pageType.CreateIndexesIfNotExists()
	// payment
	store.stores.transaction.CreateIndexesIfNotExists()
	store.stores.payment.CreateIndexesIfNotExists()
	// product
	store.stores.category.CreateIndexesIfNotExists()
	store.stores.categoryTranslation.CreateIndexesIfNotExists()
	store.stores.productType.CreateIndexesIfNotExists()
	store.stores.product.CreateIndexesIfNotExists()
	store.stores.productTranslation.CreateIndexesIfNotExists()
	store.stores.productChannelListing.CreateIndexesIfNotExists()
	store.stores.productVariant.CreateIndexesIfNotExists()
	store.stores.productVariantTranslation.CreateIndexesIfNotExists()
	store.stores.productVariantChannelListing.CreateIndexesIfNotExists()
	store.stores.digitalContent.CreateIndexesIfNotExists()
	store.stores.digitalContentUrl.CreateIndexesIfNotExists()
	store.stores.productMedia.CreateIndexesIfNotExists()
	store.stores.variantMedia.CreateIndexesIfNotExists()
	store.stores.collectionProduct.CreateIndexesIfNotExists()
	store.stores.collection.CreateIndexesIfNotExists()
	store.stores.collectionChannelListing.CreateIndexesIfNotExists()
	store.stores.collectionTranslation.CreateIndexesIfNotExists()
	// shipping
	store.stores.shippingMethodTranslation.CreateIndexesIfNotExists()
	store.stores.shippingMethodChannelListing.CreateIndexesIfNotExists()
	store.stores.shippingMethodPostalCodeRule.CreateIndexesIfNotExists()
	store.stores.shippingMethod.CreateIndexesIfNotExists()
	store.stores.shippingZone.CreateIndexesIfNotExists()
	// warehouse
	store.stores.warehouse.CreateIndexesIfNotExists()
	store.stores.stock.CreateIndexesIfNotExists()
	store.stores.allocation.CreateIndexesIfNotExists()
	// wishlist
	store.stores.wishlist.CreateIndexesIfNotExists()
	store.stores.wishlistItem.CreateIndexesIfNotExists()
	// plugin
	store.stores.pluginConfig.CreateIndexesIfNotExists()
	// compliance
	store.stores.compliance.CreateIndexesIfNotExists()
	// attribute
	store.stores.attribute.CreateIndexesIfNotExists()
	store.stores.attributeTranslation.CreateIndexesIfNotExists()
	store.stores.attributeValue.CreateIndexesIfNotExists()
	store.stores.attributeValueTranslation.CreateIndexesIfNotExists()
	store.stores.assignedPageAttributeValue.CreateIndexesIfNotExists()
	store.stores.assignedPageAttribute.CreateIndexesIfNotExists()
	store.stores.attributePage.CreateIndexesIfNotExists()
	store.stores.assignedVariantAttributeValue.CreateIndexesIfNotExists()
	store.stores.assignedVariantAttribute.CreateIndexesIfNotExists()
	store.stores.attributeVariant.CreateIndexesIfNotExists()
	store.stores.assignedProductAttributeValue.CreateIndexesIfNotExists()
	store.stores.assignedProductAttribute.CreateIndexesIfNotExists()
	store.stores.attributeProduct.CreateIndexesIfNotExists()
	// file info
	store.stores.fileInfo.CreateIndexesIfNotExists()
	// upload session
	store.stores.uploadSession.CreateIndexesIfNotExists()
}

// account
func (ss *SqlStore) Address() store.AddressStore {
	return ss.stores.address
}
func (ss *SqlStore) User() store.UserStore {
	return ss.stores.user
}
func (ss *SqlStore) App() store.AppStore {
	return ss.stores.app
}

// common
func (ss *SqlStore) Session() store.SessionStore {
	return ss.stores.session
}
func (ss *SqlStore) Audit() store.AuditStore {
	return ss.stores.audit
}
func (ss *SqlStore) ClusterDiscovery() store.ClusterDiscoveryStore {
	return ss.stores.cluster
}
func (ss *SqlStore) System() store.SystemStore {
	return ss.stores.system
}
func (ss *SqlStore) Preference() store.PreferenceStore {
	return ss.stores.preference
}
func (ss *SqlStore) Token() store.TokenStore {
	return ss.stores.token
}
func (ss *SqlStore) Status() store.StatusStore {
	return ss.stores.status
}
func (ss *SqlStore) Job() store.JobStore {
	return ss.stores.job
}
func (ss *SqlStore) UserAccessToken() store.UserAccessTokenStore {
	return ss.stores.userAccessToken
}
func (ss *SqlStore) Role() store.RoleStore {
	return ss.stores.role
}
func (ss *SqlStore) TermsOfService() store.TermsOfServiceStore {
	return ss.stores.TermsOfService
}
func (ss *SqlStore) AppToken() store.AppTokenStore {
	return ss.stores.appToken
}

// channel
func (ss *SqlStore) Channel() store.ChannelStore {
	return ss.stores.channel
}

// checkout
func (ss *SqlStore) Checkout() store.CheckoutStore {
	return ss.stores.checkout
}
func (ss *SqlStore) CheckoutLine() store.CheckoutLineStore {
	return ss.stores.checkoutLine
}

// csv
func (ss *SqlStore) CsvExportEvent() store.CsvExportEventStore {
	return ss.stores.csvExportEvent
}
func (ss *SqlStore) CsvExportFile() store.CsvExportFileStore {
	return ss.stores.csvExportFile
}

// product and discount
func (ss *SqlStore) DiscountVoucher() store.DiscountVoucherStore {
	return ss.stores.discountVoucher
}
func (ss *SqlStore) VoucherChannelListing() store.VoucherChannelListingStore {
	return ss.stores.discountVoucherChannelListing
}
func (ss *SqlStore) VoucherTranslation() store.VoucherTranslationStore {
	return ss.stores.discountVoucherTranslation
}
func (ss *SqlStore) DiscountVoucherCustomer() store.DiscountVoucherCustomerStore {
	return ss.stores.discountVoucherCustomer
}
func (ss *SqlStore) DiscountSale() store.DiscountSaleStore {
	return ss.stores.discountSale
}
func (ss *SqlStore) DiscountSaleChannelListing() store.DiscountSaleChannelListingStore {
	return ss.stores.discountSaleChannelListing
}
func (ss *SqlStore) OrderDiscount() store.OrderDiscountStore {
	return ss.stores.orderDiscount
}
func (ss *SqlStore) DiscountSaleTranslation() store.DiscountSaleTranslationStore {
	return ss.stores.discountSaleTranslation
}

// gift card
func (ss *SqlStore) GiftCard() store.GiftCardStore {
	return ss.stores.giftCard
}

// menu
func (ss *SqlStore) Menu() store.MenuStore {
	return ss.stores.menu
}
func (ss *SqlStore) MenuItemTranslation() store.MenuItemTranslationStore {
	return ss.stores.menuItemTranslation
}

// invoice
func (ss *SqlStore) InvoiceEvent() store.InvoiceEventStore {
	return ss.stores.invoiceEvent
}

// order
func (ss *SqlStore) Order() store.OrderStore {
	return ss.stores.order
}
func (ss *SqlStore) OrderLine() store.OrderLineStore {
	return ss.stores.orderLine
}
func (ss *SqlStore) Fulfillment() store.FulfillmentStore {
	return ss.stores.fulfillment
}
func (ss *SqlStore) FulfillmentLine() store.FulfillmentLineStore {
	return ss.stores.fulfillmentLine
}
func (ss *SqlStore) OrderEvent() store.OrderEventStore {
	return ss.stores.orderEvent
}

// page
func (ss *SqlStore) Page() store.PageStore {
	return ss.stores.page
}
func (ss *SqlStore) PageType() store.PageTypeStore {
	return ss.stores.pageType
}
func (ss *SqlStore) PageTranslation() store.PageTranslationStore {
	return ss.stores.pageTranslation
}

// payment
func (ss *SqlStore) Payment() store.PaymentStore {
	return ss.stores.payment
}
func (ss *SqlStore) PaymentTransaction() store.PaymentTransactionStore {
	return ss.stores.transaction
}

// product
func (ss *SqlStore) Category() store.CategoryStore {
	return ss.stores.category
}
func (ss *SqlStore) CategoryTranslation() store.CategoryTranslationStore {
	return ss.stores.categoryTranslation
}
func (ss *SqlStore) ProductType() store.ProductTypeStore {
	return ss.stores.productType
}
func (ss *SqlStore) Product() store.ProductStore {
	return ss.stores.product
}
func (ss *SqlStore) ProductTranslation() store.ProductTranslationStore {
	return ss.stores.productTranslation
}
func (ss *SqlStore) ProductChannelListing() store.ProductChannelListingStore {
	return ss.stores.productChannelListing
}
func (ss *SqlStore) ProductVariant() store.ProductVariantStore {
	return ss.stores.productVariant
}
func (ss *SqlStore) ProductVariantTranslation() store.ProductVariantTranslationStore {
	return ss.stores.productVariantTranslation
}
func (ss *SqlStore) ProductVariantChannelListing() store.ProductVariantChannelListingStore {
	return ss.stores.productVariantChannelListing
}
func (ss *SqlStore) DigitalContent() store.DigitalContentStore {
	return ss.stores.digitalContent
}
func (ss *SqlStore) DigitalContentUrl() store.DigitalContentUrlStore {
	return ss.stores.digitalContentUrl
}
func (ss *SqlStore) ProductMedia() store.ProductMediaStore {
	return ss.stores.productMedia
}
func (ss *SqlStore) VariantMedia() store.VariantMediaStore {
	return ss.stores.variantMedia
}
func (ss *SqlStore) CollectionProduct() store.CollectionProductStore {
	return ss.stores.collectionProduct
}
func (ss *SqlStore) Collection() store.CollectionStore {
	return ss.stores.collection
}
func (ss *SqlStore) CollectionChannelListing() store.CollectionChannelListingStore {
	return ss.stores.collectionChannelListing
}
func (ss *SqlStore) CollectionTranslation() store.CollectionTranslationStore {
	return ss.stores.collectionTranslation
}

// shipping
func (ss *SqlStore) ShippingMethodTranslation() store.ShippingMethodTranslationStore {
	return ss.stores.shippingMethodTranslation
}
func (ss *SqlStore) ShippingMethodChannelListing() store.ShippingMethodChannelListingStore {
	return ss.stores.shippingMethodChannelListing
}
func (ss *SqlStore) ShippingMethodPostalCodeRule() store.ShippingMethodPostalCodeRuleStore {
	return ss.stores.shippingMethodPostalCodeRule
}
func (ss *SqlStore) ShippingMethod() store.ShippingMethodStore {
	return ss.stores.shippingMethod
}
func (ss *SqlStore) ShippingZone() store.ShippingZoneStore {
	return ss.stores.shippingZone
}

// warehouse
func (ss *SqlStore) Allocation() store.AllocationStore {
	return ss.stores.allocation
}
func (ss *SqlStore) Warehouse() store.WarehouseStore {
	return ss.stores.warehouse
}
func (ss *SqlStore) Stock() store.StockStore {
	return ss.stores.stock
}

// wishlist
func (ss *SqlStore) Wishlist() store.WishlistStore {
	return ss.stores.wishlist
}
func (ss *SqlStore) WishlistItem() store.WishlistItemStore {
	return ss.stores.wishlistItem
}

// plugin
func (ss *SqlStore) PluginConfiguration() store.PluginConfigurationStore {
	return ss.stores.pluginConfig
}

// compliance
func (ss *SqlStore) Compliance() store.ComplianceStore {
	return ss.stores.compliance
}

// attribute
func (ss *SqlStore) Attribute() store.AttributeStore {
	return ss.stores.attribute
}
func (ss *SqlStore) AttributeTranslation() store.AttributeTranslationStore {
	return ss.stores.attributeTranslation
}
func (ss *SqlStore) AttributeValue() store.AttributeValueStore {
	return ss.stores.attributeValue
}
func (ss *SqlStore) AttributeValueTranslation() store.AttributeValueTranslationStore {
	return ss.stores.attributeValueTranslation
}
func (ss *SqlStore) AssignedPageAttributeValue() store.AssignedPageAttributeValueStore {
	return ss.stores.assignedPageAttributeValue
}
func (ss *SqlStore) AssignedPageAttribute() store.AssignedPageAttributeStore {
	return ss.stores.assignedPageAttribute
}
func (ss *SqlStore) AttributePage() store.AttributePageStore {
	return ss.stores.attributePage
}
func (ss *SqlStore) AssignedVariantAttributeValue() store.AssignedVariantAttributeValueStore {
	return ss.stores.assignedVariantAttributeValue
}
func (ss *SqlStore) AssignedVariantAttribute() store.AssignedVariantAttributeStore {
	return ss.stores.assignedVariantAttribute
}
func (ss *SqlStore) AttributeVariant() store.AttributeVariantStore {
	return ss.stores.attributeVariant
}
func (ss *SqlStore) AssignedProductAttributeValue() store.AssignedProductAttributeValueStore {
	return ss.stores.assignedProductAttributeValue
}
func (ss *SqlStore) AssignedProductAttribute() store.AssignedProductAttributeStore {
	return ss.stores.assignedProductAttribute
}
func (ss *SqlStore) AttributeProduct() store.AttributeProductStore {
	return ss.stores.attributeProduct
}

// file info
func (ss *SqlStore) FileInfo() store.FileInfoStore {
	return ss.stores.fileInfo
}

// upload session
func (ss *SqlStore) UploadSession() store.UploadSessionStore {
	return ss.stores.uploadSession
}
