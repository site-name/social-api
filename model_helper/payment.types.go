package model_helper

import (
	"github.com/site-name/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/sitename/sitename/modules/util"
)

type PaymentErrorCode string

const (
	BILLING_ADDRESS_NOT_SET     PaymentErrorCode = "billing_address_not_set"
	GRAPHQL_ERROR               PaymentErrorCode = "graphql_error"
	INVALID                     PaymentErrorCode = "invalid"
	NOT_FOUND                   PaymentErrorCode = "not_found"
	REQUIRED                    PaymentErrorCode = "required"
	UNIQUE                      PaymentErrorCode = "unique"
	PARTIAL_PAYMENT_NOT_ALLOWED PaymentErrorCode = "partial_payment_not_allowed"
	SHIPPING_ADDRESS_NOT_SET    PaymentErrorCode = "shipping_address_not_set"
	INVALID_SHIPPING_METHOD     PaymentErrorCode = "invalid_shipping_method"
	SHIPPING_METHOD_NOT_SET     PaymentErrorCode = "shipping_method_not_set"
	PAYMENT_ERROR               PaymentErrorCode = "payment_error"
	NOT_SUPPORTED_GATEWAY       PaymentErrorCode = "not_supported_gateway"
	CHANNEL_INACTIVE            PaymentErrorCode = "channel_inactive"
)

const (
	ErrPayment = "app.payment.payment_error.app_error"
)

type PaymentError struct {
	Where   string
	Message string
	Code    PaymentErrorCode
}

func NewPaymentError(where, message string, code PaymentErrorCode) *PaymentError {
	return &PaymentError{
		Where:   where,
		Message: message,
		Code:    code,
	}
}

func (p *PaymentError) Error() string {
	return p.Where + ": " + p.Message
}

type GatewayError struct {
	Where   string
	Message string
}

func (g *GatewayError) Error() string {
	return g.Where + ": " + g.Message
}

// Uniform way to represent payment method information.
type PaymentMethodInfo struct {
	First4   *string
	Last4    *string
	ExpYear  *int
	ExpMonth *int
	Brand    *string
	Name     *string
	Type     *string
}

// for storing gateway response.
// Used for unifying the representation of gateway response.
// It is required to communicate between Sitename and given payment gateway.
type GatewayResponse struct {
	IsSucess                    bool
	ActionRequired              bool
	Kind                        model.TransactionKind
	Amount                      decimal.Decimal
	Currency                    model.Currency
	TransactionID               string
	Error                       string
	CustomerID                  string
	PaymentMethodInfo           *PaymentMethodInfo
	RawResponse                 model_types.JSONString
	ActionRequiredData          model_types.JSONString
	TransactionAlreadyProcessed bool
	SearchableKey               string
	PspReference                string
}

type AddressData struct {
	FirstName      string
	LastName       string
	CompanyName    string
	StreetAddress1 string
	StreetAddress2 string
	City           string
	CityArea       string
	PostalCode     string
	Country        model.CountryCode
	CountryArea    string
	Phone          string
}

// AddressDataFromAddress convert *Address to *AddressData
func AddressDataFromAddress(a *model.Address) *AddressData {
	return &AddressData{
		FirstName:      a.FirstName,
		LastName:       a.LastName,
		CompanyName:    a.CompanyName,
		StreetAddress1: a.StreetAddress1,
		StreetAddress2: a.StreetAddress2,
		City:           a.City,
		CityArea:       a.CityArea,
		PostalCode:     a.PostalCode,
		Country:        a.Country,
		CountryArea:    a.CountryArea,
		Phone:          a.Phone,
	}
}

// Dataclass for storing all payment information.
// Used for unifying the representation of data.
// It is required to communicate between Saleor and given payment gateway.
type PaymentData struct {
	Gateway            string
	Amount             decimal.Decimal
	Currency           model.Currency
	Billing            *AddressData // can be bil
	Shipping           *AddressData // can be nil
	PaymentID          string       // payment's Token property
	OrderID            *string      // can be nil
	CustomerIpAddress  *string      // can be nil
	CustomerEmail      string
	Token              *string // can be nil
	CustomerID         *string // can be nil
	ReuseSource        bool
	Data               model_types.JSONString   // can be nil
	GraphqlPaymentID   string                   // default to payment's Token
	GraphqlCustomerID  *string                  // can be nil
	StorePaymentMethod model.StorePaymentMethod // default to StorePaymentMethodEnum_NONE ("none")
	PaymentMetadata    model_types.JSONString   // can be nil
}

// Dataclass for payment gateway token fetching customization.
type TokenConfig struct {
	CustomerID string
}

// Dataclass for storing gateway config data.
// Used for unifying the representation of config data.
// It is required to communicate between Saleor and given payment gateway.
type GatewayConfig struct {
	GatewayName         string
	AutoCapture         bool
	SupportedCurrencies string
	ConnectionParams    model_types.JSONString
	StoreCustomer       bool
	Require3dSecure     bool
}

// Dataclass for storing information about stored payment sources in gateways.
type CustomerSource struct {
	Id             string
	Gateway        string
	CreditCardInfo *PaymentMethodInfo
	Metadata       map[string]string
}

// Dataclass for storing information about a payment gateway.
type PaymentGateway struct {
	Id         string
	Name       string
	Currencies util.AnyArray[string]
	Config     []model_types.JSONString
}

type InitializedPaymentResponse struct {
	Gateway string
	Name    string
	Data    any
}

type PaymentInterface interface {
	ListPaymentGateWays(currency *string, checkout *model.Checkout, channelSlug *string, activeOnly bool) []*PaymentGateway
	AuthorizePayment(gateway string, paymentInformation *PaymentData, channelSlug *string) *GatewayResponse
	CapturePayment(gateway string, paymentInformation *PaymentData, channelSlug *string) *GatewayResponse
	RefundPayent(gateway string, paymentInformation *PaymentData, channelSlug *string) *GatewayResponse
	VoidPayment(gateway string, paymentInformation *PaymentData, channelSlug *string) *GatewayResponse
	ConfirmPayment(gateway string, paymentInformation *PaymentData, channelSlug *string) *GatewayResponse
	TokenIsRequiredAsPaymentInput(gateway string, channelSlug *string) bool
	ProcessPayment(gateway string, paymentInformation *PaymentData, channelSlug *string) *GatewayResponse
	GetClientToken(gateway string, tokenConfig *TokenConfig, channelSlug *string) string
	ListPaymentSources(gateway string, customerId string, channelSlug *string) []*CustomerSource
}
