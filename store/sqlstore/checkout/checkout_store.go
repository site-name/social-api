package checkout

import (
	"database/sql"
	"fmt"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SqlCheckoutStore struct {
	store.Store
}

func NewSqlCheckoutStore(sqlStore store.Store) store.CheckoutStore {
	return &SqlCheckoutStore{sqlStore}
}

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

func (cs *SqlCheckoutStore) FilterByOption(option model_helper.CheckoutFilterOptions) (model.CheckoutSlice, error) {
	return model.Checkouts(option.Conditions...).All(cs.GetReplica())
}

func (cs *SqlCheckoutStore) FetchCheckoutLinesAndPrefetchRelatedValue(checkout model.Checkout) (model_helper.CheckoutLineInfos, error) {
	// please refer to file checkout_store_sql.md for details
	checkoutLines, err := model.CheckoutLines(
		model.CheckoutLineWhere.CheckoutID.EQ(checkout.Token),
		qm.Load(fmt.Sprintf(
			"%s.%s.%s.%s",
			model.CheckoutLineRels.Variant,
			model.ProductVariantRels.Product,
			model.ProductRels.ProductCollections,
			model.ProductCollectionRels.Collection,
		)),
		qm.Load(fmt.Sprintf(
			"%s.%s.%s",
			model.CheckoutLineRels.Variant,
			model.ProductVariantRels.VariantProductVariantChannelListings,
			model.ProductVariantChannelListingRels.Channel,
		)),
	).All(cs.GetReplica())
	if err != nil {
		return nil, err
	}

	var result model_helper.CheckoutLineInfos

	for _, checkoutLine := range checkoutLines {
		if checkoutLine == nil {
			continue
		}

		checkoutLineInfo := &model_helper.CheckoutLineInfo{
			Line: *checkoutLine,
		}

		if checkoutLine.R != nil {
			productVariant := checkoutLine.R.Variant

			if productVariant != nil {
				checkoutLineInfo.Variant = *productVariant //

				if productVariant.R != nil {
					product := productVariant.R.Product

					if product != nil {
						checkoutLineInfo.Product = *product //

						if product.R != nil {
							productCollections := product.R.ProductCollections

							if len(productCollections) > 0 {
								var collections model.CollectionSlice

								for _, prdCl := range productCollections {
									if prdCl.R != nil && prdCl.R.Collection != nil {
										collections = append(collections, prdCl.R.Collection)
									}
								}

								checkoutLineInfo.Collections = collections //
							}
						}
					}

					var productVariantChannelListing *model.ProductVariantChannelListing = nil

					for _, listing := range productVariant.R.VariantProductVariantChannelListings {
						if listing != nil && listing.ChannelID == checkout.ChannelID {
							productVariantChannelListing = listing
						}
					}
					if productVariantChannelListing == nil {
						continue
					}

					checkoutLineInfo.ChannelListing = *productVariantChannelListing //

					result = append(result, checkoutLineInfo)
				}
			}
		}
	}

	return result, nil
}

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
