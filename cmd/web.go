package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	context2 "github.com/gorilla/context"
	"github.com/sitename/sitename/modules/graceful"
	"github.com/sitename/sitename/modules/log"
	"github.com/sitename/sitename/modules/setting"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/routers"
	"github.com/sitename/sitename/routers/routes"
	"github.com/urfave/cli"
	"gopkg.in/ini.v1"
)

var CmdWeb = cli.Command{
	Name:        "web",
	Usage:       "Start Sitename web server",
	Description: "Sitename web server is the only thing you need to run",
	Action:      runWeb,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "port, p",
			Value: "3000",
			Usage: "Temporary port number to prevent conflict",
		},
		cli.StringFlag{
			Name:  "pid, P",
			Value: setting.PIDFile,
			Usage: "Custom pid file path",
		},
	},
}

func runHTTPRedirector() {
	source := fmt.Sprintf("%s:%s", setting.HTTPAddr, setting.PortToRedirect)
	dest := strings.TrimSuffix(setting.AppURL, "/")
	log.Info("Redirecting: %s to %s", source, dest)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		target := dest + r.URL.Path
		if len(r.URL.RawQuery) > 0 {
			target += "?" + r.URL.RawQuery
		}
		http.Redirect(w, r, target, http.StatusTemporaryRedirect)
	})

	var err = runHTTP("tcp", source, "HTTP Redirector", context2.ClearHandler(handler))

	if err != nil {
		log.Fatal("Failed to start port redirection: %v", err)
	}
}

func runWeb(ctx *cli.Context) error {
	managerCtx, cancel := context.WithCancel(context.Background())
	graceful.InitManager(managerCtx)
	defer cancel()

	if os.Getpid() > 1 && len(os.Getenv("LISTEN_FDS")) > 0 {
		log.Info("Restarting Sitename on PID: %d from parent PID: %d", os.Getpid(), os.Getppid())
	} else {
		log.Info("Starting Sitename on PID: %d", os.Getpid())
	}

	// Set pid file setting
	if ctx.IsSet("pid") {
		setting.PIDFile = ctx.String("pid")
		setting.WritePIDFile = true
	}

	if setting.EnablePprof {
		go func() {
			log.Info("Starting pprof server on localhost:6060")
			log.Info("%v", http.ListenAndServe("localhost:6060", nil))
		}()
	}

	log.Info("Global init")
	// Perform global initialization
	routers.GlobalInit(graceful.GetManager().HammerContext())

	// Override the provided port number within the configuration
	if ctx.IsSet("port") {
		if err := setPort(ctx.String("port")); err != nil {
			return err
		}
	}

	NoInstallListener()
	if setting.EnablePprof {
		go func() {
			log.Info("Starting pport server on localhost:6060")
			log.Info("%v", http.ListenAndServe("localhost:6060", nil))
		}()
	}

	log.Info("Global init")
	// perform global initialization
	routers.GlobalInit(graceful.GetManager().HammerContext())

	// Override the provided port number within the configuration
	if ctx.IsSet("port") {
		if err := setPort(ctx.String("port")); err != nil {
			return err
		}
	}

	// Setup chi routes
	c := routes.NormalRoutes()
	err := listen(c, true)
	<-graceful.GetManager().Done()
	log.Info("PID: %d Sitename Web Finished", os.Getpid())
	log.Close()
	return err
}

func setPort(port string) error {
	setting.AppURL = strings.Replace(setting.AppURL, setting.HTTPPort, port, 1)
	setting.HTTPPort = port

	switch setting.Protocol {
	case setting.UnixSocket:
	case setting.FCGI:
	case setting.FCGIUnix:
	default:
		// save LOCAL_ROOT_URL if port changed
		cfg := ini.Empty()
		isFile, err := util.IsFile(setting.CustomConf)
		if err != nil {
			log.Fatal("Unable to check if %s is a file", err)
		}
		if isFile {
			// Keeps custom settings if there is already something.
			if err := cfg.Append(setting.CustomConf); err != nil {
				return fmt.Errorf("Failed to load custom conf '%s': %v", setting.CustomConf, err)
			}
		}

		defaultLocalURL := string(setting.Protocol) + "://"
		if setting.HTTPAddr == "0.0.0.0" {
			defaultLocalURL += "localhost"
		} else {
			defaultLocalURL += setting.HTTPAddr
		}
		defaultLocalURL += ":" + setting.HTTPPort + "/"

		cfg.Section("server").Key("LOCAL_ROOT_URL").SetValue(defaultLocalURL)
		if err := cfg.SaveTo(setting.CustomConf); err != nil {
			return fmt.Errorf("Error saving generated JWT Secret to custom config: %v", err)
		}
	}
	return nil
}

func listen(m http.Handler, handleRedirector bool) error {
	listenAddr := setting.HTTPAddr
	if setting.Protocol != setting.UnixSocket && setting.Protocol != setting.FCGIUnix {
		listenAddr = net.JoinHostPort(listenAddr, setting.HTTPPort)
	}
	log.Info("Listen: %v://%s%s", setting.Protocol, listenAddr, setting.AppSubURL)

	if setting.LFS.StartServer {
		log.Info("LFS server enabled")
	}

	var err error
	switch setting.Protocol {
	case setting.HTTP:
		if handleRedirector {
			NoHTTPRedirector()
		}
		err = runHTTP("tcp", listenAddr, "Web", context2.ClearHandler(m))
	case setting.HTTPS:
		if setting.EnableLetsEncrypt {
			err = runLetsEncrypt(listenAddr, setting.Domain, setting.LetsEncryptDirectory, setting.LetsEncryptEmail, context2.ClearHandler(m))
			break
		}
		if handleRedirector {
			if setting.RedirectOtherPort {
				go runHTTPRedirector()
			} else {
				NoHTTPRedirector()
			}
		}
		err = runHTTPS("tcp", listenAddr, "Web", setting.CertFile, setting.KeyFile, context2.ClearHandler(m))
	// case setting.FCGI:
	// 	if handleRedirector {
	// 		NoHTTPRedirector()
	// 	}
	// 	err = runFCGI("tcp", listenAddr, "FCGI Web", context2.ClearHandler(m))
	// case setting.UnixSocket:
	// 	if handleRedirector {
	// 		NoHTTPRedirector()
	// 	}
	// 	err = runHTTP("unix", listenAddr, "Web", context2.ClearHandler(m))
	// case setting.FCGIUnix:
	// 	if handleRedirector {
	// 		NoHTTPRedirector()
	// 	}
	// 	err = runFCGI("unix", listenAddr, "Web", context2.ClearHandler(m))
	default:
		log.Fatal("Invalid protocol: %s", setting.Protocol)
	}

	if err != nil {
		log.Critical("Failed to start server: %v", err)
	}
	log.Info("HTTP Listener: %s Closed", listenAddr)
	return err
}
