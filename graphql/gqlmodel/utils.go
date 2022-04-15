package gqlmodel

import (
	"strings"

	"github.com/site-name/i18naddress"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
)

// MapToGraphqlMetaDataItems converts a map of key-value into a slice of graphql MetadataItems
func MapToGraphqlMetaDataItems(m map[string]string) []*MetadataItem {
	if m == nil {
		return []*MetadataItem{}
	}

	res := []*MetadataItem{}
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
	if r == nil {
		return nil
	}

	return &AddressValidationData{
		CountryCode:        &r.CountryCode,
		CountryName:        &r.CountryName,
		AddressFormat:      &r.AddressFormat,
		AddressLatinFormat: &r.AddressLatinFormat,
		AllowedFields:      getAllowedFieldsCamelCase(r.AllowedFields),
		RequiredFields:     getFieldsToCamelCase(&r.RequiredFields),
		UpperFields:        getFieldsToCamelCase(&r.UpperFields),
		CountryAreaType:    &r.CountryAreaType,
		CountryAreaChoices: choicesToChoiceValues(r.CountryAreaChoices),
		CityType:           &r.CityType,
		CityChoices:        choicesToChoiceValues(r.CityChoices),
		CityAreaType:       &r.CityAreaType,
		CityAreaChoices:    choicesToChoiceValues(r.CityAreaChoices),
		PostalCodeType:     &r.PostalCodeType,
		PostalCodeMatchers: util.StringSliceToStringPointerSlice(i18naddress.RegexesToStrings(r.PostalCodeMatchers)),
		PostalCodeExamples: util.StringSliceToStringPointerSlice(r.PostalCodeExamples),
		PostalCodePrefix:   &r.PostalCodePrefix,
	}
}

// choicesToChoiceValues convert [][2]string => []*ChoiceValue
func choicesToChoiceValues(choices [][2]string) []*ChoiceValue {
	res := make([]*ChoiceValue, len(choices))

	for i := range choices {
		res[i] = &ChoiceValue{
			Raw:     &choices[i][0],
			Verbose: &choices[i][1],
		}
	}

	return res
}

// toCamelCase converts "the_snake" => "theSnake"
func toCamelCase(snakeStr string) string {
	splitSnake := strings.Split(strings.ToLower(snakeStr), "_")

	res := splitSnake[0]
	if splitSnake[0] == "_" || splitSnake[0] == "" {
		res = "_"
	}

	for _, str := range splitSnake[1:] {
		if trimmed := strings.TrimSpace(str); trimmed == "" {
			res += "_"
		} else if len(trimmed) == 1 {
			res += strings.ToUpper(trimmed)
		} else {
			res += strings.ToUpper(string(trimmed[0])) + trimmed[1:]
		}
	}

	return res
}

func validationFieldToCamelCase(name string) string {
	name = toCamelCase(name)
	if name == "streetAddress" {
		return "streetAddress1"
	}
	return name
}

func getFieldsToCamelCase(fields *[]string) []*string {
	res := make([]*string, len(*fields))

	for i := range *fields {
		res[i] = model.NewString(validationFieldToCamelCase((*fields)[i]))
	}

	return res
}

func getAllowedFieldsCamelCase(allowedFields []string) []*string {
	res := []*string{}

	for i := range allowedFields {
		convStr := validationFieldToCamelCase(allowedFields[i])
		res = append(res, &convStr)
		if convStr == "streetAddress1" {
			res = append(res, model.NewString("streetAddress2"))
		}
	}

	return res
}
