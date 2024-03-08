package attribute

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (a *ServiceAttribute) AttributeValuesOfAttribute(attributeID string) (model.AttributeValueSlice, *model_helper.AppError) {
	return a.FilterAttributeValuesByOptions(model_helper.AttributeValueFilterOptions{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.AttributeValueWhere.AttributeID.EQ(attributeID),
		),
	})
}

func (s *ServiceAttribute) FilterAttributeValuesByOptions(option model_helper.AttributeValueFilterOptions) (model.AttributeValueSlice, *model_helper.AppError) {
	values, err := s.srv.Store.AttributeValue().FilterByOptions(option)
	if err != nil {
		return nil, model_helper.NewAppError("FilterAttributeValuesByOptions", "app.attribute.error_finding_attribute_values_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return values, nil
}

func (a *ServiceAttribute) BulkUpsertAttributeValue(transaction boil.ContextTransactor, values model.AttributeValueSlice) (model.AttributeValueSlice, *model_helper.AppError) {
	values, err := a.srv.Store.AttributeValue().Upsert(transaction, values)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}

		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}

		return nil, model_helper.NewAppError("BulkUpsertAttributeValue", "app.attribute.error_upserting_attribute_values.app_error", nil, err.Error(), statusCode)
	}

	return values, nil
}

type Reordering struct {
	Values     model.AttributeValueSlice
	Operations map[string]*int
	Field      string

	// Will contain the original data, before sorting.
	// This will be useful to look for the sort orders that
	// actually were changed
	OldSortMap map[string]*int

	cachedOrderedNodeMap  map[string]*int
	cachedAttributeValues model.AttributeValueSlice

	// Will contain the list of keys kept
	// in correct order in accordance to their sort order
	OrderedPKs util.AnyArray[string]

	s      *ServiceAttribute
	runned bool // to make sure that the method `orderedNodeMap` only run once
}

func (s *ServiceAttribute) newReordering(values model.AttributeValueSlice, operations map[string]*int, field string) *Reordering {
	return &Reordering{
		Values:     values,
		Operations: operations,
		Field:      field,
		s:          s,
	}
}

func (r *Reordering) orderedNodeMap(transaction boil.ContextTransactor) (map[string]*int, *model_helper.AppError) {
	if !r.runned { // check if runned or not
		// indicate runned
		r.runned = true

		valueIDs := lo.Map(r.Values, func(a *model.AttributeValue, _ int) string { return a.ID })

		attributeValues, appErr := r.s.FilterAttributeValuesByOptions(model_helper.AttributeValueFilterOptions{
			CommonQueryOptions: model_helper.NewCommonQueryOptions(
				qm.OrderBy(fmt.Sprintf("%s %s NULLS LAST", model.AttributeValueColumns.SortOrder, model_helper.ASC)),
				model.AttributeValueWhere.ID.IN(valueIDs),
			),
		})
		if appErr != nil {
			return nil, appErr
		}

		// cached
		r.cachedAttributeValues = attributeValues

		// orderingMap has keys are attribute value ids
		var orderingMap = make(map[string]*int)
		for _, value := range attributeValues {
			orderingMap[value.ID] = value.SortOrder.Int
		}

		// copy
		r.OldSortMap = make(map[string]*int)
		for key, value := range orderingMap {
			r.OldSortMap[key] = value
		}

		r.OrderedPKs = lo.Map(attributeValues, func(a *model.AttributeValue, _ int) string { return a.ID })

		previousSortOrder := 0

		// Add sort order to null values
		for key, sortOrder := range orderingMap {
			if sortOrder != nil {
				previousSortOrder = *sortOrder
				continue
			}

			previousSortOrder++
			orderingMap[key] = model_helper.GetPointerOfValue(previousSortOrder)
		}

		// cache
		r.cachedOrderedNodeMap = make(map[string]*int)
		for key, value := range orderingMap {
			r.cachedOrderedNodeMap[key] = value
		}

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
	targetPos = max(0, targetPos)
	targetPos = min(len(r.OrderedPKs)-1, targetPos)

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
		move = model_helper.GetPointerOfValue(1)
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
	s.OrderedPKs = s.OrderedPKs.Remove(pk)
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
		if sortOrder != nil && !(start <= *sortOrder && *sortOrder <= end) {
			continue
		}

		r.cachedOrderedNodeMap[pk] = model_helper.GetPointerOfValue(valueToAdd + *sortOrder)
	}
}

func (r *Reordering) commit(transaction boil.ContextTransactor) *model_helper.AppError {
	// Do nothing if nothing was done
	if len(r.OldSortMap) == 0 {
		return nil
	}

	copiedAttributeValues := model_helper.DeepCopyAttributeValueSlice(r.cachedAttributeValues)
	var attributeValuesMap = lo.SliceToMap(copiedAttributeValues, func(a *model.AttributeValue) (string, *model.AttributeValue) { return a.ID, a })

	changed := false

	orderedNodeMap, _ := r.orderedNodeMap(transaction)
	for pk, sortOrder := range orderedNodeMap {
		oldSortOrder, exist := r.OldSortMap[pk]
		if exist && oldSortOrder != nil && sortOrder != nil && *oldSortOrder != *sortOrder {
			attributeValuesMap[pk].SortOrder.Int = sortOrder
			changed = true
		}
	}

	if !changed {
		return nil
	}

	_, appErr := r.s.BulkUpsertAttributeValue(transaction, copiedAttributeValues)
	return appErr
}

func (r *Reordering) Run(transaction boil.ContextTransactor) *model_helper.AppError {
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

	return r.commit(transaction)
}

func (s *ServiceAttribute) PerformReordering(values model.AttributeValueSlice, operations map[string]*int) *model_helper.AppError {
	transaction, err := s.srv.Store.GetMaster().BeginTx(context.Background(), nil)
	if err != nil {
		return model_helper.NewAppError("PerformOrdering", model_helper.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer s.srv.Store.FinalizeTransaction(transaction)

	appErr := s.newReordering(values, operations, "moves").Run(transaction)
	if appErr != nil {
		return appErr
	}

	err = transaction.Commit()
	if err != nil {
		return model_helper.NewAppError("PerformReordering", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (s *ServiceAttribute) DeleteAttributeValues(tx boil.ContextTransactor, ids []string) (int64, *model_helper.AppError) {
	numDeleted, err := s.srv.Store.AttributeValue().Delete(tx, ids)
	if err != nil {
		return 0, model_helper.NewAppError("DeleteAttributeValues", "app.attribute.error_delete_attribute_values_by_ids.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return numDeleted, nil
}
