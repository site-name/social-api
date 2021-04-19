package setting

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

// DatabaseSetting holds settings for database
type DatabaseSetting struct {
	Host                     string
	Name                     string
	User                     string
	Passwd                   string
	Schema                   string
	SSLMode                  string
	Path                     string
	LogSQL                   bool
	Charset                  string
	Timeout                  int // seconds
	DBConnectRetries         int
	DBConnectBackoff         time.Duration
	MaxIdleConns             int
	MaxOpenConns             int
	ConnMaxLifetime          time.Duration
	ConnMaxIdleTime          time.Duration
	IterateBufferSize        int
	DataSource               string
	DataSourceReplicas       []string
	DataSourceSearchReplicas []string
	ReplicaLagSettings       []*ReplicaLagSetting
}

type ReplicaLagSetting struct {
	DataSource       string
	QueryAbsoluteLag string
	QueryTimeLag     string
}

var (
	// Database holds the database settings
	Database = DatabaseSetting{
		Timeout:           500,
		IterateBufferSize: 50,
	}
)

// InitDBConfig loads the database settings
func InitDBConfig() {
	sec := Cfg.Section("database")

	defaultCharset := "utf8"

	Database.Host = sec.Key("HOST").String()
	Database.Name = sec.Key("NAME").MustString("sitename")
	Database.User = sec.Key("USER").MustString("minh")
	if len(Database.Passwd) == 0 {
		Database.Passwd = sec.Key("PASSWD").MustString("anhyeuem98")
	}
	Database.Schema = sec.Key("SCHEMA").MustString("postgres")
	Database.SSLMode = sec.Key("SSL_MODE").MustString("disable")
	Database.Charset = sec.Key("CHARSET").In(defaultCharset, []string{"utf8", "utf8mb4"})
	Database.Path = sec.Key("PATH").MustString(filepath.Join(AppDataPath, "gitea.db"))
	Database.Timeout = sec.Key("SQLITE_TIMEOUT").MustInt(500)
	Database.MaxIdleConns = sec.Key("MAX_IDLE_CONNS").MustInt(2)
	Database.ConnMaxLifetime = sec.Key("CONN_MAX_LIFE_TIME").MustDuration(0)
	Database.MaxOpenConns = sec.Key("MAX_OPEN_CONNS").MustInt(0)
	Database.IterateBufferSize = sec.Key("ITERATE_BUFFER_SIZE").MustInt(50)
	Database.LogSQL = sec.Key("LOG_SQL").MustBool(true)
	Database.DBConnectRetries = sec.Key("DB_RETRIES").MustInt(10)
	Database.DBConnectBackoff = sec.Key("DB_RETRY_BACKOFF").MustDuration(3 * time.Second)

	Database.ConnMaxIdleTime = sec.Key("CONN_MAX_IDLE_TIME").MustDuration(0) // NOTE: database idle time is only supported in go >= 1.15
	Database.DataSource, _ = DBConnStr()                                     // error is nil

	replicaSec := Cfg.Section("database-replica")
	// defaults settings for replications
	Database.DataSourceReplicas = replicaSec.Key("DB_SOURCE_REPLICAS").Strings(",")
	Database.DataSourceSearchReplicas = replicaSec.Key("DB_SOURCE_SEARCH_REPLICAS").Strings(",")
	Database.ReplicaLagSettings = []*ReplicaLagSetting{} // NOTE: fix me
}

// DBConnStr returns database connection string
func DBConnStr() (string, error) {
	var Param = "?"
	if strings.Contains(Database.Name, Param) {
		Param = "&"
	}
	connStr := getPostgreSQLConnectionString(Database.Host, Database.User, Database.Passwd, Database.Name, Param, Database.SSLMode)
	return connStr, nil
}

// parsePostgreSQLHostPort parses given input in various forms defined in
// https://www.postgresql.org/docs/current/static/libpq-connect.html#LIBPQ-CONNSTRING
// and returns proper host and port number.
func parsePostgreSQLHostPort(info string) (string, string) {
	host, port := "127.0.0.1", "5432"
	if strings.Contains(info, ":") && !strings.HasSuffix(info, "]") {
		idx := strings.LastIndex(info, ":")
		host = info[:idx]
		port = info[idx+1:]
	} else if len(info) > 0 {
		host = info
	}
	return host, port
}

func getPostgreSQLConnectionString(dbHost, dbUser, dbPasswd, dbName, dbParam, dbsslMode string) (connStr string) {
	host, port := parsePostgreSQLHostPort(dbHost)
	if host[0] == '/' { // looks like a unix socket
		connStr = fmt.Sprintf("postgres://%s:%s@:%s/%s%ssslmode=%s&host=%s",
			url.PathEscape(dbUser), url.PathEscape(dbPasswd), port, dbName, dbParam, dbsslMode, host)
	} else {
		connStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s%ssslmode=%s",
			url.PathEscape(dbUser), url.PathEscape(dbPasswd), host, port, dbName, dbParam, dbsslMode)
	}
	return
}
