package opentracing

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"io"
	"mime/multipart"
	"net/http"
	"reflect"
	"time"

	"github.com/opentracing/opentracing-go/ext"
	spanlog "github.com/opentracing/opentracing-go/log"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/request"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	modelAudit "github.com/sitename/sitename/model/audit"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/modules/audit"
	"github.com/sitename/sitename/modules/filestore"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/timezones"
	"github.com/sitename/sitename/services/httpservice"
	"github.com/sitename/sitename/services/imageproxy"
	"github.com/sitename/sitename/services/searchengine"
	"github.com/sitename/sitename/services/tracing"
	"github.com/sitename/sitename/store"
)

type OpenTracingAppLayer struct {
	app app.AppIface

	srv *app.Server

	log              *slog.Logger
	notificationsLog *slog.Logger

	cluster      einterfaces.ClusterInterface
	compliance   einterfaces.ComplianceInterface
	searchEngine *searchengine.Broker
	ldap         einterfaces.LdapInterface
	metrics      einterfaces.MetricsInterface
	httpService  httpservice.HTTPService
	imageProxy   *imageproxy.ImageProxy
	timezones    *timezones.Timezones
	// notification     einterfaces.NotificationInterface
	saml einterfaces.SamlInterface
	// messageExport    einterfaces.MessageExportInterface
	dataRetention    einterfaces.DataRetentionInterface
	accountMigration einterfaces.AccountMigrationInterface

	ctx context.Context
}

func (a *OpenTracingAppLayer) Account() sub_app_iface.AccountApp {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.Account")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.Account()

	return resultVar0
}

func (a *OpenTracingAppLayer) ActivateMfa(userID string, token string) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.ActivateMfa")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.ActivateMfa(userID, token)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) AddConfigListener(listener func(*model.Config, *model.Config)) string {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.AddConfigListener")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.AddConfigListener(listener)

	return resultVar0
}

func (a *OpenTracingAppLayer) AddSessionToCache(s *model.Session) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.AddSessionToCache")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	a.app.AddSessionToCache(s)
}

func (a *OpenTracingAppLayer) AdjustImage(file io.Reader) (*bytes.Buffer, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.AdjustImage")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.AdjustImage(file)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) AppendFile(fr io.Reader, path string) (int64, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.AppendFile")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.AppendFile(fr, path)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) AsymmetricSigningKey() *ecdsa.PrivateKey {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.AsymmetricSigningKey")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.AsymmetricSigningKey()

	return resultVar0
}

func (a *OpenTracingAppLayer) AttachDeviceId(sessionID string, deviceID string, expiresAt int64) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.AttachDeviceId")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.AttachDeviceId(sessionID, deviceID, expiresAt)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) AttachSessionCookies(c *request.Context, w http.ResponseWriter, r *http.Request) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.AttachSessionCookies")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	a.app.AttachSessionCookies(c, w, r)
}

func (a *OpenTracingAppLayer) Attribute() sub_app_iface.AttributeApp {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.Attribute")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.Attribute()

	return resultVar0
}

func (a *OpenTracingAppLayer) AuthenticateUserForLogin(c *request.Context, id string, loginId string, password string, mfaToken string, cwsToken string, ldapOnly bool) (user *account.User, err *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.AuthenticateUserForLogin")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.AuthenticateUserForLogin(c, id, loginId, password, mfaToken, cwsToken, ldapOnly)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) Channel() sub_app_iface.ChannelApp {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.Channel")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.Channel()

	return resultVar0
}

func (a *OpenTracingAppLayer) CheckForClientSideCert(r *http.Request) (string, string, string) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.CheckForClientSideCert")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1, resultVar2 := a.app.CheckForClientSideCert(r)

	return resultVar0, resultVar1, resultVar2
}

func (a *OpenTracingAppLayer) CheckMandatoryS3Fields(settings *model.FileSettings) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.CheckMandatoryS3Fields")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.CheckMandatoryS3Fields(settings)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) CheckPasswordAndAllCriteria(user *account.User, password string, mfaToken string) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.CheckPasswordAndAllCriteria")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.CheckPasswordAndAllCriteria(user, password, mfaToken)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) CheckProviderAttributes(user *account.User, patch *account.UserPatch) string {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.CheckProviderAttributes")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.CheckProviderAttributes(user, patch)

	return resultVar0
}

