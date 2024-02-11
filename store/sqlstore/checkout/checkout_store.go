package checkout

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlCheckoutStore struct {
	store.Store
}

func NewSqlCheckoutStore(sqlStore store.Store) store.CheckoutStore {
	return &SqlCheckoutStore{sqlStore}
}

// Upsert depends on given checkout's Token property to decide to update or insert it
func (cs *SqlCheckoutStore) Upsert(tx boil.ContextTransactor, checkouts model.CheckoutSlice) (model.CheckoutSlice, error) {
	if tx == nil {
		tx = cs.GetMaster()
	}

	for _, checkout := range checkouts {
		if checkout == nil {
			continue
		}

		var isSaving bool

		if checkout.Token == "" {
			model_helper.CheckoutPreSave(checkout)
			isSaving = true
		} else {
			model_helper.CheckoutPreUpdate(checkout)
		}

		if err := model_helper.CheckoutIsValid(*checkout); err != nil {
			return nil, err
		}

		var err error
		if isSaving {
			err = checkout.Insert(tx, boil.Infer())
		} else {
			_, err = checkout.Update(tx, boil.Blacklist(model.CheckoutColumns.CreatedAt))
		}

		if err != nil {
			return nil, err
		}
	}

	return checkouts, nil
}

// GetByOption finds and returns 1 checkout based on given option
func (cs *SqlCheckoutStore) GetByOption(option model_helper.CheckoutFilterOptions) (*model.Checkout, error) {
	checkout, err := model.Checkouts(option.Conditions...).One(cs.GetReplica())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.Checkouts, "options")
		}
		return nil, err
	}

	return checkout, nil
}

// FilterByOption finds and returns a list of checkout based on given option
func (cs *SqlCheckoutStore) FilterByOption(option model_helper.CheckoutFilterOptions) (model.CheckoutSlice, error) {
	return model.Checkouts(option.Conditions...).All(cs.GetReplica())
}

