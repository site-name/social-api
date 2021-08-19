package product

import (
	"database/sql"
	"strings"

	"github.com/Masterminds/squirrel"
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

func (ps *SqlProductChannelListingStore) ModelFields() []string {
	return []string{
		"ProductChannelListings.Id",
		"ProductChannelListings.ProductID",
		"ProductChannelListings.ChannelID",
		"ProductChannelListings.VisibleInListings",
		"ProductChannelListings.AvailableForPurchase",
		"ProductChannelListings.Currency",
		"ProductChannelListings.DiscountedPriceAmount",
		"ProductChannelListings.CreateAt",
		"ProductChannelListings.PublicationDate",
		"ProductChannelListings.IsPublished",
	}
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
	var res product_and_discount.ProductChannelListing
	err := ps.GetReplica().SelectOne(&res, "SELECT * FROM "+store.ProductChannelListingTableName+" WHERE Id = :ID", map[string]interface{}{"ID": listingID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ProductChannelListingTableName, listingID)
		}
		return nil, errors.Wrapf(err, "failed to find product channel listing with id=%s", listingID)
	}

	return &res, nil
}

// FilterByOption finds and returns all ProductChannelListings filtered by given option
func (ps *SqlProductChannelListingStore) FilterByOption(option *product_and_discount.ProductChannelListingFilterOption) ([]*product_and_discount.ProductChannelListing, error) {
	query := ps.
		GetQueryBuilder().
		Select(ps.ModelFields()...).
		From(store.ProductChannelListingTableName).
		OrderBy(store.TableOrderingMap[store.ProductChannelListingTableName])

	// parse option
	if option.ProductID != nil {
		query = query.Where(option.ProductID.ToSquirrel("ProductChannelListings.ProductID"))
	}
	if option.ChannelID != nil {
		query = query.Where(option.ChannelID.ToSquirrel("ProductChannelListings.ChannelID"))
	}
	if option.ChannelSlug != nil {
		query = query.
			InnerJoin(store.ChannelTableName + " ON (Channels.Id = ProductChannelListings.ChannelID)").
			Where(squirrel.Eq{"Channels.ChannelSlug": *option.ChannelSlug})
	}
	if option.VisibleInListings != nil {
		query = query.Where(squirrel.Eq{"ProductChannelListings.VisibleInListings": *option.VisibleInListings})
	}
	if pur := option.AvailableForPurchase; pur != nil {
		query = query.Where(pur.ToSquirrel("ProductChannelListings.AvailableForPurchase"))
	}
	if option.Currency != nil {
		query = query.Where(option.Currency.ToSquirrel("ProductChannelListings.Currency"))
	}

	if option.ProductVariantsId != nil {
		query = query.
			InnerJoin(store.ProductTableName + " ON (Products.Id = ProductChannelListings.ProductID)").
			InnerJoin(store.ProductVariantTableName + " ON (ProductVariants.ProductID = Products.Id)").
			Where(option.ProductVariantsId.ToSquirrel("ProductVariants.Id"))
	}
	if option.PublicationDate != nil {
		query = query.Where(option.PublicationDate.ToSquirrel("ProductChannelListings.PublicationDate"))
	}
	if option.IsPublished != nil {
		query = query.Where(squirrel.Eq{"ProductChannelListings.IsPublished": *option.IsPublished})
	}

	sqlString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var listings []*product_and_discount.ProductChannelListing
	if _, err = ps.GetReplica().Select(&listings, sqlString, args...); err != nil {
		return nil, errors.Wrap(err, "failed to find product channel listings with given option")
	}

	return listings, nil
}
