package file

import (
	"database/sql"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlUploadSessionStore struct {
	store.Store
}

func NewSqlUploadSessionStore(sqlStore store.Store) store.UploadSessionStore {
	return &SqlUploadSessionStore{sqlStore}
}

func (us SqlUploadSessionStore) Upsert(session model.UploadSession) (*model.UploadSession, error) {
	isSaving := session.ID == ""
	if isSaving {
		model_helper.UploadSessionPreSave(&session)
	} else {
		model_helper.UploadSessionCommonPre(&session)
	}

	if err := model_helper.UploadSessionIsValid(session); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		err = session.Insert(us.GetMaster(), boil.Infer())
	} else {
		_, err = session.Update(us.GetMaster(), boil.Blacklist(model.UploadSessionColumns.CreatedAt))
	}
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (us SqlUploadSessionStore) Get(id string) (*model.UploadSession, error) {
	session, err := model.FindUploadSession(us.GetReplica(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.UploadSessions, id)
		}
		return nil, err
	}

	return session, nil
}

func (us SqlUploadSessionStore) FindAll(options model_helper.UploadSessionFilterOption) (model.UploadSessionSlice, error) {
	return model.UploadSessions(options.Conditions...).All(us.GetReplica())
}

func (us SqlUploadSessionStore) Delete(id string) error {
	_, err := model.UploadSessions(model.UploadSessionWhere.ID.EQ(id)).DeleteAll(us.GetMaster())
	return err
}
