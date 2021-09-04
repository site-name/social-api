package sub_app_iface

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/mattermost/gorp"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app/request"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/model/checkout"
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
	"github.com/sitename/sitename/modules/filestore"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/modules/plugin"
	"github.com/sitename/sitename/store"
)

// GiftCardApp defines methods for giftcard app
type GiftcardService interface {
	GetGiftCard(id string) (*giftcard.GiftCard, *model.AppError)                                                  // GetGiftCard returns a giftcard with given id
	GiftcardsByCheckout(checkoutToken string) ([]*giftcard.GiftCard, *model.AppError)                             // GiftcardsByCheckout returns all giftcards belong to given checkout
	PromoCodeIsGiftCard(code string) (bool, *model.AppError)                                                      // PromoCodeIsGiftCard checks whether there is giftcard with given code
	ToggleGiftcardStatus(giftCard *giftcard.GiftCard) *model.AppError                                             // ToggleGiftcardStatus set status of given giftcard to inactive/active
	RemoveGiftcardCodeFromCheckout(ckout *checkout.Checkout, giftcardCode string) *model.AppError                 // RemoveGiftcardCodeFromCheckout drops a relation between giftcard and checkout
	AddGiftcardCodeToCheckout(ckout *checkout.Checkout, promoCode string) *model.AppError                         // AddGiftcardCodeToCheckout adds giftcard data to checkout by code.
	CreateOrderGiftcardRelation(orderGiftCard *giftcard.OrderGiftCard) (*giftcard.OrderGiftCard, *model.AppError) // CreateOrderGiftcardRelation takes an order-giftcard relation instance then save it
	UpsertGiftcard(giftcard *giftcard.GiftCard) (*giftcard.GiftCard, *model.AppError)                             // UpsertGiftcard depends on given giftcard's Id to decide saves or updates it
}

// PaymentService defines methods for payment sub app
type PaymentService interface {
	PaymentsByOption(option *payment.PaymentFilterOption) ([]*payment.Payment, *model.AppError)                                                                                                                                                          // PaymentsByOption returns all payments that satisfy given option
	GetLastOrderPayment(orderID string) (*payment.Payment, *model.AppError)                                                                                                                                                                              // GetLastOrderPayment get most recent payment made for given order
	GetAllPaymentTransactions(paymentID string) ([]*payment.PaymentTransaction, *model.AppError)                                                                                                                                                         // GetAllPaymentTransactions returns all transactions belong to given payment
	GetLastPaymentTransaction(paymentID string) (*payment.PaymentTransaction, *model.AppError)                                                                                                                                                           // GetLastPaymentTransaction return most recent transaction made for given payment
	PaymentIsAuthorized(paymentID string) (bool, *model.AppError)                                                                                                                                                                                        // PaymentIsAuthorized checks if given payment is authorized
	PaymentGetAuthorizedAmount(pm *payment.Payment) (*goprices.Money, *model.AppError)                                                                                                                                                                   // PaymentGetAuthorizedAmount calculates authorized amount
	PaymentCanVoid(pm *payment.Payment) (bool, *model.AppError)                                                                                                                                                                                          // PaymentCanVoid check if payment can void
	CreatePaymentInformation(payment *payment.Payment, paymentToken *string, amount *decimal.Decimal, customerId *string, storeSource bool, additionalData map[string]string) (*payment.PaymentData, *model.AppError)                                    // Extract order information along with payment details. Returns information required to process payment and additional billing/shipping addresses for optional fraud-prevention mechanisms.
	GetAlreadyProcessedTransaction(paymentID string, gatewayResponse *payment.GatewayResponse) (*payment.PaymentTransaction, *model.AppError)                                                                                                            // GetAlreadyProcessedTransaction returns most recent processed transaction made for given payment
	SaveTransaction(transaction *payment.PaymentTransaction) (*payment.PaymentTransaction, *model.AppError)                                                                                                                                              // SaveTransaction save new payment transaction into database
	CreateTransaction(paymentID string, kind string, paymentInformation *payment.PaymentData, actionRequired bool, gatewayResponse *payment.GatewayResponse, errorMsg string, isSuccess bool) (*payment.PaymentTransaction, *model.AppError)             // CreatePaymentTransaction save new payment transaction into database and returns it
	GetAlreadyProcessedTransactionOrCreateNewTransaction(paymentID, kind string, paymentInformation *payment.PaymentData, actionRequired bool, gatewayResponse *payment.GatewayResponse, errorMsg string) (*payment.PaymentTransaction, *model.AppError) // GetAlreadyProcessedTransactionOrCreateNewTransaction either create new transaction or get already processed transaction
	GetPaymentToken(payMent *payment.Payment) (string, *payment.PaymentError, *model.AppError)                                                                                                                                                           // get first transaction that belongs to given payment and has kind of "auth", IsSuccess is true
	GatewayPostProcess(transaction *payment.PaymentTransaction, payment *payment.Payment) *model.AppError                                                                                                                                                // GatewayPostProcess
	UpdateTransaction(transaction *payment.PaymentTransaction) (*payment.PaymentTransaction, *model.AppError)                                                                                                                                            // UpdateTransaction updates given transaction and returns updated on
	CreateOrUpdatePayment(pm *payment.Payment) (*payment.Payment, *model.AppError)                                                                                                                                                                       // CreateOrUpdatePayment depends on whether given payment's Id is set or not to decide to update/save payment
	UpdatePayment(pm *payment.Payment, gatewayResponse *payment.GatewayResponse) *model.AppError                                                                                                                                                         // UpdatePayment updates given payment based on given `gatewayResponse`
	StoreCustomerId(userID string, gateway string, customerID string) *model.AppError                                                                                                                                                                    // StoreCustomerId process
	GetSubTotal(orderLines []*order.OrderLine, fallbackCurrency string) (*goprices.TaxedMoney, *model.AppError)                                                                                                                                          // GetSubTotal adds up all Total prices of given order lines
	// CreatePayment Create a payment instance.
	//
	// This method is responsible for creating payment instances that works for
	// both Django views and GraphQL mutations.
	//
	// NOTE: `customerIpAddress`, `paymentToken`, `returnUrl` and `externalReference` can be empty
	//
	// `extraData`, `ckout`, `ord` can be nil
	CreatePayment(gateway string, total *decimal.Decimal, currency string, email string, customerIpAddress string, paymentToken string, extraData map[string]string, ckout *checkout.Checkout, ord *order.Order, returnUrl string, externalReference string) (*payment.Payment, *payment.PaymentError, *model.AppError)
	CleanAuthorize(payMent *payment.Payment) *payment.PaymentError                   // CleanAuthorize Check if payment can be authorized
	CleanCapture(pm *payment.Payment, amount decimal.Decimal) *payment.PaymentError  // CleanCapture Check if payment can be captured.
	FetchCustomerId(user *account.User, gateway string) (string, *model.AppError)    // FetchCustomerId Retrieve users customer_id stored for desired gateway.
	ValidateGatewayResponse(response *payment.GatewayResponse) *payment.GatewayError // ValidateGatewayResponse Validate response to be a correct format for Saleor to process.
}

