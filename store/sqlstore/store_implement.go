// Code generated by "make store-layers"
// DO NOT EDIT

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
	"github.com/sitename/sitename/store/sqlstore/external_services"
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
	"github.com/sitename/sitename/store/sqlstore/shipping"
	"github.com/sitename/sitename/store/sqlstore/shop"
	"github.com/sitename/sitename/store/sqlstore/system"
	"github.com/sitename/sitename/store/sqlstore/warehouse"
	"github.com/sitename/sitename/store/sqlstore/wishlist"
)

type SqlStoreStores struct {
	address                       store.AddressStore
	allocation                    store.AllocationStore
	app                           store.AppStore
	appToken                      store.AppTokenStore
	assignedPageAttribute         store.AssignedPageAttributeStore
	assignedPageAttributeValue    store.AssignedPageAttributeValueStore
	assignedProductAttribute      store.AssignedProductAttributeStore
	assignedProductAttributeValue store.AssignedProductAttributeValueStore
	attribute                     store.AttributeStore
	attributePage                 store.AttributePageStore
	attributeTranslation          store.AttributeTranslationStore
	attributeValue                store.AttributeValueStore
	attributeValueTranslation     store.AttributeValueTranslationStore
	audit                         store.AuditStore
	category                      store.CategoryStore
	categoryTranslation           store.CategoryTranslationStore
	channel                       store.ChannelStore
	checkout                      store.CheckoutStore
	checkoutLine                  store.CheckoutLineStore
	clusterDiscovery              store.ClusterDiscoveryStore
	collection                    store.CollectionStore
	collectionChannelListing      store.CollectionChannelListingStore
	collectionProduct             store.CollectionProductStore
	collectionTranslation         store.CollectionTranslationStore
	compliance                    store.ComplianceStore
	csvExportEvent                store.CsvExportEventStore
	csvExportFile                 store.CsvExportFileStore
	customProductAttribute        store.CustomProductAttributeStore
	customerEvent                 store.CustomerEventStore
	customerNote                  store.CustomerNoteStore
	digitalContent                store.DigitalContentStore
	digitalContentUrl             store.DigitalContentUrlStore
	discountSale                  store.DiscountSaleStore
	discountSaleChannelListing    store.DiscountSaleChannelListingStore
	discountSaleTranslation       store.DiscountSaleTranslationStore
	discountVoucher               store.DiscountVoucherStore
	fileInfo                      store.FileInfoStore
	fulfillment                   store.FulfillmentStore
	fulfillmentLine               store.FulfillmentLineStore
	giftCard                      store.GiftCardStore
	giftcardEvent                 store.GiftcardEventStore
	invoice                       store.InvoiceStore
	invoiceEvent                  store.InvoiceEventStore
	job                           store.JobStore
	menu                          store.MenuStore
	menuItem                      store.MenuItemStore
	menuItemTranslation           store.MenuItemTranslationStore
	openExchangeRate              store.OpenExchangeRateStore
	order                         store.OrderStore
	orderDiscount                 store.OrderDiscountStore
	orderEvent                    store.OrderEventStore
	orderLine                     store.OrderLineStore
	page                          store.PageStore
	pageTranslation               store.PageTranslationStore
	pageType                      store.PageTypeStore
	payment                       store.PaymentStore
	paymentTransaction            store.PaymentTransactionStore
	plugin                        store.PluginStore
	pluginConfiguration           store.PluginConfigurationStore
	preference                    store.PreferenceStore
	preorderAllocation            store.PreorderAllocationStore
	product                       store.ProductStore
	productChannelListing         store.ProductChannelListingStore
	productMedia                  store.ProductMediaStore
	productTranslation            store.ProductTranslationStore
	productType                   store.ProductTypeStore
	productVariant                store.ProductVariantStore
	productVariantChannelListing  store.ProductVariantChannelListingStore
	productVariantTranslation     store.ProductVariantTranslationStore
	role                          store.RoleStore
	session                       store.SessionStore
	shippingMethod                store.ShippingMethodStore
	shippingMethodChannelListing  store.ShippingMethodChannelListingStore
	shippingMethodPostalCodeRule  store.ShippingMethodPostalCodeRuleStore
	shippingMethodTranslation     store.ShippingMethodTranslationStore
	shippingZone                  store.ShippingZoneStore
	shopStaff                     store.ShopStaffStore
	shopTranslation               store.ShopTranslationStore
	staffNotificationRecipient    store.StaffNotificationRecipientStore
	status                        store.StatusStore
	stock                         store.StockStore
	system                        store.SystemStore
	termsOfService                store.TermsOfServiceStore
	token                         store.TokenStore
	uploadSession                 store.UploadSessionStore
	user                          store.UserStore
	userAccessToken               store.UserAccessTokenStore
	vat                           store.VatStore
	voucherChannelListing         store.VoucherChannelListingStore
	voucherCustomer               store.VoucherCustomerStore
	voucherTranslation            store.VoucherTranslationStore
	warehouse                     store.WarehouseStore
	wishlist                      store.WishlistStore
	wishlistItem                  store.WishlistItemStore
}

