package app

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"hash/maphash"
	"html/template"
	"net"
	"net/http"
	"net/http/pprof"
	"net/url"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	sentry "github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"github.com/site-name/decimal"
	"github.com/sitename/sitename/app/email"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/audit"
	"github.com/sitename/sitename/modules/config"
	"github.com/sitename/sitename/modules/i18n"
	"github.com/sitename/sitename/modules/jobs"
	"github.com/sitename/sitename/modules/jobs/active_users"
	"github.com/sitename/sitename/modules/mail"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/sitename/sitename/modules/plugin"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/templates"
	"github.com/sitename/sitename/modules/timezones"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/modules/util/api"
	"github.com/sitename/sitename/services/awsmeter"
	"github.com/sitename/sitename/services/cache"
	"github.com/sitename/sitename/services/httpservice"
	"github.com/sitename/sitename/services/imageproxy"
	"github.com/sitename/sitename/services/searchengine"
	"github.com/sitename/sitename/services/searchengine/bleveengine"
	"github.com/sitename/sitename/services/searchengine/bleveengine/indexer"
	"github.com/sitename/sitename/services/tracing"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/localcachelayer"
	"github.com/sitename/sitename/store/retrylayer"
	"github.com/sitename/sitename/store/searchlayer"
	"github.com/sitename/sitename/store/sqlstore"
	"github.com/sitename/sitename/store/timerlayer"
	"golang.org/x/crypto/acme/autocert"
)

// declaring this as var to allow overriding in tests
var SentryDSN = "placeholder_sentry_dsn"

const (
	SessionsCleanupBatchSize                        = 1000
	TimeToWaitForConnectionsToCloseOnServerShutdown = time.Second
)

type Server struct {
	sqlStore *sqlstore.SqlStore
	Store    store.Store
	// WebSocketRouter *WebSocketRouter

	// RootRouter is the starting point for all HTTP requests to the server.
	RootRouter *mux.Router
	// Router is the starting point for all web, api4 and ws requests to the server. It differs
	// from RootRouter only if the SiteURL contains a /subpath.
	Router *mux.Router

	Server      *http.Server
	ListenAddr  *net.TCPAddr
	RateLimiter *RateLimiter
	Busy        *Busy

	metricsServer *http.Server
	metricsRouter *mux.Router
	metricsLock   sync.Mutex

	didFinishListen chan struct{}

	goroutineCount      int32
	goroutineExitSignal chan struct{}

	PluginsEnvironment     *plugin.Environment
	PluginConfigListenerId string
	PluginsLock            sync.RWMutex

	EmailService *email.Service

	// hubs     []*Hub
	hashSeed maphash.Seed

	// PushNotificationsHub   PushNotificationsHub
	pushNotificationClient *http.Client // TODO: move this to it's own package

	runEssentialJobs bool
	Jobs             *jobs.JobServer

	clusterLeaderListeners sync.Map

	timezones *timezones.Timezones

	newStore func() (store.Store, error)

	htmlTemplateWatcher *templates.Container
	// seenPendingPostIdsCache cache.Cache
	// licenseListenerId       string
	// searchLicenseListenerId string
	// loggerLicenseListenerId string
	openGraphDataCache      cache.Cache
	configListenerId        string
	clusterLeaderListenerId string
	searchConfigListenerId  string
	ConfigStore             *config.Store
	postActionCookieSecret  []byte

	// pluginCommands     []*PluginCommand
	// pluginCommandsLock sync.RWMutex

	asymmetricSigningKey atomic.Value
	clientConfig         atomic.Value
	clientConfigHash     atomic.Value
	limitedClientConfig  atomic.Value

	// telemetryService *telemetry.TelemetryService
	// serviceMux sync.RWMutex
	// remoteClusterService remotecluster.RemoteClusterServiceIFace
	// sharedChannelService SharedChannelServiceIFace

	phase2PermissionsMigrationComplete bool

	HTTPService httpservice.HTTPService
	ImageProxy  *imageproxy.ImageProxy

	Audit            *audit.Audit
	Log              *slog.Logger
	NotificationsLog *slog.Logger

	joinCluster       bool
	startMetrics      bool
	startSearchEngine bool
	// skipPostInit      bool

	SearchEngine *searchengine.Broker

	AccountMigration einterfaces.AccountMigrationInterface
	Cluster          einterfaces.ClusterInterface
	Compliance       einterfaces.ComplianceInterface
	DataRetention    einterfaces.DataRetentionInterface
	Ldap             einterfaces.LdapInterface
	Metrics          einterfaces.MetricsInterface
	Saml             einterfaces.SamlInterface
	// MessageExport    einterfaces.MessageExportInterface
	// Cloud            einterfaces.CloudInterface
	// Notification     einterfaces.NotificationInterface

	CacheProvider cache.Provider

	tracer *tracing.Tracer

	// featureFlagSynchronizer      *featureflag.Synchronizer
	featureFlagStop              chan struct{}
	featureFlagStopped           chan struct{}
	featureFlagSynchronizerMutex sync.Mutex

	ExchangeRateMap sync.Map // this is cache for storing currency exchange rates. Keys are strings, values are float64

	// these are sub services
	Account   sub_app_iface.AccountService
	Order     sub_app_iface.OrderService
	Payment   sub_app_iface.PaymentService
	Giftcard  sub_app_iface.GiftcardService
	Checkout  sub_app_iface.CheckoutService
	Product   sub_app_iface.ProductService
	Warehouse sub_app_iface.WarehouseService
	Wishlist  sub_app_iface.WishlistService
	Webhook   sub_app_iface.WebhookService
	Shipping  sub_app_iface.ShippingService
	Discount  sub_app_iface.DiscountService
	Menu      sub_app_iface.MenuService
	Csv       sub_app_iface.CsvService
	Page      sub_app_iface.PageService
	Seo       sub_app_iface.SeoService
	Attribute sub_app_iface.AttributeService
	Channel   sub_app_iface.ChannelService
	Invoice   sub_app_iface.InvoiceService
	File      sub_app_iface.FileService
	Plugin    sub_app_iface.PluginService
	Shop      sub_app_iface.ShopService
}

