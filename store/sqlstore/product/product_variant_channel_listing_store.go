package product

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

type SqlProductVariantChannelListingStore struct {
	store.Store
}

func NewSqlProductVariantChannelListingStore(s store.Store) store.ProductVariantChannelListingStore {
	return &SqlProductVariantChannelListingStore{s}
}

func (ps *SqlProductVariantChannelListingStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
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

func (ps *SqlProductVariantChannelListingStore) ScanFields(listing product_and_discount.ProductVariantChannelListing) []interface{} {
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
func (ps *SqlProductVariantChannelListingStore) Save(variantChannelListing *product_and_discount.ProductVariantChannelListing) (*product_and_discount.ProductVariantChannelListing, error) {
	variantChannelListing.PreSave()
	if err := variantChannelListing.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.ProductVariantChannelListingTableName + "(" + ps.ModelFields("").Join(",") + ") VALUES (" + ps.ModelFields(":").Join(",") + ")"
	_, err := ps.GetMasterX().NamedExec(query, variantChannelListing)
	if err != nil {
		if ps.IsUniqueConstraintError(err, []string{"VariantID", "ChannelID", "productvariantchannellistings_variantid_channelid_key"}) {
			return nil, store.NewErrNotFound(store.ProductVariantChannelListingTableName, variantChannelListing.Id)
		}
		return nil, errors.Wrapf(err, "failed to save product variant channel listing with id=%s", variantChannelListing.Id)
	}

	return variantChannelListing, nil
}

// Get finds and returns 1 product variant channel listing based on given variantChannelListingID
func (ps *SqlProductVariantChannelListingStore) Get(variantChannelListingID string) (*product_and_discount.ProductVariantChannelListing, error) {
	var res product_and_discount.ProductVariantChannelListing

	err := ps.GetReplicaX().Get(&res, "SELECT * FROM "+store.ProductVariantChannelListingTableName+" WHERE Id = ?", variantChannelListingID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ProductVariantChannelListingTableName, variantChannelListingID)
		}
		return nil, errors.Wrapf(err, "failed to find product variant channel listing with id=%s", variantChannelListingID)
	}

	return &res, nil
}

