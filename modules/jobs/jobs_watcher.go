package jobs

import (
	"math/rand"
	"time"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
)

// Default polling interval for jobs termination.
// (Defining as `var` rather than `const` allows tests to lower the interval.)
var DefaultWatcherPollingInterval = 15000

type Watcher struct {
	srv     *JobServer
	workers *Workers

	stop            chan struct{}
	stopped         chan struct{}
	pollingInterval int
}

func (srv *JobServer) MakeWatcher(workers *Workers, pollingInterval int) *Watcher {
	return &Watcher{
		stop:            make(chan struct{}),
		stopped:         make(chan struct{}),
		pollingInterval: pollingInterval,
		workers:         workers,
		srv:             srv,
	}
}

func (watcher *Watcher) Start() {
	slog.Debug("Watcher Started")

	// Delay for some random number of milliseconds before starting to ensure that multiple
	// instances of the jobserver  don't poll at a time too close to each other.
	rand.Seed(time.Now().UnixNano())
	<-time.After(time.Duration(rand.Intn(watcher.pollingInterval)) * time.Millisecond)

	defer func() {
		slog.Debug("Watcher Finished")
		close(watcher.stopped)
	}()

	for {
		select {
		case <-watcher.stop:
			slog.Debug("Watcher: Received stop signal.")
			return
		case <-time.After(time.Duration(watcher.pollingInterval) * time.Millisecond):
			watcher.PollAndNotify()
		}
	}
}

func (watcher *Watcher) Stop() {
	slog.Debug("Watcher Stopping")
	close(watcher.stop)
	<-watcher.stopped

	watcher.stop = make(chan struct{})
	watcher.stopped = make(chan struct{})
}

// sitting there waiting for new
func (watcher *Watcher) PollAndNotify() {
	jobs, err := watcher.srv.Store.Job().GetAllByStatus(model.JobstatusPending.String())
	if err != nil {
		slog.Error("Error occured getting all pending statuses.", slog.Err(err))
		return
	}

	for _, job := range jobs {
		worker := watcher.workers.Get(job.Type.String())
		if worker != nil {
			select {
			case worker.JobChannel() <- *job:
			default:
			}
		}
	}
}