func (a *OpenTracingAppLayer) CheckRolesExist(roleNames []string) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.CheckRolesExist")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.CheckRolesExist(roleNames)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) CheckUserAllAuthenticationCriteria(user *account.User, mfaToken string) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.CheckUserAllAuthenticationCriteria")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.CheckUserAllAuthenticationCriteria(user, mfaToken)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) CheckUserMfa(user *account.User, token string) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.CheckUserMfa")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.CheckUserMfa(user, token)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) CheckUserPostflightAuthenticationCriteria(user *account.User) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.CheckUserPostflightAuthenticationCriteria")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.CheckUserPostflightAuthenticationCriteria(user)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) CheckUserPreflightAuthenticationCriteria(user *account.User, mfaToken string) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.CheckUserPreflightAuthenticationCriteria")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.CheckUserPreflightAuthenticationCriteria(user, mfaToken)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) Checkout() sub_app_iface.CheckoutApp {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.Checkout")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.Checkout()

	return resultVar0
}

func (a *OpenTracingAppLayer) ClearSessionCacheForUser(userID string) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.ClearSessionCacheForUser")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	a.app.ClearSessionCacheForUser(userID)
}

func (a *OpenTracingAppLayer) ClearSessionCacheForUserSkipClusterSend(userID string) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.ClearSessionCacheForUserSkipClusterSend")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	a.app.ClearSessionCacheForUserSkipClusterSend(userID)
}

func (a *OpenTracingAppLayer) ClientConfig() map[string]string {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.ClientConfig")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.ClientConfig()

	return resultVar0
}

func (a *OpenTracingAppLayer) ClientConfigHash() string {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.ClientConfigHash")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.ClientConfigHash()

	return resultVar0
}

func (a *OpenTracingAppLayer) ClientConfigWithComputed() map[string]string {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.ClientConfigWithComputed")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.ClientConfigWithComputed()

	return resultVar0
}

func (a *OpenTracingAppLayer) Config() *model.Config {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.Config")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.Config()

	return resultVar0
}

func (a *OpenTracingAppLayer) CopyFileInfos(userID string, fileIDs []string) ([]string, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.CopyFileInfos")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.CopyFileInfos(userID, fileIDs)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) CreateGuest(c *request.Context, user *account.User) (*account.User, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.CreateGuest")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.CreateGuest(c, user)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) CreatePasswordRecoveryToken(userID string, email string) (*model.Token, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.CreatePasswordRecoveryToken")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.CreatePasswordRecoveryToken(userID, email)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) CreateRole(role *model.Role) (*model.Role, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.CreateRole")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.CreateRole(role)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) CreateSession(session *model.Session) (*model.Session, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.CreateSession")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.CreateSession(session)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) CreateUploadSession(us *model.UploadSession) (*model.UploadSession, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.CreateUploadSession")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.CreateUploadSession(us)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) CreateUser(c *request.Context, user *account.User) (*account.User, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.CreateUser")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.CreateUser(c, user)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) CreateUserAccessToken(token *account.UserAccessToken) (*account.UserAccessToken, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.CreateUserAccessToken")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.CreateUserAccessToken(token)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) CreateUserAsAdmin(c *request.Context, user *account.User, redirect string) (*account.User, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.CreateUserAsAdmin")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.CreateUserAsAdmin(c, user, redirect)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) CreateUserFromSignup(c *request.Context, user *account.User, redirect string) (*account.User, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.CreateUserFromSignup")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.CreateUserFromSignup(c, user, redirect)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) CreateUserWithToken(c *request.Context, user *account.User, token *model.Token) (*account.User, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.CreateUserWithToken")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.CreateUserWithToken(c, user, token)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) CreateZipFileAndAddFiles(fileBackend filestore.FileBackend, fileDatas []model.FileData, zipFileName string, directory string) error {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.CreateZipFileAndAddFiles")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.CreateZipFileAndAddFiles(fileBackend, fileDatas, zipFileName, directory)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) Csv() sub_app_iface.CsvApp {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.Csv")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.Csv()

	return resultVar0
}

func (a *OpenTracingAppLayer) DBHealthCheckDelete() error {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.DBHealthCheckDelete")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.DBHealthCheckDelete()

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) DBHealthCheckWrite() error {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.DBHealthCheckWrite")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.DBHealthCheckWrite()

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) DeactivateGuests(c *request.Context) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.DeactivateGuests")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.DeactivateGuests(c)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) DeactivateMfa(userID string) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.DeactivateMfa")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.DeactivateMfa(userID)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) DeleteToken(token *model.Token) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.DeleteToken")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.DeleteToken(token)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) DisableUserAccessToken(token *account.UserAccessToken) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.DisableUserAccessToken")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.DisableUserAccessToken(token)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) DoAdvancedPermissionsMigration() {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.DoAdvancedPermissionsMigration")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	a.app.DoAdvancedPermissionsMigration()
}