// FilterbyOption finds and returns all product variant channel listings filterd using given option
func (ps *SqlProductVariantChannelListingStore) FilterbyOption(transaction store_iface.SqlxTxExecutor, option *product_and_discount.ProductVariantChannelListingFilterOption) ([]*product_and_discount.ProductVariantChannelListing, error) {
	var runner store_iface.SqlxExecutor = ps.GetReplicaX()
	if transaction != nil {
		runner = transaction
	}

	// NOTE: In the scan fields creation below, the order of fields must be identical to the order of select fiels
	selectFields := ps.ModelFields(store.ProductVariantChannelListingTableName + ".")
	if option.SelectRelatedChannel {
		selectFields = append(selectFields, ps.Channel().ModelFields(store.ChannelTableName+".")...)
	}
	if option.AnnotateAvailablePreorderQuantity {
		selectFields = append(selectFields, `ProductVariantChannelListings.PreorderQuantityThreshold - COALESCE( SUM( PreorderAllocations.Quantity ), 0) AS availablePreorderQuantity`)
	}
	if option.AnnotatePreorderQuantityAllocated {
		selectFields = append(selectFields, `COALESCE( SUM( PreorderAllocations.Quantity ), 0) AS preorderQuantityAllocated`)
	}

	query := ps.GetQueryBuilder().
		Select(selectFields...).
		From(store.ProductVariantChannelListingTableName).
		OrderBy(store.TableOrderingMap[store.ProductVariantChannelListingTableName])

	// parse option
	if option.SelectForUpdate {
		var forUpdateOf string
		if option.SelectForUpdateOf != "" {
			forUpdateOf = " OF " + option.SelectForUpdateOf
		}
		query = query.Suffix("FOR UPDATE" + forUpdateOf)
	}
	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.VariantID != nil {
		query = query.Where(option.VariantID)
	}
	if option.ChannelID != nil {
		query = query.Where(option.ChannelID)
	}
	if option.PriceAmount != nil {
		query = query.Where(option.PriceAmount)
	}
	if option.VariantProductID != nil {
		query = query.
			InnerJoin(store.ProductVariantTableName + " ON (ProductVariants.Id = ProductVariantChannelListings.variantID)").
			Where(option.VariantProductID)
	}
	if option.SelectRelatedChannel {
		query = query.InnerJoin(store.ChannelTableName + " ON (Channels.Id = ProductVariants.ChannelID)")
	}

	var groupBy []string

	if option.AnnotateAvailablePreorderQuantity || option.AnnotatePreorderQuantityAllocated {
		query = query.LeftJoin(store.PreOrderAllocationTableName + " ON (PreorderAllocations.ProductVariantChannelListingID = ProductVariantChannelListings.Id)")
		groupBy = append(groupBy, "ProductVariantChannelListings.Id")

		if option.SelectRelatedChannel {
			groupBy = append(groupBy, "Channels.Id")
		}
	}

	if len(groupBy) > 0 {
		query = query.GroupBy(groupBy...)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterbyOption_ToSql")
	}

	rows, err := runner.QueryX(queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find product variant channel listings")
	}

	var (
		res                       []*product_and_discount.ProductVariantChannelListing
		chanNel                   channel.Channel
		variantChannelListing     product_and_discount.ProductVariantChannelListing
		availablePreorderQuantity int
		preorderQuantityAllocated int
		scanFields                = ps.ScanFields(variantChannelListing) // order of fields must be identical to select fields above
	)
	if option.SelectRelatedChannel {
		scanFields = append(scanFields, ps.Channel().ScanFields(chanNel)...)
	}
	if option.AnnotateAvailablePreorderQuantity {
		scanFields = append(scanFields, &availablePreorderQuantity)
	}
	if option.AnnotatePreorderQuantityAllocated {
		scanFields = append(scanFields, &preorderQuantityAllocated)
	}

	for rows.Next() {
		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row of product variant channel listing")
		}

		if option.SelectRelatedChannel {
			variantChannelListing.Channel = chanNel.DeepCopy()
		}
		if option.AnnotateAvailablePreorderQuantity {
			var copied_availablePreorderQuantity int = availablePreorderQuantity
			variantChannelListing.Set_availablePreorderQuantity(copied_availablePreorderQuantity)
		}
		if option.AnnotatePreorderQuantityAllocated {
			var copied_preorderQuantityAllocated int = preorderQuantityAllocated
			variantChannelListing.Set_preorderQuantityAllocated(copied_preorderQuantityAllocated)
		}
		res = append(res, variantChannelListing.DeepCopy())
	}

	if err := rows.Close(); err != nil {
		return nil, errors.Wrap(err, "failed to close rows")
	}

	return res, nil
}

// BulkUpsert performs bulk upsert given product variant channel listings then returns them
func (ps *SqlProductVariantChannelListingStore) BulkUpsert(transaction store_iface.SqlxTxExecutor, variantChannelListings []*product_and_discount.ProductVariantChannelListing) ([]*product_and_discount.ProductVariantChannelListing, error) {
	var (
		executor    store_iface.SqlxExecutor = ps.GetMasterX()
		saveQuery                            = "INSERT INTO " + store.ProductVariantChannelListingTableName + "(" + ps.ModelFields("").Join(",") + ") VALUES (" + ps.ModelFields(":").Join(",") + ")"
		updateQuery                          = "UPDATE " + store.ProductVariantChannelListingTableName + " SET " + ps.
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
				return nil, store.NewErrInvalidInput(store.ProductVariantChannelListingTableName, "VariantID/ChannelID", "duplicate")
			}
			return nil, errors.Wrapf(err, "failed to upsert a product variant channel listing with id=%s", listing.Id)
		}
		if numUpdated > 1 {
			return nil, errors.Errorf("%d product variant channel listing(s) with id=%s was/were updated instead of 1", numUpdated, listing.Id)
		}
	}

	return variantChannelListings, nil
}
