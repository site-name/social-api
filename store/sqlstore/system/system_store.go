package system

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/sitename/sitename/model"
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
	if err := s.GetMaster().Create(system).Error; err != nil {
		return errors.Wrapf(err, "failed to save system property with name=%s", system.Name)
	}
	return nil
}

func (s *SqlSystemStore) SaveOrUpdate(system *model.System) error {
	query, args, err := s.GetQueryBuilder().
		Insert("Systems").
		Columns("Name", "Value").
		Values(system.Name, system.Value).
		SuffixExpr(squirrel.Expr("ON CONFLICT (name) DO UPDATE SET Value = ?", system.Value)).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "SaveOrUpdate_ToSql")
	}

	err = s.GetMaster().Raw(query, args...).Error
	if err != nil {
		return errors.Wrap(err, "failed to upsert system property")
	}
	return nil
}

func (s *SqlSystemStore) SaveOrUpdateWithWarnMetricHandling(system *model.System) error {
	if err := s.SaveOrUpdate(system); err != nil {
		return err
	}

	if strings.HasPrefix(system.Name, model.WarnMetricStatusStorePrefix) &&
		(system.Value == model.WarnMetricStatusRunonce || system.Value == model.WarnMetricStatusLimitReached) {
		if err := s.SaveOrUpdate(&model.System{
			Name:  model.SystemWarnMetricLastRunTimestampKey,
			Value: strconv.FormatInt(util.MillisFromTime(time.Now()), 10),
		}); err != nil {
			return errors.Wrapf(err, "failed to save system property with name=%s", model.SystemWarnMetricLastRunTimestampKey)
		}
	}

	return nil
}

func (s *SqlSystemStore) Update(system *model.System) error {
	if err := s.GetMaster().Raw("UPDATE Systems SET Value=? WHERE Name=?", system.Value, system.Name).Error; err != nil {
		return errors.Wrapf(err, "failed to update system property with name=%s", system.Name)
	}
	return nil
}

func (s *SqlSystemStore) Get() (model.StringMap, error) {
	var systems []model.System
	props := make(model.StringMap)
	if err := s.GetReplica().Find(&systems).Error; err != nil {
		return nil, errors.Wrap(err, "failed to system properties")
	}
	for _, prop := range systems {
		props[prop.Name] = prop.Value
	}

	return props, nil
}

func (s *SqlSystemStore) GetByName(name string) (*model.System, error) {
	var system model.System
	if err := s.GetMaster().First(&system, "Name = ?", name).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound("System", fmt.Sprintf("name=%s", system.Name))
		}
		return nil, errors.Wrapf(err, "failed to get system property with name=%s", system.Name)
	}

	return &system, nil
}

func (s *SqlSystemStore) PermanentDeleteByName(name string) (*model.System, error) {
	if err := s.GetMaster().Raw("DELETE FROM Systems WHERE Name = ?", name).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to permanent delete system property with name=%s", name)
	}

	return &model.System{Name: name}, nil
}

// InsertIfExists inserts a given system value if it does not already exist. If a value
// already exists, it returns the old one, else returns the new one.
func (s *SqlSystemStore) InsertIfExists(system *model.System) (*model.System, error) {
	tx := s.GetMaster().Begin(&sql.TxOptions{
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
