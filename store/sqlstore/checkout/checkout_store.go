package checkout

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/mattermost/gorp"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlCheckoutStore struct {
	store.Store
}

func NewSqlCheckoutStore(sqlStore store.Store) store.CheckoutStore {
	cs := &SqlCheckoutStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(checkout.Checkout{}, store.CheckoutTableName).SetKeys(false, "Token")
		table.ColMap("Token").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("UserID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ShopID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ChannelID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("BillingAddressID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ShippingAddressID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ShippingMethodID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("CollectionPointID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("DiscountName").SetMaxSize(checkout.CHECKOUT_DISCOUNT_NAME_MAX_LENGTH)
		table.ColMap("TranslatedDiscountName").SetMaxSize(checkout.CHECKOUT_TRANSLATED_DISCOUNT_NAME_MAX_LENGTH)
		table.ColMap("VoucherCode").SetMaxSize(checkout.CHECKOUT_VOUCHER_CODE_MAX_LENGTH)
		table.ColMap("TrackingCode").SetMaxSize(checkout.CHECKOUT_TRACKING_CODE_MAX_LENGTH)
		table.ColMap("Country").SetMaxSize(model.SINGLE_COUNTRY_CODE_MAX_LENGTH)
	}

	return cs
}

func (cs *SqlCheckoutStore) CreateIndexesIfNotExists() {
	cs.CreateIndexIfNotExists("idx_checkouts_userid", store.CheckoutTableName, "UserID")
	cs.CreateIndexIfNotExists("idx_checkouts_token", store.CheckoutTableName, "Token")
	cs.CreateIndexIfNotExists("idx_checkouts_channelid", store.CheckoutTableName, "ChannelID")
	cs.CreateIndexIfNotExists("idx_checkouts_billing_address_id", store.CheckoutTableName, "BillingAddressID")
	cs.CreateIndexIfNotExists("idx_checkouts_shipping_address_id", store.CheckoutTableName, "ShippingAddressID")
	cs.CreateIndexIfNotExists("idx_checkouts_shipping_method_id", store.CheckoutTableName, "ShippingMethodID")

	cs.CreateForeignKeyIfNotExists(store.CheckoutTableName, "ShopID", store.ShopTableName, "Id", true)
	cs.CreateForeignKeyIfNotExists(store.CheckoutTableName, "UserID", store.UserTableName, "Id", true)
	cs.CreateForeignKeyIfNotExists(store.CheckoutTableName, "ChannelID", store.ChannelTableName, "Id", false)
	cs.CreateForeignKeyIfNotExists(store.CheckoutTableName, "BillingAddressID", store.AddressTableName, "Id", false)
	cs.CreateForeignKeyIfNotExists(store.CheckoutTableName, "ShippingAddressID", store.AddressTableName, "Id", false)
	cs.CreateForeignKeyIfNotExists(store.CheckoutTableName, "ShippingMethodID", store.ShippingMethodTableName, "Id", false)
	cs.CreateForeignKeyIfNotExists(store.CheckoutTableName, "CollectionPointID", store.WarehouseTableName, "Id", false)
}

// Upsert depends on given checkout's Token property to decide to update or insert it
func (cs *SqlCheckoutStore) Upsert(ckout *checkout.Checkout) (*checkout.Checkout, error) {
	var isSaving bool

	if ckout.Token == "" {
		isSaving = true
		ckout.PreSave()
	} else {
		ckout.PreUpdate()
	}

	if err := ckout.IsValid(); err != nil {
		return nil, err
	}

	var (
		err         error
		numUpdated  int64
		oldCheckout *checkout.Checkout
	)
	if isSaving {
		err = cs.GetMaster().Insert(ckout)
	} else {
		// validate if checkout exist
		oldCheckout, err = cs.Get(ckout.Token)
		if err != nil {
			return nil, err
		}

		// set fields that CANNOT be changed
		ckout.BillingAddressID = oldCheckout.BillingAddressID
		ckout.ShippingAddressID = oldCheckout.ShippingAddressID
		ckout.Token = oldCheckout.Token

		// update checkout
		numUpdated, err = cs.GetMaster().Update(ckout)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "failed to upsert checkout with token=%s", ckout.Token)
	}
	if numUpdated > 1 {
		return nil, errors.Errorf("multiple checkouts were updated: %d instead of 1", numUpdated)
	}

	return ckout, nil
}

// Get finds a checkout with given token (checkouts use tokens(uuids) as primary keys)
func (cs *SqlCheckoutStore) Get(token string) (*checkout.Checkout, error) {
	var ckout *checkout.Checkout
	err := cs.GetReplica().SelectOne(
		&ckout,
		`SELECT * FROM `+store.CheckoutTableName+` WHERE Token = :Token`,
		map[string]interface{}{
			"Token": token,
		},
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.CheckoutTableName, token)
		}
		return nil, errors.Wrapf(err, "failed to find checkout with token=%s", token)
	}

	return ckout, nil
}

type checkoutStatement string

const (
	delete checkoutStatement = "delete"
	slect  checkoutStatement = "select"
)

