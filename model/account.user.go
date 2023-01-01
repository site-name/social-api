package model

import (
	"encoding/json"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/sitename/sitename/modules/timezones"
	"github.com/sitename/sitename/modules/util"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/language"
)

// constants used in package account
const (
	ME                            = "me"
	PUSH_NOTIFY_PROP              = "push"
	EMAIL_NOTIFY_PROP             = "email"
	USER_NOTIFY_MENTION           = "mention"
	MENTION_KEYS_NOTIFY_PROP      = "mention_keys"
	USER_NICKNAME_MAX_RUNES       = 64
	USER_POSITION_MAX_RUNES       = 128
	USER_FIRST_NAME_MAX_RUNES     = 64
	USER_LAST_NAME_MAX_RUNES      = 64
	USER_AUTH_DATA_MAX_LENGTH     = 128
	USER_PASSWORD_MAX_LENGTH      = 72
	USER_HASH_PASSWORD_MAX_LENGTH = 128
	USER_LOCALE_MAX_LENGTH        = 5
	USER_TIMEZONE_MAX_RUNES       = 256
	USER_LANGUAGE_CODE_MAX_LENGTH = 10
)

// types for addresses
const (
	ADDRESS_TYPE_SHIPPING = "shipping"
	ADDRESS_TYPE_BILLING  = "billing"
)

// User contains the details about the user.
// This struct's serializer methods are auto-generated. If a new field is added/removed,
// please run make gen-serialized.
type User struct {
	Id                       string    `json:"id"`
	Email                    string    `json:"email"`      // unique
	Username                 string    `json:"username"`   // unique
	FirstName                string    `json:"first_name"` // can be empty
	LastName                 string    `json:"last_name"`  // can be empty
	DefaultShippingAddressID *string   `json:"default_shipping_address,omitempty"`
	DefaultBillingAddressID  *string   `json:"default_billing_address,omitempty"`
	Password                 string    `json:"password,omitempty"`
	AuthData                 *string   `json:"auth_data,omitempty"`
	AuthService              string    `json:"auth_service"`
	EmailVerified            bool      `json:"email_verified,omitempty"`
	Nickname                 string    `json:"nickname"`
	Roles                    string    `json:"roles"`
	Props                    StringMap `json:"props,omitempty"`
	NotifyProps              StringMap `json:"notify_props,omitempty"`
	LastPasswordUpdate       int64     `json:"last_password_update,omitempty"`
	LastPictureUpdate        int64     `json:"last_picture_update,omitempty"`
	FailedAttempts           int       `json:"failed_attempts,omitempty"`
	Locale                   string    `json:"locale"` // user's language
	Timezone                 StringMap `json:"timezone"`
	MfaActive                bool      `json:"mfa_active,omitempty"`
	MfaSecret                string    `json:"mfa_secret,omitempty"`
	CreateAt                 int64     `json:"create_at,omitempty"`
	UpdateAt                 int64     `json:"update_at,omitempty"`
	DeleteAt                 int64     `json:"delete_at"`
	IsActive                 bool      `json:"is_active"`
	Note                     *string   `json:"note"`
	JwtTokenKey              string    `json:"jwt_token_key"`
	LastActivityAt           int64     `json:"last_activity_at,omitempty"`
	TermsOfServiceId         string    `json:"terms_of_service_id,omitempty"`
	TermsOfServiceCreateAt   int64     `json:"terms_of_service_create_at,omitempty"`
	DisableWelcomeEmail      bool      `json:"disable_welcome_email"`
	ModelMetadata
}

// UserMap is a map from a userId to a user object.
// It is used to generate methods which can be used for fast serialization/de-serialization.
type UserMap map[string]*User

//msgp:ignore UserUpdate
type UserUpdate struct {
	Old *User
	New *User
}

//msgp:ignore UserPatch
type UserPatch struct {
	Username    *string   `json:"username"`
	Password    *string   `json:"password,omitempty"`
	Nickname    *string   `json:"nickname"`
	FirstName   *string   `json:"first_name"`
	LastName    *string   `json:"last_name"`
	Email       *string   `json:"email"`
	Locale      *string   `json:"locale"`
	Timezone    StringMap `json:"timezone"`
	NotifyProps StringMap `json:"notify_props,omitempty"`
}

//msgp:ignore UserAuth
type UserAuth struct {
	Password    string  `json:"password,omitempty"`
	AuthData    *string `json:"auth_data,omitempty"`
	AuthService string  `json:"auth_service,omitempty"`
}