func (a *OpenTracingAppLayer) DoAppMigrations() {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.DoAppMigrations")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	a.app.DoAppMigrations()
}

func (a *OpenTracingAppLayer) DoLogin(c *request.Context, w http.ResponseWriter, r *http.Request, user *account.User, deviceID string, isMobile bool, isOAuthUser bool, isSaml bool) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.DoLogin")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.DoLogin(c, w, r, user, deviceID, isMobile, isOAuthUser, isSaml)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) DoPermissionsMigrations() error {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.DoPermissionsMigrations")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.DoPermissionsMigrations()

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) DoSystemConsoleRolesCreationMigration() {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.DoSystemConsoleRolesCreationMigration")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	a.app.DoSystemConsoleRolesCreationMigration()
}

func (a *OpenTracingAppLayer) DoUploadFile(c *request.Context, now time.Time, rawTeamId string, rawChannelId string, rawUserId string, rawFilename string, data []byte) (*model.FileInfo, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.DoUploadFile")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.DoUploadFile(c, now, rawTeamId, rawChannelId, rawUserId, rawFilename, data)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) DoUploadFileExpectModification(c *request.Context, now time.Time, rawTeamId string, rawChannelId string, rawUserId string, rawFilename string, data []byte) (*model.FileInfo, []byte, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.DoUploadFileExpectModification")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1, resultVar2 := a.app.DoUploadFileExpectModification(c, now, rawTeamId, rawChannelId, rawUserId, rawFilename, data)

	if resultVar2 != nil {
		span.LogFields(spanlog.Error(resultVar2))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1, resultVar2
}

func (a *OpenTracingAppLayer) DoubleCheckPassword(user *account.User, password string) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.DoubleCheckPassword")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.DoubleCheckPassword(user, password)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) EnableUserAccessToken(token *account.UserAccessToken) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.EnableUserAccessToken")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.EnableUserAccessToken(token)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) EnvironmentConfig(filter func(reflect.StructField) bool) map[string]interface{} {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.EnvironmentConfig")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.EnvironmentConfig(filter)

	return resultVar0
}

func (a *OpenTracingAppLayer) ExportPermissions(w io.Writer) error {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.ExportPermissions")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.ExportPermissions(w)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) ExtendSessionExpiryIfNeeded(session *model.Session) bool {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.ExtendSessionExpiryIfNeeded")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.ExtendSessionExpiryIfNeeded(session)

	return resultVar0
}

func (a *OpenTracingAppLayer) ExtractContentFromFileInfo(fileInfo *model.FileInfo) error {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.ExtractContentFromFileInfo")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.ExtractContentFromFileInfo(fileInfo)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) FileBackend() (filestore.FileBackend, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.FileBackend")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.FileBackend()

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) FileExists(path string) (bool, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.FileExists")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.FileExists(path)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) FileModTime(path string) (time.Time, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.FileModTime")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.FileModTime(path)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) FileReader(path string) (filestore.ReadCloseSeeker, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.FileReader")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.FileReader(path)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) FileSize(path string) (int64, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.FileSize")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.FileSize(path)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GenerateMfaSecret(userID string) (*model.MfaSecret, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GenerateMfaSecret")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GenerateMfaSecret(userID)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GeneratePublicLink(siteURL string, info *model.FileInfo) string {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GeneratePublicLink")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.GeneratePublicLink(siteURL, info)

	return resultVar0
}

func (a *OpenTracingAppLayer) GetAudits(userID string, limit int) (modelAudit.Audits, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetAudits")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetAudits(userID, limit)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetAuditsPage(userID string, page int, perPage int) (modelAudit.Audits, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetAuditsPage")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetAuditsPage(userID, page, perPage)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetCloudSession(token string) (*model.Session, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetCloudSession")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetCloudSession(token)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetClusterId() string {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetClusterId")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.GetClusterId()

	return resultVar0
}

func (a *OpenTracingAppLayer) GetConfigFile(name string) ([]byte, error) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetConfigFile")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetConfigFile(name)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetCookieDomain() string {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetCookieDomain")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.GetCookieDomain()

	return resultVar0
}

func (a *OpenTracingAppLayer) GetDefaultProfileImage(user *account.User) ([]byte, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetDefaultProfileImage")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetDefaultProfileImage(user)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetEnvironmentConfig(filter func(reflect.StructField) bool) map[string]interface{} {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetEnvironmentConfig")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.GetEnvironmentConfig(filter)

	return resultVar0
}

