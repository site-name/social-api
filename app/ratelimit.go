package app

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
	"github.com/throttled/throttled"
	"github.com/throttled/throttled/store/memstore"
)

type RateLimiter struct {
	throttledRateLimiter *throttled.GCRARateLimiter
	useAuth              bool
	useIP                bool
	header               string
	trustedProxyIPHeader []string
}

// NewRateLimiter creates new RateLimiter
func NewRateLimiter(settings *model.RateLimitSettings, trustedProxyIPHeader []string) (*RateLimiter, error) {
	store, err := memstore.New(*settings.MemoryStoreSize)
	if err != nil {
		return nil, errors.Wrap(err, "api.server.start_server.rate_limiting_memory_store")
	}

	quota := throttled.RateQuota{
		MaxRate:  throttled.PerSec(*settings.PerSec),
		MaxBurst: *settings.MaxBurst,
	}

	throttledRateLimiter, err := throttled.NewGCRARateLimiter(store, quota)
	if err != nil {
		return nil, errors.Wrap(err, "api.server.start_server.rate_limiting_rate_limiter")
	}

	return &RateLimiter{
		throttledRateLimiter: throttledRateLimiter,
		useAuth:              *settings.VaryByUser,
		useIP:                *settings.VaryByRemoteAddr,
		header:               settings.VaryByHeader,
		trustedProxyIPHeader: trustedProxyIPHeader,
	}, nil
}

func (rl *RateLimiter) GenerateKey(r *http.Request) string {
	key := ""
	if rl.useAuth {
		token, tokenLocation := ParseAuthTokenFromRequest(r)
		if tokenLocation != TokenLocationNotFound {
			key += token
		} else if rl.useIP {
			key += util.GetIPAddress(r, rl.trustedProxyIPHeader)
		}
	} else if rl.useIP {
		key += util.GetIPAddress(r, rl.trustedProxyIPHeader)
	}

	// Note that most of the time the user won't have to set this because the utils.GetIpAddress above tries the
	// most common headers anyway.
	if rl.header != "" {
		key += strings.ToLower(r.Header.Get(rl.header))
	}

	return key
}

func (rl *RateLimiter) RateLimitWriter(key string, w http.ResponseWriter) bool {
	limited, context, err := rl.throttledRateLimiter.RateLimit(key, 1)
	if err != nil {
		slog.Error("Internal server error when rate limiting. Rate limiting broken", slog.Err(err))
		return false
	}

	setRateLimitHeaders(w, context)

	if limited {
		slog.Debug("Denied due to throttling settings code=429", slog.String("key", key))
		http.Error(w, "limit exceeded", http.StatusTooManyRequests)
	}

	return limited
}

func (rl *RateLimiter) UserIdRateLimit(userID string, w http.ResponseWriter) bool {
	if rl.useAuth {
		return rl.RateLimitWriter(userID, w)
	}

	return false
}

func (rl *RateLimiter) RateLimitHandler(wrap http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := rl.GenerateKey(r)

		if !rl.RateLimitWriter(key, w) {
			wrap.ServeHTTP(w, r)
		}
	})
}

// Copied from https://github.com/throttled/throttled http.go
func setRateLimitHeaders(w http.ResponseWriter, context throttled.RateLimitResult) {
	if v := context.Limit; v >= 0 {
		w.Header().Add("X-RateLimit-Limit", strconv.Itoa(v))
	}

	if v := context.Remaining; v >= 0 {
		w.Header().Add("X-RateLimit-Remaining", strconv.Itoa(v))
	}

	if v := context.ResetAfter; v >= 0 {
		vi := int(math.Ceil(v.Seconds()))
		w.Header().Add("X-RateLimit-Reset", strconv.Itoa(vi))
	}

	if v := context.RetryAfter; v >= 0 {
		vi := int(math.Ceil(v.Seconds()))
		w.Header().Add("Retry-After", strconv.Itoa(vi))
	}
}
