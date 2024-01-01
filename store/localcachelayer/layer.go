package localcachelayer

import (
	"runtime"
	"time"

	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/services/cache"
	"github.com/sitename/sitename/store"
)

const (
	ReactionCacheSize = 20000
	ReactionCacheSec  = 30 * 60

	RoleCacheSize = 20000
	RoleCacheSec  = 30 * 60

	SchemeCacheSize = 20000
	SchemeCacheSec  = 30 * 60

	FileInfoCacheSize = 25000
	FileInfoCacheSec  = 30 * 60

	// ChannelGuestCountCacheSize = model.CHANNEL_CACHE_SIZE
	// ChannelGuestCountCacheSec  = 30 * 60

	WebhookCacheSize = 25000
	WebhookCacheSec  = 15 * 60

	EmojiCacheSize = 5000
	EmojiCacheSec  = 30 * 60

	// ChannelPinnedPostsCounsCacheSize = model.CHANNEL_CACHE_SIZE
	// ChannelPinnedPostsCountsCacheSec = 30 * 60

	// ChannelMembersCountsCacheSize = model.CHANNEL_CACHE_SIZE
	// ChannelMembersCountsCacheSec  = 30 * 60

	LastPostsCacheSize = 20000
	LastPostsCacheSec  = 30 * 60

	TermsOfServiceCacheSize = 20000
	TermsOfServiceCacheSec  = 30 * 60
	LastPostTimeCacheSize   = 25000
	LastPostTimeCacheSec    = 15 * 60

	UserProfileByIDCacheSize = 20000
	UserProfileByIDSec       = 30 * 60

	// ProfilesInChannelCacheSize  = model.CHANNEL_CACHE_SIZE
	// PROFILES_IN_ChannelCacheSec = 15 * 60

	TeamCacheSize = 20000
	TeamCacheSec  = 30 * 60

	ChannelCacheSec = 15 * 60 // 15 mins

	CategoryCacheSize = 25000
	CategoryCacheSec  = 30 * 60
)

var clearCacheMessageData = []byte("")

type LocalCacheStore struct {
	store.Store
	metrics einterfaces.MetricsInterface
	cluster einterfaces.ClusterInterface

	user                  *LocalCacheUserStore
	userProfileByIdsCache cache.Cache

	role                 LocalCacheRoleStore
	roleCache            cache.Cache
	rolePermissionsCache cache.Cache
}

func NewLocalCacheLayer(baseStore store.Store, metrics einterfaces.MetricsInterface, cluster einterfaces.ClusterInterface, cacheProvider cache.Provider) (localCacheStore LocalCacheStore, err error) {
	localCacheStore = LocalCacheStore{
		Store:   baseStore,
		cluster: cluster,
		metrics: metrics,
	}

	// Users
	if localCacheStore.userProfileByIdsCache, err = cacheProvider.NewCache(&cache.CacheOptions{
		Size:                   UserProfileByIDCacheSize,
		Name:                   "UserProfileByIds",
		DefaultExpiry:          UserProfileByIDSec * time.Second,
		InvalidateClusterEvent: model_helper.ClusterEventInvalidateCacheForProfileByIds,
		Striped:                true,
		StripedBuckets:         max(runtime.NumCPU()-1, 1),
	}); err != nil {
		return
	}
	localCacheStore.user = &LocalCacheUserStore{
		UserStore:                     baseStore.User(),
		rootStore:                     &localCacheStore,
		userProfileByIdsInvalidations: make(map[string]bool),
	}

	// Roles
	if localCacheStore.roleCache, err = cacheProvider.NewCache(&cache.CacheOptions{
		Size:                   RoleCacheSize,
		Name:                   "Role",
		DefaultExpiry:          RoleCacheSec * time.Second,
		InvalidateClusterEvent: model_helper.ClusterEventInvalidateCacheForRoles,
		Striped:                true,
		StripedBuckets:         max(runtime.NumCPU()-1, 1),
	}); err != nil {
		return
	}
	if localCacheStore.rolePermissionsCache, err = cacheProvider.NewCache(&cache.CacheOptions{
		Size:                   RoleCacheSize,
		Name:                   "RolePermission",
		DefaultExpiry:          RoleCacheSec * time.Second,
		InvalidateClusterEvent: model_helper.ClusterEventInvalidateCacheForRolePermissions,
	}); err != nil {
		return
	}
	localCacheStore.role = LocalCacheRoleStore{RoleStore: baseStore.Role(), rootStore: &localCacheStore}

	if cluster != nil {
		cluster.RegisterClusterMessageHandler(model_helper.ClusterEventInvalidateCacheForRoles, localCacheStore.role.handleClusterInvalidateRole)
		cluster.RegisterClusterMessageHandler(model_helper.ClusterEventInvalidateCacheForProfileByIds, localCacheStore.user.handleClusterInvalidateScheme)
		cluster.RegisterClusterMessageHandler(model_helper.ClusterEventInvalidateCacheForRolePermissions, localCacheStore.role.handleClusterInvalidateRolePermissions)
		cluster.RegisterClusterMessageHandler(model_helper.ClusterEventInvalidateCacheForProfileByIds, localCacheStore.user.handleClusterInvalidateScheme)
	}

	return
}

func (s LocalCacheStore) User() store.UserStore {
	return s.user
}

func (s LocalCacheStore) Role() store.RoleStore {
	return s.role
}

func (s LocalCacheStore) DropAllTables() {
	s.Invalidate()
	s.Store.DropAllTables()
}

func (s *LocalCacheStore) doInvalidateCacheCluster(cache cache.Cache, key string) {
	cache.Remove(key)
	if s.cluster != nil {
		msg := &model_helper.ClusterMessage{
			Event:    cache.GetInvalidateClusterEvent(),
			SendType: model_helper.ClusterSendBestEffort,
			Data:     []byte(key),
		}
		s.cluster.SendClusterMessage(msg)
	}
}

func (s *LocalCacheStore) doStandardAddToCache(cache cache.Cache, key string, value interface{}) {
	cache.SetWithDefaultExpiry(key, value)
}

func (s *LocalCacheStore) doStandardReadCache(cache cache.Cache, key string, value interface{}) error {
	err := cache.Get(key, value)
	if err == nil {
		if s.metrics != nil {
			s.metrics.IncrementMemCacheHitCounter(cache.Name())
		}
		return nil
	}
	if s.metrics != nil {
		s.metrics.IncrementMemCacheMissCounter(cache.Name())
	}
	return err
}

func (s *LocalCacheStore) doClearCacheCluster(cache cache.Cache) {
	cache.Purge()
	if s.cluster != nil {
		msg := &model_helper.ClusterMessage{
			Event:    cache.GetInvalidateClusterEvent(),
			SendType: model_helper.ClusterSendBestEffort,
			Data:     clearCacheMessageData,
		}
		s.cluster.SendClusterMessage(msg)
	}
}

func (s *LocalCacheStore) Invalidate() {
	s.doClearCacheCluster(s.userProfileByIdsCache)
	s.doClearCacheCluster(s.roleCache)
}
