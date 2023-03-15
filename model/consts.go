package model

import (
	"regexp"
	"unsafe"
)

// constants for access http(s) requests's headers
const (
	HeaderRequestId            = "X-Request-ID"
	HeaderVersionId            = "X-Version-ID"
	HEADER_CLUSTER_ID          = "X-Cluster-ID"
	HEADER_ETAG_SERVER         = "ETag"
	HEADER_ETAG_CLIENT         = "If-None-Match"
	HEADER_FORWARDED           = "X-Forwarded-For"
	HEADER_REAL_IP             = "X-Real-IP"
	HEADER_FORWARDED_PROTO     = "X-Forwarded-Proto"
	HEADER_TOKEN               = "token"
	HEADER_CSRF_TOKEN          = "X-CSRF-Token"
	HEADER_BEARER              = "BEARER"
	HEADER_AUTH                = "Authorization"
	HEADER_CLOUD_TOKEN         = "X-Cloud-Token"
	HEADER_REMOTECLUSTER_TOKEN = "X-RemoteCluster-Token"
	HEADER_REMOTECLUSTER_ID    = "X-RemoteCluster-Id"
	HEADER_REQUESTED_WITH      = "X-Requested-With"
	HEADER_REQUESTED_WITH_XML  = "XMLHttpRequest"
	HEADER_RANGE               = "Range"
	STATUS                     = "status"
	STATUS_OK                  = "OK"
	STATUS_FAIL                = "FAIL"
	STATUS_UNHEALTHY           = "UNHEALTHY"
	STATUS_REMOVE              = "REMOVE"
	CLIENT_DIR                 = "client"
)

type TimePeriodType string

// time period types
const (
	DAY   TimePeriodType = "day"
	WEEK  TimePeriodType = "week"
	MONTH TimePeriodType = "month"
	YEAR  TimePeriodType = "year"
)

var TimePeriodMap = map[TimePeriodType]string{
	DAY:   "Day",
	WEEK:  "Week",
	MONTH: "Month",
	YEAR:  "Year",
}

// some default values for model fields
const (
	TimeZone                       = "UTC"
	USER_AUTH_SERVICE_EMAIL        = "email"
	DEFAULT_CURRENCY               = "USD"
	USER_NAME_MAX_LENGTH           = 64
	USER_EMAIL_MAX_LENGTH          = 128
	USER_NAME_MIN_LENGTH           = 1
	CURRENCY_CODE_MAX_LENGTH       = 3
	LANGUAGE_CODE_MAX_LENGTH       = 5
	WEIGHT_UNIT_MAX_LENGTH         = 5
	URL_LINK_MAX_LENGTH            = 200
	SINGLE_COUNTRY_CODE_MAX_LENGTH = 5
	IP_ADDRESS_MAX_LENGTH          = 39
	DEFAULT_LOCALE                 = "en" // this is default language also
	DEFAULT_COUNTRY                = CountryCodeUs
)

var (
	Countries                     map[CountryCode]string // countries supported by app
	Languages                     map[string]string      // Languages supported by app
	MULTIPLE_COUNTRIES_MAX_LENGTH int                    // some model"s country fields contains multiple countries
	ReservedName                  []string               // usernames that can only be used by system
	ValidUsernameChars            *regexp.Regexp         // regexp for username validation
	RestrictedUsernames           map[string]bool        // usernames that cannot be used
)

type CountryCode string

func (c CountryCode) IsValid() bool {
	return Countries[c] != ""
}

func (c CountryCode) String() string {
	return *(*string)(unsafe.Pointer(&c))
}

