package jobs

import (
	"errors"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/services/configservice"
)

// Workers contains some works that can be ran
type Workers struct {
	ConfigService configservice.ConfigService
	Watcher       *Watcher

	DataRetention            model.Worker
	MessageExport            model.Worker
	ElasticsearchIndexing    model.Worker
	ElasticsearchAggregation model.Worker
	LdapSync                 model.Worker
	Migrations               model.Worker
	Plugins                  model.Worker
	BleveIndexing            model.Worker
	ExpiryNotify             model.Worker
	ProductNotices           model.Worker
	ActiveUsers              model.Worker
	ImportProcess            model.Worker
	ImportDelete             model.Worker
	ExportProcess            model.Worker
	ExportDelete             model.Worker
	Cloud                    model.Worker
	ResendInvitationEmail    model.Worker

	listenerId string
	running    bool
}

var (
	ErrWorkersNotRunning    = errors.New("job workers are not running")
	ErrWorkersRunning       = errors.New("job workers are running")
	ErrWorkersUninitialized = errors.New("job workers are not initialized")
)

// Start starts the workers. This call is not safe for concurrent use.
// Synchronization should be implemented by the caller.
func (workers *Workers) Start() {
	slog.Info("Starting workers")
	if workers.DataRetention != nil &&
		(*workers.ConfigService.Config().DataRetentionSettings.EnableMessageDeletion ||
			*workers.ConfigService.Config().DataRetentionSettings.EnableFileDeletion) {
		go workers.DataRetention.Run()
	}

	if workers.MessageExport != nil && *workers.ConfigService.Config().MessageExportSettings.EnableExport {
		go workers.MessageExport.Run()
	}

	if workers.ElasticsearchIndexing != nil && *workers.ConfigService.Config().ElasticsearchSettings.EnableIndexing {
		go workers.ElasticsearchIndexing.Run()
	}

	if workers.ElasticsearchAggregation != nil && *workers.ConfigService.Config().ElasticsearchSettings.EnableIndexing {
		go workers.ElasticsearchAggregation.Run()
	}

	if workers.LdapSync != nil && *workers.ConfigService.Config().LdapSettings.EnableSync {
		go workers.LdapSync.Run()
	}

	if workers.Migrations != nil {
		go workers.Migrations.Run()
	}

	if workers.Plugins != nil {
		go workers.Plugins.Run()
	}

	if workers.BleveIndexing != nil &&
		*workers.ConfigService.Config().BleveSettings.EnableIndexing &&
		*workers.ConfigService.Config().BleveSettings.IndexDir != "" {
		go workers.BleveIndexing.Run()
	}

	if workers.ExpiryNotify != nil {
		go workers.ExpiryNotify.Run()
	}

	if workers.ActiveUsers != nil {
		go workers.ActiveUsers.Run()
	}

	if workers.ProductNotices != nil {
		go workers.ProductNotices.Run()
	}

	if workers.ImportProcess != nil {
		go workers.ImportProcess.Run()
	}

	if workers.ImportDelete != nil {
		go workers.ImportDelete.Run()
	}

	if workers.ExportProcess != nil {
		go workers.ExportProcess.Run()
	}

	if workers.ExportDelete != nil {
		go workers.ExportDelete.Run()
	}

	if workers.Cloud != nil {
		go workers.Cloud.Run()
	}

	if workers.ResendInvitationEmail != nil {
		go workers.ResendInvitationEmail.Run()
	}

	go workers.Watcher.Start()

	workers.listenerId = workers.ConfigService.AddConfigListener(workers.handleConfigChange)
	workers.running = true
}

