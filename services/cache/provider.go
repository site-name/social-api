package cache

import (
	"time"

	"github.com/sitename/sitename/model/cluster"
)

// CacheOptions contains options for initializaing a cache
type CacheOptions struct {
	Size                   int
	DefaultExpiry          time.Duration
	Name                   string
	InvalidateClusterEvent cluster.ClusterEvent
	Striped                bool
	StripedBuckets         int
}

// Provider is a provider for Cache
type Provider interface {
	// NewCache creates a new cache with given options.
	NewCache(opts *CacheOptions) (Cache, error)
	// Connect opens a new connection to the cache using specific provider parameters.
	Connect() error
	// Close releases any resources used by the cache provider.
	Close() error
}

type cacheProvider struct{}

// NewProvider creates a new CacheProvider
func NewProvider() Provider {
	return &cacheProvider{}
}

// NewCache creates a new cache with given opts
func (c *cacheProvider) NewCache(opts *CacheOptions) (Cache, error) {
	if opts.Striped {
		return NewLRUStriped(LRUOptions{
			Name:                   opts.Name,
			Size:                   opts.Size,
			DefaultExpiry:          opts.DefaultExpiry,
			InvalidateClusterEvent: opts.InvalidateClusterEvent,
			StripedBuckets:         opts.StripedBuckets,
		})
	}
	return NewLRU(LRUOptions{
		Name:                   opts.Name,
		Size:                   opts.Size,
		DefaultExpiry:          opts.DefaultExpiry,
		InvalidateClusterEvent: opts.InvalidateClusterEvent,
	}), nil
}

// Connect opens a new connection to the cache using specific provider parameters.
func (c *cacheProvider) Connect() error {
	return nil
}

// Close releases any resources used by the cache provider.
func (c *cacheProvider) Close() error {
	return nil
}
