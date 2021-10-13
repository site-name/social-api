package audit

import (
	"fmt"

	"github.com/sitename/sitename/modules/slog"
)

const (
	DefMaxQueueSize = 1000

	KeyAPIPath   = "api_path"
	KeyEvent     = "event"
	KeyStatus    = "status"
	KeyUserID    = "user_id"
	KeySessionID = "session_id"
	KeyClient    = "client"
	KeyIPAddress = "ip_address"
	KeyClusterID = "cluster_id"

	Success = "success"
	Attempt = "attempt"
	Fail    = "fail"
)

type Audit struct {
	logger *slog.Logger

	// OnQueueFull is called on an attempt to add an audit record to a full queue.
	// Return true to drop record, or false to block until there is room in queue.
	OnQueueFull func(qname string, maxQueueSize int) bool

	// OnError is called when an error occurs while writing an audit record.
	OnError func(err error)
}

func (a *Audit) Init(maxQueueSize int) {
	a.logger, _ = slog.NewLogger(
		slog.MaxQueueSize(maxQueueSize),
		slog.OnLoggerError(a.onLoggerError),
		slog.OnQueueFull(a.onQueueFull),
		slog.OnTargetQueueFull(a.onTargetQueueFull),
	)
}

func (a *Audit) LogRecord(level slog.Level, rec Record) {
	flds := []slog.Field{
		slog.String(KeyAPIPath, rec.APIPath),
		slog.String(KeyEvent, rec.Event),
		slog.String(KeyStatus, rec.Status),
		slog.String(KeyUserID, rec.UserID),
		slog.String(KeySessionID, rec.SessionID),
		slog.String(KeyClient, rec.Client),
		slog.String(KeyIPAddress, rec.IPAddress),
	}

	for k, v := range rec.Meta {
		flds = append(flds, slog.Any(k, v))
	}
	a.logger.Log(level, "", flds...)
}

// Log emits an audit record based on minimum required info.
func (a *Audit) Log(level slog.Level, path string, evt string, status string, userID string, sessionID string, meta Meta) {
	a.LogRecord(level, Record{
		APIPath:   path,
		Event:     evt,
		Status:    status,
		UserID:    userID,
		SessionID: sessionID,
		Meta:      meta,
	})
}

func (a *Audit) Configure(cfg slog.LoggerConfiguration) error {
	return a.logger.ConfigureTargets(cfg, nil)
}

// Flush attempts to write all queued audit records to all targets.
func (a *Audit) Flush() error {
	err := a.logger.Flush()
	if err != nil {
		a.onLoggerError(err)
	}
	return err
}

// Shutdown cleanly stops the audit engine after making best efforts to flush all targets.
func (a *Audit) Shutdown() error {
	err := a.logger.Shutdown()
	if err != nil {
		a.onLoggerError(err)
	}
	return err
}

func (a *Audit) onQueueFull(rec *slog.LogRec, maxQueueSize int) bool {
	if a.OnQueueFull != nil {
		return a.OnQueueFull("main", maxQueueSize)
	}
	slog.Error("Audit logging queue full, dropping record.", slog.Int("queueSize", maxQueueSize))
	return true
}

func (a *Audit) onTargetQueueFull(target slog.Target, rec *slog.LogRec, maxQueueSize int) bool {
	if a.OnQueueFull != nil {
		return a.OnQueueFull(fmt.Sprintf("%v", target), maxQueueSize)
	}
	slog.Error("Audit logging queue full for target, dropping record.", slog.Any("target", target), slog.Int("queueSize", maxQueueSize))
	return true
}

func (a *Audit) onLoggerError(err error) {
	if a.OnError != nil {
		a.OnError(err)
	}
}
