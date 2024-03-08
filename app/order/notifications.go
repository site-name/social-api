package order

import (
	"net/url"

	"github.com/mattermost/squirrel"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
)

func (s *ServiceOrder) getDefaultImagesPayload(images model.ProductMedias) model_types.JSONString {
	// NOTE:
	// TODO: implement me
	return nil
}

func (s *ServiceOrder) getProductAttributes(product *model.Product) ([]model_types.JSONString, *model_helper.AppError) {
	assignedPrdAttributes, appErr := s.srv.AttributeService().AssignedProductAttributesByOption(&model.AssignedProductAttributeFilterOption{
		Conditions: squirrel.Expr(model.AssignedProductAttributeTableName+".ProductID = ?", product.Id),
		Preloads: []string{
			"Values",
			"AttributeProduct.Attribute",
		},
	})
	if appErr != nil {
		return nil, appErr
	}

	res := []model_types.JSONString{}
	for _, attr := range assignedPrdAttributes {
		data := model_types.JSONString{}

		if attr.AttributeProduct != nil && attr.AttributeProduct.Attribute != nil {
			data["assignment"] = model_types.JSONString{
				"attribute": model_types.JSONString{
					"slug": attr.AttributeProduct.Attribute.Slug,
					"name": attr.AttributeProduct.Attribute.Name,
				},
			}
		}
		if len(attr.Values) > 0 {
			data["values"] = lo.Map(attr.Values, func(value *model.AttributeValue, _ int) model_types.JSONString {
				return model_types.JSONString{
					"name":     value.Name,
					"value":    value.Value,
					"slug":     value.Slug,
					"file_url": model.GetValueOfpointerOrNil(value.FileUrl),
				}
			})
		}
	}

	return res, nil
}

func (s *ServiceOrder) getProductPayload(product *model.Product) (model_types.JSONString, *model_helper.AppError) {
	productMedias, appErr := s.srv.ProductService().ProductMediasByOption(&model.ProductMediaFilterOption{
		Conditions: squirrel.Expr(model.ProductMediaTableName+".ProductID = ?", product.Id),
	})
	if appErr != nil {
		return nil, appErr
	}

	images := lo.Filter(productMedias, func(item *model.ProductMedia, _ int) bool { return item != nil && item.Type == model.IMAGE })

	attributes, appErr := s.getProductAttributes(product)
	if appErr != nil {
		return nil, appErr
	}

	// NOTE:
	// TODO: add image field to result below

	res := model_types.JSONString{
		"id":         product.Id,
		"attributes": attributes,
		"weight":     product.WeightString(),
	}
	res.Merge(s.getDefaultImagesPayload(images))

	return res, nil
}

func (s *ServiceOrder) getProductVariantPayload(variant *model.ProductVariant) (model_types.JSONString, *model_helper.AppError) {
	productMedias := variant.ProductMedias

	if len(productMedias) == 0 {
		var appErr *model_helper.AppError
		productMedias, appErr = s.srv.ProductService().ProductMediasByOption(&model.ProductMediaFilterOption{
			VariantID: squirrel.Expr(model.ProductVariantMediaTableName+".product_variant_id = ?", variant.Id),
		})
		if appErr != nil {
			return nil, appErr
		}
	}

	imageMedias := lo.Filter(productMedias, func(item *model.ProductMedia, index int) bool { return item != nil && item.Type == model.IMAGE })

	res := model_types.JSONString{
		"id":                        variant.Id,
		"weight":                    variant.WeightString(),
		"is_preorder":               variant.IsPreorderActive(),
		"preorder_global_threshold": model.GetValueOfpointerOrNil(variant.PreOrderGlobalThreshold),
		"preorder_end_date":         model.GetValueOfpointerOrNil(variant.PreorderEndDate),
	}
	res.Merge(s.getDefaultImagesPayload(imageMedias))

	return res, nil
}

