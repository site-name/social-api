package cluster

import (
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
)

type sqlClusterDiscoveryStore struct {
	store.Store
}

func NewSqlClusterDiscoveryStore(sqlStore store.Store) store.ClusterDiscoveryStore {
	return &sqlClusterDiscoveryStore{sqlStore}
}

func (s sqlClusterDiscoveryStore) Save(ClusterDiscovery model.ClusterDiscovery) error {
	model_helper.ClusterDiscoveryPreSave(&ClusterDiscovery)
	if err := model_helper.ClusterDiscoveryIsValid(ClusterDiscovery); err != nil {
		return err
	}

	return ClusterDiscovery.Insert(s.GetMaster(), boil.Infer())
}

func (s sqlClusterDiscoveryStore) Delete(ClusterDiscovery model.ClusterDiscovery) (bool, error) {
	numDeleted, err := model.ClusterDiscoveries(
		model.ClusterDiscoveryWhere.Type.EQ(ClusterDiscovery.Type),
		model.ClusterDiscoveryWhere.ClusterName.EQ(ClusterDiscovery.ClusterName),
		model.ClusterDiscoveryWhere.HostName.EQ(ClusterDiscovery.HostName),
	).DeleteAll(s.GetMaster())
	return numDeleted > 0, err
}

func (s sqlClusterDiscoveryStore) Exists(ClusterDiscovery model.ClusterDiscovery) (bool, error) {
	return model.ClusterDiscoveries(
		model.ClusterDiscoveryWhere.Type.EQ(ClusterDiscovery.Type),
		model.ClusterDiscoveryWhere.ClusterName.EQ(ClusterDiscovery.ClusterName),
		model.ClusterDiscoveryWhere.HostName.EQ(ClusterDiscovery.HostName),
	).Exists(s.GetReplica())
}

func (s sqlClusterDiscoveryStore) GetAll(ClusterDiscoveryType, clusterName string) (model.ClusterDiscoverySlice, error) {
	return model.ClusterDiscoveries(
		model.ClusterDiscoveryWhere.Type.EQ(ClusterDiscoveryType),
		model.ClusterDiscoveryWhere.ClusterName.EQ(clusterName),
		model.ClusterDiscoveryWhere.LastPingAt.GT(model_helper.GetMillis()-model_helper.CDS_OFFLINE_AFTER_MILLIS),
	).All(s.GetReplica())
}

func (s sqlClusterDiscoveryStore) SetLastPingAt(ClusterDiscovery model.ClusterDiscovery) error {
	_, err := model.ClusterDiscoveries(
		model.ClusterDiscoveryWhere.Type.EQ(ClusterDiscovery.Type),
		model.ClusterDiscoveryWhere.ClusterName.EQ(ClusterDiscovery.ClusterName),
		model.ClusterDiscoveryWhere.HostName.EQ(ClusterDiscovery.HostName),
	).UpdateAll(s.GetMaster(), model.M{model.ClusterDiscoveryColumns.LastPingAt: model_helper.GetMillis()})
	return err
}

func (s sqlClusterDiscoveryStore) Cleanup() error {
	_, err := model.ClusterDiscoveries(
		model.ClusterDiscoveryWhere.LastPingAt.LT(model_helper.GetMillis() - model_helper.CDS_OFFLINE_AFTER_MILLIS),
	).DeleteAll(s.GetMaster())
	return err
}