//msgp:ignore UserForIndexing
type UserForIndexing struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Roles     string `json:"roles"`
	CreateAt  int64  `json:"create_at"`
	DeleteAt  int64  `json:"delete_at"`
}

//msgp:ignore UserSlice
type UserSlice []*User

func (u UserSlice) Usernames() []string {
	usernames := make([]string, len(u))
	for i, user := range u {
		usernames[i] = user.Username
	}
	sort.Strings(usernames)
	return usernames
}

// IDs returns slice of uuids from slice of users
func (u UserSlice) IDs() []string {
	ids := make([]string, len(u))
	for i, user := range u {
		ids[i] = user.Id
	}
	return ids
}

func (u UserSlice) FilterByActive(active bool) UserSlice {
	var matches []*User
	for _, user := range u {
		if user.DeleteAt == 0 && active {
			matches = append(matches, user)
		} else if user.DeleteAt != 0 && !active {
			matches = append(matches, user)
		}
	}
	return UserSlice(matches)
}

func (u UserSlice) FilterByID(ids []string) UserSlice {
	var matches []*User
	for _, user := range u {
		for _, id := range ids {
			if id == user.Id {
				matches = append(matches, user)
			}
		}
	}
	return UserSlice(matches)
}

func (u UserSlice) FilterWithoutID(ids []string) UserSlice {
	var keep []*User
	for _, user := range u {
		present := false
		for _, id := range ids {
			if id == user.Id {
				present = true
			}
		}
		if !present {
			keep = append(keep, user)
		}
	}
	return UserSlice(keep)
}

func (u *User) DeepCopy() *User {
	copyUser := *u

	if u.DefaultShippingAddressID != nil {
		copyUser.DefaultShippingAddressID = NewString(*u.DefaultShippingAddressID)
	}
	if u.DefaultBillingAddressID != nil {
		copyUser.DefaultBillingAddressID = NewString(*u.DefaultBillingAddressID)
	}
	if u.AuthData != nil {
		copyUser.AuthData = NewString(*u.AuthData)
	}
	copyUser.NotifyProps = u.NotifyProps.DeepCopy()
	copyUser.Props = u.Props.DeepCopy()
	copyUser.Timezone = u.Timezone.DeepCopy()
	if u.Note != nil {
		copyUser.Note = NewString(*u.Note)
	}

	return &copyUser
}

// IsValid validates the user and returns an error if it isn't configured
// correctly.
func (u *User) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"user.is_valid.%s.app_error",
		"user_id=",
		"User.IsValid")

	if !IsValidId(u.Id) {
		return outer("id", nil)
	}
	if u.CreateAt == 0 {
		return outer("create_at", &u.Id)
	}
	if u.UpdateAt == 0 {
		return outer("update_at", &u.Id)
	}
	if !IsValidUsername(u.Username) {
		return outer("username", &u.Id)
	}
	if len(u.Email) > USER_EMAIL_MAX_LENGTH || u.Email == "" || !IsValidEmail(u.Email) {
		return outer("email", &u.Id)
	}
	if utf8.RuneCountInString(u.Nickname) > USER_NICKNAME_MAX_RUNES {
		return outer("nickname", &u.Id)
	}
	if u.FirstName != "" && !IsValidNamePart(u.FirstName, FirstName) {
		return outer("first_name", &u.Id)
	}
	if u.LastName != "" && !IsValidNamePart(u.LastName, LastName) {
		return outer("last_name", &u.Id)
	}
	if u.AuthData != nil && len(*u.AuthData) > USER_AUTH_DATA_MAX_LENGTH {
		return outer("auth_data", &u.Id)
	}
	if u.AuthData != nil && *u.AuthData != "" && u.AuthService == "" {
		return outer("auth_data_type", &u.Id)
	}
	if u.Password != "" && u.AuthData != nil && *u.AuthData != "" {
		return outer("auth_data_pwd", &u.Id)
	}
	if len(u.Password) > USER_PASSWORD_MAX_LENGTH {
		return outer("password_limit", &u.Id)
	}
	if tag, err := language.Parse(u.Locale); err != nil || !strings.EqualFold(tag.String(), u.Locale) {
		return outer("locale", &u.Id)
	}
	if len(u.Timezone) > 0 {
		if tzJson, err := json.Marshal(u.Timezone); err != nil {
			return NewAppError("User.IsValid", "user.is_valid.marshal.app_error", nil, err.Error(), http.StatusInternalServerError)
		} else if utf8.RuneCount(tzJson) > USER_TIMEZONE_MAX_RUNES {
			return outer("timezone_limit", &u.Id)
		}
	}

	return nil
}

