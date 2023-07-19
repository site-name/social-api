package account

import (
	"errors"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlCustomerNoteStore struct {
	store.Store
}

var customerNoteModelFields = util.AnyArray[string]{
	"Id",
	"UserID",
	"Date",
	"Content",
	"IsPublic",
	"CustomerID",
}

func NewSqlCustomerNoteStore(s store.Store) store.CustomerNoteStore {
	return &SqlCustomerNoteStore{s}
}

func (cs *SqlCustomerNoteStore) ModelFields(prefix string) util.AnyArray[string] {
	if prefix == "" {
		return customerNoteModelFields
	}

	return customerNoteModelFields.Map(func(_ int, item string) string {
		return prefix + item
	})
}

func (cs *SqlCustomerNoteStore) Save(note *model.CustomerNote) (*model.CustomerNote, error) {
	err := cs.GetMaster().Create(note).Error
	if err != nil {
		return nil, err
	}
	return note, err
}

func (cs *SqlCustomerNoteStore) Get(id string) (*model.CustomerNote, error) {
	var res model.CustomerNote
	err := cs.GetReplica().First(&res, "Id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.CustomerNoteTableName, id)
		}
		return nil, err
	}
	return &res, nil
}