// commonFilterQueryBuilder is common function, used to build checkout(s) filter queries.
func (cs *SqlCheckoutStore) commonFilterQueryBuilder(option *checkout.CheckoutFilterOption, statementType checkoutStatement) interface{} {
	andCondition := squirrel.And{}
	// parse option
	if option.Token != nil {
		andCondition = append(andCondition, option.Token.ToSquirrel("Token"))
	}
	if option.UserID != nil {
		andCondition = append(andCondition, option.UserID.ToSquirrel("UserID"))
	}
	if option.ChannelID != nil {
		andCondition = append(andCondition, option.ChannelID.ToSquirrel("ChannelID"))
	}

	if statementType == slect {
		return cs.GetQueryBuilder().Select("*").From(store.CheckoutTableName).Where(andCondition)
	}
	return cs.GetQueryBuilder().Delete(store.CheckoutTableName).Where(andCondition)
}

// GetByOption finds and returns 1 checkout based on given option
func (cs *SqlCheckoutStore) GetByOption(option *checkout.CheckoutFilterOption) (*checkout.Checkout, error) {
	queryString, args, err := cs.commonFilterQueryBuilder(option, slect).(squirrel.SelectBuilder).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetbyOption_ToSql")
	}

	var res *checkout.Checkout
	err = cs.GetReplica().SelectOne(&res, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.CheckoutTableName, "option")
		}
		return nil, errors.Wrap(err, "failed to find checkout woth given option")
	}

	return res, nil
}

// FilterByOption finds and returns a list of checkout based on given option
func (cs *SqlCheckoutStore) FilterByOption(option *checkout.CheckoutFilterOption) ([]*checkout.Checkout, error) {
	queryString, args, err := cs.commonFilterQueryBuilder(option, slect).(squirrel.SelectBuilder).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res []*checkout.Checkout
	_, err = cs.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find checkouts by given option")
	}

	return res, nil
}

