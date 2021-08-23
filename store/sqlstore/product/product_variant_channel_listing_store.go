package product

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlProductVariantChannelListingStore struct {
	store.Store
}

func NewSqlProductVariantChannelListingStore(s store.Store) store.ProductVariantChannelListingStore {
	pvcls := &SqlProductVariantChannelListingStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.ProductVariantChannelListing{}, store.ProductVariantChannelListingTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("VariantID").SetMaxSize(store.UUID_MAX_LENGTH).SetNotNull(true)
		table.ColMap("ChannelID").SetMaxSize(store.UUID_MAX_LENGTH).SetNotNull(true)
		table.ColMap("Currency").SetMaxSize(model.CURRENCY_CODE_MAX_LENGTH)

		table.SetUniqueTogether("VariantID", "ChannelID")
	}
	return pvcls
}

func (ps *SqlProductVariantChannelListingStore) CreateIndexesIfNotExists() {
	ps.CreateForeignKeyIfNotExists(store.ProductVariantChannelListingTableName, "VariantID", store.ProductVariantTableName, "Id", true)
	ps.CreateForeignKeyIfNotExists(store.ProductVariantChannelListingTableName, "ChannelID", store.ChannelTableName, "Id", true)
}

func (ps *SqlProductVariantChannelListingStore) ModelFields() []string {
	return []string{
		"ProductVariantChannelListings.Id",
		"ProductVariantChannelListings.VariantID",
		"ProductVariantChannelListings.ChannelID",
		"ProductVariantChannelListings.Currency",
		"ProductVariantChannelListings.PriceAmount",
		"ProductVariantChannelListings.CostPriceAmount",
		"ProductVariantChannelListings.CreateAt",
	}
}

// Save insert given value into database then returns it with an error
func (ps *SqlProductVariantChannelListingStore) Save(variantChannelListing *product_and_discount.ProductVariantChannelListing) (*product_and_discount.ProductVariantChannelListing, error) {
	variantChannelListing.PreSave()
	if err := variantChannelListing.IsValid(); err != nil {
		return nil, err
	}

	err := ps.GetMaster().Insert(variantChannelListing)
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
	err := ps.GetReplica().SelectOne(&res, "SELECT * FROM "+store.ProductVariantChannelListingTableName+" WHERE Id = :ID", map[string]interface{}{"ID": variantChannelListingID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ProductVariantChannelListingTableName, variantChannelListingID)
		}
		return nil, errors.Wrapf(err, "failed to find product variant channel listing with id=%s", variantChannelListingID)
	}

	return &res, nil
}

// FilterbyOption finds and returns all product variant channel listings filterd using given option
func (ps *SqlProductVariantChannelListingStore) FilterbyOption(option *product_and_discount.ProductVariantChannelListingFilterOption) ([]*product_and_discount.ProductVariantChannelListing, error) {
	query := ps.GetQueryBuilder().
		Select(ps.ModelFields()...).
		From(store.ProductVariantChannelListingTableName).
		OrderBy(store.TableOrderingMap[store.ProductVariantChannelListingTableName])

	// parse option
	if option.Id != nil {
		query = query.Where(option.Id.ToSquirrel("ProductVariantChannelListings.Id"))
	}
	if option.VariantID != nil {
		query = query.Where(option.VariantID.ToSquirrel("ProductVariantChannelListings.VariantID"))
	}
	if option.ChannelID != nil {
		query = query.Where(option.ChannelID.ToSquirrel("ProductVariantChannelListings.ChannelID"))
	}

	if option.PriceAmount != nil {
		query = query.Where(option.PriceAmount.ToSquirrel("ProductVariantChannelListings.PriceAmount"))
	}
	if option.VariantProductID != nil {
		query = query.
			InnerJoin(store.ProductVariantTableName + " ON (ProductVariants.Id = ProductVariantChannelListings.variantID)").
			Where(option.VariantProductID.ToSquirrel("ProductVariants.ProductID"))
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterbyOption_ToSql")
	}

	var res []*product_and_discount.ProductVariantChannelListing
	_, err = ps.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find product variant channel listings by given option")
	}

	return res, nil
}
