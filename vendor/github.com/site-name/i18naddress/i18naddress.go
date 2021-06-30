package i18naddress

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

// InvalidCodeErr indicate given country code is invalid
type InvalidCodeErr struct {
	value interface{}
	msg   string
}

func (i *InvalidCodeErr) Error() string {
	return fmt.Sprintf(i.msg, i.value)
}

func newInvalidCodeErr(value interface{}) *InvalidCodeErr {
	return &InvalidCodeErr{
		msg:   "%s is not a valid code",
		value: value,
	}
}

var (
	VALID_COUNTRY_CODE   *regexp.Regexp    // regexp for checking country code
	FORMAT_REGEX         *regexp.Regexp    //
	VALIDATION_DATA_PATH string            // path to .json files
	FIELD_MAPPING        map[string]string // short name representations
	KNOWN_FIELDS         []string          //
	json                 jsoniter.API      // fast json
)

func init() {
	json = jsoniter.ConfigCompatibleWithStandardLibrary
	VALID_COUNTRY_CODE = regexp.MustCompile(`^\w{2,3}$`)
	FORMAT_REGEX = regexp.MustCompile(`%([ACDNOSXZ])`)
	VALIDATION_DATA_PATH = "/%s.json"

	FIELD_MAPPING = map[string]string{
		"A": "street_address",
		"C": "city",
		"D": "city_area",
		"N": "name",
		"O": "company_name",
		"S": "country_area",
		"X": "sorting_code",
		"Z": "postal_code",
	}

	KNOWN_FIELDS = []string{"country_code"}
	for _, value := range FIELD_MAPPING {
		KNOWN_FIELDS = append(KNOWN_FIELDS, value)
	}

	sort.Strings(KNOWN_FIELDS)
}

func LoadValidationData(countryCode string) (io.Reader, error) {
	if countryCode == "" {
		countryCode = "all"
	}

	if !VALID_COUNTRY_CODE.MatchString(countryCode) {
		return nil, newInvalidCodeErr(countryCode)
	}

	path := fmt.Sprintf(VALIDATION_DATA_PATH, strings.ToLower(countryCode))

	file, err := assets.Open(path)
	if err != nil {
		return nil, newInvalidCodeErr(countryCode)
	}

	return file, nil
}

func makeChoices(rules map[string]string, translated bool) [][2]string {
	subKeys, ok := rules["sub_keys"]
	if !ok {
		return [][2]string{}
	}

	choices := [][2]string{}
	splitSubKeys := strings.Split(subKeys, "~")

	subNames, ok := rules["sub_names"]
	if ok {
		splitSubNames := strings.Split(subNames, "~")
		for i := 0; i < max(len(splitSubKeys), len(splitSubNames)); i++ {
			if trimmedName := strings.TrimSpace(splitSubNames[i]); trimmedName != "" {
				choices = append(choices, [2]string{splitSubKeys[i], trimmedName})
			}
		}
	} else if !translated {
		for _, key := range splitSubKeys {
			choices = append(choices, [2]string{key, key})
		}
	}

	if !translated {
		subLNames, ok := rules["sub_lnames"]
		if ok {
			splitSubLNames := strings.Split(subLNames, "~")
			for i := 0; i < max(len(splitSubKeys), len(splitSubLNames)); i++ {
				if trimmedName := strings.TrimSpace(splitSubLNames[i]); trimmedName != "" {
					choices = append(choices, [2]string{splitSubKeys[i], trimmedName})
				}
			}
		}

		subLFNames, ok := rules["sub_lfnames"]
		if ok {
			splitSubLFNames := strings.Split(subLFNames, "~")
			for i := 0; i < max(len(splitSubKeys), len(splitSubLFNames)); i++ {
				if trimmedName := strings.TrimSpace(splitSubLFNames[i]); trimmedName != "" {
					choices = append(choices, [2]string{splitSubKeys[i], trimmedName})
				}
			}
		}
	}

	return choices
}

func compactChoices(choices [][2]string) *[][2]string {
	valueMap := make(map[string][]string)
	for _, choice := range choices {
		if _, found := valueMap[choice[0]]; !found {
			valueMap[choice[0]] = []string{}
		}
		valueMap[choice[0]] = append(valueMap[choice[0]], choice[1])
	}

	res := [][2]string{}
	for key, values := range valueMap {
		values = *filterDuplicate(&values)
		for _, value := range values {
			res = append(res, [2]string{key, value})
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i][0] < res[j][0]
	})

	return &res
}

