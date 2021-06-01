package sqlstore

import "github.com/sitename/sitename/store"

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
}

// setup tables before performing database migration
func (store *SqlStore) setupTables() {
	// account
	store.stores.user = newSqlUserStore(store, store.metrics) // metrics is already set in caller
	store.stores.address = newSqlAddressStore(store)
	// common
	store.stores.audit = newSqlAuditStore(store)
	store.stores.cluster = newSqlClusterDiscoveryStore(store)
	store.stores.session = newSqlSessionStore(store)
	store.stores.system = newSqlSystemStore(store)
	store.stores.preference = newSqlPreferenceStore(store)
	store.stores.token = newSqlTokenStore(store)
	store.stores.status = newSqlStatusStore(store)
	store.stores.job = newSqlJobStore(store)
	store.stores.userAccessToken = newSqlUserAccessTokenStore(store)
	store.stores.TermsOfService = newSqlTermsOfServiceStore(store, store.metrics) // metrics is already set in caller
	store.stores.role = newSqlRoleStore(store)
	store.stores.app = newAppSqlStore(store)
	store.stores.appToken = newSqlAppTokenStore(store)
	// channel
	store.stores.channel = newSqlChannelStore(store)
	// checkout
	store.stores.checkout = newSqlCheckoutStore(store)
	store.stores.checkoutLine = newSqlCheckoutLineStore(store)
	// csv
	store.stores.csvExportEvent = newSqlCsvExportEventStore(store)
	store.stores.csvExportFile = newSqlCsvExportFileStore(store)
	// product and discount
	store.stores.discountVoucher = newSqlVoucherStore(store)
	store.stores.discountVoucherChannelListing = newSqlVoucherChannelListingStore(store)
	store.stores.discountVoucherTranslation = newSqlVoucherTranslationStore(store)
	store.stores.discountVoucherCustomer = newSqlVoucherCustomerStore(store)
	store.stores.discountSale = newSqlDiscountSaleStore(store)
	store.stores.discountSaleChannelListing = newSqlSaleChannelListingStore(store)
	store.stores.discountSaleTranslation = newSqlDiscountSaleTranslationStore(store)
	store.stores.orderDiscount = newSqlOrderDiscountStore(store)
	// gift card
	store.stores.giftCard = newSqlGiftCardStore(store)
	// invoice
	store.stores.invoiceEvent = newSqlInvoiceEventStore(store)
	// menu
	store.stores.menu = newSqlMenuStore(store)
	store.stores.menuItemTranslation = newSqlMenuItemTranslationStore(store)
	// order
	store.stores.order = newSqlOrderStore(store)
	store.stores.orderLine = newSqlOrderLineStore(store)
	store.stores.fulfillment = newSqlFulfillmentStore(store)
	store.stores.fulfillmentLine = newSqlFulfillmentLineStore(store)
	store.stores.orderEvent = newSqlOrderEventStore(store)
	// page
	store.stores.page = newSqlPageStore(store)
	store.stores.pageTranslation = newSqlPageTranslationStore(store)
	store.stores.pageType = newSqlPageTypeStore(store)
	// payment
	store.stores.payment = newSqlPaymentStore(store)
	store.stores.transaction = newSqlPaymentTransactionStore(store)
	// product
	store.stores.category = newSqlCategoryStore(store)
	store.stores.categoryTranslation = newSqlCategoryTranslationStore(store)
	store.stores.productType = newSqlProductTypeStore(store)
	store.stores.product = newSqlProductStore(store)
	store.stores.productTranslation = newSqlProductTranslationStore(store)
	store.stores.productChannelListing = newSqlProductChannelListingStore(store)
	store.stores.productVariant = newSqlProductVariantStore(store)
	store.stores.productVariantTranslation = newSqlProductVariantTranslationStore(store)
	store.stores.productVariantChannelListing = newSqlProductVariantChannelListingStore(store)
	store.stores.digitalContent = newSqlDigitalContentStore(store)
	store.stores.digitalContentUrl = newSqlDigitalContentUrlStore(store)
	store.stores.productMedia = newSqlProductMediaStore(store)
	store.stores.variantMedia = newSqlVariantMediaStore(store)
	store.stores.collectionProduct = newSqlCollectionProductStore(store)
	store.stores.collection = newSqlCollectionStore(store)
	store.stores.collectionChannelListing = newSqlCollectionChannelListingStore(store)
	store.stores.collectionTranslation = newSqlCollectionTranslationStore(store)
	// shipping
	store.stores.shippingMethodTranslation = newSqlShippingMethodTranslationStore(store)
	store.stores.shippingMethodChannelListing = newSqlShippingMethodChannelListingStore(store)
	store.stores.shippingMethodPostalCodeRule = newSqlShippingMethodPostalCodeRuleStore(store)
	store.stores.shippingMethod = newSqlShippingMethodStore(store)
	store.stores.shippingZone = newSqlShippingZoneStore(store)
	// warehouse
	store.stores.warehouse = newSqlWareHouseStore(store)
	store.stores.stock = newSqlStockStore(store)
	store.stores.allocation = newSqlAllocationStore(store)
	// wishlist
	store.stores.wishlist = newSqlWishlistStore(store)
	store.stores.wishlistItem = newSqlWishlistItemStore(store)
	// plugin
	store.stores.pluginConfig = newSqlPluginConfigurationStore(store)
	// compliance
	store.stores.compliance = newSqlComplianceStore(store)
	// attribute
	store.stores.attribute = newSqlAttributeStore(store)
	store.stores.attributeTranslation = newSqlAttributeTranslationStore(store)
	store.stores.attributeValue = newSqlAttributeValueStore(store)
	store.stores.attributeValueTranslation = newSqlAttributeValueTranslationStore(store)
	store.stores.assignedPageAttributeValue = newSqlAssignedPageAttributeValueStore(store)
	store.stores.assignedPageAttribute = newSqlAssignedPageAttributeStore(store)
	store.stores.attributePage = newSqlAttributePageStore(store)
	store.stores.assignedVariantAttributeValue = newSqlAssignedVariantAttributeValueStore(store)
	store.stores.assignedVariantAttribute = newSqlAssignedVariantAttributeStore(store)
	store.stores.attributeVariant = newSqlAttributeVariantStore(store)
	store.stores.assignedProductAttributeValue = newSqlAssignedProductAttributeValueStore(store)
	store.stores.assignedProductAttribute = newSqlAssignedProductAttributeStore(store)
	store.stores.attributeProduct = newSqlAttributeProductStore(store)
	// file info
	store.stores.fileInfo = newSqlFileInfoStore(store, store.metrics)
}