// FetchCheckoutLinesAndPrefetchRelatedValue Fetch checkout lines as CheckoutLineInfo objects.
func (cs *SqlCheckoutStore) FetchCheckoutLinesAndPrefetchRelatedValue(checkout *model.Checkout) ([]*model.CheckoutLineInfo, error) {
	// please refer to file checkout_store_sql.md for details

	// fetch checkout lines:
	var checkoutLines model.CheckoutLines

	err := cs.GetReplica().Order("CreateAt ASC").Find(&checkoutLines, "CheckoutID = ?", checkout.Token).Error
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find checkout lines belong to checkout with token=%s", checkout.Token)
	}
	productVariantIDs := checkoutLines.VariantIDs()

	// fetch product variants
	var (
		productIDs        []string
		productVariantMap = map[string]*model.ProductVariant{} // productVariantMap has keys are product variant ids
	)
	// check if we can proceed:
	if len(productVariantIDs) > 0 {
		var productVariants model.ProductVariants
		err = cs.GetReplica().Find(&productVariants, "Id = ?", productVariantIDs).Error
		if err != nil {
			return nil, errors.Wrap(err, "failed to find product variants")
		}
		for _, variant := range productVariants {
			productIDs = append(productIDs, variant.ProductID)
			productVariantMap[variant.Id] = variant
		}
	}

	// fetch products
	var (
		products       model.Products
		productTypeIDs []string
		productMap     = map[string]*model.Product{} // productMap has keys are product ids
	)
	// check if we can proceed:
	if len(productIDs) > 0 {
		err = cs.GetReplica().Find(&products, "Id IN ?", productIDs).Error
		if err != nil {
			return nil, errors.Wrap(err, "failed to finds products")
		}
		for _, product := range products {
			productTypeIDs = append(productTypeIDs, product.ProductTypeID)
			productMap[product.Id] = product
		}
	}

	// fetch product collections
	var (
		collectionXs []*struct {
			model.Collection
			PrefetchRelatedValProductID string
		}
		collectionsByProducts = map[string]model.Collections{} // collectionsByProducts has keys are product ids
	)
	// check if we can proceed
	if len(productIDs) > 0 {
		err = cs.GetReplica().
			Table(model.CollectionTableName).
			Where("ProductCollections.ProductID IN ?", productIDs).
			Select("Collections.*", "ProductCollections.ProductID AS PrefetchRelatedValProductID").
			Joins("INNER JOIN " + model.CollectionProductRelationTableName + " ON ProductCollections.CollectionID = Collections.Id").
			Scan(&collectionXs).Error

		if err != nil {
			return nil, errors.Wrap(err, "failed to find collections")
		}
		for _, collectionX := range collectionXs {
			collectionsByProducts[collectionX.PrefetchRelatedValProductID] = append(collectionsByProducts[collectionX.PrefetchRelatedValProductID], &collectionX.Collection)
		}
	}

	// fetch product variant channel listing
	var (
		productVariantChannelListings                 []*model.ProductVariantChannelListing
		channelIDs                                    []string
		productVariantChannelListingsByProductVariant = map[string][]*model.ProductVariantChannelListing{} // productVariantChannelListingsByProductVariant has keys are product variant ids
	)
	// check if we can proceed:
	if len(productVariantIDs) > 0 {
		err := cs.GetReplica().Find(&productVariantChannelListings, "VariantID IN ?", productVariantIDs).Error
		if err != nil {
			return nil, errors.Wrap(err, "failed to find product variant channel listing")
		}
		for _, listing := range productVariantChannelListings {
			channelIDs = append(channelIDs, listing.ChannelID)
			productVariantChannelListingsByProductVariant[listing.VariantID] = append(productVariantChannelListingsByProductVariant[listing.VariantID], listing)
		}
	}

	// fetch channels
	var channels []*model.Channel
	// check if we can proceed
	if len(channelIDs) > 0 {
		err = cs.GetReplica().Find(&channels, "Id in ? ORDER BY Slug ASC", channelIDs).Error
		if err != nil {
			return nil, errors.Wrap(err, "failed to find channels")
		}
	}

	// fetch product types
	var (
		productTypes   []*model.ProductType
		productTypeMap = map[string]*model.ProductType{} // productTypeMap has keys are product type ids
	)
	// check if we can proceed
	if len(productTypeIDs) > 0 {
		err = cs.GetReplica().Find(&productTypes, "Id IN ?", productTypeIDs).Error
		if err != nil {
			return nil, errors.Wrap(err, "failed to finds product types")
		}
		for _, productType := range productTypes {
			productTypeMap[productType.Id] = productType
		}
	}

	var checkoutLineInfos []*model.CheckoutLineInfo

	for _, checkoutLine := range checkoutLines {
		productVariant := productVariantMap[checkoutLine.VariantID]

		if productVariant != nil {
			var variantChannelListing *model.ProductVariantChannelListing
			for _, listing := range productVariantChannelListingsByProductVariant[productVariant.Id] {
				if listing.ChannelID == checkout.ChannelID {
					variantChannelListing = listing
				}
			}

			// FIXME: Temporary solution to pass type checks. Figure out how to handle case
			// when variant channel listing is not defined for a checkout line.
			if variantChannelListing == nil {
				continue
			}

			product := productMap[productVariant.ProductID]
			if product != nil {
				productType := productTypeMap[product.ProductTypeID]
				collections := collectionsByProducts[product.Id]

				if productType != nil && collections != nil {
					checkoutLineInfos = append(checkoutLineInfos, &model.CheckoutLineInfo{
						Line:           *checkoutLine,
						Variant:        *productVariant,
						ChannelListing: *variantChannelListing,
						Product:        *product,
						ProductType:    *productType,
						Collections:    collections,
					})
				}
			}
		}
	}

	return checkoutLineInfos, nil
}

// DeleteCheckoutsByOption deletes checkout row(s) from database, filtered using given option.
// It returns an error indicating if the operation was performed successfully.
func (cs *SqlCheckoutStore) Delete(transaction boil.ContextTransactor, tokens []string) error {
	if transaction == nil {
		transaction = cs.GetMaster()
	}

	_, err := model.Checkouts(model.CheckoutWhere.Token.IN(tokens)).DeleteAll(transaction)
	return err
}

func (cs *SqlCheckoutStore) CountCheckouts(options model_helper.CheckoutFilterOptions) (int64, error) {
	return model.Checkouts(options.Conditions...).Count(cs.GetReplica())
}
