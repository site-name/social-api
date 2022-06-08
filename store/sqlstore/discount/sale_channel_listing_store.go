package discount

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlSaleChannelListingStore struct {
	store.Store
}

func NewSqlDiscountSaleChannelListingStore(sqlStore store.Store) store.DiscountSaleChannelListingStore {
	scls := &SqlSaleChannelListingStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.SaleChannelListing{}, store.SaleChannelListingTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH).SetNotNull(true)
		table.ColMap("SaleID").SetMaxSize(store.UUID_MAX_LENGTH).SetNotNull(false)
		table.ColMap("ChannelID").SetMaxSize(store.UUID_MAX_LENGTH).SetNotNull(true)

		table.SetUniqueTogether("SaleID", "ChannelID")
	}

	return scls
}

func (scls *SqlSaleChannelListingStore) ModelFields() []string {
	return []string{
		"SaleChannelListings.Id",
		"SaleChannelListings.SaleID",
		"SaleChannelListings.ChannelID",
		"SaleChannelListings.DiscountValue",
		"SaleChannelListings.Currency",
		"SaleChannelListings.CreateAt",
	}
}

func (scls *SqlSaleChannelListingStore) ScanFields(listing product_and_discount.SaleChannelListing) []interface{} {
	return []interface{}{
		&listing.Id,
		&listing.SaleID,
		&listing.ChannelID,
		&listing.DiscountValue,
		&listing.Currency,
		&listing.CreateAt,
	}
}

func (scls *SqlSaleChannelListingStore) CreateIndexesIfNotExists() {
	scls.CreateForeignKeyIfNotExists(store.SaleChannelListingTableName, "SaleID", store.SaleTableName, "Id", true)
	scls.CreateForeignKeyIfNotExists(store.SaleChannelListingTableName, "ChannelID", store.ChannelTableName, "Id", true)
}

// Save insert given instance into database then returns it
func (scls *SqlSaleChannelListingStore) Save(saleChannelListing *product_and_discount.SaleChannelListing) (*product_and_discount.SaleChannelListing, error) {
	saleChannelListing.PreSave()
	if err := saleChannelListing.IsValid(); err != nil {
		return nil, err
	}

	err := scls.GetMaster().Insert(saleChannelListing)
	if err != nil {
		if scls.IsUniqueConstraintError(err, []string{"SaleID", "ChannelID", "salechannellistings_saleid_channelid_key"}) {
			return nil, store.NewErrInvalidInput(store.SaleChannelListingTableName, "SaleID/ChannelID", "duplicate")
		}
		return nil, errors.Wrapf(err, "failed to save sale channel listing with id=%s", saleChannelListing.Id)
	}

	return saleChannelListing, nil
}

// Get finds and returns sale channel listing with given id
func (scls *SqlSaleChannelListingStore) Get(saleChannelListingID string) (*product_and_discount.SaleChannelListing, error) {
	var res product_and_discount.SaleChannelListing
	err := scls.GetReplica().SelectOne(
		&res,
		"SELECT * FROM "+store.SaleChannelListingTableName+" WHERE Id = :ID ORDER BY :OrderBy",
		map[string]interface{}{
			"ID":      saleChannelListingID,
			"OrderBy": store.TableOrderingMap[store.SaleChannelListingTableName],
		},
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.SaleChannelListingTableName, saleChannelListingID)
		}
		return nil, errors.Wrapf(err, "failed to find sale channel listing with id=%s", saleChannelListingID)
	}

	return &res, nil
}

// SaleChannelListingsWithOption finds a list of sale channel listings plus foreign channel slugs
func (scls *SqlSaleChannelListingStore) SaleChannelListingsWithOption(option *product_and_discount.SaleChannelListingFilterOption) (
	[]*struct {
		product_and_discount.SaleChannelListing
		ChannelSlug string
	},
	error,
) {

	query := scls.GetQueryBuilder().
		Select(scls.ModelFields()...).
		Column("Channels.Slug AS ChannelSlug").
		From(store.SaleChannelListingTableName).
		InnerJoin(store.ChannelTableName + " ON (Channels.Id = SaleChannelListings.ChannelID)").
		OrderBy(store.TableOrderingMap[store.SaleChannelListingTableName])

	// parse filter option
	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.SaleID != nil {
		query = query.Where(option.SaleID)
	}
	if option.ChannelID != nil {
		query = query.Where(option.ChannelID)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "SaleChannelListingsWithOption_ToSql")
	}

	var res []*struct {
		product_and_discount.SaleChannelListing
		ChannelSlug string
	}

	_, err = scls.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find sale channel listing with given option")
	}

	return res, nil
}