// NewServer create new system server
func NewServer(options ...Option) (*Server, error) {
	s := &Server{
		goroutineExitSignal: make(chan struct{}, 1),
		RootRouter:          mux.NewRouter(),
		hashSeed:            maphash.MakeSeed(),
	}

	for _, option := range options {
		if err := option(s); err != nil {
			return nil, errors.Wrap(err, "failed to apply option")
		}
	}

	if s.ConfigStore == nil {
		innerStore, err := config.NewFileStore("config.json", true)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load config")
		}
		configStoree, err := config.NewStoreFromBacking(innerStore, nil, false)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load config")
		}

		s.ConfigStore = configStoree
	}

	if err := s.initLogging(); err != nil {
		slog.Error("Could not initiate logging", slog.Err(err))
	}

	// This is called after initLogging() to avoid a race condition.
	slog.Info("Server is initializing...", slog.String("go_version", runtime.Version()))

	// It is important to initialize the hub only after the global logger is set
	// to avoid race conditions while logging from inside the hub.
	// app := New(ServerConnector(s))
	// app.HubStart()

	if *s.Config().LogSettings.EnableDiagnostics && *s.Config().LogSettings.EnableSentry {
		if strings.Contains(SentryDSN, "placeholder") {
			slog.Warn("Sentry reporting is enabled, but SENTRY_DSN is not set. Disabling reporting.")
		} else {
			if err := sentry.Init(sentry.ClientOptions{
				Dsn:              SentryDSN,
				Release:          model_helper.BuildHash,
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
	model_helper.AppErrorInit(i18n.T)

	searchEngine := searchengine.NewBroker(s.Config())
	bleveEngine := bleveengine.NewBleveEngine(s.Config(), s.Jobs)
	if err := bleveEngine.Start(); err != nil {
		return nil, err
	}
	searchEngine.RegisterBleveEngine(bleveEngine)
	s.SearchEngine = searchEngine

	// at the moment we only have this implementation
	// in the future the cache provider will be built based on the loaded config
	// this must be created before registering sub services
	s.CacheProvider = cache.NewProvider()
	if err := s.CacheProvider.Connect(); err != nil {
		return nil, errors.Wrapf(err, "Unable to connect to cache provider")
	}

	var err error
	if s.openGraphDataCache, err = s.CacheProvider.NewCache(&cache.CacheOptions{
		Size: openGraphMetadataCacheSize,
	}); err != nil {
		return nil, errors.Wrap(err, "Unable to create opengraphdata cache")
	}

	// s.createPushNotificationsHub()

	if err2 := i18n.InitTranslations(s.Config().LocalizationSettings.DefaultServerLocale.String(), s.Config().LocalizationSettings.DefaultClientLocale.String()); err2 != nil {
		return nil, errors.Wrapf(err2, "unable to load Mattermost translation files")
	}

	s.initEnterprise()

	if s.newStore == nil {
		s.newStore = func() (store.Store, error) {
			s.sqlStore = sqlstore.New(s.Config().SqlSettings, s.Metrics)

			lcl, err2 := localcachelayer.NewLocalCacheLayer(
				retrylayer.New(s.sqlStore),
				s.Metrics,
				s.Cluster,
				s.CacheProvider,
			)
			if err2 != nil {
				return nil, errors.Wrap(err2, "cannot create local cache layer")
			}

			searchStore := searchlayer.NewSearchLayer(
				lcl,
				s.SearchEngine,
				s.Config(),
			)

			s.AddConfigListener(func(prevCfg, cfg *model_helper.Config) {
				searchStore.UpdateConfig(cfg)
			})

			return timerlayer.New(
				searchStore,
				s.Metrics,
			), nil
		}
	}

	templatesDir, ok := templates.GetTemplateDirectory()
	if !ok {
		return nil, errors.New("Failed find server templates in \"templates\" directory or SN_SERVER_PATH")
	}
	htmlTemplateWatcher, errorsChan, err2 := templates.NewWithWatcher(templatesDir)
	if err2 != nil {
		return nil, errors.Wrap(err2, "cannot initialize server templates")
	}
	s.Go(func() {
		for err2 := range errorsChan {
			slog.Warn("Server templates error", slog.Err(err2))
		}
	})
	s.htmlTemplateWatcher = htmlTemplateWatcher

	s.Store, err = s.newStore()
	if err != nil {
		return nil, errors.Wrap(err, "cannot create store")
	}

	// This enterprise init should happen after the store is set
	// but we don't want to move the s.initEnterprise() call because
	// we had side-effects with that in the past and needs further
	// investigation
	// if cloudInterface != nil {
	// 	s.Cloud = cloudInterface(s)
	// }

	// s.telemetryService = telemetry.New(s, s.Store, s.SearchEngine, s.Log)

	emailService, err := email.NewService(email.ServiceConfig{
		ConfigFn:          s.Config,
		GoFn:              s.Go,
		TemplateContainer: s.TemplatesContainer(),
		Store:             s.Store,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "unable to initialize email service")
	}
	s.EmailService = emailService

	// s.setupFeatureFlags()

	// initialize job server
	s.initJobs()

	s.clusterLeaderListenerId = s.AddClusterLeaderChangedListener(func() {
		slog.Info("Cluster leader changed. Determining if job schedulers should be running:", slog.Bool("isLeader", s.IsLeader()))
		if s.Jobs != nil {
			s.Jobs.HandleClusterLeaderChange(s.IsLeader())
		}
		// s.setupFeatureFlags()
	})

	if s.joinCluster && s.Cluster != nil {
		s.registerClusterHandlers()
		s.Cluster.StartInterNodeCommunication()
	}

	if err = s.ensureAsymmetricSigningKey(); err != nil {
		return nil, errors.Wrapf(err, "unable to ensure asymmetric signing key")
	}

	if err = s.ensurePostActionCookieSecret(); err != nil {
		return nil, errors.Wrapf(err, "unable to ensure PostAction cookie secret")
	}

	if err = s.ensureInstallationDate(); err != nil {
		return nil, errors.Wrapf(err, "unable to ensure installation date")
	}

	if err = s.ensureFirstServerRunTimestamp(); err != nil {
		return nil, errors.Wrapf(err, "unable to ensure first run timestamp")
	}

	s.regenerateClientConfig()

	subPath, err := model_helper.GetSubpathFromConfig(s.Config())
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse SiteURL subpath")
	}
	s.Router = s.RootRouter.PathPrefix(subPath).Subrouter()

	pluginsRoute := s.Router.PathPrefix("/plugins/{plugin_id:[A-Za-z0-9\\_\\-\\.]+}").Subrouter()
	pluginsRoute.HandleFunc("", s.ServePluginRequest)
	pluginsRoute.HandleFunc("/public/{public_file:.*}", s.ServePluginPublicRequest)
	pluginsRoute.HandleFunc("/{anything:.*}", s.ServePluginRequest)

	// If configured with a subpath, redirect 404s at the root back into the subpath.
	if subPath != "/" {
		s.RootRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.URL.Path = path.Join(subPath, r.URL.Path)
			http.Redirect(w, r, r.URL.String(), http.StatusFound)
		})
	}

	// s.WebSocketRouter = &WebSocketRouter{
	// 	handlers: make(map[string]webSocketHandler),
	// 	app:      fakeApp,
	// }

	mailConfig := s.MailServiceConfig()
	if nErr := mail.TestConnection(mailConfig); nErr != nil {
		slog.Error("Mail server connection test is failed", slog.Err(nErr))
	}

	if _, err = url.ParseRequestURI(*s.Config().ServiceSettings.SiteURL); err != nil {
		slog.Error("SiteURL must be set. Some features will operate incorrectly if the SiteURL is not set.")
	}

	s.timezones = timezones.New()
	// Start email batching because it's not like the other jobs
	s.AddConfigListener(func(_, _ *model_helper.Config) {
		s.EmailService.InitEmailBatching()
	})

	// Start plugin health check job
	pluginsEnvironment := s.PluginsEnvironment
	if pluginsEnvironment != nil {
		pluginsEnvironment.InitPluginHealthCheckJob(*s.Config().PluginSettings.Enable && *s.Config().PluginSettings.EnableHealthCheck)
	}
	s.AddConfigListener(func(_, c *model_helper.Config) {
		s.PluginsLock.RLock()
		pluginsEnvironment := s.PluginsEnvironment
		s.PluginsLock.RUnlock()
		if pluginsEnvironment != nil {
			pluginsEnvironment.InitPluginHealthCheckJob(*s.Config().PluginSettings.Enable && *c.PluginSettings.EnableHealthCheck)
		}
	})

	logCurrentVersion := fmt.Sprintf(
		"Current version is %v (%v/%v/%v/%v)",
		model_helper.CurrentVersion,
		model_helper.BuildNumber,
		model_helper.BuildDate,
		model_helper.BuildHash,
		model_helper.BuildHashEnterprise,
	)
	slog.Info(
		logCurrentVersion,
		slog.String("current_version", model_helper.CurrentVersion),
		slog.String("build_number", model_helper.BuildNumber),
		slog.String("build_date", model_helper.BuildDate),
		slog.String("build_hash", model_helper.BuildHash),
		slog.String("build_hash_enterprise", model_helper.BuildHashEnterprise),
	)

	pwd, _ := os.Getwd()
	slog.Info("Printing current working", slog.String("directory", pwd))
	slog.Info("Loading config", slog.String("source", s.ConfigStore.String()))

	s.checkPushNotificationServerUrl()

	s.ReloadConfig()

	if s.Audit == nil {
		s.Audit = &audit.Audit{}
		s.Audit.Init(audit.DefMaxQueueSize)
		if err = s.configureAudit(s.Audit, true); err != nil {
			slog.Error("Error configuring audit", slog.Err(err))
		}
	}

	s.enableLoggingMetrics()

	// Enable developer settings if this is a "dev" build
	if model_helper.BuildNumber == "dev" {
		s.UpdateConfig(func(cfg *model_helper.Config) { *cfg.ServiceSettings.EnableDeveloper = true })
	}

	if err = s.Store.Status().ResetAll(); err != nil {
		slog.Error("Error to reset the server status.", slog.Err(err))
	}

	if s.startMetrics {
		s.SetupMetricsServer()
	}

	s.SearchEngine.UpdateConfig(s.Config())
	s.searchConfigListenerId = s.StartSearchEngine()

	// app := New(ServerConnector(s))
	// c := request.EmptyContext()

	if s.runEssentialJobs {
		// s.Go(func() {
		// 	runCheckAdminSupportStatusJob(app, c)
		// 	runDNDStatusExpireJob(app)
		// })
		s.runJobs()
	}

	// register all sub services, must go after store creation
	if err = s.registerSubServices(); err != nil {
		return nil, err
	}

	s.doAppMigrations()

	return s, nil
}

