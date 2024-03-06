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

type SqlProductChannelListingStore struct {
	store.Store
}

func NewSqlProductChannelListingStore(s store.Store) store.ProductChannelListingStore {
	return &SqlProductChannelListingStore{s}
}

func (ps *SqlProductChannelListingStore) Upsert(transaction boil.ContextTransactor, listings model.ProductChannelListingSlice) (model.ProductChannelListingSlice, error) {
	if transaction == nil {
		transaction = ps.GetMaster()
	}

	for _, listing := range listings {
		if listing == nil {
			continue
		}

		isSaving := listing.ID == ""
		if isSaving {
			model_helper.ProductChannelListingPreSave(listing)
		} else {
			model_helper.ProductChannelListingCommonPre(listing)
		}

		if err := model_helper.ProductChannelListingIsValid(*listing); err != nil {
			return nil, err
		}

		var err error
		if isSaving {
			err = listing.Insert(transaction, boil.Infer())
		} else {
			_, err = listing.Update(transaction, boil.Blacklist(
				model.ProductChannelListingColumns.CreatedAt,
			))
		}

		if err != nil {
			if ps.IsUniqueConstraintError(err, []string{model.ProductChannelListingTableColumns.ProductID, model.ProductChannelListingColumns.ChannelID, "product_channel_listings_product_id_channel_id_key"}) {
				return nil, store.NewErrInvalidInput(model.TableNames.ProductChannelListings, "ProductID/ChannelID", "duplicate")
			}
			return nil, err
		}
	}

	return listings, nil
}

func (ps *SqlProductChannelListingStore) Get(listingID string) (*model.ProductChannelListing, error) {
	listing, err := model.FindProductChannelListing(ps.GetReplica(), listingID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.ProductChannelListings, listingID)
		}
		return nil, err
	}

	return listing, nil
}

func (ps *SqlProductChannelListingStore) FilterByOption(option model_helper.ProductChannelListingFilterOption) (model.ProductChannelListingSlice, error) {
	conds := option.Conditions
	if option.RelatedChannelConditions != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Channels, model.ChannelTableColumns.ID, model.ProductChannelListingTableColumns.ChannelID)),
			option.RelatedChannelConditions,
		)
	}
	if option.ProductVariantID != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Products, model.ProductTableColumns.ID, model.ProductChannelListingTableColumns.ProductID)),
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ProductVariants, model.ProductVariantTableColumns.ProductID, model.ProductTableColumns.ID)),
			option.ProductVariantID,
		)
	}
	for _, load := range option.Preloads {
		conds = append(conds, qm.Load(load))
	}

	return model.ProductChannelListings(conds...).All(ps.GetReplica())
}
