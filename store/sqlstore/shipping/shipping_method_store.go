package shipping

import (
	"database/sql"
	"fmt"

	"github.com/mattermost/squirrel"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SqlShippingMethodStore struct {
	store.Store
}

func NewSqlShippingMethodStore(s store.Store) store.ShippingMethodStore {
	return &SqlShippingMethodStore{s}
}

func (s *SqlShippingMethodStore) Upsert(transaction boil.ContextTransactor, method model.ShippingMethod) (*model.ShippingMethod, error) {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	isSaving := method.ID == ""
	if isSaving {
		model_helper.ShippingMethodPreSave(&method)
	} else {
		model_helper.ShippingMethodCommonPre(&method)
	}

	if err := model_helper.ShippingMethodIsValid(method); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		err = method.Insert(transaction, boil.Infer())
	} else {
		_, err = method.Update(transaction, boil.Infer())
	}
	if err != nil {
		return nil, err
	}

	return &method, nil
}

func (s *SqlShippingMethodStore) Get(methodID string) (*model.ShippingMethod, error) {
	method, err := model.FindShippingMethod(s.GetReplica(), methodID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.ShippingMethods, methodID)
		}
		return nil, err
	}

	return method, nil
}

func (ss *SqlShippingMethodStore) ApplicableShippingMethods(price goprices.Money, channelID string, weight measurement.Weight, countryCode model.CountryCode, productIDs []string) (model.ShippingMethodSlice, error) {
	var forExcludedProductQuery squirrel.Sqlizer = squirrel.Expr("(1=1)")
	if len(productIDs) > 0 {
		forExcludedProductQuery = ss.
			GetQueryBuilder(squirrel.Question).
			Select(`(1) AS "a"`).
			From(model.TableNames.ShippingMethodExcludedProducts).
			Where(squirrel.Eq{
				model.ShippingMethodExcludedProductTableColumns.ProductID:        productIDs,
				model.ShippingMethodExcludedProductTableColumns.ShippingMethodID: model.ShippingMethodTableColumns.ID,
			}).
			Limit(1).
			Prefix("NOT EXISTS (").
			Suffix(")")
	}

	countryCodeExpr := "%" + countryCode + "%"

	shippingMethodIdSelectQuery := ss.
		GetQueryBuilder(squirrel.Question).
		Select(model.ShippingMethodTableColumns.ID).
		From(model.TableNames.ShippingMethods).
		InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ShippingMethodChannelListings, model.ShippingMethodChannelListingTableColumns.ShippingMethodID, model.ShippingMethodTableColumns.ID)).
		InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ShippingZones, model.ShippingZoneTableColumns.ID, model.ShippingMethodTableColumns.ShippingZoneID)).
		InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ShippingZoneChannels, model.ShippingZoneChannelTableColumns.ShippingZoneID, model.ShippingZoneTableColumns.ID)).
		Where(squirrel.Eq{
			model.ShippingMethodChannelListingTableColumns.ChannelID: channelID,
			model.ShippingMethodChannelListingTableColumns.Currency:  price.Currency,
			model.ShippingZoneChannelTableColumns.ChannelID:          channelID,
			model.ShippingMethodTableColumns.Type:                    model.ShippingMethodTypePrice,
		}).
		Where(squirrel.ILike{
			model.ShippingZoneTableColumns.Countries: countryCodeExpr,
		}).
		Where(forExcludedProductQuery)

	shippingMethodIdSelectQuery = ss.
		GetQueryBuilder(squirrel.Question).
		Select(model.ShippingMethodChannelListingTableColumns.ShippingMethodID).
		From(model.TableNames.ShippingMethodChannelListings).
		Where(squirrel.Eq{
			model.ShippingMethodChannelListingTableColumns.ChannelID: channelID,
		}).
		Where(squirrel.Expr(
			model.ShippingMethodChannelListingTableColumns.ShippingMethodID+" IN ?",
			shippingMethodIdSelectQuery,
		)).
		Where(squirrel.LtOrEq{
			model.ShippingMethodChannelListingTableColumns.MinimumOrderPriceAmount: price.Amount,
		}).
		Where(squirrel.Or{
			squirrel.Eq{
				model.ShippingMethodChannelListingTableColumns.MaximumOrderPriceAmount: nil,
			},
			squirrel.GtOrEq{
				model.ShippingMethodChannelListingTableColumns.MaximumOrderPriceAmount: price.Amount,
			},
		})

	queryMods := []qm.QueryMod{
		qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ShippingMethodChannelListings, model.ShippingMethodChannelListingTableColumns.ShippingMethodID, model.ShippingMethodTableColumns.ID)),
		qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ShippingZones, model.ShippingZoneTableColumns.ID, model.ShippingMethodTableColumns.ShippingZoneID)),
		qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ShippingZoneChannels, model.ShippingZoneChannelTableColumns.ShippingZoneID, model.ShippingZoneTableColumns.ID)),
		qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ShippingMethodPostalCodeRules, model.ShippingMethodPostalCodeRuleTableColumns.ShippingMethodID, model.ShippingMethodTableColumns.ID)),
		qm.Select(model.TableNames.ShippingMethods + ".*"),
		qm.Select(
			fmt.Sprintf(
				`(
					SELECT
						%[1]s
					FROM
						%[2]s
					WHERE (
						%[3]s = '%[4]s'
						AND %[5]s = %[6]s
					)
				) AS PriceAmount`,
				model.ShippingMethodChannelListingTableColumns.PriceAmount, // 1
				model.TableNames.ShippingMethodChannelListings,             // 2
				model.ShippingMethodChannelListingTableColumns.ChannelID,   // 3
				channelID, // 4
				model.ShippingMethodChannelListingTableColumns.ShippingMethodID, // 5
				model.ShippingMethodTableColumns.ID,                             // 6
			),
		),
		qm.Load(model.ShippingMethodRels.ShippingZone),
		qm.Load(model.ShippingMethodRels.ShippingMethodPostalCodeRules),
		qm.OrderBy("PriceAmount ASC"),
		model_helper.Or{
			squirrel.And{
				squirrel.Eq{
					model.ShippingMethodChannelListingTableColumns.ChannelID: channelID,
					model.ShippingMethodChannelListingTableColumns.Currency:  price.Currency,
					model.ShippingZoneChannelTableColumns.ChannelID:          channelID,
					model.ShippingMethodTableColumns.Type:                    model.ShippingMethodTypePrice,
				},
				squirrel.ILike{model.ShippingZoneTableColumns.Countries: countryCodeExpr},
				squirrel.Expr(model.ShippingMethodTableColumns.ID+" IN ?", shippingMethodIdSelectQuery),
			},
			//
			squirrel.And{
				squirrel.Eq{
					model.ShippingMethodChannelListingTableColumns.ChannelID: channelID,
					model.ShippingMethodChannelListingTableColumns.Currency:  price.Currency,
					model.ShippingZoneChannelTableColumns.ChannelID:          channelID,
					model.ShippingMethodTableColumns.Type:                    model.ShippingMethodTypeWeight,
				},
				squirrel.ILike{
					model.ShippingZoneTableColumns.Countries: countryCodeExpr,
				},
				squirrel.Or{
					squirrel.LtOrEq{model.ShippingMethodTableColumns.MinimumOrderWeight: weight.Amount},
					squirrel.Eq{model.ShippingMethodTableColumns.MinimumOrderWeight: nil},
				},
				squirrel.Or{
					squirrel.GtOrEq{model.ShippingMethodTableColumns.MaximumOrderWeight: weight.Amount},
					squirrel.Eq{model.ShippingMethodTableColumns.MaximumOrderWeight: nil},
				},
				forExcludedProductQuery,
			},
		},
	}

	return model.ShippingMethods(queryMods...).All(ss.GetReplica())
}

