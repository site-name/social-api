package setting

import (
	"time"

	"code.gitea.io/gitea/modules/log"
)

var (
	// CORSConfig defines CORS settings
	CORSConfig = struct {
		Enabled          bool
		Scheme           string
		AllowDomain      []string
		AllowSubdomain   bool
		Methods          []string
		MaxAge           time.Duration
		AllowCredentials bool
	}{
		Enabled: false,
		MaxAge:  10 * time.Minute,
	}
)

func newCORSService() {
	sec := Cfg.Section("cors")
	if err := sec.MapTo(&CORSConfig); err != nil {
		log.Fatal("Failed to map cors settings: %v", err)
	}

	if CORSConfig.Enabled {
		log.Info("CORS Service Enabled")
	}
}
