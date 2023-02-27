package app

import (
	"net/http"
	"strings"

	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app/request"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/plugin"
)

type TokenLocation int

const (
	TokenLocationNotFound TokenLocation = iota
	TokenLocationHeader
	TokenLocationCookie
)

// String implements fmt Stringer interface
func (tl TokenLocation) String() string {
	switch tl {
	case TokenLocationNotFound:
		return "Not Found"
	case TokenLocationHeader:
		return "Header"
	case TokenLocationCookie:
		return "Cookie"
	default:
		return "Unknown"
	}
}

// ParseAuthTokenFromRequest reads header "Authorization" from request's header, then parses it into token and token location
func ParseAuthTokenFromRequest(r *http.Request) (string, TokenLocation) {
	authHeader := r.Header.Get(model.HEADER_AUTH)

	// Attempt to parse the token from the cookie
	if cookie, err := r.Cookie(model.SESSION_COOKIE_TOKEN); err == nil {
		return cookie.Value, TokenLocationCookie
	}

	// Parse the token from the header
	if len(authHeader) > 6 && strings.ToUpper(authHeader[0:6]) == model.HEADER_BEARER {
		// Default session token
		return authHeader[7:], TokenLocationHeader
	}

	if len(authHeader) > 5 && strings.ToLower(authHeader[0:5]) == model.HEADER_TOKEN {
		// OAuth token
		return authHeader[6:], TokenLocationHeader
	}

	return "", TokenLocationNotFound
}

// PluginContext
func PluginContext(c *request.Context) *plugin.Context {
	return &plugin.Context{
		RequestId:      c.RequestId(),
		SessionId:      c.Session().Id,
		IpAddress:      c.IpAddress(),
		AcceptLanguage: c.AcceptLanguage(),
		UserAgent:      c.UserAgent(),
	}
}

// ToLocalCurrency performs convert given price to local currency
//
// NOTE: `price` must be either *Money, *MoneyRange, *TaxedMoney, *TaxedMoneyRange
func (a *Server) ToLocalCurrency(price interface{}, currency string) (interface{}, *model.AppError) {
	// validate if currency exchange is enabled
	if a.Config().ThirdPartySettings.OpenExchangeRateApiKey == nil {
		return nil, model.NewAppError("ToLocalCurrency", "app.setting.currency_conversion_disabled.app_error", nil, "", http.StatusNotAcceptable)
	}

	// validate price is valid:
	var fromCurrency string

	switch t := price.(type) {
	case goprices.Currencyable:
		fromCurrency = t.MyCurrency()

	default:
		return nil, model.NewAppError("ToLocalCurrency", InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "price"}, "price is not Money type", http.StatusBadRequest)
	}
	// validate provided currency is valid:
	currency = strings.ToUpper(currency)
	if goprices.CurrenciesMap[currency] == "" {
		return nil, model.NewAppError("ToLocalCurrency", InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "currency"}, "unknown currency", http.StatusBadRequest)
	}

	if !strings.EqualFold(currency, fromCurrency) {
		return a.ExchangeCurrency(price, currency, nil)
	}

	return price, nil
}