func matchChoices(value string, choices [][2]string) string {
	if value != "" {
		value = strings.TrimSpace(value)
	}
	for _, choice := range choices {
		if strings.EqualFold(choice[0], value) {
			return choice[0]
		}
		if strings.EqualFold(choice[1], value) {
			return choice[0]
		}
	}

	return ""
}

func loadCountryData(countryCode string) (map[string]string, map[string]map[string]string, error) {

	if strings.EqualFold(countryCode, "zz") {
		return nil, nil, newInvalidCodeErr(countryCode)
	}

	countryCode = strings.ToUpper(countryCode)

	reader, err := LoadValidationData("zz")
	if err != nil {
		return nil, nil, err
	}
	database := make(map[string]map[string]string)

	err = json.NewDecoder(reader).Decode(&database)
	if err != nil {
		return nil, nil, err
	}

	countryData := database["ZZ"]

	delete(database, "ZZ") // since this map has 1 key-value pair

	reader, err = LoadValidationData(strings.ToLower(countryCode))
	if err != nil {
		return nil, nil, err
	}

	err = json.NewDecoder(reader).Decode(&database)
	if err != nil {
		return nil, nil, err
	}
	for key, value := range database[countryCode] {
		countryData[key] = value
	}

	return countryData, database, nil
}

type Params struct {
	CountryArea   string
	CountryCode   string
	City          string
	CityArea      string
	PostalCode    string
	StreetAddress string
	SortingCode   string
}

