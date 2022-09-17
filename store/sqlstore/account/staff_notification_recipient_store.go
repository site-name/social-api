package account

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlStaffNotificationRecipientStore struct {
	store.Store
}

var staffNotificationRecipientModelFields = model.AnyArray[string]{
	"Id",
	"UserID",
	"StaffEmail",
	"Active",
}

func (ss *SqlStaffNotificationRecipientStore) ModelFields(prefix string) model.AnyArray[string] {
	if prefix == "" {
		return staffNotificationRecipientModelFields
	}

	return staffNotificationRecipientModelFields.Map(func(_ int, item string) string {
		return prefix + item
	})
}

func NewSqlStaffNotificationRecipientStore(s store.Store) store.StaffNotificationRecipientStore {
	return &SqlStaffNotificationRecipientStore{s}
}

func (ss *SqlStaffNotificationRecipientStore) Save(record *model.StaffNotificationRecipient) (*model.StaffNotificationRecipient, error) {
	record.PreSave()
	if err := record.IsValid(); err != nil {
		return nil, err
	}
	query := "INSERT INTO " + store.StaffNotificationRecipientTableName + " (" + ss.ModelFields("").Join(",") + ") VALUES (" + ss.ModelFields(":").Join(",") + ")"
	if _, err := ss.GetMasterX().NamedExec(query, record); err != nil {
		return nil, errors.Wrapf(err, "failed to save StaffNotificationRecipient with Id=%s", record.Id)
	}

	return record, nil
}

func (ss *SqlStaffNotificationRecipientStore) Get(id string) (*model.StaffNotificationRecipient, error) {
	var res model.StaffNotificationRecipient

	err := ss.GetReplicaX().Get(&res, "SELECT * FROM "+store.StaffNotificationRecipientTableName+" WHERE Id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.StaffNotificationRecipientTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find StaffNotificationRecipient with Id=%s", id)
	}

	return &res, nil
}
