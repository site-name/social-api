package shipping

import (
	"regexp"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/shipping"
)

type checkPostalCodeFunc func(code string, start string, end string) bool

var (
	ukPostalCodePattern    *regexp.Regexp
	irishPostalCodePattern *regexp.Regexp
	countryFuncMap         map[string]checkPostalCodeFunc
)

func init() {
	ukPostalCodePattern = regexp.MustCompile(`^([A-Z]{1,2})([0-9]+)([A-Z]?) ?([0-9][A-Z]{2})$`)
	irishPostalCodePattern = regexp.MustCompile(`([\dA-Z]{3}) ?([\dA-Z]{4})`)
	countryFuncMap = map[string]checkPostalCodeFunc{
		"GB": CheckUkPostalCode,    // United Kingdom
		"IM": CheckUkPostalCode,    // Isle of Man
		"GG": CheckUkPostalCode,    // Guernsey
		"JE": CheckUkPostalCode,    // Jersey
		"IE": CheckIRishPostalCode, // Ireland
	}
}

func GroupValues(pattern *regexp.Regexp, values ...string) {

	for _, value := range values {

	}
}

func CompareValues(code string, start string, end string) bool {
	if code == "" || start == "" {
		return false
	}
	if end == "" {
		return start <= code
	}
	return code >= start && code <= end
}

// Check postal code for uk, split the code by regex.
// Example postal codes: BH20 2BC  (UK), IM16 7HF  (Isle of Man).
func CheckUkPostalCode(code string, start string, end string) bool {
	GroupValues(ukPostalCodePattern, code, start, end)

	return CompareValues()
}

// Check postal code for Ireland, split the code by regex.
// Example postal codes: A65 2F0A, A61 2F0G.
func CheckIRishPostalCode(code string, start string, end string) bool {
	GroupValues(irishPostalCodePattern, code, start, end)

	return CompareValues()
}

// Fallback for any country not present in country_func_map.
// Perform simple lexicographical comparison without splitting to sections.
func CheckAnyPostalCode(code string, start string, end string) bool {
	return CompareValues(code, start, end)
}

func CheckPostalCodeInRange(countryCode string, postalCode string, start string, end string) bool {
	fun, exist := countryFuncMap[countryCode]
	if !exist || fun == nil {
		fun = CheckAnyPostalCode
	}

	return fun(postalCode, start, end)
}

func CheckShippingMethodForPostalCode(customerShippingAddress *account.Address, method *shipping.ShippingMethod) map[*shipping.ShippingMethodPostalCodeRule]bool {
	result := map[*shipping.ShippingMethodPostalCodeRule]bool{}

	for _, rule := range method.ShippingMethodPostalCodeRules {
		result[rule] = CheckPostalCodeInRange(customerShippingAddress.Country, customerShippingAddress.PostalCode, rule.Start, rule.End)
	}

	return result
}

// IsShippingMethodApplicableForPostalCode Return if shipping method is applicable with the postal code rules.
func IsShippingMethodApplicableForPostalCode(customerShippingAddress *account.Address, shippingMethod *shipping.ShippingMethod) bool {
	result := CheckShippingMethodForPostalCode(customerShippingAddress, shippingMethod)

	resultLength := len(result)
	if resultLength == 0 {
		return true
	}

	var (
		numberOfInclude     int
		numberOfExclude     int
		atLeastOneValueTrue = false // all rules inclusion_type == 'include'
		allValueAreFalse    = true  // all rules inclusion_type == 'exclude'
	)

	for key, value := range result {
		switch key.InclusionType {
		case shipping.INCLUDE:
			numberOfInclude++
		case shipping.EXCLUDE:
			numberOfExclude++
		}

		if value {
			atLeastOneValueTrue = true
			allValueAreFalse = false
		}
	}

	return numberOfInclude == resultLength && atLeastOneValueTrue || (numberOfExclude == resultLength && allValueAreFalse)

	return false
}

// FilterShippingMethodsByPostalCodeRules Filter shipping methods for given address by postal code rules.
func (a *AppShipping) FilterShippingMethodsByPostalCodeRules(shippingMethods []*shipping.ShippingMethod, shippingAddressID string) ([]*shipping.ShippingMethod, *model.AppError) {

	shippingAddress, appErr := a.AccountApp().AddressById(shippingAddressID)
	if appErr != nil {
		return nil, appErr
	}

	for i, shippingMethod := range shippingMethods {
		if !IsShippingMethodApplicableForPostalCode(shippingAddress, shippingMethod) {
			shippingMethods = append(shippingMethods[:i], shippingMethods[i+1:]...)
		}
	}

	return shippingMethods, nil
}