func (s *Server) ClientConfigHash() string {
	return s.clientConfigHash.Load().(string)
}

func (s *Server) initJobs() {
	s.Jobs = jobs.NewJobServer(s, s.Store, s.Metrics)

	if jobsDataRetentionJobInterface != nil {
		builder := jobsDataRetentionJobInterface(s)
		s.Jobs.RegisterJobType(model.JobTypeDataRetention, builder.MakeWorker(), builder.MakeScheduler())
	}

	if jobsMessageExportJobInterface != nil {
		builder := jobsMessageExportJobInterface(s)
		s.Jobs.RegisterJobType(model.JobTypeMessageExport, builder.MakeWorker(), builder.MakeScheduler())
	}

	if jobsElasticsearchAggregatorInterface != nil {
		builder := jobsElasticsearchAggregatorInterface(s)
		s.Jobs.RegisterJobType(model.JobTypeElasticsearchPostAggregation, builder.MakeWorker(), builder.MakeScheduler())
	}

	if jobsElasticsearchIndexerInterface != nil {
		builder := jobsElasticsearchIndexerInterface(s)
		s.Jobs.RegisterJobType(model.JobTypeElasticsearchPostIndexing, builder.MakeWorker(), nil)
	}

	if jobsLdapSyncInterface != nil {
		builder := jobsLdapSyncInterface(New())
		s.Jobs.RegisterJobType(model.JobTypeLdapSync, builder.MakeWorker(), builder.MakeScheduler())
	}

	s.Jobs.RegisterJobType(
		model.JobTypeBlevePostIndexing,
		indexer.MakeWorker(s.Jobs, s.SearchEngine.BleveEngine.(*bleveengine.BleveEngine)),
		nil,
	)

	s.Jobs.RegisterJobType(
		model.JobTypeActiveUsers,
		active_users.MakeWorker(s.Jobs, s.Store, func() einterfaces.MetricsInterface { return s.Metrics }),
		active_users.MakeScheduler(s.Jobs),
	)

	// s.Jobs.RegisterJobType(
	// 	model.JobTypeMigrations,
	// 	migrations.MakeWorker(s.Jobs, s.Store),
	// 	migrations.MakeScheduler(s.Jobs, s.Store),
	// )

	// s.Jobs.RegisterJobType(
	// 	model.JobTypeResendInvitationEmail,
	// 	resend_invitation_email.MakeWorker(s.Jobs, New(ServerConnector(s), s.Store, s.telemetryService),
	// 		nil,
	// 	),
	// )
}

