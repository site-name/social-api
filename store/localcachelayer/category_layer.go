package localcachelayer

import (
	"bytes"
	"context"
	"sync"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/sqlstore"
)

type LocalCacheCategoryStore struct {
	store.CategoryStore
	categoryByIdMut           sync.Mutex
	rootStore                 *LocalCacheStore
	categoryByIdInvalidations map[string]bool
}

func (l *LocalCacheCategoryStore) handleClusterInvalidateCategoryById(msg *model.ClusterMessage) {
	if bytes.Equal(msg.Data, clearCacheMessageData) {
		l.rootStore.categoryByIdsCache.Purge()
	} else {
		l.categoryByIdMut.Lock()
		l.categoryByIdInvalidations[string(msg.Data)] = true
		l.categoryByIdMut.Unlock()
		l.rootStore.categoryByIdsCache.Remove(string(msg.Data))
	}
}

func (l *LocalCacheCategoryStore) Get(ctx context.Context, id string, allowFromCache bool) (*model.Category, error) {
	if allowFromCache {
		cate, ok := l.getFromCacheById(id)
		if ok {
			return cate, nil
		}
	}

	// if it was invalidated, then we need to query master
	l.categoryByIdMut.Lock()
	if l.categoryByIdInvalidations[id] {
		ctx = sqlstore.WithMaster(ctx)
		delete(l.categoryByIdInvalidations, id)
	}
	l.categoryByIdMut.Unlock()

	cate, err := l.CategoryStore.Get(ctx, id, allowFromCache)
	if allowFromCache && err == nil {
		l.addToCache(cate)
	}
	return cate, err
}

func (l *LocalCacheCategoryStore) addToCache(cate *model.Category) {
	l.rootStore.doStandardAddToCache(l.rootStore.categoryByIdsCache, cate.Id, cate)
}

func (l *LocalCacheCategoryStore) UpdateCategoryCache(categories model.Categories, allowFromCache bool) {
	if allowFromCache {
		for _, cate := range categories {
			l.addToCache(cate)
		}
		return
	}

	l.CategoryStore.UpdateCategoryCache(categories, allowFromCache)
}

func (l *LocalCacheCategoryStore) getFromCacheById(id string) (*model.Category, bool) {
	var cate *model.Category
	err := l.rootStore.doStandardReadCache(l.rootStore.categoryByIdsCache, id, &cate)
	if err == nil {
		return cate, true
	}
	return nil, false
}
