package interfaces

import (
	"github.com/sitename/sitename/model"
)

type ImportProcessInterface interface {
	MakeWorker() model.Worker
}
