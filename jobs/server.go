package jobs

import (
	"sync"
	"time"

	"github.com/sitename/sitename/einterfaces"
	ejobs "github.com/sitename/sitename/einterfaces/jobs"
	tjobs "github.com/sitename/sitename/jobs/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/services/configservice"
	"github.com/sitename/sitename/store"
)

type JobServer struct {
	ConfigService           configservice.ConfigService
	Store                   store.Store
	metrics                 einterfaces.MetricsInterface
	DataRetentionJob        ejobs.DataRetentionJobInterface
	MessageExportJob        ejobs.MessageExportJobInterface
	ElasticsearchAggregator ejobs.ElasticsearchAggregatorInterface
	ElasticsearchIndexer    tjobs.IndexerJobInterface
	LdapSync                ejobs.LdapSyncInterface
	Migrations              tjobs.MigrationsJobInterface
	Plugins                 tjobs.PluginsJobInterface
	BleveIndexer            tjobs.IndexerJobInterface
	ExpiryNotify            tjobs.ExpiryNotifyJobInterface
	ProductNotices          tjobs.ProductNoticesJobInterface
	ActiveUsers             tjobs.ActiveUsersJobInterface
	ImportProcess           tjobs.ImportProcessInterface
	ImportDelete            tjobs.ImportDeleteInterface
	ExportProcess           tjobs.ExportProcessInterface
	ExportDelete            tjobs.ExportDeleteInterface
	Cloud                   ejobs.CloudJobInterface
	ResendInvitationEmails  ejobs.ResendInvitationEmailJobInterface

	// mut is used to protect the following fields from concurrent access.
	mut        sync.Mutex
	workers    *Workers
	schedulers *Schedulers
}

func NewJobServer(configService configservice.ConfigService, store store.Store, metrics einterfaces.MetricsInterface) *JobServer {
	return &JobServer{
		ConfigService: configService,
		Store:         store,
		metrics:       metrics,
	}
}

func (srv *JobServer) MakeWatcher(workers *Workers, pollingInterval int) *Watcher {
	return &Watcher{
		stop:            make(chan struct{}),
		stopped:         make(chan struct{}),
		pollingInterval: pollingInterval,
		workers:         workers,
		srv:             srv,
	}
}

// InitWorkers initializes all the registered workers
func (srv *JobServer) InitWorkers() error {
	srv.mut.Lock()
	defer srv.mut.Unlock()
	if srv.workers != nil && srv.workers.running {
		return ErrWorkersRunning
	}

	slog.Debug("Initialising workers.")

	workers := &Workers{
		ConfigService: srv.ConfigService,
	}
	workers.Watcher = srv.MakeWatcher(workers, DefaultWatcherPollingInterval)
	if srv.DataRetentionJob != nil {
		workers.DataRetention = srv.DataRetentionJob.MakeWorker()
	}
	if srv.MessageExportJob != nil {
		workers.MessageExport = srv.MessageExportJob.MakeWorker()
	}
	if srv.ElasticsearchIndexer != nil {
		workers.ElasticsearchIndexing = srv.ElasticsearchIndexer.MakeWorker()
	}
	if srv.ElasticsearchAggregator != nil {
		workers.ElasticsearchAggregation = srv.ElasticsearchAggregator.MakeWorker()
	}
	if srv.LdapSync != nil {
		workers.LdapSync = srv.LdapSync.MakeWorker()
	}
	if srv.Migrations != nil {
		workers.Migrations = srv.Migrations.MakeWorker()
	}
	if srv.Plugins != nil {
		workers.Plugins = srv.Plugins.MakeWorker()
	}
	if srv.BleveIndexer != nil {
		workers.BleveIndexing = srv.BleveIndexer.MakeWorker()
	}
	if srv.ExpiryNotify != nil {
		workers.ExpiryNotify = srv.ExpiryNotify.MakeWorker()
	}
	if srv.ActiveUsers != nil {
		workers.ActiveUsers = srv.ActiveUsers.MakeWorker()
	}
	if srv.ProductNotices != nil {
		workers.ProductNotices = srv.ProductNotices.MakeWorker()
	}
	if srv.ImportProcess != nil {
		workers.ImportProcess = srv.ImportProcess.MakeWorker()
	}
	if srv.ImportDelete != nil {
		workers.ImportDelete = srv.ImportDelete.MakeWorker()
	}
	if srv.ExportProcess != nil {
		workers.ExportProcess = srv.ExportProcess.MakeWorker()
	}
	if srv.ExportDelete != nil {
		workers.ExportDelete = srv.ExportDelete.MakeWorker()
	}
	if srv.Cloud != nil {
		workers.Cloud = srv.Cloud.MakeWorker()
	}
	if srv.ResendInvitationEmails != nil {
		workers.ResendInvitationEmail = srv.ResendInvitationEmails.MakeWorker()
	}

	srv.workers = workers
	return nil
}