// CheckoutService
type CheckoutService interface {
	CheckoutByOption(option *checkout.CheckoutFilterOption) (*checkout.Checkout, *model.AppError)                                                                                                                                            // CheckoutByOption returns a checkout filtered by given option
	CheckoutsByOption(option *checkout.CheckoutFilterOption) ([]*checkout.Checkout, *model.AppError)                                                                                                                                         // CheckoutsByOption returns a list of checkouts, filtered by given option
	FetchCheckoutLines(checkout *checkout.Checkout) ([]*checkout.CheckoutLineInfo, *model.AppError)                                                                                                                                          // CheckoutLineShippingRequired(checkoutLine *checkout.CheckoutLine) (bool, *model.AppError) // CheckoutLineShippingRequired check if given checkout line's product variant requires shipping
	CheckVariantInStock(ckout *checkout.Checkout, variant *product_and_discount.ProductVariant, channelSlug string, quantity int, replace, checkQuantity bool) (int, *checkout.CheckoutLine, *model.AppError)                                // CheckVariantInStock checks if given variant is already in stock
	CheckoutShippingRequired(checkoutToken string) (bool, *model.AppError)                                                                                                                                                                   // CheckoutShippingRequired checks if given checkout requires shipping
	CheckoutCountry(checkout *checkout.Checkout) (string, *model.AppError)                                                                                                                                                                   // CheckoutCountry returns country code for given checkout
	CheckoutSetCountry(checkout *checkout.Checkout, newCountryCode string) *model.AppError                                                                                                                                                   // CheckoutSetCountry set new country code for checkout
	UpsertCheckout(ckout *checkout.Checkout) (*checkout.Checkout, *model.AppError)                                                                                                                                                           // UpsertCheckout updates or inserts given checkout and returns it
	GetCustomerEmail(checkout *checkout.Checkout) (string, *model.AppError)                                                                                                                                                                  // GetCustomerEmail returns either checkout owner's email or checkout's Email property
	CheckoutTotalGiftCardsBalance(checkout *checkout.Checkout) (*goprices.Money, *model.AppError)                                                                                                                                            // CheckoutTotalGiftCardsBalance returns giftcards balance money
	CheckoutLineWithVariant(checkout *checkout.Checkout, productVariantID string) (*checkout.CheckoutLine, *model.AppError)                                                                                                                  // CheckoutLineWithVariant return a checkout line of given checkout, that checkout line has VariantID of given product variant id
	AddVariantToCheckout(checkoutInfo *checkout.CheckoutInfo, variant *product_and_discount.ProductVariant, quantity int, replace bool, checkQuantity bool) (*checkout.Checkout, *model.AppError)                                            // AddVariantToCheckout adds a product variant to given checkout. If `replace`, any previous quantity is discarded instead of added to
	CheckoutLinesByCheckoutToken(checkoutToken string) ([]*checkout.CheckoutLine, *model.AppError)                                                                                                                                           // CheckoutLinesByCheckoutToken finds checkout lines that belong to given checkout
	DeleteCheckoutLines(checkoutLineIDs []string) *model.AppError                                                                                                                                                                            // DeleteCheckoutLines deletes all checkout lines by given uuid list
	UpsertCheckoutLine(checkoutLine *checkout.CheckoutLine) (*checkout.CheckoutLine, *model.AppError)                                                                                                                                        // Upsert creates or updates given checkout line and returns it with an error
	CalculateCheckoutQuantity(lineInfos []*checkout.CheckoutLineInfo) (int, *model.AppError)                                                                                                                                                 // CalculateCheckoutQuantity returns total sum of quantity of check out lines in given checkout infos
	AddVariantsToCheckout(ckout *checkout.Checkout, variants []*product_and_discount.ProductVariant, quantities []int, channelSlug string, skipStockCheck, replace bool) (*checkout.Checkout, *warehouse.InsufficientStock, *model.AppError) // AddVariantsToCheckout add variants to checkout
	ChangeBillingAddressInCheckout(ckout *checkout.Checkout, address *account.Address) *model.AppError                                                                                                                                       // ChangeBillingAddressInCheckout update billing address of given checkout
	CheckoutTotalWeight(checkoutLineInfos []*checkout.CheckoutLineInfo) (*measurement.Weight, *model.AppError)                                                                                                                               // CheckoutTotalWeight calculate total weight for given checkout lines (these lines belong to a single checkout)
}

