package checkout

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

type SqlCheckoutStore struct {
	store.Store
}

func NewSqlCheckoutStore(sqlStore store.Store) store.CheckoutStore {
	return &SqlCheckoutStore{sqlStore}
}

func (cs *SqlCheckoutStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Token",
		"CreateAt",
		"UpdateAt",
		"UserID",
		"Email",
		"Quantity",
		"ChannelID",
		"BillingAddressID",
		"ShippingAddressID",
		"ShippingMethodID",
		"CollectionPointID",
		"Note",
		"Currency",
		"Country",
		"DiscountAmount",
		"DiscountName",
		"TranslatedDiscountName",
		"VoucherCode",
		"RedirectURL",
		"TrackingCode",
		"LanguageCode",
		"Metadata",
		"PrivateMetadata",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (cs *SqlCheckoutStore) ScanFields(checkOut *model.Checkout) []interface{} {
	return []interface{}{
		&checkOut.Token,
		&checkOut.CreateAt,
		&checkOut.UpdateAt,
		&checkOut.UserID,
		&checkOut.Email,
		&checkOut.Quantity,
		&checkOut.ChannelID,
		&checkOut.BillingAddressID,
		&checkOut.ShippingAddressID,
		&checkOut.ShippingMethodID,
		&checkOut.CollectionPointID,
		&checkOut.Note,
		&checkOut.Currency,
		&checkOut.Country,
		&checkOut.DiscountAmount,
		&checkOut.DiscountName,
		&checkOut.TranslatedDiscountName,
		&checkOut.VoucherCode,
		&checkOut.RedirectURL,
		&checkOut.TrackingCode,
		&checkOut.LanguageCode,
		&checkOut.Metadata,
		&checkOut.PrivateMetadata,
	}
}

// Upsert depends on given checkout's Token property to decide to update or insert it
func (cs *SqlCheckoutStore) Upsert(transaction store_iface.SqlxTxExecutor, checkouts []*model.Checkout) ([]*model.Checkout, error) {
	runner := cs.GetMasterX()
	if transaction != nil {
		runner = transaction
	}
	saveQuery := "INSERT INTO " + store.CheckoutTableName + " (" + cs.ModelFields("").Join(",") + ") VALUES (" + cs.ModelFields(":").Join(",") + ")"
	updateQuery := "UPDATE " + store.CheckoutTableName + " SET " + cs.
		ModelFields("").
		Map(func(_ int, s string) string {
			return s + "=:" + s
		}).
		Join(",") + " WHERE Token=:Token"

	for _, checkout := range checkouts {
		isSaving := false

		if !model.IsValidId(checkout.Token) {
			isSaving = true
			checkout.Token = ""
			checkout.PreSave()
		} else {
			checkout.PreUpdate()
		}

		appErr := checkout.IsValid()
		if appErr != nil {
			return nil, appErr
		}

		var err error

		if isSaving {
			_, err = runner.NamedExec(saveQuery, checkout)
		} else {

			var oldCheckout model.Checkout
			eror := runner.Get(&oldCheckout, "SELECT * FROM "+store.CheckoutTableName+" WHERE Token = $1", checkout.Token)
			if eror != nil {
				return nil, eror
			}

			// keep uneditable field intact
			checkout.BillingAddressID = oldCheckout.BillingAddressID
			checkout.ShippingAddressID = oldCheckout.ShippingAddressID

			_, err = runner.NamedExec(updateQuery, checkout)
		}

		if err != nil {
			return nil, errors.Wrap(err, "failed to upsert checkout")
		}
	}

	return checkouts, nil
}

type checkoutStatement string

const (
	delete checkoutStatement = "delete"
	slect  checkoutStatement = "select"
)

// commonFilterQueryBuilder is common function, used to build checkout(s) filter queries.
func (cs *SqlCheckoutStore) commonFilterQueryBuilder(option *model.CheckoutFilterOption, statementType checkoutStatement) interface{} {
	andCondition := squirrel.And{}
	// parse option
	if option.Token != nil {
		andCondition = append(andCondition, option.Token)
	}
	if option.UserID != nil {
		andCondition = append(andCondition, option.UserID)
	}
	if option.ChannelID != nil {
		andCondition = append(andCondition, option.ChannelID)
	}
	if option.ChannelIsActive != nil {
		andCondition = append(andCondition, squirrel.Expr("Channels.IsActive = ?", *option.ChannelIsActive))
	}
	if option.ShippingMethodID != nil {
		andCondition = append(andCondition, option.ShippingMethodID)
	}

	if statementType == slect {
		selectFields := cs.ModelFields(store.CheckoutTableName + ".")
		if option.SelectRelatedChannel {
			selectFields = append(selectFields, cs.Channel().ModelFields(store.ChannelTableName+".")...)
		}
		if option.SelectRelatedBillingAddress {
			selectFields = append(selectFields, cs.Address().ModelFields(store.AddressTableName+".")...)
		}

		query := cs.GetQueryBuilder().
			Select(selectFields...).
			From(store.CheckoutTableName).
			Where(andCondition)

		if option.SelectRelatedChannel || option.ChannelIsActive != nil {
			query = query.InnerJoin(store.ChannelTableName + " ON Checkouts.ChannelID = Channels.Id")
		}

		if option.Limit > 0 {
			query = query.Limit(uint64(option.Limit))
		}

		return query
	}

	return cs.GetQueryBuilder().Delete(store.CheckoutTableName).Where(andCondition)
}

// GetByOption finds and returns 1 checkout based on given option
func (cs *SqlCheckoutStore) GetByOption(option *model.CheckoutFilterOption) (*model.Checkout, error) {
	option.Limit = 0 // no limit

	var (
		res            model.Checkout
		channel        model.Channel
		billingAddress model.Address
		scanFields     = cs.ScanFields(&res)
	)
	if option.SelectRelatedChannel {
		scanFields = append(scanFields, cs.Channel().ScanFields(&channel)...)
	}
	if option.SelectRelatedBillingAddress {
		scanFields = append(scanFields, cs.Address().ScanFields(&billingAddress)...)
	}

	query, args, err := cs.commonFilterQueryBuilder(option, slect).(squirrel.SelectBuilder).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetByOption_ToSql")
	}

	err = cs.GetReplicaX().QueryRowX(query, args...).Scan(scanFields...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.CheckoutTableName, "option")
		}
		return nil, errors.Wrap(err, "failed to scan a checkout")
	}

	if option.SelectRelatedChannel {
		res.SetChannel(&channel)
	}
	if option.SelectRelatedBillingAddress {
		res.SetBilingAddress(&billingAddress)
	}

	return &res, nil
}

