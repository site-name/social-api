package model

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
	"gorm.io/gorm"
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

type StringMap = StringMAP

// Session contains the user session details.
// This struct's serializer methods are auto-generated. If a new field is added/removed,
// please run make gen-serialized.
type Session struct {
	Id             UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	Token          string    `json:"token" gorm:"type:uuid;default:gen_random_uuid();column:Token;index:token_key"` // index, uuid
	CreateAt       int64     `json:"create_at" gorm:"type:bigint;autoCreateTime:milli;column:CreateAt"`
	ExpiresAt      int64     `json:"expires_at" gorm:"type:bigint;column:ExpiresAt"`
	LastActivityAt int64     `json:"last_activity_at" gorm:"type:bigint;column:LastActivityAt;autoUpdateTime:milli;autoCreateTime:milli"`
	UserId         UUID      `json:"user_id" gorm:"type:uuid;column:UserId;index:userid_key"` // uuid, index
	DeviceId       string    `json:"device_id" gorm:"type:varchar(512);column:DeviceId"`
	Roles          string    `json:"roles" gorm:"type:varchar(256);column:Roles"`
	IsOAuth        bool      `json:"is_oauth" gorm:"column:IsOAuth"`
	ExpiredNotify  bool      `json:"expired_notify" gorm:"column:ExpiredNotify"`
	Props          StringMap `json:"props" gorm:"type:jsonb;column:Props"`
	Local          bool      `json:"local" gorm:"-"` // this field is populated at some point
}

func (s *Session) BeforeCreate(_ *gorm.DB) error { s.commonPre(); return s.IsValid() }
func (s *Session) BeforeUpdate(_ *gorm.DB) error { s.commonPre(); return s.IsValid() }
func (s *Session) TableName() string             { return SessionTableName }

// Returns true if the session is unrestricted, which should grant it
// with all permissions. This is used for local mode sessions
func (s *Session) IsUnrestricted() bool {
	return s.Local
}

func (s *Session) DeepCopy() *Session {
	copySession := *s

	if s.Props != nil {
		copySession.Props = s.Props.DeepCopy()
	}

	return &copySession
}

func (s *Session) ToJSON() string {
	return ModelToJson(s)
}

func (s *Session) commonPre() {
	if s.Props == nil {
		s.Props = make(map[string]string)
	}
}

func (s *Session) IsValid() *AppError {
	if !IsValidId(s.UserId) {
		return NewAppError("Session.IsValid", "model.session.is_valid.user_id.app_error", nil, "", http.StatusBadRequest)
	}
	return nil
}

// Sanitize sets session's Token to empty string
func (s *Session) Sanitize() {
	s.Token = ""
}

// IsExpired checks if GetMillis() > s.ExpiresAt
func (s *Session) IsExpired() bool {
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
//				   cases where the new ExpiresAt is not relative to CreateAt.
func (s *Session) SetExpireInDays(days int) {
	if s.CreateAt == 0 {
		s.ExpiresAt = GetMillis() + (1000 * 60 * 60 * 24 * int64(days))
	} else {
		s.ExpiresAt = s.CreateAt + (1000 * 60 * 60 * 24 * int64(days))
	}
}

// AddProp adds given value to session's Props with key of given key
func (s *Session) AddProp(key string, value string) {
	if s.Props == nil {
		s.Props = make(map[string]string)
	}

	s.Props[key] = value
}

// IsMobileApp check if current session has non-empty `DeviceId` field or IsMobile()
func (s *Session) IsMobileApp() bool {
	return s.DeviceId != "" || s.IsMobile()
}

// IsMobile checks if the Props field has an item is ("isMobile": "true")
func (s *Session) IsMobile() bool {
	val, ok := s.Props[USER_AUTH_SERVICE_IS_MOBILE]
	return ok && strings.EqualFold(val, "true")
}

func (s *Session) IsSaml() bool {
	val, ok := s.Props[USER_AUTH_SERVICE_IS_SAML]
	if !ok {
		return false
	}
	isSaml, err := strconv.ParseBool(val)
	if err != nil {
		slog.Debug("Error parsing boolean property from Session", slog.Err(err))
		return false
	}
	return isSaml
}

func (s *Session) IsOAuthUser() bool {
	val, ok := s.Props[USER_AUTH_SERVICE_IS_OAUTH]
	if !ok {
		return false
	}
	isOAuthUser, err := strconv.ParseBool(val)
	if err != nil {
		slog.Debug("Error parsing boolean property from Session", slog.Err(err))
		return false
	}
	return isOAuthUser
}

func (s *Session) IsSSOLogin() bool {
	return s.IsOAuthUser() || s.IsSaml()
}

// GetUserRoles turns current session's Roles into a slice of strings
func (s *Session) GetUserRoles() util.AnyArray[string] {
	return strings.Fields(s.Roles)
}

// GenerateCSRF simply generates new UUID, then add that uuid to its "Props" with key is "csrf". Finally returns that token
func (s *Session) GenerateCSRF() string {
	token := NewId()
	s.AddProp("csrf", token)
	return token
}

// get value with key of "csrf" from session's Props
func (s *Session) GetCSRF() string {
	if s.Props == nil {
		return ""
	}

	return s.Props["csrf"]
}

func SessionsToJson(o []*Session) string {
	return ModelToJson(o)
}
