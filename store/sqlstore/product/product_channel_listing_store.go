package product

import (
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
	"gorm.io/gorm"
)

type SqlProductChannelListingStore struct {
	store.Store
}

func NewSqlProductChannelListingStore(s store.Store) store.ProductChannelListingStore {
	return &SqlProductChannelListingStore{s}
}

func (ps *SqlProductChannelListingStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
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

func (ps *SqlProductChannelListingStore) ScanFields(prd *model.ProductChannelListing) []interface{} {
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
func (ps *SqlProductChannelListingStore) BulkUpsert(transaction store_iface.SqlxExecutor, listings []*model.ProductChannelListing) ([]*model.ProductChannelListing, error) {
	runner := ps.GetMasterX()
	if transaction != nil {
		runner = transaction
	}

	var (
		saveQuery   = "INSERT INTO " + model.ProductChannelListingTableName + "(" + ps.ModelFields("").Join(",") + ") VALUES (" + ps.ModelFields(":").Join(",") + ")"
		updateQuery = "UPDATE " + model.ProductChannelListingTableName + " SET " + ps.
				ModelFields("").
				Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"
	)

	for _, listing := range listings {
		var (
			isSaving bool = false
			err      error
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
			_, err = runner.NamedExec(saveQuery, listing)

		} else {
			_, err = runner.NamedExec(updateQuery, listing)
		}

		if err != nil {
			if ps.IsUniqueConstraintError(err, []string{"ProductID", "ChannelID", "productchannellistings_productid_channelid_key"}) {
				return nil, store.NewErrInvalidInput(model.ProductChannelListingTableName, "ProductID/ChannelID", "duplicate")
			}
			return nil, errors.Wrapf(err, "failed to upsert product channel listing with id=%s", listing.Id)
		}
	}

	return listings, nil
}

func (ps *SqlProductChannelListingStore) Get(listingID string) (*model.ProductChannelListing, error) {
	var res model.ProductChannelListing

	err := ps.GetReplicaX().Get(&res, "SELECT * FROM "+model.ProductChannelListingTableName+" WHERE Id = ?", listingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.ProductChannelListingTableName, listingID)
		}
		return nil, errors.Wrapf(err, "failed to find product channel listing with id=%s", listingID)
	}

	return &res, nil
}

// FilterByOption finds and returns all ProductChannelListings filtered by given option
func (ps *SqlProductChannelListingStore) FilterByOption(option *model.ProductChannelListingFilterOption) ([]*model.ProductChannelListing, error) {
	query := ps.
		GetQueryBuilder().
		Select(ps.ModelFields(model.ProductChannelListingTableName + ".")...).
		From(model.ProductChannelListingTableName)

	// parse option
	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.ProductID != nil {
		query = query.Where(option.ProductID)
	}
	if option.ChannelID != nil {
		query = query.Where(option.ChannelID)
	}
	if option.ChannelSlug != nil && *option.ChannelSlug != "" {
		query = query.
			InnerJoin(model.ChannelTableName + " ON Channels.Id = ProductChannelListings.ChannelID").
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
			InnerJoin(model.ProductTableName + " ON Products.Id = ProductChannelListings.ProductID").
			InnerJoin(model.ProductVariantTableName + " ON ProductVariants.ProductID = Products.Id").
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

	var listings model.ProductChannelListings
	if err = ps.GetReplicaX().Select(&listings, sqlString, args...); err != nil {
		return nil, errors.Wrap(err, "failed to find product channel listings with given option")
	}

	// check if we need prefetch channels and founded `listings` is not empty
	if option.PrefetchChannel && len(listings) > 0 {

		channelIDs := listings.ChannelIDs()
		if len(channelIDs) > 0 {
			channels, err := ps.Channel().FilterByOption(&model.ChannelFilterOption{
				Id: squirrel.Eq{model.ChannelTableName + ".Id": channelIDs},
			})
			if err != nil {
				return nil, errors.Wrap(err, "failed to find channels by IDs")
			}

		outerLoop:
			for _, listing := range listings {
				for _, chanNel := range channels {
					if listing.ChannelID == chanNel.Id {
						listing.SetChannel(chanNel)
						continue outerLoop
					}
				}
			}
		}
	}

	return listings, nil
}
