//go:generate go run layer_generators/main.go

package store

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/mattermost/gorp"
	"github.com/shopspring/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/app"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/model/audit"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/model/compliance"
	"github.com/sitename/sitename/model/csv"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/modules/mfa"
)

const (
	UUID_MAX_LENGTH = 36 // max length for all tables's Id fields, since google's uuid generates ids have length of 36
)

type StoreResult struct {
	Data interface{}
	NErr error // NErr a temporary field used by the new code for the AppError migration. This will later become Err when the entire store is migrated.
}

type Indexer interface {
	CreateIndexesIfNotExists() // CreateIndexesIfNotExists creates indexes for table fields
}

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
	DiscountVoucherCustomer() DiscountVoucherCustomerStore             //
	VoucherTranslation() VoucherTranslationStore                       //
	DiscountSale() DiscountSaleStore                                   //
	DiscountSaleTranslation() DiscountSaleTranslationStore             //
	DiscountSaleChannelListing() DiscountSaleChannelListingStore       //
	OrderDiscount() OrderDiscountStore                                 //
	GiftCard() GiftCardStore                                           // giftcard
	InvoiceEvent() InvoiceEventStore                                   // invoice
	Menu() MenuStore                                                   // menu
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
	Warehouse() WarehouseStore                                         // warehouse
	Stock() StockStore                                                 //
	Allocation() AllocationStore                                       //
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
	AssignedVariantAttributeValue() AssignedVariantAttributeValueStore //
	AssignedVariantAttribute() AssignedVariantAttributeStore           //
	AttributeVariant() AttributeVariantStore                           //
	AssignedProductAttributeValue() AssignedProductAttributeValueStore //
	AssignedProductAttribute() AssignedProductAttributeStore           //
	AttributeProduct() AttributeProductStore                           //
	FileInfo() FileInfoStore                                           //
	UploadSession() UploadSessionStore                                 // upload session
}

type UploadSessionStore interface {
	Indexer
	Save(session *model.UploadSession) (*model.UploadSession, error)
	Update(session *model.UploadSession) error
	Get(id string) (*model.UploadSession, error)
	GetForUser(userID string) ([]*model.UploadSession, error)
	Delete(id string) error
}

// fileinfo
type FileInfoStore interface {
	Indexer
	Save(info *model.FileInfo) (*model.FileInfo, error)
	Upsert(info *model.FileInfo) (*model.FileInfo, error)
	Get(id string) (*model.FileInfo, error)
	GetFromMaster(id string) (*model.FileInfo, error)
	GetByIds(ids []string) ([]*model.FileInfo, error)
	GetByPath(path string) (*model.FileInfo, error)
	GetForUser(userID string) ([]*model.FileInfo, error)
	GetWithOptions(page, perPage int, opt *model.GetFileInfosOptions) ([]*model.FileInfo, error)
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
		Indexer
		Save(attr *attribute.Attribute) (*attribute.Attribute, error)
		Get(id string) (*attribute.Attribute, error)
		GetAttributesByIds(ids []string) ([]*attribute.Attribute, error)
		GetProductAndVariantHeaders(ids []string) ([]string, error)
	}
	AttributeTranslationStore interface {
		Indexer
	}
	AttributeValueStore interface {
		Indexer
	}
	AttributeValueTranslationStore interface {
		Indexer
	}
	AssignedPageAttributeValueStore interface {
		Indexer
	}
	AssignedPageAttributeStore interface {
		Indexer
	}
	AttributePageStore interface {
		Indexer
	}
	AssignedVariantAttributeValueStore interface {
		Indexer
	}
	AssignedVariantAttributeStore interface {
		Indexer
	}
	AttributeVariantStore interface {
		Indexer
	}
	AssignedProductAttributeValueStore interface {
		Indexer
	}
	AssignedProductAttributeStore interface {
		Indexer
	}
	AttributeProductStore interface {
		Indexer
	}
)

// compliance
type ComplianceStore interface {
	Indexer
	Save(compliance *compliance.Compliance) (*compliance.Compliance, error)
	Update(compliance *compliance.Compliance) (*compliance.Compliance, error)
	Get(id string) (*compliance.Compliance, error)
	GetAll(offset, limit int) (compliance.Compliances, error)
	ComplianceExport(compliance *compliance.Compliance, cursor compliance.ComplianceExportCursor, limit int) ([]*compliance.CompliancePost, compliance.ComplianceExportCursor, error)
	MessageExport(cursor compliance.MessageExportCursor, limit int) ([]*compliance.MessageExport, compliance.MessageExportCursor, error)
}

