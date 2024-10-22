package shipping

import (
	"regexp"

	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
)

type checkPostalCodeFunc func(code, start, end string) bool

var (
	ukPostalCodePattern    = regexp.MustCompile(`^([A-Z]{1,2})([0-9]+)([A-Z]?) ?([0-9][A-Z]{2})$`) // ukPostalCodePattern to check againts United Kingdom postal codes
	irishPostalCodePattern = regexp.MustCompile(`([\dA-Z]{3}) ?([\dA-Z]{4})`)                      // irishPostalCodePattern to check againts ireland postal codes
	countryFuncMap         = map[model.CountryCode]checkPostalCodeFunc{
		model.CountryCodeGB: CheckUkPostalCode,    // United Kingdom
		model.CountryCodeIM: CheckUkPostalCode,    // Isle of Man
		model.CountryCodeGG: CheckUkPostalCode,    // Guernsey
		model.CountryCodeJE: CheckUkPostalCode,    // Jersey
		model.CountryCodeIE: CheckIRishPostalCode, // Ireland
	}
)

func GroupValues(pattern *regexp.Regexp, values ...string) {
	// TODO: implement me
	panic("not implt")
}

func CompareValues(code, start, end string) bool {
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
func CheckUkPostalCode(code, start, end string) bool {
	GroupValues(ukPostalCodePattern, code, start, end)

	return CompareValues(code, start, end)
}

// Check postal code for Ireland, split the code by regex.
// Example postal codes: A65 2F0A, A61 2F0G.
func CheckIRishPostalCode(code, start, end string) bool {
	GroupValues(irishPostalCodePattern, code, start, end)

	return CompareValues(code, start, end)
}

// Fallback for any country not present in country_func_map.
// Perform simple lexicographical comparison without splitting to sections.
func CheckAnyPostalCode(code, start, end string) bool {
	return CompareValues(code, start, end)
}

func CheckPostalCodeInRange(countryCode model.CountryCode, postalCode, start, end string) bool {
	fun, exist := countryFuncMap[countryCode]
	if !exist {
		fun = CheckAnyPostalCode
	}

	return fun(postalCode, start, end)
}

func CheckShippingMethodForPostalCode(customerShippingAddress model.Address, method model.ShippingMethod) map[*model.ShippingMethodPostalCodeRule]bool {
	var rules model.ShippingMethodPostalCodeRuleSlice
	if method.R != nil {
		rules = method.R.ShippingMethodPostalCodeRules
	}
	return lo.SliceToMap(
		rules,
		func(rule *model.ShippingMethodPostalCodeRule) (*model.ShippingMethodPostalCodeRule, bool) {
			return rule, CheckPostalCodeInRange(customerShippingAddress.Country, customerShippingAddress.PostalCode, rule.Start, rule.End)
		},
	)
}

// IsShippingMethodApplicableForPostalCode Return if shipping method is applicable with the postal code rules.
func IsShippingMethodApplicableForPostalCode(customerShippingAddress model.Address, shippingMethod model.ShippingMethod) bool {
	result := CheckShippingMethodForPostalCode(customerShippingAddress, shippingMethod)

	resultLength := len(result)
	if resultLength == 0 {
		return true
	}

	var (
		numberOfInclude     int
		numberOfExclude     int
		atLeastOneValueTrue = false // all rules's inclusion_type == 'include'
		allValueAreFalse    = true  // all rules's inclusion_type == 'exclude'
	)

	for key, value := range result {
		switch key.InclusionType {
		case model.InclusionTypeInclude:
			numberOfInclude++
		case model.InclusionTypeExclude:
			numberOfExclude++
		}

		if value {
			atLeastOneValueTrue = true
			allValueAreFalse = false
		}
	}

	return numberOfInclude == resultLength && atLeastOneValueTrue || (numberOfExclude == resultLength && allValueAreFalse)
}

func (a *ServiceShipping) FilterShippingMethodsByPostalCodeRules(shippingMethods model.ShippingMethodSlice, shippingAddress model.Address) model.ShippingMethodSlice {
	return lo.Filter(shippingMethods, func(method *model.ShippingMethod, _ int) bool {
		return method != nil && IsShippingMethodApplicableForPostalCode(shippingAddress, *method)
	})
}
