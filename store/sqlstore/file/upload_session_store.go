package file

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/file"
	"github.com/sitename/sitename/store"
)

const (
	UploadSessionTableName = "UploadSessions"
)

type SqlUploadSessionStore struct {
	store.Store
}

func NewSqlUploadSessionStore(sqlStore store.Store) store.UploadSessionStore {
	s := &SqlUploadSessionStore{sqlStore}
	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(file.UploadSession{}, UploadSessionTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Type").SetMaxSize(32)
		table.ColMap("UserID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("FileName").SetMaxSize(256)
		table.ColMap("Path").SetMaxSize(512)
	}
	return s
}

func (us SqlUploadSessionStore) CreateIndexesIfNotExists() {
	us.CreateIndexIfNotExists("idx_uploadsessions_user_id", UploadSessionTableName, "Type")
	us.CreateIndexIfNotExists("idx_uploadsessions_create_at", UploadSessionTableName, "CreateAt")
	us.CreateIndexIfNotExists("idx_uploadsessions_user_id", UploadSessionTableName, "UserID")
}

func (us *SqlUploadSessionStore) Save(session *file.UploadSession) (*file.UploadSession, error) {
	session.PreSave()
	if err := session.IsValid(); err != nil {
		return nil, err
	}
	if err := us.GetMaster().Insert(session); err != nil {
		return nil, errors.Wrap(err, "SqlUploadSessionStore.Save: failed to insert")
	}
	return session, nil
}

func (us *SqlUploadSessionStore) Update(session *file.UploadSession) error {
	if err := session.IsValid(); err != nil {
		return errors.Wrap(err, "SqlUploadSessionStore.Update: validation failed")
	}
	if _, err := us.GetMaster().Update(session); err != nil {
		if err == sql.ErrNoRows {
			return store.NewErrNotFound("UploadSession", session.Id)
		}
		return errors.Wrapf(err, "SqlUploadSessionStore.Update: failed to update session with id=%s", session.Id)
	}
	return nil
}

func (us SqlUploadSessionStore) Get(id string) (*file.UploadSession, error) {
	var session *file.UploadSession
	if err := us.GetReplica().SelectOne(&session, "SELECT * FROM "+UploadSessionTableName+" WHERE Id = :Id", map[string]interface{}{"Id": id}); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("UploadSession", id)
		}
		return nil, errors.Wrapf(err, "SqlUploadSessionStore.Get: failed to select session with id=%s", id)
	}
	return session, nil
}

func (us *SqlUploadSessionStore) GetForUser(userId string) ([]*file.UploadSession, error) {
	var sessions []*file.UploadSession

	if _, err := us.GetReplica().Select(
		&sessions,
		"SELECT * FROM "+UploadSessionTableName+" WHERE UserId = :UserId ORDER BY CreateAt ASC",
		map[string]interface{}{"UserId": userId},
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(UploadSessionTableName, "UserId="+userId)
		}
		return nil, errors.Wrap(err, "failed to find upload session for user id="+userId)
	}
	return sessions, nil
}

func (us *SqlUploadSessionStore) Delete(id string) error {

	if _, err := us.GetMaster().Exec("DELETE FROM "+UploadSessionTableName+" WHERE Id = :Id", map[string]interface{}{"Id": id}); err != nil {
		return errors.Wrap(err, "failed to delete upload session with id="+id)
	}

	return nil
}
