package shipping

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/store"
)

type SqlShippingMethodStore struct {
	store.Store
}

var (
	ShippingMethodSelectFields = []string{
		"SM.Id",
		"SM.Name",
		"SM.Type",
		"SM.ShippingZoneID",
		"SM.MinimumOrderWeight",
		"SM.MaximumOrderWeight",
		"SM.WeightUnit",
		"SM.MaximumDeliveryDays",
		"SM.MinimumDeliveryDays",
		"SM.Description",
		"SM.Metadata",
		"SM.PrivateMetadata",
	}
)

func NewSqlShippingMethodStore(s store.Store) store.ShippingMethodStore {
	smls := &SqlShippingMethodStore{s}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(shipping.ShippingMethod{}, store.ShippingMethodTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ShippingZoneID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(shipping.SHIPPING_METHOD_NAME_MAX_LENGTH)
		table.ColMap("Type").SetMaxSize(shipping.SHIPPING_METHOD_TYPE_MAX_LENGTH)
	}
	return smls
}

func (s *SqlShippingMethodStore) CreateIndexesIfNotExists() {
	s.CreateIndexIfNotExists("idx_shipping_methods_name", store.ShippingMethodTableName, "Name")
	s.CreateIndexIfNotExists("idx_shipping_methods_name_lower_textpattern", store.ShippingMethodTableName, "lower(Name) text_pattern_ops")

	s.CreateForeignKeyIfNotExists(store.ShippingMethodTableName, "ShippingZoneID", store.ShippingZoneTableName, "Id", true)
}

// Upsert bases on given method's Id to decide update or insert it
func (s *SqlShippingMethodStore) Upsert(method *shipping.ShippingMethod) (*shipping.ShippingMethod, error) {
	method.PreSave()
	if err := method.IsValid(); err != nil {
		return nil, err
	}

	err := s.GetMaster().Insert(method)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to upsert shipping method with id=%s", method.Id)
	}

	return method, nil
}

// Get finds and returns a shipping method with given id
func (s *SqlShippingMethodStore) Get(methodID string) (*shipping.ShippingMethod, error) {
	result, err := s.GetReplica().Get(shipping.ShippingMethod{}, methodID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ShippingMethodTableName, methodID)
		}
		return nil, errors.Wrapf(err, "failed to find shipping method with id=%s", methodID)
	}

	return result.(*shipping.ShippingMethod), nil
}

// ShippingMethodsByOption finds and returns a list of shipping methods that satisfy given filtering option
func (s *SqlShippingMethodStore) ShippingMethodsByOption(option *shipping.ShippingMethodFilterOption) ([]*shipping.ShippingMethod, error) {
	query := s.
		GetQueryBuilder().
		Select(ShippingMethodSelectFields...).
		From(store.ShippingMethodTableName + " AS SM")

	// check type:
	if option.Type != nil {
		query = query.Where(option.Type.ToSquirrel("SM.Type"))
	}

	// check shipping zone channel
	// if option.ShippingZoneChannelSlug != nil {
	// 	query = query.
	// 		InnerJoin(store.ShippingZoneTableName + " AS SZ ON (SZ.Id = SM.ShippingZoneID)").
	// 		InnerJoin(store.ShippingZoneChannelTableName + " AS SZCn ON (SZ.Id = SZCn.ShippingZoneID)").
	// 		InnerJoin(store.ChannelTableName + " AS Cn ON (Cn.Id = SZCn.ChannelID)").
	// 		Where(option.ShippingZoneChannelSlug.ToSquirrel("Cn.Slug"))
	// }

	// check channel listing
	// if option.ChannelListingsChannelSlug != nil {
	// 	query = query.
	// 		InnerJoin(store.ShippingMethodChannelListingTableName + " AS SMCL ON (SMCL.ShippingMethodID = SM.Id)").
	// 		InnerJoin(store.ChannelTableName + " AS CN ON ")
	// }
}
