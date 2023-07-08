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
	LanguageCodeEnumAf           LanguageCodeEnum = "AF"
	LanguageCodeEnumAfNa         LanguageCodeEnum = "AF_NA"
	LanguageCodeEnumId           LanguageCodeEnum = "ID"
	LanguageCodeEnumAfZa         LanguageCodeEnum = "AF_ZA"
	LanguageCodeEnumAgq          LanguageCodeEnum = "AGQ"
	LanguageCodeEnumAgqCm        LanguageCodeEnum = "AGQ_CM"
	LanguageCodeEnumAk           LanguageCodeEnum = "AK"
	LanguageCodeEnumAkGh         LanguageCodeEnum = "AK_GH"
	LanguageCodeEnumAm           LanguageCodeEnum = "AM"
	LanguageCodeEnumAmEt         LanguageCodeEnum = "AM_ET"
	LanguageCodeEnumAr           LanguageCodeEnum = "AR"
	LanguageCodeEnumArAe         LanguageCodeEnum = "AR_AE"
	LanguageCodeEnumArBh         LanguageCodeEnum = "AR_BH"
	LanguageCodeEnumArDj         LanguageCodeEnum = "AR_DJ"
	LanguageCodeEnumArDz         LanguageCodeEnum = "AR_DZ"
	LanguageCodeEnumArEg         LanguageCodeEnum = "AR_EG"
	LanguageCodeEnumArEh         LanguageCodeEnum = "AR_EH"
	LanguageCodeEnumArEr         LanguageCodeEnum = "AR_ER"
	LanguageCodeEnumArIl         LanguageCodeEnum = "AR_IL"
	LanguageCodeEnumArIq         LanguageCodeEnum = "AR_IQ"
	LanguageCodeEnumArJo         LanguageCodeEnum = "AR_JO"
	LanguageCodeEnumArKm         LanguageCodeEnum = "AR_KM"
	LanguageCodeEnumArKw         LanguageCodeEnum = "AR_KW"
	LanguageCodeEnumArLb         LanguageCodeEnum = "AR_LB"
	LanguageCodeEnumArLy         LanguageCodeEnum = "AR_LY"
	LanguageCodeEnumArMa         LanguageCodeEnum = "AR_MA"
	LanguageCodeEnumArMr         LanguageCodeEnum = "AR_MR"
	LanguageCodeEnumArOm         LanguageCodeEnum = "AR_OM"
	LanguageCodeEnumArPs         LanguageCodeEnum = "AR_PS"
	LanguageCodeEnumArQa         LanguageCodeEnum = "AR_QA"
	LanguageCodeEnumArSa         LanguageCodeEnum = "AR_SA"
	LanguageCodeEnumArSd         LanguageCodeEnum = "AR_SD"
	LanguageCodeEnumArSo         LanguageCodeEnum = "AR_SO"
	LanguageCodeEnumArSs         LanguageCodeEnum = "AR_SS"
	LanguageCodeEnumArSy         LanguageCodeEnum = "AR_SY"
	LanguageCodeEnumArTd         LanguageCodeEnum = "AR_TD"
	LanguageCodeEnumArTn         LanguageCodeEnum = "AR_TN"
	LanguageCodeEnumArYe         LanguageCodeEnum = "AR_YE"
	LanguageCodeEnumAs           LanguageCodeEnum = "AS"
	LanguageCodeEnumAsIn         LanguageCodeEnum = "AS_IN"
	LanguageCodeEnumAsa          LanguageCodeEnum = "ASA"
	LanguageCodeEnumAsaTz        LanguageCodeEnum = "ASA_TZ"
	LanguageCodeEnumAst          LanguageCodeEnum = "AST"
	LanguageCodeEnumAstEs        LanguageCodeEnum = "AST_ES"
	LanguageCodeEnumAz           LanguageCodeEnum = "AZ"
	LanguageCodeEnumAzCyrl       LanguageCodeEnum = "AZ_CYRL"
	LanguageCodeEnumAzCyrlAz     LanguageCodeEnum = "AZ_CYRL_AZ"
	LanguageCodeEnumAzLatn       LanguageCodeEnum = "AZ_LATN"
	LanguageCodeEnumAzLatnAz     LanguageCodeEnum = "AZ_LATN_AZ"
	LanguageCodeEnumBas          LanguageCodeEnum = "BAS"
	LanguageCodeEnumBasCm        LanguageCodeEnum = "BAS_CM"
	LanguageCodeEnumBe           LanguageCodeEnum = "BE"
	LanguageCodeEnumBeBy         LanguageCodeEnum = "BE_BY"
	LanguageCodeEnumBem          LanguageCodeEnum = "BEM"
	LanguageCodeEnumBemZm        LanguageCodeEnum = "BEM_ZM"
	LanguageCodeEnumBez          LanguageCodeEnum = "BEZ"
	LanguageCodeEnumBezTz        LanguageCodeEnum = "BEZ_TZ"
	LanguageCodeEnumBg           LanguageCodeEnum = "BG"
	LanguageCodeEnumBgBg         LanguageCodeEnum = "BG_BG"
	LanguageCodeEnumBm           LanguageCodeEnum = "BM"
	LanguageCodeEnumBmMl         LanguageCodeEnum = "BM_ML"
	LanguageCodeEnumBn           LanguageCodeEnum = "BN"
	LanguageCodeEnumBnBd         LanguageCodeEnum = "BN_BD"
	LanguageCodeEnumBnIn         LanguageCodeEnum = "BN_IN"
	LanguageCodeEnumBo           LanguageCodeEnum = "BO"
	LanguageCodeEnumBoCn         LanguageCodeEnum = "BO_CN"
	LanguageCodeEnumBoIn         LanguageCodeEnum = "BO_IN"
	LanguageCodeEnumBr           LanguageCodeEnum = "BR"
	LanguageCodeEnumBrFr         LanguageCodeEnum = "BR_FR"
	LanguageCodeEnumBrx          LanguageCodeEnum = "BRX"
	LanguageCodeEnumBrxIn        LanguageCodeEnum = "BRX_IN"
	LanguageCodeEnumBs           LanguageCodeEnum = "BS"
	LanguageCodeEnumBsCyrl       LanguageCodeEnum = "BS_CYRL"
	LanguageCodeEnumBsCyrlBa     LanguageCodeEnum = "BS_CYRL_BA"
	LanguageCodeEnumBsLatn       LanguageCodeEnum = "BS_LATN"
	LanguageCodeEnumBsLatnBa     LanguageCodeEnum = "BS_LATN_BA"
	LanguageCodeEnumCa           LanguageCodeEnum = "CA"
	LanguageCodeEnumCaAd         LanguageCodeEnum = "CA_AD"
	LanguageCodeEnumCaEs         LanguageCodeEnum = "CA_ES"
	LanguageCodeEnumCaEsValencia LanguageCodeEnum = "CA_ES_VALENCIA"
	LanguageCodeEnumCaFr         LanguageCodeEnum = "CA_FR"
	LanguageCodeEnumCaIt         LanguageCodeEnum = "CA_IT"
	LanguageCodeEnumCcp          LanguageCodeEnum = "CCP"
	LanguageCodeEnumCcpBd        LanguageCodeEnum = "CCP_BD"
	LanguageCodeEnumCcpIn        LanguageCodeEnum = "CCP_IN"
	LanguageCodeEnumCe           LanguageCodeEnum = "CE"
	LanguageCodeEnumCeRu         LanguageCodeEnum = "CE_RU"
	LanguageCodeEnumCeb          LanguageCodeEnum = "CEB"
	LanguageCodeEnumCebPh        LanguageCodeEnum = "CEB_PH"
	LanguageCodeEnumCgg          LanguageCodeEnum = "CGG"
	LanguageCodeEnumCggUg        LanguageCodeEnum = "CGG_UG"
	LanguageCodeEnumChr          LanguageCodeEnum = "CHR"
	LanguageCodeEnumChrUs        LanguageCodeEnum = "CHR_US"
	LanguageCodeEnumCkb          LanguageCodeEnum = "CKB"
	LanguageCodeEnumCkbIq        LanguageCodeEnum = "CKB_IQ"
	LanguageCodeEnumCkbIr        LanguageCodeEnum = "CKB_IR"
	LanguageCodeEnumCs           LanguageCodeEnum = "CS"
	LanguageCodeEnumCsCz         LanguageCodeEnum = "CS_CZ"
	LanguageCodeEnumCu           LanguageCodeEnum = "CU"
	LanguageCodeEnumCuRu         LanguageCodeEnum = "CU_RU"
	LanguageCodeEnumCy           LanguageCodeEnum = "CY"
	LanguageCodeEnumCyGb         LanguageCodeEnum = "CY_GB"
	LanguageCodeEnumDa           LanguageCodeEnum = "DA"
	LanguageCodeEnumDaDk         LanguageCodeEnum = "DA_DK"
	LanguageCodeEnumDaGl         LanguageCodeEnum = "DA_GL"
	LanguageCodeEnumDav          LanguageCodeEnum = "DAV"
	LanguageCodeEnumDavKe        LanguageCodeEnum = "DAV_KE"
	LanguageCodeEnumDe           LanguageCodeEnum = "DE"
	LanguageCodeEnumDeAt         LanguageCodeEnum = "DE_AT"
	LanguageCodeEnumDeBe         LanguageCodeEnum = "DE_BE"
	LanguageCodeEnumDeCh         LanguageCodeEnum = "DE_CH"
	LanguageCodeEnumDeDe         LanguageCodeEnum = "DE_DE"
	LanguageCodeEnumDeIt         LanguageCodeEnum = "DE_IT"
	LanguageCodeEnumDeLi         LanguageCodeEnum = "DE_LI"
	LanguageCodeEnumDeLu         LanguageCodeEnum = "DE_LU"
	LanguageCodeEnumDje          LanguageCodeEnum = "DJE"
	LanguageCodeEnumDjeNe        LanguageCodeEnum = "DJE_NE"
	LanguageCodeEnumDsb          LanguageCodeEnum = "DSB"
	LanguageCodeEnumDsbDe        LanguageCodeEnum = "DSB_DE"
	LanguageCodeEnumDua          LanguageCodeEnum = "DUA"
	LanguageCodeEnumDuaCm        LanguageCodeEnum = "DUA_CM"
	LanguageCodeEnumDyo          LanguageCodeEnum = "DYO"
	LanguageCodeEnumDyoSn        LanguageCodeEnum = "DYO_SN"
	LanguageCodeEnumDz           LanguageCodeEnum = "DZ"
	LanguageCodeEnumDzBt         LanguageCodeEnum = "DZ_BT"
	LanguageCodeEnumEbu          LanguageCodeEnum = "EBU"
	LanguageCodeEnumEbuKe        LanguageCodeEnum = "EBU_KE"
	LanguageCodeEnumEe           LanguageCodeEnum = "EE"
	LanguageCodeEnumEeGh         LanguageCodeEnum = "EE_GH"
	LanguageCodeEnumEeTg         LanguageCodeEnum = "EE_TG"
	LanguageCodeEnumEl           LanguageCodeEnum = "EL"
	LanguageCodeEnumElCy         LanguageCodeEnum = "EL_CY"
	LanguageCodeEnumElGr         LanguageCodeEnum = "EL_GR"
	LanguageCodeEnumEn           LanguageCodeEnum = "EN"
	LanguageCodeEnumEnAe         LanguageCodeEnum = "EN_AE"
	LanguageCodeEnumEnAg         LanguageCodeEnum = "EN_AG"
	LanguageCodeEnumEnAi         LanguageCodeEnum = "EN_AI"
	LanguageCodeEnumEnAs         LanguageCodeEnum = "EN_AS"
	LanguageCodeEnumEnAt         LanguageCodeEnum = "EN_AT"
	LanguageCodeEnumEnAu         LanguageCodeEnum = "EN_AU"
	LanguageCodeEnumEnBb         LanguageCodeEnum = "EN_BB"
	LanguageCodeEnumEnBe         LanguageCodeEnum = "EN_BE"
	LanguageCodeEnumEnBi         LanguageCodeEnum = "EN_BI"
	LanguageCodeEnumEnBm         LanguageCodeEnum = "EN_BM"
	LanguageCodeEnumEnBs         LanguageCodeEnum = "EN_BS"
	LanguageCodeEnumEnBw         LanguageCodeEnum = "EN_BW"
	LanguageCodeEnumEnBz         LanguageCodeEnum = "EN_BZ"
	LanguageCodeEnumEnCa         LanguageCodeEnum = "EN_CA"
	LanguageCodeEnumEnCc         LanguageCodeEnum = "EN_CC"
	LanguageCodeEnumEnCh         LanguageCodeEnum = "EN_CH"
	LanguageCodeEnumEnCk         LanguageCodeEnum = "EN_CK"
	LanguageCodeEnumEnCm         LanguageCodeEnum = "EN_CM"
	LanguageCodeEnumEnCx         LanguageCodeEnum = "EN_CX"
	LanguageCodeEnumEnCy         LanguageCodeEnum = "EN_CY"
	LanguageCodeEnumEnDe         LanguageCodeEnum = "EN_DE"
	LanguageCodeEnumEnDg         LanguageCodeEnum = "EN_DG"
	LanguageCodeEnumEnDk         LanguageCodeEnum = "EN_DK"
	LanguageCodeEnumEnDm         LanguageCodeEnum = "EN_DM"
	LanguageCodeEnumEnEr         LanguageCodeEnum = "EN_ER"
	LanguageCodeEnumEnFi         LanguageCodeEnum = "EN_FI"
	LanguageCodeEnumEnFj         LanguageCodeEnum = "EN_FJ"
	LanguageCodeEnumEnFk         LanguageCodeEnum = "EN_FK"
	LanguageCodeEnumEnFm         LanguageCodeEnum = "EN_FM"
	LanguageCodeEnumEnGb         LanguageCodeEnum = "EN_GB"
	LanguageCodeEnumEnGd         LanguageCodeEnum = "EN_GD"
	LanguageCodeEnumEnGg         LanguageCodeEnum = "EN_GG"
	LanguageCodeEnumEnGh         LanguageCodeEnum = "EN_GH"
	LanguageCodeEnumEnGi         LanguageCodeEnum = "EN_GI"
	LanguageCodeEnumEnGm         LanguageCodeEnum = "EN_GM"
	LanguageCodeEnumEnGu         LanguageCodeEnum = "EN_GU"
	LanguageCodeEnumEnGy         LanguageCodeEnum = "EN_GY"
	LanguageCodeEnumEnHk         LanguageCodeEnum = "EN_HK"
	LanguageCodeEnumEnIe         LanguageCodeEnum = "EN_IE"
	LanguageCodeEnumEnIl         LanguageCodeEnum = "EN_IL"
	LanguageCodeEnumEnIm         LanguageCodeEnum = "EN_IM"
	LanguageCodeEnumEnIn         LanguageCodeEnum = "EN_IN"
	LanguageCodeEnumEnIo         LanguageCodeEnum = "EN_IO"
	LanguageCodeEnumEnJe         LanguageCodeEnum = "EN_JE"
	LanguageCodeEnumEnJm         LanguageCodeEnum = "EN_JM"
	LanguageCodeEnumEnKe         LanguageCodeEnum = "EN_KE"
	LanguageCodeEnumEnKi         LanguageCodeEnum = "EN_KI"
	LanguageCodeEnumEnKn         LanguageCodeEnum = "EN_KN"
	LanguageCodeEnumEnKy         LanguageCodeEnum = "EN_KY"
	LanguageCodeEnumEnLc         LanguageCodeEnum = "EN_LC"
	LanguageCodeEnumEnLr         LanguageCodeEnum = "EN_LR"
	LanguageCodeEnumEnLs         LanguageCodeEnum = "EN_LS"
	LanguageCodeEnumEnMg         LanguageCodeEnum = "EN_MG"
	LanguageCodeEnumEnMh         LanguageCodeEnum = "EN_MH"
	LanguageCodeEnumEnMo         LanguageCodeEnum = "EN_MO"
	LanguageCodeEnumEnMp         LanguageCodeEnum = "EN_MP"
	LanguageCodeEnumEnMs         LanguageCodeEnum = "EN_MS"
	LanguageCodeEnumEnMt         LanguageCodeEnum = "EN_MT"
	LanguageCodeEnumEnMu         LanguageCodeEnum = "EN_MU"
	LanguageCodeEnumEnMw         LanguageCodeEnum = "EN_MW"
	LanguageCodeEnumEnMy         LanguageCodeEnum = "EN_MY"
	LanguageCodeEnumEnNa         LanguageCodeEnum = "EN_NA"
	LanguageCodeEnumEnNf         LanguageCodeEnum = "EN_NF"
	LanguageCodeEnumEnNg         LanguageCodeEnum = "EN_NG"
	LanguageCodeEnumEnNl         LanguageCodeEnum = "EN_NL"
	LanguageCodeEnumEnNr         LanguageCodeEnum = "EN_NR"
	LanguageCodeEnumEnNu         LanguageCodeEnum = "EN_NU"
	LanguageCodeEnumEnNz         LanguageCodeEnum = "EN_NZ"
	LanguageCodeEnumEnPg         LanguageCodeEnum = "EN_PG"
	LanguageCodeEnumEnPh         LanguageCodeEnum = "EN_PH"
	LanguageCodeEnumEnPk         LanguageCodeEnum = "EN_PK"
	LanguageCodeEnumEnPn         LanguageCodeEnum = "EN_PN"
	LanguageCodeEnumEnPr         LanguageCodeEnum = "EN_PR"
	LanguageCodeEnumEnPw         LanguageCodeEnum = "EN_PW"
	LanguageCodeEnumEnRw         LanguageCodeEnum = "EN_RW"
	LanguageCodeEnumEnSb         LanguageCodeEnum = "EN_SB"
	LanguageCodeEnumEnSc         LanguageCodeEnum = "EN_SC"
	LanguageCodeEnumEnSd         LanguageCodeEnum = "EN_SD"
	LanguageCodeEnumEnSe         LanguageCodeEnum = "EN_SE"
	LanguageCodeEnumEnSg         LanguageCodeEnum = "EN_SG"
	LanguageCodeEnumEnSh         LanguageCodeEnum = "EN_SH"
	LanguageCodeEnumEnSi         LanguageCodeEnum = "EN_SI"
	LanguageCodeEnumEnSl         LanguageCodeEnum = "EN_SL"
	LanguageCodeEnumEnSs         LanguageCodeEnum = "EN_SS"
	LanguageCodeEnumEnSx         LanguageCodeEnum = "EN_SX"
	LanguageCodeEnumEnSz         LanguageCodeEnum = "EN_SZ"
	LanguageCodeEnumEnTc         LanguageCodeEnum = "EN_TC"
	LanguageCodeEnumEnTk         LanguageCodeEnum = "EN_TK"
	LanguageCodeEnumEnTo         LanguageCodeEnum = "EN_TO"
	LanguageCodeEnumEnTt         LanguageCodeEnum = "EN_TT"
	LanguageCodeEnumEnTv         LanguageCodeEnum = "EN_TV"
	LanguageCodeEnumEnTz         LanguageCodeEnum = "EN_TZ"
	LanguageCodeEnumEnUg         LanguageCodeEnum = "EN_UG"
	LanguageCodeEnumEnUm         LanguageCodeEnum = "EN_UM"
	LanguageCodeEnumEnUs         LanguageCodeEnum = "EN_US"
	LanguageCodeEnumEnVc         LanguageCodeEnum = "EN_VC"
	LanguageCodeEnumEnVg         LanguageCodeEnum = "EN_VG"
	LanguageCodeEnumEnVi         LanguageCodeEnum = "EN_VI"
	LanguageCodeEnumEnVu         LanguageCodeEnum = "EN_VU"
	LanguageCodeEnumEnWs         LanguageCodeEnum = "EN_WS"
	LanguageCodeEnumEnZa         LanguageCodeEnum = "EN_ZA"
	LanguageCodeEnumEnZm         LanguageCodeEnum = "EN_ZM"
	LanguageCodeEnumEnZw         LanguageCodeEnum = "EN_ZW"
	LanguageCodeEnumEo           LanguageCodeEnum = "EO"
	LanguageCodeEnumEs           LanguageCodeEnum = "ES"
	LanguageCodeEnumEsAr         LanguageCodeEnum = "ES_AR"
	LanguageCodeEnumEsBo         LanguageCodeEnum = "ES_BO"
	LanguageCodeEnumEsBr         LanguageCodeEnum = "ES_BR"
	LanguageCodeEnumEsBz         LanguageCodeEnum = "ES_BZ"
	LanguageCodeEnumEsCl         LanguageCodeEnum = "ES_CL"
	LanguageCodeEnumEsCo         LanguageCodeEnum = "ES_CO"
	LanguageCodeEnumEsCr         LanguageCodeEnum = "ES_CR"
	LanguageCodeEnumEsCu         LanguageCodeEnum = "ES_CU"
	LanguageCodeEnumEsDo         LanguageCodeEnum = "ES_DO"
	LanguageCodeEnumEsEa         LanguageCodeEnum = "ES_EA"
	LanguageCodeEnumEsEc         LanguageCodeEnum = "ES_EC"
	LanguageCodeEnumEsEs         LanguageCodeEnum = "ES_ES"
	LanguageCodeEnumEsGq         LanguageCodeEnum = "ES_GQ"
	LanguageCodeEnumEsGt         LanguageCodeEnum = "ES_GT"
	LanguageCodeEnumEsHn         LanguageCodeEnum = "ES_HN"
	LanguageCodeEnumEsIc         LanguageCodeEnum = "ES_IC"
	LanguageCodeEnumEsMx         LanguageCodeEnum = "ES_MX"
	LanguageCodeEnumEsNi         LanguageCodeEnum = "ES_NI"
	LanguageCodeEnumEsPa         LanguageCodeEnum = "ES_PA"
	LanguageCodeEnumEsPe         LanguageCodeEnum = "ES_PE"
	LanguageCodeEnumEsPh         LanguageCodeEnum = "ES_PH"
	LanguageCodeEnumEsPr         LanguageCodeEnum = "ES_PR"
	LanguageCodeEnumEsPy         LanguageCodeEnum = "ES_PY"
	LanguageCodeEnumEsSv         LanguageCodeEnum = "ES_SV"
	LanguageCodeEnumEsUs         LanguageCodeEnum = "ES_US"
	LanguageCodeEnumEsUy         LanguageCodeEnum = "ES_UY"
	LanguageCodeEnumEsVe         LanguageCodeEnum = "ES_VE"
	LanguageCodeEnumEt           LanguageCodeEnum = "ET"
	LanguageCodeEnumEtEe         LanguageCodeEnum = "ET_EE"
	LanguageCodeEnumEu           LanguageCodeEnum = "EU"
	LanguageCodeEnumEuEs         LanguageCodeEnum = "EU_ES"
	LanguageCodeEnumEwo          LanguageCodeEnum = "EWO"
	LanguageCodeEnumEwoCm        LanguageCodeEnum = "EWO_CM"
	LanguageCodeEnumFa           LanguageCodeEnum = "FA"
	LanguageCodeEnumFaAf         LanguageCodeEnum = "FA_AF"
	LanguageCodeEnumFaIr         LanguageCodeEnum = "FA_IR"
	LanguageCodeEnumFf           LanguageCodeEnum = "FF"
	LanguageCodeEnumFfAdlm       LanguageCodeEnum = "FF_ADLM"
	LanguageCodeEnumFfAdlmBf     LanguageCodeEnum = "FF_ADLM_BF"
	LanguageCodeEnumFfAdlmCm     LanguageCodeEnum = "FF_ADLM_CM"
	LanguageCodeEnumFfAdlmGh     LanguageCodeEnum = "FF_ADLM_GH"
	LanguageCodeEnumFfAdlmGm     LanguageCodeEnum = "FF_ADLM_GM"
	LanguageCodeEnumFfAdlmGn     LanguageCodeEnum = "FF_ADLM_GN"
	LanguageCodeEnumFfAdlmGw     LanguageCodeEnum = "FF_ADLM_GW"
	LanguageCodeEnumFfAdlmLr     LanguageCodeEnum = "FF_ADLM_LR"
	LanguageCodeEnumFfAdlmMr     LanguageCodeEnum = "FF_ADLM_MR"
	LanguageCodeEnumFfAdlmNe     LanguageCodeEnum = "FF_ADLM_NE"
	LanguageCodeEnumFfAdlmNg     LanguageCodeEnum = "FF_ADLM_NG"
	LanguageCodeEnumFfAdlmSl     LanguageCodeEnum = "FF_ADLM_SL"
	LanguageCodeEnumFfAdlmSn     LanguageCodeEnum = "FF_ADLM_SN"
	LanguageCodeEnumFfLatn       LanguageCodeEnum = "FF_LATN"
	LanguageCodeEnumFfLatnBf     LanguageCodeEnum = "FF_LATN_BF"
	LanguageCodeEnumFfLatnCm     LanguageCodeEnum = "FF_LATN_CM"
	LanguageCodeEnumFfLatnGh     LanguageCodeEnum = "FF_LATN_GH"
	LanguageCodeEnumFfLatnGm     LanguageCodeEnum = "FF_LATN_GM"
	LanguageCodeEnumFfLatnGn     LanguageCodeEnum = "FF_LATN_GN"
	LanguageCodeEnumFfLatnGw     LanguageCodeEnum = "FF_LATN_GW"
	LanguageCodeEnumFfLatnLr     LanguageCodeEnum = "FF_LATN_LR"
	LanguageCodeEnumFfLatnMr     LanguageCodeEnum = "FF_LATN_MR"
	LanguageCodeEnumFfLatnNe     LanguageCodeEnum = "FF_LATN_NE"
	LanguageCodeEnumFfLatnNg     LanguageCodeEnum = "FF_LATN_NG"
	LanguageCodeEnumFfLatnSl     LanguageCodeEnum = "FF_LATN_SL"
	LanguageCodeEnumFfLatnSn     LanguageCodeEnum = "FF_LATN_SN"
	LanguageCodeEnumFi           LanguageCodeEnum = "FI"
	LanguageCodeEnumFiFi         LanguageCodeEnum = "FI_FI"
	LanguageCodeEnumFil          LanguageCodeEnum = "FIL"
	LanguageCodeEnumFilPh        LanguageCodeEnum = "FIL_PH"
	LanguageCodeEnumFo           LanguageCodeEnum = "FO"
	LanguageCodeEnumFoDk         LanguageCodeEnum = "FO_DK"
	LanguageCodeEnumFoFo         LanguageCodeEnum = "FO_FO"
	LanguageCodeEnumFr           LanguageCodeEnum = "FR"
	LanguageCodeEnumFrBe         LanguageCodeEnum = "FR_BE"
	LanguageCodeEnumFrBf         LanguageCodeEnum = "FR_BF"
	LanguageCodeEnumFrBi         LanguageCodeEnum = "FR_BI"
	LanguageCodeEnumFrBj         LanguageCodeEnum = "FR_BJ"
	LanguageCodeEnumFrBl         LanguageCodeEnum = "FR_BL"
	LanguageCodeEnumFrCa         LanguageCodeEnum = "FR_CA"
	LanguageCodeEnumFrCd         LanguageCodeEnum = "FR_CD"
	LanguageCodeEnumFrCf         LanguageCodeEnum = "FR_CF"
	LanguageCodeEnumFrCg         LanguageCodeEnum = "FR_CG"
	LanguageCodeEnumFrCh         LanguageCodeEnum = "FR_CH"
	LanguageCodeEnumFrCi         LanguageCodeEnum = "FR_CI"
	LanguageCodeEnumFrCm         LanguageCodeEnum = "FR_CM"
	LanguageCodeEnumFrDj         LanguageCodeEnum = "FR_DJ"
	LanguageCodeEnumFrDz         LanguageCodeEnum = "FR_DZ"
	LanguageCodeEnumFrFr         LanguageCodeEnum = "FR_FR"
	LanguageCodeEnumFrGa         LanguageCodeEnum = "FR_GA"
	LanguageCodeEnumFrGf         LanguageCodeEnum = "FR_GF"
	LanguageCodeEnumFrGn         LanguageCodeEnum = "FR_GN"
	LanguageCodeEnumFrGp         LanguageCodeEnum = "FR_GP"
	LanguageCodeEnumFrGq         LanguageCodeEnum = "FR_GQ"
	LanguageCodeEnumFrHt         LanguageCodeEnum = "FR_HT"
	LanguageCodeEnumFrKm         LanguageCodeEnum = "FR_KM"
	LanguageCodeEnumFrLu         LanguageCodeEnum = "FR_LU"
	LanguageCodeEnumFrMa         LanguageCodeEnum = "FR_MA"
	LanguageCodeEnumFrMc         LanguageCodeEnum = "FR_MC"
	LanguageCodeEnumFrMf         LanguageCodeEnum = "FR_MF"
	LanguageCodeEnumFrMg         LanguageCodeEnum = "FR_MG"
	LanguageCodeEnumFrMl         LanguageCodeEnum = "FR_ML"
	LanguageCodeEnumFrMq         LanguageCodeEnum = "FR_MQ"
	LanguageCodeEnumFrMr         LanguageCodeEnum = "FR_MR"
	LanguageCodeEnumFrMu         LanguageCodeEnum = "FR_MU"
	LanguageCodeEnumFrNc         LanguageCodeEnum = "FR_NC"
	LanguageCodeEnumFrNe         LanguageCodeEnum = "FR_NE"
	LanguageCodeEnumFrPf         LanguageCodeEnum = "FR_PF"
	LanguageCodeEnumFrPm         LanguageCodeEnum = "FR_PM"
	LanguageCodeEnumFrRe         LanguageCodeEnum = "FR_RE"
	LanguageCodeEnumFrRw         LanguageCodeEnum = "FR_RW"
	LanguageCodeEnumFrSc         LanguageCodeEnum = "FR_SC"
	LanguageCodeEnumFrSn         LanguageCodeEnum = "FR_SN"
	LanguageCodeEnumFrSy         LanguageCodeEnum = "FR_SY"
	LanguageCodeEnumFrTd         LanguageCodeEnum = "FR_TD"
	LanguageCodeEnumFrTg         LanguageCodeEnum = "FR_TG"
	LanguageCodeEnumFrTn         LanguageCodeEnum = "FR_TN"
	LanguageCodeEnumFrVu         LanguageCodeEnum = "FR_VU"
	LanguageCodeEnumFrWf         LanguageCodeEnum = "FR_WF"
	LanguageCodeEnumFrYt         LanguageCodeEnum = "FR_YT"
	LanguageCodeEnumFur          LanguageCodeEnum = "FUR"
	LanguageCodeEnumFurIt        LanguageCodeEnum = "FUR_IT"
	LanguageCodeEnumFy           LanguageCodeEnum = "FY"
	LanguageCodeEnumFyNl         LanguageCodeEnum = "FY_NL"
	LanguageCodeEnumGa           LanguageCodeEnum = "GA"
	LanguageCodeEnumGaGb         LanguageCodeEnum = "GA_GB"
	LanguageCodeEnumGaIe         LanguageCodeEnum = "GA_IE"
	LanguageCodeEnumGd           LanguageCodeEnum = "GD"
	LanguageCodeEnumGdGb         LanguageCodeEnum = "GD_GB"
	LanguageCodeEnumGl           LanguageCodeEnum = "GL"
	LanguageCodeEnumGlEs         LanguageCodeEnum = "GL_ES"
	LanguageCodeEnumGsw          LanguageCodeEnum = "GSW"
	LanguageCodeEnumGswCh        LanguageCodeEnum = "GSW_CH"
	LanguageCodeEnumGswFr        LanguageCodeEnum = "GSW_FR"
	LanguageCodeEnumGswLi        LanguageCodeEnum = "GSW_LI"
	LanguageCodeEnumGu           LanguageCodeEnum = "GU"
	LanguageCodeEnumGuIn         LanguageCodeEnum = "GU_IN"
	LanguageCodeEnumGuz          LanguageCodeEnum = "GUZ"
	LanguageCodeEnumGuzKe        LanguageCodeEnum = "GUZ_KE"
	LanguageCodeEnumGv           LanguageCodeEnum = "GV"
	LanguageCodeEnumGvIm         LanguageCodeEnum = "GV_IM"
	LanguageCodeEnumHa           LanguageCodeEnum = "HA"
	LanguageCodeEnumHaGh         LanguageCodeEnum = "HA_GH"
	LanguageCodeEnumHaNe         LanguageCodeEnum = "HA_NE"
	LanguageCodeEnumHaNg         LanguageCodeEnum = "HA_NG"
	LanguageCodeEnumHaw          LanguageCodeEnum = "HAW"
	LanguageCodeEnumHawUs        LanguageCodeEnum = "HAW_US"
	LanguageCodeEnumHe           LanguageCodeEnum = "HE"
	LanguageCodeEnumHeIl         LanguageCodeEnum = "HE_IL"
	LanguageCodeEnumHi           LanguageCodeEnum = "HI"
	LanguageCodeEnumHiIn         LanguageCodeEnum = "HI_IN"
	LanguageCodeEnumHr           LanguageCodeEnum = "HR"
	LanguageCodeEnumHrBa         LanguageCodeEnum = "HR_BA"
	LanguageCodeEnumHrHr         LanguageCodeEnum = "HR_HR"
	LanguageCodeEnumHsb          LanguageCodeEnum = "HSB"
	LanguageCodeEnumHsbDe        LanguageCodeEnum = "HSB_DE"
	LanguageCodeEnumHu           LanguageCodeEnum = "HU"
	LanguageCodeEnumHuHu         LanguageCodeEnum = "HU_HU"
	LanguageCodeEnumHy           LanguageCodeEnum = "HY"
	LanguageCodeEnumHyAm         LanguageCodeEnum = "HY_AM"
	LanguageCodeEnumIa           LanguageCodeEnum = "IA"
	LanguageCodeEnumString       LanguageCodeEnum = "STRING"
	LanguageCodeEnumIDID         LanguageCodeEnum = "ID_ID"
	LanguageCodeEnumIg           LanguageCodeEnum = "IG"
	LanguageCodeEnumIgNg         LanguageCodeEnum = "IG_NG"
	LanguageCodeEnumIi           LanguageCodeEnum = "II"
	LanguageCodeEnumIiCn         LanguageCodeEnum = "II_CN"
	LanguageCodeEnumIs           LanguageCodeEnum = "IS"
	LanguageCodeEnumIsIs         LanguageCodeEnum = "IS_IS"
	LanguageCodeEnumIt           LanguageCodeEnum = "IT"
	LanguageCodeEnumItCh         LanguageCodeEnum = "IT_CH"
	LanguageCodeEnumItIt         LanguageCodeEnum = "IT_IT"
	LanguageCodeEnumItSm         LanguageCodeEnum = "IT_SM"
	LanguageCodeEnumItVa         LanguageCodeEnum = "IT_VA"
	LanguageCodeEnumJa           LanguageCodeEnum = "JA"
	LanguageCodeEnumJaJp         LanguageCodeEnum = "JA_JP"
	LanguageCodeEnumJgo          LanguageCodeEnum = "JGO"
	LanguageCodeEnumJgoCm        LanguageCodeEnum = "JGO_CM"
	LanguageCodeEnumJmc          LanguageCodeEnum = "JMC"
	LanguageCodeEnumJmcTz        LanguageCodeEnum = "JMC_TZ"
	LanguageCodeEnumJv           LanguageCodeEnum = "JV"
	LanguageCodeEnumJvID         LanguageCodeEnum = "JV_ID"
	LanguageCodeEnumKa           LanguageCodeEnum = "KA"
	LanguageCodeEnumKaGe         LanguageCodeEnum = "KA_GE"
	LanguageCodeEnumKab          LanguageCodeEnum = "KAB"
	LanguageCodeEnumKabDz        LanguageCodeEnum = "KAB_DZ"
	LanguageCodeEnumKam          LanguageCodeEnum = "KAM"
	LanguageCodeEnumKamKe        LanguageCodeEnum = "KAM_KE"
	LanguageCodeEnumKde          LanguageCodeEnum = "KDE"
	LanguageCodeEnumKdeTz        LanguageCodeEnum = "KDE_TZ"
	LanguageCodeEnumKea          LanguageCodeEnum = "KEA"
	LanguageCodeEnumKeaCv        LanguageCodeEnum = "KEA_CV"
	LanguageCodeEnumKhq          LanguageCodeEnum = "KHQ"
	LanguageCodeEnumKhqMl        LanguageCodeEnum = "KHQ_ML"
	LanguageCodeEnumKi           LanguageCodeEnum = "KI"
	LanguageCodeEnumKiKe         LanguageCodeEnum = "KI_KE"
	LanguageCodeEnumKk           LanguageCodeEnum = "KK"
	LanguageCodeEnumKkKz         LanguageCodeEnum = "KK_KZ"
	LanguageCodeEnumKkj          LanguageCodeEnum = "KKJ"
	LanguageCodeEnumKkjCm        LanguageCodeEnum = "KKJ_CM"
	LanguageCodeEnumKl           LanguageCodeEnum = "KL"
	LanguageCodeEnumKlGl         LanguageCodeEnum = "KL_GL"
	LanguageCodeEnumKln          LanguageCodeEnum = "KLN"
	LanguageCodeEnumKlnKe        LanguageCodeEnum = "KLN_KE"
	LanguageCodeEnumKm           LanguageCodeEnum = "KM"
	LanguageCodeEnumKmKh         LanguageCodeEnum = "KM_KH"
	LanguageCodeEnumKn           LanguageCodeEnum = "KN"
	LanguageCodeEnumKnIn         LanguageCodeEnum = "KN_IN"
	LanguageCodeEnumKo           LanguageCodeEnum = "KO"
	LanguageCodeEnumKoKp         LanguageCodeEnum = "KO_KP"
	LanguageCodeEnumKoKr         LanguageCodeEnum = "KO_KR"
	LanguageCodeEnumKok          LanguageCodeEnum = "KOK"
	LanguageCodeEnumKokIn        LanguageCodeEnum = "KOK_IN"
	LanguageCodeEnumKs           LanguageCodeEnum = "KS"
	LanguageCodeEnumKsArab       LanguageCodeEnum = "KS_ARAB"
	LanguageCodeEnumKsArabIn     LanguageCodeEnum = "KS_ARAB_IN"
	LanguageCodeEnumKsb          LanguageCodeEnum = "KSB"
	LanguageCodeEnumKsbTz        LanguageCodeEnum = "KSB_TZ"
	LanguageCodeEnumKsf          LanguageCodeEnum = "KSF"
	LanguageCodeEnumKsfCm        LanguageCodeEnum = "KSF_CM"
	LanguageCodeEnumKsh          LanguageCodeEnum = "KSH"
	LanguageCodeEnumKshDe        LanguageCodeEnum = "KSH_DE"
	LanguageCodeEnumKu           LanguageCodeEnum = "KU"
	LanguageCodeEnumKuTr         LanguageCodeEnum = "KU_TR"
	LanguageCodeEnumKw           LanguageCodeEnum = "KW"
	LanguageCodeEnumKwGb         LanguageCodeEnum = "KW_GB"
	LanguageCodeEnumKy           LanguageCodeEnum = "KY"
	LanguageCodeEnumKyKg         LanguageCodeEnum = "KY_KG"
	LanguageCodeEnumLag          LanguageCodeEnum = "LAG"
	LanguageCodeEnumLagTz        LanguageCodeEnum = "LAG_TZ"
	LanguageCodeEnumLb           LanguageCodeEnum = "LB"
	LanguageCodeEnumLbLu         LanguageCodeEnum = "LB_LU"
	LanguageCodeEnumLg           LanguageCodeEnum = "LG"
	LanguageCodeEnumLgUg         LanguageCodeEnum = "LG_UG"
	LanguageCodeEnumLkt          LanguageCodeEnum = "LKT"
	LanguageCodeEnumLktUs        LanguageCodeEnum = "LKT_US"
	LanguageCodeEnumLn           LanguageCodeEnum = "LN"
	LanguageCodeEnumLnAo         LanguageCodeEnum = "LN_AO"
	LanguageCodeEnumLnCd         LanguageCodeEnum = "LN_CD"
	LanguageCodeEnumLnCf         LanguageCodeEnum = "LN_CF"
	LanguageCodeEnumLnCg         LanguageCodeEnum = "LN_CG"
	LanguageCodeEnumLo           LanguageCodeEnum = "LO"
	LanguageCodeEnumLoLa         LanguageCodeEnum = "LO_LA"
	LanguageCodeEnumLrc          LanguageCodeEnum = "LRC"
	LanguageCodeEnumLrcIq        LanguageCodeEnum = "LRC_IQ"
	LanguageCodeEnumLrcIr        LanguageCodeEnum = "LRC_IR"
	LanguageCodeEnumLt           LanguageCodeEnum = "LT"
	LanguageCodeEnumLtLt         LanguageCodeEnum = "LT_LT"
	LanguageCodeEnumLu           LanguageCodeEnum = "LU"
	LanguageCodeEnumLuCd         LanguageCodeEnum = "LU_CD"
	LanguageCodeEnumLuo          LanguageCodeEnum = "LUO"
	LanguageCodeEnumLuoKe        LanguageCodeEnum = "LUO_KE"
	LanguageCodeEnumLuy          LanguageCodeEnum = "LUY"
	LanguageCodeEnumLuyKe        LanguageCodeEnum = "LUY_KE"
	LanguageCodeEnumLv           LanguageCodeEnum = "LV"
	LanguageCodeEnumLvLv         LanguageCodeEnum = "LV_LV"
	LanguageCodeEnumMai          LanguageCodeEnum = "MAI"
	LanguageCodeEnumMaiIn        LanguageCodeEnum = "MAI_IN"
	LanguageCodeEnumMas          LanguageCodeEnum = "MAS"
	LanguageCodeEnumMasKe        LanguageCodeEnum = "MAS_KE"
	LanguageCodeEnumMasTz        LanguageCodeEnum = "MAS_TZ"
	LanguageCodeEnumMer          LanguageCodeEnum = "MER"
	LanguageCodeEnumMerKe        LanguageCodeEnum = "MER_KE"
	LanguageCodeEnumMfe          LanguageCodeEnum = "MFE"
	LanguageCodeEnumMfeMu        LanguageCodeEnum = "MFE_MU"
	LanguageCodeEnumMg           LanguageCodeEnum = "MG"
	LanguageCodeEnumMgMg         LanguageCodeEnum = "MG_MG"
	LanguageCodeEnumMgh          LanguageCodeEnum = "MGH"
	LanguageCodeEnumMghMz        LanguageCodeEnum = "MGH_MZ"
	LanguageCodeEnumMgo          LanguageCodeEnum = "MGO"
	LanguageCodeEnumMgoCm        LanguageCodeEnum = "MGO_CM"
	LanguageCodeEnumMi           LanguageCodeEnum = "MI"
	LanguageCodeEnumMiNz         LanguageCodeEnum = "MI_NZ"
	LanguageCodeEnumMk           LanguageCodeEnum = "MK"
	LanguageCodeEnumMkMk         LanguageCodeEnum = "MK_MK"
	LanguageCodeEnumMl           LanguageCodeEnum = "ML"
	LanguageCodeEnumMlIn         LanguageCodeEnum = "ML_IN"
	LanguageCodeEnumMn           LanguageCodeEnum = "MN"
	LanguageCodeEnumMnMn         LanguageCodeEnum = "MN_MN"
	LanguageCodeEnumMni          LanguageCodeEnum = "MNI"
	LanguageCodeEnumMniBeng      LanguageCodeEnum = "MNI_BENG"
	LanguageCodeEnumMniBengIn    LanguageCodeEnum = "MNI_BENG_IN"
	LanguageCodeEnumMr           LanguageCodeEnum = "MR"
	LanguageCodeEnumMrIn         LanguageCodeEnum = "MR_IN"
	LanguageCodeEnumMs           LanguageCodeEnum = "MS"
	LanguageCodeEnumMsBn         LanguageCodeEnum = "MS_BN"
	LanguageCodeEnumMsID         LanguageCodeEnum = "MS_ID"
	LanguageCodeEnumMsMy         LanguageCodeEnum = "MS_MY"
	LanguageCodeEnumMsSg         LanguageCodeEnum = "MS_SG"
	LanguageCodeEnumMt           LanguageCodeEnum = "MT"
	LanguageCodeEnumMtMt         LanguageCodeEnum = "MT_MT"
	LanguageCodeEnumMua          LanguageCodeEnum = "MUA"
	LanguageCodeEnumMuaCm        LanguageCodeEnum = "MUA_CM"
	LanguageCodeEnumMy           LanguageCodeEnum = "MY"
	LanguageCodeEnumMyMm         LanguageCodeEnum = "MY_MM"
	LanguageCodeEnumMzn          LanguageCodeEnum = "MZN"
	LanguageCodeEnumMznIr        LanguageCodeEnum = "MZN_IR"
	LanguageCodeEnumNaq          LanguageCodeEnum = "NAQ"
	LanguageCodeEnumNaqNa        LanguageCodeEnum = "NAQ_NA"
	LanguageCodeEnumNb           LanguageCodeEnum = "NB"
	LanguageCodeEnumNbNo         LanguageCodeEnum = "NB_NO"
	LanguageCodeEnumNbSj         LanguageCodeEnum = "NB_SJ"
	LanguageCodeEnumNd           LanguageCodeEnum = "ND"
	LanguageCodeEnumNdZw         LanguageCodeEnum = "ND_ZW"
	LanguageCodeEnumNds          LanguageCodeEnum = "NDS"
	LanguageCodeEnumNdsDe        LanguageCodeEnum = "NDS_DE"
	LanguageCodeEnumNdsNl        LanguageCodeEnum = "NDS_NL"
	LanguageCodeEnumNe           LanguageCodeEnum = "NE"
	LanguageCodeEnumNeIn         LanguageCodeEnum = "NE_IN"
	LanguageCodeEnumNeNp         LanguageCodeEnum = "NE_NP"
	LanguageCodeEnumNl           LanguageCodeEnum = "NL"
	LanguageCodeEnumNlAw         LanguageCodeEnum = "NL_AW"
	LanguageCodeEnumNlBe         LanguageCodeEnum = "NL_BE"
	LanguageCodeEnumNlBq         LanguageCodeEnum = "NL_BQ"
	LanguageCodeEnumNlCw         LanguageCodeEnum = "NL_CW"
	LanguageCodeEnumNlNl         LanguageCodeEnum = "NL_NL"
	LanguageCodeEnumNlSr         LanguageCodeEnum = "NL_SR"
	LanguageCodeEnumNlSx         LanguageCodeEnum = "NL_SX"
	LanguageCodeEnumNmg          LanguageCodeEnum = "NMG"
	LanguageCodeEnumNmgCm        LanguageCodeEnum = "NMG_CM"
	LanguageCodeEnumNn           LanguageCodeEnum = "NN"
	LanguageCodeEnumNnNo         LanguageCodeEnum = "NN_NO"
	LanguageCodeEnumNnh          LanguageCodeEnum = "NNH"
	LanguageCodeEnumNnhCm        LanguageCodeEnum = "NNH_CM"
	LanguageCodeEnumNus          LanguageCodeEnum = "NUS"
	LanguageCodeEnumNusSs        LanguageCodeEnum = "NUS_SS"
	LanguageCodeEnumNyn          LanguageCodeEnum = "NYN"
	LanguageCodeEnumNynUg        LanguageCodeEnum = "NYN_UG"
	LanguageCodeEnumOm           LanguageCodeEnum = "OM"
	LanguageCodeEnumOmEt         LanguageCodeEnum = "OM_ET"
	LanguageCodeEnumOmKe         LanguageCodeEnum = "OM_KE"
	LanguageCodeEnumOr           LanguageCodeEnum = "OR"
	LanguageCodeEnumOrIn         LanguageCodeEnum = "OR_IN"
	LanguageCodeEnumOs           LanguageCodeEnum = "OS"
	LanguageCodeEnumOsGe         LanguageCodeEnum = "OS_GE"
	LanguageCodeEnumOsRu         LanguageCodeEnum = "OS_RU"
	LanguageCodeEnumPa           LanguageCodeEnum = "PA"
	LanguageCodeEnumPaArab       LanguageCodeEnum = "PA_ARAB"
	LanguageCodeEnumPaArabPk     LanguageCodeEnum = "PA_ARAB_PK"
	LanguageCodeEnumPaGuru       LanguageCodeEnum = "PA_GURU"
	LanguageCodeEnumPaGuruIn     LanguageCodeEnum = "PA_GURU_IN"
	LanguageCodeEnumPcm          LanguageCodeEnum = "PCM"
	LanguageCodeEnumPcmNg        LanguageCodeEnum = "PCM_NG"
	LanguageCodeEnumPl           LanguageCodeEnum = "PL"
	LanguageCodeEnumPlPl         LanguageCodeEnum = "PL_PL"
	LanguageCodeEnumPrg          LanguageCodeEnum = "PRG"
	LanguageCodeEnumPs           LanguageCodeEnum = "PS"
	LanguageCodeEnumPsAf         LanguageCodeEnum = "PS_AF"
	LanguageCodeEnumPsPk         LanguageCodeEnum = "PS_PK"
	LanguageCodeEnumPt           LanguageCodeEnum = "PT"
	LanguageCodeEnumPtAo         LanguageCodeEnum = "PT_AO"
	LanguageCodeEnumPtBr         LanguageCodeEnum = "PT_BR"
	LanguageCodeEnumPtCh         LanguageCodeEnum = "PT_CH"
	LanguageCodeEnumPtCv         LanguageCodeEnum = "PT_CV"
	LanguageCodeEnumPtGq         LanguageCodeEnum = "PT_GQ"
	LanguageCodeEnumPtGw         LanguageCodeEnum = "PT_GW"
	LanguageCodeEnumPtLu         LanguageCodeEnum = "PT_LU"
	LanguageCodeEnumPtMo         LanguageCodeEnum = "PT_MO"
	LanguageCodeEnumPtMz         LanguageCodeEnum = "PT_MZ"
	LanguageCodeEnumPtPt         LanguageCodeEnum = "PT_PT"
	LanguageCodeEnumPtSt         LanguageCodeEnum = "PT_ST"
	LanguageCodeEnumPtTl         LanguageCodeEnum = "PT_TL"
	LanguageCodeEnumQu           LanguageCodeEnum = "QU"
	LanguageCodeEnumQuBo         LanguageCodeEnum = "QU_BO"
	LanguageCodeEnumQuEc         LanguageCodeEnum = "QU_EC"
	LanguageCodeEnumQuPe         LanguageCodeEnum = "QU_PE"
	LanguageCodeEnumRm           LanguageCodeEnum = "RM"
	LanguageCodeEnumRmCh         LanguageCodeEnum = "RM_CH"
	LanguageCodeEnumRn           LanguageCodeEnum = "RN"
	LanguageCodeEnumRnBi         LanguageCodeEnum = "RN_BI"
	LanguageCodeEnumRo           LanguageCodeEnum = "RO"
	LanguageCodeEnumRoMd         LanguageCodeEnum = "RO_MD"
	LanguageCodeEnumRoRo         LanguageCodeEnum = "RO_RO"
	LanguageCodeEnumRof          LanguageCodeEnum = "ROF"
	LanguageCodeEnumRofTz        LanguageCodeEnum = "ROF_TZ"
	LanguageCodeEnumRu           LanguageCodeEnum = "RU"
	LanguageCodeEnumRuBy         LanguageCodeEnum = "RU_BY"
	LanguageCodeEnumRuKg         LanguageCodeEnum = "RU_KG"
	LanguageCodeEnumRuKz         LanguageCodeEnum = "RU_KZ"
	LanguageCodeEnumRuMd         LanguageCodeEnum = "RU_MD"
	LanguageCodeEnumRuRu         LanguageCodeEnum = "RU_RU"
	LanguageCodeEnumRuUa         LanguageCodeEnum = "RU_UA"
	LanguageCodeEnumRw           LanguageCodeEnum = "RW"
	LanguageCodeEnumRwRw         LanguageCodeEnum = "RW_RW"
	LanguageCodeEnumRwk          LanguageCodeEnum = "RWK"
	LanguageCodeEnumRwkTz        LanguageCodeEnum = "RWK_TZ"
	LanguageCodeEnumSah          LanguageCodeEnum = "SAH"
	LanguageCodeEnumSahRu        LanguageCodeEnum = "SAH_RU"
	LanguageCodeEnumSaq          LanguageCodeEnum = "SAQ"
	LanguageCodeEnumSaqKe        LanguageCodeEnum = "SAQ_KE"
	LanguageCodeEnumSat          LanguageCodeEnum = "SAT"
	LanguageCodeEnumSatOlck      LanguageCodeEnum = "SAT_OLCK"
	LanguageCodeEnumSatOlckIn    LanguageCodeEnum = "SAT_OLCK_IN"
	LanguageCodeEnumSbp          LanguageCodeEnum = "SBP"
	LanguageCodeEnumSbpTz        LanguageCodeEnum = "SBP_TZ"
	LanguageCodeEnumSd           LanguageCodeEnum = "SD"
	LanguageCodeEnumSdArab       LanguageCodeEnum = "SD_ARAB"
	LanguageCodeEnumSdArabPk     LanguageCodeEnum = "SD_ARAB_PK"
	LanguageCodeEnumSdDeva       LanguageCodeEnum = "SD_DEVA"
	LanguageCodeEnumSdDevaIn     LanguageCodeEnum = "SD_DEVA_IN"
	LanguageCodeEnumSe           LanguageCodeEnum = "SE"
	LanguageCodeEnumSeFi         LanguageCodeEnum = "SE_FI"
	LanguageCodeEnumSeNo         LanguageCodeEnum = "SE_NO"
	LanguageCodeEnumSeSe         LanguageCodeEnum = "SE_SE"
	LanguageCodeEnumSeh          LanguageCodeEnum = "SEH"
	LanguageCodeEnumSehMz        LanguageCodeEnum = "SEH_MZ"
	LanguageCodeEnumSes          LanguageCodeEnum = "SES"
	LanguageCodeEnumSesMl        LanguageCodeEnum = "SES_ML"
	LanguageCodeEnumSg           LanguageCodeEnum = "SG"
	LanguageCodeEnumSgCf         LanguageCodeEnum = "SG_CF"
	LanguageCodeEnumShi          LanguageCodeEnum = "SHI"
	LanguageCodeEnumShiLatn      LanguageCodeEnum = "SHI_LATN"
	LanguageCodeEnumShiLatnMa    LanguageCodeEnum = "SHI_LATN_MA"
	LanguageCodeEnumShiTfng      LanguageCodeEnum = "SHI_TFNG"
	LanguageCodeEnumShiTfngMa    LanguageCodeEnum = "SHI_TFNG_MA"
	LanguageCodeEnumSi           LanguageCodeEnum = "SI"
	LanguageCodeEnumSiLk         LanguageCodeEnum = "SI_LK"
	LanguageCodeEnumSk           LanguageCodeEnum = "SK"
	LanguageCodeEnumSkSk         LanguageCodeEnum = "SK_SK"
	LanguageCodeEnumSl           LanguageCodeEnum = "SL"
	LanguageCodeEnumSlSi         LanguageCodeEnum = "SL_SI"
	LanguageCodeEnumSmn          LanguageCodeEnum = "SMN"
	LanguageCodeEnumSmnFi        LanguageCodeEnum = "SMN_FI"
	LanguageCodeEnumSn           LanguageCodeEnum = "SN"
	LanguageCodeEnumSnZw         LanguageCodeEnum = "SN_ZW"
	LanguageCodeEnumSo           LanguageCodeEnum = "SO"
	LanguageCodeEnumSoDj         LanguageCodeEnum = "SO_DJ"
	LanguageCodeEnumSoEt         LanguageCodeEnum = "SO_ET"
	LanguageCodeEnumSoKe         LanguageCodeEnum = "SO_KE"
	LanguageCodeEnumSoSo         LanguageCodeEnum = "SO_SO"
	LanguageCodeEnumSq           LanguageCodeEnum = "SQ"
	LanguageCodeEnumSqAl         LanguageCodeEnum = "SQ_AL"
	LanguageCodeEnumSqMk         LanguageCodeEnum = "SQ_MK"
	LanguageCodeEnumSqXk         LanguageCodeEnum = "SQ_XK"
	LanguageCodeEnumSr           LanguageCodeEnum = "SR"
	LanguageCodeEnumSrCyrl       LanguageCodeEnum = "SR_CYRL"
	LanguageCodeEnumSrCyrlBa     LanguageCodeEnum = "SR_CYRL_BA"
	LanguageCodeEnumSrCyrlMe     LanguageCodeEnum = "SR_CYRL_ME"
	LanguageCodeEnumSrCyrlRs     LanguageCodeEnum = "SR_CYRL_RS"
	LanguageCodeEnumSrCyrlXk     LanguageCodeEnum = "SR_CYRL_XK"
	LanguageCodeEnumSrLatn       LanguageCodeEnum = "SR_LATN"
	LanguageCodeEnumSrLatnBa     LanguageCodeEnum = "SR_LATN_BA"
	LanguageCodeEnumSrLatnMe     LanguageCodeEnum = "SR_LATN_ME"
	LanguageCodeEnumSrLatnRs     LanguageCodeEnum = "SR_LATN_RS"
	LanguageCodeEnumSrLatnXk     LanguageCodeEnum = "SR_LATN_XK"
	LanguageCodeEnumSu           LanguageCodeEnum = "SU"
	LanguageCodeEnumSuLatn       LanguageCodeEnum = "SU_LATN"
	LanguageCodeEnumSuLatnID     LanguageCodeEnum = "SU_LATN_ID"
	LanguageCodeEnumSv           LanguageCodeEnum = "SV"
	LanguageCodeEnumSvAx         LanguageCodeEnum = "SV_AX"
	LanguageCodeEnumSvFi         LanguageCodeEnum = "SV_FI"
	LanguageCodeEnumSvSe         LanguageCodeEnum = "SV_SE"
	LanguageCodeEnumSw           LanguageCodeEnum = "SW"
	LanguageCodeEnumSwCd         LanguageCodeEnum = "SW_CD"
	LanguageCodeEnumSwKe         LanguageCodeEnum = "SW_KE"
	LanguageCodeEnumSwTz         LanguageCodeEnum = "SW_TZ"
	LanguageCodeEnumSwUg         LanguageCodeEnum = "SW_UG"
	LanguageCodeEnumTa           LanguageCodeEnum = "TA"
	LanguageCodeEnumTaIn         LanguageCodeEnum = "TA_IN"
	LanguageCodeEnumTaLk         LanguageCodeEnum = "TA_LK"
	LanguageCodeEnumTaMy         LanguageCodeEnum = "TA_MY"
	LanguageCodeEnumTaSg         LanguageCodeEnum = "TA_SG"
	LanguageCodeEnumTe           LanguageCodeEnum = "TE"
	LanguageCodeEnumTeIn         LanguageCodeEnum = "TE_IN"
	LanguageCodeEnumTeo          LanguageCodeEnum = "TEO"
	LanguageCodeEnumTeoKe        LanguageCodeEnum = "TEO_KE"
	LanguageCodeEnumTeoUg        LanguageCodeEnum = "TEO_UG"
	LanguageCodeEnumTg           LanguageCodeEnum = "TG"
	LanguageCodeEnumTgTj         LanguageCodeEnum = "TG_TJ"
	LanguageCodeEnumTh           LanguageCodeEnum = "TH"
	LanguageCodeEnumThTh         LanguageCodeEnum = "TH_TH"
	LanguageCodeEnumTi           LanguageCodeEnum = "TI"
	LanguageCodeEnumTiEr         LanguageCodeEnum = "TI_ER"
	LanguageCodeEnumTiEt         LanguageCodeEnum = "TI_ET"
	LanguageCodeEnumTk           LanguageCodeEnum = "TK"
	LanguageCodeEnumTkTm         LanguageCodeEnum = "TK_TM"
	LanguageCodeEnumTo           LanguageCodeEnum = "TO"
	LanguageCodeEnumToTo         LanguageCodeEnum = "TO_TO"
	LanguageCodeEnumTr           LanguageCodeEnum = "TR"
	LanguageCodeEnumTrCy         LanguageCodeEnum = "TR_CY"
	LanguageCodeEnumTrTr         LanguageCodeEnum = "TR_TR"
	LanguageCodeEnumTt           LanguageCodeEnum = "TT"
	LanguageCodeEnumTtRu         LanguageCodeEnum = "TT_RU"
	LanguageCodeEnumTwq          LanguageCodeEnum = "TWQ"
	LanguageCodeEnumTwqNe        LanguageCodeEnum = "TWQ_NE"
	LanguageCodeEnumTzm          LanguageCodeEnum = "TZM"
	LanguageCodeEnumTzmMa        LanguageCodeEnum = "TZM_MA"
	LanguageCodeEnumUg           LanguageCodeEnum = "UG"
	LanguageCodeEnumUgCn         LanguageCodeEnum = "UG_CN"
	LanguageCodeEnumUk           LanguageCodeEnum = "UK"
	LanguageCodeEnumUkUa         LanguageCodeEnum = "UK_UA"
	LanguageCodeEnumUr           LanguageCodeEnum = "UR"
	LanguageCodeEnumUrIn         LanguageCodeEnum = "UR_IN"
	LanguageCodeEnumUrPk         LanguageCodeEnum = "UR_PK"
	LanguageCodeEnumUz           LanguageCodeEnum = "UZ"
	LanguageCodeEnumUzArab       LanguageCodeEnum = "UZ_ARAB"
	LanguageCodeEnumUzArabAf     LanguageCodeEnum = "UZ_ARAB_AF"
	LanguageCodeEnumUzCyrl       LanguageCodeEnum = "UZ_CYRL"
	LanguageCodeEnumUzCyrlUz     LanguageCodeEnum = "UZ_CYRL_UZ"
	LanguageCodeEnumUzLatn       LanguageCodeEnum = "UZ_LATN"
	LanguageCodeEnumUzLatnUz     LanguageCodeEnum = "UZ_LATN_UZ"
	LanguageCodeEnumVai          LanguageCodeEnum = "VAI"
	LanguageCodeEnumVaiLatn      LanguageCodeEnum = "VAI_LATN"
	LanguageCodeEnumVaiLatnLr    LanguageCodeEnum = "VAI_LATN_LR"
	LanguageCodeEnumVaiVaii      LanguageCodeEnum = "VAI_VAII"
	LanguageCodeEnumVaiVaiiLr    LanguageCodeEnum = "VAI_VAII_LR"
	LanguageCodeEnumVi           LanguageCodeEnum = "VI"
	LanguageCodeEnumViVn         LanguageCodeEnum = "VI_VN"
	LanguageCodeEnumVo           LanguageCodeEnum = "VO"
	LanguageCodeEnumVun          LanguageCodeEnum = "VUN"
	LanguageCodeEnumVunTz        LanguageCodeEnum = "VUN_TZ"
	LanguageCodeEnumWae          LanguageCodeEnum = "WAE"
	LanguageCodeEnumWaeCh        LanguageCodeEnum = "WAE_CH"
	LanguageCodeEnumWo           LanguageCodeEnum = "WO"
	LanguageCodeEnumWoSn         LanguageCodeEnum = "WO_SN"
	LanguageCodeEnumXh           LanguageCodeEnum = "XH"
	LanguageCodeEnumXhZa         LanguageCodeEnum = "XH_ZA"
	LanguageCodeEnumXog          LanguageCodeEnum = "XOG"
	LanguageCodeEnumXogUg        LanguageCodeEnum = "XOG_UG"
	LanguageCodeEnumYav          LanguageCodeEnum = "YAV"
	LanguageCodeEnumYavCm        LanguageCodeEnum = "YAV_CM"
	LanguageCodeEnumYi           LanguageCodeEnum = "YI"
	LanguageCodeEnumYo           LanguageCodeEnum = "YO"
	LanguageCodeEnumYoBj         LanguageCodeEnum = "YO_BJ"
	LanguageCodeEnumYoNg         LanguageCodeEnum = "YO_NG"
	LanguageCodeEnumYue          LanguageCodeEnum = "YUE"
	LanguageCodeEnumYueHans      LanguageCodeEnum = "YUE_HANS"
	LanguageCodeEnumYueHansCn    LanguageCodeEnum = "YUE_HANS_CN"
	LanguageCodeEnumYueHant      LanguageCodeEnum = "YUE_HANT"
	LanguageCodeEnumYueHantHk    LanguageCodeEnum = "YUE_HANT_HK"
	LanguageCodeEnumZgh          LanguageCodeEnum = "ZGH"
	LanguageCodeEnumZghMa        LanguageCodeEnum = "ZGH_MA"
	LanguageCodeEnumZh           LanguageCodeEnum = "ZH"
	LanguageCodeEnumZhHans       LanguageCodeEnum = "ZH_HANS"
	LanguageCodeEnumZhHansCn     LanguageCodeEnum = "ZH_HANS_CN"
	LanguageCodeEnumZhHansHk     LanguageCodeEnum = "ZH_HANS_HK"
	LanguageCodeEnumZhHansMo     LanguageCodeEnum = "ZH_HANS_MO"
	LanguageCodeEnumZhHansSg     LanguageCodeEnum = "ZH_HANS_SG"
	LanguageCodeEnumZhHant       LanguageCodeEnum = "ZH_HANT"
	LanguageCodeEnumZhHantHk     LanguageCodeEnum = "ZH_HANT_HK"
	LanguageCodeEnumZhHantMo     LanguageCodeEnum = "ZH_HANT_MO"
	LanguageCodeEnumZhHantTw     LanguageCodeEnum = "ZH_HANT_TW"
	LanguageCodeEnumZu           LanguageCodeEnum = "ZU"
	LanguageCodeEnumZuZa         LanguageCodeEnum = "ZU_ZA"
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
