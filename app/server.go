package app

import (
	"hash/maphash"
	"net"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
	"github.com/sitename/sitename/config"
	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/jobs"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/services/cache"
	"github.com/sitename/sitename/services/telemetry"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/sqlstore"
)

const (
	SessionsCleanupBatchSize = 1000
)

type Server struct {
	// RootRouter is the starting point for all HTTP requests to the server.
	RootRouter                         *mux.Router
	sqlStore                           *sqlstore.SqlStore
	Store                              store.Store
	AppInitializedOnce                 sync.Once
	Server                             *http.Server
	ListenAddr                         *net.TCPAddr
	localModeServer                    *http.Server
	metricsServer                      *http.Server
	metricsRouter                      *mux.Router
	metricsLock                        sync.Mutex
	goroutineCount                     int32
	hashSeed                           maphash.Seed
	goroutineExitSignal                chan struct{}
	didFinishListen                    chan struct{}
	pushNotificationClient             *http.Client // TODO: move this to it's own package
	runEssentialJobs                   bool
	clusterLeaderListeners             sync.Map
	newStore                           func() (store.Store, error)
	configListenerId                   string
	licenseListenerId                  string
	logListenerId                      string
	clusterLeaderListenerId            string
	searchConfigListenerId             string
	searchLicenseListenerId            string
	loggerLicenseListenerId            string
	configStore                        *config.Store
	advancedLogListenerCleanup         func()
	pluginCommandsLock                 sync.RWMutex
	asymmetricSigningKey               atomic.Value
	clientConfig                       atomic.Value
	clientConfigHash                   atomic.Value
	limitedClientConfig                atomic.Value
	phase2PermissionsMigrationComplete bool
	joinCluster                        bool
	startMetrics                       bool
	startSearchEngine                  bool
	uploadLockMapMut                   sync.Mutex // These are used to prevent concurrent upload requests
	uploadLockMap                      map[string]bool
	featureFlagSynchronizer            *config.FeatureFlagSynchronizer
	featureFlagStop                    chan struct{}
	featureFlagStopped                 chan struct{}
	featureFlagSynchronizerMutex       sync.Mutex
	Log                                *slog.Logger
	NotificationsLog                   *slog.Logger
	sessionCache                       cache.Cache
	seenPendingPostIdsCache            cache.Cache
	statusCache                        cache.Cache
	telemetryService                   *telemetry.TelemetryService
	licenseValue                       atomic.Value
	Jobs                               *jobs.JobServer
	Metrics                            einterfaces.MetricsInterface
	RateLimiter                        *RateLimiter
	Busy                               *Busy
	EmailService                       *EmailService
}

// func NewServer(options ...Option) (*Server, error) {

// }

func (s *Server) runJobs() {
	s.Go(func() {
		runSecurityJob(s)
	})
	s.Go(func() {
		firstRun, err := s.getFirstServerRunTimestamp()
		if err != nil {
			slog.Warn("Fetching time of first server run failed. Setting to 'now'.")
			s.ensureFirstServerRunTimestamp()
			firstRun = util.MillisFromTime(time.Now())
		}
		s.telemetryService.RunTelemetryJob(firstRun)
	})
	s.Go(func() {
		runSessionCleanupJob(s)
	})
	s.Go(func() {
		runTokenCleanupJob(s)
	})
	// s.Go(func() {
	// 	runCommandWebhookCleanupJob(s)
	// })

	// if complianceI := s.Com
}

// Go creates a goroutine, but maintains a record of it to ensure that execution completes before
// the server is shutdown.
func (s *Server) Go(f func()) {
	atomic.AddInt32(&s.goroutineCount, 1)

	go func() {
		f()

		atomic.AddInt32(&s.goroutineCount, -1)
		select {
		case s.goroutineExitSignal <- struct{}{}:
		default:
		}
	}()
}

// func runCommandWebhookCleanupJob(s *Server) {
// 	doCommandWebhookCleanup(s)
// 	model.CreateRecurringTask("Command Hook Cleanup", func() {
// 		doCommandWebhookCleanup(s)
// 	}, time.Hour*1)
// }

// func doCommandWebhookCleanup(s *Server) {
// 	s.Store.CommandWebhook().Cleanup()
// }

func runTokenCleanupJob(s *Server) {
	doTokenCleanup(s)
	model.CreateRecurringTask("Token Cleanup", func() {
		doTokenCleanup(s)
	}, time.Hour*1)
}

func doTokenCleanup(s *Server) {
	s.Store.Token().Cleanup()
}

func runSecurityJob(s *Server) {
	doSecurity(s)
	model.CreateRecurringTask("Security", func() {
		doSecurity(s)
	}, time.Hour*4)
}

func runSessionCleanupJob(s *Server) {
	doSessionCleanup(s)
	model.CreateRecurringTask("Session Cleanup", func() {
		doSessionCleanup(s)
	}, time.Hour*24)
}

func doSessionCleanup(s *Server) {
	s.Store.Session().Cleanup(model.GetMillis(), SessionsCleanupBatchSize)
}

func doSecurity(s *Server) {
	s.DoSecurityUpdateCheck()
}

func (s *Server) TelemetryId() string {
	if s.telemetryService == nil {
		return ""
	}

	return s.telemetryService.TelemetryID
}

func (s *Server) License() *model.License {
	license, _ := s.licenseValue.Load().(*model.License)
	return license
}

func (s *Server) getFirstServerRunTimestamp() (int64, *model.AppError) {
	systemData, err := s.Store.System().GetByName(model.SYSTEM_FIRST_SERVER_RUN_TIMESTAMP_KEY)
	if err != nil {
		return 0, model.NewAppError("getFirstServerRunTimestamp", "app.system.get_by_name.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	value, err := strconv.ParseInt(systemData.Value, 10, 64)
	if err != nil {
		return 0, model.NewAppError("getFirstServerRunTimestamp", "app.system_install_date.parse_int.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return value, nil
}

func (s *Server) ensureFirstServerRunTimestamp() error {
	_, appErr := s.getFirstServerRunTimestamp()
	if appErr == nil {
		return nil
	}

	if err := s.Store.System().SaveOrUpdate(&model.System{
		Name:  model.SYSTEM_FIRST_SERVER_RUN_TIMESTAMP_KEY,
		Value: strconv.FormatInt(util.MillisFromTime(time.Now()), 10),
	}); err != nil {
		return err
	}
	return nil
}
