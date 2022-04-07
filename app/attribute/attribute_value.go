package attribute

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/mattermost/gorp"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

func (a *ServiceAttribute) AttributeValuesOfAttribute(attributeID string) ([]*attribute.AttributeValue, *model.AppError) {
	attrValues, err := a.srv.Store.AttributeValue().FilterByOptions(attribute.AttributeValueFilterOptions{
		AttributeID: squirrel.Eq{store.AttributeValueTableName + ".AttributeID": attributeID},
	})
	var (
		statusCode = 0
		errMsg     string
	)
	if err != nil {
		errMsg = err.Error()
		statusCode = http.StatusInternalServerError
	}
	if len(attrValues) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("AttributeValuesOfAttribute", "app.attribute.error_finding_attribute_values_by_attribute_id.app_error", nil, errMsg, statusCode)
	}

	return attrValues, nil
}

func (s *ServiceAttribute) FilterAttributeValuesByOptions(option attribute.AttributeValueFilterOptions) (attribute.AttributeValues, *model.AppError) {
	values, err := s.srv.Store.AttributeValue().FilterByOptions(option)

	var (
		statusCode = 0
		errMsg     = ""
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errMsg = err.Error()
	} else if len(values) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("FilterAttributeValuesByOptions", "app.attribute.error_finding_attribute_values_by_options.app_error", nil, errMsg, statusCode)
	}

	return values, nil
}

// UpsertAttributeValue insderts or updates given attribute value then returns it
func (a *ServiceAttribute) UpsertAttributeValue(attrValue *attribute.AttributeValue) (*attribute.AttributeValue, *model.AppError) {
	attrValue, err := a.srv.Store.AttributeValue().Upsert(attrValue)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}

		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}

		return nil, model.NewAppError("UpsertAttributeValue", "app.attribute.error_upserting_attribute_value.app_error", nil, err.Error(), statusCode)
	}

	return attrValue, nil
}

type Reordering struct {
	Values     attribute.AttributeValues
	Operations map[string]*int
	Field      string

	// Will contain the original data, before sorting.
	// This will be useful to look for the sort orders that
	// actually were changed
	OldSortMap map[string]*int

	// Will contain the list of keys kept
	// in correct order in accordance to their sort order
	OrderedPKs []string

	s      *ServiceAttribute
	runned bool
}

func (s *ServiceAttribute) NewReordering(values attribute.AttributeValues, operations map[string]*int, field string) *Reordering {
	return &Reordering{
		Values:     values,
		Operations: operations,
		Field:      field,
		s:          s,
	}
}

func (r *Reordering) OrderedNodeMap(transaction *gorp.Transaction) (map[string]*int, *model.AppError) {
	if !r.runned { // check if runned or not
		attributeValues, appErr := r.s.FilterAttributeValuesByOptions(attribute.AttributeValueFilterOptions{
			Transaction:     transaction,
			OrderBy:         store.AttributeValueTableName + ".Id ASC, " + store.AttributeValueTableName + ".SortOrder ASC NULLS LAST",
			Id:              squirrel.Eq{store.AttributeValueTableName + ".Id": r.Values.IDs()},
			SelectForUpdate: true,
		})
		if appErr != nil {
			return nil, appErr
		}

		var orderingMap = map[string]*int{}
		for _, value := range attributeValues {
			orderingMap[value.Id] = value.SortOrder
		}

		// copy
		r.OldSortMap = make(map[string]*int)
		for key, value := range orderingMap {
			r.OldSortMap[key] = value
		}

		r.OrderedPKs = attributeValues.IDs()

		previousSortOrder := 0

		// Add sort order to null values
		for key, value := range orderingMap {
			if value != nil {
				previousSortOrder = *value
				continue
			}

			previousSortOrder++
			i := previousSortOrder
			orderingMap[key] = &i
		}

		// indicate runned
		r.runned = true

		return orderingMap, nil
	}

	return r.OldSortMap, nil
}

func (s *Reordering) ProcessMoveOperation(pk string, move *int) {
	oldSortOrder, _ := s.OrderedNodeMap(nil)

	// skip if nothing to do
	if move != nil && *move == 0 {
		return
	}
	if move == nil {
		move = model.NewInt(1)
	}

}

func (r *Reordering) Run(transaction *gorp.Transaction) *model.AppError {
	for key, value := range r.Operations {
		// skip operation if it was deleted in concurrence
		orderedNodeMap, appErr := r.OrderedNodeMap(transaction)
		if appErr != nil {
			return appErr
		}

		if _, ok := orderedNodeMap[key]; !ok {
			continue
		}

		r.ProcessMoveOperation(key, value)
	}
}