// func (s *Server) TelemetryId() string {
// 	if s.telemetryService == nil {
// 		return ""
// 	}
// 	return s.telemetryService.TelemetryID
// }

// initLogging initializes and configures the logger. This may be called more than once.
func (s *Server) initLogging() error {
	var err error
	if s.Log == nil {
		s.Log, err = slog.NewLogger()
		if err != nil {
			return err
		}
	}

	// create notification logger if needed
	if s.NotificationsLog == nil {
		l, err := slog.NewLogger()
		if err != nil {
			return err
		}
		s.NotificationsLog = l.With(slog.String("logSource", "notifications"))
	}

	if err := s.configureLogger("logging", s.Log, &s.Config().LogSettings, s.ConfigStore, config.GetLogFileLocation); err != nil {
		// if the config is locked then a unit test has already configured and locked the logger; not an error.
		if !errors.Is(err, slog.ErrConfigurationLock) {
			// revert to default logger if the config is invalid
			slog.InitGlobalLogger(nil)
			return err
		}
	}

	// Redirect default Go logger to app logger.
	s.Log.RedirectStdLog(slog.LvlStdLog)

	// Use the app logger as the global logger (eventually remove all instances of global logging).
	slog.InitGlobalLogger(s.Log)

	notificationLogSettings := config.GetLogSettingsFromNotificationsLogSettings(&s.Config().NotificationLogSettings)
	if err := s.configureLogger("notification logging", s.NotificationsLog, notificationLogSettings, s.ConfigStore, config.GetNotificationsLogFileLocation); err != nil {
		if !errors.Is(err, slog.ErrConfigurationLock) {
			slog.Error("Error configuring notification logger", slog.Err(err))
			return err
		}
	}
	return nil
}

