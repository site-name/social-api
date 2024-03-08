package warehouse

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mattermost/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SqlWareHouseStore struct {
	store.Store
}

func NewSqlWarehouseStore(s store.Store) store.WarehouseStore {
	return &SqlWareHouseStore{s}
}

func (ws *SqlWareHouseStore) Upsert(wh model.Warehouse) (*model.Warehouse, error) {
	isSaving := wh.ID == ""
	if isSaving {
		model_helper.WarehousePreSave(&wh)
	} else {
		model_helper.WarehousePreUpdate(&wh)
	}

	if err := model_helper.WarehouseIsValid(wh); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		err = wh.Insert(ws.GetMaster(), boil.Infer())
	} else {
		_, err = wh.Update(ws.GetMaster(), boil.Blacklist(model.WarehouseColumns.CreatedAt))
	}

	if err != nil {
		if ws.IsUniqueConstraintError(err, []string{model.WarehouseColumns.Slug, "warehouses_slug_key"}) {
			return nil, store.NewErrInvalidInput(model.TableNames.Warehouses, model.WarehouseColumns.Slug, wh.Slug)
		}
		return nil, err
	}

	return &wh, nil
}

func (ws *SqlWareHouseStore) commonQueryBuilder(option model_helper.WarehouseFilterOption) []qm.QueryMod {
	conds := option.Conditions

	for _, load := range option.Preloads {
		conds = append(conds, qm.Load(load))
	}
	if option.ShippingZoneId != nil ||
		option.ShippingZoneCountries != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.WarehouseShippingZones, model.WarehouseTableColumns.ID, model.WarehouseShippingZoneTableColumns.WarehouseID)),
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ShippingZones, model.ShippingZoneTableColumns.ID, model.WarehouseShippingZoneTableColumns.ShippingZoneID)),
		)

		if option.ShippingZoneId != nil {
			conds = append(conds, option.ShippingZoneId)
		}
		if option.ShippingZoneCountries != nil {
			conds = append(conds, option.ShippingZoneCountries)
		}
	}
	if option.Search != "" {
		expr := "%" + option.Search + "%"

		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Addresses, model.WarehouseTableColumns.AddressID, model.AddressTableColumns.ID)),
			model_helper.Or{
				squirrel.ILike{model.WarehouseTableColumns.Name: expr},
				squirrel.ILike{model.WarehouseTableColumns.Email: expr},

				squirrel.ILike{model.AddressTableColumns.CompanyName: expr},
				squirrel.ILike{model.AddressTableColumns.StreetAddress1: expr},
				squirrel.ILike{model.AddressTableColumns.StreetAddress2: expr},
				squirrel.ILike{model.AddressTableColumns.City: expr},
				squirrel.ILike{model.AddressTableColumns.PostalCode: expr},
				squirrel.ILike{model.AddressTableColumns.Phone: expr},
			},
		)
	}

	return conds
}

func (wh *SqlWareHouseStore) FilterByOprion(option model_helper.WarehouseFilterOption) (model.WarehouseSlice, error) {
	conds := wh.commonQueryBuilder(option)
	return model.Warehouses(conds...).All(wh.GetReplica())
}

func (ws *SqlWareHouseStore) WarehouseByStockID(stockID string) (*model.Warehouse, error) {
	warehouse, err := model.Warehouses(
		qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Stocks, model.StockTableColumns.WarehouseID, model.WarehouseTableColumns.ID)),
		model.StockWhere.ID.EQ(stockID),
	).One(ws.GetReplica())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.Warehouses, "StockID="+stockID)
		}
		return nil, err
	}

	return warehouse, nil
}

func (ws *SqlWareHouseStore) ApplicableForClickAndCollectNoQuantityCheck(checkoutLines model.CheckoutLineSlice, country model.CountryCode) (model.WarehouseSlice, error) {
	variantIDs := make([]string, len(checkoutLines))
	for idx, line := range checkoutLines {
		variantIDs[idx] = line.VariantID
	}

	stocks, err := ws.Stock().FilterByOption(model_helper.StockFilterOption{
		Preloads: []string{
			model.StockRels.ProductVariant,
		},
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.StockWhere.ProductVariantID.IN(variantIDs),
		),
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to find stocks")
	}

	return ws.forCountryLinesAndStocks(checkoutLines, stocks, country)
}

func (w *SqlWareHouseStore) Delete(transaction boil.ContextTransactor, ids []string) error {
	if transaction == nil {
		transaction = w.GetMaster()
	}

	_, err := model.Warehouses(model.WarehouseWhere.ID.IN(ids)).DeleteAll(transaction)
	return err
}

func (s *SqlWareHouseStore) WarehouseShipingZonesByCountryCodeAndChannelID(countryCode, channelID string) (model.WarehouseShippingZoneSlice, error) {
	query := s.
		GetQueryBuilder().
		Select(model.TableNames.WarehouseShippingZones + ".*")

	if countryCode != "" {
		shippingZoneQuery := s.
			GetQueryBuilder(squirrel.Question).
			Select(`(1) AS "a"`).
			Prefix("EXISTS (").
			Suffix(")").
			From(model.TableNames.ShippingZones).
			Where(squirrel.ILike{
				model.ShippingZoneTableColumns.Countries: "%" + countryCode + "%",
			}).
			Where(squirrel.Eq{
				model.ShippingZoneTableColumns.ID: model.WarehouseShippingZoneTableColumns.ShippingZoneID,
			}).
			Limit(1)

		query = query.Where(shippingZoneQuery)
	}

	if channelID != "" {
		channelQuery := s.
			GetQueryBuilder(squirrel.Question).
			Select(`(1) AS "a"`).
			Prefix("EXISTS (").
			Suffix(")").
			From(model.TableNames.Channels).
			Where(squirrel.Eq{
				model.ChannelTableColumns.ID: channelID,
			}).
			Where(squirrel.Eq{
				model.ChannelTableColumns.ID: model.ShippingZoneChannelTableColumns.ChannelID,
			}).
			Limit(1)

		shippingZoneChannelQuery := s.
			GetQueryBuilder(squirrel.Question).
			Select(`(1) AS "a"`).
			Prefix("EXISTS (").
			Suffix(")").
			From(model.TableNames.ShippingZoneChannels).
			Where(channelQuery).
			Where(squirrel.Eq{
				model.ShippingZoneChannelTableColumns.ShippingZoneID: model.WarehouseShippingZoneTableColumns.ShippingZoneID,
			}).
			Limit(1)

		query = query.Where(shippingZoneChannelQuery)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByCountryCodeAndChannelID_ToSql")
	}

	var res model.WarehouseShippingZoneSlice
	err = queries.Raw(queryString, args...).Bind(context.Background(), s.GetReplica(), &res)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find warehouse shipping zones by options")
	}

	return res, nil
}

func (ws *SqlWareHouseStore) ApplicableForClickAndCollectCheckoutLines(checkoutLines model.CheckoutLineSlice, country model.CountryCode) (model.WarehouseSlice, error) {
	panic("not implemented")
}

func (s *SqlWareHouseStore) ApplicableForClickAndCollectOrderLines(orderLines model.OrderLineSlice, country model.CountryCode) (model.WarehouseSlice, error) {
	panic("not implemented")
}

func (ws *SqlWareHouseStore) forCountryLinesAndStocks(checkoutLines model.CheckoutLineSlice, stocks model.StockSlice, country model.CountryCode) (model.WarehouseSlice, error) {
	panic("not implemented")
}
