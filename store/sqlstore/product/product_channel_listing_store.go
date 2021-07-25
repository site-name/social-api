package product

import (
	"database/sql"
	"strings"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlProductChannelListingStore struct {
	store.Store
}

var (
	// ProductChannelListingDuplicateKeys is used to catch crud duplicate errors
	ProductChannelListingDuplicateKeys = []string{"ProductID", "ChannelID", strings.ToLower(store.ProductChannelListingTableName) + "_productid_channelid_key"}
)

func NewSqlProductChannelListingStore(s store.Store) store.ProductChannelListingStore {
	pcls := &SqlProductChannelListingStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.ProductChannelListing{}, store.ProductChannelListingTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ChannelID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Currency").SetMaxSize(model.CURRENCY_CODE_MAX_LENGTH)

		table.SetUniqueTogether("ProductID", "ChannelID")
	}
	return pcls
}

func (ps *SqlProductChannelListingStore) CreateIndexesIfNotExists() {
	ps.CreateIndexIfNotExists("idx_productchannellistings_puplication_date", store.ProductChannelListingTableName, "PublicationDate")
	ps.CreateForeignKeyIfNotExists(store.ProductChannelListingTableName, "ProductID", store.ProductTableName, "Id", true)
	ps.CreateForeignKeyIfNotExists(store.ProductChannelListingTableName, "ChannelID", store.ChannelTableName, "Id", true)
}

func (ps *SqlProductChannelListingStore) Save(listing *product_and_discount.ProductChannelListing) (*product_and_discount.ProductChannelListing, error) {
	listing.PreSave()
	if err := listing.IsValid(); err != nil {
		return nil, err
	}

	if err := ps.GetMaster().Insert(listing); err != nil {
		if ps.IsUniqueConstraintError(err, ProductChannelListingDuplicateKeys) {
			return nil, store.NewErrInvalidInput(store.ProductChannelListingTableName, "ProductID/ChannelID", listing.ProductID+"/"+listing.ChannelID)
		}
		return nil, errors.Wrapf(err, "failed to save product channel listing with id=%s", listing.Id)
	}

	return listing, nil
}

func (ps *SqlProductChannelListingStore) Get(listingID string) (*product_and_discount.ProductChannelListing, error) {
	result, err := ps.GetReplica().Get(product_and_discount.ProductChannelListing{}, listingID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ProductChannelListingTableName, listingID)
		}
		return nil, errors.Wrapf(err, "failed to find product channel listing with id=%s", listingID)
	}

	return result.(*product_and_discount.ProductChannelListing), nil
}