// NOTE: given order line should have `ProductVariant` field preloaded by caller(s)
func (s *ServiceOrder) getOrderLinePayload(line *model.OrderLine) (model_types.JSONString, *model_helper.AppError) {
	orderLineIsDigital, appErr := s.OrderLineIsDigital(line)
	if appErr != nil {
		return nil, appErr
	}

	var digitalUrl string

	if orderLineIsDigital {
		_, appErr := s.srv.ProductService().DigitalContentURLSByOptions(&model.DigitalContentUrlFilterOptions{
			Conditions: squirrel.Expr(model.DigitalContentURLTableName+".LineID = ?", line.Id),
		})
		if appErr != nil {
			return nil, appErr
		}

		slog.Debug("please construct url for digital content url")
		// TODO:
		// Add a step to construct URL for digitalContentUrl found
	}

	var variantId, productId string
	if line.ProductVariant != nil {
		variantId = line.ProductVariant.Id
		productId = line.ProductVariant.ProductID

	} else if line.VariantID != nil {
		variantId = *line.VariantID

		variant, appErr := s.srv.ProductService().ProductVariantById(*line.VariantID)
		if appErr != nil {
			return nil, appErr
		}
		productId = variant.ProductID
	}

	line.PopulateNonDbFields() // this call is needed

	translatedProductName := line.TranslatedProductName
	translatedVariantName := line.TranslatedVariantName
	if translatedProductName == "" {
		translatedProductName = line.ProductName
	}
	if translatedVariantName == "" {
		translatedVariantName = line.VariantName
	}

	// these evaluations below helps prevent nil pointer dereference
	var unitTaxAmount, totalTaxAmount *decimal.Decimal
	if line.UnitPrice != nil {
		unitTaxAmount = &line.UnitPrice.Tax().Amount
	}
	if line.TotalPrice != nil {
		totalTaxAmount = &line.TotalPrice.Tax().Amount
	}

	return model_types.JSONString{
		"id":                      line.Id,
		"product":                 productId, // type: ignore
		"product_name":            line.ProductName,
		"translated_product_name": translatedProductName,
		"variant_name":            line.VariantName,
		"variant":                 variantId, // type: ignore
		"translated_variant_name": translatedVariantName,
		"quantity":                line.Quantity,
		"quantity_fulfilled":      line.QuantityFulfilled,
		"currency":                line.Currency,
		"is_shipping_required":    line.IsShippingRequired,
		"is_digital":              orderLineIsDigital,
		"digital_url":             digitalUrl, // TODO: implement this
		"unit_discount_type":      line.UnitDiscountType,
		"unit_tax_amount":         model.GetValueOfpointerOrNil(unitTaxAmount),
		"total_tax_amount":        model.GetValueOfpointerOrNil(totalTaxAmount),
		"total_gross_amount":      model.GetValueOfpointerOrNil(line.TotalPriceGrossAmount),
		"total_net_amount":        model.GetValueOfpointerOrNil(line.TotalPriceNetAmount),
		"tax_rate":                model.GetValueOfpointerOrNil(line.TaxRate),
		"product_sku":             model.GetValueOfpointerOrNil(line.ProductSku),
		"product_variant_id":      model.GetValueOfpointerOrNil(line.ProductVariantID),
		"unit_price_net_amount":   model.GetValueOfpointerOrNil(line.UnitPriceNetAmount),
		"unit_price_gross_amount": model.GetValueOfpointerOrNil(line.UnitPriceGrossAmount),
		"unit_discount_value":     model.GetValueOfpointerOrNil(line.UnitDiscountValue),
		"unit_discount_reason":    model.GetValueOfpointerOrNil(line.UnitDiscountReason),
		"unit_discount_amount":    model.GetValueOfpointerOrNil(line.UnitDiscountAmount),
	}, nil
}

