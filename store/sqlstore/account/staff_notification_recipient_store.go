package account

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlStaffNotificationRecipientStore struct {
	store.Store
}

func NewSqlStaffNotificationRecipientStore(s store.Store) store.StaffNotificationRecipientStore {
	return &SqlStaffNotificationRecipientStore{s}
}

func (ss *SqlStaffNotificationRecipientStore) Save(record model.StaffNotificationRecipient) (*model.StaffNotificationRecipient, error) {
	if err := model_helper.StaffNotificationRecipientIsValid(record); err != nil {
		return nil, err
	}
	err := record.Insert(ss.GetMaster(), boil.Infer())
	if err != nil {
		if ss.IsUniqueConstraintError(err, []string{"staff_email_unique"}) {
			return nil, store.NewErrInvalidInput(model.TableNames.StaffNotificationRecipients, "StaffEmail", record.StaffEmail)
		}
		return nil, err
	}
	return &record, nil
}

func (s *SqlStaffNotificationRecipientStore) FilterByOptions(options model_helper.StaffNotificationRecipientFilterOptions) (model.StaffNotificationRecipientSlice, error) {
	return model.StaffNotificationRecipients(options.Conditions...).All(s.GetReplica())
}