const (
	CountryCodeAf CountryCode = "AF"
	CountryCodeId CountryCode = "ID"
	CountryCodeAx CountryCode = "AX"
	CountryCodeAl CountryCode = "AL"
	CountryCodeDz CountryCode = "DZ"
	CountryCodeAs CountryCode = "AS"
	CountryCodeAd CountryCode = "AD"
	CountryCodeAo CountryCode = "AO"
	CountryCodeAi CountryCode = "AI"
	CountryCodeAq CountryCode = "AQ"
	CountryCodeAg CountryCode = "AG"
	CountryCodeAr CountryCode = "AR"
	CountryCodeAm CountryCode = "AM"
	CountryCodeAw CountryCode = "AW"
	CountryCodeAu CountryCode = "AU"
	CountryCodeAt CountryCode = "AT"
	CountryCodeAz CountryCode = "AZ"
	CountryCodeBs CountryCode = "BS"
	CountryCodeBh CountryCode = "BH"
	CountryCodeBd CountryCode = "BD"
	CountryCodeBb CountryCode = "BB"
	CountryCodeBy CountryCode = "BY"
	CountryCodeBe CountryCode = "BE"
	CountryCodeBz CountryCode = "BZ"
	CountryCodeBj CountryCode = "BJ"
	CountryCodeBm CountryCode = "BM"
	CountryCodeBt CountryCode = "BT"
	CountryCodeBo CountryCode = "BO"
	CountryCodeBq CountryCode = "BQ"
	CountryCodeBa CountryCode = "BA"
	CountryCodeBw CountryCode = "BW"
	CountryCodeBv CountryCode = "BV"
	CountryCodeBr CountryCode = "BR"
	CountryCodeIo CountryCode = "IO"
	CountryCodeBn CountryCode = "BN"
	CountryCodeBg CountryCode = "BG"
	CountryCodeBf CountryCode = "BF"
	CountryCodeBi CountryCode = "BI"
	CountryCodeCv CountryCode = "CV"
	CountryCodeKh CountryCode = "KH"
	CountryCodeCm CountryCode = "CM"
	CountryCodeCa CountryCode = "CA"
	CountryCodeKy CountryCode = "KY"
	CountryCodeCf CountryCode = "CF"
	CountryCodeTd CountryCode = "TD"
	CountryCodeCl CountryCode = "CL"
	CountryCodeCn CountryCode = "CN"
	CountryCodeCx CountryCode = "CX"
	CountryCodeCc CountryCode = "CC"
	CountryCodeCo CountryCode = "CO"
	CountryCodeKm CountryCode = "KM"
	CountryCodeCg CountryCode = "CG"
	CountryCodeCd CountryCode = "CD"
	CountryCodeCk CountryCode = "CK"
	CountryCodeCr CountryCode = "CR"
	CountryCodeCi CountryCode = "CI"
	CountryCodeHr CountryCode = "HR"
	CountryCodeCu CountryCode = "CU"
	CountryCodeCw CountryCode = "CW"
	CountryCodeCy CountryCode = "CY"
	CountryCodeCz CountryCode = "CZ"
	CountryCodeDk CountryCode = "DK"
	CountryCodeDj CountryCode = "DJ"
	CountryCodeDm CountryCode = "DM"
	CountryCodeDo CountryCode = "DO"
	CountryCodeEc CountryCode = "EC"
	CountryCodeEg CountryCode = "EG"
	CountryCodeSv CountryCode = "SV"
	CountryCodeGq CountryCode = "GQ"
	CountryCodeEr CountryCode = "ER"
	CountryCodeEe CountryCode = "EE"
	CountryCodeSz CountryCode = "SZ"
	CountryCodeEt CountryCode = "ET"
	CountryCodeEu CountryCode = "EU"
	CountryCodeFk CountryCode = "FK"
	CountryCodeFo CountryCode = "FO"
	CountryCodeFj CountryCode = "FJ"
	CountryCodeFi CountryCode = "FI"
	CountryCodeFr CountryCode = "FR"
	CountryCodeGf CountryCode = "GF"
	CountryCodePf CountryCode = "PF"
	CountryCodeTf CountryCode = "TF"
	CountryCodeGa CountryCode = "GA"
	CountryCodeGm CountryCode = "GM"
	CountryCodeGe CountryCode = "GE"
	CountryCodeDe CountryCode = "DE"
	CountryCodeGh CountryCode = "GH"
	CountryCodeGi CountryCode = "GI"
	CountryCodeGr CountryCode = "GR"
	CountryCodeGl CountryCode = "GL"
	CountryCodeGd CountryCode = "GD"
	CountryCodeGp CountryCode = "GP"
	CountryCodeGu CountryCode = "GU"
	CountryCodeGt CountryCode = "GT"
	CountryCodeGg CountryCode = "GG"
	CountryCodeGn CountryCode = "GN"
	CountryCodeGw CountryCode = "GW"
	CountryCodeGy CountryCode = "GY"
	CountryCodeHt CountryCode = "HT"
	CountryCodeHm CountryCode = "HM"
	CountryCodeVa CountryCode = "VA"
	CountryCodeHn CountryCode = "HN"
	CountryCodeHk CountryCode = "HK"
	CountryCodeHu CountryCode = "HU"
	CountryCodeIs CountryCode = "IS"
	CountryCodeIn CountryCode = "IN"
	CountryCodeIr CountryCode = "IR"
	CountryCodeIq CountryCode = "IQ"
	CountryCodeIe CountryCode = "IE"
	CountryCodeIm CountryCode = "IM"
	CountryCodeIl CountryCode = "IL"
	CountryCodeIt CountryCode = "IT"
	CountryCodeJm CountryCode = "JM"
	CountryCodeJp CountryCode = "JP"
	CountryCodeJe CountryCode = "JE"
	CountryCodeJo CountryCode = "JO"
	CountryCodeKz CountryCode = "KZ"
	CountryCodeKe CountryCode = "KE"
	CountryCodeKi CountryCode = "KI"
	CountryCodeKw CountryCode = "KW"
	CountryCodeKg CountryCode = "KG"
	CountryCodeLa CountryCode = "LA"
	CountryCodeLv CountryCode = "LV"
	CountryCodeLb CountryCode = "LB"
	CountryCodeLs CountryCode = "LS"
	CountryCodeLr CountryCode = "LR"
	CountryCodeLy CountryCode = "LY"
	CountryCodeLi CountryCode = "LI"
	CountryCodeLt CountryCode = "LT"
	CountryCodeLu CountryCode = "LU"
	CountryCodeMo CountryCode = "MO"
	CountryCodeMg CountryCode = "MG"
	CountryCodeMw CountryCode = "MW"
	CountryCodeMy CountryCode = "MY"
	CountryCodeMv CountryCode = "MV"
	CountryCodeMl CountryCode = "ML"
	CountryCodeMt CountryCode = "MT"
	CountryCodeMh CountryCode = "MH"
	CountryCodeMq CountryCode = "MQ"
	CountryCodeMr CountryCode = "MR"
	CountryCodeMu CountryCode = "MU"
	CountryCodeYt CountryCode = "YT"
	CountryCodeMx CountryCode = "MX"
	CountryCodeFm CountryCode = "FM"
	CountryCodeMd CountryCode = "MD"
	CountryCodeMc CountryCode = "MC"
	CountryCodeMn CountryCode = "MN"
	CountryCodeMe CountryCode = "ME"
	CountryCodeMs CountryCode = "MS"
	CountryCodeMa CountryCode = "MA"
	CountryCodeMz CountryCode = "MZ"
	CountryCodeMm CountryCode = "MM"
	CountryCodeNa CountryCode = "NA"
	CountryCodeNr CountryCode = "NR"
	CountryCodeNp CountryCode = "NP"
	CountryCodeNl CountryCode = "NL"
	CountryCodeNc CountryCode = "NC"
	CountryCodeNz CountryCode = "NZ"
	CountryCodeNi CountryCode = "NI"
	CountryCodeNe CountryCode = "NE"
	CountryCodeNg CountryCode = "NG"
	CountryCodeNu CountryCode = "NU"
	CountryCodeNf CountryCode = "NF"
	CountryCodeKp CountryCode = "KP"
	CountryCodeMk CountryCode = "MK"
	CountryCodeMp CountryCode = "MP"
	CountryCodeNo CountryCode = "NO"
	CountryCodeOm CountryCode = "OM"
	CountryCodePk CountryCode = "PK"
	CountryCodePw CountryCode = "PW"
	CountryCodePs CountryCode = "PS"
	CountryCodePa CountryCode = "PA"
	CountryCodePg CountryCode = "PG"
	CountryCodePy CountryCode = "PY"
	CountryCodePe CountryCode = "PE"
	CountryCodePh CountryCode = "PH"
	CountryCodePn CountryCode = "PN"
	CountryCodePl CountryCode = "PL"
	CountryCodePt CountryCode = "PT"
	CountryCodePr CountryCode = "PR"
	CountryCodeQa CountryCode = "QA"
	CountryCodeRe CountryCode = "RE"
	CountryCodeRo CountryCode = "RO"
	CountryCodeRu CountryCode = "RU"
	CountryCodeRw CountryCode = "RW"
	CountryCodeBl CountryCode = "BL"
	CountryCodeSh CountryCode = "SH"
	CountryCodeKn CountryCode = "KN"
	CountryCodeLc CountryCode = "LC"
	CountryCodeMf CountryCode = "MF"
	CountryCodePm CountryCode = "PM"
	CountryCodeVc CountryCode = "VC"
	CountryCodeWs CountryCode = "WS"
	CountryCodeSm CountryCode = "SM"
	CountryCodeSt CountryCode = "ST"
	CountryCodeSa CountryCode = "SA"
	CountryCodeSn CountryCode = "SN"
	CountryCodeRs CountryCode = "RS"
	CountryCodeSc CountryCode = "SC"
	CountryCodeSl CountryCode = "SL"
	CountryCodeSg CountryCode = "SG"
	CountryCodeSx CountryCode = "SX"
	CountryCodeSk CountryCode = "SK"
	CountryCodeSi CountryCode = "SI"
	CountryCodeSb CountryCode = "SB"
	CountryCodeSo CountryCode = "SO"
	CountryCodeZa CountryCode = "ZA"
	CountryCodeGs CountryCode = "GS"
	CountryCodeKr CountryCode = "KR"
	CountryCodeSs CountryCode = "SS"
	CountryCodeEs CountryCode = "ES"
	CountryCodeLk CountryCode = "LK"
	CountryCodeSd CountryCode = "SD"
	CountryCodeSr CountryCode = "SR"
	CountryCodeSj CountryCode = "SJ"
	CountryCodeSe CountryCode = "SE"
	CountryCodeCh CountryCode = "CH"
	CountryCodeSy CountryCode = "SY"
	CountryCodeTw CountryCode = "TW"
	CountryCodeTj CountryCode = "TJ"
	CountryCodeTz CountryCode = "TZ"
	CountryCodeTh CountryCode = "TH"
	CountryCodeTl CountryCode = "TL"
	CountryCodeTg CountryCode = "TG"
	CountryCodeTk CountryCode = "TK"
	CountryCodeTo CountryCode = "TO"
	CountryCodeTt CountryCode = "TT"
	CountryCodeTn CountryCode = "TN"
	CountryCodeTr CountryCode = "TR"
	CountryCodeTm CountryCode = "TM"
	CountryCodeTc CountryCode = "TC"
	CountryCodeTv CountryCode = "TV"
	CountryCodeUg CountryCode = "UG"
	CountryCodeUa CountryCode = "UA"
	CountryCodeAe CountryCode = "AE"
	CountryCodeGb CountryCode = "GB"
	CountryCodeUm CountryCode = "UM"
	CountryCodeUs CountryCode = "US"
	CountryCodeUy CountryCode = "UY"
	CountryCodeUz CountryCode = "UZ"
	CountryCodeVu CountryCode = "VU"
	CountryCodeVe CountryCode = "VE"
	CountryCodeVn CountryCode = "VN"
	CountryCodeVg CountryCode = "VG"
	CountryCodeVi CountryCode = "VI"
	CountryCodeWf CountryCode = "WF"
	CountryCodeEh CountryCode = "EH"
	CountryCodeYe CountryCode = "YE"
	CountryCodeZm CountryCode = "ZM"
	CountryCodeZw CountryCode = "ZW"
)