func (s *ServiceOrder) getLinesPayload(orderLines model.OrderLineSlice) ([]model_types.JSONString, *model_helper.AppError) {
	// if some order line(s) don't have ProductVariant field populated, then populate them
	if lo.SomeBy(orderLines, func(item *model.OrderLine) bool { return item != nil && item.ProductVariant == nil }) {
		var appErr *model_helper.AppError
		orderLines, appErr = s.srv.OrderService().OrderLinesByOption(&model.OrderLineFilterOption{
			Conditions: squirrel.Eq{model.OrderLineTableName + ".Id": orderLines.IDs()},
			Preload:    []string{"ProductVariant"},
		})
		if appErr != nil {
			return nil, appErr
		}
	}

	var res = make([]model_types.JSONString, 0, len(orderLines))
	for _, line := range orderLines {
		value, appErr := s.getOrderLinePayload(line)
		if appErr != nil {
			return nil, appErr
		}

		res = append(res, value)
	}

	return res, nil
}

func getAddressPayload(address *model.Address) model_types.JSONString {
	if address == nil {
		return nil
	}

	return model_types.JSONString{
		"first_name":       address.FirstName,
		"last_name":        address.LastName,
		"company_name":     address.CompanyName,
		"street_address_1": address.StreetAddress1,
		"street_address_2": address.StreetAddress2,
		"city":             address.City,
		"city_area":        address.CityArea,
		"postal_code":      address.PostalCode,
		"country":          address.Country,
		"country_area":     address.CountryArea,
		"phone":            address.Phone,
	}
}

func (s *ServiceOrder) getDiscountsPayload(order *model.Order) (model_types.JSONString, *model_helper.AppError) {
	orderDiscounts, appErr := s.srv.DiscountService().OrderDiscountsByOption(&model.OrderDiscountFilterOption{
		Conditions: squirrel.Expr(model.OrderDiscountTableName+".OrderID = ?", order.Id),
	})
	if appErr != nil {
		return nil, appErr
	}

	var (
		allDiscounts                           = make([]model_types.JSONString, len(orderDiscounts))
		voucherDiscount model_types.JSONString = nil
		discountAmount                         = decimal.NewFromInt(0)
	)

	for idx, orderDiscount := range orderDiscounts {
		discountObj := model_types.JSONString{
			"type":            orderDiscount.Type,
			"value_type":      orderDiscount.ValueType,
			"value":           model.GetValueOfpointerOrNil(orderDiscount.Value),
			"amount_value":    model.GetValueOfpointerOrNil(orderDiscount.AmountValue),
			"name":            model.GetValueOfpointerOrNil(orderDiscount.Name),
			"translated_name": model.GetValueOfpointerOrNil(orderDiscount.TranslatedName),
			"reason":          model.GetValueOfpointerOrNil(orderDiscount.Reason),
		}

		allDiscounts[idx] = discountObj
		if orderDiscount.Type == model.VOUCHER {
			voucherDiscount = discountObj
		}

		if orderDiscount.AmountValue != nil {
			discountAmount = discountAmount.Add(*orderDiscount.AmountValue)
		}
	}

	return model_types.JSONString{
		"voucher_discount": voucherDiscount,
		"discounts":        allDiscounts,
		"discount_amount":  discountAmount,
	}, nil
}

