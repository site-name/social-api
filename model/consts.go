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

func (t TimePeriodType) IsValid() bool {
	switch t {
	case DAY, WEEK, MONTH, YEAR:
		return true
	default:
		return false
	}
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
	DEFAULT_LOCALE                 = LanguageCodeEnumEn // this is default language also
	DEFAULT_COUNTRY                = CountryCodeUs
)

var (
	Countries                     map[CountryCode]string      // countries supported by app
	Languages                     map[LanguageCodeEnum]string // Languages supported by app
	MULTIPLE_COUNTRIES_MAX_LENGTH int                         // some model"s country fields contains multiple countries
	ReservedName                  []string                    // usernames that can only be used by system
	ValidUsernameChars            *regexp.Regexp              // regexp for username validation
	RestrictedUsernames           map[string]bool             // usernames that cannot be used
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

type LanguageCodeEnum string

func (l LanguageCodeEnum) IsValid() bool {
	return Languages[l] != ""
}

func (l LanguageCodeEnum) String() string {
	return *(*string)(unsafe.Pointer(&l))
}

const (
	LanguageCodeEnumAf           LanguageCodeEnum = "af"
	LanguageCodeEnumAfNa         LanguageCodeEnum = "af-na"
	LanguageCodeEnumId           LanguageCodeEnum = "id"
	LanguageCodeEnumAfZa         LanguageCodeEnum = "af-za"
	LanguageCodeEnumAgq          LanguageCodeEnum = "agq"
	LanguageCodeEnumAgqCm        LanguageCodeEnum = "agq-cm"
	LanguageCodeEnumAk           LanguageCodeEnum = "ak"
	LanguageCodeEnumAkGh         LanguageCodeEnum = "ak-gh"
	LanguageCodeEnumAm           LanguageCodeEnum = "am"
	LanguageCodeEnumAmEt         LanguageCodeEnum = "am-et"
	LanguageCodeEnumAr           LanguageCodeEnum = "ar"
	LanguageCodeEnumArAe         LanguageCodeEnum = "ar-ae"
	LanguageCodeEnumArBh         LanguageCodeEnum = "ar-bh"
	LanguageCodeEnumArDj         LanguageCodeEnum = "ar-dj"
	LanguageCodeEnumArDz         LanguageCodeEnum = "ar-dz"
	LanguageCodeEnumArEg         LanguageCodeEnum = "ar-eg"
	LanguageCodeEnumArEh         LanguageCodeEnum = "ar-eh"
	LanguageCodeEnumArEr         LanguageCodeEnum = "ar-er"
	LanguageCodeEnumArIl         LanguageCodeEnum = "ar-il"
	LanguageCodeEnumArIq         LanguageCodeEnum = "ar-iq"
	LanguageCodeEnumArJo         LanguageCodeEnum = "ar-jo"
	LanguageCodeEnumArKm         LanguageCodeEnum = "ar-km"
	LanguageCodeEnumArKw         LanguageCodeEnum = "ar-kw"
	LanguageCodeEnumArLb         LanguageCodeEnum = "ar-lb"
	LanguageCodeEnumArLy         LanguageCodeEnum = "ar-ly"
	LanguageCodeEnumArMa         LanguageCodeEnum = "ar-ma"
	LanguageCodeEnumArMr         LanguageCodeEnum = "ar-mr"
	LanguageCodeEnumArOm         LanguageCodeEnum = "ar-om"
	LanguageCodeEnumArPs         LanguageCodeEnum = "ar-ps"
	LanguageCodeEnumArQa         LanguageCodeEnum = "ar-qa"
	LanguageCodeEnumArSa         LanguageCodeEnum = "ar-sa"
	LanguageCodeEnumArSd         LanguageCodeEnum = "ar-sd"
	LanguageCodeEnumArSo         LanguageCodeEnum = "ar-so"
	LanguageCodeEnumArSs         LanguageCodeEnum = "ar-ss"
	LanguageCodeEnumArSy         LanguageCodeEnum = "ar-sy"
	LanguageCodeEnumArTd         LanguageCodeEnum = "ar-td"
	LanguageCodeEnumArTn         LanguageCodeEnum = "ar-tn"
	LanguageCodeEnumArYe         LanguageCodeEnum = "ar-ye"
	LanguageCodeEnumAs           LanguageCodeEnum = "as"
	LanguageCodeEnumAsIn         LanguageCodeEnum = "as-in"
	LanguageCodeEnumAsa          LanguageCodeEnum = "asa"
	LanguageCodeEnumAsaTz        LanguageCodeEnum = "asa-tz"
	LanguageCodeEnumAst          LanguageCodeEnum = "ast"
	LanguageCodeEnumAstEs        LanguageCodeEnum = "ast-es"
	LanguageCodeEnumAz           LanguageCodeEnum = "az"
	LanguageCodeEnumAzCyrl       LanguageCodeEnum = "az-cyrl"
	LanguageCodeEnumAzCyrlAz     LanguageCodeEnum = "az-cyrl-az"
	LanguageCodeEnumAzLatn       LanguageCodeEnum = "az-latn"
	LanguageCodeEnumAzLatnAz     LanguageCodeEnum = "az-latn-az"
	LanguageCodeEnumBas          LanguageCodeEnum = "bas"
	LanguageCodeEnumBasCm        LanguageCodeEnum = "bas-cm"
	LanguageCodeEnumBe           LanguageCodeEnum = "be"
	LanguageCodeEnumBeBy         LanguageCodeEnum = "be-by"
	LanguageCodeEnumBem          LanguageCodeEnum = "bem"
	LanguageCodeEnumBemZm        LanguageCodeEnum = "bem-zm"
	LanguageCodeEnumBez          LanguageCodeEnum = "bez"
	LanguageCodeEnumBezTz        LanguageCodeEnum = "bez-tz"
	LanguageCodeEnumBg           LanguageCodeEnum = "bg"
	LanguageCodeEnumBgBg         LanguageCodeEnum = "bg-bg"
	LanguageCodeEnumBm           LanguageCodeEnum = "bm"
	LanguageCodeEnumBmMl         LanguageCodeEnum = "bm-ml"
	LanguageCodeEnumBn           LanguageCodeEnum = "bn"
	LanguageCodeEnumBnBd         LanguageCodeEnum = "bn-bd"
	LanguageCodeEnumBnIn         LanguageCodeEnum = "bn-in"
	LanguageCodeEnumBo           LanguageCodeEnum = "bo"
	LanguageCodeEnumBoCn         LanguageCodeEnum = "bo-cn"
	LanguageCodeEnumBoIn         LanguageCodeEnum = "bo-in"
	LanguageCodeEnumBr           LanguageCodeEnum = "br"
	LanguageCodeEnumBrFr         LanguageCodeEnum = "br-fr"
	LanguageCodeEnumBrx          LanguageCodeEnum = "brx"
	LanguageCodeEnumBrxIn        LanguageCodeEnum = "brx-in"
	LanguageCodeEnumBs           LanguageCodeEnum = "bs"
	LanguageCodeEnumBsCyrl       LanguageCodeEnum = "bs-cyrl"
	LanguageCodeEnumBsCyrlBa     LanguageCodeEnum = "bs-cyrl-ba"
	LanguageCodeEnumBsLatn       LanguageCodeEnum = "bs-latn"
	LanguageCodeEnumBsLatnBa     LanguageCodeEnum = "bs-latn-ba"
	LanguageCodeEnumCa           LanguageCodeEnum = "ca"
	LanguageCodeEnumCaAd         LanguageCodeEnum = "ca-ad"
	LanguageCodeEnumCaEs         LanguageCodeEnum = "ca-es"
	LanguageCodeEnumCaEsValencia LanguageCodeEnum = "ca-es-valencia"
	LanguageCodeEnumCaFr         LanguageCodeEnum = "ca-fr"
	LanguageCodeEnumCaIt         LanguageCodeEnum = "ca-it"
	LanguageCodeEnumCcp          LanguageCodeEnum = "ccp"
	LanguageCodeEnumCcpBd        LanguageCodeEnum = "ccp-bd"
	LanguageCodeEnumCcpIn        LanguageCodeEnum = "ccp-in"
	LanguageCodeEnumCe           LanguageCodeEnum = "ce"
	LanguageCodeEnumCeRu         LanguageCodeEnum = "ce-ru"
	LanguageCodeEnumCeb          LanguageCodeEnum = "ceb"
	LanguageCodeEnumCebPh        LanguageCodeEnum = "ceb-ph"
	LanguageCodeEnumCgg          LanguageCodeEnum = "cgg"
	LanguageCodeEnumCggUg        LanguageCodeEnum = "cgg-ug"
	LanguageCodeEnumChr          LanguageCodeEnum = "chr"
	LanguageCodeEnumChrUs        LanguageCodeEnum = "chr-us"
	LanguageCodeEnumCkb          LanguageCodeEnum = "ckb"
	LanguageCodeEnumCkbIq        LanguageCodeEnum = "ckb-iq"
	LanguageCodeEnumCkbIr        LanguageCodeEnum = "ckb-ir"
	LanguageCodeEnumCs           LanguageCodeEnum = "cs"
	LanguageCodeEnumCsCz         LanguageCodeEnum = "cs-cz"
	LanguageCodeEnumCu           LanguageCodeEnum = "cu"
	LanguageCodeEnumCuRu         LanguageCodeEnum = "cu-ru"
	LanguageCodeEnumCy           LanguageCodeEnum = "cy"
	LanguageCodeEnumCyGb         LanguageCodeEnum = "cy-gb"
	LanguageCodeEnumDa           LanguageCodeEnum = "da"
	LanguageCodeEnumDaDk         LanguageCodeEnum = "da-dk"
	LanguageCodeEnumDaGl         LanguageCodeEnum = "da-gl"
	LanguageCodeEnumDav          LanguageCodeEnum = "dav"
	LanguageCodeEnumDavKe        LanguageCodeEnum = "dav-ke"
	LanguageCodeEnumDe           LanguageCodeEnum = "de"
	LanguageCodeEnumDeAt         LanguageCodeEnum = "de-at"
	LanguageCodeEnumDeBe         LanguageCodeEnum = "de-be"
	LanguageCodeEnumDeCh         LanguageCodeEnum = "de-ch"
	LanguageCodeEnumDeDe         LanguageCodeEnum = "de-de"
	LanguageCodeEnumDeIt         LanguageCodeEnum = "de-it"
	LanguageCodeEnumDeLi         LanguageCodeEnum = "de-li"
	LanguageCodeEnumDeLu         LanguageCodeEnum = "de-lu"
	LanguageCodeEnumDje          LanguageCodeEnum = "dje"
	LanguageCodeEnumDjeNe        LanguageCodeEnum = "dje-ne"
	LanguageCodeEnumDsb          LanguageCodeEnum = "dsb"
	LanguageCodeEnumDsbDe        LanguageCodeEnum = "dsb-de"
	LanguageCodeEnumDua          LanguageCodeEnum = "dua"
	LanguageCodeEnumDuaCm        LanguageCodeEnum = "dua-cm"
	LanguageCodeEnumDyo          LanguageCodeEnum = "dyo"
	LanguageCodeEnumDyoSn        LanguageCodeEnum = "dyo-sn"
	LanguageCodeEnumDz           LanguageCodeEnum = "dz"
	LanguageCodeEnumDzBt         LanguageCodeEnum = "dz-bt"
	LanguageCodeEnumEbu          LanguageCodeEnum = "ebu"
	LanguageCodeEnumEbuKe        LanguageCodeEnum = "ebu-ke"
	LanguageCodeEnumEe           LanguageCodeEnum = "ee"
	LanguageCodeEnumEeGh         LanguageCodeEnum = "ee-gh"
	LanguageCodeEnumEeTg         LanguageCodeEnum = "ee-tg"
	LanguageCodeEnumEl           LanguageCodeEnum = "el"
	LanguageCodeEnumElCy         LanguageCodeEnum = "el-cy"
	LanguageCodeEnumElGr         LanguageCodeEnum = "el-gr"
	LanguageCodeEnumEn           LanguageCodeEnum = "en"
	LanguageCodeEnumEnAe         LanguageCodeEnum = "en-ae"
	LanguageCodeEnumEnAg         LanguageCodeEnum = "en-ag"
	LanguageCodeEnumEnAi         LanguageCodeEnum = "en-ai"
	LanguageCodeEnumEnAs         LanguageCodeEnum = "en-as"
	LanguageCodeEnumEnAt         LanguageCodeEnum = "en-at"
	LanguageCodeEnumEnAu         LanguageCodeEnum = "en-au"
	LanguageCodeEnumEnBb         LanguageCodeEnum = "en-bb"
	LanguageCodeEnumEnBe         LanguageCodeEnum = "en-be"
	LanguageCodeEnumEnBi         LanguageCodeEnum = "en-bi"
	LanguageCodeEnumEnBm         LanguageCodeEnum = "en-bm"
	LanguageCodeEnumEnBs         LanguageCodeEnum = "en-bs"
	LanguageCodeEnumEnBw         LanguageCodeEnum = "en-bw"
	LanguageCodeEnumEnBz         LanguageCodeEnum = "en-bz"
	LanguageCodeEnumEnCa         LanguageCodeEnum = "en-ca"
	LanguageCodeEnumEnCc         LanguageCodeEnum = "en-cc"
	LanguageCodeEnumEnCh         LanguageCodeEnum = "en-ch"
	LanguageCodeEnumEnCk         LanguageCodeEnum = "en-ck"
	LanguageCodeEnumEnCm         LanguageCodeEnum = "en-cm"
	LanguageCodeEnumEnCx         LanguageCodeEnum = "en-cx"
	LanguageCodeEnumEnCy         LanguageCodeEnum = "en-cy"
	LanguageCodeEnumEnDe         LanguageCodeEnum = "en-de"
	LanguageCodeEnumEnDg         LanguageCodeEnum = "en-dg"
	LanguageCodeEnumEnDk         LanguageCodeEnum = "en-dk"
	LanguageCodeEnumEnDm         LanguageCodeEnum = "en-dm"
	LanguageCodeEnumEnEr         LanguageCodeEnum = "en-er"
	LanguageCodeEnumEnFi         LanguageCodeEnum = "en-fi"
	LanguageCodeEnumEnFj         LanguageCodeEnum = "en-fj"
	LanguageCodeEnumEnFk         LanguageCodeEnum = "en-fk"
	LanguageCodeEnumEnFm         LanguageCodeEnum = "en-fm"
	LanguageCodeEnumEnGb         LanguageCodeEnum = "en-gb"
	LanguageCodeEnumEnGd         LanguageCodeEnum = "en-gd"
	LanguageCodeEnumEnGg         LanguageCodeEnum = "en-gg"
	LanguageCodeEnumEnGh         LanguageCodeEnum = "en-gh"
	LanguageCodeEnumEnGi         LanguageCodeEnum = "en-gi"
	LanguageCodeEnumEnGm         LanguageCodeEnum = "en-gm"
	LanguageCodeEnumEnGu         LanguageCodeEnum = "en-gu"
	LanguageCodeEnumEnGy         LanguageCodeEnum = "en-gy"
	LanguageCodeEnumEnHk         LanguageCodeEnum = "en-hk"
	LanguageCodeEnumEnIe         LanguageCodeEnum = "en-ie"
	LanguageCodeEnumEnIl         LanguageCodeEnum = "en-il"
	LanguageCodeEnumEnIm         LanguageCodeEnum = "en-im"
	LanguageCodeEnumEnIn         LanguageCodeEnum = "en-in"
	LanguageCodeEnumEnIo         LanguageCodeEnum = "en-io"
	LanguageCodeEnumEnJe         LanguageCodeEnum = "en-je"
	LanguageCodeEnumEnJm         LanguageCodeEnum = "en-jm"
	LanguageCodeEnumEnKe         LanguageCodeEnum = "en-ke"
	LanguageCodeEnumEnKi         LanguageCodeEnum = "en-ki"
	LanguageCodeEnumEnKn         LanguageCodeEnum = "en-kn"
	LanguageCodeEnumEnKy         LanguageCodeEnum = "en-ky"
	LanguageCodeEnumEnLc         LanguageCodeEnum = "en-lc"
	LanguageCodeEnumEnLr         LanguageCodeEnum = "en-lr"
	LanguageCodeEnumEnLs         LanguageCodeEnum = "en-ls"
	LanguageCodeEnumEnMg         LanguageCodeEnum = "en-mg"
	LanguageCodeEnumEnMh         LanguageCodeEnum = "en-mh"
	LanguageCodeEnumEnMo         LanguageCodeEnum = "en-mo"
	LanguageCodeEnumEnMp         LanguageCodeEnum = "en-mp"
	LanguageCodeEnumEnMs         LanguageCodeEnum = "en-ms"
	LanguageCodeEnumEnMt         LanguageCodeEnum = "en-mt"
	LanguageCodeEnumEnMu         LanguageCodeEnum = "en-mu"
	LanguageCodeEnumEnMw         LanguageCodeEnum = "en-mw"
	LanguageCodeEnumEnMy         LanguageCodeEnum = "en-my"
	LanguageCodeEnumEnNa         LanguageCodeEnum = "en-na"
	LanguageCodeEnumEnNf         LanguageCodeEnum = "en-nf"
	LanguageCodeEnumEnNg         LanguageCodeEnum = "en-ng"
	LanguageCodeEnumEnNl         LanguageCodeEnum = "en-nl"
	LanguageCodeEnumEnNr         LanguageCodeEnum = "en-nr"
	LanguageCodeEnumEnNu         LanguageCodeEnum = "en-nu"
	LanguageCodeEnumEnNz         LanguageCodeEnum = "en-nz"
	LanguageCodeEnumEnPg         LanguageCodeEnum = "en-pg"
	LanguageCodeEnumEnPh         LanguageCodeEnum = "en-ph"
	LanguageCodeEnumEnPk         LanguageCodeEnum = "en-pk"
	LanguageCodeEnumEnPn         LanguageCodeEnum = "en-pn"
	LanguageCodeEnumEnPr         LanguageCodeEnum = "en-pr"
	LanguageCodeEnumEnPw         LanguageCodeEnum = "en-pw"
	LanguageCodeEnumEnRw         LanguageCodeEnum = "en-rw"
	LanguageCodeEnumEnSb         LanguageCodeEnum = "en-sb"
	LanguageCodeEnumEnSc         LanguageCodeEnum = "en-sc"
	LanguageCodeEnumEnSd         LanguageCodeEnum = "en-sd"
	LanguageCodeEnumEnSe         LanguageCodeEnum = "en-se"
	LanguageCodeEnumEnSg         LanguageCodeEnum = "en-sg"
	LanguageCodeEnumEnSh         LanguageCodeEnum = "en-sh"
	LanguageCodeEnumEnSi         LanguageCodeEnum = "en-si"
	LanguageCodeEnumEnSl         LanguageCodeEnum = "en-sl"
	LanguageCodeEnumEnSs         LanguageCodeEnum = "en-ss"
	LanguageCodeEnumEnSx         LanguageCodeEnum = "en-sx"
	LanguageCodeEnumEnSz         LanguageCodeEnum = "en-sz"
	LanguageCodeEnumEnTc         LanguageCodeEnum = "en-tc"
	LanguageCodeEnumEnTk         LanguageCodeEnum = "en-tk"
	LanguageCodeEnumEnTo         LanguageCodeEnum = "en-to"
	LanguageCodeEnumEnTt         LanguageCodeEnum = "en-tt"
	LanguageCodeEnumEnTv         LanguageCodeEnum = "en-tv"
	LanguageCodeEnumEnTz         LanguageCodeEnum = "en-tz"
	LanguageCodeEnumEnUg         LanguageCodeEnum = "en-ug"
	LanguageCodeEnumEnUm         LanguageCodeEnum = "en-um"
	LanguageCodeEnumEnUs         LanguageCodeEnum = "en-us"
	LanguageCodeEnumEnVc         LanguageCodeEnum = "en-vc"
	LanguageCodeEnumEnVg         LanguageCodeEnum = "en-vg"
	LanguageCodeEnumEnVi         LanguageCodeEnum = "en-vi"
	LanguageCodeEnumEnVu         LanguageCodeEnum = "en-vu"
	LanguageCodeEnumEnWs         LanguageCodeEnum = "en-ws"
	LanguageCodeEnumEnZa         LanguageCodeEnum = "en-za"
	LanguageCodeEnumEnZm         LanguageCodeEnum = "en-zm"
	LanguageCodeEnumEnZw         LanguageCodeEnum = "en-zw"
	LanguageCodeEnumEo           LanguageCodeEnum = "eo"
	LanguageCodeEnumEs           LanguageCodeEnum = "es"
	LanguageCodeEnumEsAr         LanguageCodeEnum = "es-ar"
	LanguageCodeEnumEsBo         LanguageCodeEnum = "es-bo"
	LanguageCodeEnumEsBr         LanguageCodeEnum = "es-br"
	LanguageCodeEnumEsBz         LanguageCodeEnum = "es-bz"
	LanguageCodeEnumEsCl         LanguageCodeEnum = "es-cl"
	LanguageCodeEnumEsCo         LanguageCodeEnum = "es-co"
	LanguageCodeEnumEsCr         LanguageCodeEnum = "es-cr"
	LanguageCodeEnumEsCu         LanguageCodeEnum = "es-cu"
	LanguageCodeEnumEsDo         LanguageCodeEnum = "es-do"
	LanguageCodeEnumEsEa         LanguageCodeEnum = "es-ea"
	LanguageCodeEnumEsEc         LanguageCodeEnum = "es-ec"
	LanguageCodeEnumEsEs         LanguageCodeEnum = "es-es"
	LanguageCodeEnumEsGq         LanguageCodeEnum = "es-gq"
	LanguageCodeEnumEsGt         LanguageCodeEnum = "es-gt"
	LanguageCodeEnumEsHn         LanguageCodeEnum = "es-hn"
	LanguageCodeEnumEsIc         LanguageCodeEnum = "es-ic"
	LanguageCodeEnumEsMx         LanguageCodeEnum = "es-mx"
	LanguageCodeEnumEsNi         LanguageCodeEnum = "es-ni"
	LanguageCodeEnumEsPa         LanguageCodeEnum = "es-pa"
	LanguageCodeEnumEsPe         LanguageCodeEnum = "es-pe"
	LanguageCodeEnumEsPh         LanguageCodeEnum = "es-ph"
	LanguageCodeEnumEsPr         LanguageCodeEnum = "es-pr"
	LanguageCodeEnumEsPy         LanguageCodeEnum = "es-py"
	LanguageCodeEnumEsSv         LanguageCodeEnum = "es-sv"
	LanguageCodeEnumEsUs         LanguageCodeEnum = "es-us"
	LanguageCodeEnumEsUy         LanguageCodeEnum = "es-uy"
	LanguageCodeEnumEsVe         LanguageCodeEnum = "es-ve"
	LanguageCodeEnumEt           LanguageCodeEnum = "et"
	LanguageCodeEnumEtEe         LanguageCodeEnum = "et-ee"
	LanguageCodeEnumEu           LanguageCodeEnum = "eu"
	LanguageCodeEnumEuEs         LanguageCodeEnum = "eu-es"
	LanguageCodeEnumEwo          LanguageCodeEnum = "ewo"
	LanguageCodeEnumEwoCm        LanguageCodeEnum = "ewo-cm"
	LanguageCodeEnumFa           LanguageCodeEnum = "fa"
	LanguageCodeEnumFaAf         LanguageCodeEnum = "fa-af"
	LanguageCodeEnumFaIr         LanguageCodeEnum = "fa-ir"
	LanguageCodeEnumFf           LanguageCodeEnum = "ff"
	LanguageCodeEnumFfAdlm       LanguageCodeEnum = "ff-adlm"
	LanguageCodeEnumFfAdlmBf     LanguageCodeEnum = "ff-adlm-bf"
	LanguageCodeEnumFfAdlmCm     LanguageCodeEnum = "ff-adlm-cm"
	LanguageCodeEnumFfAdlmGh     LanguageCodeEnum = "ff-adlm-gh"
	LanguageCodeEnumFfAdlmGm     LanguageCodeEnum = "ff-adlm-gm"
	LanguageCodeEnumFfAdlmGn     LanguageCodeEnum = "ff-adlm-gn"
	LanguageCodeEnumFfAdlmGw     LanguageCodeEnum = "ff-adlm-gw"
	LanguageCodeEnumFfAdlmLr     LanguageCodeEnum = "ff-adlm-lr"
	LanguageCodeEnumFfAdlmMr     LanguageCodeEnum = "ff-adlm-mr"
	LanguageCodeEnumFfAdlmNe     LanguageCodeEnum = "ff-adlm-ne"
	LanguageCodeEnumFfAdlmNg     LanguageCodeEnum = "ff-adlm-ng"
	LanguageCodeEnumFfAdlmSl     LanguageCodeEnum = "ff-adlm-sl"
	LanguageCodeEnumFfAdlmSn     LanguageCodeEnum = "ff-adlm-sn"
	LanguageCodeEnumFfLatn       LanguageCodeEnum = "ff-latn"
	LanguageCodeEnumFfLatnBf     LanguageCodeEnum = "ff-latn-bf"
	LanguageCodeEnumFfLatnCm     LanguageCodeEnum = "ff-latn-cm"
	LanguageCodeEnumFfLatnGh     LanguageCodeEnum = "ff-latn-gh"
	LanguageCodeEnumFfLatnGm     LanguageCodeEnum = "ff-latn-gm"
	LanguageCodeEnumFfLatnGn     LanguageCodeEnum = "ff-latn-gn"
	LanguageCodeEnumFfLatnGw     LanguageCodeEnum = "ff-latn-gw"
	LanguageCodeEnumFfLatnLr     LanguageCodeEnum = "ff-latn-lr"
	LanguageCodeEnumFfLatnMr     LanguageCodeEnum = "ff-latn-mr"
	LanguageCodeEnumFfLatnNe     LanguageCodeEnum = "ff-latn-ne"
	LanguageCodeEnumFfLatnNg     LanguageCodeEnum = "ff-latn-ng"
	LanguageCodeEnumFfLatnSl     LanguageCodeEnum = "ff-latn-sl"
	LanguageCodeEnumFfLatnSn     LanguageCodeEnum = "ff-latn-sn"
	LanguageCodeEnumFi           LanguageCodeEnum = "fi"
	LanguageCodeEnumFiFi         LanguageCodeEnum = "fi-fi"
	LanguageCodeEnumFil          LanguageCodeEnum = "fil"
	LanguageCodeEnumFilPh        LanguageCodeEnum = "fil-ph"
	LanguageCodeEnumFo           LanguageCodeEnum = "fo"
	LanguageCodeEnumFoDk         LanguageCodeEnum = "fo-dk"
	LanguageCodeEnumFoFo         LanguageCodeEnum = "fo-fo"
	LanguageCodeEnumFr           LanguageCodeEnum = "fr"
	LanguageCodeEnumFrBe         LanguageCodeEnum = "fr-be"
	LanguageCodeEnumFrBf         LanguageCodeEnum = "fr-bf"
	LanguageCodeEnumFrBi         LanguageCodeEnum = "fr-bi"
	LanguageCodeEnumFrBj         LanguageCodeEnum = "fr-bj"
	LanguageCodeEnumFrBl         LanguageCodeEnum = "fr-bl"
	LanguageCodeEnumFrCa         LanguageCodeEnum = "fr-ca"
	LanguageCodeEnumFrCd         LanguageCodeEnum = "fr-cd"
	LanguageCodeEnumFrCf         LanguageCodeEnum = "fr-cf"
	LanguageCodeEnumFrCg         LanguageCodeEnum = "fr-cg"
	LanguageCodeEnumFrCh         LanguageCodeEnum = "fr-ch"
	LanguageCodeEnumFrCi         LanguageCodeEnum = "fr-ci"
	LanguageCodeEnumFrCm         LanguageCodeEnum = "fr-cm"
	LanguageCodeEnumFrDj         LanguageCodeEnum = "fr-dj"
	LanguageCodeEnumFrDz         LanguageCodeEnum = "fr-dz"
	LanguageCodeEnumFrFr         LanguageCodeEnum = "fr-fr"
	LanguageCodeEnumFrGa         LanguageCodeEnum = "fr-ga"
	LanguageCodeEnumFrGf         LanguageCodeEnum = "fr-gf"
	LanguageCodeEnumFrGn         LanguageCodeEnum = "fr-gn"
	LanguageCodeEnumFrGp         LanguageCodeEnum = "fr-gp"
	LanguageCodeEnumFrGq         LanguageCodeEnum = "fr-gq"
	LanguageCodeEnumFrHt         LanguageCodeEnum = "fr-ht"
	LanguageCodeEnumFrKm         LanguageCodeEnum = "fr-km"
	LanguageCodeEnumFrLu         LanguageCodeEnum = "fr-lu"
	LanguageCodeEnumFrMa         LanguageCodeEnum = "fr-ma"
	LanguageCodeEnumFrMc         LanguageCodeEnum = "fr-mc"
	LanguageCodeEnumFrMf         LanguageCodeEnum = "fr-mf"
	LanguageCodeEnumFrMg         LanguageCodeEnum = "fr-mg"
	LanguageCodeEnumFrMl         LanguageCodeEnum = "fr-ml"
	LanguageCodeEnumFrMq         LanguageCodeEnum = "fr-mq"
	LanguageCodeEnumFrMr         LanguageCodeEnum = "fr-mr"
	LanguageCodeEnumFrMu         LanguageCodeEnum = "fr-mu"
	LanguageCodeEnumFrNc         LanguageCodeEnum = "fr-nc"
	LanguageCodeEnumFrNe         LanguageCodeEnum = "fr-ne"
	LanguageCodeEnumFrPf         LanguageCodeEnum = "fr-pf"
	LanguageCodeEnumFrPm         LanguageCodeEnum = "fr-pm"
	LanguageCodeEnumFrRe         LanguageCodeEnum = "fr-re"
	LanguageCodeEnumFrRw         LanguageCodeEnum = "fr-rw"
	LanguageCodeEnumFrSc         LanguageCodeEnum = "fr-sc"
	LanguageCodeEnumFrSn         LanguageCodeEnum = "fr-sn"
	LanguageCodeEnumFrSy         LanguageCodeEnum = "fr-sy"
	LanguageCodeEnumFrTd         LanguageCodeEnum = "fr-td"
	LanguageCodeEnumFrTg         LanguageCodeEnum = "fr-tg"
	LanguageCodeEnumFrTn         LanguageCodeEnum = "fr-tn"
	LanguageCodeEnumFrVu         LanguageCodeEnum = "fr-vu"
	LanguageCodeEnumFrWf         LanguageCodeEnum = "fr-wf"
	LanguageCodeEnumFrYt         LanguageCodeEnum = "fr-yt"
	LanguageCodeEnumFur          LanguageCodeEnum = "fur"
	LanguageCodeEnumFurIt        LanguageCodeEnum = "fur-it"
	LanguageCodeEnumFy           LanguageCodeEnum = "fy"
	LanguageCodeEnumFyNl         LanguageCodeEnum = "fy-nl"
	LanguageCodeEnumGa           LanguageCodeEnum = "ga"
	LanguageCodeEnumGaGb         LanguageCodeEnum = "ga-gb"
	LanguageCodeEnumGaIe         LanguageCodeEnum = "ga-ie"
	LanguageCodeEnumGd           LanguageCodeEnum = "gd"
	LanguageCodeEnumGdGb         LanguageCodeEnum = "gd-gb"
	LanguageCodeEnumGl           LanguageCodeEnum = "gl"
	LanguageCodeEnumGlEs         LanguageCodeEnum = "gl-es"
	LanguageCodeEnumGsw          LanguageCodeEnum = "gsw"
	LanguageCodeEnumGswCh        LanguageCodeEnum = "gsw-ch"
	LanguageCodeEnumGswFr        LanguageCodeEnum = "gsw-fr"
	LanguageCodeEnumGswLi        LanguageCodeEnum = "gsw-li"
	LanguageCodeEnumGu           LanguageCodeEnum = "gu"
	LanguageCodeEnumGuIn         LanguageCodeEnum = "gu-in"
	LanguageCodeEnumGuz          LanguageCodeEnum = "guz"
	LanguageCodeEnumGuzKe        LanguageCodeEnum = "guz-ke"
	LanguageCodeEnumGv           LanguageCodeEnum = "gv"
	LanguageCodeEnumGvIm         LanguageCodeEnum = "gv-im"
	LanguageCodeEnumHa           LanguageCodeEnum = "ha"
	LanguageCodeEnumHaGh         LanguageCodeEnum = "ha-gh"
	LanguageCodeEnumHaNe         LanguageCodeEnum = "ha-ne"
	LanguageCodeEnumHaNg         LanguageCodeEnum = "ha-ng"
	LanguageCodeEnumHaw          LanguageCodeEnum = "haw"
	LanguageCodeEnumHawUs        LanguageCodeEnum = "haw-us"
	LanguageCodeEnumHe           LanguageCodeEnum = "he"
	LanguageCodeEnumHeIl         LanguageCodeEnum = "he-il"
	LanguageCodeEnumHi           LanguageCodeEnum = "hi"
	LanguageCodeEnumHiIn         LanguageCodeEnum = "hi-in"
	LanguageCodeEnumHr           LanguageCodeEnum = "hr"
	LanguageCodeEnumHrBa         LanguageCodeEnum = "hr-ba"
	LanguageCodeEnumHrHr         LanguageCodeEnum = "hr-hr"
	LanguageCodeEnumHsb          LanguageCodeEnum = "hsb"
	LanguageCodeEnumHsbDe        LanguageCodeEnum = "hsb-de"
	LanguageCodeEnumHu           LanguageCodeEnum = "hu"
	LanguageCodeEnumHuHu         LanguageCodeEnum = "hu-hu"
	LanguageCodeEnumHy           LanguageCodeEnum = "hy"
	LanguageCodeEnumHyAm         LanguageCodeEnum = "hy-am"
	LanguageCodeEnumIa           LanguageCodeEnum = "ia"
	LanguageCodeEnumString       LanguageCodeEnum = "string"
	LanguageCodeEnumIDID         LanguageCodeEnum = "id-id"
	LanguageCodeEnumIg           LanguageCodeEnum = "ig"
	LanguageCodeEnumIgNg         LanguageCodeEnum = "ig-ng"
	LanguageCodeEnumIi           LanguageCodeEnum = "ii"
	LanguageCodeEnumIiCn         LanguageCodeEnum = "ii-cn"
	LanguageCodeEnumIs           LanguageCodeEnum = "is"
	LanguageCodeEnumIsIs         LanguageCodeEnum = "is-is"
	LanguageCodeEnumIt           LanguageCodeEnum = "it"
	LanguageCodeEnumItCh         LanguageCodeEnum = "it-ch"
	LanguageCodeEnumItIt         LanguageCodeEnum = "it-it"
	LanguageCodeEnumItSm         LanguageCodeEnum = "it-sm"
	LanguageCodeEnumItVa         LanguageCodeEnum = "it-va"
	LanguageCodeEnumJa           LanguageCodeEnum = "ja"
	LanguageCodeEnumJaJp         LanguageCodeEnum = "ja-jp"
	LanguageCodeEnumJgo          LanguageCodeEnum = "jgo"
	LanguageCodeEnumJgoCm        LanguageCodeEnum = "jgo-cm"
	LanguageCodeEnumJmc          LanguageCodeEnum = "jmc"
	LanguageCodeEnumJmcTz        LanguageCodeEnum = "jmc-tz"
	LanguageCodeEnumJv           LanguageCodeEnum = "jv"
	LanguageCodeEnumJvID         LanguageCodeEnum = "jv-id"
	LanguageCodeEnumKa           LanguageCodeEnum = "ka"
	LanguageCodeEnumKaGe         LanguageCodeEnum = "ka-ge"
	LanguageCodeEnumKab          LanguageCodeEnum = "kab"
	LanguageCodeEnumKabDz        LanguageCodeEnum = "kab-dz"
	LanguageCodeEnumKam          LanguageCodeEnum = "kam"
	LanguageCodeEnumKamKe        LanguageCodeEnum = "kam-ke"
	LanguageCodeEnumKde          LanguageCodeEnum = "kde"
	LanguageCodeEnumKdeTz        LanguageCodeEnum = "kde-tz"
	LanguageCodeEnumKea          LanguageCodeEnum = "kea"
	LanguageCodeEnumKeaCv        LanguageCodeEnum = "kea-cv"
	LanguageCodeEnumKhq          LanguageCodeEnum = "khq"
	LanguageCodeEnumKhqMl        LanguageCodeEnum = "khq-ml"
	LanguageCodeEnumKi           LanguageCodeEnum = "ki"
	LanguageCodeEnumKiKe         LanguageCodeEnum = "ki-ke"
	LanguageCodeEnumKk           LanguageCodeEnum = "kk"
	LanguageCodeEnumKkKz         LanguageCodeEnum = "kk-kz"
	LanguageCodeEnumKkj          LanguageCodeEnum = "kkj"
	LanguageCodeEnumKkjCm        LanguageCodeEnum = "kkj-cm"
	LanguageCodeEnumKl           LanguageCodeEnum = "kl"
	LanguageCodeEnumKlGl         LanguageCodeEnum = "kl-gl"
	LanguageCodeEnumKln          LanguageCodeEnum = "kln"
	LanguageCodeEnumKlnKe        LanguageCodeEnum = "kln-ke"
	LanguageCodeEnumKm           LanguageCodeEnum = "km"
	LanguageCodeEnumKmKh         LanguageCodeEnum = "km-kh"
	LanguageCodeEnumKn           LanguageCodeEnum = "kn"
	LanguageCodeEnumKnIn         LanguageCodeEnum = "kn-in"
	LanguageCodeEnumKo           LanguageCodeEnum = "ko"
	LanguageCodeEnumKoKp         LanguageCodeEnum = "ko-kp"
	LanguageCodeEnumKoKr         LanguageCodeEnum = "ko-kr"
	LanguageCodeEnumKok          LanguageCodeEnum = "kok"
	LanguageCodeEnumKokIn        LanguageCodeEnum = "kok-in"
	LanguageCodeEnumKs           LanguageCodeEnum = "ks"
	LanguageCodeEnumKsArab       LanguageCodeEnum = "ks-arab"
	LanguageCodeEnumKsArabIn     LanguageCodeEnum = "ks-arab-in"
	LanguageCodeEnumKsb          LanguageCodeEnum = "ksb"
	LanguageCodeEnumKsbTz        LanguageCodeEnum = "ksb-tz"
	LanguageCodeEnumKsf          LanguageCodeEnum = "ksf"
	LanguageCodeEnumKsfCm        LanguageCodeEnum = "ksf-cm"
	LanguageCodeEnumKsh          LanguageCodeEnum = "ksh"
	LanguageCodeEnumKshDe        LanguageCodeEnum = "ksh-de"
	LanguageCodeEnumKu           LanguageCodeEnum = "ku"
	LanguageCodeEnumKuTr         LanguageCodeEnum = "ku-tr"
	LanguageCodeEnumKw           LanguageCodeEnum = "kw"
	LanguageCodeEnumKwGb         LanguageCodeEnum = "kw-gb"
	LanguageCodeEnumKy           LanguageCodeEnum = "ky"
	LanguageCodeEnumKyKg         LanguageCodeEnum = "ky-kg"
	LanguageCodeEnumLag          LanguageCodeEnum = "lag"
	LanguageCodeEnumLagTz        LanguageCodeEnum = "lag-tz"
	LanguageCodeEnumLb           LanguageCodeEnum = "lb"
	LanguageCodeEnumLbLu         LanguageCodeEnum = "lb-lu"
	LanguageCodeEnumLg           LanguageCodeEnum = "lg"
	LanguageCodeEnumLgUg         LanguageCodeEnum = "lg-ug"
	LanguageCodeEnumLkt          LanguageCodeEnum = "lkt"
	LanguageCodeEnumLktUs        LanguageCodeEnum = "lkt-us"
	LanguageCodeEnumLn           LanguageCodeEnum = "ln"
	LanguageCodeEnumLnAo         LanguageCodeEnum = "ln-ao"
	LanguageCodeEnumLnCd         LanguageCodeEnum = "ln-cd"
	LanguageCodeEnumLnCf         LanguageCodeEnum = "ln-cf"
	LanguageCodeEnumLnCg         LanguageCodeEnum = "ln-cg"
	LanguageCodeEnumLo           LanguageCodeEnum = "lo"
	LanguageCodeEnumLoLa         LanguageCodeEnum = "lo-la"
	LanguageCodeEnumLrc          LanguageCodeEnum = "lrc"
	LanguageCodeEnumLrcIq        LanguageCodeEnum = "lrc-iq"
	LanguageCodeEnumLrcIr        LanguageCodeEnum = "lrc-ir"
	LanguageCodeEnumLt           LanguageCodeEnum = "lt"
	LanguageCodeEnumLtLt         LanguageCodeEnum = "lt-lt"
	LanguageCodeEnumLu           LanguageCodeEnum = "lu"
	LanguageCodeEnumLuCd         LanguageCodeEnum = "lu-cd"
	LanguageCodeEnumLuo          LanguageCodeEnum = "luo"
	LanguageCodeEnumLuoKe        LanguageCodeEnum = "luo-ke"
	LanguageCodeEnumLuy          LanguageCodeEnum = "luy"
	LanguageCodeEnumLuyKe        LanguageCodeEnum = "luy-ke"
	LanguageCodeEnumLv           LanguageCodeEnum = "lv"
	LanguageCodeEnumLvLv         LanguageCodeEnum = "lv-lv"
	LanguageCodeEnumMai          LanguageCodeEnum = "mai"
	LanguageCodeEnumMaiIn        LanguageCodeEnum = "mai-in"
	LanguageCodeEnumMas          LanguageCodeEnum = "mas"
	LanguageCodeEnumMasKe        LanguageCodeEnum = "mas-ke"
	LanguageCodeEnumMasTz        LanguageCodeEnum = "mas-tz"
	LanguageCodeEnumMer          LanguageCodeEnum = "mer"
	LanguageCodeEnumMerKe        LanguageCodeEnum = "mer-ke"
	LanguageCodeEnumMfe          LanguageCodeEnum = "mfe"
	LanguageCodeEnumMfeMu        LanguageCodeEnum = "mfe-mu"
	LanguageCodeEnumMg           LanguageCodeEnum = "mg"
	LanguageCodeEnumMgMg         LanguageCodeEnum = "mg-mg"
	LanguageCodeEnumMgh          LanguageCodeEnum = "mgh"
	LanguageCodeEnumMghMz        LanguageCodeEnum = "mgh-mz"
	LanguageCodeEnumMgo          LanguageCodeEnum = "mgo"
	LanguageCodeEnumMgoCm        LanguageCodeEnum = "mgo-cm"
	LanguageCodeEnumMi           LanguageCodeEnum = "mi"
	LanguageCodeEnumMiNz         LanguageCodeEnum = "mi-nz"
	LanguageCodeEnumMk           LanguageCodeEnum = "mk"
	LanguageCodeEnumMkMk         LanguageCodeEnum = "mk-mk"
	LanguageCodeEnumMl           LanguageCodeEnum = "ml"
	LanguageCodeEnumMlIn         LanguageCodeEnum = "ml-in"
	LanguageCodeEnumMn           LanguageCodeEnum = "mn"
	LanguageCodeEnumMnMn         LanguageCodeEnum = "mn-mn"
	LanguageCodeEnumMni          LanguageCodeEnum = "mni"
	LanguageCodeEnumMniBeng      LanguageCodeEnum = "mni-beng"
	LanguageCodeEnumMniBengIn    LanguageCodeEnum = "mni-beng-in"
	LanguageCodeEnumMr           LanguageCodeEnum = "mr"
	LanguageCodeEnumMrIn         LanguageCodeEnum = "mr-in"
	LanguageCodeEnumMs           LanguageCodeEnum = "ms"
	LanguageCodeEnumMsBn         LanguageCodeEnum = "ms-bn"
	LanguageCodeEnumMsID         LanguageCodeEnum = "ms-id"
	LanguageCodeEnumMsMy         LanguageCodeEnum = "ms-my"
	LanguageCodeEnumMsSg         LanguageCodeEnum = "ms-sg"
	LanguageCodeEnumMt           LanguageCodeEnum = "mt"
	LanguageCodeEnumMtMt         LanguageCodeEnum = "mt-mt"
	LanguageCodeEnumMua          LanguageCodeEnum = "mua"
	LanguageCodeEnumMuaCm        LanguageCodeEnum = "mua-cm"
	LanguageCodeEnumMy           LanguageCodeEnum = "my"
	LanguageCodeEnumMyMm         LanguageCodeEnum = "my-mm"
	LanguageCodeEnumMzn          LanguageCodeEnum = "mzn"
	LanguageCodeEnumMznIr        LanguageCodeEnum = "mzn-ir"
	LanguageCodeEnumNaq          LanguageCodeEnum = "naq"
	LanguageCodeEnumNaqNa        LanguageCodeEnum = "naq-na"
	LanguageCodeEnumNb           LanguageCodeEnum = "nb"
	LanguageCodeEnumNbNo         LanguageCodeEnum = "nb-no"
	LanguageCodeEnumNbSj         LanguageCodeEnum = "nb-sj"
	LanguageCodeEnumNd           LanguageCodeEnum = "nd"
	LanguageCodeEnumNdZw         LanguageCodeEnum = "nd-zw"
	LanguageCodeEnumNds          LanguageCodeEnum = "nds"
	LanguageCodeEnumNdsDe        LanguageCodeEnum = "nds-de"
	LanguageCodeEnumNdsNl        LanguageCodeEnum = "nds-nl"
	LanguageCodeEnumNe           LanguageCodeEnum = "ne"
	LanguageCodeEnumNeIn         LanguageCodeEnum = "ne-in"
	LanguageCodeEnumNeNp         LanguageCodeEnum = "ne-np"
	LanguageCodeEnumNl           LanguageCodeEnum = "nl"
	LanguageCodeEnumNlAw         LanguageCodeEnum = "nl-aw"
	LanguageCodeEnumNlBe         LanguageCodeEnum = "nl-be"
	LanguageCodeEnumNlBq         LanguageCodeEnum = "nl-bq"
	LanguageCodeEnumNlCw         LanguageCodeEnum = "nl-cw"
	LanguageCodeEnumNlNl         LanguageCodeEnum = "nl-nl"
	LanguageCodeEnumNlSr         LanguageCodeEnum = "nl-sr"
	LanguageCodeEnumNlSx         LanguageCodeEnum = "nl-sx"
	LanguageCodeEnumNmg          LanguageCodeEnum = "nmg"
	LanguageCodeEnumNmgCm        LanguageCodeEnum = "nmg-cm"
	LanguageCodeEnumNn           LanguageCodeEnum = "nn"
	LanguageCodeEnumNnNo         LanguageCodeEnum = "nn-no"
	LanguageCodeEnumNnh          LanguageCodeEnum = "nnh"
	LanguageCodeEnumNnhCm        LanguageCodeEnum = "nnh-cm"
	LanguageCodeEnumNus          LanguageCodeEnum = "nus"
	LanguageCodeEnumNusSs        LanguageCodeEnum = "nus-ss"
	LanguageCodeEnumNyn          LanguageCodeEnum = "nyn"
	LanguageCodeEnumNynUg        LanguageCodeEnum = "nyn-ug"
	LanguageCodeEnumOm           LanguageCodeEnum = "om"
	LanguageCodeEnumOmEt         LanguageCodeEnum = "om-et"
	LanguageCodeEnumOmKe         LanguageCodeEnum = "om-ke"
	LanguageCodeEnumOr           LanguageCodeEnum = "or"
	LanguageCodeEnumOrIn         LanguageCodeEnum = "or-in"
	LanguageCodeEnumOs           LanguageCodeEnum = "os"
	LanguageCodeEnumOsGe         LanguageCodeEnum = "os-ge"
	LanguageCodeEnumOsRu         LanguageCodeEnum = "os-ru"
	LanguageCodeEnumPa           LanguageCodeEnum = "pa"
	LanguageCodeEnumPaArab       LanguageCodeEnum = "pa-arab"
	LanguageCodeEnumPaArabPk     LanguageCodeEnum = "pa-arab-pk"
	LanguageCodeEnumPaGuru       LanguageCodeEnum = "pa-guru"
	LanguageCodeEnumPaGuruIn     LanguageCodeEnum = "pa-guru-in"
	LanguageCodeEnumPcm          LanguageCodeEnum = "pcm"
	LanguageCodeEnumPcmNg        LanguageCodeEnum = "pcm-ng"
	LanguageCodeEnumPl           LanguageCodeEnum = "pl"
	LanguageCodeEnumPlPl         LanguageCodeEnum = "pl-pl"
	LanguageCodeEnumPrg          LanguageCodeEnum = "prg"
	LanguageCodeEnumPs           LanguageCodeEnum = "ps"
	LanguageCodeEnumPsAf         LanguageCodeEnum = "ps-af"
	LanguageCodeEnumPsPk         LanguageCodeEnum = "ps-pk"
	LanguageCodeEnumPt           LanguageCodeEnum = "pt"
	LanguageCodeEnumPtAo         LanguageCodeEnum = "pt-ao"
	LanguageCodeEnumPtBr         LanguageCodeEnum = "pt-br"
	LanguageCodeEnumPtCh         LanguageCodeEnum = "pt-ch"
	LanguageCodeEnumPtCv         LanguageCodeEnum = "pt-cv"
	LanguageCodeEnumPtGq         LanguageCodeEnum = "pt-gq"
	LanguageCodeEnumPtGw         LanguageCodeEnum = "pt-gw"
	LanguageCodeEnumPtLu         LanguageCodeEnum = "pt-lu"
	LanguageCodeEnumPtMo         LanguageCodeEnum = "pt-mo"
	LanguageCodeEnumPtMz         LanguageCodeEnum = "pt-mz"
	LanguageCodeEnumPtPt         LanguageCodeEnum = "pt-pt"
	LanguageCodeEnumPtSt         LanguageCodeEnum = "pt-st"
	LanguageCodeEnumPtTl         LanguageCodeEnum = "pt-tl"
	LanguageCodeEnumQu           LanguageCodeEnum = "qu"
	LanguageCodeEnumQuBo         LanguageCodeEnum = "qu-bo"
	LanguageCodeEnumQuEc         LanguageCodeEnum = "qu-ec"
	LanguageCodeEnumQuPe         LanguageCodeEnum = "qu-pe"
	LanguageCodeEnumRm           LanguageCodeEnum = "rm"
	LanguageCodeEnumRmCh         LanguageCodeEnum = "rm-ch"
	LanguageCodeEnumRn           LanguageCodeEnum = "rn"
	LanguageCodeEnumRnBi         LanguageCodeEnum = "rn-bi"
	LanguageCodeEnumRo           LanguageCodeEnum = "ro"
	LanguageCodeEnumRoMd         LanguageCodeEnum = "ro-md"
	LanguageCodeEnumRoRo         LanguageCodeEnum = "ro-ro"
	LanguageCodeEnumRof          LanguageCodeEnum = "rof"
	LanguageCodeEnumRofTz        LanguageCodeEnum = "rof-tz"
	LanguageCodeEnumRu           LanguageCodeEnum = "ru"
	LanguageCodeEnumRuBy         LanguageCodeEnum = "ru-by"
	LanguageCodeEnumRuKg         LanguageCodeEnum = "ru-kg"
	LanguageCodeEnumRuKz         LanguageCodeEnum = "ru-kz"
	LanguageCodeEnumRuMd         LanguageCodeEnum = "ru-md"
	LanguageCodeEnumRuRu         LanguageCodeEnum = "ru-ru"
	LanguageCodeEnumRuUa         LanguageCodeEnum = "ru-ua"
	LanguageCodeEnumRw           LanguageCodeEnum = "rw"
	LanguageCodeEnumRwRw         LanguageCodeEnum = "rw-rw"
	LanguageCodeEnumRwk          LanguageCodeEnum = "rwk"
	LanguageCodeEnumRwkTz        LanguageCodeEnum = "rwk-tz"
	LanguageCodeEnumSah          LanguageCodeEnum = "sah"
	LanguageCodeEnumSahRu        LanguageCodeEnum = "sah-ru"
	LanguageCodeEnumSaq          LanguageCodeEnum = "saq"
	LanguageCodeEnumSaqKe        LanguageCodeEnum = "saq-ke"
	LanguageCodeEnumSat          LanguageCodeEnum = "sat"
	LanguageCodeEnumSatOlck      LanguageCodeEnum = "sat-olck"
	LanguageCodeEnumSatOlckIn    LanguageCodeEnum = "sat-olck-in"
	LanguageCodeEnumSbp          LanguageCodeEnum = "sbp"
	LanguageCodeEnumSbpTz        LanguageCodeEnum = "sbp-tz"
	LanguageCodeEnumSd           LanguageCodeEnum = "sd"
	LanguageCodeEnumSdArab       LanguageCodeEnum = "sd-arab"
	LanguageCodeEnumSdArabPk     LanguageCodeEnum = "sd-arab-pk"
	LanguageCodeEnumSdDeva       LanguageCodeEnum = "sd-deva"
	LanguageCodeEnumSdDevaIn     LanguageCodeEnum = "sd-deva-in"
	LanguageCodeEnumSe           LanguageCodeEnum = "se"
	LanguageCodeEnumSeFi         LanguageCodeEnum = "se-fi"
	LanguageCodeEnumSeNo         LanguageCodeEnum = "se-no"
	LanguageCodeEnumSeSe         LanguageCodeEnum = "se-se"
	LanguageCodeEnumSeh          LanguageCodeEnum = "seh"
	LanguageCodeEnumSehMz        LanguageCodeEnum = "seh-mz"
	LanguageCodeEnumSes          LanguageCodeEnum = "ses"
	LanguageCodeEnumSesMl        LanguageCodeEnum = "ses-ml"
	LanguageCodeEnumSg           LanguageCodeEnum = "sg"
	LanguageCodeEnumSgCf         LanguageCodeEnum = "sg-cf"
	LanguageCodeEnumShi          LanguageCodeEnum = "shi"
	LanguageCodeEnumShiLatn      LanguageCodeEnum = "shi-latn"
	LanguageCodeEnumShiLatnMa    LanguageCodeEnum = "shi-latn-ma"
	LanguageCodeEnumShiTfng      LanguageCodeEnum = "shi-tfng"
	LanguageCodeEnumShiTfngMa    LanguageCodeEnum = "shi-tfng-ma"
	LanguageCodeEnumSi           LanguageCodeEnum = "si"
	LanguageCodeEnumSiLk         LanguageCodeEnum = "si-lk"
	LanguageCodeEnumSk           LanguageCodeEnum = "sk"
	LanguageCodeEnumSkSk         LanguageCodeEnum = "sk-sk"
	LanguageCodeEnumSl           LanguageCodeEnum = "sl"
	LanguageCodeEnumSlSi         LanguageCodeEnum = "sl-si"
	LanguageCodeEnumSmn          LanguageCodeEnum = "smn"
	LanguageCodeEnumSmnFi        LanguageCodeEnum = "smn-fi"
	LanguageCodeEnumSn           LanguageCodeEnum = "sn"
	LanguageCodeEnumSnZw         LanguageCodeEnum = "sn-zw"
	LanguageCodeEnumSo           LanguageCodeEnum = "so"
	LanguageCodeEnumSoDj         LanguageCodeEnum = "so-dj"
	LanguageCodeEnumSoEt         LanguageCodeEnum = "so-et"
	LanguageCodeEnumSoKe         LanguageCodeEnum = "so-ke"
	LanguageCodeEnumSoSo         LanguageCodeEnum = "so-so"
	LanguageCodeEnumSq           LanguageCodeEnum = "sq"
	LanguageCodeEnumSqAl         LanguageCodeEnum = "sq-al"
	LanguageCodeEnumSqMk         LanguageCodeEnum = "sq-mk"
	LanguageCodeEnumSqXk         LanguageCodeEnum = "sq-xk"
	LanguageCodeEnumSr           LanguageCodeEnum = "sr"
	LanguageCodeEnumSrCyrl       LanguageCodeEnum = "sr-cyrl"
	LanguageCodeEnumSrCyrlBa     LanguageCodeEnum = "sr-cyrl-ba"
	LanguageCodeEnumSrCyrlMe     LanguageCodeEnum = "sr-cyrl-me"
	LanguageCodeEnumSrCyrlRs     LanguageCodeEnum = "sr-cyrl-rs"
	LanguageCodeEnumSrCyrlXk     LanguageCodeEnum = "sr-cyrl-xk"
	LanguageCodeEnumSrLatn       LanguageCodeEnum = "sr-latn"
	LanguageCodeEnumSrLatnBa     LanguageCodeEnum = "sr-latn-ba"
	LanguageCodeEnumSrLatnMe     LanguageCodeEnum = "sr-latn-me"
	LanguageCodeEnumSrLatnRs     LanguageCodeEnum = "sr-latn-rs"
	LanguageCodeEnumSrLatnXk     LanguageCodeEnum = "sr-latn-xk"
	LanguageCodeEnumSu           LanguageCodeEnum = "su"
	LanguageCodeEnumSuLatn       LanguageCodeEnum = "su-latn"
	LanguageCodeEnumSuLatnID     LanguageCodeEnum = "su-latn-id"
	LanguageCodeEnumSv           LanguageCodeEnum = "sv"
	LanguageCodeEnumSvAx         LanguageCodeEnum = "sv-ax"
	LanguageCodeEnumSvFi         LanguageCodeEnum = "sv-fi"
	LanguageCodeEnumSvSe         LanguageCodeEnum = "sv-se"
	LanguageCodeEnumSw           LanguageCodeEnum = "sw"
	LanguageCodeEnumSwCd         LanguageCodeEnum = "sw-cd"
	LanguageCodeEnumSwKe         LanguageCodeEnum = "sw-ke"
	LanguageCodeEnumSwTz         LanguageCodeEnum = "sw-tz"
	LanguageCodeEnumSwUg         LanguageCodeEnum = "sw-ug"
	LanguageCodeEnumTa           LanguageCodeEnum = "ta"
	LanguageCodeEnumTaIn         LanguageCodeEnum = "ta-in"
	LanguageCodeEnumTaLk         LanguageCodeEnum = "ta-lk"
	LanguageCodeEnumTaMy         LanguageCodeEnum = "ta-my"
	LanguageCodeEnumTaSg         LanguageCodeEnum = "ta-sg"
	LanguageCodeEnumTe           LanguageCodeEnum = "te"
	LanguageCodeEnumTeIn         LanguageCodeEnum = "te-in"
	LanguageCodeEnumTeo          LanguageCodeEnum = "teo"
	LanguageCodeEnumTeoKe        LanguageCodeEnum = "teo-ke"
	LanguageCodeEnumTeoUg        LanguageCodeEnum = "teo-ug"
	LanguageCodeEnumTg           LanguageCodeEnum = "tg"
	LanguageCodeEnumTgTj         LanguageCodeEnum = "tg-tj"
	LanguageCodeEnumTh           LanguageCodeEnum = "th"
	LanguageCodeEnumThTh         LanguageCodeEnum = "th-th"
	LanguageCodeEnumTi           LanguageCodeEnum = "ti"
	LanguageCodeEnumTiEr         LanguageCodeEnum = "ti-er"
	LanguageCodeEnumTiEt         LanguageCodeEnum = "ti-et"
	LanguageCodeEnumTk           LanguageCodeEnum = "tk"
	LanguageCodeEnumTkTm         LanguageCodeEnum = "tk-tm"
	LanguageCodeEnumTo           LanguageCodeEnum = "to"
	LanguageCodeEnumToTo         LanguageCodeEnum = "to-to"
	LanguageCodeEnumTr           LanguageCodeEnum = "tr"
	LanguageCodeEnumTrCy         LanguageCodeEnum = "tr-cy"
	LanguageCodeEnumTrTr         LanguageCodeEnum = "tr-tr"
	LanguageCodeEnumTt           LanguageCodeEnum = "tt"
	LanguageCodeEnumTtRu         LanguageCodeEnum = "tt-ru"
	LanguageCodeEnumTwq          LanguageCodeEnum = "twq"
	LanguageCodeEnumTwqNe        LanguageCodeEnum = "twq-ne"
	LanguageCodeEnumTzm          LanguageCodeEnum = "tzm"
	LanguageCodeEnumTzmMa        LanguageCodeEnum = "tzm-ma"
	LanguageCodeEnumUg           LanguageCodeEnum = "ug"
	LanguageCodeEnumUgCn         LanguageCodeEnum = "ug-cn"
	LanguageCodeEnumUk           LanguageCodeEnum = "uk"
	LanguageCodeEnumUkUa         LanguageCodeEnum = "uk-ua"
	LanguageCodeEnumUr           LanguageCodeEnum = "ur"
	LanguageCodeEnumUrIn         LanguageCodeEnum = "ur-in"
	LanguageCodeEnumUrPk         LanguageCodeEnum = "ur-pk"
	LanguageCodeEnumUz           LanguageCodeEnum = "uz"
	LanguageCodeEnumUzArab       LanguageCodeEnum = "uz-arab"
	LanguageCodeEnumUzArabAf     LanguageCodeEnum = "uz-arab-af"
	LanguageCodeEnumUzCyrl       LanguageCodeEnum = "uz-cyrl"
	LanguageCodeEnumUzCyrlUz     LanguageCodeEnum = "uz-cyrl-uz"
	LanguageCodeEnumUzLatn       LanguageCodeEnum = "uz-latn"
	LanguageCodeEnumUzLatnUz     LanguageCodeEnum = "uz-latn-uz"
	LanguageCodeEnumVai          LanguageCodeEnum = "vai"
	LanguageCodeEnumVaiLatn      LanguageCodeEnum = "vai-latn"
	LanguageCodeEnumVaiLatnLr    LanguageCodeEnum = "vai-latn-lr"
	LanguageCodeEnumVaiVaii      LanguageCodeEnum = "vai-vaii"
	LanguageCodeEnumVaiVaiiLr    LanguageCodeEnum = "vai-vaii-lr"
	LanguageCodeEnumVi           LanguageCodeEnum = "vi"
	LanguageCodeEnumViVn         LanguageCodeEnum = "vi-vn"
	LanguageCodeEnumVo           LanguageCodeEnum = "vo"
	LanguageCodeEnumVun          LanguageCodeEnum = "vun"
	LanguageCodeEnumVunTz        LanguageCodeEnum = "vun-tz"
	LanguageCodeEnumWae          LanguageCodeEnum = "wae"
	LanguageCodeEnumWaeCh        LanguageCodeEnum = "wae-ch"
	LanguageCodeEnumWo           LanguageCodeEnum = "wo"
	LanguageCodeEnumWoSn         LanguageCodeEnum = "wo-sn"
	LanguageCodeEnumXh           LanguageCodeEnum = "xh"
	LanguageCodeEnumXhZa         LanguageCodeEnum = "xh-za"
	LanguageCodeEnumXog          LanguageCodeEnum = "xog"
	LanguageCodeEnumXogUg        LanguageCodeEnum = "xog-ug"
	LanguageCodeEnumYav          LanguageCodeEnum = "yav"
	LanguageCodeEnumYavCm        LanguageCodeEnum = "yav-cm"
	LanguageCodeEnumYi           LanguageCodeEnum = "yi"
	LanguageCodeEnumYo           LanguageCodeEnum = "yo"
	LanguageCodeEnumYoBj         LanguageCodeEnum = "yo-bj"
	LanguageCodeEnumYoNg         LanguageCodeEnum = "yo-ng"
	LanguageCodeEnumYue          LanguageCodeEnum = "yue"
	LanguageCodeEnumYueHans      LanguageCodeEnum = "yue-hans"
	LanguageCodeEnumYueHansCn    LanguageCodeEnum = "yue-hans-cn"
	LanguageCodeEnumYueHant      LanguageCodeEnum = "yue-hant"
	LanguageCodeEnumYueHantHk    LanguageCodeEnum = "yue-hant-hk"
	LanguageCodeEnumZgh          LanguageCodeEnum = "zgh"
	LanguageCodeEnumZghMa        LanguageCodeEnum = "zgh-ma"
	LanguageCodeEnumZh           LanguageCodeEnum = "zh"
	LanguageCodeEnumZhHans       LanguageCodeEnum = "zh-hans"
	LanguageCodeEnumZhHansCn     LanguageCodeEnum = "zh-hans-cn"
	LanguageCodeEnumZhHansHk     LanguageCodeEnum = "zh-hans-hk"
	LanguageCodeEnumZhHansMo     LanguageCodeEnum = "zh-hans-mo"
	LanguageCodeEnumZhHansSg     LanguageCodeEnum = "zh-hans-sg"
	LanguageCodeEnumZhHant       LanguageCodeEnum = "zh-hant"
	LanguageCodeEnumZhHantHk     LanguageCodeEnum = "zh-hant-hk"
	LanguageCodeEnumZhHantMo     LanguageCodeEnum = "zh-hant-mo"
	LanguageCodeEnumZhHantTw     LanguageCodeEnum = "zh-hant-tw"
	LanguageCodeEnumZu           LanguageCodeEnum = "zu"
	LanguageCodeEnumZuZa         LanguageCodeEnum = "zu-za"
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

	Languages = map[LanguageCodeEnum]string{
		LanguageCodeEnumAf:           "Afrikaans",
		LanguageCodeEnumAfNa:         "Afrikaans (Namibia)",
		LanguageCodeEnumAfZa:         "Afrikaans (South Africa)",
		LanguageCodeEnumAgq:          "Aghem",
		LanguageCodeEnumAgqCm:        "Aghem (Cameroon)",
		LanguageCodeEnumAk:           "Akan",
		LanguageCodeEnumAkGh:         "Akan (Ghana)",
		LanguageCodeEnumAm:           "Amharic",
		LanguageCodeEnumAmEt:         "Amharic (Ethiopia)",
		LanguageCodeEnumAr:           "Arabic",
		LanguageCodeEnumArAe:         "Arabic (United Arab Emirates)",
		LanguageCodeEnumArBh:         "Arabic (Bahrain)",
		LanguageCodeEnumArDj:         "Arabic (Djibouti)",
		LanguageCodeEnumArDz:         "Arabic (Algeria)",
		LanguageCodeEnumArEg:         "Arabic (Egypt)",
		LanguageCodeEnumArEh:         "Arabic (Western Sahara)",
		LanguageCodeEnumArEr:         "Arabic (Eritrea)",
		LanguageCodeEnumArIl:         "Arabic (Israel)",
		LanguageCodeEnumArIq:         "Arabic (Iraq)",
		LanguageCodeEnumArJo:         "Arabic (Jordan)",
		LanguageCodeEnumArKm:         "Arabic (Comoros)",
		LanguageCodeEnumArKw:         "Arabic (Kuwait)",
		LanguageCodeEnumArLb:         "Arabic (Lebanon)",
		LanguageCodeEnumArLy:         "Arabic (Libya)",
		LanguageCodeEnumArMa:         "Arabic (Morocco)",
		LanguageCodeEnumArMr:         "Arabic (Mauritania)",
		LanguageCodeEnumArOm:         "Arabic (Oman)",
		LanguageCodeEnumArPs:         "Arabic (Palestinian Territories)",
		LanguageCodeEnumArQa:         "Arabic (Qatar)",
		LanguageCodeEnumArSa:         "Arabic (Saudi Arabia)",
		LanguageCodeEnumArSd:         "Arabic (Sudan)",
		LanguageCodeEnumArSo:         "Arabic (Somalia)",
		LanguageCodeEnumArSs:         "Arabic (South Sudan)",
		LanguageCodeEnumArSy:         "Arabic (Syria)",
		LanguageCodeEnumArTd:         "Arabic (Chad)",
		LanguageCodeEnumArTn:         "Arabic (Tunisia)",
		LanguageCodeEnumArYe:         "Arabic (Yemen)",
		LanguageCodeEnumAs:           "Assamese",
		LanguageCodeEnumAsIn:         "Assamese (India)",
		LanguageCodeEnumAsa:          "Asu",
		LanguageCodeEnumAsaTz:        "Asu (Tanzania)",
		LanguageCodeEnumAst:          "Asturian",
		LanguageCodeEnumAstEs:        "Asturian (Spain)",
		LanguageCodeEnumAz:           "Azerbaijani",
		LanguageCodeEnumAzCyrl:       "Azerbaijani (Cyrillic)",
		LanguageCodeEnumAzCyrlAz:     "Azerbaijani (Cyrillic, Azerbaijan)",
		LanguageCodeEnumAzLatn:       "Azerbaijani (Latin)",
		LanguageCodeEnumAzLatnAz:     "Azerbaijani (Latin, Azerbaijan)",
		LanguageCodeEnumBas:          "Basaa",
		LanguageCodeEnumBasCm:        "Basaa (Cameroon)",
		LanguageCodeEnumBe:           "Belarusian",
		LanguageCodeEnumBeBy:         "Belarusian (Belarus)",
		LanguageCodeEnumBem:          "Bemba",
		LanguageCodeEnumBemZm:        "Bemba (Zambia)",
		LanguageCodeEnumBez:          "Bena",
		LanguageCodeEnumBezTz:        "Bena (Tanzania)",
		LanguageCodeEnumBg:           "Bulgarian",
		LanguageCodeEnumBgBg:         "Bulgarian (Bulgaria)",
		LanguageCodeEnumBm:           "Bambara",
		LanguageCodeEnumBmMl:         "Bambara (Mali)",
		LanguageCodeEnumBn:           "Bangla",
		LanguageCodeEnumBnBd:         "Bangla (Bangladesh)",
		LanguageCodeEnumBnIn:         "Bangla (India)",
		LanguageCodeEnumBo:           "Tibetan",
		LanguageCodeEnumBoCn:         "Tibetan (China)",
		LanguageCodeEnumBoIn:         "Tibetan (India)",
		LanguageCodeEnumBr:           "Breton",
		LanguageCodeEnumBrFr:         "Breton (France)",
		LanguageCodeEnumBrx:          "Bodo",
		LanguageCodeEnumBrxIn:        "Bodo (India)",
		LanguageCodeEnumBs:           "Bosnian",
		LanguageCodeEnumBsCyrl:       "Bosnian (Cyrillic)",
		LanguageCodeEnumBsCyrlBa:     "Bosnian (Cyrillic, Bosnia & Herzegovina)",
		LanguageCodeEnumBsLatn:       "Bosnian (Latin)",
		LanguageCodeEnumBsLatnBa:     "Bosnian (Latin, Bosnia & Herzegovina)",
		LanguageCodeEnumCa:           "Catalan",
		LanguageCodeEnumCaAd:         "Catalan (Andorra)",
		LanguageCodeEnumCaEs:         "Catalan (Spain)",
		LanguageCodeEnumCaEsValencia: "Catalan (Spain, Valencian)",
		LanguageCodeEnumCaFr:         "Catalan (France)",
		LanguageCodeEnumCaIt:         "Catalan (Italy)",
		LanguageCodeEnumCcp:          "Chakma",
		LanguageCodeEnumCcpBd:        "Chakma (Bangladesh)",
		LanguageCodeEnumCcpIn:        "Chakma (India)",
		LanguageCodeEnumCe:           "Chechen",
		LanguageCodeEnumCeRu:         "Chechen (Russia)",
		LanguageCodeEnumCeb:          "Cebuano",
		LanguageCodeEnumCebPh:        "Cebuano (Philippines)",
		LanguageCodeEnumCgg:          "Chiga",
		LanguageCodeEnumCggUg:        "Chiga (Uganda)",
		LanguageCodeEnumChr:          "Cherokee",
		LanguageCodeEnumChrUs:        "Cherokee (United States)",
		LanguageCodeEnumCkb:          "Central Kurdish",
		LanguageCodeEnumCkbIq:        "Central Kurdish (Iraq)",
		LanguageCodeEnumCkbIr:        "Central Kurdish (Iran)",
		LanguageCodeEnumCs:           "Czech",
		LanguageCodeEnumCsCz:         "Czech (Czechia)",
		LanguageCodeEnumCu:           "Church Slavic",
		LanguageCodeEnumCuRu:         "Church Slavic (Russia)",
		LanguageCodeEnumCy:           "Welsh",
		LanguageCodeEnumCyGb:         "Welsh (United Kingdom)",
		LanguageCodeEnumDa:           "Danish",
		LanguageCodeEnumDaDk:         "Danish (Denmark)",
		LanguageCodeEnumDaGl:         "Danish (Greenland)",
		LanguageCodeEnumDav:          "Taita",
		LanguageCodeEnumDavKe:        "Taita (Kenya)",
		LanguageCodeEnumDe:           "German",
		LanguageCodeEnumDeAt:         "German (Austria)",
		LanguageCodeEnumDeBe:         "German (Belgium)",
		LanguageCodeEnumDeCh:         "German (Switzerland)",
		LanguageCodeEnumDeDe:         "German (Germany)",
		LanguageCodeEnumDeIt:         "German (Italy)",
		LanguageCodeEnumDeLi:         "German (Liechtenstein)",
		LanguageCodeEnumDeLu:         "German (Luxembourg)",
		LanguageCodeEnumDje:          "Zarma",
		LanguageCodeEnumDjeNe:        "Zarma (Niger)",
		LanguageCodeEnumDsb:          "Lower Sorbian",
		LanguageCodeEnumDsbDe:        "Lower Sorbian (Germany)",
		LanguageCodeEnumDua:          "Duala",
		LanguageCodeEnumDuaCm:        "Duala (Cameroon)",
		LanguageCodeEnumDyo:          "Jola-Fonyi",
		LanguageCodeEnumDyoSn:        "Jola-Fonyi (Senegal)",
		LanguageCodeEnumDz:           "Dzongkha",
		LanguageCodeEnumDzBt:         "Dzongkha (Bhutan)",
		LanguageCodeEnumEbu:          "Embu",
		LanguageCodeEnumEbuKe:        "Embu (Kenya)",
		LanguageCodeEnumEe:           "Ewe",
		LanguageCodeEnumEeGh:         "Ewe (Ghana)",
		LanguageCodeEnumEeTg:         "Ewe (Togo)",
		LanguageCodeEnumEl:           "Greek",
		LanguageCodeEnumElCy:         "Greek (Cyprus)",
		LanguageCodeEnumElGr:         "Greek (Greece)",
		LanguageCodeEnumEn:           "English",
		LanguageCodeEnumEnAe:         "English (United Arab Emirates)",
		LanguageCodeEnumEnAg:         "English (Antigua & Barbuda)",
		LanguageCodeEnumEnAi:         "English (Anguilla)",
		LanguageCodeEnumEnAs:         "English (American Samoa)",
		LanguageCodeEnumEnAt:         "English (Austria)",
		LanguageCodeEnumEnAu:         "English (Australia)",
		LanguageCodeEnumEnBb:         "English (Barbados)",
		LanguageCodeEnumEnBe:         "English (Belgium)",
		LanguageCodeEnumEnBi:         "English (Burundi)",
		LanguageCodeEnumEnBm:         "English (Bermuda)",
		LanguageCodeEnumEnBs:         "English (Bahamas)",
		LanguageCodeEnumEnBw:         "English (Botswana)",
		LanguageCodeEnumEnBz:         "English (Belize)",
		LanguageCodeEnumEnCa:         "English (Canada)",
		LanguageCodeEnumEnCc:         "English (Cocos (Keeling) Islands)",
		LanguageCodeEnumEnCh:         "English (Switzerland)",
		LanguageCodeEnumEnCk:         "English (Cook Islands)",
		LanguageCodeEnumEnCm:         "English (Cameroon)",
		LanguageCodeEnumEnCx:         "English (Christmas Island)",
		LanguageCodeEnumEnCy:         "English (Cyprus)",
		LanguageCodeEnumEnDe:         "English (Germany)",
		LanguageCodeEnumEnDg:         "English (Diego Garcia)",
		LanguageCodeEnumEnDk:         "English (Denmark)",
		LanguageCodeEnumEnDm:         "English (Dominica)",
		LanguageCodeEnumEnEr:         "English (Eritrea)",
		LanguageCodeEnumEnFi:         "English (Finland)",
		LanguageCodeEnumEnFj:         "English (Fiji)",
		LanguageCodeEnumEnFk:         "English (Falkland Islands)",
		LanguageCodeEnumEnFm:         "English (Micronesia)",
		LanguageCodeEnumEnGb:         "English (United Kingdom)",
		LanguageCodeEnumEnGd:         "English (Grenada)",
		LanguageCodeEnumEnGg:         "English (Guernsey)",
		LanguageCodeEnumEnGh:         "English (Ghana)",
		LanguageCodeEnumEnGi:         "English (Gibraltar)",
		LanguageCodeEnumEnGm:         "English (Gambia)",
		LanguageCodeEnumEnGu:         "English (Guam)",
		LanguageCodeEnumEnGy:         "English (Guyana)",
		LanguageCodeEnumEnHk:         "English (Hong Kong SAR China)",
		LanguageCodeEnumEnIe:         "English (Ireland)",
		LanguageCodeEnumEnIl:         "English (Israel)",
		LanguageCodeEnumEnIm:         "English (Isle of Man)",
		LanguageCodeEnumEnIn:         "English (India)",
		LanguageCodeEnumEnIo:         "English (British Indian Ocean Territory)",
		LanguageCodeEnumEnJe:         "English (Jersey)",
		LanguageCodeEnumEnJm:         "English (Jamaica)",
		LanguageCodeEnumEnKe:         "English (Kenya)",
		LanguageCodeEnumEnKi:         "English (Kiribati)",
		LanguageCodeEnumEnKn:         "English (St. Kitts & Nevis)",
		LanguageCodeEnumEnKy:         "English (Cayman Islands)",
		LanguageCodeEnumEnLc:         "English (St. Lucia)",
		LanguageCodeEnumEnLr:         "English (Liberia)",
		LanguageCodeEnumEnLs:         "English (Lesotho)",
		LanguageCodeEnumEnMg:         "English (Madagascar)",
		LanguageCodeEnumEnMh:         "English (Marshall Islands)",
		LanguageCodeEnumEnMo:         "English (Macao SAR China)",
		LanguageCodeEnumEnMp:         "English (Northern Mariana Islands)",
		LanguageCodeEnumEnMs:         "English (Montserrat)",
		LanguageCodeEnumEnMt:         "English (Malta)",
		LanguageCodeEnumEnMu:         "English (Mauritius)",
		LanguageCodeEnumEnMw:         "English (Malawi)",
		LanguageCodeEnumEnMy:         "English (Malaysia)",
		LanguageCodeEnumEnNa:         "English (Namibia)",
		LanguageCodeEnumEnNf:         "English (Norfolk Island)",
		LanguageCodeEnumEnNg:         "English (Nigeria)",
		LanguageCodeEnumEnNl:         "English (Netherlands)",
		LanguageCodeEnumEnNr:         "English (Nauru)",
		LanguageCodeEnumEnNu:         "English (Niue)",
		LanguageCodeEnumEnNz:         "English (New Zealand)",
		LanguageCodeEnumEnPg:         "English (Papua New Guinea)",
		LanguageCodeEnumEnPh:         "English (Philippines)",
		LanguageCodeEnumEnPk:         "English (Pakistan)",
		LanguageCodeEnumEnPn:         "English (Pitcairn Islands)",
		LanguageCodeEnumEnPr:         "English (Puerto Rico)",
		LanguageCodeEnumEnPw:         "English (Palau)",
		LanguageCodeEnumEnRw:         "English (Rwanda)",
		LanguageCodeEnumEnSb:         "English (Solomon Islands)",
		LanguageCodeEnumEnSc:         "English (Seychelles)",
		LanguageCodeEnumEnSd:         "English (Sudan)",
		LanguageCodeEnumEnSe:         "English (Sweden)",
		LanguageCodeEnumEnSg:         "English (Singapore)",
		LanguageCodeEnumEnSh:         "English (St. Helena)",
		LanguageCodeEnumEnSi:         "English (Slovenia)",
		LanguageCodeEnumEnSl:         "English (Sierra Leone)",
		LanguageCodeEnumEnSs:         "English (South Sudan)",
		LanguageCodeEnumEnSx:         "English (Sint Maarten)",
		LanguageCodeEnumEnSz:         "English (Eswatini)",
		LanguageCodeEnumEnTc:         "English (Turks & Caicos Islands)",
		LanguageCodeEnumEnTk:         "English (Tokelau)",
		LanguageCodeEnumEnTo:         "English (Tonga)",
		LanguageCodeEnumEnTt:         "English (Trinidad & Tobago)",
		LanguageCodeEnumEnTv:         "English (Tuvalu)",
		LanguageCodeEnumEnTz:         "English (Tanzania)",
		LanguageCodeEnumEnUg:         "English (Uganda)",
		LanguageCodeEnumEnUm:         "English (U.S. Outlying Islands)",
		LanguageCodeEnumEnUs:         "English (United States)",
		LanguageCodeEnumEnVc:         "English (St. Vincent & Grenadines)",
		LanguageCodeEnumEnVg:         "English (British Virgin Islands)",
		LanguageCodeEnumEnVi:         "English (U.S. Virgin Islands)",
		LanguageCodeEnumEnVu:         "English (Vanuatu)",
		LanguageCodeEnumEnWs:         "English (Samoa)",
		LanguageCodeEnumEnZa:         "English (South Africa)",
		LanguageCodeEnumEnZm:         "English (Zambia)",
		LanguageCodeEnumEnZw:         "English (Zimbabwe)",
		LanguageCodeEnumEo:           "Esperanto",
		LanguageCodeEnumEs:           "Spanish",
		LanguageCodeEnumEsAr:         "Spanish (Argentina)",
		LanguageCodeEnumEsBo:         "Spanish (Bolivia)",
		LanguageCodeEnumEsBr:         "Spanish (Brazil)",
		LanguageCodeEnumEsBz:         "Spanish (Belize)",
		LanguageCodeEnumEsCl:         "Spanish (Chile)",
		LanguageCodeEnumEsCo:         "Spanish (Colombia)",
		LanguageCodeEnumEsCr:         "Spanish (Costa Rica)",
		LanguageCodeEnumEsCu:         "Spanish (Cuba)",
		LanguageCodeEnumEsDo:         "Spanish (Dominican Republic)",
		LanguageCodeEnumEsEa:         "Spanish (Ceuta & Melilla)",
		LanguageCodeEnumEsEc:         "Spanish (Ecuador)",
		LanguageCodeEnumEsEs:         "Spanish (Spain)",
		LanguageCodeEnumEsGq:         "Spanish (Equatorial Guinea)",
		LanguageCodeEnumEsGt:         "Spanish (Guatemala)",
		LanguageCodeEnumEsHn:         "Spanish (Honduras)",
		LanguageCodeEnumEsIc:         "Spanish (Canary Islands)",
		LanguageCodeEnumEsMx:         "Spanish (Mexico)",
		LanguageCodeEnumEsNi:         "Spanish (Nicaragua)",
		LanguageCodeEnumEsPa:         "Spanish (Panama)",
		LanguageCodeEnumEsPe:         "Spanish (Peru)",
		LanguageCodeEnumEsPh:         "Spanish (Philippines)",
		LanguageCodeEnumEsPr:         "Spanish (Puerto Rico)",
		LanguageCodeEnumEsPy:         "Spanish (Paraguay)",
		LanguageCodeEnumEsSv:         "Spanish (El Salvador)",
		LanguageCodeEnumEsUs:         "Spanish (United States)",
		LanguageCodeEnumEsUy:         "Spanish (Uruguay)",
		LanguageCodeEnumEsVe:         "Spanish (Venezuela)",
		LanguageCodeEnumEt:           "Estonian",
		LanguageCodeEnumEtEe:         "Estonian (Estonia)",
		LanguageCodeEnumEu:           "Basque",
		LanguageCodeEnumEuEs:         "Basque (Spain)",
		LanguageCodeEnumEwo:          "Ewondo",
		LanguageCodeEnumEwoCm:        "Ewondo (Cameroon)",
		LanguageCodeEnumFa:           "Persian",
		LanguageCodeEnumFaAf:         "Persian (Afghanistan)",
		LanguageCodeEnumFaIr:         "Persian (Iran)",
		LanguageCodeEnumFf:           "Fulah",
		LanguageCodeEnumFfAdlm:       "Fulah (Adlam)",
		LanguageCodeEnumFfAdlmBf:     "Fulah (Adlam, Burkina Faso)",
		LanguageCodeEnumFfAdlmCm:     "Fulah (Adlam, Cameroon)",
		LanguageCodeEnumFfAdlmGh:     "Fulah (Adlam, Ghana)",
		LanguageCodeEnumFfAdlmGm:     "Fulah (Adlam, Gambia)",
		LanguageCodeEnumFfAdlmGn:     "Fulah (Adlam, Guinea)",
		LanguageCodeEnumFfAdlmGw:     "Fulah (Adlam, Guinea-Bissau)",
		LanguageCodeEnumFfAdlmLr:     "Fulah (Adlam, Liberia)",
		LanguageCodeEnumFfAdlmMr:     "Fulah (Adlam, Mauritania)",
		LanguageCodeEnumFfAdlmNe:     "Fulah (Adlam, Niger)",
		LanguageCodeEnumFfAdlmNg:     "Fulah (Adlam, Nigeria)",
		LanguageCodeEnumFfAdlmSl:     "Fulah (Adlam, Sierra Leone)",
		LanguageCodeEnumFfAdlmSn:     "Fulah (Adlam, Senegal)",
		LanguageCodeEnumFfLatn:       "Fulah (Latin)",
		LanguageCodeEnumFfLatnBf:     "Fulah (Latin, Burkina Faso)",
		LanguageCodeEnumFfLatnCm:     "Fulah (Latin, Cameroon)",
		LanguageCodeEnumFfLatnGh:     "Fulah (Latin, Ghana)",
		LanguageCodeEnumFfLatnGm:     "Fulah (Latin, Gambia)",
		LanguageCodeEnumFfLatnGn:     "Fulah (Latin, Guinea)",
		LanguageCodeEnumFfLatnGw:     "Fulah (Latin, Guinea-Bissau)",
		LanguageCodeEnumFfLatnLr:     "Fulah (Latin, Liberia)",
		LanguageCodeEnumFfLatnMr:     "Fulah (Latin, Mauritania)",
		LanguageCodeEnumFfLatnNe:     "Fulah (Latin, Niger)",
		LanguageCodeEnumFfLatnNg:     "Fulah (Latin, Nigeria)",
		LanguageCodeEnumFfLatnSl:     "Fulah (Latin, Sierra Leone)",
		LanguageCodeEnumFfLatnSn:     "Fulah (Latin, Senegal)",
		LanguageCodeEnumFi:           "Finnish",
		LanguageCodeEnumFiFi:         "Finnish (Finland)",
		LanguageCodeEnumFil:          "Filipino",
		LanguageCodeEnumFilPh:        "Filipino (Philippines)",
		LanguageCodeEnumFo:           "Faroese",
		LanguageCodeEnumFoDk:         "Faroese (Denmark)",
		LanguageCodeEnumFoFo:         "Faroese (Faroe Islands)",
		LanguageCodeEnumFr:           "French",
		LanguageCodeEnumFrBe:         "French (Belgium)",
		LanguageCodeEnumFrBf:         "French (Burkina Faso)",
		LanguageCodeEnumFrBi:         "French (Burundi)",
		LanguageCodeEnumFrBj:         "French (Benin)",
		LanguageCodeEnumFrBl:         "French (St. Barth\u00e9lemy)",
		LanguageCodeEnumFrCa:         "French (Canada)",
		LanguageCodeEnumFrCd:         "French (Congo - Kinshasa)",
		LanguageCodeEnumFrCf:         "French (Central African Republic)",
		LanguageCodeEnumFrCg:         "French (Congo - Brazzaville)",
		LanguageCodeEnumFrCh:         "French (Switzerland)",
		LanguageCodeEnumFrCi:         "French (C\u00f4te d\u2019Ivoire)",
		LanguageCodeEnumFrCm:         "French (Cameroon)",
		LanguageCodeEnumFrDj:         "French (Djibouti)",
		LanguageCodeEnumFrDz:         "French (Algeria)",
		LanguageCodeEnumFrFr:         "French (France)",
		LanguageCodeEnumFrGa:         "French (Gabon)",
		LanguageCodeEnumFrGf:         "French (French Guiana)",
		LanguageCodeEnumFrGn:         "French (Guinea)",
		LanguageCodeEnumFrGp:         "French (Guadeloupe)",
		LanguageCodeEnumFrGq:         "French (Equatorial Guinea)",
		LanguageCodeEnumFrHt:         "French (Haiti)",
		LanguageCodeEnumFrKm:         "French (Comoros)",
		LanguageCodeEnumFrLu:         "French (Luxembourg)",
		LanguageCodeEnumFrMa:         "French (Morocco)",
		LanguageCodeEnumFrMc:         "French (Monaco)",
		LanguageCodeEnumFrMf:         "French (St. Martin)",
		LanguageCodeEnumFrMg:         "French (Madagascar)",
		LanguageCodeEnumFrMl:         "French (Mali)",
		LanguageCodeEnumFrMq:         "French (Martinique)",
		LanguageCodeEnumFrMr:         "French (Mauritania)",
		LanguageCodeEnumFrMu:         "French (Mauritius)",
		LanguageCodeEnumFrNc:         "French (New Caledonia)",
		LanguageCodeEnumFrNe:         "French (Niger)",
		LanguageCodeEnumFrPf:         "French (French Polynesia)",
		LanguageCodeEnumFrPm:         "French (St. Pierre & Miquelon)",
		LanguageCodeEnumFrRe:         "French (R\u00e9union)",
		LanguageCodeEnumFrRw:         "French (Rwanda)",
		LanguageCodeEnumFrSc:         "French (Seychelles)",
		LanguageCodeEnumFrSn:         "French (Senegal)",
		LanguageCodeEnumFrSy:         "French (Syria)",
		LanguageCodeEnumFrTd:         "French (Chad)",
		LanguageCodeEnumFrTg:         "French (Togo)",
		LanguageCodeEnumFrTn:         "French (Tunisia)",
		LanguageCodeEnumFrVu:         "French (Vanuatu)",
		LanguageCodeEnumFrWf:         "French (Wallis & Futuna)",
		LanguageCodeEnumFrYt:         "French (Mayotte)",
		LanguageCodeEnumFur:          "Friulian",
		LanguageCodeEnumFurIt:        "Friulian (Italy)",
		LanguageCodeEnumFy:           "Western Frisian",
		LanguageCodeEnumFyNl:         "Western Frisian (Netherlands)",
		LanguageCodeEnumGa:           "Irish",
		LanguageCodeEnumGaGb:         "Irish (United Kingdom)",
		LanguageCodeEnumGaIe:         "Irish (Ireland)",
		LanguageCodeEnumGd:           "Scottish Gaelic",
		LanguageCodeEnumGdGb:         "Scottish Gaelic (United Kingdom)",
		LanguageCodeEnumGl:           "Galician",
		LanguageCodeEnumGlEs:         "Galician (Spain)",
		LanguageCodeEnumGsw:          "Swiss German",
		LanguageCodeEnumGswCh:        "Swiss German (Switzerland)",
		LanguageCodeEnumGswFr:        "Swiss German (France)",
		LanguageCodeEnumGswLi:        "Swiss German (Liechtenstein)",
		LanguageCodeEnumGu:           "Gujarati",
		LanguageCodeEnumGuIn:         "Gujarati (India)",
		LanguageCodeEnumGuz:          "Gusii",
		LanguageCodeEnumGuzKe:        "Gusii (Kenya)",
		LanguageCodeEnumGv:           "Manx",
		LanguageCodeEnumGvIm:         "Manx (Isle of Man)",
		LanguageCodeEnumHa:           "Hausa",
		LanguageCodeEnumHaGh:         "Hausa (Ghana)",
		LanguageCodeEnumHaNe:         "Hausa (Niger)",
		LanguageCodeEnumHaNg:         "Hausa (Nigeria)",
		LanguageCodeEnumHaw:          "Hawaiian",
		LanguageCodeEnumHawUs:        "Hawaiian (United States)",
		LanguageCodeEnumHe:           "Hebrew",
		LanguageCodeEnumHeIl:         "Hebrew (Israel)",
		LanguageCodeEnumHi:           "Hindi",
		LanguageCodeEnumHiIn:         "Hindi (India)",
		LanguageCodeEnumHr:           "Croatian",
		LanguageCodeEnumHrBa:         "Croatian (Bosnia & Herzegovina)",
		LanguageCodeEnumHrHr:         "Croatian (Croatia)",
		LanguageCodeEnumHsb:          "Upper Sorbian",
		LanguageCodeEnumHsbDe:        "Upper Sorbian (Germany)",
		LanguageCodeEnumHu:           "Hungarian",
		LanguageCodeEnumHuHu:         "Hungarian (Hungary)",
		LanguageCodeEnumHy:           "Armenian",
		LanguageCodeEnumHyAm:         "Armenian (Armenia)",
		LanguageCodeEnumIa:           "Interlingua",
		LanguageCodeEnumId:           "Indonesian",
		LanguageCodeEnumIDID:         "Indonesian (Indonesia)",
		LanguageCodeEnumIg:           "Igbo",
		LanguageCodeEnumIgNg:         "Igbo (Nigeria)",
		LanguageCodeEnumIi:           "Sichuan Yi",
		LanguageCodeEnumIiCn:         "Sichuan Yi (China)",
		LanguageCodeEnumIs:           "Icelandic",
		LanguageCodeEnumIsIs:         "Icelandic (Iceland)",
		LanguageCodeEnumIt:           "Italian",
		LanguageCodeEnumItCh:         "Italian (Switzerland)",
		LanguageCodeEnumItIt:         "Italian (Italy)",
		LanguageCodeEnumItSm:         "Italian (San Marino)",
		LanguageCodeEnumItVa:         "Italian (Vatican City)",
		LanguageCodeEnumJa:           "Japanese",
		LanguageCodeEnumJaJp:         "Japanese (Japan)",
		LanguageCodeEnumJgo:          "Ngomba",
		LanguageCodeEnumJgoCm:        "Ngomba (Cameroon)",
		LanguageCodeEnumJmc:          "Machame",
		LanguageCodeEnumJmcTz:        "Machame (Tanzania)",
		LanguageCodeEnumJv:           "Javanese",
		LanguageCodeEnumJvID:         "Javanese (Indonesia)",
		LanguageCodeEnumKa:           "Georgian",
		LanguageCodeEnumKaGe:         "Georgian (Georgia)",
		LanguageCodeEnumKab:          "Kabyle",
		LanguageCodeEnumKabDz:        "Kabyle (Algeria)",
		LanguageCodeEnumKam:          "Kamba",
		LanguageCodeEnumKamKe:        "Kamba (Kenya)",
		LanguageCodeEnumKde:          "Makonde",
		LanguageCodeEnumKdeTz:        "Makonde (Tanzania)",
		LanguageCodeEnumKea:          "Kabuverdianu",
		LanguageCodeEnumKeaCv:        "Kabuverdianu (Cape Verde)",
		LanguageCodeEnumKhq:          "Koyra Chiini",
		LanguageCodeEnumKhqMl:        "Koyra Chiini (Mali)",
		LanguageCodeEnumKi:           "Kikuyu",
		LanguageCodeEnumKiKe:         "Kikuyu (Kenya)",
		LanguageCodeEnumKk:           "Kazakh",
		LanguageCodeEnumKkKz:         "Kazakh (Kazakhstan)",
		LanguageCodeEnumKkj:          "Kako",
		LanguageCodeEnumKkjCm:        "Kako (Cameroon)",
		LanguageCodeEnumKl:           "Kalaallisut",
		LanguageCodeEnumKlGl:         "Kalaallisut (Greenland)",
		LanguageCodeEnumKln:          "Kalenjin",
		LanguageCodeEnumKlnKe:        "Kalenjin (Kenya)",
		LanguageCodeEnumKm:           "Khmer",
		LanguageCodeEnumKmKh:         "Khmer (Cambodia)",
		LanguageCodeEnumKn:           "Kannada",
		LanguageCodeEnumKnIn:         "Kannada (India)",
		LanguageCodeEnumKo:           "Korean",
		LanguageCodeEnumKoKp:         "Korean (North Korea)",
		LanguageCodeEnumKoKr:         "Korean (South Korea)",
		LanguageCodeEnumKok:          "Konkani",
		LanguageCodeEnumKokIn:        "Konkani (India)",
		LanguageCodeEnumKs:           "Kashmiri",
		LanguageCodeEnumKsArab:       "Kashmiri (Arabic)",
		LanguageCodeEnumKsArabIn:     "Kashmiri (Arabic, India)",
		LanguageCodeEnumKsb:          "Shambala",
		LanguageCodeEnumKsbTz:        "Shambala (Tanzania)",
		LanguageCodeEnumKsf:          "Bafia",
		LanguageCodeEnumKsfCm:        "Bafia (Cameroon)",
		LanguageCodeEnumKsh:          "Colognian",
		LanguageCodeEnumKshDe:        "Colognian (Germany)",
		LanguageCodeEnumKu:           "Kurdish",
		LanguageCodeEnumKuTr:         "Kurdish (Turkey)",
		LanguageCodeEnumKw:           "Cornish",
		LanguageCodeEnumKwGb:         "Cornish (United Kingdom)",
		LanguageCodeEnumKy:           "Kyrgyz",
		LanguageCodeEnumKyKg:         "Kyrgyz (Kyrgyzstan)",
		LanguageCodeEnumLag:          "Langi",
		LanguageCodeEnumLagTz:        "Langi (Tanzania)",
		LanguageCodeEnumLb:           "Luxembourgish",
		LanguageCodeEnumLbLu:         "Luxembourgish (Luxembourg)",
		LanguageCodeEnumLg:           "Ganda",
		LanguageCodeEnumLgUg:         "Ganda (Uganda)",
		LanguageCodeEnumLkt:          "Lakota",
		LanguageCodeEnumLktUs:        "Lakota (United States)",
		LanguageCodeEnumLn:           "Lingala",
		LanguageCodeEnumLnAo:         "Lingala (Angola)",
		LanguageCodeEnumLnCd:         "Lingala (Congo - Kinshasa)",
		LanguageCodeEnumLnCf:         "Lingala (Central African Republic)",
		LanguageCodeEnumLnCg:         "Lingala (Congo - Brazzaville)",
		LanguageCodeEnumLo:           "Lao",
		LanguageCodeEnumLoLa:         "Lao (Laos)",
		LanguageCodeEnumLrc:          "Northern Luri",
		LanguageCodeEnumLrcIq:        "Northern Luri (Iraq)",
		LanguageCodeEnumLrcIr:        "Northern Luri (Iran)",
		LanguageCodeEnumLt:           "Lithuanian",
		LanguageCodeEnumLtLt:         "Lithuanian (Lithuania)",
		LanguageCodeEnumLu:           "Luba-Katanga",
		LanguageCodeEnumLuCd:         "Luba-Katanga (Congo - Kinshasa)",
		LanguageCodeEnumLuo:          "Luo",
		LanguageCodeEnumLuoKe:        "Luo (Kenya)",
		LanguageCodeEnumLuy:          "Luyia",
		LanguageCodeEnumLuyKe:        "Luyia (Kenya)",
		LanguageCodeEnumLv:           "Latvian",
		LanguageCodeEnumLvLv:         "Latvian (Latvia)",
		LanguageCodeEnumMai:          "Maithili",
		LanguageCodeEnumMaiIn:        "Maithili (India)",
		LanguageCodeEnumMas:          "Masai",
		LanguageCodeEnumMasKe:        "Masai (Kenya)",
		LanguageCodeEnumMasTz:        "Masai (Tanzania)",
		LanguageCodeEnumMer:          "Meru",
		LanguageCodeEnumMerKe:        "Meru (Kenya)",
		LanguageCodeEnumMfe:          "Morisyen",
		LanguageCodeEnumMfeMu:        "Morisyen (Mauritius)",
		LanguageCodeEnumMg:           "Malagasy",
		LanguageCodeEnumMgMg:         "Malagasy (Madagascar)",
		LanguageCodeEnumMgh:          "Makhuwa-Meetto",
		LanguageCodeEnumMghMz:        "Makhuwa-Meetto (Mozambique)",
		LanguageCodeEnumMgo:          "Meta\u02bc",
		LanguageCodeEnumMgoCm:        "Meta\u02bc (Cameroon)",
		LanguageCodeEnumMi:           "Maori",
		LanguageCodeEnumMiNz:         "Maori (New Zealand)",
		LanguageCodeEnumMk:           "Macedonian",
		LanguageCodeEnumMkMk:         "Macedonian (North Macedonia)",
		LanguageCodeEnumMl:           "Malayalam",
		LanguageCodeEnumMlIn:         "Malayalam (India)",
		LanguageCodeEnumMn:           "Mongolian",
		LanguageCodeEnumMnMn:         "Mongolian (Mongolia)",
		LanguageCodeEnumMni:          "Manipuri",
		LanguageCodeEnumMniBeng:      "Manipuri (Bangla)",
		LanguageCodeEnumMniBengIn:    "Manipuri (Bangla, India)",
		LanguageCodeEnumMr:           "Marathi",
		LanguageCodeEnumMrIn:         "Marathi (India)",
		LanguageCodeEnumMs:           "Malay",
		LanguageCodeEnumMsBn:         "Malay (Brunei)",
		LanguageCodeEnumMsID:         "Malay (Indonesia)",
		LanguageCodeEnumMsMy:         "Malay (Malaysia)",
		LanguageCodeEnumMsSg:         "Malay (Singapore)",
		LanguageCodeEnumMt:           "Maltese",
		LanguageCodeEnumMtMt:         "Maltese (Malta)",
		LanguageCodeEnumMua:          "Mundang",
		LanguageCodeEnumMuaCm:        "Mundang (Cameroon)",
		LanguageCodeEnumMy:           "Burmese",
		LanguageCodeEnumMyMm:         "Burmese (Myanmar (Burma))",
		LanguageCodeEnumMzn:          "Mazanderani",
		LanguageCodeEnumMznIr:        "Mazanderani (Iran)",
		LanguageCodeEnumNaq:          "Nama",
		LanguageCodeEnumNaqNa:        "Nama (Namibia)",
		LanguageCodeEnumNb:           "Norwegian Bokm\u00e5l",
		LanguageCodeEnumNbNo:         "Norwegian Bokm\u00e5l (Norway)",
		LanguageCodeEnumNbSj:         "Norwegian Bokm\u00e5l (Svalbard & Jan Mayen)",
		LanguageCodeEnumNd:           "North Ndebele",
		LanguageCodeEnumNdZw:         "North Ndebele (Zimbabwe)",
		LanguageCodeEnumNds:          "Low German",
		LanguageCodeEnumNdsDe:        "Low German (Germany)",
		LanguageCodeEnumNdsNl:        "Low German (Netherlands)",
		LanguageCodeEnumNe:           "Nepali",
		LanguageCodeEnumNeIn:         "Nepali (India)",
		LanguageCodeEnumNeNp:         "Nepali (Nepal)",
		LanguageCodeEnumNl:           "Dutch",
		LanguageCodeEnumNlAw:         "Dutch (Aruba)",
		LanguageCodeEnumNlBe:         "Dutch (Belgium)",
		LanguageCodeEnumNlBq:         "Dutch (Caribbean Netherlands)",
		LanguageCodeEnumNlCw:         "Dutch (Cura\u00e7ao)",
		LanguageCodeEnumNlNl:         "Dutch (Netherlands)",
		LanguageCodeEnumNlSr:         "Dutch (Suriname)",
		LanguageCodeEnumNlSx:         "Dutch (Sint Maarten)",
		LanguageCodeEnumNmg:          "Kwasio",
		LanguageCodeEnumNmgCm:        "Kwasio (Cameroon)",
		LanguageCodeEnumNn:           "Norwegian Nynorsk",
		LanguageCodeEnumNnNo:         "Norwegian Nynorsk (Norway)",
		LanguageCodeEnumNnh:          "Ngiemboon",
		LanguageCodeEnumNnhCm:        "Ngiemboon (Cameroon)",
		LanguageCodeEnumNus:          "Nuer",
		LanguageCodeEnumNusSs:        "Nuer (South Sudan)",
		LanguageCodeEnumNyn:          "Nyankole",
		LanguageCodeEnumNynUg:        "Nyankole (Uganda)",
		LanguageCodeEnumOm:           "Oromo",
		LanguageCodeEnumOmEt:         "Oromo (Ethiopia)",
		LanguageCodeEnumOmKe:         "Oromo (Kenya)",
		LanguageCodeEnumOr:           "Odia",
		LanguageCodeEnumOrIn:         "Odia (India)",
		LanguageCodeEnumOs:           "Ossetic",
		LanguageCodeEnumOsGe:         "Ossetic (Georgia)",
		LanguageCodeEnumOsRu:         "Ossetic (Russia)",
		LanguageCodeEnumPa:           "Punjabi",
		LanguageCodeEnumPaArab:       "Punjabi (Arabic)",
		LanguageCodeEnumPaArabPk:     "Punjabi (Arabic, Pakistan)",
		LanguageCodeEnumPaGuru:       "Punjabi (Gurmukhi)",
		LanguageCodeEnumPaGuruIn:     "Punjabi (Gurmukhi, India)",
		LanguageCodeEnumPcm:          "Nigerian Pidgin",
		LanguageCodeEnumPcmNg:        "Nigerian Pidgin (Nigeria)",
		LanguageCodeEnumPl:           "Polish",
		LanguageCodeEnumPlPl:         "Polish (Poland)",
		LanguageCodeEnumPrg:          "Prussian",
		LanguageCodeEnumPs:           "Pashto",
		LanguageCodeEnumPsAf:         "Pashto (Afghanistan)",
		LanguageCodeEnumPsPk:         "Pashto (Pakistan)",
		LanguageCodeEnumPt:           "Portuguese",
		LanguageCodeEnumPtAo:         "Portuguese (Angola)",
		LanguageCodeEnumPtBr:         "Portuguese (Brazil)",
		LanguageCodeEnumPtCh:         "Portuguese (Switzerland)",
		LanguageCodeEnumPtCv:         "Portuguese (Cape Verde)",
		LanguageCodeEnumPtGq:         "Portuguese (Equatorial Guinea)",
		LanguageCodeEnumPtGw:         "Portuguese (Guinea-Bissau)",
		LanguageCodeEnumPtLu:         "Portuguese (Luxembourg)",
		LanguageCodeEnumPtMo:         "Portuguese (Macao SAR China)",
		LanguageCodeEnumPtMz:         "Portuguese (Mozambique)",
		LanguageCodeEnumPtPt:         "Portuguese (Portugal)",
		LanguageCodeEnumPtSt:         "Portuguese (S\u00e3o Tom\u00e9 & Pr\u00edncipe)",
		LanguageCodeEnumPtTl:         "Portuguese (Timor-Leste)",
		LanguageCodeEnumQu:           "Quechua",
		LanguageCodeEnumQuBo:         "Quechua (Bolivia)",
		LanguageCodeEnumQuEc:         "Quechua (Ecuador)",
		LanguageCodeEnumQuPe:         "Quechua (Peru)",
		LanguageCodeEnumRm:           "Romansh",
		LanguageCodeEnumRmCh:         "Romansh (Switzerland)",
		LanguageCodeEnumRn:           "Rundi",
		LanguageCodeEnumRnBi:         "Rundi (Burundi)",
		LanguageCodeEnumRo:           "Romanian",
		LanguageCodeEnumRoMd:         "Romanian (Moldova)",
		LanguageCodeEnumRoRo:         "Romanian (Romania)",
		LanguageCodeEnumRof:          "Rombo",
		LanguageCodeEnumRofTz:        "Rombo (Tanzania)",
		LanguageCodeEnumRu:           "Russian",
		LanguageCodeEnumRuBy:         "Russian (Belarus)",
		LanguageCodeEnumRuKg:         "Russian (Kyrgyzstan)",
		LanguageCodeEnumRuKz:         "Russian (Kazakhstan)",
		LanguageCodeEnumRuMd:         "Russian (Moldova)",
		LanguageCodeEnumRuRu:         "Russian (Russia)",
		LanguageCodeEnumRuUa:         "Russian (Ukraine)",
		LanguageCodeEnumRw:           "Kinyarwanda",
		LanguageCodeEnumRwRw:         "Kinyarwanda (Rwanda)",
		LanguageCodeEnumRwk:          "Rwa",
		LanguageCodeEnumRwkTz:        "Rwa (Tanzania)",
		LanguageCodeEnumSah:          "Sakha",
		LanguageCodeEnumSahRu:        "Sakha (Russia)",
		LanguageCodeEnumSaq:          "Samburu",
		LanguageCodeEnumSaqKe:        "Samburu (Kenya)",
		LanguageCodeEnumSat:          "Santali",
		LanguageCodeEnumSatOlck:      "Santali (Ol Chiki)",
		LanguageCodeEnumSatOlckIn:    "Santali (Ol Chiki, India)",
		LanguageCodeEnumSbp:          "Sangu",
		LanguageCodeEnumSbpTz:        "Sangu (Tanzania)",
		LanguageCodeEnumSd:           "Sindhi",
		LanguageCodeEnumSdArab:       "Sindhi (Arabic)",
		LanguageCodeEnumSdArabPk:     "Sindhi (Arabic, Pakistan)",
		LanguageCodeEnumSdDeva:       "Sindhi (Devanagari)",
		LanguageCodeEnumSdDevaIn:     "Sindhi (Devanagari, India)",
		LanguageCodeEnumSe:           "Northern Sami",
		LanguageCodeEnumSeFi:         "Northern Sami (Finland)",
		LanguageCodeEnumSeNo:         "Northern Sami (Norway)",
		LanguageCodeEnumSeSe:         "Northern Sami (Sweden)",
		LanguageCodeEnumSeh:          "Sena",
		LanguageCodeEnumSehMz:        "Sena (Mozambique)",
		LanguageCodeEnumSes:          "Koyraboro Senni",
		LanguageCodeEnumSesMl:        "Koyraboro Senni (Mali)",
		LanguageCodeEnumSg:           "Sango",
		LanguageCodeEnumSgCf:         "Sango (Central African Republic)",
		LanguageCodeEnumShi:          "Tachelhit",
		LanguageCodeEnumShiLatn:      "Tachelhit (Latin)",
		LanguageCodeEnumShiLatnMa:    "Tachelhit (Latin, Morocco)",
		LanguageCodeEnumShiTfng:      "Tachelhit (Tifinagh)",
		LanguageCodeEnumShiTfngMa:    "Tachelhit (Tifinagh, Morocco)",
		LanguageCodeEnumSi:           "Sinhala",
		LanguageCodeEnumSiLk:         "Sinhala (Sri Lanka)",
		LanguageCodeEnumSk:           "Slovak",
		LanguageCodeEnumSkSk:         "Slovak (Slovakia)",
		LanguageCodeEnumSl:           "Slovenian",
		LanguageCodeEnumSlSi:         "Slovenian (Slovenia)",
		LanguageCodeEnumSmn:          "Inari Sami",
		LanguageCodeEnumSmnFi:        "Inari Sami (Finland)",
		LanguageCodeEnumSn:           "Shona",
		LanguageCodeEnumSnZw:         "Shona (Zimbabwe)",
		LanguageCodeEnumSo:           "Somali",
		LanguageCodeEnumSoDj:         "Somali (Djibouti)",
		LanguageCodeEnumSoEt:         "Somali (Ethiopia)",
		LanguageCodeEnumSoKe:         "Somali (Kenya)",
		LanguageCodeEnumSoSo:         "Somali (Somalia)",
		LanguageCodeEnumSq:           "Albanian",
		LanguageCodeEnumSqAl:         "Albanian (Albania)",
		LanguageCodeEnumSqMk:         "Albanian (North Macedonia)",
		LanguageCodeEnumSqXk:         "Albanian (Kosovo)",
		LanguageCodeEnumSr:           "Serbian",
		LanguageCodeEnumSrCyrl:       "Serbian (Cyrillic)",
		LanguageCodeEnumSrCyrlBa:     "Serbian (Cyrillic, Bosnia & Herzegovina)",
		LanguageCodeEnumSrCyrlMe:     "Serbian (Cyrillic, Montenegro)",
		LanguageCodeEnumSrCyrlRs:     "Serbian (Cyrillic, Serbia)",
		LanguageCodeEnumSrCyrlXk:     "Serbian (Cyrillic, Kosovo)",
		LanguageCodeEnumSrLatn:       "Serbian (Latin)",
		LanguageCodeEnumSrLatnBa:     "Serbian (Latin, Bosnia & Herzegovina)",
		LanguageCodeEnumSrLatnMe:     "Serbian (Latin, Montenegro)",
		LanguageCodeEnumSrLatnRs:     "Serbian (Latin, Serbia)",
		LanguageCodeEnumSrLatnXk:     "Serbian (Latin, Kosovo)",
		LanguageCodeEnumSu:           "Sundanese",
		LanguageCodeEnumSuLatn:       "Sundanese (Latin)",
		LanguageCodeEnumSuLatnID:     "Sundanese (Latin, Indonesia)",
		LanguageCodeEnumSv:           "Swedish",
		LanguageCodeEnumSvAx:         "Swedish (\u00c5land Islands)",
		LanguageCodeEnumSvFi:         "Swedish (Finland)",
		LanguageCodeEnumSvSe:         "Swedish (Sweden)",
		LanguageCodeEnumSw:           "Swahili",
		LanguageCodeEnumSwCd:         "Swahili (Congo - Kinshasa)",
		LanguageCodeEnumSwKe:         "Swahili (Kenya)",
		LanguageCodeEnumSwTz:         "Swahili (Tanzania)",
		LanguageCodeEnumSwUg:         "Swahili (Uganda)",
		LanguageCodeEnumTa:           "Tamil",
		LanguageCodeEnumTaIn:         "Tamil (India)",
		LanguageCodeEnumTaLk:         "Tamil (Sri Lanka)",
		LanguageCodeEnumTaMy:         "Tamil (Malaysia)",
		LanguageCodeEnumTaSg:         "Tamil (Singapore)",
		LanguageCodeEnumTe:           "Telugu",
		LanguageCodeEnumTeIn:         "Telugu (India)",
		LanguageCodeEnumTeo:          "Teso",
		LanguageCodeEnumTeoKe:        "Teso (Kenya)",
		LanguageCodeEnumTeoUg:        "Teso (Uganda)",
		LanguageCodeEnumTg:           "Tajik",
		LanguageCodeEnumTgTj:         "Tajik (Tajikistan)",
		LanguageCodeEnumTh:           "Thai",
		LanguageCodeEnumThTh:         "Thai (Thailand)",
		LanguageCodeEnumTi:           "Tigrinya",
		LanguageCodeEnumTiEr:         "Tigrinya (Eritrea)",
		LanguageCodeEnumTiEt:         "Tigrinya (Ethiopia)",
		LanguageCodeEnumTk:           "Turkmen",
		LanguageCodeEnumTkTm:         "Turkmen (Turkmenistan)",
		LanguageCodeEnumTo:           "Tongan",
		LanguageCodeEnumToTo:         "Tongan (Tonga)",
		LanguageCodeEnumTr:           "Turkish",
		LanguageCodeEnumTrCy:         "Turkish (Cyprus)",
		LanguageCodeEnumTrTr:         "Turkish (Turkey)",
		LanguageCodeEnumTt:           "Tatar",
		LanguageCodeEnumTtRu:         "Tatar (Russia)",
		LanguageCodeEnumTwq:          "Tasawaq",
		LanguageCodeEnumTwqNe:        "Tasawaq (Niger)",
		LanguageCodeEnumTzm:          "Central Atlas Tamazight",
		LanguageCodeEnumTzmMa:        "Central Atlas Tamazight (Morocco)",
		LanguageCodeEnumUg:           "Uyghur",
		LanguageCodeEnumUgCn:         "Uyghur (China)",
		LanguageCodeEnumUk:           "Ukrainian",
		LanguageCodeEnumUkUa:         "Ukrainian (Ukraine)",
		LanguageCodeEnumUr:           "Urdu",
		LanguageCodeEnumUrIn:         "Urdu (India)",
		LanguageCodeEnumUrPk:         "Urdu (Pakistan)",
		LanguageCodeEnumUz:           "Uzbek",
		LanguageCodeEnumUzArab:       "Uzbek (Arabic)",
		LanguageCodeEnumUzArabAf:     "Uzbek (Arabic, Afghanistan)",
		LanguageCodeEnumUzCyrl:       "Uzbek (Cyrillic)",
		LanguageCodeEnumUzCyrlUz:     "Uzbek (Cyrillic, Uzbekistan)",
		LanguageCodeEnumUzLatn:       "Uzbek (Latin)",
		LanguageCodeEnumUzLatnUz:     "Uzbek (Latin, Uzbekistan)",
		LanguageCodeEnumVai:          "Vai",
		LanguageCodeEnumVaiLatn:      "Vai (Latin)",
		LanguageCodeEnumVaiLatnLr:    "Vai (Latin, Liberia)",
		LanguageCodeEnumVaiVaii:      "Vai (Vai)",
		LanguageCodeEnumVaiVaiiLr:    "Vai (Vai, Liberia)",
		LanguageCodeEnumVi:           "Vietnamese",
		LanguageCodeEnumViVn:         "Vietnamese (Vietnam)",
		LanguageCodeEnumVo:           "Volap\u00fck",
		LanguageCodeEnumVun:          "Vunjo",
		LanguageCodeEnumVunTz:        "Vunjo (Tanzania)",
		LanguageCodeEnumWae:          "Walser",
		LanguageCodeEnumWaeCh:        "Walser (Switzerland)",
		LanguageCodeEnumWo:           "Wolof",
		LanguageCodeEnumWoSn:         "Wolof (Senegal)",
		LanguageCodeEnumXh:           "Xhosa",
		LanguageCodeEnumXhZa:         "Xhosa (South Africa)",
		LanguageCodeEnumXog:          "Soga",
		LanguageCodeEnumXogUg:        "Soga (Uganda)",
		LanguageCodeEnumYav:          "Yangben",
		LanguageCodeEnumYavCm:        "Yangben (Cameroon)",
		LanguageCodeEnumYi:           "Yiddish",
		LanguageCodeEnumYo:           "Yoruba",
		LanguageCodeEnumYoBj:         "Yoruba (Benin)",
		LanguageCodeEnumYoNg:         "Yoruba (Nigeria)",
		LanguageCodeEnumYue:          "Cantonese",
		LanguageCodeEnumYueHans:      "Cantonese (Simplified)",
		LanguageCodeEnumYueHansCn:    "Cantonese (Simplified, China)",
		LanguageCodeEnumYueHant:      "Cantonese (Traditional)",
		LanguageCodeEnumYueHantHk:    "Cantonese (Traditional, Hong Kong SAR China)",
		LanguageCodeEnumZgh:          "Standard Moroccan Tamazight",
		LanguageCodeEnumZghMa:        "Standard Moroccan Tamazight (Morocco)",
		LanguageCodeEnumZh:           "Chinese",
		LanguageCodeEnumZhHans:       "Chinese (Simplified)",
		LanguageCodeEnumZhHansCn:     "Chinese (Simplified, China)",
		LanguageCodeEnumZhHansHk:     "Chinese (Simplified, Hong Kong SAR China)",
		LanguageCodeEnumZhHansMo:     "Chinese (Simplified, Macao SAR China)",
		LanguageCodeEnumZhHansSg:     "Chinese (Simplified, Singapore)",
		LanguageCodeEnumZhHant:       "Chinese (Traditional)",
		LanguageCodeEnumZhHantHk:     "Chinese (Traditional, Hong Kong SAR China)",
		LanguageCodeEnumZhHantMo:     "Chinese (Traditional, Macao SAR China)",
		LanguageCodeEnumZhHantTw:     "Chinese (Traditional, Taiwan)",
		LanguageCodeEnumZu:           "Zulu",
		LanguageCodeEnumZuZa:         "Zulu (South Africa)",
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
	FirstName NamePart = "first" // "first"
	LastName  NamePart = "last"  // "last"
)

// TaxType is for unifying tax type object that comes from tax gateway
type TaxType struct {
	Code         string
	Descriptiton string
}