// PreSave will set the Id and Username if missing.  It will also fill
// in the CreateAt, UpdateAt times.  It will also hash the password.  It should
// be run before saving the user to the db.
func (u *User) PreSave() {
	if u.Id == "" {
		u.Id = NewId()
	}
	if u.Username == "" {
		u.Username = NewId()
	}

	u.CreateAt = GetMillis()
	u.UpdateAt = u.CreateAt
	u.LastPasswordUpdate = u.CreateAt
	u.MfaActive = false

	if u.NotifyProps == nil || len(u.NotifyProps) == 0 {
		u.SetDefaultNotifications()
	}

	if u.Timezone == nil {
		u.Timezone = timezones.DefaultUserTimezone()
	}
	if u.Password != "" {
		u.Password = HashPassword(u.Password)
	}
}

func (u *User) commonPre() {
	u.Username = SanitizeUnicode(u.Username)
	u.FirstName = SanitizeUnicode(u.FirstName)
	u.LastName = SanitizeUnicode(u.LastName)
	u.Nickname = SanitizeUnicode(u.Nickname)
	u.Username = NormalizeUsername(u.Username)
	u.Email = NormalizeEmail(u.Email)

	if u.AuthData != nil && *u.AuthData == "" {
		u.AuthData = nil
	}
	if u.Props == nil {
		u.Props = make(map[string]string)
	}
	if u.Locale == "" {
		u.Locale = DEFAULT_LOCALE
	} else {
		u.Locale = strings.ToLower(u.Locale)
	}
}

// PreUpdate should be run before updating the user in the db.
func (u *User) PreUpdate() {
	u.UpdateAt = GetMillis()

	if u.NotifyProps == nil || len(u.NotifyProps) == 0 {
		u.SetDefaultNotifications()
	} else if _, ok := u.NotifyProps[MENTION_KEYS_NOTIFY_PROP]; ok {
		// Remove any blank mention keys
		splitKeys := strings.Split(u.NotifyProps[MENTION_KEYS_NOTIFY_PROP], ",")
		goodKeys := []string{}
		for _, key := range splitKeys {
			if key != "" {
				goodKeys = append(goodKeys, strings.ToLower(key))
			}
		}
		u.NotifyProps[MENTION_KEYS_NOTIFY_PROP] = strings.Join(goodKeys, ",")
	}
}

func (u *User) IsSSOUser() bool {
	return u.AuthService != "" && u.AuthService != USER_AUTH_SERVICE_EMAIL
}

// IsLDAPUser checks if user's AuthService == USER_AUTH_SERVICE_LDAP = "ldap"
func (u *User) IsLDAPUser() bool {
	return u.AuthService == USER_AUTH_SERVICE_LDAP
}

// IsSAMLUser checks if user's AuthService == USER_AUTH_SERVICE_SAML = "saml"
func (u *User) IsSAMLUser() bool {
	return u.AuthService == USER_AUTH_SERVICE_SAML
}

func (u *User) Patch(patch *UserPatch) {
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

	if patch.Locale != nil {
		u.Locale = *patch.Locale
	}

	if patch.Timezone != nil {
		u.Timezone = patch.Timezone
	}
}

// ToJson convert a User to a json.JSON string
func (u *User) ToJSON() string {
	return ModelToJson(u)
}

func (u *UserPatch) ToJSON() string {
	return ModelToJson(u)
}

func (u *UserAuth) ToJSON() string {
	return ModelToJson(u)
}

// Generate a valid strong etag so the browser can cache the results
func (u *User) Etag(showFullName, showEmail bool) string {
	return Etag(u.Id, u.UpdateAt, u.TermsOfServiceId, u.TermsOfServiceCreateAt, showFullName, showEmail)
}