// GetValidationRules get validation rules for given params.
//
// there are 3 error types might be returned:
//
// nil, *InvalidCodeErr, json (encoding/decoding) error
func GetValidationRules(address *Params) (*ValidationRules, error) {
	countryCode := strings.ToUpper(address.CountryCode)
	countryData, database, err := loadCountryData(countryCode)
	if err != nil {
		return nil, err
	}
	countryName := countryData["name"]
	addressFormat := countryData["fmt"]
	addressLatinFormat, ok := countryData["lfmt"]
	if !ok {
		addressLatinFormat = addressFormat
	}

	formatFields := FORMAT_REGEX.FindAll([]byte(addressFormat), -1)

	allowedFields := []string{}
	for _, field := range formatFields {
		allowedFields = append(allowedFields, FIELD_MAPPING[string(field[1])])
	}
	allowedFields = *(filterDuplicate(&allowedFields))

	requiredFields := []string{}
	for _, r := range countryData["require"] {
		requiredFields = append(requiredFields, FIELD_MAPPING[string(r)])
	}
	requiredFields = *(filterDuplicate(&requiredFields))

	upperFields := []string{}
	for _, r := range countryData["upper"] {
		upperFields = append(upperFields, FIELD_MAPPING[string(r)])
	}
	upperFields = *(filterDuplicate(&upperFields))

	var languages []string
	if lingos, exist := countryData["languages"]; exist {
		languages = strings.Split(lingos, "~")
	}

	postalCodeMatchers := []*regexp.Regexp{}
	if stringInSlice("postal_code", &allowedFields) {
		if zip, exist := countryData["zip"]; exist {
			postalCodeMatchers = append(postalCodeMatchers, regexp.MustCompile(fmt.Sprintf("^%s$", zip)))
		}
	}

	var postalCodeExamples []string
	if zipex, exist := countryData["zipex"]; exist {
		postalCodeExamples = strings.Split(zipex, ",")
	}

	var (
		countryAreaChoices = [][2]string{}
		cityChoices        = [][2]string{}
		cityAreaChoices    = [][2]string{}
		countryAreaType    = countryData["state_name_type"]
		cityType           = countryData["locality_name_type"]
		cityAreaType       = countryData["sublocality_name_type"]
		postalCodeType     = countryData["zip_name_type"]
		postalCodePrefix   = countryData["postprefix"] // if not exist, value will be ""
		countryArea        string
		city               string
		cityArea           string
	)

	if countryValue, ok := database[countryCode]; ok {
		if _, ok = countryData["sub_keys"]; ok {
			for _, language := range languages {
				isDefaultLanguage := language == countryData["lang"] && strings.TrimSpace(language) != ""

				localizedCountryData := countryValue
				if !isDefaultLanguage {
					localizedCountryData = database[fmt.Sprintf("%s--%s", countryCode, language)]
				}
				localizedCountryAreaChoices := makeChoices(localizedCountryData, false)
				countryAreaChoices = append(countryAreaChoices, localizedCountryAreaChoices...)
				existingChoice := countryArea != ""
				matchedCountryArea := matchChoices(address.CountryArea, localizedCountryAreaChoices)
				countryArea = matchedCountryArea

				var matchedCity string

				if matchedCountryArea != "" {
					// 3rd level of data is for cities
					var countryAreaData map[string]string
					if isDefaultLanguage {
						countryAreaData = database[fmt.Sprintf("%s/%s", countryCode, countryArea)]
					} else {
						countryAreaData = database[fmt.Sprintf("%s/%s--%s", countryCode, countryArea, language)]
					}

					if !existingChoice {
						if zip, exist := countryAreaData["zip"]; exist {
							postalCodeMatchers = append(postalCodeMatchers, regexp.MustCompile("^"+zip))
						}
						if zipex, exist := countryAreaData["zipex"]; exist {
							postalCodeExamples = strings.Split(zipex, ",")
						}
					}

					if _, exist := countryAreaData["sub_keys"]; exist {
						localizedCityChoices := makeChoices(countryAreaData, false)
						cityChoices = append(cityChoices, localizedCityChoices...)
						existingChoice = city != ""
						matchedCity = matchChoices(address.City, localizedCityChoices)
						city = matchedCity
					}

					if matchedCity != "" {
						// 4th level of data is for dependent sublocalities
						cityData := database[fmt.Sprintf("%s/%s/%s", countryCode, countryArea, city)]
						if !isDefaultLanguage {
							cityData = database[fmt.Sprintf("%s/%s/%s--%s", countryCode, countryArea, city, language)]
						}

						if !existingChoice {
							if zip, exist := cityData["zip"]; exist {
								postalCodeMatchers = append(postalCodeMatchers, regexp.MustCompile("^"+zip))
							}
							if zipex, exist := cityData["zipex"]; exist {
								postalCodeExamples = strings.Split(zipex, ",")
							}
						}

						if _, exist := cityData["sub_keys"]; exist {
							localizedCityAreaChoices := makeChoices(cityData, false)
							cityAreaChoices = append(cityAreaChoices, localizedCityAreaChoices...)
							existingChoice = cityArea != ""
							matchedCityArea := matchChoices(address.CityArea, localizedCityAreaChoices)

							if matchedCityArea != "" {
								cityAreaData := database[fmt.Sprintf("%s/%s/%s/%s", countryCode, countryArea, city, matchedCityArea)]
								if !isDefaultLanguage {
									cityAreaData = database[fmt.Sprintf("%s/%s/%s/%s--%s", countryCode, countryArea, city, matchedCityArea, language)]
								}

								if !existingChoice {
									if zip, exist := cityAreaData["zip"]; exist {
										postalCodeMatchers = append(postalCodeMatchers, regexp.MustCompile("^"+zip))
									}
									if zipex, exist := cityAreaData["zipex"]; exist {
										postalCodeExamples = strings.Split(zipex, ",")
									}
								}
							}
						}
					}
				}
			}
		}

		countryAreaChoices = *(compactChoices(countryAreaChoices))
		cityChoices = *(compactChoices(cityChoices))
		cityAreaChoices = *(compactChoices(cityAreaChoices))
	}

	return &ValidationRules{
		CountryCode:        countryCode,
		CountryName:        countryName,
		AddressFormat:      addressFormat,
		AddressLatinFormat: addressLatinFormat,
		AllowedFields:      allowedFields,
		RequiredFields:     requiredFields,
		UpperFields:        upperFields,
		CountryAreaType:    countryAreaType,
		CountryAreaChoices: countryAreaChoices,
		CityType:           cityType,
		CityChoices:        cityChoices,
		CityAreaType:       cityAreaType,
		CityAreaChoices:    cityAreaChoices,
		PostalCodeType:     postalCodeType,
		PostalCodeMatchers: postalCodeMatchers,
		PostalCodeExamples: postalCodeExamples,
		PostalCodePrefix:   postalCodePrefix,
	}, nil
}