func (s *ServiceOrder) getDefaultOrderPayload(order *model.Order, redirectUrl *string) (model_types.JSONString, *model_helper.AppError) {
	var orderDetailsUrl string

	if redirectUrl != nil {
		_, err := url.Parse(*redirectUrl)
		if err == nil {
			params := url.Values{"token": []string{order.Token}}
			orderDetailsUrl, _ = util.PrepareUrl(params, *redirectUrl)
		}
	}

	orderSubTotal, appErr := s.OrderSubTotal(order)
	if appErr != nil {
		return nil, appErr
	}

	tax := order.TotalGrossAmount.Sub(*order.TotalNetAmount)

	orderLines, appErr := s.OrderLinesByOption(&model.OrderLineFilterOption{
		Conditions: squirrel.Expr(model.OrderLineTableName+".OrderID = ?", order.Id),
		Preload:    []string{"ProductVariant"},
	})
	if appErr != nil {
		return nil, appErr
	}

	customerEmail, appErr := s.CustomerEmail(order)
	if appErr != nil {
		return nil, appErr
	}

	// find order shipping/billing addresses
	var orderShippingAddress, orderBilingAddress *model.Address
	orderAddressIds := make([]string, 0, 2)
	if order.ShippingAddressID != nil {
		orderAddressIds = append(orderAddressIds, *order.ShippingAddressID)
	}
	if order.BillingAddressID != nil {
		orderAddressIds = append(orderAddressIds, *order.BillingAddressID)
	}
	if len(orderAddressIds) > 0 {
		orderAddresses, appErr := s.srv.AccountService().AddressesByOption(&model.AddressFilterOption{
			Conditions: squirrel.Eq{model.AddressTableName + ".Id": orderAddressIds},
		})
		if appErr != nil {
			return nil, appErr
		}

		for _, addr := range orderAddresses {
			if order.BillingAddressID != nil && *order.BillingAddressID == addr.Id {
				orderBilingAddress = addr
			} else if order.ShippingAddressID != nil && *order.ShippingAddressID == addr.Id {
				orderShippingAddress = addr
			}
		}
	}

	orderLinesPayload, appErr := s.getLinesPayload(orderLines)
	if appErr != nil {
		return nil, appErr
	}

	orderDiscountPayload, appErr := s.getDiscountsPayload(order)
	if appErr != nil {
		return nil, appErr
	}

	orderPayload := model_types.JSONString{
		// "discount_amount": order.Discount, // TODO: check this
		"id":                              order.Id,
		"token":                           order.Token,
		"display_gross_prices":            model.GetValueOfpointerOrNil(order.DisplayGrossPrices),
		"currency":                        order.Currency,
		"total_gross_amount":              model.GetValueOfpointerOrNil(order.TotalGrossAmount),
		"total_net_amount":                model.GetValueOfpointerOrNil(order.TotalNetAmount),
		"undiscounted_total_gross_amount": model.GetValueOfpointerOrNil(order.UnDiscountedTotalGrossAmount),
		"undiscounted_total_net_amount":   model.GetValueOfpointerOrNil(order.UnDiscountedTotalNetAmount),
		"status":                          order.Status,
		"metadata":                        order.Metadata,
		"private_metadata":                order.PrivateMetadata,
		"user_id":                         model.GetValueOfpointerOrNil(order.UserID),
		"language_code":                   order.LanguageCode,

		"channel_id":                  order.ChannelID,
		"created":                     order.CreateAt,
		"shipping_price_net_amount":   order.ShippingPriceNetAmount,
		"shipping_price_gross_amount": order.ShippingPriceGrossAmount,
		"order_details_url":           orderDetailsUrl,
		"email":                       customerEmail,
		"subtotal_gross_amount":       orderSubTotal.Gross.Amount,
		"subtotal_net_amount":         orderSubTotal.Net.Amount,
		"tax_amount":                  tax,
		"lines":                       orderLinesPayload,
		"billing_address":             getAddressPayload(orderBilingAddress),
		"shipping_address":            getAddressPayload(orderShippingAddress),
		"shipping_method_name":        order.ShippingMethodName,
	}

	orderPayload.Merge(orderDiscountPayload)

	return orderPayload, nil
}

func (s *ServiceOrder) getDefaultFulfillmentLinePayload(line *model.FulfillmentLine) (model_types.JSONString, *model_helper.AppError) {
	orderLine := line.OrderLine

	if orderLine == nil {
		var appErr *model_helper.AppError
		orderLine, appErr = s.OrderLineById(line.OrderLineID)
		if appErr != nil {
			return nil, appErr
		}
	}

	orderLinePayload, appErr := s.getOrderLinePayload(orderLine)
	if appErr != nil {
		return nil, appErr
	}

	return model_types.JSONString{
		"id":         line.Id,
		"order_line": orderLinePayload,
		"quantity":   line.Quantity,
	}, nil
}

