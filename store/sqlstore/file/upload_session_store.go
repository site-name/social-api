package file

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlUploadSessionStore struct {
	store.Store
}

func NewSqlUploadSessionStore(sqlStore store.Store) store.UploadSessionStore {
	return &SqlUploadSessionStore{sqlStore}
}

func (s *SqlUploadSessionStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"Type",
		"CreateAt",
		"UserID",
		"FileName",
		"Path",
		"FileSize",
		"FileOffset",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (us *SqlUploadSessionStore) Save(session *model.UploadSession) (*model.UploadSession, error) {
	err := us.GetMaster().Create(session).Error
	if err != nil {
		return nil, errors.Wrap(err, "SqlUploadSessionStore.Save: failed to insert")
	}
	return session, nil
}

func (us *SqlUploadSessionStore) Update(session *model.UploadSession) error {
	err := us.GetMaster().Updates(session).Error
	if err != nil {
		return errors.Wrapf(err, "SqlUploadSessionStore.Update: failed to update session with id=%s", session.Id)
	}
	return nil
}

func (us SqlUploadSessionStore) Get(id string) (*model.UploadSession, error) {
	var session model.UploadSession
	if err := us.GetReplica().First(&session, "Id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.UploadSessionTableName, id)
		}
		return nil, errors.Wrapf(err, "SqlUploadSessionStore.Get: failed to select session with id=%s", id)
	}
	return &session, nil
}

func (us *SqlUploadSessionStore) GetForUser(userId string) ([]*model.UploadSession, error) {
	var sessions []*model.UploadSession
	err := us.GetReplica().Order("CreateAt ASC").Find(&sessions, "UserID = ?", userId).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find upload session for user id="+userId)
	}
	return sessions, nil
}

func (us *SqlUploadSessionStore) Delete(id string) error {
	if err := us.GetMaster().Raw("DELETE FROM "+model.UploadSessionTableName+" WHERE Id = ?", id).Error; err != nil {
		return errors.Wrap(err, "failed to delete upload session with id="+id)
	}

	return nil
}