func (a *OpenTracingAppLayer) GetFile(fileID string) ([]byte, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetFile")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetFile(fileID)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetFileInfo(fileID string) (*model.FileInfo, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetFileInfo")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetFileInfo(fileID)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetFileInfos(page int, perPage int, opt *model.GetFileInfosOptions) ([]*model.FileInfo, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetFileInfos")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetFileInfos(page, perPage, opt)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetFilteredUsersStats(options *account.UserCountOptions) (*account.UsersStats, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetFilteredUsersStats")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetFilteredUsersStats(options)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetPasswordRecoveryToken(token string) (*model.Token, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetPasswordRecoveryToken")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetPasswordRecoveryToken(token)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetProfileImage(user *account.User) ([]byte, bool, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetProfileImage")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1, resultVar2 := a.app.GetProfileImage(user)

	if resultVar2 != nil {
		span.LogFields(spanlog.Error(resultVar2))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1, resultVar2
}

func (a *OpenTracingAppLayer) GetRole(id string) (*model.Role, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetRole")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetRole(id)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetRoleByName(ctx context.Context, name string) (*model.Role, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetRoleByName")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetRoleByName(ctx, name)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetRolesByNames(names []string) ([]*model.Role, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetRolesByNames")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetRolesByNames(names)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetSanitizeOptions(asAdmin bool) map[string]bool {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetSanitizeOptions")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.GetSanitizeOptions(asAdmin)

	return resultVar0
}

func (a *OpenTracingAppLayer) GetSanitizedConfig() *model.Config {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetSanitizedConfig")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.GetSanitizedConfig()

	return resultVar0
}

func (a *OpenTracingAppLayer) GetSession(token string) (*model.Session, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetSession")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetSession(token)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetSessionById(sessionID string) (*model.Session, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetSessionById")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetSessionById(sessionID)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetSessionLengthInMillis(session *model.Session) int64 {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetSessionLengthInMillis")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.GetSessionLengthInMillis(session)

	return resultVar0
}

func (a *OpenTracingAppLayer) GetSessions(userID string) ([]*model.Session, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetSessions")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetSessions(userID)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetSiteURL() string {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetSiteURL")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.GetSiteURL()

	return resultVar0
}

func (a *OpenTracingAppLayer) GetStatus(userID string) (*model.Status, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetStatus")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetStatus(userID)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetStatusFromCache(userID string) *model.Status {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetStatusFromCache")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.GetStatusFromCache(userID)

	return resultVar0
}

func (a *OpenTracingAppLayer) GetTotalUsersStats() (*account.UsersStats, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetTotalUsersStats")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetTotalUsersStats()

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetUploadSession(uploadId string) (*model.UploadSession, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetUploadSession")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetUploadSession(uploadId)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetUploadSessionsForUser(userID string) ([]*model.UploadSession, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetUploadSessionsForUser")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetUploadSessionsForUser(userID)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetUser(userID string) (*account.User, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetUser")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetUser(userID)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetUserAccessToken(tokenID string, sanitize bool) (*account.UserAccessToken, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetUserAccessToken")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetUserAccessToken(tokenID, sanitize)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetUserAccessTokens(page int, perPage int) ([]*account.UserAccessToken, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetUserAccessTokens")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetUserAccessTokens(page, perPage)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetUserAccessTokensForUser(userID string, page int, perPage int) ([]*account.UserAccessToken, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetUserAccessTokensForUser")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetUserAccessTokensForUser(userID, page, perPage)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetUserByAuth(authData *string, authService string) (*account.User, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetUserByAuth")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetUserByAuth(authData, authService)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetUserByEmail(email string) (*account.User, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetUserByEmail")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetUserByEmail(email)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetUserByUsername(username string) (*account.User, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetUserByUsername")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetUserByUsername(username)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetUserForLogin(id string, loginId string) (*account.User, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetUserForLogin")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetUserForLogin(id, loginId)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetUsers(options *account.UserGetOptions) ([]*account.User, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetUsers")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetUsers(options)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetUsersByIds(userIDs []string, options *store.UserGetByIdsOpts) ([]*account.User, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetUsersByIds")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetUsersByIds(userIDs, options)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetUsersByUsernames(usernames []string, asAdmin bool) ([]*account.User, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetUsersByUsernames")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetUsersByUsernames(usernames, asAdmin)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetVerifyEmailToken(token string) (*model.Token, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetVerifyEmailToken")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetVerifyEmailToken(token)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetWarnMetricsStatus() (map[string]*model.WarnMetricStatus, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetWarnMetricsStatus")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetWarnMetricsStatus()

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) Giftcard() sub_app_iface.GiftcardApp {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.Giftcard")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.Giftcard()

	return resultVar0
}

func (a *OpenTracingAppLayer) Handle404(w http.ResponseWriter, r *http.Request) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.Handle404")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	a.app.Handle404(w, r)
}

