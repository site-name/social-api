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
func (a *Server) ToLocalCurrency(price goprices.Currencier, destCurrency string) (goprices.Currencier, *model_helper.AppError) {
	if a.Config().ThirdPartySettings.OpenExchangeRateApiKey == nil {
		return nil, model_helper.NewAppError("ToLocalCurrency", "app.setting.currency_conversion_disabled.app_error", nil, "", http.StatusNotAcceptable)
	}

	fromCurrency := price.GetCurrency()
	destCurrency = strings.ToUpper(destCurrency)

	if goprices.CurrenciesMap[destCurrency] == "" {
		return nil, model_helper.NewAppError("ToLocalCurrency", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "currency"}, "unknown currency", http.StatusBadRequest)
	}

	if !strings.EqualFold(destCurrency, fromCurrency) {
		return a.ExchangeCurrency(price, destCurrency, nil)
	}

	return price, nil
}

// ExchangeCurrency Exchanges Money, TaxedMoney and their ranges to the specified currency.
// get_rate parameter is a callable taking single argument (target currency)
// that returns proper conversion rate
//
// `base` must be either Money, MoneyRange, TaxedMoney, TaxedMoneyRange.
// `conversionRate` can be nil
func (a *Server) ExchangeCurrency(base goprices.Currencier, toCurrency string, conversionRate *decimal.Decimal) (goprices.Currencier, *model_helper.AppError) {
	var appErr *model_helper.AppError

	if !strings.EqualFold(base.GetCurrency(), model_helper.DEFAULT_CURRENCY.String()) &&
		!strings.EqualFold(toCurrency, model_helper.DEFAULT_CURRENCY.String()) {
		base, appErr = a.ExchangeCurrency(base, model_helper.DEFAULT_CURRENCY.String(), conversionRate)
		if appErr != nil {
			return nil, appErr
		}
	}

	if conversionRate == nil {
		conversionRate, appErr = a.GetConversionRate(base.GetCurrency(), toCurrency)
		if appErr != nil {
			return nil, appErr
		}
	}

	switch t := base.(type) {
	case goprices.Money:
		newAmount := t.GetAmount().Mul(*conversionRate)
		money, _ := goprices.NewMoneyFromDecimal(newAmount, toCurrency)
		return money, nil

	case goprices.MoneyRange:
		newStart, appErr := a.ExchangeCurrency(t.GetStart(), toCurrency, conversionRate)
		if appErr != nil {
			return nil, appErr
		}
		newStop, appErr := a.ExchangeCurrency(t.GetStop(), toCurrency, conversionRate)
		if appErr != nil {
			return nil, appErr
		}
		res, _ := goprices.NewMoneyRange(newStart.(goprices.Money), newStop.(goprices.Money))
		return res, nil

	case goprices.TaxedMoney:
		newNet, appErr := a.ExchangeCurrency(t.GetNet(), toCurrency, conversionRate)
		if appErr != nil {
			return nil, appErr
		}
		newGross, appErr := a.ExchangeCurrency(t.GetGross(), toCurrency, conversionRate)
		if appErr != nil {
			return nil, appErr
		}
		res, _ := goprices.NewTaxedMoney(newNet.(goprices.Money), newGross.(goprices.Money))
		return res, nil

	case goprices.TaxedMoneyRange:
		newStart, appErr := a.ExchangeCurrency(t.GetStart(), toCurrency, conversionRate)
		if appErr != nil {
			return nil, appErr
		}
		newStop, appErr := a.ExchangeCurrency(t.GetStop(), toCurrency, conversionRate)
		if appErr != nil {
			return nil, appErr
		}
		res, _ := goprices.NewTaxedMoneyRange(newStart.(goprices.TaxedMoney), newStop.(goprices.TaxedMoney))
		return res, nil

	default:
		return nil, nil
	}
}

// decimalOne holds value of 1.0
var decimalOne = decimal.NewFromInt(1)

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
	if exist && value != nil {
		exchangeRate, ok := value.(*model.OpenExchangeRate)
		if ok && exchangeRate != nil && !exchangeRate.Rate.IsNil() {
			rate = *exchangeRate.Rate.Decimal
		}
	} else {
		exchangeRatesFromDatabase, err := a.Store.OpenExchangeRate().GetAll()
		if err != nil {
			return nil, model_helper.NewAppError("GetConversionRate", "app.currency.error_finding_conversion_rates.app_error", nil, err.Error(), http.StatusInternalServerError)
		}

		for _, exchangeRate := range exchangeRatesFromDatabase {
			if exchangeRate.ToCurrency.String() == rateCurrency {
				rate = *exchangeRate.Rate.Decimal
				break
			}
		}
	}

	if reverseRate {
		rate = decimalOne.Div(rate)
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
