package warehouse

import (
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlWarehouseShippingZoneStore struct {
	store.Store
}

func NewSqlWarehouseShippingZoneStore(s store.Store) store.WarehouseShippingZoneStore {
	return &SqlWarehouseShippingZoneStore{s}
}

func (ws *SqlWarehouseShippingZoneStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"WarehouseID",
		"ShippingZoneID",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// Save inserts given warehouse-shipping zone relation into database
func (ws *SqlWarehouseShippingZoneStore) Save(warehouseShippingZone *model.WarehouseShippingZone) (*model.WarehouseShippingZone, error) {
	warehouseShippingZone.PreSave()
	if err := warehouseShippingZone.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.WarehouseShippingZoneTableName + "(" + ws.ModelFields("").Join(",") + ") VALUES (" + ws.ModelFields(":").Join(",") + ")"
	_, err := ws.GetMasterX().NamedExec(query, warehouseShippingZone)
	if err != nil {
		if ws.IsUniqueConstraintError(err, []string{"WarehouseID", "ShippingZoneID", "warehouseshippingzones_warehouseid_shippingzoneid_key"}) {
			return nil, store.NewErrInvalidInput(store.WarehouseShippingZoneTableName, "WarehouseID/ShippingZoneID", "duplicate")
		}
		return nil, errors.Wrapf(err, "failed to save warehouse-shipping zone relation with id=%s", warehouseShippingZone.Id)
	}

	return warehouseShippingZone, nil
}

func (s *SqlWarehouseShippingZoneStore) FilterByCountryCodeAndChannelID(countryCode, channelID string) ([]*model.WarehouseShippingZone, error) {
	countryCode = strings.ToUpper(countryCode)

	query := s.
		GetQueryBuilder().
		Select(s.ModelFields(store.WarehouseShippingZoneTableName + ".")...)

	if countryCode != "" {
		shippingZoneQuery := s.
			GetQueryBuilder(squirrel.Question).
			Select(`(1) AS "a"`).
			Prefix("EXISTS (").
			Suffix(")").
			From(store.ShippingZoneTableName).
			Where("ShippingZones.Countries::text LIKE ?", "%"+countryCode+"%").
			Where("ShippingZones.Id = WarehouseShippingZones.ShippingZoneID").
			Limit(1)

		query = query.Where(shippingZoneQuery)
	}

	if channelID != "" {
		channelQuery := s.
			GetQueryBuilder(squirrel.Question).
			Select(`(1) AS "a"`).
			Prefix("EXISTS (").
			Suffix(")").
			From(store.ChannelTableName).
			Where("Channels.Id = ?", channelID).
			Where("Channels.Id = ShippingZoneChannels.ChannelID").
			Limit(1)

		shippingZoneChannelQuery := s.
			GetQueryBuilder(squirrel.Question).
			Select(`(1) AS "a"`).
			Prefix("EXISTS (").
			Suffix(")").
			From(store.ShippingZoneChannelTableName).
			Where(channelQuery).
			Where("ShippingZoneChannels.ShippingZoneID = WarehouseShippingZones.ShippingZoneID").
			Limit(1)

		query = query.Where(shippingZoneChannelQuery)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByCountryCodeAndChannelID_ToSql")
	}

	var res []*model.WarehouseShippingZone
	err = s.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find warehouse shipping zones by options")
	}

	return res, nil
}

func (s *SqlWarehouseShippingZoneStore) FilterByOptions(options *model.WarehouseShippingZoneFilterOption) ([]*model.WarehouseShippingZone, error) {
	query := s.GetQueryBuilder().Select("*").From(store.WarehouseShippingZoneTableName)

	if options.WarehouseID != nil {
		query = query.Where(options.WarehouseID)
	}
	if options.ShippingZoneID != nil {
		query = query.Where(options.ShippingZoneID)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var res []*model.WarehouseShippingZone
	err = s.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find warehouse shipping zones by given options")
	}

	return res, nil
}
