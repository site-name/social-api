package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/sitename/sitename/modules/setting"
	"github.com/urfave/cli"
)

var (
	// Sitename version
	Version = "development"
	// Build tags used
	Tags = ""
	// Make version used
	MakeVersion = ""

	originalAppHelpTemplate        = ""
	originalCommandHelpTemplate    = ""
	originalSubcommandHelpTemplate = ""
)

func init() {
	setting.AppVer = Version
	setting.AppBuiltWith = formatBuiltWith()
	setting.AppStartTime = time.Now().UTC()

	originalAppHelpTemplate = cli.AppHelpTemplate
	originalCommandHelpTemplate = cli.CommandHelpTemplate
	originalSubcommandHelpTemplate = cli.SubcommandHelpTemplate
}

func main() {
	app := cli.NewApp()
	app.Name = "Sitename"
	app.Usage = "selling platform"
	app.Description = "Selling platform"
	app.Version = Version + formatBuiltWith()
	app.Commands = []cli.Command{}

	setting.SetCustomPathAndConf("", "", "")
	setAppHelpTemplates()

	defaultFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "custom-path, C",
			Value: setting.CustomPath,
			Usage: "Custom path file path",
		},
		cli.StringFlag{
			Name:  "config, c",
			Value: setting.CustomConf,
			Usage: "Custom configuration file path",
		},
		cli.VersionFlag,
		cli.StringFlag{
			Name:  "work-path, w",
			Value: setting.AppWorkPath,
			Usage: "Set the sitename working path",
		},
	}

	// Set the default to be equivalent to cmdWeb and add the default flags
	app.Flags = append(app.Flags, defaultFlags...)

	// Add functions to set these paths and these flags to the commands
	app.Before = establishCustomPath
}

func establishCustomPath(ctx *cli.Context) error {
	var (
		providedCustom   string
		providedConf     string
		providedWorkPath string
	)

	currentCtx := ctx
	for {
		if len(providedCustom) != 0 && len(providedConf) != 0 && len(providedWorkPath) != 0 {
			break
		}
		if currentCtx == nil {
			break
		}
		if currentCtx.IsSet("custom-path") && len(providedCustom) == 0 {
			providedCustom = currentCtx.String("cusom-path")
		}
		if currentCtx.IsSet("config") && len(providedConf) == 0 {
			providedConf = currentCtx.String("config")
		}
		if currentCtx.IsSet("work-path") && len(providedWorkPath) == 0 {
			providedWorkPath = currentCtx.String("work-path")
		}
		currentCtx = currentCtx.Parent()
	}
	setting.SetCustomPathAndConf(providedCustom, providedConf, providedWorkPath)
	setAppHelpTemplates()

	if ctx.IsSet("version") {
		cli.ShowVersion(ctx)
		os.Exit(0)
	}

	return nil
}

func setAppHelpTemplates() {
	cli.AppHelpTemplate = adjustHelpTemplate(originalAppHelpTemplate)
	cli.CommandHelpTemplate = adjustHelpTemplate(originalCommandHelpTemplate)
	cli.SubcommandHelpTemplate = adjustHelpTemplate(originalSubcommandHelpTemplate)
}

func adjustHelpTemplate(originalTemplate string) string {
	overrided := ""
	if _, ok := os.LookupEnv("SITENAME_CUSTOM"); ok {
		overrided = "(SITENAME_CUSTOM)"
	}

	return fmt.Sprintf(`%s
DEFAULT CONFIGURATION:
     CustomPath:  %s %s
     CustomConf:  %s
     AppPath:     %s
     AppWorkPath: %s

`, originalTemplate, setting.CustomPath, overrided, setting.CustomConf, setting.AppPath, setting.AppWorkPath)
}

func formatBuiltWith() string {
	version := runtime.Version()
	if len(version) > 0 {
		version = MakeVersion + ", " + version
	}
	if len(Tags) == 0 {
		return " built with " + version
	}

	return " built with " + version + " : " + strings.ReplaceAll(Tags, " ", ", ")
}
