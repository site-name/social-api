package commands

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/request"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/config"
	"github.com/sitename/sitename/modules/i18n"
	"github.com/sitename/sitename/modules/util"
	"github.com/spf13/cobra"
)

func initDBCommandContextCobra(command *cobra.Command, readOnlyConfigStore bool) (*app.App, error) {
	a, err := initDBCommandContext(getConfigDSN(command, config.GetEnvironment()), readOnlyConfigStore)
	if err != nil {
		// Returning an error just prints the usage message, so actually panic
		panic(err)
	}

	a.Srv().PluginService().InitPlugins(&request.Context{}, *a.Config().PluginSettings.Directory, *a.Config().PluginSettings.ClientDirectory)
	a.DoAppMigrations()

	return a, nil
}

// InitDBCommandContextCobra
func InitDBCommandContextCobra(command *cobra.Command) (*app.App, error) {
	return initDBCommandContextCobra(command, true)
}

func InitDBCommandContextCobraReadWrite(command *cobra.Command) (*app.App, error) {
	return initDBCommandContextCobra(command, false)
}

func initDBCommandContext(configDSN string, readOnlyConfigStore bool) (*app.App, error) {
	if err := util.TranslationsPreInit(); err != nil {
		return nil, err
	}
	model_helper.AppErrorInit(i18n.T)

	s, err := app.NewServer(
		app.Config(configDSN, readOnlyConfigStore, nil),
		app.StartSearchEngine,
	)
	if err != nil {
		return nil, err
	}

	a := app.New(app.ServerConnector(s))

	return a, nil
}