func (a *OpenTracingAppLayer) HandleImages(previewPathList []string, thumbnailPathList []string, fileData [][]byte) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.HandleImages")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	a.app.HandleImages(previewPathList, thumbnailPathList, fileData)
}

func (a *OpenTracingAppLayer) HandleMessageExportConfig(cfg *model.Config, appCfg *model.Config) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.HandleMessageExportConfig")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	a.app.HandleMessageExportConfig(cfg, appCfg)
}

func (a *OpenTracingAppLayer) HasPermissionTo(askingUserId string, permission *model.Permission) bool {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.HasPermissionTo")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.HasPermissionTo(askingUserId, permission)

	return resultVar0
}

func (a *OpenTracingAppLayer) HasPermissionToUser(askingUserId string, userID string) bool {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.HasPermissionToUser")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.HasPermissionToUser(askingUserId, userID)

	return resultVar0
}

func (a *OpenTracingAppLayer) InvalidateCacheForUser(userID string) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.InvalidateCacheForUser")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	a.app.InvalidateCacheForUser(userID)
}

func (a *OpenTracingAppLayer) Invoice() sub_app_iface.InvoiceApp {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.Invoice")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.Invoice()

	return resultVar0
}

func (a *OpenTracingAppLayer) IsFirstUserAccount() bool {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.IsFirstUserAccount")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.IsFirstUserAccount()

	return resultVar0
}

func (a *OpenTracingAppLayer) IsLeader() bool {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.IsLeader")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.IsLeader()

	return resultVar0
}

func (a *OpenTracingAppLayer) IsPasswordValid(password string) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.IsPasswordValid")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.IsPasswordValid(password)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) IsUserSignUpAllowed() *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.IsUserSignUpAllowed")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.IsUserSignUpAllowed()

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) IsUserSignupAllowed() *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.IsUserSignupAllowed")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.IsUserSignupAllowed()

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) IsUsernameTaken(name string) bool {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.IsUsernameTaken")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.IsUsernameTaken(name)

	return resultVar0
}

func (a *OpenTracingAppLayer) LimitedClientConfig() map[string]string {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.LimitedClientConfig")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.LimitedClientConfig()

	return resultVar0
}

func (a *OpenTracingAppLayer) LimitedClientConfigWithComputed() map[string]string {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.LimitedClientConfigWithComputed")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.LimitedClientConfigWithComputed()

	return resultVar0
}

func (a *OpenTracingAppLayer) ListDirectory(path string) ([]string, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.ListDirectory")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.ListDirectory(path)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) LogAuditRec(rec *audit.Record, err error) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.LogAuditRec")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	a.app.LogAuditRec(rec, err)
}

func (a *OpenTracingAppLayer) LogAuditRecWithLevel(rec *audit.Record, level slog.LogLevel, err error) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.LogAuditRecWithLevel")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	a.app.LogAuditRecWithLevel(rec, level, err)
}

func (a *OpenTracingAppLayer) MakeAuditRecord(event string, initialStatus string) *audit.Record {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.MakeAuditRecord")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.MakeAuditRecord(event, initialStatus)

	return resultVar0
}

func (a *OpenTracingAppLayer) MakePermissionError(s *model.Session, permissions []*model.Permission) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.MakePermissionError")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.MakePermissionError(s, permissions)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) Menu() sub_app_iface.MenuApp {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.Menu")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.Menu()

	return resultVar0
}

func (a *OpenTracingAppLayer) MoveFile(oldPath string, newPath string) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.MoveFile")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.MoveFile(oldPath, newPath)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) NewClusterDiscoveryService() *app.ClusterDiscoveryService {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.NewClusterDiscoveryService")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.NewClusterDiscoveryService()

	return resultVar0
}

func (a *OpenTracingAppLayer) NotifyAndSetWarnMetricAck(warnMetricId string, sender *account.User, forceAck bool, isBot bool) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.NotifyAndSetWarnMetricAck")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.NotifyAndSetWarnMetricAck(warnMetricId, sender, forceAck, isBot)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) Order() sub_app_iface.OrderApp {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.Order")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.Order()

	return resultVar0
}

func (a *OpenTracingAppLayer) OriginChecker() func(*http.Request) bool {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.OriginChecker")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.OriginChecker()

	return resultVar0
}

