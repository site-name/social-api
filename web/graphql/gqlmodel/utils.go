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

	return &AddressValidationData{
		CountryCode:        &r.CountryCode,
		CountryName:        &r.CountryName,
		AddressFormat:      &r.AddressFormat,
		AddressLatinFormat: &r.AddressLatinFormat,
		AllowedFields:      StringSliceToStringPointerSlice(r.AllowedFields),
		RequiredFields:     StringSliceToStringPointerSlice(r.RequiredFields),
		UpperFields:        StringSliceToStringPointerSlice(r.UpperFields),
		CountryAreaType:    &r.CountryAreaType,
		CountryAreaChoices: ChoicesToChoiceValues(r.CountryAreaChoices),
		CityType:           &r.CityType,
		CityChoices:        ChoicesToChoiceValues(r.CityChoices),
		CityAreaType:       &r.CityAreaType,
		CityAreaChoices:    ChoicesToChoiceValues(r.CityAreaChoices),
		PostalCodeType:     &r.PostalCodeType,
		PostalCodeMatchers: StringSliceToStringPointerSlice(*i18naddress.RegexesToStrings(r.PostalCodeMatchers)),
		PostalCodeExamples: StringSliceToStringPointerSlice(r.PostalCodeExamples),
		PostalCodePrefix:   &r.PostalCodePrefix,
	}
}

// ChoicesToChoiceValues convert [][2]string => []*ChoiceValue
func ChoicesToChoiceValues(choices [][2]string) []*ChoiceValue {
	res := make([]*ChoiceValue, len(choices))

	for _, choice := range choices {
		res = append(res, &ChoiceValue{
			Raw:     &choice[0],
			Verbose: &choice[1],
		})
	}

	return res
}

// StringSliceToStringPointerSlice convert []string => []*string
func StringSliceToStringPointerSlice(s []string) []*string {
	res := make([]*string, len(s))

	for _, str := range s {
		res = append(res, &str)
	}

	return res
}