// setup tables before performing database migration
func (store *SqlStore) setupStores() {
	store.stores = &SqlStoreStores{
		address:                       account.NewSqlAddressStore(store),
		allocation:                    warehouse.NewSqlAllocationStore(store),
		app:                           app.NewSqlAppStore(store),
		appToken:                      app.NewSqlAppTokenStore(store),
		assignedPageAttribute:         attribute.NewSqlAssignedPageAttributeStore(store),
		assignedPageAttributeValue:    attribute.NewSqlAssignedPageAttributeValueStore(store),
		assignedProductAttribute:      attribute.NewSqlAssignedProductAttributeStore(store),
		assignedProductAttributeValue: attribute.NewSqlAssignedProductAttributeValueStore(store),
		attribute:                     attribute.NewSqlAttributeStore(store),
		attributePage:                 attribute.NewSqlAttributePageStore(store),
		attributeTranslation:          attribute.NewSqlAttributeTranslationStore(store),
		attributeValue:                attribute.NewSqlAttributeValueStore(store),
		attributeValueTranslation:     attribute.NewSqlAttributeValueTranslationStore(store),
		audit:                         audit.NewSqlAuditStore(store),
		category:                      product.NewSqlCategoryStore(store),
		categoryTranslation:           product.NewSqlCategoryTranslationStore(store),
		channel:                       channel.NewSqlChannelStore(store),
		checkout:                      checkout.NewSqlCheckoutStore(store),
		checkoutLine:                  checkout.NewSqlCheckoutLineStore(store),
		clusterDiscovery:              cluster.NewSqlClusterDiscoveryStore(store),
		collection:                    product.NewSqlCollectionStore(store),
		collectionChannelListing:      product.NewSqlCollectionChannelListingStore(store),
		collectionProduct:             product.NewSqlCollectionProductStore(store),
		collectionTranslation:         product.NewSqlCollectionTranslationStore(store),
		compliance:                    compliance.NewSqlComplianceStore(store),
		csvExportEvent:                csv.NewSqlCsvExportEventStore(store),
		csvExportFile:                 csv.NewSqlCsvExportFileStore(store),
		customProductAttribute:        attribute.NewSqlCustomProductAttributeStore(store),
		customerEvent:                 account.NewSqlCustomerEventStore(store),
		customerNote:                  account.NewSqlCustomerNoteStore(store),
		digitalContent:                product.NewSqlDigitalContentStore(store),
		digitalContentUrl:             product.NewSqlDigitalContentUrlStore(store),
		discountSale:                  discount.NewSqlDiscountSaleStore(store),
		discountSaleChannelListing:    discount.NewSqlDiscountSaleChannelListingStore(store),
		discountSaleTranslation:       discount.NewSqlDiscountSaleTranslationStore(store),
		discountVoucher:               discount.NewSqlDiscountVoucherStore(store),
		fileInfo:                      file.NewSqlFileInfoStore(store, store.metrics),
		fulfillment:                   order.NewSqlFulfillmentStore(store),
		fulfillmentLine:               order.NewSqlFulfillmentLineStore(store),
		giftCard:                      giftcard.NewSqlGiftCardStore(store),
		giftcardEvent:                 giftcard.NewSqlGiftcardEventStore(store),
		invoice:                       invoice.NewSqlInvoiceStore(store),
		invoiceEvent:                  invoice.NewSqlInvoiceEventStore(store),
		job:                           job.NewSqlJobStore(store),
		menu:                          menu.NewSqlMenuStore(store),
		menuItem:                      menu.NewSqlMenuItemStore(store),
		menuItemTranslation:           menu.NewSqlMenuItemTranslationStore(store),
		openExchangeRate:              external_services.NewSqlOpenExchangeRateStore(store),
		order:                         order.NewSqlOrderStore(store),
		orderDiscount:                 discount.NewSqlOrderDiscountStore(store),
		orderEvent:                    order.NewSqlOrderEventStore(store),
		orderLine:                     order.NewSqlOrderLineStore(store),
		page:                          page.NewSqlPageStore(store),
		pageTranslation:               page.NewSqlPageTranslationStore(store),
		pageType:                      page.NewSqlPageTypeStore(store),
		payment:                       payment.NewSqlPaymentStore(store),
		paymentTransaction:            payment.NewSqlPaymentTransactionStore(store),
		plugin:                        plugin.NewSqlPluginStore(store),
		pluginConfiguration:           plugin.NewSqlPluginConfigurationStore(store),
		preference:                    preference.NewSqlPreferenceStore(store),
		preorderAllocation:            warehouse.NewSqlPreorderAllocationStore(store),
		product:                       product.NewSqlProductStore(store),
		productChannelListing:         product.NewSqlProductChannelListingStore(store),
		productMedia:                  product.NewSqlProductMediaStore(store),
		productTranslation:            product.NewSqlProductTranslationStore(store),
		productType:                   product.NewSqlProductTypeStore(store),
		productVariant:                product.NewSqlProductVariantStore(store),
		productVariantChannelListing:  product.NewSqlProductVariantChannelListingStore(store),
		productVariantTranslation:     product.NewSqlProductVariantTranslationStore(store),
		role:                          account.NewSqlRoleStore(store),
		session:                       account.NewSqlSessionStore(store),
		shippingMethod:                shipping.NewSqlShippingMethodStore(store),
		shippingMethodChannelListing:  shipping.NewSqlShippingMethodChannelListingStore(store),
		shippingMethodPostalCodeRule:  shipping.NewSqlShippingMethodPostalCodeRuleStore(store),
		shippingMethodTranslation:     shipping.NewSqlShippingMethodTranslationStore(store),
		shippingZone:                  shipping.NewSqlShippingZoneStore(store),
		shopStaff:                     shop.NewSqlShopStaffStore(store),
		shopTranslation:               shop.NewSqlShopTranslationStore(store),
		staffNotificationRecipient:    account.NewSqlStaffNotificationRecipientStore(store),
		status:                        account.NewSqlStatusStore(store),
		stock:                         warehouse.NewSqlStockStore(store),
		system:                        system.NewSqlSystemStore(store),
		termsOfService:                account.NewSqlTermsOfServiceStore(store, store.metrics),
		token:                         account.NewSqlTokenStore(store),
		uploadSession:                 file.NewSqlUploadSessionStore(store),
		user:                          account.NewSqlUserStore(store, store.metrics),
		userAccessToken:               account.NewSqlUserAccessTokenStore(store),
		vat:                           shop.NewSqlVatStore(store),
		voucherChannelListing:         discount.NewSqlVoucherChannelListingStore(store),
		voucherCustomer:               discount.NewSqlVoucherCustomerStore(store),
		voucherTranslation:            discount.NewSqlVoucherTranslationStore(store),
		warehouse:                     warehouse.NewSqlWarehouseStore(store),
		wishlist:                      wishlist.NewSqlWishlistStore(store),
		wishlistItem:                  wishlist.NewSqlWishlistItemStore(store),
	}
}

