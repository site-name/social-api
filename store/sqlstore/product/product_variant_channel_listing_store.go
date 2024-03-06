package product

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlProductVariantChannelListingStore struct {
	store.Store
}

func NewSqlProductVariantChannelListingStore(s store.Store) store.ProductVariantChannelListingStore {
	return &SqlProductVariantChannelListingStore{s}
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

func (ps *SqlProductVariantChannelListingStore) FilterbyOption(option model_helper.ProductVariantChannelListingFilterOption) (model.ProductVariantChannelListingSlice, error) {
	// NOTE: In the scan fields creation below, the order of fields must be identical to the order of select fiels
	selectFields := []string{model.ProductVariantChannelListingTableName + ".*"}
	if option.SelectRelatedChannel {
		selectFields = append(selectFields, model.ChannelTableName+".*")
	}
	if option.SelectRelatedProductVariant {
		selectFields = append(selectFields, model.ProductVariantTableName+".*")
	}
	if option.AnnotateAvailablePreorderQuantity {
		selectFields = append(selectFields, `ProductVariantChannelListings.PreorderQuantityThreshold - COALESCE( SUM( PreorderAllocations.Quantity ), 0) AS availablePreorderQuantity`)
	}
	if option.AnnotatePreorderQuantityAllocated {
		selectFields = append(selectFields, `COALESCE( SUM( PreorderAllocations.Quantity ), 0) AS preorderQuantityAllocated`)
	}

	query := ps.GetQueryBuilder().
		Select(selectFields...).
		From(model.ProductVariantChannelListingTableName).Where(option.Conditions).Where(option.VariantProductID)

	var groupBy []string

	// parse option
	if option.SelectForUpdate && option.Transaction != nil {
		var forUpdateOf string
		if option.SelectForUpdateOf != "" {
			forUpdateOf = " OF " + option.SelectForUpdateOf
		}
		query = query.Suffix("FOR UPDATE" + forUpdateOf)
	}

	if option.SelectRelatedChannel {
		query = query.
			InnerJoin(model.ChannelTableName + " ON Channels.Id = ProductVariantChannelListings.ChannelID")
		groupBy = append(groupBy, "Channels.Id")
	}
	if option.SelectRelatedProductVariant || option.VariantProductID != nil {
		query = query.
			InnerJoin(model.ProductVariantTableName + " ON ProductVariants.Id = ProductVariantChannelListings.variantID")
	}

	if option.AnnotateAvailablePreorderQuantity ||
		option.AnnotatePreorderQuantityAllocated {
		query = query.LeftJoin(model.PreOrderAllocationTableName + " ON PreorderAllocations.ProductVariantChannelListingID = ProductVariantChannelListings.Id")
		groupBy = append(groupBy, "ProductVariantChannelListings.Id")
	}

	if len(groupBy) > 0 {
		query = query.GroupBy(groupBy...)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterbyOption_ToSql")
	}

	rows, err := ps.GetReplica().Raw(queryString, args...).Rows()
	if err != nil {
		return nil, errors.Wrap(err, "failed to find product variant channel listings")
	}
	defer rows.Close()

	var res model.ProductVariantChannelListings

	for rows.Next() {
		var (
			channel                   model.Channel
			variantChannelListing     model.ProductVariantChannelListing
			availablePreorderQuantity int
			preorderQuantityAllocated int
			scanFields                = ps.ScanFields(&variantChannelListing) // order of fields must be identical to select fields above
			variant                   model.ProductVariant
		)
		if option.SelectRelatedChannel {
			scanFields = append(scanFields, ps.Channel().ScanFields(&channel)...)
		}
		if option.SelectRelatedProductVariant {
			scanFields = append(scanFields, ps.ProductVariant().ScanFields(&variant)...)
		}
		if option.AnnotateAvailablePreorderQuantity {
			scanFields = append(scanFields, &availablePreorderQuantity)
		}
		if option.AnnotatePreorderQuantityAllocated {
			scanFields = append(scanFields, &preorderQuantityAllocated)
		}

		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row of product variant channel listing")
		}

		if option.SelectRelatedChannel {
			variantChannelListing.SetChannel(&channel)
		}
		if option.SelectRelatedProductVariant {
			variantChannelListing.SetVariant(&variant)
		}
		if option.AnnotateAvailablePreorderQuantity {
			variantChannelListing.Set_availablePreorderQuantity(availablePreorderQuantity)
		}
		if option.AnnotatePreorderQuantityAllocated {
			variantChannelListing.Set_preorderQuantityAllocated(preorderQuantityAllocated)
		}
		res = append(res, &variantChannelListing)
	}

	return res, nil
}

func (ps *SqlProductVariantChannelListingStore) Upsert(transaction boil.ContextTransactor, variantChannelListings model.ProductVariantChannelListingSlice) (model.ProductVariantChannelListingSlice, error) {
	if transaction == nil {
		transaction = ps.GetMaster()
	}

	for _, listing := range variantChannelListings {
		err := transaction.Save(listing).Error

		if err != nil {
			if ps.IsUniqueConstraintError(err, []string{"VariantID", "ChannelID", "productvariantchannellistings_variantid_channelid_key"}) {
				return nil, store.NewErrInvalidInput(model.ProductVariantChannelListingTableName, "VariantID/ChannelID", "duplicate")
			}
			return nil, errors.Wrapf(err, "failed to upsert a product variant channel listing with id=%s", listing.Id)
		}
	}

	return variantChannelListings, nil
}
