package csv

import (
	"sync"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
)

type ServiceCsv struct {
	srv *app.Server

	sync.WaitGroup
	sync.Mutex
}

type ServiceCsvConfig struct {
	Server *app.Server
}

func NewServiceCsv(config *ServiceCsvConfig) sub_app_iface.CsvService {
	return &ServiceCsv{
		srv: config.Server,
	}
}
