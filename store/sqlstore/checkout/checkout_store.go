package checkout

import (
	"database/sql"
	"strings"

	"github.com/Masterminds/squirrel"
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
		table.ColMap("ChannelID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("BillingAddressID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ShippingAddressID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ShippingMethodID").SetMaxSize(store.UUID_MAX_LENGTH)
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

	cs.CreateForeignKeyIfNotExists(store.CheckoutTableName, "UserID", store.UserTableName, "Id", true)
	cs.CreateForeignKeyIfNotExists(store.CheckoutTableName, "ChannelID", store.ChannelTableName, "Id", false)
	cs.CreateForeignKeyIfNotExists(store.CheckoutTableName, "BillingAddressID", store.AddressTableName, "Id", false)
	cs.CreateForeignKeyIfNotExists(store.CheckoutTableName, "ShippingAddressID", store.AddressTableName, "Id", false)
	cs.CreateForeignKeyIfNotExists(store.CheckoutTableName, "ShippingMethodID", store.ShippingMethodTableName, "Id", false)
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

func (cs *SqlCheckoutStore) CheckoutsByUserID(userID string, channelActive bool) ([]*checkout.Checkout, error) {
	var checkouts []*checkout.Checkout

	query := cs.GetQueryBuilder().
		Select("*").
		From(store.CheckoutTableName).
		InnerJoin(store.ChannelTableName + " ON (Channels.Id = Checkouts.ChannelID)")

	condition := squirrel.And{
		squirrel.Eq{"Checkouts.UserID": userID},
	}

	if channelActive {
		condition = append(condition, squirrel.Eq{"Channels.IsActive": true})
	} else {
		condition = append(condition, squirrel.NotEq{"Channels.IsActive": true})
	}

	queryString, args, err := query.Where(condition).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "CheckoutsByUserID_ToSql")
	}

	_, err = cs.GetReplica().Select(&checkouts, queryString, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find checkouts for user with Id=%s", userID)
	}

	return checkouts, nil
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
		checkoutLines   []*checkout.CheckoutLine
		checkoutLineIDs []string
	)
	_, err = tx.Select(
		&checkoutLines,
		"SELECT * FROM CheckoutLines WHERE CheckoutID = :CheckoutID ORDER BY :OrderBy",
		map[string]interface{}{
			"CheckoutID": ckout.Token,
			"OrderBy":    store.TableOrderingMap[store.CheckoutLineTableName],
		},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find checkout lines belong to checkout with token=%s", ckout.Token)
	}

	for _, line := range checkoutLines {
		checkoutLineIDs = append(checkoutLineIDs, line.Id)
	}

	// fetch product variants
	var (
		productVariants   []*product_and_discount.ProductVariant
		productIDs        []string
		productVariantIDs []string
	)
	_, err = tx.Select(
		&productVariants,
		"SELECT * FROM ProductVariants WHERE Id IN :IDs ORDER BY :OrderBy",
		map[string]interface{}{
			"IDs":     checkoutLineIDs,
			"OrderBy": store.TableOrderingMap[store.ProductVariantTableName],
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find product variants")
	}
	for _, variant := range productVariants {
		productIDs = append(productIDs, variant.ProductID)
		productVariantIDs = append(productVariantIDs, variant.Id)
	}

	// fetch products
	var (
		products       []*product_and_discount.Product
		productTypeIDs []string
	)
	_, err = tx.Select(
		&products,
		"SELECT * FROM Products WHERE Id IN :IDs ORDER BY :OrderBy",
		map[string]interface{}{
			"IDs":     productVariantIDs,
			"OrderBy": store.TableOrderingMap[store.ProductTableName],
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to finds products")
	}
	for _, prd := range products {
		productTypeIDs = append(productTypeIDs, prd.ProductTypeID)
	}

	// fetch product collections
	var (
		collectionXs []*struct {
			PrefetchRelatedValProductID string
			product_and_discount.Collection
		}
	)
	_, err = tx.Select(
		&collectionXs,
		`SELECT 
			ProductCollections.ProductID AS PrefetchRelatedValProductID, `+strings.Join(cs.Collection().ModelFields(), ", ")+`
		FROM
			Collections
		INNER JOIN ProductCollections ON (
			ProductCollections.CollectionID = Collections.Id
		)
		WHERE 
			ProductCollections.ProductID IN :IDs
		ORDER BY :OrderBy`,
		map[string]interface{}{
			"IDs":     productIDs,
			"OrderBy": store.TableOrderingMap[store.ProductCollectionTableName],
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find collections")
	}

	// fetch product variant channel listing
	var (
		productVariantChannelListings []*product_and_discount.ProductVariantChannelListing
		channelIDs                    []string
	)
	_, err = tx.Select(
		&productVariantChannelListings,
		"SELECT * FROM ProductVariantChannelListings WHERE VariantID IN :IDs ORDER BY :OrderBy",
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
	}

	// fetch channels
	var (
		channels []*channel.Channel
	)
	_, err = tx.Select(
		&channels,
		"SELECT * FROM Channels WHERE Id in :IDs ORDER BY :OrderBy",
		map[string]interface{}{
			"IDs":     channelIDs,
			"OrderBy": store.TableOrderingMap[store.ChannelTableName],
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find channels")
	}

	// fetch product types
	var (
		productTypes []*product_and_discount.ProductType
	)
	_, err = tx.Select(
		&productTypes,
		"SELECT * FROM ProductTypes WHERE Id IN :IDs ORDER BY :OrderBy",
		map[string]interface{}{
			"IDs":     productTypeIDs,
			"OrderBy": store.TableOrderingMap[store.ProductTypeTableName],
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to finds product types")
	}
}
