package payment

import (
	"github.com/shopspring/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/checkout"
)

// Uniform way to represent payment method information.
type PaymentMethodInfo struct {
	Last4    string
	ExpYear  uint16
	ExpMonth uint8
	Brand    string
	Name     string
	Type     string
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
	Error                       string
	CustomerID                  string
	PaymentMethodInfo           *PaymentMethodInfo
	RawResponse                 model.StringMap
	ActionRequiredData          map[string]string
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
	Country        string
	CountryArea    string
	Phone          string
}

// AddressDataFromAddress convert *account.Address to *AddressData
func AddressDataFromAddress(a *account.Address) *AddressData {
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
	Gateway           string
	Amount            decimal.Decimal
	Currency          string
	Billing           *AddressData
	Shipping          *AddressData
	PaymentID         string
	GraphqlPaymentID  string // note
	OrderID           string
	CustomerIpAddress string
	CustomerEmail     string
	Token             *string
	CustomerID        *string
	ReuseSource       bool
	Data              model.StringMap
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
	ConnectionParams    model.StringInterface
	StoreCustomer       bool
	Require3dSecure     bool
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
	Config     []model.StringInterface
}

type InitializedPaymentResponse struct {
	Gateway string
	Name    string
	Data    interface{}
}

type PaymentInterface interface {
	ListPaymentGateWays(currency *string, checkout *checkout.Checkout, channelSlug *string, activeOnly bool) []*PaymentGateway
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