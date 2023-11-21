package interfaces

import (
	"github.com/sitename/sitename/model_helper"
)

type CsvExportInterface interface {
	MakeWorker() model_helper.Worker
}
