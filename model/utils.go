package model

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nyaruka/phonenumbers"
	"github.com/sitename/sitename/modules/i18n"
	"github.com/sitename/sitename/modules/json"
	"github.com/sitename/sitename/modules/slog"
)

const (
	LOWERCASE_LETTERS = "abcdefghijklmnopqrstuvwxyz"
	UPPERCASE_LETTERS = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	SYMBOLS           = " !\"\\#$%&'()*+,-./:;<=>?@[]^_`|~"
	MB                = 1 << 20
	NUMBERS           = "0123456789"
)

var (
	encoding = base32.NewEncoding("ybndrfg8ejkmcpqxot1uwisza345h769")
)

type StringInterface map[string]interface{}
type StringArray []string

// Remove removes input from the array
func (sa StringArray) Remove(input string) StringArray {
	for index := range sa {
		if sa[index] == input {
			ret := make(StringArray, 0, len(sa)-1)
			ret = append(ret, sa[:index]...)
			return append(ret, sa[index+1:]...)
		}
	}
	return sa
}

var translateFunc i18n.TranslateFunc
var translateFuncOnce sync.Once

func AppErrorInit(t i18n.TranslateFunc) {
	translateFuncOnce.Do(func() {
		translateFunc = t
	})
}

// check if array of strings contains given input
func (sa StringArray) Contains(input string) bool {
	for index := range sa {
		if sa[index] == input {
			return true
		}
	}

	return false
}

// Equals checks if two arrays of strings have same length and contains the same elements at each index
func (sa StringArray) Equals(input StringArray) bool {
	if len(sa) != len(input) {
		return false
	}
	for index := range sa {
		if sa[index] != input[index] {
			return false
		}
	}

	return true
}

// StringInterfaceFromJson decodes input data in to a map with keys are strings and values are interface{}
func StringInterfaceFromJson(data io.Reader) map[string]interface{} {
	var objMap map[string]interface{}
	ModelFromJson(&objMap, data)
	return objMap
}

// GetMillis is a convenience method to get milliseconds since epoch.
func GetMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// GetMillisForTime is a convenience method to get milliseconds since epoch for provided Time.
func GetMillisForTime(thisTime time.Time) int64 {
	return thisTime.UnixNano() / int64(time.Millisecond)
}

// GetTimeForMillis is a convenience method to get time.Time for milliseconds since epoch.
func GetTimeForMillis(millis int64) time.Time {
	return time.Unix(0, millis*int64(time.Millisecond))
}

// GetStartOfDayMillis is a convenience method to get milliseconds since epoch for provided date's start of day
func GetStartOfDayMillis(thisTime time.Time, timeZoneOffset int) int64 {
	localSearchTimeZone := time.FixedZone("Local Search Time Zone", timeZoneOffset)
	resultTime := time.Date(thisTime.Year(), thisTime.Month(), thisTime.Day(), 0, 0, 0, 0, localSearchTimeZone)
	return GetMillisForTime(resultTime)
}

// GetEndOfDayMillis is a convenience method to get milliseconds since epoch for provided date's end of day
func GetEndOfDayMillis(thisTime time.Time, timeZoneOffset int) int64 {
	localSearchTimeZone := time.FixedZone("Local Search Time Zone", timeZoneOffset)
	resultTime := time.Date(thisTime.Year(), thisTime.Month(), thisTime.Day(), 23, 59, 59, 999999999, localSearchTimeZone)
	return GetMillisForTime(resultTime)
}

// AppError represents error caused while the system is operating
type AppError struct {
	Id            string `json:"id"`
	Message       string `json:"message"`               // Message to be display to the end user without debugging information
	DetailedError string `json:"detailed_error"`        // Internal error string to help the developer
	RequestId     string `json:"request_id,omitempty"`  // The RequestId that's also set in the header
	StatusCode    int    `json:"status_code,omitempty"` // The http status code
	Where         string `json:"-"`                     // The function where it happened in the form of Struct.Func
	IsOAuth       bool   `json:"is_oauth,omitempty"`    // Whether the error is OAuth specific
	params        map[string]interface{}
}

func (er *AppError) Error() string {
	return er.Where + ": " + er.Message + ", " + er.DetailedError
}

