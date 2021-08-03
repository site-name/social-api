package account

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/store"
)

type SqlStaffNotificationRecipientStore struct {
	store.Store
}

func NewSqlStaffNotificationRecipientStore(s store.Store) store.StaffNotificationRecipientStore {
	ss := &SqlStaffNotificationRecipientStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(account.StaffNotificationRecipient{}, store.StaffNotificationRecipientTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("UserID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("StaffEmail").SetMaxSize(model.USER_EMAIL_MAX_LENGTH)
	}

	return ss
}

func (ss *SqlStaffNotificationRecipientStore) CreateIndexesIfNotExists() {
	ss.CreateForeignKeyIfNotExists(store.StaffNotificationRecipientTableName, "UserID", store.UserTableName, "Id", true)
}

func (ss *SqlStaffNotificationRecipientStore) Save(record *account.StaffNotificationRecipient) (*account.StaffNotificationRecipient, error) {
	record.PreSave()
	if err := record.IsValid(); err != nil {
		return nil, err
	}

	if err := ss.GetMaster().Insert(record); err != nil {
		return nil, errors.Wrapf(err, "failed to save StaffNotificationRecipient with Id=%s", record.Id)
	}

	return record, nil
}

func (ss *SqlStaffNotificationRecipientStore) Get(id string) (*account.StaffNotificationRecipient, error) {
	result, err := ss.GetReplica().Get(account.StaffNotificationRecipient{}, id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find StaffNotificationRecipient with Id=%s", id)
	}

	return result.(*account.StaffNotificationRecipient), nil
}
