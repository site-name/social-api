package system

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlSystemStore struct {
	store.Store
}

func NewSqlSystemStore(sqlStore store.Store) store.SystemStore {
	return &SqlSystemStore{sqlStore}
}

func (s *SqlSystemStore) Save(system model.System) error {
	return system.Insert(s.GetMaster(), boil.Infer())
}

func (s *SqlSystemStore) SaveOrUpdate(system model.System) error {
	return system.Upsert(s.GetMaster(), true, []string{model.SystemColumns.Name}, boil.Infer(), boil.Infer())
}

func (s *SqlSystemStore) SaveOrUpdateWithWarnMetricHandling(system model.System) error {
	if err := s.SaveOrUpdate(system); err != nil {
		return err
	}

	if strings.HasPrefix(system.Name, model_helper.WarnMetricStatusStorePrefix) &&
		(system.Value == model_helper.WarnMetricStatusRunonce || system.Value == model_helper.WarnMetricStatusLimitReached) {
		if err := s.SaveOrUpdate(model.System{
			Name:  model_helper.SystemWarnMetricLastRunTimestampKey,
			Value: strconv.FormatInt(util.MillisFromTime(time.Now()), 10),
		}); err != nil {
			return errors.Wrapf(err, "failed to save system property with name=%s", model_helper.SystemWarnMetricLastRunTimestampKey)
		}
	}

	return nil
}

func (s *SqlSystemStore) Update(system model.System) error {
	_, err := model.
		Systems(model.SystemWhere.Name.EQ(system.Name)).
		UpdateAll(s.GetMaster(), model.M{
			model.SystemColumns.Value: system.Value,
		})
	return err
}

func (s *SqlSystemStore) Get() (map[string]string, error) {
	systems, err := model.Systems().All(s.GetReplica())
	if err != nil {
		return nil, err
	}

	res := map[string]string{}
	for i := range systems {
		system := systems[i]
		res[system.Name] = system.Value
	}
	return res, nil
}

func (s *SqlSystemStore) GetByName(name string) (*model.System, error) {
	system, err := model.Systems(model.SystemWhere.Name.EQ(name)).One(s.GetReplica())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.NewErrNotFound("System", fmt.Sprintf("name=%s", system.Name))
		}
		return nil, err
	}
	return system, nil
}

func (s *SqlSystemStore) PermanentDeleteByName(name string) (*model.System, error) {
	_, err := (&model.System{Name: name}).Delete(s.GetMaster())
	if err != nil {
		return nil, err
	}
	return &model.System{Name: name}, nil
}

func (s *SqlSystemStore) InsertIfExists(system model.System) (*model.System, error) {
	tx, err := s.GetMaster().BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return nil, errors.Wrap(err, "begin_transaction")
	}
	defer s.FinalizeTransaction(tx)

	origSystem, err := model.FindSystem(tx, system.Name)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if origSystem.Value != "" {
		// Already a value exists, return that.
		return origSystem, nil
	}

	err = system.Insert(tx, boil.Infer())
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "commit_transaction")
	}
	return &system, nil
}
