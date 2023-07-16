package account

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlStaffNotificationRecipientStore struct {
	store.Store
}

var staffNotificationRecipientModelFields = util.AnyArray[string]{
	"Id",
	"UserID",
	"StaffEmail",
	"Active",
}

func (ss *SqlStaffNotificationRecipientStore) ModelFields(prefix string) util.AnyArray[string] {
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
	err := ss.GetMaster().Create(record).Error
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (s *SqlStaffNotificationRecipientStore) FilterByOptions(options *model.StaffNotificationRecipientFilterOptions) ([]*model.StaffNotificationRecipient, error) {
	var res []*model.StaffNotificationRecipient
	err := s.GetReplica().Select(&res, store.BuildSqlizer(options.Conditions)...).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