// FilterByOption finds and returns a list of checkout based on given option
func (cs *SqlCheckoutStore) FilterByOption(option *model.CheckoutFilterOption) ([]*model.Checkout, error) {
	query, args, err := cs.commonFilterQueryBuilder(option, slect).(squirrel.SelectBuilder).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	rows, err := cs.GetReplicaX().QueryX(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find checkouts with given options")
	}

	var res []*model.Checkout

	for rows.Next() {
		var (
			checkout       model.Checkout
			channel        model.Channel
			billingAddress model.Address
			scanFields     = cs.ScanFields(&checkout)
		)
		if option.SelectRelatedChannel {
			scanFields = append(scanFields, cs.Channel().ScanFields(&channel)...)
		}
		if option.SelectRelatedBillingAddress {
			scanFields = append(scanFields, cs.Address().ScanFields(&billingAddress)...)
		}

		if err := rows.Scan(scanFields...); err != nil {
			return nil, errors.Wrap(err, "failed to scan a row of checkout")
		}

		if option.SelectRelatedChannel {
			checkout.SetChannel(&channel)
		}
		if option.SelectRelatedBillingAddress {
			checkout.SetBilingAddress(&billingAddress)
		}

		res = append(res, &checkout)
	}

	if err := rows.Close(); err != nil {
		return nil, errors.Wrap(err, "failed to close rows")
	}

	return res, nil
}

