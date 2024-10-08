package app

import (
	"errors"
	"fmt"
	"net/http"
	"os/user"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/audit"
	"github.com/sitename/sitename/modules/config"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
)

const (
	RestLevelID        = 240
	RestContentLevelID = 241
	RestPermsLevelID   = 242
	CLILevelID         = 243
)

var (
	LevelAPI     = slog.LvlAuditAPI
	LevelContent = slog.LvlAuditContent
	LevelPerms   = slog.LvlAuditPerms
	LevelCLI     = slog.LvlAuditCLI
)

func (a *App) GetAudits(userID string, limit int) (model.AuditSlice, *model_helper.AppError) {
	audits, err := a.Srv().Store.Audit().Get(userID, 0, limit)
	if err != nil {
		var outErr *store.ErrOutOfBounds
		switch {
		case errors.As(err, &outErr):
			return nil, model_helper.NewAppError("GetAudits", "app.audit.get.limit.app_error", nil, err.Error(), http.StatusBadRequest)
		default:
			return nil, model_helper.NewAppError("GetAudits", "app.audit.get.finding.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}
	return audits, nil
}

func (a *App) GetAuditsPage(userID string, page int, perPage int) (model.AuditSlice, *model_helper.AppError) {
	audits, err := a.Srv().Store.Audit().Get(userID, page*perPage, perPage)
	if err != nil {
		var outErr *store.ErrOutOfBounds
		switch {
		case errors.As(err, &outErr):
			return nil, model_helper.NewAppError("GetAuditsPage", "app.audit.get.limit.app_error", nil, err.Error(), http.StatusBadRequest)
		default:
			return nil, model_helper.NewAppError("GetAuditsPage", "app.audit.get.finding.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}
	return audits, nil
}

// LogAuditRec logs an audit record using default LvlAuditCLI.
func (a *App) LogAuditRec(rec *audit.Record, err error) {
	a.LogAuditRecWithLevel(rec, slog.LvlAuditCLI, err)
}

// LogAuditRecWithLevel logs an audit record using specified Level.
func (a *App) LogAuditRecWithLevel(rec *audit.Record, level slog.Level, err error) {
	if rec == nil {
		return
	}
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			rec.AddMeta("err", appErr.Error())
			rec.AddMeta("code", appErr.StatusCode)
		} else {
			rec.AddMeta("err", err)
		}
		rec.Fail()
	}
	a.Srv().Audit.LogRecord(level, *rec)
}

// MakeAuditRecord creates a audit record pre-populated with defaults.
func (a *App) MakeAuditRecord(event string, initialStatus string) *audit.Record {
	var userID string
	user, err := user.Current()
	if err == nil {
		userID = fmt.Sprintf("%s:%s", user.Uid, user.Username)
	}

	rec := &audit.Record{
		APIPath:   "",
		Event:     event,
		Status:    initialStatus,
		UserID:    userID,
		SessionID: "",
		Client:    fmt.Sprintf("server %s-%s", model_helper.BuildNumber, model_helper.BuildHash),
		IPAddress: "",
		Meta:      audit.Meta{audit.KeyClusterID: a.GetClusterId()},
	}
	rec.AddMetaTypeConverter(model_helper.AuditModelTypeConv)

	return rec
}
func (s *Server) configureAudit(adt *audit.Audit, bAllowAdvancedLogging bool) error {
	adt.OnQueueFull = s.onAuditTargetQueueFull
	adt.OnError = s.onAuditError

	var logConfigSrc config.LogConfigSrc
	dsn := *s.Config().ExperimentalAuditSettings.AdvancedLoggingConfig
	if bAllowAdvancedLogging && dsn != "" {
		var err error
		logConfigSrc, err = config.NewLogConfigSrc(dsn, s.ConfigStore)
		if err != nil {
			return fmt.Errorf("invalid config source for audit, %w", err)
		}
		slog.Debug("Loaded audit configuration", slog.String("source", dsn))
	}

	// ExperimentalAuditSettings provides basic file audit (E0, E10); logConfigSrc provides advanced config (E20).
	cfg, err := config.MloggerConfigFromAuditConfig(s.Config().ExperimentalAuditSettings, logConfigSrc)
	if err != nil {
		return fmt.Errorf("invalid config for audit, %w", err)
	}

	return adt.Configure(cfg)
}

func (s *Server) onAuditTargetQueueFull(qname string, maxQSize int) bool {
	slog.Error("Audit queue full, dropping record.", slog.String("qname", qname), slog.Int("queueSize", maxQSize))
	return true // drop it
}

func (s *Server) onAuditError(err error) {
	slog.Error("Audit Error", slog.Err(err))
}
