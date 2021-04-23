package app

import (
	"fmt"
	"hash/maphash"
	"net"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	sentry "github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/config"
	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/jobs"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/i18n"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/templates"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/services/cache"
	"github.com/sitename/sitename/services/httpservice"
	"github.com/sitename/sitename/services/imageproxy"
	"github.com/sitename/sitename/services/searchengine"

	"github.com/sitename/sitename/services/searchengine/bleveengine"
	"github.com/sitename/sitename/services/telemetry"
	"github.com/sitename/sitename/services/tracing"
	"github.com/sitename/sitename/store"

	"github.com/sitename/sitename/store/localcachelayer"
	"github.com/sitename/sitename/store/retrylayer"
	"github.com/sitename/sitename/store/searchlayer"
	"github.com/sitename/sitename/store/sqlstore"
)

// declaring this as var to allow overriding in tests
var SentryDSN = "placeholder_sentry_dsn"

const (
	SessionsCleanupBatchSize = 1000
)

type Server struct {
	// RootRouter is the starting point for all HTTP requests to the server.
	RootRouter         *mux.Router
	sqlStore           *sqlstore.SqlStore
	Store              store.Store
	AppInitializedOnce sync.Once

	Server      *http.Server
	ListenAddr  *net.TCPAddr
	RateLimiter *RateLimiter
	Busy        *Busy

	localModeServer *http.Server

	metricsServer *http.Server
	metricsRouter *mux.Router
	metricsLock   sync.Mutex

	goroutineCount      int32
	hashSeed            maphash.Seed
	goroutineExitSignal chan struct{}

	didFinishListen chan struct{}

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
	EmailService                       *EmailService
	htmlTemplateWatcher                *templates.Container
	Ldap                               einterfaces.LdapInterface
	tracer                             *tracing.Tracer
	HTTPService                        httpservice.HTTPService
	ImageProxy                         *imageproxy.ImageProxy
	SearchEngine                       *searchengine.Broker
	CacheProvider                      cache.Provider
	Cluster                            einterfaces.ClusterInterface
}

func NewServer(options ...Option) (*Server, error) {
	rootRouter := mux.NewRouter()

	s := &Server{
		goroutineExitSignal: make(chan struct{}, 1),
		RootRouter:          rootRouter,
		hashSeed:            maphash.MakeSeed(),
		uploadLockMap:       map[string]bool{},
	}

	for _, option := range options {
		if err := option(s); err != nil {
			return nil, errors.Wrap(err, "failed to apply option")
		}
	}

	if s.configStore == nil {
		innerStore, err := config.NewFileStore("config.json", true)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load config")
		}
		configStoree, err := config.NewStoreFromBacking(innerStore, nil, false)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load config")
		}

		s.configStore = configStoree
	}

	if err := s.initLogging(); err != nil {
		slog.Error("Could not initiate logging", slog.Err(err))
	}

	// This is called after initLogging() to avoid a race condition.
	slog.Info("Server is initializing...", slog.String("go_version", runtime.Version()))

	// It is important to initialize the hub only after the global logger is set
	// to avoid race conditions while logging from inside the hub.
	// fakeApp := New(ServerConnector(s))
	// fakeApp.HubStart()

	if *s.Config().LogSettings.EnableDiagnostics && *s.Config().LogSettings.EnableSentry {
		if strings.Contains(SentryDSN, "placeholder") {
			slog.Warn("Sentry reporting is enabled, but SENTRY_DSN is not set. Disabling reporting.")
		} else {
			if err := sentry.Init(sentry.ClientOptions{
				Dsn:              SentryDSN,
				Release:          model.BuildHash,
				AttachStacktrace: true,
				BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
					// sanitize data sent to sentry to reduce exposure of PII
					if event.Request != nil {
						event.Request.Cookies = ""
						event.Request.QueryString = ""
						event.Request.Headers = nil
						event.Request.Data = ""
					}
					return event
				},
			}); err != nil {
				slog.Warn("Sentry could not be initiated, propably bad DSN?", slog.Err(err))
			}
		}
	}

	if *s.Config().ServiceSettings.EnableOpenTracing {
		tracer, err := tracing.New()
		if err != nil {
			return nil, err
		}
		s.tracer = tracer
	}

	s.HTTPService = httpservice.MakeHTTPService(s)
	s.pushNotificationClient = s.HTTPService.MakeClient(true)
	s.ImageProxy = imageproxy.MakeImageProxy(s, s.HTTPService, s.Log)

	if err := util.TranslationsPreInit(); err != nil {
		return nil, errors.Wrapf(err, "unable to load Sitename translation files")
	}
	model.AppErrorInit(i18n.T)

	searchEngine := searchengine.NewBroker(s.Config(), s.Jobs)
	bleveEngine := bleveengine.NewBleveEngine(s.Config(), s.Jobs)
	if err := bleveEngine.Start(); err != nil {
		return nil, err
	}
	searchEngine.RegisterBleveEngine(bleveEngine)
	s.SearchEngine = searchEngine

	// at the moment we only have this implementation
	// in the future the cache provider will be built based on the loaded config
	s.CacheProvider = cache.NewProvider()
	if err := s.CacheProvider.Connect(); err != nil {
		return nil, errors.Wrapf(err, "Unable to connect to cache provider")
	}

	var err error
	if s.sessionCache, err = s.CacheProvider.NewCache(&cache.CacheOptions{
		Size:           model.SESSION_CACHE_SIZE,
		Striped:        true,
		StripedBuckets: util.MaxInt(runtime.NumCPU()-1, 1),
	}); err != nil {
		return nil, errors.Wrap(err, "Unable to create session cache")
	}
	if s.seenPendingPostIdsCache, err = s.CacheProvider.NewCache(&cache.CacheOptions{
		Size:           model.STATUS_CACHE_SIZE,
		Striped:        true,
		StripedBuckets: util.MaxInt(runtime.NumCPU()-1, 1),
	}); err != nil {
		return nil, errors.Wrap(err, "Unable to create status cache")
	}

	// s.createPushNotificationsHub()

	if err2 := i18n.InitTranslations(*s.Config().LocalizationSettings.DefaultClientLocale, *s.Config().LocalizationSettings.DefaultClientLocale); err2 != nil {
		return nil, errors.Wrapf(err2, "unable to load Sitename translation files")
	}

	// s.initEnterprise()

	if s.newStore == nil {
		s.newStore = func() (store.Store, error) {
			s.sqlStore = sqlstore.New(s.Config().SqlSettings, s.Metrics)
			if s.sqlStore.DriverName() == model.DATABASE_DRIVER_POSTGRES {
				ver, err2 := s.sqlStore.GetDbVersion(true)
				if err != nil {
					return nil, errors.Wrap(err2, "cannot get DB version")
				}
				intVer, err2 := strconv.Atoi(ver)
				if err2 != nil {
					return nil, errors.Wrap(err2, "cannot parse DB version")
				}
				if intVer < sqlstore.MinimumRequiredPostgresVersion {
					return nil, fmt.Errorf("minimum required postgres version is %s; found %s", sqlstore.VersionString(sqlstore.MinimumRequiredPostgresVersion), sqlstore.VersionString(intVer))
				}
			}

			lcl, err2 := localcachelayer.NewLocalCacheLayer(
				retrylayer.New(s.sqlStore),
				s.Metrics,
				s.Cluster,
				s.CacheProvider,
			)
			if err2 != nil {
				return nil, errors.Wrap(err2, "cannot create local cache layer")
			}

			searchStore := searchlayer.NewSearchLayer
		}
	}

	return nil, nil
}