// FetchCheckoutLinesAndPrefetchRelatedValue Fetch checkout lines as CheckoutLineInfo objects.
func (cs *SqlCheckoutStore) FetchCheckoutLinesAndPrefetchRelatedValue(ckout *model.Checkout) ([]*model.CheckoutLineInfo, error) {
	// please refer to file checkout_store_sql.md for details

	// fetch checkout lines:
	var (
		checkoutLines     model.CheckoutLines
		productVariantIDs []string
	)
	err := cs.GetReplicaX().Select(
		&checkoutLines,
		"SELECT * FROM "+store.CheckoutLineTableName+" WHERE CheckoutID = ? ORDER BY ?",
		ckout.Token,
		store.TableOrderingMap[store.CheckoutLineTableName],
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find checkout lines belong to checkout with token=%s", ckout.Token)
	}
	productVariantIDs = checkoutLines.VariantIDs()

	// fetch product variants
	var (
		productVariants []*model.ProductVariant
		productIDs      []string
		// productVariantMap has keys are product variant ids
		productVariantMap = map[string]*model.ProductVariant{}
	)
	// check if we can proceed:
	if len(productVariantIDs) > 0 {

		queryString, args, _ := cs.GetQueryBuilder().Select("*").
			From(store.ProductVariantTableName).
			Where(squirrel.Eq{store.ProductVariantTableName + ".Id": productVariantIDs}).
			OrderBy(store.TableOrderingMap[store.ProductVariantTableName]).
			ToSql()

		err = cs.GetReplicaX().Select(
			&productVariants,
			queryString, args...,
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
		products       []*model.Product
		productTypeIDs []string
		// productMap has keys are product ids
		productMap = map[string]*model.Product{}
	)
	// check if we can proceed:
	if len(productIDs) > 0 {
		query, args, _ := cs.GetQueryBuilder().Select("*").
			From(store.ProductTableName).
			Where(squirrel.Eq{store.ProductTableName + ".Id": productIDs}).
			OrderBy(store.TableOrderingMap[store.ProductTableName]).ToSql()

		err = cs.GetReplicaX().Select(
			&products,
			query, args...,
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
			model.Collection
			PrefetchRelatedValProductID string
		}
		// collectionsByProducts has keys are product ids
		collectionsByProducts = map[string][]*model.Collection{}
	)
	// check if we can proceed
	if len(productIDs) > 0 {
		query, args, _ := cs.GetQueryBuilder().
			Select(cs.Collection().ModelFields(store.CollectionTableName + ".")...).
			From(store.CollectionTableName).
			Column("ProductCollections.ProductID AS PrefetchRelatedValProductID"). // extra collumn
			InnerJoin(store.CollectionProductRelationTableName + " ON (ProductCollections.CollectionID = Collections.Id)").
			Where(squirrel.Eq{"ProductCollections.ProductID": productIDs}).
			OrderBy(store.TableOrderingMap[store.CollectionTableName]).
			ToSql()

		err = cs.GetReplicaX().Select(&collectionXs, query, args...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find collections")
		}
		for _, collectionX := range collectionXs {
			collectionsByProducts[collectionX.PrefetchRelatedValProductID] = append(collectionsByProducts[collectionX.PrefetchRelatedValProductID], &collectionX.Collection)
		}
	}

	// fetch product variant channel listing
	var (
		productVariantChannelListings []*model.ProductVariantChannelListing
		channelIDs                    []string
		// productVariantChannelListingsByProductVariant has keys are product variant ids
		productVariantChannelListingsByProductVariant = map[string][]*model.ProductVariantChannelListing{}
	)
	// check if we can proceed:
	if len(productVariantIDs) > 0 {
		query, args, _ := cs.GetQueryBuilder().Select("*").
			From(store.ProductVariantChannelListingTableName).
			Where(squirrel.Eq{store.ProductVariantChannelListingTableName + ".VariantID": productVariantIDs}).
			OrderBy(store.TableOrderingMap[store.ProductVariantChannelListingTableName]).ToSql()

		err = cs.GetReplicaX().Select(
			&productVariantChannelListings,
			query, args...,
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
	var channels []*model.Channel
	// check if we can proceed
	if len(channelIDs) > 0 {
		err = cs.GetReplicaX().Select(
			&channels,
			"SELECT * FROM "+store.ChannelTableName+" WHERE Id in ? ORDER BY ?",
			channelIDs,
			store.TableOrderingMap[store.ChannelTableName],
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find channels")
		}
	}

	// fetch product types
	var (
		productTypes []*model.ProductType
		// productTypeMap has keys are product type ids
		productTypeMap = map[string]*model.ProductType{}
	)
	// check if we can proceed
	if len(productTypeIDs) > 0 {
		query, args, _ := cs.GetQueryBuilder().
			Select("*").
			From(store.ProductTypeTableName).
			Where(squirrel.Eq{"Id": productTypeIDs}).
			OrderBy(store.TableOrderingMap[store.ProductTypeTableName]).
			ToSql()
		err = cs.GetReplicaX().Select(
			&productTypes,
			query, args...,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to finds product types")
		}
		for _, productType := range productTypes {
			productTypeMap[productType.Id] = productType
		}
	}

	var checkoutLineInfos []*model.CheckoutLineInfo

	for _, checkoutLine := range checkoutLines {
		productVariant := productVariantMap[checkoutLine.VariantID]

		if productVariant != nil {
			var variantChannelListing *model.ProductVariantChannelListing
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
					checkoutLineInfos = append(checkoutLineInfos, &model.CheckoutLineInfo{
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
func (cs *SqlCheckoutStore) DeleteCheckoutsByOption(transaction store_iface.SqlxTxExecutor, option *model.CheckoutFilterOption) error {
	var runner = cs.GetMasterX()
	if transaction != nil {
		runner = transaction
	}

	query, args, err := cs.commonFilterQueryBuilder(option, delete).(squirrel.DeleteBuilder).ToSql()
	if err != nil {
		return errors.Wrap(err, "DeleteCheckoutsByOption_ToSql")
	}

	_, err = runner.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to delete checkout(s) by given options")
	}

	return nil
}

func (cs *SqlCheckoutStore) CountCheckouts(options *model.CheckoutFilterOption) (int64, error) {
	options.Limit = 0 // no limit

	query := cs.commonFilterQueryBuilder(options, slect).(squirrel.SelectBuilder)
	queryString, args, err := cs.GetQueryBuilder().Select("COUNT(*)").FromSelect(query, "count").ToSql()

	if err != nil {
		return 0, errors.Wrap(err, "CountCheckouts_ToSql")
	}

	var count int64
	err = cs.GetReplicaX().Get(&count, queryString, args...)
	if err != nil {
		return 0, errors.Wrap(err, "failed to count number of checkouts")
	}

	return count, err
}
