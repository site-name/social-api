package jobs

import (
	"errors"

	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/services/configservice"
)

// Workers contains some works that can be ran
type Workers struct {
	ConfigService configservice.ConfigService
	Watcher       *Watcher

	workers map[string]model_helper.Worker

	listenerId string
	running    bool
}

var (
	ErrWorkersNotRunning    = errors.New("job workers are not running")
	ErrWorkersRunning       = errors.New("job workers are running")
	ErrWorkersUninitialized = errors.New("job workers are not initialized")
)

func NewWorkers(configService configservice.ConfigService) *Workers {
	return &Workers{
		ConfigService: configService,
		workers:       make(map[string]model_helper.Worker),
	}
}

func (workers *Workers) AddWorker(name string, worker model_helper.Worker) {
	workers.workers[name] = worker
}

func (workers *Workers) Get(name string) model_helper.Worker {
	return workers.workers[name]
}

// Start starts the workers. This call is not safe for concurrent use.
// Synchronization should be implemented by the caller.
func (workers *Workers) Start() {
	slog.Info("Starting workers")

	for _, w := range workers.workers {
		if w.IsEnabled(workers.ConfigService.Config()) {
			go w.Run()
		}
	}

	go workers.Watcher.Start()

	workers.listenerId = workers.ConfigService.AddConfigListener(workers.handleConfigChange)
	workers.running = true
}

func (workers *Workers) handleConfigChange(oldConfig *model_helper.Config, newConfig *model_helper.Config) {
	slog.Debug("Workers received config change.")

	for _, w := range workers.workers {
		if w.IsEnabled(oldConfig) && !w.IsEnabled(newConfig) {
			w.Stop()
		}
		if !w.IsEnabled(oldConfig) && w.IsEnabled(newConfig) {
			go w.Run()
		}
	}
}

// Stop stops the workers. This call is not safe for concurrent use.
// Synchronization should be implemented by the caller.
func (workers *Workers) Stop() {
	workers.ConfigService.RemoveConfigListener(workers.listenerId)

	workers.Watcher.Stop()

	for _, w := range workers.workers {
		if w.IsEnabled(workers.ConfigService.Config()) {
			w.Stop()
		}
	}

	workers.running = false

	slog.Info("Stopped workers")
}
