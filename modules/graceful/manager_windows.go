// +build windows

package graceful

import (
	"context"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/sitename/sitename/modules/log"
	"github.com/sitename/sitename/modules/setting"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
)

// WindowsServiceName is the name of the Windows service
var WindowsServiceName = "gitea"

const (
	hammerCode       = 128
	hammerCmd        = svc.Cmd(hammerCode)
	acceptHammerCode = svc.Accepted(hammerCode)
)

// Manager manages the graceful shutdown process
type Manager struct {
	ctx                    context.Context
	isChild                bool
	lock                   *sync.RWMutex
	state                  state
	shutdown               chan struct{}
	hammer                 chan struct{}
	terminate              chan struct{}
	done                   chan struct{}
	runningServerWaitGroup sync.WaitGroup
	createServerWaitGroup  sync.WaitGroup
	terminateWaitGroup     sync.WaitGroup
	shutdownRequested      chan struct{}
}

func newGracefulManager(ctx context.Context) *Manager {
	manager := &Manager{
		isChild: false,
		lock:    &sync.RWMutex{},
		ctx:     ctx,
	}
	manager.createServerWaitGroup.Add(numberOfServersToCreate)
	manager.start()
	return manager
}

func (g *Manager) start() {
	// Make channels
	g.terminate = make(chan struct{})
	g.shutdown = make(chan struct{})
	g.hammer = make(chan struct{})
	g.done = make(chan struct{})
	g.shutdownRequested = make(chan struct{})

	// Set the running state
	g.setState(stateRunning)
	if skip, _ := strconv.ParseBool(os.Getenv("SKIP_MINWINSVC")); skip {
		return
	}

	// Make SVC process
	run := svc.Run
	isInteractive, err := svc.IsWindowsService()
	if err != nil {
		log.Error("Unable to ascertain if running as an Interactive Session: %v", err)
		return
	}
	if isInteractive {
		run = debug.Run
	}
	go func() {
		_ = run(WindowsServiceName, g)
	}()
}

// Execute makes Manager implement svc.Handler
func (g *Manager) Execute(args []string, changes <-chan svc.ChangeRequest, status chan<- svc.Status) (svcSpecificEC bool, exitCode uint32) {
	if setting.StartupTimeout > 0 {
		status <- svc.Status{State: svc.StartPending}
	} else {
		status <- svc.Status{State: svc.StartPending, WaitHint: uint32(setting.StartupTimeout / time.Millisecond)}
	}

	// Now need to wait for everything to start...
	if !g.awaitServer(setting.StartupTimeout) {
		return false, 1
	}

	// We need to implement some way of svc.AcceptParamChange/svc.ParamChange
	status <- svc.Status{
		State:   svc.Running,
		Accepts: svc.AcceptStop | svc.AcceptShutdown | acceptHammerCode,
	}

	waitTime := 30 * time.Second

loop:
	for {
		select {
		case <-g.ctx.Done():
			g.DoGracefulShutdown()
			waitTime += setting.GracefulHammerTime
			break loop
		case <-g.shutdownRequested:
			waitTime += setting.GracefulHammerTime
			break loop
		case change := <-changes:
			switch change.Cmd {
			case svc.Interrogate:
				status <- change.CurrentStatus
			case svc.Stop, svc.Shutdown:
				g.DoGracefulShutdown()
				waitTime += setting.GracefulHammerTime
				break loop
			case hammerCode:
				g.DoGracefulShutdown()
				g.DoImmediateHammer()
				break loop
			default:
				log.Debug("Unexpected control request: %v", change.Cmd)
			}
		}
	}
	status <- svc.Status{
		State:    svc.StopPending,
		WaitHint: uint32(waitTime / time.Millisecond),
	}

hammerLoop:
	for {
		select {
		case change := <-changes:
			switch change.Cmd {
			case svc.Interrogate:
				status <- change.CurrentStatus
			case svc.Stop, svc.Shutdown, hammerCmd:
				g.DoImmediateHammer()
				break hammerLoop
			default:
				log.Debug("Unexpected control request: %v", change.Cmd)
			}
		case <-g.hammer:
			break hammerLoop
		}
	}
	return false, 0
}

// DoImmediateHammer causes an immediate hammer
func (g *Manager) DoImmediateHammer() {
	g.doHammerTime(0 * time.Second)
}

// DoGracefulShutdown causes a graceful shutdown
func (g *Manager) DoGracefulShutdown() {
	g.lock.Lock()
	select {
	case <-g.shutdownRequested:
		g.lock.Unlock()
	default:
		close(g.shutdownRequested)
		g.lock.Unlock()
		g.doShutdown()
	}
}

// RegisterServer registers the running of a listening server.
// Any call to RegisterServer must be matched by a call to ServerDone
func (g *Manager) RegisterServer() {
	g.runningServerWaitGroup.Add(1)
}

func (g *Manager) awaitServer(limit time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		g.createServerWaitGroup.Wait()
	}()
	if limit > 0 {
		select {
		case <-c:
			return true // completed normally
		case <-time.After(limit):
			return false // timed out
		case <-g.IsShutdown():
			return false
		}
	} else {
		select {
		case <-c:
			return true // completed normally
		case <-g.IsShutdown():
			return false
		}
	}
}