// NewAppError returns new app error with given parameters
func NewAppError(where, id string, params map[string]interface{}, details string, status int) *AppError {
	appErr := new(AppError)
	appErr.Id = id
	appErr.params = params
	appErr.Message = id
	appErr.Where = where
	appErr.DetailedError = details
	appErr.StatusCode = status
	appErr.IsOAuth = false

	return appErr
}

// common function for creating model.AppError type
//
// Example:
// 	modelInstance := &ModelType{
//		Id: "dsdsdre984jf8se990834",
// 		Name: "Hello World",
//	}
// 	outer := CreateAppError("model.collection.is_valid.%s.app_error", "collection_id=", "Collection.IsValid")
//	return outer("create_at", modelInstance.Id)
// NOTE: This is applied for errors with status code "http.StatusBadRequest (400)" only
func CreateAppErrorForModel(format, detailKey, where string) func(fieldName string, typeId *string) *AppError {
	var id, details string
	return func(fieldName string, typeId *string) *AppError {
		id = fmt.Sprintf(format, fieldName)
		if !strings.EqualFold(fieldName, "id") && typeId != nil {
			details = detailKey + *typeId
		}

		return NewAppError(where, id, nil, details, http.StatusBadRequest)
	}
}

// Encodes database models to json string format
func ModelToJson(model interface{}) string {
	bytes, err := json.JSON.Marshal(&model)
	if err != nil {
		slog.Error("Error encoding model", slog.Any("type", fmt.Sprintf("%T", model)))
		return ""
	}
	return string(bytes)
}

// Decodes json string into model.
//
// If decoding process encounter error, model will be nil
func ModelFromJson(model interface{}, data io.Reader) {
	err := json.JSON.NewDecoder(data).Decode(&model)
	if err != nil {
		slog.Error("Error decoding json", slog.String("error", err.Error()))
		model = nil
	}
}

func (a *AppError) ToJson() string {
	return ModelToJson(a)
}

func AppErrorFromJSon(data io.Reader) *AppError {
	str := ""
	bytes, err := ioutil.ReadAll(data)
	if err != nil {
		str = err.Error()
	} else {
		str = string(bytes)
	}

	decoder := json.JSON.NewDecoder(strings.NewReader(str))
	var er AppError
	err = decoder.Decode(&er)
	if err != nil {
		return NewAppError("AppErrorFromJson", "model.utils.decode_json.app_error", nil, "body: "+str, http.StatusInternalServerError)
	}
	return &er
}

// CopyStringMap make a new map, copy all the key-value pairs from original map into the new map
func CopyStringMap(originalMap map[string]string) map[string]string {
	copyMap := make(map[string]string)
	for k, v := range originalMap {
		copyMap[k] = v
	}
	return copyMap
}

func GetPreferredTimezone(timezone StringMap) string {
	if timezone["useAutomaticTimezone"] == "true" {
		return timezone["automaticTimezone"]
	}

	return timezone["manualTimezone"]
}

// IsValidID check if given value is a valid uuid or not
func IsValidId(value string) bool {
	_, err := uuid.Parse(value)
	return err == nil
}

// check if s is lower-cased
func IsLower(s string) bool {
	return strings.ToLower(s) == s
}

// check if given email is valid email
func IsValidEmail(email string) bool {
	if addr, err := mail.ParseAddress(email); err != nil {
		return false
	} else if addr.Name != "" {
		// mail.ParseAddress accepts input of the form "Billy Bob <billy@example.com>" which we don't allow
		return false
	}

	return true
}

// Copied from https://golang.org/src/net/dnsclient.go#L119
func IsDomainName(s string) bool {
	// See RFC 1035, RFC 3696.
	// Presentation format has dots before every label except the first, and the
	// terminal empty label is optional here because we assume fully-qualified
	// (absolute) input. We must therefore reserve space for the first and last
	// labels' length octets in wire format, where they are necessary and the
	// maximum total length is 255.
	// So our _effective_ maximum is 253, but 254 is not rejected if the last
	// character is a dot.
	l := len(s)
	if l == 0 || l > 254 || l == 254 && s[l-1] != '.' {
		return false
	}

	last := byte('.')
	ok := false // Ok once we've seen a letter.
	partlen := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		default:
			return false
		case 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_':
			ok = true
			partlen++
		case '0' <= c && c <= '9':
			// fine
			partlen++
		case c == '-':
			// Byte before dash cannot be dot.
			if last == '.' {
				return false
			}
			partlen++
		case c == '.':
			// Byte before dot cannot be dot, dash.
			if last == '.' || last == '-' {
				return false
			}
			if partlen > 63 || partlen == 0 {
				return false
			}
			partlen = 0
		}
		last = c
	}
	if last == '-' || partlen > 63 {
		return false
	}

	return ok
}

