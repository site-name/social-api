package account

import (
	"database/sql"

	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
)

type SqlStatusStore struct {
	store.Store
}

func NewSqlStatusStore(sqlStore store.Store) store.StatusStore {
	return &SqlStatusStore{sqlStore}
}

func (s SqlStatusStore) Upsert(status model.Status) (*model.Status, error) {
	if err := model_helper.StatusIsValid(status); err != nil {
		return nil, err
	}

	err := status.Upsert(
		s.GetMaster(),
		true,
		[]string{model.StatusColumns.UserID},
		boil.Blacklist(model.StatusColumns.UserID),
		boil.Infer(),
	)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func (s *SqlStatusStore) Get(userId string) (*model.Status, error) {
	status, err := model.FindStatus(s.GetReplica(), userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.Status, userId)
		}
		return nil, err
	}
	return status, nil
}

func (s *SqlStatusStore) GetByIds(userIds []string) (model.StatusSlice, error) {
	return model.Statuses(model.StatusWhere.UserID.IN(userIds)).All(s.GetReplica())
}

func (s *SqlStatusStore) ResetAll() error {
	_, err := model.Statuses(model.StatusWhere.Manual.EQ(false)).UpdateAll(s.GetMaster(), model.M{
		model.StatusColumns.Status: model_helper.STATUS_OFFLINE,
	})
	return err
}

func (s *SqlStatusStore) GetTotalActiveUsersCount() (int64, error) {
	var lastActivityTime = model_helper.GetMillis() - (1000 * 60 * 60 * 24)

	return model.Statuses(model.StatusWhere.LastActivityAt.GT(lastActivityTime)).Count(s.GetReplica())
}

func (s *SqlStatusStore) UpdateLastActivityAt(userId string, lastActivityAt int64) error {
	_, err := model.Statuses(model.StatusWhere.UserID.EQ(userId)).UpdateAll(s.GetMaster(), model.M{
		model.StatusColumns.LastActivityAt: lastActivityAt,
	})
	return err
}