// configureLogger applies the specified configuration to a logger.
func (s *Server) configureLogger(name string, logger *slog.Logger, logSettings *model_helper.LogSettings, configStore *config.Store, getPath func(string) string) error {
	// Advanced logging is E20 only, however logging must be initialized before the license
	// file is loaded.  If no valid E20 license exists then advanced logging will be
	// shutdown once license is loaded/checked.
	var err error
	dsn := *logSettings.AdvancedLoggingConfig
	var logConfigSrc config.LogConfigSrc
	if dsn != "" {
		logConfigSrc, err = config.NewLogConfigSrc(dsn, configStore)
		if err != nil {
			return fmt.Errorf("invalid config source for %s, %w", name, err)
		}
		slog.Info("Loaded configuration for "+name, slog.String("source", dsn))
	}

	cfg, err := config.MloggerConfigFromLoggerConfig(logSettings, logConfigSrc, getPath)
	if err != nil {
		return fmt.Errorf("invalid config source for %s, %w", name, err)
	}

	if err := logger.ConfigureTargets(cfg, nil); err != nil {
		return fmt.Errorf("invalid config for %s, %w", name, err)
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
	// s.Go(func() {
	// 	firstRun, err := s.getFirstServerRunTimestamp()
	// 	if err != nil {
	// 		slog.Warn("Fetching time of first server run failed. Setting to 'now'.")
	// 		s.ensureFirstServerRunTimestamp()
	// 		firstRun = util.MillisFromTime(time.Now())
	// 	}
	// 	s.telemetryService.RunTelemetryJob(firstRun)
	// })
	// s.Go(func() {
	// 	runCommandWebhookCleanupJob(s)
	// })
	s.Go(func() {
		runSessionCleanupJob(s)
	})
	s.Go(func() {
		runTokenCleanupJob(s)
	})

	if s.Compliance != nil {
		s.Compliance.StartComplianceDailyJob()
	}

	if *s.Config().JobSettings.RunJobs && s.Jobs != nil {
		if err := s.Jobs.StartWorkers(); err != nil {
			slog.Error("Failed to start job server workers", slog.Err(err))
		}
	}
	if *s.Config().JobSettings.RunScheduler && s.Jobs != nil {
		if err := s.Jobs.StartSchedulers(); err != nil {
			slog.Error("Failed to start job server schedulers", slog.Err(err))
		}
	}

	if *s.Config().ServiceSettings.EnableAWSMetering {
		runReportToAWSMeterJob(s)
	}

	// check if we can run periodic task on fetching currency rate
	if setting := s.Config().ThirdPartySettings; setting.OpenExchangeRateApiKey != nil &&
		setting.OpenExchangeRecuringDurationHours != nil &&
		setting.OpenExchangeApiEndpoint != nil {
		runFetchingCurrencyExchangeRateJob(s, *setting.OpenExchangeRateApiKey, *setting.OpenExchangeRecuringDurationHours, *setting.OpenExchangeApiEndpoint)
	}
}

// runFetchingCurrencyExchangeRateJob every 2 hours it performs:
//
// Fetching exchange rates from external service, then upsate in a cache map and upsert in the database
func runFetchingCurrencyExchangeRateJob(s *Server, apiKey string, recuringHours int, apiEndPoint string) {
	var (
		client = s.HTTPService.MakeClient(true)
		params = url.Values{
			"app_id": []string{apiKey},
			"base":   []string{model_helper.DEFAULT_CURRENCY.String()}, // units other than USD require service subsciption
		}
		responseValue struct {
			Disclaimer string                     `json:"disclaimer,omitempty"`
			License    string                     `json:"license,omitempty"`
			TimeStamp  int64                      `json:"timestamp,omitempty"`
			Base       string                     `json:"base"`
			Rates      map[model.Currency]float64 `json:"rates"`
		}
	)

	apiEndPoint = apiEndPoint + "?" + params.Encode()

	fetchFun := func() {
		req, err := http.NewRequest(http.MethodGet, apiEndPoint, nil)
		if err != nil {
			s.Log.Error("Error creating http request to fetch currency exchange rate", slog.Err(err))
			return
		}
		response, err := client.Do(req)
		if err != nil {
			s.Log.Error("Error fetching exchange rates", slog.Err(err))
			return
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			s.Log.Error("Returned exchange response status code was not 200", slog.Int("status", response.StatusCode))
			return
		}
		// process data
		err = json.NewDecoder(response.Body).Decode(&responseValue)
		if err != nil {
			s.Log.Error("Error parsing currency exchange response body", slog.Err(err))
			return
		}

		var exchangeRateInstances model.OpenExchangeRateSlice
		for currency, rate := range responseValue.Rates {
			exchangeRate := &model.OpenExchangeRate{
				ToCurrency: currency,
				Rate:       model_types.NewNullDecimal(decimal.NewFromFloat(rate)),
			}
			s.ExchangeRateMap.Store(currency, exchangeRate)
			exchangeRateInstances = append(exchangeRateInstances, exchangeRate)
		}
		// update rates in database
		if err := s.upsertCurrencyExchangeRates(exchangeRateInstances); err != nil {
			s.Log.Error("Failed to upsert exchange rates", slog.Err(err))
			return
		}

		s.Log.Info("Successfully fetched and set currency exchange rates", slog.Any("base_currency", model_helper.DEFAULT_CURRENCY))
	}

	// first run
	fetchFun()
	model_helper.CreateRecurringTask("Collect and set currency exchange rates", fetchFun, time.Duration(recuringHours)*time.Hour)
}

func (s *Server) upsertCurrencyExchangeRates(newRates []*model.OpenExchangeRate) error {
	_, err := s.Store.OpenExchangeRate().BulkUpsert(newRates)
	return err
}

// Global app options that should be applied to apps created by this server
func (s *Server) AppOptions() []AppOption {
	return []AppOption{
		ServerConnector(s),
	}
}

// Return Database type (postgres or mysql) and current version of Mattermost
func (s *Server) DatabaseTypeAndMattermostVersion() (string, string) {
	mattermostVersion, _ := s.Store.System().GetByName("Version")
	return *s.Config().SqlSettings.DriverName, mattermostVersion.Value
}

func runReportToAWSMeterJob(s *Server) {
	model_helper.CreateRecurringTask("Collect and send usage report to AWS Metering Service", func() {
		doReportUsageToAWSMeteringService(s)
	}, time.Hour*model_helper.AwsMeteringReportInterval)
}

func doReportUsageToAWSMeteringService(s *Server) {
	awsMeter := awsmeter.New(s.Store, s.Config())
	if awsMeter == nil {
		slog.Error("Cannot obtain instance of AWS Metering Service.")
		return
	}

	dimensions := []string{model_helper.AwsMeteringDimensionUsageHrs}
	reports := awsMeter.GetUserCategoryUsage(dimensions, time.Now().UTC(), time.Now().Add(-model_helper.AwsMeteringReportInterval*time.Hour).UTC())
	awsMeter.ReportUserCategoryUsage(reports)
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
// 	model_helper.CreateRecurringTask("Command Hook Cleanup", func() {
// 		doCommandWebhookCleanup(s)
// 	}, time.Hour*1)
// }

// func doCommandWebhookCleanup(s *Server) {
// 	s.Store.CommandWebhook().Cleanup()
// }

func runTokenCleanupJob(s *Server) {
	doTokenCleanup(s)
	model_helper.CreateRecurringTask("Token Cleanup", func() {
		doTokenCleanup(s)
	}, time.Hour*1)
}

func doTokenCleanup(s *Server) {
	s.Store.Token().Cleanup()
}

func runSecurityJob(s *Server) {
	doSecurity(s)
	model_helper.CreateRecurringTask("Security", func() {
		doSecurity(s)
	}, time.Hour*4)
}

func runSessionCleanupJob(s *Server) {
	doSessionCleanup(s)
	model_helper.CreateRecurringTask("Session Cleanup", func() {
		doSessionCleanup(s)
	}, time.Hour*24)
}

func doSessionCleanup(s *Server) {
	s.Store.Session().Cleanup(model_helper.GetMillis(), SessionsCleanupBatchSize)
}

func doSecurity(s *Server) {
	s.DoSecurityUpdateCheck()
}

// func (s *Server) TelemetryId() string {
// 	if s.telemetryService == nil {
// 		return ""
// 	}

// 	return s.telemetryService.TelemetryID
// }

func (s *Server) getFirstServerRunTimestamp() (int64, *model_helper.AppError) {
	systemData, err := s.Store.System().GetByName(model_helper.SystemFirstServerRunTimestampKey)
	if err != nil {
		return 0, model_helper.NewAppError("getFirstServerRunTimestamp", "app.system.get_by_name.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	value, err := strconv.ParseInt(systemData.Value, 10, 64)
	if err != nil {
		return 0, model_helper.NewAppError("getFirstServerRunTimestamp", "app.system_install_date.parse_int.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return value, nil
}

func (s *Server) checkPushNotificationServerUrl() {
	notificationServer := *s.Config().EmailSettings.PushNotificationServer
	if strings.HasPrefix(notificationServer, "http://") {
		slog.Warn("Your push notification server is configured with HTTP. For improved security, update to HTTPS in your configuration.")
	}
}

func (s *Server) enableLoggingMetrics() {
	if s.Metrics == nil {
		return
	}

	s.Log.SetMetricsCollector(s.Metrics.GetLoggerMetricsCollector(), slog.DefaultMetricsUpdateFreqMillis)

	// logging config needs to be reloaded when metrics collector is added or changed.
	if err := s.initLogging(); err != nil {
		slog.Error("Error re-configuring logging for metrics")
		return
	}

	slog.Debug("Logging metrics enabled")
}

func (s *Server) StopHTTPServer() {
	if s.Server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), TimeToWaitForConnectionsToCloseOnServerShutdown)
		defer cancel()

		didShutDown := false
		for s.didFinishListen != nil && !didShutDown {
			if err := s.Server.Shutdown(ctx); err != nil {
				slog.Warn("Unable to shutdown server", slog.Err(err))
			}
			timer := time.NewTimer(time.Millisecond * 50)
			select {
			case <-s.didFinishListen:
				didShutDown = true
			case <-timer.C:
			}
			timer.Stop()
		}
		s.Server.Close()
		s.Server = nil
	}
}

// Shutdown  turn off system's server
func (s *Server) Shutdown() {
	slog.Info("Stopping Server...")

	defer sentry.Flush(2 * time.Second)

	// s.HubStop()
	// s.ShutDownPlugins()
	// s.RemoveLicenseListener(s.licenseListenerId)
	// s.RemoveLicenseListener(s.loggerLicenseListenerId)
	// s.RemoveClusterLeaderChangedListener(s.clusterLeaderListenerId)

	if s.tracer != nil {
		if err := s.tracer.Close(); err != nil {
			slog.Warn("Unable to cleanly shutdown opentracing client", slog.Err(err))
		}
	}

	// err := s.telemetryService.Shutdown()
	// if err != nil {
	// 	slog.Warn("Unable to cleanly shutdown telemetry client", slog.Err(err))
	// }

	// if s.remoteClusterService != nil {
	// 	if err = s.remoteClusterService.Shutdown(); err != nil {
	// 		slog.Error("Error shutting down intercluster services", slog.Err(err))
	// 	}
	// }

	s.StopHTTPServer()
	// s.stopLocalModeServer()

	// Push notification hub needs to be shutdown after HTTP server
	// to prevent stray requests from generating a push notification after it's shut down.

	// s.StopPushNotificationsHubWorkers()
	s.htmlTemplateWatcher.Close()

	s.WaitForGoroutines()

	s.RemoveConfigListener(s.configListenerId)
	s.stopSearchEngine()

	s.Audit.Shutdown()

	// s.stopFeatureFlagUpdateJob()
	s.ConfigStore.Close()

	if s.Cluster != nil {
		s.Cluster.StopInterNodeCommunication()
	}

	s.StopMetricsServer()

	var err error
	if s.Jobs != nil {
		// For simplicity we don't check if workers and schedulers are active
		// before stopping them as both calls essentially become no-ops
		// if nothing is running.
		if err = s.Jobs.StopWorkers(); err != nil && !errors.Is(err, jobs.ErrWorkersNotRunning) {
			slog.Warn("Failed to stop job server workers", slog.Err(err))
		}
		if err = s.Jobs.StopSchedulers(); err != nil && !errors.Is(err, jobs.ErrSchedulersNotRunning) {
			slog.Warn("Failed to stop job server schedulers", slog.Err(err))
		}
	}

	if s.Store != nil {
		s.Store.Close()
	}

	if s.CacheProvider != nil {
		if err = s.CacheProvider.Close(); err != nil {
			slog.Warn("Unable to cleanly shutdown cache", slog.Err(err))
		}
	}

	slog.Info("Server stopped")

	// shutdown main and notification loggers which will flush any remaining log records.
	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), time.Second*15)
	defer timeoutCancel()
	if err = s.NotificationsLog.ShutdownWithTimeout(timeoutCtx); err != nil {
		fmt.Fprintf(os.Stderr, "Error shutting down notification logger: %v", err)
	}
	if err = s.Log.ShutdownWithTimeout(timeoutCtx); err != nil {
		fmt.Fprintf(os.Stderr, "Error shutting down main logger: %v", err)
	}
}

func (s *Server) Restart() error {
	// percentage, err := s.UpgradeToE0Status()
	// if err != nil || percentage != 100 {
	// 	return errors.Wrap(err, "unable to restart because the system has not been upgraded")
	// }
	s.Shutdown()

	argv0, err := exec.LookPath(os.Args[0])
	if err != nil {
		return err
	}

	if _, err = os.Stat(argv0); err != nil {
		return err
	}

	slog.Info("Restarting server")
	return syscall.Exec(argv0, os.Args, os.Environ())
}

var corsAllowedMethods = []string{
	http.MethodPost,
	http.MethodGet,
	http.MethodOptions,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
}

// golang.org/x/crypto/acme/autocert/autocert.go
func handleHTTPRedirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "HEAD" {
		http.Error(w, "Use HTTPS", http.StatusBadRequest)
		return
	}
	target := "https://" + stripPort(r.Host) + r.URL.RequestURI()
	http.Redirect(w, r, target, http.StatusFound)
}

