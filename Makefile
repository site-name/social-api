.PHONY: build package run stop run-client run-server run-haserver stop-haserver stop-client stop-server restart restart-server restart-client restart-haserver start-docker clean-dist clean nuke check-style check-client-style check-server-style check-unit-tests test dist prepare-enteprise run-client-tests setup-run-client-tests cleanup-run-client-tests test-client build-linux build-osx build-windows internal-test-web-client vet run-server-for-web-client-tests diff-config prepackaged-plugins prepackaged-binaries test-server test-server-ee test-server-quick test-server-race start-docker-check migrations-bindata new-migration migration-prereqs

ROOT := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

ifeq ($(OS), Windows_NT)
	PLATFORM := Windows
else
	PLATFORM := $(shell uname)
endif

# Set an environment variable on Linux used to resolve `docker.host.internal` inconsistencies with
# docker. This can be reworked once https://github.com/docker/for-linux/issues/264 is resolved
# satisfactorily.

ifeq ($(PLATFORM),linux)
	export IS_LINUX = -linux
else
	export IS_LINUX =
endif

IS_CT ?= false
# Build Flags
BUILD_NUMBER ?= $(BUILD_NUMBER:)
BUILD_DATE = $(shell date -u)
BUILD_HASH = $(shell git rev-parse HEAD)

# If we don't set the build number it defaults to dev
ifeq ($(BUILD_NUMBER),)
	BUILD_NUMBER := dev
endif
BUILD_ENTERPRISE_DIR ?= ../enterprise
BUILD_ENTERPRISE ?= true
BUILD_ENTERPRISE_READY = false
BUILD_TYPE_NAME = team
BUILD_HASH_ENTERPRISE = none

ifneq ($(wildcard $(BUILD_ENTERPRISE_DIR)/.),)
	ifeq ($(BUILD_ENTERPRISE),true)
		BUILD_ENTERPRISE_READY = true
		BUILD_TYPE_NAME = enterprise
		BUILD_HASH_ENTERPRISE = $(shell cd $(BUILD_ENTERPRISE_DIR) && git rev-parse HEAD)
	else
		BUILD_ENTERPRISE_READY = false
		BUILD_TYPE_NAME = team
	endif
else
	BUILD_ENTERPRISE_READY = false
	BUILD_TYPE_NAME = team
endif

BUILD_WEBAPP_DIR ?= ../sitename-webapp
BUILD_CLIENT = false
BUILD_HASH_CLIENT = independant
ifneq ($(wildcard $(BUILD_WEBAPP_DIR)/.),)
	ifeq ($(BUILD_CLIENT),true)
		BUILD_CLIENT = true
		BUILD_HASH_CLIENT = $(shell cd $(BUILD_WEBAPP_DIR) && git rev-parse HEAD)
	else
		BUILD_CLIENT = false
	endif
else
	BUILD_CLIENT = false
endif

# We need current user's UID for `run-haserver` so docker compose does not run server
# as root and mess up file permissions for devs. When running like this HOME will be blank
# and docker will add '/', so we need to set the go-build cache location or we'll get 
# permission errors on build as it tries to create a cache in filesystem root.
export CURRENT_UID = $(shell id -u):$(shell id -g)
ifeq ($(HOME),/)
	export XDG_CACHE_HOME = /tmp/go-cache/
endif

# Go Flags
GOFLAGS ?= $(GOFLAGS:)
# We need to export GOBIN to allow it to be set
# for processes spawned from the Makefile
export GOBIN ?= $(PWD)/bin
GO=go
DELVE=dlv

GO_MAJOR_VERSION = $(shell $(GO) version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f1)
GO_MINOR_VERSION = $(shell $(GO) version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)
MINIMUM_SUPPORTED_GO_MAJOR_VERSION = 1
MINIMUM_SUPPORTED_GO_MINOR_VERSION = 15
GO_VERSION_VALIDATION_ERR_MSG = Golang version is not supported, please update to at least $(MINIMUM_SUPPORTED_GO_MAJOR_VERSION).$(MINIMUM_SUPPORTED_GO_MINOR_VERSION)


# GOOS/GOARCH of the build host, used to determine whether we're cross-compiling or not
BUILDER_GOOS_GOARCH="$(shell $(GO) env GOOS)_$(shell $(GO) env GOARCH)"

PLATFORM_FILES="./cmd/sitename/main.go"

# Output paths
DIST_ROOT=dist
DIST_PATH=$(DIST_ROOT)/sitename

# Tests
TESTS=.

# Packages lists
TE_PACKAGES=$(shell $(GO) list ./... | grep -v ./data)

TEMPLATES_DIR=templates

