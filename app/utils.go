package app

import (
	"net/http"
	"strings"

	"github.com/sitename/sitename/app/request"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/plugin"
)

type TokenLocation int

const (
	TokenLocationNotFound TokenLocation = iota
	TokenLocationHeader
	TokenLocationCookie
	TokenLocationQueryString
	TokenLocationCloudHeader
	TokenLocationRemoteClusterHeader
)

// String implements fmt Stringer interface
func (tl TokenLocation) String() string {
	switch tl {
	case TokenLocationNotFound:
		return "Not Found"
	case TokenLocationHeader:
		return "Header"
	case TokenLocationCookie:
		return "Cookie"
	case TokenLocationQueryString:
		return "QueryString"
	case TokenLocationCloudHeader:
		return "CloudHeader"
	case TokenLocationRemoteClusterHeader:
		return "RemoteClusterHeader"
	default:
		return "Unknown"
	}
}

// ParseAuthTokenFromRequest reads header "Authorization" from request's header, then parses it into token and token location
func ParseAuthTokenFromRequest(r *http.Request) (string, TokenLocation) {
	authHeader := r.Header.Get(model.HEADER_AUTH)

	// Attempt to parse the token from the cookie
	if cookie, err := r.Cookie(model.SESSION_COOKIE_TOKEN); err == nil {
		return cookie.Value, TokenLocationCookie
	}

	// Parse the token from the header
	if len(authHeader) > 6 && strings.ToUpper(authHeader[0:6]) == model.HEADER_BEARER {
		// Default session token
		return authHeader[7:], TokenLocationHeader
	}

	if len(authHeader) > 5 && strings.ToLower(authHeader[0:5]) == model.HEADER_TOKEN {
		// OAuth token
		return authHeader[6:], TokenLocationHeader
	}

	// Attempt to parse token out of the query string
	if token := r.URL.Query().Get("access_token"); token != "" {
		return token, TokenLocationQueryString
	}

	if token := r.Header.Get(model.HEADER_CLOUD_TOKEN); token != "" {
		return token, TokenLocationCloudHeader
	}

	if token := r.Header.Get(model.HEADER_REMOTECLUSTER_TOKEN); token != "" {
		return token, TokenLocationRemoteClusterHeader
	}

	return "", TokenLocationNotFound
}

// PluginContext
func PluginContext(c *request.Context) *plugin.Context {
	return &plugin.Context{
		RequestId:      c.RequestId(),
		SessionId:      c.Session().Id,
		IpAddress:      c.IpAddress(),
		AcceptLanguage: c.AcceptLanguage(),
		UserAgent:      c.UserAgent(),
	}
}
