module github.com/sitename/sitename

go 1.16

require (
	code.sajari.com/docconv v1.1.1-0.20200701232649-d9ea05fbd50a
	github.com/99designs/gqlgen v0.13.0
	github.com/HdrHistogram/hdrhistogram-go v1.1.0 // indirect
	github.com/Masterminds/squirrel v1.5.0
	github.com/advancedlogic/GoOse v0.0.0-20210708011750-e3d1acc33807 // indirect
	github.com/avct/uasurfer v0.0.0-20191028135549-26b5daa857f1
	github.com/aws/aws-sdk-go v1.40.7
	github.com/blang/semver v3.5.1+incompatible
	github.com/blevesearch/bleve v1.0.14
	github.com/cespare/xxhash/v2 v2.1.1
	github.com/dgryski/dgoogauth v0.0.0-20190221195224-5a805980a5f3
	github.com/disintegration/imaging v1.6.2
	github.com/dyatlov/go-opengraph v0.0.0-20210112100619-dae8665a5b09
	github.com/fortytw2/leaktest v1.3.0 // indirect
	github.com/francoispqt/gojay v1.2.13
	github.com/fsnotify/fsnotify v1.4.9
	github.com/getsentry/sentry-go v0.11.0
	github.com/go-redis/redis/v8 v8.11.0 // indirect
	github.com/go-sql-driver/mysql v1.6.0
	github.com/golang-migrate/migrate/v4 v4.14.1
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/google/uuid v1.3.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/gosimple/slug v1.9.0
	github.com/h2non/go-is-svg v0.0.0-20160927212452-35e8c4b0612c
	github.com/hako/durafmt v0.0.0-20210608085754-5c1018a4e16b
	github.com/hashicorp/go-hclog v0.16.2
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/go-plugin v1.4.2
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/memberlist v0.2.4
	github.com/jaytaylor/html2text v0.0.0-20200412013138-3577fbdbcff7
	github.com/jmoiron/sqlx v1.3.4
	github.com/json-iterator/go v1.1.11
	github.com/ledongthuc/pdf v0.0.0-20210621053716-e28cb8259002
	github.com/lib/pq v1.10.2
	github.com/mattermost/go-i18n v1.11.0
	github.com/mattermost/gorp v1.6.2-0.20210419141818-0904a6a388d3
	github.com/mattermost/gosaml2 v0.3.3
	github.com/mattermost/gziphandler v0.0.1
	github.com/mattermost/ldap v0.0.0-20201202150706-ee0e6284187d
	github.com/mattermost/logr v1.0.13
	github.com/mattermost/rsc v0.0.0-20160330161541-bbaefb05eaa0
	github.com/mholt/archiver/v3 v3.5.0
	github.com/minio/minio-go/v7 v7.0.12
	github.com/nyaruka/phonenumbers v1.0.70
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/olivere/elastic v6.2.36+incompatible // indirect
	github.com/oov/psd v0.0.0-20210618170533-9fb823ddb631
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/rs/cors v1.8.0
	github.com/russellhaering/goxmldsig v1.1.0 // indirect
	github.com/rwcarlsen/goexif v0.0.0-20190401172101-9e8deecbddbd
	github.com/sirupsen/logrus v1.8.1
	github.com/site-name/go-prices v0.0.0-20210722081319-cce771fbc863
	github.com/site-name/i18naddress v0.0.0-20210716143658-e5dba1b87ab3
	github.com/smartystreets/assertions v1.2.0 // indirect
	github.com/spf13/cobra v1.2.1
	github.com/splitio/go-client/v6 v6.1.0
	github.com/stretchr/testify v1.7.0
	github.com/throttled/throttled v2.2.5+incompatible
	github.com/tinylib/msgp v1.1.6
	github.com/tylerb/graceful v1.2.15
	github.com/uber/jaeger-client-go v2.29.1+incompatible
	github.com/uber/jaeger-lib v2.4.1+incompatible
	github.com/vektah/gqlparser/v2 v2.2.0
	github.com/vmihailenco/msgpack/v5 v5.3.4
	github.com/wiggin77/merror v1.0.3
	github.com/wiggin77/srslog v1.0.1
	go.uber.org/zap v1.18.1
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97
	golang.org/x/image v0.0.0-20210628002857-a66eb6448b8d
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/text v0.3.6
	golang.org/x/tools v0.1.5
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/mail.v2 v2.3.1
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/olivere/elastic.v6 v6.2.36
	gopkg.in/yaml.v2 v2.4.0
	willnorris.com/go/imageproxy v0.10.0
)

// Hack to prevent the willf/bitset module from being upgraded to 1.2.0.
// They changed the module path from github.com/willf/bitset to
// github.com/bits-and-blooms/bitset and a couple of dependent repos are yet
// to update their module paths.
exclude (
	github.com/RoaringBitmap/roaring v0.7.0
	github.com/RoaringBitmap/roaring v0.7.1
	github.com/willf/bitset v1.2.0
)
