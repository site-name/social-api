package interfaces

import (
	"github.com/sitename/sitename/model_helper"
)

type IndexerJobInterface interface {
	MakeWorker() model_helper.Worker
}
