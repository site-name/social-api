package app

import (
	"github.com/sitename/sitename/einterfaces"
	ejobs "github.com/sitename/sitename/einterfaces/jobs"
	tjobs "github.com/sitename/sitename/modules/jobs/interfaces"
	"github.com/sitename/sitename/services/searchengine"
)

var accountMigrationInterface func(*App) einterfaces.AccountMigrationInterface

func RegisterAccountMigrationInterface(f func(*App) einterfaces.AccountMigrationInterface) {
	accountMigrationInterface = f
}

var clusterInterface func(*Server) einterfaces.ClusterInterface

func RegisterClusterInterface(f func(*Server) einterfaces.ClusterInterface) {
	clusterInterface = f
}

var complianceInterface func(*Server) einterfaces.ComplianceInterface

func RegisterComplianceInterface(f func(*Server) einterfaces.ComplianceInterface) {
	complianceInterface = f
}

var jobsLdapSyncInterface func(*App) ejobs.LdapSyncInterface

func RegisterJobsLdapSyncInterface(f func(*App) ejobs.LdapSyncInterface) {
	jobsLdapSyncInterface = f
}

var jobsDataRetentionJobInterface func(*Server) ejobs.DataRetentionJobInterface

func RegisterJobsDataRetentionJobInterface(f func(*Server) ejobs.DataRetentionJobInterface) {
	jobsDataRetentionJobInterface = f
}

var jobsMessageExportJobInterface func(*Server) ejobs.MessageExportJobInterface

func RegisterJobsMessageExportJobInterface(f func(*Server) ejobs.MessageExportJobInterface) {
	jobsMessageExportJobInterface = f
}

var jobsElasticsearchAggregatorInterface func(*Server) ejobs.ElasticsearchAggregatorInterface

func RegisterJobsElasticsearchAggregatorInterface(f func(*Server) ejobs.ElasticsearchAggregatorInterface) {
	jobsElasticsearchAggregatorInterface = f
}

var jobsElasticsearchIndexerInterface func(*Server) tjobs.IndexerJobInterface

func RegisterJobsElasticsearchIndexerInterface(f func(*Server) tjobs.IndexerJobInterface) {
	jobsElasticsearchIndexerInterface = f
}

var jobsBleveIndexerInterface func(*Server) tjobs.IndexerJobInterface

func RegisterJobsBleveIndexerInterface(f func(*Server) tjobs.IndexerJobInterface) {
	jobsBleveIndexerInterface = f
}

var jobsMigrationsInterface func(*Server) tjobs.MigrationsJobInterface

func RegisterJobsMigrationsJobInterface(f func(*Server) tjobs.MigrationsJobInterface) {
	jobsMigrationsInterface = f
}

var jobsResendInvitationEmailInterface func(*App) ejobs.ResendInvitationEmailJobInterface

// RegisterJobsResendInvitationEmailInterface is used to register or initialize the jobsResendInvitationEmailInterface
func RegisterJobsResendInvitationEmailInterface(f func(*App) ejobs.ResendInvitationEmailJobInterface) {
	jobsResendInvitationEmailInterface = f
}

var metricsInterface func(*Server) einterfaces.MetricsInterface

func RegisterMetricsInterface(f func(*Server) einterfaces.MetricsInterface) {
	metricsInterface = f
}

var elasticsearchInterface func(*Server) searchengine.SearchEngineInterface

func RegisterElasticsearchInterface(f func(*Server) searchengine.SearchEngineInterface) {
	elasticsearchInterface = f
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
	// if dataRetentionInterface != nil {
	// 	s.DataRetention = dataRetentionInterface(s)
	// }
	if clusterInterface != nil {
		s.Cluster = clusterInterface(s)
	}
	if elasticsearchInterface != nil {
		s.SearchEngine.RegisterElasticsearchEngine(elasticsearchInterface(s))
	}
}
