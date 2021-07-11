package sub_app_iface

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app/request"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/file"
	"github.com/sitename/sitename/model/giftcard"
	"github.com/sitename/sitename/model/menu"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/model/wishlist"
	"github.com/sitename/sitename/modules/filestore"
	"github.com/sitename/sitename/modules/plugin"
	"github.com/sitename/sitename/store"
)

// GiftCardApp defines methods for giftcard app
type GiftcardApp interface {
	GetGiftCard(id string) (*giftcard.GiftCard, *model.AppError)                   // GetGiftCard returns a giftcard with given id
	GiftcardsByCheckout(checkoutID string) ([]*giftcard.GiftCard, *model.AppError) // GiftcardsByCheckout returns all giftcards belong to given checkout
	GiftcardsByOrder(orderID string) ([]*giftcard.GiftCard, *model.AppError)       // GiftcardsByOrder returns all giftcards belong to given order
}

// PaymentApp defines methods for payment app
type PaymentApp interface {
	GetAllPaymentsByOrderId(orderID string) ([]*payment.Payment, *model.AppError)                // GetAllPaymentsByOrderId returns all payments that belong to order with given orderID
	GetLastOrderPayment(orderID string) (*payment.Payment, *model.AppError)                      // GetLastOrderPayment get most recent payment made for given order
	GetAllPaymentTransactions(paymentID string) ([]*payment.PaymentTransaction, *model.AppError) // GetAllPaymentTransactions returns all transactions belong to given payment
	GetLastPaymentTransaction(paymentID string) (*payment.PaymentTransaction, *model.AppError)   // GetLastPaymentTransaction return most recent transaction made for given payment
	PaymentIsAuthorized(paymentID string) (bool, *model.AppError)                                // PaymentIsAuthorized checks if given payment is authorized
	PaymentGetAuthorizedAmount(pm *payment.Payment) (*goprices.Money, *model.AppError)           // PaymentGetAuthorizedAmount calculates authorized amount
	PaymentCanVoid(pm *payment.Payment) (bool, *model.AppError)                                  // PaymentCanVoid check if payment can void
	// Extract order information along with payment details. Returns information required to process payment and additional billing/shipping addresses for optional fraud-prevention mechanisms.
	CreatePaymentInformation(payment *payment.Payment, paymentToken *string, amount *decimal.Decimal, customerId *string, storeSource bool, additionalData map[string]string) (*payment.PaymentData, *model.AppError)
	GetAlreadyProcessedTransaction(paymentID string, gatewayResponse *payment.GatewayResponse) (*payment.PaymentTransaction, *model.AppError) // GetAlreadyProcessedTransaction returns most recent processed transaction made for given payment
	// CreatePayment creates new payment inside database with given data and returned it
	CreatePayment(gateway, currency, email, customerIpAddress, paymentToken, returnUrl, externalReference string, total decimal.Decimal, extraData map[string]string, checkOut *checkout.Checkout, orDer *order.Order) (*payment.Payment, *model.AppError)
	SavePayment(payment *payment.Payment) (*payment.Payment, *model.AppError)                               // SavePayment save new payment into database
	SaveTransaction(transaction *payment.PaymentTransaction) (*payment.PaymentTransaction, *model.AppError) // SaveTransaction save new payment transaction into database
	// CreatePaymentTransaction save new payment transaction into database and returns it
	CreatePaymentTransaction(paymentID string, kind string, paymentInformation *payment.PaymentData, actionRequired bool, gatewayResponse *payment.GatewayResponse, errorMsg string, isSuccess bool) (*payment.PaymentTransaction, *model.AppError)
	// GetAlreadyProcessedTransactionOrCreateNewTransaction either create new transaction or get already processed transaction
	GetAlreadyProcessedTransactionOrCreateNewTransaction(paymentID, kind string, paymentInformation *payment.PaymentData, actionRequired bool, gatewayResponse *payment.GatewayResponse, errorMsg string) (*payment.PaymentTransaction, *model.AppError)
	CleanCapture(payment *payment.Payment, amount decimal.Decimal) *model.AppError    // CleanCapture Checks if payment can be captured.
	GetPaymentToken(paymentID string) (string, *model.AppError)                       // get first transaction that belongs to given payment and has kind of "auth", IsSuccess is true
	GetAllPaymentsByCheckout(checkoutID string) ([]*payment.Payment, *model.AppError) // GetAllPaymentsByCheckout returns all payments have been made for given checkout
}

