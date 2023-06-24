package account

import (
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
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

func (s *SqlStaffNotificationRecipientStore) FilterByOptions(options *model.StaffNotificationRecipientFilterOptions) ([]*model.StaffNotificationRecipient, error) {
	query := s.GetQueryBuilder().
		Select(s.ModelFields(store.StaffNotificationRecipientTableName + ".")...).
		From(store.StaffNotificationRecipientTableName)

	for _, opt := range []squirrel.Sqlizer{
		options.Id, options.Active, options.UserID, options.StaffEmail,
	} {
		if opt != nil {
			query = query.Where(opt)
		}
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var res []*model.StaffNotificationRecipient
	err = s.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find staff notification recipients")
	}

	return res, nil
}