// initLogging initializes and configures the logger. This may be called more than once.
func (s *Server) initLogging() error {
	if s.Log == nil {
		s.Log = slog.NewLogger(util.MloggerConfigFromLoggerConfig(&s.Config().LogSettings, util.GetLogFileLocation))
	}

	// Use this app logger as the global logger (eventually remove all instances of global logging).
	// This is deferred because a copy is made of the logger and it must be fully configured before
	// the copy is made.
	defer slog.InitGlobalLogger(s.Log)

	// Redirect default Go logger to this logger.
	defer slog.RedirectStdLog(s.Log)

	if s.NotificationsLog == nil {
		notificationLogSettings := util.GetLogSettingsFromNotificationsLogSettings(&s.Config().NotificationLogSettings)
		s.NotificationsLog = slog.NewLogger(util.MloggerConfigFromLoggerConfig(notificationLogSettings, util.GetNotificationsLogFileLocation)).
			WithCallerSkip(1).With(slog.String("logSource", "notifications"))
	}

	if s.logListenerId != "" {
		s.RemoveConfigListener(s.logListenerId)
	}
	s.logListenerId = s.AddConfigListener(func(_, after *model.Config) {
		s.Log.ChangeLevels(util.MloggerConfigFromLoggerConfig(&after.LogSettings, util.GetLogFileLocation))

		notificationLogSettings := util.GetLogSettingsFromNotificationsLogSettings(&after.NotificationLogSettings)
		s.NotificationsLog.ChangeLevels(util.MloggerConfigFromLoggerConfig(notificationLogSettings, util.GetNotificationsLogFileLocation))
	})

	// Configure advanced logging.
	// Advanced logging is E20 only, however logging must be initialized before the license
	// file is loaded.  If no valid E20 license exists then advanced logging will be
	// shutdown once license is loaded/checked.
	if *s.Config().LogSettings.AdvancedLoggingConfig != "" {
		dsn := *s.Config().LogSettings.AdvancedLoggingConfig
		isJson := config.IsJsonMap(dsn)

		// If this is a file based config we need the full path so it can be watched.
		if !isJson && strings.HasPrefix(s.configStore.String(), "file://") && !filepath.IsAbs(dsn) {
			configPath := strings.TrimPrefix(s.configStore.String(), "file://")
			dsn = filepath.Join(filepath.Dir(configPath), dsn)
		}

		cfg, err := config.NewLogConfigSrc(dsn, isJson, s.configStore)
		if err != nil {
			return fmt.Errorf("invalid advanced logging config, %w", err)
		}

		if err := s.Log.ConfigAdvancedLogging(cfg.Get()); err != nil {
			return fmt.Errorf("error configuring advanced logging, %w", err)
		}

		if !isJson {
			slog.Info("Loaded advanced logging config", slog.String("source", dsn))
		}

		listenerId := cfg.AddListener(func(_, newCfg slog.LogTargetCfg) {
			if err := s.Log.ConfigAdvancedLogging(newCfg); err != nil {
				slog.Error("Error re-configuring advanced logging", slog.Err(err))
			} else {
				slog.Info("Re-configured advanced logging")
			}
		})

		// In case initLogging is called more than once.
		if s.advancedLogListenerCleanup != nil {
			s.advancedLogListenerCleanup()
		}

		s.advancedLogListenerCleanup = func() {
			cfg.RemoveListener(listenerId)
		}
	}
	return nil
}

func (s *Server) TemplatesContainer() *templates.Container {
	return s.htmlTemplateWatcher
}

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
