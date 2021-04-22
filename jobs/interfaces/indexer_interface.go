package interfaces

import (
	"github.com/sitename/sitename/model"
)

type IndexerJobInterface interface {
	MakeWorker() model.Worker
}