func (ss *SqlShippingMethodStore) commonQueryBuilder(options model_helper.ShippingMethodFilterOption) []qm.QueryMod {
	conds := options.Conditions
	if options.ShippingZoneChannelSlug != nil || options.ShippingZoneCountries != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ShippingZones, model.ShippingZoneTableColumns.ID, model.ShippingMethodTableColumns.ShippingZoneID)),
		)

		if options.ShippingZoneCountries != nil {
			conds = append(conds, options.ShippingZoneCountries)
		}

		if options.ShippingZoneChannelSlug != nil {
			conds = append(
				conds,
				qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ShippingZoneChannels, model.ShippingZoneChannelTableColumns.ShippingZoneID, model.ShippingZoneTableColumns.ID)),
				qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Channels, model.ChannelTableColumns.ID, model.ShippingZoneChannelTableColumns.ChannelID)),
				options.ShippingZoneChannelSlug,
			)
		}
	}

	for _, load := range options.Load {
		conds = append(conds, qm.Load(load))
	}

	return conds
}

func (ss *SqlShippingMethodStore) FilterByOptions(options model_helper.ShippingMethodFilterOption) (model.ShippingMethodSlice, error) {
	conds := ss.commonQueryBuilder(options)
	return model.ShippingMethods(conds...).All(ss.GetReplica())
}

func (s *SqlShippingMethodStore) Delete(transaction boil.ContextTransactor, ids []string) error {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	_, err := model.ShippingMethods(model.ShippingMethodWhere.ID.IN(ids)).DeleteAll(transaction)
	return err
}