// Remove any private data from the user object
//
// options's keys can be "email", "fullname", "passwordupdate", "authservice" OR Nothing
func (u *User) Sanitize(options map[string]bool) {
	u.Password = ""
	u.AuthData = NewString("")
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

// Remove any input data from the user object that is not user controlled
func (u *User) SanitizeInput(isAdmin bool) {
	if !isAdmin {
		u.AuthData = NewString("")
		u.AuthService = ""
		u.EmailVerified = false
	}
	u.LastPasswordUpdate = 0
	u.LastPictureUpdate = 0
	u.FailedAttempts = 0
	u.MfaActive = false
	u.MfaSecret = ""
}

// SetDefaultNotifications set default values for user's NotifyProps attribute
func (u *User) SetDefaultNotifications() {
	u.NotifyProps = make(map[string]string)
	u.NotifyProps[EMAIL_NOTIFY_PROP] = "true"
	u.NotifyProps[PUSH_NOTIFY_PROP] = USER_NOTIFY_MENTION
}

func (u *User) UpdateMentionKeysFromUsername(oldUsername string) {
	nonUsernameKeys := []string{}
	for _, key := range u.GetMentionKeys() {
		if key != oldUsername && key != "@"+oldUsername {
			nonUsernameKeys = append(nonUsernameKeys, key)
		}
	}

	u.NotifyProps[MENTION_KEYS_NOTIFY_PROP] = ""
	if len(nonUsernameKeys) > 0 {
		u.NotifyProps[MENTION_KEYS_NOTIFY_PROP] += "," + strings.Join(nonUsernameKeys, ",")
	}
}

func (u *User) GetMentionKeys() []string {
	var keys []string
	for _, key := range strings.Split(u.NotifyProps[MENTION_KEYS_NOTIFY_PROP], ",") {
		trimmedKey := strings.TrimSpace(key)
		if trimmedKey == "" {
			continue
		}
		keys = append(keys, trimmedKey)
	}

	return keys
}

// ClearNonProfileFields reset user's password, authData, MfaSecret, EmailVerified,
// LastPasswordUpdate, FailedAttempts to their default values
func (u *User) ClearNonProfileFields() {
	u.Password = ""
	u.AuthData = NewString("")
	u.MfaSecret = ""
	u.EmailVerified = false
	u.LastPasswordUpdate = 0
	u.FailedAttempts = 0
}

func (u *User) SanitizeProfile(options map[string]bool) {
	u.ClearNonProfileFields()

	u.Sanitize(options)
}

func (u *User) GetFullName() string {
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

func (u *User) getDisplayName(baseName, nameFormat string) string {
	displayName := baseName

	if nameFormat == SHOW_NICKNAME_FULLNAME {
		if u.Nickname != "" {
			displayName = u.Nickname
		} else if fullName := u.GetFullName(); fullName != "" {
			displayName = fullName
		}
	} else if nameFormat == SHOW_FULLNAME {
		if fullName := u.GetFullName(); fullName != "" {
			displayName = fullName
		}
	}

	return displayName
}

func (u *User) GetDisplayName(nameFormat string) string {
	displayName := u.Username

	return u.getDisplayName(displayName, nameFormat)
}

func (u *User) GetDisplayNameWithPrefix(nameFormat, prefix string) string {
	displayName := prefix + u.Username

	return u.getDisplayName(displayName, nameFormat)
}

func (u *User) GetRoles() []string {
	return strings.Fields(u.Roles)
}

// GetRawRoles return user's raw roles
func (u *User) GetRawRoles() string {
	return u.Roles
}

// IsValidUserRoles checks if user's roles are both valid
func IsValidUserRoles(userRoles string) bool {
	roles := strings.Fields(strings.TrimSpace(userRoles))

	for _, r := range roles {
		if !IsValidRoleName(r) {
			return false
		}
	}

	// Exclude just the system_admin role explicitly to prevent mistakes
	if len(roles) == 1 && roles[0] == "system_admin" {
		return false
	}

	return true
}

// Make sure you acually want to use this function. In context.go there are functions to check permissions
// This function should not be used to check permissions.
func (u *User) IsInRole(inRole string) bool {
	return IsInRole(u.Roles, inRole)
}

// Make sure you acually want to use this function. In context.go there are functions to check permissions
//
// This function should not be used to check permissions.
//
// IsGuest checks if user's roles contains "system_guest"
func (u *User) IsGuest() bool {
	return IsInRole(u.Roles, SystemGuestRoleId)
}

// IsSystemAdmin checks if user's roles contains "system_admin"
func (u *User) IsSystemAdmin() bool {
	return IsInRole(u.Roles, SystemAdminRoleId)
}

// Make sure you acually want to use this function. In context.go there are functions to check permissions
//
// This function should not be used to check permissions.
func IsInRole(userRoles string, inRole string) bool {
	roles := strings.Split(userRoles, " ")

	for _, r := range roles {
		if r == inRole {
			return true
		}
	}

	return false
}

// IsOAuthUser checks if user is authenticated via google or open oauth systems
func (u *User) IsOAuthUser() bool {
	return u.AuthService == SERVICE_GOOGLE || u.AuthService == SERVICE_OPENID
}

func (u *User) ToPatch() *UserPatch {
	return &UserPatch{
		Username:  &u.Username,
		Password:  &u.Password,
		Nickname:  &u.Nickname,
		FirstName: &u.FirstName,
		LastName:  &u.LastName,
		Email:     &u.Email,
		Locale:    &u.Locale,
		Timezone:  u.Timezone,
	}
}

// set value for user's given fieldName.
//
// fieldName can be either: "FirstName" | "LastName" | "Nickname" | "Email" | "Username"
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

// UserFromJson will decode the input and return a User
func UserFromJson(data io.Reader) *User {
	var user *User
	ModelFromJson(&user, data)
	return user
}

func UserPatchFromJson(data io.Reader) *UserPatch {
	var user UserPatch
	ModelFromJson(&user, data)
	return &user
}

func UserAuthFromJson(data io.Reader) *UserAuth {
	var user UserAuth
	ModelFromJson(&user, data)
	return &user
}

func UserMapToJson(u map[string]*User) string {
	return ModelToJson(&u)
}

func UserMapFromJson(data io.Reader) map[string]*User {
	var users map[string]*User
	ModelFromJson(&users, data)
	return users
}

func UserListToJson(u []*User) string {
	return ModelToJson(&u)
}

func UserListFromJson(data io.Reader) []*User {
	var users []*User
	ModelFromJson(&users, data)
	return users
}

// HashPassword generates a hash using the bcrypt.GenerateFromPassword
func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		panic(err)
	}

	return string(hash)
}

