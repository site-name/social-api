package storetest

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/storetest/mocks"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var _ store.Store = (*Store)(nil)

type Store struct {
	context context.Context

	AddressStore         mocks.AddressStore
	UserStore            mocks.UserStore
	SessionStore         mocks.SessionStore
	UserAccessTokenStore mocks.UserAccessTokenStore
	TermsOfServiceStore  mocks.TermsOfServiceStore
	RoleStore            mocks.RoleStore

	AttributeStore             mocks.AttributeStore
	AttributeValueStore        mocks.AttributeValueStore
	AssignedPageAttributeStore mocks.AssignedPageAttributeStore

	AppStore      mocks.AppStore
	AppTokenStore mocks.AppTokenStore

	ChannelStore mocks.ChannelStore

	AllocationStore mocks.AllocationStore
	WarehouseStore  mocks.WarehouseStore
	StockStore      mocks.StockStore

	AuditStore            mocks.AuditStore
	ClusterDiscoveryStore mocks.ClusterDiscoveryStore
	ComplianceStore       mocks.ComplianceStore
	SystemStore           mocks.SystemStore
	PreferenceStore       mocks.PreferenceStore
	TokenStore            mocks.TokenStore
	StatusStore           mocks.StatusStore
	FileInfoStore         mocks.FileInfoStore
	UploadSessionStore    mocks.UploadSessionStore
	JobStore              mocks.JobStore
	PluginStore           mocks.PluginStore
}

func (s *Store) SetContext(ctx context.Context) { s.context = ctx }
func (s *Store) Context() context.Context       { return s.context }
func (s *Store) User() store.UserStore          { return &s.UserStore }
func (s *Store) Address() store.AddressStore    { return &s.AddressStore }
func (s *Store) Session() store.SessionStore    { return &s.SessionStore }

func (s *Store) Allocation() store.AllocationStore { return &s.AllocationStore }
func (s *Store) Warehouse() store.WarehouseStore   { return &s.WarehouseStore }
func (s *Store) Stock() store.StockStore           { return &s.StockStore }

func (s *Store) App() store.AppStore           { return &s.AppStore }
func (s *Store) AppToken() store.AppTokenStore { return &s.AppTokenStore }

func (*Store) AssignedPageAttribute() store.AssignedPageAttributeStore {
	panic("unimplemented")
}

func (*Store) AssignedPageAttributeValue() store.AssignedPageAttributeValueStore {
	panic("unimplemented")
}

func (*Store) AssignedProductAttribute() store.AssignedProductAttributeStore {
	panic("unimplemented")
}

func (*Store) AssignedProductAttributeValue() store.AssignedProductAttributeValueStore {
	panic("unimplemented")
}

func (*Store) AssignedVariantAttribute() store.AssignedVariantAttributeStore {
	panic("unimplemented")
}

func (*Store) AssignedVariantAttributeValue() store.AssignedVariantAttributeValueStore {
	panic("unimplemented")
}

func (*Store) Attribute() store.AttributeStore {
	panic("unimplemented")
}

func (*Store) AttributePage() store.AttributePageStore {
	panic("unimplemented")
}

func (*Store) AttributeProduct() store.AttributeProductStore {
	panic("unimplemented")
}

func (*Store) AttributeTranslation() store.AttributeTranslationStore {
	panic("unimplemented")
}

func (*Store) AttributeValue() store.AttributeValueStore {
	panic("unimplemented")
}

func (*Store) AttributeValueTranslation() store.AttributeValueTranslationStore {
	panic("unimplemented")
}

func (*Store) AttributeVariant() store.AttributeVariantStore {
	panic("unimplemented")
}

func (*Store) Audit() store.AuditStore {
	panic("unimplemented")
}

func (*Store) Category() store.CategoryStore {
	panic("unimplemented")
}

func (*Store) CategoryTranslation() store.CategoryTranslationStore {
	panic("unimplemented")
}

func (*Store) Channel() store.ChannelStore {
	panic("unimplemented")
}

func (*Store) CheckIntegrity() <-chan model.IntegrityCheckResult {
	return make(chan model.IntegrityCheckResult)
}

func (*Store) Checkout() store.CheckoutStore {
	panic("unimplemented")
}

func (*Store) CheckoutLine() store.CheckoutLineStore {
	panic("unimplemented")
}

func (*Store) Close() {}

func (*Store) ClusterDiscovery() store.ClusterDiscoveryStore {
	panic("unimplemented")
}

func (*Store) Collection() store.CollectionStore {
	panic("unimplemented")
}

func (*Store) CollectionChannelListing() store.CollectionChannelListingStore {
	panic("unimplemented")
}

func (*Store) CollectionProduct() store.CollectionProductStore {
	panic("unimplemented")
}

func (*Store) CollectionTranslation() store.CollectionTranslationStore {
	panic("unimplemented")
}