// CheckoutApp
type CheckoutApp interface {
	CheckoutbyToken(checkoutToken string) (*checkout.Checkout, *model.AppError) // CheckoutbyToken returns 1 checkout by its token (checkout's pripary key)
	// CheckoutLineShippingRequired(checkoutLine *checkout.CheckoutLine) (bool, *model.AppError) // CheckoutLineShippingRequired check if given checkout line's product variant requires shipping
	FetchCheckoutLines(checkout *checkout.Checkout) ([]*checkout.CheckoutLineInfo, *model.AppError)
	CheckVariantInStock(variant *product_and_discount.ProductVariant, ckout *checkout.Checkout, channelSlug string, quantity *uint, replace, checkQuantity bool) (uint, *checkout.CheckoutLine, *model.AppError)
	CheckoutShippingRequired(checkoutToken string) (bool, *model.AppError)                                                  // CheckoutShippingRequired checks if given checkout requires shipping
	CheckoutsByUser(userID string, channelActive bool) ([]*checkout.Checkout, *model.AppError)                              // CheckoutsByUser returns a list of checkouts belong to given user.
	CheckoutByUser(userID string) (*checkout.Checkout, *model.AppError)                                                     // CheckoutByUser returns a checkout that is active and belongs to given user
	CheckoutCountry(checkout *checkout.Checkout) (string, *model.AppError)                                                  // CheckoutCountry returns country code for given checkout
	CheckoutSetCountry(checkout *checkout.Checkout, newCountryCode string) *model.AppError                                  // CheckoutSetCountry set new country code for checkout
	UpdateCheckout(checkout *checkout.Checkout) (*checkout.Checkout, *model.AppError)                                       // UpdateCheckout updates given checkout and returns it
	GetCustomerEmail(checkout *checkout.Checkout) (string, *model.AppError)                                                 // GetCustomerEmail returns either checkout owner's email or checkout's Email property
	CheckoutTotalGiftCardsBalance(checkout *checkout.Checkout) (*goprices.Money, *model.AppError)                           // CheckoutTotalGiftCardsBalance returns giftcards balance money
	CheckoutLineWithVariant(checkout *checkout.Checkout, productVariantID string) (*checkout.CheckoutLine, *model.AppError) // CheckoutLineWithVariant return a checkout line of given checkout, that checkout line has VariantID of given product variant id
}

// CheckoutApp
type AccountApp interface {
	AddressById(id string) (*account.Address, *model.AppError)                                    // GetAddressById returns address with given id. If not found returns nil and concret error
	UserById(ctx context.Context, userID string) (*account.User, *model.AppError)                 // GetUserById get user from database with given userId
	CustomerEventsByUser(userID string) ([]*account.CustomerEvent, *model.AppError)               // CustomerEventsByUser returns all customer event(s) belong to given user
	AddressesByUserId(userID string) ([]*account.Address, *model.AppError)                        // AddressesByUserId returns list of address(es) (if found) that belong to given user
	UserSetDefaultAddress(userID, addressID, addressType string) (*account.User, *model.AppError) // UserSetDefaultAddress set given address to be default for given user
	AddressDeleteForUser(userID, addressID string) *model.AppError                                // AddressDeleteForUser deletes relationship between given user and address
	UserByEmail(email string) (*account.User, *model.AppError)                                    // UserByEmail try finding user with given email and returns that user
	// CommonCustomerCreateEvent is common method for creating customer events
	CommonCustomerCreateEvent(userID *string, orderID *string, eventType string, params model.StringInterface) (*account.CustomerEvent, *model.AppError)
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
	GetDefaultProfileImage(user *account.User) ([]byte, *model.AppError)
	SetDefaultProfileImage(user *account.User) *model.AppError
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
	GetStatus(userID string) (*model.Status, *model.AppError)
	GetStatusFromCache(userID string) *model.Status
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
	GetSession(token string) (*model.Session, *model.AppError)
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
}

type ProductApp interface {
	ProductVariantById(id string) (*product_and_discount.ProductVariant, *model.AppError)                    // ProductVariantById returns a product variants with given id
	ProductTypesByCheckoutToken(checkoutToken string) ([]*product_and_discount.ProductType, *model.AppError) // ProductTypesByCheckoutToken returns all product types related to given checkout
}

type WishlistApp interface {
	CreateWishlist(userID string) (*wishlist.Wishlist, *model.AppError)                      // CreateWishlist creates new wishlist for given user and returns it
	WishlistByUserID(userID string) (*wishlist.Wishlist, *model.AppError)                    // WishlistByUserID returns a wishlist belongs to given user
	WishlistItemsByWishlistID(wishlistID string) ([]*wishlist.WishlistItem, *model.AppError) // WishlistItemsByWishlistID returns a list of wishlist items that belong to given wishlist
}

type AttributeApp interface {
}

type InvoiceApp interface {
}

type ChannelApp interface {
	// GetChannelBySlug returns a channel (if found) from database with given slug
	GetChannelBySlug(slug string) (*channel.Channel, *model.AppError)
	// GetDefaultChannel get random channel that is active
	GetDefaultActiveChannel() (*channel.Channel, *model.AppError)
	// CleanChannel performs:
	//
	// 1) If given slug is not nil, try getting a channel with that slug.
	//   +) if found, check if channel is active
	//
	// 2) If given slug if nil, it try
	CleanChannel(channelSlug *string) (*channel.Channel, *model.AppError)
}

