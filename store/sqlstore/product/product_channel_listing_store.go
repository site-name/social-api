package product

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlProductChannelListingStore struct {
	store.Store
}

func NewSqlProductChannelListingStore(s store.Store) store.ProductChannelListingStore {
	return &SqlProductChannelListingStore{s}
}

func (ps *SqlProductChannelListingStore) ModelFields(prefix string) model.StringArray {
	res := model.StringArray{
		"Id",
		"ProductID",
		"ChannelID",
		"VisibleInListings",
		"AvailableForPurchase",
		"Currency",
		"DiscountedPriceAmount",
		"CreateAt",
		"PublicationDate",
		"IsPublished",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (ps *SqlProductChannelListingStore) ScanFields(prd product_and_discount.ProductChannelListing) []interface{} {
	return []interface{}{
		&prd.Id,
		&prd.ProductID,
		&prd.ChannelID,
		&prd.VisibleInListings,
		&prd.AvailableForPurchase,
		&prd.Currency,
		&prd.DiscountedPriceAmount,
		&prd.CreateAt,
		&prd.PublicationDate,
		&prd.IsPublished,
	}
}

// BulkUpsert performs bulk upsert on given product channel listings
func (ps *SqlProductChannelListingStore) BulkUpsert(listings []*product_and_discount.ProductChannelListing) ([]*product_and_discount.ProductChannelListing, error) {
	transaction, err := ps.GetMasterX().Beginx()
	if err != nil {
		return nil, errors.Wrap(err, "transaction_begin")
	}
	defer store.FinalizeTransaction(transaction)

	var (
		saveQuery   = "INSERT INTO " + store.ProductChannelListingTableName + "(" + ps.ModelFields("").Join(",") + ") VALUES (" + ps.ModelFields(":").Join(",") + ")"
		updateQuery = "UPDATE " + store.ProductChannelListingTableName + " SET " + ps.
				ModelFields("").
				Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"
	)

	for _, listing := range listings {
		var (
			isSaving   bool
			numUpdated int64
		)

		if listing.Id == "" {
			listing.PreSave()
			isSaving = true
		} else {
			listing.PreUpdate()
		}

		if err := listing.IsValid(); err != nil {
			return nil, err
		}

		if isSaving {
			_, err = transaction.NamedExec(saveQuery, listing)

		} else {
			var result sql.Result
			result, err = transaction.NamedExec(updateQuery, listing)
			if err == nil && result != nil {
				numUpdated, _ = result.RowsAffected()
			}
		}

		if err != nil {
			if ps.IsUniqueConstraintError(err, []string{"ProductID", "ChannelID", "productchannellistings_productid_channelid_key"}) {
				return nil, store.NewErrInvalidInput(store.ProductChannelListingTableName, "ProductID/ChannelID", "duplicate")
			}
			return nil, errors.Wrapf(err, "failed to upsert product channel listing with id=%s", listing.Id)
		}
		if numUpdated > 1 {
			return nil, errors.New("multiple listings were updated: %d instead of 1")
		}
	}

	if err := transaction.Commit(); err != nil {
		return nil, errors.Wrap(err, "transaction_commit")
	}

	return listings, nil
}

func (ps *SqlProductChannelListingStore) Get(listingID string) (*product_and_discount.ProductChannelListing, error) {
	var res product_and_discount.ProductChannelListing

	err := ps.GetReplicaX().Get(&res, "SELECT * FROM "+store.ProductChannelListingTableName+" WHERE Id = ?", listingID)
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
		Select(ps.ModelFields(store.ProductChannelListingTableName + ".")...).
		From(store.ProductChannelListingTableName).
		OrderBy(store.TableOrderingMap[store.ProductChannelListingTableName])

	// parse option
	if option.ProductID != nil {
		query = query.Where(option.ProductID)
	}
	if option.ChannelID != nil {
		query = query.Where(option.ChannelID)
	}
	if option.ChannelSlug != nil && *option.ChannelSlug != "" {
		query = query.
			InnerJoin(store.ChannelTableName + " ON Channels.Id = ProductChannelListings.ChannelID").
			Where(squirrel.Eq{"Channels.ChannelSlug": *option.ChannelSlug})
	}
	if option.VisibleInListings != nil {
		query = query.Where(squirrel.Eq{"ProductChannelListings.VisibleInListings": *option.VisibleInListings})
	}
	if option.AvailableForPurchase != nil {
		query = query.Where(option.AvailableForPurchase)
	}
	if option.Currency != nil {
		query = query.Where(option.Currency)
	}

	if option.ProductVariantsId != nil {
		query = query.
			InnerJoin(store.ProductTableName + " ON (Products.Id = ProductChannelListings.ProductID)").
			InnerJoin(store.ProductVariantTableName + " ON (ProductVariants.ProductID = Products.Id)").
			Where(option.ProductVariantsId)
	}
	if option.PublicationDate != nil {
		query = query.Where(option.PublicationDate)
	}
	if option.IsPublished != nil {
		query = query.Where(squirrel.Eq{"ProductChannelListings.IsPublished": *option.IsPublished})
	}

	sqlString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var listings product_and_discount.ProductChannelListings
	if err = ps.GetReplicaX().Select(&listings, sqlString, args...); err != nil {
		return nil, errors.Wrap(err, "failed to find product channel listings with given option")
	}

	// check if we need prefetch channels and founded `listings` is not empty
	if option.PrefetchChannel && len(listings) > 0 {

		channelIDs := listings.ChannelIDs()
		if len(channelIDs) > 0 {
			channels, err := ps.Channel().FilterByOption(&channel.ChannelFilterOption{
				Id: squirrel.Eq{store.ChannelTableName + ".Id": channelIDs},
			})
			if err != nil {
				return nil, errors.Wrap(err, "failed to find channels by IDs")
			}

		outerLoop:
			for _, listing := range listings {
				for _, chanNel := range channels {
					if listing.ChannelID == chanNel.Id {
						listing.Channel = chanNel
						continue outerLoop
					}
				}
			}
		}
	}

	return listings, nil
}
