package main

import (
	"os"

	// Enterprise Deps
	_ "github.com/gorilla/handlers"
	_ "github.com/hako/durafmt"
	_ "github.com/hashicorp/memberlist"
	_ "github.com/mattermost/gosaml2"
	_ "github.com/mattermost/ldap"
	_ "github.com/mattermost/rsc/qr"
	_ "github.com/prometheus/client_golang/prometheus"
	_ "github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sitename/sitename/cmd/sitename/commands"
	_ "github.com/sitename/sitename/model_helper"
	_ "github.com/sitename/sitename/modules/imports"
	_ "github.com/tylerb/graceful"
	_ "gopkg.in/olivere/elastic.v6"
)

func main() {
	if err := commands.Run(os.Args[1:]); err != nil {
		os.Exit(1)
	}
}
