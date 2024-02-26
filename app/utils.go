package app

import (
	"net/http"
	"strings"

	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app/request"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
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
	// Attempt to parse the token from the cookie
	if cookie, err := r.Cookie(model_helper.SESSION_COOKIE_TOKEN); err == nil {
		return cookie.Value, TokenLocationCookie
	}

	authHeader := r.Header.Get(model_helper.HeaderAuth)
	// Parse the token from the header
	if len(authHeader) > 6 && strings.ToUpper(authHeader[0:6]) == model_helper.HeaderBearer {
		// Default session token
		return authHeader[7:], TokenLocationHeader
	}

	if len(authHeader) > 5 && strings.ToLower(authHeader[0:5]) == model_helper.HeaderToken {
		// OAuth token
		return authHeader[6:], TokenLocationHeader
	}

	return "", TokenLocationNotFound
}

// PluginContext
func PluginContext(c request.Context) *plugin.Context {
	return &plugin.Context{
		RequestId:      c.RequestId(),
		SessionId:      c.Session().ID,
		IpAddress:      c.IpAddress(),
		AcceptLanguage: c.AcceptLanguage(),
		UserAgent:      c.UserAgent(),
	}
}

// ToLocalCurrency performs convert given price to local currency
//
// NOTE: `price` must be either *Money, *MoneyRange, *TaxedMoney, *TaxedMoneyRange
func (a *Server) ToLocalCurrency(price interface{}, currency string) (interface{}, *model_helper.AppError) {
	// validate if currency exchange is enabled
	if a.Config().ThirdPartySettings.OpenExchangeRateApiKey == nil {
		return nil, model_helper.NewAppError("ToLocalCurrency", "app.setting.currency_conversion_disabled.app_error", nil, "", http.StatusNotAcceptable)
	}

	// validate price is valid:
	var fromCurrency string

	switch t := price.(type) {
	case goprices.Currencier:
		fromCurrency = t.GetCurrency()

	default:
		return nil, model_helper.NewAppError("ToLocalCurrency", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "price"}, "price is not Money type", http.StatusBadRequest)
	}
	// validate provided currency is valid:
	currency = strings.ToUpper(currency)
	if goprices.CurrenciesMap[currency] == "" {
		return nil, model_helper.NewAppError("ToLocalCurrency", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "currency"}, "unknown currency", http.StatusBadRequest)
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
func (a *Server) ExchangeCurrency(base interface{}, toCurrency string, conversionRate *decimal.Decimal) (interface{}, *model_helper.AppError) {
	var appErr *model_helper.AppError

	impl, ok := base.(goprices.Currencier)
	if ok {
		if !strings.EqualFold(impl.GetCurrency(), model_helper.DEFAULT_CURRENCY.String()) &&
			!strings.EqualFold(toCurrency, model_helper.DEFAULT_CURRENCY.String()) {
			base, appErr = a.ExchangeCurrency(base, model_helper.DEFAULT_CURRENCY.String(), conversionRate)
			if appErr != nil {
				return nil, appErr
			}
		}
	} else {
		return nil, model_helper.NewAppError("ExchangeCurrency", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "base"}, "", http.StatusBadRequest)
	}

	if conversionRate == nil {
		conversionRate, appErr = a.GetConversionRate(impl.GetCurrency(), toCurrency)
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
func (a *Server) GetConversionRate(fromCurrency string, toCurrency string) (*decimal.Decimal, *model_helper.AppError) {
	fromCurrency = strings.ToUpper(fromCurrency)
	toCurrency = strings.ToUpper(toCurrency)

	var (
		reverseRate  bool
		rateCurrency string
	)
	if toCurrency == model_helper.DEFAULT_CURRENCY.String() {
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
			return nil, model_helper.NewAppError("GetConversionRate", "app.currency.error_finding_conversion_rates.app_error", nil, err.Error(), http.StatusInternalServerError)
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
	if r.Header.Get(model_helper.HeaderForwardedProto) == "https" || r.TLS != nil {
		return "https"
	}
	return "http"
}

func (s *Server) GetSiteContext() map[string]any {
	settings := s.Config().ServiceSettings
	return map[string]any{
		"domain":    settings.SiteURL,
		"site_name": settings.SiteName,
	}
}