// ExchangeCurrency Exchanges Money, TaxedMoney and their ranges to the specified currency.
// get_rate parameter is a callable taking single argument (target currency)
// that returns proper conversion rate
//
// `base` must be either *Money, *MoneyRange, *TaxedMoney, *TaxedMoneyRange.
// `conversionrate` can be nil
//
// NOTE: `base` and `toCurrency` must be validated before given to me.
func (a *Server) ExchangeCurrency(base interface{}, toCurrency string, conversionRate *decimal.Decimal) (interface{}, *model.AppError) {
	var appErr *model.AppError

	impl, ok := base.(goprices.Currencyable)
	if ok {
		if !strings.EqualFold(impl.MyCurrency(), model.DEFAULT_CURRENCY) &&
			!strings.EqualFold(toCurrency, model.DEFAULT_CURRENCY) {
			base, appErr = a.ExchangeCurrency(base, model.DEFAULT_CURRENCY, conversionRate)
			if appErr != nil {
				return nil, appErr
			}
		}
	} else {
		return nil, model.NewAppError("ExchangeCurrency", InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "base"}, "", http.StatusBadRequest)
	}

	if conversionRate == nil {
		conversionRate, appErr = a.GetConversionRate(impl.MyCurrency(), toCurrency)
	}
	if appErr != nil {
		return nil, appErr
	}

	switch t := base.(type) {
	case *goprices.Money:
		newAmount := t.Amount.Mul(*conversionRate)
		return &goprices.Money{
			Amount:   newAmount,
			Currency: toCurrency,
		}, nil

	case *goprices.MoneyRange:
		newStart, appErr := a.ExchangeCurrency(t.Start, toCurrency, conversionRate)
		if appErr != nil {
			return nil, appErr
		}
		newStop, appErr := a.ExchangeCurrency(t.Stop, toCurrency, conversionRate)
		if appErr != nil {
			return nil, appErr
		}
		res, _ := goprices.NewMoneyRange(newStart.(*goprices.Money), newStop.(*goprices.Money))
		return res, nil

	case *goprices.TaxedMoney:
		newNet, appErr := a.ExchangeCurrency(t.Net, toCurrency, conversionRate)
		if appErr != nil {
			return nil, appErr
		}
		newGross, appErr := a.ExchangeCurrency(t.Gross, toCurrency, conversionRate)
		if appErr != nil {
			return nil, appErr
		}
		res, _ := goprices.NewTaxedMoney(newNet.(*goprices.Money), newGross.(*goprices.Money))
		return res, nil

	case *goprices.TaxedMoneyRange:
		newStart, appErr := a.ExchangeCurrency(t.Start, toCurrency, conversionRate)
		if appErr != nil {
			return nil, appErr
		}
		newStop, appErr := a.ExchangeCurrency(t.Stop, toCurrency, conversionRate)
		if appErr != nil {
			return nil, appErr
		}
		res, _ := goprices.NewTaxedMoneyRange(newStart.(*goprices.TaxedMoney), newStop.(*goprices.TaxedMoney))
		return res, nil

	default:
		return nil, nil
	}
}

// GetConversionRate get conversion rate to use in exchange.
// It first try getting exchange rate from cache and returns the found value. If nothing found, it try finding from database
func (a *Server) GetConversionRate(fromCurrency string, toCurrency string) (*decimal.Decimal, *model.AppError) {
	fromCurrency = strings.ToUpper(fromCurrency)
	toCurrency = strings.ToUpper(toCurrency)

	var (
		reverseRate  bool
		rateCurrency string
	)
	if toCurrency == model.DEFAULT_CURRENCY {
		rateCurrency = fromCurrency
		reverseRate = true
	} else {
		rateCurrency = toCurrency
	}

	var rate decimal.Decimal
	// try get rate from the cache first, if not found, find in database
	value, exist := a.ExchangeRateMap.Load(rateCurrency)
	if exist {
		rate = *(value.(*model.OpenExchangeRate).Rate)
	} else {
		exchangeRatesFromDatabase, err := a.Store.OpenExchangeRate().GetAll()
		if err != nil {
			return nil, model.NewAppError("GetConversionRate", "app.currency.error_finding_conversion_rates.app_error", nil, err.Error(), http.StatusInternalServerError)
		}

		for _, exchangeRate := range exchangeRatesFromDatabase {
			if exchangeRate.ToCurrency == rateCurrency {
				rate = *exchangeRate.Rate
				break
			}
		}
	}

	if reverseRate {
		rate = decimal.NewFromInt(1).Div(rate)
	}

	return &rate, nil
}

// GetProtocol returns request's protocol
func GetProtocol(r *http.Request) string {
	if r.Header.Get(model.HEADER_FORWARDED_PROTO) == "https" || r.TLS != nil {
		return "https"
	}
	return "http"
}
