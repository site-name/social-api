package cluster

import (
	"github.com/pkg/errors"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type sqlClusterDiscoveryStore struct {
	store.Store
}

func NewSqlClusterDiscoveryStore(sqlStore store.Store) store.ClusterDiscoveryStore {
	return &sqlClusterDiscoveryStore{sqlStore}
}

func (s sqlClusterDiscoveryStore) Save(ClusterDiscovery *model.ClusterDiscovery) error {
	if err := s.GetMaster().Create(ClusterDiscovery).Error; err != nil {
		return errors.Wrap(err, "failed to save ClusterDiscovery")
	}
	return nil
}

func (s sqlClusterDiscoveryStore) Delete(ClusterDiscovery *model.ClusterDiscovery) (bool, error) {
	res := s.GetMaster().
		Raw("DELETE FROM "+model.ClusterDiscoveryTableName+" WHERE Type = ? AND ClusterName = ? AND Hostname = ?", ClusterDiscovery.Type, ClusterDiscovery.ClusterName, ClusterDiscovery.Hostname)
	if res.Error != nil {
		return false, errors.Wrap(res.Error, "failed to delete ClusterDiscovery")
	}
	return res.RowsAffected != 0, nil
}

func (s sqlClusterDiscoveryStore) Exists(ClusterDiscovery *model.ClusterDiscovery) (bool, error) {
	var count int64
	err := s.GetMaster().Raw("SELECT COUNT(*) FROM "+model.ClusterDiscoveryTableName+" WHERE Type = ? AND ClusterName = ? AND Hostname = ?", ClusterDiscovery.Type, ClusterDiscovery.ClusterName, ClusterDiscovery.Hostname).Scan(&count).Error
	if err != nil {
		return false, errors.Wrap(err, "failed to count ClusterDiscovery")
	}
	return count != 0, nil
}

func (s sqlClusterDiscoveryStore) GetAll(ClusterDiscoveryType, clusterName string) ([]*model.ClusterDiscovery, error) {
	var list []*model.ClusterDiscovery
	err := s.GetMaster().
		Find(&list, "Type = ? AND ClusterName = ? AND ListPingAt > ?", ClusterDiscoveryType, clusterName, model.GetMillis()-model.CDS_OFFLINE_AFTER_MILLIS).Error
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find ClusterDiscovery")
	}
	return list, nil
}

func (s sqlClusterDiscoveryStore) SetLastPingAt(ClusterDiscovery *model.ClusterDiscovery) error {
	err := s.GetMaster().Raw("UPDATE "+model.ClusterDiscoveryTableName+" SET LastPingAt = ? WHERE Type = ? AND ClusterName = ? AND Hostname = ?", model.GetMillis(), ClusterDiscovery.Type, ClusterDiscovery.ClusterName, ClusterDiscovery.Hostname).Error
	if err != nil {
		return errors.Wrap(err, "failed to update ClusterDiscovery")
	}
	return nil
}

func (s sqlClusterDiscoveryStore) Cleanup() error {
	err := s.GetMaster().Raw("DELETE FROM "+model.ClusterDiscoveryTableName+" WHERE LastPingAt < ?", model.GetMillis()-model.CDS_OFFLINE_AFTER_MILLIS).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete ClusterDiscoveries")
	}
	return nil
}
