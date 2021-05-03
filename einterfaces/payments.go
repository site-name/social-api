package einterfaces

import (
	"github.com/shopspring/decimal"
	"github.com/sitename/sitename/model"
)

// Uniform way to represent payment method information.
type PaymentMethodInfo struct {
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
	Kind                        string
	Amount                      decimal.Decimal
	Currency                    string
	TransactionID               string
	Error                       *string
	CustomerID                  *string
	PaymentMethodInfo           *PaymentMethodInfo
	RawResponse                 *model.StringMap
	ActionRequiredData          interface{}
	TransactionAlreadyProcessed *bool
	SearchableKey               *string
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
	Country        string
	CountryArea    string
	Phone          string
}

// Dataclass for storing all payment information.
// Used for unifying the representation of data.
// It is required to communicate between Saleor and given payment gateway.
type PaymentData struct {
	Amount            decimal.Decimal
	Currency          string
	Billing           *AddressData
	Shipping          *AddressData
	PaymentID         string
	GraphqlPaymentID  string
	OrderID           *string
	CustomerIpAddress *string
	CustomerEmail     string
	Token             *string
	CustomerID        *string
	ReuseSource       bool
	Data              *model.StringMap
}

// Dataclass for payment gateway token fetching customization.
type TokenConfig struct {
	CustomerID *string
}

// Dataclass for storing gateway config data.
// Used for unifying the representation of config data.
// It is required to communicate between Saleor and given payment gateway.
type GatewayConfig struct {
	GatewayName         string
	AutoCapture         bool
	SupportedCurrencies string
	ConnectionParams    model.StringInterface
	StoreCustomer       *bool
	Require3dSecure     *bool
}

// Dataclass for storing information about stored payment sources in gateways.
type CustomerSource struct {
	Id             string
	Gateway        string
	CreditCardInfo *PaymentMethodInfo
}

// Dataclass for storing information about a payment gateway.
type PaymentGateway struct {
	Id         string
	Name       string
	Currencies model.StringArray
	Config     []*model.StringInterface
}

type InitializedPaymentResponse struct {
	Gateway string
	Name    string
	Data    interface{}
}

type PaymentInterface interface {
	ListPaymentGateWays(currency string, checkout *model.Checkout, activeOnly bool) []*PaymentGateway
	AuthorizePayment(gateway string, paymentInformation *PaymentData) *GatewayResponse
	CapturePayment(gateway string, paymentInformation *PaymentData) *GatewayResponse
	RefundPayent(gateway string, paymentInformation *PaymentData) *GatewayResponse
	VoidPayment(gateway string, paymentInformation *PaymentData) *GatewayResponse
	ConfirmPayment(gateway string, paymentInformation *PaymentData) *GatewayResponse
	TokenIsRequiredAsPaymentInput(gateway string) bool
	ProcessPayment(gateway string, paymentInformation *PaymentData) *GatewayResponse
	GetClientToken(gateway string, tokenConfig *TokenConfig) string
	ListPaymentSources(gateway string, customerId string) []*CustomerSource
}
