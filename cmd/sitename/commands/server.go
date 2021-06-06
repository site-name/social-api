package commands

import (
	"bytes"
	"net"
	"os"
	"os/signal"
	"runtime/debug"
	"runtime/pprof"
	"syscall"

	"github.com/pkg/errors"
	// "github.com/sitename/sitename/api"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/modules/config"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/web"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:          "server",
	Short:        "Run the Mattermost server",
	RunE:         serverCmdF,
	SilenceUsage: true,
}

func init() {
	RootCmd.AddCommand(serverCmd)
	RootCmd.RunE = serverCmdF
}

func serverCmdF(command *cobra.Command, args []string) error {
	disableConfigWatch, _ := command.Flags().GetBool("disableconfigwatch")

	interruptChan := make(chan os.Signal, 1)

	if err := util.TranslationsPreInit(); err != nil {
		return errors.Wrap(err, "unable to load Sitename translation files")
	}

	customDefaults, err := loadCustomDefaults()
	if err != nil {
		slog.Warn("Error loading custom configuration defaults: " + err.Error())
	}

	configStore, err := config.NewStoreFromDSN(getConfigDSN(command, config.GetEnvironment()), !disableConfigWatch, false, customDefaults)
	if err != nil {
		return errors.Wrap(err, "failed to load configuration")
	}
	defer configStore.Close()

	return runServer(configStore, interruptChan)
}

func runServer(configStore *config.Store, interruptChan chan os.Signal) error {
	// Setting the highest traceback level from the code.
	// This is done to print goroutines from all threads (see golang.org/issue/13161)
	// and also preserve a crash dump for later investigation.
	debug.SetTraceback("crash")

	options := []app.Option{
		app.ConfigStore(configStore),
		app.RunEssentialJobs,
		app.JoinCluster,
		app.StartSearchEngine,
		app.StartMetrics,
	}
	server, err := app.NewServer(options...)
	if err != nil {
		slog.Critical(err.Error())
		return err
	}
	defer server.Shutdown()
	// We add this after shutdown so that it can be called
	// before server shutdown happens as it can close
	// the advanced logger and prevent the mlog call from working properly.
	defer func() {
		// A panic pass-through layer which just logs it
		// and sends it upwards.
		if x := recover(); x != nil {
			var buf bytes.Buffer
			pprof.Lookup("goroutine").WriteTo(&buf, 2)
			slog.Critical("A panic occurred",
				slog.Any("error", x),
				slog.String("stack", buf.String()))
			panic(x)
		}
	}()

	a := app.New(app.ServerConnector(server))
	// api.Init(a, server.RootRouter)
	web.New(a, server.RootRouter)

	serverErr := server.Start()
	if serverErr != nil {
		slog.Critical(serverErr.Error())
		return serverErr
	}

	// If we allow testing then listen for manual testing URL hits
	// if *server.Config().ServiceSettings.EnableTesting {
	// 	manualtesting.Init(api)
	// }

	notifyReady()

	// wait for kill signal before attempting to gracefully shutdown
	// the running service
	signal.Notify(interruptChan, syscall.SIGINT, syscall.SIGTERM)
	<-interruptChan
	return nil
}

func notifyReady() {
	// If the environment vars provide a systemd notification socket,
	// notify systemd that the server is ready.
	systemdSocket := os.Getenv("NOTIFY_SOCKET")
	if systemdSocket != "" {
		slog.Info("Sending systemd READY notification.")

		err := sendSystemdReadyNotification(systemdSocket)
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

func sendSystemdReadyNotification(socketPath string) error {
	msg := "READY=1"
	addr := &net.UnixAddr{
		Name: socketPath,
		Net:  "unixgram",
	}
	conn, err := net.DialUnix(addr.Net, nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write([]byte(msg))
	return err
}