// AccountService
type AccountService interface {
	AddressById(id string) (*account.Address, *model.AppError)                                                                                           // GetAddressById returns address with given id. If not found returns nil and concret error
	AddressesByOption(option *account.AddressFilterOption) ([]*account.Address, *model.AppError)                                                         // AddressesByOption returns a list of addresses by given option
	UserById(ctx context.Context, userID string) (*account.User, *model.AppError)                                                                        // GetUserById get user from database with given userId
	CustomerEventsByUser(userID string) ([]*account.CustomerEvent, *model.AppError)                                                                      // CustomerEventsByUser returns all customer event(s) belong to given user
	AddressesByUserId(userID string) ([]*account.Address, *model.AppError)                                                                               // AddressesByUserId returns list of address(es) (if found) that belong to given user
	UserSetDefaultAddress(userID, addressID, addressType string) (*account.User, *model.AppError)                                                        // UserSetDefaultAddress set given address to be default for given user
	AddressDeleteForUser(userID, addressID string) *model.AppError                                                                                       // AddressDeleteForUser deletes relationship between given user and address
	UserByEmail(email string) (*account.User, *model.AppError)                                                                                           // UserByEmail try finding user with given email and returns that user
	CommonCustomerCreateEvent(userID *string, orderID *string, eventType string, params model.StringInterface) (*account.CustomerEvent, *model.AppError) // CommonCustomerCreateEvent is common method for creating customer events
	// CreateUserFromSignup create new user with user input information by:
	//
	// 1) Checks if user signup is allowed
	//
	// 2) call to CreateUser
	//
	// 3) sends verification email to given email
	CreateUserFromSignup(c *request.Context, user *account.User, redirect string) (*account.User, *model.AppError)
	CreateUser(c *request.Context, user *account.User) (*account.User, *model.AppError) // CreateUser
	GetVerifyEmailToken(token string) (*model.Token, *model.AppError)
	VerifyEmailFromToken(userSuppliedTokenString string) *model.AppError
	IsUserSignUpAllowed() *model.AppError
	VerifyUserEmail(userID, email string) *model.AppError
	GetUserByUsername(username string) (*account.User, *model.AppError)
	IsUsernameTaken(name string) bool
	GetUserByAuth(authData *string, authService string) (*account.User, *model.AppError)
	GetUsers(options *account.UserGetOptions) ([]*account.User, *model.AppError)
	GenerateMfaSecret(userID string) (*model.MfaSecret, *model.AppError)
	DeactivateMfa(userID string) *model.AppError
	ActivateMfa(userID, token string) *model.AppError
	GetProfileImage(user *account.User) ([]byte, bool, *model.AppError)
	GetDefaultProfileImage(user *account.User) ([]byte, *model.AppError) // GetDefaultProfileImage generate user's default prifile image (first character of their first name)
	SetDefaultProfileImage(user *account.User) *model.AppError           // SetDefaultProfileImage sets default profile image for given user
	SetProfileImage(userID string, imageData *multipart.FileHeader) *model.AppError
	SetProfileImageFromMultiPartFile(userID string, f multipart.File) *model.AppError
	AdjustImage(file io.Reader) (*bytes.Buffer, *model.AppError)
	SetProfileImageFromFile(userID string, file io.Reader) *model.AppError
	UpdateActive(c *request.Context, user *account.User, active bool) (*account.User, *model.AppError)
	UpdateHashedPasswordByUserId(userID, newHashedPassword string) *model.AppError
	UpdateHashedPassword(user *account.User, newHashedPassword string) *model.AppError
	UpdateUserRolesWithUser(user *account.User, newRoles string, sendWebSocketEvent bool) (*account.User, *model.AppError)
	PermanentDeleteAllUsers(c *request.Context) *model.AppError
	UpdateUser(user *account.User, sendNotifications bool) (*account.User, *model.AppError)
	SendEmailVerification(user *account.User, newEmail, redirect string) *model.AppError
	GetStatus(userID string) (*account.Status, *model.AppError)
	GetStatusFromCache(userID string) *account.Status
	SearchUsers(props *account.UserSearch, options *account.UserSearchOptions) ([]*account.User, *model.AppError)
	PermanentDeleteUser(c *request.Context, user *account.User) *model.AppError
	UpdatePasswordAsUser(userID, currentPassword, newPassword string) *model.AppError
	UpdatePassword(user *account.User, newPassword string) *model.AppError
	UpdatePasswordSendEmail(user *account.User, newPassword, method string) *model.AppError
	UpdatePasswordByUserIdSendEmail(userID, newPassword, method string) *model.AppError
	GetPasswordRecoveryToken(token string) (*model.Token, *model.AppError)
	ResetPasswordFromToken(userSuppliedTokenString, newPassword string) *model.AppError
	GetUsersByIds(userIDs []string, options *store.UserGetByIdsOpts) ([]*account.User, *model.AppError)
	GetUsersByUsernames(usernames []string, asAdmin bool) ([]*account.User, *model.AppError)
	GetTotalUsersStats() (*account.UsersStats, *model.AppError)
	GetFilteredUsersStats(options *account.UserCountOptions) (*account.UsersStats, *model.AppError)
	UpdateUserRoles(userID string, newRoles string, sendWebSocketEvent bool) (*account.User, *model.AppError)
	SendPasswordReset(email string, siteURL string) (bool, *model.AppError)
	CheckProviderAttributes(user *account.User, patch *account.UserPatch) string
	CreatePasswordRecoveryToken(userID, email string) (*model.Token, *model.AppError)
	UpdateUserAsUser(user *account.User, asAdmin bool) (*account.User, *model.AppError)
	UpdateUserAuth(userID string, userAuth *account.UserAuth) (*account.UserAuth, *model.AppError)
	UpdateMfa(activate bool, userID, token string) *model.AppError
	GetUserTermsOfService(userID string) (*account.UserTermsOfService, *model.AppError)
	SaveUserTermsOfService(userID, termsOfServiceId string, accepted bool) *model.AppError
	DeleteToken(token *model.Token) *model.AppError
	IsFirstUserAccount() bool
	GetSanitizeOptions(asAdmin bool) map[string]bool
	SanitizeProfile(user *account.User, asAdmin bool)
	CreateUserAsAdmin(c *request.Context, user *account.User, redirect string) (*account.User, *model.AppError)
	CreateUserWithToken(c *request.Context, user *account.User, token *model.Token) (*account.User, *model.AppError)
	GetSession(token string) (*model.Session, *model.AppError)                         // GetSession finds a session with given token and returns it
	GetSessionById(sessionID string) (*model.Session, *model.AppError)                 // GetSessionById finds a sessionw ith given id and returns it
	AttachDeviceId(sessionID string, deviceID string, expiresAt int64) *model.AppError // AttachDeviceId add device id to given session and returns updated session
	GetCloudSession(token string) (*model.Session, *model.AppError)
	ReturnSessionToPool(session *model.Session)
	SessionHasPermissionTo(session *model.Session, permission *model.Permission) bool
	MakePermissionError(s *model.Session, permissions []*model.Permission) *model.AppError
	ExtendSessionExpiryIfNeeded(session *model.Session) bool
	AttachSessionCookies(c *request.Context, w http.ResponseWriter, r *http.Request)
	AuthenticateUserForLogin(c *request.Context, id, loginId, password, mfaToken, cwsToken string, ldapOnly bool) (user *account.User, err *model.AppError)
	DoLogin(c *request.Context, w http.ResponseWriter, r *http.Request, user *account.User, deviceID string, isMobile, isOAuthUser, isSaml bool) *model.AppError
	CheckForClientSideCert(r *http.Request) (string, string, string)
	HasPermissionTo(askingUserId string, permission *model.Permission) bool
	UpdateUserActive(c *request.Context, userID string, active bool) *model.AppError // UpdateUserActive updates given user's status to ...
	GetPreferencesForUser(userID string) (model.Preferences, *model.AppError)
	GetPreferenceByCategoryForUser(userID string, category string) (model.Preferences, *model.AppError)
	GetPreferenceByCategoryAndNameForUser(userID string, category string, preferenceName string) (*model.Preference, *model.AppError)
	UpdatePreferences(userID string, preferences model.Preferences) *model.AppError
	DeletePreferences(userID string, preferences model.Preferences) *model.AppError
	AddStatusCacheSkipClusterSend(status *account.Status)
	GetUserStatusesByIds(userIDs []string) ([]*account.Status, *model.AppError) // GetUserStatusesByIds tries getting statuses from cache, if any cache for an user not found, it finds in database
	AddStatusCache(status *account.Status)
	StatusByID(statusID string) (*account.Status, *model.AppError)
	StatusesByIDs(statusIDs []string) ([]*account.Status, *model.AppError)
	CreateUserAccessToken(token *account.UserAccessToken) (*account.UserAccessToken, *model.AppError)       // CreateUserAccessToken creates new user access token for user
	GetUserAccessToken(tokenID string, sanitize bool) (*account.UserAccessToken, *model.AppError)           // GetUserAccessToken get access token for user
	RevokeUserAccessToken(token *account.UserAccessToken) *model.AppError                                   // RevokeUserAccessToken
	SetStatusOnline(userID string, manual bool)                                                             // SetStatusOnline sets given user's status to online
	SetStatusOffline(userID string, manual bool)                                                            // SetStatusOffline sets user's status to offline
	DeleteAddresses(addressIDs []string) *model.AppError                                                    // DeleteAddress deletes given address and returns an error
	UpsertAddress(transaction *gorp.Transaction, addr *account.Address) (*account.Address, *model.AppError) // UpsertAddress inserts or updates given address by checking its Id attribute
	UserByOrderId(orderID string) (*account.User, *model.AppError)                                          // UserByOrderId returns an user who owns given order
	InvalidateCacheForUser(userID string)                                                                   // InvalidateCacheForUser invalidates cache for given user
	ClearAllUsersSessionCacheLocal()                                                                        // ClearAllUsersSessionCacheLocal purges current `*ServiceAccount` sessionCache
}

