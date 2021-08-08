package shipping

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/store"
)

type SqlShippingMethodChannelListingStore struct {
	store.Store
}

func NewSqlShippingMethodChannelListingStore(s store.Store) store.ShippingMethodChannelListingStore {
	smls := &SqlShippingMethodChannelListingStore{s}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(shipping.ShippingMethodChannelListing{}, store.ShippingMethodChannelListingTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ShippingMethodID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ChannelID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Currency").SetMaxSize(model.CURRENCY_CODE_MAX_LENGTH)

		table.SetUniqueTogether("ShippingMethodID", "ChannelID")
	}
	return smls
}

func (s *SqlShippingMethodChannelListingStore) CreateIndexesIfNotExists() {
	s.CreateForeignKeyIfNotExists(store.ShippingMethodChannelListingTableName, "ShippingMethodID", store.ShippingMethodTableName, "Id", true)
	s.CreateForeignKeyIfNotExists(store.ShippingMethodChannelListingTableName, "ChannelID", store.ChannelTableName, "Id", true)
}

// Upsert depends on given listing's Id to decide whether to save or update the listing
func (s *SqlShippingMethodChannelListingStore) Upsert(listing *shipping.ShippingMethodChannelListing) (*shipping.ShippingMethodChannelListing, error) {
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
		err        error
		numUpdated int64
	)
	if isSaving {
		err = s.GetMaster().Insert(listing)
	} else {
		_, err = s.Get(listing.Id)
		if err != nil {
			return nil, err
		}

		numUpdated, err = s.GetMaster().Update(listing)
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
func (s *SqlShippingMethodChannelListingStore) Get(listingID string) (*shipping.ShippingMethodChannelListing, error) {
	var res shipping.ShippingMethodChannelListing
	err := s.GetReplica().SelectOne(&res, "SELECT * FROM "+store.ShippingMethodChannelListingTableName+" WHERE Id = :ID", map[string]interface{}{"ID": listingID})
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
func (s *SqlShippingMethodChannelListingStore) FilterByOption(option *shipping.ShippingMethodChannelListingFilterOption) ([]*shipping.ShippingMethodChannelListing, error) {
	query := s.GetQueryBuilder().
		Select("*").
		From(store.ShippingMethodChannelListingTableName).
		OrderBy(store.TableOrderingMap[store.ShippingMethodChannelListingTableName])

	if option.ShippingMethodID != nil {
		query = query.Where(option.ShippingMethodID.ToSquirrel("ShippingMethodID"))
	}
	if option.ChannelID != nil {
		query = query.Where(option.ChannelID.ToSquirrel("ChannelID"))
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_tosql")
	}

	var res []*shipping.ShippingMethodChannelListing
	_, err = s.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find shipping method channel listings by option")
	}

	return res, nil
}
