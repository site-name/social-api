package attribute

// AssociateAttributeValuesToInstance Assign given attribute values to a product or variant.
// Note: be award this function invokes the “set“ method on the instance's
// attribute association. Meaning any values already assigned or concurrently
// assigned will be overridden by this call.
//
// `instance` must be either *Product or *ProductVariant or *Page.
// `attributeID` must be ID of processing `Attribute`
//
// Returned any must be either: `*AssignedProductAttribute` or `*AssignedPageAttribute`
// func (a *ServiceAttribute) AssociateAttributeValuesToInstance(instance any, attributeID string, values model.AttributeValueSlice) (any, *model_helper.AppError) {
// 	switch instance.(type) {
// 	case *model.Product, *model.ProductVariant, *model.Page:
// 	default:
// 		return nil, model_helper.NewAppError("AssociateAttributeValuesToInstance", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "instance"}, "", http.StatusBadRequest)
// 	}

// 	if lo.SomeBy(values, func(item *model.AttributeValue) bool { return item != nil && item.AttributeID != attributeID }) {
// 		return nil, model_helper.NewAppError("AssociateAttributeValuesToInstance", "app.attribute.attribute_does_not_own_values.app_error", nil, "given attribute does not own all given attribute values", http.StatusNotAcceptable)
// 	}

// 	valueIDs := lo.Map(values, func(item *model.AttributeValue, _ int) string { return item.ID })

// 	// Associate the attribute and the passed values
// 	assignment, appErr := a.associateAttributeToInstance(instance, attributeID)
// 	if appErr != nil {
// 		return nil, appErr
// 	}

// 	// save in bulk value relationships:
// 	switch v := assignment.(type) {
// 	case *model.AssignedProductAttribute:
// 		_, err := a.srv.Store.AssignedProductAttributeValue().SaveInBulk(v.ID, valueIDs)
// 		if err != nil {
// 			if appErr, ok := err.(*model_helper.AppError); ok {
// 				return nil, appErr
// 			}
// 			statusCode := http.StatusInternalServerError
// 			if _, ok := err.(*store.ErrInvalidInput); ok {
// 				statusCode = http.StatusNotFound
// 			}
// 			return nil, model_helper.NewAppError("AssociateAttributeValuesToInstance", "app.attribute.error_creating_assigned_product_attribute_values.app_error", nil, err.Error(), statusCode)
// 		}

// 	case *model.AssignedPageAttribute:
// 		_, err := a.srv.Store.AssignedPageAttributeValue().SaveInBulk(v.ID, valueIDs)
// 		if err != nil {
// 			if appErr, ok := err.(*model_helper.AppError); ok {
// 				return nil, appErr
// 			}
// 			statusCode := http.StatusInternalServerError
// 			if _, ok := err.(*store.ErrInvalidInput); ok {
// 				statusCode = http.StatusNotFound
// 			}
// 			return nil, model_helper.NewAppError("AssociateAttributeValuesToInstance", "app.attribute.error_creating_assigned_page_attribute_values.app_error", nil, err.Error(), statusCode)
// 		}

// 	default:
// 		return nil, model_helper.NewAppError("AssociateAttributeValuesToInstance", "app.attribute.unknown_returned_assignment.app_error", nil, "", http.StatusNotAcceptable)
// 	}

// 	// sort assigned attribute values:
// 	appErr = a.sortAssignedAttributeValues(instance, assignment, valueIDs)
// 	if appErr != nil {
// 		return nil, appErr
// 	}

// 	return assignment, nil
// }

// associateAttributeToInstance associates given attribute to given instance
//
// NOTE:
//
// `instance` must be either `*product.Product` or `*model.Page`
//
// returned any is either:
//
//	+) *AssignedProductAttribute
//	+) *AssignedPageAttribute
// func (a *ServiceAttribute) associateAttributeToInstance(instance any, attributeID string) (any, *model_helper.AppError) {
// 	switch v := instance.(type) {
// 	case *model.Product:
// 		attributeProduct, appErr := a.AttributeProductByOption(&model.AttributeProductFilterOption{
// 			Conditions: squirrel.Eq{
// 				model.AttributeProductTableName + ".ProductTypeID": v.ProductTypeID,
// 				model.AttributeProductTableName + ".AttributeID":   attributeID,
// 			},
// 		})
// 		if appErr != nil {
// 			return nil, appErr
// 		}

// 		return a.GetOrCreateAssignedProductAttribute(&model.AssignedProductAttribute{
// 			ProductID:    v.ID,
// 			AssignmentID: attributeProduct.Id,
// 		})

// 	case *model.Page:
// 		attributePage, appErr := a.AttributePageByOption(&model.AttributePageFilterOption{
// 			Conditions: squirrel.Eq{
// 				model.AttributePageTableName + ".AttributeID": attributeID,
// 				model.AttributePageTableName + ".PageTypeID":  v.PageTypeID,
// 			},
// 		})
// 		if appErr != nil {
// 			return nil, appErr
// 		}

// 		return a.GetOrCreateAssignedPageAttribute(&model.AssignedPageAttribute{
// 			PageID:       v.Id,
// 			AssignmentID: attributePage.Id,
// 		})

// 	default:
// 		return nil, model_helper.NewAppError("associateAttributeToInstance", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "instance"}, "", http.StatusBadRequest)
// 	}
// }

