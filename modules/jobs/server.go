package jobs

import (
	"sync"

	"github.com/sitename/sitename/einterfaces"
	ejobs "github.com/sitename/sitename/einterfaces/jobs"
	"github.com/sitename/sitename/model"
	tjobs "github.com/sitename/sitename/modules/jobs/interfaces"
	"github.com/sitename/sitename/services/configservice"
	"github.com/sitename/sitename/store"
)

type JobServer struct {
	ConfigService configservice.ConfigService
	Store         store.Store
	metrics       einterfaces.MetricsInterface

	DataRetentionJob        ejobs.DataRetentionJobInterface         //
	MessageExportJob        ejobs.MessageExportJobInterface         //
	ElasticsearchAggregator ejobs.ElasticsearchAggregatorInterface  //
	ElasticsearchIndexer    tjobs.IndexerJobInterface               //
	LdapSync                ejobs.LdapSyncInterface                 //
	Migrations              tjobs.MigrationsJobInterface            //
	Plugins                 tjobs.PluginsJobInterface               //
	BleveIndexer            tjobs.IndexerJobInterface               //
	ExpiryNotify            tjobs.ExpiryNotifyJobInterface          //
	ProductNotices          tjobs.ProductNoticesJobInterface        //
	ActiveUsers             tjobs.ActiveUsersJobInterface           //
	ImportProcess           tjobs.ImportProcessInterface            //
	ImportDelete            tjobs.ImportDeleteInterface             //
	ExportProcess           tjobs.ExportProcessInterface            //
	ExportDelete            tjobs.ExportDeleteInterface             //
	Cloud                   ejobs.CloudJobInterface                 //
	ResendInvitationEmails  ejobs.ResendInvitationEmailJobInterface //
	CsvExport               tjobs.CsvExportInterface                // csv export

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

func (srv *JobServer) Config() *model.Config {
	return srv.ConfigService.Config()
}

func (srv *JobServer) StartWorkers() error {
	srv.mut.Lock()
	defer srv.mut.Unlock()

	// check if either workers are not initialized or they are already running. If true -> return
	if srv.workers == nil {
		return ErrWorkersUninitialized
	} else if srv.workers.running {
		return ErrWorkersRunning
	}

	// starts workers
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
