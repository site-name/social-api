package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/samber/lo"
	"github.com/sitename/sitename/modules/timezones"
	"github.com/sitename/sitename/modules/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// constants used in package account
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

// User contains the details about the user.
// This struct's serializer methods are auto-generated. If a new field is added/removed,
// please run make gen-serialized.
type User struct {
	Id                       string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	Email                    string    `json:"email" gorm:"type:varchar(128);unique:users_email_key;index:users_email_index_key;column:Email"`            // unique; varchar(128)
	Username                 string    `json:"username" gorm:"type:varchar(64);unique:users_username_key;index:users_username_index_key;column:Username"` // unique; varchar(64)
	FirstName                string    `json:"first_name" gorm:"type:varchar(64);index:users_firstname_key;column:FirstName"`                             // can be empty, varchar(64)
	LastName                 string    `json:"last_name" gorm:"type:varchar(64);index:users_lastname_key;column:LastName"`                                // can be empty, varchar(64)
	DefaultShippingAddressID *string   `json:"default_shipping_address,omitempty" gorm:"type:uuid;column:DefaultShippingAddressID"`
	DefaultBillingAddressID  *string   `json:"default_billing_address,omitempty" gorm:"type:uuid;column:DefaultBillingAddressID"`
	Password                 string    `json:"password,omitempty" gorm:"column:Password;type:varchar(128)"`  // varchar(128)
	AuthData                 *string   `json:"auth_data,omitempty" gorm:"type:varchar(128);column:AuthData"` // varchar(128)
	AuthService              string    `json:"auth_service" gorm:"type:varchar(20);column:AuthService"`      // varchar(20)
	EmailVerified            bool      `json:"email_verified,omitempty" gorm:"column:EmailVerified"`
	Nickname                 string    `json:"nickname" gorm:"type:varchar(64);column:Nickname"` // varchar(64)
	Roles                    string    `json:"roles" gorm:"type:varchar(200);column:Roles"`      // varchar(200)
	Props                    StringMap `json:"props,omitempty" gorm:"type:jsonb;column:Props"`
	NotifyProps              StringMap `json:"notify_props,omitempty" gorm:"type:jsonb;column:NotifyProps"`
	LastPasswordUpdate       int64     `json:"last_password_update,omitempty" gorm:"column:LastPasswordUpdate"`
	LastPictureUpdate        int64     `json:"last_picture_update,omitempty" gorm:"column:LastPictureUpdate"`
	FailedAttempts           int       `json:"failed_attempts,omitempty" gorm:"column:FailedAttempts"`
	Locale                   string    `json:"locale" gorm:"type:varchar(5);column:Locale"` // user's language; varchar(5); E.g EN
	Timezone                 StringMap `json:"timezone" gorm:"type:jsonb;column:Timezone"`
	MfaActive                bool      `json:"mfa_active,omitempty" gorm:"column:MfaActive"`
	MfaSecret                string    `json:"mfa_secret,omitempty" gorm:"type:varchar(100);column:MfaSecret"`  // varchar(100)
	CreateAt                 int64     `json:"create_at,omitempty" gorm:"autoCreateTime:milli;column:CreateAt"` // read and create only
	UpdateAt                 int64     `json:"update_at,omitempty" gorm:"autoUpdateTime:milli;column:UpdateAt"`
	DeleteAt                 int64     `json:"delete_at" gorm:"type:bigint;column:DeleteAt"`
	IsActive                 bool      `json:"is_active" gorm:"column:IsActive"`
	Note                     *string   `json:"note" gorm:"type:varchar(500);column:Note"`                 // varchar(500)
	JwtTokenKey              string    `json:"jwt_token_key" gorm:"type:varchar(200);column:JwtTokenKey"` // varchar(200)
	LastActivityAt           int64     `json:"last_activity_at,omitempty" gorm:"type:bigint;column:LastActivityAt"`
	TermsOfServiceId         string    `json:"terms_of_service_id,omitempty" gorm:"type:uuid;column:TermsOfServiceId"`
	TermsOfServiceCreateAt   int64     `json:"terms_of_service_create_at,omitempty" gorm:"type:bigint;column:TermsOfServiceCreateAt"`
	DisableWelcomeEmail      bool      `json:"disable_welcome_email" gorm:"column:DisableWelcomeEmail"`
	ModelMetadata

	Addresses                   []*Address                    `json:"-" gorm:"many2many:UserAddresses"`
	CustomerEvents              []*CustomerEvent              `json:"-" gorm:"foreignKey:UserID"`
	CustomerNotes               []*CustomerNote               `json:"-" gorm:"foreignKey:UserID"` // notes that this user has made
	NotesOnMe                   []*CustomerNote               `json:"-" gorm:"foreignKey:CustomerID"`
	StaffNotificationRecipients []*StaffNotificationRecipient `json:"-" gorm:"foreignKey:UserID"`

	// NOTE: field(s) below is/are used for sorting
	OrderCount int `json:"-" gorm:"-"`
}

func (u *User) BeforeCreate(_ *gorm.DB) error { u.PreSave(); return u.IsValid() }
func (u *User) BeforeUpdate(_ *gorm.DB) error { u.PreUpdate(); return u.IsValid() }
func (u *User) TableName() string             { return UserTableName }

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
		copyUser.DefaultShippingAddressID = GetPointerOfValue(*u.DefaultShippingAddressID)
	}
	if u.DefaultBillingAddressID != nil {
		copyUser.DefaultBillingAddressID = GetPointerOfValue(*u.DefaultBillingAddressID)
	}
	if u.AuthData != nil {
		copyUser.AuthData = GetPointerOfValue(*u.AuthData)
	}
	copyUser.NotifyProps = u.NotifyProps.DeepCopy()
	copyUser.Props = u.Props.DeepCopy()
	copyUser.Timezone = u.Timezone.DeepCopy()
	if u.Note != nil {
		copyUser.Note = GetPointerOfValue(*u.Note)
	}

	return &copyUser
}

