package account

import (
	"io"
	"net/http"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/json"
	"github.com/sitename/sitename/modules/timezones"
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
	USER_AUTH_SERVICE_EMAIL       = "email"
	USER_NICKNAME_MAX_RUNES       = 64
	USER_POSITION_MAX_RUNES       = 128
	USER_FIRST_NAME_MAX_RUNES     = 64
	USER_LAST_NAME_MAX_RUNES      = 64
	USER_AUTH_DATA_MAX_LENGTH     = 128
	USER_PASSWORD_MAX_LENGTH      = 72
	USER_HASH_PASSWORD_MAX_LENGTH = 128
	USER_LOCALE_MAX_LENGTH        = 5
	USER_TIMEZONE_MAX_RUNES       = 256
)

// NOTE: don't delete this
type StringMap map[string]string

// Address contains information belong to the address
// NOTE: don't delete this
type Address_ struct {
	Id             string  `json:"id"`
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	CompanyName    *string `json:"company_name,omitempty"`
	StreetAddress1 string  `json:"street_address_1,omitempty"`
	StreetAddress2 *string `json:"street_address_2,omitempty"`
	City           string  `json:"city"`
	CityArea       *string `json:"city_area,omitempty"`
	PostalCode     string  `json:"postal_code"`
	Country        string  `json:"country"` // one country name only
	CountryArea    string  `json:"country_area"`
	Phone          string  `json:"phone"`
	CreateAt       int64   `json:"create_at,omitempty"`
	UpdateAt       int64   `json:"update_at,omitempty"`
}

// NOTE: don't delete this
type ModelMetadata struct {
	// Id              string    `json:"string,omitempty"`
	Metadata        StringMap `json:"metadata"`
	PrivateMetadata StringMap `json:"private_metadata"`

	// mutex is used for safe access concurrenly
	// mutex sync.RWMutex `json:"-" db:"-"`
}

// User contains the details about the user.
// This struct's serializer methods are auto-generated. If a new field is added/removed,
// please run make gen-serialized.
type User struct {
	Id                       string      `json:"id"`
	Email                    string      `json:"email"`
	Username                 string      `json:"username"`
	FirstName                string      `json:"first_name"`
	LastName                 string      `json:"last_name"`
	DefaultShippingAddressID *string     `json:"default_shipping_address,omitempty"`
	DefaultBillingAddressID  *string     `json:"default_billing_address,omitempty"`
	Password                 string      `json:"password,omitempty"`
	AuthData                 *string     `json:"auth_data,omitempty"`
	AuthService              string      `json:"auth_service"`
	EmailVerified            bool        `json:"email_verified,omitempty"`
	Nickname                 string      `json:"nickname"`
	Roles                    string      `json:"roles"`
	Props                    StringMap   `json:"props,omitempty"`
	NotifyProps              StringMap   `json:"notify_props,omitempty"`
	LastPasswordUpdate       int64       `json:"last_password_update,omitempty"`
	LastPictureUpdate        int64       `json:"last_picture_update,omitempty"`
	FailedAttempts           int         `json:"failed_attempts,omitempty"`
	Locale                   string      `json:"locale"`
	Timezone                 StringMap   `json:"timezone"`
	MfaActive                bool        `json:"mfa_active,omitempty"`
	MfaSecret                string      `json:"mfa_secret,omitempty"`
	CreateAt                 int64       `json:"create_at,omitempty"`
	UpdateAt                 int64       `json:"update_at,omitempty"`
	DeleteAt                 int64       `json:"delete_at"`
	IsStaff                  bool        `json:"is_staff"`
	IsActive                 bool        `json:"is_active"`
	Note                     *string     `json:"note"`
	Addresses                []*Address_ `json:"addresses" db:"-"`
	JwtTokenKey              string      `json:"jwt_token_key"`
	LastActivityAt           int64       `db:"-" json:"last_activity_at,omitempty"`
	TermsOfServiceId         string      `db:"-" json:"terms_of_service_id,omitempty"`
	TermsOfServiceCreateAt   int64       `db:"-" json:"terms_of_service_create_at,omitempty"`
	DisableWelcomeEmail      bool        `db:"-" json:"disable_welcome_email"`
	ModelMetadata
}

// UserMap is a map from a userId to a user object.
// It is used to generate methods which can be used for fast serialization/de-serialization.
type UserMap map[string]*User

type UserUpdate struct {
	Old *User
	New *User
}

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

type UserAuth struct {
	Password    string  `json:"password,omitempty"`
	AuthData    *string `json:"auth_data,omitempty"`
	AuthService string  `json:"auth_service,omitempty"`
}

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

type UserSlice []*User

func (u UserSlice) Usernames() []string {
	usernames := []string{}
	for _, user := range u {
		usernames = append(usernames, user.Username)
	}
	sort.Strings(usernames)
	return usernames
}

