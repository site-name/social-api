package discount

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlSaleChannelListingStore struct {
	store.Store
}

func NewSqlDiscountSaleChannelListingStore(sqlStore store.Store) store.DiscountSaleChannelListingStore {
	return &SqlSaleChannelListingStore{sqlStore}
}

func (scls *SqlSaleChannelListingStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
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

func (scls *SqlSaleChannelListingStore) ScanFields(listing *model.SaleChannelListing) []interface{} {
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

	query := "INSERT INTO " + model.SaleChannelListingTableName + "(" + scls.ModelFields("").Join(",") + ") VALUES (" + scls.ModelFields(":").Join(",") + ")"
	_, err := scls.GetMasterX().NamedExec(query, saleChannelListing)
	if err != nil {
		if scls.IsUniqueConstraintError(err, []string{"SaleID", "ChannelID", "salechannellistings_saleid_channelid_key"}) {
			return nil, store.NewErrInvalidInput(model.SaleChannelListingTableName, "SaleID/ChannelID", "duplicate")
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
		"SELECT * FROM "+model.SaleChannelListingTableName+" WHERE Id = ?",
		saleChannelListingID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.SaleChannelListingTableName, saleChannelListingID)
		}
		return nil, errors.Wrapf(err, "failed to find sale channel listing with id=%s", saleChannelListingID)
	}

	return &res, nil
}

// SaleChannelListingsWithOption finds a list of sale channel listings plus foreign channel slugs
func (scls *SqlSaleChannelListingStore) SaleChannelListingsWithOption(option *model.SaleChannelListingFilterOption) ([]*model.SaleChannelListing, error) {

	selectFields := scls.ModelFields(model.SaleChannelListingTableName + ".")
	if option.SelectRelatedChannel {
		selectFields = append(selectFields, scls.Channel().ModelFields(model.ChannelTableName+".")...)
	}

	query := scls.GetQueryBuilder().
		Select(selectFields...).
		From(model.SaleChannelListingTableName)

	if option.SelectRelatedChannel {
		query = query.InnerJoin(model.ChannelTableName + " ON (Channels.Id = SaleChannelListings.ChannelID)")
	}

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

	rows, err := scls.GetReplicaX().QueryX(queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find sale channel listing with given option")
	}
	defer rows.Close()

	var res []*model.SaleChannelListing

	for rows.Next() {
		var listing model.SaleChannelListing
		var channel model.Channel
		var scanFields = scls.ScanFields(&listing)
		if option.SelectRelatedChannel {
			scanFields = append(scanFields, scls.Channel().ScanFields(&channel)...)
		}

		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan sale channel listing")
		}

		if option.SelectRelatedChannel {
			listing.SetChannel(&channel)
		}
		res = append(res, &listing)
	}

	return res, nil
}