// FetchCheckoutLinesAndPrefetchRelatedValue Fetch checkout lines as CheckoutLineInfo objects.
func (cs *SqlCheckoutStore) FetchCheckoutLinesAndPrefetchRelatedValue(ckout *checkout.Checkout) ([]*checkout.CheckoutLineInfo, error) {
	// please refer to file checkout_store_sql.md for details

	// finds all checkout lines belong to given checkout:
	tx, err := cs.GetReplica().Begin()
	if err != nil {
		return nil, errors.Wrap(err, "transaction_begin")
	}
	defer store.FinalizeTransaction(tx)

	// fetch checkout lines:
	var (
		checkoutLines     checkout.CheckoutLines
		productVariantIDs []string
	)
	_, err = tx.Select(
		&checkoutLines,
		"SELECT * FROM "+store.CheckoutLineTableName+" WHERE CheckoutID = :CheckoutID ORDER BY :OrderBy",
		map[string]interface{}{
			"CheckoutID": ckout.Token,
			"OrderBy":    store.TableOrderingMap[store.CheckoutLineTableName],
		},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find checkout lines belong to checkout with token=%s", ckout.Token)
	}
	productVariantIDs = checkoutLines.VariantIDs()

	// fetch product variants
	var (
		productVariants []*product_and_discount.ProductVariant
		productIDs      []string
		// productVariantMap has keys are product variant ids
		productVariantMap = map[string]*product_and_discount.ProductVariant{}
	)
	// check if we can proceed:
	if len(productVariantIDs) > 0 {
		_, err = tx.Select(
			&productVariants,
			"SELECT * FROM "+store.ProductVariantTableName+" WHERE Id IN :IDs ORDER BY :OrderBy",
			map[string]interface{}{
				"IDs":     productVariantIDs,
				"OrderBy": store.TableOrderingMap[store.ProductVariantTableName],
			},
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find product variants")
		}
		for _, variant := range productVariants {
			productIDs = append(productIDs, variant.ProductID)
			productVariantMap[variant.Id] = variant
		}
	}

	// fetch products
	var (
		products       []*product_and_discount.Product
		productTypeIDs []string
		// productMap has keys are product ids
		productMap = map[string]*product_and_discount.Product{}
	)
	// check if we can proceed:
	if len(productIDs) > 0 {
		_, err = tx.Select(
			&products,
			"SELECT * FROM "+store.ProductTableName+" WHERE Id IN :IDs ORDER BY :OrderBy",
			map[string]interface{}{
				"IDs":     productIDs,
				"OrderBy": store.TableOrderingMap[store.ProductTableName],
			},
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to finds products")
		}
		for _, product := range products {
			productTypeIDs = append(productTypeIDs, product.ProductTypeID)
			productMap[product.Id] = product
		}
	}

	// fetch product collections
	var (
		collectionXs []*struct {
			product_and_discount.Collection
			PrefetchRelatedValProductID string
		}
		// collectionsByProducts has keys are product ids
		collectionsByProducts = map[string][]*product_and_discount.Collection{}
	)
	// check if we can proceed
	if len(productIDs) > 0 {
		query, args, _ := cs.GetQueryBuilder().
			Select(cs.Collection().ModelFields()...).
			From(store.ProductCollectionTableName).
			Column(squirrel.Alias(squirrel.Expr("ProductCollections.ProductID"), "PrefetchRelatedValProductID")). // extra collumn
			InnerJoin(store.CollectionProductRelationTableName + " ON (ProductCollections.CollectionID = Collections.Id)").
			Where(squirrel.Eq{
				"ProductCollections.ProductID": productIDs,
			}).
			OrderBy(store.TableOrderingMap[store.ProductCollectionTableName]).
			ToSql()

		_, err = tx.Select(&collectionXs, query, args...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find collections")
		}
		for _, collectionX := range collectionXs {
			collectionsByProducts[collectionX.PrefetchRelatedValProductID] = append(collectionsByProducts[collectionX.PrefetchRelatedValProductID], &collectionX.Collection)
		}
	}

	// fetch product variant channel listing
	var (
		productVariantChannelListings []*product_and_discount.ProductVariantChannelListing
		channelIDs                    []string
		// productVariantChannelListingsByProductVariant has keys are product variant ids
		productVariantChannelListingsByProductVariant = map[string][]*product_and_discount.ProductVariantChannelListing{}
	)
	// check if we can proceed:
	if len(productVariantIDs) > 0 {
		_, err = tx.Select(
			&productVariantChannelListings,
			"SELECT * FROM "+store.ProductVariantChannelListingTableName+" WHERE VariantID IN :IDs ORDER BY :OrderBy",
			map[string]interface{}{
				"IDs":     productVariantIDs,
				"OrderBy": store.TableOrderingMap[store.ProductVariantChannelListingTableName],
			},
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find product variant channel listing")
		}
		for _, listing := range productVariantChannelListings {
			channelIDs = append(channelIDs, listing.ChannelID)
			productVariantChannelListingsByProductVariant[listing.VariantID] = append(productVariantChannelListingsByProductVariant[listing.VariantID], listing)
		}
	}

	// fetch channels
	var channels []*channel.Channel
	// check if we can proceed
	if len(channelIDs) > 0 {
		_, err = tx.Select(
			&channels,
			"SELECT * FROM "+store.ChannelTableName+" WHERE Id in :IDs ORDER BY :OrderBy",
			map[string]interface{}{
				"IDs":     channelIDs,
				"OrderBy": store.TableOrderingMap[store.ChannelTableName],
			},
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find channels")
		}
	}

	// fetch product types
	var (
		productTypes []*product_and_discount.ProductType
		// productTypeMap has keys are product type ids
		productTypeMap = map[string]*product_and_discount.ProductType{}
	)
	// check if we can proceed
	if len(productTypeIDs) > 0 {
		_, err = tx.Select(
			&productTypes,
			"SELECT * FROM "+store.ProductTypeTableName+" WHERE Id IN :IDs ORDER BY :OrderBy",
			map[string]interface{}{
				"IDs":     productTypeIDs,
				"OrderBy": store.TableOrderingMap[store.ProductTypeTableName],
			},
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to finds product types")
		}
		for _, productType := range productTypes {
			productTypeMap[productType.Id] = productType
		}
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "transaction_commit")
	}

	var checkoutLineInfos []*checkout.CheckoutLineInfo

	for _, checkoutLine := range checkoutLines {
		productVariant := productVariantMap[checkoutLine.VariantID]

		if productVariant != nil {
			var variantChannelListing *product_and_discount.ProductVariantChannelListing
			for _, listing := range productVariantChannelListingsByProductVariant[productVariant.Id] {
				if listing.ChannelID == ckout.ChannelID {
					variantChannelListing = listing
				}
			}

			// FIXME: Temporary solution to pass type checks. Figure out how to handle case
			// when variant channel listing is not defined for a checkout line.
			if variantChannelListing == nil {
				continue
			}

			product := productMap[productVariant.ProductID]
			if product != nil {
				productType := productTypeMap[product.ProductTypeID]
				collections := collectionsByProducts[product.Id]

				if productType != nil && collections != nil {
					checkoutLineInfos = append(checkoutLineInfos, &checkout.CheckoutLineInfo{
						Line:           *checkoutLine,
						Variant:        *productVariant,
						ChannelListing: *variantChannelListing,
						Product:        *product,
						ProductType:    *productType,
						Collections:    collections,
					})
				}
			}
		}
	}

	return checkoutLineInfos, nil
}

// DeleteCheckoutsByOption deletes checkout row(s) from database, filtered using given option.
// It returns an error indicating if the operation was performed successfully.
func (cs *SqlCheckoutStore) DeleteCheckoutsByOption(transaction *gorp.Transaction, option *checkout.CheckoutFilterOption) error {
	var runner squirrel.BaseRunner = cs.GetMaster()
	if transaction != nil {
		runner = transaction
	}
	_, err := cs.commonFilterQueryBuilder(option, delete).(squirrel.DeleteBuilder).RunWith(runner).Exec()
	if err != nil {
		return errors.Wrap(err, "failed to delete checkout(s) by given options")
	}

	return nil
}
