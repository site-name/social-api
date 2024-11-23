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
	$(GO) install github.com/reflog/struct2interface@v0.6.1
	$(GOBIN)/struct2interface -f "app" -o "app/app_iface.go" -p "app" -s "App" -i "AppIface" -t ./app/layer_generators/app_iface.go.tmpl
	$(GO) run ./app/layer_generators -in ./app/app_iface.go -out ./app/opentracing/opentracing_layer.go -template ./app/layer_generators/opentracing_layer.go.tmpl
	$(GOBIN)/struct2interface -f "app/checkout" -o "app/sub_app_iface/checkout_iface.go" -p "checkout" -s "ServiceCheckout" -i "CheckoutService" -t ./app/layer_generators/checkout_iface.go.tmpl
	$(GOBIN)/struct2interface -f "app/account" -o "app/sub_app_iface/account_iface.go" -p "account" -s "ServiceAccount" -i "AccountService" -t ./app/layer_generators/account_iface.go.tmpl
	$(GOBIN)/struct2interface -f "app/attribute" -o "app/sub_app_iface/attribute_iface.go" -p "attribute" -s "ServiceAttribute" -i "AttributeService" -t ./app/layer_generators/attribute_iface.go.tmpl
	$(GOBIN)/struct2interface -f "app/channel" -o "app/sub_app_iface/channel_iface.go" -p "channel" -s "ServiceChannel" -i "ChannelService" -t ./app/layer_generators/channel_iface.go.tmpl
	$(GOBIN)/struct2interface -f "app/csv" -o "app/sub_app_iface/csv_iface.go" -p "csv" -s "ServiceCsv" -i "CsvService" -t ./app/layer_generators/csv_iface.go.tmpl
	$(GOBIN)/struct2interface -f "app/discount" -o "app/sub_app_iface/discount_iface.go" -p "discount" -s "ServiceDiscount" -i "DiscountService" -t ./app/layer_generators/discount_iface.go.tmpl
	$(GOBIN)/struct2interface -f "app/file" -o "app/sub_app_iface/file_iface.go" -p "file" -s "ServiceFile" -i "FileService" -t ./app/layer_generators/file_iface.go.tmpl
	$(GOBIN)/struct2interface -f "app/giftcard" -o "app/sub_app_iface/giftcard_iface.go" -p "giftcard" -s "ServiceGiftcard" -i "GiftcardService" -t ./app/layer_generators/giftcard_iface.go.tmpl
	$(GOBIN)/struct2interface -f "app/invoice" -o "app/sub_app_iface/invoice_iface.go" -p "invoice" -s "ServiceInvoice" -i "InvoiceService" -t ./app/layer_generators/invoice_iface.go.tmpl
	$(GOBIN)/struct2interface -f "app/menu" -o "app/sub_app_iface/menu_iface.go" -p "menu" -s "ServiceMenu" -i "MenuService" -t ./app/layer_generators/menu_iface.go.tmpl
	$(GOBIN)/struct2interface -f "app/order" -o "app/sub_app_iface/order_iface.go" -p "order" -s "ServiceOrder" -i "OrderService" -t ./app/layer_generators/order_iface.go.tmpl
	$(GOBIN)/struct2interface -f "app/page" -o "app/sub_app_iface/page_iface.go" -p "page" -s "ServicePage" -i "PageService" -t ./app/layer_generators/page_iface.go.tmpl
	$(GOBIN)/struct2interface -f "app/payment" -o "app/sub_app_iface/payment_iface.go" -p "payment" -s "ServicePayment" -i "PaymentService" -t ./app/layer_generators/payment_iface.go.tmpl
	$(GOBIN)/struct2interface -f "app/plugin" -o "app/sub_app_iface/plugin_iface.go" -p "plugin" -s "ServicePlugin" -i "PluginService" -t ./app/layer_generators/plugin_iface.go.tmpl
	$(GOBIN)/struct2interface -f "app/product" -o "app/sub_app_iface/product_iface.go" -p "product" -s "ServiceProduct" -i "ProductService" -t ./app/layer_generators/product_iface.go.tmpl
	$(GOBIN)/struct2interface -f "app/seo" -o "app/sub_app_iface/seo_iface.go" -p "seo" -s "ServiceSeo" -i "SeoService" -t ./app/layer_generators/seo_iface.go.tmpl
	$(GOBIN)/struct2interface -f "app/shipping" -o "app/sub_app_iface/shipping_iface.go" -p "shipping" -s "ServiceShipping" -i "ShippingService" -t ./app/layer_generators/shipping_iface.go.tmpl
	$(GOBIN)/struct2interface -f "app/shop" -o "app/sub_app_iface/shop_iface.go" -p "shop" -s "ServiceShop" -i "ShopService" -t ./app/layer_generators/shop_iface.go.tmpl
	$(GOBIN)/struct2interface -f "app/warehouse" -o "app/sub_app_iface/warehouse_iface.go" -p "warehouse" -s "ServiceWarehouse" -i "WarehouseService" -t ./app/layer_generators/warehouse_iface.go.tmpl
	$(GOBIN)/struct2interface -f "app/webhook" -o "app/sub_app_iface/webhook_iface.go" -p "webhook" -s "ServiceWebhook" -i "WebhookService" -t ./app/layer_generators/webhook_iface.go.tmpl
	$(GOBIN)/struct2interface -f "app/wishlist" -o "app/sub_app_iface/wishlist_iface.go" -p "wishlist" -s "ServiceWishlist" -i "WishlistService" -t ./app/layer_generators/wishlist_iface.go.tmpl
	$(GOBIN)/struct2interface -f "app/plugin" -o "app/plugin/interfaces/plugin_manager_iface.go" -p "plugin" -s "PluginManager" -i "PluginManagerInterface" -t ./app/layer_generators/plugin_manager_iface.go.tmpl

