package localcachelayer

import (
	"runtime"
	"time"

	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
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

	ClearCacheMessageData = ""

	ChannelCacheSec = 15 * 60 // 15 mins
)

type LocalCacheStore struct {
	store.Store
	metrics einterfaces.MetricsInterface
	cluster einterfaces.ClusterInterface

	user                  *LocalCacheUserStore
	userProfileByIdsCache cache.Cache
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
		InvalidateClusterEvent: model.CLUSTER_EVENT_INVALIDATE_CACHE_FOR_PROFILE_BY_IDS,
		Striped:                true,
		StripedBuckets:         util.Max(runtime.NumCPU()-1, 1),
	}); err != nil {
		return
	}
	localCacheStore.user = &LocalCacheUserStore{
		UserStore:                     baseStore.User(),
		rootStore:                     &localCacheStore,
		userProfileByIdsInvalidations: make(map[string]bool),
	}

	if cluster != nil {
		cluster.RegisterClusterMessageHandler(model.CLUSTER_EVENT_INVALIDATE_CACHE_FOR_PROFILE_BY_IDS, localCacheStore.user.handleClusterInvalidateScheme)
	}

	return
}

func (s LocalCacheStore) User() store.UserStore {
	return s.user
}

func (s LocalCacheStore) DropAllTables() {
	s.Invalidate()
	s.Store.DropAllTables()
}

func (s *LocalCacheStore) doClearCacheCluster(cache cache.Cache) {
	cache.Purge()
	if s.cluster != nil {
		msg := &model.ClusterMessage{
			Event:    cache.GetInvalidateClusterEvent(),
			SendType: model.CLUSTER_SEND_BEST_EFFORT,
			Data:     ClearCacheMessageData,
		}
		s.cluster.SendClusterMessage(msg)
	}
}

func (s *LocalCacheStore) Invalidate() {
	s.doClearCacheCluster(s.userProfileByIdsCache)
}
