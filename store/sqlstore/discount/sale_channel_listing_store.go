package discount

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlSaleChannelListingStore struct {
	store.Store
}

func NewSqlDiscountSaleChannelListingStore(sqlStore store.Store) store.DiscountSaleChannelListingStore {
	return &SqlSaleChannelListingStore{sqlStore}
}

func (scls *SqlSaleChannelListingStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
		"Id",
		"SaleID",
		"ChannelID",
		"DiscountValue",
		"Currency",
		"CreateAt",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (scls *SqlSaleChannelListingStore) ScanFields(listing model.SaleChannelListing) []interface{} {
	return []interface{}{
		&listing.Id,
		&listing.SaleID,
		&listing.ChannelID,
		&listing.DiscountValue,
		&listing.Currency,
		&listing.CreateAt,
	}
}

// Save insert given instance into database then returns it
func (scls *SqlSaleChannelListingStore) Save(saleChannelListing *model.SaleChannelListing) (*model.SaleChannelListing, error) {
	saleChannelListing.PreSave()
	if err := saleChannelListing.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.SaleChannelListingTableName + "(" + scls.ModelFields("").Join(",") + ") VALUES (" + scls.ModelFields(":").Join(",") + ")"
	_, err := scls.GetMasterX().NamedExec(query, saleChannelListing)
	if err != nil {
		if scls.IsUniqueConstraintError(err, []string{"SaleID", "ChannelID", "salechannellistings_saleid_channelid_key"}) {
			return nil, store.NewErrInvalidInput(store.SaleChannelListingTableName, "SaleID/ChannelID", "duplicate")
		}
		return nil, errors.Wrapf(err, "failed to save sale channel listing with id=%s", saleChannelListing.Id)
	}

	return saleChannelListing, nil
}

// Get finds and returns sale channel listing with given id
func (scls *SqlSaleChannelListingStore) Get(saleChannelListingID string) (*model.SaleChannelListing, error) {
	var res model.SaleChannelListing

	err := scls.GetReplicaX().Get(
		&res,
		"SELECT * FROM "+store.SaleChannelListingTableName+" WHERE Id = ? ORDER BY ?",
		saleChannelListingID,
		store.TableOrderingMap[store.SaleChannelListingTableName],
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
func (scls *SqlSaleChannelListingStore) SaleChannelListingsWithOption(option *model.SaleChannelListingFilterOption) (
	[]*struct {
		model.SaleChannelListing
		ChannelSlug string
	},
	error,
) {

	query := scls.GetQueryBuilder().
		Select(scls.ModelFields(store.SaleChannelListingTableName + ".")...).
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
		model.SaleChannelListing
		ChannelSlug string
	}

	err = scls.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find sale channel listing with given option")
	}

	return res, nil
}