//plugin
type PluginConfigurationStore interface {
	Indexer
}

// wishlist
type (
	WishlistStore interface {
		Indexer
	}
	WishlistItemStore interface {
		Indexer
	}
)

// warehouse
type (
	WarehouseStore interface {
		Indexer
		Save(wh *warehouse.WareHouse) (*warehouse.WareHouse, error)
		Get(id string) (*warehouse.WareHouse, error)
		GetWarehousesHeaders(ids []string) ([]string, error)
	}
	StockStore interface {
		Indexer
	}
	AllocationStore interface {
		Indexer
	}
)

// shipping
type (
	ShippingZoneStore interface {
		Indexer
	}
	ShippingMethodStore interface {
		Indexer
	}
	ShippingMethodPostalCodeRuleStore interface {
		Indexer
	}
	ShippingMethodChannelListingStore interface {
		Indexer
	}
	ShippingMethodTranslationStore interface {
		Indexer
	}
)

// product
type (
	CollectionTranslationStore interface {
		Indexer
	}
	CollectionChannelListingStore interface {
		Indexer
	}
	CollectionStore interface {
		Indexer
	}
	CollectionProductStore interface {
		Indexer
	}
	VariantMediaStore interface {
		Indexer
	}
	ProductMediaStore interface {
		Indexer
	}
	DigitalContentUrlStore interface {
		Indexer
	}
	DigitalContentStore interface {
		Indexer
	}
	ProductVariantChannelListingStore interface {
		Indexer
	}
	ProductVariantTranslationStore interface {
		Indexer
	}
	ProductVariantStore interface {
		Indexer
	}
	ProductChannelListingStore interface {
		Indexer
	}
	ProductTranslationStore interface {
		Indexer
	}
	ProductTypeStore interface {
		Indexer
	}
	CategoryTranslationStore interface {
		Indexer
	}
	CategoryStore interface {
		Indexer
	}
	ProductStore interface {
		Indexer
		Save(prd *product_and_discount.Product) (*product_and_discount.Product, error)
		Get(id string) (*product_and_discount.Product, error)
		GetProductsByIds(ids []string) ([]*product_and_discount.Product, error)
		// FilterProducts(filterInput *webmodel.ProductFilterInput) ([]*product_and_discount.Product, error)
	}
)

// payment
type (
	PaymentStore interface {
		Indexer
		Save(*payment.Payment) (*payment.Payment, error)                                        // Save save payment instance into database
		Get(string) (*payment.Payment, error)                                                   // Get returns a payment with given id
		GetPaymentsByOrderID(orderID string) ([]*payment.Payment, error)                        // GetPaymentsByOrderID returns all payments that belong to given order
		PaymentExistWithOptions(opts *payment.PaymentFilterOpts) (paymentExist bool, err error) // FilterWithOptions filter order's payments based on given options
	}
	PaymentTransactionStore interface {
		Indexer
		Save(*payment.PaymentTransaction) (*payment.PaymentTransaction, error)     // Save inserts new payment transaction into database
		Get(id string) (*payment.PaymentTransaction, error)                        // Get returns a payment transaction with given id
		GetAllByPaymentID(paymentID string) ([]*payment.PaymentTransaction, error) // GetAllByPaymentID returns a slice of payment transaction(s) that belong to given payment
	}
)

// page
type (
	PageTypeStore interface {
		Indexer
	}
	PageTranslationStore interface {
		Indexer
	}
	PageStore interface {
		Indexer
	}
)

// order
type (
	OrderLineStore interface {
		Indexer
		Save(*order.OrderLine) (*order.OrderLine, error)            // Save save given order line instance into database and returns it
		Get(id string) (*order.OrderLine, error)                    // Get returns a order line with id of given id
		GetAllByOrderID(orderID string) ([]*order.OrderLine, error) // GetAllByOrderID returns a slice of order lines that belong to given order
	}
	OrderStore interface {
		Indexer
		Save(*order.Order) (*order.Order, error)                             // Save insert an order into database and returns that order if success
		Get(id string) (*order.Order, error)                                 // Get find order in database with given id
		Update(order *order.Order) (*order.Order, error)                     // Update update order
		UpdateTotalPaid(orderId string, newTotalPaid *decimal.Decimal) error // updateTotalPaid update total paid amount of given order
	}
	OrderEventStore interface {
		Indexer
	}
	FulfillmentLineStore interface {
		Indexer
		Save(fulfillmentLine *order.FulfillmentLine) (*order.FulfillmentLine, error)
		Get(id string) (*order.FulfillmentLine, error)
	}
	FulfillmentStore interface {
		Indexer
		Save(fulfillment *order.Fulfillment) (*order.Fulfillment, error)
		Get(id string) (*order.Fulfillment, error)
		FilterByExcludeStatuses(orderID string, excludeStatuses []string) (exist bool, err error) // FilterByExcludeStatuses check if there is at least 1 fulfillment belong to given order and have status differnt than given statuses.
	}
)

