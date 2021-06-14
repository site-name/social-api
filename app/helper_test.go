package app

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/sitename/sitename/app/request"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/modules/config"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/localcachelayer"
)

type TestHelper struct {
	App        *App
	Context    *request.Context
	Server     *Server
	BasicUser  *account.User
	BasicUser2 *account.User

	SystemAdminUser   *account.User
	LogBuffer         *bytes.Buffer
	IncludeCacheLayer bool
	tempWorkspace     string
}

func setupTestHelper(dbStore store.Store, enterprise bool, includeCacheLayer bool, tb testing.TB) *TestHelper {
	tempWorkspace, err := ioutil.TempDir("", "apptest")
	if err != nil {
		panic(err)
	}

	configStore := config.NewTestMemoryStore()

	config := configStore.Get()
	*config.PluginSettings.Directory = filepath.Join(tempWorkspace, "plugins")
	*config.PluginSettings.ClientDirectory = filepath.Join(tempWorkspace, "webapp")
	*config.PluginSettings.AutomaticPrepackagedPlugins = false
	*config.LogSettings.EnableSentry = false // disable error reporting during tests
	*config.AnnouncementSettings.AdminNoticesEnabled = false
	*config.AnnouncementSettings.UserNoticesEnabled = false
	configStore.Set(config)

	buffer := &bytes.Buffer{}

	var options []Option
	options = append(options, ConfigStore(configStore))
	if includeCacheLayer {
		// Adds the cache layer to the test store
		options = append(options, StoreOverride(func(s *Server) store.Store {
			lcl, err2 := localcachelayer.NewLocalCacheLayer(dbStore, s.Metrics, s.Cluster, s.CacheProvider)
			if err2 != nil {
				panic(err2)
			}
			return lcl
		}))
	} else {
		options = append(options, StoreOverride(dbStore))
	}
	options = append(options, SetLogger(slog.NewTestingLogger(tb, buffer)))

	s, err := NewServer(options...)
	if err != nil {
		panic(err)
	}

	th := &TestHelper{
		App:               New(ServerConnector(s)),
		Context:           &request.Context{},
		Server:            s,
		LogBuffer:         buffer,
		IncludeCacheLayer: includeCacheLayer,
	}

	// th.App.UpdateConfig(func(cfg *model.Config) { *cfg.TeamSettings.MaxUsersPerTeam = 50 })
	th.App.UpdateConfig(func(cfg *model.Config) { *cfg.RateLimitSettings.Enable = false })
	prevListenAddress := *th.App.Config().ServiceSettings.ListenAddress
	th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ServiceSettings.ListenAddress = ":0" })
	serverErr := th.Server.Start()
	if serverErr != nil {
		panic(serverErr)
	}

	th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ServiceSettings.ListenAddress = prevListenAddress })

	th.App.Srv().SearchEngine = mainHelper.SearchEngine

	th.App.Srv().Store.MarkSystemRanUnitTests()

	// th.App.UpdateConfig(func(cfg *model.Config) { *cfg.TeamSettings.EnableOpenServer = true })

	// Disable strict password requirements for test
	th.App.UpdateConfig(func(cfg *model.Config) {
		*cfg.PasswordSettings.MinimumLength = 5
		*cfg.PasswordSettings.Lowercase = false
		*cfg.PasswordSettings.Uppercase = false
		*cfg.PasswordSettings.Symbol = false
		*cfg.PasswordSettings.Number = false
	})

	if enterprise {
		// th.App.Srv().SetLicense(model.NewTestLicense())
		th.App.Srv().Jobs.InitWorkers()
		th.App.Srv().Jobs.InitSchedulers()
	} else {
		// th.App.Srv().SetLicense(nil)
	}

	if th.tempWorkspace == "" {
		th.tempWorkspace = tempWorkspace
	}

	return th
}

func Setup(tb testing.TB) *TestHelper {
	if testing.Short() {
		tb.SkipNow()
	}
	dbStore := mainHelper.GetStore()
	dbStore.DropAllTables()
	dbStore.MarkSystemRanUnitTests()
	mainHelper.PreloadMigrations()

	return setupTestHelper(dbStore, false, true, tb)
}

