package opentracing

import (
	"context"
	"crypto/ecdsa"
	"io"
	"net/http"
	"reflect"

	"github.com/opentracing/opentracing-go/ext"
	spanlog "github.com/opentracing/opentracing-go/log"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/audit"
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

	cluster          einterfaces.ClusterInterface
	compliance       einterfaces.ComplianceInterface
	searchEngine     *searchengine.Broker
	ldap             einterfaces.LdapInterface
	metrics          einterfaces.MetricsInterface
	httpService      httpservice.HTTPService
	imageProxy       *imageproxy.ImageProxy
	timezones        *timezones.Timezones
	saml             einterfaces.SamlInterface
	dataRetention    einterfaces.DataRetentionInterface
	accountMigration einterfaces.AccountMigrationInterface

	// messageExport    einterfaces.MessageExportInterface
	// notification     einterfaces.NotificationInterface

	ctx context.Context
}

func (a *OpenTracingAppLayer) AccountService() sub_app_iface.AccountService {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.AccountService")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.AccountService()

	return resultVar0
}

func (a *OpenTracingAppLayer) AddConfigListener(listener func(*model_helper.Config, *model_helper.Config)) string {
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

func (a *OpenTracingAppLayer) AttributeService() sub_app_iface.AttributeService {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.AttributeService")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.AttributeService()

	return resultVar0
}

func (a *OpenTracingAppLayer) ChannelService() sub_app_iface.ChannelService {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.ChannelService")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.ChannelService()

	return resultVar0
}

func (a *OpenTracingAppLayer) CheckoutService() sub_app_iface.CheckoutService {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.CheckoutService")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.CheckoutService()

	return resultVar0
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

func (a *OpenTracingAppLayer) Config() *model_helper.Config {
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

func (a *OpenTracingAppLayer) CsvService() sub_app_iface.CsvService {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.CsvService")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.CsvService()

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

func (a *OpenTracingAppLayer) DiscountService() sub_app_iface.DiscountService {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.DiscountService")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.DiscountService()

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

func (a *OpenTracingAppLayer) EnvironmentConfig(filter func(reflect.StructField) bool) map[string]app.any {
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

func (a *OpenTracingAppLayer) FileService() sub_app_iface.FileService {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.FileService")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.FileService()

	return resultVar0
}

func (a *OpenTracingAppLayer) GetAudits(userID string, limit int) (model.AuditSlice, *model_helper.AppError) {
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

func (a *OpenTracingAppLayer) GetAuditsPage(userID string, page int, perPage int) (model.AuditSlice, *model_helper.AppError) {
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

func (a *OpenTracingAppLayer) GetClusterStatus() []*model_helper.ClusterInfo {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetClusterStatus")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.GetClusterStatus()

	return resultVar0
}

func (a *OpenTracingAppLayer) GetComplianceFile(job *model.Compliance) ([]byte, *model_helper.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetComplianceFile")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetComplianceFile(job)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetComplianceReport(reportID string) (*model.Compliance, *model_helper.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetComplianceReport")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetComplianceReport(reportID)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetComplianceReports(page int, perPage int) (model.ComplianceSlice, *model_helper.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetComplianceReports")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetComplianceReports(page, perPage)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
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

func (a *OpenTracingAppLayer) GetEnvironmentConfig(filter func(reflect.StructField) bool) map[string]app.any {
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

func (a *OpenTracingAppLayer) GetLogs(page int, perPage int) ([]string, *model_helper.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetLogs")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetLogs(page, perPage)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetLogsSkipSend(page int, perPage int) ([]string, *model_helper.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetLogsSkipSend")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetLogsSkipSend(page, perPage)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetOpenGraphMetadata(requestURL string) ([]byte, error) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetOpenGraphMetadata")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetOpenGraphMetadata(requestURL)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetSanitizedConfig() *model_helper.Config {
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

func (a *OpenTracingAppLayer) GetSystemInstallDate() (int64, *model_helper.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GetSystemInstallDate")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.GetSystemInstallDate()

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) GetWarnMetricsStatus() (map[string]*model_helper.WarnMetricStatus, *model_helper.AppError) {
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

func (a *OpenTracingAppLayer) GiftcardService() sub_app_iface.GiftcardService {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.GiftcardService")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.GiftcardService()

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

func (a *OpenTracingAppLayer) HandleMessageExportConfig(cfg *model_helper.Config, appCfg *model_helper.Config) {
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

func (a *OpenTracingAppLayer) ImageProxyAdder() func(string) string {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.ImageProxyAdder")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.ImageProxyAdder()

	return resultVar0
}

func (a *OpenTracingAppLayer) ImageProxyRemover() func(string) string {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.ImageProxyRemover")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.ImageProxyRemover()

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

func (a *OpenTracingAppLayer) InvoiceService() sub_app_iface.InvoiceService {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.InvoiceService")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.InvoiceService()

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

func (a *OpenTracingAppLayer) LogAuditRecWithLevel(rec *audit.Record, level slog.Level, err error) {
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

func (a *OpenTracingAppLayer) MenuService() sub_app_iface.MenuService {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.MenuService")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.MenuService()

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

func (a *OpenTracingAppLayer) NotifyAndSetWarnMetricAck(warnMetricId string, sender model.User, forceAck bool, isBot bool) *model_helper.AppError {
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

func (a *OpenTracingAppLayer) OrderService() sub_app_iface.OrderService {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.OrderService")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.OrderService()

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

func (a *OpenTracingAppLayer) PageService() sub_app_iface.PageService {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.PageService")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.PageService()

	return resultVar0
}

func (a *OpenTracingAppLayer) PaymentService() sub_app_iface.PaymentService {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.PaymentService")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.PaymentService()

	return resultVar0
}

func (a *OpenTracingAppLayer) PluginService() sub_app_iface.PluginService {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.PluginService")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.PluginService()

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

func (a *OpenTracingAppLayer) ProductService() sub_app_iface.ProductService {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.ProductService")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.ProductService()

	return resultVar0
}

func (a *OpenTracingAppLayer) Publish(message *model_helper.WebSocketEvent) {
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

func (a *OpenTracingAppLayer) ResetPermissionsSystem() *model_helper.AppError {
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

func (a *OpenTracingAppLayer) SaveComplianceReport(job model.Compliance) (*model.Compliance, *model_helper.AppError) {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.SaveComplianceReport")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0, resultVar1 := a.app.SaveComplianceReport(job)

	if resultVar1 != nil {
		span.LogFields(spanlog.Error(resultVar1))
		ext.Error.Set(span, true)
	}

	return resultVar0, resultVar1
}

func (a *OpenTracingAppLayer) SaveConfig(newCfg *model_helper.Config, sendConfigChangeClusterMessage bool) (*model_helper.Config, *model_helper.Config, *model_helper.AppError) {
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

func (a *OpenTracingAppLayer) SeoService() sub_app_iface.SeoService {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.SeoService")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.SeoService()

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

func (a *OpenTracingAppLayer) ShippingService() sub_app_iface.ShippingService {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.ShippingService")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.ShippingService()

	return resultVar0
}

func (a *OpenTracingAppLayer) UpdateConfig(f func(*model_helper.Config)) {
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

func (a *OpenTracingAppLayer) WarehouseService() sub_app_iface.WarehouseService {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.WarehouseService")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.WarehouseService()

	return resultVar0
}

func (a *OpenTracingAppLayer) WebhookService() sub_app_iface.WebhookService {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.WebhookService")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.WebhookService()

	return resultVar0
}

func (a *OpenTracingAppLayer) WishlistService() sub_app_iface.WishlistService {
	origCtx := a.ctx
	span, newCtx := tracing.StartSpanWithParentByContext(a.ctx, "app.WishlistService")

	a.ctx = newCtx
	a.app.Srv().Store.SetContext(newCtx)
	defer func() {
		a.app.Srv().Store.SetContext(origCtx)
		a.ctx = origCtx
	}()

	defer span.Finish()
	resultVar0 := a.app.WishlistService()

	return resultVar0
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
	newApp.metrics = childApp.Metrics()
	newApp.saml = childApp.Saml()
	newApp.httpService = childApp.HTTPService()
	newApp.imageProxy = childApp.ImageProxy()
	newApp.timezones = childApp.Timezones()

	// newApp.notification = childApp.Notification()
	// newApp.messageExport = childApp.MessageExport()

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
func (a *OpenTracingAppLayer) Metrics() einterfaces.MetricsInterface {
	return a.metrics
}

//	func (a *OpenTracingAppLayer) MessageExport() einterfaces.MessageExportInterface {
//		return a.messageExport
//	}
//
//	func (a *OpenTracingAppLayer) Notification() einterfaces.NotificationInterface {
//		return a.notification
//	}
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
