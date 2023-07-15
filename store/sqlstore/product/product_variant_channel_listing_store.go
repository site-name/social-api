package product

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

type SqlProductVariantChannelListingStore struct {
	store.Store
}

func NewSqlProductVariantChannelListingStore(s store.Store) store.ProductVariantChannelListingStore {
	return &SqlProductVariantChannelListingStore{s}
}

func (ps *SqlProductVariantChannelListingStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"VariantID",
		"ChannelID",
		"Currency",
		"PriceAmount",
		"CostPriceAmount",
		"PreorderQuantityThreshold",
		"CreateAt",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (ps *SqlProductVariantChannelListingStore) ScanFields(listing *model.ProductVariantChannelListing) []interface{} {
	return []interface{}{
		&listing.Id,
		&listing.VariantID,
		&listing.ChannelID,
		&listing.Currency,
		&listing.PriceAmount,
		&listing.CostPriceAmount,
		&listing.PreorderQuantityThreshold,
		&listing.CreateAt,
	}
}

// Save insert given value into database then returns it with an error
func (ps *SqlProductVariantChannelListingStore) Save(variantChannelListing *model.ProductVariantChannelListing) (*model.ProductVariantChannelListing, error) {
	variantChannelListing.PreSave()
	if err := variantChannelListing.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + model.ProductVariantChannelListingTableName + "(" + ps.ModelFields("").Join(",") + ") VALUES (" + ps.ModelFields(":").Join(",") + ")"
	_, err := ps.GetMasterX().NamedExec(query, variantChannelListing)
	if err != nil {
		if ps.IsUniqueConstraintError(err, []string{"VariantID", "ChannelID", "productvariantchannellistings_variantid_channelid_key"}) {
			return nil, store.NewErrNotFound(model.ProductVariantChannelListingTableName, variantChannelListing.Id)
		}
		return nil, errors.Wrapf(err, "failed to save product variant channel listing with id=%s", variantChannelListing.Id)
	}

	return variantChannelListing, nil
}

// Get finds and returns 1 product variant channel listing based on given variantChannelListingID
func (ps *SqlProductVariantChannelListingStore) Get(variantChannelListingID string) (*model.ProductVariantChannelListing, error) {
	var res model.ProductVariantChannelListing

	err := ps.GetReplicaX().Get(&res, "SELECT * FROM "+model.ProductVariantChannelListingTableName+" WHERE Id = ?", variantChannelListingID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.ProductVariantChannelListingTableName, variantChannelListingID)
		}
		return nil, errors.Wrapf(err, "failed to find product variant channel listing with id=%s", variantChannelListingID)
	}

	return &res, nil
}

// FilterbyOption finds and returns all product variant channel listings filterd using given option
func (ps *SqlProductVariantChannelListingStore) FilterbyOption(option *model.ProductVariantChannelListingFilterOption) ([]*model.ProductVariantChannelListing, error) {
	// NOTE: In the scan fields creation below, the order of fields must be identical to the order of select fiels
	selectFields := ps.ModelFields(model.ProductVariantChannelListingTableName + ".")
	if option.SelectRelatedChannel {
		selectFields = append(selectFields, ps.Channel().ModelFields(model.ChannelTableName+".")...)
	}
	if option.SelectRelatedProductVariant {
		selectFields = append(selectFields, ps.ProductVariant().ModelFields(model.ProductVariantTableName+".")...)
	}
	if option.AnnotateAvailablePreorderQuantity {
		selectFields = append(selectFields, `ProductVariantChannelListings.PreorderQuantityThreshold - COALESCE( SUM( PreorderAllocations.Quantity ), 0) AS availablePreorderQuantity`)
	}
	if option.AnnotatePreorderQuantityAllocated {
		selectFields = append(selectFields, `COALESCE( SUM( PreorderAllocations.Quantity ), 0) AS preorderQuantityAllocated`)
	}

	query := ps.GetQueryBuilder().
		Select(selectFields...).
		From(model.ProductVariantChannelListingTableName)

	var groupBy []string

	// parse option
	if option.SelectForUpdate {
		var forUpdateOf string
		if option.SelectForUpdateOf != "" {
			forUpdateOf = " OF " + option.SelectForUpdateOf
		}
		query = query.Suffix("FOR UPDATE" + forUpdateOf)
	}

	for _, opt := range []squirrel.Sqlizer{
		option.Id,
		option.VariantID,
		option.ChannelID,
		option.PriceAmount,
		option.VariantProductID,
	} {
		if opt != nil {
			query = query.Where(opt)
		}
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

	rows, err := ps.GetReplicaX().QueryX(queryString, args...)
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

// BulkUpsert performs bulk upsert given product variant channel listings then returns them
func (ps *SqlProductVariantChannelListingStore) BulkUpsert(transaction store_iface.SqlxExecutor, variantChannelListings []*model.ProductVariantChannelListing) ([]*model.ProductVariantChannelListing, error) {
	var (
		executor    store_iface.SqlxExecutor = ps.GetMasterX()
		saveQuery                            = "INSERT INTO " + model.ProductVariantChannelListingTableName + "(" + ps.ModelFields("").Join(",") + ") VALUES (" + ps.ModelFields(":").Join(",") + ")"
		updateQuery                          = "UPDATE " + model.ProductVariantChannelListingTableName + " SET " + ps.
				ModelFields("").
				Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"
	)
	if transaction != nil {
		executor = transaction
	}

	for _, listing := range variantChannelListings {
		var (
			isSaving   bool
			numUpdated int64
			err        error
		)

		if !model.IsValidId(listing.Id) {
			listing.PreSave()
			isSaving = true
		}

		if err := listing.IsValid(); err != nil {
			return nil, err
		}

		if isSaving {
			_, err = executor.NamedExec(saveQuery, listing)

		} else {
			var result sql.Result
			result, err = executor.NamedExec(updateQuery, listing)
			if err == nil && result != nil {
				numUpdated, _ = result.RowsAffected()
			}
		}

		if err != nil {
			if ps.IsUniqueConstraintError(err, []string{"VariantID", "ChannelID", "productvariantchannellistings_variantid_channelid_key"}) {
				return nil, store.NewErrInvalidInput(model.ProductVariantChannelListingTableName, "VariantID/ChannelID", "duplicate")
			}
			return nil, errors.Wrapf(err, "failed to upsert a product variant channel listing with id=%s", listing.Id)
		}
		if numUpdated > 1 {
			return nil, errors.Errorf("%d product variant channel listing(s) with id=%s was/were updated instead of 1", numUpdated, listing.Id)
		}
	}

	return variantChannelListings, nil
}
