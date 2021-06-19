module github.com/sitename/sitename

go 1.16

require (
	code.sajari.com/docconv v1.1.1-0.20200701232649-d9ea05fbd50a
	github.com/99designs/gqlgen v0.13.0
	github.com/HdrHistogram/hdrhistogram-go v1.1.0 // indirect
	github.com/Masterminds/squirrel v1.5.0
	github.com/NYTimes/gziphandler v1.1.1
	github.com/advancedlogic/GoOse v0.0.0-20200830213114-1225d531e0ad // indirect
	github.com/avct/uasurfer v0.0.0-20191028135549-26b5daa857f1
	github.com/aws/aws-sdk-go v1.38.64
	github.com/blang/semver v3.5.1+incompatible
	github.com/blevesearch/bleve v1.0.14
	github.com/cespare/xxhash/v2 v2.1.1
	github.com/dgryski/dgoogauth v0.0.0-20190221195224-5a805980a5f3
	github.com/disintegration/imaging v1.6.2
	github.com/dyatlov/go-opengraph v0.0.0-20210112100619-dae8665a5b09
	github.com/francoispqt/gojay v1.2.13
	github.com/fsnotify/fsnotify v1.4.9
	github.com/getsentry/sentry-go v0.11.0
	github.com/go-redis/redis/v8 v8.10.0 // indirect
	github.com/go-sql-driver/mysql v1.6.0
	github.com/golang-migrate/migrate/v4 v4.14.1
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/golang/snappy v0.0.3 // indirect
	github.com/google/uuid v1.2.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/gosimple/slug v1.9.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/jaytaylor/html2text v0.0.0-20200412013138-3577fbdbcff7
	github.com/jmoiron/sqlx v1.3.4
	github.com/json-iterator/go v1.1.11
	github.com/ledongthuc/pdf v0.0.0-20200323191019-23c5852adbd2
	github.com/lib/pq v1.10.2
	github.com/mattermost/go-i18n v1.11.0
	github.com/mattermost/gorp v1.6.2-0.20210419141818-0904a6a388d3
	github.com/mattermost/ldap v0.0.0-20201202150706-ee0e6284187d
	github.com/mattermost/logr v1.0.13
	github.com/mattermost/rsc v0.0.0-20160330161541-bbaefb05eaa0
	github.com/mholt/archiver/v3 v3.5.0
	github.com/minio/minio-go/v7 v7.0.11
	github.com/nyaruka/phonenumbers v1.0.70
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/oov/psd v0.0.0-20210618170533-9fb823ddb631
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pelletier/go-toml v1.9.3 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/rs/cors v1.7.0
	github.com/rwcarlsen/goexif v0.0.0-20190401172101-9e8deecbddbd
	github.com/shopspring/decimal v1.2.0
	github.com/sirupsen/logrus v1.8.1
	github.com/site-name/go-prices v0.0.0-20210616032024-5891e0c6d6c8
	github.com/smartystreets/assertions v1.2.0 // indirect
	github.com/spf13/cobra v1.1.3
	github.com/splitio/go-client/v6 v6.1.0
	github.com/stretchr/testify v1.7.0
	github.com/throttled/throttled v2.2.5+incompatible
	github.com/tinylib/msgp v1.1.5
	github.com/uber/jaeger-client-go v2.29.1+incompatible
	github.com/uber/jaeger-lib v2.4.1+incompatible
	github.com/vektah/gqlparser/v2 v2.2.0 // indirect
	github.com/vmihailenco/msgpack/v5 v5.3.4
	github.com/wiggin77/merror v1.0.3
	github.com/wiggin77/srslog v1.0.1
	go.uber.org/zap v1.17.0
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e
	golang.org/x/image v0.0.0-20210607152325-775e3b0c77b9
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22 // indirect
	golang.org/x/text v0.3.6
	golang.org/x/tools v0.1.3
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/mail.v2 v2.3.1
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v2 v2.4.0
	willnorris.com/go/imageproxy v0.10.0
)

replace github.com/NYTimes/gziphandler v1.1.1 => github.com/agnivade/gziphandler v1.1.2-0.20200815170021-7481835cb745

replace github.com/dyatlov/go-opengraph => github.com/agnivade/go-opengraph v0.0.0-20201221052033-34e69ee2a627

// Hack to prevent the willf/bitset module from being upgraded to 1.2.0.
// They changed the module path from github.com/willf/bitset to
// github.com/bits-and-blooms/bitset and a couple of dependent repos are yet
// to update their module paths.
exclude (
	github.com/RoaringBitmap/roaring v0.7.0
	github.com/RoaringBitmap/roaring v0.7.1
	github.com/willf/bitset v1.2.0
)