// golang.org/x/crypto/acme/autocert/autocert.go
func stripPort(hostport string) string {
	host, _, err := net.SplitHostPort(hostport)
	if err != nil {
		return hostport
	}
	return net.JoinHostPort(host, "443")
}

func (s *Server) Start() error {
	slog.Info("Starting Server...")

	var handler http.Handler = s.RootRouter

	if *s.Config().LogSettings.EnableDiagnostics && *s.Config().LogSettings.EnableSentry && !strings.Contains(SentryDSN, "placeholder") {
		sentryHandler := sentryhttp.New(sentryhttp.Options{
			Repanic: true,
		})
		handler = sentryHandler.Handle(handler)
	}

	if allowedOrigins := *s.Config().ServiceSettings.AllowCorsFrom; allowedOrigins != "" {
		exposedCorsHeaders := *s.Config().ServiceSettings.CorsExposedHeaders
		allowCredentials := *s.Config().ServiceSettings.CorsAllowCredentials
		debug := *s.Config().ServiceSettings.CorsDebug
		corsWrapper := cors.New(cors.Options{
			AllowedOrigins:   strings.Fields(allowedOrigins),
			AllowedMethods:   corsAllowedMethods,
			AllowedHeaders:   []string{"*"},
			ExposedHeaders:   strings.Fields(exposedCorsHeaders),
			MaxAge:           86400,
			AllowCredentials: allowCredentials,
			Debug:            debug,
		})

		// If we have debugging of CORS turned on then forward messages to logs
		if debug {
			corsWrapper.Log = s.Log.With(slog.String("source", "cors")).StdLogger(slog.LvlDebug)
		}

		handler = corsWrapper.Handler(handler)
	}

	if *s.Config().RateLimitSettings.Enable {
		slog.Info("RateLimiter is enabled")

		rateLimiter, err := NewRateLimiter(&s.Config().RateLimitSettings, s.Config().ServiceSettings.TrustedProxyIPHeader)
		if err != nil {
			return err
		}

		s.RateLimiter = rateLimiter
		handler = rateLimiter.RateLimitHandler(handler)
	}
	s.Busy = NewBusy(s.Cluster)

	// Creating a logger for logging errors from http.Server at error level

	s.Server = &http.Server{
		Handler:      handler,
		ReadTimeout:  time.Duration(*s.Config().ServiceSettings.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(*s.Config().ServiceSettings.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(*s.Config().ServiceSettings.IdleTimeout) * time.Second,
		ErrorLog:     s.Log.With(slog.String("source", "httpserver")).StdLogger(slog.LvlError),
	}

	addr := *s.Config().ServiceSettings.ListenAddress
	if addr == "" {
		if *s.Config().ServiceSettings.ConnectionSecurity == model_helper.CONN_SECURITY_TLS {
			addr = ":https"
		} else {
			addr = ":http"
		}
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrapf(err, i18n.T("api.server.start_server.starting.critical"), err)
	}
	s.ListenAddr = listener.Addr().(*net.TCPAddr)

	logListeningPort := fmt.Sprintf("Server is listening on %v", listener.Addr().String())
	slog.Info(logListeningPort, slog.String("address", listener.Addr().String()))

	m := &autocert.Manager{
		Cache:  autocert.DirCache(*s.Config().ServiceSettings.LetsEncryptCertificateCacheFile),
		Prompt: autocert.AcceptTOS,
	}

	if *s.Config().ServiceSettings.Forward80To443 {
		if host, port, err := net.SplitHostPort(addr); err != nil {
			slog.Error("Unable to setup forwarding", slog.Err(err))
		} else if port != "443" {
			return fmt.Errorf(i18n.T("api.server.start_server.forward80to443.enabled_but_listening_on_wrong_port"), port)
		} else {
			httpListenAddress := net.JoinHostPort(host, "http")

			if *s.Config().ServiceSettings.UseLetsEncrypt {
				server := &http.Server{
					Addr:     httpListenAddress,
					Handler:  m.HTTPHandler(nil),
					ErrorLog: s.Log.With(slog.String("source", "le_forwarder_server")).StdLogger(slog.LvlError),
				}
				go server.ListenAndServe()
			} else {
				go func() {
					redirectListener, err := net.Listen("tcp", httpListenAddress)
					if err != nil {
						slog.Error("Unable to setup forwarding", slog.Err(err))
						return
					}
					defer redirectListener.Close()

					server := &http.Server{
						Handler:  http.HandlerFunc(handleHTTPRedirect),
						ErrorLog: s.Log.With(slog.String("source", "forwarder_server")).StdLogger(slog.LvlError),
					}
					server.Serve(redirectListener)
				}()
			}
		}
	} else if *s.Config().ServiceSettings.UseLetsEncrypt {
		return errors.New(i18n.T("api.server.start_server.forward80to443.disabled_while_using_lets_encrypt"))
	}

	s.didFinishListen = make(chan struct{})
	go func() {
		var err error
		if *s.Config().ServiceSettings.ConnectionSecurity == model_helper.CONN_SECURITY_TLS {

			tlsConfig := &tls.Config{
				PreferServerCipherSuites: true,
				CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			}

			switch *s.Config().ServiceSettings.TLSMinVer {
			case "1.0":
				tlsConfig.MinVersion = tls.VersionTLS10
			case "1.1":
				tlsConfig.MinVersion = tls.VersionTLS11
			default:
				tlsConfig.MinVersion = tls.VersionTLS12
			}

			defaultCiphers := []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			}

			if len(s.Config().ServiceSettings.TLSOverwriteCiphers) == 0 {
				tlsConfig.CipherSuites = defaultCiphers
			} else {
				var cipherSuites []uint16
				for _, cipher := range s.Config().ServiceSettings.TLSOverwriteCiphers {
					value, ok := model_helper.ServerTLSSupportedCiphers[cipher]

					if !ok {
						slog.Warn("Unsupported cipher passed", slog.String("cipher", cipher))
						continue
					}

					cipherSuites = append(cipherSuites, value)
				}

				if len(cipherSuites) == 0 {
					slog.Warn("No supported ciphers passed, fallback to default cipher suite")
					cipherSuites = defaultCiphers
				}

				tlsConfig.CipherSuites = cipherSuites
			}

			certFile := ""
			keyFile := ""

			if *s.Config().ServiceSettings.UseLetsEncrypt {
				tlsConfig.GetCertificate = m.GetCertificate
				tlsConfig.NextProtos = append(tlsConfig.NextProtos, "h2")
			} else {
				certFile = *s.Config().ServiceSettings.TLSCertFile
				keyFile = *s.Config().ServiceSettings.TLSKeyFile
			}

			s.Server.TLSConfig = tlsConfig
			err = s.Server.ServeTLS(listener, certFile, keyFile)
		} else {
			err = s.Server.Serve(listener)
		}

		if err != nil && err != http.ErrServerClosed {
			slog.Critical("Error starting server", slog.Err(err))
			time.Sleep(time.Second)
		}

		close(s.didFinishListen)
	}()

	return nil
}

func (a *App) OriginChecker() func(*http.Request) bool {
	if allowed := *a.Config().ServiceSettings.AllowCorsFrom; allowed != "" {
		if allowed != "*" {
			siteURL, err := url.Parse(*a.Config().ServiceSettings.SiteURL)
			if err == nil {
				siteURL.Path = ""
				allowed += " " + siteURL.String()
			}
		}
		return api.OriginChecker(allowed)
	}

	return nil
}

// WaitForGoroutines blocks until all goroutines created by App.Go exit.
func (s *Server) WaitForGoroutines() {
	for atomic.LoadInt32(&s.goroutineCount) != 0 {
		<-s.goroutineExitSignal
	}
}

func (s *Server) stopSearchEngine() {
	s.RemoveConfigListener(s.searchConfigListenerId)
	// s.RemoveLicenseListener(s.searchLicenseListenerId)
	if s.SearchEngine != nil && s.SearchEngine.ElasticsearchEngine != nil && s.SearchEngine.ElasticsearchEngine.IsActive() {
		s.SearchEngine.ElasticsearchEngine.Stop()
	}
	if s.SearchEngine != nil && s.SearchEngine.BleveEngine != nil && s.SearchEngine.BleveEngine.IsActive() {
		s.SearchEngine.BleveEngine.Stop()
	}
}

func (s *Server) SetupMetricsServer() {
	if !*s.Config().MetricsSettings.Enable {
		return
	}

	s.StopMetricsServer()

	if err := s.InitMetricsRouter(); err != nil {
		slog.Error("Error initiating metrics router", slog.Err(err))
	}

	if s.Metrics != nil {
		s.Metrics.Register()
	}

	s.startMetricsServer()
}

func (s *Server) startMetricsServer() {
	var notify chan struct{}
	s.metricsLock.Lock()
	defer func() {
		if notify != nil {
			<-notify
		}
		s.metricsLock.Unlock()
	}()

	l, err := net.Listen("tcp", *s.Config().MetricsSettings.ListenAddress)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	notify = make(chan struct{})
	s.metricsServer = &http.Server{
		Handler:      handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(s.metricsRouter),
		ReadTimeout:  time.Duration(*s.Config().ServiceSettings.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(*s.Config().ServiceSettings.WriteTimeout) * time.Second,
	}

	go func() {
		close(notify)
		if err := s.metricsServer.Serve(l); err != nil && err != http.ErrServerClosed {
			slog.Critical(err.Error())
		}
	}()

	s.Log.Info("Metrics and profiling server is started", slog.String("address", l.Addr().String()))
}

func (s *Server) StopMetricsServer() {
	s.metricsLock.Lock()
	defer s.metricsLock.Unlock()

	if s.metricsServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), TimeToWaitForConnectionsToCloseOnServerShutdown)
		defer cancel()

		s.metricsServer.Shutdown(ctx)
		s.Log.Info("Metrics and profiling server is stopping")
	}
}

