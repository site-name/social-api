package model_helper

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
	"golang.org/x/crypto/bcrypt"
)

const (
	USER_NAME_PART_MAX_RUNES = 64
)

const (
	SESSION_COOKIE_TOKEN              = "SNAUTHTOKEN"
	SESSION_COOKIE_USER               = "SNUSERID"
	SESSION_COOKIE_CSRF               = "SNCSRF"
	SESSION_CACHE_SIZE                = 35000
	SESSION_PROP_PLATFORM             = "platform"
	SESSION_PROP_OS                   = "os"
	SESSION_PROP_BROWSER              = "browser"
	SESSION_PROP_TYPE                 = "type"
	SESSION_PROP_USER_ACCESS_TOKEN_ID = "user_access_token_id"
	SESSION_PROP_IS_BOT               = "is_bot"
	SESSION_PROP_IS_BOT_VALUE         = "true"
	SESSION_TYPE_USER_ACCESS_TOKEN    = "UserAccessToken"
	SESSION_TYPE_CLOUD_KEY            = "CloudKey"
	SESSION_TYPE_REMOTECLUSTER_TOKEN  = "RemoteClusterToken"
	SESSION_PROP_IS_GUEST             = "is_guest"
	SESSION_ACTIVITY_TIMEOUT          = 1000 * 60 * 5 // 5 minutes
	SESSION_USER_ACCESS_TOKEN_EXPIRY  = 100 * 365     // 100 years
)

const (
	STATUS_OUT_OF_OFFICE   = "ooo"
	STATUS_OFFLINE         = "offline"
	STATUS_ONLINE          = "online"
	STATUS_CACHE_SIZE      = SESSION_CACHE_SIZE
	STATUS_CHANNEL_TIMEOUT = 20000  // 20 seconds
	STATUS_MIN_UPDATE_TIME = 120000 // 2 minutes
)

// Options for counting users
type UserCountOptions struct {
	// Should include deleted users (of any type)
	IncludeDeleted bool
	// Exclude regular users
	ExcludeRegularUsers bool
	Roles               []string
}

var (
	SysconsoleAncillaryPermissions        map[string][]string // SysconsoleAncillaryPermissions maps the non-sysconsole permissions required by each sysconsole view.
	SystemManagerDefaultPermissions       util.AnyArray[string]
	SystemUserManagerDefaultPermissions   util.AnyArray[string]
	SystemReadOnlyAdminDefaultPermissions util.AnyArray[string]
	BuiltInSchemeManagedRoleIDs           util.AnyArray[string]
	NewSystemRoleIDs                      util.AnyArray[string]
)

const (
	SystemUserRoleId            = "system_user"
	SystemAdminRoleId           = "system_admin"
	SystemUserAccessTokenRoleId = "system_user_access_token"
	SystemUserManagerRoleId     = "system_user_manager"
	SystemReadOnlyAdminRoleId   = "system_read_only_admin"
	SystemManagerRoleId         = "system_manager"
	ShopAdminRoleId             = "shop_admin"
	ShopStaffRoleId             = "shop_staff"
)

// ----- address ---------

func AddressFullname(a model.Address) string {
	return fmt.Sprintf("%s %s", a.FirstName, a.LastName)
}

func AddressString(a model.Address) string {
	if a.CompanyName != "" {
		return fmt.Sprintf("%s - %s", a.CompanyName, AddressFullname(a))
	}
	return AddressFullname(a)
}

func AddressCommonPre(a *model.Address) {
	if strings.TrimSpace(a.FirstName) == "" {
		a.FirstName = "first_name"
	}
	if strings.TrimSpace(a.LastName) == "" {
		a.LastName = "last_name"
	}
	a.FirstName = SanitizeUnicode(CleanNamePart(a.FirstName))
	a.LastName = SanitizeUnicode(CleanNamePart(a.LastName))
	if a.Country.IsValid() != nil {
		a.Country = DEFAULT_COUNTRY
	}
}

