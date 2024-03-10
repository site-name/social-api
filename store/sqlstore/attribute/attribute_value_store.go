package attribute

import (
	"database/sql"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SqlAttributeValueStore struct {
	store.Store
}

func NewSqlAttributeValueStore(s store.Store) store.AttributeValueStore {
	return &SqlAttributeValueStore{s}
}

func (as *SqlAttributeValueStore) Upsert(transaction boil.ContextTransactor, values model.AttributeValueSlice) (model.AttributeValueSlice, error) {
	if transaction == nil {
		transaction = as.GetMaster()
	}

	for _, value := range values {
		if value == nil {
			continue
		}

		isSaving := value.ID == ""
		if isSaving {
			model_helper.AttributeValuePreSave(value)
		} else {
			model_helper.AttributeValueCommonPre(value)
		}

		if err := model_helper.AttributeValueIsValid(*value); err != nil {
			return nil, err
		}

		var err error
		if isSaving {
			err = value.Insert(transaction, boil.Infer())
		} else {
			_, err = value.Update(transaction, boil.Infer())
		}
		if err != nil {
			if as.IsUniqueConstraintError(err, []string{"attribute_values_slug_attribute_id_key", model.AttributeValueColumns.AttributeID, model.AttributeValueColumns.Slug}) {
				return nil, store.NewErrInvalidInput(model.TableNames.AttributeValues, model.AttributeValueColumns.Slug+"/"+model.AttributeValueColumns.AttributeID, "unique")
			}
			return nil, err
		}
	}

	return values, nil
}

func (as *SqlAttributeValueStore) Get(id string) (*model.AttributeValue, error) {
	av, err := model.FindAttributeValue(as.GetReplica(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.AttributeValues, id)
		}
		return nil, err
	}

	return av, nil
}

func (as *SqlAttributeValueStore) FilterByOptions(option model_helper.AttributeValueFilterOptions) (model.AttributeValueSlice, error) {
	conds := option.Conditions
	for _, load := range option.Preloads {
		conds = append(conds, qm.Load(load))
	}
	return model.AttributeValues(conds...).All(as.GetReplica())
}

func (as *SqlAttributeValueStore) Delete(tx boil.ContextTransactor, ids []string) (int64, error) {
	if tx == nil {
		tx = as.GetMaster()
	}

	return model.AttributeValues(model.AttributeValueWhere.ID.IN(ids)).DeleteAll(tx)
}

func (as *SqlAttributeValueStore) Count(options model_helper.AttributeValueFilterOptions) (int64, error) {
	return model.AttributeValues(options.Conditions...).Count(as.GetReplica())
}
