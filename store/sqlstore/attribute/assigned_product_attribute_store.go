package attribute

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SqlAssignedProductAttributeStore struct {
	store.Store
}

func NewSqlAssignedProductAttributeStore(s store.Store) store.AssignedProductAttributeStore {
	return &SqlAssignedProductAttributeStore{s}
}

// func (as *SqlAssignedProductAttributeStore) Save(assignedProductAttribute model.AssignedProductAttribute) (*model.AssignedProductAttribute, error) {
// 	panic("unimplemented")
// }

// func (as *SqlAssignedProductAttributeStore) Get(id string) (*model.AssignedProductAttribute, error) {
// 	record, err := model.FindAssignedProductAttribute(as.GetReplica(), id)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, store.NewErrNotFound(model.TableNames.AssignedProductAttributes, id)
// 		}
// 		return nil, err
// 	}

// 	return record, nil
// }

func (as *SqlAssignedProductAttributeStore) commonQueryBuilder(options model_helper.AssignedProductAttributeFilterOption) []qm.QueryMod {
	conds := options.Conditions
	if options.AttributeProduct_Attribute_VisibleInStoreFront != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.CategoryAttributes, model.CategoryAttributeTableColumns.ID, model.AssignedProductAttributeTableColumns.AssignmentID)),
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Attributes, model.AttributeTableColumns.ID, model.CategoryAttributeTableColumns.AttributeID)),
			model.AttributeWhere.VisibleInStorefront.EQ(*options.AttributeProduct_Attribute_VisibleInStoreFront),
		)
	}
	return conds
}

func (as *SqlAssignedProductAttributeStore) GetWithOption(option model_helper.AssignedProductAttributeFilterOption) (*model.AssignedProductAttribute, error) {
	conds := as.commonQueryBuilder(option)
	record, err := model.AssignedProductAttributes(conds...).One(as.GetReplica())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.AssignedProductAttributes, "options")
		}
		return nil, errors.Wrap(err, "GetWithOption")
	}

	return record, nil
}

func (as *SqlAssignedProductAttributeStore) FilterByOptions(options model_helper.AssignedProductAttributeFilterOption) (model.AssignedProductAttributeSlice, error) {
	conds := as.commonQueryBuilder(options)
	return model.AssignedProductAttributes(conds...).All(as.GetReplica())
}