func (s *ServiceOrder) getDefaultFulfillmentPayload(order *model.Order, fulfillment *model.Fulfillment) (model_types.JSONString, *model_helper.AppError) {
	fulfillmentLines, appErr := s.FulfillmentLinesByOption(&model.FulfillmentLineFilterOption{
		Conditions: squirrel.Expr(model.FulfillmentLineTableName+".FulfillmentID = ?", fulfillment.Id),
		Preloads:   []string{"OrderLine"},
	})
	if appErr != nil {
		return nil, appErr
	}

	// TODO: check performance this loop
	var digitalLinesPayloads, physicalLinesPayloads []model_types.JSONString
	for _, line := range fulfillmentLines {
		orderLineIsDigital, appErr := s.OrderLineIsDigital(line.OrderLine)
		if appErr != nil {
			return nil, appErr
		}

		fulfillmentLinePayload, appErr := s.getDefaultFulfillmentLinePayload(line)
		if appErr != nil {
			return nil, appErr
		}

		if orderLineIsDigital {
			digitalLinesPayloads = append(digitalLinesPayloads, fulfillmentLinePayload)
			continue
		}
		physicalLinesPayloads = append(physicalLinesPayloads, fulfillmentLinePayload)
	}

	orderPayload, appErr := s.getDefaultOrderPayload(order, order.RedirectUrl)
	if appErr != nil {
		return nil, appErr
	}

	customerEmail, appErr := s.CustomerEmail(order)
	if appErr != nil {
		return nil, appErr
	}

	res := model_types.JSONString{
		"order": orderPayload,
		"fulfillment": model_types.JSONString{
			"tracking_number":        fulfillment.TrackingNumber,
			"is_tracking_number_url": fulfillment.IsTrackingNumberURL(),
		},
		"physical_lines":  physicalLinesPayloads,
		"digital_lines":   digitalLinesPayloads,
		"recipient_email": customerEmail,
	}
	res.Merge(s.srv.GetSiteContext())

	return res, nil
}

// SendPaymentConfirmation sends notification with the payment confirmation
func (s *ServiceOrder) SendPaymentConfirmation(order model.Order, manager interfaces.PluginManagerInterface) *model_helper.AppError {
	panic("not implemented")
}

func (s *ServiceOrder) SendOrderCancelledConfirmation(order *model.Order, user *model.User, _, manager interfaces.PluginManagerInterface) *model_helper.AppError {
	panic("not implemented")
}

// SendOrderConfirmation sends notification with order confirmation
func (s *ServiceOrder) SendOrderConfirmation(order *model.Order, redirectURL string, manager interfaces.PluginManagerInterface) *model_helper.AppError {
	panic("not implemented")
}

// SendFulfillmentConfirmationToCustomer
//
// NOTE: user can be nil
func (s *ServiceOrder) SendFulfillmentConfirmationToCustomer(order *model.Order, fulfillment *model.Fulfillment, user *model.User, _, manager interfaces.PluginManagerInterface) *model_helper.AppError {
	panic("not implemented")
}

// SendOrderConfirmed Send email which tells customer that order has been confirmed
func (s *ServiceOrder) SendOrderConfirmed(order model.Order, user *model.User, _ any, manager interfaces.PluginManagerInterface) {

}

func (s *ServiceOrder) SendOrderRefundedConfirmation(order model.Order, user *model.User, _ any, amount decimal.Decimal, currency string, manager interfaces.PluginManagerInterface) *model_helper.AppError {
	panic("not implemented")
}

func (s *ServiceOrder) SendFulfillmentUpdate(order *model.Order, fulfillment *model.Fulfillment, manager interfaces.PluginManagerInterface) *model_helper.AppError {
	fulfillmentPayload, appErr := s.getDefaultFulfillmentPayload(order, fulfillment)
	if appErr != nil {
		return appErr
	}

	_, appErr = manager.Notify(model.ORDER_FULFILLMENT_UPDATE, fulfillmentPayload, order.ChannelID, "")
	return appErr
}
