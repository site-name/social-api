package account

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/store"
)

type SqlStaffNotificationRecipientStore struct {
	store.Store
}

const (
	staffNotificationRecipientTableName = "StaffNotificationRecipients"
)

func NewSqlStaffNotificationRecipientStore(s store.Store) store.StaffNotificationRecipientStore {
	ss := &SqlStaffNotificationRecipientStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(account.StaffNotificationRecipient{}, staffNotificationRecipientTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("UserID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("StaffEmail").SetMaxSize(model.USER_EMAIL_MAX_LENGTH)
	}

	return ss
}

func (ss *SqlStaffNotificationRecipientStore) CreateIndexesIfNotExists() {}

func (ss *SqlStaffNotificationRecipientStore) Save(record *account.StaffNotificationRecipient) (*account.StaffNotificationRecipient, error) {
	record.PreSave()
	if err := ss.GetMaster().Insert(record); err != nil {
		return nil, errors.Wrapf(err, "failed to save StaffNotificationRecipient with Id=%s", record.Id)
	}

	return record, nil
}

func (ss *SqlStaffNotificationRecipientStore) Get(id string) (*account.StaffNotificationRecipient, error) {
	var record account.StaffNotificationRecipient
	err := ss.GetReplica().SelectOne(&record, "SELECT * FROM "+staffNotificationRecipientTableName+" WHERE Id = :ID", map[string]interface{}{"ID": id})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(staffNotificationRecipientTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find StaffNotificationRecipient with Id=%s", id)
	}

	return &record, nil
}