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
