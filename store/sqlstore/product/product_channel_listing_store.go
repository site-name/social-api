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
	pcls := &SqlProductChannelListingStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.ProductChannelListing{}, pcls.TableName("")).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ChannelID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Currency").SetMaxSize(model.CURRENCY_CODE_MAX_LENGTH)

		table.SetUniqueTogether("ProductID", "ChannelID")
	}
	return pcls
}

func (ps *SqlProductChannelListingStore) TableName(withField string) string {
	name := "ProductChannelListings"
	if withField != "" {
		name += "." + withField
	}

	return name
}

func (ps *SqlProductChannelListingStore) OrderBy() string {
	return "CreateAt ASC"
}

func (ps *SqlProductChannelListingStore) ModelFields() []string {
	return []string{
		"ProductChannelListings.Id",
		"ProductChannelListings.ProductID",
		"ProductChannelListings.ChannelID",
		"ProductChannelListings.VisibleInListings",
		"ProductChannelListings.AvailableForPurchase",
		"ProductChannelListings.Currency",
		"ProductChannelListings.DiscountedPriceAmount",
		"ProductChannelListings.CreateAt",
		"ProductChannelListings.PublicationDate",
		"ProductChannelListings.IsPublished",
	}
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

func (ps *SqlProductChannelListingStore) CreateIndexesIfNotExists() {
	ps.CreateIndexIfNotExists("idx_productchannellistings_puplication_date", ps.TableName(""), "PublicationDate")
	ps.CreateForeignKeyIfNotExists(ps.TableName(""), "ProductID", store.ProductTableName, "Id", true)
	ps.CreateForeignKeyIfNotExists(ps.TableName(""), "ChannelID", store.ChannelTableName, "Id", true)
}

// BulkUpsert performs bulk upsert on given product channel listings
func (ps *SqlProductChannelListingStore) BulkUpsert(listings []*product_and_discount.ProductChannelListing) ([]*product_and_discount.ProductChannelListing, error) {
	transaction, err := ps.GetMaster().Begin()
	if err != nil {
		return nil, errors.Wrap(err, "transaction_begin")
	}
	defer store.FinalizeTransaction(transaction)

	var (
		oldListing product_and_discount.ProductChannelListing
		numUpdated int64
		isSaving   bool
	)
	for _, listing := range listings {
		isSaving = false
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
			err = transaction.Insert(listing)
		} else {
			err = transaction.SelectOne(&oldListing, "SELECT * FROM "+ps.TableName("")+" WHERE Id = :ID", map[string]interface{}{"ID": listing.Id})
			if err != nil {
				if err == sql.ErrNoRows {
					return nil, store.NewErrNotFound(ps.TableName(""), listing.Id)
				}
				return nil, errors.Wrapf(err, "failed to find product channel listing with id=%s", listing.Id)
			}

			listing.CreateAt = oldListing.CreateAt
			numUpdated, err = transaction.Update(listing)
		}

		if err != nil {
			if ps.IsUniqueConstraintError(err, []string{"ProductID", "ChannelID", "productchannellistings_productid_channelid_key"}) {
				return nil, store.NewErrInvalidInput(ps.TableName(""), "ProductID/ChannelID", "duplicate")
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
	err := ps.GetReplica().SelectOne(&res, "SELECT * FROM "+ps.TableName("")+" WHERE Id = :ID", map[string]interface{}{"ID": listingID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(ps.TableName(""), listingID)
		}
		return nil, errors.Wrapf(err, "failed to find product channel listing with id=%s", listingID)
	}

	return &res, nil
}

// FilterByOption finds and returns all ProductChannelListings filtered by given option
func (ps *SqlProductChannelListingStore) FilterByOption(option *product_and_discount.ProductChannelListingFilterOption) ([]*product_and_discount.ProductChannelListing, error) {
	query := ps.
		GetQueryBuilder().
		Select(ps.ModelFields()...).
		From(ps.TableName("")).
		OrderBy(ps.OrderBy())

	// parse option
	if option.ProductID != nil {
		query = query.Where(option.ProductID.ToSquirrel("ProductChannelListings.ProductID"))
	}
	if option.ChannelID != nil {
		query = query.Where(option.ChannelID.ToSquirrel("ProductChannelListings.ChannelID"))
	}
	if option.ChannelSlug != nil {
		query = query.
			InnerJoin(store.ChannelTableName + " ON (Channels.Id = ProductChannelListings.ChannelID)").
			Where(squirrel.Eq{"Channels.ChannelSlug": *option.ChannelSlug})
	}
	if option.VisibleInListings != nil {
		query = query.Where(squirrel.Eq{"ProductChannelListings.VisibleInListings": *option.VisibleInListings})
	}
	if pur := option.AvailableForPurchase; pur != nil {
		query = query.Where(pur.ToSquirrel("ProductChannelListings.AvailableForPurchase"))
	}
	if option.Currency != nil {
		query = query.Where(option.Currency.ToSquirrel("ProductChannelListings.Currency"))
	}

	if option.ProductVariantsId != nil {
		query = query.
			InnerJoin(store.ProductTableName + " ON (Products.Id = ProductChannelListings.ProductID)").
			InnerJoin(store.ProductVariantTableName + " ON (ProductVariants.ProductID = Products.Id)").
			Where(option.ProductVariantsId.ToSquirrel("ProductVariants.Id"))
	}
	if option.PublicationDate != nil {
		query = query.Where(option.PublicationDate.ToSquirrel("ProductChannelListings.PublicationDate"))
	}
	if option.IsPublished != nil {
		query = query.Where(squirrel.Eq{"ProductChannelListings.IsPublished": *option.IsPublished})
	}

	sqlString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var listings []*product_and_discount.ProductChannelListing
	if _, err = ps.GetReplica().Select(&listings, sqlString, args...); err != nil {
		return nil, errors.Wrap(err, "failed to find product channel listings with given option")
	}

	// check if we need prefetch channels and founded `listings` is not empty
	if option.PrefetchChannel && len(listings) > 0 {

		channelIDs := product_and_discount.ProductChannelListings(listings).GetIDs(product_and_discount.ChannelIDs)
		if len(channelIDs) > 0 {
			channels, err := ps.Channel().FilterByOption(&channel.ChannelFilterOption{
				Id: squirrel.Eq{ps.Channel().TableName("Id"): channelIDs},
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