func init() {
	// borrowed from django_countries
	Countries = map[CountryCode]string{
		CountryCodeAf: "Afghanistan",
		CountryCodeAx: "Aland Islands",
		CountryCodeAl: "Albania",
		CountryCodeDz: "Algeria",
		CountryCodeAs: "American Samoa",
		CountryCodeAd: "Andorra",
		CountryCodeAo: "Angola",
		CountryCodeAi: "Anguilla",
		CountryCodeAq: "Antarctica",
		CountryCodeAg: "Antigua and Barbuda",
		CountryCodeAr: "Argentina",
		CountryCodeAm: "Armenia",
		CountryCodeAw: "Aruba",
		CountryCodeAu: "Australia",
		CountryCodeAt: "Austria",
		CountryCodeAz: "Azerbaijan",
		CountryCodeBs: "Bahamas",
		CountryCodeBh: "Bahrain",
		CountryCodeBd: "Bangladesh",
		CountryCodeBb: "Barbados",
		CountryCodeBy: "Belarus",
		CountryCodeBe: "Belgium",
		CountryCodeBz: "Belize",
		CountryCodeBj: "Benin",
		CountryCodeBm: "Bermuda",
		CountryCodeBt: "Bhutan",
		CountryCodeBo: "Bolivia",
		CountryCodeBq: "Bonaire, Sint Eustatius and Saba",
		CountryCodeBa: "Bosnia and Herzegovina",
		CountryCodeBw: "Botswana",
		CountryCodeBv: "Bouvet Island",
		CountryCodeBr: "Brazil",
		CountryCodeIo: "British Indian Ocean Territory",
		CountryCodeBn: "Brunei",
		CountryCodeBg: "Bulgaria",
		CountryCodeBf: "Burkina Faso",
		CountryCodeBi: "Burundi",
		CountryCodeCv: "Cabo Verde",
		CountryCodeKh: "Cambodia",
		CountryCodeCm: "Cameroon",
		CountryCodeCa: "Canada",
		CountryCodeKy: "Cayman Islands",
		CountryCodeCf: "Central African Republic",
		CountryCodeTd: "Chad",
		CountryCodeCl: "Chile",
		CountryCodeCn: "China",
		CountryCodeCx: "Christmas Island",
		CountryCodeCc: "Cocos (Keeling) Islands",
		CountryCodeCo: "Colombia",
		CountryCodeKm: "Comoros",
		CountryCodeCg: "Congo",
		CountryCodeCd: "Congo (the Democratic Republic of the)",
		CountryCodeCk: "Cook Islands",
		CountryCodeCr: "Costa Rica",
		CountryCodeCi: "C\u00f4te d'Ivoire",
		CountryCodeHr: "Croatia",
		CountryCodeCu: "Cuba",
		CountryCodeCw: "Cura\u00e7ao",
		CountryCodeCy: "Cyprus",
		CountryCodeCz: "Czechia",
		CountryCodeDk: "Denmark",
		CountryCodeDj: "Djibouti",
		CountryCodeDm: "Dominica",
		CountryCodeDo: "Dominican Republic",
		CountryCodeEc: "Ecuador",
		CountryCodeEg: "Egypt",
		CountryCodeSv: "El Salvador",
		CountryCodeGq: "Equatorial Guinea",
		CountryCodeEr: "Eritrea",
		CountryCodeEe: "Estonia",
		CountryCodeSz: "Eswatini",
		CountryCodeEt: "Ethiopia",
		CountryCodeFk: "Falkland Islands (Malvinas)",
		CountryCodeFo: "Faroe Islands",
		CountryCodeFj: "Fiji",
		CountryCodeFi: "Finland",
		CountryCodeFr: "France",
		CountryCodeGf: "French Guiana",
		CountryCodePf: "French Polynesia",
		CountryCodeTf: "French Southern Territories",
		CountryCodeGa: "Gabon",
		CountryCodeGm: "Gambia",
		CountryCodeGe: "Georgia",
		CountryCodeDe: "Germany",
		CountryCodeGh: "Ghana",
		CountryCodeGi: "Gibraltar",
		CountryCodeGr: "Greece",
		CountryCodeGl: "Greenland",
		CountryCodeGd: "Grenada",
		CountryCodeGp: "Guadeloupe",
		CountryCodeGu: "Guam",
		CountryCodeGt: "Guatemala",
		CountryCodeGg: "Guernsey",
		CountryCodeGn: "Guinea",
		CountryCodeGw: "Guinea-Bissau",
		CountryCodeGy: "Guyana",
		CountryCodeHt: "Haiti",
		CountryCodeHm: "Heard Island and McDonald Islands",
		CountryCodeVa: "Holy See",
		CountryCodeHn: "Honduras",
		CountryCodeHk: "Hong Kong",
		CountryCodeHu: "Hungary",
		CountryCodeIs: "Iceland",
		CountryCodeIn: "India",
		CountryCodeId: "Indonesia",
		CountryCodeIr: "Iran",
		CountryCodeIq: "Iraq",
		CountryCodeIe: "Ireland",
		CountryCodeIm: "Isle of Man",
		CountryCodeIl: "Israel",
		CountryCodeIt: "Italy",
		CountryCodeJm: "Jamaica",
		CountryCodeJp: "Japan",
		CountryCodeJe: "Jersey",
		CountryCodeJo: "Jordan",
		CountryCodeKz: "Kazakhstan",
		CountryCodeKe: "Kenya",
		CountryCodeKi: "Kiribati",
		CountryCodeKp: "North Korea",
		CountryCodeKr: "South Korea",
		CountryCodeKw: "Kuwait",
		CountryCodeKg: "Kyrgyzstan",
		CountryCodeLa: "Laos",
		CountryCodeLv: "Latvia",
		CountryCodeLb: "Lebanon",
		CountryCodeLs: "Lesotho",
		CountryCodeLr: "Liberia",
		CountryCodeLy: "Libya",
		CountryCodeLi: "Liechtenstein",
		CountryCodeLt: "Lithuania",
		CountryCodeLu: "Luxembourg",
		CountryCodeMo: "Macao",
		CountryCodeMg: "Madagascar",
		CountryCodeMw: "Malawi",
		CountryCodeMy: "Malaysia",
		CountryCodeMv: "Maldives",
		CountryCodeMl: "Mali",
		CountryCodeMt: "Malta",
		CountryCodeMh: "Marshall Islands",
		CountryCodeMq: "Martinique",
		CountryCodeMr: "Mauritania",
		CountryCodeMu: "Mauritius",
		CountryCodeYt: "Mayotte",
		CountryCodeMx: "Mexico",
		CountryCodeFm: "Micronesia (Federated States of)",
		CountryCodeMd: "Moldova",
		CountryCodeMc: "Monaco",
		CountryCodeMn: "Mongolia",
		CountryCodeMe: "Montenegro",
		CountryCodeMs: "Montserrat",
		CountryCodeMa: "Morocco",
		CountryCodeMz: "Mozambique",
		CountryCodeMm: "Myanmar",
		CountryCodeNa: "Namibia",
		CountryCodeNr: "Nauru",
		CountryCodeNp: "Nepal",
		CountryCodeNl: "Netherlands",
		CountryCodeNc: "New Caledonia",
		CountryCodeNz: "New Zealand",
		CountryCodeNi: "Nicaragua",
		CountryCodeNe: "Niger",
		CountryCodeNg: "Nigeria",
		CountryCodeNu: "Niue",
		CountryCodeNf: "Norfolk Island",
		CountryCodeMk: "North Macedonia",
		CountryCodeMp: "Northern Mariana Islands",
		CountryCodeNo: "Norway",
		CountryCodeOm: "Oman",
		CountryCodePk: "Pakistan",
		CountryCodePw: "Palau",
		CountryCodePs: "Palestine, State of",
		CountryCodePa: "Panama",
		CountryCodePg: "Papua New Guinea",
		CountryCodePy: "Paraguay",
		CountryCodePe: "Peru",
		CountryCodePh: "Philippines",
		CountryCodePn: "Pitcairn",
		CountryCodePl: "Poland",
		CountryCodePt: "Portugal",
		CountryCodePr: "Puerto Rico",
		CountryCodeQa: "Qatar",
		CountryCodeRe: "R\u00e9union",
		CountryCodeRo: "Romania",
		CountryCodeRu: "Russia",
		CountryCodeRw: "Rwanda",
		CountryCodeBl: "Saint Barth\u00e9lemy",
		CountryCodeSh: "Saint Helena, Ascension and Tristan da Cunha",
		CountryCodeKn: "Saint Kitts and Nevis",
		CountryCodeLc: "Saint Lucia",
		CountryCodeMf: "Saint Martin (French part)",
		CountryCodePm: "Saint Pierre and Miquelon",
		CountryCodeVc: "Saint Vincent and the Grenadines",
		CountryCodeWs: "Samoa",
		CountryCodeSm: "San Marino",
		CountryCodeSt: "Sao Tome and Principe",
		CountryCodeSa: "Saudi Arabia",
		CountryCodeSn: "Senegal",
		CountryCodeRs: "Serbia",
		CountryCodeSc: "Seychelles",
		CountryCodeSl: "Sierra Leone",
		CountryCodeSg: "Singapore",
		CountryCodeSx: "Sint Maarten (Dutch part)",
		CountryCodeSk: "Slovakia",
		CountryCodeSi: "Slovenia",
		CountryCodeSb: "Solomon Islands",
		CountryCodeSo: "Somalia",
		CountryCodeZa: "South Africa",
		CountryCodeGs: "South Georgia and the South Sandwich Islands",
		CountryCodeSs: "South Sudan",
		CountryCodeEs: "Spain",
		CountryCodeLk: "Sri Lanka",
		CountryCodeSd: "Sudan",
		CountryCodeSr: "Suriname",
		CountryCodeSj: "Svalbard and Jan Mayen",
		CountryCodeSe: "Sweden",
		CountryCodeCh: "Switzerland",
		CountryCodeSy: "Syria",
		CountryCodeTw: "Taiwan",
		CountryCodeTj: "Tajikistan",
		CountryCodeTz: "Tanzania",
		CountryCodeTh: "Thailand",
		CountryCodeTl: "Timor-Leste",
		CountryCodeTg: "Togo",
		CountryCodeTk: "Tokelau",
		CountryCodeTo: "Tonga",
		CountryCodeTt: "Trinidad and Tobago",
		CountryCodeTn: "Tunisia",
		CountryCodeTr: "Turkey",
		CountryCodeTm: "Turkmenistan",
		CountryCodeTc: "Turks and Caicos Islands",
		CountryCodeTv: "Tuvalu",
		CountryCodeUg: "Uganda",
		CountryCodeUa: "Ukraine",
		CountryCodeAe: "United Arab Emirates",
		CountryCodeGb: "United Kingdom",
		CountryCodeUm: "United States Minor Outlying Islands",
		CountryCodeUs: "United States of America",
		CountryCodeUy: "Uruguay",
		CountryCodeUz: "Uzbekistan",
		CountryCodeVu: "Vanuatu",
		CountryCodeVe: "Venezuela",
		CountryCodeVn: "Vietnam",
		CountryCodeVg: "Virgin Islands (British)",
		CountryCodeVi: "Virgin Islands (U.S.)",
		CountryCodeWf: "Wallis and Futuna",
		CountryCodeEh: "Western Sahara",
		CountryCodeYe: "Yemen",
		CountryCodeZm: "Zambia",
		CountryCodeZw: "Zimbabwe",
		CountryCodeEu: "European Union",
	}
	Languages = map[string]string{
		"af":             "Afrikaans",
		"af-na":          "Afrikaans (Namibia)",
		"af-za":          "Afrikaans (South Africa)",
		"agq":            "Aghem",
		"agq-cm":         "Aghem (Cameroon)",
		"ak":             "Akan",
		"ak-gh":          "Akan (Ghana)",
		"am":             "Amharic",
		"am-et":          "Amharic (Ethiopia)",
		"ar":             "Arabic",
		"ar-ae":          "Arabic (United Arab Emirates)",
		"ar-bh":          "Arabic (Bahrain)",
		"ar-dj":          "Arabic (Djibouti)",
		"ar-dz":          "Arabic (Algeria)",
		"ar-eg":          "Arabic (Egypt)",
		"ar-eh":          "Arabic (Western Sahara)",
		"ar-er":          "Arabic (Eritrea)",
		"ar-il":          "Arabic (Israel)",
		"ar-iq":          "Arabic (Iraq)",
		"ar-jo":          "Arabic (Jordan)",
		"ar-km":          "Arabic (Comoros)",
		"ar-kw":          "Arabic (Kuwait)",
		"ar-lb":          "Arabic (Lebanon)",
		"ar-ly":          "Arabic (Libya)",
		"ar-ma":          "Arabic (Morocco)",
		"ar-mr":          "Arabic (Mauritania)",
		"ar-om":          "Arabic (Oman)",
		"ar-ps":          "Arabic (Palestinian Territories)",
		"ar-qa":          "Arabic (Qatar)",
		"ar-sa":          "Arabic (Saudi Arabia)",
		"ar-sd":          "Arabic (Sudan)",
		"ar-so":          "Arabic (Somalia)",
		"ar-ss":          "Arabic (South Sudan)",
		"ar-sy":          "Arabic (Syria)",
		"ar-td":          "Arabic (Chad)",
		"ar-tn":          "Arabic (Tunisia)",
		"ar-ye":          "Arabic (Yemen)",
		"as":             "Assamese",
		"as-in":          "Assamese (India)",
		"asa":            "Asu",
		"asa-tz":         "Asu (Tanzania)",
		"ast":            "Asturian",
		"ast-es":         "Asturian (Spain)",
		"az":             "Azerbaijani",
		"az-cyrl":        "Azerbaijani (Cyrillic)",
		"az-cyrl-az":     "Azerbaijani (Cyrillic, Azerbaijan)",
		"az-latn":        "Azerbaijani (Latin)",
		"az-latn-az":     "Azerbaijani (Latin, Azerbaijan)",
		"bas":            "Basaa",
		"bas-cm":         "Basaa (Cameroon)",
		"be":             "Belarusian",
		"be-by":          "Belarusian (Belarus)",
		"bem":            "Bemba",
		"bem-zm":         "Bemba (Zambia)",
		"bez":            "Bena",
		"bez-tz":         "Bena (Tanzania)",
		"bg":             "Bulgarian",
		"bg-bg":          "Bulgarian (Bulgaria)",
		"bm":             "Bambara",
		"bm-ml":          "Bambara (Mali)",
		"bn":             "Bangla",
		"bn-bd":          "Bangla (Bangladesh)",
		"bn-in":          "Bangla (India)",
		"bo":             "Tibetan",
		"bo-cn":          "Tibetan (China)",
		"bo-in":          "Tibetan (India)",
		"br":             "Breton",
		"br-fr":          "Breton (France)",
		"brx":            "Bodo",
		"brx-in":         "Bodo (India)",
		"bs":             "Bosnian",
		"bs-cyrl":        "Bosnian (Cyrillic)",
		"bs-cyrl-ba":     "Bosnian (Cyrillic, Bosnia & Herzegovina)",
		"bs-latn":        "Bosnian (Latin)",
		"bs-latn-ba":     "Bosnian (Latin, Bosnia & Herzegovina)",
		"ca":             "Catalan",
		"ca-ad":          "Catalan (Andorra)",
		"ca-es":          "Catalan (Spain)",
		"ca-es-valencia": "Catalan (Spain, Valencian)",
		"ca-fr":          "Catalan (France)",
		"ca-it":          "Catalan (Italy)",
		"ccp":            "Chakma",
		"ccp-bd":         "Chakma (Bangladesh)",
		"ccp-in":         "Chakma (India)",
		"ce":             "Chechen",
		"ce-ru":          "Chechen (Russia)",
		"ceb":            "Cebuano",
		"ceb-ph":         "Cebuano (Philippines)",
		"cgg":            "Chiga",
		"cgg-ug":         "Chiga (Uganda)",
		"chr":            "Cherokee",
		"chr-us":         "Cherokee (United States)",
		"ckb":            "Central Kurdish",
		"ckb-iq":         "Central Kurdish (Iraq)",
		"ckb-ir":         "Central Kurdish (Iran)",
		"cs":             "Czech",
		"cs-cz":          "Czech (Czechia)",
		"cu":             "Church Slavic",
		"cu-ru":          "Church Slavic (Russia)",
		"cy":             "Welsh",
		"cy-gb":          "Welsh (United Kingdom)",
		"da":             "Danish",
		"da-dk":          "Danish (Denmark)",
		"da-gl":          "Danish (Greenland)",
		"dav":            "Taita",
		"dav-ke":         "Taita (Kenya)",
		"de":             "German",
		"de-at":          "German (Austria)",
		"de-be":          "German (Belgium)",
		"de-ch":          "German (Switzerland)",
		"de-de":          "German (Germany)",
		"de-it":          "German (Italy)",
		"de-li":          "German (Liechtenstein)",
		"de-lu":          "German (Luxembourg)",
		"dje":            "Zarma",
		"dje-ne":         "Zarma (Niger)",
		"dsb":            "Lower Sorbian",
		"dsb-de":         "Lower Sorbian (Germany)",
		"dua":            "Duala",
		"dua-cm":         "Duala (Cameroon)",
		"dyo":            "Jola-Fonyi",
		"dyo-sn":         "Jola-Fonyi (Senegal)",
		"dz":             "Dzongkha",
		"dz-bt":          "Dzongkha (Bhutan)",
		"ebu":            "Embu",
		"ebu-ke":         "Embu (Kenya)",
		"ee":             "Ewe",
		"ee-gh":          "Ewe (Ghana)",
		"ee-tg":          "Ewe (Togo)",
		"el":             "Greek",
		"el-cy":          "Greek (Cyprus)",
		"el-gr":          "Greek (Greece)",
		"en":             "English",
		"en-ae":          "English (United Arab Emirates)",
		"en-ag":          "English (Antigua & Barbuda)",
		"en-ai":          "English (Anguilla)",
		"en-as":          "English (American Samoa)",
		"en-at":          "English (Austria)",
		"en-au":          "English (Australia)",
		"en-bb":          "English (Barbados)",
		"en-be":          "English (Belgium)",
		"en-bi":          "English (Burundi)",
		"en-bm":          "English (Bermuda)",
		"en-bs":          "English (Bahamas)",
		"en-bw":          "English (Botswana)",
		"en-bz":          "English (Belize)",
		"en-ca":          "English (Canada)",
		"en-cc":          "English (Cocos (Keeling) Islands)",
		"en-ch":          "English (Switzerland)",
		"en-ck":          "English (Cook Islands)",
		"en-cm":          "English (Cameroon)",
		"en-cx":          "English (Christmas Island)",
		"en-cy":          "English (Cyprus)",
		"en-de":          "English (Germany)",
		"en-dg":          "English (Diego Garcia)",
		"en-dk":          "English (Denmark)",
		"en-dm":          "English (Dominica)",
		"en-er":          "English (Eritrea)",
		"en-fi":          "English (Finland)",
		"en-fj":          "English (Fiji)",
		"en-fk":          "English (Falkland Islands)",
		"en-fm":          "English (Micronesia)",
		"en-gb":          "English (United Kingdom)",
		"en-gd":          "English (Grenada)",
		"en-gg":          "English (Guernsey)",
		"en-gh":          "English (Ghana)",
		"en-gi":          "English (Gibraltar)",
		"en-gm":          "English (Gambia)",
		"en-gu":          "English (Guam)",
		"en-gy":          "English (Guyana)",
		"en-hk":          "English (Hong Kong SAR China)",
		"en-ie":          "English (Ireland)",
		"en-il":          "English (Israel)",
		"en-im":          "English (Isle of Man)",
		"en-in":          "English (India)",
		"en-io":          "English (British Indian Ocean Territory)",
		"en-je":          "English (Jersey)",
		"en-jm":          "English (Jamaica)",
		"en-ke":          "English (Kenya)",
		"en-ki":          "English (Kiribati)",
		"en-kn":          "English (St. Kitts & Nevis)",
		"en-ky":          "English (Cayman Islands)",
		"en-lc":          "English (St. Lucia)",
		"en-lr":          "English (Liberia)",
		"en-ls":          "English (Lesotho)",
		"en-mg":          "English (Madagascar)",
		"en-mh":          "English (Marshall Islands)",
		"en-mo":          "English (Macao SAR China)",
		"en-mp":          "English (Northern Mariana Islands)",
		"en-ms":          "English (Montserrat)",
		"en-mt":          "English (Malta)",
		"en-mu":          "English (Mauritius)",
		"en-mw":          "English (Malawi)",
		"en-my":          "English (Malaysia)",
		"en-na":          "English (Namibia)",
		"en-nf":          "English (Norfolk Island)",
		"en-ng":          "English (Nigeria)",
		"en-nl":          "English (Netherlands)",
		"en-nr":          "English (Nauru)",
		"en-nu":          "English (Niue)",
		"en-nz":          "English (New Zealand)",
		"en-pg":          "English (Papua New Guinea)",
		"en-ph":          "English (Philippines)",
		"en-pk":          "English (Pakistan)",
		"en-pn":          "English (Pitcairn Islands)",
		"en-pr":          "English (Puerto Rico)",
		"en-pw":          "English (Palau)",
		"en-rw":          "English (Rwanda)",
		"en-sb":          "English (Solomon Islands)",
		"en-sc":          "English (Seychelles)",
		"en-sd":          "English (Sudan)",
		"en-se":          "English (Sweden)",
		"en-sg":          "English (Singapore)",
		"en-sh":          "English (St. Helena)",
		"en-si":          "English (Slovenia)",
		"en-sl":          "English (Sierra Leone)",
		"en-ss":          "English (South Sudan)",
		"en-sx":          "English (Sint Maarten)",
		"en-sz":          "English (Eswatini)",
		"en-tc":          "English (Turks & Caicos Islands)",
		"en-tk":          "English (Tokelau)",
		"en-to":          "English (Tonga)",
		"en-tt":          "English (Trinidad & Tobago)",
		"en-tv":          "English (Tuvalu)",
		"en-tz":          "English (Tanzania)",
		"en-ug":          "English (Uganda)",
		"en-um":          "English (U.S. Outlying Islands)",
		"en-us":          "English (United States)",
		"en-vc":          "English (St. Vincent & Grenadines)",
		"en-vg":          "English (British Virgin Islands)",
		"en-vi":          "English (U.S. Virgin Islands)",
		"en-vu":          "English (Vanuatu)",
		"en-ws":          "English (Samoa)",
		"en-za":          "English (South Africa)",
		"en-zm":          "English (Zambia)",
		"en-zw":          "English (Zimbabwe)",
		"eo":             "Esperanto",
		"es":             "Spanish",
		"es-ar":          "Spanish (Argentina)",
		"es-bo":          "Spanish (Bolivia)",
		"es-br":          "Spanish (Brazil)",
		"es-bz":          "Spanish (Belize)",
		"es-cl":          "Spanish (Chile)",
		"es-co":          "Spanish (Colombia)",
		"es-cr":          "Spanish (Costa Rica)",
		"es-cu":          "Spanish (Cuba)",
		"es-do":          "Spanish (Dominican Republic)",
		"es-ea":          "Spanish (Ceuta & Melilla)",
		"es-ec":          "Spanish (Ecuador)",
		"es-es":          "Spanish (Spain)",
		"es-gq":          "Spanish (Equatorial Guinea)",
		"es-gt":          "Spanish (Guatemala)",
		"es-hn":          "Spanish (Honduras)",
		"es-ic":          "Spanish (Canary Islands)",
		"es-mx":          "Spanish (Mexico)",
		"es-ni":          "Spanish (Nicaragua)",
		"es-pa":          "Spanish (Panama)",
		"es-pe":          "Spanish (Peru)",
		"es-ph":          "Spanish (Philippines)",
		"es-pr":          "Spanish (Puerto Rico)",
		"es-py":          "Spanish (Paraguay)",
		"es-sv":          "Spanish (El Salvador)",
		"es-us":          "Spanish (United States)",
		"es-uy":          "Spanish (Uruguay)",
		"es-ve":          "Spanish (Venezuela)",
		"et":             "Estonian",
		"et-ee":          "Estonian (Estonia)",
		"eu":             "Basque",
		"eu-es":          "Basque (Spain)",
		"ewo":            "Ewondo",
		"ewo-cm":         "Ewondo (Cameroon)",
		"fa":             "Persian",
		"fa-af":          "Persian (Afghanistan)",
		"fa-ir":          "Persian (Iran)",
		"ff":             "Fulah",
		"ff-adlm":        "Fulah (Adlam)",
		"ff-adlm-bf":     "Fulah (Adlam, Burkina Faso)",
		"ff-adlm-cm":     "Fulah (Adlam, Cameroon)",
		"ff-adlm-gh":     "Fulah (Adlam, Ghana)",
		"ff-adlm-gm":     "Fulah (Adlam, Gambia)",
		"ff-adlm-gn":     "Fulah (Adlam, Guinea)",
		"ff-adlm-gw":     "Fulah (Adlam, Guinea-Bissau)",
		"ff-adlm-lr":     "Fulah (Adlam, Liberia)",
		"ff-adlm-mr":     "Fulah (Adlam, Mauritania)",
		"ff-adlm-ne":     "Fulah (Adlam, Niger)",
		"ff-adlm-ng":     "Fulah (Adlam, Nigeria)",
		"ff-adlm-sl":     "Fulah (Adlam, Sierra Leone)",
		"ff-adlm-sn":     "Fulah (Adlam, Senegal)",
		"ff-latn":        "Fulah (Latin)",
		"ff-latn-bf":     "Fulah (Latin, Burkina Faso)",
		"ff-latn-cm":     "Fulah (Latin, Cameroon)",
		"ff-latn-gh":     "Fulah (Latin, Ghana)",
		"ff-latn-gm":     "Fulah (Latin, Gambia)",
		"ff-latn-gn":     "Fulah (Latin, Guinea)",
		"ff-latn-gw":     "Fulah (Latin, Guinea-Bissau)",
		"ff-latn-lr":     "Fulah (Latin, Liberia)",
		"ff-latn-mr":     "Fulah (Latin, Mauritania)",
		"ff-latn-ne":     "Fulah (Latin, Niger)",
		"ff-latn-ng":     "Fulah (Latin, Nigeria)",
		"ff-latn-sl":     "Fulah (Latin, Sierra Leone)",
		"ff-latn-sn":     "Fulah (Latin, Senegal)",
		"fi":             "Finnish",
		"fi-fi":          "Finnish (Finland)",
		"fil":            "Filipino",
		"fil-ph":         "Filipino (Philippines)",
		"fo":             "Faroese",
		"fo-dk":          "Faroese (Denmark)",
		"fo-fo":          "Faroese (Faroe Islands)",
		"fr":             "French",
		"fr-be":          "French (Belgium)",
		"fr-bf":          "French (Burkina Faso)",
		"fr-bi":          "French (Burundi)",
		"fr-bj":          "French (Benin)",
		"fr-bl":          "French (St. Barth\u00e9lemy)",
		"fr-ca":          "French (Canada)",
		"fr-cd":          "French (Congo - Kinshasa)",
		"fr-cf":          "French (Central African Republic)",
		"fr-cg":          "French (Congo - Brazzaville)",
		"fr-ch":          "French (Switzerland)",
		"fr-ci":          "French (C\u00f4te d\u2019Ivoire)",
		"fr-cm":          "French (Cameroon)",
		"fr-dj":          "French (Djibouti)",
		"fr-dz":          "French (Algeria)",
		"fr-fr":          "French (France)",
		"fr-ga":          "French (Gabon)",
		"fr-gf":          "French (French Guiana)",
		"fr-gn":          "French (Guinea)",
		"fr-gp":          "French (Guadeloupe)",
		"fr-gq":          "French (Equatorial Guinea)",
		"fr-ht":          "French (Haiti)",
		"fr-km":          "French (Comoros)",
		"fr-lu":          "French (Luxembourg)",
		"fr-ma":          "French (Morocco)",
		"fr-mc":          "French (Monaco)",
		"fr-mf":          "French (St. Martin)",
		"fr-mg":          "French (Madagascar)",
		"fr-ml":          "French (Mali)",
		"fr-mq":          "French (Martinique)",
		"fr-mr":          "French (Mauritania)",
		"fr-mu":          "French (Mauritius)",
		"fr-nc":          "French (New Caledonia)",
		"fr-ne":          "French (Niger)",
		"fr-pf":          "French (French Polynesia)",
		"fr-pm":          "French (St. Pierre & Miquelon)",
		"fr-re":          "French (R\u00e9union)",
		"fr-rw":          "French (Rwanda)",
		"fr-sc":          "French (Seychelles)",
		"fr-sn":          "French (Senegal)",
		"fr-sy":          "French (Syria)",
		"fr-td":          "French (Chad)",
		"fr-tg":          "French (Togo)",
		"fr-tn":          "French (Tunisia)",
		"fr-vu":          "French (Vanuatu)",
		"fr-wf":          "French (Wallis & Futuna)",
		"fr-yt":          "French (Mayotte)",
		"fur":            "Friulian",
		"fur-it":         "Friulian (Italy)",
		"fy":             "Western Frisian",
		"fy-nl":          "Western Frisian (Netherlands)",
		"ga":             "Irish",
		"ga-gb":          "Irish (United Kingdom)",
		"ga-ie":          "Irish (Ireland)",
		"gd":             "Scottish Gaelic",
		"gd-gb":          "Scottish Gaelic (United Kingdom)",
		"gl":             "Galician",
		"gl-es":          "Galician (Spain)",
		"gsw":            "Swiss German",
		"gsw-ch":         "Swiss German (Switzerland)",
		"gsw-fr":         "Swiss German (France)",
		"gsw-li":         "Swiss German (Liechtenstein)",
		"gu":             "Gujarati",
		"gu-in":          "Gujarati (India)",
		"guz":            "Gusii",
		"guz-ke":         "Gusii (Kenya)",
		"gv":             "Manx",
		"gv-im":          "Manx (Isle of Man)",
		"ha":             "Hausa",
		"ha-gh":          "Hausa (Ghana)",
		"ha-ne":          "Hausa (Niger)",
		"ha-ng":          "Hausa (Nigeria)",
		"haw":            "Hawaiian",
		"haw-us":         "Hawaiian (United States)",
		"he":             "Hebrew",
		"he-il":          "Hebrew (Israel)",
		"hi":             "Hindi",
		"hi-in":          "Hindi (India)",
		"hr":             "Croatian",
		"hr-ba":          "Croatian (Bosnia & Herzegovina)",
		"hr-hr":          "Croatian (Croatia)",
		"hsb":            "Upper Sorbian",
		"hsb-de":         "Upper Sorbian (Germany)",
		"hu":             "Hungarian",
		"hu-hu":          "Hungarian (Hungary)",
		"hy":             "Armenian",
		"hy-am":          "Armenian (Armenia)",
		"ia":             "Interlingua",
		"id":             "Indonesian",
		"id-id":          "Indonesian (Indonesia)",
		"ig":             "Igbo",
		"ig-ng":          "Igbo (Nigeria)",
		"ii":             "Sichuan Yi",
		"ii-cn":          "Sichuan Yi (China)",
		"is":             "Icelandic",
		"is-is":          "Icelandic (Iceland)",
		"it":             "Italian",
		"it-ch":          "Italian (Switzerland)",
		"it-it":          "Italian (Italy)",
		"it-sm":          "Italian (San Marino)",
		"it-va":          "Italian (Vatican City)",
		"ja":             "Japanese",
		"ja-jp":          "Japanese (Japan)",
		"jgo":            "Ngomba",
		"jgo-cm":         "Ngomba (Cameroon)",
		"jmc":            "Machame",
		"jmc-tz":         "Machame (Tanzania)",
		"jv":             "Javanese",
		"jv-id":          "Javanese (Indonesia)",
		"ka":             "Georgian",
		"ka-ge":          "Georgian (Georgia)",
		"kab":            "Kabyle",
		"kab-dz":         "Kabyle (Algeria)",
		"kam":            "Kamba",
		"kam-ke":         "Kamba (Kenya)",
		"kde":            "Makonde",
		"kde-tz":         "Makonde (Tanzania)",
		"kea":            "Kabuverdianu",
		"kea-cv":         "Kabuverdianu (Cape Verde)",
		"khq":            "Koyra Chiini",
		"khq-ml":         "Koyra Chiini (Mali)",
		"ki":             "Kikuyu",
		"ki-ke":          "Kikuyu (Kenya)",
		"kk":             "Kazakh",
		"kk-kz":          "Kazakh (Kazakhstan)",
		"kkj":            "Kako",
		"kkj-cm":         "Kako (Cameroon)",
		"kl":             "Kalaallisut",
		"kl-gl":          "Kalaallisut (Greenland)",
		"kln":            "Kalenjin",
		"kln-ke":         "Kalenjin (Kenya)",
		"km":             "Khmer",
		"km-kh":          "Khmer (Cambodia)",
		"kn":             "Kannada",
		"kn-in":          "Kannada (India)",
		"ko":             "Korean",
		"ko-kp":          "Korean (North Korea)",
		"ko-kr":          "Korean (South Korea)",
		"kok":            "Konkani",
		"kok-in":         "Konkani (India)",
		"ks":             "Kashmiri",
		"ks-arab":        "Kashmiri (Arabic)",
		"ks-arab-in":     "Kashmiri (Arabic, India)",
		"ksb":            "Shambala",
		"ksb-tz":         "Shambala (Tanzania)",
		"ksf":            "Bafia",
		"ksf-cm":         "Bafia (Cameroon)",
		"ksh":            "Colognian",
		"ksh-de":         "Colognian (Germany)",
		"ku":             "Kurdish",
		"ku-tr":          "Kurdish (Turkey)",
		"kw":             "Cornish",
		"kw-gb":          "Cornish (United Kingdom)",
		"ky":             "Kyrgyz",
		"ky-kg":          "Kyrgyz (Kyrgyzstan)",
		"lag":            "Langi",
		"lag-tz":         "Langi (Tanzania)",
		"lb":             "Luxembourgish",
		"lb-lu":          "Luxembourgish (Luxembourg)",
		"lg":             "Ganda",
		"lg-ug":          "Ganda (Uganda)",
		"lkt":            "Lakota",
		"lkt-us":         "Lakota (United States)",
		"ln":             "Lingala",
		"ln-ao":          "Lingala (Angola)",
		"ln-cd":          "Lingala (Congo - Kinshasa)",
		"ln-cf":          "Lingala (Central African Republic)",
		"ln-cg":          "Lingala (Congo - Brazzaville)",
		"lo":             "Lao",
		"lo-la":          "Lao (Laos)",
		"lrc":            "Northern Luri",
		"lrc-iq":         "Northern Luri (Iraq)",
		"lrc-ir":         "Northern Luri (Iran)",
		"lt":             "Lithuanian",
		"lt-lt":          "Lithuanian (Lithuania)",
		"lu":             "Luba-Katanga",
		"lu-cd":          "Luba-Katanga (Congo - Kinshasa)",
		"luo":            "Luo",
		"luo-ke":         "Luo (Kenya)",
		"luy":            "Luyia",
		"luy-ke":         "Luyia (Kenya)",
		"lv":             "Latvian",
		"lv-lv":          "Latvian (Latvia)",
		"mai":            "Maithili",
		"mai-in":         "Maithili (India)",
		"mas":            "Masai",
		"mas-ke":         "Masai (Kenya)",
		"mas-tz":         "Masai (Tanzania)",
		"mer":            "Meru",
		"mer-ke":         "Meru (Kenya)",
		"mfe":            "Morisyen",
		"mfe-mu":         "Morisyen (Mauritius)",
		"mg":             "Malagasy",
		"mg-mg":          "Malagasy (Madagascar)",
		"mgh":            "Makhuwa-Meetto",
		"mgh-mz":         "Makhuwa-Meetto (Mozambique)",
		"mgo":            "Meta\u02bc",
		"mgo-cm":         "Meta\u02bc (Cameroon)",
		"mi":             "Maori",
		"mi-nz":          "Maori (New Zealand)",
		"mk":             "Macedonian",
		"mk-mk":          "Macedonian (North Macedonia)",
		"ml":             "Malayalam",
		"ml-in":          "Malayalam (India)",
		"mn":             "Mongolian",
		"mn-mn":          "Mongolian (Mongolia)",
		"mni":            "Manipuri",
		"mni-beng":       "Manipuri (Bangla)",
		"mni-beng-in":    "Manipuri (Bangla, India)",
		"mr":             "Marathi",
		"mr-in":          "Marathi (India)",
		"ms":             "Malay",
		"ms-bn":          "Malay (Brunei)",
		"ms-id":          "Malay (Indonesia)",
		"ms-my":          "Malay (Malaysia)",
		"ms-sg":          "Malay (Singapore)",
		"mt":             "Maltese",
		"mt-mt":          "Maltese (Malta)",
		"mua":            "Mundang",
		"mua-cm":         "Mundang (Cameroon)",
		"my":             "Burmese",
		"my-mm":          "Burmese (Myanmar (Burma))",
		"mzn":            "Mazanderani",
		"mzn-ir":         "Mazanderani (Iran)",
		"naq":            "Nama",
		"naq-na":         "Nama (Namibia)",
		"nb":             "Norwegian Bokm\u00e5l",
		"nb-no":          "Norwegian Bokm\u00e5l (Norway)",
		"nb-sj":          "Norwegian Bokm\u00e5l (Svalbard & Jan Mayen)",
		"nd":             "North Ndebele",
		"nd-zw":          "North Ndebele (Zimbabwe)",
		"nds":            "Low German",
		"nds-de":         "Low German (Germany)",
		"nds-nl":         "Low German (Netherlands)",
		"ne":             "Nepali",
		"ne-in":          "Nepali (India)",
		"ne-np":          "Nepali (Nepal)",
		"nl":             "Dutch",
		"nl-aw":          "Dutch (Aruba)",
		"nl-be":          "Dutch (Belgium)",
		"nl-bq":          "Dutch (Caribbean Netherlands)",
		"nl-cw":          "Dutch (Cura\u00e7ao)",
		"nl-nl":          "Dutch (Netherlands)",
		"nl-sr":          "Dutch (Suriname)",
		"nl-sx":          "Dutch (Sint Maarten)",
		"nmg":            "Kwasio",
		"nmg-cm":         "Kwasio (Cameroon)",
		"nn":             "Norwegian Nynorsk",
		"nn-no":          "Norwegian Nynorsk (Norway)",
		"nnh":            "Ngiemboon",
		"nnh-cm":         "Ngiemboon (Cameroon)",
		"nus":            "Nuer",
		"nus-ss":         "Nuer (South Sudan)",
		"nyn":            "Nyankole",
		"nyn-ug":         "Nyankole (Uganda)",
		"om":             "Oromo",
		"om-et":          "Oromo (Ethiopia)",
		"om-ke":          "Oromo (Kenya)",
		"or":             "Odia",
		"or-in":          "Odia (India)",
		"os":             "Ossetic",
		"os-ge":          "Ossetic (Georgia)",
		"os-ru":          "Ossetic (Russia)",
		"pa":             "Punjabi",
		"pa-arab":        "Punjabi (Arabic)",
		"pa-arab-pk":     "Punjabi (Arabic, Pakistan)",
		"pa-guru":        "Punjabi (Gurmukhi)",
		"pa-guru-in":     "Punjabi (Gurmukhi, India)",
		"pcm":            "Nigerian Pidgin",
		"pcm-ng":         "Nigerian Pidgin (Nigeria)",
		"pl":             "Polish",
		"pl-pl":          "Polish (Poland)",
		"prg":            "Prussian",
		"ps":             "Pashto",
		"ps-af":          "Pashto (Afghanistan)",
		"ps-pk":          "Pashto (Pakistan)",
		"pt":             "Portuguese",
		"pt-ao":          "Portuguese (Angola)",
		"pt-br":          "Portuguese (Brazil)",
		"pt-ch":          "Portuguese (Switzerland)",
		"pt-cv":          "Portuguese (Cape Verde)",
		"pt-gq":          "Portuguese (Equatorial Guinea)",
		"pt-gw":          "Portuguese (Guinea-Bissau)",
		"pt-lu":          "Portuguese (Luxembourg)",
		"pt-mo":          "Portuguese (Macao SAR China)",
		"pt-mz":          "Portuguese (Mozambique)",
		"pt-pt":          "Portuguese (Portugal)",
		"pt-st":          "Portuguese (S\u00e3o Tom\u00e9 & Pr\u00edncipe)",
		"pt-tl":          "Portuguese (Timor-Leste)",
		"qu":             "Quechua",
		"qu-bo":          "Quechua (Bolivia)",
		"qu-ec":          "Quechua (Ecuador)",
		"qu-pe":          "Quechua (Peru)",
		"rm":             "Romansh",
		"rm-ch":          "Romansh (Switzerland)",
		"rn":             "Rundi",
		"rn-bi":          "Rundi (Burundi)",
		"ro":             "Romanian",
		"ro-md":          "Romanian (Moldova)",
		"ro-ro":          "Romanian (Romania)",
		"rof":            "Rombo",
		"rof-tz":         "Rombo (Tanzania)",
		"ru":             "Russian",
		"ru-by":          "Russian (Belarus)",
		"ru-kg":          "Russian (Kyrgyzstan)",
		"ru-kz":          "Russian (Kazakhstan)",
		"ru-md":          "Russian (Moldova)",
		"ru-ru":          "Russian (Russia)",
		"ru-ua":          "Russian (Ukraine)",
		"rw":             "Kinyarwanda",
		"rw-rw":          "Kinyarwanda (Rwanda)",
		"rwk":            "Rwa",
		"rwk-tz":         "Rwa (Tanzania)",
		"sah":            "Sakha",
		"sah-ru":         "Sakha (Russia)",
		"saq":            "Samburu",
		"saq-ke":         "Samburu (Kenya)",
		"sat":            "Santali",
		"sat-olck":       "Santali (Ol Chiki)",
		"sat-olck-in":    "Santali (Ol Chiki, India)",
		"sbp":            "Sangu",
		"sbp-tz":         "Sangu (Tanzania)",
		"sd":             "Sindhi",
		"sd-arab":        "Sindhi (Arabic)",
		"sd-arab-pk":     "Sindhi (Arabic, Pakistan)",
		"sd-deva":        "Sindhi (Devanagari)",
		"sd-deva-in":     "Sindhi (Devanagari, India)",
		"se":             "Northern Sami",
		"se-fi":          "Northern Sami (Finland)",
		"se-no":          "Northern Sami (Norway)",
		"se-se":          "Northern Sami (Sweden)",
		"seh":            "Sena",
		"seh-mz":         "Sena (Mozambique)",
		"ses":            "Koyraboro Senni",
		"ses-ml":         "Koyraboro Senni (Mali)",
		"sg":             "Sango",
		"sg-cf":          "Sango (Central African Republic)",
		"shi":            "Tachelhit",
		"shi-latn":       "Tachelhit (Latin)",
		"shi-latn-ma":    "Tachelhit (Latin, Morocco)",
		"shi-tfng":       "Tachelhit (Tifinagh)",
		"shi-tfng-ma":    "Tachelhit (Tifinagh, Morocco)",
		"si":             "Sinhala",
		"si-lk":          "Sinhala (Sri Lanka)",
		"sk":             "Slovak",
		"sk-sk":          "Slovak (Slovakia)",
		"sl":             "Slovenian",
		"sl-si":          "Slovenian (Slovenia)",
		"smn":            "Inari Sami",
		"smn-fi":         "Inari Sami (Finland)",
		"sn":             "Shona",
		"sn-zw":          "Shona (Zimbabwe)",
		"so":             "Somali",
		"so-dj":          "Somali (Djibouti)",
		"so-et":          "Somali (Ethiopia)",
		"so-ke":          "Somali (Kenya)",
		"so-so":          "Somali (Somalia)",
		"sq":             "Albanian",
		"sq-al":          "Albanian (Albania)",
		"sq-mk":          "Albanian (North Macedonia)",
		"sq-xk":          "Albanian (Kosovo)",
		"sr":             "Serbian",
		"sr-cyrl":        "Serbian (Cyrillic)",
		"sr-cyrl-ba":     "Serbian (Cyrillic, Bosnia & Herzegovina)",
		"sr-cyrl-me":     "Serbian (Cyrillic, Montenegro)",
		"sr-cyrl-rs":     "Serbian (Cyrillic, Serbia)",
		"sr-cyrl-xk":     "Serbian (Cyrillic, Kosovo)",
		"sr-latn":        "Serbian (Latin)",
		"sr-latn-ba":     "Serbian (Latin, Bosnia & Herzegovina)",
		"sr-latn-me":     "Serbian (Latin, Montenegro)",
		"sr-latn-rs":     "Serbian (Latin, Serbia)",
		"sr-latn-xk":     "Serbian (Latin, Kosovo)",
		"su":             "Sundanese",
		"su-latn":        "Sundanese (Latin)",
		"su-latn-id":     "Sundanese (Latin, Indonesia)",
		"sv":             "Swedish",
		"sv-ax":          "Swedish (\u00c5land Islands)",
		"sv-fi":          "Swedish (Finland)",
		"sv-se":          "Swedish (Sweden)",
		"sw":             "Swahili",
		"sw-cd":          "Swahili (Congo - Kinshasa)",
		"sw-ke":          "Swahili (Kenya)",
		"sw-tz":          "Swahili (Tanzania)",
		"sw-ug":          "Swahili (Uganda)",
		"ta":             "Tamil",
		"ta-in":          "Tamil (India)",
		"ta-lk":          "Tamil (Sri Lanka)",
		"ta-my":          "Tamil (Malaysia)",
		"ta-sg":          "Tamil (Singapore)",
		"te":             "Telugu",
		"te-in":          "Telugu (India)",
		"teo":            "Teso",
		"teo-ke":         "Teso (Kenya)",
		"teo-ug":         "Teso (Uganda)",
		"tg":             "Tajik",
		"tg-tj":          "Tajik (Tajikistan)",
		"th":             "Thai",
		"th-th":          "Thai (Thailand)",
		"ti":             "Tigrinya",
		"ti-er":          "Tigrinya (Eritrea)",
		"ti-et":          "Tigrinya (Ethiopia)",
		"tk":             "Turkmen",
		"tk-tm":          "Turkmen (Turkmenistan)",
		"to":             "Tongan",
		"to-to":          "Tongan (Tonga)",
		"tr":             "Turkish",
		"tr-cy":          "Turkish (Cyprus)",
		"tr-tr":          "Turkish (Turkey)",
		"tt":             "Tatar",
		"tt-ru":          "Tatar (Russia)",
		"twq":            "Tasawaq",
		"twq-ne":         "Tasawaq (Niger)",
		"tzm":            "Central Atlas Tamazight",
		"tzm-ma":         "Central Atlas Tamazight (Morocco)",
		"ug":             "Uyghur",
		"ug-cn":          "Uyghur (China)",
		"uk":             "Ukrainian",
		"uk-ua":          "Ukrainian (Ukraine)",
		"ur":             "Urdu",
		"ur-in":          "Urdu (India)",
		"ur-pk":          "Urdu (Pakistan)",
		"uz":             "Uzbek",
		"uz-arab":        "Uzbek (Arabic)",
		"uz-arab-af":     "Uzbek (Arabic, Afghanistan)",
		"uz-cyrl":        "Uzbek (Cyrillic)",
		"uz-cyrl-uz":     "Uzbek (Cyrillic, Uzbekistan)",
		"uz-latn":        "Uzbek (Latin)",
		"uz-latn-uz":     "Uzbek (Latin, Uzbekistan)",
		"vai":            "Vai",
		"vai-latn":       "Vai (Latin)",
		"vai-latn-lr":    "Vai (Latin, Liberia)",
		"vai-vaii":       "Vai (Vai)",
		"vai-vaii-lr":    "Vai (Vai, Liberia)",
		"vi":             "Vietnamese",
		"vi-vn":          "Vietnamese (Vietnam)",
		"vo":             "Volap\u00fck",
		"vun":            "Vunjo",
		"vun-tz":         "Vunjo (Tanzania)",
		"wae":            "Walser",
		"wae-ch":         "Walser (Switzerland)",
		"wo":             "Wolof",
		"wo-sn":          "Wolof (Senegal)",
		"xh":             "Xhosa",
		"xh-za":          "Xhosa (South Africa)",
		"xog":            "Soga",
		"xog-ug":         "Soga (Uganda)",
		"yav":            "Yangben",
		"yav-cm":         "Yangben (Cameroon)",
		"yi":             "Yiddish",
		"yo":             "Yoruba",
		"yo-bj":          "Yoruba (Benin)",
		"yo-ng":          "Yoruba (Nigeria)",
		"yue":            "Cantonese",
		"yue-hans":       "Cantonese (Simplified)",
		"yue-hans-cn":    "Cantonese (Simplified, China)",
		"yue-hant":       "Cantonese (Traditional)",
		"yue-hant-hk":    "Cantonese (Traditional, Hong Kong SAR China)",
		"zgh":            "Standard Moroccan Tamazight",
		"zgh-ma":         "Standard Moroccan Tamazight (Morocco)",
		"zh":             "Chinese",
		"zh-hans":        "Chinese (Simplified)",
		"zh-hans-cn":     "Chinese (Simplified, China)",
		"zh-hans-hk":     "Chinese (Simplified, Hong Kong SAR China)",
		"zh-hans-mo":     "Chinese (Simplified, Macao SAR China)",
		"zh-hans-sg":     "Chinese (Simplified, Singapore)",
		"zh-hant":        "Chinese (Traditional)",
		"zh-hant-hk":     "Chinese (Traditional, Hong Kong SAR China)",
		"zh-hant-mo":     "Chinese (Traditional, Macao SAR China)",
		"zh-hant-tw":     "Chinese (Traditional, Taiwan)",
		"zu":             "Zulu",
		"zu-za":          "Zulu (South Africa)",
	}
	ReservedName = []string{
		"admin",
		"api",
		"channel",
		"claim",
		"error",
		"files",
		"help",
		"landing",
		"login",
		"mfa",
		"oauth",
		"plug",
		"plugins",
		"post",
		"signup",
		"sitename",
	}
	ValidUsernameChars = regexp.MustCompile(`^[a-z0-9\.\-_]+$`)
	RestrictedUsernames = map[string]bool{
		"all":      true,
		"channel":  true,
		"sitename": true,
		"system":   true,
		"admin":    true,
	}

	for code := range Countries {
		MULTIPLE_COUNTRIES_MAX_LENGTH += len(code)
	}
	MULTIPLE_COUNTRIES_MAX_LENGTH += len(Countries) - 1
}

// NamePart is an Enum
type NamePart string

// two name parts
const (
	FirstName NamePart = "first"
	LastName  NamePart = "last"
)

// TaxType is for unifying tax type object that comes from tax gateway
type TaxType struct {
	Code         string
	Descriptiton string
}