func SetupWithoutPreloadMigrations(tb testing.TB) *TestHelper {
	if testing.Short() {
		tb.SkipNow()
	}
	dbStore := mainHelper.GetStore()
	dbStore.DropAllTables()
	dbStore.MarkSystemRanUnitTests()

	return setupTestHelper(dbStore, false, true, tb)
}

// func SetupWithStoreMock(tb testing.TB) *TestHelper {
// 	mockStore := testlib.GetMockStoreForSetupFunctions()
// 	th := setupTestHelper(mockStore, false, false, tb)
// 	emptyMockStore := mocks.Store{}
// 	emptyMockStore.On("Close").Return(nil)
// 	th.App.Srv().Store = &emptyMockStore
// 	return th
// }

var initBasicOnce sync.Once
var userCache struct {
	SystemAdminUser *account.User
	BasicUser       *account.User
	BasicUser2      *account.User
}

func (th *TestHelper) InitBasic() *TestHelper {
	// create users once and cache them because password hashing is slow
	initBasicOnce.Do(func() {
		th.SystemAdminUser = th.CreateUser()
		th.App.UpdateUserRoles(th.SystemAdminUser.Id, model.SYSTEM_USER_ROLE_ID+" "+model.SYSTEM_ADMIN_ROLE_ID, false)
		th.SystemAdminUser, _ = th.App.GetUser(th.SystemAdminUser.Id)
		userCache.SystemAdminUser = th.SystemAdminUser.DeepCopy()

		th.BasicUser = th.CreateUser()
		th.BasicUser, _ = th.App.GetUser(th.BasicUser.Id)
		userCache.BasicUser = th.BasicUser.DeepCopy()

		th.BasicUser2 = th.CreateUser()
		th.BasicUser2, _ = th.App.GetUser(th.BasicUser2.Id)
		userCache.BasicUser2 = th.BasicUser2.DeepCopy()
	})
	// restore cached users
	th.SystemAdminUser = userCache.SystemAdminUser.DeepCopy()
	th.BasicUser = userCache.BasicUser.DeepCopy()
	th.BasicUser2 = userCache.BasicUser2.DeepCopy()
	mainHelper.GetSQLStore().GetMaster().Insert(th.SystemAdminUser, th.BasicUser, th.BasicUser2)

	// th.BasicTeam = th.CreateTeam()

	// th.LinkUserToTeam(th.BasicUser, th.BasicTeam)
	// th.LinkUserToTeam(th.BasicUser2, th.BasicTeam)
	// th.BasicChannel = th.CreateChannel(th.BasicTeam)
	// th.BasicPost = th.CreatePost(th.BasicChannel)
	return th
}

func (*TestHelper) MakeEmail() string {
	return "success_" + model.NewId() + "@simulator.amazonses.com"
}

func (th *TestHelper) CreateUser() *account.User {
	return th.CreateUserOrGuest(false)
}

func (th *TestHelper) CreateGuest() *account.User {
	return th.CreateUserOrGuest(true)
}

func (th *TestHelper) CreateUserOrGuest(guest bool) *account.User {
	id := model.NewId()

	user := &account.User{
		Email:         "success+" + id + "@simulator.amazonses.com",
		Username:      "un_" + id,
		Nickname:      "nn_" + id,
		Password:      "Password1",
		EmailVerified: true,
	}

	util.DisableDebugLogForTest()
	var err *model.AppError
	if guest {
		if user, err = th.App.CreateGuest(th.Context, user); err != nil {
			panic(err)
		}
	} else {
		if user, err = th.App.CreateUser(th.Context, user); err != nil {
			panic(err)
		}
	}
	util.EnableDebugLogForTest()
	return user
}

func (th *TestHelper) TearDown() {
	if th.IncludeCacheLayer {
		// Clean all the caches
		th.App.Srv().InvalidateAllCaches()
	}
	th.ShutdownApp()
	if th.tempWorkspace != "" {
		os.RemoveAll(th.tempWorkspace)
	}
}

func (th *TestHelper) ShutdownApp() {
	done := make(chan bool)
	go func() {
		th.Server.Shutdown()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(30 * time.Second):
		// panic instead of fatal to terminate all tests in this package, otherwise the
		// still running App could spuriously fail subsequent tests.
		panic("failed to shutdown App within 30 seconds")
	}
}