// IDs returns slice of uuids from slice of users
func (u UserSlice) IDs() []string {
	ids := []string{}
	for _, user := range u {
		ids = append(ids, user.Id)
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
	if u.AuthData != nil {
		copyUser.AuthData = model.NewString(*u.AuthData)
	}
	if u.NotifyProps != nil {
		copyUser.NotifyProps = model.CopyStringMap(u.NotifyProps)
	}
	if u.Timezone != nil {
		copyUser.Timezone = model.CopyStringMap(u.Timezone)
	}

	return &copyUser
}

// IsValid validates the user and returns an error if it isn't configured
// correctly.
func (u *User) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.user.is_valid.%s.app_error",
		"user_id=",
		"User.IsValid")

	if !model.IsValidId(u.Id) {
		return outer("id", nil)
	}
	if u.CreateAt == 0 {
		return outer("create_at", &u.Id)
	}
	if u.UpdateAt == 0 {
		return outer("update_at", &u.Id)
	}
	if !model.IsValidUsername(u.Username) {
		return outer("username", &u.Id)
	}
	if len(u.Email) > model.USER_EMAIL_MAX_LENGTH || u.Email == "" || !model.IsValidEmail(u.Email) {
		return outer("email", &u.Id)
	}
	if utf8.RuneCountInString(u.Nickname) > USER_NICKNAME_MAX_RUNES {
		return outer("nickname", &u.Id)
	}
	if !IsValidNamePart(u.FirstName, model.FirstName) {
		return outer("first_name", &u.Id)
	}
	if !IsValidNamePart(u.LastName, model.LastName) {
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
		if tzJson, err := json.JSON.Marshal(u.Timezone); err != nil {
			return model.NewAppError("User.IsValid", "model.user.is_valid.marshal.app_error", nil, err.Error(), http.StatusInternalServerError)
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
		u.Id = model.NewId()
	}
	if u.Username == "" {
		u.Username = model.NewId()
	}
	if u.AuthData != nil && *u.AuthData == "" {
		u.AuthData = nil
	}
	u.Username = model.SanitizeUnicode(u.Username)
	u.FirstName = model.SanitizeUnicode(u.FirstName)
	u.LastName = model.SanitizeUnicode(u.LastName)
	u.Nickname = model.SanitizeUnicode(u.Nickname)
	u.Username = model.NormalizeUsername(u.Username)
	u.Email = model.NormalizeEmail(u.Email)
	u.CreateAt = model.GetMillis()
	u.UpdateAt = u.CreateAt
	u.LastPasswordUpdate = u.CreateAt
	u.MfaActive = false

	if u.Props == nil {
		u.Props = make(map[string]string)
	}
	if u.NotifyProps == nil || len(u.NotifyProps) == 0 {
		u.SetDefaultNotifications()
	}
	if u.Locale == "" {
		u.Locale = model.DEFAULT_LOCALE
	}
	if u.Timezone == nil {
		u.Timezone = timezones.DefaultUserTimezone()
	}
	if u.Password != "" {
		u.Password = HashPassword(u.Password)
	}
}

// PreUpdate should be run before updating the user in the db.
func (u *User) PreUpdate() {
	u.Username = model.SanitizeUnicode(u.Username)
	u.FirstName = model.SanitizeUnicode(u.FirstName)
	u.LastName = model.SanitizeUnicode(u.LastName)
	u.Nickname = model.SanitizeUnicode(u.Nickname)
	u.Username = model.NormalizeUsername(u.Username)
	u.Email = model.NormalizeEmail(u.Email)
	u.UpdateAt = model.GetMillis()

	if u.AuthData != nil && *u.AuthData == "" {
		u.AuthData = nil
	}

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

// IsLDAPUser checks if user's AuthService is "ldap"
func (u *User) IsLDAPUser() bool {
	return u.AuthService == model.USER_AUTH_SERVICE_LDAP
}

// IsSAMLUser checks if user's AuthService is "saml"
func (u *User) IsSAMLUser() bool {
	return u.AuthService == model.USER_AUTH_SERVICE_SAML
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
func (u *User) ToJson() string {
	return model.ModelToJson(u)
}

func (u *UserPatch) ToJson() string {
	return model.ModelToJson(u)
}

func (u *UserAuth) ToJson() string {
	return model.ModelToJson(u)
}

// Generate a valid strong etag so the browser can cache the results
func (u *User) Etag(showFullName, showEmail bool) string {
	return model.Etag(u.Id, u.UpdateAt, u.TermsOfServiceId, u.TermsOfServiceCreateAt, showFullName, showEmail)
}

// Remove any private data from the user object
//
// options's keys can be "email", "fullname", "passwordupdate", "authservice" OR Nothing
func (u *User) Sanitize(options map[string]bool) {
	u.Password = ""
	u.AuthData = model.NewString("")
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
		u.AuthData = model.NewString("")
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
	u.AuthData = model.NewString("")
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

	if nameFormat == model.SHOW_NICKNAME_FULLNAME {
		if u.Nickname != "" {
			displayName = u.Nickname
		} else if fullName := u.GetFullName(); fullName != "" {
			displayName = fullName
		}
	} else if nameFormat == model.SHOW_FULLNAME {
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
		if !model.IsValidRoleName(r) {
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
	return IsInRole(u.Roles, model.SYSTEM_GUEST_ROLE_ID)
}

// IsSystemAdmin checks if user's roles contains "system_admin"
func (u *User) IsSystemAdmin() bool {
	return IsInRole(u.Roles, model.SYSTEM_ADMIN_ROLE_ID)
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
	return u.AuthService == model.SERVICE_GOOGLE || u.AuthService == model.SERVICE_OPENID
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
	model.ModelFromJson(&user, data)
	return user
}

func UserPatchFromJson(data io.Reader) *UserPatch {
	var user UserPatch
	model.ModelFromJson(&user, data)
	return &user
}

func UserAuthFromJson(data io.Reader) *UserAuth {
	var user UserAuth
	model.ModelFromJson(&user, data)
	return &user
}

func UserMapToJson(u map[string]*User) string {
	return model.ModelToJson(&u)
}

func UserMapFromJson(data io.Reader) map[string]*User {
	var users map[string]*User
	model.ModelFromJson(&users, data)
	return users
}

func UserListToJson(u []*User) string {
	return model.ModelToJson(&u)
}

func UserListFromJson(data io.Reader) []*User {
	var users []*User
	model.ModelFromJson(&users, data)
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

// ComparePassword checks if the hash and given password are matches
func ComparePassword(hash string, password string) bool {

	if password == "" || hash == "" {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
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