type ProductService interface {
	ProductVariantById(id string) (*product_and_discount.ProductVariant, *model.AppError)                                                                           // ProductVariantById returns a product variants with given id
	ProductTypesByCheckoutToken(checkoutToken string) ([]*product_and_discount.ProductType, *model.AppError)                                                        // ProductTypesByCheckoutToken returns all product types related to given checkout
	ProductById(productID string) (*product_and_discount.Product, *model.AppError)                                                                                  // ProductById returns a product with id of given id
	ProductChannelListingsByOption(option *product_and_discount.ProductChannelListingFilterOption) ([]*product_and_discount.ProductChannelListing, *model.AppError) // ProductChannelListingsByOption finds product channel listings by given options
	CollectionsByVoucherID(voucherID string) ([]*product_and_discount.Collection, *model.AppError)                                                                  // CollectionsByVoucherID finds all collections that have relationships with given voucher
	ProductsByVoucherID(voucherID string) ([]*product_and_discount.Product, *model.AppError)                                                                        // ProductsByVoucherID finds all products that have relationships with given voucher
	ProductsRequireShipping(productIDs []string) (bool, *model.AppError)                                                                                            // ProductsRequireShipping checks if at least 1 product require shipping, then return true, false otherwise
	// ProductVariantGetPrice
	ProductVariantGetPrice(product *product_and_discount.Product, collections []*product_and_discount.Collection, channel *channel.Channel, channelListing *product_and_discount.ProductVariantChannelListing, discounts []*product_and_discount.DiscountInfo) (*goprices.Money, *model.AppError)
	ProductVariantIsDigital(productVariantID string) (bool, *model.AppError)                                                                   // ProductVariantIsDigital finds product type that related to given product variant and check if that product type is digital and does not require shipping
	GetDefaultDigitalContentSettings(shop *shop.Shop) *shop.ShopDefaultDigitalContentSettings                                                  // GetDefaultDigitalContentSettings takes a shop and returns some setting of the shop
	CollectionsByProductID(productID string) ([]*product_and_discount.Collection, *model.AppError)                                             // CollectionsByProductID finds and returns all collections related to given product
	ProductVariantByOrderLineID(orderLineID string) (*product_and_discount.ProductVariant, *model.AppError)                                    // ProductVariantByOrderLineID returns a product variant by given order line id
	ProductVariantsByOption(option *product_and_discount.ProductVariantFilterOption) ([]*product_and_discount.ProductVariant, *model.AppError) // ProductVariantsByOption returns a list of product variants satisfy given option
	UpsertDigitalContentURL(contentURL *product_and_discount.DigitalContentUrl) (*product_and_discount.DigitalContentUrl, *model.AppError)     // UpsertDigitalContentURL create a digital content url then returns it
	ProductsByOption(option *product_and_discount.ProductFilterOption) ([]*product_and_discount.Product, *model.AppError)                      // ProductsByOption returns a list of products that satisfy given option
	ProductByOption(option *product_and_discount.ProductFilterOption) (*product_and_discount.Product, *model.AppError)                         // ProductByOption returns 1 product that satisfy given option
	ProductVariantGetWeight(productVariantID string) (*measurement.Weight, *model.AppError)                                                    // ProductVariantGetWeight returns weight of given product variant
	CategoriesByOption(option *product_and_discount.CategoryFilterOption) ([]*product_and_discount.Category, *model.AppError)                  // CategoriesByOption returns all categories that satisfy given option
	CategoryByOption(option *product_and_discount.CategoryFilterOption) (*product_and_discount.Category, *model.AppError)                      // CategoryByOption returns 1 category that satisfies given option
	DigitalContentbyOption(option *product_and_discount.DigitalContenetFilterOption) (*product_and_discount.DigitalContent, *model.AppError)   // DigitalContentbyOption returns 1 digital content filtered using given option
}

