package gqlmodel

import "github.com/site-name/i18naddress"

// MapToGraphqlMetaDataItems converts a map of key-value into a slice of graphql MetadataItems
func MapToGraphqlMetaDataItems(m map[string]string) []*MetadataItem {
	if m == nil {
		return []*MetadataItem{}
	}

	res := make([]*MetadataItem, len(m))
	for key, value := range m {
		res = append(res, &MetadataItem{Key: key, Value: value})
	}

	return res
}

// MetaDataToStringMap converts a slice of *MetadataInput || *MetadataItem to map[string]string.
//
// Other types will result in an empty map
func MetaDataToStringMap(metaList interface{}) map[string]string {
	res := make(map[string]string)

	switch t := metaList.(type) {
	case []*MetadataInput:
		for _, input := range t {
			res[input.Key] = input.Value
		}
	case []*MetadataItem:
		for _, item := range t {
			res[item.Key] = item.Value
		}
	default:
		return res
	}

	return res
}

// I18nAddressValidationRulesToGraphql convert *i18naddress.ValidationRules to *AddressValidationData
func I18nAddressValidationRulesToGraphql(r *i18naddress.ValidationRules) *AddressValidationData {

	allowedFields := StringSliceToStringPointerSlice(r.AllowedFields)
	requiredFields := StringSliceToStringPointerSlice(r.RequiredFields)
	upperFields := StringSliceToStringPointerSlice(r.UpperFields)
	postalCodeMatchers := StringSliceToStringPointerSlice(*i18naddress.RegexesToStrings(r.PostalCodeMatchers))
	postalCodeExamples := StringSliceToStringPointerSlice(r.PostalCodeExamples)

	countryAreaChoices := ChoicesToChoiceValues(r.CountryAreaChoices)
	cityChoices := ChoicesToChoiceValues(r.CityChoices)
	cityAreaChoices := ChoicesToChoiceValues(r.CityAreaChoices)

	return &AddressValidationData{
		CountryCode:        &r.CountryCode,
		CountryName:        &r.CountryName,
		AddressFormat:      &r.AddressFormat,
		AddressLatinFormat: &r.AddressLatinFormat,
		AllowedFields:      allowedFields,
		RequiredFields:     requiredFields,
		UpperFields:        upperFields,
		CountryAreaType:    &r.CountryAreaType,
		CountryAreaChoices: countryAreaChoices,
		CityType:           &r.CityType,
		CityChoices:        cityChoices,
		CityAreaType:       &r.CityAreaType,
		CityAreaChoices:    cityAreaChoices,
		PostalCodeType:     &r.PostalCodeType,
		PostalCodeMatchers: postalCodeMatchers,
		PostalCodeExamples: postalCodeExamples,
		PostalCodePrefix:   &r.PostalCodePrefix,
	}
}

// ChoicesToChoiceValues convert [][2]string => []*ChoiceValue
func ChoicesToChoiceValues(choices [][2]string) []*ChoiceValue {
	res := []*ChoiceValue{}

	for i := range choices {
		res = append(res, &ChoiceValue{
			Raw:     &choices[i][0],
			Verbose: &choices[i][1],
		})
	}

	return res
}

// StringSliceToStringPointerSlice convert []string => []*string
func StringSliceToStringPointerSlice(s []string) []*string {
	res := []*string{}

	for i := range s {
		res = append(res, &s[i])
	}

	return res
}
