package app

import (
	"github.com/sitename/sitename/einterfaces"
	ejobs "github.com/sitename/sitename/einterfaces/jobs"
	tjobs "github.com/sitename/sitename/modules/jobs/interfaces"
	"github.com/sitename/sitename/services/searchengine"
)

// normal job
var (
	accountMigrationInterface            func(*App) einterfaces.AccountMigrationInterface     //
	jobsLdapSyncInterface                func(*App) ejobs.LdapSyncInterface                   //
	jobsResendInvitationEmailInterface   func(*App) ejobs.ResendInvitationEmailJobInterface   //
	jobsDataRetentionJobInterface        func(*Server) ejobs.DataRetentionJobInterface        //
	jobsMessageExportJobInterface        func(*Server) ejobs.MessageExportJobInterface        //
	jobsElasticsearchAggregatorInterface func(*Server) ejobs.ElasticsearchAggregatorInterface //
	jobsElasticsearchIndexerInterface    func(*Server) tjobs.IndexerJobInterface              //
	jobsBleveIndexerInterface            func(*Server) tjobs.IndexerJobInterface              //
	jobsMigrationsInterface              func(*Server) tjobs.MigrationsJobInterface           //
	csvExportInterface                   func(*Server) tjobs.CsvExportInterface               // csv export work
	jobsPluginsInterface                 func(*Server) tjobs.PluginsJobInterface              //
	jobsExtractContentInterface          func(*Server) tjobs.ExtractContentInterface          //
)

func RegisterJobsExtractContentInterface(f func(*Server) tjobs.ExtractContentInterface) {
	jobsExtractContentInterface = f
}

func RegisterJobsPluginsJobInterface(f func(*Server) tjobs.PluginsJobInterface) {
	jobsPluginsInterface = f
}

func RegisterAccountMigrationInterface(f func(*App) einterfaces.AccountMigrationInterface) {
	accountMigrationInterface = f
}

func RegisterJobsLdapSyncInterface(f func(*App) ejobs.LdapSyncInterface) {
	jobsLdapSyncInterface = f
}

func RegisterJobsDataRetentionJobInterface(f func(*Server) ejobs.DataRetentionJobInterface) {
	jobsDataRetentionJobInterface = f
}

func RegisterJobsMessageExportJobInterface(f func(*Server) ejobs.MessageExportJobInterface) {
	jobsMessageExportJobInterface = f
}

func RegisterJobsElasticsearchAggregatorInterface(f func(*Server) ejobs.ElasticsearchAggregatorInterface) {
	jobsElasticsearchAggregatorInterface = f
}

func RegisterJobsElasticsearchIndexerInterface(f func(*Server) tjobs.IndexerJobInterface) {
	jobsElasticsearchIndexerInterface = f
}

func RegisterJobsBleveIndexerInterface(f func(*Server) tjobs.IndexerJobInterface) {
	jobsBleveIndexerInterface = f
}

func RegisterJobsMigrationsJobInterface(f func(*Server) tjobs.MigrationsJobInterface) {
	jobsMigrationsInterface = f
}

// RegisterJobsResendInvitationEmailInterface is used to register or initialize the jobsResendInvitationEmailInterface
func RegisterJobsResendInvitationEmailInterface(f func(*App) ejobs.ResendInvitationEmailJobInterface) {
	jobsResendInvitationEmailInterface = f
}

func RegisterCsvExportInterface(f func(*Server) tjobs.CsvExportInterface) {
	csvExportInterface = f
}

var jobsActiveUsersInterface func(*Server) tjobs.ActiveUsersJobInterface

func RegisterJobsActiveUsersInterface(f func(*Server) tjobs.ActiveUsersJobInterface) {
	jobsActiveUsersInterface = f
}

// enterprise jobs -----------------

var (
	complianceInterface    func(*Server) einterfaces.ComplianceInterface
	elasticsearchInterface func(*Server) searchengine.SearchEngineInterface
	clusterInterface       func(*Server) einterfaces.ClusterInterface
	dataRetentionInterface func(*Server) einterfaces.DataRetentionInterface
	metricsInterface       func(*Server) einterfaces.MetricsInterface
)

func RegisterComplianceInterface(f func(*Server) einterfaces.ComplianceInterface) {
	complianceInterface = f
}

func RegisterElasticsearchInterface(f func(*Server) searchengine.SearchEngineInterface) {
	elasticsearchInterface = f
}

func RegisterClusterInterface(f func(*Server) einterfaces.ClusterInterface) {
	clusterInterface = f
}

func RegisterDataRetentionInterface(f func(*Server) einterfaces.DataRetentionInterface) {
	dataRetentionInterface = f
}

func RegisterMetricsInterface(f func(*Server) einterfaces.MetricsInterface) {
	metricsInterface = f
}

func (s *Server) initEnterprise() {
	if metricsInterface != nil {
		s.Metrics = metricsInterface(s)
	}
	if complianceInterface != nil {
		s.Compliance = complianceInterface(s)
	}
	// if messageExportInterface != nil {
	// 	s.MessageExport = messageExportInterface(s)
	// }
	if dataRetentionInterface != nil {
		s.DataRetention = dataRetentionInterface(s)
	}
	if clusterInterface != nil {
		s.Cluster = clusterInterface(s)
	}
	if elasticsearchInterface != nil {
		s.SearchEngine.RegisterElasticsearchEngine(elasticsearchInterface(s))
	}
}