func (a *OpenTracingAppLayer) Page() sub_app_iface.PageApp {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.Page")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.Page()

	return resultVar0
}

func (a *OpenTracingAppLayer) PatchRole(role *model.Role, patch *model.RolePatch) (*model.Role, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.PatchRole")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.PatchRole(role, patch)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) Payment() sub_app_iface.PaymentApp {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.Payment")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.Payment()

	return resultVar0
}

func (a *OpenTracingAppLayer) PermanentDeleteAllUsers(c *request.Context) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.PermanentDeleteAllUsers")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.PermanentDeleteAllUsers(c)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) PermanentDeleteUser(c *request.Context, user *account.User) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.PermanentDeleteUser")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.PermanentDeleteUser(c, user)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) PostActionCookieSecret() []byte {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.PostActionCookieSecret")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.PostActionCookieSecret()

	return resultVar0
}

func (a *OpenTracingAppLayer) Product() sub_app_iface.ProductApp {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.Product")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.Product()

	return resultVar0
}

func (a *OpenTracingAppLayer) ProductVariantById(id string) (*product_and_discount.ProductVariant, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.ProductVariantById")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.ProductVariantById(id)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) Publish(message *model.WebSocketEvent) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.Publish")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	a.app.Publish(message)
}

func (a *OpenTracingAppLayer) ReadFile(path string) ([]byte, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.ReadFile")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.ReadFile(path)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) ReloadConfig() error {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.ReloadConfig")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.ReloadConfig()

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) RemoveConfigListener(id string) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.RemoveConfigListener")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	a.app.RemoveConfigListener(id)
}

func (a *OpenTracingAppLayer) RemoveDirectory(path string) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.RemoveDirectory")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.RemoveDirectory(path)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) RemoveFile(path string) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.RemoveFile")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.RemoveFile(path)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) ResetPasswordFromToken(userSuppliedTokenString string, newPassword string) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.ResetPasswordFromToken")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.ResetPasswordFromToken(userSuppliedTokenString, newPassword)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) ResetPermissionsSystem() *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.ResetPermissionsSystem")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.ResetPermissionsSystem()

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) RevokeAllSessions(userID string) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.RevokeAllSessions")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.RevokeAllSessions(userID)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) RevokeSession(session *model.Session) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.RevokeSession")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.RevokeSession(session)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) RevokeSessionById(sessionID string) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.RevokeSessionById")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.RevokeSessionById(sessionID)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) RevokeSessionsForDeviceId(userID string, deviceID string, currentSessionId string) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.RevokeSessionsForDeviceId")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.RevokeSessionsForDeviceId(userID, deviceID, currentSessionId)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) RevokeUserAccessToken(token *account.UserAccessToken) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.RevokeUserAccessToken")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.RevokeUserAccessToken(token)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) RolesGrantPermission(roleNames []string, permissionId string) bool {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.RolesGrantPermission")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.RolesGrantPermission(roleNames, permissionId)

	return resultVar0
}

func (a *OpenTracingAppLayer) SanitizeProfile(user *account.User, asAdmin bool) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.SanitizeProfile")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	a.app.SanitizeProfile(user, asAdmin)
}

func (a *OpenTracingAppLayer) SaveConfig(newCfg *model.Config, sendConfigChangeClusterMessage bool) (*model.Config, *model.Config, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.SaveConfig")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1, resultVar2 := a.app.SaveConfig(newCfg, sendConfigChangeClusterMessage)

	if resultVar2 != nil {
		span.LogFields(spanlog.Error(resultVar2))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1, resultVar2
}

func (a *OpenTracingAppLayer) SearchEngine() *searchengine.Broker {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.SearchEngine")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.SearchEngine()

	return resultVar0
}

func (a *OpenTracingAppLayer) SearchUserAccessTokens(term string) ([]*account.UserAccessToken, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.SearchUserAccessTokens")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.SearchUserAccessTokens(term)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) SearchUsers(props *account.UserSearch, options *account.UserSearchOptions) ([]*account.User, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.SearchUsers")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.SearchUsers(props, options)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) SendEmailVerification(user *account.User, newEmail string, redirect string) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.SendEmailVerification")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.SendEmailVerification(user, newEmail, redirect)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) SendPasswordReset(email string, siteURL string) (bool, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.SendPasswordReset")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.SendPasswordReset(email, siteURL)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) Seo() sub_app_iface.SeoApp {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.Seo")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.Seo()

	return resultVar0
}

