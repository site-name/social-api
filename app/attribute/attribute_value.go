package attribute

import (
	"net/http"
	"sort"

	"github.com/Masterminds/squirrel"
	"github.com/mattermost/gorp"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/modules/util"
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

func (a *ServiceAttribute) BulkUpsertAttributeValue(transaction *gorp.Transaction, values attribute.AttributeValues) (attribute.AttributeValues, *model.AppError) {
	values, err := a.srv.Store.AttributeValue().BulkUpsert(transaction, values)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}

		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}

		return nil, model.NewAppError("BulkUpsertAttributeValue", "app.attribute.error_upserting_attribute_values.app_error", nil, err.Error(), statusCode)
	}

	return values, nil
}

type Reordering struct {
	Values     attribute.AttributeValues
	Operations map[string]*int
	Field      string

	// Will contain the original data, before sorting.
	// This will be useful to look for the sort orders that
	// actually were changed
	OldSortMap map[string]*int

	cachedOrderedNodeMap  map[string]*int
	cachedAttributeValues attribute.AttributeValues

	// Will contain the list of keys kept
	// in correct order in accordance to their sort order
	OrderedPKs []string

	s      *ServiceAttribute
	runned bool // to make sure that the method `orderedNodeMap` only run once
}

func (s *ServiceAttribute) newReordering(values attribute.AttributeValues, operations map[string]*int, field string) *Reordering {
	return &Reordering{
		Values:     values,
		Operations: operations,
		Field:      field,
		s:          s,
	}
}

func (r *Reordering) orderedNodeMap(transaction *gorp.Transaction) (map[string]*int, *model.AppError) {
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

		// cached
		r.cachedAttributeValues = attributeValues

		// orderingMap has keys are attribute value ids
		var orderingMap = make(map[string]*int)
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
		for key, sortOrder := range orderingMap {
			if sortOrder != nil {
				previousSortOrder = *sortOrder
				continue
			}

			previousSortOrder++
			orderingMap[key] = model.NewInt(previousSortOrder)
		}

		// cache
		r.cachedOrderedNodeMap = make(map[string]*int)
		for key, value := range orderingMap {
			r.cachedOrderedNodeMap[key] = value
		}
		// indicate runned
		r.runned = true

		return orderingMap, nil
	}

	return r.cachedOrderedNodeMap, nil
}

func (r *Reordering) calculateNewSortOrder(pk string, move int) (int, int, int) {
	// Retrieve the position of the node to move
	nodePos := sort.SearchStrings(r.OrderedPKs, pk)

	// Set the target position from the current position
	// of the node + the relative position to move from
	targetPos := nodePos + move

	// Make sure we are not getting out of bounds
	targetPos = util.Max(0, targetPos)
	targetPos = util.Min(len(r.OrderedPKs)-1, targetPos)

	// Retrieve the target node and its sort order
	var (
		targetPk          = r.OrderedPKs[targetPos]
		orderedNodeMap, _ = r.orderedNodeMap(nil)
		targetPosition    = orderedNodeMap[targetPk]
	)

	// Return the new position
	return nodePos, targetPos, *targetPosition
}

func (s *Reordering) processMoveOperation(pk string, move *int) {
	var (
		orderedNodeMap, _ = s.orderedNodeMap(nil)
		oldSortOrder      = orderedNodeMap[pk]
	)

	// skip if nothing to do
	if move != nil && *move == 0 {
		return
	}
	if move == nil {
		move = model.NewInt(1)
	}

	_, targetPos, newSortOrder := s.calculateNewSortOrder(pk, *move) // move is non-nil now

	// Determine how we should shift for this operation
	var (
		shift  int
		range_ [2]int
	)
	if *move > 0 {
		shift = -1
		range_ = [2]int{*oldSortOrder + 1, newSortOrder}
	} else {
		shift = 1
		range_ = [2]int{newSortOrder, *oldSortOrder - 1}
	}

	// Shift the sort orders within the moving range
	s.addToSortValueIfInRange(shift, range_[0], range_[1])

	// Update the sort order of the node to move
	s.cachedOrderedNodeMap[pk] = &newSortOrder

	// Reorder the pk list
	s.OrderedPKs = util.RemoveStringFromSlice(pk, s.OrderedPKs)
	s.OrderedPKs = append( // <=> list.insert() in python3
		s.OrderedPKs[0:targetPos],
		append(
			[]string{pk},
			s.OrderedPKs[targetPos:]...,
		)...,
	)
}

func (r *Reordering) addToSortValueIfInRange(valueToAdd int, start int, end int) {
	orderedNodeMap, _ := r.orderedNodeMap(nil)
	for pk, sortOrder := range orderedNodeMap {
		if sortOrder != nil {
			if !(start <= *sortOrder && *sortOrder <= end) {
				continue
			}

			r.cachedOrderedNodeMap[pk] = model.NewInt(valueToAdd + *sortOrder)
		}
	}
}

func (r *Reordering) commit(transaction *gorp.Transaction) *model.AppError {
	// Do nothing if nothing was done
	if len(r.OldSortMap) == 0 {
		return nil
	}

	var attributeValuesMap = make(map[string]*attribute.AttributeValue)
	for _, item := range r.cachedAttributeValues {
		attributeValuesMap[item.Id] = item
	}

	changed := false

	orderedNodeMap, _ := r.orderedNodeMap(nil)
	for pk, sortOrder := range orderedNodeMap {
		if sortOrder != nil && r.OldSortMap[pk] != nil && *sortOrder != *(r.OldSortMap[pk]) {
			if !changed {
				changed = true
			}
			attributeValuesMap[pk].SortOrder = sortOrder
		}
	}

	if !changed {
		return nil
	}

	_, appErr := r.s.BulkUpsertAttributeValue(transaction, r.cachedAttributeValues)
	return appErr
}

func (r *Reordering) Run(transaction *gorp.Transaction) *model.AppError {
	for key, move := range r.Operations {
		// skip operation if it was deleted in concurrence
		orderedNodeMap, appErr := r.orderedNodeMap(transaction)
		if appErr != nil {
			return appErr
		}

		if _, ok := orderedNodeMap[key]; !ok {
			continue
		}

		r.processMoveOperation(key, move)
	}

	appErr := r.commit(transaction)
	if appErr != nil {
		return appErr
	}

	return nil
}

func (s *ServiceAttribute) PerformReordering(values attribute.AttributeValues, operations map[string]*int) *model.AppError {
	transaction, err := s.srv.Store.GetMaster().Begin()
	if err != nil {
		return model.NewAppError("PerformReordering", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer s.srv.Store.FinalizeTransaction(transaction)

	appErr := s.newReordering(values, operations, "moves").Run(transaction)
	if appErr != nil {
		return appErr
	}

	err = transaction.Commit()
	if err != nil {
		return model.NewAppError("PerformReordering", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (s *ServiceAttribute) DeleteAttributeValues(ids ...string) (int64, *model.AppError) {
	numDeleted, err := s.srv.Store.AttributeValue().Delete(ids...)
	if err != nil {
		return 0, model.NewAppError("DeleteAttributeValues", "app.attribute.error_delete_attribute_values_by_ids.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return numDeleted, nil
}
