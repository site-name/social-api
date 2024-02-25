package model_helper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/timezones"
	"github.com/sitename/sitename/modules/util"
	"golang.org/x/crypto/bcrypt"
)

const (
	SESSION_COOKIE_TOKEN              = "SNAUTHTOKEN"
	SESSION_COOKIE_USER               = "SNUSERID"
	SESSION_COOKIE_CSRF               = "SNCSRF"
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
	SESSION_CACHE_SIZE                = 35000
	SESSION_ACTIVITY_TIMEOUT          = 1000 * 60 * 5 // 5 minutes
	SESSION_USER_ACCESS_TOKEN_EXPIRY  = 100 * 365     // 100 years
	USER_NICK_NAME_MAX_RUNES          = 64
	USER_AUTH_DATA_MAX_LENGTH         = 128
	USER_ROLES_MAX_LENGTH             = 256
)

const (
	USER_AUTH_SERVICE_LDAP       = "ldap"
	LDAP_PUBLIC_CERTIFICATE_NAME = "ldap-public.crt"
	LDAP_PRIVATE_KEY_NAME        = "ldap-private.key"
)

const (
	STATUS_OUT_OF_OFFICE   = "ooo"
	STATUS_OFFLINE         = "offline"
	STATUS_ONLINE          = "online"
	STATUS_CACHE_SIZE      = SESSION_CACHE_SIZE
	STATUS_CHANNEL_TIMEOUT = 20000  // 20 seconds
	STATUS_MIN_UPDATE_TIME = 120000 // 2 minutes
)

const (
	ME                        = "me"
	PUSH_NOTIFY_PROP          = "push"
	EMAIL_NOTIFY_PROP         = "email"
	USER_NOTIFY_MENTION       = "mention"
	MENTION_KEYS_NOTIFY_PROP  = "mention_keys"
	USER_FIRST_NAME_MAX_RUNES = 64
	USER_LAST_NAME_MAX_RUNES  = 64
	USER_TIMEZONE_MAX_RUNES   = 256
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
	SystemGuestRoleId           = "system_guest"
)

// ----- address ---------

type AddressTypeEnum string

const (
	ADDRESS_TYPE_SHIPPING AddressTypeEnum = "shipping"
	ADDRESS_TYPE_BILLING  AddressTypeEnum = "billing"
)

func (e AddressTypeEnum) IsValid() bool {
	switch e {
	case ADDRESS_TYPE_SHIPPING, ADDRESS_TYPE_BILLING:
		return true
	}
	return false
}

func AddressFullname(a model.Address) string {
	return fmt.Sprintf("%s %s", a.FirstName, a.LastName)
}

func AddressString(a model.Address) string {
	if a.CompanyName != "" {
		return fmt.Sprintf("%s - %s", a.CompanyName, AddressFullname(a))
	}
	return AddressFullname(a)
}

func AddressPreSave(a *model.Address) {
	if a.ID == "" {
		a.ID = NewId()

	}
	addressCommonPre(a)
	a.CreatedAt = GetMillis()
	a.UpdatedAt = a.CreatedAt
}