func (*Store) Compliance() store.ComplianceStore {
	panic("unimplemented")
}

func (*Store) CsvExportEvent() store.CsvExportEventStore {
	panic("unimplemented")
}

func (*Store) CsvExportFile() store.CsvExportFileStore {
	panic("unimplemented")
}

func (*Store) CustomerEvent() store.CustomerEventStore {
	panic("unimplemented")
}

func (*Store) CustomerNote() store.CustomerNoteStore {
	panic("unimplemented")
}

func (*Store) DBXFromContext(ctx context.Context) *gorm.DB {
	panic("unimplemented")
}

func (*Store) DigitalContent() store.DigitalContentStore {
	panic("unimplemented")
}

func (*Store) DigitalContentUrl() store.DigitalContentUrlStore {
	panic("unimplemented")
}

func (*Store) DiscountSale() store.DiscountSaleStore {
	panic("unimplemented")
}

func (*Store) DiscountSaleChannelListing() store.DiscountSaleChannelListingStore {
	panic("unimplemented")
}

func (*Store) DiscountSaleTranslation() store.DiscountSaleTranslationStore {
	panic("unimplemented")
}

func (*Store) DiscountVoucher() store.DiscountVoucherStore {
	panic("unimplemented")
}

func (*Store) DropAllTables() {}

func (*Store) FileInfo() store.FileInfoStore {
	panic("unimplemented")
}

func (*Store) Fulfillment() store.FulfillmentStore {
	panic("unimplemented")
}

func (*Store) FulfillmentLine() store.FulfillmentLineStore {
	panic("unimplemented")
}

func (*Store) GetDbVersion(numerical bool) (string, error) {
	return "", nil
}

func (*Store) GetMaster(noTimeout ...bool) *gorm.DB {
	panic("unimplemented")
}

func (*Store) GetQueryBuilder(placeholderFormats ...squirrel.PlaceholderFormat) squirrel.StatementBuilderType {
	panic("unimplemented")
}

// GetReplica implements store.Store.
func (*Store) GetReplica(noTimeout ...bool) *gorm.DB {
	panic("unimplemented")
}

// GiftCard implements store.Store.
func (*Store) GiftCard() store.GiftCardStore {
	panic("unimplemented")
}

// GiftcardEvent implements store.Store.
func (*Store) GiftcardEvent() store.GiftcardEventStore {
	panic("unimplemented")
}

// Invoice implements store.Store.
func (*Store) Invoice() store.InvoiceStore {
	panic("unimplemented")
}

// InvoiceEvent implements store.Store.
func (*Store) InvoiceEvent() store.InvoiceEventStore {
	panic("unimplemented")
}

// IsUniqueConstraintError implements store.Store.
func (*Store) IsUniqueConstraintError(err error, indexNames []string) bool {
	panic("unimplemented")
}

// Job implements store.Store.
func (*Store) Job() store.JobStore {
	panic("unimplemented")
}

func (*Store) LockToMaster() {}

func (*Store) MarkSystemRanUnitTests() {}

// Menu implements store.Store.
func (*Store) Menu() store.MenuStore {
	panic("unimplemented")
}

// MenuItem implements store.Store.
func (*Store) MenuItem() store.MenuItemStore {
	panic("unimplemented")
}

// MenuItemTranslation implements store.Store.
func (*Store) MenuItemTranslation() store.MenuItemTranslationStore {
	panic("unimplemented")
}

// OpenExchangeRate implements store.Store.
func (*Store) OpenExchangeRate() store.OpenExchangeRateStore {
	panic("unimplemented")
}

// Order implements store.Store.
func (*Store) Order() store.OrderStore {
	panic("unimplemented")
}

// OrderDiscount implements store.Store.
func (*Store) OrderDiscount() store.OrderDiscountStore {
	panic("unimplemented")
}

// OrderEvent implements store.Store.
func (*Store) OrderEvent() store.OrderEventStore {
	panic("unimplemented")
}

// OrderLine implements store.Store.
func (*Store) OrderLine() store.OrderLineStore {
	panic("unimplemented")
}

// Page implements store.Store.
func (*Store) Page() store.PageStore {
	panic("unimplemented")
}

// PageTranslation implements store.Store.
func (*Store) PageTranslation() store.PageTranslationStore {
	panic("unimplemented")
}

// PageType implements store.Store.
func (*Store) PageType() store.PageTypeStore {
	panic("unimplemented")
}

// Payment implements store.Store.
func (*Store) Payment() store.PaymentStore {
	panic("unimplemented")
}

// PaymentTransaction implements store.Store.
func (*Store) PaymentTransaction() store.PaymentTransactionStore {
	panic("unimplemented")
}

// Plugin implements store.Store.
func (*Store) Plugin() store.PluginStore {
	panic("unimplemented")
}

