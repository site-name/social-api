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

func (ps *SqlProductChannelListingStore) FilterByOption(option *product_and_discount.ProductChannelListingFilterOption) ([]*product_and_discount.ProductChannelListing, error) {
	if option == nil {
		return nil, nil
	}

	query := ps.
		GetQueryBuilder().
		Select("*").
		From(store.ProductChannelListingTableName + " AS PCL").
		OrderBy("PCL.CreateAt ASC") // since Ids are UUIDs, so order create at is a good option

	// check product id
	if option.ProductID != nil {
		query = query.Where(option.ProductID.ToSquirrel("PCL.ProductID"))
	}

	// check channel id
	if option.ChannelID != nil {
		query = query.Where(option.ChannelID.ToSquirrel("PCL.ChannelID"))
	}

	// check channel slug
	if option.ChannelSlug != nil {
		query = query.
			Where(squirrel.Eq{"Cn.ChannelSlug": *option.ChannelSlug}).
			InnerJoin(store.ChannelTableName + " AS Cn ON (Cn.Id = PCL.ChannelID)")
	}

	// check visible in listing
	if option.VisibleInListings != nil {
		query = query.Where(squirrel.Eq{"PCL.VisibleInListings": *option.VisibleInListings})
	}

	// check available for purchase
	if pur := option.AvailableForPurchase; pur != nil {
		query = query.Where(pur.ToSquirrel("PCL.AvailableForPurchase"))
	}

	// check currency
	if option.Currency != nil {
		query = query.Where(option.Currency.ToSquirrel("PCL.Currency"))
	}

	// check product variant
	if option.ProductVariantsId != nil {
		query = query.
			InnerJoin(store.ProductTableName + " AS P ON (P.Id = PCL.ProductID)").
			InnerJoin(store.ProductVariantTableName + " AS PV ON (PV.ProductID = P.Id)").
			Where(option.ProductVariantsId.ToSquirrel("PV.Id"))
	}

	// check publish
	if option.PublicationDate != nil {
		query = query.Where(option.PublicationDate.ToSquirrel("PCL.PublicationDate"))
	}

	if option.IsPublished != nil {
		query = query.Where(squirrel.Eq{"PCL.IsPublished": *option.IsPublished})
	}

	sqlString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "sql to string")
	}

	var listings []*product_and_discount.ProductChannelListing
	if _, err = ps.GetReplica().Select(&listings, sqlString, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ProductChannelListingTableName, "")
		}
		return nil, errors.Wrap(err, "failed to find product channel listings with given option")
	}

	return listings, nil
}