type ValidationRules struct {
	CountryCode        string
	CountryName        string
	AddressFormat      string
	AddressLatinFormat string
	AllowedFields      []string
	RequiredFields     []string
	UpperFields        []string
	CountryAreaType    string
	CountryAreaChoices [][2]string
	CityType           string
	CityChoices        [][2]string
	CityAreaType       string
	CityAreaChoices    [][2]string
	PostalCodeType     string
	PostalCodeMatchers []*regexp.Regexp
	PostalCodeExamples []string
	PostalCodePrefix   string
}

func (v *ValidationRules) String() string {
	return fmt.Sprintf(
		"ValidationRules("+
			"country_code=%v, "+
			"country_name=%v, "+
			"address_format=%v, "+
			"address_latin_format=%v, "+
			"allowed_fields=%v, "+
			"required_fields=%v, "+
			"upper_fields=%v, "+
			"country_area_type=%v, "+
			"country_area_choices=%v, "+
			"city_type=%v, "+
			"city_choices=%v, "+
			"city_area_type=%v, "+
			"city_area_choices=%v, "+
			"postal_code_type=%v, "+
			"postal_code_matchers=%v, "+
			"postal_code_examples=%v, "+
			"postal_code_prefix=%v)",
		v.CountryCode,
		v.CountryName,
		v.AddressFormat,
		v.AddressLatinFormat,
		v.AllowedFields,
		v.RequiredFields,
		v.UpperFields,
		v.CountryAreaType,
		v.CountryAreaChoices,
		v.CityType,
		v.CityChoices,
		v.CityAreaType,
		v.CityAreaChoices,
		v.PostalCodeType,
		*RegexesToStrings(v.PostalCodeMatchers),
		v.PostalCodeExamples,
		v.PostalCodePrefix)
}

func (p *Params) GetProperty(name string) string {
	switch strings.ToLower(name) {
	case "country_area":
		return p.CountryArea
	case "city":
		return p.City
	case "city_area":
		return p.CityArea
	case "country_code":
		return p.CountryCode
	case "postal_code":
		return p.PostalCode
	case "street_address":
		return p.StreetAddress
	case "sorting_code":
		return p.SortingCode
	default:
		return ""
	}
}

func (p *Params) Patch(name string, value string) {
	switch strings.ToLower(name) {
	case "country_area":
		p.CountryArea = value
	case "city":
		p.City = value
	case "city_area":
		p.CityArea = value
	case "country_code":
		p.CountryCode = value
	case "postal_code":
		p.PostalCode = value
	case "street_address":
		p.StreetAddress = value
	case "sorting_code":
		p.SortingCode = value
	default:
		return
	}
}

// normalizeField
func normalizeField(name string, rules *ValidationRules, data *Params, choices [][2]string, errors map[string]string) {
	value := data.GetProperty(name)

	if stringInSlice(name, &rules.UpperFields) && value != "" {
		value = strings.ToUpper(value)
		data.Patch(name, value)
	}

	if !stringInSlice(name, &rules.AllowedFields) {
		data.Patch(name, "")
	} else if value == "" && stringInSlice(name, &rules.RequiredFields) {
		errors[name] = "required"
	} else if len(choices) > 0 {
		if value != "" || stringInSlice(name, &rules.RequiredFields) {
			value = matchChoices(value, choices)
			if value != "" {
				data.Patch(name, value)
			} else {
				errors[name] = "invalid"
			}
		}
	}

	if value == "" {
		data.Patch(name, value)
	}
}

func (p *Params) Copy() *Params {
	var newP Params = *p
	return &newP
}

