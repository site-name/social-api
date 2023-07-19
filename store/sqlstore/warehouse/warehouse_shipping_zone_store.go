package warehouse

import (
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
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
func (ws *SqlWarehouseShippingZoneStore) Save(transaction *gorm.DB, warehouseShippingZones []*model.WarehouseShippingZone) ([]*model.WarehouseShippingZone, error) {
	runner := ws.GetMaster()
	if transaction != nil {
		runner = transaction
	}
	query := "INSERT INTO " + model.WarehouseShippingZoneTableName + "(" + ws.ModelFields("").Join(",") + ") VALUES (" + ws.ModelFields(":").Join(",") + ")"

	for _, relation := range warehouseShippingZones {
		relation.PreSave()

		appErr := relation.IsValid()
		if appErr != nil {
			return nil, appErr
		}

		_, err := runner.NamedExec(query, relation)
		if err != nil {
			if ws.IsUniqueConstraintError(err, []string{"WarehouseID", "ShippingZoneID", "warehouseshippingzones_warehouseid_shippingzoneid_key"}) {
				return nil, store.NewErrInvalidInput(model.WarehouseShippingZoneTableName, "WarehouseID/ShippingZoneID", "duplicate")
			}
			return nil, errors.Wrapf(err, "failed to save warehouse-shipping zone relation with id=%s", relation.Id)
		}
	}

	return warehouseShippingZones, nil
}

func (s *SqlWarehouseShippingZoneStore) FilterByCountryCodeAndChannelID(countryCode, channelID string) ([]*model.WarehouseShippingZone, error) {
	countryCode = strings.ToUpper(countryCode)

	query := s.
		GetQueryBuilder().
		Select(s.ModelFields(model.WarehouseShippingZoneTableName + ".")...)

	if countryCode != "" {
		shippingZoneQuery := s.
			GetQueryBuilder(squirrel.Question).
			Select(`(1) AS "a"`).
			Prefix("EXISTS (").
			Suffix(")").
			From(model.ShippingZoneTableName).
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
			From(model.ChannelTableName).
			Where("Channels.Id = ?", channelID).
			Where("Channels.Id = ShippingZoneChannels.ChannelID").
			Limit(1)

		shippingZoneChannelQuery := s.
			GetQueryBuilder(squirrel.Question).
			Select(`(1) AS "a"`).
			Prefix("EXISTS (").
			Suffix(")").
			From(model.ShippingZoneChannelTableName).
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
	err = s.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find warehouse shipping zones by options")
	}

	return res, nil
}

func (s *SqlWarehouseShippingZoneStore) FilterByOptions(options *model.WarehouseShippingZoneFilterOption) ([]*model.WarehouseShippingZone, error) {
	query := s.
		GetQueryBuilder().
		Select(s.ModelFields(model.WarehouseShippingZoneTableName + ".")...).
		From(model.WarehouseShippingZoneTableName)

	if options.Conditions != nil {
		query = query.Where(options.Conditions)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var res []*model.WarehouseShippingZone
	err = s.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find warehouse shipping zones by given options")
	}

	return res, nil
}

func (s *SqlWarehouseShippingZoneStore) Delete(transaction *gorm.DB, options *model.WarehouseShippingZoneFilterOption) error {
	if options == nil || options.Conditions == nil {
		return errors.New("please provide valid options to delete")
	}
	query, args, err := s.GetQueryBuilder().Delete(model.WarehouseShippingZoneTableName).Where(options.Conditions).ToSql()
	if err != nil {
		return errors.Wrap(err, "Delete_ToSql")
	}

	runner := s.GetMaster()
	if transaction != nil {
		runner = transaction
	}

	_, err = runner.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to delete warehouse shipping zones by options")
	}
	return nil
}