i18n-check: ## Exit on empty translation strings and translation source strings
	$(GO) get -modfile=go.tools.mod github.com/mattermost/mattermost-utilities/mmgotool
	$(GOBIN)/mmgotool i18n clean-empty --portal-dir="" --check
	$(GOBIN)/mmgotool i18n check-empty-src --portal-dir=""

store-layers: ## Generate layers for the store
	$(GO) generate $(GOFLAGS) ./store

migration-prereqs: ## Builds prerequisite packages for migrations
	$(GO) get -modfile=go.tools.mod github.com/golang-migrate/migrate/v4/cmd/migrate

new-migration: ## Creates a new migration. Run with make new-migration name=<>
	$(GO) install github.com/mattermost/morph/cmd/morph@master
	@echo "Generating new migration for mysql"
	$(GOBIN)/morph generate $(name) --driver mysql --dir db/migrations --sequence

	@echo "Generating new migration for postgres"
	$(GOBIN)/morph generate $(name) --driver postgres --dir db/migrations --sequence

	@echo "When you are done writing your migration, run 'make migrations-bindata'"

# migrations-bindata: ## Generates bindata migrations
# 	$(GO) get -d -modfile=go.tools.mod github.com/go-bindata/go-bindata/...

# 	@echo Generating bindata for migrations
# 	$(GO) generate $(GOFLAGS) ./db/migrations/

filestore-mocks: ## Creates mock files.
	$(GO) install github.com/vektra/mockery/v2/...@v2.10.4
	$(GOBIN)/mockery --dir shared/filestore --all --output shared/filestore/mocks --note 'Regenerate this file using `make filestore-mocks`.'

einterfaces-mocks: ## Creates mock files for einterfaces.
	$(GO) install github.com/vektra/mockery/v2/...@v2.10.4
	$(GOBIN)/mockery --dir einterfaces --all --output einterfaces/mocks --note 'Regenerate this file using `make einterfaces-mocks`.'

gen-serialized: ## Generates serialization methods for hot structs
  #This tool only works at a file level, not at a package level.
  #There will be some warnings about "unresolved identifiers",
  #but that is because of the above problem. Since we are generating
  #methods for all the relevant files at a package level, all
  #identifiers will be resolved. An alternative to remove the warnings
  #would be to temporarily move all the structs to the same file,
  #but that involves a lot of manual work.
	$(GO) install github.com/tinylib/msgp
	$(GOBIN)/msgp -file=./model/sessions.go -tests=false -o=./model/session_serial_gen.go
	$(GOBIN)/msgp -file=./model/users.go -tests=false -o=./model/user_serial_gen.go

update-dependencies: ## Uses go get -u to update all the dependencies while holding back any that require it.
	@echo Updating Dependencies

	# Update all dependencies (does not update across major versions)
	$(GO) get -u ./...

	# Tidy up
	$(GO) mod tidy

pluginapi: ## Generates api and hooks glue code for plugins
	$(GO) run modules/plugin/interface_generator/main.go

server:
	@echo Starting the server...
	$(GO) run cmd/sitename/main.go server

store-mocks: ## Creates mock files.
	$(GO) install github.com/vektra/mockery/v2/...@v2.23.2
	$(GOBIN)/mockery --dir store --name ".*Store" --output store/storetest/mocks --note 'Regenerate this file using `make store-mocks`.'

# add ./bin directory to system PATH
export PATH := $(GOBIN):$(PATH)

model-gen:
	$(GO) install github.com/volatiletech/sqlboiler/v4@latest
	$(GO) install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest

	$(GOBIN)/sqlboiler psql

migrate:
	$(GO) install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	$(GOBIN)/migrate -path db/migrations/postgres -database postgres://sitename:sitename@localhost:5432/sitename?sslmode=disable up

searchengine-mocks: ## Creates mock files for searchengines.
	$(GO) install github.com/vektra/mockery/v2/...@v2.10.4
	$(GOBIN)/mockery --dir services/searchengine --all --output services/searchengine/mocks --note 'Regenerate this file using `make searchengine-mocks`.'

categories:
	$(GO) generate $(GOFLAGS) ./model/generate

i18n-extract: ## Extract strings for translation from the source code
	$(GO) install github.com/mattermost/mattermost-utilities/mmgotool@fdf2cd651b261bcd511a32da33dd76febedd44a8
	$(GOBIN)/mmgotool i18n extract --portal-dir=""

graphql-gen:
	$(GO) run api/schemas/gen.go
