package searchengine

import (
	"github.com/mattermost/mattermost-server/jobs"
	"github.com/sitename/sitename/model"
)

type Broker struct {
	cfg                 *model.Config
	jobServer           *jobs.JobServer
	ElasticsearchEngine SearchEngineInterface
	BleveEngine         SearchEngineInterface
}