func (s *Server) InitMetricsRouter() error {
	s.metricsRouter = mux.NewRouter()
	runtime.SetBlockProfileRate(*s.Config().MetricsSettings.BlockProfileRate)

	metricsPage := `
			<html>
				<body>{{if .}}
					<div><a href="/metrics">Metrics</a></div>{{end}}
					<div><a href="/debug/pprof/">Profiling Root</a></div>
					<div><a href="/debug/pprof/cmdline">Profiling Command Line</a></div>
					<div><a href="/debug/pprof/symbol">Profiling Symbols</a></div>
					<div><a href="/debug/pprof/goroutine">Profiling Goroutines</a></div>
					<div><a href="/debug/pprof/heap">Profiling Heap</a></div>
					<div><a href="/debug/pprof/threadcreate">Profiling Threads</a></div>
					<div><a href="/debug/pprof/block">Profiling Blocking</a></div>
					<div><a href="/debug/pprof/trace">Profiling Execution Trace</a></div>
					<div><a href="/debug/pprof/profile">Profiling CPU</a></div>
				</body>
			</html>
		`
	metricsPageTmpl, err := template.New("page").Parse(metricsPage)
	if err != nil {
		return errors.Wrap(err, "failed to create template")
	}

	rootHandler := func(w http.ResponseWriter, r *http.Request) {
		metricsPageTmpl.Execute(w, s.Metrics != nil)
	}

	s.metricsRouter.HandleFunc("/", rootHandler)
	s.metricsRouter.StrictSlash(true)

	s.metricsRouter.Handle("/debug", http.RedirectHandler("/", http.StatusMovedPermanently))
	s.metricsRouter.HandleFunc("/debug/pprof/", pprof.Index)
	s.metricsRouter.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	s.metricsRouter.HandleFunc("/debug/pprof/profile", pprof.Profile)
	s.metricsRouter.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	s.metricsRouter.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// Manually add support for paths linked to by index page at /debug/pprof/
	s.metricsRouter.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	s.metricsRouter.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	s.metricsRouter.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	s.metricsRouter.Handle("/debug/pprof/block", pprof.Handler("block"))

	return nil
}