func AddressIsValid(a *model.Address) *AppError {
	if !IsValidNamePart(a.FirstName) {
		return NewAppError("Address.IsValid", "model.address.is_valid.first_name.app_error", nil, "please provide valid first name", http.StatusBadRequest)
	}
	if !IsValidNamePart(a.LastName) {
		return NewAppError("Address.IsValid", "model.address.is_valid.last_name.app_error", nil, "please provide valid last name", http.StatusBadRequest)
	}
	if !IsAllNumbers(a.PostalCode) {
		return NewAppError("Address.IsValid", "model.address.is_valid.postal_code.app_error", nil, "please provide valid postal code", http.StatusBadRequest)
	}
	if a.Country.IsValid() != nil {
		return NewAppError("Address.IsValid", "model.address.is_valid.country.app_error", nil, "please provide valid country code", http.StatusBadRequest)
	}
	if str, ok := util.ValidatePhoneNumber(a.Phone, a.Country.String()); !ok {
		return NewAppError("Address.IsValid", "model.address.is_valid.phone.app_error", nil, "please provide valid phone", http.StatusBadRequest)
	} else {
		a.Phone = str
	}

	return nil
}

func AddressObfuscate(a model.Address) model.Address {
	res := a
	res.FirstName = util.ObfuscateString(res.FirstName, false)
	res.LastName = util.ObfuscateString(res.LastName, false)
	res.CompanyName = util.ObfuscateString(res.CompanyName, false)
	res.StreetAddress1 = util.ObfuscateString(res.StreetAddress1, false)
	res.StreetAddress2 = util.ObfuscateString(res.StreetAddress2, false)
	res.Phone = util.ObfuscateString(res.Phone, true)
	return res
}

// CleanNamePart should be used to clean first/last name of users or addresses, ...
func CleanNamePart(nameValue string) string {
	name := NormalizeUsername(strings.Replace(nameValue, " ", "-", -1))
	for _, value := range ReservedName {
		if name == value {
			name = strings.Replace(name, value, "", -1)
		}
	}
	name = strings.TrimSpace(name)
	for _, c := range name {
		char := string(c)
		if !ValidUsernameChars.MatchString(char) {
			name = strings.Replace(nameValue, char, "-", -1)
		}
	}
	name = strings.Trim(name, "-")

	if !IsValidNamePart(name) {
		name = "a" + strings.ReplaceAll(NewRandomString(8), "-", "")
	}

	return name
}

func IsValidNamePart(name string) bool {
	if utf8.RuneCountInString(name) > USER_NAME_PART_MAX_RUNES {
		return false
	}
	if !ValidUsernameChars.MatchString(name) {
		return false
	}
	return !RestrictedUsernames[name]
}

// ----------- customer note ----------

func CustomerNoteIsValid(c *model.CustomerNote) *AppError {
	if !c.UserID.IsNil() && !IsValidId(*c.UserID.String) {
		return NewAppError("CustomerNote.IsValid", "model.customer_note.is_valid.user_id.app_error", nil, "please provide valid user id", http.StatusBadRequest)
	}
	if !IsValidId(c.CustomerID) {
		return NewAppError("CustomerNote.IsValid", "model.customer_note.is_valid.customer_id.app_error", nil, "please provide valid customer id", http.StatusBadRequest)
	}
	return nil
}

// -------- session ---------

func SessionIsUnrestricted(s *model.Session) bool {
	return s.Local
}

func SessionIsValid(s *model.Session) *AppError {
	if !IsValidId(s.UserID) {
		return NewAppError("Session.IsValid", "model.session.is_valid.user_id.app_error", nil, "", http.StatusBadRequest)
	}
	return nil
}

func SessionSanitize(s *model.Session) {
	s.Token = ""
}

func SessionIsExpired(s model.Session) bool {
	if s.ExpiresAt <= 0 {
		return false
	}
	if GetMillis() > s.ExpiresAt {
		return true
	}
	return false
}