type WishlistService interface {
	UpsertWishlist(wishList *wishlist.Wishlist) (*wishlist.Wishlist, *model.AppError)                            // UpsertWishlist inserts a new wishlist instance into database with given userID
	WishlistByOption(option *wishlist.WishlistFilterOption) (*wishlist.Wishlist, *model.AppError)                // WishlistByOption returns 1 wishlist filtered by given option
	WishlistItemByOption(option *wishlist.WishlistItemFilterOption) (*wishlist.WishlistItem, *model.AppError)    // WishlistItemByOption returns 1 wishlist item filtered using given option
	WishlistItemsByOption(option *wishlist.WishlistItemFilterOption) ([]*wishlist.WishlistItem, *model.AppError) // WishlistItemsByOption returns a slice of wishlist items filtered using given option
}

type AttributeService interface {
	AttributeValuesOfAttribute(attributeID string) ([]*attribute.AttributeValue, *model.AppError) // AttributeValuesOfAttribute finds all attribute values of given attribute, it may return an app-error indicates error occured. returned error could be either (*store.ErrNotFound or system error)
	// AssociateAttributeValuesToInstance assigns given attribute values to a product or variant.
	//
	// `instance` must be either `*product.Product` or `*product.ProductVariant` or `*page.Page`
	//
	// `attributeID` must be ID of processing `Attribute`
	//
	// Returned interface{} must be either: `*AssignedProductAttribute` or `*AssignedVariantAttribute` or `*AssignedPageAttribute`
	AssociateAttributeValuesToInstance(instance interface{}, attributeID string, values []*attribute.AttributeValue) (interface{}, *model.AppError)
	AttributeByID(id string) (*attribute.Attribute, *model.AppError)                                                                         // AttributeByID finds attribute with given id
	AttributeBySlug(slug string) (*attribute.Attribute, *model.AppError)                                                                     // AttributeBySlug finds an attribute with given slug
	AssignedVariantAttributesByOption(option *attribute.AssignedVariantAttributeFilterOption) ([]*attribute.AssignedVariantAttribute, error) // AssignedVariantAttributesByOption returns a list of assigned variant attributes filtered by given options
	AttributesByOption(option *attribute.AttributeFilterOption) ([]*attribute.Attribute, *model.AppError)                                    // AttributesByOption returns a list of attributes filtered using given options
}

type InvoiceService interface {
}

type ChannelService interface {
	GetChannelBySlug(slug string) (*channel.Channel, *model.AppError) // GetChannelBySlug returns a channel (if found) from database with given slug
	GetDefaultActiveChannel() (*channel.Channel, *model.AppError)     // GetDefaultChannel get random channel that is active
	// CleanChannel performs:
	//
	// 1) If given slug is not nil, try getting a channel with that slug.
	//   +) if found, check if channel is active
	//
	// 2) If given slug if nil, it try
	CleanChannel(channelSlug *string) (*channel.Channel, *model.AppError)
	ValidateChannel(channelSlug string) (*channel.Channel, *model.AppError)                     // ValidateChannel finds a channel with given slug, then check if the channel is active. If no channel found or found but not active, return an error
	GetDefaultChannelSlugOrGraphqlError() (string, *model.AppError)                             // GetDefaultChannelSlugOrGraphqlError returns a default channel slug
	ChannelsByOption(option *channel.ChannelFilterOption) ([]*channel.Channel, *model.AppError) // ChannelsByOption returns a list of channels by given options
	ChannelByOption(option *channel.ChannelFilterOption) (*channel.Channel, *model.AppError)    // ChannelByOption returns a channel that satisfies given options
}

