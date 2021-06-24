package cluster

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/cluster"
	"github.com/sitename/sitename/store"
)

type sqlClusterDiscoveryStore struct {
	store.Store
}

func NewSqlClusterDiscoveryStore(sqlStore store.Store) store.ClusterDiscoveryStore {
	s := &sqlClusterDiscoveryStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(cluster.ClusterDiscovery{}, "ClusterDiscovery").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Type").SetMaxSize(64)
		table.ColMap("ClusterName").SetMaxSize(64)
		table.ColMap("Hostname").SetMaxSize(512)
	}

	return s
}

func (*sqlClusterDiscoveryStore) CreateIndexesIfNotExists() {}

func (s sqlClusterDiscoveryStore) Save(ClusterDiscovery *cluster.ClusterDiscovery) error {
	ClusterDiscovery.PreSave()
	if err := ClusterDiscovery.IsValid(); err != nil {
		return err
	}

	if err := s.GetMaster().Insert(ClusterDiscovery); err != nil {
		return errors.Wrap(err, "failed to save ClusterDiscovery")
	}
	return nil
}

func (s sqlClusterDiscoveryStore) Delete(ClusterDiscovery *cluster.ClusterDiscovery) (bool, error) {
	query := s.GetQueryBuilder().
		Delete("ClusterDiscovery").
		Where(sq.Eq{"Type": ClusterDiscovery.Type}).
		Where(sq.Eq{"ClusterName": ClusterDiscovery.ClusterName}).
		Where(sq.Eq{"Hostname": ClusterDiscovery.Hostname})

	queryString, args, err := query.ToSql()
	if err != nil {
		return false, errors.Wrap(err, "cluster_discovery_tosql")
	}

	count, err := s.GetMaster().SelectInt(queryString, args...)
	if err != nil {
		return false, errors.Wrap(err, "failed to delete ClusterDiscovery")
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (s sqlClusterDiscoveryStore) Exists(ClusterDiscovery *cluster.ClusterDiscovery) (bool, error) {
	query := s.GetQueryBuilder().
		Select("COUNT(*)").
		From("ClusterDiscovery").
		Where(sq.Eq{"Type": ClusterDiscovery.Type}).
		Where(sq.Eq{"ClusterName": ClusterDiscovery.ClusterName}).
		Where(sq.Eq{"Hostname": ClusterDiscovery.Hostname})

	queryString, args, err := query.ToSql()
	if err != nil {
		return false, errors.Wrap(err, "cluster_discovery_tosql")
	}

	count, err := s.GetMaster().SelectInt(queryString, args...)
	if err != nil {
		return false, errors.Wrap(err, "failed to count ClusterDiscovery")
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (s sqlClusterDiscoveryStore) GetAll(ClusterDiscoveryType, clusterName string) ([]*cluster.ClusterDiscovery, error) {
	query := s.GetQueryBuilder().
		Select("*").
		From("ClusterDiscovery").
		Where(sq.Eq{"Type": ClusterDiscoveryType}).
		Where(sq.Eq{"ClusterName": clusterName}).
		Where(sq.Gt{"LastPingAt": model.GetMillis() - cluster.CDS_OFFLINE_AFTER_MILLIS})

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "cluster_discovery_tosql")
	}

	var list []*cluster.ClusterDiscovery
	if _, err := s.GetMaster().Select(&list, queryString, args...); err != nil {
		return nil, errors.Wrapf(err, "failed to find ClusterDiscovery")
	}
	return list, nil
}

func (s sqlClusterDiscoveryStore) SetLastPingAt(ClusterDiscovery *cluster.ClusterDiscovery) error {
	query := s.GetQueryBuilder().
		Update("ClusterDiscovery").
		Set("LastPingAt", model.GetMillis()).
		Where(sq.Eq{"Type": ClusterDiscovery.Type}).
		Where(sq.Eq{"ClusterName": ClusterDiscovery.ClusterName}).
		Where(sq.Eq{"Hostname": ClusterDiscovery.Hostname})

	queryString, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "cluster_discovery_tosql")
	}

	if _, err := s.GetMaster().Exec(queryString, args...); err != nil {
		return errors.Wrap(err, "failed to update ClusterDiscovery")
	}
	return nil
}

func (s sqlClusterDiscoveryStore) Cleanup() error {
	query := s.GetQueryBuilder().
		Delete("ClusterDiscovery").
		Where(sq.Lt{"LastPingAt": model.GetMillis() - cluster.CDS_OFFLINE_AFTER_MILLIS})

	queryString, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "cluster_discovery_tosql")
	}

	if _, err := s.GetMaster().Exec(queryString, args...); err != nil {
		return errors.Wrap(err, "failed to delete ClusterDiscoveries")
	}
	return nil
}