func NormalizeAddress(address *Params) (p *Params, errorMap map[string]string) {
	errors := make(map[string]string)

	rules, err := GetValidationRules(address)
	if err != nil {
		errors["country_code"] = "invalid"
	}

	cleanedData := address.Copy()
	if cleanedData.CountryCode == "" {
		errors["country_code"] = "required"
	} else {
		cleanedData.CountryCode = strings.ToUpper(cleanedData.CountryCode)
	}
	normalizeField("country_area", rules, cleanedData, rules.CountryAreaChoices, errors)
	normalizeField("city", rules, cleanedData, rules.CityChoices, errors)
	normalizeField("city_area", rules, cleanedData, rules.CityAreaChoices, errors)
	normalizeField("postal_code", rules, cleanedData, [][2]string{}, errors)

	if cleanedData.PostalCode != "" && len(rules.PostalCodeMatchers) > 0 {
		for _, matcher := range rules.PostalCodeMatchers {
			if !matcher.MatchString(cleanedData.PostalCode) {
				errors["postal_code"] = "invalid"
				break
			}
		}
	}

	normalizeField("street_address", rules, cleanedData, [][2]string{}, errors)
	normalizeField("sorting_code", rules, cleanedData, [][2]string{}, errors)

	if len(errors) > 0 {
		return nil, errors
	}

	return cleanedData, nil
}

func formatAddressLine(lineFormat string, address *Params, rules *ValidationRules) string {
	getField := func(name string) string {
		value := address.GetProperty(name)
		if stringInSlice(name, &rules.UpperFields) {
			value = strings.ToUpper(value)
		}

		return value
	}

	replacements := make(map[string]string, len(FIELD_MAPPING))
	for key, value := range FIELD_MAPPING {
		replacements[fmt.Sprintf("%%%s", key)] = getField(value)
	}

	fields := regexp.MustCompile("(%.)").Split(lineFormat, -1)
	for i, field := range fields {
		if repl, exist := replacements[field]; exist {
			fields[i] = repl
		}
	}

	return strings.TrimSpace(strings.Join(fields, ""))
}

func GetFieldOrder(address *Params, latin bool) ([][]string, error) {
	rules, err := GetValidationRules(address)
	if err != nil {
		return nil, err
	}
	addressFormat := rules.AddressLatinFormat
	if !latin {
		addressFormat = rules.AddressFormat
	}

	addressLines := strings.Split(addressFormat, "%n")
	replacements := make(map[string]string, len(FIELD_MAPPING))
	for key, value := range FIELD_MAPPING {
		replacements[fmt.Sprintf("%%%s", key)] = value
	}

	allLines := [][]string{}
	for _, line := range addressLines {
		fields := regexp.MustCompile("(%.)").Split(line, -1)
		singleLine := []string{}
		for _, field := range fields {
			singleLine = append(singleLine, replacements[field])
		}
		singleLine = filterSlice(singleLine, func(s string) bool {
			return s != ""
		})
		allLines = append(allLines, singleLine)
	}

	return allLines, nil
}

func FormatAddress(address *Params, latin bool) (string, error) {
	rules, err := GetValidationRules(address)
	if err != nil {
		return "", err
	}

	addressFormat := rules.AddressLatinFormat
	if !latin {
		addressFormat = rules.AddressFormat
	}

	addressLines := []string{}
	for _, lf := range strings.Split(addressFormat, "%n") {
		addressLines = append(addressLines, formatAddressLine(lf, address, rules))
	}
	addressLines = append(addressLines, rules.CountryName)
	addressLines = filterSlice(addressLines, func(s string) bool {
		return s != ""
	})

	return strings.Join(addressLines, "\n"), nil
}

type ErrorMap map[string]string

func (e ErrorMap) Error() string {
	if e == nil || len(e) == 0 {
		return "{}"
	}
	b, _ := json.Marshal(e)
	return string(b)
}

// func LatinizeAddress(address *Params, normalized bool) (interface{}, error) {
// 	if !normalized {
// 		address, errMap := NormalizeAddress(address)
// 		if errMap != nil || len(errMap) > 0 {
// 			return nil, ErrorMap(errMap)
// 		}
// 	}

// 	cleanedData := address.Copy()
// 	countryCode := address.GetProperty("country_code")
// 	countryCode = strings.ToUpper(countryCode)
// 	_, database, err := loadCountryData(countryCode)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if countryCode != "" {
// 		countryArea := address.GetProperty("country_area")
// 		if countryArea != "" {
// 			key = countryCode + "/" + countryArea
// 			country
// 		}
// 	}
// }