// Deprecated: SetExpireInDays is deprecated and should not be used.
//
//	            Use (*App).SetSessionExpireInDays instead which handles the
//				   cases where the new ExpiresAt is not relative to CreatedAt.
func SessionSetExpireInDays(s *model.Session, days int) {
	if s.CreatedAt == 0 {
		s.ExpiresAt = GetMillis() + (1000 * 60 * 60 * 24 * int64(days))
	} else {
		s.ExpiresAt = s.CreatedAt + (1000 * 60 * 60 * 24 * int64(days))
	}
}

// AddProp adds given value to session's Props with key of given key
func SessionAddProp(s *model.Session, key string, value string) {
	if s.Props == nil {
		s.Props = map[string]any{}
	}
	s.Props[key] = value
}

// IsMobileApp check if current session has non-empty `DeviceId` field or IsMobile()
func SessionIsMobileApp(s model.Session) bool {
	return s.DeviceID != "" || SessionIsMobile(s)
}

// IsMobile checks if the Props field has an item is ("isMobile": "true")
func SessionIsMobile(s model.Session) bool {
	if s.Props == nil || len(s.Props) == 0 {
		return false
	}
	val, ok := s.Props[USER_AUTH_SERVICE_IS_MOBILE]
	return ok && strings.EqualFold(val.(string), "true")
}

func SessionIsSaml(s model.Session) bool {
	if s.Props == nil || len(s.Props) == 0 {
		return false
	}
	val, ok := s.Props[USER_AUTH_SERVICE_IS_SAML]
	if !ok {
		return false
	}
	isSaml, err := strconv.ParseBool(val.(string))
	if err != nil {
		slog.Debug("Error parsing boolean property from Session", slog.Err(err))
		return false
	}
	return isSaml
}

func SessionIsOAuthUser(s model.Session) bool {
	if s.Props == nil || len(s.Props) == 0 {
		return false
	}

	val, ok := s.Props[USER_AUTH_SERVICE_IS_OAUTH]
	if !ok {
		return false
	}
	isOAuthUser, err := strconv.ParseBool(val.(string))
	if err != nil {
		slog.Debug("Error parsing boolean property from Session", slog.Err(err))
		return false
	}
	return isOAuthUser
}

func SessionIsSSOLogin(s model.Session) bool {
	return SessionIsOAuthUser(s) || SessionIsSaml(s)
}

// GetUserRoles turns current session's Roles into a slice of strings
func SessionGetUserRoles(s model.Session) util.AnyArray[string] {
	return strings.Fields(s.Roles)
}

// GenerateCSRF simply generates new UUID, then add that uuid to its "Props" with key is "csrf". Finally returns that token
func SessionGenerateCSRF(s *model.Session) string {
	if s.Props == nil {
		s.Props = map[string]any{}
	}
	token := NewId()
	s.Props["csrf"] = token
	return token
}

// get value with key of "csrf" from session's Props
func SessionGetCSRF(s model.Session) string {
	if s.Props == nil || len(s.Props) == 0 {
		return ""
	}
	value, ok := s.Props["csrf"]
	if !ok {
		return ""
	}

	return value.(string)
}

// HashPassword generates a hash using the bcrypt.GenerateFromPassword
func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		panic(err)
	}

	return string(hash)
}

// --------- user access token --------

func UserAccessTokenIsValid(t *model.UserAccessToken) *AppError {
	if !IsValidId(t.Token) {
		return NewAppError("UserAccessToken.IsValid", "model.user_access_token.is_valid.token.app_error", nil, "please provide valid token", http.StatusBadRequest)
	}
	if !IsValidId(t.UserID) {
		return NewAppError("UserAccessToken.IsValid", "model.user_access_token.is_valid.user_id.app_error", nil, "please provide valid user id", http.StatusBadRequest)
	}
	return nil
}

func UserAccessTokenCommonPre(t *model.UserAccessToken) {
	if t.Token == "" {
		t.Token = NewId()
	}
}

func IsSSOUser(u model.User) bool {
	return u.AuthService != "" && u.AuthService != USER_AUTH_SERVICE_EMAIL
}