// SanitizeUnicode will remove undesirable Unicode characters from a string.
func SanitizeUnicode(s string) string {
	return strings.Map(filterBlocklist, s)
}

// filterBlocklist returns `r` if it is not in the blocklist, otherwise drop (-1).
// Blocklist is taken from https://www.w3.org/TR/unicode-xml/#Charlist
func filterBlocklist(r rune) rune {
	const drop = -1
	switch r {
	case '\u0340', '\u0341': // clones of grave and acute; deprecated in Unicode
		return drop
	case '\u17A3', '\u17D3': // obsolete characters for Khmer; deprecated in Unicode
		return drop
	case '\u2028', '\u2029': // line and paragraph separator
		return drop
	case '\u202A', '\u202B', '\u202C', '\u202D', '\u202E': // BIDI embedding controls
		return drop
	case '\u206A', '\u206B': // activate/inhibit symmetric swapping; deprecated in Unicode
		return drop
	case '\u206C', '\u206D': // activate/inhibit Arabic form shaping; deprecated in Unicode
		return drop
	case '\u206E', '\u206F': // activate/inhibit national digit shapes; deprecated in Unicode
		return drop
	case '\uFFF9', '\uFFFA', '\uFFFB': // interlinear annotation characters
		return drop
	case '\uFEFF': // byte order mark
		return drop
	case '\uFFFC': // object replacement character
		return drop
	}

	// Scoping for musical notation
	if r >= 0x0001D173 && r <= 0x0001D17A {
		return drop
	}

	// Language tag code points
	if r >= 0x000E0000 && r <= 0x000E007F {
		return drop
	}

	return r
}