type WarehouseService interface {
	CheckStockQuantity(variant *product_and_discount.ProductVariant, countryCode string, channelSlug string, quantity int) (*warehouse.InsufficientStock, *model.AppError)          // Validate if there is stock available for given variant in given country. If so - returns None. If there is less stock then required raise InsufficientStock exception.
	CheckStockQuantityBulk(variants product_and_discount.ProductVariants, countryCode string, quantities []int, channelSlug string) (*warehouse.InsufficientStock, *model.AppError) // Validate if there is stock available for given variants in given country. It raises InsufficientStock: when there is not enough items in stock for a variant
	IsProductInStock(productID string, countryCode string, channelSlug string) (bool, *model.AppError)                                                                              // IsProductInStock
	GetOrderLinesWithTrackInventory(orderLineInfos []*order.OrderLineData) []*order.OrderLineData                                                                                   // GetOrderLinesWithTrackInventory Return order lines with variants with track inventory set to True
	DecreaseAllocations(lineInfos []*order.OrderLineData) (*warehouse.InsufficientStock, *model.AppError)                                                                           // DecreaseAllocations Decreate allocations for provided order lines.
	WarehousesByOption(option *warehouse.WarehouseFilterOption) ([]*warehouse.WareHouse, *model.AppError)                                                                           // WarehouseByOption returns a list of warehouses based on given option
	AllocationsByOption(transaction *gorp.Transaction, option *warehouse.AllocationFilterOption) ([]*warehouse.Allocation, *model.AppError)                                         // AllocationsByOption returns all warehouse allocations filtered based on given option
	WarehouseByStockID(stockID string) (*warehouse.WareHouse, *model.AppError)                                                                                                      // WarehouseByStockID returns a warehouse that owns the given stock
	// IncreaseStock Increse stock quantity for given `order_line` in a given warehouse.
	//
	// Function lock for update stock and allocations related to given `order_line`
	// in a given warehouse. If the stock exists, increase the stock quantity
	// by given value. If not exist create a stock with the given quantity. This function
	// can create the allocation for increased quantity in stock by passing True
	// to `allocate` argument. If the order line has the allocation in this stock
	// function increase `quantity_allocated`. If allocation does not exist function
	// create a new allocation for this order line in this stock.
	//
	// NOTE: allocate is default to false
	IncreaseStock(orderLine *order.OrderLine, warehouse *warehouse.WareHouse, quantity int, allocate bool) *model.AppError
	// DeallocateStock Deallocate stocks for given `order_lines`.
	//
	// Function lock for update stocks and allocations related to given `order_lines`.
	// Iterate over allocations sorted by `stock.pk` and deallocate as many items
	// as needed of available in stock for order line, until deallocated all required
	// quantity for the order line. If there is less quantity in stocks then
	// raise an exception.
	DeallocateStock(orderLineDatas []*order.OrderLineData) (*warehouse.AllocationError, *model.AppError)
	// Decrease stocks quantities for given `order_lines` in given warehouses.
	//
	// Function deallocate as many quantities as requested if order_line has less quantity
	// from requested function deallocate whole quantity. Next function try to find the
	// stock in a given warehouse, if stock not exists or have not enough stock,
	// the function raise InsufficientStock exception. When the stock has enough quantity
	// function decrease it by given value.
	// If update_stocks is False, allocations will decrease but stocks quantities
	// will stay unmodified (case of unconfirmed order editing).
	//
	// updateStocks default to true
	DecreaseStock(orderLineInfos []*order.OrderLineData, updateStocks bool) (*warehouse.InsufficientStock, *model.AppError)
	IncreaseAllocations(lineInfos []*order.OrderLineData, channelSlug string) (*warehouse.InsufficientStock, *model.AppError) // IncreaseAllocations ncrease allocation for order lines with appropriate quantity
	GetStockById(stockID string) (*warehouse.Stock, *model.AppError)                                                          // GetStockById takes options for filtering 1 stock
	FilterStocksForChannel(option *warehouse.StockFilterForChannelOption) ([]*warehouse.Stock, *model.AppError)               // FilterStocksForChannel returns a slice of stocks that filtered using given options
}

type DiscountService interface {
	VouchersByOption(option *product_and_discount.VoucherFilterOption) ([]*product_and_discount.Voucher, *model.AppError)                         // VouchersByOption finds all vouchers with given option then returns them
	ValidateMinSpent(voucher *product_and_discount.Voucher, value *goprices.TaxedMoney, channelID string) (*model.NotApplicable, *model.AppError) // ValidateMinSpent validates if the order cost at least a specific amount of money
	ValidateOncePerCustomer(voucher *product_and_discount.Voucher, customerEmail string) (*model.NotApplicable, *model.AppError)                  // ValidateOncePerCustomer checks to make sure each customer has ONLY 1 time usage with 1 voucher
	// GetDiscountAmountFor checks given voucher's `DiscountValueType` and returns according discount calculator function
	//
	//  price.(type) == *Money || *MoneyRange || *TaxedMoney || *TaxedMoneyRange
	//
	// NOTE: the returning interface's type should be identical to given price's type
	GetDiscountAmountFor(voucher *product_and_discount.Voucher, price interface{}, channelID string) (interface{}, *model.AppError)
	// FilterSalesByOption should be used to filter active or expired sales
	// refer: saleor/discount/models.SaleQueryset for details
	FilterSalesByOption(option *product_and_discount.SaleFilterOption) ([]*product_and_discount.Sale, *model.AppError)
	PromoCodeIsVoucher(code string) (bool, *model.AppError)                                                                       // PromoCodeIsVoucher checks if given code is belong to a voucher
	ValidateVoucherOnlyForStaff(voucher *product_and_discount.Voucher, customerID string) (*model.NotApplicable, *model.AppError) // ValidateVoucherOnlyForStaff validate if voucher is only for staff
	// CalculateDiscountedPrice Return minimum product's price of all prices with discounts applied
	//
	// `discounts` is optional
	CalculateDiscountedPrice(product *product_and_discount.Product, price *goprices.Money, collections []*product_and_discount.Collection, discounts []*product_and_discount.DiscountInfo, channeL *channel.Channel) (*goprices.Money, *model.AppError)
	OrderDiscountsByOption(option *product_and_discount.OrderDiscountFilterOption) ([]*product_and_discount.OrderDiscount, *model.AppError)                      // OrderDiscountsByOption filters and returns order discounts with given option
	UpsertOrderDiscount(transaction *gorp.Transaction, orderDiscount *product_and_discount.OrderDiscount) (*product_and_discount.OrderDiscount, *model.AppError) // UpsertOrderDiscount updates or inserts given order discount
	ValidateVoucherInOrder(ord *order.Order) (*model.NotApplicable, *model.AppError)                                                                             // ValidateVoucherInOrder validates order has voucher and the voucher satisfies all requirements
	VoucherById(voucherID string) (*product_and_discount.Voucher, *model.AppError)                                                                               // VoucherById finds and returns a voucher with given id
	GetProductsVoucherDiscount(voucher *product_and_discount.Voucher, prices []*goprices.Money, channelID string) (*goprices.Money, *model.AppError)             // GetProductsVoucherDiscount Calculate discount value for a voucher of product or category type
	BulkDeleteOrderDiscounts(orderDiscountIDs []string) *model.AppError                                                                                          // BulkDeleteOrderDiscounts performs bulk delete given order discounts
	FetchActiveDiscounts() ([]*product_and_discount.DiscountInfo, *model.AppError)                                                                               // FetchActiveDiscounts returns discounts that are activated
}

