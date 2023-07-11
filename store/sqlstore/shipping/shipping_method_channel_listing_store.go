package shipping

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

type SqlShippingMethodChannelListingStore struct {
	store.Store
}

func NewSqlShippingMethodChannelListingStore(s store.Store) store.ShippingMethodChannelListingStore {
	return &SqlShippingMethodChannelListingStore{s}
}

func (s *SqlShippingMethodChannelListingStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"ShippingMethodID",
		"ChannelID",
		"MinimumOrderPriceAmount",
		"Currency",
		"MaximumOrderPriceAmount",
		"PriceAmount",
		"CreateAt",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// Upsert depends on given listing's Id to decide whether to save or update the listing
func (s *SqlShippingMethodChannelListingStore) Upsert(transaction store_iface.SqlxTxExecutor, listings model.ShippingMethodChannelListings) (model.ShippingMethodChannelListings, error) {
	var (
		saveQuery   = "INSERT INTO " + store.ShippingMethodChannelListingTableName + "(" + s.ModelFields("").Join(",") + ") VALUES (" + s.ModelFields(":").Join(",") + ")"
		updateQuery = "UPDATE " + store.ShippingMethodChannelListingTableName + " SET " + s.
				ModelFields("").
				Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"
		runner = s.GetMasterX()
	)
	if transaction != nil {
		runner = transaction
	}

	for _, listing := range listings {
		var isSaving bool

		if listing.Id == "" {
			isSaving = true
			listing.PreSave()
		} else {
			listing.PreUpdate()
		}

		if err := listing.IsValid(); err != nil {
			return nil, err
		}

		var (
			err    error
			result sql.Result
		)
		if isSaving {
			result, err = runner.NamedExec(saveQuery, listing)
		} else {
			result, err = runner.NamedExec(updateQuery, listing)
		}

		if err != nil {
			if s.IsUniqueConstraintError(err, []string{"ShippingMethodID", "ChannelID", "shippingmethodchannellistings_shippingmethodid_channelid_key"}) {
				return nil, store.NewErrInvalidInput(store.ShippingMethodChannelListingTableName, "ShippingMethodID/ChannelID", listing.ShippingMethodID+"/"+listing.ChannelID)
			}
			return nil, errors.Wrapf(err, "failed to upsert shipping method channel listing with id=%s", listing.Id)
		}

		numUpserted, _ := result.RowsAffected()
		if numUpserted > 1 {
			return nil, errors.Errorf("%d shipping method channel listing(s) upserted instead of 1", numUpserted)
		}
	}

	return listings, nil
}

// Get finds a shipping method channel listing with given listingID
func (s *SqlShippingMethodChannelListingStore) Get(listingID string) (*model.ShippingMethodChannelListing, error) {
	var res model.ShippingMethodChannelListing
	err := s.GetReplicaX().Get(&res, "SELECT * FROM "+store.ShippingMethodChannelListingTableName+" WHERE Id = ?", listingID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ShippingMethodChannelListingTableName, listingID)
		}
		return nil, errors.Wrapf(err, "failed to find shipping method channel listing with id=%s", listingID)
	}

	res.PopulateNonDbFields()
	return &res, nil
}

// FilterByOption returns a list of shipping method channel listings based on given option. result sorted by creation time ASC
func (s *SqlShippingMethodChannelListingStore) FilterByOption(option *model.ShippingMethodChannelListingFilterOption) ([]*model.ShippingMethodChannelListing, error) {
	query := s.GetQueryBuilder().
		Select(s.ModelFields(store.ShippingMethodChannelListingTableName + ".")...).
		From(store.ShippingMethodChannelListingTableName)

	// parse filter option
	for _, opt := range []squirrel.Sqlizer{option.ShippingMethodID, option.ChannelID, option.Id} {
		if opt != nil {
			query = query.Where(opt)
		}
	}
	if option.ChannelSlug != nil {
		query = query.
			InnerJoin(store.ChannelTableName + " ON Channels.Id = ShippingMethodChannelListings.ChannelID").
			Where(option.ChannelSlug)
	}
	if option.ShippingMethod_ShippingZoneID_Inner != nil {
		query = query.
			InnerJoin(store.ShippingMethodTableName + " ON ShippingMethods.Id = ShippingMethodChannelListings.ShippingMethodID").
			InnerJoin(store.ShippingZoneTableName + " ON ShippingZones.Id = ShippingMethods.ShippingZoneID").
			Where(option.ShippingMethod_ShippingZoneID_Inner)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_tosql")
	}

	var res []*model.ShippingMethodChannelListing
	err = s.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find shipping method channel listings by option")
	}

	return res, nil
}

func (s *SqlShippingMethodChannelListingStore) BulkDelete(transaction store_iface.SqlxTxExecutor, options *model.ShippingMethodChannelListingFilterOption) error {
	query := s.GetQueryBuilder().Delete(store.ShippingMethodChannelListingTableName)
	for _, opt := range []squirrel.Sqlizer{options.ShippingMethodID, options.ChannelID, options.Id} {
		if opt != nil {
			query = query.Where(opt)
		}
	}

	queryStr, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "BulkDelete_ToSql")
	}

	runner := s.GetReplicaX()
	if transaction != nil {
		runner = transaction
	}

	_, err = runner.Exec(queryStr, args...)
	if err != nil {
		return errors.Wrap(err, "failed to delete shipping method channel listings")
	}

	return nil
}
