package discount

import (
	"database/sql"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlSaleChannelListingStore struct {
	store.Store
}

func NewSqlDiscountSaleChannelListingStore(sqlStore store.Store) store.DiscountSaleChannelListingStore {
	return &SqlSaleChannelListingStore{sqlStore}
}

// Save insert given instance into database then returns it
func (scls *SqlSaleChannelListingStore) Upsert(transaction boil.ContextTransactor, listings model.SaleChannelListingSlice) (model.SaleChannelListingSlice, error) {
	if transaction == nil {
		transaction = scls.GetMaster()
	}

	for _, listing := range listings {
		if listing == nil {
			continue
		}

		isSaving := false
		if listing.ID == "" {
			isSaving = true
			model_helper.SaleChannelListingPreSave(listing)
		} else {
			model_helper.SaleChannelListingPreUpdate(listing)
		}

		if err := model_helper.SaleChannelListingIsValid(*listing); err != nil {
			return nil, err
		}

		var err error
		if isSaving {
			err = listing.Insert(transaction, boil.Infer())
		} else {
			_, err = listing.Update(transaction, boil.Blacklist(model.SaleChannelListingColumns.CreatedAt))
		}

		if err != nil {
			return nil, err
		}
	}

	return listings, nil
}

// Get finds and returns sale channel listing with given id
func (scls *SqlSaleChannelListingStore) Get(id string) (*model.SaleChannelListing, error) {
	listing, err := model.FindSaleChannelListing(scls.GetReplica(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.SaleChannelListings, id)
		}
		return nil, err
	}

	return listing, nil
}

// SaleChannelListingsWithOption finds a list of sale channel listings plus foreign channel slugs
func (scls *SqlSaleChannelListingStore) FilterByOptions(option model_helper.SaleChannelListingFilterOption) (model.SaleChannelListingSlice, error) {
	return model.SaleChannelListings(option.Conditions...).All(scls.GetReplica())
}

func (s *SqlSaleChannelListingStore) Delete(transaction boil.ContextTransactor, ids []string) error {
	if transaction == nil {
		transaction = s.GetMaster()
	}
	_, err := model.SaleChannelListings(model.SaleChannelListingWhere.ID.IN(ids)).DeleteAll(transaction)
	return err
}
