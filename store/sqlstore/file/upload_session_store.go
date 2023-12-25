package file

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gorm.io/gorm"
)

type SqlUploadSessionStore struct {
	store.Store
}

func NewSqlUploadSessionStore(sqlStore store.Store) store.UploadSessionStore {
	return &SqlUploadSessionStore{sqlStore}
}

func (us SqlUploadSessionStore) Save(session model.UploadSession) (*model.UploadSession, error) {
	err := session.Insert(us.Context(), us.GetMaster(), boil.Infer())
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (us SqlUploadSessionStore) Update(session model.UploadSession) error {
	_, err := session.Update(us.Context(), us.GetMaster(), boil.Infer())
	return err
}

func (us SqlUploadSessionStore) Get(id string) (*model.UploadSession, error) {
	session, err := model.FindUploadSession(us.Context(), us.GetReplica(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.TableNames.UploadSessions, id)
		}
		return nil, err
	}

	return session, nil
}

func (us SqlUploadSessionStore) FindAll(mods ...qm.QueryMod) (model.UploadSessionSlice, error) {
	return model.UploadSessions(mods...).All(us.Context(), us.GetReplica())
}

func (us SqlUploadSessionStore) Delete(id string) error {
	session := model.UploadSession{ID: id}
	_, err := session.Delete(us.Context(), us.GetMaster())
	return err
}