type MenuItemTranslationStore interface {
	Indexer
}

type MenuStore interface {
	Indexer
}

type InvoiceEventStore interface {
	Indexer
}

type GiftCardStore interface {
	Indexer
}

type OrderDiscountStore interface {
	Indexer
}

type DiscountSaleTranslationStore interface {
	Indexer
}

type DiscountSaleChannelListingStore interface {
	Indexer
}

type DiscountSaleStore interface {
	Indexer
}

type VoucherTranslationStore interface {
	Indexer
}

type DiscountVoucherCustomerStore interface {
	Indexer
}

type VoucherChannelListingStore interface {
	Indexer
}

type DiscountVoucherStore interface {
	Indexer
}

// csv
type (
	CsvExportEventStore interface {
		Indexer
		Save(event *csv.ExportEvent) (*csv.ExportEvent, error)
	}
	CsvExportFileStore interface {
		Indexer
		Save(file *csv.ExportFile) (*csv.ExportFile, error)
		Get(id string) (*csv.ExportFile, error)
	}
)

// checkout
type (
	CheckoutLineStore interface {
		Indexer
	}
	CheckoutStore interface {
		Indexer
	}
)

// channel
type ChannelStore interface {
	Indexer
	Save(ch *channel.Channel) (*channel.Channel, error)
	Get(id string) (*channel.Channel, error)                                         // Get returns channel by given id
	GetBySlug(slug string) (*channel.Channel, error)                                 // GetBySlug returns channel by given slug
	GetChannelsByIdsAndOrder(ids []string, order string) ([]*channel.Channel, error) //
	GetRandomActiveChannel() (*channel.Channel, error)                               // GetRandomActiveChannel get an abitrary channel that is active
}

// app
type (
	AppTokenStore interface {
		Indexer
		Save(appToken *app.AppToken) (*app.AppToken, error)
	}

	AppStore interface {
		Indexer
		Save(app *app.App) (*app.App, error)
	}
)

type ClusterDiscoveryStore interface {
	Indexer
	Save(discovery *model.ClusterDiscovery) error
	Delete(discovery *model.ClusterDiscovery) (bool, error)
	Exists(discovery *model.ClusterDiscovery) (bool, error)
	GetAll(discoveryType, clusterName string) ([]*model.ClusterDiscovery, error)
	SetLastPingAt(discovery *model.ClusterDiscovery) error
	Cleanup() error
}

type AuditStore interface {
	Indexer
	Save(audit *audit.Audit) error
	Get(userID string, offset int, limit int) (audit.Audits, error)
	PermanentDeleteByUser(userID string) error
}

type TermsOfServiceStore interface {
	Indexer
	Save(termsOfService *model.TermsOfService) (*model.TermsOfService, error)
	GetLatest(allowFromCache bool) (*model.TermsOfService, error)
	Get(id string, allowFromCache bool) (*model.TermsOfService, error)
}

