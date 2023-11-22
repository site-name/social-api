package jobs

import (
	"github.com/sitename/sitename/model_helper"
)

type ElasticsearchIndexerInterface interface {
	MakeWorker() model_helper.Worker
}

type ElasticsearchAggregatorInterface interface {
	MakeWorker() model_helper.Worker
	MakeScheduler() model_helper.Scheduler
}
