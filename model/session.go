package model

import (
	"io"
	"strconv"
	"strings"

	"github.com/sitename/sitename/modules/json"
	"github.com/sitename/sitename/modules/slog"
)

const (
	SESSION_COOKIE_TOKEN              = "MMAUTHTOKEN"
	SESSION_COOKIE_USER               = "MMUSERID"
	SESSION_COOKIE_CSRF               = "MMCSRF"
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
	USER_AUTH_SERVICE_SAML            = "saml"
	USER_AUTH_SERVICE_SAML_TEXT       = "SAML"
	USER_AUTH_SERVICE_IS_SAML         = "isSaml"
	USER_AUTH_SERVICE_IS_MOBILE       = "isMobile"
	USER_AUTH_SERVICE_IS_OAUTH        = "isOAuthUser"
)

// Session contains the user session details.
// This struct's serializer methods are auto-generated. If a new field is added/removed,
// please run make gen-serialized.
type Session struct {
	Id             string    `json:"id"`
	Token          string    `json:"token"`
	CreateAt       int64     `json:"create_at"`
	ExpiresAt      int64     `json:"expires_at"`
	LastActivityAt int64     `json:"last_activity_at"`
	UserId         string    `json:"user_id"`
	DeviceId       string    `json:"device_id"`
	Roles          string    `json:"roles"`
	IsOAuth        bool      `json:"is_oauth"`
	ExpiredNotify  bool      `json:"expired_notify"`
	Props          StringMap `json:"props"`
	Local          bool      `json:"local" db:"-"`
	// TeamMembers    []*TeamMember `json:"team_members" db:"-"`
}

// Returns true if the session is unrestricted, which should grant it
// with all permissions. This is used for local mode sessions
func (s *Session) IsUnrestricted() bool {
	return s.Local
}

func (s *Session) DeepCopy() *Session {
	copySession := *s

	if s.Props != nil {
		copySession.Props = CopyStringMap(s.Props)
	}

	return &copySession
}

func (s *Session) ToJson() string {
	b, _ := json.JSON.Marshal(s)
	return string(b)
}

func SessionFromJson(data io.Reader) *Session {
	var s *Session
	json.JSON.NewDecoder(data).Decode(s)
	return s
}

func (s *Session) PreSave() {
	if s.Id == "" {
		s.Id = NewId()
	}

	if s.Token == "" {
		s.Token = NewId()
	}

	s.CreateAt = GetMillis()
	s.LastActivityAt = s.CreateAt

	if s.Props == nil {
		s.Props = make(map[string]string)
	}
}

func (s *Session) Sanitize() {
	s.Token = ""
}

func (s *Session) IsExpired() bool {

	if s.ExpiresAt <= 0 {
		return false
	}

	if GetMillis() > s.ExpiresAt {
		return true
	}

	return false
}

func (s *Session) AddProp(key string, value string) {

	if s.Props == nil {
		s.Props = make(map[string]string)
	}

	s.Props[key] = value
}

func (s *Session) IsMobileApp() bool {
	return s.DeviceId != "" || s.IsMobile()
}

func (s *Session) IsMobile() bool {
	val, ok := s.Props[USER_AUTH_SERVICE_IS_MOBILE]
	if !ok {
		return false
	}
	isMobile, err := strconv.ParseBool(val)
	if err != nil {
		slog.Debug("Error parsing boolean property from Session", slog.Err(err))
		return false
	}
	return isMobile
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
	return s.IsOAuthUser()
}

func (s *Session) GetUserRoles() []string {
	return strings.Fields(s.Roles)
}

func (s *Session) GenerateCSRF() string {
	token := NewId()
	s.AddProp("csrf", token)
	return token
}

func (s *Session) GetCSRF() string {
	if s.Props == nil {
		return ""
	}

	return s.Props["csrf"]
}

func SessionsToJson(o []*Session) string {
	b, err := json.JSON.Marshal(o)
	if err != nil {
		return "[]"
	}
	return string(b)
}

func SessionsFromJson(data io.Reader) []*Session {
	var o []*Session
	json.JSON.NewDecoder(data).Decode(&o)
	return o
}