// MakeNonNil sets empty value for user's Props and NotifyProps fields
func (u *User) MakeNonNil() {
	if u.Props == nil {
		u.Props = make(map[string]string)
	}

	if u.NotifyProps == nil {
		u.NotifyProps = make(map[string]string)
	}
}

var (
	_ util.Hashable = (*User)(nil)
)

func (u *User) GetId() string {
	return u.Id
}

func (u *User) GetPassword() string {
	return u.Password
}

func (u *User) GetLastLogin() time.Time {
	return util.TimeFromMillis(u.LastActivityAt)
}

func (u *User) GetEmail() string {
	return u.Email
}

type StringMAP map[string]string

func (s StringMAP) DeepCopy() StringMAP {
	if s == nil {
		return nil
	}

	res := StringMAP{}
	for key, value := range s {
		res[key] = value
	}

	return res
}

func (s StringMAP) Merge(other StringMAP) StringMAP {
	for key, value := range other {
		s[key] = value
	}

	return s
}

// Common abstract model for other models to inherit from
type ModelMetadata struct {
	Metadata        StringMAP `json:"metadata"`
	PrivateMetadata StringMAP `json:"private_metadata"`
}

func (m *ModelMetadata) PopulateFields() {
	if m.PrivateMetadata == nil {
		m.PrivateMetadata = make(map[string]string)
	}
	if m.Metadata == nil {
		m.Metadata = make(map[string]string)
	}
}

func (p ModelMetadata) DeepCopy() ModelMetadata {
	return ModelMetadata{
		p.Metadata.DeepCopy(),
		p.PrivateMetadata.DeepCopy(),
	}
}

type WhichMeta string

const (
	PrivateMetadata WhichMeta = "private"
	Metadata        WhichMeta = "metadata"
)

func (p *ModelMetadata) GetValueFromMeta(key string, defaultValue string, which WhichMeta) string {

	if which == PrivateMetadata { // get from private metadata
		if p.PrivateMetadata == nil {
			return defaultValue
		}

		if vl, ok := p.PrivateMetadata[key]; ok {
			return vl
		}
	} else if which == Metadata { // get from metadata
		if p.Metadata == nil {
			return defaultValue
		}

		if vl, ok := p.Metadata[key]; ok {
			return vl
		}
	}

	return defaultValue
}

func (p *ModelMetadata) StoreValueInMeta(items map[string]string, which WhichMeta) {

	if which == PrivateMetadata {
		if p.PrivateMetadata == nil {
			p.PrivateMetadata = make(map[string]string)
		}

		for k, vl := range items {
			p.PrivateMetadata[k] = vl
		}
	} else if which == Metadata {
		if p.Metadata == nil {
			p.Metadata = make(map[string]string)
		}

		for k, vl := range items {
			p.Metadata[k] = vl
		}
	}
}

func (p *ModelMetadata) ClearMeta(which WhichMeta) {

	if which == PrivateMetadata {
		for k := range p.PrivateMetadata {
			delete(p.PrivateMetadata, k)
		}
	} else if which == Metadata {
		for k := range p.Metadata {
			delete(p.Metadata, k)
		}
	}
}

func (p *ModelMetadata) DeleteValueFromMeta(key string, which WhichMeta) {

	if which == PrivateMetadata {
		delete(p.PrivateMetadata, key)
	} else if which == Metadata {
		delete(p.Metadata, key)
	}
}
