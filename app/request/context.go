package request

import (
	"context"
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/i18n"
)

type Context struct {
	t              i18n.TranslateFunc
	session        model.Session
	requestId      string
	ipAddress      string
	path           string
	userAgent      string
	acceptLanguage string
	context        context.Context
	request        *http.Request
}

func NewContext(ctx context.Context, requestId, ipAddress, path, userAgent, acceptLanguage string, session model.Session, t i18n.TranslateFunc) *Context {
	return &Context{
		t:              t,
		session:        session,
		requestId:      requestId,
		ipAddress:      ipAddress,
		path:           path,
		userAgent:      userAgent,
		acceptLanguage: acceptLanguage,
		context:        ctx,
	}
}

func EmptyContext() *Context {
	return &Context{
		t:       i18n.T,
		context: context.Background(),
	}
}

func (c *Context) T(translationID string, args ...any) string {
	return c.t(translationID, args...)
}
func (c *Context) Session() *model.Session {
	return &c.session
}
func (c *Context) RequestId() string {
	return c.requestId
}
func (c *Context) IpAddress() string {
	return c.ipAddress
}
func (c *Context) Path() string {
	return c.path
}
func (c *Context) UserAgent() string {
	return c.userAgent
}
func (c *Context) AcceptLanguage() string {
	return c.acceptLanguage
}

func (c *Context) Context() context.Context {
	return c.context
}

// SetSession sets c's session to given session
func (c *Context) SetSession(s *model.Session) {
	c.session = *s
}

func (c *Context) SetT(t i18n.TranslateFunc) {
	c.t = t
}
func (c *Context) SetRequestId(s string) {
	c.requestId = s
}
func (c *Context) SetIpAddress(s string) {
	c.ipAddress = s
}
func (c *Context) SetUserAgent(s string) {
	c.userAgent = s
}
func (c *Context) SetAcceptLanguage(s string) {
	c.acceptLanguage = s
}
func (c *Context) SetPath(s string) {
	c.path = s
}
func (c *Context) SetContext(ctx context.Context) {
	c.context = ctx
}

func (c *Context) GetT() i18n.TranslateFunc {
	return c.t
}

func (c *Context) SetRequest(r *http.Request) {
	c.request = r
}

func (c *Context) GetRequest() *http.Request {
	return c.request
}
