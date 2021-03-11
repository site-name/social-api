package setting

import (
	"strings"
	"time"

	"code.gitea.io/gitea/modules/log"
)

// Cache represents cache settings
type Cache struct {
	Enabled  bool
	Adapter  string
	Interval int
	Conn     string
	TTL      time.Duration `init:"ITEM_TTL"`
}

var (
	// CacheService the global cache
	CacheService = struct {
		Cache `init:"cache"`
	}{
		Cache: Cache{
			Enabled:  true,
			Adapter:  "memory",
			Interval: 60,
			TTL:      16 * time.Hour,
		},
	}
)

// MemcacheMaxTTL represents the maximum memcache TTL
const MemcacheMaxTTL = 30 * 24 * time.Hour

func newCacheService() {
	sec := Cfg.Section("cache")
	if err := sec.MapTo(&CacheService); err != nil {
		log.Fatal("Failed to map Cache settings: %v", err)
	}

	CacheService.Adapter = sec.Key("ADAPTER").In("memory", []string{"memory", "redis", "memcache"})
	switch CacheService.Adapter {
	case "memory":
	case "redis", "memcache":
		CacheService.Conn = strings.Trim(sec.Key("HOST").String(), "\" ")
	case "": // disable cache
		CacheService.Enabled = false
	default:
		log.Fatal("Unknown cache adapter: %s", CacheService.Adapter)
	}

	if CacheService.Enabled {
		log.Info("Cache Service Enabled")
	} else {
		log.Warn("Cache Service Disabled so that captcha disabled too")
		// captcha depends on cache service
		Service.EnableCaptcha = false
	}
}

// TTLSeconds returns the TTLSeconds or unix timestamp for memcache
func (c Cache) TTLSeconds() int64 {
	if c.Adapter == "memcache" && c.TTL > MemcacheMaxTTL {
		return time.Now().Add(c.TTL).Unix()
	}
	return int64(c.TTL.Seconds())
}