// performs database indexing
func (store *SqlStore) indexingTableFields() {
	// account
	store.stores.user.(*SqlUserStore).createIndexesIfNotExists()
	store.stores.address.(*SqlAddressStore).createIndexesIfNotExists()
	// common
	store.stores.audit.(*SqlAuditStore).createIndexesIfNotExists()
	store.stores.session.(*SqlSessionStore).createIndexesIfNotExists()
	store.stores.system.(*SqlSystemStore).createIndexesIfNotExists()
	// preference
	store.stores.preference.(*SqlPreferenceStore).createIndexesIfNotExists()
	store.stores.preference.(*SqlPreferenceStore).deleteUnusedFeatures()

	store.stores.token.(*SqlTokenStore).createIndexesIfNotExists()
	store.stores.status.(*SqlStatusStore).createIndexesIfNotExists()
	store.stores.job.(*SqlJobStore).createIndexesIfNotExists()
	store.stores.userAccessToken.(*SqlUserAccessTokenStore).createIndexesIfNotExists()
	store.stores.TermsOfService.(*SqlTermsOfServiceStore).createIndexesIfNotExists()
	store.stores.app.(*SqlAppStore).createIndexesIfNotExists()
	store.stores.appToken.(*SqlAppTokenStore).createIndexesIfNotExists()
	// channel
	store.stores.channel.(*SqlChannelStore).createIndexesIfNotExists()
	// checkout
	store.stores.checkout.(*SqlCheckoutStore).createIndexesIfNotExists()
	store.stores.checkoutLine.(*SqlCheckoutLineStore).createIndexesIfNotExists()
	// csv
	store.stores.csvExportEvent.(*SqlCsvExportEventStore).createIndexesIfNotExists()
	store.stores.csvExportFile.(*SqlCsvExportFileStore).createIndexesIfNotExists()
	// product and discount
	store.stores.discountVoucher.(*SqlVoucherStore).createIndexesIfNotExists()
	store.stores.discountVoucherChannelListing.(*SqlVoucherChannelListingStore).createIndexesIfNotExists()
	store.stores.discountVoucherTranslation.(*SqlVoucherTranslationStore).createIndexesIfNotExists()
	store.stores.discountSale.(*SqlDiscountSaleStore).createIndexesIfNotExists()
	store.stores.discountSaleChannelListing.(*SqlSaleChannelListingStore).createIndexesIfNotExists()
	store.stores.discountVoucherCustomer.(*SqlVoucherCustomerStore).createIndexesIfNotExists()
	store.stores.discountSaleTranslation.(*SqlDiscountSaleTranslationStore).createIndexesIfNotExists()
	store.stores.orderDiscount.(*SqlOrderDiscountStore).createIndexesIfNotExists()
	// gift card
	store.stores.giftCard.(*SqlGiftCardStore).createIndexesIfNotExists()
	// invoice
	store.stores.invoiceEvent.(*SqlInvoiceEventStore).createIndexesIfNotExists()
	// menu
	store.stores.menu.(*SqlMenuStore).createIndexesIfNotExists()
	store.stores.menuItemTranslation.(*SqlMenuItemTranslationStore).createIndexesIfNotExists()
	// order
	store.stores.order.(*SqlOrderStore).createIndexesIfNotExists()
	store.stores.orderLine.(*SqlOrderLineStore).createIndexesIfNotExists()
	store.stores.fulfillment.(*SqlFulfillmentStore).createIndexesIfNotExists()
	store.stores.fulfillmentLine.(*SqlFulfillmentLineStore).createIndexesIfNotExists()
	store.stores.orderEvent.(*SqlOrderEventStore).createIndexesIfNotExists()
	// page
	store.stores.page.(*SqlPageStore).createIndexesIfNotExists()
	store.stores.pageTranslation.(*SqlPageTranslationStore).createIndexesIfNotExists()
	store.stores.pageType.(*SqlPageTypeStore).createIndexesIfNotExists()
	// payment
	store.stores.transaction.(*SqlPaymentTransactionStore).createIndexesIfNotExists()
	store.stores.payment.(*SqlPaymentStore).createIndexesIfNotExists()
	// product
	store.stores.category.(*SqlCategoryStore).createIndexesIfNotExists()
	store.stores.categoryTranslation.(*SqlCategoryTranslationStore).createIndexesIfNotExists()
	store.stores.productType.(*SqlProductTypeStore).createIndexesIfNotExists()
	store.stores.product.(*SqlProductStore).createIndexesIfNotExists()
	store.stores.productTranslation.(*SqlProductTranslationStore).createIndexesIfNotExists()
	store.stores.productChannelListing.(*SqlProductChannelListingStore).createIndexesIfNotExists()
	store.stores.productVariant.(*SqlProductVariantStore).createIndexesIfNotExists()
	store.stores.productVariantTranslation.(*SqlProductVariantTranslationStore).createIndexesIfNotExists()
	store.stores.productVariantChannelListing.(*SqlProductVariantChannelListingStore).createIndexesIfNotExists()
	store.stores.digitalContent.(*SqlDigitalContentStore).createIndexesIfNotExists()
	store.stores.digitalContentUrl.(*SqlDigitalContentUrlStore).createIndexesIfNotExists()
	store.stores.productMedia.(*SqlProductMediaStore).createIndexesIfNotExists()
	store.stores.variantMedia.(*SqlVariantMediaStore).createIndexesIfNotExists()
	store.stores.collectionProduct.(*SqlCollectionProductStore).createIndexesIfNotExists()
	store.stores.collection.(*SqlCollectionStore).createIndexesIfNotExists()
	store.stores.collectionChannelListing.(*SqlCollectionChannelListingStore).createIndexesIfNotExists()
	store.stores.collectionTranslation.(*SqlCollectionTranslationStore).createIndexesIfNotExists()
	// shipping
	store.stores.shippingMethodTranslation.(*SqlShippingMethodTranslationStore).createIndexesIfNotExists()
	store.stores.shippingMethodChannelListing.(*SqlShippingMethodChannelListingStore).createIndexesIfNotExists()
	store.stores.shippingMethodPostalCodeRule.(*SqlShippingMethodPostalCodeRuleStore).createIndexesIfNotExists()
	store.stores.shippingMethod.(*SqlShippingMethodStore).createIndexesIfNotExists()
	store.stores.shippingZone.(*SqlShippingZoneStore).createIndexesIfNotExists()
	// warehouse
	store.stores.warehouse.(*SqlWareHouseStore).createIndexesIfNotExists()
	store.stores.stock.(*SqlStockStore).createIndexesIfNotExists()
	store.stores.allocation.(*SqlAllocationStore).createIndexesIfNotExists()
	// wishlist
	store.stores.wishlist.(*SqlWishlistStore).createIndexesIfNotExists()
	store.stores.wishlistItem.(*SqlWishlistItemStore).createIndexesIfNotExists()
	// plugin
	store.stores.pluginConfig.(*SqlPluginConfigurationStore).createIndexesIfNotExists()
	// compliance
	store.stores.compliance.(*SqlComplianceStore).createIndexesIfNotExists()
	// attribute
	store.stores.attribute.(*SqlAttributeStore).createIndexesIfNotExists()
	store.stores.attributeTranslation.(*SqlAttributeTranslationStore).createIndexesIfNotExists()
	store.stores.attributeValue.(*SqlAttributeValueStore).createIndexesIfNotExists()
	store.stores.attributeValueTranslation.(*SqlAttributeValueTranslationStore).createIndexesIfNotExists()
	store.stores.assignedPageAttributeValue.(*SqlAssignedPageAttributeValueStore).createIndexesIfNotExists()
	store.stores.assignedPageAttribute.(*SqlAssignedPageAttributeStore).createIndexesIfNotExists()
	store.stores.attributePage.(*SqlAttributePageStore).createIndexesIfNotExists()
	store.stores.assignedVariantAttributeValue.(*SqlAssignedVariantAttributeValueStore).createIndexesIfNotExists()
	store.stores.assignedVariantAttribute.(*SqlAssignedVariantAttributeStore).createIndexesIfNotExists()
	store.stores.attributeVariant.(*SqlAttributeVariantStore).createIndexesIfNotExists()
	store.stores.assignedProductAttributeValue.(*SqlAssignedProductAttributeValueStore).createIndexesIfNotExists()
	store.stores.assignedProductAttribute.(*SqlAssignedProductAttributeStore).createIndexesIfNotExists()
	store.stores.attributeProduct.(*SqlAttributeProductStore).createIndexesIfNotExists()
	// file info
	store.stores.fileInfo.(*SqlFileInfoStore).createIndexesIfNotExists()
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
func (ss *SqlStore) VoucherStore() store.DiscountVoucherStore {
	return ss.stores.discountVoucher
}
func (ss *SqlStore) VoucherChannelListing() store.VoucherChannelListingStore {
	return ss.stores.discountVoucherChannelListing
}
func (ss *SqlStore) VoucherTranslation() store.VoucherTranslationStore {
	return ss.stores.discountVoucherTranslation
}
func (ss *SqlStore) VoucherCustomer() store.DiscountVoucherCustomerStore {
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