type OrderService interface {
	// OrderShippingIsRequired checks if an order requires ship or not by:
	//
	// 1) Find all child order lines that belong to given order
	//
	// 2) iterates over resulting slice to check if at least one order line requires shipping
	OrderShippingIsRequired(orderID string) (bool, *model.AppError)
	OrderTotalQuantity(orderID string) (int, *model.AppError)                                                      // OrderTotalQuantity return total quantity of given order
	UpdateOrderTotalPaid(transaction *gorp.Transaction, orderID string) *model.AppError                            // UpdateOrderTotalPaid update given order's total paid amount
	OrderIsPreAuthorized(orderID string) (bool, *model.AppError)                                                   // OrderIsPreAuthorized checks if order is pre-authorized
	OrderIsCaptured(orderID string) (bool, *model.AppError)                                                        // OrderIsCaptured checks if given order is captured
	OrderSubTotal(order *order.Order) (*goprices.TaxedMoney, *model.AppError)                                      // OrderSubTotal returns sum of TotalPrice of all order lines that belong to given order
	OrderCanCancel(ord *order.Order) (bool, *model.AppError)                                                       // OrderCanCalcel checks if given order can be canceled
	OrderCanCapture(ord *order.Order, payment *payment.Payment) (bool, *model.AppError)                            // OrderCanCapture checks if given order can capture.
	OrderCanVoid(ord *order.Order, payment *payment.Payment) (bool, *model.AppError)                               // OrderCanVoid checks if given order can void
	OrderCanRefund(ord *order.Order, payment *payment.Payment) (bool, *model.AppError)                             // OrderCanRefund checks if order can refund
	CanMarkOrderAsPaid(ord *order.Order, payments []*payment.Payment) (bool, *model.AppError)                      // CanMarkOrderAsPaid checks if given order can be marked as paid.
	OrderTotalAuthorized(ord *order.Order) (*goprices.Money, *model.AppError)                                      // OrderTotalAuthorized returns order's total authorized amount
	GetOrderCountryCode(ord *order.Order) (string, *model.AppError)                                                // GetOrderCountryCode is helper function, returns contry code of given order
	OrderLineById(id string) (*order.OrderLine, *model.AppError)                                                   // OrderLineById returns order line with id of given id
	OrderById(id string) (*order.Order, *model.AppError)                                                           // OrderById returns order with id of given id
	CustomerEmail(ord *order.Order) (string, *model.AppError)                                                      // CustomerEmail try finding order's owner's email. If order has no user or error occured during the finding process, returns order's UserEmail property instead
	OrderLinesByOption(option *order.OrderLineFilterOption) ([]*order.OrderLine, *model.AppError)                  // OrderLinesByOption returns a list of order lines by given option
	AnAddressOfOrder(orderID string, whichAddressID order.WhichOrderAddressID) (*account.Address, *model.AppError) // AnAddressOfOrder returns shipping address of given order if presents
	OrderLineIsDigital(orderLine *order.OrderLine) (bool, *model.AppError)                                         // OrderLineIsDigital Check if a variant is digital and contains digital content.
}

type MenuService interface {
	MenuById(id string) (*menu.Menu, *model.AppError)     // MenuById returns a menu with given id
	MenuByName(name string) (*menu.Menu, *model.AppError) // MenuByName returns a menu with given name
	MenuBySlug(slug string) (*menu.Menu, *model.AppError) // MenuBySlug returns a menu with given slug
}

type AppApp interface {
}

type CsvService interface {
}

type ShopService interface {
	ShopById(shopID string) (*shop.Shop, *model.AppError)                                                         // ShopById finds shop by given id
	ShopStaffRelationByShopIDAndStaffID(shopID string, staffID string) (*shop.ShopStaffRelation, *model.AppError) // ShopStaffRelationByShopIDAndStaffID finds a shop-staff relation and returns it
}

type ShippingService interface {
	ShippingMethodChannelListingsByOption(option *shipping.ShippingMethodChannelListingFilterOption) ([]*shipping.ShippingMethodChannelListing, *model.AppError)                                                  // ShippingMethodChannelListingsByOption returns a list of shipping method channel listings by given option
	ApplicableShippingMethodsForCheckout(ckout *checkout.Checkout, channelID string, price *goprices.Money, countryCode string, lines []*checkout.CheckoutLineInfo) ([]*shipping.ShippingMethod, *model.AppError) // ApplicableShippingMethodsForCheckout finds all applicable shipping methods for given checkout, based on given additional arguments
	ApplicableShippingMethodsForOrder(oder *order.Order, channelID string, price *goprices.Money, countryCode string, lines []*checkout.CheckoutLineInfo) ([]*shipping.ShippingMethod, *model.AppError)           // ApplicableShippingMethodsForOrder finds all applicable shippingmethods for given order, based on other arguments passed in
	DefaultShippingZoneExists(shippingZoneID string) ([]*shipping.ShippingZone, *model.AppError)                                                                                                                  // DefaultShippingZoneExists returns all shipping zones that have Ids differ than given shippingZoneID and has `Default` properties equal to true
	GetCountriesWithoutShippingZone() ([]string, *model.AppError)                                                                                                                                                 // GetCountriesWithoutShippingZone Returns country codes that are not assigned to any shipping zone.
	ShippingZonesByOption(option *shipping.ShippingZoneFilterOption) ([]*shipping.ShippingZone, *model.AppError)                                                                                                  // ShippingZonesByOption returns all shipping zones that satisfy given options
	ShippingMethodByOption(option *shipping.ShippingMethodFilterOption) (*shipping.ShippingMethod, *model.AppError)                                                                                               // ShippingMethodByOption returns a shipping method with given options
}