type WarehouseApp interface {
	CheckStockQuantity(variant *product_and_discount.ProductVariant, countryCode string, channelSlug string, quantity uint) (*warehouse.InsufficientStock, *model.AppError)
	CheckStockQuantityBulk(variants []*product_and_discount.ProductVariant, countryCode string, quantities []uint, channelSlug string) (*warehouse.InsufficientStock, *model.AppError)
	IsProductInStock(productID string, countryCode string, channelSlug string) (bool, *model.AppError)
}

type DiscountApp interface {
}

type OrderApp interface {
	GetAllOrderLinesByOrderId(orderID string) ([]*order.OrderLine, *model.AppError) // GetAllOrderLinesByOrderId returns a slice of order lines that belong to given order
	// OrderShippingIsRequired checks if an order requires ship or not by:
	//
	// 1) Find all child order lines that belong to given order
	//
	// 2) iterates over resulting slice to check if at least one order line requires shipping
	OrderShippingIsRequired(orderID string) (bool, *model.AppError)
	OrderTotalQuantity(orderID string) (int, *model.AppError)                                   // OrderTotalQuantity return total quantity of given order
	UpdateOrderTotalPaid(orderID string) *model.AppError                                        // UpdateOrderTotalPaid update given order's total paid amount
	OrderIsPreAuthorized(orderID string) (bool, *model.AppError)                                // OrderIsPreAuthorized checks if order is pre-authorized
	OrderIsCaptured(orderID string) (bool, *model.AppError)                                     // OrderIsCaptured checks if given order is captured
	OrderSubTotal(orderID string, orderCurrency string) (*goprices.TaxedMoney, *model.AppError) // OrderSubTotal returns sum of TotalPrice of all order lines that belong to given order
	OrderCanCancel(ord *order.Order) (bool, *model.AppError)                                    // OrderCanCalcel checks if given order can be canceled
	OrderCanCapture(ord *order.Order, payment *payment.Payment) (bool, *model.AppError)         // OrderCanCapture checks if given order can capture.
	OrderCanVoid(ord *order.Order, payment *payment.Payment) (bool, *model.AppError)            // OrderCanVoid checks if given order can void
	OrderCanRefund(ord *order.Order, payments []*payment.Payment) (bool, *model.AppError)       // OrderCanRefund checks if order can refund
	CanMarkOrderAsPaid(ord *order.Order, payments []*payment.Payment) (bool, *model.AppError)   // CanMarkOrderAsPaid checks if given order can be marked as paid.
	OrderTotalAuthorized(ord *order.Order) (*goprices.Money, *model.AppError)                   // OrderTotalAuthorized returns order's total authorized amount
	GetOrderCountryCode(ord *order.Order) (string, *model.AppError)                             // GetOrderCountryCode is helper function, returns contry code of given order
	OrderLineById(id string) (*order.OrderLine, *model.AppError)                                // OrderLineById returns order line with id of given id
	OrderById(id string) (*order.Order, *model.AppError)                                        // OrderById returns order with id of given id

}

type MenuApp interface {
	MenuById(id string) (*menu.Menu, *model.AppError)     // MenuById returns a menu with given id
	MenuByName(name string) (*menu.Menu, *model.AppError) // MenuByName returns a menu with given name
	MenuBySlug(slug string) (*menu.Menu, *model.AppError) // MenuBySlug returns a menu with given slug
}

type AppApp interface {
}

type CsvApp interface {
}

type SiteApp interface {
}

type ShippingApp interface {
}

type WebhookApp interface {
}

type PageApp interface {
}

type SeoApp interface {
}

type FileApp interface {
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
}

type PluginApp interface {
	// GetPluginsEnvironment returns the plugin environment for use if plugins are enabled and
	// initialized.
	//
	// To get the plugins environment when the plugins are disabled, manually acquire the plugins
	// lock instead.
	GetPluginsEnvironment() *plugin.Environment
	SetPluginsEnvironment(pluginsEnvironment *plugin.Environment) // SetPluginsEnvironment set plugins environment for server
	SyncPluginsActiveState()                                      // SyncPluginsActiveState
	ServeInterPluginRequest(w http.ResponseWriter, r *http.Request, sourcePluginId, destinationPluginId string)
}

type PreferenceApp interface {
	GetPreferencesForUser(userID string) (model.Preferences, *model.AppError)
	GetPreferenceByCategoryForUser(userID string, category string) (model.Preferences, *model.AppError)
	GetPreferenceByCategoryAndNameForUser(userID string, category string, preferenceName string) (*model.Preference, *model.AppError)
	UpdatePreferences(userID string, preferences model.Preferences) *model.AppError
	DeletePreferences(userID string, preferences model.Preferences) *model.AppError
}
