package localcachelayer

import (
	"sync"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type LocalCacheUserStore struct {
	store.UserStore
	rootStore                     *LocalCacheStore
	userProfileByIdsMut           sync.Mutex
	userProfileByIdsInvalidations map[string]bool
}

func (s *LocalCacheUserStore) handleClusterInvalidateScheme(msg *model.ClusterMessage) {
	if msg.Data == ClearCacheMessageData {
		s.rootStore.userProfileByIdsCache.Purge()
	} else {
		s.userProfileByIdsMut.Lock()
		s.userProfileByIdsInvalidations[msg.Data] = true
		s.userProfileByIdsMut.Unlock()
		s.rootStore.userProfileByIdsCache.Remove(msg.Data)
	}
}