type WebhookService interface {
}

type PageService interface {
}

type SeoService interface {
}

type FileService interface {
	FileBackend() (filestore.FileBackend, *model.AppError)
	CheckMandatoryS3Fields(settings *model.FileSettings) *model.AppError
	TestFileStoreConnection() *model.AppError
	TestFileStoreConnectionWithConfig(settings *model.FileSettings) *model.AppError
	ReadFile(path string) ([]byte, *model.AppError)
	FileReader(path string) (filestore.ReadCloseSeeker, *model.AppError)
	FileExists(path string) (bool, *model.AppError)
	FileSize(path string) (int64, *model.AppError)
	FileModTime(path string) (time.Time, *model.AppError)
	MoveFile(oldPath, newPath string) *model.AppError
	WriteFile(fr io.Reader, path string) (int64, *model.AppError)
	AppendFile(fr io.Reader, path string) (int64, *model.AppError)
	RemoveFile(path string) *model.AppError
	ListDirectory(path string) ([]string, *model.AppError)
	RemoveDirectory(path string) *model.AppError
	GeneratePublicLink(siteURL string, info *file.FileInfo) string
	DoUploadFile(c *request.Context, now time.Time, rawTeamId string, rawChannelId string, rawUserId string, rawFilename string, data []byte) (*file.FileInfo, *model.AppError)
	DoUploadFileExpectModification(c *request.Context, now time.Time, rawTeamId string, rawChannelId string, rawUserId string, rawFilename string, data []byte) (*file.FileInfo, []byte, *model.AppError)
	HandleImages(previewPathList []string, thumbnailPathList []string, fileData [][]byte)
	GetFileInfos(page, perPage int, opt *file.GetFileInfosOptions) ([]*file.FileInfo, *model.AppError)
	GetFileInfo(fileID string) (*file.FileInfo, *model.AppError)
	GetFile(fileID string) ([]byte, *model.AppError)
	CopyFileInfos(userID string, fileIDs []string) ([]string, *model.AppError)
	CreateZipFileAndAddFiles(fileBackend filestore.FileBackend, fileDatas []model.FileData, zipFileName, directory string) error
	ExtractContentFromFileInfo(fileInfo *file.FileInfo) error
	GetUploadSessionsForUser(userID string) ([]*file.UploadSession, *model.AppError)
	UploadData(c *request.Context, us *file.UploadSession, rd io.Reader) (*file.FileInfo, *model.AppError)
	GetUploadSession(uploadId string) (*file.UploadSession, *model.AppError)
	DownloadFromURL(downloadURL string) ([]byte, error) // DownloadFromURL
}

type PluginService interface {
	// GetPluginsEnvironment returns the plugin environment for use if plugins are enabled and
	// initialized.
	//
	// To get the plugins environment when the plugins are disabled, manually acquire the plugins
	// lock instead.
	GetPluginsEnvironment() (*plugin.Environment, *model.AppError)                                                                 // GetPluginsEnvironment return plugin environment of Server
	SetPluginsEnvironment(pluginsEnvironment *plugin.Environment)                                                                  // SetPluginsEnvironment set plugins environment for server
	SyncPluginsActiveState()                                                                                                       // SyncPluginsActiveState
	ServeInterPluginRequest(w http.ResponseWriter, r *http.Request, sourcePluginId, destinationPluginId string)                    // ServeInterPluginRequest
	InitPlugins(c *request.Context, pluginDir, webappPluginDir string)                                                             // InitPlugins creates new plugin api
	GetPlugins() (*plugins.PluginsResponse, *model.AppError)                                                                       // GetPlugins returns active/inactive plugins
	GetMarketplacePlugins(filter *plugins.MarketplacePluginFilter) ([]*plugins.MarketplacePlugin, *model.AppError)                 // GetMarketplacePlugins returns a list of available plugins on marketplace
	EnablePlugin(id string) *model.AppError                                                                                        // EnablePlugin will set the config for an installed plugin to enabled, triggering asynchronous activation if inactive anywhere in the cluster. Notifies cluster peers through config change.
	DisablePlugin(id string) *model.AppError                                                                                       // DisablePlugin will set the config for an installed plugin to disabled, triggering deactivation if active. Notifies cluster peers through config change.
	RemovePlugin(id string) *model.AppError                                                                                        // RemovePlugin removes given plugin
	GetPluginStatus(id string) (*plugins.PluginStatus, *model.AppError)                                                            // GetPluginStatus returns status for given plugin
	InstallPlugin(pluginFile io.ReadSeeker, replace bool) (*plugins.Manifest, *model.AppError)                                     // InstallPlugin installs plugins from given file
	SetPluginKeyWithOptions(pluginID string, key string, value []byte, options plugins.PluginKVSetOptions) (bool, *model.AppError) // SetPluginKeyWithOptions
	SetPluginKey(pluginID string, key string, value []byte) *model.AppError                                                        // SetPluginKey
	CompareAndSetPluginKey(pluginID string, key string, oldValue, newValue []byte) (bool, *model.AppError)                         //
	CompareAndDeletePluginKey(pluginID string, key string, oldValue []byte) (bool, *model.AppError)
	GetPluginKey(pluginID string, key string) ([]byte, *model.AppError)
	DeletePluginKey(pluginID string, key string) *model.AppError
	DeleteAllKeysForPlugin(pluginID string) *model.AppError
	DeleteAllExpiredPluginKeys() *model.AppError
	ListPluginKeys(pluginID string, page, perPage int) ([]string, *model.AppError)
	SetPluginKeyWithExpiry(pluginID string, key string, value []byte, expireInSeconds int64) *model.AppError
}
