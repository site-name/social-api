package sqlstore

import "github.com/sitename/sitename/store"

type SqlStoreStores struct {
	user                          store.UserStore                       // account models
	address                       store.AddressStore                    //
	audit                         store.AuditStore                      // common
	cluster                       store.ClusterDiscoveryStore           //
	session                       store.SessionStore                    //
	system                        store.SystemStore                     //
	preference                    store.PreferenceStore                 //
	token                         store.TokenStore                      //
	status                        store.StatusStore                     //
	job                           store.JobStore                        //
	userAccessToken               store.UserAccessTokenStore            //
	role                          store.RoleStore                       //
	TermsOfService                store.TermsOfServiceStore             //
	app                           store.AppStore                        //
	appToken                      store.AppTokenStore                   //
	channel                       store.ChannelStore                    // channel models
	checkout                      store.CheckoutStore                   // checkout models
	checkoutLine                  store.CheckoutLineStore               //
	csvExportEvent                store.CsvExportEventStore             // csv models
	discountVoucher               store.DiscountVoucherStore            // product and discount models
	discountVoucherChannelListing store.VoucherChannelListingStore      //
	discountVoucherTranslation    store.VoucherTranslationStore         //
	discountVoucherCustomer       store.DiscountVoucherCustomerStore    //
	discountSale                  store.DiscountSaleStore               //
	discountSaleTranslation       store.DiscountSaleTranslationStore    //
	discountSaleChannelListing    store.DiscountSaleChannelListingStore //
	orderDiscount                 store.OrderDiscountStore              //
	giftCard                      store.GiftCardStore                   // gift card models
	invoiceEvent                  store.InvoiceEventStore               // invoice models
	menu                          store.MenuStore                       // menu models
	menuItemTranslation           store.MenuItemTranslationStore        //
	order                         store.OrderStore                      // order models
	orderLine                     store.OrderLineStore                  //
	fulfillment                   store.FulfillmentStore                //
	fulfillmentLine               store.FulfillmentLineStore            //
	orderEvent                    store.OrderEventStore                 //
	page                          store.PageStore                       // page models
	pageType                      store.PageTypeStore                   //
	pageTranslation               store.PageTranslationStore            //
	payment                       store.PaymentStore                    // payment models
	transaction                   store.PaymentTransactionStore         //
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
	store.stores.TermsOfService.(SqlTermsOfServiceStore).createIndexesIfNotExists()
	store.stores.app.(*SqlAppStore).createIndexesIfNotExists()
	store.stores.appToken.(*SqlAppTokenStore).createIndexesIfNotExists()
	// channel
	store.stores.channel.(*SqlChannelStore).createIndexesIfNotExists()
	// checkout
	store.stores.checkout.(*SqlCheckoutStore).createIndexesIfNotExists()
	store.stores.checkoutLine.(*SqlCheckoutLineStore).createIndexesIfNotExists()
	// csv
	store.stores.csvExportEvent.(*SqlCsvExportEventStore).createIndexesIfNotExists()
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
