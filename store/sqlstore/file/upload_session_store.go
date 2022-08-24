package file

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/file"
	"github.com/sitename/sitename/store"
)

type SqlUploadSessionStore struct {
	store.Store
}

func NewSqlUploadSessionStore(sqlStore store.Store) store.UploadSessionStore {
	return &SqlUploadSessionStore{sqlStore}
}

func (s *SqlUploadSessionStore) ModelFields(prefix string) model.StringArray {
	res := model.StringArray{
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

func (us *SqlUploadSessionStore) Save(session *file.UploadSession) (*file.UploadSession, error) {
	session.PreSave()
	if err := session.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.UploadSessionTableName + "(" + us.ModelFields("").Join(",") + ") VALUES (" + us.ModelFields(":").Join(",") + ")"
	if _, err := us.GetMasterX().NamedExec(query, session); err != nil {
		return nil, errors.Wrap(err, "SqlUploadSessionStore.Save: failed to insert")
	}
	return session, nil
}

func (us *SqlUploadSessionStore) Update(session *file.UploadSession) error {
	if err := session.IsValid(); err != nil {
		return err
	}

	query := "UPDATE " + store.UploadSessionTableName + " SET " + us.
		ModelFields("").
		Map(func(_ int, s string) string {
			return s + "=:" + s
		}).
		Join(",") + " WHERE Id=:Id"

	if _, err := us.GetMasterX().NamedExec(query, session); err != nil {
		return errors.Wrapf(err, "SqlUploadSessionStore.Update: failed to update session with id=%s", session.Id)
	}
	return nil
}

func (us SqlUploadSessionStore) Get(id string) (*file.UploadSession, error) {
	var session *file.UploadSession
	if err := us.GetReplicaX().Get(&session, "SELECT * FROM "+store.UploadSessionTableName+" WHERE Id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.UploadSessionTableName, id)
		}
		return nil, errors.Wrapf(err, "SqlUploadSessionStore.Get: failed to select session with id=%s", id)
	}
	return session, nil
}

func (us *SqlUploadSessionStore) GetForUser(userId string) ([]*file.UploadSession, error) {
	var sessions []*file.UploadSession

	if err := us.GetReplicaX().Select(
		&sessions,
		"SELECT * FROM "+store.UploadSessionTableName+" WHERE UserId = ? ORDER BY CreateAt ASC",
		userId,
	); err != nil {
		return nil, errors.Wrap(err, "failed to find upload session for user id="+userId)
	}
	return sessions, nil
}

func (us *SqlUploadSessionStore) Delete(id string) error {
	if _, err := us.GetMasterX().Exec("DELETE FROM "+store.UploadSessionTableName+" WHERE Id = ?", id); err != nil {
		return errors.Wrap(err, "failed to delete upload session with id="+id)
	}

	return nil
}