func (a *OpenTracingAppLayer) SessionCacheLength() int {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.SessionCacheLength")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.SessionCacheLength()

	return resultVar0
}

func (a *OpenTracingAppLayer) SessionHasPermissionTo(session *model.Session, permission *model.Permission) bool {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.SessionHasPermissionTo")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.SessionHasPermissionTo(session, permission)

	return resultVar0
}

func (a *OpenTracingAppLayer) SessionHasPermissionToAny(session *model.Session, permissions []*model.Permission) bool {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.SessionHasPermissionToAny")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.SessionHasPermissionToAny(session, permissions)

	return resultVar0
}

func (a *OpenTracingAppLayer) SessionHasPermissionToUser(session *model.Session, userID string) bool {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.SessionHasPermissionToUser")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.SessionHasPermissionToUser(session, userID)

	return resultVar0
}

func (a *OpenTracingAppLayer) SetDefaultProfileImage(user *account.User) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.SetDefaultProfileImage")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.SetDefaultProfileImage(user)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) SetPhase2PermissionsMigrationStatus(isComplete bool) error {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.SetPhase2PermissionsMigrationStatus")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.SetPhase2PermissionsMigrationStatus(isComplete)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) SetProfileImage(userID string, imageData *multipart.FileHeader) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.SetProfileImage")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.SetProfileImage(userID, imageData)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) SetProfileImageFromFile(userID string, file io.Reader) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.SetProfileImageFromFile")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.SetProfileImageFromFile(userID, file)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) SetProfileImageFromMultiPartFile(userID string, file multipart.File) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.SetProfileImageFromMultiPartFile")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.SetProfileImageFromMultiPartFile(userID, file)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) SetSessionExpireInDays(session *model.Session, days int) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.SetSessionExpireInDays")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	a.app.SetSessionExpireInDays(session, days)
}

func (a *OpenTracingAppLayer) Shipping() sub_app_iface.ShippingApp {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.Shipping")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.Shipping()

	return resultVar0
}

func (a *OpenTracingAppLayer) Site() sub_app_iface.SiteApp {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.Site")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.Site()

	return resultVar0
}

func (a *OpenTracingAppLayer) TestFileStoreConnection() *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.TestFileStoreConnection")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.TestFileStoreConnection()

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) TestFileStoreConnectionWithConfig(settings *model.FileSettings) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.TestFileStoreConnectionWithConfig")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.TestFileStoreConnectionWithConfig(settings)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) UpdateActive(c *request.Context, user *account.User, active bool) (*account.User, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.UpdateActive")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.UpdateActive(c, user, active)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) UpdateConfig(f func(*model.Config)) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.UpdateConfig")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	a.app.UpdateConfig(f)
}

func (a *OpenTracingAppLayer) UpdateHashedPassword(user *account.User, newHashedPassword string) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.UpdateHashedPassword")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.UpdateHashedPassword(user, newHashedPassword)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) UpdateHashedPasswordByUserId(userID string, newHashedPassword string) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.UpdateHashedPasswordByUserId")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.UpdateHashedPasswordByUserId(userID, newHashedPassword)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) UpdateLastActivityAtIfNeeded(session model.Session) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.UpdateLastActivityAtIfNeeded")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	a.app.UpdateLastActivityAtIfNeeded(session)
}

func (a *OpenTracingAppLayer) UpdateMfa(activate bool, userID string, token string) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.UpdateMfa")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.UpdateMfa(activate, userID, token)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) UpdatePassword(user *account.User, newPassword string) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.UpdatePassword")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.UpdatePassword(user, newPassword)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) UpdatePasswordAsUser(userID string, currentPassword string, newPassword string) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.UpdatePasswordAsUser")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.UpdatePasswordAsUser(userID, currentPassword, newPassword)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) UpdatePasswordByUserIdSendEmail(userID string, newPassword string, method string) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.UpdatePasswordByUserIdSendEmail")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.UpdatePasswordByUserIdSendEmail(userID, newPassword, method)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) UpdatePasswordSendEmail(user *account.User, newPassword string, method string) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.UpdatePasswordSendEmail")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.UpdatePasswordSendEmail(user, newPassword, method)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) UpdateRole(role *model.Role) (*model.Role, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.UpdateRole")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.UpdateRole(role)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) UpdateUser(user *account.User, sendNotifications bool) (*account.User, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.UpdateUser")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.UpdateUser(user, sendNotifications)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) UpdateUserAsUser(user *account.User, asAdmin bool) (*account.User, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.UpdateUserAsUser")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.UpdateUserAsUser(user, asAdmin)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) UpdateUserAuth(userID string, userAuth *account.UserAuth) (*account.UserAuth, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.UpdateUserAuth")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.UpdateUserAuth(userID, userAuth)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) UpdateUserRoles(userID string, newRoles string, sendWebSocketEvent bool) (*account.User, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.UpdateUserRoles")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.UpdateUserRoles(userID, newRoles, sendWebSocketEvent)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) UpdateUserRolesWithUser(user *account.User, newRoles string, sendWebSocketEvent bool) (*account.User, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.UpdateUserRolesWithUser")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.UpdateUserRolesWithUser(user, newRoles, sendWebSocketEvent)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) UploadData(c *request.Context, us *model.UploadSession, rd io.Reader) (*model.FileInfo, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.UploadData")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.UploadData(c, us, rd)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) VerifyEmailFromToken(userSuppliedTokenString string) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.VerifyEmailFromToken")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.VerifyEmailFromToken(userSuppliedTokenString)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) VerifyUserEmail(userID string, email string) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.VerifyUserEmail")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.VerifyUserEmail(userID, email)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
}