// sortAssignedAttributeValues Sorts assigned attribute values based on values list order
//
// (instance - assignment) must be provided by pair like this:
//
//	+) *Product        - *AssignedProductAttribute
//
//	+) *ProductVariant - *AssignedVariantAttribute
//
//	+) *Page           - *AssignedPageAttribute
// func (a *ServiceAttribute) sortAssignedAttributeValues(instance any, assignment any, valueIDs []string) *model_helper.AppError {
// 	if instance == nil || assignment == nil || len(valueIDs) == 0 {
// 		return model_helper.NewAppError("sortAssignedAttributeValues", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "assignment or instance or valuesIDs"}, "", http.StatusBadRequest)
// 	}

// 	switch instance.(type) {
// 	case *model.Product:
// 		if assignmentValue, ok := assignment.(*model.AssignedProductAttribute); !ok {
// 			assignedProductAttrValues, attrValues, err := a.srv.Store.AssignedProductAttributeValue().SelectForSort(assignmentValue.ID)
// 			// err can be *store.ErrNotFound or system error
// 			if err != nil {
// 				return model_helper.NewAppError("sortAssignedAttributeValues", "app.attribute.select_assigned_product_attribute_values_for_sort.app_error", nil, err.Error(), http.StatusInternalServerError)
// 			}
// 			// NOTE: this sort can be done since len(assignedProductAttrValues) == len(attrValues)
// 			sort.Slice(assignedProductAttrValues, func(i, j int) bool {
// 				return sort.SearchStrings(valueIDs, attrValues[i].ID) <= sort.SearchStrings(valueIDs, attrValues[j].ID)
// 			})
// 			for i, value := range assignedProductAttrValues {
// 				value.SortOrder = &i
// 			}
// 			// update if database:
// 			if err = a.srv.Store.AssignedProductAttributeValue().UpdateInBulk(assignedProductAttrValues); err != nil {
// 				return model_helper.NewAppError("sortAssignedAttributeValues", "app.attribute.error_updating_assigned_product_attribute_values.app_error", nil, err.Error(), http.StatusInternalServerError)
// 			}
// 		}
// 		// other types are not accepted and returns an error:
// 		return model_helper.NewAppError("sortAssignedAttributeValues", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "assignment"}, "", http.StatusBadRequest)

// 	case *model.ProductVariant:
// 		if assignmentValue, ok := assignment.(*model.AssignedVariantAttribute); ok {
// 			assignedVariantAttrValues, attrValues, err := a.srv.Store.AssignedVariantAttributeValue().SelectForSort(assignmentValue.ID)
// 			// err can be *store.ErrNotFound or system error
// 			if err != nil {
// 				return model_helper.NewAppError("sortAssignedAttributeValues", "app.attribute.select_assigned_variant_attribute_values_for_sort.app_error", nil, err.Error(), http.StatusInternalServerError)
// 			}
// 			// NOTE: this sort can be done since len(assignedVariantAttrValues) == len(attrValues)
// 			sort.Slice(assignedVariantAttrValues, func(i, j int) bool {
// 				return sort.SearchStrings(valueIDs, attrValues[i].ID) <= sort.SearchStrings(valueIDs, attrValues[j].ID)
// 			})
// 			for i, value := range assignedVariantAttrValues {
// 				value.SortOrder = &i
// 			}
// 			// update if database:
// 			if err = a.srv.Store.AssignedVariantAttributeValue().UpdateInBulk(assignedVariantAttrValues); err != nil {
// 				return model_helper.NewAppError("sortAssignedAttributeValues", "app.attribute.error_updating_assigned_variant_attribute_values.app_error", nil, err.Error(), http.StatusInternalServerError)
// 			}
// 		}
// 		// other types are not accepted and returns an error:
// 		return model_helper.NewAppError("sortAssignedAttributeValues", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "assignment"}, "", http.StatusBadRequest)

// 	case *model.Page:
// 		if assignmentValue, ok := assignment.(*model.AssignedPageAttribute); ok {
// 			assignedPageAttrValues, attrValues, err := a.srv.Store.AssignedPageAttributeValue().SelectForSort(assignmentValue.ID)
// 			// err can be *store.ErrNotFound or system error
// 			if err != nil {
// 				return model_helper.NewAppError("sortAssignedAttributeValues", "app.attribute.select_assigned_page_attribute_values_for_sort.app_error", nil, err.Error(), http.StatusInternalServerError)
// 			}
// 			// NOTE: this sort can be done since len(assignedPageAttrValues) == len(attrValues)
// 			sort.Slice(assignedPageAttrValues, func(i, j int) bool {
// 				return sort.SearchStrings(valueIDs, attrValues[i].ID) <= sort.SearchStrings(valueIDs, attrValues[j].ID)
// 			})
// 			for i, value := range assignedPageAttrValues {
// 				value.SortOrder = &i
// 			}
// 			// update if database:
// 			if err = a.srv.Store.AssignedPageAttributeValue().UpdateInBulk(assignedPageAttrValues); err != nil {
// 				return model_helper.NewAppError("sortAssignedAttributeValues", "app.attribute.error_updating_assigned_page_attribute_values.app_error", nil, err.Error(), http.StatusInternalServerError)
// 			}
// 		}
// 		// other types are not accepted and returns an error:
// 		return model_helper.NewAppError("sortAssignedAttributeValues", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "assignment"}, "", http.StatusBadRequest)

// 	default:
// 		return model_helper.NewAppError("sortAssignedAttributeValues", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "instance"}, "", http.StatusBadRequest)
// 	}
// }
