package plugin

import (
	"fmt"
	"sync"

	"github.com/hashicorp/go-plugin"
	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model/plugins"
	"github.com/sitename/sitename/modules/slog"
)

type supervisor struct {
	lock        sync.RWMutex
	client      *plugin.Client
	hooks       Hooks
	implemented [TotalHooksID]bool
	pid         int
}

func newSupervisor(pluginInfo *plugins.BundleInfo, apiImpl API, driver Driver, parentLogger *slog.Logger, metrics einterfaces.MetricsInterface) (retSupervisor *supervisor, retErr error) {
	// sup := supervisor{}
	// defer func() {
	// 	if retErr != nil {
	// 		sup.Shutdown()
	// 	}
	// }()

	// wrappedLogger := pluginInfo.WrapLogger(parentLogger)

	// hclogAdaptedLogger := &hclogAdapter{
	// 	wrappedLogger: wrappedLogger.WithCallerSkip(1),
	// 	extrasKey:     "wrapped_extras",
	// }

	// pluginMap := map[string]plugin.Plugin{
	// 	"hooks": &hooksPlugin{
	// 		log:        wrappedLogger,
	// 		driverImpl: driver,
	// 		apiImpl: &apiTimerLayer{
	// 			pluginInfo.Manifest.Id,
	// 			apiImpl,
	// 			metrics,
	// 		},
	// 	},
	// }
}

func (sup *supervisor) Shutdown() {
	sup.lock.RLock()
	defer sup.lock.RUnlock()
	if sup.client != nil {
		sup.client.Kill()
	}
}

func (sup *supervisor) Hooks() Hooks {
	sup.lock.RLock()
	defer sup.lock.RUnlock()
	return sup.hooks
}

// PerformHealthCheck checks the plugin through an an RPC ping.
func (sup *supervisor) PerformHealthCheck() error {
	// No need for a lock here because Ping is read-locked.
	if pingErr := sup.Ping(); pingErr != nil {
		for pingFails := 1; pingFails < HealthCheckNumRestartsLimit; pingFails++ {
			pingErr = sup.Ping()
			if pingErr == nil {
				break
			}
		}
		if pingErr != nil {
			return fmt.Errorf("plugin RPC connection is not responding")
		}
	}

	return nil
}

// Ping checks that the RPC connection with the plugin is alive and healthy.
func (sup *supervisor) Ping() error {
	sup.lock.RLock()
	defer sup.lock.RUnlock()
	client, err := sup.client.Client()
	if err != nil {
		return err
	}

	return client.Ping()
}

func (sup *supervisor) Implements(hookId int) bool {
	sup.lock.RLock()
	defer sup.lock.RUnlock()
	return sup.implemented[hookId]
}
