package system

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"gorm.io/gorm"

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

func (s *SqlSystemStore) Save(system *model.System) error {
	return system.Insert(s.Context(), s.GetMaster(), boil.Infer())
}

func (s *SqlSystemStore) SaveOrUpdate(system *model.System) error {
	return system.Upsert(s.Context(), s.GetMaster(), true, []string{model.SystemColumns.Name}, boil.Infer(), boil.Infer())
}

func (s *SqlSystemStore) SaveOrUpdateWithWarnMetricHandling(system *model.System) error {
	if err := s.SaveOrUpdate(system); err != nil {
		return err
	}

	if strings.HasPrefix(system.Name, model_helper.WarnMetricStatusStorePrefix) &&
		(system.Value == model_helper.WarnMetricStatusRunonce || system.Value == model_helper.WarnMetricStatusLimitReached) {
		if err := s.SaveOrUpdate(&model.System{
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
		UpdateAll(s.Context(), s.GetMaster(), model.M{
			model.SystemColumns.Value: system.Value,
		})
	return err
}

func (s *SqlSystemStore) Get() (map[string]string, error) {
	systems, err := model.Systems().All(s.Context(), s.GetReplica())
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
	system, err := model.Systems(model.SystemWhere.Name.EQ(name)).One(s.Context(), s.GetReplica())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.NewErrNotFound("System", fmt.Sprintf("name=%s", system.Name))
		}
		return nil, err
	}
	return system, nil
}

func (s *SqlSystemStore) PermanentDeleteByName(name string) (*model.System, error) {
	_, err := (&model.System{Name: name}).Delete(s.Context(), s.GetMaster())
	if err != nil {
		return nil, err
	}
	return &model.System{Name: name}, nil
}

// InsertIfExists inserts a given system value if it does not already exist. If a value
// already exists, it returns the old one, else returns the new one.
func (s *SqlSystemStore) InsertIfExists(system *model.System) (*model.System, error) {
	tx := s.GetMaster().BeginTx(s.Context(), &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	defer s.FinalizeTransaction(tx)

	var origSystem model.System
	if err := tx.First(&origSystem, `Name = ?`, system.Name).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.Wrapf(err, "failed to get system property with name=%s", system.Name)
	}

	if origSystem.Value != "" {
		// Already a value exists, return that.
		return &origSystem, nil
	}

	// Key does not exist, need to insert.
	if err := tx.Create(system).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to save system property with name=%s", system.Name)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.Wrap(err, "commit_transaction")
	}
	return system, nil
}