func (s *Server) StartSearchEngine() string {
	if s.SearchEngine.ElasticsearchEngine != nil && s.SearchEngine.ElasticsearchEngine.IsActive() {
		s.Go(func() {
			if err := s.SearchEngine.ElasticsearchEngine.Start(); err != nil {
				s.Log.Error(err.Error())
			}
		})
	}

	configListenerId := s.AddConfigListener(func(oldConfig *model_helper.Config, newConfig *model_helper.Config) {
		if s.SearchEngine == nil {
			return
		}
		s.SearchEngine.UpdateConfig(newConfig)

		if s.SearchEngine.ElasticsearchEngine != nil && !*oldConfig.ElasticsearchSettings.EnableIndexing && *newConfig.ElasticsearchSettings.EnableIndexing {
			s.Go(func() {
				if err := s.SearchEngine.ElasticsearchEngine.Start(); err != nil {
					slog.Error(err.Error())
				}
			})
		} else if s.SearchEngine.ElasticsearchEngine != nil && *oldConfig.ElasticsearchSettings.EnableIndexing && !*newConfig.ElasticsearchSettings.EnableIndexing {
			s.Go(func() {
				if err := s.SearchEngine.ElasticsearchEngine.Stop(); err != nil {
					slog.Error(err.Error())
				}
			})
		} else if s.SearchEngine.ElasticsearchEngine != nil && *oldConfig.ElasticsearchSettings.Password != *newConfig.ElasticsearchSettings.Password || *oldConfig.ElasticsearchSettings.Username != *newConfig.ElasticsearchSettings.Username || *oldConfig.ElasticsearchSettings.ConnectionUrl != *newConfig.ElasticsearchSettings.ConnectionUrl || *oldConfig.ElasticsearchSettings.Sniff != *newConfig.ElasticsearchSettings.Sniff {
			s.Go(func() {
				if *oldConfig.ElasticsearchSettings.EnableIndexing {
					if err := s.SearchEngine.ElasticsearchEngine.Stop(); err != nil {
						slog.Error(err.Error())
					}
					if err := s.SearchEngine.ElasticsearchEngine.Start(); err != nil {
						slog.Error(err.Error())
					}
				}
			})
		}
	})

	return configListenerId
}

func (s *Server) GetSiteURL() string {
	return *s.Config().ServiceSettings.SiteURL
}

// GetCookieDomain
func (s *Server) GetCookieDomain() string {
	if *s.Config().ServiceSettings.AllowCookiesForSubdomains {
		if siteURL, err := url.Parse(*s.Config().ServiceSettings.SiteURL); err == nil {
			return siteURL.Hostname()
		}
	}
	return ""
}
