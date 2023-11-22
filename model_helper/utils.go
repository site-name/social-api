package model_helper

import (
	"cmp"
	"crypto/rand"
	"crypto/sha256"
	"database/sql/driver"
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/mail"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode"
	"unsafe"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/site-name/decimal"
	"github.com/sitename/sitename/modules/i18n"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util/fileutils"
	"github.com/uber/jaeger-client-go/utils"
)

const (
	LOWERCASE_LETTERS = "abcdefghijklmnopqrstuvwxyz"
	UPPERCASE_LETTERS = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	SYMBOLS           = " !\"\\#$%&'()*+,-./:;<=>?@[]^_`|~"
	MB                = 1 << 20
	NUMBERS           = "0123456789"
)

var encoding = base32.NewEncoding("ybndrfg8ejkmcpqxot1uwisza345h769")

const (
	MinIdLength  = 3
	MaxIdLength  = 190
	ValidIdRegex = `^[a-zA-Z0-9-_\.]+$`
)

// ValidId constrains the set of valid plugin identifiers:
//
//	^[a-zA-Z0-9-_\.]+
var validId *regexp.Regexp = regexp.MustCompile(ValidIdRegex)

// IsSamlFile checks if filename is a SAML file.
func IsSamlFile(saml *SamlSettings, filename string) bool {
	return filename == *saml.PublicCertificateFile || filename == *saml.PrivateKeyFile || filename == *saml.IdpCertificateFile
}

type StringInterface map[string]any
type StringInterfaces []StringInterface

func (s StringInterface) DeepCopy() StringInterface {
	if s == nil {
		return nil
	}

	res := StringInterface{}
	for key, value := range s {
		res[key] = value
	}

	return res
}

// Get trys finding and returns the value associated with given key.
func (s StringInterface) Get(key string, defaultValue ...any) any {
	if vl, ok := s[key]; ok {
		return vl
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return nil
}

// Pop trys finding and returns the value associated with given key.
// If the key does not exist:
//
// 1) Check if any default value given, returns the first value
//
// 2) returns nil
//
// Also delete the key-value from the map if found
func (s StringInterface) Pop(key string, defaultValue ...any) any {
	v := s.Get(key)
	delete(s, key)
	return v
}

func (si *StringInterface) Scan(value any) error {
	if value == nil {
		return nil
	}

	switch t := value.(type) {
	case []byte:
		return json.Unmarshal(t, si)
	case string:
		return json.Unmarshal([]byte(t), si)
	default:
		return fmt.Errorf("unsupported value type: %T", value)
	}
}

func (si StringInterface) Value() (driver.Value, error) {
	return json.Marshal(si)
}

func (si *StringInterfaces) Scan(value any) error {
	if value == nil {
		return nil
	}

	switch t := value.(type) {
	case []byte:
		return json.Unmarshal(t, si)
	case string:
		return json.Unmarshal([]byte(t), si)
	default:
		return fmt.Errorf("unsupported value type: %T", value)
	}
}

func (si StringInterfaces) Value() (driver.Value, error) {
	return json.Marshal(si)
}

// GetPointerOfValue takes a primitive value, returns pointer to it.
//
// E.g:
//
//	2 => &2 // ^.^
func GetPointerOfValue[T any](value T) *T {
	return &value
}

func CopyPointer[T cmp.Ordered | decimal.Decimal | time.Time | bool](value *T) *T {
	if value == nil {
		return nil
	}

	newPointer := *(*T)(unsafe.Pointer(value))
	return &newPointer
}

var translateFunc i18n.TranslateFunc
var translateFuncOnce sync.Once

// init translation function for translation app error
func AppErrorInit(t i18n.TranslateFunc) {
	translateFuncOnce.Do(func() {
		translateFunc = t
	})
}

// GetMillis is a convenience method to get milliseconds since epoch, utc time
func GetMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// GetTimeForMillis is a convenience method to get time.Time for milliseconds since epoch.
func GetTimeForMillis(millis int64) time.Time {
	return time.Unix(0, millis*int64(time.Millisecond))
}

// GetStartOfDayMillis is a convenience method to get milliseconds since epoch for provided date's start of day
func GetStartOfDayMillis(thisTime time.Time, timeZoneOffset int) int64 {
	localSearchTimeZone := time.FixedZone("Local Search Time Zone", timeZoneOffset)
	resultTime := time.Date(thisTime.Year(), thisTime.Month(), thisTime.Day(), 0, 0, 0, 0, localSearchTimeZone)
	return utils.TimeToMicrosecondsSinceEpochInt64(resultTime)
}

// GetEndOfDayMillis is a convenience method to get milliseconds since epoch for provided date's end of day
func GetEndOfDayMillis(thisTime time.Time, timeZoneOffset int) int64 {
	localSearchTimeZone := time.FixedZone("Local Search Time Zone", timeZoneOffset)
	resultTime := time.Date(thisTime.Year(), thisTime.Month(), thisTime.Day(), 23, 59, 59, 999999999, localSearchTimeZone)
	return utils.TimeToMicrosecondsSinceEpochInt64(resultTime)
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
	appErr.Translate(translateFunc)
	return appErr
}

// common function for creating model.AppError type
//
// Example:
//
//		collection := &Collection{
//			Id: "dsdsdre984jf8se990834",
//			Name: "Hello World",
//		}
//		outer := CreateAppErrorForModel(
//				"model.collection.is_valid.%s.app_error",
//	     "collection_id=",
//	     "Collection.IsValid",
//	 )
//		return outer("name", &collection.Id)
//
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
	bytes, err := json.Marshal(&model)
	if err != nil {
		return ""
	}
	return string(bytes)
}

