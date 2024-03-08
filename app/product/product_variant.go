package product

import (
	"net/http"

	"github.com/mattermost/squirrel"
	"github.com/samber/lo"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// ProductVariantById finds product variant by given id
func (a *ServiceProduct) ProductVariantById(id string) (*model.ProductVariant, *model_helper.AppError) {
	variant, err := a.srv.Store.ProductVariant().Get(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model_helper.NewAppError("ProductVariantById", "app.product.product_variant_missing.app_error", nil, err.Error(), statusCode)
	}

	return variant, nil
}

// ProductVariantGetPrice returns price
func (a *ServiceProduct) ProductVariantGetPrice(
	productVariant *model.ProductVariant,
	product model.Product,
	collections []*model.Collection,
	channel model.Channel,
	channelListing *model.ProductVariantChannelListing,
	discounts []*model_helper.DiscountInfo, // optional
) (*goprices.Money, *model_helper.AppError) {
	return a.srv.DiscountService().CalculateDiscountedPrice(product, channelListing.Price, collections, discounts, channel, productVariant.Id)
}

// ProductVariantIsDigital finds product type that related to given product variant and check if that product type is digital and does not require shipping
func (a *ServiceProduct) ProductVariantIsDigital(productVariantID string) (bool, *model_helper.AppError) {
	productType, err := a.srv.Store.ProductType().ProductTypeByProductVariantID(productVariantID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return false, model_helper.NewAppError("ProductVariantIsDigital", "app.product.product_type_by_product_variant_id.app_error", nil, err.Error(), statusCode)
	}

	return *productType.IsDigital && !*productType.IsShippingRequired, nil
}

// ProductVariantByOrderLineID returns a product variant by given order line id
func (a *ServiceProduct) ProductVariantByOrderLineID(orderLineID string) (*model.ProductVariant, *model_helper.AppError) {
	productVariant, err := a.srv.Store.ProductVariant().GetByOrderLineID(orderLineID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model_helper.NewAppError("ProductVariantByOrderLineID", "app.product.error_finding_product_variant_by_order_line_id.app_error", nil, err.Error(), statusCode)
	}

	return productVariant, nil
}

// ProductVariantsByOption returns a list of product variants satisfy given option
func (a *ServiceProduct) ProductVariantsByOption(option *model.ProductVariantFilterOption) (model.ProductVariants, *model_helper.AppError) {
	productVariants, err := a.srv.Store.ProductVariant().FilterByOption(option)
	if err != nil {
		return nil, model_helper.NewAppError("ProductVariantsByOption", "app.product.error_finding_product_variants_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return productVariants, nil
}

// ProductVariantGetWeight returns weight of given product variant
func (a *ServiceProduct) ProductVariantGetWeight(productVariantID string) (*measurement.Weight, *model_helper.AppError) {
	weight, err := a.srv.Store.ProductVariant().GetWeight(productVariantID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model_helper.NewAppError("ProductVariantGetWeight", "app.product.error_getting_product_variant_weight.app_error", nil, err.Error(), statusCode)
	}

	return weight, nil
}

// DisplayProduct return display text for given product variant
//
// `translated` default to false
func (a *ServiceProduct) DisplayProduct(productVariant *model.ProductVariant, translated bool) (stringm *model_helper.AppError) {
	panic("not implt")
}

// ProductVariantsAvailableInChannel returns product variants based on given channel slug
func (a *ServiceProduct) ProductVariantsAvailableInChannel(channelSlug string) ([]*model.ProductVariant, *model_helper.AppError) {
	productVariants, appErr := a.ProductVariantsByOption(&model.ProductVariantFilterOption{
		RelatedProductVariantChannelListingConditions: squirrel.NotEq{model.ProductVariantChannelListingTableName + "." + model.ProductVariantChannelListingColumnPriceAmount: nil},
		ProductVariantChannelListingChannelSlug:       squirrel.Eq{model.ChannelTableName + "." + model.ChannelColumnSlug: channelSlug},
	})

	if appErr != nil {
		return nil, appErr
	}

	return productVariants, nil
}

// UpsertProductVariant tells store to upsert given product variant and returns it
func (s *ServiceProduct) UpsertProductVariant(transaction boil.ContextTransactor, variant *model.ProductVariant) (*model.ProductVariant, *model_helper.AppError) {
	upsertedVariant, err := s.srv.Store.ProductVariant().Save(transaction, variant)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}
		var statusCode = http.StatusInternalServerError

		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		} else if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}
		return nil, model_helper.NewAppError("UpsertProductVariant", "app.product.error_upserting_product_variant.app_error", nil, err.Error(), statusCode)
	}

	return upsertedVariant, nil
}

