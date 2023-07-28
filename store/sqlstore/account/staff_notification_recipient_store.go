package account

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlStaffNotificationRecipientStore struct {
	store.Store
}

func NewSqlStaffNotificationRecipientStore(s store.Store) store.StaffNotificationRecipientStore {
	return &SqlStaffNotificationRecipientStore{s}
}

func (ss *SqlStaffNotificationRecipientStore) Save(record *model.StaffNotificationRecipient) (*model.StaffNotificationRecipient, error) {
	err := ss.GetMaster().Create(record).Error
	if err != nil {
		if ss.IsUniqueConstraintError(err, []string{"staff_notification_recipients_staff_email_key", "StaffEmail"}) {
			return nil, store.NewErrInvalidInput(model.StaffNotificationRecipientTableName, "StaffEmail", record.StaffEmail)
		}
		return nil, errors.Wrap(err, "failed to create the record")
	}
	return record, nil
}

func (s *SqlStaffNotificationRecipientStore) FilterByOptions(options *model.StaffNotificationRecipientFilterOptions) ([]*model.StaffNotificationRecipient, error) {
	var res []*model.StaffNotificationRecipient
	err := s.GetReplica().Find(&res, store.BuildSqlizer(options.Conditions)...).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
