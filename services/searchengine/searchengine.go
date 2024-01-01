package searchengine

import (
	"github.com/sitename/sitename/model_helper"
)

type Broker struct {
	cfg                 *model_helper.Config
	ElasticsearchEngine SearchEngineInterface
	BleveEngine         SearchEngineInterface
}

func NewBroker(cfg *model_helper.Config) *Broker {
	return &Broker{
		cfg: cfg,
	}
}

func (seb *Broker) RegisterElasticsearchEngine(es SearchEngineInterface) {
	seb.ElasticsearchEngine = es
}

func (seb *Broker) RegisterBleveEngine(be SearchEngineInterface) {
	seb.BleveEngine = be
}

func (seb *Broker) UpdateConfig(cfg *model_helper.Config) *model_helper.AppError {
	seb.cfg = cfg
	if seb.ElasticsearchEngine != nil {
		seb.ElasticsearchEngine.UpdateConfig(cfg)
	}

	if seb.BleveEngine != nil {
		seb.BleveEngine.UpdateConfig(cfg)
	}

	return nil
}

func (seb *Broker) GetActiveEngines() []SearchEngineInterface {
	engines := []SearchEngineInterface{}
	if seb.ElasticsearchEngine != nil && seb.ElasticsearchEngine.IsActive() {
		engines = append(engines, seb.ElasticsearchEngine)
	}
	if seb.BleveEngine != nil && seb.BleveEngine.IsActive() {
		engines = append(engines, seb.BleveEngine)
	}
	return engines
}