type PreferenceStore interface {
	Indexer
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
	Indexer
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
	Indexer
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
		Indexer                                                              // CreateIndexesIfNotExists creates indexes for table if needed
		Save(address *account.Address) (*account.Address, error)             // Save saves address into database
		Get(addressID string) (*account.Address, error)                      // Get returns an Address with given addressID is exist
		GetAddressesByIDs(addressesIDs []string) ([]*account.Address, error) // GetAddressesByIDs returns a slice of Addresses with given slice of id strings
		GetAddressesByUserID(userID string) ([]*account.Address, error)      // GetAddressesByUserID returns slice of addresses belong to given user
	}
	UserStore interface {
		mfa.Store                                                                     // for multifactor authentication
		Indexer                                                                       //
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
		InvalidateProfileCacheForUser(userID string) // NOTE: maybe need a look
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

		GetUnreadCount(userID string) (int64, error) // gonna be removed

		// GetTeamGroupUsers(teamID string) ([]*model.User, error)
		// GetProfileByGroupChannelIdsForUser(userID string, channelIds []string) (map[string][]*model.User, error)
		// GetEtagForProfilesNotInTeam(teamID string) string
		// GetChannelGroupUsers(channelID string) ([]*model.User, error)
		// GetUnreadCountForChannel(userID string, channelID string) (int64, error)
		// GetAnyUnreadPostCountForChannel(userID string, channelID string) (int64, error)
		// GetRecentlyActiveUsersForTeam(teamID string, offset, limit int, viewRestrictions *model.ViewUsersRestrictions) ([]*model.User, error)
		// GetNewUsersForTeam(teamID string, offset, limit int, viewRestrictions *model.ViewUsersRestrictions) ([]*model.User, error)
		// SearchNotInTeam(notInTeamID string, term string, options *model.UserSearchOptions) ([]*model.User, error)
		// SearchInChannel(channelID string, term string, options *model.UserSearchOptions) ([]*model.User, error)
		// SearchNotInChannel(teamID string, channelID string, term string, options *model.UserSearchOptions) ([]*model.User, error)
		// SearchWithoutTeam(term string, options *model.UserSearchOptions) ([]*model.User, error)
		// SearchInGroup(groupID string, term string, options *model.UserSearchOptions) ([]*model.User, error)
		// InvalidateProfilesInChannelCacheByUser(userID string)
		// InvalidateProfilesInChannelCache(channelID string)
		// GetProfilesInChannel(options *model.UserGetOptions) ([]*model.User, error)
		// GetProfilesInChannelByStatus(options *model.UserGetOptions) ([]*model.User, error)
		// GetAllProfilesInChannel(ctx context.Context, channelID string, allowFromCache bool) (map[string]*model.User, error)
		// GetProfilesNotInChannel(teamID string, channelId string, groupConstrained bool, offset int, limit int, viewRestrictions *model.ViewUsersRestrictions) ([]*model.User, error)
		// GetProfilesWithoutTeam(options *model.UserGetOptions) ([]*model.User, error)
		// GetProfiles(options *model.UserGetOptions) ([]*model.User, error)
		// AnalyticsActiveCount(time int64, options model.UserCountOptions) (int64, error)
		// GetProfilesNotInTeam(teamID string, groupConstrained bool, offset int, limit int, viewRestrictions *model.ViewUsersRestrictions) ([]*model.User, error)
		// AutocompleteUsersInChannel(teamID, channelID, term string, options *model.UserSearchOptions) (*model.UserAutocompleteInChannel, error)
	}
	TokenStore interface {
		Indexer
		Save(recovery *model.Token) error
		Delete(token string) error
		GetByToken(token string) (*model.Token, error)
		Cleanup()
		RemoveAllTokensByType(tokenType string) error
	}
	UserAccessTokenStore interface {
		Indexer
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
		Indexer
		Save(*account.UserAddress) (*account.UserAddress, error)
	}
	CustomerEventStore interface {
		Indexer
		Save(*account.CustomerEvent) (*account.CustomerEvent, error)
		Get(id string) (*account.CustomerEvent, error)
		Count() (int64, error)
		GetEventsByUserID(userID string) ([]*account.CustomerEvent, error) // get list of customer event belongs to given id
	}
	StaffNotificationRecipientStore interface {
		Indexer
		Save(*account.StaffNotificationRecipient) (*account.StaffNotificationRecipient, error)
		Get(string) (*account.StaffNotificationRecipient, error)
	}
	CustomerNoteStore interface {
		Indexer
	}
)

type SystemStore interface {
	Indexer
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
	Indexer
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
	UpdateLastActivityAt(sessionID string, time int64) error                    //
	UpdateRoles(userID string, roles string) (string, error)                    // UpdateRoles updates roles for all sessions that have userId of given userID,
	UpdateDeviceId(id string, deviceID string, expiresAt int64) (string, error) //
	UpdateProps(session *model.Session) error
	AnalyticsSessionCount() (int64, error)
	Cleanup(expiryTime int64, batchSize int64)
}

type RoleStore interface {
	Indexer
	Save(role *model.Role) (*model.Role, error)
	Get(roleID string) (*model.Role, error)
	GetAll() ([]*model.Role, error)
	GetByName(ctx context.Context, name string) (*model.Role, error)
	GetByNames(names []string) ([]*model.Role, error)
	Delete(roleID string) (*model.Role, error)
	PermanentDeleteAll() error
	ChannelHigherScopedPermissions(roleNames []string) (map[string]*model.RolePermissions, error)
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