// PluginConfiguration implements store.Store.
func (*Store) PluginConfiguration() store.PluginConfigurationStore {
	panic("unimplemented")
}

// Preference implements store.Store.
func (*Store) Preference() store.PreferenceStore {
	panic("unimplemented")
}

// PreorderAllocation implements store.Store.
func (*Store) PreorderAllocation() store.PreorderAllocationStore {
	panic("unimplemented")
}

// Product implements store.Store.
func (*Store) Product() store.ProductStore {
	panic("unimplemented")
}

// ProductChannelListing implements store.Store.
func (*Store) ProductChannelListing() store.ProductChannelListingStore {
	panic("unimplemented")
}

// ProductMedia implements store.Store.
func (*Store) ProductMedia() store.ProductMediaStore {
	panic("unimplemented")
}

// ProductTranslation implements store.Store.
func (*Store) ProductTranslation() store.ProductTranslationStore {
	panic("unimplemented")
}

// ProductType implements store.Store.
func (*Store) ProductType() store.ProductTypeStore {
	panic("unimplemented")
}

// ProductVariant implements store.Store.
func (*Store) ProductVariant() store.ProductVariantStore {
	panic("unimplemented")
}

// ProductVariantChannelListing implements store.Store.
func (*Store) ProductVariantChannelListing() store.ProductVariantChannelListingStore {
	panic("unimplemented")
}

// ProductVariantTranslation implements store.Store.
func (*Store) ProductVariantTranslation() store.ProductVariantTranslationStore {
	panic("unimplemented")
}

func (*Store) ReplicaLagAbs() error  { return nil }
func (*Store) ReplicaLagTime() error { return nil }

// Role implements store.Store.
func (*Store) Role() store.RoleStore {
	panic("unimplemented")
}

// ShippingMethod implements store.Store.
func (*Store) ShippingMethod() store.ShippingMethodStore {
	panic("unimplemented")
}

// ShippingMethodChannelListing implements store.Store.
func (*Store) ShippingMethodChannelListing() store.ShippingMethodChannelListingStore {
	panic("unimplemented")
}

// ShippingMethodPostalCodeRule implements store.Store.
func (*Store) ShippingMethodPostalCodeRule() store.ShippingMethodPostalCodeRuleStore {
	panic("unimplemented")
}

// ShippingMethodTranslation implements store.Store.
func (*Store) ShippingMethodTranslation() store.ShippingMethodTranslationStore {
	panic("unimplemented")
}

// ShippingZone implements store.Store.
func (*Store) ShippingZone() store.ShippingZoneStore {
	panic("unimplemented")
}

// ShopStaff implements store.Store.
func (*Store) ShopStaff() store.ShopStaffStore {
	panic("unimplemented")
}

// ShopTranslation implements store.Store.
func (*Store) ShopTranslation() store.ShopTranslationStore {
	panic("unimplemented")
}

// StaffNotificationRecipient implements store.Store.
func (*Store) StaffNotificationRecipient() store.StaffNotificationRecipientStore {
	panic("unimplemented")
}

// Status implements store.Store.
func (*Store) Status() store.StatusStore {
	panic("unimplemented")
}

// System implements store.Store.
func (*Store) System() store.SystemStore {
	panic("unimplemented")
}

// TermsOfService implements store.Store.
func (*Store) TermsOfService() store.TermsOfServiceStore {
	panic("unimplemented")
}

// Token implements store.Store.
func (*Store) Token() store.TokenStore {
	panic("unimplemented")
}

func (*Store) UnlockFromMaster() {}

// UploadSession implements store.Store.
func (*Store) UploadSession() store.UploadSessionStore {
	panic("unimplemented")
}

// UserAccessToken implements store.Store.
func (*Store) UserAccessToken() store.UserAccessTokenStore {
	panic("unimplemented")
}

// Vat implements store.Store.
func (*Store) Vat() store.VatStore {
	panic("unimplemented")
}

// VoucherChannelListing implements store.Store.
func (*Store) VoucherChannelListing() store.VoucherChannelListingStore {
	panic("unimplemented")
}

// VoucherCustomer implements store.Store.
func (*Store) VoucherCustomer() store.VoucherCustomerStore {
	panic("unimplemented")
}

// VoucherTranslation implements store.Store.
func (*Store) VoucherTranslation() store.VoucherTranslationStore {
	panic("unimplemented")
}

// Wishlist implements store.Store.
func (*Store) Wishlist() store.WishlistStore {
	panic("unimplemented")
}

// WishlistItem implements store.Store.
func (*Store) WishlistItem() store.WishlistItemStore {
	panic("unimplemented")
}

func (s *Store) AssertExpectations(t mock.TestingT) bool {
	return mock.AssertExpectationsForObjects(t,
		&s.UserStore,
		&s.AddressStore,
	)
}