# Prepares the enterprise build if exists. The IGNORE stuff is a hack to get the Makefile to execute the commands outside a target
ifeq ($(BUILD_ENTERPRISE_READY),true)
	IGNORE:=$(shell echo Enterprise build selected, preparing)
	IGNORE:=$(shell rm -f imports/imports.go)
	IGNORE:=$(shell cp $(BUILD_ENTERPRISE_DIR)/imports/imports.go imports/)
	IGNORE:=$(shell rm -f enterprise)
	IGNORE:=$(shell ln -s $(BUILD_ENTERPRISE_DIR) enterprise)
else
	IGNORE:=$(shell rm -f imports/imports.go)
endif

EE_PACKAGES=$(shell $(GO) list ./enterprise/...)

ifeq ($(BUILD_ENTERPRISE_READY),true)
ALL_PACKAGES=$(TE_PACKAGES) $(EE_PACKAGES)
else
ALL_PACKAGES=$(TE_PACKAGES)
endif

-include config.override.mk
include config.mk
# include build/*.mk

RUN_IN_BACKGROUND ?=
ifeq ($(RUN_SERVER_IN_BACKGROUND),true)
	RUN_IN_BACKGROUND := &
endif

app-layers: ## Extract interface from App struct
	$(GO) get -modfile=go.tools.mod github.com/reflog/struct2interface
	$(GOBIN)/struct2interface -f "app" -o "app/app_iface.go" -p "app" -s "App" -i "AppIface" -t ./app/layer_generators/app_iface.go.tmpl
	$(GO) run ./app/layer_generators -in ./app/app_iface.go -out ./app/opentracing/opentracing_layer.go -template ./app/layer_generators/opentracing_layer.go.tmpl

i18n-extract: ## Extract strings for translation from the source code
	$(GO) get -modfile=go.tools.mod github.com/mattermost/mattermost-utilities/mmgotool
	$(GOBIN)/mmgotool i18n extract --portal-dir=""

i18n-check: ## Exit on empty translation strings and translation source strings
	$(GO) get -modfile=go.tools.mod github.com/mattermost/mattermost-utilities/mmgotool
	$(GOBIN)/mmgotool i18n clean-empty --portal-dir="" --check
	$(GOBIN)/mmgotool i18n check-empty-src --portal-dir=""

store-layers: ## Generate layers for the store
	$(GOFLAGS)
	$(GO) generate $(GOFLAGS) ./store

migration-prereqs: ## Builds prerequisite packages for migrations
	$(GO) get -modfile=go.tools.mod github.com/golang-migrate/migrate/v4/cmd/migrate

new-migration: migration-prereqs ## Creates a new migration
	@echo "Generating new migration for postgres"
	$(GOBIN)/migrate create -ext sql -dir db/migrations/postgres -seq $(name)

	@echo "When you are done writing your migration, run 'make migrations'"

migrations-bindata: ## Generates bindata migrations
	$(GO) get -modfile=go.tools.mod github.com/go-bindata/go-bindata/...

	@echo Generating bindata for migrations
	$(GO) generate $(GOFLAGS) ./db/migrations/

filestore-mocks: ## Creates mock files.
	$(GO) get -modfile=go.tools.mod github.com/vektra/mockery/...
	$(GOBIN)/mockery -dir modules/filestore -all -output modules/filestore/mocks -note 'Regenerate this file using `make filestore-mocks`.'

einterfaces-mocks: ## Creates mock files for einterfaces.
	$(GO) get -modfile=go.tools.mod github.com/vektra/mockery/...
	$(GOBIN)/mockery -dir einterfaces -all -output einterfaces/mocks -note 'Regenerate this file using `make einterfaces-mocks`.'

searchengine-mocks: ## Creates mock files for searchengines.
	$(GO) get -modfile=go.tools.mod github.com/vektra/mockery/...
	$(GOBIN)/mockery -dir services/searchengine -all -output services/searchengine/mocks -note 'Regenerate this file using `make searchengine-mocks`.'

gen-serialized: ## Generates serialization methods for hot structs
	# This tool only works at a file level, not at a package level.
	# There will be some warnings about "unresolved identifiers",
	# but that is because of the above problem. Since we are generating
	# methods for all the relevant files at a package level, all
	# identifiers will be resolved. An alternative to remove the warnings
	# would be to temporarily move all the structs to the same file,
	# but that involves a lot of manual work.
	$(GO) get -modfile=go.tools.mod github.com/tinylib/msgp
	$(GOBIN)/msgp -file=./model/session.go -tests=false -o=./model/session_serial_gen.go
	$(GOBIN)/msgp -file=./model/account/user.go -tests=false -o=./model/account/user_serial_gen.go

gqlgen:
	$(GO) get github.com/99designs/gqlgen
	$(GO) run github.com/99designs/gqlgen
	@echo Gqlgen has done generating.

update-dependencies: ## Uses go get -u to update all the dependencies while holding back any that require it.
	@echo Updating Dependencies

	# Update all dependencies (does not update across major versions)
	$(GO) get -u ./...

	# Tidy up
	$(GO) mod tidy

pluginapi: ## Generates api and hooks glue code for plugins
	$(GO) generate $(GOFLAGS) ./modules/plugin
