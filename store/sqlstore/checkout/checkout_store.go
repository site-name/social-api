package checkout

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlCheckoutStore struct {
	store.Store
}

func NewSqlCheckoutStore(sqlStore store.Store) store.CheckoutStore {
	return &SqlCheckoutStore{sqlStore}
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
func (cs *SqlCheckoutStore) Upsert(transaction *gorm.DB, checkouts []*model.Checkout) ([]*model.Checkout, error) {
	if transaction == nil {
		transaction = cs.GetMaster()
	}

	for _, checkout := range checkouts {
		var err error
		if checkout.Token == "" {
			err = transaction.Create(checkout).Error
		} else {
			checkout.ShippingAddressID = model.NewPrimitive("") // prevent update
			checkout.BillingAddressID = model.NewPrimitive("")  // prevent update
			checkout.CreateAt = 0                               // prevent update

			err = transaction.Model(checkout).Updates(checkout).Error
		}

		if err != nil {
			return nil, errors.Wrap(err, "failed to upsert a checkout")
		}
	}

	return checkouts, nil
}

func (cs *SqlCheckoutStore) commonFilterQueryBuilder(option *model.CheckoutFilterOption) squirrel.SelectBuilder {
	andCondition := squirrel.And{}
	// parse option
	if option.Conditions != nil {
		andCondition = append(andCondition, option.Conditions)
	}
	if option.ChannelIsActive != nil {
		andCondition = append(andCondition, option.ChannelIsActive)
	}

	selectFields := []string{model.CheckoutTableName + ".*"}
	if option.SelectRelatedChannel {
		selectFields = append(selectFields, model.ChannelTableName+".*")
	}
	if option.SelectRelatedBillingAddress {
		selectFields = append(selectFields, model.AddressTableName+".*")
	}
	if option.SelectRelatedUser {
		selectFields = append(selectFields, model.UserTableName+".*")
	}

	query := cs.GetQueryBuilder().
		Select(selectFields...).
		From(model.CheckoutTableName).
		Where(andCondition)

	if option.SelectRelatedChannel || option.ChannelIsActive != nil {
		query = query.InnerJoin(model.ChannelTableName + " ON Checkouts.ChannelID = Channels.Id")
	}
	if option.SelectRelatedBillingAddress {
		query = query.InnerJoin(model.AddressTableName + " ON Checkouts.BillingAddressID = Addresses.Id")
	}
	if option.SelectRelatedUser {
		query = query.InnerJoin(model.UserTableName + " ON Users.Id = Checkouts.UserID")
	}

	return query
}

// GetByOption finds and returns 1 checkout based on given option
func (cs *SqlCheckoutStore) GetByOption(option *model.CheckoutFilterOption) (*model.Checkout, error) {
	option.GraphqlPaginationValues.Limit = 0

	var (
		res            model.Checkout
		channel        model.Channel
		billingAddress model.Address
		user           model.User
		scanFields     = cs.ScanFields(&res)
	)
	if option.SelectRelatedChannel {
		scanFields = append(scanFields, cs.Channel().ScanFields(&channel)...)
	}
	if option.SelectRelatedBillingAddress {
		scanFields = append(scanFields, cs.Address().ScanFields(&billingAddress)...)
	}
	if option.SelectRelatedUser {
		scanFields = append(scanFields, cs.User().ScanFields(&user)...)
	}

	query, args, err := cs.commonFilterQueryBuilder(option).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetByOption_ToSql")
	}

	err = cs.GetReplica().Raw(query, args...).Row().Scan(scanFields)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.NewErrNotFound(model.CheckoutTableName, "option")
		}
		return nil, errors.Wrap(err, "failed to scan a checkout")
	}

	if option.SelectRelatedChannel {
		res.SetChannel(&channel)
	}
	if option.SelectRelatedBillingAddress {
		res.SetBilingAddress(&billingAddress)
	}
	if option.SelectRelatedUser {
		res.SetUser(&user)
	}

	return &res, nil
}

