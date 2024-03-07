package product

import (
	"database/sql"
	"fmt"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SqlProductVariantChannelListingStore struct {
	store.Store
}

func NewSqlProductVariantChannelListingStore(s store.Store) store.ProductVariantChannelListingStore {
	return &SqlProductVariantChannelListingStore{s}
}

func (ps *SqlProductVariantChannelListingStore) Upsert(transaction boil.ContextTransactor, variantChannelListings model.ProductVariantChannelListingSlice) (model.ProductVariantChannelListingSlice, error) {
	if transaction == nil {
		transaction = ps.GetMaster()
	}

	for _, listing := range variantChannelListings {
		if listing == nil {
			continue
		}

		isSaving := listing.ID == ""
		if isSaving {
			model_helper.ProductVariantChannelListingPreSave(listing)
		} else {
			model_helper.ProductVariantChannelListingCommonPre(listing)
		}

		if err := model_helper.ProductVariantChannelListingIsValid(*listing); err != nil {
			return nil, err
		}

		var err error
		if isSaving {
			err = listing.Insert(transaction, boil.Infer())
		} else {
			_, err = listing.Update(transaction, boil.Blacklist(
				model.ProductVariantChannelListingColumns.CreatedAt,
			))
		}

		if err != nil {
			if ps.IsUniqueConstraintError(err, []string{"product_variant_channel_listings_variant_id_channel_id_key", model.ProductVariantChannelListingColumns.VariantID}) {
				return nil, store.NewErrInvalidInput(model.TableNames.ProductVariantChannelListings, "VariantID/ChannelID", "duplicate")
			}
			return nil, err
		}
	}

	return variantChannelListings, nil
}

func (ps *SqlProductVariantChannelListingStore) Get(variantChannelListingID string) (*model.ProductVariantChannelListing, error) {
	listing, err := model.FindProductVariantChannelListing(ps.GetReplica(), variantChannelListingID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.ProductVariantChannelListings, variantChannelListingID)
		}
		return nil, err
	}

	return listing, nil
}

func (ps *SqlProductVariantChannelListingStore) commonQueryBuilder(option model_helper.ProductVariantChannelListingFilterOption) []qm.QueryMod {
	conds := option.Conditions
	conds = append(conds, qm.Select(model.TableNames.ProductVariantChannelListings+".*"))

	for _, load := range option.Preloads {
		conds = append(conds, qm.Load(load))
	}
	if option.VariantProductID != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ProductVariants, model.ProductVariantTableColumns.ID, model.ProductVariantChannelListingTableColumns.VariantID)),
			option.VariantProductID,
		)
	}

	var annotations = model_helper.AnnotationAggregator{}
	if option.AnnotateAvailablePreorderQuantity {
		annotations[model_helper.ProductVariantChannelListingAnnotationKeys.AvailablePreorderQuantity] = fmt.Sprintf("%s - COALESCE( SUM( %s ), 0)", model.ProductVariantChannelListingTableColumns.PreorderQuantityThreshold, model.PreorderAllocationTableColumns.Quantity)
	}
	if option.AnnotatePreorderQuantityAllocated {
		annotations[model_helper.ProductVariantChannelListingAnnotationKeys.PreorderQuantityAllocated] = fmt.Sprintf("COALESCE( SUM( %s ), 0)", model.PreorderAllocationTableColumns.Quantity)
	}
	conds = append(conds, annotations)

	return conds
}

func (ps *SqlProductVariantChannelListingStore) FilterbyOption(option model_helper.ProductVariantChannelListingFilterOption) (model.ProductVariantChannelListingSlice, error) {
	conds := ps.commonQueryBuilder(option)
	return model.ProductVariantChannelListings(conds...).All(ps.GetReplica())
}