func (ss *SqlStore) Address() store.AddressStore {
	return ss.stores.address
}

func (ss *SqlStore) Allocation() store.AllocationStore {
	return ss.stores.allocation
}

func (ss *SqlStore) App() store.AppStore {
	return ss.stores.app
}

func (ss *SqlStore) AppToken() store.AppTokenStore {
	return ss.stores.appToken
}

func (ss *SqlStore) AssignedPageAttribute() store.AssignedPageAttributeStore {
	return ss.stores.assignedPageAttribute
}

func (ss *SqlStore) AssignedPageAttributeValue() store.AssignedPageAttributeValueStore {
	return ss.stores.assignedPageAttributeValue
}

func (ss *SqlStore) AssignedProductAttribute() store.AssignedProductAttributeStore {
	return ss.stores.assignedProductAttribute
}

func (ss *SqlStore) AssignedProductAttributeValue() store.AssignedProductAttributeValueStore {
	return ss.stores.assignedProductAttributeValue
}

func (ss *SqlStore) Attribute() store.AttributeStore {
	return ss.stores.attribute
}

func (ss *SqlStore) AttributePage() store.AttributePageStore {
	return ss.stores.attributePage
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

func (ss *SqlStore) Audit() store.AuditStore {
	return ss.stores.audit
}

func (ss *SqlStore) Category() store.CategoryStore {
	return ss.stores.category
}

func (ss *SqlStore) CategoryTranslation() store.CategoryTranslationStore {
	return ss.stores.categoryTranslation
}

func (ss *SqlStore) Channel() store.ChannelStore {
	return ss.stores.channel
}

func (ss *SqlStore) Checkout() store.CheckoutStore {
	return ss.stores.checkout
}

func (ss *SqlStore) CheckoutLine() store.CheckoutLineStore {
	return ss.stores.checkoutLine
}

func (ss *SqlStore) ClusterDiscovery() store.ClusterDiscoveryStore {
	return ss.stores.clusterDiscovery
}

func (ss *SqlStore) Collection() store.CollectionStore {
	return ss.stores.collection
}

func (ss *SqlStore) CollectionChannelListing() store.CollectionChannelListingStore {
	return ss.stores.collectionChannelListing
}

func (ss *SqlStore) CollectionProduct() store.CollectionProductStore {
	return ss.stores.collectionProduct
}

func (ss *SqlStore) CollectionTranslation() store.CollectionTranslationStore {
	return ss.stores.collectionTranslation
}

func (ss *SqlStore) Compliance() store.ComplianceStore {
	return ss.stores.compliance
}

func (ss *SqlStore) CsvExportEvent() store.CsvExportEventStore {
	return ss.stores.csvExportEvent
}

func (ss *SqlStore) CsvExportFile() store.CsvExportFileStore {
	return ss.stores.csvExportFile
}

func (ss *SqlStore) CustomProductAttribute() store.CustomProductAttributeStore {
	return ss.stores.customProductAttribute
}

func (ss *SqlStore) CustomerEvent() store.CustomerEventStore {
	return ss.stores.customerEvent
}

func (ss *SqlStore) CustomerNote() store.CustomerNoteStore {
	return ss.stores.customerNote
}

func (ss *SqlStore) DigitalContent() store.DigitalContentStore {
	return ss.stores.digitalContent
}

func (ss *SqlStore) DigitalContentUrl() store.DigitalContentUrlStore {
	return ss.stores.digitalContentUrl
}

func (ss *SqlStore) DiscountSale() store.DiscountSaleStore {
	return ss.stores.discountSale
}

func (ss *SqlStore) DiscountSaleChannelListing() store.DiscountSaleChannelListingStore {
	return ss.stores.discountSaleChannelListing
}

func (ss *SqlStore) DiscountSaleTranslation() store.DiscountSaleTranslationStore {
	return ss.stores.discountSaleTranslation
}

func (ss *SqlStore) DiscountVoucher() store.DiscountVoucherStore {
	return ss.stores.discountVoucher
}

func (ss *SqlStore) FileInfo() store.FileInfoStore {
	return ss.stores.fileInfo
}

func (ss *SqlStore) Fulfillment() store.FulfillmentStore {
	return ss.stores.fulfillment
}

func (ss *SqlStore) FulfillmentLine() store.FulfillmentLineStore {
	return ss.stores.fulfillmentLine
}

func (ss *SqlStore) GiftCard() store.GiftCardStore {
	return ss.stores.giftCard
}

func (ss *SqlStore) GiftcardEvent() store.GiftcardEventStore {
	return ss.stores.giftcardEvent
}

func (ss *SqlStore) Invoice() store.InvoiceStore {
	return ss.stores.invoice
}

func (ss *SqlStore) InvoiceEvent() store.InvoiceEventStore {
	return ss.stores.invoiceEvent
}

func (ss *SqlStore) Job() store.JobStore {
	return ss.stores.job
}

func (ss *SqlStore) Menu() store.MenuStore {
	return ss.stores.menu
}

func (ss *SqlStore) MenuItem() store.MenuItemStore {
	return ss.stores.menuItem
}

func (ss *SqlStore) MenuItemTranslation() store.MenuItemTranslationStore {
	return ss.stores.menuItemTranslation
}

func (ss *SqlStore) OpenExchangeRate() store.OpenExchangeRateStore {
	return ss.stores.openExchangeRate
}

func (ss *SqlStore) Order() store.OrderStore {
	return ss.stores.order
}

func (ss *SqlStore) OrderDiscount() store.OrderDiscountStore {
	return ss.stores.orderDiscount
}

func (ss *SqlStore) OrderEvent() store.OrderEventStore {
	return ss.stores.orderEvent
}

func (ss *SqlStore) OrderLine() store.OrderLineStore {
	return ss.stores.orderLine
}

func (ss *SqlStore) Page() store.PageStore {
	return ss.stores.page
}

func (ss *SqlStore) PageTranslation() store.PageTranslationStore {
	return ss.stores.pageTranslation
}

func (ss *SqlStore) PageType() store.PageTypeStore {
	return ss.stores.pageType
}

func (ss *SqlStore) Payment() store.PaymentStore {
	return ss.stores.payment
}

func (ss *SqlStore) PaymentTransaction() store.PaymentTransactionStore {
	return ss.stores.paymentTransaction
}

func (ss *SqlStore) Plugin() store.PluginStore {
	return ss.stores.plugin
}

func (ss *SqlStore) PluginConfiguration() store.PluginConfigurationStore {
	return ss.stores.pluginConfiguration
}

func (ss *SqlStore) Preference() store.PreferenceStore {
	return ss.stores.preference
}

func (ss *SqlStore) PreorderAllocation() store.PreorderAllocationStore {
	return ss.stores.preorderAllocation
}

func (ss *SqlStore) Product() store.ProductStore {
	return ss.stores.product
}

func (ss *SqlStore) ProductChannelListing() store.ProductChannelListingStore {
	return ss.stores.productChannelListing
}

func (ss *SqlStore) ProductMedia() store.ProductMediaStore {
	return ss.stores.productMedia
}

func (ss *SqlStore) ProductTranslation() store.ProductTranslationStore {
	return ss.stores.productTranslation
}

func (ss *SqlStore) ProductType() store.ProductTypeStore {
	return ss.stores.productType
}

func (ss *SqlStore) ProductVariant() store.ProductVariantStore {
	return ss.stores.productVariant
}

func (ss *SqlStore) ProductVariantChannelListing() store.ProductVariantChannelListingStore {
	return ss.stores.productVariantChannelListing
}

func (ss *SqlStore) ProductVariantTranslation() store.ProductVariantTranslationStore {
	return ss.stores.productVariantTranslation
}

func (ss *SqlStore) Role() store.RoleStore {
	return ss.stores.role
}

func (ss *SqlStore) Session() store.SessionStore {
	return ss.stores.session
}

func (ss *SqlStore) ShippingMethod() store.ShippingMethodStore {
	return ss.stores.shippingMethod
}

func (ss *SqlStore) ShippingMethodChannelListing() store.ShippingMethodChannelListingStore {
	return ss.stores.shippingMethodChannelListing
}

func (ss *SqlStore) ShippingMethodPostalCodeRule() store.ShippingMethodPostalCodeRuleStore {
	return ss.stores.shippingMethodPostalCodeRule
}

func (ss *SqlStore) ShippingMethodTranslation() store.ShippingMethodTranslationStore {
	return ss.stores.shippingMethodTranslation
}

func (ss *SqlStore) ShippingZone() store.ShippingZoneStore {
	return ss.stores.shippingZone
}

func (ss *SqlStore) ShopStaff() store.ShopStaffStore {
	return ss.stores.shopStaff
}

func (ss *SqlStore) ShopTranslation() store.ShopTranslationStore {
	return ss.stores.shopTranslation
}

func (ss *SqlStore) StaffNotificationRecipient() store.StaffNotificationRecipientStore {
	return ss.stores.staffNotificationRecipient
}

func (ss *SqlStore) Status() store.StatusStore {
	return ss.stores.status
}

func (ss *SqlStore) Stock() store.StockStore {
	return ss.stores.stock
}

func (ss *SqlStore) System() store.SystemStore {
	return ss.stores.system
}

func (ss *SqlStore) TermsOfService() store.TermsOfServiceStore {
	return ss.stores.termsOfService
}

func (ss *SqlStore) Token() store.TokenStore {
	return ss.stores.token
}

func (ss *SqlStore) UploadSession() store.UploadSessionStore {
	return ss.stores.uploadSession
}

func (ss *SqlStore) User() store.UserStore {
	return ss.stores.user
}

func (ss *SqlStore) UserAccessToken() store.UserAccessTokenStore {
	return ss.stores.userAccessToken
}

func (ss *SqlStore) Vat() store.VatStore {
	return ss.stores.vat
}

func (ss *SqlStore) VoucherChannelListing() store.VoucherChannelListingStore {
	return ss.stores.voucherChannelListing
}

func (ss *SqlStore) VoucherCustomer() store.VoucherCustomerStore {
	return ss.stores.voucherCustomer
}

func (ss *SqlStore) VoucherTranslation() store.VoucherTranslationStore {
	return ss.stores.voucherTranslation
}

func (ss *SqlStore) Warehouse() store.WarehouseStore {
	return ss.stores.warehouse
}

func (ss *SqlStore) Wishlist() store.WishlistStore {
	return ss.stores.wishlist
}

func (ss *SqlStore) WishlistItem() store.WishlistItemStore {
	return ss.stores.wishlistItem
}
