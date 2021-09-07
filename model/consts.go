package model

import "regexp"

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
	DEFAULT_COUNTRY                = "US"
)

var (
	Countries                     map[string]string // countries supported by app
	Languages                     map[string]string // Languages supported by app
	MULTIPLE_COUNTRIES_MAX_LENGTH int               // some model"s country fields contains multiple countries
	ReservedName                  []string          // usernames that can only be used by system
	ValidUsernameChars            *regexp.Regexp    // regexp for username validation
	RestrictedUsernames           map[string]bool   // usernames that cannot be used
)

func init() {
	// borrowed from django_countries
	Countries = map[string]string{
		"AF": "Afghanistan",
		"AX": "Åland Islands",
		"AL": "Albania",
		"DZ": "Algeria",
		"AS": "American Samoa",
		"AD": "Andorra",
		"AO": "Angola",
		"AI": "Anguilla",
		"AQ": "Antarctica",
		"AG": "Antigua and Barbuda",
		"AR": "Argentina",
		"AM": "Armenia",
		"AW": "Aruba",
		"AU": "Australia",
		"AT": "Austria",
		"AZ": "Azerbaijan",
		"BS": "Bahamas",
		"BH": "Bahrain",
		"BD": "Bangladesh",
		"BB": "Barbados",
		"BY": "Belarus",
		"BE": "Belgium",
		"BZ": "Belize",
		"BJ": "Benin",
		"BM": "Bermuda",
		"BT": "Bhutan",
		"BO": "Bolivia",
		"BQ": "Bonaire, Sint Eustatius and Saba",
		"BA": "Bosnia and Herzegovina",
		"BW": "Botswana",
		"BV": "Bouvet Island",
		"BR": "Brazil",
		"IO": "British Indian Ocean Territory",
		"BN": "Brunei",
		"BG": "Bulgaria",
		"BF": "Burkina Faso",
		"BI": "Burundi",
		"CV": "Cabo Verde",
		"KH": "Cambodia",
		"CM": "Cameroon",
		"CA": "Canada",
		"KY": "Cayman Islands",
		"CF": "Central African Republic",
		"TD": "Chad",
		"CL": "Chile",
		"CN": "China",
		"CX": "Christmas Island",
		"CC": "Cocos (Keeling) Islands",
		"CO": "Colombia",
		"KM": "Comoros",
		"CG": "Congo",
		"CD": "Congo (the Democratic Republic of the)",
		"CK": "Cook Islands",
		"CR": "Costa Rica",
		"CI": "Côte d'Ivoire",
		"HR": "Croatia",
		"CU": "Cuba",
		"CW": "Curaçao",
		"CY": "Cyprus",
		"CZ": "Czechia",
		"DK": "Denmark",
		"DJ": "Djibouti",
		"DM": "Dominica",
		"DO": "Dominican Republic",
		"EC": "Ecuador",
		"EG": "Egypt",
		"SV": "El Salvador",
		"GQ": "Equatorial Guinea",
		"ER": "Eritrea",
		"EE": "Estonia",
		"SZ": "Eswatini",
		"ET": "Ethiopia",
		"FK": "Falkland Islands (Malvinas)",
		"FO": "Faroe Islands",
		"FJ": "Fiji",
		"FI": "Finland",
		"FR": "France",
		"GF": "French Guiana",
		"PF": "French Polynesia",
		"TF": "French Southern Territories",
		"GA": "Gabon",
		"GM": "Gambia",
		"GE": "Georgia",
		"DE": "Germany",
		"GH": "Ghana",
		"GI": "Gibraltar",
		"GR": "Greece",
		"GL": "Greenland",
		"GD": "Grenada",
		"GP": "Guadeloupe",
		"GU": "Guam",
		"GT": "Guatemala",
		"GG": "Guernsey",
		"GN": "Guinea",
		"GW": "Guinea-Bissau",
		"GY": "Guyana",
		"HT": "Haiti",
		"HM": "Heard Island and McDonald Islands",
		"VA": "Holy See",
		"HN": "Honduras",
		"HK": "Hong Kong",
		"HU": "Hungary",
		"IS": "Iceland",
		"IN": "India",
		"ID": "Indonesia",
		"IR": "Iran",
		"IQ": "Iraq",
		"IE": "Ireland",
		"IM": "Isle of Man",
		"IL": "Israel",
		"IT": "Italy",
		"JM": "Jamaica",
		"JP": "Japan",
		"JE": "Jersey",
		"JO": "Jordan",
		"KZ": "Kazakhstan",
		"KE": "Kenya",
		"KI": "Kiribati",
		"KP": "North Korea",
		"KR": "South Korea",
		"KW": "Kuwait",
		"KG": "Kyrgyzstan",
		"LA": "Laos",
		"LV": "Latvia",
		"LB": "Lebanon",
		"LS": "Lesotho",
		"LR": "Liberia",
		"LY": "Libya",
		"LI": "Liechtenstein",
		"LT": "Lithuania",
		"LU": "Luxembourg",
		"MO": "Macao",
		"MG": "Madagascar",
		"MW": "Malawi",
		"MY": "Malaysia",
		"MV": "Maldives",
		"ML": "Mali",
		"MT": "Malta",
		"MH": "Marshall Islands",
		"MQ": "Martinique",
		"MR": "Mauritania",
		"MU": "Mauritius",
		"YT": "Mayotte",
		"MX": "Mexico",
		"FM": "Micronesia (Federated States of)",
		"MD": "Moldova",
		"MC": "Monaco",
		"MN": "Mongolia",
		"ME": "Montenegro",
		"MS": "Montserrat",
		"MA": "Morocco",
		"MZ": "Mozambique",
		"MM": "Myanmar",
		"NA": "Namibia",
		"NR": "Nauru",
		"NP": "Nepal",
		"NL": "Netherlands",
		"NC": "New Caledonia",
		"NZ": "New Zealand",
		"NI": "Nicaragua",
		"NE": "Niger",
		"NG": "Nigeria",
		"NU": "Niue",
		"NF": "Norfolk Island",
		"MK": "North Macedonia",
		"MP": "Northern Mariana Islands",
		"NO": "Norway",
		"OM": "Oman",
		"PK": "Pakistan",
		"PW": "Palau",
		"PS": "Palestine, State of",
		"PA": "Panama",
		"PG": "Papua New Guinea",
		"PY": "Paraguay",
		"PE": "Peru",
		"PH": "Philippines",
		"PN": "Pitcairn",
		"PL": "Poland",
		"PT": "Portugal",
		"PR": "Puerto Rico",
		"QA": "Qatar",
		"RE": "Réunion",
		"RO": "Romania",
		"RU": "Russia",
		"RW": "Rwanda",
		"BL": "Saint Barthélemy",
		"SH": "Saint Helena, Ascension and Tristan da Cunha",
		"KN": "Saint Kitts and Nevis",
		"LC": "Saint Lucia",
		"MF": "Saint Martin (French part)",
		"PM": "Saint Pierre and Miquelon",
		"VC": "Saint Vincent and the Grenadines",
		"WS": "Samoa",
		"SM": "San Marino",
		"ST": "Sao Tome and Principe",
		"SA": "Saudi Arabia",
		"SN": "Senegal",
		"RS": "Serbia",
		"SC": "Seychelles",
		"SL": "Sierra Leone",
		"SG": "Singapore",
		"SX": "Sint Maarten (Dutch part)",
		"SK": "Slovakia",
		"SI": "Slovenia",
		"SB": "Solomon Islands",
		"SO": "Somalia",
		"ZA": "South Africa",
		"GS": "South Georgia and the South Sandwich Islands",
		"SS": "South Sudan",
		"ES": "Spain",
		"LK": "Sri Lanka",
		"SD": "Sudan",
		"SR": "Suriname",
		"SJ": "Svalbard and Jan Mayen",
		"SE": "Sweden",
		"CH": "Switzerland",
		"SY": "Syria",
		"TW": "Taiwan",
		"TJ": "Tajikistan",
		"TZ": "Tanzania",
		"TH": "Thailand",
		"TL": "Timor-Leste",
		"TG": "Togo",
		"TK": "Tokelau",
		"TO": "Tonga",
		"TT": "Trinidad and Tobago",
		"TN": "Tunisia",
		"TR": "Turkey",
		"TM": "Turkmenistan",
		"TC": "Turks and Caicos Islands",
		"TV": "Tuvalu",
		"UG": "Uganda",
		"UA": "Ukraine",
		"AE": "United Arab Emirates",
		"GB": "United Kingdom",
		"UM": "United States Minor Outlying Islands",
		"US": "United States of America",
		"UY": "Uruguay",
		"UZ": "Uzbekistan",
		"VU": "Vanuatu",
		"VE": "Venezuela",
		"VN": "Vietnam",
		"VG": "Virgin Islands (British)",
		"VI": "Virgin Islands (U.S.)",
		"WF": "Wallis and Futuna",
		"EH": "Western Sahara",
		"YE": "Yemen",
		"ZM": "Zambia",
		"ZW": "Zimbabwe",
		"EU": "European Union",
	}
	Languages = map[string]string{
		"ar":      "Arabic",
		"az":      "Azerbaijani",
		"bg":      "Bulgarian",
		"bn":      "Bengali",
		"ca":      "Catalan",
		"cs":      "Czech",
		"da":      "Danish",
		"de":      "German",
		"el":      "Greek",
		"en":      "English",
		"es":      "Spanish",
		"es-co":   "Colombian Spanish",
		"et":      "Estonian",
		"fa":      "Persian",
		"fi":      "Finnish",
		"fr":      "French",
		"hi":      "Hindi",
		"hu":      "Hungarian",
		"hy":      "Armenian",
		"id":      "Indonesian",
		"is":      "Icelandic",
		"it":      "Italian",
		"ja":      "Japanese",
		"ka":      "Georgian",
		"km":      "Khmer",
		"ko":      "Korean",
		"lt":      "Lithuanian",
		"mn":      "Mongolian",
		"my":      "Burmese",
		"nb":      "Norwegian",
		"nl":      "Dutch",
		"pl":      "Polish",
		"pt":      "Portuguese",
		"pt-br":   "Brazilian Portuguese",
		"ro":      "Romanian",
		"ru":      "Russian",
		"sk":      "Slovak",
		"sl":      "Slovenian",
		"sq":      "Albanian",
		"sr":      "Serbian",
		"sv":      "Swedish",
		"sw":      "Swahili",
		"ta":      "Tamil",
		"th":      "Thai",
		"tr":      "Turkish",
		"uk":      "Ukrainian",
		"vi":      "Vietnamese",
		"zh-hans": "Simplified Chinese",
		"zh-hant": "Traditional Chinese",
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