func (a *OpenTracingAppLayer) Warehouse() sub_app_iface.WarehouseApp {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.Warehouse")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.Warehouse()

	return resultVar0
}

func (a *OpenTracingAppLayer) Webhook() sub_app_iface.WebhookApp {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.Webhook")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.Webhook()

	return resultVar0
}

func (a *OpenTracingAppLayer) Wishlist() sub_app_iface.WishlistApp {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.Wishlist")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.Wishlist()

	return resultVar0
}

func (a *OpenTracingAppLayer) WriteFile(fr io.Reader, path string) (int64, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.WriteFile")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.WriteFile(fr, path)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func NewOpenTracingAppLayer(childApp app.AppIface, ctx context.Context) *OpenTracingAppLayer {
	newApp := OpenTracingAppLayer{
		app: childApp,
		ctx: ctx,
	}

	newApp.srv = childApp.Srv()
	newApp.log = childApp.Log()
	newApp.notificationsLog = childApp.NotificationsLog()
	newApp.accountMigration = childApp.AccountMigration()
	newApp.cluster = childApp.Cluster()
	newApp.compliance = childApp.Compliance()
	newApp.dataRetention = childApp.DataRetention()
	newApp.searchEngine = childApp.SearchEngine()
	newApp.ldap = childApp.Ldap()
	// newApp.messageExport = childApp.MessageExport()
	newApp.metrics = childApp.Metrics()
	// newApp.notification = childApp.Notification()
	newApp.saml = childApp.Saml()
	newApp.httpService = childApp.HTTPService()
	newApp.imageProxy = childApp.ImageProxy()
	newApp.timezones = childApp.Timezones()

	return &newApp
}

func (a *OpenTracingAppLayer) Srv() *app.Server {
	return a.srv
}
func (a *OpenTracingAppLayer) Log() *slog.Logger {
	return a.log
}
func (a *OpenTracingAppLayer) NotificationsLog() *slog.Logger {
	return a.notificationsLog
}
func (a *OpenTracingAppLayer) AccountMigration() einterfaces.AccountMigrationInterface {
	return a.accountMigration
}
func (a *OpenTracingAppLayer) Cluster() einterfaces.ClusterInterface {
	return a.cluster
}
func (a *OpenTracingAppLayer) Compliance() einterfaces.ComplianceInterface {
	return a.compliance
}
func (a *OpenTracingAppLayer) DataRetention() einterfaces.DataRetentionInterface {
	return a.dataRetention
}
func (a *OpenTracingAppLayer) Ldap() einterfaces.LdapInterface {
	return a.ldap
}

// func (a *OpenTracingAppLayer) MessageExport() einterfaces.MessageExportInterface {
// 	return a.messageExport
// }
func (a *OpenTracingAppLayer) Metrics() einterfaces.MetricsInterface {
	return a.metrics
}

// func (a *OpenTracingAppLayer) Notification() einterfaces.NotificationInterface {
// 	return a.notification
// }
func (a *OpenTracingAppLayer) Saml() einterfaces.SamlInterface {
	return a.saml
}
func (a *OpenTracingAppLayer) HTTPService() httpservice.HTTPService {
	return a.httpService
}
func (a *OpenTracingAppLayer) ImageProxy() *imageproxy.ImageProxy {
	return a.imageProxy
}
func (a *OpenTracingAppLayer) Timezones() *timezones.Timezones {
	return a.timezones
}
func (a *OpenTracingAppLayer) SetServer(srv *app.Server) {
	a.srv = srv
}
