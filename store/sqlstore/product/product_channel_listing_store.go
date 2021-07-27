package product

import (
	"database/sql"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlProductChannelListingStore struct {
	store.Store
}

var (
	// ProductChannelListingDuplicateKeys is used to catch crud duplicate errors
	ProductChannelListingDuplicateKeys = []string{"ProductID", "ChannelID", strings.ToLower(store.ProductChannelListingTableName) + "_productid_channelid_key"}
)

func NewSqlProductChannelListingStore(s store.Store) store.ProductChannelListingStore {
	pcls := &SqlProductChannelListingStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.ProductChannelListing{}, store.ProductChannelListingTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ChannelID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Currency").SetMaxSize(model.CURRENCY_CODE_MAX_LENGTH)

		table.SetUniqueTogether("ProductID", "ChannelID")
	}
	return pcls
}

func (ps *SqlProductChannelListingStore) CreateIndexesIfNotExists() {
	ps.CreateIndexIfNotExists("idx_productchannellistings_puplication_date", store.ProductChannelListingTableName, "PublicationDate")
	ps.CreateForeignKeyIfNotExists(store.ProductChannelListingTableName, "ProductID", store.ProductTableName, "Id", true)
	ps.CreateForeignKeyIfNotExists(store.ProductChannelListingTableName, "ChannelID", store.ChannelTableName, "Id", true)
}

func (ps *SqlProductChannelListingStore) Save(listing *product_and_discount.ProductChannelListing) (*product_and_discount.ProductChannelListing, error) {
	listing.PreSave()
	if err := listing.IsValid(); err != nil {
		return nil, err
	}

	if err := ps.GetMaster().Insert(listing); err != nil {
		if ps.IsUniqueConstraintError(err, ProductChannelListingDuplicateKeys) {
			return nil, store.NewErrInvalidInput(store.ProductChannelListingTableName, "ProductID/ChannelID", listing.ProductID+"/"+listing.ChannelID)
		}
		return nil, errors.Wrapf(err, "failed to save product channel listing with id=%s", listing.Id)
	}

	return listing, nil
}

func (ps *SqlProductChannelListingStore) Get(listingID string) (*product_and_discount.ProductChannelListing, error) {
	result, err := ps.GetReplica().Get(product_and_discount.ProductChannelListing{}, listingID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ProductChannelListingTableName, listingID)
		}
		return nil, errors.Wrapf(err, "failed to find product channel listing with id=%s", listingID)
	}

	return result.(*product_and_discount.ProductChannelListing), nil
}

func (ps *SqlProductChannelListingStore) FilterByOption(option *product_and_discount.ProductChannelListingFilterOption) ([]*product_and_discount.ProductChannelListing, error) {
	if option == nil {
		return nil, nil
	}

	andCondition := squirrel.And{}

	query := ps.
		GetQueryBuilder().
		Select("*").
		From(store.ProductChannelListingTableName + " AS PCL")

	// check product id
	if option.ProductID != nil {
		var eq interface{}
		if model.IsValidId(option.ProductID.Eq) {
			eq = option.ProductID.Eq
		} else if len(option.ProductID.In) > 0 {
			eq = option.ProductID.In
		}
		andCondition = append(andCondition, squirrel.Eq{"PCL.ProductID": eq})
	}

	// check channel id
	if option.ChannelID != nil {
		var eq interface{}
		if model.IsValidId(option.ChannelID.Eq) {
			eq = option.ChannelID.Eq
		} else if len(option.ChannelID.In) > 0 {
			eq = option.ChannelID.In
		}
		andCondition = append(andCondition, squirrel.Eq{"PCL.ChannelID": eq})
	}

	// check channel slug
	if option.ChannelSlug != nil {
		andCondition = append(andCondition, squirrel.Eq{"Cn.ChannelSlug": *option.ChannelSlug})

		query = query.InnerJoin(store.ChannelTableName + " AS Cn ON (Cn.Id = PCL.ChannelID)")
	}

	// check visible in listing
	if option.VisibleInListings != nil {
		andCondition = append(andCondition, squirrel.Eq{"PCL.VisibleInListings": *option.VisibleInListings})
	}

	// check available for purchase
	if pur := option.AvailableForPurchase; pur != nil {
		andCondition = append(andCondition, pur.ToSquirrelCondition("PCL.AvailableForPurchase")...)
	}

	// check currency
	if option.Currency != nil {
		var eq interface{}
		// GetCurrencyPrecision() can check if a currency is valid too
		if _, err := goprices.GetCurrencyPrecision(option.Currency.Eq); err == nil {
			eq = option.Currency.Eq
		} else if len(option.Currency.In) > 0 {
			for i, cur := range option.Currency.In {
				if _, err := goprices.GetCurrencyPrecision(cur); err != nil {
					option.Currency.In = append(option.Currency.In[:i], option.ProductID.In[i+1:]...)
				}
			}
			eq = option.Currency.In
		}
		andCondition = append(andCondition, squirrel.Eq{"PCL.Currency": eq})
	}

	// check product variant
	if option.ProductVariantsId != nil {
		var eq interface{}
		if model.IsValidId(option.ProductVariantsId.Eq) {
			eq = option.ProductVariantsId.Eq
		} else if len(option.ProductVariantsId.In) > 0 {
			eq = option.ProductVariantsId.In
		}
		andCondition = append(andCondition, squirrel.Eq{"PV.Id": eq})

		query = query.
			InnerJoin(store.ProductTableName + " AS P ON (P.Id = PCL.ProductID)").
			InnerJoin(store.ProductVariantTableName + " AS PV ON (PV.ProductID = P.Id)")
	}

	// check publish
	if option.PublicationDate != nil {
		andCondition = append(andCondition, option.PublicationDate.ToSquirrelCondition("PCL.PublicationDate")...)
	}

	if option.IsPublished != nil {
		andCondition = append(andCondition, squirrel.Eq{"PCL.IsPublished": *option.IsPublished})
	}

	sqlString, args, err := query.
		Where(andCondition).
		OrderBy("PCL.CreateAt ASC").
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "sql to string")
	}

	var listings []*product_and_discount.ProductChannelListing
	if _, err = ps.GetReplica().Select(&listings, sqlString, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ProductChannelListingTableName, "")
		}
		return nil, errors.Wrap(err, "failed to find product channel listings with given option")
	}

	return listings, nil
}
