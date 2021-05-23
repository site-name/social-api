package interfaces

import (
	"github.com/sitename/sitename/model"
)

type CsvExportInterface interface {
	MakeWorker() model.Worker
}