// UniqueStrings returns a unique subset of the string slice provided.
func UniqueStrings(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

// IsValidAlphaNum checks if s contains only ASCII characters
func IsValidAlphaNum(s string) bool {
	validAlphaNum := regexp.MustCompile(`^[a-z0-9]+([a-z\-0-9]+|(__)?)[a-z0-9]+$`)

	return validAlphaNum.MatchString(s)
}

// IsAllNumbers checks is string s contains only ASCII digits
func IsAllNumbers(s string) bool {
	validNumbers := regexp.MustCompile("^[0-9]+$")
	return validNumbers.MatchString(s)
}

func IsValidAlphaNumHyphenUnderscore(s string, withFormat bool) bool {
	if withFormat {
		validAlphaNumHyphenUnderscore := regexp.MustCompile(`^[a-z0-9]+([a-z\-\_0-9]+|(__)?)[a-z0-9]+$`)
		return validAlphaNumHyphenUnderscore.MatchString(s)
	}

	validSimpleAlphaNumHyphenUnderscore := regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`)
	return validSimpleAlphaNumHyphenUnderscore.MatchString(s)
}

func AsStringBoolMap(list []string) map[string]bool {
	listMap := map[string]bool{}
	for _, p := range list {
		listMap[p] = true
	}
	return listMap
}

// NewRandomString returns a random string of the given length.
// The resulting entropy will be (5 * length) bits.
func NewRandomString(length int) string {
	data := make([]byte, 1+(length*5/8))
	rand.Read(data)
	return encoding.EncodeToString(data)[:length]
}

// NewId generate new uuid string value
func NewId() string {
	return uuid.NewString()
}

// MapToJson converts a map to a json string
func MapToJson(objmap map[string]string) string {
	return ModelToJson(objmap)
}

// MapBoolToJson converts a map to a json string
func MapBoolToJson(objmap map[string]bool) string {
	return ModelToJson(objmap)
}

// MapFromJson will decode the key/value pair map
func MapFromJson(data io.Reader) map[string]string {
	decoder := json.JSON.NewDecoder(data)

	var objmap map[string]string
	if err := decoder.Decode(&objmap); err != nil {
		return make(map[string]string)
	}
	return objmap
}

// Conserts Jsonify string aray
func ArrayToJson(objmap []string) string {
	return ModelToJson(objmap)
}

func ArrayFromJson(data io.Reader) []string {
	decoder := json.JSON.NewDecoder(data)

	var objmap []string
	if err := decoder.Decode(&objmap); err != nil {
		return make([]string, 0)
	}
	return objmap
}

func StringInterfaceToJson(objmap map[string]interface{}) string {
	return ModelToJson(objmap)
}

func Etag(parts ...interface{}) string {
	etag := CurrentVersion

	for _, part := range parts {
		etag += fmt.Sprintf(".%v", part)
	}

	return etag
}

// Check is rawURL is a valid URL or not
func IsValidHttpUrl(rawURL string) bool {
	if strings.Index(rawURL, "http://") != 0 && strings.Index(rawURL, "https://") != 0 {
		return false
	}

	if u, err := url.ParseRequestURI(rawURL); err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func IsSafeLink(link *string) bool {
	if link != nil {
		if IsValidHttpUrl(*link) {
			return true
		} else if strings.HasPrefix(*link, "/") {
			return true
		} else {
			return false
		}
	}

	return true
}

// Check is rawURL is valid websocket url or not
//
// Valid websocket URLs must be prefixed with ws:// or wss://
func IsValidWebsocketUrl(rawUrl string) bool {
	if strings.Index(rawUrl, "ws://") != 0 && strings.Index(rawUrl, "wss://") != 0 {
		return false
	}
	if _, err := url.ParseRequestURI(rawUrl); err != nil {
		return false
	}

	return true
}

// NormalizeEmail is borrowed from django's BaseUserManager class
func NormalizeEmail(email string) string {
	splitEmail := strings.Split(email, "@")
	if len(splitEmail) != 2 {
		return email
	}

	return splitEmail[0] + "@" + strings.ToLower(splitEmail[1])
}

func GetServerIpAddress(iface string) string {
	var addrs []net.Addr
	if iface == "" {
		var err error
		addrs, err = net.InterfaceAddrs()
		if err != nil {
			return ""
		}
	} else {
		interfaces, err := net.Interfaces()
		if err != nil {
			return ""
		}
		for _, i := range interfaces {
			if i.Name == iface {
				addrs, err = i.Addrs()
				if err != nil {
					return ""
				}
				break
			}
		}
	}

	for _, addr := range addrs {

		if ip, ok := addr.(*net.IPNet); ok && !ip.IP.IsLoopback() && !ip.IP.IsLinkLocalUnicast() && !ip.IP.IsLinkLocalMulticast() {
			if ip.IP.To4() != nil {
				return ip.IP.String()
			}
		}
	}

	return ""
}

// Translates AppError to user's locale
func (er *AppError) Translate(T i18n.TranslateFunc) {
	if T == nil {
		er.Message = er.Id
		return
	}

	if er.params == nil {
		er.Message = T(er.Id)
	} else {
		er.Message = T(er.Id, er.params)
	}
}

func (er *AppError) SystemMessage(T i18n.TranslateFunc) string {
	if er.params == nil {
		return T(er.Id)
	}
	return T(er.Id, er.params)
}

// IsValidPhoneNumber checks if number is valid
func IsValidPhoneNumber(phone, countryCode string) bool {
	_, err := phonenumbers.Parse(phone, countryCode)
	return err == nil
}

// checkif username is valid
func IsValidUsername(s string) bool {
	if len(s) < USER_NAME_MIN_LENGTH || len(s) > USER_NAME_MAX_LENGTH {
		return false
	}
	if !ValidUsernameChars.MatchString(s) {
		return false
	}
	_, found := RestrictedUsernames[s]

	return !found
}

func NormalizeUsername(username string) string {
	return strings.ToLower(username)
}

// makes sure uname does not violate system standard naming rules
func CleanUsername(uname string) string {
	s := NormalizeUsername(strings.Replace(uname, " ", "-", -1))
	for _, value := range ReservedName {
		if s == value {
			s = strings.Replace(s, value, "", -1)
		}
	}
	s = strings.TrimSpace(s)
	for _, c := range s {
		char := fmt.Sprintf("%c", c)
		if !ValidUsernameChars.MatchString(char) {
			s = strings.Replace(s, char, "-", -1)
		}
	}
	s = strings.Trim(s, "-")

	if !IsValidUsername(s) {
		slog.Info("generating new username")
		s = "a" + uuid.New().String()
	}

	return s
}
