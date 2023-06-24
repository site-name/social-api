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
func (s *SqlShippingMethodChannelListingStore) Upsert(listing *model.ShippingMethodChannelListing) (*model.ShippingMethodChannelListing, error) {
	var isSaving bool
	if !model.IsValidId(listing.Id) {
		listing.Id = ""
		isSaving = true
		listing.PreSave()
	} else {
		listing.PreUpdate()
	}

	if err := listing.IsValid(); err != nil {
		return nil, err
	}

	var (
		err        error
		numUpdated int64
	)
	if isSaving {
		query := "INSERT INTO " + store.ShippingMethodChannelListingTableName + "(" + s.ModelFields("").Join(",") + ") VALUES (" + s.ModelFields(":").Join(",") + ")"
		_, err = s.GetMasterX().NamedExec(query, listing)

	} else {
		query := "UPDATE " + store.ShippingMethodChannelListingTableName + " SET " + s.
			ModelFields("").
			Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"
		var result sql.Result

		result, err = s.GetMasterX().NamedExec(query, listing)
		if err == nil && result != nil {
			numUpdated, _ = result.RowsAffected()
		}
	}

	if err != nil {
		if s.IsUniqueConstraintError(err, []string{"ShippingMethodID", "ChannelID", "shippingmethodchannellistings_shippingmethodid_channelid_key"}) {
			return nil, store.NewErrInvalidInput(store.ShippingMethodChannelListingTableName, "ShippingMethodID/ChannelID", listing.ShippingMethodID+"/"+listing.ChannelID)
		}
		return nil, errors.Wrapf(err, "failed to upsert shipping method channel listing with id=%s", listing.Id)
	}

	if numUpdated > 1 {
		return nil, errors.Errorf("multiple shipping method channel listings were updated: %d instead of 1", numUpdated)
	}

	listing.PopulateNonDbFields()
	return listing, nil
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
	for _, opt := range []squirrel.Sqlizer{option.ShippingMethodID, option.ChannelID} {
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

func (s *SqlShippingMethodChannelListingStore) BulkDelete(transaction store_iface.SqlxTxExecutor, ids []string) error {
	runner := s.GetReplicaX()
	if transaction != nil {
		runner = transaction
	}

	query, args, err := s.GetQueryBuilder().
		Delete(store.ShippingMethodChannelListingTableName).
		Where(squirrel.Eq{store.ShippingMethodChannelListingTableName + ".Id": ids}).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "BulkDelete_ToSql")
	}

	_, err = runner.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to delete shipping method channel listings")
	}

	return nil
}