// InitSchedulers inits all job schedulers
func (srv *JobServer) InitSchedulers() error {
	srv.mut.Lock()
	defer srv.mut.Unlock()
	if srv.schedulers != nil && srv.schedulers.running {
		return ErrSchedulersRunning
	}
	slog.Debug("Initialising schedulers.")

	schedulers := &Schedulers{
		stop:                 make(chan bool),
		stopped:              make(chan bool),
		configChanged:        make(chan *model.Config),
		clusterLeaderChanged: make(chan bool),
		jobs:                 srv,
		isLeader:             true,
	}

	if srv.DataRetentionJob != nil {
		schedulers.schedulers = append(schedulers.schedulers, srv.DataRetentionJob.MakeScheduler())
	}
	if srv.MessageExportJob != nil {
		schedulers.schedulers = append(schedulers.schedulers, srv.MessageExportJob.MakeScheduler())
	}
	if srv.ElasticsearchAggregator != nil {
		schedulers.schedulers = append(schedulers.schedulers, srv.ElasticsearchAggregator.MakeScheduler())
	}
	if srv.LdapSync != nil {
		schedulers.schedulers = append(schedulers.schedulers, srv.LdapSync.MakeScheduler())
	}
	if srv.Migrations != nil {
		schedulers.schedulers = append(schedulers.schedulers, srv.Migrations.MakeScheduler())
	}
	if srv.Plugins != nil {
		schedulers.schedulers = append(schedulers.schedulers, srv.Plugins.MakeScheduler())
	}
	if srv.ExpiryNotify != nil {
		schedulers.schedulers = append(schedulers.schedulers, srv.ExpiryNotify.MakeScheduler())
	}
	if srv.ActiveUsers != nil {
		schedulers.schedulers = append(schedulers.schedulers, srv.ActiveUsers.MakeScheduler())
	}
	if srv.ProductNotices != nil {
		schedulers.schedulers = append(schedulers.schedulers, srv.ProductNotices.MakeScheduler())
	}
	if srv.Cloud != nil {
		schedulers.schedulers = append(schedulers.schedulers, srv.Cloud.MakeScheduler())
	}
	if srv.ResendInvitationEmails != nil {
		schedulers.schedulers = append(schedulers.schedulers, srv.ResendInvitationEmails.MakeScheduler())
	}
	if srv.ImportDelete != nil {
		schedulers.schedulers = append(schedulers.schedulers, srv.ImportDelete.MakeScheduler())
	}
	if srv.ExportDelete != nil {
		schedulers.schedulers = append(schedulers.schedulers, srv.ExportDelete.MakeScheduler())
	}

	schedulers.nextRunTimes = make([]*time.Time, len(schedulers.schedulers))
	srv.schedulers = schedulers

	return nil
}

func (srv *JobServer) Config() *model.Config {
	return srv.ConfigService.Config()
}

func (srv *JobServer) StartWorkers() error {
	srv.mut.Lock()
	defer srv.mut.Unlock()
	if srv.workers == nil {
		return ErrWorkersUninitialized
	} else if srv.workers.running {
		return ErrWorkersRunning
	}
	srv.workers.Start()
	return nil
}

func (srv *JobServer) StartSchedulers() error {
	srv.mut.Lock()
	defer srv.mut.Unlock()
	if srv.schedulers == nil {
		return ErrSchedulersUninitialized
	} else if srv.schedulers.running {
		return ErrSchedulersRunning
	}
	srv.schedulers.Start()
	return nil
}

func (srv *JobServer) StopWorkers() error {
	srv.mut.Lock()
	defer srv.mut.Unlock()
	if srv.workers == nil {
		return ErrWorkersUninitialized
	} else if !srv.workers.running {
		return ErrWorkersNotRunning
	}
	srv.workers.Stop()
	return nil
}

func (srv *JobServer) StopSchedulers() error {
	srv.mut.Lock()
	defer srv.mut.Unlock()
	if srv.schedulers == nil {
		return ErrSchedulersUninitialized
	} else if !srv.schedulers.running {
		return ErrSchedulersNotRunning
	}
	srv.schedulers.Stop()
	return nil
}

func (srv *JobServer) HandleClusterLeaderChange(isLeader bool) {
	srv.mut.Lock()
	defer srv.mut.Unlock()
	if srv.schedulers != nil {
		srv.schedulers.HandleClusterLeaderChange(isLeader)
	}
}