func addressCommonPre(a *model.Address) {
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

func AddressPreUpdate(a *model.Address) {
	addressCommonPre(a)
	a.UpdatedAt = GetMillis()
}

func AddressIsValid(a model.Address) *AppError {
	if !IsValidId(a.ID) {
		return NewAppError("Address.IsValid", "model.address.is_valid.id.app_error", nil, "please provide valid id", http.StatusBadRequest)
	}
	if !IsValidId(a.UserID) {
		return NewAppError("Address.IsValid", "model.address.is_valid.user_id.app_error", nil, "please provide valid user id", http.StatusBadRequest)
	}
	if a.CreatedAt <= 0 {
		return NewAppError("Address.IsValid", "model.address.is_valid.created_at.app_error", nil, "", http.StatusBadRequest)
	}
	if a.UpdatedAt <= 0 {
		return NewAppError("Address.IsValid", "model.address.is_valid.updated_at.app_error", nil, "", http.StatusBadRequest)
	}
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
	if utf8.RuneCountInString(name) > USER_FIRST_NAME_MAX_RUNES {
		return false
	}
	if !ValidUsernameChars.MatchString(name) {
		return false
	}
	return !RestrictedUsernames[name]
}

// ----------- customer note ----------

func CustomerNoteIsValid(c model.CustomerNote) *AppError {
	if !c.UserID.IsNil() && !IsValidId(*c.UserID.String) {
		return NewAppError("CustomerNote.IsValid", "model.customer_note.is_valid.user_id.app_error", nil, "please provide valid user id", http.StatusBadRequest)
	}
	if !IsValidId(c.CustomerID) {
		return NewAppError("CustomerNote.IsValid", "model.customer_note.is_valid.customer_id.app_error", nil, "please provide valid customer id", http.StatusBadRequest)
	}
	return nil
}

// -------- session ---------

func SessionPreSave(s *model.Session) {
	if s.Token == "" {
		s.Token = NewId()
	}
	s.CreatedAt = GetMillis()
	s.LastActivityAt = s.CreatedAt
	if s.Props == nil {
		s.Props = model_types.JsonMap{}
	}
}

func SessionIsUnrestricted(s model.Session) bool {
	return s.Local
}

func SessionIsValid(s model.Session) *AppError {
	if !IsValidId(s.UserID) {
		return NewAppError("Session.IsValid", "model.session.is_valid.user_id.app_error", nil, "", http.StatusBadRequest)
	}
	if s.CreatedAt == 0 {
		return NewAppError("Session.IsValid", "model.session.is_valid.create_at.app_error", nil, "", http.StatusBadRequest)
	}

	if len(s.Roles) > USER_ROLES_MAX_LENGTH {
		return NewAppError("Session.IsValid", "model.session.is_valid.roles_limit.app_error",
			map[string]any{"Limit": USER_ROLES_MAX_LENGTH}, "session_id="+s.ID, http.StatusBadRequest)
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
	return GetMillis() > s.ExpiresAt
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

func UserAccessTokenIsValid(t model.UserAccessToken) *AppError {
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
	t.Description = SanitizeUnicode(t.Description)
}

// ------------ user ---------------

type UserPatch struct {
	Username    *string             `json:"username"`
	Password    *string             `json:"password,omitempty"`
	Nickname    *string             `json:"nickname"`
	FirstName   *string             `json:"first_name"`
	LastName    *string             `json:"last_name"`
	Email       *string             `json:"email"`
	Locale      *model.LanguageCode `json:"locale"`
	Timezone    model_types.JsonMap `json:"timezone"`
	NotifyProps model_types.JsonMap `json:"notify_props,omitempty"`
}

type UserAuth struct {
	Password    string  `json:"password,omitempty"`
	AuthData    *string `json:"auth_data,omitempty"`
	AuthService string  `json:"auth_service,omitempty"`
}

var _ util.Hashable = (*UserWrapper)(nil)

type UserWrapper struct {
	model.User
}

// GetEmail implements util.Hashable.
func (u UserWrapper) GetEmail() string {
	return u.Email
}

// GetId implements util.Hashable.
func (u UserWrapper) GetId() string {
	return u.ID
}

// GetLastLogin implements util.Hashable.
func (u UserWrapper) GetLastLogin() time.Time {
	return util.TimeFromMillis(u.LastActivityAt)
}

// GetPassword implements util.Hashable.
func (u UserWrapper) GetPassword() string {
	return u.Password
}

func UserIsLDAP(u model.User) bool {
	return u.AuthService == USER_AUTH_SERVICE_LDAP
}

func UserIsSAML(u model.User) bool {
	return u.AuthService == USER_AUTH_SERVICE_SAML
}

func UserIsSSO(u model.User) bool {
	return u.AuthService != "" && u.AuthService != USER_AUTH_SERVICE_EMAIL
}

func UserIsOauth(u model.User) bool {
	return u.AuthService == SERVICE_GOOGLE || u.AuthService == SERVICE_OPENID
}

func UserPreSave(u *model.User) {
	if u.ID == "" {
		u.ID = NewId()
	}
	userCommonPre(u)
	if u.Password != "" {
		u.Password = HashPassword(u.Password)
	}

	u.CreatedAt = GetMillis()
	u.UpdatedAt = u.CreatedAt
}

func UserPreUpdate(u *model.User) {
	if _, ok := u.NotifyProps[MENTION_KEYS_NOTIFY_PROP]; ok {
		// Remove any blank mention keys
		splitKeys := strings.Split(u.NotifyProps[MENTION_KEYS_NOTIFY_PROP].(string), ",")
		goodKeys := []string{}
		for _, key := range splitKeys {
			if key != "" {
				goodKeys = append(goodKeys, strings.ToLower(key))
			}
		}
		u.NotifyProps[MENTION_KEYS_NOTIFY_PROP] = strings.Join(goodKeys, ",")
	}
	userCommonPre(u)
	u.UpdatedAt = GetMillis()
}

func userCommonPre(u *model.User) {
	if u.NotifyProps == nil || len(u.NotifyProps) == 0 {
		UserSetDefaultNotifications(u)
	}
	u.Username = SanitizeUnicode(u.Username)
	u.FirstName = SanitizeUnicode(u.FirstName)
	u.LastName = SanitizeUnicode(u.LastName)
	u.Nickname = SanitizeUnicode(u.Nickname)
	u.Username = NormalizeUsername(u.Username)
	u.Email = NormalizeEmail(u.Email)

	if !u.AuthData.IsNil() && *u.AuthData.String == "" {
		u.AuthData.String = nil
	}
	if u.Props == nil {
		u.Props = model_types.JsonMap{}
	}

	if model.LanguageCode(u.Locale).IsValid() != nil {
		u.Locale = DEFAULT_LOCALE
	}
	if u.Timezone == nil {
		u.Timezone = timezones.DefaultUserTimezone()
	}
}

func UserSetDefaultNotifications(u *model.User) {
	u.NotifyProps = model_types.JsonMap{
		EMAIL_NOTIFY_PROP: "true",
		PUSH_NOTIFY_PROP:  USER_NOTIFY_MENTION,
	}
}

func PatchUser(u *model.User, patch UserPatch) {
	if patch.Username != nil {
		u.Username = *patch.Username
	}
	if patch.Nickname != nil {
		u.Nickname = *patch.Nickname
	}
	if patch.FirstName != nil {
		u.FirstName = *patch.FirstName
	}
	if patch.NotifyProps != nil {
		u.NotifyProps = patch.NotifyProps
	}
	if patch.LastName != nil {
		u.LastName = *patch.LastName
	}
	if patch.Email != nil {
		u.Email = *patch.Email
	}
	if patch.Locale != nil && patch.Locale.IsValid() == nil {
		u.Locale = *patch.Locale
	}
	if patch.Timezone != nil {
		u.Timezone = patch.Timezone
	}
}

func UserEtag(u model.User, showFullName, showEmail string) string {
	return Etag(u.ID, u.UpdatedAt, u.TermsOfServiceID, u.TermsOfServiceCreatedAt, showFullName, showEmail)
}

func UserSanitize(u *model.User, options map[string]bool) {
	u.Password = ""
	u.AuthData.String = GetPointerOfValue("")
	u.MfaSecret = ""

	if len(options) != 0 && !options["email"] {
		u.Email = ""
	}
	if len(options) != 0 && !options["fullname"] {
		u.FirstName = ""
		u.LastName = ""
	}
	if len(options) != 0 && !options["passwordupdate"] {
		u.LastPasswordUpdate = 0
	}
	if len(options) != 0 && !options["authservice"] {
		u.AuthService = ""
	}
}

func UserSanitizeInput(u *model.User, isAdmin bool) {
	if !isAdmin {
		u.AuthData.String = GetPointerOfValue("")
		u.AuthService = ""
		u.EmailVerified = false
	}
	u.LastPasswordUpdate = 0
	u.LastPictureUpdate = 0
	u.FailedAttempts = 0
	u.MfaActive = false
	u.MfaSecret = ""
}

func UserUpdateMentionKeysFromUsername(u *model.User, oldUsername string) {
	nonUsernameKeys := []string{}
	for _, key := range UserGetMentionKeys(*u) {
		if key != oldUsername && key != "@"+oldUsername {
			nonUsernameKeys = append(nonUsernameKeys, key)
		}
	}

	u.NotifyProps[MENTION_KEYS_NOTIFY_PROP] = ""
	if len(nonUsernameKeys) > 0 {
		u.NotifyProps[MENTION_KEYS_NOTIFY_PROP] = u.NotifyProps[MENTION_KEYS_NOTIFY_PROP].(string) + "," + strings.Join(nonUsernameKeys, ",")
	}
}

func UserGetMentionKeys(u model.User) []string {
	var keys []string
	for _, key := range strings.Split(u.NotifyProps[MENTION_KEYS_NOTIFY_PROP].(string), ",") {
		trimmedKey := strings.TrimSpace(key)
		if trimmedKey == "" {
			continue
		}
		keys = append(keys, trimmedKey)
	}

	return keys
}

func UserClearNonProfileFields(u *model.User) {
	u.Password = ""
	u.AuthData.String = GetPointerOfValue("")
	u.MfaSecret = ""
	u.EmailVerified = false
	u.LastPasswordUpdate = 0
	u.FailedAttempts = 0
}

func UserSanitizeProfile(u *model.User, options map[string]bool) {
	UserClearNonProfileFields(u)
	UserSanitize(u, options)
}

func UserGetFullName(u model.User) string {
	if u.FirstName != "" && u.LastName != "" {
		return u.FirstName + " " + u.LastName
	} else if u.FirstName != "" {
		return u.FirstName
	} else if u.LastName != "" {
		return u.LastName
	} else {
		return ""
	}
}

func getUserDisplayName(u model.User, baseName, nameFormat string) string {
	displayName := baseName

	if nameFormat == SHOW_NICKNAME_FULLNAME {
		if u.Nickname != "" {
			displayName = u.Nickname
		} else if fullName := UserGetFullName(u); fullName != "" {
			displayName = fullName
		}
	} else if nameFormat == SHOW_FULLNAME {
		if fullName := UserGetFullName(u); fullName != "" {
			displayName = fullName
		}
	}

	return displayName
}

func UserGetDisplayName(u model.User, nameFormat string) string {
	displayName := u.Username
	return getUserDisplayName(u, displayName, nameFormat)
}

func UserGetDisplayNameWithPrefix(u model.User, nameFormat, prefix string) string {
	displayName := prefix + u.Username

	return getUserDisplayName(u, displayName, nameFormat)
}

func UserGetRoles(u model.User) util.AnyArray[string] {
	return strings.Fields(u.Roles)
}

func UserGetRawRoles(u model.User) string {
	return u.Roles
}

func IsValidUserRoles(userRoles string) bool {
	roles := strings.Fields(strings.TrimSpace(userRoles))

	for _, r := range roles {
		if !IsValidRoleName(r) {
			return false
		}
	}

	// Exclude just the system_admin role explicitly to prevent mistakes
	if len(roles) == 1 && roles[0] == SystemAdminRoleId {
		return false
	}

	return true
}

func UserToPatch(u model.User) UserPatch {
	return UserPatch{
		Username:  &u.Username,
		Password:  &u.Password,
		Nickname:  &u.Nickname,
		FirstName: &u.FirstName,
		LastName:  &u.LastName,
		Email:     &u.Email,
		Locale:    GetPointerOfValue(model.LanguageCode(u.Locale)),
		Timezone:  u.Timezone,
	}
}

func (u *UserPatch) SetField(fieldName string, fieldValue string) {
	switch fieldName {
	case "FirstName":
		u.FirstName = &fieldValue
	case "LastName":
		u.LastName = &fieldValue
	case "Nickname":
		u.Nickname = &fieldValue
	case "Email":
		u.Email = &fieldValue
	case "Username":
		u.Username = &fieldValue
	}
}

// UserSearchOptions captures internal parameters derived from the user's permissions and a
// UserSearch request.
type UserSearchOptions struct {
	// IsAdmin tracks whether or not the search is being conducted by an administrator.
	IsAdmin bool
	// AllowEmails allows search to examine the emails of users.
	AllowEmails bool
	// AllowFullNames allows search to examine the full names of users, vs. just usernames and nicknames.
	AllowFullNames bool
	// AllowInactive configures whether or not to return inactive users in the search results.
	AllowInactive bool
	// Narrows the search to the group constrained users
	// GroupConstrained bool
	// Limit limits the total number of results returned.
	Limit int
	// Filters for the given role
	Role string
	// Filters for users that have any of the given system roles
	Roles []string
}

type UserGetOptions struct {
	Inactive bool
	// Filters the active users
	Active bool
	// Filters for the given role
	Role string
	// Filters for users matching any of the given system wide roles
	Roles []string
	// Filters for users matching any of the given channel roles, must be used with InChannelId
	// Sorting option
	Sort string
	// Restrict to search in a list of teams and channels
	// ViewRestrictions *ViewUsersRestrictions
	// Page
	Page int
	// Page size
	PerPage int
}

func UserMakeNonNil(u *model.User) {
	if u.Props == nil {
		u.Props = model_types.JsonMap{}
	}
	if u.NotifyProps == nil {
		u.NotifyProps = model_types.JsonMap{}
	}
}

func UserIsValid(u model.User) *AppError {
	if u.CreatedAt <= 0 {
		return InvalidUserError(model.UserColumns.CreatedAt, u.ID, u.CreatedAt)
	}
	if u.UpdatedAt <= 0 {
		return InvalidUserError(model.UserColumns.UpdatedAt, u.ID, u.UpdatedAt)
	}
	if !IsValidUsername(u.Username) {
		return InvalidUserError(model.UserColumns.Username, u.ID, u.Username)
	}
	if len(u.Email) > USER_EMAIL_MAX_LENGTH || u.Email == "" || !IsValidEmail(u.Email) {
		return InvalidUserError(model.UserColumns.Email, u.ID, u.Email)
	}
	if utf8.RuneCountInString(u.Nickname) > USER_NICK_NAME_MAX_RUNES {
		return InvalidUserError(model.UserColumns.Nickname, u.ID, u.Nickname)
	}
	if utf8.RuneCountInString(u.FirstName) > USER_FIRST_NAME_MAX_RUNES {
		return InvalidUserError(model.UserColumns.FirstName, u.ID, u.FirstName)
	}
	if utf8.RuneCountInString(u.LastName) > USER_LAST_NAME_MAX_RUNES {
		return InvalidUserError(model.UserColumns.LastName, u.ID, u.LastName)
	}
	if !u.AuthData.IsNil() && len(*u.AuthData.String) > USER_AUTH_DATA_MAX_LENGTH {
		return InvalidUserError(model.UserColumns.AuthData, u.ID, u.AuthData)
	}
	if !u.AuthData.IsNil() && *u.AuthData.String != "" && u.AuthService == "" {
		return InvalidUserError("auth_data_type", u.ID, *u.AuthData.String+" "+u.AuthService)
	}
	if u.Password != "" && !u.AuthData.IsNil() && *u.AuthData.String != "" {
		return InvalidUserError("auth_data_pwd", u.ID, *u.AuthData.String)
	}
	if model.LanguageCode(u.Locale).IsValid() != nil {
		return InvalidUserError(model.UserColumns.Locale, u.ID, u.Locale)
	}
	if len(u.Timezone) > 0 {
		if tzJSON, err := json.Marshal(u.Timezone); err != nil {
			return NewAppError("User.IsValid", "model.user.is_valid.marshal.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		} else if utf8.RuneCount(tzJSON) > USER_TIMEZONE_MAX_RUNES {
			return InvalidUserError("timezone_limit", u.ID, u.Timezone)
		}
	}
	if len(u.Roles) > USER_ROLES_MAX_LENGTH {
		return NewAppError("User.IsValid", "model.user.is_valid.roles_limit.app_error",
			map[string]any{"Limit": USER_ROLES_MAX_LENGTH}, "user_id="+u.ID+" roles_limit="+u.Roles, http.StatusBadRequest)
	}
	if !u.DefaultBillingAddressID.IsNil() && !IsValidId(*u.DefaultBillingAddressID.String) {
		return InvalidUserError(model.UserColumns.DefaultBillingAddressID, u.ID, *u.DefaultBillingAddressID.String)
	}
	if !u.DefaultShippingAddressID.IsNil() && !IsValidId(*u.DefaultShippingAddressID.String) {
		return InvalidUserError(model.UserColumns.DefaultShippingAddressID, u.ID, *u.DefaultBillingAddressID.String)
	}
	if !IsValidId(u.TermsOfServiceID) {
		return InvalidUserError(model.UserColumns.TermsOfServiceID, u.ID, u.TermsOfServiceID)
	}

	return nil
}

func InvalidUserError(fieldName, userId string, fieldValue any) *AppError {
	id := fmt.Sprintf("model.user.is_valid.%s.app_error", fieldName)
	details := ""
	if userId != "" {
		details = "user_id=" + userId
	}
	details += fmt.Sprintf(" %s=%v", fieldName, fieldValue)
	return NewAppError("User.IsValid", id, nil, details, http.StatusBadRequest)
}

type UserForIndexing struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Roles     string `json:"roles"`
	CreatedAt int64  `json:"create_at"`
	DeletedAt int64  `json:"delete_at"`
}

type UserUpdate struct {
	Old *model.User
	New *model.User
}

// --------------- customer event

func CustomerEventIsValid(ce model.CustomerEvent) *AppError {
	if err := ce.Type.IsValid(); err != nil {
		return NewAppError("CustomerEventIsValid", "model.customer_event.is_valid.type.app_error", nil, err.Error(), http.StatusBadRequest)
	}
	if !ce.OrderID.IsNil() && !IsValidId(*ce.OrderID.String) {
		return NewAppError("CustomerEventIsValid", "model.customer_event.is_valid.order_id.app_error", nil, "invalid order id", http.StatusBadRequest)
	}
	if !ce.UserID.IsNil() && !IsValidId(*ce.UserID.String) {
		return NewAppError("CustomerEventIsValid", "model.customer_event.is_valid.user_id.app_error", nil, "invalid user id", http.StatusBadRequest)
	}
	return nil
}

// ----------------- staff notification recipient --------------------

func StaffNotificationRecipientIsValid(s model.StaffNotificationRecipient) *AppError {
	if !s.UserID.IsNil() && !IsValidId(*s.UserID.String) {
		return NewAppError("StaffNotificationRecipientIsValid", "model.staff_notification_recipient.is_valid.user_id.app_error", nil, "invalid user id", http.StatusBadRequest)
	}
	if !s.StaffEmail.IsNil() && !IsValidEmail(*s.StaffEmail.String) {
		return NewAppError("StaffNotificationRecipientIsValid", "model.staff_notification_recipient.is_valid.staff_email.app_error", nil, "invalid staff email", http.StatusBadRequest)
	}
	return nil
}

// ------------------------ status -------------------

func StatusIsValid(s model.Status) *AppError {
	if !IsValidId(s.UserID) {
		return NewAppError("StatusIsValid", "model.status.is_valid.user_id.app_error", nil, "invalid user id", http.StatusBadRequest)
	}
	return nil
}

// ------------------ term of service ----------------

func TermsOfServicePreSave(t *model.TermsOfService) {
	if t.CreatedAt == 0 {
		t.CreatedAt = GetMillis()
	}
}

func TermsOfServiceIsValid(t model.TermsOfService) *AppError {
	if t.CreatedAt == 0 {
		return NewAppError("TermsOfServiceIsValid", "model.terms_of_service.is_valid.created_at.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(t.UserID) {
		return NewAppError("TermsOfServiceIsValid", "model.terms_of_service.is_valid.user_id.app_error", nil, "invalid user id", http.StatusBadRequest)
	}

	return nil
}

type MfaSecret struct {
	Secret string `json:"secret"`
	QRCode string `json:"qr_code"`
}

// UserSearch captures the parameters provided by a client for initiating a user search.
type UserSearch struct {
	Term          string   `json:"term"`
	AllowInactive bool     `json:"allow_inactive"`
	Limit         int      `json:"limit"`
	Role          string   `json:"role"`
	Roles         []string `json:"roles"`
}

// Make sure you acually want to use this function. In context.go there are functions to check permissions
//
// This function should not be used to check permissions.
func IsInRole(userRoles string, inRole string) bool {
	return strings.Contains(userRoles, inRole)
}

type UsersStats struct {
	TotalUsersCount int64 `json:"total_users_count"`
}

type UserFilterOptions struct {
	CommonQueryOptions
}