// FilterByOption finds and returns a list of checkout based on given option
func (cs *SqlCheckoutStore) FilterByOption(option *model.CheckoutFilterOption) (int64, []*model.Checkout, error) {
	filterQuery := cs.commonFilterQueryBuilder(option)

	var totalCount int64
	if option.CountTotal {
		countQuery, args, err := cs.GetQueryBuilder().Select("COUNT (*)").FromSelect(filterQuery, "subquery").ToSql()
		if err != nil {
			return 0, nil, errors.Wrap(err, "FilterByOption_Count_ToSql")
		}

		err = cs.GetReplica().Raw(countQuery, args...).Scan(&totalCount).Error
		if err != nil {
			return 0, nil, errors.Wrap(err, "failed to count total number of checkouts by given options")
		}
	}

	// apply pagination
	option.GraphqlPaginationValues.AddPaginationToSelectBuilderIfNeeded(&filterQuery)

	query, args, err := filterQuery.ToSql()
	if err != nil {
		return 0, nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	rows, err := cs.GetReplica().Raw(query, args...).Rows()
	if err != nil {
		return 0, nil, errors.Wrap(err, "failed to find checkouts with given options")
	}
	defer rows.Close()

	var res []*model.Checkout

	for rows.Next() {
		var (
			checkout       model.Checkout
			channel        model.Channel
			billingAddress model.Address
			user           model.User
			scanFields     = cs.ScanFields(&checkout)
		)
		if option.SelectRelatedChannel {
			scanFields = append(scanFields, cs.Channel().ScanFields(&channel)...)
		}
		if option.SelectRelatedBillingAddress {
			scanFields = append(scanFields, cs.Address().ScanFields(&billingAddress)...)
		}
		if option.SelectRelatedUser {
			scanFields = append(scanFields, cs.User().ScanFields(&user)...)
		}

		if err := rows.Scan(scanFields...); err != nil {
			return 0, nil, errors.Wrap(err, "failed to scan a row of checkout")
		}

		if option.SelectRelatedChannel {
			checkout.SetChannel(&channel)
		}
		if option.SelectRelatedBillingAddress {
			checkout.SetBilingAddress(&billingAddress)
		}
		if option.SelectRelatedUser {
			checkout.SetUser(&user)
		}

		res = append(res, &checkout)
	}

	return totalCount, res, nil
}

// FetchCheckoutLinesAndPrefetchRelatedValue Fetch checkout lines as CheckoutLineInfo objects.
func (cs *SqlCheckoutStore) FetchCheckoutLinesAndPrefetchRelatedValue(checkout *model.Checkout) ([]*model.CheckoutLineInfo, error) {
	// please refer to file checkout_store_sql.md for details

	// fetch checkout lines:
	var checkoutLines model.CheckoutLines

	err := cs.GetReplica().Order("CreateAt ASC").Find(&checkoutLines, "CheckoutID = ?", checkout.Token).Error
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find checkout lines belong to checkout with token=%s", checkout.Token)
	}
	productVariantIDs := checkoutLines.VariantIDs()

	// fetch product variants
	var (
		productIDs        []string
		productVariantMap = map[string]*model.ProductVariant{} // productVariantMap has keys are product variant ids
	)
	// check if we can proceed:
	if len(productVariantIDs) > 0 {
		var productVariants model.ProductVariants
		err = cs.GetReplica().Find(&productVariants, "Id = ?", productVariantIDs).Error
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
		products       model.Products
		productTypeIDs []string
		productMap     = map[string]*model.Product{} // productMap has keys are product ids
	)
	// check if we can proceed:
	if len(productIDs) > 0 {
		err = cs.GetReplica().Find(&products, "Id IN ?", productIDs).Error
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
		collectionsByProducts = map[string]model.Collections{} // collectionsByProducts has keys are product ids
	)
	// check if we can proceed
	if len(productIDs) > 0 {
		err = cs.GetReplica().
			Table(model.CollectionTableName).
			Where("ProductCollections.ProductID IN ?", productIDs).
			Select("Collections.*", "ProductCollections.ProductID AS PrefetchRelatedValProductID").
			Joins("INNER JOIN " + model.CollectionProductRelationTableName + " ON ProductCollections.CollectionID = Collections.Id").
			Scan(&collectionXs).Error

		if err != nil {
			return nil, errors.Wrap(err, "failed to find collections")
		}
		for _, collectionX := range collectionXs {
			collectionsByProducts[collectionX.PrefetchRelatedValProductID] = append(collectionsByProducts[collectionX.PrefetchRelatedValProductID], &collectionX.Collection)
		}
	}

	// fetch product variant channel listing
	var (
		productVariantChannelListings                 []*model.ProductVariantChannelListing
		channelIDs                                    []string
		productVariantChannelListingsByProductVariant = map[string][]*model.ProductVariantChannelListing{} // productVariantChannelListingsByProductVariant has keys are product variant ids
	)
	// check if we can proceed:
	if len(productVariantIDs) > 0 {
		err := cs.GetReplica().Find(&productVariantChannelListings, "VariantID IN ?", productVariantIDs).Error
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
		err = cs.GetReplica().Find(&channels, "Id in ? ORDER BY Slug ASC", channelIDs).Error
		if err != nil {
			return nil, errors.Wrap(err, "failed to find channels")
		}
	}

	// fetch product types
	var (
		productTypes   []*model.ProductType
		productTypeMap = map[string]*model.ProductType{} // productTypeMap has keys are product type ids
	)
	// check if we can proceed
	if len(productTypeIDs) > 0 {
		err = cs.GetReplica().Find(&productTypes, "Id IN ?", productTypeIDs).Error
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
				if listing.ChannelID == checkout.ChannelID {
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
func (cs *SqlCheckoutStore) DeleteCheckoutsByOption(transaction *gorm.DB, option *model.CheckoutFilterOption) error {
	if transaction == nil {
		transaction = cs.GetMaster()
	}

	query, args, err := cs.GetQueryBuilder().Delete(model.CheckoutTableName).Where(option.Conditions).ToSql()
	if err != nil {
		return errors.Wrap(err, "DeleteCheckoutsByOption_ToSql")
	}

	err = transaction.Raw(query, args...).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete checkout(s) by given options")
	}

	return nil
}

func (cs *SqlCheckoutStore) CountCheckouts(options *model.CheckoutFilterOption) (int64, error) {
	db := cs.GetReplica().Table(model.CheckoutTableName)

	conditions := squirrel.And{
		options.Conditions,
	}
	if options.ChannelIsActive != nil {
		conditions = append(conditions, options.ChannelIsActive)
		db = db.Joins("INNER JOIN Channels ON Checkouts.ChannelID = Channels.Id")
	}

	condStr, args, err := conditions.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "CountCheckouts_ToSql")
	}
	var count int64
	return count, db.Where(condStr, args...).Raw("SELECT COUNT(*) FROM " + model.CheckoutTableName).Scan(&count).Error
}
