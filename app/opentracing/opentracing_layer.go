package opentracing

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"io"
	"mime/multipart"
	"net/http"
	"reflect"

	"github.com/opentracing/opentracing-go/ext"
	spanlog "github.com/opentracing/opentracing-go/log"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/modules/filestore"
	"github.com/sitename/sitename/modules/i18n"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/timezones"
	"github.com/sitename/sitename/services/httpservice"
	"github.com/sitename/sitename/services/imageproxy"
	"github.com/sitename/sitename/services/searchengine"
	"github.com/sitename/sitename/services/tracing"
)

type OpenTracingAppLayer struct {
	app app.AppIface

	srv *app.Server

	log              *slog.Logger
	notificationsLog *slog.Logger

	t              i18n.TranslateFunc
	session        model.Session
	requestId      string
	ipAddress      string
	path           string
	userAgent      string
	acceptLanguage string

	// accountMigration einterfaces.AccountMigrationInterface
	cluster    einterfaces.ClusterInterface
	compliance einterfaces.ComplianceInterface
	// dataRetention    einterfaces.DataRetentionInterface
	searchEngine *searchengine.Broker
	ldap         einterfaces.LdapInterface
	// messageExport    einterfaces.MessageExportInterface
	metrics einterfaces.MetricsInterface
	// notification     einterfaces.NotificationInterface
	// saml             einterfaces.SamlInterface

	httpService httpservice.HTTPService
	imageProxy  *imageproxy.ImageProxy
	timezones   *timezones.Timezones

	context context.Context
	ctx     context.Context
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

func (a *OpenTracingAppLayer) CreateGuest(user *account.User) (*account.User, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.CreateGuest")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.CreateGuest(user)

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

func (a *OpenTracingAppLayer) CreateUser(user *account.User) (*account.User, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.CreateUser")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.CreateUser(user)

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

func (a *OpenTracingAppLayer) CreateUserAsAdmin(user *account.User, redirect string) (*account.User, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.CreateUserAsAdmin")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.CreateUserAsAdmin(user, redirect)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) CreateUserFromSignup(user *account.User, redirect string) (*account.User, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.CreateUserFromSignup")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.CreateUserFromSignup(user, redirect)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
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

func (a *OpenTracingAppLayer) DeactivateGuests() *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.DeactivateGuests")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.DeactivateGuests()

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

func (a *OpenTracingAppLayer) InitServer() {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.InitServer")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	a.app.InitServer()
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

func (a *OpenTracingAppLayer) SaveConfig(newCfg *model.Config, sendConfigChangeClusterMessage bool) *model.AppError {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.SaveConfig")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.SaveConfig(newCfg, sendConfigChangeClusterMessage)

	if resultVar0 != nil {
		span.LogFields(spanlog.Error(resultVar0))
		ext.Error.Set(span, true)
	}

	return resultVar0
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

func (a *OpenTracingAppLayer) UpdateActive(user *account.User, active bool) (*account.User, *model.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.UpdateActive")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.UpdateActive(user, active)

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
	newApp.t = childApp.GetT()
	if childApp.Session() != nil {
		newApp.session = *childApp.Session()
	}
	newApp.requestId = childApp.RequestId()
	newApp.ipAddress = childApp.IpAddress()
	newApp.path = childApp.Path()
	newApp.userAgent = childApp.UserAgent()
	newApp.acceptLanguage = childApp.AcceptLanguage()
	//newApp.accountMigration = childApp.AccountMigration()
	newApp.cluster = childApp.Cluster()
	newApp.compliance = childApp.Compliance()
	//newApp.dataRetention = childApp.DataRetention()
	newApp.searchEngine = childApp.SearchEngine()
	newApp.ldap = childApp.Ldap()
	//newApp.messageExport = childApp.MessageExport()
	newApp.metrics = childApp.Metrics()
	//newApp.notification = childApp.Notification()
	//newApp.saml = childApp.Saml()
	newApp.httpService = childApp.HTTPService()
	newApp.imageProxy = childApp.ImageProxy()
	newApp.timezones = childApp.Timezones()
	newApp.context = childApp.Context()

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
func (a *OpenTracingAppLayer) T(translationID string, args ...interface{}) string {
	return a.t(translationID, args...)
}
func (a *OpenTracingAppLayer) Session() *model.Session {
	return &a.session
}
func (a *OpenTracingAppLayer) RequestId() string {
	return a.requestId
}
func (a *OpenTracingAppLayer) IpAddress() string {
	return a.ipAddress
}
func (a *OpenTracingAppLayer) Path() string {
	return a.path
}
func (a *OpenTracingAppLayer) UserAgent() string {
	return a.userAgent
}
func (a *OpenTracingAppLayer) AcceptLanguage() string {
	return a.acceptLanguage
}

//func (a *OpenTracingAppLayer) AccountMigration() einterfaces.AccountMigrationInterface {
//	return a.accountMigration
//}
func (a *OpenTracingAppLayer) Cluster() einterfaces.ClusterInterface {
	return a.cluster
}
func (a *OpenTracingAppLayer) Compliance() einterfaces.ComplianceInterface {
	return a.compliance
}

//func (a *OpenTracingAppLayer) DataRetention() einterfaces.DataRetentionInterface {
//	return a.dataRetention
//}
func (a *OpenTracingAppLayer) Ldap() einterfaces.LdapInterface {
	return a.ldap
}

//func (a *OpenTracingAppLayer) MessageExport() einterfaces.MessageExportInterface {
//	return a.messageExport
//}
func (a *OpenTracingAppLayer) Metrics() einterfaces.MetricsInterface {
	return a.metrics
}

//func (a *OpenTracingAppLayer) Notification() einterfaces.NotificationInterface {
//	return a.notification
//}
//func (a *OpenTracingAppLayer) Saml() einterfaces.SamlInterface {
//	return a.saml
//}
func (a *OpenTracingAppLayer) HTTPService() httpservice.HTTPService {
	return a.httpService
}
func (a *OpenTracingAppLayer) ImageProxy() *imageproxy.ImageProxy {
	return a.imageProxy
}
func (a *OpenTracingAppLayer) Timezones() *timezones.Timezones {
	return a.timezones
}
func (a *OpenTracingAppLayer) Context() context.Context {
	return a.context
}
func (a *OpenTracingAppLayer) SetSession(sess *model.Session) {
	a.session = *sess
}
func (a *OpenTracingAppLayer) SetT(t i18n.TranslateFunc) {
	a.t = t
}
func (a *OpenTracingAppLayer) SetRequestId(str string) {
	a.requestId = str
}
func (a *OpenTracingAppLayer) SetIpAddress(str string) {
	a.ipAddress = str
}
func (a *OpenTracingAppLayer) SetUserAgent(str string) {
	a.userAgent = str
}
func (a *OpenTracingAppLayer) SetAcceptLanguage(str string) {
	a.acceptLanguage = str
}
func (a *OpenTracingAppLayer) SetPath(str string) {
	a.path = str
}
func (a *OpenTracingAppLayer) SetContext(c context.Context) {
	a.context = c
}
func (a *OpenTracingAppLayer) SetServer(srv *app.Server) {
	a.srv = srv
}
func (a *OpenTracingAppLayer) GetT() i18n.TranslateFunc {
	return a.t
}