// Decodes json string into model.
//
// If decoding process encounter error, model will be nil
func ModelFromJson(model interface{}, data io.Reader) error {
	return json.NewDecoder(data).Decode(&model)
}

func (a *AppError) ToJSON() string {
	return ModelToJson(a)
}

func AppErrorFromJSon(data io.Reader) *AppError {
	str := ""
	bytes, err := io.ReadAll(data)
	if err != nil {
		str = err.Error()
	} else {
		str = string(bytes)
	}

	decoder := json.NewDecoder(strings.NewReader(str))
	var er AppError
	err = decoder.Decode(&er)
	if err != nil {
		return NewAppError("AppErrorFromJson", "model.utils.decode_json.app_error", nil, "body: "+str, http.StatusInternalServerError)
	}
	return &er
}

// Merge addkey-value pairs from `other` to current StringInterface
func (m StringInterface) Merge(other StringInterface) {
	for key, value := range other {
		m[key] = value
	}
}

func GetPreferredTimezone(timezone Map[string, string]) string {
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

// NewRandomString returns a random string of the given length.
// The resulting entropy will be (5 * length) bits.
func NewRandomString(length int) string {
	data := make([]byte, 1+(length*5/8))
	rand.Read(data)
	return encoding.EncodeToString(data)[:length]
}

// NewPromoCode generates promo codes for Voucher, Giftcard. Codes have format of XXXX-XXXX-XXXX (length 14)
func NewPromoCode() string {
	raw := uuid.NewString()
	raw = strings.Map(func(r rune) rune {
		switch {
		case r == '-':
			return 0x0
		case unicode.IsLetter(r):
			return unicode.ToUpper(r)
		default:
			return r
		}
	}, raw)
	return fmt.Sprintf("%s-%s-%s", raw[:4], raw[4:8], raw[8:12])
}

// code must has format of XXXX-XXXX-XXXX to satisfy this pattern
var PromoCodeRegex = regexp.MustCompile(`^([A-Z0-9]{4})-([A-Z0-9]{4})-([A-Z0-9]{4})$`)

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
	decoder := json.NewDecoder(data)

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
	decoder := json.NewDecoder(data)

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
func IsValidHTTPURL(rawURL string) bool {
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
		if IsValidHTTPURL(*link) {
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
		return
	}

	if er.params == nil {
		er.Message = T(er.Id)
	} else {
		er.Message = T(er.Id, er.params)
	}
}

// SystemMessage need an translate func,
func (er *AppError) SystemMessage(T i18n.TranslateFunc) string {
	if er.params == nil {
		return T(er.Id)
	}
	return T(er.Id, er.params)
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

// ValidateStoreFrontUrl checks if given url is valid and is one of
// CORS whitelisted items.
func ValidateStoreFrontUrl(config *Config, urlValue string) *AppError {
	// try check if provided redirect url is valid
	parsedRedirectUrl, err := url.Parse(urlValue)
	if err != nil {
		return NewAppError("ValidateStoreFrontUrl", "app.provided_url_invalid.app_error", map[string]interface{}{"Value": urlValue}, "", http.StatusBadRequest)
	}

	allowedClients := strings.Fields(*config.ServiceSettings.AllowCorsFrom)
	for _, validClient := range allowedClients {
		if validClient == parsedRedirectUrl.Host {
			return nil
		}
	}
	return NewAppError("ValidateStoreFrontUrl", "app.provided_url_invalid.app_error", map[string]interface{}{"Value": urlValue}, "", http.StatusBadRequest)
}

//	{
//		"blocks": [
//			{
//				"data": {
//					"text": "There is life in outer space. This vibrant light speed yellow paint brings life to any surface. Goes on easy and dries at light speed."
//				},
//				"type": "paragraph"
//			}
//		]
//	}
func DraftJSContentToRawText(content StringInterface, sep string) string {
	if sep == "" {
		sep = "/n"
	}

	if content == nil {
		return ""
	}

	blocks, ok := content["blocks"]
	if !ok {
		return ""
	}

	blocksSlice, ok := blocks.([]StringInterface)
	if !ok {
		return ""
	}

	paragraphs := []string{}

	for _, block := range blocksSlice {
		data, ok := block["data"]
		if !ok {
			continue
		}

		dataMap, ok := data.(Map[string, string])
		if !ok {
			continue
		}

		text, ok := dataMap["text"]
		if !ok {
			continue
		}

		paragraphs = append(paragraphs, text)
	}

	return strings.Join(paragraphs, sep)
}

// getSubpathScript renders the inline script that defines window.publicPath to change how webpack loads assets.
func getSubpathScript(subpath string) string {
	if subpath == "" {
		subpath = "/"
	}

	newPath := path.Join(subpath, "static") + "/"

	return fmt.Sprintf("window.publicPath='%s'", newPath)
}

// GetSubpathScriptHash computes the script-src addition required for the subpath script to bypass CSP protections.
func GetSubpathScriptHash(subpath string) string {
	// No hash is required for the default subpath.
	if subpath == "" || subpath == "/" {
		return ""
	}

	scriptHash := sha256.Sum256([]byte(getSubpathScript(subpath)))

	return fmt.Sprintf(" 'sha256-%s'", base64.StdEncoding.EncodeToString(scriptHash[:]))
}

// UpdateAssetsSubpathInDir rewrites assets in the given directory to assume the application is
// hosted at the given subpath instead of at the root. No changes are written unless necessary.
func UpdateAssetsSubpathInDir(subpath, directory string) error {
	if subpath == "" {
		subpath = "/"
	}

	staticDir, found := fileutils.FindDir(directory)
	if !found {
		return errors.New("failed to find client dir")
	}

	staticDir, err := filepath.EvalSymlinks(staticDir)
	if err != nil {
		return errors.Wrapf(err, "failed to resolve symlinks to %s", staticDir)
	}

	rootHTMLPath := filepath.Join(staticDir, "root.html")
	oldRootHTML, err := os.ReadFile(rootHTMLPath)
	if err != nil {
		return errors.Wrap(err, "failed to open root.html")
	}

	oldSubpath := "/"

	// Determine if a previous subpath had already been rewritten into the assets.
	reWebpackPublicPathScript := regexp.MustCompile("window.publicPath='([^']+/)static/'")
	alreadyRewritten := false
	if matches := reWebpackPublicPathScript.FindStringSubmatch(string(oldRootHTML)); matches != nil {
		oldSubpath = matches[1]
		alreadyRewritten = true
	}

	pathToReplace := path.Join(oldSubpath, "static") + "/"
	newPath := path.Join(subpath, "static") + "/"

	slog.Debug("Rewriting static assets", slog.String("from_subpath", oldSubpath), slog.String("to_subpath", subpath))

	newRootHTML := string(oldRootHTML)

	reCSP := regexp.MustCompile(`<meta http-equiv="Content-Security-Policy" content="script-src 'self' cdn.rudderlabs.com/ js.stripe.com/v3([^"]*)">`)
	if results := reCSP.FindAllString(newRootHTML, -1); len(results) == 0 {
		return fmt.Errorf("failed to find 'Content-Security-Policy' meta tag to rewrite")
	}

	newRootHTML = reCSP.ReplaceAllLiteralString(newRootHTML, fmt.Sprintf(
		`<meta http-equiv="Content-Security-Policy" content="script-src 'self' cdn.rudderlabs.com/ js.stripe.com/v3%s">`,
		GetSubpathScriptHash(subpath),
	))

	// Rewrite the root.html references to `/static/*` to include the given subpath.
	// This potentially includes a previously injected inline script that needs to
	// be updated (and isn't covered by the cases above).
	newRootHTML = strings.Replace(newRootHTML, pathToReplace, newPath, -1)

	if alreadyRewritten && subpath == "/" {
		// Remove the injected script since no longer required. Note that the rewrite above
		// will have affected the script, so look for the new subpath, not the old one.
		oldScript := getSubpathScript(subpath)
		newRootHTML = strings.Replace(newRootHTML, fmt.Sprintf("</style><script>%s</script>", oldScript), "</style>", 1)

	} else if !alreadyRewritten && subpath != "/" {
		// Otherwise, inject the script to define `window.publicPath`.
		script := getSubpathScript(subpath)
		newRootHTML = strings.Replace(newRootHTML, "</style>", fmt.Sprintf("</style><script>%s</script>", script), 1)
	}

	// Write out the updated root.html.
	if err = os.WriteFile(rootHTMLPath, []byte(newRootHTML), 0); err != nil {
		return errors.Wrapf(err, "failed to update root.html with subpath %s", subpath)
	}

	// Rewrite the manifest.json and *.css references to `/static/*` (or a previously rewritten subpath).
	err = filepath.Walk(staticDir, func(walkPath string, info os.FileInfo, err error) error {
		if filepath.Base(walkPath) == "manifest.json" || filepath.Ext(walkPath) == ".css" {
			old, err := os.ReadFile(walkPath)
			if err != nil {
				return errors.Wrapf(err, "failed to open %s", walkPath)
			}
			new := strings.Replace(string(old), pathToReplace, newPath, -1)
			if err = os.WriteFile(walkPath, []byte(new), 0); err != nil {
				return errors.Wrapf(err, "failed to update %s with subpath %s", walkPath, subpath)
			}
		}

		return nil
	})
	if err != nil {
		return errors.Wrapf(err, "error walking %s", staticDir)
	}

	return nil
}

// UpdateAssetsSubpath rewrites assets in the /client directory to assume the application is hosted
// at the given subpath instead of at the root. No changes are written unless necessary.
func UpdateAssetsSubpath(subpath string) error {
	return UpdateAssetsSubpathInDir(subpath, CLIENT_DIR)
}

// UpdateAssetsSubpathFromConfig uses UpdateAssetsSubpath and any path defined in the SiteURL.
func UpdateAssetsSubpathFromConfig(config *Config) error {
	// Don't rewrite in development environments, since webpack in developer mode constantly
	// updates the assets and must be configured separately.
	if BuildNumber == "dev" {
		slog.Debug("Skipping update to assets subpath since dev build")
		return nil
	}

	// Similarly, don't rewrite during a CI build, when the assets may not even be present.
	if os.Getenv("IS_CI") == "true" {
		slog.Debug("Skipping update to assets subpath since CI build")
		return nil
	}

	subpath, err := GetSubpathFromConfig(config)
	if err != nil {
		return err
	}

	return UpdateAssetsSubpath(subpath)
}

// GetSubpathFromConfig returns subpath from given config's ServiceSettings.SiteURL
func GetSubpathFromConfig(config *Config) (string, error) {
	if config == nil {
		return "", errors.New("no config provided")
	} else if config.ServiceSettings.SiteURL == nil {
		return "/", nil
	}

	u, err := url.Parse(*config.ServiceSettings.SiteURL)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse SiteURL from config")
	}

	if u.Path == "" {
		return "/", nil
	}

	return path.Clean(u.Path), nil
}

// ordering in "ASC" | "DESC" order
type OrderDirection string

const (
	DESC OrderDirection = "DESC"
	ASC  OrderDirection = "ASC"
)

func (d OrderDirection) String() string {
	return string(d)
}

// E.g
//
//	var name = "hello"
//	var namePtr *string = nil
//
//	GetValueOfpointerOrNil(&name) == "hello"
//	GetValueOfpointerOrNil(namePtr) == nil
func GetValueOfpointerOrNil[T any](pointer *T) any {
	if pointer == nil {
		return nil
	}

	return *pointer
}

// type StringMap map[string]string

// func (s StringMap) DeepCopy() StringMap {
// 	if s == nil {
// 		return nil
// 	}
// 	res := StringMap{}
// 	for key, value := range s {
// 		res[key] = value
// 	}

// 	return res
// }

// func (s StringMap) Merge(other StringMap) StringMap {
// 	for key, value := range other {
// 		s[key] = value
// 	}
// 	return s
// }

// func (m StringMap) Pop(key string, defaultValue ...string) string {
// 	v := m.Get(key, defaultValue...)
// 	delete(m, key)
// 	return v
// }

// func (m StringMap) Get(key string, defaultValue ...string) string {
// 	if v, ok := m[key]; ok {
// 		return v
// 	}
// 	if len(defaultValue) > 0 {
// 		return defaultValue[0]
// 	}
// 	return ""
// }

// func (m StringMap) Set(key, value string) {
// 	m[key] = value
// }

// func (m StringMap) Keys() []string {
// 	return lo.Keys(m)
// }

// func (m StringMap) Values() []string {
// 	return lo.Values(m)
// }

// // Scan converts database column value to StringMap
// func (m *StringMap) Scan(value any) error {
// 	if value == nil {
// 		return nil
// 	}

// 	buf, ok := value.([]byte)
// 	if ok {
// 		return json.Unmarshal(buf, m)
// 	}

// 	str, ok := value.(string)
// 	if ok {
// 		return json.Unmarshal([]byte(str), m)
// 	}

// 	return errors.New("received value is neither a byte slice nor string")
// }

// const (
// 	maxPropSizeBytes = 1024 * 1024
// )

// var ErrMaxPropSizeExceeded = fmt.Errorf("max prop size of %d exceeded", maxPropSizeBytes)

// // Value converts StringMap to database value
// func (m StringMap) Value() (driver.Value, error) {
// 	sz := 0
// 	for k := range m {
// 		sz += len(k) + len(m[k])
// 		if sz > maxPropSizeBytes {
// 			return nil, ErrMaxPropSizeExceeded
// 		}
// 	}

// 	buf, err := json.Marshal(m)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// Key wasn't found. We fall back to the default case.
// 	return string(buf), nil
// }

type Map[K cmp.Ordered, V any] map[K]V

func (m *Map[K, V]) Add(key K, value V) {
	(*m)[key] = value
}

func (m *Map[K, V]) Pop(key K, defaultValue ...V) V {
	value, ok := (*m)[key]
	if ok {
		delete(*m, key)
		return value
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	var res V
	return res
}

func (m Map[K, V]) Get(key K, defaultValue ...V) V {
	value, ok := m[key]
	if ok {
		return value
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	var res V
	return res
}

func (m *Map[K, V]) Merge(other Map[K, V]) {
	for key, value := range other {
		(*m)[key] = value
	}
}