func (s *ServiceProduct) DeleteProductVariants(variantIds []string, requesterID string) (int64, *model_helper.AppError) {
	// find all draft order lines related to given variants
	orderLines, appErr := s.srv.OrderService().OrderLinesByOption(&model.OrderLineFilterOption{
		Conditions:             squirrel.Eq{model.OrderLineTableName + "." + model.OrderLineColumnVariantID: variantIds},
		RelatedOrderConditions: squirrel.Eq{model.OrderTableName + "." + model.OrderColumnStatus: model.ORDER_STATUS_DRAFT},
		Preload:                []string{"Order"},
	})
	if appErr != nil {
		return 0, appErr
	}

	// begin tx
	tx := s.srv.Store.GetMaster().Begin()
	if tx.Error != nil {
		return 0, model_helper.NewAppError("DeleteProductVariants", model_helper.ErrorCreatingTransactionErrorID, nil, tx.Error.Error(), http.StatusInternalServerError)
	}

	// create order events on order lines
	orderOrderLinesMap := map[string]model.OrderLineSlice{}
	orders := model.Orders{}
	for _, line := range orderLines {
		_, exist := orderOrderLinesMap[line.OrderID]
		if !exist {
			orders = append(orders, line.Order)
		}
		orderOrderLinesMap[line.OrderID] = append(orderOrderLinesMap[line.OrderID], line)
	}

	for orderID, orderLines := range orderOrderLinesMap {
		quantityOrderLines := lo.Map(orderLines, func(item *model.OrderLine, _ int) *model.QuantityOrderLine {
			return &model.QuantityOrderLine{Quantity: item.Quantity, OrderLine: item}
		})

		_, appErr = s.srv.OrderService().CommonCreateOrderEvent(tx, &model.OrderEventOption{
			OrderID: orderID,
			UserID:  &requesterID,
			Type:    model.ORDER_EVENT_TYPE_ORDER_LINE_VARIANT_DELETED,
			Parameters: model_types.JSONString{
				"lines": s.srv.OrderService().LinesPerQuantityToLineObjectList(quantityOrderLines),
			},
		})
		if appErr != nil {
			return 0, appErr
		}
	}

	// actually delete variants, related draft order lines and related attribute values
	numDeleted, err := s.srv.Store.ProductVariant().Delete(tx, variantIds)
	if err != nil {
		return 0, model_helper.NewAppError("DeleteProductVariants", "app.product.error_deleting_variants.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	var finishTransaction = func() *model_helper.AppError {
		// commit
		defer s.srv.Store.FinalizeTransaction(tx)
		err = tx.Commit().Error
		if err != nil {
			return model_helper.NewAppError("DeleteProductVariants", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
		}
		return nil
	}

	// perform recalculate orders
	if len(orders) > 0 {
		s.srv.Go(func() {
			for _, order := range orders {
				appErr := s.srv.OrderService().RecalculateOrder(tx, order, model_types.JSONString{})
				if appErr != nil {
					slog.Error("failed to recalculate order after deleting product variants", slog.String("orderID", order.Id), slog.Err(appErr))
				}
			}

			appErr := finishTransaction()
			if appErr != nil {
				slog.Error("failed to finish transaction", slog.Err(appErr))
			}
		})
	} else {
		appErr = finishTransaction()
		if appErr != nil {
			return 0, appErr
		}
	}

	pluginMng := s.srv.PluginService().GetPluginManager()
	for _, variantID := range variantIds {
		_, appErr = pluginMng.ProductVariantDeleted(model.ProductVariant{Id: variantID})
		if appErr != nil {
			return 0, appErr
		}
	}

	return numDeleted, nil
}

func (s *ServiceProduct) ToggleVariantRelations(variants model.ProductVariants, medias model.ProductMedias, sales model.Sales, vouchers model.Vouchers, wishlistItems model.WishlistItems, isDelete bool) *model_helper.AppError {
	// create tx:
	tx := s.srv.Store.GetMaster().Begin()
	if tx.Error != nil {
		return model_helper.NewAppError("ToggleVariantRelations", model_helper.ErrorCreatingTransactionErrorID, nil, tx.Error.Error(), http.StatusInternalServerError)
	}
	defer s.srv.Store.FinalizeTransaction(tx)

	err := s.srv.Store.
		ProductVariant().
		ToggleProductVariantRelations(
			tx,
			variants,
			medias,
			sales,
			vouchers,
			wishlistItems,
			isDelete,
		)
	if err != nil {
		return model_helper.NewAppError("ToggleVariantRelations", "app.product.toggle_variant_relations.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// commit tx
	err = tx.Commit().Error
	if err != nil {
		return model_helper.NewAppError("ToggleVariantRelations", model_helper.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	pluginMng := s.srv.PluginService().GetPluginManager()
	for _, variant := range variants {
		_, appErr := pluginMng.ProductVariantUpdated(*variant)
		if appErr != nil {
			return appErr
		}
	}

	return nil
}

// func (s *ServiceProduct) GetProductVariantsForRequester(id, sku, channelIdOrSlug string, requesterIsShopStaff bool) (model.ProductVariants, *model_helper.AppError) {
// 	s.srv.AccountService().Per
// }