// IsValid validates the user and returns an error if it isn't configured
// correctly.
func (u *User) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.user.is_valid.%s.app_error",
		"user_id=",
		"User.IsValid")

	if !IsValidUsername(u.Username) {
		return outer("username", &u.Id)
	}
	if u.Email == "" || !IsValidEmail(u.Email) {
		return outer("email", &u.Id)
	}
	if u.FirstName != "" && !IsValidNamePart(u.FirstName, FirstName) {
		return outer("first_name", &u.Id)
	}
	if u.LastName != "" && !IsValidNamePart(u.LastName, LastName) {
		return outer("last_name", &u.Id)
	}
	if u.AuthData != nil && *u.AuthData != "" && u.AuthService == "" {
		return outer("auth_data_type", &u.Id)
	}
	if u.Password != "" && u.AuthData != nil && *u.AuthData != "" {
		return outer("auth_data_pwd", &u.Id)
	}
	if !LanguageCodeEnum(u.Locale).IsValid() {
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
	if u.Password != "" {
		u.Password = HashPassword(u.Password)
	}

	u.commonPre()
}

func (u *User) commonPre() {
	if u.NotifyProps == nil || len(u.NotifyProps) == 0 {
		u.SetDefaultNotifications()
	}
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
	if !LanguageCodeEnum(u.Locale).IsValid() {
		u.Locale = DEFAULT_LOCALE.String()
	}
	if u.Timezone == nil {
		u.Timezone = timezones.DefaultUserTimezone()
	}
}

// PreUpdate should be run before updating the user in the db.
func (u *User) PreUpdate() {
	if _, ok := u.NotifyProps[MENTION_KEYS_NOTIFY_PROP]; ok {
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

	u.commonPre()
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
	if patch.Locale != nil && LanguageCodeEnum(*patch.Locale).IsValid() {
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
	u.AuthData = GetPointerOfValue("")
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
		u.AuthData = GetPointerOfValue("")
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
	u.AuthData = GetPointerOfValue("")
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

func (u *User) GetRoles() util.AnyArray[string] {
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

// IsSystemAdmin checks if user's roles contains "system_admin"
func (u *User) IsSystemAdmin() bool {
	return IsInRole(u.Roles, SystemAdminRoleId)
}

// Make sure you acually want to use this function. In context.go there are functions to check permissions
//
// This function should not be used to check permissions.
func IsInRole(userRoles string, inRole string) bool {
	return strings.Contains(userRoles, inRole)
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

func (m StringMAP) Pop(key string, defaultValue ...string) string {
	v := m.Get(key, defaultValue...)
	delete(m, key)
	return v
}

func (m StringMAP) Get(key string, defaultValue ...string) string {
	if v, ok := m[key]; ok {
		return v
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

func (m StringMAP) Set(key, value string) {
	m[key] = value
}

func (m StringMAP) Keys() []string {
	return lo.MapToSlice(m, func(k string, _ string) string { return k })
}

func (m StringMap) Values() []string {
	return lo.MapToSlice(m, func(_ string, v string) string { return v })
}

// Scan converts database column value to StringMap
func (m *StringMAP) Scan(value any) error {
	if value == nil {
		return nil
	}

	buf, ok := value.([]byte)
	if ok {
		return json.Unmarshal(buf, m)
	}

	str, ok := value.(string)
	if ok {
		return json.Unmarshal([]byte(str), m)
	}

	return errors.New("received value is neither a byte slice nor string")
}

const (
	maxPropSizeBytes = 1024 * 1024
)

var ErrMaxPropSizeExceeded = fmt.Errorf("max prop size of %d exceeded", maxPropSizeBytes)

// Value converts StringMap to database value
func (m StringMAP) Value() (driver.Value, error) {
	sz := 0
	for k := range m {
		sz += len(k) + len(m[k])
		if sz > maxPropSizeBytes {
			return nil, ErrMaxPropSizeExceeded
		}
	}

	buf, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	// Key wasn't found. We fall back to the default case.
	return string(buf), nil
}

// Common abstract model for other models to inherit from
type ModelMetadata struct {
	Metadata        StringMAP `json:"metadata,omitempty" gorm:"type:jsonb;column:Metadata"`
	PrivateMetadata StringMAP `json:"private_metadata,omitempty" gorm:"type:jsonb;column:PrivateMetadata"`
}

const (
	ModelMetadataColumnMetadata        = "Metadata"
	ModelMetadataColumnPrivateMetadata = "PrivateMetadata"
)

// PopulateFields checks if PrivateMetadata or Metadata is nil, if yes then assign them empty maps
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

func (m *ModelMetadata) SetMetadata(key, value string) {
	m.PopulateFields()
	m.Metadata[key] = value
}

func (m *ModelMetadata) DelMetadata(key string) bool {
	m.PopulateFields()
	_, ok := m.Metadata[key]
	if ok {
		delete(m.Metadata, key)
	}
	return ok
}

func (m *ModelMetadata) SetPrivateMetadata(key, value string) {
	m.PopulateFields()
	m.PrivateMetadata[key] = value
}

func (m *ModelMetadata) DelPrivateMetadata(key string) bool {
	m.PopulateFields()
	_, ok := m.PrivateMetadata[key]
	if ok {
		delete(m.PrivateMetadata, key)
	}
	return ok
}
