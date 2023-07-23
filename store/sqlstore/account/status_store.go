package account

import (
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlStatusStore struct {
	store.Store
}

func NewSqlStatusStore(sqlStore store.Store) store.StatusStore {
	return &SqlStatusStore{sqlStore}
}

func (s SqlStatusStore) SaveOrUpdate(status *model.Status) error {
	err := s.GetMaster().Save(status).Error
	if err != nil {
		return errors.Wrap(err, "failed to upsert status")
	}
	return nil
}

func (s *SqlStatusStore) Get(userId string) (*model.Status, error) {
	var status model.Status

	if err := s.GetReplica().First(&status, `UserId = ?`, userId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound("Status", fmt.Sprintf("userId=%s", userId))
		}
		return nil, errors.Wrapf(err, "failed to get Status with userId=%s", userId)
	}
	return &status, nil
}

func (s *SqlStatusStore) GetByIds(userIds []string) ([]*model.Status, error) {
	var statuses []*model.Status
	err := s.GetReplica().Find("UserId IN ?", userIds).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find Statuses")
	}

	return statuses, nil
}

func (s *SqlStatusStore) ResetAll() error {
	if err := s.GetMaster().Raw("UPDATE Status SET Status = ? WHERE Manual = false", model.STATUS_OFFLINE).Error; err != nil {
		return errors.Wrap(err, "failed to update Statuses")
	}
	return nil
}

func (s *SqlStatusStore) GetTotalActiveUsersCount() (int64, error) {

	var (
		time  = model.GetMillis() - (1000 * 60 * 60 * 24)
		count int64
	)
	err := s.GetReplica().Raw("SELECT COUNT(UserId) FROM Status WHERE LastActivityAt > ?", time).Scan(&count).Error
	if err != nil {
		return count, errors.Wrap(err, "failed to count active users")
	}
	return count, nil
}

func (s *SqlStatusStore) UpdateLastActivityAt(userId string, lastActivityAt int64) error {
	if err := s.GetMaster().Raw("UPDATE Status SET LastActivityAt = ? WHERE UserId = ?", lastActivityAt, userId).Error; err != nil {
		return errors.Wrapf(err, "failed to update last activity for userId=%s", userId)
	}

	return nil
}