func (workers *Workers) handleConfigChange(oldConfig *model.Config, newConfig *model.Config) {
	slog.Debug("Workers received config change.")

	if workers.DataRetention != nil {
		if (!*oldConfig.DataRetentionSettings.EnableMessageDeletion && !*oldConfig.DataRetentionSettings.EnableFileDeletion) &&
			(*newConfig.DataRetentionSettings.EnableMessageDeletion || *newConfig.DataRetentionSettings.EnableFileDeletion) {
			// check if:
			// 1) old message deletion and file deletion configurations are false
			// 2) new message deletion or file deletion configurations are true
			// then allow the "DataRetension" worker to run:
			go workers.DataRetention.Run()
		} else if (*oldConfig.DataRetentionSettings.EnableMessageDeletion || *oldConfig.DataRetentionSettings.EnableFileDeletion) &&
			(!*newConfig.DataRetentionSettings.EnableMessageDeletion && !*newConfig.DataRetentionSettings.EnableFileDeletion) {
			// check if:
			// 1) old message deletion or file deletion configurations are true
			// 2) new message deletion and file deletion configurations are false
			// then allow the "DataRetension" worker to run:
			// old configurations are the same when start this worker but new configurations does now allow this worker to run
			workers.DataRetention.Stop()
		}
	}

	if workers.MessageExport != nil {
		if !*oldConfig.MessageExportSettings.EnableExport && *newConfig.MessageExportSettings.EnableExport {
			// check if new configuration allows this worker to continue running
			go workers.MessageExport.Run()
		} else if *oldConfig.MessageExportSettings.EnableExport && !*newConfig.MessageExportSettings.EnableExport {
			workers.MessageExport.Stop()
		}
	}

	if workers.ElasticsearchIndexing != nil {
		if !*oldConfig.ElasticsearchSettings.EnableIndexing && *newConfig.ElasticsearchSettings.EnableIndexing {
			go workers.ElasticsearchIndexing.Run()
		} else if *oldConfig.ElasticsearchSettings.EnableIndexing && !*newConfig.ElasticsearchSettings.EnableIndexing {
			workers.ElasticsearchIndexing.Stop()
		}
	}

	if workers.ElasticsearchAggregation != nil {
		if !*oldConfig.ElasticsearchSettings.EnableIndexing && *newConfig.ElasticsearchSettings.EnableIndexing {
			go workers.ElasticsearchAggregation.Run()
		} else if *oldConfig.ElasticsearchSettings.EnableIndexing && !*newConfig.ElasticsearchSettings.EnableIndexing {
			workers.ElasticsearchAggregation.Stop()
		}
	}

	if workers.LdapSync != nil {
		if !*oldConfig.LdapSettings.EnableSync && *newConfig.LdapSettings.EnableSync {
			go workers.LdapSync.Run()
		} else if *oldConfig.LdapSettings.EnableSync && !*newConfig.LdapSettings.EnableSync {
			workers.LdapSync.Stop()
		}
	}

	if workers.BleveIndexing != nil {
		if !*oldConfig.BleveSettings.EnableIndexing && *newConfig.BleveSettings.EnableIndexing {
			go workers.BleveIndexing.Run()
		} else if *oldConfig.BleveSettings.EnableIndexing && !*newConfig.BleveSettings.EnableIndexing {
			workers.BleveIndexing.Stop()
		}
	}
}

// Stop stops the workers. This call is not safe for concurrent use.
// Synchronization should be implemented by the caller.
func (workers *Workers) Stop() {
	workers.ConfigService.RemoveConfigListener(workers.listenerId)

	workers.Watcher.Stop()

	if workers.DataRetention != nil &&
		(*workers.ConfigService.Config().DataRetentionSettings.EnableMessageDeletion ||
			*workers.ConfigService.Config().DataRetentionSettings.EnableFileDeletion) {
		workers.DataRetention.Stop()
	}

	if workers.MessageExport != nil && *workers.ConfigService.Config().MessageExportSettings.EnableExport {
		workers.MessageExport.Stop()
	}

	if workers.ElasticsearchIndexing != nil && *workers.ConfigService.Config().ElasticsearchSettings.EnableIndexing {
		workers.ElasticsearchIndexing.Stop()
	}

	if workers.ElasticsearchAggregation != nil && *workers.ConfigService.Config().ElasticsearchSettings.EnableIndexing {
		workers.ElasticsearchAggregation.Stop()
	}

	if workers.LdapSync != nil && *workers.ConfigService.Config().LdapSettings.EnableSync {
		workers.LdapSync.Stop()
	}

	if workers.Migrations != nil {
		workers.Migrations.Stop()
	}

	if workers.Plugins != nil {
		workers.Plugins.Stop()
	}

	if workers.BleveIndexing != nil && *workers.ConfigService.Config().BleveSettings.EnableIndexing {
		workers.BleveIndexing.Stop()
	}

	if workers.ExpiryNotify != nil {
		workers.ExpiryNotify.Stop()
	}

	if workers.ActiveUsers != nil {
		workers.ActiveUsers.Stop()
	}

	if workers.ProductNotices != nil {
		workers.ProductNotices.Stop()
	}

	if workers.ImportProcess != nil {
		workers.ImportProcess.Stop()
	}

	if workers.ImportDelete != nil {
		workers.ImportDelete.Stop()
	}

	if workers.ExportProcess != nil {
		workers.ExportProcess.Stop()
	}

	if workers.ExportDelete != nil {
		workers.ExportDelete.Stop()
	}

	if workers.Cloud != nil {
		workers.Cloud.Stop()
	}

	if workers.ResendInvitationEmail != nil {
		workers.ResendInvitationEmail.Stop()
	}

	workers.running = false

	slog.Info("Stopped workers")
}